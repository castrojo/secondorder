-- agents
CREATE TABLE IF NOT EXISTS agents (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    archetype_slug TEXT NOT NULL DEFAULT 'other',
    model TEXT NOT NULL DEFAULT 'sonnet',
    working_dir TEXT NOT NULL DEFAULT '.',
    max_turns INTEGER NOT NULL DEFAULT 50,
    timeout_sec INTEGER NOT NULL DEFAULT 600,
    heartbeat_enabled INTEGER NOT NULL DEFAULT 0,
    heartbeat_cron TEXT NOT NULL DEFAULT '',
    chrome_enabled INTEGER NOT NULL DEFAULT 0,
    reports_to TEXT,
    review_agent_id TEXT,
    active INTEGER NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- issues
CREATE TABLE IF NOT EXISTS issues (
    id TEXT PRIMARY KEY,
    key TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'todo',
    priority INTEGER NOT NULL DEFAULT 0,
    assignee_agent_id TEXT,
    parent_issue_key TEXT,
    work_block_id TEXT,
    started_at DATETIME,
    completed_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (assignee_agent_id) REFERENCES agents(id),
    FOREIGN KEY (parent_issue_key) REFERENCES issues(key)
);

-- runs
CREATE TABLE IF NOT EXISTS runs (
    id TEXT PRIMARY KEY,
    agent_id TEXT NOT NULL,
    issue_key TEXT,
    mode TEXT NOT NULL DEFAULT 'task',
    status TEXT NOT NULL DEFAULT 'running',
    stdout TEXT NOT NULL DEFAULT '',
    diff TEXT NOT NULL DEFAULT '',
    input_tokens INTEGER NOT NULL DEFAULT 0,
    output_tokens INTEGER NOT NULL DEFAULT 0,
    cache_read_tokens INTEGER NOT NULL DEFAULT 0,
    cache_create_tokens INTEGER NOT NULL DEFAULT 0,
    total_cost_usd REAL NOT NULL DEFAULT 0,
    started_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (agent_id) REFERENCES agents(id)
);

-- comments
CREATE TABLE IF NOT EXISTS comments (
    id TEXT PRIMARY KEY,
    issue_key TEXT NOT NULL,
    agent_id TEXT,
    author TEXT NOT NULL DEFAULT 'Board',
    body TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (issue_key) REFERENCES issues(key)
);

-- approvals
CREATE TABLE IF NOT EXISTS approvals (
    id TEXT PRIMARY KEY,
    issue_key TEXT NOT NULL,
    requested_by TEXT NOT NULL,
    reviewer_id TEXT,
    status TEXT NOT NULL DEFAULT 'pending',
    comment TEXT NOT NULL DEFAULT '',
    resolved_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- api_keys
CREATE TABLE IF NOT EXISTS api_keys (
    id TEXT PRIMARY KEY,
    agent_id TEXT NOT NULL,
    key_hash TEXT NOT NULL UNIQUE,
    prefix TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    revoked_at DATETIME,
    FOREIGN KEY (agent_id) REFERENCES agents(id)
);

-- labels
CREATE TABLE IF NOT EXISTS labels (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    color TEXT NOT NULL DEFAULT '#6366f1'
);

CREATE TABLE IF NOT EXISTS issue_labels (
    issue_id TEXT NOT NULL,
    label_id TEXT NOT NULL,
    PRIMARY KEY (issue_id, label_id)
);

-- activity_log
CREATE TABLE IF NOT EXISTS activity_log (
    id TEXT PRIMARY KEY,
    action TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    entity_id TEXT NOT NULL,
    agent_id TEXT,
    details TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- cost_events
CREATE TABLE IF NOT EXISTS cost_events (
    id TEXT PRIMARY KEY,
    run_id TEXT NOT NULL,
    agent_id TEXT NOT NULL,
    input_tokens INTEGER NOT NULL DEFAULT 0,
    output_tokens INTEGER NOT NULL DEFAULT 0,
    total_cost_usd REAL NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- budget_policies
CREATE TABLE IF NOT EXISTS budget_policies (
    id TEXT PRIMARY KEY,
    agent_id TEXT NOT NULL UNIQUE,
    daily_token_limit INTEGER NOT NULL DEFAULT 0,
    daily_cost_limit REAL NOT NULL DEFAULT 0,
    active INTEGER NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (agent_id) REFERENCES agents(id)
);

-- secrets
CREATE TABLE IF NOT EXISTS secrets (
    id TEXT PRIMARY KEY,
    key TEXT NOT NULL UNIQUE,
    value TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- skills
CREATE TABLE IF NOT EXISTS skills (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    agent_id TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (agent_id) REFERENCES agents(id)
);

-- agent_config_revisions
CREATE TABLE IF NOT EXISTS agent_config_revisions (
    id TEXT PRIMARY KEY,
    agent_id TEXT NOT NULL,
    config TEXT NOT NULL,
    changed_by TEXT NOT NULL DEFAULT 'system',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (agent_id) REFERENCES agents(id)
);

-- run_events
CREATE TABLE IF NOT EXISTS run_events (
    id TEXT PRIMARY KEY,
    run_id TEXT NOT NULL,
    event_type TEXT NOT NULL,
    data TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (run_id) REFERENCES runs(id)
);

-- work_blocks
CREATE TABLE IF NOT EXISTS work_blocks (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    goal TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'active',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME
);

-- indexes
CREATE INDEX IF NOT EXISTS idx_issues_status ON issues(status);
CREATE INDEX IF NOT EXISTS idx_issues_assignee ON issues(assignee_agent_id);
CREATE INDEX IF NOT EXISTS idx_issues_parent ON issues(parent_issue_key);
CREATE INDEX IF NOT EXISTS idx_runs_agent ON runs(agent_id);
CREATE INDEX IF NOT EXISTS idx_runs_issue ON runs(issue_key);
CREATE INDEX IF NOT EXISTS idx_comments_issue ON comments(issue_key);
CREATE INDEX IF NOT EXISTS idx_api_keys_hash ON api_keys(key_hash);
CREATE INDEX IF NOT EXISTS idx_cost_events_agent ON cost_events(agent_id);
CREATE INDEX IF NOT EXISTS idx_cost_events_created ON cost_events(created_at);

-- schema_migrations tracking
CREATE TABLE IF NOT EXISTS schema_migrations (
    version INTEGER PRIMARY KEY,
    applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
