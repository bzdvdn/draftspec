package project

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"draftspec/src/internal/config"
)

func TestInitializeCreatesWorkspaceAndAgentTargets(t *testing.T) {
	root := t.TempDir()

	_, err := Initialize(root, InitOptions{
		InitGit:      false,
		DefaultLang:  "en",
		AgentTargets: []string{"claude", "cursor"},
	})
	if err != nil {
		t.Fatalf("Initialize returned error: %v", err)
	}

	cfg, err := config.Load(root)
	if err != nil {
		t.Fatalf("config.Load returned error: %v", err)
	}

	if got, want := cfg.Language.Docs, "en"; got != want {
		t.Fatalf("docs language = %q, want %q", got, want)
	}
	if got, want := strings.Join(cfg.Agents.Targets, ","), "claude,cursor"; got != want {
		t.Fatalf("agent targets = %q, want %q", got, want)
	}

	required := []string{
		filepath.Join(root, ".draftspec", "draftspec.yaml"),
		filepath.Join(root, ".draftspec", "constitution.md"),
		filepath.Join(root, ".draftspec", "memory.md"),
		filepath.Join(root, "AGENTS.md"),
		filepath.Join(root, ".claude", "commands", "draftspec.inspect.md"),
		filepath.Join(root, ".cursor", "rules", "draftspec-inspect.mdc"),
	}
	for _, path := range required {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected %s to exist: %v", path, err)
		}
	}
}

func TestAddRemoveAndCleanupAgents(t *testing.T) {
	root := t.TempDir()

	_, err := Initialize(root, InitOptions{
		InitGit:      false,
		DefaultLang:  "en",
		AgentTargets: []string{"claude"},
	})
	if err != nil {
		t.Fatalf("Initialize returned error: %v", err)
	}

	_, err = AddAgents(root, AddAgentsOptions{Targets: []string{"cursor"}})
	if err != nil {
		t.Fatalf("AddAgents returned error: %v", err)
	}

	cursorPath := filepath.Join(root, ".cursor", "rules", "draftspec-inspect.mdc")
	if _, err := os.Stat(cursorPath); err != nil {
		t.Fatalf("expected cursor agent file after AddAgents: %v", err)
	}

	_, err = RemoveAgents(root, RemoveAgentsOptions{Targets: []string{"cursor"}})
	if err != nil {
		t.Fatalf("RemoveAgents returned error: %v", err)
	}

	if _, err := os.Stat(cursorPath); !os.IsNotExist(err) {
		t.Fatalf("expected cursor agent file to be removed, got err=%v", err)
	}

	if err := os.MkdirAll(filepath.Dir(cursorPath), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(cursorPath, []byte("orphan"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	result, err := CleanupAgents(root)
	if err != nil {
		t.Fatalf("CleanupAgents returned error: %v", err)
	}
	if len(result.Messages) == 0 || !strings.Contains(strings.Join(result.Messages, "\n"), "removed orphaned agent artifact") {
		t.Fatalf("expected cleanup message, got %v", result.Messages)
	}
	if _, err := os.Stat(cursorPath); !os.IsNotExist(err) {
		t.Fatalf("expected orphaned cursor file to be removed, got err=%v", err)
	}

	list, err := ListAgents(root)
	if err != nil {
		t.Fatalf("ListAgents returned error: %v", err)
	}
	if got, want := strings.Join(list.Targets, ","), "claude"; got != want {
		t.Fatalf("enabled targets = %q, want %q", got, want)
	}
}

func TestCleanupAgentsNoop(t *testing.T) {
	root := t.TempDir()

	_, err := Initialize(root, InitOptions{InitGit: false, DefaultLang: "en"})
	if err != nil {
		t.Fatalf("Initialize returned error: %v", err)
	}

	result, err := CleanupAgents(root)
	if err != nil {
		t.Fatalf("CleanupAgents returned error: %v", err)
	}
	if len(result.Messages) != 1 || result.Messages[0] != "no orphaned agent artifacts found" {
		t.Fatalf("unexpected cleanup messages: %v", result.Messages)
	}
}
