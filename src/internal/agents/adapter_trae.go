package agents

import (
	"fmt"
	"strings"
)

type traeAdapter struct{}

func (traeAdapter) Target() string { return "trae" }

func (traeAdapter) Render(commands []CommandDefinition, language string) ([]File, error) {
	return []File{{
		Path:    ".trae/project_rules.md",
		Content: renderTraeCommands(commands, language),
		Mode:    0o644,
	}}, nil
}

func (traeAdapter) Paths(commands []CommandDefinition, language string) ([]string, error) {
	return []string{".trae/project_rules.md"}, nil
}

func renderTrae(language, shell string) string {
	return renderTraeCommands(DefaultCommands(shell), language)
}

func renderTraeCommands(commands []CommandDefinition, language string) string {
	lang := normalizeLanguage(language)
	if lang == "ru" {
		var sections []string
		sections = append(sections, "# Draftspec Project Rules")
		sections = append(sections, "")
		sections = append(sections, "Используйте .draftspec как основной источник проектного контекста. Следуйте AGENTS.md и соответствующим prompt-файлам в .draftspec/templates/prompts/.")
		for _, spec := range commands {
			sections = append(sections, "")
			sections = append(sections, fmt.Sprintf("## /draftspec.%s", spec.Name))
			sections = append(sections, fmt.Sprintf("- Основной prompt: %s", spec.PromptPath))
			sections = append(sections, fmt.Sprintf("- %s", commandHint(spec.Name, lang)))
			sections = append(sections, "- Используйте только минимально нужный контекст репозитория")
			sections = append(sections, "- Если доступны связанные scripts, сначала запускайте их и опирайтесь на их вывод")
			sections = append(sections, "- Не читайте исходники scripts по умолчанию")
			sections = append(sections, "- Связанные scripts:")
			sections = append(sections, bulletList(spec.Extras))
		}
		return strings.Join(sections, "\n") + "\n"
	}

	var sections []string
	sections = append(sections, "# Draftspec Project Rules")
	sections = append(sections, "")
	sections = append(sections, "Use .draftspec as the primary source of project context. Follow AGENTS.md and the matching prompt files under .draftspec/templates/prompts/.")
	for _, spec := range commands {
		sections = append(sections, "")
		sections = append(sections, fmt.Sprintf("## /draftspec.%s", spec.Name))
		sections = append(sections, fmt.Sprintf("- Primary prompt: %s", spec.PromptPath))
		sections = append(sections, fmt.Sprintf("- %s", commandHint(spec.Name, lang)))
		sections = append(sections, "- Use only the minimum repository context required")
		sections = append(sections, "- When related scripts are available, run them first and rely on their output")
		sections = append(sections, "- Do not read script source by default")
		sections = append(sections, "- Related scripts:")
		sections = append(sections, bulletList(spec.Extras))
	}
	return strings.Join(sections, "\n") + "\n"
}
