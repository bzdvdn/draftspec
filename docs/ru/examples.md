# Примеры

На этой странице собраны реалистичные end-to-end сценарии Draftspec для одного feature package.

## Быстрые Сценарии Использования

### Новый проект

Когда проект создается с нуля, Draftspec лучше вводить как минимальный каркас проектного контекста с самого начала.

Пример:

```bash
draftspec init my-project --lang ru --agents codex
cd my-project
draftspec doctor .
```

Что делать дальше:

- сформировать `constitution` для правил проекта
- описать первую фичу через `spec`
- подготовить `plan` и `tasks`
- выполнять `implement` только от текущего task list

Практический смысл такого старта:

- команда и агент сразу работают от одного набора правил
- контекст проекта с самого начала остается явным и редактируемым
- workflow остается легким, потому что Draftspec не требует тяжелого process engine

### Уже существующий проект

Для brownfield-проекта Draftspec лучше вводить постепенно, а не пытаться сразу описать всю кодовую базу.

Пример:

```bash
cd existing-project
draftspec init . --lang ru --agents codex
draftspec doctor .
```

Рекомендуемый старт:

- сначала зафиксировать `constitution` под текущую реальность проекта
- выбрать одну активную фичу или change request
- создать spec только для нее
- переходить к plan, tasks и implement только внутри этой feature scope

Чего не стоит делать:

- не пытаться сразу документировать весь проект
- не тянуть широкий repository context, если текущая фича этого не требует

Практический смысл такого входа:

- Draftspec добавляет легкий слой дисциплины поверх уже существующей кодовой базы
- adoption идет по одной фиче за раз
- это снижает токеноемкость и уменьшает риск бюрократии

## 1. Создание Конституции для Brownfield-проекта

Пример запроса:

```text
/draftspec.constitution Python-проект в стиле DDD, разделен на API и workers, Kafka для асинхронной интеграции, ClickHouse как аналитический sink.
```

Ожидаемое поведение агента:

- прочитать prompt `.draftspec/templates/prompts/constitution.md`
- собрать только минимально нужные evidence из репозитория
- создать или обновить `.draftspec/constitution.md`
- при необходимости запустить `check-constitution.sh`

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

- сначала прочитать constitution
- создать `.draftspec/specs/partner-scheduling.md`
- записать acceptance criteria в каноническом формате `Given / When / Then`
- остальной текст держать на configured documentation language

Пример acceptance criterion:

```md
### Acceptance Criterion 1

- ID: AC-001
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

- прочитать constitution и `.draftspec/specs/partner-scheduling.md`
- проверить полноту, соответствие конституции и качество сценариев
- выпустить focused inspection report
- если отчет нужно сохранять на диск, до планирования предпочитать `.draftspec/specs/partner-scheduling.inspect.md`, а после появления plan package — `.draftspec/plans/partner-scheduling/inspect.md`
- использовать `.draftspec/templates/inspect-report.md` как канонический шаблон отчета

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

- прочитать constitution и spec
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

- AC-001 -> Task 1, Task 2
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

Эта фаза не должна читать широкий контекст репозитория без реальной необходимости.

## 7. Verify Фичи

Пример запроса:

```text
/draftspec.verify partner-scheduling
```

Ожидаемое поведение агента:

- сначала прочитать constitution и tasks
- подтвердить, что завершенные задачи достаточно соответствуют текущему состоянию реализации
- выпустить легкий verification report
- начинать с `.draftspec/scripts/verify-task-state.sh partner-scheduling`, если сначала нужно только подтвердить состояние задач
- использовать `.draftspec/templates/verify-report.md`, если отчет нужно сохранить в файл
- по умолчанию использовать `.draftspec/plans/partner-scheduling/verify.md`, если путь явно не указан

## 8. Архивация Фичи

Пример запроса:

```text
/draftspec.archive partner-scheduling --status completed --reason "implemented and merged"
```

Ожидаемое поведение агента:

- для статуса `completed` сначала запустить `.draftspec/scripts/verify-task-state.sh partner-scheduling` и остановиться, если открытые задачи еще остались
- скопировать feature package в `.draftspec/archive/partner-scheduling/<YYYY-MM-DD>/`
- записать `summary.md`

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
      contracts/
```

## 9. Сценарий Обслуживания Агентов

Практический maintenance flow для agent targets:

```bash
draftspec add-agent my-project --agents claude --agents cursor
draftspec list-agents my-project
draftspec remove-agent my-project --agents cursor
draftspec cleanup-agents my-project
draftspec doctor my-project
```

Этот сценарий полезен, когда проект со временем меняет предпочитаемый набор агентов.
