#!/bin/sh

set -eu

SPEC_SLUG="${1:-}"
errors=0

if [ -z "$SPEC_SLUG" ]; then
  echo "usage: check-plan-ready.sh <spec-slug>" >&2
  exit 1
fi

check_file() {
  path="$1"
  if [ -f "$path" ]; then
    echo "OK: $path"
  else
    echo "ERROR: missing $path"
    errors=$((errors + 1))
  fi
}

check_contains() {
  path="$1"
  pattern="$2"
  label="$3"
  if grep -Eq "$pattern" "$path"; then
    echo "OK: $label"
  else
    echo "ERROR: $label"
    errors=$((errors + 1))
  fi
}

check_file ".draftspec/constitution.md"
check_file ".draftspec/specs/$SPEC_SLUG.md"
check_file ".draftspec/templates/plan.md"
check_file ".draftspec/templates/prompts/plan.md"

SPEC_FILE=".draftspec/specs/$SPEC_SLUG.md"
if [ -f "$SPEC_FILE" ]; then
  check_contains "$SPEC_FILE" '^## (Критерии приемки|Acceptance Criteria)$' "spec has acceptance criteria section"
  check_contains "$SPEC_FILE" 'AC-[0-9][0-9][0-9]' "spec has stable acceptance IDs"
fi

if [ -x ".draftspec/scripts/inspect-spec.sh" ] && [ -f "$SPEC_FILE" ]; then
  ./.draftspec/scripts/inspect-spec.sh "$SPEC_FILE"
fi

if [ "$errors" -ne 0 ]; then
  exit 1
fi
