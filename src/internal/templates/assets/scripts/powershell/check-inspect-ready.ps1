$ErrorActionPreference = "Stop"

$Slug = if ($args.Count -gt 0) { $args[0] } else { "" }
if ([string]::IsNullOrWhiteSpace($Slug)) {
  Write-Error "usage: check-inspect-ready.ps1 <slug>"
  exit 1
}

$Root = ".draftspec"
$SpecFile = Join-Path $Root "specs/$Slug.md"
$ConstitutionFile = Join-Path $Root "constitution.md"
$missing = $false

foreach ($Path in @($ConstitutionFile, $SpecFile)) {
  if (-not (Test-Path -LiteralPath $Path)) {
    Write-Output "ERROR: missing required file: $Path"
    $missing = $true
  }
}

if ($missing) {
  exit 1
}

Write-Output "OK: inspect is ready for slug '$Slug'"
