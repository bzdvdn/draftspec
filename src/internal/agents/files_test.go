package agents

import (
	"strings"
	"testing"
)

func TestNormalizeTargets(t *testing.T) {
	targets, err := NormalizeTargets([]string{"claude", "cursor,kilocode", "claude", "trae"})
	if err != nil {
		t.Fatalf("NormalizeTargets returned error: %v", err)
	}

	if len(targets) != 4 || targets[0] != "claude" || targets[1] != "cursor" || targets[2] != "kilocode" || targets[3] != "trae" {
		t.Fatalf("unexpected normalized targets: %#v", targets)
	}
}

func TestNormalizeTargetsAll(t *testing.T) {
	targets, err := NormalizeTargets([]string{"all"})
	if err != nil {
		t.Fatalf("NormalizeTargets returned error: %v", err)
	}

	if len(targets) != 6 {
		t.Fatalf("expected 6 targets for all, got %#v", targets)
	}
}

func TestFiles(t *testing.T) {
	files, err := Files([]string{"claude", "codex", "copilot", "cursor", "kilocode", "trae"}, "en", "sh")
	if err != nil {
		t.Fatalf("Files returned error: %v", err)
	}

	if len(files) != 41 {
		t.Fatalf("expected 41 generated agent files, got %d", len(files))
	}

	required := map[string]bool{
		".claude/commands/draftspec.inspect.md":      false,
		".claude/commands/draftspec.verify.md":       false,
		".codex/prompts/draftspec.plan.md":           false,
		".github/prompts/draftspec-spec.prompt.md":   false,
		".github/prompts/draftspec-verify.prompt.md": false,
		".cursor/rules/draftspec-implement.mdc":      false,
		".cursor/rules/draftspec-verify.mdc":         false,
		".kilocode/rules/draftspec-archive.md":       false,
		".kilocode/rules/draftspec-verify.md":        false,
		".trae/project_rules.md":                     false,
	}

	for _, file := range files {
		if _, ok := required[file.Path]; ok {
			required[file.Path] = true
		}
		if file.Content == "" {
			t.Fatalf("expected non-empty content for %s", file.Path)
		}
	}

	for path, found := range required {
		if !found {
			t.Fatalf("missing generated agent file %s", path)
		}
	}
}

func TestRenderEmphasizesRunningScriptsFirst(t *testing.T) {
	spec := commandSpecs("sh")[3] // plan

	tests := []struct {
		name   string
		target string
		lang   string
		want   string
	}{
		{
			name:   "claude en",
			target: "claude",
			lang:   "en",
			want:   "run related scripts first and rely on their output; do not read script source by default",
		},
		{
			name:   "codex ru",
			target: "codex",
			lang:   "ru",
			want:   "сначала запускайте их и опирайтесь на их вывод; не читайте исходники scripts по умолчанию",
		},
		{
			name:   "copilot en",
			target: "copilot",
			lang:   "en",
			want:   "run related scripts first and rely on their output; do not read script source by default",
		},
		{
			name:   "cursor ru",
			target: "cursor",
			lang:   "ru",
			want:   "Если доступны связанные scripts, сначала запускайте их и опирайтесь на их вывод. Не читайте исходники scripts по умолчанию.",
		},
		{
			name:   "kilocode en",
			target: "kilocode",
			lang:   "en",
			want:   "When related scripts are available, run them first and rely on their output. Do not read script source by default.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, content := render(tt.target, tt.lang, spec)
			if !strings.Contains(content, tt.want) {
				t.Fatalf("expected rendered content for %s/%s to contain %q\ncontent:\n%s", tt.target, tt.lang, tt.want, content)
			}
		})
	}
}

func TestRenderTraeEmphasizesRunningScriptsFirst(t *testing.T) {
	tests := []struct {
		name string
		lang string
		want string
	}{
		{
			name: "en",
			lang: "en",
			want: "- When related scripts are available, run them first and rely on their output",
		},
		{
			name: "ru",
			lang: "ru",
			want: "- Если доступны связанные scripts, сначала запускайте их и опирайтесь на их вывод",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := renderTrae(tt.lang, "sh")
			if !strings.Contains(content, tt.want) {
				t.Fatalf("expected trae rules for %s to contain %q\ncontent:\n%s", tt.lang, tt.want, content)
			}
		})
	}
}

func TestRenderIncludesCommandHints(t *testing.T) {
	specs := map[string]commandSpec{}
	for _, spec := range commandSpecs("sh") {
		specs[spec.Name] = spec
	}

	tests := []struct {
		name   string
		target string
		lang   string
		spec   string
		want   string
	}{
		{name: "claude spec en", target: "claude", lang: "en", spec: "spec", want: "Command: `/draftspec.spec [request]`"},
		{name: "codex tasks en", target: "codex", lang: "en", spec: "tasks", want: "Command: `/draftspec.tasks [request]`"},
		{name: "copilot implement ru", target: "copilot", lang: "ru", spec: "implement", want: "Команда: `/draftspec.implement [request]`"},
		{name: "cursor verify en", target: "cursor", lang: "en", spec: "verify", want: "Command: `/draftspec.verify [request]`"},
		{name: "kilocode archive ru", target: "kilocode", lang: "ru", spec: "archive", want: "Команда: `/draftspec.archive [request]`"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, content := render(tt.target, tt.lang, specs[tt.spec])
			if !strings.Contains(content, tt.want) {
				t.Fatalf("expected rendered content for %s/%s/%s to contain %q\ncontent:\n%s", tt.target, tt.lang, tt.spec, tt.want, content)
			}
		})
	}
}

func TestRenderTraeIncludesCommandHints(t *testing.T) {
	tests := []struct {
		name string
		lang string
		want string
	}{
		{name: "en", lang: "en", want: "- Command: `/draftspec.verify [request]`"},
		{name: "ru", lang: "ru", want: "- Команда: `/draftspec.verify [request]`"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := renderTrae(tt.lang, "sh")
			if !strings.Contains(content, tt.want) {
				t.Fatalf("expected trae rules for %s to contain %q\ncontent:\n%s", tt.lang, tt.want, content)
			}
		})
	}
}

func TestRenderCodexDisallowsRawToolPayloads(t *testing.T) {
	specs := map[string]commandSpec{}
	for _, spec := range commandSpecs("sh") {
		specs[spec.Name] = spec
	}

	tests := []struct {
		name string
		lang string
		want string
	}{
		{
			name: "en",
			lang: "en",
			want: "Use tools directly through the agent runtime; do not print raw JSON/XML/tool-call payloads or expose internal reasoning about tool choice.",
		},
		{
			name: "ru",
			lang: "ru",
			want: "Используйте инструменты напрямую через runtime агента; не печатайте raw JSON/XML/tool-call payloads и не выводите внутренние рассуждения о выборе инструмента.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, content := render("codex", tt.lang, specs["plan"])
			if !strings.Contains(content, tt.want) {
				t.Fatalf("expected codex rendered content for %s to contain %q\ncontent:\n%s", tt.lang, tt.want, content)
			}
		})
	}
}
