package project

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"draftspec/src/internal/config"
	"draftspec/src/internal/gitutil"
	"draftspec/src/internal/templates"
)

type InitOptions struct {
	InitGit      bool
	DefaultLang  string
	DocsLang     string
	AgentLang    string
	CommentsLang string
}

type InitResult struct { Messages []string }

func Initialize(root string, options InitOptions) (InitResult, error) {
	root, err := filepath.Abs(root)
	if err != nil { return InitResult{}, err }
	languages, err := templates.ResolveLanguageSettings(options.DefaultLang, options.DocsLang, options.AgentLang, options.CommentsLang)
	if err != nil { return InitResult{}, err }
	cfg := config.Default()
	cfg.Language.Default = languages.Default
	cfg.Language.Docs = languages.Docs
	cfg.Language.Agent = languages.Agent
	cfg.Language.Comments = languages.Comments
	var messages []string
	if options.InitGit {
		created, err := gitutil.EnsureRepository(root)
		if err != nil { return InitResult{}, err }
		if created { messages = append(messages, "initialized git repository") } else { messages = append(messages, "kept existing git repository") }
	} else {
		messages = append(messages, "skipped git repository initialization")
	}
	draftspecDir, err := cfg.DraftspecDir(root); if err != nil { return InitResult{}, err }
	specsDir, err := cfg.SpecsDir(root); if err != nil { return InitResult{}, err }
	plansDir, err := cfg.PlansDir(root); if err != nil { return InitResult{}, err }
	archiveDir, err := cfg.ArchiveDir(root); if err != nil { return InitResult{}, err }
	templatesDir, err := cfg.TemplatesDir(root); if err != nil { return InitResult{}, err }
	scriptsDir, err := cfg.ScriptsDir(root); if err != nil { return InitResult{}, err }
	subdirs := []string{draftspecDir, specsDir, plansDir, archiveDir, templatesDir, filepath.Join(templatesDir, "prompts"), filepath.Join(templatesDir, "contracts"), filepath.Join(templatesDir, "archive"), scriptsDir}
	for _, dir := range subdirs {
		if err := os.MkdirAll(dir, 0o755); err != nil { return InitResult{}, err }
		messages = append(messages, fmt.Sprintf("ensured directory %s", rel(root, dir)))
	}
	files, err := templates.Files(languages)
	if err != nil { return InitResult{}, err }
	for _, file := range files {
		target := filepath.Join(draftspecDir, file.TargetPath)
		written, err := writeIfMissing(target, file.Content, file.Mode)
		if err != nil { return InitResult{}, err }
		if written { messages = append(messages, fmt.Sprintf("created %s", rel(root, target))) } else { messages = append(messages, fmt.Sprintf("kept existing %s", rel(root, target))) }
	}
	messages = append(messages, fmt.Sprintf("configured languages: docs=%s agent=%s comments=%s", cfg.Language.Docs, cfg.Language.Agent, cfg.Language.Comments))
	agentsPath := filepath.Join(root, "AGENTS.md")
	snippetPath := filepath.Join(templatesDir, "agents-snippet.md")
	changed, err := ensureAgentsSnippet(agentsPath, snippetPath)
	if err != nil { return InitResult{}, err }
	if changed { messages = append(messages, "updated AGENTS.md with Draftspec guidance") } else { messages = append(messages, "kept existing AGENTS.md Draftspec guidance") }
	return InitResult{Messages: messages}, nil
}

func writeIfMissing(path, content string, mode os.FileMode) (bool, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil { return false, err }
	if _, err := os.Stat(path); err == nil { return false, nil } else if !errors.Is(err, os.ErrNotExist) { return false, err }
	return true, os.WriteFile(path, []byte(content), mode)
}

func ensureAgentsSnippet(path, snippetPath string) (bool, error) {
	snippetBytes, err := os.ReadFile(snippetPath)
	if err != nil { return false, err }
	snippet := string(snippetBytes)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) { return true, os.WriteFile(path, []byte(snippet), 0o644) } else if err != nil { return false, err }
	content, err := os.ReadFile(path)
	if err != nil { return false, err }
	if strings.Contains(string(content), "## Draftspec") { return false, nil }
	var builder strings.Builder
	builder.Write(content)
	if len(content) > 0 && !strings.HasSuffix(string(content), "\n") { builder.WriteString("\n") }
	builder.WriteString("\n")
	builder.WriteString(snippet)
	return true, os.WriteFile(path, []byte(builder.String()), 0o644)
}

func rel(root, target string) string {
	relative, err := filepath.Rel(root, target)
	if err != nil { return target }
	return relative
}
