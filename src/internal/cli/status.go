package cli

import (
	"encoding/json"
	"fmt"

	"draftspec/src/internal/status"
	"github.com/spf13/cobra"
)

func newStatusCmd() *cobra.Command {
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "status <slug> [path]",
		Short: "Show feature workflow status for one slug",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			root := "."
			if len(args) == 2 {
				root = args[1]
			}

			result, err := status.Check(root, args[0])
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

			fmt.Fprintf(cmd.OutOrStdout(), "slug: %s\n", result.Slug)
			fmt.Fprintf(cmd.OutOrStdout(), "phase: %s\n", result.Phase)
			fmt.Fprintf(cmd.OutOrStdout(), "spec_exists: %t\n", result.SpecExists)
			fmt.Fprintf(cmd.OutOrStdout(), "plan_exists: %t\n", result.PlanExists)
			fmt.Fprintf(cmd.OutOrStdout(), "tasks_exists: %t\n", result.TasksExists)
			if result.TasksExists {
				fmt.Fprintf(cmd.OutOrStdout(), "tasks: total=%d completed=%d open=%d\n", result.TasksTotal, result.TasksCompleted, result.TasksOpen)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "archived: %t\n", result.Archived)
			if result.ReadyFor != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "ready_for: %s\n", result.ReadyFor)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "blocked: %t\n", result.Blocked)
			return nil
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output feature status as JSON")
	return cmd
}
