# Prompt проверки реализации Draftspec

Вы проверяете один feature package после выполнения задач.

## Goal

Подтвердите, что реализованная работа достаточно согласована с задачами и правилами проекта, чтобы безопасно двигаться дальше.

## Flags

`--deep`: режим полной валидации реализации — читает все plan artifacts и проверяет реальный код для каждой завершённой задачи и acceptance criterion, а не только структурные проверки. Выдаёт comprehensive отчёт с per-AC evidence. Без этого флага verification остаётся структурным и cheap по умолчанию.

## Phase Contract

Inputs: `.draftspec/constitution.md`, `.draftspec/plans/<slug>/tasks.md`; spec, plan, код — только для подтверждения конкретных выводов (или все артефакты в режиме `--deep`).
Outputs: отчет с verdict (`pass`, `concerns` или `blocked`) в чате; сохраняется в `.draftspec/plans/<slug>/verify.md` по запросу.
Stop if: slug неоднозначен, tasks.md отсутствует, или verdict потребовал бы выдумывать факты о реализации.

## Load First

Всегда сначала прочитайте:

- `.draftspec/constitution.summary.md` если присутствует; иначе `.draftspec/constitution.md`
- `.draftspec/plans/<slug>/tasks.md`

## Load If Present

Читайте когда конкретная проверка ссылается на содержимое этих файлов (например, задача заявляет выполнение `AC-*`, или `DEC-*` ограничивает форму реализации):

- `.draftspec/specs/<slug>/summary.md` (или `spec.md`) — при проверке acceptance coverage или выравнивания task↔AC
- `.draftspec/plans/<slug>/plan.md` — когда задача ссылается на `DEC-*` или архитектурное решение, требующее подтверждения
- `.draftspec/plans/<slug>/data-model.md` — когда задача затрагивает persisted state или форму сущности
- `.draftspec/plans/<slug>/contracts/` — когда задача затрагивает API или event boundaries
- `.draftspec/plans/<slug>/research.md` — только если проверка зависит от зафиксированного trade-off или finding по внешней зависимости
- code files — только конкретные файлы из `Touches:` задачи, нужные для подтверждения реализации

## Do Not Read By Default

- нерелевантные области кода
- широкую историю репозитория
- архивы, если текущая проверка явно от них не зависит

## Stop Conditions

Остановитесь и задайте уточняющий вопрос только если:

- slug неоднозначен
- файл задач отсутствует
- без этого пришлось бы выдумывать факты о реализации
- запрошенный вывод потребовал бы широкого обхода репозитория вместо сфокусированных evidence по этой feature package
- implementation claim нельзя подтвердить по текущим tasks, planning artifacts и точечной проверке кода

## Rules

- Начинайте с `tasks.md` как verification entrypoint.
- Если доступен `/.draftspec/scripts/check-verify-ready.*`, предпочитайте его как cheap first pass перед чтением более глубоких артефактов.
- Используйте `/.draftspec/scripts/verify-task-state.*` только как fallback, когда phase readiness wrapper недоступен.
- Предпочитайте вывод helper scripts чтению их исходников.
- Не читайте `/.draftspec/scripts/*` по умолчанию, если только не отлаживаете сам script, не работаете над самим Draftspec или пользователь явно не просит проанализировать script logic.
- Предпочитайте подтверждение конкретных implementation claims вместо широкого субъективного review.
- Относитесь к verify как к журналу evidence, а не к ритуалу успокоения.
- Проверяйте, что завершенные задачи согласованы с текущим состоянием feature package.
- Проверяйте, что незавершенные задачи не противоречат заявлению о полной готовности фичи.
- Проверяйте согласованность acceptance-to-task coverage, если в `tasks.md` есть секция `Acceptance Coverage`.
- Если `tasks.md` использует task IDs вроде `T1.1`, ссылайтесь на них напрямую в checks, findings и выводах.
- Если evidence частичны, но явного противоречия нет, предпочитайте `concerns`, а не `pass`.
- Держите verification структурным и cheap-by-default.
- Когда `--deep` присутствует в `$ARGUMENTS`, переключитесь в режим полной валидации:
  - Читайте все plan artifacts (`plan.md`, `data-model.md`, `contracts/`, `research.md`).
  - Для каждой завершённой задачи читайте реальные файлы реализации из `Touches:` и подтверждайте, что работа соответствует описанию задачи.
  - Для каждого `AC-*` проследите по коду, что acceptance criterion выполнен с конкретными evidence.
  - Секция `## Scope` должна указывать `mode: deep` и перечислять все проверенные surfaces.
  - Секция `## Not Verified` должна быть минимальной или `none` — deep режим предполагает тщательность.
