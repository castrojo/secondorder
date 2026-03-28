package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/msoedov/thelastorg/internal/db"
	"github.com/msoedov/thelastorg/internal/models"
	"github.com/google/uuid"
)

type API struct {
	db   *db.DB
	sse  *SSEHub
	wake func(agent *models.Agent, issue *models.Issue)
}

func NewAPI(database *db.DB, sse *SSEHub, wake func(*models.Agent, *models.Issue)) *API {
	return &API{db: database, sse: sse, wake: wake}
}

// Auth middleware extracts agent from API key
func (a *API) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			http.Error(w, `{"error":"missing api key"}`, http.StatusUnauthorized)
			return
		}
		hash := sha256.Sum256([]byte(token))
		keyHash := hex.EncodeToString(hash[:])
		agent, err := a.db.GetAgentByAPIKey(keyHash)
		if err != nil {
			http.Error(w, `{"error":"invalid api key"}`, http.StatusUnauthorized)
			return
		}
		r = r.WithContext(withAgent(r.Context(), agent))
		next(w, r)
	}
}

func (a *API) Inbox(w http.ResponseWriter, r *http.Request) {
	agent := agentFromContext(r.Context())
	issues, err := a.db.GetAgentInbox(agent.ID)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonOK(w, issues)
}

func (a *API) GetIssue(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	issue, err := a.db.GetIssue(key)
	if err != nil {
		jsonError(w, "issue not found", http.StatusNotFound)
		return
	}
	comments, _ := a.db.ListComments(key)
	children, _ := a.db.GetChildIssues(key)

	jsonOK(w, map[string]any{
		"issue":    issue,
		"comments": comments,
		"children": children,
	})
}

func (a *API) CheckoutIssue(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	agent := agentFromContext(r.Context())

	var body struct {
		AgentID          string   `json:"agentId"`
		ExpectedStatuses []string `json:"expectedStatuses"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "invalid body", http.StatusBadRequest)
		return
	}
	if len(body.ExpectedStatuses) == 0 {
		body.ExpectedStatuses = []string{"todo", "backlog"}
	}

	if err := a.db.CheckoutIssue(key, agent.ID, body.ExpectedStatuses); err != nil {
		jsonError(w, err.Error(), http.StatusConflict)
		return
	}

	a.db.LogActivity("checkout", "issue", key, &agent.ID, "")
	jsonOK(w, map[string]string{"status": "checked_out"})
}

func (a *API) UpdateIssue(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	issue, err := a.db.GetIssue(key)
	if err != nil {
		jsonError(w, "issue not found", http.StatusNotFound)
		return
	}

	var body struct {
		Status  string `json:"status"`
		Comment string `json:"comment"`
		Title   string `json:"title"`
		Description string `json:"description"`
		Priority *int   `json:"priority"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "invalid body", http.StatusBadRequest)
		return
	}

	agent := agentFromContext(r.Context())

	if body.Status != "" {
		issue.Status = body.Status
	}
	if body.Title != "" {
		issue.Title = body.Title
	}
	if body.Description != "" {
		issue.Description = body.Description
	}
	if body.Priority != nil {
		issue.Priority = *body.Priority
	}

	if err := a.db.UpdateIssue(issue); err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add comment if provided
	if body.Comment != "" {
		agentName := "Board"
		if agent != nil {
			agentName = agent.Name
		}
		comment := &models.Comment{
			ID:        uuid.New().String(),
			IssueKey:  key,
			AgentID:   ptrStr(agent),
			Author:    agentName,
			Body:      body.Comment,
		}
		a.db.CreateComment(comment)

		// SSE broadcast
		data, _ := json.Marshal(map[string]string{
			"issue_key": key,
			"author":    agentName,
			"body":      body.Comment,
		})
		a.sse.Broadcast("comment", string(data))
		if a.wake != nil && agent != nil {
			go a.notifyOnComment(key, agentName, body.Comment)
		}
	}

	a.db.LogActivity("update", "issue", key, ptrStr(agent), body.Status)

	// Wake chain on status change
	if body.Status == models.StatusDone || body.Status == models.StatusBlocked || body.Status == models.StatusInReview {
		if agent != nil && a.wake != nil {
			go a.wakeReviewerForIssue(agent.ID, key)
		}
	}

	jsonOK(w, issue)
}

