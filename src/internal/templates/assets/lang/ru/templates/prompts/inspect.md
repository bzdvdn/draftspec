# Prompt проверки Draftspec

Вы проверяете один пакет фичи на согласованность и качество.

## Goal

Сформируйте сфокусированный отчет проверки для одной фичи, не расширяя scope.

## Phase Contract

Inputs: `.draftspec/constitution.md`, `.draftspec/specs/<slug>.md`; опционально `plan.md`, `tasks.md`, если они существуют.
Outputs: `.draftspec/specs/<slug>.inspect.md` с verdict `pass`, `concerns` или `blocked`.
Stop if: slug неоднозначен, spec отсутствует, или отчет потребовал бы выдумывать продуктовый intent.

## Load First

Всегда сначала прочитайте:

- `.draftspec/constitution.md`
- `.draftspec/specs/<slug>.md`

## Load Only If Needed

Читайте это только если файл существует и реально влияет на проверку:

- `.draftspec/plans/<slug>/plan.md`
- `.draftspec/plans/<slug>/tasks.md`

## Do Not Read By Default

- `.draftspec/plans/<slug>/data-model.md`
- `.draftspec/plans/<slug>/contracts/`
- `.draftspec/plans/<slug>/research.md`
- нерелевантные спецификации
- нерелевантные plan packages
- широкую историю репозитория
- implementation-файлы, если они не нужны для подтверждения конкретного вывода по согласованности

## Stop Conditions

Остановитесь и задайте минимальный уточняющий вопрос только если:

- целевой slug неоднозначен
- сама спецификация отсутствует
- без этого пришлось бы выдумывать продуктовый intent

## Rules

- Сначала проверяйте соответствие конституции.
- Если доступен `/.draftspec/scripts/check-inspect-ready.*`, предпочитайте его как cheap first pass перед углублением в артефакты.
- Используйте `/.draftspec/scripts/inspect-spec.*` только как fallback, когда phase readiness wrapper недоступен.
- Предпочитайте вывод helper scripts чтению их исходников.
- Не читайте `/.draftspec/scripts/*` по умолчанию, если только не отлаживаете сам script, не работаете над самим Draftspec или пользователь явно не просит проанализировать script logic.
- Проверяйте полноту и ясность спецификации.
- Проверяйте `constitution <-> spec`: спецификация не должна противоречить явным ограничениям конституции, workflow-правилам или language policy.
- Считайте technology names, framework choices, library lists или version pins в спецификации `Warning`, если они явно не выглядят как user requirement, repository constraint или внешний compatibility contract.
- Каждый критерий приемки в спецификации ДОЛЖЕН иметь явный формат Given/When/Then. Маркеры `Given`, `When`, `Then` остаются каноническими независимо от языка документации. Отсутствие G/W/T — `Error`, а не `Suggestion`.
- Если `tasks.md` существует, проверяйте, что каждый критерий приемки из spec покрыт хотя бы одной задачей. Непокрытый критерий — `Error`.
- Если `tasks.md` использует task IDs вроде `T1.1`, предпочитайте traceability-формулировки с прямыми ссылками на эти task IDs.
- Предпочитайте самый дешевый inspection scope: `constitution.md` и `spec.md`, затем `plan.md`, затем `tasks.md`, и только после этого более глубокие plan artifacts, если они нужны для подтверждения конкретного вывода.
- Если `plan.md` отсутствует, не расширяйте проверку на optional plan artifacts или implementation code.
- Если артефакты планирования уже существуют, проверяйте согласованность между spec, plan, data model, contracts и tasks.
- Когда существует `plan.md`, сначала проверяйте `spec <-> plan` consistency, не читая более глубокие plan artifacts без необходимости.
- Проверяйте `spec <-> plan`: план должен сохранять цель фичи, отражать major acceptance-critical behavior и не добавлять необоснованные новые workstreams.
- Если `tasks.md` существует, проверяйте `plan <-> tasks`: фазы и task IDs должны отражать intent плана без явных пропусков по acceptance-critical behavior.
- Считайте `spec.md` и `plan.md` обязательными входами для дешевой проверки согласованности плана.
- Читайте `data-model.md` или `contracts/` только если `plan.md` явно на них опирается или без них нельзя подтвердить конкретный consistency claim.
- Проверяйте `Goal Alignment`: plan не должен менять основную цель фичи, выраженную в spec.
- Проверяйте `Scope Expansion`: plan не должен вводить крупные новые workstreams, компоненты или integration surfaces, которых нет в spec.
- Проверяйте `Acceptance Coverage at Plan Level`: major acceptance-critical behavior из spec должно быть отражено в намерении плана, даже до появления tasks.
- Проверяйте `Constitution Consistency`: plan не должен нарушать правила конституции или архитектурные ограничения.
- Проверяйте `Artifact Justification`: если plan вводит `data-model.md` или `contracts/`, необходимость этих артефактов должна быть оправдана spec.
- Не превращайте это в широкий design review. Предпочитайте ловить явный drift, а не оценивать качество архитектуры целиком.
- Если записываете отчет в файл, держите его на настроенном языке документации проекта.
- Предпочитайте конкретные находки вместо общих советов.
- По умолчанию делайте compact report в разговоре: всегда включайте `Verdict`, включайте `Errors`, `Warnings` и `Next Step`, если они не пусты, а `Questions`, `Suggestions` и `Traceability` — только когда они действительно добавляют сигнал.
- Полный sectioned report используйте только если пользователь явно просит полный отчет или если отчет сохраняется в файл.
- Если отчет сохраняется в файл, добавляйте сверху machine-readable metadata block с полями `report_type`, `slug`, `status`, `docs_language` и `generated_at`.
- Используйте такую структуру отчета:
  - YAML-подобный metadata block в начале
  - `# Inspect Report: <slug>`
  - `## Scope`
  - `## Verdict`
  - `## Errors`
  - `## Warnings`
  - `## Questions`
  - `## Suggestions`
  - `## Traceability`
  - `## Next Step`
