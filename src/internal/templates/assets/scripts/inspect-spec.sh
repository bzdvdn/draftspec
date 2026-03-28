#!/bin/sh

set -eu

SPEC_FILE="${1:-}"
TASKS_FILE="${2:-}"
CONFIG_FILE="${DRAFTSPEC_CONFIG:-.draftspec/draftspec.yaml}"

if [ -z "$SPEC_FILE" ]; then
  echo "usage: inspect-spec.sh <spec-file> [tasks-file]"
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

extract_section() {
  file_path="$1"
  section="$2"
  awk -v section="$section" '
    $0 == "## " section { in_section=1; next }
    in_section && /^## / { exit }
    in_section { print }
  ' "$file_path"
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
    COVERAGE="Покрытие критериев приемки"
    ;;
  *)
    GOAL="Goal"
    CONTEXT="Context"
    REQUIREMENTS="Requirements"
    ACCEPTANCE="Acceptance Criteria"
    QUESTIONS="Open Questions"
    COVERAGE="Acceptance Coverage"
    ;;
esac

errors=0
warnings=0

error() {
  echo "ERROR: $1"
  errors=$((errors + 1))
}

warn() {
  echo "WARN: $1"
  warnings=$((warnings + 1))
}

ok() {
  echo "OK: $1"
}

check_required_section() {
  section="$1"
  if grep -q "^## $section$" "$SPEC_FILE"; then
    ok "$section"
  else
    error "missing required section: $section"
  fi
}

check_optional_section() {
  section="$1"
  if grep -q "^## $section$" "$SPEC_FILE"; then
    ok "$section"
  else
    warn "missing section: $section"
  fi
}

check_required_section "$GOAL"
check_optional_section "$CONTEXT"
check_required_section "$REQUIREMENTS"
check_required_section "$ACCEPTANCE"
check_optional_section "$QUESTIONS"

acceptance_body="$(extract_section "$SPEC_FILE" "$ACCEPTANCE")"
if [ -z "$(printf '%s' "$acceptance_body" | tr -d '[:space:]')" ]; then
  error "empty acceptance criteria section"
else
  printf '%s
' "$acceptance_body" | grep -q 'Given' && ok 'Given marker found' || error 'missing Given marker in acceptance criteria'
  printf '%s
' "$acceptance_body" | grep -q 'When' && ok 'When marker found' || error 'missing When marker in acceptance criteria'
  printf '%s
' "$acceptance_body" | grep -q 'Then' && ok 'Then marker found' || error 'missing Then marker in acceptance criteria'
fi

criteria_count=$(printf '%s
' "$acceptance_body" | grep -c '^### ' || true)
if [ "$criteria_count" -gt 0 ]; then
  ok "acceptance criteria count: $criteria_count"
else
  warn "no explicit acceptance criterion headings found"
fi

if [ -n "$TASKS_FILE" ] && [ -f "$TASKS_FILE" ]; then
  if grep -q "^## $COVERAGE$" "$TASKS_FILE"; then
    ok "$COVERAGE"
    coverage_body="$(extract_section "$TASKS_FILE" "$COVERAGE")"
    coverage_lines=$(printf '%s
' "$coverage_body" | grep -c -- '->' || true)
    if [ "$criteria_count" -gt 0 ] && [ "$coverage_lines" -lt "$criteria_count" ]; then
      error "acceptance coverage entries ($coverage_lines) are fewer than acceptance criteria ($criteria_count)"
    else
      ok "acceptance coverage entries: $coverage_lines"
    fi
  else
    error "tasks file is missing required section: $COVERAGE"
  fi
fi

echo "SUMMARY: errors=$errors warnings=$warnings"

if [ "$errors" -ne 0 ]; then
  exit 1
fi
