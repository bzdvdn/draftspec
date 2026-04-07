#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: required command not found: $1" >&2
    exit 1
  fi
}

require_cmd go
require_cmd vhs

cd "$ROOT_DIR"

mkdir -p bin
go build -o bin/draftspec ./src/cmd/draftspec

rm -rf demo/_work

vhs demo/quick.tape

echo "ok: wrote demo/draftspec-demo.gif"

