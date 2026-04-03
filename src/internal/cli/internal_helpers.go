package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"draftspec/src/internal/config"
	"draftspec/src/internal/featurepaths"
	"github.com/spf13/cobra"
)

func newInternalListOpenTasksCmd() *cobra.Command {
	var root string
	cmd := &cobra.Command{
		Use:           "list-open-tasks <slug>",
		Hidden:        true,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(root)
			if err != nil {
				return err
			}
			tasksDisplay := filepath.ToSlash(filepath.Join(cfg.Paths.PlansDir, args[0], "tasks.md"))
			tasksPath := filepath.Join(root, filepath.FromSlash(tasksDisplay))
			content, err := os.ReadFile(tasksPath)
			if err != nil {
				if os.IsNotExist(err) {
					return newExitError(1, fmt.Sprintf("tasks file not found: %s", tasksDisplay))
				}
				return err
			}
			for _, line := range strings.Split(string(content), "\n") {
				if strings.HasPrefix(line, "- [ ]") {
					fmt.Fprintln(cmd.OutOrStdout(), line)
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "Draftspec project root")
	return cmd
}

func newInternalListSpecsCmd() *cobra.Command {
	var root string
	cmd := &cobra.Command{
		Use:           "list-specs [specs-dir]",
		Hidden:        true,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			specsDir := ".draftspec/specs"
			if len(args) == 1 {
				specsDir = args[0]
			}
			if !filepath.IsAbs(specsDir) {
				specsDir = filepath.Join(root, filepath.FromSlash(specsDir))
			}
			if _, err := os.Stat(specsDir); err != nil {
				if os.IsNotExist(err) {
					return newExitError(1, fmt.Sprintf("specs directory not found: %s", filepath.ToSlash(argsOrDefault(args, ".draftspec/specs"))))
				}
				return err
			}
			names, err := featurepaths.ListSpecSlugs(specsDir)
			if err != nil {
				return err
			}
			for _, name := range names {
				fmt.Fprintln(cmd.OutOrStdout(), name)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "Draftspec project root")
	return cmd
}

func newInternalShowSpecCmd() *cobra.Command {
	var root string
	cmd := &cobra.Command{
		Use:           "show-spec <spec-name> [specs-dir]",
		Hidden:        true,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			specsDirDisplay := ".draftspec/specs"
			if len(args) == 2 {
				specsDirDisplay = args[1]
			}
			specsDir := specsDirDisplay
			if !filepath.IsAbs(specsDir) {
				specsDir = filepath.Join(root, filepath.FromSlash(specsDir))
			}
			specFileDisplay := filepath.ToSlash(filepath.Join(specsDirDisplay, args[0], "spec.md"))
			specFilePath, legacy := featurepaths.ResolveSpec(specsDir, args[0])
			if legacy {
				specFileDisplay = filepath.ToSlash(filepath.Join(specsDirDisplay, args[0]+".md"))
			}
			content, err := os.ReadFile(specFilePath)
			if err != nil {
				if os.IsNotExist(err) {
					return newExitError(1, fmt.Sprintf("spec not found: %s", specFileDisplay))
				}
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), string(content))
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "Draftspec project root")
	return cmd
}

func newInternalLinkAgentsCmd() *cobra.Command {
	var root string
	cmd := &cobra.Command{
		Use:           "link-agents [agents-file] [snippet-file]",
		Hidden:        true,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			agentsFile := "AGENTS.md"
			if len(args) >= 1 {
				agentsFile = args[0]
			}
			snippetFile := ".draftspec/templates/agents-snippet.md"
			if len(args) == 2 {
				snippetFile = args[1]
			}

			agentsPath := agentsFile
			if !filepath.IsAbs(agentsPath) {
				agentsPath = filepath.Join(root, filepath.FromSlash(agentsFile))
			}
			snippetPath := snippetFile
			if !filepath.IsAbs(snippetPath) {
				snippetPath = filepath.Join(root, filepath.FromSlash(snippetFile))
			}

			if _, err := os.Stat(agentsPath); os.IsNotExist(err) {
				if err := os.WriteFile(agentsPath, nil, 0o644); err != nil {
					return err
				}
			}

			contentBytes, err := os.ReadFile(agentsPath)
			if err != nil {
				return err
			}
			content := string(contentBytes)
			if strings.Contains(content, "## Draftspec") {
				fmt.Fprintf(cmd.OutOrStdout(), "Draftspec block already present in %s\n", agentsFile)
				return nil
			}

			snippetBytes, err := os.ReadFile(snippetPath)
			if err != nil {
				return err
			}
			if len(content) > 0 && !strings.HasSuffix(content, "\n") {
				content += "\n"
			}
			content += "\n" + string(snippetBytes)
			if err := os.WriteFile(agentsPath, []byte(content), 0o644); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Draftspec block added to %s\n", agentsFile)
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "Draftspec project root")
	return cmd
}

func argsOrDefault(values []string, fallback string) string {
	if len(values) == 0 {
		return fallback
	}
	return values[0]
}