- В секции `## Verdict` ДОЛЖНО использоваться одно из значений: `pass`, `concerns`, `blocked`.
- Используйте `pass`, когда ошибок нет и остаются только незначительные предупреждения или предупреждений нет совсем.
- Используйте `concerns`, когда по фиче можно двигаться дальше, но warnings, пробелы traceability или открытые вопросы желательно закрыть в ближайшее время.
- Используйте `blocked`, если конфликт с конституцией, отсутствие продуктового intent, отсутствие Given/When/Then в acceptance criteria, непокрытые acceptance criteria или крупные противоречия между `spec` и `plan` не позволяют безопасно продолжать следующую фазу.
- `## Traceability` должна кратко показывать, как acceptance criteria связаны с задачами, если `tasks.md` уже существует.
- Предпочитайте traceability-строки со стабильными acceptance IDs и task IDs, например `AC-001 -> T1.1, T2.1`.
- `## Next Step` должна явно говорить, можно ли безопасно продолжать к `plan`, `tasks`, или сначала нужно уточнение.
- Для `pass` указывайте точную следующую slash-команду.
- Для `concerns` явно говорите, можно ли двигаться дальше; если можно, указывайте точную следующую slash-команду.
- Для `blocked` не подсказывайте следующую фазу; вместо этого указывайте, какой refinement нужен сначала.

## Output expectations

- По умолчанию сохраняйте отчет в `.draftspec/specs/<slug>.inspect.md`, если пользователь явно не просит другой путь.
- Если пользователь задал явный путь, используйте именно его.
- Также кратко суммируйте verdict в разговоре.
- В разговоре по умолчанию предпочитайте compact report только с непустыми секциями.
- Кратко суммируйте ошибки, предупреждения, открытые вопросы, предложения и итоговый verdict.
- Завершайте разговор коротким стабильным summary block с полями `Slug`, `Status`, `Artifacts`, `Blockers` и `Next command` только если такой handoff действительно безопасен
- Когда по фиче можно безопасно продолжать работу, в `## Next Step` и в разговорной сводке указывайте точную slash-команду следующей фазы, например `/draftspec.plan <slug>` или `/draftspec.tasks <slug>`.
- Если сначала нужен refinement, говорите об этом прямо вместо подсказки следующей фазы.

## Self-Check

- Я загрузил только артефакты, нужные для этого slug?
- Я проверил каждый AC на наличие формата Given/When/Then?
- Verdict (`pass`, `concerns`, `blocked`) опирается на конкретные находки, а не общие впечатления?
- Если `tasks.md` существует, я убедился, что каждый AC покрыт хотя бы одной задачей?
- Я избежал превращения проверки в широкий design review?
- Следующая команда в Next Step соответствует verdict?
