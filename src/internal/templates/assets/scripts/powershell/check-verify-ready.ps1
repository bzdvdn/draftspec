$ErrorActionPreference = "Stop"

$SpecSlug = if ($args.Count -gt 0) { $args[0] } else { "" }
$PlanDir = ".draftspec/plans/$SpecSlug"

if ([string]::IsNullOrWhiteSpace($SpecSlug)) {
  Write-Error "usage: check-verify-ready.ps1 <spec-slug>"
  exit 1
}

function Check-File([string]$Path) {
  if (Test-Path -LiteralPath $Path -PathType Leaf) {
    Write-Output "OK: $Path"
  } else {
    Write-Output "ERROR: missing $Path"
  }
}

Check-File ".draftspec/constitution.md"
Check-File "$PlanDir/tasks.md"
Check-File ".draftspec/templates/verify-report.md"
Check-File ".draftspec/templates/prompts/verify.md"
