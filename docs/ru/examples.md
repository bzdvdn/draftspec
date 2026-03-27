# Примеры

На этой странице собраны реалистичные end-to-end сценарии Draftspec для одного feature package.

## 1. Создание Конституции для Brownfield-проекта

Пример запроса:

```text
/draftspec.constitution Python-проект в стиле DDD, разделен на API и workers, Kafka для асинхронной интеграции, ClickHouse как аналитический sink.
```

Ожидаемое поведение агента:

- прочитать prompt `.draftspec/templates/prompts/constitution.md`
- собрать только минимально нужные evidence из репозитория
- создать или обновить `.draftspec/constitution.md`
- обновить `.draftspec/memory.md`
- при необходимости запустить `check-constitution.sh` и `sync-memory.sh`

Ожидаемый результат:

- архитектурные правила формализованы
- правила разработки зафиксированы явно
- конституция становится главным документом для следующих фаз

## 2. Создание Spec

Пример запроса:

```text
/draftspec.spec Добавить partner-specific расписание ingestion с override для retry policy.
```

Ожидаемое поведение агента:

- сначала прочитать constitution и memory
- создать `.draftspec/specs/partner-scheduling.md`
- записать acceptance criteria в каноническом формате `Given / When / Then`
- остальной текст держать на configured documentation language

Пример acceptance criterion:

```md
### Acceptance Criterion 1

- **Given** у партнера задана собственная retry policy
- **When** рассчитывается расписание ingestion
- **Then** worker использует partner-specific retry window вместо default policy
```

## 3. Проверка Spec через Inspect

Пример запроса:

```text
/draftspec.inspect partner-scheduling
```

Ожидаемое поведение агента:

- прочитать constitution, memory и `.draftspec/specs/partner-scheduling.md`
- проверить полноту, соответствие конституции и качество сценариев
- выпустить focused inspection report

Типовые находки:

- отсутствует failure-path сценарий
- непонятно покрытие для manual retry overrides
- есть открытый вопрос про ownership scheduler logic

## 4. Создание Plan Package

Пример запроса:

```text
/draftspec.plan partner-scheduling
```

Ожидаемое поведение агента:

- прочитать constitution, memory и spec
- создать `.draftspec/plans/partner-scheduling/plan.md`
- создать `.draftspec/plans/partner-scheduling/data-model.md`
- создать `.draftspec/plans/partner-scheduling/contracts/`
- создавать `research.md` только если действительно есть неопределенность

Типовые выходы:

- plan по integration points для scheduler
- data model для partner overrides и retry windows
- event или API contracts для обновления конфигурации

## 5. Создание Tasks

Пример запроса:

```text
/draftspec.tasks partner-scheduling
```

Ожидаемое поведение агента:

- использовать `plan.md` как decomposition entrypoint
- подтягивать spec, contracts или data model только при необходимости
- создать `.draftspec/plans/partner-scheduling/tasks.md`
- включить acceptance-to-task coverage

Пример структуры задач:

```md
## Phase 1: Data Model

- [ ] Add partner scheduling override model
- [ ] Persist retry window fields

## Acceptance Coverage

- Acceptance Criterion 1 -> Task 1, Task 2
```

## 6. Реализация Фичи

Пример запроса:

```text
/draftspec.implement partner-scheduling
```

Ожидаемое поведение агента:

- стартовать от `tasks.md`
- читать spec, plan, data model или contracts только для активной задачи
- выполнять незавершенные задачи по порядку
- обновлять `tasks.md`
- обновлять `.draftspec/memory.md`

Эта фаза не должна читать широкий контекст репозитория без реальной необходимости.

## 7. Архивация Фичи

Пример запроса:

```text
/draftspec.archive partner-scheduling --status completed --reason "implemented and merged"
```

Ожидаемое поведение агента:

- скопировать feature package в `.draftspec/archive/partner-scheduling/<YYYY-MM-DD>/`
- записать `summary.md`
- добавить короткую archived-запись в `memory.md`

Ожидаемый результат архива:

```text
.draftspec/archive/
  partner-scheduling/
    2026-03-28/
      summary.md
      spec.md
      plan.md
      tasks.md
      data-model.md
      memory-snapshot.md
      contracts/
```

## 8. Сценарий Обслуживания Агентов

Практический maintenance flow для agent targets:

```bash
draftspec add-agent my-project --agents claude --agents cursor
draftspec list-agents my-project
draftspec remove-agent my-project --agents cursor
draftspec cleanup-agents my-project
draftspec doctor my-project
```

Этот сценарий полезен, когда проект со временем меняет предпочитаемый набор агентов.
