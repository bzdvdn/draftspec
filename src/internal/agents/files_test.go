package agents

import "testing"

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
	files, err := Files([]string{"claude", "codex", "copilot", "cursor", "kilocode", "trae"}, "en")
	if err != nil {
		t.Fatalf("Files returned error: %v", err)
	}

	if len(files) != 36 {
		t.Fatalf("expected 36 generated agent files, got %d", len(files))
	}

	required := map[string]bool{
		".claude/commands/draftspec.inspect.md":    false,
		".codex/prompts/draftspec.plan.md":         false,
		".github/prompts/draftspec-spec.prompt.md": false,
		".cursor/rules/draftspec-implement.mdc":    false,
		".kilocode/rules/draftspec-archive.md":     false,
		".trae/project_rules.md":                   false,
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
