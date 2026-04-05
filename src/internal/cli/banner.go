package cli

import (
	"fmt"
	"os"
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
	color := shouldUseColor(out)

	art := strings.TrimLeft(draftspecASCII, "\n") + "\n"
	tagline := "Draftspec — Spec-Driven Development Toolkit\n"
	if cmd != nil && cmd.Root() != nil && cmd != cmd.Root() {
		tagline = fmt.Sprintf("Draftspec — %s\n", cmd.CommandPath())
	}

	if !color {
		return art + tagline + "\n"
	}

	cyan := "\x1b[36m"
	yellow := "\x1b[33m"
	reset := "\x1b[0m"
	return cyan + art + reset + yellow + tagline + reset + "\n"
}

func shouldUseColor(w any) bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	if strings.EqualFold(os.Getenv("TERM"), "dumb") {
		return false
	}
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) != 0
}
