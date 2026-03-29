$ErrorActionPreference = "Stop"

$SpecSlug = if ($args.Count -gt 0) { $args[0] } else { "" }
$PlanDir = ".draftspec/plans/$SpecSlug"
$SpecFile = ".draftspec/specs/$SpecSlug.md"
$PlanFile = "$PlanDir/plan.md"
$errors = 0

if ([string]::IsNullOrWhiteSpace($SpecSlug)) {
  Write-Error "usage: check-tasks-ready.ps1 <spec-slug>"
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

function Check-Contains([string]$Path, [string]$Pattern, [string]$Label) {
  if (Select-String -Path $Path -Pattern $Pattern -Quiet) {
    Write-Output "OK: $Label"
  } else {
    Write-Output "ERROR: $Label"
    $script:errors++
  }
}

Check-File ".draftspec/constitution.md"
Check-File $SpecFile
Check-File $PlanFile
Check-File "$PlanDir/data-model.md"
Check-File ".draftspec/templates/tasks.md"
Check-File ".draftspec/templates/prompts/tasks.md"

if (Test-Path -LiteralPath "$PlanDir/contracts" -PathType Container) {
  Write-Output "OK: $PlanDir/contracts"
} else {
  Write-Output "ERROR: missing $PlanDir/contracts"
  $script:errors++
}

if (Test-Path -LiteralPath $SpecFile -PathType Leaf) {
  Check-Contains $SpecFile 'AC-[0-9][0-9][0-9]' "spec has stable acceptance IDs"
}

if (Test-Path -LiteralPath $PlanFile -PathType Leaf) {
  Check-Contains $PlanFile 'DEC-[0-9][0-9][0-9]' "plan has stable decision IDs"
}

$InspectScript = ".draftspec/scripts/inspect-spec.ps1"
if ((Test-Path -LiteralPath $InspectScript -PathType Leaf) -and (Test-Path -LiteralPath $SpecFile -PathType Leaf)) {
  & powershell -NoProfile -ExecutionPolicy Bypass -File $InspectScript $SpecFile
}

if ($errors -ne 0) {
  exit 1
}
