# Plan: export-report

## Goal

Add a `GET /reports/export` endpoint that streams a CSV of the current filtered table view to the browser. Reuse the existing filter parsing from the table handler. Wire a frontend button that passes the current filter state as query parameters.

## Implementation Surfaces

- **`handlers/reports.go`** (existing) — add `ExportHandler` alongside the existing table handler; reuse `parseFilters()` already present in this file
- **`services/csv.go`** (new) — CSV generation service: accepts rows + column order, returns `io.Reader`
- **`middleware/auth.go`** (existing) — apply `RequireAuth` middleware already used by other report endpoints
- **`ui/components/TableToolbar`** (existing) — add "Export as CSV" button; on click, build export URL from current filter state and trigger download

## Acceptance Approach

- AC-001: `ExportHandler` calls `parseFilters()`, queries rows, passes result to `csv.Generate()`; integration test asserts row count and headers match filtered table response
- AC-002: `ExportHandler` sets `Content-Disposition: attachment; filename="report-YYYY-MM-DD.csv"` using `time.Now()` on the server side
- AC-003: `csv.Generate()` with empty row slice must produce a header-only output; unit test covers this path explicitly
- RQ-003: `RequireAuth` middleware already returns 401 for missing/invalid session; add route-level test asserting 401 on unauthenticated export request

## Decisions

### DEC-001 — Synchronous export, no streaming
**Why**: Dataset is capped at 50k rows (out of scope per spec). Synchronous generation keeps the handler simple and avoids background job infrastructure.
**Tradeoff**: Response latency scales with row count; acceptable for the stated scope.
**Affects**: AC-001, AC-003
**Validation**: Load test with 50k rows must complete within 10s on CI hardware.

### DEC-002 — Reuse parseFilters() from table handler
**Why**: Filter logic is already tested and handles all edge cases (missing params, invalid values). Duplicating it would create a maintenance gap.
**Affects**: AC-001, RQ-001
**Validation**: Existing filter unit tests cover the shared code path.

### DEC-003 — Column order driven by server-side column registry
**Why**: Visible column order is managed server-side; the client must not dictate column order to avoid inconsistency between table view and export.
**Affects**: RQ-002
**Validation**: Integration test compares CSV header order against the column registry order for the same table view.

## Data and Contracts

No new persistent state. See `data-model.md` for the `ExportRequest` value object.
Export endpoint contract: `GET /reports/export?<filter_params>` → `text/csv` with `Content-Disposition` header.

## Sequencing Notes

Phase 1 (backend) must complete before Phase 2 (frontend) can be tested end-to-end. Within Phase 1, T1.1 and T1.2 are the prerequisite for T1.3–T1.5.

## Risks

- `parseFilters()` has one known edge case with date range params (see inline TODO in `handlers/reports.go`). If that edge case affects export, T1.2 may take longer.
  - Mitigation: inspect the TODO before T1.2, escalate if the fix is non-trivial.

## Rollout and Compatibility

No migration needed. No feature flag required — the endpoint is new. The UI button is additive.

## Validation

- Unit: `csv.Generate()` with empty rows → header-only output (AC-003)
- Unit: `csv.Generate()` column order matches registry order (DEC-003)
- Integration: authenticated export with active filter → correct row count (AC-001)
- Integration: unauthenticated request → 401 (RQ-003)
- Manual: filename pattern `report-YYYY-MM-DD.csv` in browser download dialog (AC-002)
