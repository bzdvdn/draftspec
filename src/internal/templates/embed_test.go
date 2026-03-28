package templates

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
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
			},
		},
		{
			name: "russian",
			settings: LanguageSettings{
				Default:  "ru",
				Docs:     "ru",
				Agent:    "ru",
				Comments: "ru",
			},
		},
		{
			name: "mixed",
			settings: LanguageSettings{
				Default:  "en",
				Docs:     "ru",
				Agent:    "en",
				Comments: "ru",
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

			for _, required := range []string{
				"draftspec.yaml",
				"constitution.md",
				"memory.md",
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
				"scripts/check-inspect-ready.sh",
				"scripts/check-archive-ready.sh",
				"scripts/check-verify-ready.sh",
				"scripts/verify-task-state.sh",
				"scripts/verify-memory-sync.sh",
			} {
				if _, ok := targets[required]; !ok {
					t.Fatalf("expected generated file set to include %s", required)
				}
			}
		})
	}
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
