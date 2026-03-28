package scheduler

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/msoedov/thelastorg/internal/db"
	"github.com/msoedov/thelastorg/internal/models"
	log "github.com/sirupsen/logrus"
)

type Scheduler struct {
	db           *db.DB
	port         int
	archetypesDir string
	mu           sync.Mutex
	running      map[string]context.CancelFunc // runID -> cancel
	wg           sync.WaitGroup
	stopped      bool
	onRunComplete func(run *models.Run)
	onComment     func(issueKey, author, body string)
}

func New(database *db.DB, port int, archetypesDir string) *Scheduler {
	return &Scheduler{
		db:            database,
		port:          port,
		archetypesDir: archetypesDir,
		running:       make(map[string]context.CancelFunc),
	}
}

func (s *Scheduler) SetOnRunComplete(fn func(run *models.Run)) {
	s.onRunComplete = fn
}

func (s *Scheduler) SetOnComment(fn func(issueKey, author, body string)) {
	s.onComment = fn
}

func (s *Scheduler) Stop() {
	s.mu.Lock()
	s.stopped = true
	for _, cancel := range s.running {
		cancel()
	}
	s.mu.Unlock()
	s.wg.Wait()
	log.Info("scheduler: all agents stopped")
}

// WakeAgent spawns an agent for a specific issue (event-driven)
func (s *Scheduler) WakeAgent(agent *models.Agent, issue *models.Issue) {
	s.mu.Lock()
	if s.stopped {
		s.mu.Unlock()
		return
	}
	s.mu.Unlock()

	prompt := s.buildTaskPrompt(agent, issue)
	s.spawnAgent(agent, issue.Key, "task", prompt)
}

// WakeAgentHeartbeat spawns a heartbeat run for the agent
func (s *Scheduler) WakeAgentHeartbeat(agent *models.Agent) {
	s.mu.Lock()
	if s.stopped {
		s.mu.Unlock()
		return
	}
	s.mu.Unlock()

	prompt := s.buildHeartbeatPrompt(agent)
	s.spawnAgent(agent, "", "heartbeat", prompt)
}

// WakeReviewer finds and wakes the appropriate reviewer for an agent's completed issue
func (s *Scheduler) WakeReviewer(agentID, issueKey string) {
	reviewer, err := s.db.GetReviewer(agentID)
	if err != nil {
		log.WithError(err).Error("scheduler: failed to find reviewer")
		return
	}
	issue, err := s.db.GetIssue(issueKey)
	if err != nil {
		log.WithError(err).Error("scheduler: failed to get issue for reviewer")
		return
	}
	s.WakeAgent(reviewer, issue)
}

func (s *Scheduler) spawnAgent(agent *models.Agent, issueKey, mode, prompt string) {
	runID := uuid.New().String()

	run := &models.Run{
		ID:        runID,
		AgentID:   agent.ID,
		Mode:      mode,
		Status:    models.RunStatusRunning,
		StartedAt: time.Now(),
		CreatedAt: time.Now(),
	}
	if issueKey != "" {
		run.IssueKey = &issueKey
	}

	if err := s.db.CreateRun(run); err != nil {
		log.WithError(err).Error("scheduler: failed to create run")
		return
	}

	// Provision API key
	rawKey, err := s.provisionAPIKey(agent.ID)
	if err != nil {
		log.WithError(err).Error("scheduler: failed to provision API key")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(agent.TimeoutSec)*time.Second)

	s.mu.Lock()
	s.running[runID] = cancel
	s.mu.Unlock()

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer func() {
			s.mu.Lock()
			delete(s.running, runID)
			s.mu.Unlock()
			cancel()
		}()

		logEntry := log.WithFields(log.Fields{
			"agent":     agent.Name,
			"archetype": agent.ArchetypeSlug,
			"run_id":    runID,
			"model":     agent.Model,
			"mode":      mode,
			"issue_key": issueKey,
		})
		logEntry.Info("scheduler: spawning agent")

		startTime := time.Now()
		stdout, err := s.execClaude(ctx, agent, rawKey, runID, issueKey, prompt)
		elapsed := time.Since(startTime)

		status := models.RunStatusCompleted
		if err != nil {
			status = models.RunStatusFailed
			if ctx.Err() == context.DeadlineExceeded {
				status = models.RunStatusFailed
				logEntry.Warn("scheduler: agent timed out")
			} else if ctx.Err() == context.Canceled {
				status = models.RunStatusCancelled
			} else {
				logEntry.WithError(err).Error("scheduler: agent failed")
			}
		}

		// Parse token usage from stream-json output
		tokens := parseTokenUsage(stdout)

		// Capture git diff
		diff := captureGitDiff(agent.WorkingDir)

		completedAt := time.Now()
		completedRun := models.Run{
			InputTokens:       tokens.InputTokens,
			OutputTokens:      tokens.OutputTokens,
			CacheReadTokens:   tokens.CacheReadTokens,
			CacheCreateTokens: tokens.CacheCreateTokens,
			TotalCostUSD:      tokens.TotalCostUSD,
		}
		if err := s.db.CompleteRun(runID, status, stdout, diff, completedRun); err != nil {
			logEntry.WithError(err).Error("scheduler: failed to complete run")
		}

		// Record cost event
		if tokens.TotalCostUSD > 0 {
			s.db.CreateCostEvent(&models.CostEvent{
				ID:           uuid.New().String(),
				RunID:        runID,
				AgentID:      agent.ID,
				InputTokens:  tokens.InputTokens,
				OutputTokens: tokens.OutputTokens,
				TotalCostUSD: tokens.TotalCostUSD,
				CreatedAt:    time.Now(),
			})
		}

		logEntry.WithFields(log.Fields{
			"status":        status,
			"elapsed":       elapsed.Round(time.Second),
			"cost_usd":      fmt.Sprintf("%.4f", tokens.TotalCostUSD),
			"input_tokens":  tokens.InputTokens,
			"output_tokens": tokens.OutputTokens,
		}).Info("scheduler: agent completed")

		if s.onRunComplete != nil {
			finalRun := &models.Run{
				ID:          runID,
				AgentID:     agent.ID,
				IssueKey:    run.IssueKey,
				Mode:        mode,
				Status:      status,
				CompletedAt: &completedAt,
			}
			s.onRunComplete(finalRun)
		}
	}()
}

