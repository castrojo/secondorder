# QA Report: Issue Detail Empty States Validation (SO-116)

## Summary
The implementation of SO-76 (Empty states for sub-issues and runs) is complete and verified. 
Both the "Sub-issues" and "Runs" sections now render consistently in the Issue Detail view even when empty, providing clear feedback and an actionable path for adding sub-issues.

## Verified
- [x] **Empty Sections Rendering**: Both "Sub-issues" and "Runs" sections are visible on an issue with no child issues and no runs.
  - Verified by `TestIssueDetail_EmptySections` in `internal/handlers/handlers_test.go`.
- [x] **Add sub-issue CTA**: The "Add sub-issue" button is visible in the empty state and toggles the inline creation form.
  - Verified by inspection of `internal/templates/issue_detail.html` and `TestIssueDetail_EmptySections`.
- [x] **CTA Usability**: Creating a sub-issue via the inline form correctly associates it with the parent issue.
  - Verified by a new test case `TestCreateSubIssueFromUI_DetailForm` in `internal/handlers/handlers_test.go`.
- [x] **Non-empty State Rendering**: Existing sub-issues and runs still render correctly when present.
  - Verified by `TestIssueDetail_WithSections` in `internal/handlers/handlers_test.go`.
- [x] **Styling Consistency**: The empty states use dashed borders and appropriate typography matching the design spec (`bg-sf/30`, `border-dashed`, etc.).

## Findings
- No significant gaps found. The implementation closely follows the provided design spec and mockup.
- The inline form in `issue_detail.html` includes an `autofocus` attribute on the title field, which is a nice touch for usability when the "Add" button is clicked.

## Code Evidence
- Template implementation:
  [issue_detail.html](/Users/alexander/side-projects/2026/secondorder/internal/templates/issue_detail.html#L149)
- Handler logic for `parent_issue_key`:
  [ui.go](/Users/alexander/side-projects/2026/secondorder/internal/handlers/ui.go#L137)
- New verification test:
  [handlers_test.go](/Users/alexander/side-projects/2026/secondorder/internal/handlers/handlers_test.go#L1508)

## Verification
- `go test ./...`
- `bash artifact-docs/gates.sh`

Result: **All gates and tests passed.**

## Verdict
**Pass for acceptance.**
The changes satisfy all acceptance criteria defined in SO-116 and SO-76.
