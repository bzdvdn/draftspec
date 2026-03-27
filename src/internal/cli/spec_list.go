package cli

import (
	"fmt"

	"draftspec/src/internal/specs"
	"github.com/spf13/cobra"
)

func newListSpecsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list-specs",
		Short: "List available specifications",
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := specs.List(".")
			if err != nil {
				return err
			}

			for _, name := range names {
				fmt.Fprintln(cmd.OutOrStdout(), name)
			}

			return nil
		},
	}
}
