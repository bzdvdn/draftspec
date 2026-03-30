package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func executeRoot(t *testing.T, args ...string) (string, string, error) {
	t.Helper()

	cmd := NewRootCmd()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.SetArgs(args)

	err := cmd.Execute()
	return stdout.String(), stderr.String(), err
}

func TestInitCommandCreatesWorkspace(t *testing.T) {
	root := t.TempDir()

	stdout, _, err := executeRoot(t, "init", root, "--git=false", "--lang", "en", "--shell", "sh", "--agents", "claude")
	if err != nil {
		t.Fatalf("init command returned error: %v", err)
	}
	if !strings.Contains(stdout, "enabled agent targets: claude") {
		t.Fatalf("unexpected init output: %s", stdout)
	}

	required := []string{
		filepath.Join(root, ".draftspec", "draftspec.yaml"),
		filepath.Join(root, ".draftspec", "constitution.md"),
		filepath.Join(root, ".claude", "commands", "draftspec.inspect.md"),
	}
	for _, path := range required {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected %s to exist: %v", path, err)
		}
	}
}

func TestListSpecsAndShowSpecCommands(t *testing.T) {
	root := t.TempDir()

	if _, _, err := executeRoot(t, "init", root, "--git=false", "--lang", "en", "--shell", "sh"); err != nil {
		t.Fatalf("init command returned error: %v", err)
	}

	specsDir := filepath.Join(root, ".draftspec", "specs")
	if err := os.WriteFile(filepath.Join(specsDir, "alpha.md"), []byte("# Alpha\n"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(specsDir, "beta.md"), []byte("# Beta\nBody\n"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	stdout, _, err := executeRoot(t, "list-specs", root)
	if err != nil {
		t.Fatalf("list-specs command returned error: %v", err)
	}
	if strings.TrimSpace(stdout) != "alpha\nbeta" {
		t.Fatalf("unexpected list-specs output: %q", stdout)
	}

	stdout, _, err = executeRoot(t, "show-spec", "beta", root)
	if err != nil {
		t.Fatalf("show-spec command returned error: %v", err)
	}
	if stdout != "# Beta\nBody\n" {
		t.Fatalf("unexpected show-spec output: %q", stdout)
	}
}

func TestAddAgentAndDoctorCommands(t *testing.T) {
	root := t.TempDir()

	if _, _, err := executeRoot(t, "init", root, "--git=false", "--lang", "en", "--shell", "sh"); err != nil {
		t.Fatalf("init command returned error: %v", err)
	}

	stdout, _, err := executeRoot(t, "add-agent", root, "--agents", "cursor")
	if err != nil {
		t.Fatalf("add-agent command returned error: %v", err)
	}
	if !strings.Contains(stdout, "enabled agent targets: cursor") {
		t.Fatalf("unexpected add-agent output: %s", stdout)
	}

	stdout, _, err = executeRoot(t, "doctor", root)
	if err != nil {
		t.Fatalf("doctor command returned error: %v", err)
	}
	if !strings.Contains(stdout, "summary:") || !strings.Contains(stdout, "oks:") || !strings.Contains(stdout, "draftspec workspace looks healthy") {
		t.Fatalf("unexpected doctor output: %s", stdout)
	}
}

func TestDoctorCommandJSONOutput(t *testing.T) {
	root := t.TempDir()

	if _, _, err := executeRoot(t, "init", root, "--git=false", "--lang", "en", "--shell", "sh"); err != nil {
		t.Fatalf("init command returned error: %v", err)
	}

	stdout, _, err := executeRoot(t, "doctor", root, "--json")
	if err != nil {
		t.Fatalf("doctor --json returned error: %v", err)
	}

	var payload struct {
		Findings []struct {
			Level   string `json:"Level"`
			Message string `json:"Message"`
		} `json:"Findings"`
	}
	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("failed to parse doctor json output %q: %v", stdout, err)
	}
	if len(payload.Findings) == 0 {
		t.Fatalf("expected findings in doctor json output, got %q", stdout)
	}
	if payload.Findings[len(payload.Findings)-1].Level != "ok" {
		t.Fatalf("expected trailing ok finding in json output, got %+v", payload.Findings)
	}
}

func TestStatusCommandJSONOutput(t *testing.T) {
	root := t.TempDir()

	if _, _, err := executeRoot(t, "init", root, "--git=false", "--lang", "en", "--shell", "sh"); err != nil {
		t.Fatalf("init command returned error: %v", err)
	}

	specPath := filepath.Join(root, ".draftspec", "specs", "demo.md")
	if err := os.WriteFile(specPath, []byte("# Demo\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(spec) returned error: %v", err)
	}

	inspectPath := filepath.Join(root, ".draftspec", "specs", "demo.inspect.md")
	if err := os.WriteFile(inspectPath, []byte("---\nreport_type: inspect\nslug: demo\nstatus: pass\ndocs_language: en\ngenerated_at: 2026-03-30\n---\n# Inspect Report: demo\n\n## Verdict\n\n- status: pass\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(inspect) returned error: %v", err)
	}

	planDir := filepath.Join(root, ".draftspec", "plans", "demo")
	if err := os.MkdirAll(planDir, 0o755); err != nil {
		t.Fatalf("MkdirAll(planDir) returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(planDir, "plan.md"), []byte("# Demo Plan\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(plan) returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(planDir, "tasks.md"), []byte("- [x] T1.1 Done\n- [ ] T1.2 Open\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(tasks) returned error: %v", err)
	}

	stdout, _, err := executeRoot(t, "status", "demo", root, "--json")
	if err != nil {
		t.Fatalf("status --json returned error: %v", err)
	}

	var payload struct {
		Slug           string `json:"slug"`
		Phase          string `json:"phase"`
		SpecExists     bool   `json:"spec_exists"`
		InspectExists  bool   `json:"inspect_exists"`
		PlanExists     bool   `json:"plan_exists"`
		TasksExists    bool   `json:"tasks_exists"`
		VerifyExists   bool   `json:"verify_exists"`
		TasksTotal     int    `json:"tasks_total"`
		TasksCompleted int    `json:"tasks_completed"`
		TasksOpen      int    `json:"tasks_open"`
		ReadyFor       string `json:"ready_for"`
		Blocked        bool   `json:"blocked"`
	}
	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("failed to parse status json output %q: %v", stdout, err)
	}
	if payload.Slug != "demo" || payload.Phase != "implement" {
		t.Fatalf("unexpected status payload: %+v", payload)
	}
	if !payload.SpecExists || !payload.InspectExists || !payload.PlanExists || !payload.TasksExists {
		t.Fatalf("expected spec/plan/tasks to exist, got %+v", payload)
	}
	if payload.TasksTotal != 2 || payload.TasksCompleted != 1 || payload.TasksOpen != 1 {
		t.Fatalf("unexpected task counts: %+v", payload)
	}
	if payload.ReadyFor != "implement" || payload.Blocked {
		t.Fatalf("unexpected ready/block state: %+v", payload)
	}
}

func TestInitAndStatusCommandsFollowFeatureLifecycle(t *testing.T) {
	root := t.TempDir()

	if _, _, err := executeRoot(t, "init", root, "--git=false", "--lang", "en", "--shell", "sh"); err != nil {
		t.Fatalf("init command returned error: %v", err)
	}

	type statusPayload struct {
		Slug           string `json:"slug"`
		Phase          string `json:"phase"`
		SpecExists     bool   `json:"spec_exists"`
		InspectExists  bool   `json:"inspect_exists"`
		PlanExists     bool   `json:"plan_exists"`
		TasksExists    bool   `json:"tasks_exists"`
		VerifyExists   bool   `json:"verify_exists"`
		TasksTotal     int    `json:"tasks_total"`
		TasksCompleted int    `json:"tasks_completed"`
		TasksOpen      int    `json:"tasks_open"`
		ReadyFor       string `json:"ready_for"`
		Blocked        bool   `json:"blocked"`
	}

	checkStatus := func(want statusPayload) {
		t.Helper()

		stdout, _, err := executeRoot(t, "status", "demo", root, "--json")
		if err != nil {
			t.Fatalf("status --json returned error: %v", err)
		}

		var payload statusPayload
		if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
			t.Fatalf("failed to parse status json output %q: %v", stdout, err)
		}

		if payload.Slug != "demo" {
			t.Fatalf("slug = %q, want demo", payload.Slug)
		}
		if payload.Phase != want.Phase || payload.ReadyFor != want.ReadyFor || payload.Blocked != want.Blocked {
			t.Fatalf("unexpected phase payload: %+v, want %+v", payload, want)
		}
		if payload.SpecExists != want.SpecExists || payload.InspectExists != want.InspectExists || payload.PlanExists != want.PlanExists || payload.TasksExists != want.TasksExists || payload.VerifyExists != want.VerifyExists {
			t.Fatalf("unexpected artifact flags: %+v, want %+v", payload, want)
		}
		if payload.TasksTotal != want.TasksTotal || payload.TasksCompleted != want.TasksCompleted || payload.TasksOpen != want.TasksOpen {
			t.Fatalf("unexpected task counts: %+v, want %+v", payload, want)
		}
	}

	checkStatus(statusPayload{
		Phase:         "constitution",
		ReadyFor:      "spec",
		Blocked:       true,
		SpecExists:    false,
		InspectExists: false,
		PlanExists:    false,
		TasksExists:   false,
		VerifyExists:  false,
	})

	specPath := filepath.Join(root, ".draftspec", "specs", "demo.md")
	specContent := "# Feature Specification: Demo\n\n## Requirements\n- RQ-001 Support a minimal demo flow.\n\n## Acceptance Criteria\n- AC-001\n  - Given a prepared workspace\n  - When the feature lifecycle is checked\n  - Then the status should advance predictably.\n"
	if err := os.WriteFile(specPath, []byte(specContent), 0o644); err != nil {
		t.Fatalf("WriteFile(spec) returned error: %v", err)
	}

	checkStatus(statusPayload{
		Phase:         "spec",
		ReadyFor:      "inspect",
		Blocked:       false,
		SpecExists:    true,
		InspectExists: false,
		PlanExists:    false,
		TasksExists:   false,
		VerifyExists:  false,
	})

	inspectPath := filepath.Join(root, ".draftspec", "specs", "demo.inspect.md")
	inspectContent := "---\nreport_type: inspect\nslug: demo\nstatus: pass\ndocs_language: en\ngenerated_at: 2026-03-30\n---\n# Inspect Report: demo\n\n## Verdict\n\n- status: pass\n"
	if err := os.WriteFile(inspectPath, []byte(inspectContent), 0o644); err != nil {
		t.Fatalf("WriteFile(inspect) returned error: %v", err)
	}

	checkStatus(statusPayload{
		Phase:         "inspect",
		ReadyFor:      "plan",
		Blocked:       false,
		SpecExists:    true,
		InspectExists: true,
		PlanExists:    false,
		TasksExists:   false,
		VerifyExists:  false,
	})

	planDir := filepath.Join(root, ".draftspec", "plans", "demo")
	if err := os.MkdirAll(planDir, 0o755); err != nil {
		t.Fatalf("MkdirAll(planDir) returned error: %v", err)
	}
	planContent := "# Implementation Plan: Demo\n\n## Decisions\n- DEC-001 Keep the integration test minimal and deterministic.\n"
	if err := os.WriteFile(filepath.Join(planDir, "plan.md"), []byte(planContent), 0o644); err != nil {
		t.Fatalf("WriteFile(plan) returned error: %v", err)
	}

	checkStatus(statusPayload{
		Phase:         "plan",
		ReadyFor:      "tasks",
		Blocked:       false,
		SpecExists:    true,
		InspectExists: true,
		PlanExists:    true,
		TasksExists:   false,
		VerifyExists:  false,
	})

	tasksContent := "# Tasks: Demo\n\n## Phase 1: Implementation\n- [x] T1.1 Create the first slice\n- [ ] T1.2 Finish the second slice\n\n## Acceptance Coverage\n- AC-001 -> T1.1, T1.2\n"
	if err := os.WriteFile(filepath.Join(planDir, "tasks.md"), []byte(tasksContent), 0o644); err != nil {
		t.Fatalf("WriteFile(tasks) returned error: %v", err)
	}

	checkStatus(statusPayload{
		Phase:          "implement",
		ReadyFor:       "implement",
		Blocked:        false,
		SpecExists:     true,
		InspectExists:  true,
		PlanExists:     true,
		TasksExists:    true,
		VerifyExists:   false,
		TasksTotal:     2,
		TasksCompleted: 1,
		TasksOpen:      1,
	})

	completeTasks := "# Tasks: Demo\n\n## Phase 1: Implementation\n- [x] T1.1 Create the first slice\n- [x] T1.2 Finish the second slice\n\n## Acceptance Coverage\n- AC-001 -> T1.1, T1.2\n"
	if err := os.WriteFile(filepath.Join(planDir, "tasks.md"), []byte(completeTasks), 0o644); err != nil {
		t.Fatalf("WriteFile(tasks complete) returned error: %v", err)
	}

	checkStatus(statusPayload{
		Phase:          "verify",
		ReadyFor:       "verify",
		Blocked:        false,
		SpecExists:     true,
		InspectExists:  true,
		PlanExists:     true,
		TasksExists:    true,
		VerifyExists:   false,
		TasksTotal:     2,
		TasksCompleted: 2,
		TasksOpen:      0,
	})

	verifyContent := "---\nreport_type: verify\nslug: demo\nstatus: pass\ndocs_language: en\ngenerated_at: 2026-03-30\n---\n# Verify Report: demo\n\n## Verdict\n\n- status: pass\n"
	if err := os.WriteFile(filepath.Join(planDir, "verify.md"), []byte(verifyContent), 0o644); err != nil {
		t.Fatalf("WriteFile(verify) returned error: %v", err)
	}

	checkStatus(statusPayload{
		Phase:          "verify",
		ReadyFor:       "archive",
		Blocked:        false,
		SpecExists:     true,
		InspectExists:  true,
		PlanExists:     true,
		TasksExists:    true,
		VerifyExists:   true,
		TasksTotal:     2,
		TasksCompleted: 2,
		TasksOpen:      0,
	})
}

func TestFeaturesCommandSummarizesProjectWorkflow(t *testing.T) {
	root := t.TempDir()

	if _, _, err := executeRoot(t, "init", root, "--git=false", "--lang", "en", "--shell", "sh"); err != nil {
		t.Fatalf("init command returned error: %v", err)
	}

	if err := os.WriteFile(filepath.Join(root, ".draftspec", "specs", "alpha.md"), []byte("# Alpha\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(spec) returned error: %v", err)
	}

	stdout, _, err := executeRoot(t, "features", root)
	if err != nil {
		t.Fatalf("features command returned error: %v", err)
	}
	if !strings.Contains(stdout, "slug") || !strings.Contains(stdout, "issues") || !strings.Contains(stdout, "alpha") || !strings.Contains(stdout, "inspect") || !strings.Contains(stdout, "verify") || !strings.Contains(stdout, "tasks") {
		t.Fatalf("unexpected features output: %s", stdout)
	}
}

func TestFeatureCommandShowsDetailedWorkflowView(t *testing.T) {
	root := t.TempDir()

	if _, _, err := executeRoot(t, "init", root, "--git=false", "--lang", "en", "--shell", "sh"); err != nil {
		t.Fatalf("init command returned error: %v", err)
	}

	specContent := "# Demo\n\n## Goal\nx\n\n## Requirements\n- RQ-001 x\n\n## Acceptance Criteria\n### AC-001 Demo\n- Given x\n- When y\n- Then z\n"
	if err := os.WriteFile(filepath.Join(root, ".draftspec", "specs", "demo.md"), []byte(specContent), 0o644); err != nil {
		t.Fatalf("WriteFile(spec) returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, ".draftspec", "specs", "demo.inspect.md"), []byte("---\nreport_type: inspect\nslug: demo\nstatus: concerns\ndocs_language: en\ngenerated_at: 2026-03-31\n---\n# Inspect Report: demo\n\n## Verdict\n\n- status: concerns\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(inspect) returned error: %v", err)
	}

	stdout, _, err := executeRoot(t, "feature", "demo", root)
	if err != nil {
		t.Fatalf("feature command returned error: %v", err)
	}
	if !strings.Contains(stdout, "inspect_status: concerns") || !strings.Contains(stdout, "ready_for: plan") || !strings.Contains(stdout, "issues:") || !strings.Contains(stdout, "focus: write the plan package") || strings.Contains(stdout, "verify_path:") {
		t.Fatalf("unexpected feature output: %s", stdout)
	}
}

func TestFeatureCommandShowsSemanticFindings(t *testing.T) {
	root := t.TempDir()

	if _, _, err := executeRoot(t, "init", root, "--git=false", "--lang", "en", "--shell", "sh"); err != nil {
		t.Fatalf("init command returned error: %v", err)
	}

	specContent := "# Demo\n\n## Goal\nx\n\n## Requirements\n- RQ-001 x\n\n## Acceptance Criteria\n### AC-001 Demo\n- Given x\n- When y\n- Then z\n\n### AC-002 Demo\n- Given a\n- When b\n- Then c\n"
	if err := os.WriteFile(filepath.Join(root, ".draftspec", "specs", "demo.md"), []byte(specContent), 0o644); err != nil {
		t.Fatalf("WriteFile(spec) returned error: %v", err)
	}
	inspectContent := "---\nreport_type: inspect\nslug: demo\nstatus: pass\ndocs_language: en\ngenerated_at: 2026-03-31\n---\n# Inspect Report: demo\n\n## Verdict\n\n- status: pass\n"
	if err := os.WriteFile(filepath.Join(root, ".draftspec", "specs", "demo.inspect.md"), []byte(inspectContent), 0o644); err != nil {
		t.Fatalf("WriteFile(inspect) returned error: %v", err)
	}
	planDir := filepath.Join(root, ".draftspec", "plans", "demo")
	if err := os.MkdirAll(planDir, 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(planDir, "plan.md"), []byte("# Demo Plan\n\n## Acceptance Approach\n- AC-001 -> path\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(plan) returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(planDir, "tasks.md"), []byte("# Tasks\n\n## Phase 1: Foundation\n- [ ] T1.1 do\n\n## Acceptance Coverage\n- AC-001 -> T1.1\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(tasks) returned error: %v", err)
	}

	stdout, _, err := executeRoot(t, "feature", "demo", root)
	if err != nil {
		t.Fatalf("feature command returned error: %v", err)
	}
	if !strings.Contains(stdout, "warnings:") || !strings.Contains(stdout, "plan does not reference acceptance criterion AC-002") || strings.Contains(stdout, "for slug demo") {
		t.Fatalf("expected feature output to include semantic findings, got %s", stdout)
	}
}

func TestDoctorCommandPrefixesWorkspaceAndFeatureFindings(t *testing.T) {
	root := t.TempDir()

	if _, _, err := executeRoot(t, "init", root, "--git=false", "--lang", "en", "--shell", "sh"); err != nil {
		t.Fatalf("init command returned error: %v", err)
	}

	specContent := "# Demo\n\n## Goal\nx\n\n## Requirements\n- RQ-001 x\n\n## Acceptance Criteria\n### AC-001 Demo\n- Given x\n- When y\n- Then z\n"
	if err := os.WriteFile(filepath.Join(root, ".draftspec", "specs", "demo.md"), []byte(specContent), 0o644); err != nil {
		t.Fatalf("WriteFile(spec) returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, ".draftspec", "specs", "demo.inspect.md"), []byte("---\nreport_type: inspect\nslug: demo\nstatus: pass\ndocs_language: en\ngenerated_at: 2026-03-31\n---\n# Inspect Report: demo\n\n## Verdict\n\n- status: pass\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(inspect) returned error: %v", err)
	}
	planDir := filepath.Join(root, ".draftspec", "plans", "demo")
	if err := os.MkdirAll(planDir, 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(planDir, "plan.md"), []byte("# Demo Plan\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(plan) returned error: %v", err)
	}

	stdout, _, err := executeRoot(t, "doctor", root)
	if err != nil {
		t.Fatalf("doctor command returned error: %v", err)
	}
	if !strings.Contains(stdout, "[workspace]") || !strings.Contains(stdout, "[demo]") {
		t.Fatalf("expected doctor output to prefix workspace and feature findings, got %s", stdout)
	}
}

func TestFeatureCommandShowsLegacyInspectHint(t *testing.T) {
	root := t.TempDir()

	if _, _, err := executeRoot(t, "init", root, "--git=false", "--lang", "en", "--shell", "sh"); err != nil {
		t.Fatalf("init command returned error: %v", err)
	}

	legacyInspectPath := filepath.Join(root, ".draftspec", "plans", "demo", "inspect.md")
	if err := os.MkdirAll(filepath.Dir(legacyInspectPath), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	content := "---\nreport_type: inspect\nslug: demo\nstatus: concerns\ndocs_language: en\ngenerated_at: 2026-03-31\n---\n# Inspect Report: demo\n\n## Verdict\n\n- status: concerns\n"
	if err := os.WriteFile(legacyInspectPath, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(inspect) returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, ".draftspec", "specs", "demo.md"), []byte("# Demo\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(spec) returned error: %v", err)
	}

	stdout, _, err := executeRoot(t, "feature", "demo", root)
	if err != nil {
		t.Fatalf("feature command returned error: %v", err)
	}
	if !strings.Contains(stdout, "inspect_legacy: true") {
		t.Fatalf("expected feature output to show legacy inspect hint, got %s", stdout)
	}
}

func TestFeatureRepairCommandMigratesLegacyInspectReport(t *testing.T) {
	root := t.TempDir()

	if _, _, err := executeRoot(t, "init", root, "--git=false", "--lang", "en", "--shell", "sh"); err != nil {
		t.Fatalf("init command returned error: %v", err)
	}

	legacyInspectPath := filepath.Join(root, ".draftspec", "plans", "demo", "inspect.md")
	if err := os.MkdirAll(filepath.Dir(legacyInspectPath), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	content := "---\nreport_type: inspect\nslug: demo\nstatus: pass\ndocs_language: en\ngenerated_at: 2026-03-31\n---\n# Inspect Report: demo\n"
	if err := os.WriteFile(legacyInspectPath, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	stdout, _, err := executeRoot(t, "feature", "repair", "demo", root)
	if err != nil {
		t.Fatalf("feature repair command returned error: %v", err)
	}
	if !strings.Contains(stdout, "changed: true") || !strings.Contains(stdout, "move legacy inspect report") {
		t.Fatalf("unexpected feature repair output: %s", stdout)
	}
	if _, err := os.Stat(filepath.Join(root, ".draftspec", "specs", "demo.inspect.md")); err != nil {
		t.Fatalf("expected canonical inspect report after repair: %v", err)
	}
}

func TestMigrateCommandRepairsLegacyInspectReportsAcrossProject(t *testing.T) {
	root := t.TempDir()

	if _, _, err := executeRoot(t, "init", root, "--git=false", "--lang", "en", "--shell", "sh"); err != nil {
		t.Fatalf("init command returned error: %v", err)
	}

	for _, slug := range []string{"alpha", "beta"} {
		legacyInspectPath := filepath.Join(root, ".draftspec", "plans", slug, "inspect.md")
		if err := os.MkdirAll(filepath.Dir(legacyInspectPath), 0o755); err != nil {
			t.Fatalf("MkdirAll returned error: %v", err)
		}
		content := "---\nreport_type: inspect\nslug: " + slug + "\nstatus: pass\ndocs_language: en\ngenerated_at: 2026-03-31\n---\n# Inspect Report: " + slug + "\n"
		if err := os.WriteFile(legacyInspectPath, []byte(content), 0o644); err != nil {
			t.Fatalf("WriteFile returned error: %v", err)
		}
	}

	stdout, _, err := executeRoot(t, "migrate", root)
	if err != nil {
		t.Fatalf("migrate command returned error: %v", err)
	}
	if !strings.Contains(stdout, "slug: alpha") || !strings.Contains(stdout, "slug: beta") {
		t.Fatalf("unexpected migrate output: %s", stdout)
	}
	for _, slug := range []string{"alpha", "beta"} {
		if _, err := os.Stat(filepath.Join(root, ".draftspec", "specs", slug+".inspect.md")); err != nil {
			t.Fatalf("expected canonical inspect report for %s after migrate: %v", slug, err)
		}
	}
}

func TestCleanupAgentsCommandRemovesOrphanedArtifacts(t *testing.T) {
	root := t.TempDir()

	if _, _, err := executeRoot(t, "init", root, "--git=false", "--lang", "en", "--shell", "sh", "--agents", "cursor"); err != nil {
		t.Fatalf("init command returned error: %v", err)
	}
	if _, _, err := executeRoot(t, "remove-agent", root, "--agents", "cursor"); err != nil {
		t.Fatalf("remove-agent command returned error: %v", err)
	}

	orphanPath := filepath.Join(root, ".cursor", "rules", "draftspec-inspect.mdc")
	if err := os.MkdirAll(filepath.Dir(orphanPath), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(orphanPath, []byte("orphan"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	stdout, _, err := executeRoot(t, "cleanup-agents", root)
	if err != nil {
		t.Fatalf("cleanup-agents command returned error: %v", err)
	}
	if !strings.Contains(stdout, "removed orphaned agent artifact") {
		t.Fatalf("unexpected cleanup-agents output: %s", stdout)
	}
	if _, err := os.Stat(orphanPath); !os.IsNotExist(err) {
		t.Fatalf("expected orphaned file to be removed, got err=%v", err)
	}
}

func TestInitCommandRequiresShell(t *testing.T) {
	root := t.TempDir()

	_, _, err := executeRoot(t, "init", root, "--git=false", "--lang", "en")
	if err == nil {
		t.Fatal("expected init without --shell to fail")
	}
}

func TestRefreshCommandUpdatesManagedArtifacts(t *testing.T) {
	root := t.TempDir()

	if _, _, err := executeRoot(t, "init", root, "--git=false", "--lang", "en", "--shell", "sh", "--agents", "claude"); err != nil {
		t.Fatalf("init command returned error: %v", err)
	}

	promptPath := filepath.Join(root, ".draftspec", "templates", "prompts", "inspect.md")
	if err := os.WriteFile(promptPath, []byte("stale prompt"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	stdout, _, err := executeRoot(t, "refresh", root, "--shell", "powershell")
	if err != nil {
		t.Fatalf("refresh command returned error: %v", err)
	}
	if !strings.Contains(stdout, "update .draftspec/templates/prompts/inspect.md") {
		t.Fatalf("unexpected refresh output: %s", stdout)
	}

	if _, err := os.Stat(filepath.Join(root, ".draftspec", "scripts", "check-spec-ready.ps1")); err != nil {
		t.Fatalf("expected refreshed powershell script to exist: %v", err)
	}
}

func TestRefreshCommandJSONDryRunOutput(t *testing.T) {
	root := t.TempDir()

	if _, _, err := executeRoot(t, "init", root, "--git=false", "--lang", "en", "--shell", "sh"); err != nil {
		t.Fatalf("init command returned error: %v", err)
	}

	stdout, _, err := executeRoot(t, "refresh", root, "--shell", "powershell", "--dry-run", "--json")
	if err != nil {
		t.Fatalf("refresh command returned error: %v", err)
	}

	var payload struct {
		DryRun  bool     `json:"dry_run"`
		Updated []string `json:"updated"`
		Created []string `json:"created"`
	}
	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("failed to parse refresh json output %q: %v", stdout, err)
	}
	if !payload.DryRun {
		t.Fatalf("expected dry_run true in refresh json output, got %q", stdout)
	}
	if len(payload.Updated) == 0 && len(payload.Created) == 0 {
		t.Fatalf("expected refresh json to report pending changes, got %q", stdout)
	}
}

func TestInternalInspectSpecCommandUsesWorkflowBackend(t *testing.T) {
	root := t.TempDir()

	if _, _, err := executeRoot(t, "init", root, "--git=false", "--lang", "en", "--shell", "sh"); err != nil {
		t.Fatalf("init command returned error: %v", err)
	}

	specPath := filepath.Join(root, ".draftspec", "specs", "demo.md")
	specContent := "# Demo\n\n## Goal\nx\n\n## Requirements\n- RQ-001 x\n\n## Acceptance Criteria\n### AC-001 Demo\n- Given x\n- When y\n- Then z\n"
	if err := os.WriteFile(specPath, []byte(specContent), 0o644); err != nil {
		t.Fatalf("WriteFile(spec) returned error: %v", err)
	}

	stdout, _, err := executeRoot(t, "__internal", "inspect-spec", "--root", root, ".draftspec/specs/demo.md")
	if err != nil {
		t.Fatalf("internal inspect-spec command returned error: %v", err)
	}
	if !strings.Contains(stdout, "SUMMARY: errors=0") {
		t.Fatalf("unexpected internal inspect-spec output: %s", stdout)
	}
}

func TestInternalVerifyTaskStateCommandReturnsNonFatalWarnings(t *testing.T) {
	root := t.TempDir()

	if _, _, err := executeRoot(t, "init", root, "--git=false", "--lang", "en", "--shell", "sh"); err != nil {
		t.Fatalf("init command returned error: %v", err)
	}

	tasksPath := filepath.Join(root, ".draftspec", "plans", "demo", "tasks.md")
	if err := os.MkdirAll(filepath.Dir(tasksPath), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(tasksPath, []byte("- [ ] T1.1 open\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(tasks) returned error: %v", err)
	}

	stdout, _, err := executeRoot(t, "__internal", "verify-task-state", "--root", root, "demo")
	if err != nil {
		t.Fatalf("internal verify-task-state command returned error: %v", err)
	}
	if !strings.Contains(stdout, "TASKS_OPEN=1") || !strings.Contains(stdout, "WARN: open tasks remain") {
		t.Fatalf("unexpected internal verify-task-state output: %s", stdout)
	}
}

func TestInternalListOpenTasksCommandUsesCLIBackend(t *testing.T) {
	root := t.TempDir()

	if _, _, err := executeRoot(t, "init", root, "--git=false", "--lang", "en", "--shell", "sh"); err != nil {
		t.Fatalf("init command returned error: %v", err)
	}

	tasksPath := filepath.Join(root, ".draftspec", "plans", "demo", "tasks.md")
	if err := os.MkdirAll(filepath.Dir(tasksPath), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(tasksPath, []byte("- [x] T1.1 done\n- [ ] T1.2 open\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(tasks) returned error: %v", err)
	}

	stdout, _, err := executeRoot(t, "__internal", "list-open-tasks", "--root", root, "demo")
	if err != nil {
		t.Fatalf("internal list-open-tasks command returned error: %v", err)
	}
	if strings.TrimSpace(stdout) != "- [ ] T1.2 open" {
		t.Fatalf("unexpected internal list-open-tasks output: %q", stdout)
	}
}

func TestInternalLinkAgentsCommandUsesCLIBackend(t *testing.T) {
	root := t.TempDir()

	if _, _, err := executeRoot(t, "init", root, "--git=false", "--lang", "en", "--shell", "sh"); err != nil {
		t.Fatalf("init command returned error: %v", err)
	}

	agentsPath := filepath.Join(root, "CUSTOM_AGENTS.md")
	snippetPath := filepath.Join(root, ".draftspec", "templates", "agents-snippet.md")

	stdout, _, err := executeRoot(t, "__internal", "link-agents", "--root", root, "CUSTOM_AGENTS.md", ".draftspec/templates/agents-snippet.md")
	if err != nil {
		t.Fatalf("internal link-agents command returned error: %v", err)
	}
	if !strings.Contains(stdout, "Draftspec block added to CUSTOM_AGENTS.md") {
		t.Fatalf("unexpected internal link-agents output: %s", stdout)
	}

	content, err := os.ReadFile(agentsPath)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	snippet, err := os.ReadFile(snippetPath)
	if err != nil {
		t.Fatalf("ReadFile(snippet) returned error: %v", err)
	}
	if !strings.Contains(string(content), strings.TrimSpace(string(snippet))) {
		t.Fatalf("expected linked agents file to contain snippet, got %q", string(content))
	}
}
