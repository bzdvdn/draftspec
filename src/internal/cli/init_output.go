package cli

import (
	"fmt"
	"io"
	"strings"

	"draftspec/src/internal/project"
)

func printInitOutput(w io.Writer, result project.InitResult) {
	color := shouldUseColor(w)

	agentTargets := "none"
	if len(result.AgentTargets) > 0 {
		agentTargets = strings.Join(result.AgentTargets, ", ")
	}

	fmt.Fprintf(w, "Selected script type: %s\n", result.Shell)
	fmt.Fprintf(w, "Configured languages: docs=%s agent=%s comments=%s\n", result.DocsLang, result.AgentLang, result.CommentsLang)
	fmt.Fprintf(w, "enabled agent targets: %s\n\n", agentTargets)

	fmt.Fprintln(w, "Initialize Draftspec Project")
	steps := initSteps(result)
	for i, step := range steps {
		last := i == len(steps)-1
		printStepLine(w, step, last, color)
	}
	fmt.Fprintln(w)

	if color {
		fmt.Fprintln(w, "\x1b[32mProject ready.\x1b[0m")
	} else {
		fmt.Fprintln(w, "Project ready.")
	}
	fmt.Fprintln(w)

	printPanel(w, "Agent Folder Security", []string{
		"Some agents may store credentials, auth tokens, or other private artifacts",
		"inside project-level agent folders (e.g. " + stylePath(".claude/", color) + ", " + stylePath(".cursor/", color) + ", " + stylePath(".kilocode/", color) + ").",
		"Consider adding them (or parts of them) to .gitignore to prevent leaks.",
	})

	printPanel(w, "Next Steps", []string{
		"1. You're already in the project directory.",
		"2. Start using slash commands with your AI agent:",
		"   2.1 " + styleCmd("/draftspec.constitution", color) + "  - Establish project principles",
		"   2.2 " + styleCmd("/draftspec.spec", color) + "          - Create a baseline specification",
		"   2.3 " + styleCmd("/draftspec.plan", color) + "          - Create an implementation plan",
		"   2.4 " + styleCmd("/draftspec.tasks", color) + "         - Generate actionable tasks",
		"   2.5 " + styleCmd("/draftspec.implement", color) + "     - Execute implementation",
	})

	printPanel(w, "Enhancement Commands", []string{
		styleCmd("/draftspec.challenge", color) + " (optional) - adversarial review of spec/plan",
		styleCmd("/draftspec.scope", color) + " (optional)     - scope boundary check",
		styleCmd("/draftspec.recap", color) + " (optional)     - recap active features",
	})

	printPanel(w, "Useful CLI Commands", []string{
		styleCmd("draftspec doctor .", color) + "        - workspace health check",
		styleCmd("draftspec list-specs .", color) + "    - list active specs",
		styleCmd("draftspec check <slug> .", color) + "  - feature status",
		styleCmd("draftspec dashboard .", color) + "     - visual dashboard",
	})
}

type initStep struct {
	Label  string
	Status string // ok, kept, updated, skipped
	Detail string
}

func initSteps(result project.InitResult) []initStep {
	createdCount := len(result.Created)
	keptCount := len(result.Kept)

	templateDetail := fmt.Sprintf("%d created, %d kept", createdCount, keptCount)

	agentsStatus := "skipped"
	agentsDetail := "no targets"
	if len(result.AgentTargets) > 0 {
		agentsStatus = "ok"
		agentsDetail = strings.Join(result.AgentTargets, ", ")
	}

	agentsSnippetStatus := "kept"
	if result.AgentsSnippetChanged {
		agentsSnippetStatus = "updated"
	}

	gitStatus := result.GitRepoStatus
	if gitStatus == "" {
		gitStatus = "ok"
	}

	return []initStep{
		{Label: "Ensure .draftspec directory structure", Status: "ok"},
		{Label: "Install managed templates", Status: "ok", Detail: templateDetail},
		{Label: "Link Draftspec block in AGENTS.md", Status: agentsSnippetStatus},
		{Label: "Generate agent target artifacts", Status: agentsStatus, Detail: agentsDetail},
		{Label: "Initialize git repository", Status: gitStatus},
		{Label: "Finalize", Status: "ok", Detail: "project ready"},
	}
}

func printStepLine(w io.Writer, step initStep, last bool, color bool) {
	connector := "├─"
	if last {
		connector = "└─"
	}

	status := renderStepStatus(step.Status, color)
	suffix := renderStepSuffix(step.Status, step.Detail)
	fmt.Fprintf(w, "  %s %s %s (%s)\n", connector, status, step.Label, suffix)
}

func renderStepStatus(status string, color bool) string {
	switch status {
	case "ok":
		return colorize("●", "\x1b[32m", color)
	case "updated":
		return colorize("●", "\x1b[36m", color)
	case "kept":
		return colorize("●", "\x1b[33m", color)
	case "skipped":
		return colorize("○", "\x1b[90m", color)
	case "initialized":
		return colorize("●", "\x1b[32m", color)
	default:
		return colorize("●", "\x1b[32m", color)
	}
}

func renderStepSuffix(status string, detail string) string {
	label := status
	if status == "initialized" {
		label = "ok"
	}
	if detail == "" {
		return label
	}
	return label + ": " + detail
}

func colorize(s string, code string, enabled bool) string {
	if !enabled {
		return s
	}
	return code + s + "\x1b[0m"
}

func styleCmd(s string, enabled bool) string {
	return colorize(s, "\x1b[36m", enabled)
}

func stylePath(s string, enabled bool) string {
	return colorize(s, "\x1b[33m", enabled)
}
