# Documentation Housekeeping (SO-112)

Date: 2026-04-08

## Completed cleanup

- Deleted `artifact-docs/SO-62-fix-gemini-runner-selection.md` because SO-62 was cancelled and the doc was left at the root level.
- Moved `artifact-docs/SO-71-always-show-assign-form.md` to `artifact-docs/product/SO-71-always-show-assign-form.md`.
- Moved `artifact-docs/qa-ux-accessibility-SO-13.md` to `artifact-docs/product/qa-ux-accessibility-SO-13.md`.
- Moved `artifact-docs/security-audit-SO-6.md` to `artifact-docs/infra/security-audit-SO-6.md`.

## Verification notes

- `artifact-docs/tech-specs/SO-47-restore-gemini-full-support.md` should be kept. Gemini support is still present in the product and later infra updates in `artifact-docs/infra/SO-92-model-updates.md` and `artifact-docs/decisions/model-migration-to-gemini-3.md` confirm the runner remained active after SO-47 shipped.
- `artifact-docs/design/brand-system.md` is still stale. Audit `05f3fb5b` already flagged it as outdated after the Pico CSS migration, and `artifact-docs/decisions/api-key-rotation.md` records that SO-113 was created on 2026-04-07 to assign the refresh to Designer.

## Operational note

- The local SO API endpoint at `http://localhost:9003` was unreachable from this environment during SO-112, so checkout and ticket status updates could not be completed via API.
