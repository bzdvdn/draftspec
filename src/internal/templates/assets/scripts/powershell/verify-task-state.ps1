$ErrorActionPreference = "Stop"

$SpecSlug = if ($args.Count -gt 0) { $args[0] } else { "" }
$TasksFile = ".draftspec/plans/$SpecSlug/tasks.md"

if ([string]::IsNullOrWhiteSpace($SpecSlug)) {
  Write-Error "usage: verify-task-state.ps1 <spec-slug>"
  exit 1
}

if (-not (Test-Path -LiteralPath $TasksFile -PathType Leaf)) {
  Write-Output "ERROR: missing $TasksFile"
  exit 1
}

$lines = Get-Content -LiteralPath $TasksFile
$total = ($lines | Where-Object { $_ -match '^- \[[ x]\]' }).Count
$completed = ($lines | Where-Object { $_ -match '^- \[x\]' }).Count
$open = ($lines | Where-Object { $_ -match '^- \[ \]' }).Count

if ($total -eq 0) {
  Write-Output "ERROR: no task checkboxes found in $TasksFile"
  exit 1
}

Write-Output "TASKS_TOTAL=$total"
Write-Output "TASKS_COMPLETED=$completed"
Write-Output "TASKS_OPEN=$open"

if ($open -gt 0) {
  Write-Output "WARN: open tasks remain in $TasksFile"
  exit 0
}

Write-Output "OK: all tasks are marked complete in $TasksFile"
