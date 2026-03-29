#!/bin/sh

set -eu

SLUG="${1:-}"
ROOT=".draftspec"
SPEC_FILE="$ROOT/specs/$SLUG.md"
CONSTITUTION_FILE="$ROOT/constitution.md"
INSPECT_REPORT_TEMPLATE="$ROOT/templates/inspect-report.md"
INSPECT_PROMPT="$ROOT/templates/prompts/inspect.md"
INSPECT_SCRIPT="$ROOT/scripts/inspect-spec.sh"
missing=0

if [ -z "$SLUG" ]; then
  echo "usage: check-inspect-ready.sh <slug>"
  exit 1
fi

for path in "$CONSTITUTION_FILE" "$SPEC_FILE" "$INSPECT_REPORT_TEMPLATE" "$INSPECT_PROMPT"; do
  if [ ! -e "$path" ]; then
    echo "ERROR: missing required file: $path"
    missing=1
  else
    echo "OK: $path"
  fi
done

if [ "$missing" -ne 0 ]; then
  exit 1
fi

if [ -x "$INSPECT_SCRIPT" ]; then
  "$INSPECT_SCRIPT" "$SPEC_FILE"
fi

echo "OK: inspect is ready for slug '$SLUG'"
