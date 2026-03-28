package templates

import (
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"draftspec/src/internal/config"
	"gopkg.in/yaml.v3"
)

//go:embed assets/scripts/* assets/lang/* assets/lang/*/* assets/lang/*/templates/* assets/lang/*/templates/prompts/* assets/lang/*/templates/contracts/* assets/lang/*/templates/archive/*
var embedded embed.FS

type File struct {
	TargetPath string
	Content    string
	Mode       fs.FileMode
}

type LanguageSettings struct {
	Default      string
	Docs         string
	Agent        string
	Comments     string
	AgentTargets []string
}

func ResolveLanguageSettings(defaultLang, docsLang, agentLang, commentsLang string) (LanguageSettings, error) {
	base, err := normalizeLanguage(defaultLang)
	if err != nil {
		return LanguageSettings{}, fmt.Errorf("resolve default language: %w", err)
	}
	docs := base
	if strings.TrimSpace(docsLang) != "" {
		docs, err = normalizeLanguage(docsLang)
		if err != nil {
			return LanguageSettings{}, fmt.Errorf("resolve docs language: %w", err)
		}
	}
	agent := base
	if strings.TrimSpace(agentLang) != "" {
		agent, err = normalizeLanguage(agentLang)
		if err != nil {
			return LanguageSettings{}, fmt.Errorf("resolve agent language: %w", err)
		}
	}
	comments := base
	if strings.TrimSpace(commentsLang) != "" {
		comments, err = normalizeLanguage(commentsLang)
		if err != nil {
			return LanguageSettings{}, fmt.Errorf("resolve comments language: %w", err)
		}
	}
	return LanguageSettings{Default: base, Docs: docs, Agent: agent, Comments: comments}, nil
}

func Files(settings LanguageSettings) ([]File, error) {
	definitions := []struct {
		RelativePath string
		TargetPath   string
		Mode         fs.FileMode
		Language     string
	}{
		{RelativePath: "constitution.md", TargetPath: "constitution.md", Mode: 0o644, Language: settings.Docs},
		{RelativePath: "memory.md", TargetPath: "memory.md", Mode: 0o644, Language: settings.Docs},
		{RelativePath: "templates/constitution.md", TargetPath: "templates/constitution.md", Mode: 0o644, Language: settings.Docs},
		{RelativePath: "templates/spec.md", TargetPath: "templates/spec.md", Mode: 0o644, Language: settings.Docs},
		{RelativePath: "templates/plan.md", TargetPath: "templates/plan.md", Mode: 0o644, Language: settings.Docs},
		{RelativePath: "templates/tasks.md", TargetPath: "templates/tasks.md", Mode: 0o644, Language: settings.Docs},
		{RelativePath: "templates/data-model.md", TargetPath: "templates/data-model.md", Mode: 0o644, Language: settings.Docs},
		{RelativePath: "templates/contracts/api.md", TargetPath: "templates/contracts/api.md", Mode: 0o644, Language: settings.Docs},
		{RelativePath: "templates/contracts/events.md", TargetPath: "templates/contracts/events.md", Mode: 0o644, Language: settings.Docs},
		{RelativePath: "templates/archive/summary.md", TargetPath: "templates/archive/summary.md", Mode: 0o644, Language: settings.Docs},
		{RelativePath: "templates/inspect-report.md", TargetPath: "templates/inspect-report.md", Mode: 0o644, Language: settings.Docs},
		{RelativePath: "templates/verify-report.md", TargetPath: "templates/verify-report.md", Mode: 0o644, Language: settings.Docs},
		{RelativePath: "templates/memory.md", TargetPath: "templates/memory.md", Mode: 0o644, Language: settings.Docs},
		{RelativePath: "templates/agents-snippet.md", TargetPath: "templates/agents-snippet.md", Mode: 0o644, Language: settings.Agent},
		{RelativePath: "templates/prompts/constitution.md", TargetPath: "templates/prompts/constitution.md", Mode: 0o644, Language: settings.Agent},
		{RelativePath: "templates/prompts/spec.md", TargetPath: "templates/prompts/spec.md", Mode: 0o644, Language: settings.Agent},
		{RelativePath: "templates/prompts/inspect.md", TargetPath: "templates/prompts/inspect.md", Mode: 0o644, Language: settings.Agent},
		{RelativePath: "templates/prompts/plan.md", TargetPath: "templates/prompts/plan.md", Mode: 0o644, Language: settings.Agent},
		{RelativePath: "templates/prompts/tasks.md", TargetPath: "templates/prompts/tasks.md", Mode: 0o644, Language: settings.Agent},
		{RelativePath: "templates/prompts/implement.md", TargetPath: "templates/prompts/implement.md", Mode: 0o644, Language: settings.Agent},
		{RelativePath: "templates/prompts/archive.md", TargetPath: "templates/prompts/archive.md", Mode: 0o644, Language: settings.Agent},
		{RelativePath: "templates/prompts/verify.md", TargetPath: "templates/prompts/verify.md", Mode: 0o644, Language: settings.Agent},
	}
	files := make([]File, 0, len(definitions)+11)
	files = append(files, File{TargetPath: "draftspec.yaml", Content: generateConfig(settings), Mode: 0o644})
	for _, definition := range definitions {
		content, err := localizedFileContent(definition.Language, definition.RelativePath)
		if err != nil {
			return nil, err
		}
		content = applyLanguagePlaceholders(content, definition.Language, settings)
		files = append(files, File{TargetPath: definition.TargetPath, Content: content, Mode: definition.Mode})
	}
	for _, definition := range []struct {
		AssetPath, TargetPath string
		Mode                  fs.FileMode
	}{
		{"assets/scripts/inspect-spec.sh", "scripts/inspect-spec.sh", 0o755},
		{"assets/scripts/check-constitution.sh", "scripts/check-constitution.sh", 0o755},
		{"assets/scripts/check-spec-ready.sh", "scripts/check-spec-ready.sh", 0o755},
		{"assets/scripts/check-inspect-ready.sh", "scripts/check-inspect-ready.sh", 0o755},
		{"assets/scripts/check-plan-ready.sh", "scripts/check-plan-ready.sh", 0o755},
		{"assets/scripts/check-tasks-ready.sh", "scripts/check-tasks-ready.sh", 0o755},
		{"assets/scripts/check-implement-ready.sh", "scripts/check-implement-ready.sh", 0o755},
		{"assets/scripts/check-archive-ready.sh", "scripts/check-archive-ready.sh", 0o755},
		{"assets/scripts/check-verify-ready.sh", "scripts/check-verify-ready.sh", 0o755},
		{"assets/scripts/verify-task-state.sh", "scripts/verify-task-state.sh", 0o755},
		{"assets/scripts/verify-memory-sync.sh", "scripts/verify-memory-sync.sh", 0o755},
		{"assets/scripts/list-open-tasks.sh", "scripts/list-open-tasks.sh", 0o755},
		{"assets/scripts/sync-memory.sh", "scripts/sync-memory.sh", 0o755},
		{"assets/scripts/link-agents.sh", "scripts/link-agents.sh", 0o755},
		{"assets/scripts/list-specs.sh", "scripts/list-specs.sh", 0o755},
		{"assets/scripts/show-spec.sh", "scripts/show-spec.sh", 0o755},
	} {
		content, err := FileContent(definition.AssetPath)
		if err != nil {
			return nil, err
		}
		files = append(files, File{TargetPath: definition.TargetPath, Content: content, Mode: definition.Mode})
	}
	sort.Slice(files, func(i, j int) bool { return files[i].TargetPath < files[j].TargetPath })
	return files, nil
}