func (s *Scheduler) execClaude(ctx context.Context, agent *models.Agent, apiKey, runID, issueKey, prompt string) (string, error) {
	args := []string{
		"--print",
		"-p", prompt,
		"--output-format", "stream-json",
		"--verbose",
		"--dangerously-skip-permissions",
		"--max-turns", fmt.Sprintf("%d", agent.MaxTurns),
		"--model", agent.Model,
	}

	archetypeFile := filepath.Join(s.archetypesDir, agent.ArchetypeSlug+".md")
	if _, err := os.Stat(archetypeFile); err == nil {
		args = append(args, "--append-system-prompt-file", archetypeFile)
	}

	artifactDir := filepath.Join(agent.WorkingDir, "artifact-docs")
	if info, err := os.Stat(artifactDir); err == nil && info.IsDir() {
		args = append(args, "--add-dir", artifactDir)
	}

	if agent.ChromeEnabled {
		args = append(args, "--chrome")
	}

	cmd := exec.CommandContext(ctx, "claude", args...)
	cmd.Dir = agent.WorkingDir
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("THELASTORG_AGENT_ID=%s", agent.ID),
		fmt.Sprintf("THELASTORG_AGENT_NAME=%s", agent.Name),
		fmt.Sprintf("THELASTORG_RUN_ID=%s", runID),
		fmt.Sprintf("THELASTORG_API_URL=http://localhost:%d", s.port),
		fmt.Sprintf("THELASTORG_ISSUE_KEY=%s", issueKey),
		fmt.Sprintf("THELASTORG_ARTIFACT_DOCS=%s", filepath.Join(agent.WorkingDir, "artifact-docs")),
		fmt.Sprintf("TLO_API_KEY=%s", apiKey),
	)

	// Use liveWriter to stream stdout to DB
	lw := &liveWriter{
		db:       s.db,
		runID:    runID,
		interval: 2 * time.Second,
	}
	cmd.Stdout = lw
	cmd.Stderr = lw

	err := cmd.Run()
	lw.Flush()
	return lw.String(), err
}

func (s *Scheduler) provisionAPIKey(agentID string) (string, error) {
	// Revoke existing keys
	s.db.RevokeAPIKeys(agentID)

	// Generate new key
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	rawKey := "tlo_" + hex.EncodeToString(raw)

	hash := sha256.Sum256([]byte(rawKey))
	keyHash := hex.EncodeToString(hash[:])
	prefix := rawKey[:12]

	if err := s.db.CreateAPIKey(agentID, keyHash, prefix); err != nil {
		return "", err
	}

	return rawKey, nil
}

