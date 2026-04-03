# Data Model: export-report

## ExportRequest

Value object passed from `ExportHandler` to `csv.Generate()`. Not persisted.

| Field     | Type       | Notes                                      |
|-----------|------------|--------------------------------------------|
| Filters   | FilterSet  | Parsed from query params via parseFilters()|
| Columns   | []Column   | Ordered list from server-side column registry (DEC-003) |

**Invariants**:
- `Columns` must be non-nil; an empty slice is valid and produces a header-only CSV (AC-003)
- `Filters` must be a valid `FilterSet`; invalid filters are rejected by `parseFilters()` before reaching this type

## Column

| Field   | Type   | Notes                              |
|---------|--------|------------------------------------|
| Key     | string | Internal field name                |
| Header  | string | Display label used as CSV header   |

**Justification**: AC-001 (column headers in first row), RQ-002 (column order matches visible table order)

## No New Persistent Entities

This feature adds no new database tables, migrations, or persisted state.
