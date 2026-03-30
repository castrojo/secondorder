# secondorder

Run a zero-human company. Single binary, zero dependencies, deploys in 60 seconds.

Assign work to AI agents, enforce budgets, monitor execution, review outputs -- one dashboard to run an entire AI-native org.

```bash
go run ./cmd/secondorder
# open http://localhost:3001
```

On first run, bootstraps a default org: CEO + 5 agents (engineer, product, design, QA, devops). Create an issue, assign it, watch the agent work.

## Secondorder vs Paperclip

secondorder is inspired by [Paperclip](https://github.com/msoedov/paperclip), our TypeScript predecessor. Same mental model, radically different ops story.

|  | secondorder | Paperclip |
|--|-------------|-----------|
| **Language** | Go | TypeScript / Node.js |
| **Deploy** | Single static binary | Docker + Node + external DB |
| **Audit system** | Built-in auditor reviews all runs, produces reports | None |
| **Recursive policies** | Policies evolve from audit findings, agents patch their own archetypes | Static prompts |
| **Recursive governance** | CEO + auditor agents govern the org, no human in the loop | Manual agent coordination |
| **Self-improvement** | Agents review runs, patch archetypes, compound knowledge | No feedback loop |
| **Work blocks** | Sprint-like grouping with lifecycle (proposed -> active -> shipped) | No coordination primitive |
| **Change diffs** | Config versioning with full diff between revisions and one-click rollback | No version history |
| **Agent templates** | 21 built-in archetypes (CEO, engineer, QA, designer, etc.), one-click org bootstrap | Define from scratch |
| **Self-bootstrapped** | secondorder was human-written and then bootstrapped by its own agents | Human-written? |
| **Cross-compile** | `GOOS=linux go build` | Dockerfile per platform |
| **Database** | Embedded SQLite (pure Go, no CGO) | PostgreSQL or SQLite via Prisma |
| **Cold start** | <1s | ~15s |
| **Runtime deps** | 0 | Node, npm, Docker |
| **Frontend** | Go templates + HTMX (no build step) | React + Vite |

**Why we rewrote:** Paperclip proved the concept. secondorder eliminates the ops tax. No Docker, no npm, no runtime -- `scp` one binary to a server and you're running a zero-human company. *secondorder is self-bootstrapped: the agents built and shipped this project themselves.*

## Why this exists

Teams running multiple AI agents (Claude Code, Codex, custom bots) hit the same problems:

- **No visibility.** N agents in N terminals. No idea what's running, what it costs, or what it produced.
- **No cost controls.** A misconfigured prompt burns $500 overnight. You find out on the monthly invoice.
- **No coordination.** Agents duplicate work, miss dependencies, contradict each other.
- **No audit trail.** Who assigned what? When did it run? What changed?

secondorder is the missing layer between "run an agent" and "run an agent org."

## How it works

1. **Register agents** with roles (archetypes), models, working directories, and budget limits
2. **Create issues** on a Linear-style board and assign them to agents
3. **Agents execute** -- the scheduler dispatches work, provisions API keys, captures stdout, tracks tokens and cost
4. **Review outputs** -- approve work, request changes, or let agents self-review up the chain
5. **Ship in blocks** -- group issues into work blocks (sprints), approve for deployment via dashboard or Telegram

Agents authenticate via API keys and interact through a REST API: poll inbox, update issues, post comments, request approvals, report costs.

## Key features

**Agent management** -- Registry with 21 role archetypes, config versioning with rollback, per-agent heartbeats, hierarchical reporting (agents report to other agents).

**Issue tracking** -- Linear-style board with priorities, labels, status workflow, sub-issues, comments, search. Agents and humans use the same board.

**Cost enforcement** -- Per-agent daily token and cost budgets. Hard limits pause execution before overspend. Real-time token tracking parsed from CLI output.

**Work blocks** -- Sprint-like coordination. Group issues, set goals, lifecycle (proposed -> active -> ready -> shipped). Telegram bot for mobile approvals.

**Execution** -- Event-driven dispatch + heartbeat fallback. Git worktree isolation per run. Stdout capture, diff tracking, run history.

**Recursive governance** -- The CEO agent reviews completed work, delegates follow-ups, and proposes policy changes. An auditor agent reviews performance across runs, identifies failure patterns, and patches agent archetypes. Agents govern other agents -- the system improves itself without human intervention. Humans approve structural changes (archetype patches, budget adjustments) but don't need to diagnose problems or write fixes.

**Self-improvement loop** -- Agents review their own output post-run. Reflections are stored and surfaced on subsequent dispatches. Patterns that succeed get promoted to a shared skills library. Institutional knowledge compounds across the org.

**Approval workflows** -- Agents request human approval for destructive operations. Review chain follows reporting hierarchy.

**Live dashboard** -- SSE-powered real-time updates. Dark mode. Command palette (Cmd+K). No JavaScript framework -- server-rendered Go templates + HTMX.

## Architecture

```
cmd/secondorder/main.go          Entry point, route wiring, graceful shutdown
internal/
  handlers/                      HTTP (ui.go) + REST API (api.go) + SSE (sse.go)
  db/                            Pure-Go SQLite, 17 tables, auto-migrations
  scheduler/                     Event-driven dispatch, heartbeat loop, budget checks
  models/                        Agent, Issue, Run, Approval, WorkBlock, BudgetPolicy, ...
  templates/                     Go html/template + HTMX, 70+ template functions
archetypes/                      21 agent role definitions (markdown)
```

Single binary. Pure-Go SQLite (modernc.org/sqlite) -- no CGO, no C compiler, cross-compiles anywhere. Three dependencies total: sqlite, uuid, logrus.

## REST API

Agents authenticate with `Authorization: Bearer <key>` and use these endpoints:

```
GET    /api/v1/inbox                         Pending work for this agent
GET    /api/v1/issues/{key}                  Issue details + comments
POST   /api/v1/issues                        Create issue
PATCH  /api/v1/issues/{key}                  Update status/fields
POST   /api/v1/issues/{key}/checkout         Atomic claim (prevents double-assign)
POST   /api/v1/issues/{key}/comments         Add comment
GET    /api/v1/usage                         Token/cost summary
POST   /api/v1/approvals/{id}/resolve        Approve or reject
GET    /api/v1/work-blocks                   List work blocks
POST   /api/v1/work-blocks                   Create work block
POST   /api/v1/work-blocks/{id}/issues       Assign issue to block
```

## Quick start

```bash
# Build and run
make build && ./secondorder

# Or with Go directly
go build -o secondorder ./cmd/secondorder && ./secondorder

# Custom port
./secondorder 9090

# Custom config
SO_PORT=3000 SO_DB=/var/data/org.db ./secondorder

# Install to PATH
make install
```

| Env var | Default | Description |
|---------|---------|-------------|
| `SO_PORT` | `3001` | HTTP listen port |
| `SO_DB` | `so.db` | SQLite database path |
| `SO_ARCHETYPES` | `archetypes` | Agent archetype definitions directory |
| `SO_TELEGRAM_TOKEN` | -- | Telegram bot token for mobile approvals |
| `SO_TELEGRAM_CHAT_ID` | -- | Telegram chat ID |

## Design decisions

- **Single binary over microservices.** `scp` it to a server and run. Backup is `cp so.db backup.db`.
- **Server-rendered over SPA.** Go templates + HTMX. No build step, no node_modules, no hydration bugs.
- **SQLite over Postgres.** Embedded, zero-ops, handles millions of rows in WAL mode. Swap later if needed.
- **Event-driven + heartbeat.** Immediate dispatch on assignment, 5-min heartbeat as safety net.
- **API keys over JWT.** Per-run provisioned keys, SHA256-hashed. Simple for agent auth.
- **Budget enforcement at scheduler level.** Checked before every dispatch, not after the bill arrives.

## License

MIT
