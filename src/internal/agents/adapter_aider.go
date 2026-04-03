package agents

import (
	"fmt"
	"strings"
)

type aiderAdapter struct{}

func (aiderAdapter) Target() string { return "aider" }

func (aiderAdapter) Render(commands []CommandDefinition, language string) ([]File, error) {
	return []File{{
		Path:    ".aider/CONVENTIONS.md",
		Content: renderAiderCommands(commands, language),
		Mode:    0o644,
	}}, nil
}

func (aiderAdapter) Paths(commands []CommandDefinition, language string) ([]string, error) {
	return []string{".aider/CONVENTIONS.md"}, nil
}

func renderAiderCommands(commands []CommandDefinition, language string) string {
	lang := normalizeLanguage(language)
	if lang == "ru" {
		var sections []string
		sections = append(sections, "# Draftspec Conventions")
		sections = append(sections, "")
		sections = append(sections, "Используйте `.draftspec/` как основной источник проектного контекста. Следуйте соответствующим prompt-файлам в `.draftspec/templates/prompts/`.")
		sections = append(sections, "")
		sections = append(sections, "Загружайте этот файл через `--read .aider/CONVENTIONS.md` или добавьте `read: .aider/CONVENTIONS.md` в `.aider.conf.yml`.")
		for _, cmd := range commands {
			sections = append(sections, "")
			sections = append(sections, fmt.Sprintf("## /draftspec.%s", cmd.Name))
			sections = append(sections, fmt.Sprintf("- Основной prompt: %s", cmd.PromptPath))
			sections = append(sections, fmt.Sprintf("- %s", commandHint(cmd.Name, lang)))
			sections = append(sections, "- Используйте только минимально нужный контекст репозитория")
			sections = append(sections, "- Если доступны связанные scripts, сначала запускайте их и опирайтесь на их вывод")
			sections = append(sections, "- Не читайте исходники scripts по умолчанию")
			if len(cmd.Extras) > 0 {
				sections = append(sections, "- Связанные scripts:")
				sections = append(sections, bulletList(cmd.Extras))
			}
		}
		return strings.Join(sections, "\n") + "\n"
	}

	var sections []string
	sections = append(sections, "# Draftspec Conventions")
	sections = append(sections, "")
	sections = append(sections, "Use `.draftspec/` as the primary source of project context. Follow the matching prompt files under `.draftspec/templates/prompts/`.")
	sections = append(sections, "")
	sections = append(sections, "Load this file via `--read .aider/CONVENTIONS.md` or add `read: .aider/CONVENTIONS.md` to `.aider.conf.yml`.")
	for _, cmd := range commands {
		sections = append(sections, "")
		sections = append(sections, fmt.Sprintf("## /draftspec.%s", cmd.Name))
		sections = append(sections, fmt.Sprintf("- Primary prompt: %s", cmd.PromptPath))
		sections = append(sections, fmt.Sprintf("- %s", commandHint(cmd.Name, lang)))
		sections = append(sections, "- Use only the minimum repository context required")
		sections = append(sections, "- When related scripts are available, run them first and rely on their output")
		sections = append(sections, "- Do not read script source by default")
		if len(cmd.Extras) > 0 {
			sections = append(sections, "- Related scripts:")
			sections = append(sections, bulletList(cmd.Extras))
		}
	}
	return strings.Join(sections, "\n") + "\n"
}
