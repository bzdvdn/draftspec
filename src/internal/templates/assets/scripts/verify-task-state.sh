#!/bin/sh

set -eu

SPEC_SLUG="${1:-}"
TASKS_FILE=".draftspec/plans/$SPEC_SLUG/tasks.md"

if [ -z "$SPEC_SLUG" ]; then
  echo "usage: verify-task-state.sh <spec-slug>" >&2
  exit 1
fi

if [ ! -f "$TASKS_FILE" ]; then
  echo "ERROR: missing $TASKS_FILE"
  exit 1
fi

total=$(grep -E '^- \[[ x]\]' "$TASKS_FILE" | wc -l | tr -d ' ')
completed=$(grep -E '^- \[x\]' "$TASKS_FILE" | wc -l | tr -d ' ')
open=$(grep -E '^- \[ \]' "$TASKS_FILE" | wc -l | tr -d ' ')
task_ids=$(grep -Eo 'T[0-9]+\.[0-9]+' "$TASKS_FILE" | wc -l | tr -d ' ')
coverage_lines=$(grep -E 'AC-[0-9][0-9][0-9].*->.*T[0-9]+\.[0-9]+' "$TASKS_FILE" | wc -l | tr -d ' ')

if [ "$total" -eq 0 ]; then
  echo "ERROR: no task checkboxes found in $TASKS_FILE"
  exit 1
fi

echo "TASKS_TOTAL=$total"
echo "TASKS_COMPLETED=$completed"
echo "TASKS_OPEN=$open"
echo "TASK_IDS=$task_ids"
echo "AC_COVERAGE_LINES=$coverage_lines"

if [ "$task_ids" -eq 0 ]; then
  echo "ERROR: no stable task IDs found in $TASKS_FILE"
  exit 1
fi

if [ "$coverage_lines" -eq 0 ]; then
  echo "WARN: no AC-to-task coverage lines found in $TASKS_FILE"
fi

if [ "$open" -gt 0 ]; then
  echo "WARN: open tasks remain in $TASKS_FILE"
  exit 0
fi

echo "OK: all tasks are marked complete in $TASKS_FILE"
