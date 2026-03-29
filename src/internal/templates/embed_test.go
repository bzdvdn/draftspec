package templates

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"testing"
)

func TestLanguageAssetSetsMatch(t *testing.T) {
	enFiles := languageAssetSet(t, "en")
	ruFiles := languageAssetSet(t, "ru")

	if len(enFiles) == 0 {
		t.Fatal("expected English language assets to be non-empty")
	}
	if len(ruFiles) == 0 {
		t.Fatal("expected Russian language assets to be non-empty")
	}

	if !reflect.DeepEqual(enFiles, ruFiles) {
		t.Fatalf("language asset sets differ\n\nen: %v\n\nru: %v", enFiles, ruFiles)
	}
}

func TestFilesBuildForSupportedLanguages(t *testing.T) {
	testCases := []struct {
		name     string
		settings LanguageSettings
	}{
		{
			name: "english",
			settings: LanguageSettings{
				Default:  "en",
				Docs:     "en",
				Agent:    "en",
				Comments: "en",
				Shell:    "sh",
			},
		},
		{
			name: "russian",
			settings: LanguageSettings{
				Default:  "ru",
				Docs:     "ru",
				Agent:    "ru",
				Comments: "ru",
				Shell:    "sh",
			},
		},
		{
			name: "mixed",
			settings: LanguageSettings{
				Default:  "en",
				Docs:     "ru",
				Agent:    "en",
				Comments: "ru",
				Shell:    "sh",
			},
		},
		{
			name: "powershell",
			settings: LanguageSettings{
				Default:  "en",
				Docs:     "en",
				Agent:    "en",
				Comments: "en",
				Shell:    "powershell",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			files, err := Files(tc.settings)
			if err != nil {
				t.Fatalf("Files() returned error: %v", err)
			}
			if len(files) == 0 {
				t.Fatal("expected generated file set to be non-empty")
			}

			targets := make(map[string]struct{}, len(files))
			for _, file := range files {
				if file.TargetPath == "" {
					t.Fatal("expected file target path to be non-empty")
				}
				if file.Content == "" {
					t.Fatalf("expected generated content for %s to be non-empty", file.TargetPath)
				}
				targets[file.TargetPath] = struct{}{}
			}

			requiredFiles := []string{
				"draftspec.yaml",
				"constitution.md",
				"templates/spec.md",
				"templates/plan.md",
				"templates/tasks.md",
				"templates/inspect-report.md",
				"templates/verify-report.md",
				"templates/archive/summary.md",
				"templates/prompts/spec.md",
				"templates/prompts/inspect.md",
				"templates/prompts/plan.md",
				"templates/prompts/tasks.md",
				"templates/prompts/implement.md",
				"templates/prompts/archive.md",
				"templates/prompts/verify.md",
			}
			ext := ".sh"
			if tc.settings.Shell == "powershell" {
				ext = ".ps1"
			}
			requiredFiles = append(requiredFiles,
				"scripts/check-inspect-ready"+ext,
				"scripts/check-archive-ready"+ext,
				"scripts/check-verify-ready"+ext,
				"scripts/verify-task-state"+ext,
			)
			for _, required := range requiredFiles {
				if _, ok := targets[required]; !ok {
					t.Fatalf("expected generated file set to include %s", required)
				}
			}
		})
	}
}

func TestInspectPromptDefinesCheapScopeAndVerdictRules(t *testing.T) {
	files, err := Files(LanguageSettings{
		Default:  "en",
		Docs:     "en",
		Agent:    "en",
		Comments: "en",
		Shell:    "sh",
	})
	if err != nil {
		t.Fatalf("Files() returned error: %v", err)
	}

	content := fileContentByTarget(t, files, "templates/prompts/inspect.md")
	requiredSnippets := []string{
		"Always read these first:",
		"Read these only when they exist and materially affect the inspection:",
		"Do Not Read By Default",
		"Prefer the cheapest inspection scope first",
		"Verify `constitution <-> spec`",
		"Verify `spec <-> plan`",
		"verify `plan <-> tasks`",
		"The `## Verdict` section MUST use one of: `pass`, `concerns`, `blocked`.",
		"major `spec <-> plan` contradictions",
	}
	for _, snippet := range requiredSnippets {
		if !strings.Contains(content, snippet) {
			t.Fatalf("expected inspect prompt to contain %q", snippet)
		}
	}
}

func TestImplementPromptSupportsFullRunAndScopedExecution(t *testing.T) {
	files, err := Files(LanguageSettings{
		Default:  "en",
		Docs:     "en",
		Agent:    "en",
		Comments: "en",
		Shell:    "sh",
	})
	if err != nil {
		t.Fatalf("Files() returned error: %v", err)
	}

	content := fileContentByTarget(t, files, "templates/prompts/implement.md")
	requiredSnippets := []string{
		"Default behavior: if the user does not restrict scope, execute all unfinished tasks in order.",
		"Scoped behavior: if the user explicitly provides `--phase <number>`, execute only that phase.",
		"Scoped behavior: if the user explicitly provides `--tasks <task-id-list>`, execute only those task IDs.",
		"Do not accept `--phase` and `--tasks` together in the same run.",
		"If scoped execution skips unfinished earlier work, warn about the ordering risk",
	}
	for _, snippet := range requiredSnippets {
		if !strings.Contains(content, snippet) {
			t.Fatalf("expected implement prompt to contain %q", snippet)
		}
	}
}

func TestInspectHelperScriptsSupportAcceptanceIDsAndCoverageFormat(t *testing.T) {
	testCases := []struct {
		name   string
		shell  string
		target string
	}{
		{name: "sh", shell: "sh", target: "scripts/inspect-spec.sh"},
		{name: "powershell", shell: "powershell", target: "scripts/inspect-spec.ps1"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			files, err := Files(LanguageSettings{
				Default:  "en",
				Docs:     "en",
				Agent:    "en",
				Comments: "en",
				Shell:    tc.shell,
			})
			if err != nil {
				t.Fatalf("Files() returned error: %v", err)
			}

			content := fileContentByTarget(t, files, tc.target)
			requiredSnippets := []string{
				"acceptance IDs",
				"acceptance coverage contains malformed entries",
				"AC-001 -> T1.1",
			}
			for _, snippet := range requiredSnippets {
				if !strings.Contains(content, snippet) {
					t.Fatalf("expected %s to contain %q", tc.target, snippet)
				}
			}
		})
	}
}

func fileContentByTarget(t *testing.T, files []File, target string) string {
	t.Helper()

	for _, file := range files {
		if file.TargetPath == target {
			return file.Content
		}
	}

	t.Fatalf("expected generated file set to include %s", target)
	return ""
}

func languageAssetSet(t *testing.T, language string) []string {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to resolve caller information")
	}

	root := filepath.Join(filepath.Dir(filename), "assets", "lang", language)
	var files []string

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		files = append(files, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		t.Fatalf("walk language assets for %s: %v", language, err)
	}

	sort.Strings(files)
	return files
}
