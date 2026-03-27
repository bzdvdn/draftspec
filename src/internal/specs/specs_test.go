package specs

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"draftspec/src/internal/project"
)

func TestListReturnsSortedMarkdownSpecsOnly(t *testing.T) {
	root := t.TempDir()

	_, err := project.Initialize(root, project.InitOptions{InitGit: false, DefaultLang: "en"})
	if err != nil {
		t.Fatalf("Initialize returned error: %v", err)
	}

	specsDir := filepath.Join(root, ".draftspec", "specs")
	if err := os.WriteFile(filepath.Join(specsDir, "zeta.md"), []byte("# Zeta"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(specsDir, "alpha.md"), []byte("# Alpha"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(specsDir, "notes.txt"), []byte("ignore"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
	if err := os.Mkdir(filepath.Join(specsDir, "nested"), 0o755); err != nil {
		t.Fatalf("Mkdir returned error: %v", err)
	}

	got, err := List(root)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	want := []string{"alpha", "zeta"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("List = %#v, want %#v", got, want)
	}
}

func TestShowReturnsSpecContent(t *testing.T) {
	root := t.TempDir()

	_, err := project.Initialize(root, project.InitOptions{InitGit: false, DefaultLang: "en"})
	if err != nil {
		t.Fatalf("Initialize returned error: %v", err)
	}

	specPath := filepath.Join(root, ".draftspec", "specs", "demo.md")
	content := "# Demo\n\nHello"
	if err := os.WriteFile(specPath, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	got, err := Show(root, "demo")
	if err != nil {
		t.Fatalf("Show returned error: %v", err)
	}
	if got != content {
		t.Fatalf("Show = %q, want %q", got, content)
	}
}

func TestCreateGeneratesSpecAndTasksFromTemplates(t *testing.T) {
	root := t.TempDir()

	_, err := project.Initialize(root, project.InitOptions{InitGit: false, DefaultLang: "en"})
	if err != nil {
		t.Fatalf("Initialize returned error: %v", err)
	}

	result, err := Create(root, "Partner Scheduling", CreateOptions{CreateBranch: false})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if len(result.Messages) == 0 || result.Messages[0] != "skipped spec branch creation" {
		t.Fatalf("unexpected messages: %v", result.Messages)
	}

	specPath := filepath.Join(root, ".draftspec", "specs", "partner-scheduling.md")
	tasksPath := filepath.Join(root, ".draftspec", "plans", "partner-scheduling", "tasks.md")

	specContent, err := os.ReadFile(specPath)
	if err != nil {
		t.Fatalf("ReadFile spec returned error: %v", err)
	}
	tasksContent, err := os.ReadFile(tasksPath)
	if err != nil {
		t.Fatalf("ReadFile tasks returned error: %v", err)
	}

	if !strings.Contains(string(specContent), "Partner Scheduling") {
		t.Fatalf("expected spec to contain filled title, got: %s", string(specContent))
	}
	if !strings.Contains(string(tasksContent), "Partner Scheduling") {
		t.Fatalf("expected tasks to contain filled title, got: %s", string(tasksContent))
	}
}

func TestCreateFailsOnEmptySlug(t *testing.T) {
	root := t.TempDir()

	_, err := project.Initialize(root, project.InitOptions{InitGit: false, DefaultLang: "en"})
	if err != nil {
		t.Fatalf("Initialize returned error: %v", err)
	}

	_, err = Create(root, "---", CreateOptions{CreateBranch: false})
	if err == nil {
		t.Fatal("expected error for empty slug, got nil")
	}
}
