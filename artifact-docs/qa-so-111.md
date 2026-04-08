# QA Report: Acceptance Criteria Completeness Validator Validation (SO-111)

## Summary
The acceptance-criteria completeness validator is only partially complete.

The core validator logic works for `api`, `backend`, and generic issue types, and the API returns non-blocking `warnings` on both create and update flows. The blocking gap is the human UI warning path required by the PRD: warning messages are generated during create and edit, but the redirect/query-param plumbing is inconsistent, so the warning banner does not reliably render for humans.

## Verified
- [x] `type` is persisted on issues and exposed through the issue model flow.
- [x] `ValidateAC` enforces type-specific warnings for `api` and `backend`, plus generic list validation for other issue types.
- [x] `POST /api/v1/issues` includes `warnings` in the response.
- [x] `PATCH /api/v1/issues/{key}` includes `warnings` in the response.
- [x] Warnings are non-blocking; issue creation and updates still succeed.

## Findings
- [ ] Human create flow warning is not rendered.
  Evidence:
  `createIssueUI` appends `&warning=...` to the redirect URL, but `ListIssues` does not pass `warning` into the template data.
  References:
  [ui.go](/Users/alexander/side-projects/2026/secondorder/internal/handlers/ui.go#L168), [ui.go](/Users/alexander/side-projects/2026/secondorder/internal/handlers/ui.go#L196), [ui.go](/Users/alexander/side-projects/2026/secondorder/internal/handlers/ui.go#L101), [issues.html](/Users/alexander/side-projects/2026/secondorder/internal/templates/issues.html#L26)

- [ ] Human edit flow warning is not rendered.
  Evidence:
  `updateIssueUI` redirects with `flash=warning&msg=...`, but `IssueDetail` only reads `warning` and the template only renders `.Warning`.
  References:
  [ui.go](/Users/alexander/side-projects/2026/secondorder/internal/handlers/ui.go#L391), [ui.go](/Users/alexander/side-projects/2026/secondorder/internal/handlers/ui.go#L396), [ui.go](/Users/alexander/side-projects/2026/secondorder/internal/handlers/ui.go#L271), [issue_detail.html](/Users/alexander/side-projects/2026/secondorder/internal/templates/issue_detail.html#L24)

## Code Evidence
- Validator rules:
  [ac_validator.go](/Users/alexander/side-projects/2026/secondorder/internal/validator/ac_validator.go#L9)
- API warning response on update:
  [api.go](/Users/alexander/side-projects/2026/secondorder/internal/handlers/api.go#L213)
- API warning response on create:
  [api.go](/Users/alexander/side-projects/2026/secondorder/internal/handlers/api.go#L540)
- Existing tests cover basic UI creation success and backend/API flows, but not warning-banner rendering:
  [handlers_test.go](/Users/alexander/side-projects/2026/secondorder/internal/handlers/handlers_test.go#L367)

## Verification
- `GOPROXY=off GOCACHE=$(pwd)/.gocache GOMODCACHE=/Users/alexander/gocode/pkg/mod go test ./...`

Result: passed

## Verdict
Fail for acceptance.

Reason:
The PRD requires that the human dashboard show a non-blocking warning when an issue is saved with incomplete acceptance criteria. The backend and API portions satisfy that requirement, but the current UI implementation does not complete the warning-display loop for humans.

## Notes
- I could not use the issue REST API over HTTP from this sandbox because binding `localhost:9003` is blocked here. As with prior QA artifacts, board bookkeeping must be applied against the local `so.db` state instead of the REST surface.
