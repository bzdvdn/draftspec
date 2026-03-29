$ErrorActionPreference = "Stop"

$Slug = if ($args.Count -gt 0) { $args[0] } else { "" }
if ([string]::IsNullOrWhiteSpace($Slug)) {
  Write-Error "usage: check-inspect-ready.ps1 <slug>"
  exit 1
}

$Root = ".draftspec"
$SpecFile = Join-Path $Root "specs/$Slug.md"
$ConstitutionFile = Join-Path $Root "constitution.md"
$InspectReportTemplate = Join-Path $Root "templates/inspect-report.md"
$InspectPrompt = Join-Path $Root "templates/prompts/inspect.md"
$InspectScript = Join-Path $Root "scripts/inspect-spec.ps1"
$missing = $false

foreach ($Path in @($ConstitutionFile, $SpecFile, $InspectReportTemplate, $InspectPrompt)) {
  if (-not (Test-Path -LiteralPath $Path)) {
    Write-Output "ERROR: missing required file: $Path"
    $missing = $true
  } else {
    Write-Output "OK: $Path"
  }
}

if ($missing) {
  exit 1
}

if (Test-Path -LiteralPath $InspectScript -PathType Leaf) {
  & powershell -NoProfile -ExecutionPolicy Bypass -File $InspectScript $SpecFile
}

Write-Output "OK: inspect is ready for slug '$Slug'"
