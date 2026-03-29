$ErrorActionPreference = "Stop"

$SpecSlug = if ($args.Count -gt 0) { $args[0] } else { "" }
$PlanDir = ".draftspec/plans/$SpecSlug"
$SpecFile = ".draftspec/specs/$SpecSlug.md"
$TasksFile = "$PlanDir/tasks.md"
$errors = 0

if ([string]::IsNullOrWhiteSpace($SpecSlug)) {
  Write-Error "usage: check-implement-ready.ps1 <spec-slug>"
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
Check-File "$PlanDir/plan.md"
Check-File $TasksFile
Check-File "$PlanDir/data-model.md"
Check-File ".draftspec/templates/prompts/implement.md"

if (Test-Path -LiteralPath "$PlanDir/contracts" -PathType Container) {
  Write-Output "OK: $PlanDir/contracts"
} else {
  Write-Output "ERROR: missing $PlanDir/contracts"
  $script:errors++
}

if (Test-Path -LiteralPath $TasksFile -PathType Leaf) {
  Check-Contains $TasksFile '^## (Покрытие критериев приемки|Acceptance Coverage)$' "tasks include acceptance coverage section"
  Check-Contains $TasksFile 'T[0-9]+\.[0-9]+' "tasks include phase-scoped task IDs"
  Check-Contains $TasksFile 'AC-[0-9][0-9][0-9].*->.*T[0-9]+\.[0-9]+' "tasks include AC-to-task coverage lines"
}

$InspectScript = ".draftspec/scripts/inspect-spec.ps1"
if ((Test-Path -LiteralPath $InspectScript -PathType Leaf) -and (Test-Path -LiteralPath $SpecFile -PathType Leaf) -and (Test-Path -LiteralPath $TasksFile -PathType Leaf)) {
  & powershell -NoProfile -ExecutionPolicy Bypass -File $InspectScript $SpecFile $TasksFile
}

$CheckConstitution = ".draftspec/scripts/check-constitution.ps1"
if (Test-Path -LiteralPath $CheckConstitution -PathType Leaf) {
  & powershell -NoProfile -ExecutionPolicy Bypass -File $CheckConstitution ".draftspec/constitution.md"
}

if ($errors -ne 0) {
  exit 1
}
