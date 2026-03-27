#!/bin/sh

set -eu

SLUG="${1:-}"

if [ -z "$SLUG" ]; then
  echo "usage: check-inspect-ready.sh <slug>"
  exit 1
fi

ROOT=".draftspec"
SPEC_FILE="$ROOT/specs/$SLUG.md"
CONSTITUTION_FILE="$ROOT/constitution.md"
MEMORY_FILE="$ROOT/memory.md"

missing=0

for path in "$CONSTITUTION_FILE" "$MEMORY_FILE" "$SPEC_FILE"; do
  if [ ! -e "$path" ]; then
    echo "ERROR: missing required file: $path"
    missing=1
  fi
done

if [ "$missing" -ne 0 ]; then
  exit 1
fi

echo "OK: inspect is ready for slug '$SLUG'"
