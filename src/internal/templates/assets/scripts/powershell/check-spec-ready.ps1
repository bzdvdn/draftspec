$ErrorActionPreference = "Stop"
$errors = 0

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
Check-File ".draftspec/templates/spec.md"
Check-File ".draftspec/templates/prompts/spec.md"

if (Test-Path -LiteralPath ".draftspec/templates/spec.md" -PathType Leaf) {
  Check-Contains ".draftspec/templates/spec.md" '^## (Требования|Requirements)$' "spec template has requirements section"
  Check-Contains ".draftspec/templates/spec.md" 'RQ-[0-9][0-9][0-9]' "spec template includes requirement IDs"
  Check-Contains ".draftspec/templates/spec.md" 'AC-[0-9][0-9][0-9]' "spec template includes acceptance IDs"
  Check-Contains ".draftspec/templates/spec.md" 'Given' "spec template includes Given marker"
  Check-Contains ".draftspec/templates/spec.md" 'When' "spec template includes When marker"
  Check-Contains ".draftspec/templates/spec.md" 'Then' "spec template includes Then marker"
}

if ($errors -ne 0) {
  exit 1
}
