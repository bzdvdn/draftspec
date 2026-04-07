package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newHelpCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "help [command]",
		Short: "Show help for any command",
		Long: `Shows help for any draftspec command.

Examples:
  draftspec help
  draftspec help init
  draftspec help refresh`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			target, _, err := cmd.Root().Find(args)
			if target == nil || err != nil {
				unknown := strings.Join(args, " ")
				if unknown == "" {
					unknown = "<root>"
				}
				fmt.Fprintf(cmd.ErrOrStderr(), "Unknown help topic %q\n", unknown)
				_ = cmd.Root().Usage()
				return newExitError(1, "")
			}

			title := "draftspec help"
			if target != cmd.Root() {
				title = "draftspec help: " + target.CommandPath()
			}
			printPanel(cmd.OutOrStdout(), title, []string{
				"Tip: add " + styleCmd(cmd.OutOrStdout(), "--help") + " to any command.",
			})

			target.InitDefaultHelpFlag()
			target.InitDefaultVersionFlag()
			return target.Help()
		},
	}
	return cmd
}
