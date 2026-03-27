package project

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"draftspec/src/internal/agents"
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
	AgentTargets []string
}

type InitResult struct{ Messages []string }

type AddAgentsOptions struct {
	Targets   []string
	AgentLang string
}

type AddAgentsResult struct{ Messages []string }

type RemoveAgentsOptions struct {
	Targets []string
}

type RemoveAgentsResult struct{ Messages []string }

type ListAgentsResult struct {
	Targets []string
}

type CleanupAgentsResult struct{ Messages []string }

func Initialize(root string, options InitOptions) (InitResult, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return InitResult{}, err
	}
	languages, err := templates.ResolveLanguageSettings(options.DefaultLang, options.DocsLang, options.AgentLang, options.CommentsLang)
	if err != nil {
		return InitResult{}, err
	}
	normalizedAgentTargets, err := agents.NormalizeTargets(options.AgentTargets)
	if err != nil {
		return InitResult{}, err
	}
	languages.AgentTargets = normalizedAgentTargets
	cfg := config.Default()
	cfg.Language.Default = languages.Default
	cfg.Language.Docs = languages.Docs
	cfg.Language.Agent = languages.Agent
	cfg.Language.Comments = languages.Comments
	cfg.Agents.Targets = normalizedAgentTargets
	var messages []string
	if options.InitGit {
		created, err := gitutil.EnsureRepository(root)
		if err != nil {
			return InitResult{}, err
		}
		if created {
			messages = append(messages, "initialized git repository")
		} else {
			messages = append(messages, "kept existing git repository")
		}
	} else {
		messages = append(messages, "skipped git repository initialization")
	}
	draftspecDir, err := cfg.DraftspecDir(root)
	if err != nil {
		return InitResult{}, err
	}
	specsDir, err := cfg.SpecsDir(root)
	if err != nil {
		return InitResult{}, err
	}
	plansDir, err := cfg.PlansDir(root)
	if err != nil {
		return InitResult{}, err
	}
	archiveDir, err := cfg.ArchiveDir(root)
	if err != nil {
		return InitResult{}, err
	}
	templatesDir, err := cfg.TemplatesDir(root)
	if err != nil {
		return InitResult{}, err
	}
	scriptsDir, err := cfg.ScriptsDir(root)
	if err != nil {
		return InitResult{}, err
	}
	subdirs := []string{draftspecDir, specsDir, plansDir, archiveDir, templatesDir, filepath.Join(templatesDir, "prompts"), filepath.Join(templatesDir, "contracts"), filepath.Join(templatesDir, "archive"), scriptsDir}
	for _, dir := range subdirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return InitResult{}, err
		}
		messages = append(messages, fmt.Sprintf("ensured directory %s", rel(root, dir)))
	}
	files, err := templates.Files(languages)
	if err != nil {
		return InitResult{}, err
	}
	for _, file := range files {
		target := filepath.Join(draftspecDir, file.TargetPath)
		written, err := writeIfMissing(target, file.Content, file.Mode)
		if err != nil {
			return InitResult{}, err
		}
		if written {
			messages = append(messages, fmt.Sprintf("created %s", rel(root, target)))
		} else {
			messages = append(messages, fmt.Sprintf("kept existing %s", rel(root, target)))
		}
	}
	messages = append(messages, fmt.Sprintf("configured languages: docs=%s agent=%s comments=%s", cfg.Language.Docs, cfg.Language.Agent, cfg.Language.Comments))
	agentsPath := filepath.Join(root, "AGENTS.md")
	snippetPath := filepath.Join(templatesDir, "agents-snippet.md")
	changed, err := ensureAgentsSnippet(agentsPath, snippetPath)
	if err != nil {
		return InitResult{}, err
	}
	if changed {
		messages = append(messages, "updated AGENTS.md with Draftspec guidance")
	} else {
		messages = append(messages, "kept existing AGENTS.md Draftspec guidance")
	}
	for _, message := range ensureAgentFiles(root, normalizedAgentTargets, languages.Agent) {
		messages = append(messages, message)
	}
	if len(normalizedAgentTargets) > 0 {
		messages = append(messages, fmt.Sprintf("enabled agent targets: %s", strings.Join(normalizedAgentTargets, ", ")))
	} else {
		messages = append(messages, "enabled agent targets: none")
	}
	return InitResult{Messages: messages}, nil
}

