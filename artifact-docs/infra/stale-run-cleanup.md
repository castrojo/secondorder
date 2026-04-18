# Stale Run Cleanup on Startup (SO-21)

## Problem

When the secondorder service is killed mid-run, rows in the `runs` table remain
stuck at `status=running` indefinitely. This causes:

- Activity feed clutter (running jobs that will never complete)
- The stuck-issue recovery loop attempting to re-wake agents for runs that are
  still marked "running" but have no live process

## Solution

`internal/db/queries.go` — `CleanupStaleRuns(cutoff time.Duration)`

- Queries for all `runs` where `status='running'` AND `started_at < (now - cutoff)`
- Returns the slice of matching run IDs so callers can emit per-run log lines
- Updates those runs to `status='failed'` with `completed_at=datetime('now')` in one statement
- A cutoff of 0 matches all running runs regardless of age (used by the backward-compat wrapper)

`internal/scheduler/scheduler.go` — `RecoverStuckIssues()`

- Calls `CleanupStaleRuns(10 * time.Minute)` **before** heartbeat loops and HTTP listener start
- Logs `scheduler: cleaned up stale run <run_id>` for every affected run (AC2)
- Logs a summary `scheduler: stale run cleanup complete count=N` when N > 0

### Cutoff: 10 minutes

Chosen per CEO recommendation. Runs younger than 10 minutes may belong to a
legitimately running process (e.g. a recently restarted concurrent worker).
Runs older than 10 minutes from a crashed process are safe to fail.

## Backward compatibility

`MarkStaleRunsFailed()` is preserved and now delegates to `CleanupStaleRuns(0)`,
maintaining the zero-cutoff ("fail everything") behaviour used in existing tests.

## Startup order (main.go)

```
RecoverStuckIssues()      ← CleanupStaleRuns runs here (line 305)
StartHeartbeatLoop()      ← scheduler begins accepting work (line 310)
StartAPIKeyExpiryLoop()   ← (line 313)
srv.ListenAndServe()      ← HTTP open (line 325)
```

AC3 is satisfied: cleanup completes synchronously before any new work is accepted.

## Tests added (`internal/db/db_test.go`)

| Test | What it covers |
|------|---------------|
| `TestCleanupStaleRunsAllRunning` | cutoff=0 fails all running runs, returns all IDs |
| `TestCleanupStaleRunsCutoffFiltersRecent` | 10-min cutoff spares recent run, fails old one |
| `TestCleanupStaleRunsNoneRunning` | no running rows → empty slice, no error |
| `TestMarkStaleRunsFailedBackwardCompat` | MarkStaleRunsFailed wrapper still works |

All pass: `go test ./...` — 0 failures.
