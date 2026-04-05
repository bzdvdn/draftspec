$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
& "$ScriptDir/run-draftspec.ps1" trace $args
