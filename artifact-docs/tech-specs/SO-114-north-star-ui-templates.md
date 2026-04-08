# SO-114: North Star Metrics fields in UI templates

## Scope

Confirmed that the Work Block UI supports North Star metrics in the user-facing templates:

- `internal/templates/work_blocks.html`
  - Propose form includes `north_star_metric` and `north_star_target`.
- `internal/templates/work_block_detail.html`
  - Detail header displays the North Star metric and target when present.
  - Edit modal includes `north_star_metric` and `north_star_target`.

## Verification

Added UI regression tests in `internal/handlers/work_blocks_ui_test.go` covering:

- GET `/work-blocks` renders both North Star inputs in the propose form.
- GET `/work-blocks/{id}` renders:
  - the North Star display row in the header
  - both North Star inputs in the edit modal
- GET `/work-blocks/{id}` omits the North Star summary row when no metric is set.

This closes the QA gap from `artifact-docs/qa-so-109.md` by asserting the template output directly.

Verification run:

- `go test ./internal/handlers -run 'Test(ListWorkBlocksRendersNorthStarInputs|WorkBlockDetailRendersNorthStarDisplayAndEditInputs|WorkBlockDetailHidesNorthStarDisplayWhenMetricEmpty)$'`
- `go test ./...`
