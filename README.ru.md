# draftspec

Русская версия: [English README](README.md)

`draftspec` — это легкий файловый каркас проектного контекста для людей и development agents.

Он хранит намерение проекта, спецификации, плановые артефакты и декомпозицию задач в простых файлах, не превращаясь в жесткий process engine.

Первый релиз намеренно оптимизирован под low overhead и реальную работу: узкий default context, минимально обязательные артефакты, строгую дисциплину workflow без heavyweight orchestration и branch-based collaboration, которая одинаково хорошо подходит и для solo, и для team workflow.

## Позиционирование

Draftspec — это lean SDD kit для реальных кодовых баз.

- строже, чем OpenSpec, в phase discipline и согласованности артефактов
- легче, чем Spec Kit, по default context, workflow surface и artifact overhead
- оптимизирован для agent-first workflow с узкой загрузкой контекста
- сохраняет strictness через templates, entrypoints и readiness checks, а не через разрастание процесса

Коротко: Draftspec старается быть настолько строгим, насколько это practically useful, оставаясь при этом достаточно легким для повседневной работы.

## Draftspec vs OpenSpec vs Spec Kit

| Dimension | Draftspec | OpenSpec | Spec Kit |
| --- | --- | --- | --- |
| Workflow style | Строгая цепочка фаз с узким контекстом | Более fluid artifact-guided workflow | Более thorough multi-step SDD workflow |
| Default context size | Самый маленький по умолчанию | Средний | Самый большой |
| Artifact overhead | Низкий | Средний | Высокий |
| Phase discipline | Высокая | Средняя | Максимальная |
| Brownfield ergonomics | Высокая | Высокая | Средняя |
| Team collaboration model | Branch-first, feature-local artifacts | Change-folder oriented | Branch and workflow heavy |
| Shared mutable state | Избегается по дизайну | Низкий | Зависит от setup |
| Best fit | Lean strict SDD для реальных кодовых баз | Flexible SDD-lite для быстрой итерации | Полноценный rigorous SDD toolkit |

Коротко: Draftspec занимает место между OpenSpec и Spec Kit: строже OpenSpec, легче Spec Kit и лучше приспособлен для branch-based collaboration с минимальным default context.

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
- Specs используют канонические маркеры `Given / When / Then` независимо от языка документации.
- Agent workflows должны читать только минимально нужный контекст.
- Strictness обеспечивается phase entrypoints, templates и readiness checks, а не большими prompt contexts.
- Agent-facing `/draftspec.spec` работает branch-first: от `feature/<slug>`, поддерживает `--name` с optional `--slug` / `--branch` и сохраняет приоритет явных `name:` / `slug:` в prompt-файлах.
- `draftspec init` требует явный `--shell` и генерирует только одно семейство scripts: `sh` или `powershell`.
- Генерируемые docs и prompts поддерживают английский и русский.

## Быстрый Пример

```bash
draftspec init my-project --lang ru --shell sh --agents claude --agents codex
draftspec refresh my-project --shell powershell --dry-run
draftspec doctor my-project
draftspec doctor my-project --json
```

Для более подробного входа смотри:

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
