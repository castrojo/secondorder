# SO-87: Stall Detector — Flag In-Progress Issues with No Activity

## Problem

The CEO agent has historically detected stalled agent runs manually by noticing that an issue
has been `in_progress` for multiple sessions without new output. This detection is slow (takes
1+ sessions to notice) and forces the CEO into manual escalation loops.

**Trigger events identified in audit 82deb29f:**
- SO-47, SO-64, SO-65, SO-72 — each stalled 3–13 sessions before CEO noticed

---

## Solution

Automated stall detection via:
1. A **DB query** that identifies idle `in_progress` issues
2. A **background scheduler loop** that checks periodically and alerts the CEO
3. A **dashboard panel** showing flagged issues in amber
4. A **configurable threshold** via the `settings` table

---

## Architecture

### 1. DB Query — `GetStalledIssues(threshold time.Duration)`

**Location:** `internal/db/queries.go`

```go
func (d *DB) GetStalledIssues(threshold time.Duration) ([]models.StalledIssue, error)
```

**Logic:**
- Returns all issues where `status = 'in_progress'`
- AND `MAX(last_comment.created_at, issue.updated_at) < now - threshold`
- Uses a correlated subquery: `COALESCE(MAX(comments.created_at), issues.updated_at) < cutoff`
- Results ordered by last-activity ASC (oldest first)

**Cutoff definition:** The "last activity" timestamp is the maximum of:
- The most recent `comments.created_at` for that issue key
- `issues.updated_at` (when there are no comments, or comments are even older)

This means: adding a new comment on an issue resets its stall clock.

### 2. Model — `StalledIssue`

**Location:** `internal/models/models.go`

```go
type StalledIssue struct {
    Issue         models.Issue
    IdleDuration  time.Duration
    LastCommentAt *time.Time   // nil if no comments ever
}
```

### 3. Scheduler Loop — `StartStallDetectionLoop`

**Location:** `internal/scheduler/scheduler.go`

```go
func (s *Scheduler) StartStallDetectionLoop(interval, threshold time.Duration)
```

- Runs on a `time.NewTicker(interval)` (default: 30 min)
- Calls `GetStalledIssues(threshold)` on each tick
- For each newly-detected stalled issue:
  - Logs at `WARN` level (satisfies AC1 via log pipeline)
  - Posts a `Stall Detector` comment on the issue (satisfies AC2 — CEO sees it in inbox)
- Deduplicates notifications via an in-memory `map[string]bool` keyed by issue key
  (one notification per stall event per process lifetime — avoids spam)

**Wired in** `cmd/secondorder/main.go` after `StartAPIKeyExpiryLoop`:

```go
{
    thresholdHours := 4.0
    if val, err := database.GetSetting("stall_threshold_hours"); err == nil {
        var h float64
        if _, err2 := fmt.Sscanf(val, "%f", &h); err2 == nil && h > 0 {
            thresholdHours = h
        }
    }
    sched.StartStallDetectionLoop(30*time.Minute, time.Duration(thresholdHours*float64(time.Hour)))
}
```

### 4. Dashboard Panel — "Stalled Issues"

**Location:** `internal/templates/dashboard.html`

- Appears above the Recent Issues / Agents grid
- Only visible when `len(StalledIssues) > 0`
- Amber colour scheme (`border-amber-500/30`, `bg-amber-500/20`) for at-a-glance urgency
- Shows: issue key, title, assignee name, idle duration in hours
- HTMX auto-refresh on `run-complete` and `issue-changed` events

**Dashboard handler** (`internal/handlers/ui.go`) reads `stall_threshold_hours` from settings,
calls `GetStalledIssues(threshold)`, and passes `StalledIssues` + `StallThresholdH` to template.

### 5. Configuration — `stall_threshold_hours` setting

**Location:** `internal/db/migrations/019_stall_detection.sql`

```sql
INSERT OR IGNORE INTO settings (key, value) VALUES ('stall_threshold_hours', '4');
```

- Default: **4 hours** (≈ 2 sessions at 2h each) — matches AC3
- Configurable via the Settings UI (Settings → instance settings or direct DB update)
- Read at both startup (scheduler loop) and per-request (dashboard handler)

---

## Acceptance Criteria Mapping

| AC | Implementation |
|----|---------------|
| AC1: Stalled issues appear in a "stalled" view | Dashboard panel in amber, HTMX-refreshed |
| AC2: CEO receives notification | `Stall Detector` comment posted on stalled issue |
| AC3: Threshold configurable (default 4h) | `stall_threshold_hours` setting, migration 019 |

---

## Files Changed

| File | Change |
|------|--------|
| `internal/db/queries.go` | Added `GetStalledIssues()` and `scanIssueWithExtra()` |
| `internal/db/migrations/019_stall_detection.sql` | New migration seeds `stall_threshold_hours=4` |
| `internal/db/db_test.go` | 7 new tests for `GetStalledIssues` + settings |
| `internal/models/models.go` | Added `StalledIssue` struct |
| `internal/scheduler/scheduler.go` | Added `StartStallDetectionLoop` and `runStallDetection` |
| `internal/handlers/ui.go` | Dashboard handler reads and passes stalled issues |
| `internal/templates/dashboard.html` | Stalled Issues panel (amber, HTMX-refreshed) |
| `cmd/secondorder/main.go` | Wires `StartStallDetectionLoop` with configurable threshold |

---

## Design Decisions

**Why in-memory deduplication (not DB)?**  
The notification is meant to alert the CEO once per stall event, not on every 30-min tick.
A DB-persisted `stall_notifications` table would add schema complexity for minimal gain. If
the process restarts and the issue is still stalled, a fresh notification is appropriate.

**Why `updated_at` as fallback instead of `started_at`?**  
`updated_at` reflects any activity (status changes, description edits) while `started_at` is set
once. Using the more recent timestamp makes the stall detector conservative (less false positives).

**Why not a separate `/stalled` route?**  
Embedding the stalled panel in the CEO dashboard keeps the signal in the decision-making context
where the CEO acts. A separate page would require active navigation.

**Why `time.Duration.Hours()` instead of a template helper?**  
`time.Duration.Hours()` has a value receiver and is accessible directly from Go templates via
`.IdleDuration.Hours`, avoiding the need to add a new template function.
