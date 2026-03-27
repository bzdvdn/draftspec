package cli

import (
	"bytes"
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

	stdout, _, err := executeRoot(t, "init", root, "--git=false", "--lang", "en", "--agents", "claude")
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

	if _, _, err := executeRoot(t, "init", root, "--git=false", "--lang", "en"); err != nil {
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

	if _, _, err := executeRoot(t, "init", root, "--git=false", "--lang", "en"); err != nil {
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

func TestCleanupAgentsCommandRemovesOrphanedArtifacts(t *testing.T) {
	root := t.TempDir()

	if _, _, err := executeRoot(t, "init", root, "--git=false", "--lang", "en", "--agents", "cursor"); err != nil {
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
