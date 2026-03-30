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
	if !strings.Contains(stdout, "ok: draftspec workspace looks healthy") {
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
		PlanExists     bool   `json:"plan_exists"`
		TasksExists    bool   `json:"tasks_exists"`
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
	if !payload.SpecExists || !payload.PlanExists || !payload.TasksExists {
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
		PlanExists     bool   `json:"plan_exists"`
		TasksExists    bool   `json:"tasks_exists"`
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
		if payload.SpecExists != want.SpecExists || payload.PlanExists != want.PlanExists || payload.TasksExists != want.TasksExists {
			t.Fatalf("unexpected artifact flags: %+v, want %+v", payload, want)
		}
		if payload.TasksTotal != want.TasksTotal || payload.TasksCompleted != want.TasksCompleted || payload.TasksOpen != want.TasksOpen {
			t.Fatalf("unexpected task counts: %+v, want %+v", payload, want)
		}
	}

	checkStatus(statusPayload{
		Phase:      "constitution",
		ReadyFor:   "spec",
		Blocked:    true,
		SpecExists: false,
		PlanExists: false,
		TasksExists:false,
	})

	specPath := filepath.Join(root, ".draftspec", "specs", "demo.md")
	specContent := "# Feature Specification: Demo\n\n## Requirements\n- RQ-001 Support a minimal demo flow.\n\n## Acceptance Criteria\n- AC-001\n  - Given a prepared workspace\n  - When the feature lifecycle is checked\n  - Then the status should advance predictably.\n"
	if err := os.WriteFile(specPath, []byte(specContent), 0o644); err != nil {
		t.Fatalf("WriteFile(spec) returned error: %v", err)
	}

	checkStatus(statusPayload{
		Phase:      "spec",
		ReadyFor:   "plan",
		Blocked:    false,
		SpecExists: true,
		PlanExists: false,
		TasksExists:false,
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
		Phase:      "plan",
		ReadyFor:   "tasks",
		Blocked:    false,
		SpecExists: true,
		PlanExists: true,
		TasksExists:false,
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
		PlanExists:     true,
		TasksExists:    true,
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
		PlanExists:     true,
		TasksExists:    true,
		TasksTotal:     2,
		TasksCompleted: 2,
		TasksOpen:      0,
	})
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
