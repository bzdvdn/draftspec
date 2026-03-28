$ErrorActionPreference = "Stop"

$SpecSlug = if ($args.Count -gt 0) { $args[0] } else { "" }
$PlanDir = ".draftspec/plans/$SpecSlug"

if ([string]::IsNullOrWhiteSpace($SpecSlug)) {
  Write-Error "usage: check-implement-ready.ps1 <spec-slug>"
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
Check-File ".draftspec/specs/$SpecSlug.md"
Check-File "$PlanDir/plan.md"
Check-File "$PlanDir/tasks.md"
Check-File "$PlanDir/data-model.md"
Check-File ".draftspec/templates/prompts/implement.md"

if (Test-Path -LiteralPath "$PlanDir/contracts" -PathType Container) {
  Write-Output "OK: $PlanDir/contracts"
} else {
  Write-Output "ERROR: missing $PlanDir/contracts"
}

$CheckConstitution = ".draftspec/scripts/check-constitution.ps1"
if (Test-Path -LiteralPath $CheckConstitution -PathType Leaf) {
  & powershell -NoProfile -ExecutionPolicy Bypass -File $CheckConstitution ".draftspec/constitution.md"
}
