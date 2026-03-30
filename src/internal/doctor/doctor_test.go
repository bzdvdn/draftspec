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

func TestCheckErrorsWhenPlanSkipsMandatoryInspect(t *testing.T) {
	root := t.TempDir()

	_, err := project.Initialize(root, project.InitOptions{InitGit: false, DefaultLang: "en", Shell: "sh"})
	if err != nil {
		t.Fatalf("Initialize returned error: %v", err)
	}

	specPath := filepath.Join(root, ".draftspec", "specs", "demo.md")
	if err := os.WriteFile(specPath, []byte("# Demo\n\n## Goal\nx\n\n## Requirements\n- RQ-001 x\n\n## Acceptance Criteria\n### AC-001 Demo\n- **Given** x\n- **When** y\n- **Then** z\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(spec) returned error: %v", err)
	}

	planDir := filepath.Join(root, ".draftspec", "plans", "demo")
	if err := os.MkdirAll(planDir, 0o755); err != nil {
		t.Fatalf("MkdirAll(planDir) returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(planDir, "plan.md"), []byte("# Demo Plan\n\n- DEC-001 x\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(plan) returned error: %v", err)
	}

	result, err := Check(root)
	if err != nil {
		t.Fatalf("Check returned error: %v", err)
	}

	var found bool
	for _, finding := range result.Findings {
		if finding.Level == "error" && strings.Contains(finding.Message, "mandatory inspect report") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected missing inspect error, got %+v", result.Findings)
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

func TestCheckWarnsWhenDraftspecEntrypointCannotBeResolved(t *testing.T) {
	root := t.TempDir()
	t.Setenv("PATH", "")
	t.Setenv("DRAFTSPEC_BIN", "")

	_, err := project.Initialize(root, project.InitOptions{InitGit: false, DefaultLang: "en", Shell: "sh"})
	if err != nil {
		t.Fatalf("Initialize returned error: %v", err)
	}

	result, err := Check(root)
	if err != nil {
		t.Fatalf("Check returned error: %v", err)
	}

	var foundWarning bool
	for _, finding := range result.Findings {
		if finding.Level == "warning" && strings.Contains(finding.Message, "set DRAFTSPEC_BIN or add draftspec to PATH") {
			foundWarning = true
			break
		}
	}
	if !foundWarning {
		t.Fatalf("expected missing entrypoint warning, got %+v", result.Findings)
	}
}
