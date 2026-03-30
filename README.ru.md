# draftspec

Русская версия: [English README](README.md)

`draftspec` — это легкий файловый каркас проектного контекста для людей и агентов разработки.

Он хранит намерение проекта, спецификации, плановые артефакты и декомпозицию задач в простых файлах, не превращаясь в жесткий механизм управления процессом.

Первый релиз намеренно оптимизирован под низкие накладные расходы и реальную работу: узкий контекст по умолчанию, минимально обязательные артефакты, строгую дисциплину workflow без тяжеловесной оркестрации и branch-first модель collaboration, которая одинаково хорошо подходит и для одиночной, и для командной работы.

## Позиционирование

Draftspec — это легкий SDD-kit для реальных кодовых баз.

- строже, чем OpenSpec, в дисциплине фаз и согласованности артефактов
- легче, чем Spec Kit, по контексту по умолчанию, ширине workflow и артефактным накладным расходам
- оптимизирован для agent-first workflow с узкой загрузкой контекста
- сохраняет strictness через templates, entrypoints и readiness checks, а не через разрастание процесса

Коротко: Draftspec старается быть настолько строгим, насколько это практически полезно, оставаясь при этом достаточно легким для повседневной работы.

## Draftspec vs OpenSpec vs Spec Kit

| Dimension | Draftspec | OpenSpec | Spec Kit |
| --- | --- | --- | --- |
| Workflow style | Строгая цепочка фаз с узким контекстом | Более гибкий workflow вокруг артефактов | Более подробный многошаговый SDD-workflow |
| Default context size | Самый маленький по умолчанию | Средний | Самый большой |
| Artifact overhead | Низкий | Средний | Высокий |
| Phase discipline | Высокая | Средняя | Максимальная |
| Brownfield ergonomics | Высокая | Высокая | Средняя |
| Team collaboration model | Branch-first, feature-local artifacts | Модель вокруг change-folders | Тяжелее по веткам и workflow |
| Shared mutable state | Избегается по дизайну | Низкий | Зависит от конфигурации |
| Best fit | Легкий строгий SDD для реальных кодовых баз | Гибкий SDD-lite для быстрой итерации | Полноценный строгий SDD-toolkit |

Коротко: Draftspec занимает место между OpenSpec и Spec Kit: строже OpenSpec, легче Spec Kit и лучше приспособлен для branch-based collaboration с минимальным контекстом по умолчанию.

## Где Draftspec Сильнее Всего

- Узкий контекст по умолчанию. Каждая фаза должна загружать только минимально полезный объем.
- Чтение кода должно оставаться фазово-локальным и точечным: достаточно, чтобы убрать догадки, но не настолько широко, чтобы пересобирать контекст всего репозитория.
- Строгая цепочка workflow. Конституция, spec, inspect, plan, tasks и implement остаются согласованными.
- Легкая трассировка. Стабильные ID и дешевые readiness checks уменьшают перегрузку prompt-контекста.
- Удобство для brownfield-репозиториев. Draftspec хорошо работает в существующих кодовых базах без навязывания тяжелого процессного слоя.
- Branch-first collaboration. Активное состояние фичи остается локальным для самой фичи, а не размазывается по общей изменяемой памяти.
- `inspect` обязателен перед `plan`. Для каждой фичи должен сохраняться inspect-отчет, который подтверждает, что spec готова к планированию.

OpenSpec по дизайну более гибкий и хорошо подходит командам, которым нужен более свободный workflow вокруг артефактов.

Spec Kit дает более широкую и подробную workflow surface, но обычно ценой большего числа артефактов, более широкого контекста и более тяжелого процесса.

Draftspec оптимизируется под discipline per token: сильные границы workflow, минимальный контекст по умолчанию и достаточную структуру для согласованной работы агентов без превращения процесса в тяжеловесную систему.

## Документация

Расширенная документация находится в [`docs/`](docs/README.md):

- [English docs](docs/en/index.md)
- [Русская документация](docs/ru/index.md)

## Публичный CLI

```text
draftspec init [path]
draftspec refresh [path]
draftspec add-agent [path]
draftspec list-agents [path]
draftspec remove-agent [path]
draftspec cleanup-agents [path]
draftspec doctor [path]
draftspec doctor [path] --json
draftspec feature <slug> [path]
draftspec feature repair <slug> [path]
draftspec features [path]
draftspec migrate [path]
draftspec list-specs [path]
draftspec show-spec <name> [path]
```

## Workflow

```text
constitution -> spec -> inspect -> plan -> tasks -> implement -> verify -> archive
```

## Ключевые Идеи

- Конституция — главный документ проекта.
- Plan package хранит вместе `plan.md`, `tasks.md`, `data-model.md`, `contracts/` и optional `research.md`.
- `data-model.md` и `contracts/` должны оставаться компактными, но структурированными: сущности должны описывать поля, инварианты и жизненный цикл, а контракты — входы и выходы на границах системы, ошибки и предположения о доставке.
- Specs используют канонические маркеры `Given / When / Then` независимо от языка документации.
- Draftspec предпочитает стабильные ID и явные ссылки вместо повторяющихся narrative summaries: `RQ-*` для требований, `AC-*` для критериев приемки, `DEC-*` для решений плана и phase-scoped `T*` для task IDs.
- Agent workflows должны читать только минимально нужный контекст.
- Strictness обеспечивается phase entrypoints, templates, стабильной структурой артефактов и readiness checks, а не большими prompt contexts.
- Agent-facing `/draftspec.spec` работает branch-first: от `feature/<slug>`, поддерживает `--name` с optional `--slug` / `--branch` и сохраняет приоритет явных `name:` / `slug:` в prompt-файлах.
- `draftspec init` требует явный `--shell` и генерирует только одно семейство scripts: `sh` или `powershell`.
- Сгенерированный workspace включает `.draftspec/scripts/run-draftspec.*` как стабильный CLI launcher для агентов; он сначала использует `DRAFTSPEC_BIN`, а потом пытается вызвать `draftspec` из `PATH`.
- `draftspec feature repair` и `draftspec migrate` дают безопасную каноникализацию legacy-артефактов, например старых путей к inspect reports.
- Генерируемые docs и prompts поддерживают английский и русский.

## Быстрый Пример

```bash
draftspec init my-project --lang ru --shell sh --agents claude --agents codex
draftspec refresh my-project --shell powershell --dry-run
draftspec doctor my-project
draftspec doctor my-project --json
```

Для более подробного знакомства смотри:

- [Обзор](docs/ru/overview.md)
- [CLI](docs/ru/cli.md)
- [Модель workflow](docs/ru/workflow.md)
- [Примеры](docs/ru/examples.md)
- [Roadmap](docs/ru/roadmap.md)

## Разработка

```bash
go test ./...
go build -o bin/draftspec ./src/cmd/draftspec
```

Репозиторий содержит unit tests для config, project lifecycle, doctor checks, specs, templates, agents и CLI-level behavior.

## Лицензия

Проект распространяется по лицензии [MIT](LICENSE).
