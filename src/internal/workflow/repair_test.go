package workflow

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"draftspec/src/internal/project"
)

func TestRepairFeatureMovesLegacyInspectReportToCanonicalPath(t *testing.T) {
	root := t.TempDir()

	_, err := project.Initialize(root, project.InitOptions{
		InitGit:     false,
		DefaultLang: "en",
		Shell:       "sh",
	})
	if err != nil {
		t.Fatalf("Initialize returned error: %v", err)
	}

	legacyPath := filepath.Join(root, ".draftspec", "plans", "demo", "inspect.md")
	if err := os.MkdirAll(filepath.Dir(legacyPath), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	content := "---\nreport_type: inspect\nslug: demo\nstatus: pass\ndocs_language: en\ngenerated_at: 2026-03-31\n---\n# Inspect Report: demo\n"
	if err := os.WriteFile(legacyPath, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	result, err := RepairFeature(root, "demo", false)
	if err != nil {
		t.Fatalf("RepairFeature returned error: %v", err)
	}
	if !result.Changed {
		t.Fatalf("expected repair to change files, got %+v", result)
	}

	canonicalPath := filepath.Join(root, ".draftspec", "specs", "demo.inspect.md")
	if _, err := os.Stat(canonicalPath); err != nil {
		t.Fatalf("expected canonical inspect report to exist: %v", err)
	}
	if _, err := os.Stat(legacyPath); !os.IsNotExist(err) {
		t.Fatalf("expected legacy inspect report to be moved, got err=%v", err)
	}
}

func TestRepairFeatureRemovesDuplicateLegacyInspectReport(t *testing.T) {
	root := t.TempDir()

	_, err := project.Initialize(root, project.InitOptions{
		InitGit:     false,
		DefaultLang: "en",
		Shell:       "sh",
	})
	if err != nil {
		t.Fatalf("Initialize returned error: %v", err)
	}

	content := "---\nreport_type: inspect\nslug: demo\nstatus: pass\ndocs_language: en\ngenerated_at: 2026-03-31\n---\n# Inspect Report: demo\n"
	canonicalPath := filepath.Join(root, ".draftspec", "specs", "demo.inspect.md")
	if err := os.WriteFile(canonicalPath, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(canonical) returned error: %v", err)
	}
	legacyPath := filepath.Join(root, ".draftspec", "plans", "demo", "inspect.md")
	if err := os.MkdirAll(filepath.Dir(legacyPath), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(legacyPath, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(legacy) returned error: %v", err)
	}

	result, err := RepairFeature(root, "demo", false)
	if err != nil {
		t.Fatalf("RepairFeature returned error: %v", err)
	}
	if !result.Changed {
		t.Fatalf("expected duplicate cleanup to change files, got %+v", result)
	}
	if _, err := os.Stat(legacyPath); !os.IsNotExist(err) {
		t.Fatalf("expected duplicate legacy inspect report to be removed, got err=%v", err)
	}
}

func TestRepairFeatureWarnsWhenCanonicalAndLegacyDiffer(t *testing.T) {
	root := t.TempDir()

	_, err := project.Initialize(root, project.InitOptions{
		InitGit:     false,
		DefaultLang: "en",
		Shell:       "sh",
	})
	if err != nil {
		t.Fatalf("Initialize returned error: %v", err)
	}

	canonicalPath := filepath.Join(root, ".draftspec", "specs", "demo.inspect.md")
	if err := os.WriteFile(canonicalPath, []byte("canonical"), 0o644); err != nil {
		t.Fatalf("WriteFile(canonical) returned error: %v", err)
	}
	legacyPath := filepath.Join(root, ".draftspec", "plans", "demo", "inspect.md")
	if err := os.MkdirAll(filepath.Dir(legacyPath), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(legacyPath, []byte("legacy"), 0o644); err != nil {
		t.Fatalf("WriteFile(legacy) returned error: %v", err)
	}

	result, err := RepairFeature(root, "demo", false)
	if err != nil {
		t.Fatalf("RepairFeature returned error: %v", err)
	}
	if result.Changed {
		t.Fatalf("expected conflicting repair to avoid changes, got %+v", result)
	}
	if len(result.Warnings) == 0 || !strings.Contains(result.Warnings[0], "differ") {
		t.Fatalf("expected warning about differing reports, got %+v", result)
	}
}
