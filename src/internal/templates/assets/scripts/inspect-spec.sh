#!/bin/sh

set -eu

SPEC_FILE="${1:-}"
CONFIG_FILE="${DRAFTSPEC_CONFIG:-.draftspec/draftspec.yaml}"

if [ -z "$SPEC_FILE" ]; then
  echo "usage: inspect-spec.sh <spec-file>"
  exit 1
fi

if [ ! -f "$SPEC_FILE" ]; then
  echo "ERROR: spec file not found: $SPEC_FILE"
  exit 1
fi

detect_docs_lang() {
  if [ -f "$CONFIG_FILE" ]; then
    awk '
      $0 == "language:" { in_language=1; next }
      in_language && /^[^[:space:]]/ { exit }
      in_language && $1 == "docs:" { print $2; exit }
    ' "$CONFIG_FILE"
  fi
}

DOCS_LANG="$(detect_docs_lang)"
if [ -z "$DOCS_LANG" ]; then
  DOCS_LANG="en"
fi

case "$DOCS_LANG" in
  ru)
    GOAL="Цель"
    CONTEXT="Контекст"
    REQUIREMENTS="Требования"
    ACCEPTANCE="Критерии приемки"
    QUESTIONS="Открытые вопросы"
    ;;
  *)
    GOAL="Goal"
    CONTEXT="Context"
    REQUIREMENTS="Requirements"
    ACCEPTANCE="Acceptance Criteria"
    QUESTIONS="Open Questions"
    ;;
esac

check_section() {
  section="$1"
  if grep -q "^## $section$" "$SPEC_FILE"; then
    echo "OK: $section"
  else
    echo "WARN: missing section: $section"
  fi
}

check_section "$GOAL"
check_section "$CONTEXT"
check_section "$REQUIREMENTS"
check_section "$ACCEPTANCE"
check_section "$QUESTIONS"
