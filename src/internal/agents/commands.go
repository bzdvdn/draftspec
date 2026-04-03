package agents

type CommandDefinition struct {
	Name        string
	Description string
	PromptPath  string
	Extras      []string
	Optional    bool
	Category    string
}

func DefaultCommands(shell string) []CommandDefinition {
	normalizedShell := normalizeShell(shell)
	launcher := scriptPath("run-draftspec", normalizedShell)

	return []CommandDefinition{
		{
			Name:        "constitution",
			Description: "Create or update the project constitution",
			PromptPath:  ".draftspec/templates/prompts/constitution.md",
			Extras:      []string{launcher, scriptPath("check-constitution", normalizedShell)},
			Category:    "workflow",
		},
		{
			Name:        "spec",
			Description: "Create or update one feature spec",
			PromptPath:  ".draftspec/templates/prompts/spec.md",
			Extras:      []string{launcher, scriptPath("check-spec-ready", normalizedShell)},
			Category:    "workflow",
		},
		{
			Name:        "inspect",
			Description: "Inspect one feature for consistency and quality",
			PromptPath:  ".draftspec/templates/prompts/inspect.md",
			Extras: []string{
				launcher,
				scriptPath("check-inspect-ready", normalizedShell),
				scriptPath("inspect-spec", normalizedShell),
			},
			Category: "workflow",
		},
		{
			Name:        "plan",
			Description: "Create or update one feature plan package",
			PromptPath:  ".draftspec/templates/prompts/plan.md",
			Extras:      []string{launcher, scriptPath("check-plan-ready", normalizedShell)},
			Category:    "workflow",
		},
		{
			Name:        "tasks",
			Description: "Create or update tasks for one feature",
			PromptPath:  ".draftspec/templates/prompts/tasks.md",
			Extras:      []string{launcher, scriptPath("check-tasks-ready", normalizedShell)},
			Category:    "workflow",
		},
		{
			Name:        "implement",
			Description: "Implement one feature from tasks",
			PromptPath:  ".draftspec/templates/prompts/implement.md",
			Extras: []string{
				launcher,
				scriptPath("check-implement-ready", normalizedShell),
				scriptPath("list-open-tasks", normalizedShell),
			},
			Category: "workflow",
		},
		{
			Name:        "verify",
			Description: "Verify one implemented feature package",
			PromptPath:  ".draftspec/templates/prompts/verify.md",
			Extras: []string{
				launcher,
				scriptPath("check-verify-ready", normalizedShell),
				scriptPath("verify-task-state", normalizedShell),
			},
			Category: "workflow",
		},
		{
			Name:        "archive",
			Description: "Archive one feature package",
			PromptPath:  ".draftspec/templates/prompts/archive.md",
			Extras:      []string{launcher, scriptPath("check-archive-ready", normalizedShell)},
			Category:    "workflow",
		},
		{
			Name:        "handoff",
			Description: "Generate a session handoff document for one feature",
			PromptPath:  ".draftspec/templates/prompts/handoff.md",
			Extras:      []string{launcher, scriptPath("list-open-tasks", normalizedShell)},
			Category:    "workflow",
		},
		{
			Name:        "challenge",
			Description: "Adversarial review of a feature spec or plan",
			PromptPath:  ".draftspec/templates/prompts/challenge.md",
			Extras:      []string{launcher},
			Optional:    true,
			Category:    "workflow",
		},
		{
			Name:        "scope",
			Description: "Quick scope boundary check for a feature",
			PromptPath:  ".draftspec/templates/prompts/scope.md",
			Extras:      []string{launcher},
			Optional:    true,
			Category:    "workflow",
		},
		{
			Name:        "recap",
			Description: "Project-level overview of all active features and their current phase",
			PromptPath:  ".draftspec/templates/prompts/recap.md",
			Extras:      []string{launcher, scriptPath("list-specs", normalizedShell)},
			Optional:    true,
			Category:    "workflow",
		},
		{
			Name:        "hotfix",
			Description: "Create emergency fix outside the standard phase chain",
			PromptPath:  ".draftspec/templates/prompts/hotfix.md",
			Extras:      []string{launcher},
			Optional:    true,
			Category:    "workflow",
		},
	}
}
