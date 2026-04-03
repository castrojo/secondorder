# Tech Spec: Full Support for Gemini Runner (SO-47)

## Overview
This specification details the changes made to fully restore and support the Gemini runner in SecondOrder. While Gemini was partially present in the UI, several key areas were missing support, including token usage display, comprehensive validation tests, and correct handling during initial organization bootstrapping.

## Changes

### 1. Scheduler Enhancements
- **internal/scheduler/scheduler.go**: Updated the token usage parsing logic to allow the `gemini` runner to retain its parsed tokens. Previously, tokens were cleared for any runner other than `claude_code`.

### 2. UI Improvements
- **internal/templates/agent_detail.html**: Updated usage statistics display (Tokens, Cost) to show values for `gemini` and `codex` runners instead of displaying "N/A".
- **internal/templates/run_detail.html**: Updated run-level statistics display to include `gemini` and `codex` runners.

### 3. Backend Robustness
- **cmd/secondorder/main.go**: Improved `applyStartupTemplate` to ensure that when a custom runner is specified via CLI flags (e.g., `-m gemini`), agents are created with a model compatible with that runner. It now falls back to the first available model for the runner if the template's model is incompatible.

### 4. Verification & Testing
- **internal/handlers/handlers_test.go**: Added new test cases to `TestAgentUI_Validation` to verify both creation and update validation for the `gemini` runner.
- Removed a stray `check_date.go` file from the root directory which was causing `go test ./...` to fail due to a missing dependency not used by the main project.

## Verification Results
- All tests passed: `go test ./...`
- Manual verification of the agent configuration UI confirms that Gemini is present and correctly filters models via the `updateModels` JavaScript function.
