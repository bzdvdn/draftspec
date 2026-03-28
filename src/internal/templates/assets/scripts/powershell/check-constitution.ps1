$ErrorActionPreference = "Stop"

$ConstitutionFile = if ($args.Count -gt 0) { $args[0] } else { ".draftspec/constitution.md" }
$ConfigFile = if ($env:DRAFTSPEC_CONFIG) { $env:DRAFTSPEC_CONFIG } else { ".draftspec/draftspec.yaml" }

if (-not (Test-Path -LiteralPath $ConstitutionFile)) {
  Write-Output "ERROR: constitution file not found: $ConstitutionFile"
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

$DocsLang = Detect-DocsLang
if ($DocsLang -eq "ru") {
  $Sections = @("Назначение", "Ключевые принципы", "Ограничения", "Языковая политика", "Процесс разработки", "Управление", "Последнее обновление")
} else {
  $Sections = @("Purpose", "Core Principles", "Constraints", "Language Policy", "Development Workflow", "Governance", "Last Updated")
}

$Content = Get-Content -LiteralPath $ConstitutionFile
foreach ($Section in $Sections) {
  if ($Content -contains "## $Section") {
    Write-Output "OK: $Section"
  } else {
    Write-Output "ERROR: missing section: $Section"
  }
}

$principlesCount = ($Content | Where-Object { $_ -match '^### ' }).Count
if ($principlesCount -ge 5) {
  Write-Output "OK: principles count is $principlesCount"
} else {
  Write-Output "ERROR: expected at least 5 principles, found $principlesCount"
}

$placeholderFound = $false
foreach ($Line in $Content) {
  if ($Line -match '\[[A-Z0-9_][A-Z0-9_ -]*\]') {
    $placeholderFound = $true
    break
  }
}
if ($placeholderFound) {
  Write-Output "WARN: placeholder tokens remain in constitution"
} else {
  Write-Output "OK: no placeholder tokens detected"
}
