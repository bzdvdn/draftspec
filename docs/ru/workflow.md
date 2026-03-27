# Модель Workflow

## Строгая цепочка фаз

```text
constitution -> spec -> inspect -> plan -> tasks -> implement -> archive
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

### `archive`

Копирует завершенный, вытесненный, отклоненный, abandoned или deferred feature package в `.draftspec/archive/<slug>/<YYYY-MM-DD>/`.

## Зачем нужна такая цепочка

Эта модель делает инструмент строгим, но не бюрократичным:

- сначала фиксируются архитектурные и процессные законы
- затем пользовательское намерение превращается в spec
- потом появляется технический plan
- только после этого строятся tasks
- implementation идет по tasks, а не по широкой импровизации
- завершенные feature packages уходят в archive, не раздувая active memory
