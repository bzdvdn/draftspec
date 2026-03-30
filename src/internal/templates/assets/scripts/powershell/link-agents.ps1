$ErrorActionPreference = "Stop"

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
& (Join-Path $ScriptDir "run-draftspec.ps1") __internal link-agents --root . @args
exit $LASTEXITCODE
