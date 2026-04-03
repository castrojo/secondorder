# SO-44: Implementation: App version tracking

## Overview
Move the application version tracking from a plain text `VERSION` file to a Golang constant in `internal/models/version.go`.

## Changes

### Backend
- Created `internal/models/version.go` containing the `Version` constant.
- Updated `internal/handlers/ui.go` to use `models.Version` instead of reading from a file.
- Removed the `VERSION` file from the root directory.

### Frontend
- No changes required to `internal/templates/settings.html` as it already uses the `.Version` field passed from the handler.

### Testing
- Updated `internal/handlers/handlers_test.go` to verify the version is correctly displayed using the new constant.
- Added a test in `internal/models/models_test.go` to ensure the `Version` constant is defined and follows the expected format.

## Verification Results
- All tests in `internal/handlers` passed.
- All tests in `internal/models` passed.
