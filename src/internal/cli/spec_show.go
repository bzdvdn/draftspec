package cli

import (
	"fmt"

	"draftspec/src/internal/specs"
	"github.com/spf13/cobra"
)

func newShowSpecCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show-spec <name>",
		Short: "Show a specification by slug",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			content, err := specs.Show(".", args[0])
			if err != nil {
				return err
			}

			fmt.Fprint(cmd.OutOrStdout(), content)
			return nil
		},
	}
}
