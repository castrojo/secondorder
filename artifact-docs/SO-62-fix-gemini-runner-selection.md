# Fixed: Restore Gemini runner and model selection

Issue: SO-62

## Changes

### Frontend
- **internal/templates/partials.html**: Updated `updateModels` JavaScript function to support the `gemini` runner. Added the following models:
  - Gemini 2.0 Flash
  - Gemini 2.0 Flash Lite
  - Gemini 1.5 Pro
  - Gemini 1.5 Flash
- **internal/templates/agent_detail.html**: Added "Gemini" to the Runner selection dropdown in the agent edit form.
- **internal/templates/agents.html**: Added "Gemini" to the Runner selection dropdown in the new agent creation form and cleared hardcoded model options to rely on the `updateModels` initialization.

### Backend
- **internal/models/models.go**: Added `RunnerGemini` constant and updated `RunnerModels` map to include valid models for Gemini.
- **internal/handlers/ui.go**: Improved `createAgentUI` to correctly handle default models for different runners if no model is explicitly provided.

## Verification
- Ran existing tests for models and handlers: `go test ./internal/models/... ./internal/handlers/...`. All tests passed.
- Verified that `IsValidModelForRunner` now correctly validates Gemini models.
