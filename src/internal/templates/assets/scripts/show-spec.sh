#!/bin/sh

set -eu

SPEC_NAME="${1:-}"
SPECS_DIR="${2:-.draftspec/specs}"

if [ -z "$SPEC_NAME" ]; then
  echo "usage: show-spec.sh <spec-name> [specs-dir]" >&2
  exit 1
fi

SPEC_FILE="$SPECS_DIR/$SPEC_NAME.md"

if [ ! -f "$SPEC_FILE" ]; then
  echo "spec not found: $SPEC_FILE" >&2
  exit 1
fi

cat "$SPEC_FILE"
