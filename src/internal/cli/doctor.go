package cli

import (
	"fmt"

	"draftspec/src/internal/doctor"
	"github.com/spf13/cobra"
)

func newDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor [path]",
		Short: "Check Draftspec workspace health and agent target consistency",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root := "."
			if len(args) == 1 {
				root = args[0]
			}
			result, err := doctor.Check(root)
			if err != nil {
				return err
			}
			for _, finding := range result.Findings {
				fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n", finding.Level, finding.Message)
			}
			return nil
		},
	}
}
