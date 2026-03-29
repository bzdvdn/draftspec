$ErrorActionPreference = "Stop"

$SpecSlug = if ($args.Count -gt 0) { $args[0] } else { "" }
$errors = 0
if ([string]::IsNullOrWhiteSpace($SpecSlug)) {
  Write-Error "usage: check-plan-ready.ps1 <spec-slug>"
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
Check-File ".draftspec/specs/$SpecSlug.md"
Check-File ".draftspec/templates/plan.md"
Check-File ".draftspec/templates/prompts/plan.md"

$SpecFile = ".draftspec/specs/$SpecSlug.md"
if (Test-Path -LiteralPath $SpecFile -PathType Leaf) {
  Check-Contains $SpecFile '^## (Критерии приемки|Acceptance Criteria)$' "spec has acceptance criteria section"
  Check-Contains $SpecFile 'AC-[0-9][0-9][0-9]' "spec has stable acceptance IDs"
}

$InspectScript = ".draftspec/scripts/inspect-spec.ps1"
if ((Test-Path -LiteralPath $InspectScript -PathType Leaf) -and (Test-Path -LiteralPath $SpecFile -PathType Leaf)) {
  & powershell -NoProfile -ExecutionPolicy Bypass -File $InspectScript $SpecFile
}

if ($errors -ne 0) {
  exit 1
}
