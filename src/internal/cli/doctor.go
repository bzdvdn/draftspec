package cli

import (
	"encoding/json"
	"fmt"

	"draftspec/src/internal/doctor"
	"github.com/spf13/cobra"
)

func newDoctorCmd() *cobra.Command {
	var jsonOutput bool

	cmd := &cobra.Command{
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
			if jsonOutput {
				payload, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(payload))
				return nil
			}
			for _, finding := range result.Findings {
				fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n", finding.Level, finding.Message)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output doctor findings as JSON")
	return cmd
}
