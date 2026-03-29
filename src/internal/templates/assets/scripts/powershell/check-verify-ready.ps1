$ErrorActionPreference = "Stop"

$SpecSlug = if ($args.Count -gt 0) { $args[0] } else { "" }
$PlanDir = ".draftspec/plans/$SpecSlug"
$SpecFile = ".draftspec/specs/$SpecSlug.md"
$TasksFile = "$PlanDir/tasks.md"
$errors = 0

if ([string]::IsNullOrWhiteSpace($SpecSlug)) {
  Write-Error "usage: check-verify-ready.ps1 <spec-slug>"
  exit 1
}

function Check-File([string]$Path) {
  if (Test-Path -LiteralPath $Path -PathType Leaf) {
    Write-Output "OK: $Path"
  } else {
    Write-Output "ERROR: missing $Path"
    $script:errors++
  }
}

Check-File ".draftspec/constitution.md"
Check-File $SpecFile
Check-File $TasksFile
Check-File ".draftspec/templates/verify-report.md"
Check-File ".draftspec/templates/prompts/verify.md"

$InspectScript = ".draftspec/scripts/inspect-spec.ps1"
if ((Test-Path -LiteralPath $InspectScript -PathType Leaf) -and (Test-Path -LiteralPath $SpecFile -PathType Leaf) -and (Test-Path -LiteralPath $TasksFile -PathType Leaf)) {
  & powershell -NoProfile -ExecutionPolicy Bypass -File $InspectScript $SpecFile $TasksFile
}

$VerifyTaskState = ".draftspec/scripts/verify-task-state.ps1"
if ((Test-Path -LiteralPath $VerifyTaskState -PathType Leaf) -and (Test-Path -LiteralPath $TasksFile -PathType Leaf)) {
  & powershell -NoProfile -ExecutionPolicy Bypass -File $VerifyTaskState $SpecSlug
}

if ($errors -ne 0) {
  exit 1
}
