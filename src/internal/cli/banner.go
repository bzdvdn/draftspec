package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

const draftspecASCII = `
██████╗ ██████╗  █████╗ ███████╗████████╗███████╗██████╗ ███████╗ ██████╗
██╔══██╗██╔══██╗██╔══██╗██╔════╝╚══██╔══╝██╔════╝██╔══██╗██╔════╝██╔════╝
██║  ██║██████╔╝███████║█████╗     ██║   ███████╗██████╔╝█████╗  ██║
██║  ██║██╔══██╗██╔══██║██╔══╝     ██║   ╚════██║██╔═══╝ ██╔══╝  ██║
██████╔╝██║  ██║██║  ██║██║        ██║   ███████║██║     ███████╗╚██████╗
╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝        ╚═╝   ╚══════╝╚═╝     ╚══════╝ ╚═════╝
`

func init() {
	cobra.AddTemplateFunc("draftspecBanner", func(cmd *cobra.Command) string {
		return renderDraftspecBanner(cmd)
	})
}

func renderDraftspecBanner(cmd *cobra.Command) string {
	out := cmd.OutOrStdout()
	color := useColor(out)

	art := strings.TrimLeft(draftspecASCII, "\n") + "\n"
	tagline := "Draftspec — Spec-Driven Development Toolkit\n"
	if cmd != nil && cmd.Root() != nil && cmd != cmd.Root() {
		tagline = fmt.Sprintf("Draftspec — %s\n", cmd.CommandPath())
	}

	if !color {
		return art + tagline + "\n"
	}

	return ansiCyan + art + ansiReset + ansiYellow + tagline + ansiReset + "\n"
}
