package agents

import (
	"fmt"
	"path/filepath"
)

type windsurfAdapter struct{}

func (windsurfAdapter) Target() string { return "windsurf" }

func (windsurfAdapter) Render(commands []CommandDefinition, language string) ([]File, error) {
	lang := normalizeLanguage(language)
	files := make([]File, 0, len(commands))
	for _, command := range commands {
		files = append(files, File{
			Path:    filepath.ToSlash(filepath.Join(".windsurf", "rules", fmt.Sprintf("draftspec-%s.md", command.Name))),
			Content: renderWindsurf(command, lang),
			Mode:    0o644,
		})
	}
	return files, nil
}

func (windsurfAdapter) Paths(commands []CommandDefinition, language string) ([]string, error) {
	files, err := windsurfAdapter{}.Render(commands, language)
	if err != nil {
		return nil, err
	}
	paths := make([]string, 0, len(files))
	for _, file := range files {
		paths = append(paths, file.Path)
	}
	return paths, nil
}

func renderWindsurf(spec CommandDefinition, lang string) string {
	if lang == "ru" {
		return fmt.Sprintf(`---
trigger: manual
---

Следуйте файлу %q.

%s

%s

Используйте это rule, когда запрос явно относится к фазе %q или команде /draftspec.%s.

%s

- %s
%s

%s
`, spec.PromptPath, commandHint(spec.Name, lang), workflowChainHint(lang), spec.Name, spec.Name, scriptExecutionHint(lang), helpDiscoveryHint(lang), scriptListBlock(spec.Extras, lang), antiPatternHint(lang))
	}

	return fmt.Sprintf(`---
trigger: manual
---

Follow %q.

%s

%s

Use this rule when the request clearly maps to the %q phase or the /draftspec.%s command.

%s

- %s
%s

%s
`, spec.PromptPath, commandHint(spec.Name, lang), workflowChainHint(lang), spec.Name, spec.Name, scriptExecutionHint(lang), helpDiscoveryHint(lang), scriptListBlock(spec.Extras, lang), antiPatternHint(lang))
}
