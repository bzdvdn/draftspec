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

	return cmd
}
