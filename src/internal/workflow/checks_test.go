package workflow

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"draftspec/src/internal/project"
)

func TestInspectSpecValidatesAcceptanceCoverageInGo(t *testing.T) {
	root := t.TempDir()

	_, err := project.Initialize(root, project.InitOptions{
		InitGit:     false,
		DefaultLang: "en",
		Shell:       "sh",
	})
	if err != nil {
		t.Fatalf("Initialize returned error: %v", err)
	}

	specDir := filepath.Join(root, ".draftspec", "specs", "demo")
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatalf("MkdirAll(specDir) returned error: %v", err)
	}
	specPath := filepath.Join(specDir, "spec.md")
	specContent := "# Demo\n\n## Goal\nx\n\n## Requirements\n- RQ-001 x\n\n## Acceptance Criteria\n### AC-001 First\n- Given x\n- When y\n- Then z\n\n### AC-002 Second\n- Given a\n- When b\n- Then c\n"
	if err := os.WriteFile(specPath, []byte(specContent), 0o644); err != nil {
		t.Fatalf("WriteFile(spec) returned error: %v", err)
	}

	tasksPath := filepath.Join(root, ".draftspec", "plans", "demo", "tasks.md")
	if err := os.MkdirAll(filepath.Dir(tasksPath), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	tasksContent := "# Tasks\n\n## Acceptance Coverage\n- AC-001 -> T1.1\n"
	if err := os.WriteFile(tasksPath, []byte(tasksContent), 0o644); err != nil {
		t.Fatalf("WriteFile(tasks) returned error: %v", err)
	}

	result, err := InspectSpec(root, ".draftspec/specs/demo/spec.md", ".draftspec/plans/demo/tasks.md")
	if err != nil {
		t.Fatalf("InspectSpec returned error: %v", err)
	}
	if !result.Failed {
		t.Fatalf("expected InspectSpec to fail, got %+v", result)
	}
	joined := strings.Join(result.Lines, "\n")
	if !strings.Contains(joined, "acceptance coverage entries (1) are fewer than acceptance criteria (2)") {
		t.Fatalf("expected coverage mismatch in output, got %s", joined)
	}
	if !strings.Contains(joined, "SUMMARY: errors=") {
		t.Fatalf("expected summary line in output, got %s", joined)
	}
}

func TestVerifyTaskStateReportsOpenTasksWithoutFailing(t *testing.T) {
	root := t.TempDir()

	_, err := project.Initialize(root, project.InitOptions{
		InitGit:     false,
		DefaultLang: "en",
		Shell:       "sh",
	})
	if err != nil {
		t.Fatalf("Initialize returned error: %v", err)
	}

	tasksPath := filepath.Join(root, ".draftspec", "plans", "demo", "tasks.md")
	if err := os.MkdirAll(filepath.Dir(tasksPath), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	tasksContent := "- [x] T1.1 done\n- [ ] T1.2 open\n\n## Acceptance Coverage\n- AC-001 -> T1.1\n"
	if err := os.WriteFile(tasksPath, []byte(tasksContent), 0o644); err != nil {
		t.Fatalf("WriteFile(tasks) returned error: %v", err)
	}

	result, summary, err := VerifyTaskState(root, "demo")
	if err != nil {
		t.Fatalf("VerifyTaskState returned error: %v", err)
	}
	if result.Failed {
		t.Fatalf("expected open tasks to warn but not fail, got %+v", result)
	}
	if summary.Open != 1 || summary.Total != 2 {
		t.Fatalf("unexpected summary: %+v", summary)
	}
	joined := strings.Join(result.Lines, "\n")
	if !strings.Contains(joined, "TASKS_OPEN=1") || !strings.Contains(joined, "WARN: open tasks remain") {
		t.Fatalf("unexpected verify-task-state output: %s", joined)
	}
}

func TestCheckArchiveReadyBlocksCompletedArchiveWhenTasksRemainOpen(t *testing.T) {
	root := t.TempDir()

	_, err := project.Initialize(root, project.InitOptions{
		InitGit:     false,
		DefaultLang: "en",
		Shell:       "sh",
	})
	if err != nil {
		t.Fatalf("Initialize returned error: %v", err)
	}

	specDir := filepath.Join(root, ".draftspec", "specs", "demo")
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatalf("MkdirAll(specDir) returned error: %v", err)
	}
	specPath := filepath.Join(specDir, "spec.md")
	if err := os.WriteFile(specPath, []byte("# Demo\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(spec) returned error: %v", err)
	}
	tasksPath := filepath.Join(root, ".draftspec", "plans", "demo", "tasks.md")
	if err := os.MkdirAll(filepath.Dir(tasksPath), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(tasksPath, []byte("- [ ] T1.1 open\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(tasks) returned error: %v", err)
	}

	result, err := CheckArchiveReady(root, "demo", "completed", "done")
	if err != nil {
		t.Fatalf("CheckArchiveReady returned error: %v", err)
	}
	if !result.Failed {
		t.Fatalf("expected archive readiness to fail, got %+v", result)
	}
	if !strings.Contains(strings.Join(result.Lines, "\n"), "completed archive requested while open tasks remain") {
		t.Fatalf("unexpected archive output: %+v", result.Lines)
	}
}
