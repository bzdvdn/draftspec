#!/bin/sh

set -eu

AGENTS_FILE="${1:-AGENTS.md}"
SNIPPET_FILE="${2:-.draftspec/templates/agents-snippet.md}"

if [ ! -f "$AGENTS_FILE" ]; then
  : > "$AGENTS_FILE"
fi

if grep -q "^## Draftspec$" "$AGENTS_FILE"; then
  echo "Draftspec block already present in $AGENTS_FILE"
  exit 0
fi

printf "\n" >> "$AGENTS_FILE"
cat "$SNIPPET_FILE" >> "$AGENTS_FILE"

echo "Draftspec block added to $AGENTS_FILE"