- Без `--deep` углубляйтесь в более широкий implementation review только если конкретное противоречие нельзя разрешить по `tasks`, plan artifacts и сфокусированным evidence.
- Используйте простой verdict: `pass`, `concerns` или `blocked`.
- Используйте `pass`, если блокирующих проблем нет и остаются только незначительные предупреждения или их нет совсем.
- Используйте `concerns`, если по workflow можно двигаться дальше, но warnings или открытые вопросы желательно закрыть в ближайшее время.
- Используйте `blocked`, если отсутствие завершенных задач или противоречивое состояние реализации делают архивирование или заявление о завершенности небезопасным.
- Не используйте `pass`, если состояние завершенных задач не подтверждено, если остается blocking contradiction или если acceptance / implementation claims не подкреплены реально проверенными evidence.
- Если записываете результат в файл, держите его на настроенном языке документации проекта.
- Используйте `.draftspec/templates/verify-report.md` как канонический шаблон, если отчет записывается в файл.
- Если отчет сохраняется в файл, добавляйте сверху machine-readable metadata block с полями `report_type`, `slug`, `status`, `docs_language` и `generated_at`.
- Используйте такую структуру отчета:
  - YAML-подобный metadata block в начале
  - `# Verify Report: <slug>`
  - `## Scope`
  - `## Verdict`
  - `## Checks`
  - `## Errors`
  - `## Warnings`
  - `## Questions`
  - `## Not Verified`
  - `## Next Step`
- В `## Scope` фиксируйте реальный verification mode и те surfaces, которые реально проверяли.
- В `## Verdict` добавляйте `archive_readiness` и однострочное summary, объясняющее, почему verdict обоснован.
- В `## Checks` явно отражайте:
  - `task_state` с completed/open counts
  - `acceptance_evidence` для тех `AC-*`, которые действительно подтвердили
  - `implementation_alignment` с указанием конкретной проверенной surface
- В `## Not Verified` перечисляйте material claims или surfaces, которые сознательно не проверяли. Используйте `none` только если в выбранном verification scope не осталось материальных gaps.
- Держите claims ограниченными реальным scope. Если вы проверили только task state и один endpoint или file path, так и напишите, а не намекайте на полный feature validation.
- Если verification обнаруживает workflow-gap, возвращайте фичу на самую узкую предыдущую фазу, которая честно может это исправить:
  - `implement` для отсутствующей или противоречивой реализации
  - `tasks` для неполной, вводящей в заблуждение или отсутствующей декомпозиции
  - `plan`, когда intent реализации нельзя честно оценить из-за недостаточно конкретного дизайна
- Для `pass` указывайте точную archive-команду.
- Для `concerns` явно говорите, можно ли двигаться дальше; если нельзя, используйте явную return-команду для более ранней фазы.
- Для `blocked` не подсказывайте archive; завершайте сводку строкой `Return to: /draftspec.<phase> <slug>` для самой узкой честной recovery-фазы.

## Output expectations

- Выведите отчет в разговор, если пользователь не просит сохранить его в файл; если нужно сохранить без явного пути — `.draftspec/plans/<slug>/verify.md`
- Кратко суммируйте verdict, выполненные проверки, оставшиеся concerns и можно ли безопасно архивировать фичу
- Завершайте разговор summary block: `Slug`, `Status`, `Artifacts`, `Blockers` и `Next command` / `Return to`
- Если фичу можно архивировать: `Следующая команда: /draftspec.archive <slug>`; при возврате на раннюю фазу — называйте её явно со slash-командой

## Self-Check

- Каждый вывод verdict подкреплен реально проверенными evidence, а не только состоянием чекбоксов?
- Секция `Not Verified` честно отражает всё, что я не проверял?
- Следующая команда или return-фаза соответствует verdict?
