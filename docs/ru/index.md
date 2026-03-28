# Документация draftspec

`draftspec` — это легкий файловый каркас проектного контекста для людей и агентной разработки.

Документация разбита на несколько практических разделов:

- [Обзор](overview.md)
- [CLI](cli.md)
- [Модель workflow](workflow.md)
- [Архитектура](architecture.md)
- [Агенты](agents.md)
- [Языки и конфигурация](language-and-config.md)
- [Self-hosting и разработка](self-hosting.md)
- [Примеры](examples.md)
- [FAQ](faq.md)
- [Glossary](glossary.md)
- [Roadmap](roadmap.md)

## Быстрый старт

```bash
draftspec init my-project --lang ru --agents claude --agents codex
```

Это создаст:

- рабочее пространство `.draftspec/`
- project-local файлы команд для агентов, если указан `--agents`
- `AGENTS.md` с привязкой к памяти проекта и шаблонам Draftspec

Для краткого обзора продукта смотри корневые [README](../README.md) и [MVP](../MVP.md).
