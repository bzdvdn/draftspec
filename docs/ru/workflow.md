# Модель Workflow

## Строгая цепочка фаз

```text
constitution -> spec -> inspect -> plan -> tasks -> implement -> verify -> archive
```

Draftspec предполагает branch-based delivery: каждая активная фича должна разрабатываться в своей git-ветке, а общим источником истины служат feature spec и plan package, а не общий mutable memory-файл. Соглашение по умолчанию для веток — `feature/<slug>`.

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

Для agent-facing `/draftspec.spec` Draftspec должен поддерживать optional аргументы:

- `--name <название фичи>`
- `--slug <slug фичи>`
- `--branch <имя ветки>`

Семантика аргументов:

- `--name` задает каноническое имя фичи для текущего spec-запроса
- `--slug` переопределяет slug спецификации
- `--branch` переопределяет только рабочую ветку и не меняет slug спецификации

`/draftspec.spec` должен поддерживать два режима ввода:

- inline mode: имя и описание фичи передаются в одном сообщении
- staged mode: пользователь сначала передает `/draftspec.spec --name ...`, а описание фичи присылает следующим сообщением

Когда `/draftspec.spec` запускается от prompt-файла, Draftspec должен предпочитать метаданные в начале файла, например:

```text
name: Add dark mode
slug: add-dark-mode
```

Правила приоритета для slug:

1. `--slug`
2. `slug:`
3. slug, выведенный из `--name`
4. slug, выведенный из `name:`
5. safe fallback из filename или краткого user request только если он достаточно конкретен

Правила приоритета для имени фичи:

1. `--name`
2. `name:`
3. краткое имя, безопасно выведенное из user request

Если `/draftspec.spec` вызван с `--name`, но подробного описания фичи еще недостаточно для корректной спецификации, Draftspec не должен терять контекст запроса: он должен запросить недостающее описание или интерпретировать следующее сообщение пользователя как продолжение того же spec-запроса.

По умолчанию feature-ветка должна быть `feature/<slug>`. Если пользователь явно передает `--branch <name>`, Draftspec должен использовать это имя ветки вместо default, не меняя при этом slug спецификации.

Сама спецификация при этом должна оставаться branch-agnostic: рабочая ветка относится к execution context, а не к содержимому `spec.md`.

Если запрос неоднозначен, смешивает несколько фич или пытается вывести одну spec из нескольких изменений конституции, Draftspec должен остановиться и запросить одно конкретное изменение до создания ветки и spec.

### `inspect`

Проверяет качество и согласованность одной фичи. Фаза может находить отсутствующие сценарии, слабые acceptance criteria, конфликт с конституцией, drift между spec и plan или отсутствие покрытия задачами.

`inspect` обязателен перед `plan`. Планирование не должно продолжаться, пока у фичи нет сохраненного inspect report в каноническом пути.

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

Если inspection report сохраняется на диск, Draftspec должен использовать канонический путь:

- `.draftspec/specs/<slug>/inspect.md`

Используйте `.draftspec/templates/inspect-report.md` как канонический шаблон, если отчет записывается в файл.
Сохраненные inspect и verify reports должны начинаться с machine-readable metadata block с полями `report_type`, `slug`, `status`, `docs_language` и `generated_at`.

Стабильные идентификаторы критериев вроде `AC-001` делают traceability легче и проще для валидации.

Для дешевой проверки `spec <-> plan` consistency Draftspec должен предпочитать такую область анализа:

- всегда читать: `constitution.md`, `spec.md`
- читать по необходимости: `plan.md`, `tasks.md`
- читать глубже только если этого требует конкретный вывод: `data-model.md`, `contracts/`, `research.md`
- implementation code по умолчанию не читать

Цель этой проверки — поймать явный drift, а не запускать полный архитектурный review. Полезные типы проверок:

- alignment между constitution и spec
- goal alignment
- необоснованное расширение scope
- отражение acceptance-critical behavior на уровне плана
- alignment между plan и tasks, если `tasks.md` уже существует
- соответствие конституции
- оправданность более богатых plan artifacts вроде `data-model.md` и `contracts/`

