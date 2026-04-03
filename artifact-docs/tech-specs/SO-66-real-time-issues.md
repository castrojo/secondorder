# Tech Spec: Real-time updates for Issues Tab (SO-66)

## Overview
This feature implements real-time updates for the Issues Tab in the Second Order UI. It ensures that whenever an issue is created, updated, or deleted (via API or UI), the changes are instantly reflected in the issues list without requiring a manual page refresh.

## Architecture

### Backend: SSE Broadcasts
The existing Server-Sent Events (SSE) infrastructure was leveraged to broadcast issue events. The following events were added:
- `issue_created`: Broadcasted when a new issue is created. Payload includes the issue key, title, and status.
- `issue_updated`: Broadcasted when an issue is updated (status change, assignment, etc.). Payload includes the issue key, title, and status.
- `issue_deleted`: Broadcasted when an issue is deleted. Payload includes the issue key.

Broadcasts are integrated into both the REST API handlers (`internal/handlers/api.go`) and the UI post handlers (`internal/handlers/ui.go`).

### Frontend: HTMX Partial Updates
To efficiently update the UI without refreshing the entire page, HTMX partial rendering was used:
1.  **Partial Extraction**: The issue list in `internal/templates/issues.html` was extracted into a named template block `issue_list`.
2.  **HTMX Integration**: The issue list container now uses `hx-get` to fetch the updated list and `hx-trigger="issue-changed from:body"` to listen for changes.
3.  **SSE to HTMX Bridge**: The SSE client in `internal/templates/partials.html` listens for the new issue events. When an event is received:
    - A toast notification is shown to the user.
    - A custom DOM event `issue-changed` is dispatched on the document body, which triggers HTMX to refresh the issue list.

## Changes

### `internal/handlers/api.go`
- Added SSE broadcasts to `CreateIssue`, `UpdateIssue`, `DeleteIssue`, `CheckoutIssue`, `AssignIssueToBlock`, and `UnassignIssueFromBlock`.

### `internal/handlers/ui.go`
- Added SSE broadcasts to `createIssueUI`, `updateIssueUI`, `AgentAssign`, and `updateWorkBlockUI` (assign/unassign actions).
- Updated `ListIssues` to support HTMX partial rendering by returning only the `issue_list` block when the `HX-Request` header is present.

### `internal/templates/issues.html`
- Extracted the issue list into `{{define "issue_list"}}`.
- Added HTMX attributes to the list container for real-time updates.

### `internal/templates/partials.html`
- Added SSE event listeners for `issue_created`, `issue_updated`, and `issue_deleted`.
- Implemented `issue-changed` event dispatching.
- Added toast notifications for real-time issue updates.

### `internal/templates/activity.html`
- Fixed a pre-existing template syntax error (mismatched `{{if}}`/`{{else}}`) that was preventing template parsing during testing.

## Verification Results
- Automated tests in `internal/handlers/handlers_test.go` and a temporary `sse_broadcast_test.go` verified that SSE events are broadcasted correctly from all relevant handlers.
- All 35+ handler tests pass successfully.
