package agents

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type File struct {
	Path    string
	Content string
	Mode    os.FileMode
}

var supportedTargets = map[string]struct{}{
	"claude":   {},
	"codex":    {},
	"copilot":  {},
	"cursor":   {},
	"kilocode": {},
	"trae":     {},
}

func SupportedTargets() []string {
	return []string{"claude", "codex", "copilot", "cursor", "kilocode", "trae"}
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
			if _, ok := supportedTargets[target]; !ok {
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

	var files []File
	for _, target := range normalized {
		for _, file := range filesForTarget(target, language, shell) {
			files = append(files, file)
		}
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
	return filesForTarget(normalized[0], language, shell), nil
}

func PathsForTarget(target string) ([]string, error) {
	files, err := FilesForTarget(target, "en", "sh")
	if err != nil {
		return nil, err
	}
	paths := make([]string, 0, len(files))
	for _, file := range files {
		paths = append(paths, file.Path)
	}
	sort.Strings(paths)
	return paths, nil
}

func filesForTarget(target, language, shell string) []File {
	if target == "trae" {
		return []File{{Path: ".trae/project_rules.md", Content: renderTrae(language, shell), Mode: 0o644}}
	}

	var files []File
	for _, command := range commandSpecs(shell) {
		path, content := render(target, language, command)
		files = append(files, File{Path: path, Content: content, Mode: 0o644})
	}
	return files
}

type commandSpec struct {
	Name        string
	Description string
	PromptPath  string
	Extras      []string
}

func commandSpecs(shell string) []commandSpec {
	normalizedShell := normalizeShell(shell)
	launcher := scriptPath("run-draftspec", normalizedShell)
	return []commandSpec{
		{Name: "constitution", Description: "Create or update the project constitution", PromptPath: ".draftspec/templates/prompts/constitution.md", Extras: []string{launcher, scriptPath("check-constitution", normalizedShell)}},
		{Name: "spec", Description: "Create or update one feature spec", PromptPath: ".draftspec/templates/prompts/spec.md", Extras: []string{launcher, scriptPath("check-spec-ready", normalizedShell)}},
		{Name: "inspect", Description: "Inspect one feature for consistency and quality", PromptPath: ".draftspec/templates/prompts/inspect.md", Extras: []string{launcher, scriptPath("check-inspect-ready", normalizedShell), scriptPath("inspect-spec", normalizedShell)}},
		{Name: "plan", Description: "Create or update one feature plan package", PromptPath: ".draftspec/templates/prompts/plan.md", Extras: []string{launcher, scriptPath("check-plan-ready", normalizedShell)}},
		{Name: "tasks", Description: "Create or update tasks for one feature", PromptPath: ".draftspec/templates/prompts/tasks.md", Extras: []string{launcher, scriptPath("check-tasks-ready", normalizedShell)}},
		{Name: "implement", Description: "Implement one feature from tasks", PromptPath: ".draftspec/templates/prompts/implement.md", Extras: []string{launcher, scriptPath("check-implement-ready", normalizedShell), scriptPath("list-open-tasks", normalizedShell)}},
		{Name: "verify", Description: "Verify one implemented feature package", PromptPath: ".draftspec/templates/prompts/verify.md", Extras: []string{launcher, scriptPath("check-verify-ready", normalizedShell), scriptPath("verify-task-state", normalizedShell)}},
		{Name: "archive", Description: "Archive one feature package", PromptPath: ".draftspec/templates/prompts/archive.md", Extras: []string{launcher, scriptPath("check-archive-ready", normalizedShell)}},
	}
}

func scriptPath(name, shell string) string {
	ext := ".sh"
	if shell == "powershell" {
		ext = ".ps1"
	}
	return ".draftspec/scripts/" + name + ext
}

func render(target, language string, spec commandSpec) (string, string) {
	lang := normalizeLanguage(language)

	switch target {
	case "claude":
		return filepath.ToSlash(filepath.Join(".claude", "commands", fmt.Sprintf("draftspec.%s.md", spec.Name))), renderClaude(spec, lang)
	case "codex":
		return filepath.ToSlash(filepath.Join(".codex", "prompts", fmt.Sprintf("draftspec.%s.md", spec.Name))), renderCodex(spec, lang)
	case "copilot":
		return filepath.ToSlash(filepath.Join(".github", "prompts", fmt.Sprintf("draftspec-%s.prompt.md", spec.Name))), renderCopilot(spec, lang)
	case "cursor":
		return filepath.ToSlash(filepath.Join(".cursor", "rules", fmt.Sprintf("draftspec-%s.mdc", spec.Name))), renderCursor(spec, lang)
	case "kilocode":
		return filepath.ToSlash(filepath.Join(".kilocode", "rules", fmt.Sprintf("draftspec-%s.md", spec.Name))), renderKilo(spec, lang)
	default:
		panic("unsupported target")
	}
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

func renderClaude(spec commandSpec, lang string) string {
	if lang == "ru" {
		return fmt.Sprintf(`---
description: %s
argument-hint: [request]
---

Следуйте файлу %q.

%s

Аргументы пользователя:
$ARGUMENTS

Требования:
- сначала прочитайте .draftspec/constitution.md, если это требуется prompt-файлом
- используйте только минимально нужный контекст репозитория
- если доступны, сначала запускайте связанные scripts и опирайтесь на их вывод; не читайте исходники scripts по умолчанию:
%s
- обновляйте только релевантные артефакты и кратко сообщайте об итогах и блокерах
`, spec.Description, spec.PromptPath, commandHint(spec.Name, lang), bulletList(spec.Extras))
	}

	return fmt.Sprintf(`---
description: %s
argument-hint: [request]
---

Follow %q.

%s

User arguments:
$ARGUMENTS

Requirements:
- read .draftspec/constitution.md first when the prompt requires it
- use only the minimum repository context needed
- when available, run related scripts first and rely on their output; do not read script source by default:
%s
- update only the relevant artifacts and report outcomes and blockers briefly
`, spec.Description, spec.PromptPath, commandHint(spec.Name, lang), bulletList(spec.Extras))
}

func renderCodex(spec commandSpec, lang string) string {
	title := titleCase(spec.Name)
	if lang == "ru" {
		return fmt.Sprintf(`# Draftspec %s

Следуйте файлу %q.

%s

Вход пользователя: {{arguments}}

Дополнительно:
- если доступны связанные scripts, сначала запускайте их и опирайтесь на их вывод; не читайте исходники scripts по умолчанию
- %s
%s
`, title, spec.PromptPath, commandHint(spec.Name, lang), toolInvocationHint(lang), bulletList(spec.Extras))
	}

	return fmt.Sprintf(`# Draftspec %s

Follow %q.

%s

User input: {{arguments}}

Additional context:
- when related scripts are available, run them first and rely on their output; do not read script source by default
- %s
%s
`, title, spec.PromptPath, commandHint(spec.Name, lang), toolInvocationHint(lang), bulletList(spec.Extras))
}

func renderCopilot(spec commandSpec, lang string) string {
	if lang == "ru" {
		return fmt.Sprintf(`# Draftspec %s

Используйте %q как основной workflow prompt.

%s

Что нужно сделать:
- обработать запрос пользователя для одной фазы %q
- применять только минимально нужный контекст репозитория
- при необходимости сначала запускайте связанные scripts и опирайтесь на их вывод; не читайте исходники scripts по умолчанию:
%s
- кратко сообщить о результатах и блокерах
`, spec.Name, spec.PromptPath, commandHint(spec.Name, lang), spec.Name, bulletList(spec.Extras))
	}

	return fmt.Sprintf(`# Draftspec %s

Use %q as the primary workflow prompt.

%s

What to do:
- handle the user request for the %q phase
- use only the minimum repository context required
- when needed, run related scripts first and rely on their output; do not read script source by default:
%s
- report outcomes and blockers briefly
`, spec.Name, spec.PromptPath, commandHint(spec.Name, lang), spec.Name, bulletList(spec.Extras))
}

func renderCursor(spec commandSpec, lang string) string {
	if lang == "ru" {
		return fmt.Sprintf(`---
description: Draftspec %s workflow
alwaysApply: false
---

Следуйте файлу %q.

%s

Используйте эту rule, когда запрос явно относится к фазе %q или к команде /draftspec.%s.

Если доступны связанные scripts, сначала запускайте их и опирайтесь на их вывод. Не читайте исходники scripts по умолчанию.

Связанные scripts:
%s
`, spec.Name, spec.PromptPath, commandHint(spec.Name, lang), spec.Name, spec.Name, bulletList(spec.Extras))
	}

	return fmt.Sprintf(`---
description: Draftspec %s workflow
alwaysApply: false
---

Follow %q.

%s

Use this rule when the request clearly maps to the %q phase or the /draftspec.%s command.

When related scripts are available, run them first and rely on their output. Do not read script source by default.

Related scripts:
%s
`, spec.Name, spec.PromptPath, commandHint(spec.Name, lang), spec.Name, spec.Name, bulletList(spec.Extras))
}

func renderKilo(spec commandSpec, lang string) string {
	if lang == "ru" {
		return fmt.Sprintf(`# Draftspec %s

Следуйте файлу %q.

%s

Используйте это project rule, когда запрос относится к фазе %q.

Если доступны связанные scripts, сначала запускайте их и опирайтесь на их вывод. Не читайте исходники scripts по умолчанию.

Связанные scripts:
%s
`, spec.Name, spec.PromptPath, commandHint(spec.Name, lang), spec.Name, bulletList(spec.Extras))
	}

	return fmt.Sprintf(`# Draftspec %s

Follow %q.

%s

Use this project rule when the request maps to the %q phase.

When related scripts are available, run them first and rely on their output. Do not read script source by default.

Related scripts:
%s
`, spec.Name, spec.PromptPath, commandHint(spec.Name, lang), spec.Name, bulletList(spec.Extras))
}

func renderTrae(language, shell string) string {
	lang := normalizeLanguage(language)
	if lang == "ru" {
		var sections []string
		sections = append(sections, "# Draftspec Project Rules")
		sections = append(sections, "")
		sections = append(sections, "Используйте .draftspec как основной источник проектного контекста. Следуйте AGENTS.md и соответствующим prompt-файлам в .draftspec/templates/prompts/.")
		for _, spec := range commandSpecs(shell) {
			sections = append(sections, "")
			sections = append(sections, fmt.Sprintf("## /draftspec.%s", spec.Name))
			sections = append(sections, fmt.Sprintf("- Основной prompt: %s", spec.PromptPath))
			sections = append(sections, fmt.Sprintf("- %s", commandHint(spec.Name, lang)))
			sections = append(sections, "- Используйте только минимально нужный контекст репозитория")
			sections = append(sections, "- Если доступны связанные scripts, сначала запускайте их и опирайтесь на их вывод")
			sections = append(sections, "- Не читайте исходники scripts по умолчанию")
			sections = append(sections, "- Связанные scripts:")
			sections = append(sections, bulletList(spec.Extras))
		}
		return strings.Join(sections, "\n") + "\n"
	}

	var sections []string
	sections = append(sections, "# Draftspec Project Rules")
	sections = append(sections, "")
	sections = append(sections, "Use .draftspec as the primary source of project context. Follow AGENTS.md and the matching prompt files under .draftspec/templates/prompts/.")
	for _, spec := range commandSpecs(shell) {
		sections = append(sections, "")
		sections = append(sections, fmt.Sprintf("## /draftspec.%s", spec.Name))
		sections = append(sections, fmt.Sprintf("- Primary prompt: %s", spec.PromptPath))
		sections = append(sections, fmt.Sprintf("- %s", commandHint(spec.Name, lang)))
		sections = append(sections, "- Use only the minimum repository context required")
		sections = append(sections, "- When related scripts are available, run them first and rely on their output")
		sections = append(sections, "- Do not read script source by default")
		sections = append(sections, "- Related scripts:")
		sections = append(sections, bulletList(spec.Extras))
	}
	return strings.Join(sections, "\n") + "\n"
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