func FileContent(path string) (string, error) {
	content, err := embedded.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read embedded asset %s: %w", path, err)
	}
	return string(content), nil
}

func normalizeLanguage(language string) (string, error) {
	value := strings.ToLower(strings.TrimSpace(language))
	if value == "" {
		value = "en"
	}
	switch value {
	case "en", "ru":
		return value, nil
	default:
		return "", fmt.Errorf("unsupported language %q, expected en or ru", language)
	}
}

func localizedFileContent(language, relativePath string) (string, error) {
	language, err := normalizeLanguage(language)
	if err != nil {
		return "", err
	}
	return FileContent(fmt.Sprintf("assets/lang/%s/%s", language, relativePath))
}

func applyLanguagePlaceholders(content, outputLanguage string, settings LanguageSettings) string {
	return strings.NewReplacer(
		"[DEFAULT_LANGUAGE]", settings.Default,
		"[DOCS_LANGUAGE]", languageLabel(settings.Docs, outputLanguage),
		"[AGENT_LANGUAGE]", languageLabel(settings.Agent, outputLanguage),
		"[COMMENTS_LANGUAGE]", languageLabel(settings.Comments, outputLanguage),
	).Replace(content)
}

func languageLabel(code, outputLanguage string) string {
	switch strings.ToLower(strings.TrimSpace(outputLanguage)) {
	case "ru":
		switch code {
		case "ru":
			return "русский"
		case "en":
			return "английский"
		}
	default:
		switch code {
		case "ru":
			return "Russian"
		case "en":
			return "English"
		}
	}
	return code
}

func generateConfig(settings LanguageSettings) string {
	cfg := config.Default()
	cfg.Language.Default = settings.Default
	cfg.Language.Docs = settings.Docs
	cfg.Language.Agent = settings.Agent
	cfg.Language.Comments = settings.Comments
	cfg.Agents.Targets = settings.AgentTargets
	content, err := yaml.Marshal(cfg)
	if err != nil {
		panic(fmt.Sprintf("marshal generated draftspec config: %v", err))
	}
	return string(content)
}