func (s *Scheduler) buildTaskPrompt(agent *models.Agent, issue *models.Issue) string {
	comments, _ := s.db.ListComments(issue.Key)

	var commentBlock string
	start := 0
	if len(comments) > 5 {
		start = len(comments) - 5
	}
	for _, c := range comments[start:] {
		commentBlock += fmt.Sprintf("- %s (%s): %s\n", c.Author, c.CreatedAt.Format("Jan 2 15:04"), c.Body)
	}

	apiRef := workerAPIRef
	rules := workerRules
	wbContext := ""
	if agent.ArchetypeSlug == "ceo" {
		apiRef = s.buildCEOAPIRef()
		rules = ceoRules
		wbContext = s.buildWorkBlockContext()
	}

	return fmt.Sprintf(`ISSUE: %s
TITLE: %s
DESCRIPTION:
%s
%s
RECENT COMMENTS:
%s
%s

%s

BASE_URL: http://localhost:%d
`, issue.Key, issue.Title, issue.Description, wbContext, commentBlock, rules, apiRef, s.port)
}

func (s *Scheduler) buildHeartbeatPrompt(agent *models.Agent) string {
	inbox, _ := s.db.GetAgentInbox(agent.ID)

	var issueBlock string
	for _, i := range inbox {
		issueBlock += fmt.Sprintf("- [%s] %s (status: %s, priority: %d)\n", i.Key, i.Title, i.Status, i.Priority)
	}

	apiRef := workerAPIRef
	rules := workerRules
	if agent.ArchetypeSlug == "ceo" {
		apiRef = s.buildCEOAPIRef()
		rules = ceoRules

		approvals, _ := s.db.ListPendingApprovals()
		if len(approvals) > 0 {
			issueBlock += "\nPENDING REVIEWS:\n"
			for _, a := range approvals {
				issueBlock += fmt.Sprintf("- Approval %s for issue %s (requested by: %s)\n", a.ID, a.IssueKey, a.RequestedBy)
			}
		}

		issueBlock += s.buildWorkBlockContext()
	}

	return fmt.Sprintf(`HEARTBEAT CHECK - Review your inbox and take action on any pending items.

YOUR INBOX:
%s

%s

%s

BASE_URL: http://localhost:%d
`, issueBlock, rules, apiRef, s.port)
}

func (s *Scheduler) buildWorkBlockContext() string {
	wb, err := s.db.GetActiveWorkBlock()
	if err != nil {
		return ""
	}
	issues, _ := s.db.ListWorkBlockIssues(wb.ID)
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("\nACTIVE WORK BLOCK: %s (id: %s)\nGoal: %s\n", wb.Title, wb.ID, wb.Goal))
	if len(issues) > 0 {
		buf.WriteString("Block Issues:\n")
		for _, i := range issues {
			buf.WriteString(fmt.Sprintf("  - [%s] %s (status: %s)\n", i.Key, i.Title, i.Status))
		}
	}
	return buf.String()
}

func (s *Scheduler) buildCEOAPIRef() string {
	agents, _ := s.db.ListAgents()
	var roster string
	for _, a := range agents {
		if a.ArchetypeSlug == "ceo" {
			continue
		}
		roster += fmt.Sprintf("  %s (slug: %s, role: %s)\n", a.Name, a.Slug, a.ArchetypeSlug)
	}
	return fmt.Sprintf(ceoAPIRef, roster)
}

// StartHeartbeatLoop runs heartbeat checks on a timer (safety net)
func (s *Scheduler) StartHeartbeatLoop(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				s.mu.Lock()
				if s.stopped {
					s.mu.Unlock()
					return
				}
				s.mu.Unlock()
				s.runHeartbeats()
			}
		}
	}()
}

func (s *Scheduler) runHeartbeats() {
	agents, err := s.db.ListAgents()
	if err != nil {
		log.WithError(err).Error("scheduler: failed to list agents for heartbeat")
		return
	}
	for i := range agents {
		a := &agents[i]
		if !a.Active || !a.HeartbeatEnabled {
			continue
		}
		// Check if agent is over budget
		over, _ := s.db.IsAgentOverBudget(a.ID)
		if over {
			log.WithField("agent", a.Name).Warn("scheduler: agent over budget, skipping heartbeat")
			continue
		}
		s.WakeAgentHeartbeat(a)
	}
}

const workerRules = `RULES:
- You are fully autonomous. Do NOT ask questions interactively. Do NOT wait for human input.
- If you have a question or need clarification, post it as a comment on the ticket and mark the issue "blocked".
- Do NOT request approvals. Just do the work and mark done.
- Always checkout the issue first, then do the work, then update status.
- Write any documentation to the artifact-docs/ folder.`

const workerAPIRef = `TLO API (Authorization: Bearer $TLO_API_KEY):
  GET    $THELASTORG_API_URL/api/v1/inbox                              - your assigned issues
  GET    $THELASTORG_API_URL/api/v1/issues/{key}                       - issue detail + comments
  POST   $THELASTORG_API_URL/api/v1/issues/{key}/checkout              - claim issue
  PATCH  $THELASTORG_API_URL/api/v1/issues/{key}                       - update status + comment
  POST   $THELASTORG_API_URL/api/v1/issues/{key}/comments              - add comment
  POST   $THELASTORG_API_URL/api/v1/issues                             - create sub-issue
  GET    $THELASTORG_API_URL/api/v1/usage                              - your token/cost usage`

