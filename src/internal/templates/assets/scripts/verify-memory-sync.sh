#!/bin/sh

set -eu

SPEC_SLUG="${1:-}"
TASKS_FILE=".draftspec/plans/$SPEC_SLUG/tasks.md"
MEMORY_FILE=".draftspec/memory.md"

if [ -z "$SPEC_SLUG" ]; then
  echo "usage: verify-memory-sync.sh <spec-slug>" >&2
  exit 1
fi

if [ ! -f "$TASKS_FILE" ]; then
  echo "ERROR: missing $TASKS_FILE"
  exit 1
fi

if [ ! -f "$MEMORY_FILE" ]; then
  echo "ERROR: missing $MEMORY_FILE"
  exit 1
fi

total=$(grep -E '^- \[[ x]\]' "$TASKS_FILE" | wc -l | tr -d ' ')
completed=$(grep -E '^- \[x\]' "$TASKS_FILE" | wc -l | tr -d ' ')
open=$(grep -E '^- \[ \]' "$TASKS_FILE" | wc -l | tr -d ' ')

if [ "$total" -eq 0 ]; then
  echo "ERROR: no task checkboxes found in $TASKS_FILE"
  exit 1
fi

has_slug="no"
if grep -Fqi "$SPEC_SLUG" "$MEMORY_FILE"; then
  has_slug="yes"
fi

has_completion_claim="no"
if grep -Ei "($SPEC_SLUG.*(completed|archiv|ready to archive|implementation complete|all tasks done|заверш|архив|готов.*архив))|((completed|archiv|ready to archive|implementation complete|all tasks done|заверш|архив|готов.*архив).*$SPEC_SLUG)" "$MEMORY_FILE" >/dev/null 2>&1; then
  has_completion_claim="yes"
fi

echo "TASKS_TOTAL=$total"
echo "TASKS_COMPLETED=$completed"
echo "TASKS_OPEN=$open"
echo "MEMORY_HAS_SLUG=$has_slug"
echo "MEMORY_HAS_COMPLETION_CLAIM=$has_completion_claim"

if [ "$open" -gt 0 ] && [ "$has_completion_claim" = "yes" ]; then
  echo "ERROR: memory claims completion or archive readiness while open tasks remain"
  exit 1
fi

if [ "$open" -eq 0 ] && [ "$has_slug" = "no" ]; then
  echo "WARN: all tasks are complete but memory does not mention this slug"
  exit 0
fi

echo "OK: coarse memory/task sync looks acceptable"
