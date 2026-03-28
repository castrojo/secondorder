package models

import "time"

// Issue statuses
const (
	StatusTodo       = "todo"
	StatusInProgress = "in_progress"
	StatusInReview   = "in_review"
	StatusDone       = "done"
	StatusBlocked    = "blocked"
	StatusCancelled  = "cancelled"
)

// Run statuses
const (
	RunStatusRunning   = "running"
	RunStatusCompleted = "completed"
	RunStatusFailed    = "failed"
	RunStatusCancelled = "cancelled"
)

type Agent struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Slug             string    `json:"slug"`
	ArchetypeSlug    string    `json:"archetype_slug"`
	Model            string    `json:"model"`
	WorkingDir       string    `json:"working_dir"`
	MaxTurns         int       `json:"max_turns"`
	TimeoutSec       int       `json:"timeout_sec"`
	HeartbeatEnabled bool      `json:"heartbeat_enabled"`
	HeartbeatCron    string    `json:"heartbeat_cron"`
	ChromeEnabled    bool      `json:"chrome_enabled"`
	ReportsTo        *string   `json:"reports_to"`
	ReviewAgentID    *string   `json:"review_agent_id"`
	Active           bool      `json:"active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type Issue struct {
	ID              string     `json:"id"`
	Key             string     `json:"key"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	Status          string     `json:"status"`
	Priority        int        `json:"priority"`
	AssigneeAgentID *string    `json:"assignee_agent_id"`
	ParentIssueKey  *string    `json:"parent_issue_key"`
	WorkBlockID     *string    `json:"work_block_id"`
	AssigneeName    string     `json:"assignee_name,omitempty"`
	StartedAt       *time.Time `json:"started_at"`
	CompletedAt     *time.Time `json:"completed_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type Run struct {
	ID                string     `json:"id"`
	AgentID           string     `json:"agent_id"`
	IssueKey          *string    `json:"issue_key"`
	Mode              string     `json:"mode"`
	Status            string     `json:"status"`
	Stdout            string     `json:"stdout"`
	Diff              string     `json:"diff"`
	InputTokens       int64      `json:"input_tokens"`
	OutputTokens      int64      `json:"output_tokens"`
	CacheReadTokens   int64      `json:"cache_read_tokens"`
	CacheCreateTokens int64      `json:"cache_create_tokens"`
	TotalCostUSD      float64    `json:"total_cost_usd"`
	StartedAt         time.Time  `json:"started_at"`
	CompletedAt       *time.Time `json:"completed_at"`
	CreatedAt         time.Time  `json:"created_at"`
}

type Comment struct {
	ID        string    `json:"id"`
	IssueKey  string    `json:"issue_key"`
	AgentID   *string   `json:"agent_id"`
	Author    string    `json:"author"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

type Approval struct {
	ID          string     `json:"id"`
	IssueKey    string     `json:"issue_key"`
	RequestedBy string    `json:"requested_by"`
	ReviewerID  *string   `json:"reviewer_id"`
	Status      string    `json:"status"`
	Comment     string    `json:"comment"`
	ResolvedAt  *time.Time `json:"resolved_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

type APIKey struct {
	ID        string     `json:"id"`
	AgentID   string     `json:"agent_id"`
	KeyHash   string     `json:"-"`
	Prefix    string     `json:"prefix"`
	CreatedAt time.Time  `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at"`
}

type Label struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type IssueLabel struct {
	IssueID string `json:"issue_id"`
	LabelID string `json:"label_id"`
}

type ActivityLog struct {
	ID         string    `json:"id"`
	Action     string    `json:"action"`
	EntityType string    `json:"entity_type"`
	EntityID   string    `json:"entity_id"`
	AgentID    *string   `json:"agent_id"`
	Details    string    `json:"details"`
	CreatedAt  time.Time `json:"created_at"`
}

type CostEvent struct {
	ID           string    `json:"id"`
	RunID        string    `json:"run_id"`
	AgentID      string    `json:"agent_id"`
	InputTokens  int64     `json:"input_tokens"`
	OutputTokens int64     `json:"output_tokens"`
	TotalCostUSD float64   `json:"total_cost_usd"`
	CreatedAt    time.Time `json:"created_at"`
}

type BudgetPolicy struct {
	ID              string    `json:"id"`
	AgentID         string    `json:"agent_id"`
	DailyTokenLimit int64     `json:"daily_token_limit"`
	DailyCostLimit  float64   `json:"daily_cost_limit"`
	Active          bool      `json:"active"`
	CreatedAt       time.Time `json:"created_at"`
}

type Secret struct {
	ID        string    `json:"id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Skill struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	AgentID     string    `json:"agent_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type AgentConfigRevision struct {
	ID        string    `json:"id"`
	AgentID   string    `json:"agent_id"`
	Config    string    `json:"config"`
	ChangedBy string    `json:"changed_by"`
	CreatedAt time.Time `json:"created_at"`
}

type RunEvent struct {
	ID        string    `json:"id"`
	RunID     string    `json:"run_id"`
	EventType string    `json:"event_type"`
	Data      string    `json:"data"`
	CreatedAt time.Time `json:"created_at"`
}

type WorkBlock struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Goal        string     `json:"goal"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

type DashboardStats struct {
	TotalAgents    int     `json:"total_agents"`
	ActiveAgents   int     `json:"active_agents"`
	TotalIssues    int     `json:"total_issues"`
	OpenIssues     int     `json:"open_issues"`
	RunningRuns    int     `json:"running_runs"`
	TotalCostToday float64 `json:"total_cost_today"`
}
