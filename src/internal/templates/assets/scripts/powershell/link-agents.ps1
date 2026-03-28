$ErrorActionPreference = "Stop"

$AgentsFile = if ($args.Count -gt 0) { $args[0] } else { "AGENTS.md" }
$SnippetFile = if ($args.Count -gt 1) { $args[1] } else { ".draftspec/templates/agents-snippet.md" }

if (-not (Test-Path -LiteralPath $AgentsFile)) {
  New-Item -ItemType File -Path $AgentsFile -Force | Out-Null
}

$content = Get-Content -LiteralPath $AgentsFile -Raw
if ($content -match '(?m)^## Draftspec$') {
  Write-Output "Draftspec block already present in $AgentsFile"
  exit 0
}

$snippet = Get-Content -LiteralPath $SnippetFile -Raw
if ($content.Length -gt 0 -and -not $content.EndsWith("`n")) {
  $content += "`n"
}
$content += "`n" + $snippet
Set-Content -LiteralPath $AgentsFile -Value $content -NoNewline
Write-Output "Draftspec block added to $AgentsFile"
