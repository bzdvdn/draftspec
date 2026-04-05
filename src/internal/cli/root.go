package cli

import "github.com/spf13/cobra"

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "draftspec",
		Short: "A lightweight project context kit for development agents and humans",
		Long: `draftspec — specification-driven development kit for agents and humans

Strict phase chain: constitution → spec → inspect → plan → tasks → implement → verify → archive

Quick start:
  draftspec init . --lang en --shell sh --agents codex
  draftspec doctor .
  draftspec list-specs .

For agents (Kilocode/Claude/Cursor):
  /draftspec.constitution                        — create a constitution
  /draftspec.spec --name "feature name"          — create a spec
  /draftspec.spec --amend                        — targeted spec edit
  /draftspec.plan <slug> [--research|--update]   — create a plan
  /draftspec.tasks <slug>                        — decompose into tasks
  /draftspec.implement <slug>                    — implement tasks
  /draftspec.verify <slug> [--deep]              — verify AC coverage
  /draftspec.archive <slug> [--restore]          — archive/restore

Optional commands (any phase):
  /draftspec.challenge <slug> [--spec|--plan]    — adversarial review
  /draftspec.handoff [slug]                      — session handoff doc
  /draftspec.hotfix <slug>                       — emergency fix
  /draftspec.scope <slug>                        — scope boundary check
  /draftspec.recap                               — recap active features

CLI commands:
  draftspec doctor .                      — workspace health check
  draftspec list-specs .                  — list active specs
  draftspec check <slug> . [--json]       — feature status
  draftspec check . --all                 — all features table
  draftspec dashboard .                   — visual dashboard
  draftspec trace <slug> . [--tests]      — code traceability
  draftspec export <slug> . --output f.md — export artifacts

Documentation:
  README.md — overview and examples
  docs/en/ or docs/ru/ — extended documentation`,
		Version: Version,
	}

	cmd.SetHelpTemplate(`{{draftspecBanner .}}{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`)

	cmd.SetHelpCommand(newHelpCmd())

	cmd.AddCommand(newInitCmd())
	cmd.AddCommand(newRefreshCmd())
	cmd.AddCommand(newAddAgentCmd())
	cmd.AddCommand(newListAgentsCmd())
	cmd.AddCommand(newRemoveAgentCmd())
	cmd.AddCommand(newCleanupAgentsCmd())
	cmd.AddCommand(newDoctorCmd())
	cmd.AddCommand(newStatusCmd())
	cmd.AddCommand(newDashboardCmd())
	cmd.AddCommand(newFeatureCmd())
	cmd.AddCommand(newFeaturesCmd())
	cmd.AddCommand(newMigrateCmd())
	cmd.AddCommand(newListSpecsCmd())
	cmd.AddCommand(newShowSpecCmd())
	cmd.AddCommand(newCheckCmd())
	cmd.AddCommand(newTraceCmd())
	cmd.AddCommand(newDemoCmd())
	cmd.AddCommand(newExportCmd())
	cmd.AddCommand(newExploreCmd())
	cmd.AddCommand(newContextCmd())
	cmd.AddCommand(newSchemaCmd())
	cmd.AddCommand(newRiskCmd())
	cmd.AddCommand(newInternalCmd())

	return cmd
}
