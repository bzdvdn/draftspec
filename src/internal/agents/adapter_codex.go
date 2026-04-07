package agents

import (
	"fmt"
	"path/filepath"
)

type codexAdapter struct{}

func (codexAdapter) Target() string { return "codex" }

func (codexAdapter) Render(commands []CommandDefinition, language string) ([]File, error) {
	lang := normalizeLanguage(language)
	files := make([]File, 0, len(commands))
	for _, command := range commands {
		files = append(files, File{
			Path:    filepath.ToSlash(filepath.Join(".codex", "prompts", fmt.Sprintf("draftspec.%s.md", command.Name))),
			Content: renderCodex(command, lang),
			Mode:    0o644,
		})
	}
	return files, nil
}

func (codexAdapter) Paths(commands []CommandDefinition, language string) ([]string, error) {
	files, err := codexAdapter{}.Render(commands, language)
	if err != nil {
		return nil, err
	}
	paths := make([]string, 0, len(files))
	for _, file := range files {
		paths = append(paths, file.Path)
	}
	return paths, nil
}

func renderCodex(spec CommandDefinition, lang string) string {
	title := titleCase(spec.Name)
	if lang == "ru" {
		return fmt.Sprintf(`# Draftspec %s

Следуйте файлу %q.

%s

%s

Вход пользователя: {{arguments}}

Дополнительно:
- %s
- %s
- %s
%s

%s
`, title, spec.PromptPath, commandHint(spec.Name, lang), workflowChainHint(lang), scriptExecutionHint(lang), toolInvocationHint(lang), helpDiscoveryHint(lang), scriptListBlock(spec.Extras, lang), antiPatternHint(lang))
	}

	return fmt.Sprintf(`# Draftspec %s

Follow %q.

%s

%s

User input: {{arguments}}

Additional context:
- %s
- %s
- %s
%s

%s
`, title, spec.PromptPath, commandHint(spec.Name, lang), workflowChainHint(lang), scriptExecutionHint(lang), toolInvocationHint(lang), helpDiscoveryHint(lang), scriptListBlock(spec.Extras, lang), antiPatternHint(lang))
}