func AddAgents(root string, options AddAgentsOptions) (AddAgentsResult, error) {
	root, cfg, err := loadInitializedProject(root)
	if err != nil {
		return AddAgentsResult{}, err
	}

	requested, err := agents.NormalizeTargets(options.Targets)
	if err != nil {
		return AddAgentsResult{}, err
	}
	combined, err := agents.NormalizeTargets(append(cfg.Agents.Targets, requested...))
	if err != nil {
		return AddAgentsResult{}, err
	}

	agentLanguage := cfg.Language.Agent
	if strings.TrimSpace(options.AgentLang) != "" {
		agentLanguage = strings.TrimSpace(options.AgentLang)
	}

	cfg.Agents.Targets = combined
	if err := config.Save(root, cfg); err != nil {
		return AddAgentsResult{}, err
	}

	messages := []string{"updated .draftspec/draftspec.yaml with agent targets"}
	messages = append(messages, ensureAgentFiles(root, requested, agentLanguage)...)
	messages = append(messages, fmt.Sprintf("enabled agent targets: %s", strings.Join(combined, ", ")))
	return AddAgentsResult{Messages: messages}, nil
}

func RemoveAgents(root string, options RemoveAgentsOptions) (RemoveAgentsResult, error) {
	root, cfg, err := loadInitializedProject(root)
	if err != nil {
		return RemoveAgentsResult{}, err
	}

	requested, err := agents.NormalizeTargets(options.Targets)
	if err != nil {
		return RemoveAgentsResult{}, err
	}
	removeSet := make(map[string]struct{}, len(requested))
	for _, target := range requested {
		removeSet[target] = struct{}{}
	}

	var remaining []string
	for _, target := range cfg.Agents.Targets {
		if _, ok := removeSet[target]; ok {
			continue
		}
		remaining = append(remaining, target)
	}
	sort.Strings(remaining)
	cfg.Agents.Targets = remaining
	if err := config.Save(root, cfg); err != nil {
		return RemoveAgentsResult{}, err
	}

	messages := []string{"updated .draftspec/draftspec.yaml with agent targets"}
	messages = append(messages, removeAgentFiles(root, requested)...)
	if len(remaining) > 0 {
		messages = append(messages, fmt.Sprintf("enabled agent targets: %s", strings.Join(remaining, ", ")))
	} else {
		messages = append(messages, "enabled agent targets: none")
	}
	return RemoveAgentsResult{Messages: messages}, nil
}

func ListAgents(root string) (ListAgentsResult, error) {
	_, cfg, err := loadInitializedProject(root)
	if err != nil {
		return ListAgentsResult{}, err
	}
	return ListAgentsResult{Targets: append([]string(nil), cfg.Agents.Targets...)}, nil
}

func CleanupAgents(root string) (CleanupAgentsResult, error) {
	root, cfg, err := loadInitializedProject(root)
	if err != nil {
		return CleanupAgentsResult{}, err
	}

	enabledTargets := make(map[string]struct{}, len(cfg.Agents.Targets))
	for _, target := range cfg.Agents.Targets {
		enabledTargets[target] = struct{}{}
	}

	var messages []string
	removedAny := false
	for _, target := range agents.SupportedTargets() {
		if _, ok := enabledTargets[target]; ok {
			continue
		}
		paths, err := agents.PathsForTarget(target)
		if err != nil {
			return CleanupAgentsResult{}, err
		}
		for _, relPath := range paths {
			fullPath := filepath.Join(root, filepath.FromSlash(relPath))
			if _, err := os.Stat(fullPath); errors.Is(err, os.ErrNotExist) {
				continue
			} else if err != nil {
				return CleanupAgentsResult{}, err
			}
			if err := os.Remove(fullPath); err != nil {
				return CleanupAgentsResult{}, err
			}
			messages = append(messages, fmt.Sprintf("removed orphaned agent artifact %s", rel(root, fullPath)))
			removedAny = true
		}
	}

	if !removedAny {
		messages = append(messages, "no orphaned agent artifacts found")
	}

	return CleanupAgentsResult{Messages: messages}, nil
}

