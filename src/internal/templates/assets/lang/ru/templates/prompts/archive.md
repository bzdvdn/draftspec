# Prompt архивации Draftspec

Вы архивируете один feature package.

## Goal

Создайте устойчивый архивный снимок одной фичи.

## Phase Contract

Inputs: смотрите Load First и Load Only If Needed.
Outputs: архивный снимок и summary.
Stop if: смотрите Stop Conditions.

## Load First

Всегда сначала прочитайте:

- `.draftspec/specs/<slug>.md`

## Load Only If Needed

Читайте plan artifacts только чтобы сформировать `summary.md`. Если summary достаточно написать по spec — не читайте их:

- `.draftspec/plans/<slug>/plan.md`
- `.draftspec/plans/<slug>/tasks.md`
- `.draftspec/plans/<slug>/data-model.md`
- `.draftspec/plans/<slug>/contracts/`
- `.draftspec/plans/<slug>/research.md`

## Do Not Read By Default

- нерелевантные спецификации
- нерелевантные archive entries
- нерелевантные файлы репозитория

## Stop Conditions

Остановитесь и задайте минимальный уточняющий вопрос, если:

- неясен статус архивации
- отсутствует причина архивации
- целевой slug неоднозначен

## Rules

- Архивируйте в `.draftspec/archive/<slug>/<YYYY-MM-DD>/`.
- В MVP используйте copy-based archive; не удаляйте active files (`specs/<slug>.md` и `plans/<slug>/`) без явного указания пользователя.
- Записывайте `summary.md` внутри директории архива.
- Держите `summary.md` компактным. Предпочитайте статус, причину, завершенный scope и заметные отклонения вместо длинного ретроспективного повествования.
- Если plan artifacts существуют, архивируйте их вместе со spec.
- Если `research.md` не существует, не выдумывайте его.
- Если status равен `completed` и `tasks.md` существует, используйте `/.draftspec/scripts/verify-task-state.* <slug>` перед архивацией. Не заявляйте completed-архив, если обязательные задачи еще открыты.
- Используйте один из этих статусов:
  - `completed`
  - `superseded`
  - `abandoned`
  - `rejected`
  - `deferred`

## Output expectations

- Создайте архивный снимок
- Запишите или patch-обновите `summary.md`
- Кратко суммируйте архивированные файлы, статус и причину
