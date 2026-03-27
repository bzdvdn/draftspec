#!/bin/sh

set -eu

SPECS_DIR="${1:-.draftspec/specs}"

if [ ! -d "$SPECS_DIR" ]; then
  echo "specs directory not found: $SPECS_DIR" >&2
  exit 1
fi

find "$SPECS_DIR" -maxdepth 1 -type f -name '*.md' -exec basename {} .md \; | sort
