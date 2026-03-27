#!/bin/sh

set -eu

CONSTITUTION_FILE="${1:-.draftspec/constitution.md}"
CONFIG_FILE="${DRAFTSPEC_CONFIG:-.draftspec/draftspec.yaml}"

if [ ! -f "$CONSTITUTION_FILE" ]; then
  echo "ERROR: constitution file not found: $CONSTITUTION_FILE"
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
    PURPOSE="Назначение"
    CORE="Ключевые принципы"
    CONSTRAINTS="Ограничения"
    LANGUAGE_POLICY="Языковая политика"
    WORKFLOW="Процесс разработки"
    GOVERNANCE="Управление"
    LAST_UPDATED="Последнее обновление"
    ;;
  *)
    PURPOSE="Purpose"
    CORE="Core Principles"
    CONSTRAINTS="Constraints"
    LANGUAGE_POLICY="Language Policy"
    WORKFLOW="Development Workflow"
    GOVERNANCE="Governance"
    LAST_UPDATED="Last Updated"
    ;;
esac

check_section() {
  section="$1"
  if grep -q "^## $section$" "$CONSTITUTION_FILE"; then
    echo "OK: $section"
  else
    echo "ERROR: missing section: $section"
  fi
}

check_section "$PURPOSE"
check_section "$CORE"
check_section "$CONSTRAINTS"
check_section "$LANGUAGE_POLICY"
check_section "$WORKFLOW"
check_section "$GOVERNANCE"
check_section "$LAST_UPDATED"

principles_count=$(grep -c '^### ' "$CONSTITUTION_FILE" || true)
if [ "$principles_count" -ge 5 ]; then
  echo "OK: principles count is $principles_count"
else
  echo "ERROR: expected at least 5 principles, found $principles_count"
fi

if grep -q '\[[A-Z0-9_][A-Z0-9_ -]*\]' "$CONSTITUTION_FILE"; then
  echo "WARN: placeholder tokens remain in constitution"
else
  echo "OK: no placeholder tokens detected"
fi