func (a *API) CreateComment(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	agent := agentFromContext(r.Context())

	var body struct {
		Body string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Body == "" {
		jsonError(w, "body required", http.StatusBadRequest)
		return
	}

	agentName := "Board"
	if agent != nil {
		agentName = agent.Name
	}

	comment := &models.Comment{
		ID:       uuid.New().String(),
		IssueKey: key,
		AgentID:  ptrStr(agent),
		Author:   agentName,
		Body:     body.Body,
	}
	if err := a.db.CreateComment(comment); err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(map[string]string{
		"issue_key": key,
		"author":    agentName,
		"body":      body.Body,
	})
	a.sse.Broadcast("comment", string(data))

	jsonOK(w, comment)
}

func (a *API) CreateIssue(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Title          string `json:"title"`
		Description    string `json:"description"`
		AssigneeSlug   string `json:"assignee_slug"`
		ParentIssueKey string `json:"parent_issue_key"`
		Priority       int    `json:"priority"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Title == "" {
		jsonError(w, "title required", http.StatusBadRequest)
		return
	}

	key, err := a.db.NextIssueKey()
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	issue := &models.Issue{
		ID:          uuid.New().String(),
		Key:         key,
		Title:       body.Title,
		Description: body.Description,
		Status:      models.StatusTodo,
		Priority:    body.Priority,
	}

	if body.ParentIssueKey != "" {
		issue.ParentIssueKey = &body.ParentIssueKey
	}

	// Resolve assignee
	var assignee *models.Agent
	if body.AssigneeSlug != "" {
		assignee, err = a.db.GetAgentBySlug(body.AssigneeSlug)
		if err != nil {
			jsonError(w, "agent not found: "+body.AssigneeSlug, http.StatusBadRequest)
			return
		}
		issue.AssigneeAgentID = &assignee.ID
	} else {
		// Auto-assign to CEO
		ceo, err := a.db.GetCEOAgent()
		if err == nil {
			issue.AssigneeAgentID = &ceo.ID
			assignee = ceo
		}
	}

	if err := a.db.CreateIssue(issue); err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	agent := agentFromContext(r.Context())
	a.db.LogActivity("create", "issue", key, ptrStr(agent), body.Title)

	// Wake assigned agent
	if assignee != nil && a.wake != nil {
		go a.wake(assignee, issue)
	}

	w.WriteHeader(http.StatusCreated)
	jsonOK(w, issue)
}

func (a *API) ListAgents(w http.ResponseWriter, r *http.Request) {
	agents, err := a.db.ListAgents()
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Return slim view for agent API
	type slim struct {
		ID            string `json:"id"`
		Slug          string `json:"slug"`
		Name          string `json:"name"`
		ArchetypeSlug string `json:"archetype_slug"`
	}
	result := make([]slim, len(agents))
	for i, ag := range agents {
		result[i] = slim{ag.ID, ag.Slug, ag.Name, ag.ArchetypeSlug}
	}
	jsonOK(w, result)
}

func (a *API) AgentMe(w http.ResponseWriter, r *http.Request) {
	agent := agentFromContext(r.Context())
	jsonOK(w, agent)
}

func (a *API) Usage(w http.ResponseWriter, r *http.Request) {
	agent := agentFromContext(r.Context())
	todayTokens, todayCost, totalTokens, totalCost, err := a.db.GetAgentUsage(agent.ID)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonOK(w, map[string]any{
		"today_tokens": todayTokens,
		"today_cost":   todayCost,
		"total_tokens": totalTokens,
		"total_cost":   totalCost,
	})
}

func (a *API) ResolveApproval(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var body struct {
		Status  string `json:"status"`
		Comment string `json:"comment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "invalid body", http.StatusBadRequest)
		return
	}
	if body.Status != "approved" && body.Status != "rejected" {
		jsonError(w, "status must be approved or rejected", http.StatusBadRequest)
		return
	}
	if err := a.db.ResolveApproval(id, body.Status, body.Comment); err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonOK(w, map[string]string{"status": body.Status})
}

func (a *API) wakeReviewerForIssue(agentID, issueKey string) {
	reviewer, err := a.db.GetReviewer(agentID)
	if err != nil {
		return
	}
	issue, err := a.db.GetIssue(issueKey)
	if err != nil {
		return
	}
	a.wake(reviewer, issue)
}

func (a *API) notifyOnComment(issueKey, author, body string) {
	if a.sse != nil {
		data, _ := json.Marshal(map[string]string{
			"issue_key": issueKey,
			"author":    author,
			"body":      body,
		})
		a.sse.Broadcast("comment", string(data))
	}
}

// helpers

func ptrStr(agent *models.Agent) *string {
	if agent == nil {
		return nil
	}
	return &agent.ID
}

func jsonOK(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
