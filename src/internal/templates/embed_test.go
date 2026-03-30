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
		"Default to a compact report in conversation output",
		"Produce the full sectioned report only when the user explicitly asks for a full report",
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
		"the selected work would force changes across another feature package or slug",
		"the next safe step would require inventing new tasks or acceptance coverage",
		"[T1.1] started",
		"[T1.1] done",
		"[T1.1] blocked: <reason>",
		"[Phase 1] done: T1.1, T1.2",
	}
	for _, snippet := range requiredSnippets {
		if !strings.Contains(content, snippet) {
			t.Fatalf("expected implement prompt to contain %q", snippet)
		}
	}
}

func TestSpecPromptDefinesDeterministicStagedMode(t *testing.T) {
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

	content := fileContentByTarget(t, files, "templates/prompts/spec.md")
	requiredSnippets := []string{
		"keep staged mode active for the next non-command user message",
		"If the next user message begins with `/draftspec.`, staged mode is canceled",
		"If the next user message does not begin with `/draftspec.`, treat it as the continuation of the staged spec request.",
	}
	for _, snippet := range requiredSnippets {
		if !strings.Contains(content, snippet) {
			t.Fatalf("expected spec prompt to contain %q", snippet)
		}
	}
}

func TestPlanPromptDefinesConcreteResearchTriggers(t *testing.T) {
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

	content := fileContentByTarget(t, files, "templates/prompts/plan.md")
	requiredSnippets := []string{
		"Create `.draftspec/plans/<slug>/research.md` only when at least one of these is true:",
		"external system, API, or dependency",
		"multiple realistic implementation options",
		"Do not create `research.md` for generic brainstorming",
	}
	for _, snippet := range requiredSnippets {
		if !strings.Contains(content, snippet) {
			t.Fatalf("expected plan prompt to contain %q", snippet)
		}
	}
}

func TestTasksAndImplementPromptsDoNotAssumeResearchArtifact(t *testing.T) {
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

	for _, target := range []string{
		"templates/prompts/tasks.md",
		"templates/prompts/implement.md",
	} {
		content := fileContentByTarget(t, files, target)
		if !strings.Contains(content, "Do not assume `research.md` should exist;") {
			t.Fatalf("expected %s to explain that research.md is not assumed by default", target)
		}
	}
}

func TestPromptsDefineScopeTripwiresForRefinement(t *testing.T) {
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

	testCases := []struct {
		target string
		want   []string
	}{
		{
			target: "templates/prompts/spec.md",
			want: []string{
				"multiple feature slugs or multiple independent specs",
			},
		},
		{
			target: "templates/prompts/plan.md",
			want: []string{
				"cross an unclear integration or architectural boundary",
				"multiple feature packages were planned together",
			},
		},
		{
			target: "templates/prompts/tasks.md",
			want: []string{
				"span multiple feature slugs or unrelated change sets",
				"cannot be mapped to executable work without guessing",
			},
		},
		{
			target: "templates/prompts/verify.md",
			want: []string{
				"broad repository sweep instead of focused evidence",
				"cannot be confirmed from the current tasks, plan artifacts, and targeted code inspection",
			},
		},
	}

	for _, tc := range testCases {
		content := fileContentByTarget(t, files, tc.target)
		for _, snippet := range tc.want {
			if !strings.Contains(content, snippet) {
				t.Fatalf("expected %s to contain %q", tc.target, snippet)
			}
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

func TestReadinessScriptsEnforceLeanTraceabilityRules(t *testing.T) {
	testCases := []struct {
		name   string
		shell  string
		target string
		want   []string
	}{
		{
			name:   "sh spec ready",
			shell:  "sh",
			target: "scripts/check-spec-ready.sh",
			want: []string{
				"spec template includes requirement IDs",
				"spec template includes acceptance IDs",
				"spec template includes Given marker",
			},
		},
		{
			name:   "sh plan ready",
			shell:  "sh",
			target: "scripts/check-plan-ready.sh",
			want: []string{
				"spec has stable acceptance IDs",
				"./.draftspec/scripts/inspect-spec.sh",
			},
		},
		{
			name:   "sh tasks ready",
			shell:  "sh",
			target: "scripts/check-tasks-ready.sh",
			want: []string{
				"plan has stable decision IDs",
				"spec has stable acceptance IDs",
			},
		},
		{
			name:   "sh implement ready",
			shell:  "sh",
			target: "scripts/check-implement-ready.sh",
			want: []string{
				"tasks include acceptance coverage section",
				"tasks include AC-to-task coverage lines",
			},
		},
		{
			name:   "sh verify ready",
			shell:  "sh",
			target: "scripts/check-verify-ready.sh",
			want: []string{
				"./.draftspec/scripts/inspect-spec.sh",
				"./.draftspec/scripts/verify-task-state.sh",
			},
		},
		{
			name:   "powershell verify task state",
			shell:  "powershell",
			target: "scripts/verify-task-state.ps1",
			want: []string{
				"TASK_IDS=",
				"AC_COVERAGE_LINES=",
				"no AC-to-task coverage lines found",
			},
		},
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
			for _, snippet := range tc.want {
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
