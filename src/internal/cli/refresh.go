package cli

import (
	"encoding/json"
	"fmt"

	"draftspec/src/internal/project"
	"github.com/spf13/cobra"
)

func newRefreshCmd() *cobra.Command {
	var defaultLang string
	var docsLang string
	var agentLang string
	var commentsLang string
	var shell string
	var jsonOutput bool
	var dryRun bool
	var agentTargets []string
	var legacyAgentTargets []string

	cmd := &cobra.Command{
		Use:   "refresh [path]",
		Short: "Refresh generated Draftspec artifacts for an existing project without touching authored feature state",
		Long: `Refreshes Draftspec-managed artifacts inside an already initialized project.

Synchronizes:
  - .draftspec/draftspec.yaml (config)
  - managed templates/prompts/scripts inside .draftspec/
  - the managed Draftspec block in AGENTS.md
  - agent-target artifacts (.claude/, .cursor/, etc.)

Does not touch authored feature state:
  - .draftspec/specs/** and .draftspec/plans/** are not modified.

Tip: use --dry-run to preview changes without writing.`,
		Example: "  draftspec refresh .\n  draftspec refresh . --dry-run\n  draftspec refresh . --agents claude,cursor --agent-lang en",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root := "."
			if len(args) == 1 {
				root = args[0]
			}

			if !jsonOutput {
				printPanel(cmd.OutOrStdout(), "draftspec refresh", []string{
					"Sync managed files (without modifying specs/plans).",
					"Tip: add --dry-run to preview changes.",
				})
			}

			result, err := project.Refresh(root, project.RefreshOptions{
				DefaultLang:  defaultLang,
				DocsLang:     docsLang,
				AgentLang:    agentLang,
				CommentsLang: commentsLang,
				Shell:        shell,
				AgentTargets: append(agentTargets, legacyAgentTargets...),
				DryRun:       dryRun,
			})
			if err != nil {
				return err
			}

			if jsonOutput {
				payload, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(payload))
				return nil
			}

			for _, line := range result.Messages {
				fmt.Fprintln(cmd.OutOrStdout(), line)
			}

			printPanel(cmd.OutOrStdout(), "refresh summary", []string{
				fmt.Sprintf("created: %d", len(result.Created)),
				fmt.Sprintf("updated: %d", len(result.Updated)),
				fmt.Sprintf("removed: %d", len(result.Removed)),
				fmt.Sprintf("unchanged: %d", len(result.Unchanged)),
				fmt.Sprintf("dry-run: %t", result.DryRun),
			})
			return nil
		},
	}

	cmd.Flags().StringVar(&defaultLang, "lang", "", "override the base language for generated docs and prompts: en or ru")
	cmd.Flags().StringVar(&docsLang, "docs-lang", "", "override the generated documentation language: en or ru")
	cmd.Flags().StringVar(&agentLang, "agent-lang", "", "override the generated prompt and AGENTS guidance language: en or ru")
	cmd.Flags().StringVar(&commentsLang, "comments-lang", "", "override the preferred code comment language: en or ru")
	cmd.Flags().StringVar(&shell, "shell", "", "override the generated workflow script family: sh or powershell")
	cmd.Flags().StringSliceVar(&agentTargets, "agents", nil, "override enabled project-local agent targets: claude, codex, copilot, cursor, kilocode, trae, all")
	cmd.Flags().StringSliceVar(&legacyAgentTargets, "agent", nil, "deprecated alias for --agents")
	cmd.Flags().MarkHidden("agent")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show which managed files would change without writing them")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output refresh results as JSON")

	return cmd
}
