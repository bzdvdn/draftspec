package agents

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type File struct {
	Path    string
	Content string
	Mode    os.FileMode
}

func NormalizeTargets(values []string) ([]string, error) {
	if len(values) == 0 {
		return nil, nil
	}

	seen := map[string]struct{}{}
	var out []string
	for _, value := range values {
		for _, part := range strings.Split(value, ",") {
			target := strings.ToLower(strings.TrimSpace(part))
			if target == "" {
				continue
			}
			if target == "all" {
				for _, candidate := range SupportedTargets() {
					if _, ok := seen[candidate]; ok {
						continue
					}
					seen[candidate] = struct{}{}
					out = append(out, candidate)
				}
				continue
			}
			if _, ok := adapterRegistry[target]; !ok {
				return nil, fmt.Errorf("unsupported agent target %q, expected one of: claude, codex, copilot, cursor, kilocode, trae, all", target)
			}
			if _, ok := seen[target]; ok {
				continue
			}
			seen[target] = struct{}{}
			out = append(out, target)
		}
	}

	sort.Strings(out)
	return out, nil
}

func Files(targets []string, language string, shell string) ([]File, error) {
	normalized, err := NormalizeTargets(targets)
	if err != nil {
		return nil, err
	}

	commands := DefaultCommands(shell)
	var files []File
	for _, target := range normalized {
		adapter, err := adapterForTarget(target)
		if err != nil {
			return nil, err
		}
		targetFiles, err := adapter.Render(commands, language)
		if err != nil {
			return nil, err
		}
		files = append(files, targetFiles...)
	}

	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	return files, nil
}

func FilesForTarget(target, language, shell string) ([]File, error) {
	normalized, err := NormalizeTargets([]string{target})
	if err != nil {
		return nil, err
	}
	if len(normalized) == 0 {
		return nil, nil
	}

	adapter, err := adapterForTarget(normalized[0])
	if err != nil {
		return nil, err
	}
	return adapter.Render(DefaultCommands(shell), language)
}

func PathsForTarget(target string) ([]string, error) {
	adapter, err := adapterForTarget(target)
	if err != nil {
		return nil, err
	}
	paths, err := adapter.Paths(DefaultCommands("sh"), "en")
	if err != nil {
		return nil, err
	}
	sort.Strings(paths)
	return paths, nil
}

// commandSpec and commandSpecs are kept as compatibility shims while tests and
// callers continue to use the previous names.
type commandSpec = CommandDefinition

func commandSpecs(shell string) []commandSpec {
	return DefaultCommands(shell)
}

// render is kept as a narrow compatibility shim for single-command rendering.
func render(target, language string, spec CommandDefinition) (string, string, error) {
	adapter, err := adapterForTarget(target)
	if err != nil {
		return "", "", err
	}
	files, err := adapter.Render([]CommandDefinition{spec}, language)
	if err != nil {
		return "", "", err
	}
	if len(files) != 1 {
		return "", "", fmt.Errorf("expected one rendered file for target %q, got %d", target, len(files))
	}
	return files[0].Path, files[0].Content, nil
}

func normalizeLanguage(language string) string {
	lang := strings.ToLower(strings.TrimSpace(language))
	if lang == "ru" {
		return "ru"
	}
	return "en"
}

func normalizeShell(shell string) string {
	if strings.EqualFold(strings.TrimSpace(shell), "powershell") {
		return "powershell"
	}
	return "sh"
}

func scriptPath(name, shell string) string {
	ext := ".sh"
	if shell == "powershell" {
		ext = ".ps1"
	}
	return ".draftspec/scripts/" + name + ext
}

func commandHint(name, lang string) string {
	if lang == "ru" {
		return fmt.Sprintf("Команда: `/draftspec.%s [request]`", name)
	}
	return fmt.Sprintf("Command: `/draftspec.%s [request]`", name)
}

func toolInvocationHint(lang string) string {
	if lang == "ru" {
		return "Используйте инструменты напрямую через runtime агента; не печатайте raw JSON/XML/tool-call payloads и не выводите внутренние рассуждения о выборе инструмента."
	}
	return "Use tools directly through the agent runtime; do not print raw JSON/XML/tool-call payloads or expose internal reasoning about tool choice."
}

func titleCase(value string) string {
	if value == "" {
		return value
	}
	return strings.ToUpper(value[:1]) + value[1:]
}

func bulletList(items []string) string {
	lines := make([]string, 0, len(items))
	for _, item := range items {
		lines = append(lines, "- "+item)
	}
	return strings.Join(lines, "\n")
}
