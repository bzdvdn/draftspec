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
- `archive`

Каждый prompt должен:

- читать только минимально нужный контекст
- останавливаться при отсутствии prerequisites
- уважать configured documentation language и agent language
- считать конституцию документом с наивысшим приоритетом

## Команды обслуживания

Управлять agent targets лучше через публичный CLI:

```bash
draftspec add-agent my-project --agents claude --agents cursor
draftspec list-agents my-project
draftspec remove-agent my-project --agents cursor
draftspec cleanup-agents my-project
draftspec doctor my-project
```
