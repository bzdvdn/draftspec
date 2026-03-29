#!/bin/sh

set -eu

errors=0

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
check_file ".draftspec/templates/spec.md"
check_file ".draftspec/templates/prompts/spec.md"

if [ -f ".draftspec/templates/spec.md" ]; then
  check_contains ".draftspec/templates/spec.md" '^## (Требования|Requirements)$' "spec template has requirements section"
  check_contains ".draftspec/templates/spec.md" 'RQ-[0-9][0-9][0-9]' "spec template includes requirement IDs"
  check_contains ".draftspec/templates/spec.md" 'AC-[0-9][0-9][0-9]' "spec template includes acceptance IDs"
  check_contains ".draftspec/templates/spec.md" 'Given' "spec template includes Given marker"
  check_contains ".draftspec/templates/spec.md" 'When' "spec template includes When marker"
  check_contains ".draftspec/templates/spec.md" 'Then' "spec template includes Then marker"
fi

if [ "$errors" -ne 0 ]; then
  exit 1
fi
