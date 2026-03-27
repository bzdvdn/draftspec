#!/bin/sh

set -eu

SLUG="${1:-}"
STATUS="${2:-}"
REASON="${3:-}"

if [ -z "$SLUG" ]; then
  echo "usage: check-archive-ready.sh <slug> <status> <reason>"
  exit 1
fi

if [ -z "$STATUS" ]; then
  echo "ERROR: archive status is required"
  exit 1
fi

if [ -z "$REASON" ]; then
  echo "ERROR: archive reason is required"
  exit 1
fi

case "$STATUS" in
  completed|superseded|abandoned|rejected|deferred) ;;
  *)
    echo "ERROR: invalid archive status: $STATUS"
    exit 1
    ;;
esac

ROOT=".draftspec"
SPEC_FILE="$ROOT/specs/$SLUG.md"
MEMORY_FILE="$ROOT/memory.md"
ARCHIVE_DIR="$ROOT/archive/$SLUG"

missing=0
for path in "$SPEC_FILE" "$MEMORY_FILE"; do
  if [ ! -e "$path" ]; then
    echo "ERROR: missing required file: $path"
    missing=1
  fi
done

if [ "$missing" -ne 0 ]; then
  exit 1
fi

mkdir -p "$ARCHIVE_DIR"
echo "OK: archive is ready for slug '$SLUG' with status '$STATUS'"
