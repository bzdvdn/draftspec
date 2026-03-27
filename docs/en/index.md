# draftspec Docs

`draftspec` is a lightweight project context kit for development agents and humans.

This documentation is organized into a few practical guides:

- [Overview](overview.md)
- [CLI Reference](cli.md)
- [Workflow Model](workflow.md)
- [Architecture](architecture.md)
- [Agents](agents.md)
- [Language and Configuration](language-and-config.md)
- [Self-Hosting and Development](self-hosting.md)
- [Examples](examples.md)
- [FAQ](faq.md)
- [Glossary](glossary.md)

## Quick Start

```bash
draftspec init my-project --lang en --agents claude --agents codex
```

This creates:

- `.draftspec/` workspace files
- project-local agent command files when `--agents` is used
- `AGENTS.md` guidance linked to Draftspec memory and templates

For a concise product summary, see the root [README](../README.md) and [MVP](../MVP.md).
