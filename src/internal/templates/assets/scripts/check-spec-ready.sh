#!/bin/sh

set -eu

check_file() {
  path="$1"
  if [ -f "$path" ]; then
    echo "OK: $path"
  else
    echo "ERROR: missing $path"
  fi
}

check_file ".draftspec/constitution.md"
check_file ".draftspec/templates/spec.md"
check_file ".draftspec/templates/prompts/spec.md"
