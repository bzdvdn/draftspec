$ErrorActionPreference = "Stop"

$SpecFile = if ($args.Count -gt 0) { $args[0] } else { "" }
$TasksFile = if ($args.Count -gt 1) { $args[1] } else { "" }
$ConfigFile = if ($env:DRAFTSPEC_CONFIG) { $env:DRAFTSPEC_CONFIG } else { ".draftspec/draftspec.yaml" }

if ([string]::IsNullOrWhiteSpace($SpecFile)) {
  Write-Error "usage: inspect-spec.ps1 <spec-file> [tasks-file]"
  exit 1
}
if (-not (Test-Path -LiteralPath $SpecFile)) {
  Write-Output "ERROR: spec file not found: $SpecFile"
  exit 1
}

function Detect-DocsLang {
  if (-not (Test-Path -LiteralPath $ConfigFile)) {
    return "en"
  }
  $inLanguage = $false
  foreach ($Line in Get-Content -LiteralPath $ConfigFile) {
    if ($Line -eq "language:") {
      $inLanguage = $true
      continue
    }
    if ($inLanguage -and $Line -match '^[^\s]') {
      break
    }
    if ($inLanguage -and $Line -match '^\s*docs:\s*(\S+)') {
      return $Matches[1]
    }
  }
  return "en"
}

function Extract-Section([string]$FilePath, [string]$Section) {
  $inSection = $false
  $lines = @()
  foreach ($Line in Get-Content -LiteralPath $FilePath) {
    if ($Line -eq "## $Section") {
      $inSection = $true
      continue
    }
    if ($inSection -and $Line -match '^## ') {
      break
    }
    if ($inSection) {
      $lines += $Line
    }
  }
  return $lines
}

$DocsLang = Detect-DocsLang
if ($DocsLang -eq "ru") {
  $Goal = "Цель"
  $Context = "Контекст"
  $Requirements = "Требования"
  $Acceptance = "Критерии приемки"
  $Questions = "Открытые вопросы"
  $Coverage = "Покрытие критериев приемки"
} else {
  $Goal = "Goal"
  $Context = "Context"
  $Requirements = "Requirements"
  $Acceptance = "Acceptance Criteria"
  $Questions = "Open Questions"
  $Coverage = "Acceptance Coverage"
}

$errors = 0
$warnings = 0
$specContent = Get-Content -LiteralPath $SpecFile

function Add-Error([string]$Message) {
  $script:errors++
  Write-Output "ERROR: $Message"
}

function Add-Warn([string]$Message) {
  $script:warnings++
  Write-Output "WARN: $Message"
}

function Add-OK([string]$Message) {
  Write-Output "OK: $Message"
}

function Check-RequiredSection([string]$Section) {
  if ($specContent -contains "## $Section") {
    Add-OK $Section
  } else {
    Add-Error "missing required section: $Section"
  }
}

function Check-OptionalSection([string]$Section) {
  if ($specContent -contains "## $Section") {
    Add-OK $Section
  } else {
    Add-Warn "missing section: $Section"
  }
}

Check-RequiredSection $Goal
Check-OptionalSection $Context
Check-RequiredSection $Requirements
Check-RequiredSection $Acceptance
Check-OptionalSection $Questions

$acceptanceBody = Extract-Section $SpecFile $Acceptance
$acceptanceText = ($acceptanceBody -join "`n").Trim()
if ([string]::IsNullOrWhiteSpace($acceptanceText)) {
  Add-Error "empty acceptance criteria section"
} else {
  if ($acceptanceText -match 'Given') { Add-OK 'Given marker found' } else { Add-Error 'missing Given marker in acceptance criteria' }
  if ($acceptanceText -match 'When') { Add-OK 'When marker found' } else { Add-Error 'missing When marker in acceptance criteria' }
  if ($acceptanceText -match 'Then') { Add-OK 'Then marker found' } else { Add-Error 'missing Then marker in acceptance criteria' }
}

$criteriaCount = ($acceptanceBody | Where-Object { $_ -match '^### ' }).Count
if ($criteriaCount -gt 0) {
  Add-OK "acceptance criteria count: $criteriaCount"
} else {
  Add-Warn "no explicit acceptance criterion headings found"
}

$acceptanceIDCount = ($acceptanceBody | Where-Object { $_ -match 'AC-[0-9][0-9][0-9]' }).Count
if ($acceptanceIDCount -gt 0) {
  Add-OK "acceptance IDs found: $acceptanceIDCount"
} else {
  Add-Warn "no stable acceptance IDs found in acceptance criteria"
}

if (-not [string]::IsNullOrWhiteSpace($TasksFile) -and (Test-Path -LiteralPath $TasksFile)) {
  $tasksContent = Get-Content -LiteralPath $TasksFile
  if ($tasksContent -contains "## $Coverage") {
    Add-OK $Coverage
    $coverageBody = Extract-Section $TasksFile $Coverage
    $coverageLines = ($coverageBody | Where-Object { $_ -match '->' }).Count
    $malformedLines = ($coverageBody | Where-Object { $_ -match '->' -and $_ -notmatch 'AC-[0-9][0-9][0-9].*->.*T[0-9]+\.[0-9]+' }).Count
    if ($criteriaCount -gt 0 -and $coverageLines -lt $criteriaCount) {
      Add-Error "acceptance coverage entries ($coverageLines) are fewer than acceptance criteria ($criteriaCount)"
    } else {
      Add-OK "acceptance coverage entries: $coverageLines"
    }
    if ($acceptanceIDCount -gt 0 -and $coverageLines -lt $acceptanceIDCount) {
      Add-Error "acceptance coverage entries ($coverageLines) are fewer than acceptance IDs ($acceptanceIDCount)"
    }
    if ($malformedLines -gt 0) {
      Add-Error "acceptance coverage contains malformed entries; expected lines like AC-001 -> T1.1"
    } elseif ($coverageLines -gt 0) {
      Add-OK "acceptance coverage format uses AC and task IDs"
    }
  } else {
    Add-Error "tasks file is missing required section: $Coverage"
  }
}

Write-Output "SUMMARY: errors=$errors warnings=$warnings"
if ($errors -ne 0) {
  exit 1
}
