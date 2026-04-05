package cli

import (
	"fmt"
	"strings"

	"draftspec/src/internal/workflow"
	"github.com/spf13/cobra"
)

func newDashboardCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dashboard [path]",
		Short: "Visual dashboard of all features and project health",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root := "."
			if len(args) == 1 {
				root = args[0]
			}

			states, err := workflow.States(root)
			if err != nil {
				return err
			}

			if len(states) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No features found in this project.")
				return nil
			}

			w := cmd.OutOrStdout()
			fmt.Fprintln(w, "┌──────────────────────────────────────────────────────────────────────────────┐")
			fmt.Fprintln(w, "│                          DRAFTSPEC PROJECT DASHBOARD                         │")
			fmt.Fprintln(w, "└──────────────────────────────────────────────────────────────────────────────┘")
			fmt.Fprintln(w)

			fmt.Fprintln(w, "  FEATURE SLUG           PHASE        PROGRESS   STATUS    BRANCH              ")
			fmt.Fprintln(w, "  ─────────────────────  ───────────  ─────────  ────────  ────────────────────")

			for _, s := range states {
				if s.Archived {
					continue
				}

				progress := "0%"
				if s.TasksTotal > 0 {
					progress = fmt.Sprintf("%d%%", (s.TasksCompleted*100)/s.TasksTotal)
				}

				status := "READY"
				if s.Blocked {
					status = "BLOCKED"
				}

				branchInfo := s.CurrentBranch
				if s.BranchMismatch {
					branchInfo = "!! " + s.CurrentBranch
				}

				fmt.Fprintf(w, "  %-21s  %-11s  %-9s  %-8s  %-20s\n",
					truncate(s.Slug, 21),
					strings.ToUpper(s.Phase),
					progress,
					status,
					truncate(branchInfo, 20),
				)
			}

			fmt.Fprintln(w)
			fmt.Fprintln(w, "  LEGEND: !! = Branch Mismatch | READY/BLOCKED = Phase Status")
			fmt.Fprintln(w)

			return nil
		},
	}

	return cmd
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