const ceoRules = `RULES:
- You are fully autonomous. Do NOT ask questions interactively.
- Do NOT do implementation work yourself. Always delegate by creating sub-issues with assignee_slug and parent_issue_key.
- Break complex tasks into clear sub-issues with acceptance criteria.
- After delegating, mark the parent as "in_progress" and comment your plan.
- When reviewing completed work: approve, request changes via comment, or reassign.
- If blocked, post a comment and mark "blocked".
- If there is an active work block, focus your work on its goal. Assign relevant issues to the block.
- When all issues in a block are done, mark the block as "ready" via PATCH.
- To start new work, propose a work block first. A human must approve it before it becomes active.
- Only one work block can be active or proposed at a time.`

const ceoAPIRef = `TLO API (Authorization: Bearer $TLO_API_KEY):
  GET    $THELASTORG_API_URL/api/v1/inbox                              - your assigned issues
  GET    $THELASTORG_API_URL/api/v1/issues/{key}                       - issue detail + comments
  PATCH  $THELASTORG_API_URL/api/v1/issues/{key}                       - update status + comment
  POST   $THELASTORG_API_URL/api/v1/issues/{key}/comments              - add comment
  POST   $THELASTORG_API_URL/api/v1/issues                             - create & assign: {"title":"...","assignee_slug":"...","parent_issue_key":"..."}
  GET    $THELASTORG_API_URL/api/v1/agents                             - list team (slug, name, archetype)
  POST   $THELASTORG_API_URL/api/v1/approvals/{id}/resolve             - review: {"status":"approved","comment":"..."}
  GET    $THELASTORG_API_URL/api/v1/work-blocks                        - list work blocks
  GET    $THELASTORG_API_URL/api/v1/work-blocks/{id}                   - block detail + issues + metrics
  POST   $THELASTORG_API_URL/api/v1/work-blocks                        - propose block: {"title":"...","goal":"..."}
  PATCH  $THELASTORG_API_URL/api/v1/work-blocks/{id}                   - update status: {"status":"ready"}
  POST   $THELASTORG_API_URL/api/v1/work-blocks/{id}/issues            - assign issue: {"issue_key":"TLO-5"}
  DELETE $THELASTORG_API_URL/api/v1/work-blocks/{id}/issues/{key}      - unassign issue

Your team:
%s`

// liveWriter buffers stdout and flushes to DB periodically
type liveWriter struct {
	db       *db.DB
	runID    string
	interval time.Duration
	mu       sync.Mutex
	buf      strings.Builder
	lastFlush time.Time
}

func (w *liveWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	n, err := w.buf.Write(p)
	if time.Since(w.lastFlush) >= w.interval {
		w.db.UpdateRunStdout(w.runID, w.buf.String())
		w.lastFlush = time.Now()
	}
	return n, err
}

func (w *liveWriter) Flush() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.db.UpdateRunStdout(w.runID, w.buf.String())
}

func (w *liveWriter) String() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buf.String()
}

type tokenUsage struct {
	InputTokens       int64
	OutputTokens      int64
	CacheReadTokens   int64
	CacheCreateTokens int64
	TotalCostUSD      float64
}

func parseTokenUsage(stdout string) tokenUsage {
	var usage tokenUsage
	for _, line := range strings.Split(stdout, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "{") {
			continue
		}
		var msg struct {
			Type   string `json:"type"`
			Result struct {
				InputTokens              int64   `json:"input_tokens"`
				OutputTokens             int64   `json:"output_tokens"`
				CacheReadInputTokens     int64   `json:"cache_read_input_tokens"`
				CacheCreationInputTokens int64   `json:"cache_creation_input_tokens"`
				TotalCostUSD             float64 `json:"total_cost_usd"`
			} `json:"result"`
		}
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			continue
		}
		if msg.Type == "result" {
			usage.InputTokens = msg.Result.InputTokens
			usage.OutputTokens = msg.Result.OutputTokens
			usage.CacheReadTokens = msg.Result.CacheReadInputTokens
			usage.CacheCreateTokens = msg.Result.CacheCreationInputTokens
			usage.TotalCostUSD = msg.Result.TotalCostUSD
		}
	}
	return usage
}

func captureGitDiff(workingDir string) string {
	cmd := exec.Command("git", "diff", "HEAD")
	cmd.Dir = workingDir
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	diff := string(out)
	// Cap at 100KB
	if len(diff) > 100*1024 {
		diff = diff[:100*1024] + "\n... (truncated at 100KB)"
	}
	return diff
}
