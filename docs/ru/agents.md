# Агенты

## Поддерживаемые agent targets

Draftspec умеет генерировать project-local command или prompt files для:

- `claude`
- `codex`
- `copilot`
- `cursor`
- `kilocode`
- `trae`
- `all`

## Куда пишутся файлы

- Claude: `.claude/commands/`
- Codex: `.codex/prompts/`
- Copilot: `.github/prompts/`
- Cursor: `.cursor/rules/`
- Kilo Code: `.kilocode/rules/`
- Trae: `.trae/project_rules.md`

Эти generated files являются тонкой оберткой над каноническими промтами в `.draftspec/templates/prompts/`.

## Агентская дисциплина

Agent-facing workflow в Draftspec:

- `constitution`
- `spec`
- `inspect`
- `plan`
- `tasks`
- `implement`
- `verify`
- `archive`

Каждый prompt должен:

- читать только минимально нужный контекст
- останавливаться при отсутствии prerequisites
- уважать configured documentation language и agent language
- считать конституцию документом с наивысшим приоритетом

`verify` специально сделан легким:

- он стартует от `tasks.md`
- он может использовать `.draftspec/scripts/verify-task-state.sh <slug>` как дешевый helper первого прохода
- он может использовать `.draftspec/scripts/verify-memory-sync.sh <slug>` для грубых сигналов sync между memory и tasks до чтения кода
- более глубокие артефакты читаются только когда нужно подтвердить конкретный вывод
- его задача — подтвердить готовность к архивированию или следующему refine-циклу, а не превращаться в тяжелый review engine

Для `archive` статус `completed` должен переиспользовать `verify-task-state.sh` перед созданием снимка.

## Команды обслуживания

Управлять agent targets лучше через публичный CLI:

```bash
draftspec add-agent my-project --agents claude --agents cursor
draftspec list-agents my-project
draftspec remove-agent my-project --agents cursor
draftspec cleanup-agents my-project
draftspec doctor my-project
```
