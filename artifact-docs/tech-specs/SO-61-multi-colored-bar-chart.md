# Tech Spec: SO-61 Multi-colored Bar Chart

## Overview
Implement a multi-colored bar chart for the activity feed to show separate bars for updates, creation, and checkout events. This replaces the previous 2-bar chart (Created vs Completed).

## Changes

### 1. Models (`internal/models/models.go`)
- Updated `DailyStat` struct to include `Updates`, `Creations`, and `Checkouts` instead of `Created` and `Completed`.

### 2. Database (`internal/db/queries.go`)
- Modified `GetDailyActivityStats` to query `activity_log` table for granular action counts.
- Updated `LogActivity` to format `created_at` as a string (`2006-01-02 15:04:05`) to ensure SQLite `DATE()` function works correctly across different driver behaviors.

### 3. UI Template (`internal/templates/activity.html`)
- Updated the chart to render 3 bars per day.
- Added legend for Creation, Updates, and Checkout.
- Colors:
    - **Creation**: Hollow style (surface background with border).
    - **Updates**: Blue (#60a5fa).
    - **Checkout**: Amber (accent color).
- Updated summary stats at the bottom of the chart.

### 4. Tests (`internal/db/db_test.go`)
- Updated `TestGetDailyActivityStats` to reflect the new data structure and verify correct counting of the three action types.

## Verification
- Ran `go test -v ./internal/db/... -run TestGetDailyActivityStats` and confirmed it passes.
