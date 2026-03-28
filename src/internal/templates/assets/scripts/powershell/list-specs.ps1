$ErrorActionPreference = "Stop"

$SpecsDir = if ($args.Count -gt 0) { $args[0] } else { ".draftspec/specs" }

if (-not (Test-Path -LiteralPath $SpecsDir -PathType Container)) {
  Write-Error "specs directory not found: $SpecsDir"
  exit 1
}

Get-ChildItem -LiteralPath $SpecsDir -File -Filter "*.md" | Sort-Object Name | ForEach-Object {
  $_.BaseName
}