### `plan`

Создает технические артефакты для одного feature package:

- `plan.md`
- `data-model.md`
- `contracts/`
- optional `research.md`

### `tasks`

Преобразует пакет плана в исполнимые задачи. `tasks.md` лежит рядом с остальными plan artifacts внутри `.draftspec/plans/<slug>/`.

Задачи должны быть сгруппированы по фазам и использовать phase-scoped task IDs вроде `T1.1`, `T1.2` и `T2.1`.

Покрытие критериев приемки должно ссылаться на эти task IDs напрямую:

```text
AC-001 -> T1.1, T2.1
```

### `implement`

Выполняет незавершенные задачи и обновляет `tasks.md`.

Поведение по умолчанию должно оставаться полным: без явных scope-флагов Draftspec проходит по всем незавершенным задачам в порядке task list.

Выборочное выполнение допустимо, когда пользователь явно сужает scope:

- `--phase <номер>` для одной implementation-фазы
- `--tasks <список-task-id>` для конкретных задач вроде `T1.1,T2.1`

`--phase` и `--tasks` не должны использоваться вместе в одном запуске.

Если выборочное выполнение перескакивает через незавершенную более раннюю работу, Draftspec должен предупредить о риске порядка, не расширяя scope молча.

Во время реализации Draftspec должен выдавать короткие runtime progress updates каждый раз, когда начинает или завершает фазу внутри текущего execution scope.

Такие phase-status сообщения должны следовать настроенному в проекте языку общения с агентом, а не автоматически сваливаться в английский.

### `verify`

Запускает легкую post-implementation проверку, чтобы подтвердить, что завершенная работа достаточно согласована с tasks и правилами проекта для безопасного следующего шага.

Полный verification report должен использовать стабильную структуру:

- `# Verify Report: <slug>`
- `## Scope`
- `## Verdict`
- `## Checks`
- `## Errors`
- `## Warnings`
- `## Questions`
- `## Not Verified`
- `## Next Step`

Рекомендуемые детали отчета:

- `## Scope` должен фиксировать реальный verification mode, например `default` или `deep`
- `## Scope` должен перечислять конкретные surfaces, которые реально проверялись
- `## Verdict` должен включать `archive_readiness`
- `## Verdict` должен включать однострочное summary, объясняющее, почему verdict обоснован
- `## Checks` должен включать `task_state`
- `## Checks` должен включать `acceptance_evidence` для тех `AC-*`, которые действительно подтверждены
- `## Checks` должен включать `implementation_alignment`, привязанный к конкретной проверенной surface
- `## Not Verified` должен перечислять material claims или surfaces, которые сознательно не проверялись

`Verdict` должен быть одним из значений:

- `pass`
- `concerns`
- `blocked`

Рекомендуемая семантика:

- `pass`: блокирующих проблем нет; остаются только незначительные предупреждения или их нет совсем
- `concerns`: по workflow можно двигаться дальше, но warnings или открытые вопросы желательно закрыть в ближайшее время
- `blocked`: архивирование или заявление о завершенности иначе опирались бы на противоречивое состояние реализации или незавершенную обязательную работу

Если evidence частичны, но явного противоречия не найдено, предпочитайте `concerns`, а не `pass`.

Если verification report нужно сохранить на диск, Draftspec должен по умолчанию использовать `.draftspec/plans/<slug>/verify.md`.

Используйте `.draftspec/templates/verify-report.md` как канонический шаблон, если отчет записывается в файл.

Сохраненные verify reports должны начинаться с того же machine-readable metadata block, что и inspect reports: `report_type`, `slug`, `status`, `docs_language` и `generated_at`.

Если доступен `.draftspec/scripts/check-verify-ready.sh <slug>`, Draftspec должен предпочитать его как дешевую readiness-проверку перед более глубокой verification.

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
- завершенные feature packages уходят в archive, не раздувая активное рабочее пространство
