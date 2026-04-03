package cli

import (
	"fmt"

	"draftspec/src/internal/project"
	"github.com/spf13/cobra"
)

func newDemoCmd() *cobra.Command {
	var shell string
	var agentTargets []string

	cmd := &cobra.Command{
		Use:   "demo [path]",
		Short: "Create a demo workspace with pre-populated example artifacts",
		Long: `Create a demo workspace at the given path (default: ./draftspec-demo).

The demo workspace contains a fully worked example feature (export-report) at
the implement phase — spec, inspect report, plan, tasks, and data model are all
populated so you can explore the workflow and try slash commands immediately.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root := "./draftspec-demo"
			if len(args) == 1 {
				root = args[0]
			}

			result, err := project.Demo(root, project.DemoOptions{
				Shell:        shell,
				AgentTargets: agentTargets,
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

	cmd.Flags().StringVar(&shell, "shell", "sh", "shell for generated workflow scripts: sh or powershell")
	cmd.Flags().StringSliceVar(&agentTargets, "agents", nil, "generate agent command files for one or more targets: claude, codex, copilot, cursor, kilocode, trae, windsurf, roocode, aider, all")

	return cmd
}
