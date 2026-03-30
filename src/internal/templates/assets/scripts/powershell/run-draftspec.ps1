$ErrorActionPreference = "Stop"

$DraftspecBin = $env:DRAFTSPEC_BIN
if (-not [string]::IsNullOrWhiteSpace($DraftspecBin)) {
  $configuredCommand = Get-Command -Name $DraftspecBin -ErrorAction SilentlyContinue
  if ($null -ne $configuredCommand) {
    & $DraftspecBin @args
    exit $LASTEXITCODE
  }
  if (Test-Path -LiteralPath $DraftspecBin -PathType Leaf) {
    & $DraftspecBin @args
    exit $LASTEXITCODE
  }
  Write-Error "DRAFTSPEC_BIN is set but could not be resolved: $DraftspecBin. Set DRAFTSPEC_BIN to an executable path or command name, or add draftspec to PATH."
  exit 1
}

$defaultCommand = Get-Command -Name "draftspec" -ErrorAction SilentlyContinue
if ($null -ne $defaultCommand) {
  & draftspec @args
  exit $LASTEXITCODE
}

Write-Error "draftspec CLI not found. Set DRAFTSPEC_BIN to an executable path or add draftspec to PATH."
exit 1
