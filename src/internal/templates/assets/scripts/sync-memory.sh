#!/bin/sh

set -eu

CONSTITUTION_FILE="${1:-.draftspec/constitution.md}"
MEMORY_FILE="${2:-.draftspec/memory.md}"
CONFIG_FILE="${DRAFTSPEC_CONFIG:-.draftspec/draftspec.yaml}"

if [ ! -f "$CONSTITUTION_FILE" ]; then
  echo "constitution file not found: $CONSTITUTION_FILE" >&2
  exit 1
fi

if [ ! -f "$MEMORY_FILE" ]; then
  echo "memory file not found: $MEMORY_FILE" >&2
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
    PURPOSE_SECTION="Назначение"
    CORE_SECTION="Ключевые принципы"
    CONSTRAINTS_SECTION="Ограничения"
    LAST_UPDATED_SECTION="Последнее обновление"
    SUMMARY_HEADING="## Сводка по конституции"
    LAST_UPDATED_LABEL="Последнее обновление"
    PURPOSE_LABEL="Назначение"
    PRINCIPLES_LABEL="Принципы"
    CONSTRAINTS_LABEL="Ограничения"
    UNKNOWN_VALUE="неизвестно"
    UNSUMMARIZED_VALUE="еще не заполнено"
    ;;
  *)
    PURPOSE_SECTION="Purpose"
    CORE_SECTION="Core Principles"
    CONSTRAINTS_SECTION="Constraints"
    LAST_UPDATED_SECTION="Last Updated"
    SUMMARY_HEADING="## Constitution Summary"
    LAST_UPDATED_LABEL="Last Updated"
    PURPOSE_LABEL="Purpose"
    PRINCIPLES_LABEL="Principles"
    CONSTRAINTS_LABEL="Constraints"
    UNKNOWN_VALUE="unknown"
    UNSUMMARIZED_VALUE="not yet summarized"
    ;;
esac

extract_section() {
  section="$1"
  awk -v section="$section" '
    $0 == "## " section { in_section=1; next }
    /^## / && in_section { exit }
    in_section { print }
  ' "$CONSTITUTION_FILE"
}

last_updated=$(extract_section "$LAST_UPDATED_SECTION" | sed '/^[[:space:]]*$/d' | head -n 1 | tr -d '\r')
purpose=$(extract_section "$PURPOSE_SECTION" | sed '/^[[:space:]]*$/d' | paste -sd ' ' -)
principles=$(awk -v section="$CORE_SECTION" '
  $0 == "## " section { in_core=1; next }
  /^## / && in_core { exit }
  in_core && /^### / { sub(/^### /, ""); print }
' "$CONSTITUTION_FILE")
constraints=$(extract_section "$CONSTRAINTS_SECTION" | sed '/^[[:space:]]*$/d')

if [ -z "$last_updated" ]; then
  last_updated="$UNKNOWN_VALUE"
fi
if [ -z "$purpose" ]; then
  purpose="$UNSUMMARIZED_VALUE"
fi

summary_file=$(mktemp)
{
  echo "$SUMMARY_HEADING"
  echo
  echo "- $LAST_UPDATED_LABEL: $last_updated"
  echo "- $PURPOSE_LABEL: $purpose"
  echo "- $PRINCIPLES_LABEL:"
  if [ -n "$principles" ]; then
    printf '%s\n' "$principles" | sed 's/^/  - /'
  fi
  echo "- $CONSTRAINTS_LABEL:"
  if [ -n "$constraints" ]; then
    printf '%s\n' "$constraints" | sed 's/^/  /'
  fi
  echo
} > "$summary_file"

output_file=$(mktemp)
awk -v heading="$SUMMARY_HEADING" -v replacement="$summary_file" '
  BEGIN { skip=0 }
  $0 == heading {
    while ((getline line < replacement) > 0) print line
    close(replacement)
    skip=1
    next
  }
  /^## / && skip {
    skip=0
  }
  !skip { print }
' "$MEMORY_FILE" > "$output_file"

mv "$output_file" "$MEMORY_FILE"
rm -f "$summary_file"

echo "updated $MEMORY_FILE from $CONSTITUTION_FILE"
