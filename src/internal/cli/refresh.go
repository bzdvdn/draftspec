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
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root := "."
			if len(args) == 1 {
				root = args[0]
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