func writeIfMissing(path, content string, mode os.FileMode) (bool, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return false, err
	}
	if _, err := os.Stat(path); err == nil {
		return false, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return false, err
	}
	return true, os.WriteFile(path, []byte(content), mode)
}

func ensureAgentsSnippet(path, snippetPath string) (bool, error) {
	snippetBytes, err := os.ReadFile(snippetPath)
	if err != nil {
		return false, err
	}
	snippet := string(snippetBytes)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return true, os.WriteFile(path, []byte(snippet), 0o644)
	} else if err != nil {
		return false, err
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	if strings.Contains(string(content), "## Draftspec") {
		return false, nil
	}
	var builder strings.Builder
	builder.Write(content)
	if len(content) > 0 && !strings.HasSuffix(string(content), "\n") {
		builder.WriteString("\n")
	}
	builder.WriteString("\n")
	builder.WriteString(snippet)
	return true, os.WriteFile(path, []byte(builder.String()), 0o644)
}

func ensureAgentFiles(root string, targets []string, language string) []string {
	agentFiles, err := agents.Files(targets, language)
	if err != nil {
		return []string{fmt.Sprintf("skipped agent files: %v", err)}
	}
	messages := make([]string, 0, len(agentFiles))
	for _, file := range agentFiles {
		target := filepath.Join(root, filepath.FromSlash(file.Path))
		written, err := writeIfMissing(target, file.Content, file.Mode)
		if err != nil {
			messages = append(messages, fmt.Sprintf("failed %s: %v", rel(root, target), err))
			continue
		}
		if written {
			messages = append(messages, fmt.Sprintf("created %s", rel(root, target)))
		} else {
			messages = append(messages, fmt.Sprintf("kept existing %s", rel(root, target)))
		}
	}
	return messages
}

func removeAgentFiles(root string, targets []string) []string {
	messages := []string{}
	for _, target := range targets {
		paths, err := agents.PathsForTarget(target)
		if err != nil {
			messages = append(messages, fmt.Sprintf("skipped removing %s: %v", target, err))
			continue
		}
		for _, relPath := range paths {
			fullPath := filepath.Join(root, filepath.FromSlash(relPath))
			if err := os.Remove(fullPath); err != nil {
				if errors.Is(err, os.ErrNotExist) {
					messages = append(messages, fmt.Sprintf("missing %s", rel(root, fullPath)))
					continue
				}
				messages = append(messages, fmt.Sprintf("failed %s: %v", rel(root, fullPath), err))
				continue
			}
			messages = append(messages, fmt.Sprintf("removed %s", rel(root, fullPath)))
		}
	}
	return messages
}

func loadInitializedProject(root string) (string, config.Config, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return "", config.Config{}, err
	}
	cfgPath := filepath.Join(absRoot, ".draftspec", "draftspec.yaml")
	if _, err := os.Stat(cfgPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", config.Config{}, fmt.Errorf("draftspec project is not initialized at %s", absRoot)
		}
		return "", config.Config{}, err
	}
	cfg, err := config.Load(absRoot)
	if err != nil {
		return "", config.Config{}, err
	}
	return absRoot, cfg, nil
}

func rel(root, target string) string {
	relative, err := filepath.Rel(root, target)
	if err != nil {
		return target
	}
	return relative
}
