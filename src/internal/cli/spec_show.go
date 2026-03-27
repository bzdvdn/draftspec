package cli

import (
	"fmt"

	"draftspec/src/internal/specs"
	"github.com/spf13/cobra"
)

func newShowSpecCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show-spec <name> [path]",
		Short: "Show a specification by slug",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			root := "."
			if len(args) == 2 {
				root = args[1]
			}

			content, err := specs.Show(root, args[0])
			if err != nil {
				return err
			}

			fmt.Fprint(cmd.OutOrStdout(), content)
			return nil
		},
	}
}
