# Prompt архивации Draftspec

Вы архивируете один feature package.

## Goal

Создайте устойчивый архивный снимок одной фичи и обновите память проекта короткой записью об архиве.

## Load First

Всегда сначала прочитайте:

- `.draftspec/memory.md`
- `.draftspec/specs/<slug>.md`

## Load Only If Needed

Читайте это, если файлы существуют:

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
- В MVP используйте copy-based archive; не удаляйте active files без явного указания.
- Записывайте `summary.md` внутри директории архива.
- Если plan artifacts существуют, архивируйте их вместе со spec.
- Если `research.md` не существует, не выдумывайте его.
- Обновляйте `memory.md`, добавляя короткую запись в `Archived Specs` со slug, status, date и reason.
- Используйте один из этих статусов:
  - `completed`
  - `superseded`
  - `abandoned`
  - `rejected`
  - `deferred`

## Output expectations

- Создайте архивный снимок
- Запишите или patch-обновите `summary.md`
- Обновите `memory.md`
- Кратко суммируйте архивированные файлы, статус и причину
