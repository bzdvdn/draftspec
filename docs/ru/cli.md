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

### `draftspec refresh [path]`

Обновляет только Draftspec-managed generated artifacts в уже существующем проекте.

Эта команда обновляет:

- `.draftspec/draftspec.yaml`
- `.draftspec/templates/**`
- `.draftspec/scripts/**`
- project-local agent command files
- managed Draftspec block внутри `AGENTS.md`

Эта команда не обновляет:

- `.draftspec/constitution.md`
- `.draftspec/specs/**`
- `.draftspec/plans/**`
- `.draftspec/archive/**`

Примеры:

```bash
draftspec refresh my-project
draftspec refresh my-project --shell powershell --agents claude --dry-run
draftspec refresh my-project --agent-lang ru --json
```

Важные флаги:

- `--lang`, `--docs-lang`, `--agent-lang`, `--comments-lang` переопределяют существующие language settings из config
- `--shell` переопределяет семейство генерируемых workflow scripts
- `--agents` переопределяет набор включенных project-local agent targets
- `--dry-run` показывает pending changes без записи на диск
- `--json` выводит результат refresh в JSON

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

### `draftspec feature <slug> [path]`

Показывает подробную workflow-карточку одной фичи.

Текстовый вывод включает:

- текущую фазу и `ready_for`
- статус inspect и verify, если отчеты существуют
- прогресс задач, если существует `tasks.md`
- сгруппированные workflow-findings
- короткую подсказку `focus` о наиболее вероятном следующем действии

Используй `--json`, чтобы получить структурированное состояние и feature-local findings.

### `draftspec feature repair <slug> [path]`

Исправляет безопасные feature-local проблемы Draftspec.

Сейчас repair умеет, в частности, переносить legacy inspect report из:

- `.draftspec/plans/<slug>/inspect.md`

в канонический путь:

- `.draftspec/specs/<slug>.inspect.md`

Используй `--dry-run`, чтобы посмотреть изменения без применения, и `--json` для структурированного вывода.

### `draftspec features [path]`

Показывает workflow-состояние по всем найденным фичам.

Текстовый вывод суммирует:

- фазу и `ready_for`
- verdict для inspect и verify
- прогресс задач
- сгруппированные issue counts
- наличие артефактов

Используй `--json`, если нужен machine-readable output.

### `draftspec migrate [path]`

Запускает безопасные project-wide миграции Draftspec.

Сейчас основная область миграции — каноникализация legacy inspect reports по всему проекту.

### `draftspec list-specs [path]`

Показывает список spec slug'ов из `.draftspec/specs/`.

### `draftspec show-spec <name> [path]`

Печатает одну спецификацию по slug.
