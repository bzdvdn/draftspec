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
		sections = append(sections, "")
		sections = append(sections, workflowChainHint(lang))
		for _, spec := range commands {
			sections = append(sections, "")
			sections = append(sections, fmt.Sprintf("## /draftspec.%s", spec.Name))
			sections = append(sections, fmt.Sprintf("- Основной prompt: %s", spec.PromptPath))
			sections = append(sections, fmt.Sprintf("- %s", commandHint(spec.Name, lang)))
			sections = append(sections, "- Используйте только минимально нужный контекст репозитория")
			sections = append(sections, fmt.Sprintf("- %s", scriptExecutionHint(lang)))
			sections = append(sections, fmt.Sprintf("- %s", helpDiscoveryHint(lang)))
			if len(spec.Extras) > 0 {
				sections = append(sections, scriptListBlock(spec.Extras, lang))
			}
		}
		sections = append(sections, "")
		sections = append(sections, antiPatternHint(lang))
		return strings.Join(sections, "\n") + "\n"
	}

	var sections []string
	sections = append(sections, "# Draftspec Project Rules")
	sections = append(sections, "")
	sections = append(sections, "Use .draftspec as the primary source of project context. Follow AGENTS.md and the matching prompt files under .draftspec/templates/prompts/.")
	sections = append(sections, "")
	sections = append(sections, workflowChainHint(lang))
	for _, spec := range commands {
		sections = append(sections, "")
		sections = append(sections, fmt.Sprintf("## /draftspec.%s", spec.Name))
		sections = append(sections, fmt.Sprintf("- Primary prompt: %s", spec.PromptPath))
		sections = append(sections, fmt.Sprintf("- %s", commandHint(spec.Name, lang)))
		sections = append(sections, "- Use only the minimum repository context required")
		sections = append(sections, fmt.Sprintf("- %s", scriptExecutionHint(lang)))
		sections = append(sections, fmt.Sprintf("- %s", helpDiscoveryHint(lang)))
		if len(spec.Extras) > 0 {
			sections = append(sections, scriptListBlock(spec.Extras, lang))
		}
	}
	sections = append(sections, "")
	sections = append(sections, antiPatternHint(lang))
	return strings.Join(sections, "\n") + "\n"
}
