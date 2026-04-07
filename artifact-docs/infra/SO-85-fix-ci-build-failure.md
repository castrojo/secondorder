# Tech Spec - SO-85: Fix CI Build Failure (Finalized)

## Issue
The CI build was failing with the following error:
`cmd/secondorder/main.go:29:12: pattern static: no matching files found`

## Root Cause
The `go:embed static` directive in `cmd/secondorder/main.go` required the `static` directory to be present and non-empty in the same directory as the source file. While the directory existed locally, it was being ignored by git due to a broad ignore pattern in the root `.gitignore` file.

The pattern `secondorder` in `.gitignore` was intended to ignore the compiled binary but was also matching the `cmd/secondorder/` directory and all its contents, effectively excluding `cmd/secondorder/static/` and several templates from the repository.

Additionally, some temporary Go files at the root (`check_now.go`, `test_sqlite.go`) were causing build errors due to duplicate `main` package declarations.

## Solution
1. **Fixed .gitignore**: Modified `.gitignore` to use `/secondorder` instead of `secondorder`. This ensures that only the binary in the root directory is ignored, while the `cmd/secondorder/` directory and its contents are correctly tracked.
2. **Consolidated Static Assets**:
    - Created a root `static` Go package (`static/static.go`) that embeds all static assets.
    - Updated `cmd/secondorder/main.go` to import and use the central `static.FS` instead of a local copy.
    - Removed the redundant `cmd/secondorder/static` directory to ensure a single source of truth and reduce duplication.
3. **Resolved Root Build Errors**:
    - Removed temporary root Go files (`check_now.go`, `test_sqlite.go`) which were causing package `main` redeclaration errors during `go build ./...`.
4. **Ensured Asset Tracking**: Tracked all relevant static assets (including `favicon-v2.svg`) to ensure they are available in the CI environment.

## Verification Results
- Ran `go build ./...` locally: Success (all packages build correctly).
- Ran `make build`: Success.
- Verified that `static.FS` is correctly served by the application.
- Verified tracking status with `git status --ignored`: No relevant files are ignored anymore.

## Impact
This final fix ensures that all static assets are correctly tracked and available to the build system via a centralized package, and that the project build is no longer blocked by root-level script conflicts or missing embedded patterns.
