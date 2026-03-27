# Prompt анализа Draftspec

Вы анализируете один feature package на согласованность и качество.

## Goal

Сформируйте сфокусированный inspection report для одной фичи, не расширяя scope.

## Load First

Всегда сначала прочитайте:

- `.draftspec/constitution.md`
- `.draftspec/memory.md`
- `.draftspec/specs/<slug>.md`

## Load Only If Needed

Читайте это только если файл существует и реально влияет на анализ:

- `.draftspec/plans/<slug>/plan.md`
- `.draftspec/plans/<slug>/tasks.md`
- `.draftspec/plans/<slug>/data-model.md`
- `.draftspec/plans/<slug>/contracts/`
- `.draftspec/plans/<slug>/research.md`

## Do Not Read By Default

- нерелевантные спецификации
- нерелевантные plan packages
- широкую историю репозитория
- implementation-файлы, если они не нужны для подтверждения конкретного вывода по согласованности

## Stop Conditions

Остановитесь и задайте минимальный уточняющий вопрос только если:

- целевой slug неоднозначен
- сама спецификация отсутствует
- без этого анализу пришлось бы выдумывать продуктовый intent

## Rules

- Сначала проверяйте соответствие конституции.
- Проверяйте полноту и ясность спецификации.
- Если plan artifacts уже существуют, проверяйте согласованность между spec, plan, data model, contracts и tasks.
- Если записываете отчет в файл, держите его на настроенном языке документации проекта.
- Предпочитайте конкретные findings вместо общих советов.
- Используйте в отчете секции:
  - `Errors`
  - `Warnings`
  - `Questions`
  - `Suggestions`
- Предлагайте сценарии Given/When/Then только когда они реально усиливают слабые критерии приемки.

## Output expectations

- Запишите или patch-обновите inspection report для фичи
- Кратко суммируйте ошибки, предупреждения, открытые вопросы и предложения
