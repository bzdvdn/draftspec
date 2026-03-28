# Модель Workflow

## Строгая цепочка фаз

```text
constitution -> spec -> inspect -> plan -> tasks -> implement -> verify -> archive
```

## Роли фаз

### `constitution`

Определяет неоспоримые правила проекта.

Обязательные секции:

- `Purpose`
- `Core Principles`
- `Constraints`
- `Language Policy`
- `Development Workflow`
- `Governance`
- `Last Updated`

### `spec`

Описывает одну фичу как конкретную спецификацию. Acceptance criteria должны использовать канонические маркеры `Given / When / Then`, даже если остальной текст документа на русском.

### `inspect`

Проверяет качество и согласованность одной фичи. Фаза может находить отсутствующие сценарии, слабые acceptance criteria, конфликт с конституцией, drift между spec и plan или отсутствие покрытия задачами.

Полный inspection report должен использовать стабильную структуру:

- `# Inspect Report: <slug>`
- `## Scope`
- `## Verdict`
- `## Errors`
- `## Warnings`
- `## Questions`
- `## Suggestions`
- `## Traceability`
- `## Next Step`

`Verdict` должен быть одним из значений:

- `pass`
- `concerns`
- `blocked`

Рекомендуемая семантика:

- `pass`: блокирующих проблем нет; остаются только незначительные предупреждения или их нет совсем
- `concerns`: workflow можно продолжать, но warnings или открытые вопросы желательно закрыть в ближайшее время
- `blocked`: следующая фаза иначе продолжила бы работу с отсутствующей или противоречивой информацией

Если inspection report нужно сохранить на диск, Draftspec должен по умолчанию использовать такие пути:

- `.draftspec/plans/<slug>/inspect.md`, если plan package уже существует
- `.draftspec/specs/<slug>.inspect.md`, если plan package еще не существует

Используйте `.draftspec/templates/inspect-report.md` как канонический шаблон, если отчет записывается в файл.

Стабильные идентификаторы критериев вроде `AC-001` делают traceability легче и проще для валидации.

### `plan`

Создает технические артефакты для одного feature package:

- `plan.md`
- `data-model.md`
- `contracts/`
- optional `research.md`

### `tasks`

Преобразует пакет плана в исполнимые задачи. `tasks.md` лежит рядом с остальными plan artifacts внутри `.draftspec/plans/<slug>/`.

### `implement`

Выполняет незавершенные задачи, обновляет `tasks.md` и синхронизирует `memory.md`.

### `verify`

Запускает легкую post-implementation проверку, чтобы подтвердить, что завершенная работа достаточно согласована с tasks, memory и правилами проекта для безопасного следующего шага.

Полный verification report должен использовать стабильную структуру:

- `# Verify Report: <slug>`
- `## Scope`
- `## Verdict`
- `## Checks`
- `## Errors`
- `## Warnings`
- `## Questions`
- `## Next Step`

`Verdict` должен быть одним из значений:

- `pass`
- `concerns`
- `blocked`

Рекомендуемая семантика:

- `pass`: блокирующих проблем нет; остаются только незначительные предупреждения или их нет совсем
- `concerns`: по workflow можно двигаться дальше, но warnings или открытые вопросы желательно закрыть в ближайшее время
- `blocked`: архивирование или заявление о завершенности иначе опирались бы на противоречивое состояние реализации, незавершенную обязательную работу или несинхронизированную memory

Если verification report нужно сохранить на диск, Draftspec должен по умолчанию использовать `.draftspec/plans/<slug>/verify.md`.

Используйте `.draftspec/templates/verify-report.md` как канонический шаблон, если отчет записывается в файл.

Используйте `.draftspec/scripts/verify-task-state.sh <slug>` как самый дешевый helper первого прохода, когда нужно только подтвердить состояние задач.

### `archive`

Копирует завершенный, вытесненный, отклоненный, abandoned или deferred feature package в `.draftspec/archive/<slug>/<YYYY-MM-DD>/`.

Если архивирование идет со статусом `completed`, Draftspec должен сначала использовать `.draftspec/scripts/verify-task-state.sh <slug>` и считать оставшиеся открытые задачи блокером.

## Зачем нужна такая цепочка

Эта модель делает инструмент строгим, но не бюрократичным:

- сначала фиксируются архитектурные и процессные законы
- затем пользовательское намерение превращается в spec
- потом появляется технический plan
- только после этого строятся tasks
- implementation идет по tasks, а не по широкой импровизации
- легкий verify закрывает разрыв между implementation и archive
- завершенные feature packages уходят в archive, не раздувая active memory
