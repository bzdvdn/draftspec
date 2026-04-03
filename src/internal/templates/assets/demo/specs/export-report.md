---
slug: export-report
phase: implement
---

# Export Report

## Goal

Allow authenticated users to download the current filtered table view as a CSV file so they can analyse the data in external tools without manual copy-paste.

## Why Now

Support team spends ~2h/week manually exporting data for customers. Self-service export removes a recurring escalation path before the Q2 customer review.

## Primary User Flow

1. User opens a table view with filters applied.
2. User clicks "Export as CSV".
3. System generates a CSV matching the current filtered view.
4. Browser downloads the file named `report-<date>.csv`.
5. User opens the file in a spreadsheet tool and sees all visible columns and rows.

## Acceptance Criteria

### AC-001 — Filtered export
**Given** an authenticated user viewing a filtered table
**When** they click "Export as CSV"
**Then** a CSV file is downloaded containing all rows matching the active filter, with column headers in the first row

### AC-002 — File naming
**Given** a successful export
**When** the download begins
**Then** the filename follows the pattern `report-YYYY-MM-DD.csv` using the server-side date

### AC-003 — Empty result
**Given** an authenticated user viewing a table with zero matching rows
**When** they click "Export as CSV"
**Then** a CSV file is downloaded containing only the header row and no data rows

## Requirements

- RQ-001: Export must reflect active filters, not the full unfiltered dataset
- RQ-002: Column order in CSV must match the visible column order in the table
- RQ-003: The export endpoint must be authenticated; unauthenticated requests return 401

## Scope

In scope:
- CSV download for the current filtered table view
- Authentication gate on the export endpoint
- Empty-result handling (header-only CSV)

Out of scope:
- PDF or Excel export formats
- Scheduled or email-delivered exports
- Export of hidden or system columns
- Streaming for datasets larger than 50k rows
- Export history or audit log

## Non-Goals

- Custom column selection UI
- Export progress indication for large datasets

## Context

Table views are server-rendered with filter state held in query parameters. The export endpoint can reuse the existing filter parsing logic already present in the table handler.

## Edge Cases

- Filter returns 0 rows → header-only CSV (covered by AC-003)
- Session expires mid-export → 401, client shows standard auth error

## Open Questions

none
