$ErrorActionPreference = "Stop"

$Slug = if ($args.Count -gt 0) { $args[0] } else { "" }
$Status = if ($args.Count -gt 1) { $args[1] } else { "" }
$Reason = if ($args.Count -gt 2) { $args[2] } else { "" }

if ([string]::IsNullOrWhiteSpace($Slug)) {
  Write-Error "usage: check-archive-ready.ps1 <slug> <status> <reason>"
  exit 1
}
if ([string]::IsNullOrWhiteSpace($Status)) {
  Write-Output "ERROR: archive status is required"
  exit 1
}
if ([string]::IsNullOrWhiteSpace($Reason)) {
  Write-Output "ERROR: archive reason is required"
  exit 1
}

if ($Status -notin @("completed", "superseded", "abandoned", "rejected", "deferred")) {
  Write-Output "ERROR: invalid archive status: $Status"
  exit 1
}

$Root = ".draftspec"
$SpecFile = "$Root/specs/$Slug.md"
$TasksFile = "$Root/plans/$Slug/tasks.md"
$VerifyTaskState = "$Root/scripts/verify-task-state.ps1"
$ArchiveDir = "$Root/archive/$Slug"

$missing = $false
foreach ($Path in @($SpecFile)) {
  if (-not (Test-Path -LiteralPath $Path)) {
    Write-Output "ERROR: missing required file: $Path"
    $missing = $true
  }
}

if ($missing) {
  exit 1
}

if ($Status -eq "completed" -and (Test-Path -LiteralPath $TasksFile) -and (Test-Path -LiteralPath $VerifyTaskState)) {
  $taskStateOutput = & powershell -NoProfile -ExecutionPolicy Bypass -File $VerifyTaskState $Slug
  $taskStateOutput | ForEach-Object { Write-Output $_ }
  $openLine = $taskStateOutput | Where-Object { $_ -like "TASKS_OPEN=*" } | Select-Object -First 1
  if ($openLine) {
    $openCount = [int]($openLine -replace "^TASKS_OPEN=", "")
    if ($openCount -gt 0) {
      Write-Output "ERROR: completed archive requested while open tasks remain"
      exit 1
    }
  }
}

New-Item -ItemType Directory -Path $ArchiveDir -Force | Out-Null
Write-Output "OK: archive is ready for slug '$Slug' with status '$Status'"
