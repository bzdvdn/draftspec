package cli

import "github.com/spf13/cobra"

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "draftspec",
		Short: "A lightweight project context kit for development agents and humans",
	}

	cmd.AddCommand(newInitCmd())
	cmd.AddCommand(newRefreshCmd())
	cmd.AddCommand(newAddAgentCmd())
	cmd.AddCommand(newListAgentsCmd())
	cmd.AddCommand(newRemoveAgentCmd())
	cmd.AddCommand(newCleanupAgentsCmd())
	cmd.AddCommand(newDoctorCmd())
	cmd.AddCommand(newListSpecsCmd())
	cmd.AddCommand(newShowSpecCmd())

	return cmd
}
