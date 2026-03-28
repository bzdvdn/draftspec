# CLI

## Команды

### `draftspec init [path]`

Инициализирует Draftspec workspace в целевом проекте.

Примеры:

```bash
draftspec init
draftspec init my-project --lang ru --shell sh
draftspec init my-project --docs-lang ru --agent-lang en --comments-lang en --shell powershell --agents claude --agents cursor
```

Важные флаги:

- `--git` инициализирует Git-репозиторий; по умолчанию включен
- `--lang` задает базовый язык; по умолчанию `en`
- `--shell` выбирает семейство генерируемых workflow scripts; обязателен: `sh` или `powershell`
- `--docs-lang` задает язык генерируемой документации
- `--agent-lang` задает язык генерируемых промтов и guidance для агентов
- `--comments-lang` фиксирует предпочитаемый язык комментариев в коде
- `--agents` генерирует project-local agent files

### `draftspec add-agent [path]`

Добавляет один или несколько agent targets в уже инициализированный проект.

```bash
draftspec add-agent my-project --agents claude --agents codex
```

### `draftspec list-agents [path]`

Показывает включенные agent targets из `.draftspec/draftspec.yaml`.

### `draftspec remove-agent [path]`

Отключает один или несколько agent targets и удаляет их generated files.

### `draftspec cleanup-agents [path]`

Удаляет осиротевшие agent artifacts, которые больше не соответствуют включенным targets в config.

### `draftspec doctor [path]`

Проверяет здоровье workspace.

`doctor` выводит:

- `error` для отсутствующих обязательных файлов и невалидных значений config
- `warning` для orphaned agent artifacts, которые все еще лежат на диске
- `ok`, когда workspace выглядит здоровым

Используй `--json`, если нужен machine-readable output для automation и CI.

### `draftspec list-specs [path]`

Показывает список spec slug'ов из `.draftspec/specs/`.

### `draftspec show-spec <name> [path]`

Печатает одну спецификацию по slug.
