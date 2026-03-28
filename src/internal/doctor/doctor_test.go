package doctor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"draftspec/src/internal/config"
	"draftspec/src/internal/project"
)

func TestCheckHealthyWorkspace(t *testing.T) {
	root := t.TempDir()

	_, err := project.Initialize(root, project.InitOptions{InitGit: false, DefaultLang: "en", Shell: "sh", AgentTargets: []string{"claude"}})
	if err != nil {
		t.Fatalf("Initialize returned error: %v", err)
	}

	result, err := Check(root)
	if err != nil {
		t.Fatalf("Check returned error: %v", err)
	}

	if len(result.Findings) == 0 {
		t.Fatal("expected findings, got none")
	}
	if result.Findings[len(result.Findings)-1].Level != "ok" {
		t.Fatalf("last finding level = %q, want ok", result.Findings[len(result.Findings)-1].Level)
	}
}

func TestCheckWarnsAboutOrphanedAgentArtifact(t *testing.T) {
	root := t.TempDir()

	_, err := project.Initialize(root, project.InitOptions{InitGit: false, DefaultLang: "en", Shell: "sh", AgentTargets: []string{"claude", "cursor"}})
	if err != nil {
		t.Fatalf("Initialize returned error: %v", err)
	}
	_, err = project.RemoveAgents(root, project.RemoveAgentsOptions{Targets: []string{"cursor"}})
	if err != nil {
		t.Fatalf("RemoveAgents returned error: %v", err)
	}

	orphanPath := filepath.Join(root, ".cursor", "rules", "draftspec-inspect.mdc")
	if err := os.MkdirAll(filepath.Dir(orphanPath), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(orphanPath, []byte("orphan"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	result, err := Check(root)
	if err != nil {
		t.Fatalf("Check returned error: %v", err)
	}

	var foundWarning bool
	for _, finding := range result.Findings {
		if finding.Level == "warning" && strings.Contains(finding.Message, "orphaned agent artifact") {
			foundWarning = true
			break
		}
	}
	if !foundWarning {
		t.Fatalf("expected orphaned artifact warning, got %+v", result.Findings)
	}
}

func TestCheckErrorsWhenRequiredFileIsMissing(t *testing.T) {
	root := t.TempDir()

	_, err := project.Initialize(root, project.InitOptions{InitGit: false, DefaultLang: "en", Shell: "sh", AgentTargets: []string{"claude"}})
	if err != nil {
		t.Fatalf("Initialize returned error: %v", err)
	}

	missingPath := filepath.Join(root, ".draftspec", "constitution.md")
	if err := os.Remove(missingPath); err != nil {
		t.Fatalf("Remove returned error: %v", err)
	}

	result, err := Check(root)
	if err != nil {
		t.Fatalf("Check returned error: %v", err)
	}

	if len(result.Findings) == 0 || result.Findings[0].Level != "error" {
		t.Fatalf("expected first finding to be error, got %+v", result.Findings)
	}

	var foundMissing bool
	for _, finding := range result.Findings {
		if finding.Level == "error" && strings.Contains(finding.Message, "missing") && strings.Contains(finding.Message, "constitution.md") {
			foundMissing = true
			break
		}
	}
	if !foundMissing {
		t.Fatalf("expected missing constitution error, got %+v", result.Findings)
	}
}

func TestCheckHealthyPowerShellWorkspace(t *testing.T) {
	root := t.TempDir()

	_, err := project.Initialize(root, project.InitOptions{InitGit: false, DefaultLang: "en", Shell: "powershell"})
	if err != nil {
		t.Fatalf("Initialize returned error: %v", err)
	}

	result, err := Check(root)
	if err != nil {
		t.Fatalf("Check returned error: %v", err)
	}
	if result.Findings[len(result.Findings)-1].Level != "ok" {
		t.Fatalf("last finding level = %q, want ok", result.Findings[len(result.Findings)-1].Level)
	}
}

func TestCheckErrorsOnUnsupportedShell(t *testing.T) {
	root := t.TempDir()

	_, err := project.Initialize(root, project.InitOptions{InitGit: false, DefaultLang: "en", Shell: "sh"})
	if err != nil {
		t.Fatalf("Initialize returned error: %v", err)
	}

	cfg, err := config.Load(root)
	if err != nil {
		t.Fatalf("config.Load returned error: %v", err)
	}
	cfg.Runtime.Shell = "fish"
	if err := config.Save(root, cfg); err != nil {
		t.Fatalf("config.Save returned error: %v", err)
	}

	result, err := Check(root)
	if err != nil {
		t.Fatalf("Check returned error: %v", err)
	}

	var found bool
	for _, finding := range result.Findings {
		if finding.Level == "error" && strings.Contains(finding.Message, "unsupported shell") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected unsupported shell error, got %+v", result.Findings)
	}
}
