$ErrorActionPreference = "Stop"

function Check-File([string]$Path) {
  if (Test-Path -LiteralPath $Path -PathType Leaf) {
    Write-Output "OK: $Path"
  } else {
    Write-Output "ERROR: missing $Path"
  }
}

Check-File ".draftspec/constitution.md"
Check-File ".draftspec/templates/spec.md"
Check-File ".draftspec/templates/prompts/spec.md"
