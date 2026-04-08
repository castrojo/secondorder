# Tech Spec: SO-69 Fix activity chart statistics display

## Overview
The activity chart was undercounting or dropping statistics because a large portion of `activity_log.created_at` values were stored in a legacy Go string format that SQLite `DATE()` could not parse.

This fix makes the chart query resilient to both timestamp formats and adds explicit cadence/stat summary cards so the activity page always surfaces readable statistics alongside the chart.

## Root Cause
- `GetDailyActivityStats` joined daily buckets with `DATE(activity_log.created_at)`.
- Legacy rows stored `created_at` like `2026-04-03 18:46:48.170698 +0000 UTC`, which SQLite does not parse with `DATE()`.
- Those rows were skipped by the join, so the chart and totals missed real activity.
- Legacy completion events also used `details = 'completed'`, while the chart only counted `details = 'done'`.

## Changes

### 1. Robust activity date bucketing
- Updated `internal/db/queries.go` to join on `COALESCE(DATE(created_at), SUBSTR(created_at, 1, 10))`.
- This preserves support for normal SQLite-friendly timestamps and also captures legacy text timestamps by their ISO date prefix.

### 2. Backward-compatible completion counting
- Updated completion aggregation to count both legacy `completed` updates and current `done` transitions.

### 3. Visible cadence/statistics summary
- Added an `ActivityOverview` model and server-side overview computation in `internal/handlers/ui.go`.
- The activity page now shows:
  - task cadence (`active days / window`)
  - current streak
  - average actions per day
  - busiest day
  - tasks completed

### 4. UI refinements
- Updated `internal/templates/activity.html` to render the overview cards above the chart.
- Added an `Activity totals` label above the per-action totals row to make the statistics section clearer.

## Verification
- `go test ./internal/db ./internal/handlers`
- Added a DB regression test for legacy timestamp rows.
- Added handler coverage for the computed overview stats and rendered activity page sections.
