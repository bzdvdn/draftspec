$ErrorActionPreference = "Stop"

$SpecName = if ($args.Count -gt 0) { $args[0] } else { "" }
$SpecsDir = if ($args.Count -gt 1) { $args[1] } else { ".draftspec/specs" }

if ([string]::IsNullOrWhiteSpace($SpecName)) {
  Write-Error "usage: show-spec.ps1 <spec-name> [specs-dir]"
  exit 1
}

$SpecFile = Join-Path $SpecsDir "$SpecName.md"
if (-not (Test-Path -LiteralPath $SpecFile -PathType Leaf)) {
  Write-Error "spec not found: $SpecFile"
  exit 1
}

Get-Content -LiteralPath $SpecFile -Raw
