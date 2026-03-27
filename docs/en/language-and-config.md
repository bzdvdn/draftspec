# Language and Configuration

## Supported Languages

Draftspec currently supports:

- `en`
- `ru`

## Language Settings

The generated config stores four language values:

```yaml
language:
  default: en
  docs: en
  agent: en
  comments: en
```

### `default`

Base language used when the more specific values are not overridden.

### `docs`

Controls generated project documentation such as:

- `constitution.md`
- `memory.md`
- spec templates
- plan templates
- tasks templates

### `agent`

Controls generated prompts and `AGENTS.md` guidance.

### `comments`

Records the preferred language for new or edited code comments during `implement`.

## Canonical BDD Markers

Even when `docs` is set to Russian, Draftspec keeps `Given / When / Then` as canonical acceptance markers. The surrounding text should still follow the configured documentation language.

## Config File

The main config lives in `.draftspec/draftspec.yaml` and stores:

- path layout
- language settings
- enabled agent targets
- template names
- script names

Use `draftspec doctor` if the config and workspace drift apart.
