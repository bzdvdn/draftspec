# Tasks: export-report

## Phase 1 — Backend

- [x] T1.1 Add `GET /reports/export` route and `ExportHandler` stub in `handlers/reports.go`
- [x] T1.2 Wire `parseFilters()` into `ExportHandler`; add integration test asserting filtered row count matches table response
- [ ] T1.3 Implement `services/csv.go`: `Generate(rows []Row, columns []Column) io.Reader`; unit test column order and empty-rows case (AC-001, AC-002, AC-003, DEC-003)
- [ ] T1.4 Set `Content-Disposition: attachment; filename="report-YYYY-MM-DD.csv"` header in `ExportHandler` (AC-002)
- [ ] T1.5 Apply `RequireAuth` middleware to export route; add unauthenticated request test asserting 401 (RQ-003)

## Phase 2 — Frontend

- [ ] T2.1 Add "Export as CSV" button to `TableToolbar`; on click build export URL from current filter state and trigger browser download (AC-001)
