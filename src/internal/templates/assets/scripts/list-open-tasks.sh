#!/bin/sh

set -eu

SPEC_SLUG="${1:-}"
TASKS_FILE=".draftspec/plans/$SPEC_SLUG/tasks.md"

if [ -z "$SPEC_SLUG" ]; then
  echo "usage: list-open-tasks.sh <spec-slug>" >&2
  exit 1
fi

if [ ! -f "$TASKS_FILE" ]; then
  echo "tasks file not found: $TASKS_FILE" >&2
  exit 1
fi

grep '^\- \[ \]' "$TASKS_FILE" || true
