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

Используйте это rule, когда запрос явно относится к фазе %q или команде /draftspec.%s.

Если доступны связанные scripts, сначала запускайте их и опирайтесь на их вывод. Не читайте исходники scripts по умолчанию.

Связанные scripts:
%s
`, spec.PromptPath, commandHint(spec.Name, lang), spec.Name, spec.Name, bulletList(spec.Extras))
	}

	return fmt.Sprintf(`---
trigger: manual
---

Follow %q.

%s

Use this rule when the request clearly maps to the %q phase or the /draftspec.%s command.

When related scripts are available, run them first and rely on their output. Do not read script source by default.

Related scripts:
%s
`, spec.PromptPath, commandHint(spec.Name, lang), spec.Name, spec.Name, bulletList(spec.Extras))
}
