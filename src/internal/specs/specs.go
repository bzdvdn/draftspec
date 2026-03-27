package specs

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"draftspec/src/internal/config"
	"draftspec/src/internal/gitutil"
)

type CreateOptions struct {
	CreateBranch bool
}

type CreateResult struct {
	Messages []string
}

func List(root string) ([]string, error) {
	cfg, err := config.Load(root)
	if err != nil {
		return nil, err
	}

	specsDir, err := cfg.SpecsDir(root)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(specsDir)
	if err != nil {
		return nil, fmt.Errorf("read specs directory: %w", err)
	}

	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if filepath.Ext(name) != ".md" {
			continue
		}

		names = append(names, strings.TrimSuffix(name, ".md"))
	}

	sort.Strings(names)
	return names, nil
}

func Show(root, name string) (string, error) {
	cfg, err := config.Load(root)
	if err != nil {
		return "", err
	}

	specsDir, err := cfg.SpecsDir(root)
	if err != nil {
		return "", err
	}

	specPath := filepath.Join(specsDir, name+".md")
	content, err := os.ReadFile(specPath)
	if err != nil {
		return "", fmt.Errorf("read spec %q: %w", name, err)
	}

	return string(content), nil
}

func Create(root, name string, options CreateOptions) (CreateResult, error) {
	slug := slugify(name)
	if slug == "" {
		return CreateResult{}, fmt.Errorf("spec name %q produced an empty slug", name)
	}

	cfg, err := config.Load(root)
	if err != nil {
		return CreateResult{}, err
	}

	var messages []string
	if options.CreateBranch {
		message, err := gitutil.EnsureBranch(root, "spec/"+slug)
		if err != nil {
			return CreateResult{}, err
		}
		messages = append(messages, message)
	} else {
		messages = append(messages, "skipped spec branch creation")
	}

	title := titleFromSlug(slug)
	templatesDir, err := cfg.TemplatesDir(root)
	if err != nil {
		return CreateResult{}, err
	}
	specTemplatePath := filepath.Join(templatesDir, "spec.md")
	tasksTemplatePath := filepath.Join(templatesDir, "tasks.md")
	plansDir, err := cfg.PlansDir(root)
	if err != nil {
		return CreateResult{}, err
	}
	specsDir, err := cfg.SpecsDir(root)
	if err != nil {
		return CreateResult{}, err
	}

	specTemplate, err := os.ReadFile(specTemplatePath)
	if err != nil {
		return CreateResult{}, fmt.Errorf("read spec template: %w", err)
	}
	tasksTemplate, err := os.ReadFile(tasksTemplatePath)
	if err != nil {
		return CreateResult{}, fmt.Errorf("read tasks template: %w", err)
	}

	specPath := filepath.Join(specsDir, slug+".md")
	tasksPath := filepath.Join(plansDir, slug, "tasks.md")

	created, err := writeFilledTemplate(specPath, string(specTemplate), title)
	if err != nil {
		return CreateResult{}, err
	}
	if created {
		messages = append(messages, fmt.Sprintf("created %s", displayPath(root, specPath)))
	} else {
		messages = append(messages, fmt.Sprintf("kept existing %s", displayPath(root, specPath)))
	}

	created, err = writeFilledTemplate(tasksPath, string(tasksTemplate), title)
	if err != nil {
		return CreateResult{}, err
	}
	if created {
		messages = append(messages, fmt.Sprintf("created %s", displayPath(root, tasksPath)))
	} else {
		messages = append(messages, fmt.Sprintf("kept existing %s", displayPath(root, tasksPath)))
	}

	return CreateResult{Messages: messages}, nil
}

func writeFilledTemplate(path, templateContent, title string) (bool, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return false, err
	}

	if _, err := os.Stat(path); err == nil {
		return false, nil
	} else if !os.IsNotExist(err) {
		return false, err
	}

	content := strings.ReplaceAll(templateContent, "<Spec Title>", title)
	return true, os.WriteFile(path, []byte(content), 0o644)
}

func displayPath(root, target string) string {
	relative, err := filepath.Rel(root, target)
	if err != nil {
		return target
	}
	return filepath.ToSlash(relative)
}

func slugify(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	var b strings.Builder
	lastDash := false
	for _, r := range name {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
			lastDash = false
		case r == ' ' || r == '-' || r == '_' || r == '/':
			if b.Len() > 0 && !lastDash {
				b.WriteByte('-')
				lastDash = true
			}
		}
	}
	return strings.Trim(b.String(), "-")
}

func titleFromSlug(slug string) string {
	parts := strings.Split(slug, "-")
	for i, part := range parts {
		if part == "" {
			continue
		}
		runes := []rune(part)
		runes[0] = unicode.ToUpper(runes[0])
		parts[i] = string(runes)
	}
	return strings.Join(parts, " ")
}
