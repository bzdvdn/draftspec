package cli

import (
	"fmt"

	"draftspec/src/internal/project"
	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	var initGit bool
	var defaultLang string
	var docsLang string
	var agentLang string
	var commentsLang string
	var agentTargets []string
	var legacyAgentTargets []string

	cmd := &cobra.Command{
		Use:   "init [path]",
		Short: "Initialize a .draftspec workspace in the target project",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root := "."
			if len(args) == 1 {
				root = args[0]
			}

			result, err := project.Initialize(root, project.InitOptions{
				InitGit:      initGit,
				DefaultLang:  defaultLang,
				DocsLang:     docsLang,
				AgentLang:    agentLang,
				CommentsLang: commentsLang,
				AgentTargets: append(agentTargets, legacyAgentTargets...),
			})
			if err != nil {
				return err
			}

			for _, line := range result.Messages {
				fmt.Fprintln(cmd.OutOrStdout(), line)
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&initGit, "git", true, "initialize a git repository when one does not exist")
	cmd.Flags().StringVar(&defaultLang, "lang", "en", "base language for generated docs and prompts: en or ru")
	cmd.Flags().StringVar(&docsLang, "docs-lang", "", "language for generated project documentation: en or ru")
	cmd.Flags().StringVar(&agentLang, "agent-lang", "", "language for generated agent prompts and AGENTS guidance: en or ru")
	cmd.Flags().StringVar(&commentsLang, "comments-lang", "", "preferred language for code comments: en or ru")
	cmd.Flags().StringSliceVar(&agentTargets, "agents", nil, "generate project-local agent command files for one or more targets: claude, codex, copilot, cursor, kilocode, trae, all")
	cmd.Flags().StringSliceVar(&legacyAgentTargets, "agent", nil, "deprecated alias for --agents")
	cmd.Flags().MarkHidden("agent")

	return cmd
}
