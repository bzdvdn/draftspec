package project

import (
	"fmt"
	"os"
	"path/filepath"

	"draftspec/src/internal/templates"
)

type DemoOptions struct {
	Shell        string
	AgentTargets []string
}

type DemoResult struct{ Messages []string }

func Demo(root string, options DemoOptions) (DemoResult, error) {
	shell := options.Shell
	if shell == "" {
		shell = "sh"
	}

	initResult, err := Initialize(root, InitOptions{
		InitGit:  false,
		DefaultLang: "en",
		Shell:    shell,
		AgentTargets: options.AgentTargets,
	})
	if err != nil {
		return DemoResult{}, err
	}

	absRoot, err := filepath.Abs(root)
	if err != nil {
		return DemoResult{}, err
	}

	demoFiles, err := templates.DemoFiles()
	if err != nil {
		return DemoResult{}, fmt.Errorf("load demo files: %w", err)
	}

	draftspecDir := filepath.Join(absRoot, ".draftspec")
	messages := append([]string(nil), initResult.Messages...)

	for _, file := range demoFiles {
		target := filepath.Join(draftspecDir, file.TargetPath)
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return DemoResult{}, err
		}
		if err := os.WriteFile(target, []byte(file.Content), file.Mode); err != nil {
			return DemoResult{}, err
		}
		messages = append(messages, fmt.Sprintf("created demo artifact %s", rel(absRoot, target)))
	}

	messages = append(messages, "")
	messages = append(messages, "Demo workspace ready. Example feature: export-report (phase: implement, 2/6 tasks done)")
	messages = append(messages, "Try: /draftspec.scope export-report")
	messages = append(messages, "Try: /draftspec.challenge export-report")
	messages = append(messages, "Try: /draftspec.handoff export-report")

	return DemoResult{Messages: messages}, nil
}
