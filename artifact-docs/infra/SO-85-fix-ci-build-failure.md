# Tech Spec - SO-85: Fix CI Build Failure

## Issue
The CI build was failing with the following error:
`cmd/secondorder/main.go:29:12: pattern static: no matching files found`

## Root Cause
The `go:embed static` directive in `cmd/secondorder/main.go` requires the `static` directory to be present and non-empty in the same directory as the source file. While the directory existed locally, it was being ignored by git due to a broad ignore pattern in the root `.gitignore` file.

The pattern `secondorder` in `.gitignore` was intended to ignore the compiled binary but was also matching the `cmd/secondorder/` directory and all its contents, effectively excluding `cmd/secondorder/static/` and several templates from the repository.

## Solution
1. Modified `.gitignore` to use `/secondorder` instead of `secondorder`. This ensures that only the binary in the root directory is ignored, while the `cmd/secondorder/` directory and its contents are correctly tracked.
2. Verified that `cmd/secondorder/static/` and other previously ignored files are now untracked (and thus can be added to the repository) rather than ignored.

## Verification Results
- Ran `go build ./...` locally: Success.
- Ran `make build`: Success.
- Verified tracking status with `git status --ignored`: `cmd/secondorder/static/` is now listed as an untracked directory instead of an ignored one.

## Impact
This change allows the `static` assets and organization templates required for the build to be committed to the repository, which will resolve the "no matching files found" error in the CI environment.
