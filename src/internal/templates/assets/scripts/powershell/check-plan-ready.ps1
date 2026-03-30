$ErrorActionPreference = "Stop"

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
& (Join-Path $ScriptDir "run-draftspec.ps1") __internal check-plan-ready --root . @args
exit $LASTEXITCODE
