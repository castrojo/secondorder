# Tech Spec: SO-84 Fix activity chart statistics display

## Overview
Update the activity chart to include all activity types in the daily stats calculation.
Previously, the chart only tracked `create`, `update`, and `checkout` actions.

## Changes

### 1. Model Updates (`internal/models/models.go`)
- Expanded `DailyStat` struct to include new fields:
    - `AssignToBlock`
    - `Deletions`
    - `Backlog`
    - `Recovery`

### 2. Database Updates (`internal/db/queries.go`)
- Updated `GetDailyActivityStats` SQL query to aggregate new action types from `activity_log`:
    - `assign_to_block`
    - `delete`
    - `backlog`
    - `recovery`
- Updated the scanning logic to populate the new fields in `DailyStat`.

### 3. UI Template Updates (`internal/templates/activity.html`)
- Added new bar styles and legend items for the new activity types.
- Colors used:
    - Creation: Surface color (`--c-sf`)
    - Updates: Blue (`#60a5fa`)
    - Checkout: Accent color (`--c-ac`)
    - Assign: Amber (`#f59e0b`)
    - Delete: Red (`#ef4444`)
    - Backlog: Violet (`#8b5cf6`)
    - Recovery: Emerald (`#10b981`)
- Updated the summary row to show totals for all 7 activity types.

### 4. Handler Updates (`internal/handlers/ui.go`)
- Added activity logging to `submitBacklog` handler (action: `backlog`).
- Added activity logging to `WorkBlockDetail` handler for `assign_issue` action (action: `assign_to_block`).

### 5. Test Updates (`internal/db/db_test.go`)
- Updated `TestGetDailyActivityStats` to verify that all 7 activity types are correctly aggregated by date.
