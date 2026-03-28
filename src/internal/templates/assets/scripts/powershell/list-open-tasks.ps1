$ErrorActionPreference = "Stop"

$SpecSlug = if ($args.Count -gt 0) { $args[0] } else { "" }
$TasksFile = ".draftspec/plans/$SpecSlug/tasks.md"

if ([string]::IsNullOrWhiteSpace($SpecSlug)) {
  Write-Error "usage: list-open-tasks.ps1 <spec-slug>"
  exit 1
}

if (-not (Test-Path -LiteralPath $TasksFile)) {
  Write-Error "tasks file not found: $TasksFile"
  exit 1
}

Get-Content -LiteralPath $TasksFile | Where-Object { $_ -match '^- \[ \]' }
