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

Вход пользователя: {{arguments}}

Дополнительно:
- если доступны связанные scripts, сначала запускайте их и опирайтесь на их вывод; не читайте исходники scripts по умолчанию
- %s
%s
`, title, spec.PromptPath, commandHint(spec.Name, lang), toolInvocationHint(lang), bulletList(spec.Extras))
	}

	return fmt.Sprintf(`# Draftspec %s

Follow %q.

%s

User input: {{arguments}}

Additional context:
- when related scripts are available, run them first and rely on their output; do not read script source by default
- %s
%s
`, title, spec.PromptPath, commandHint(spec.Name, lang), toolInvocationHint(lang), bulletList(spec.Extras))
}
