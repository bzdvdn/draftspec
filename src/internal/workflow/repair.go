package workflow

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"draftspec/src/internal/config"
	"draftspec/src/internal/featurepaths"
)

type RepairResult struct {
	Slug     string   `json:"slug"`
	DryRun   bool     `json:"dry_run"`
	Changed  bool     `json:"changed"`
	Actions  []string `json:"actions,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}

type MigrationResult struct {
	DryRun   bool           `json:"dry_run"`
	Changed  bool           `json:"changed"`
	Results  []RepairResult `json:"results,omitempty"`
	Warnings []string       `json:"warnings,omitempty"`
}

func RepairFeature(root, slug string, dryRun bool) (RepairResult, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return RepairResult{}, err
	}
	cfg, err := config.Load(root)
	if err != nil {
		return RepairResult{}, err
	}

	specsDir, err := cfg.SpecsDir(root)
	if err != nil {
		return RepairResult{}, err
	}
	plansDir, err := cfg.PlansDir(root)
	if err != nil {
		return RepairResult{}, err
	}

	result := RepairResult{Slug: slug, DryRun: dryRun}

	if changed, warnings, err := migrateFlatSpecArtifacts(root, specsDir, slug, dryRun, &result.Actions); err != nil {
		return RepairResult{}, err
	} else {
		result.Changed = result.Changed || changed
		result.Warnings = append(result.Warnings, warnings...)
	}

	canonicalInspectPath := featurepaths.Inspect(specsDir, slug)
	legacyInspectPath := filepath.Join(plansDir, slug, "inspect.md")
	if changed, warnings, err := migrateInspectReport(root, slug, canonicalInspectPath, legacyInspectPath, dryRun, &result.Actions); err != nil {
		return RepairResult{}, err
	} else {
		result.Changed = result.Changed || changed
		result.Warnings = append(result.Warnings, warnings...)
	}

	if !result.Changed && len(result.Warnings) == 0 {
		result.Warnings = append(result.Warnings, fmt.Sprintf("no safe migrations were needed for slug %s", slug))
	}
	return result, nil
}

func MigrateProject(root string, dryRun bool) (MigrationResult, error) {
	states, err := States(root)
	if err != nil {
		return MigrationResult{}, err
	}

	result := MigrationResult{DryRun: dryRun}
	for _, state := range states {
		repair, err := RepairFeature(root, state.Slug, dryRun)
		if err != nil {
			return MigrationResult{}, err
		}
		if repair.Changed || len(repair.Warnings) > 0 {
			result.Results = append(result.Results, repair)
		}
		if repair.Changed {
			result.Changed = true
		}
	}
	sort.Slice(result.Results, func(i, j int) bool {
		return result.Results[i].Slug < result.Results[j].Slug
	})
	if len(result.Results) == 0 {
		result.Warnings = append(result.Warnings, "no safe migrations were needed")
	}
	return result, nil
}

func displayPath(root, path string) string {
	if rel, err := filepath.Rel(root, path); err == nil {
		return filepath.ToSlash(rel)
	}
	return filepath.ToSlash(path)
}

func migrateFlatSpecArtifacts(root, specsDir, slug string, dryRun bool, actions *[]string) (bool, []string, error) {
	changed := false
	var warnings []string

	for _, artifact := range featurepaths.Artifacts(specsDir, slug) {
		canonicalExists := fileExists(artifact.CanonicalPath)
		legacyExists := fileExists(artifact.LegacyPath)

		switch {
		case !canonicalExists && !legacyExists:
			continue
		case !canonicalExists && legacyExists:
			*actions = append(*actions, fmt.Sprintf("move legacy %s from %s to %s", artifact.Name, displayPath(root, artifact.LegacyPath), displayPath(root, artifact.CanonicalPath)))
			changed = true
			if dryRun {
				continue
			}
			if err := os.MkdirAll(filepath.Dir(artifact.CanonicalPath), 0o755); err != nil {
				return false, nil, err
			}
			if err := os.Rename(artifact.LegacyPath, artifact.CanonicalPath); err != nil {
				return false, nil, fmt.Errorf("move legacy %s for slug %s: %w", artifact.Name, slug, err)
			}
		case canonicalExists && legacyExists:
			canonicalContent, err := os.ReadFile(artifact.CanonicalPath)
			if err != nil {
				return false, nil, fmt.Errorf("read canonical %s for slug %s: %w", artifact.Name, slug, err)
			}
			legacyContent, err := os.ReadFile(artifact.LegacyPath)
			if err != nil {
				return false, nil, fmt.Errorf("read legacy %s for slug %s: %w", artifact.Name, slug, err)
			}
			if !bytes.Equal(canonicalContent, legacyContent) {
				warnings = append(warnings, fmt.Sprintf("canonical and legacy %s differ for slug %s; resolve manually", artifact.Name, slug))
				continue
			}
			*actions = append(*actions, fmt.Sprintf("remove duplicate legacy %s %s", artifact.Name, displayPath(root, artifact.LegacyPath)))
			changed = true
			if dryRun {
				continue
			}
			if err := os.Remove(artifact.LegacyPath); err != nil {
				return false, nil, fmt.Errorf("remove duplicate legacy %s for slug %s: %w", artifact.Name, slug, err)
			}
		}
	}

	if !dryRun {
		_ = removeEmptyDir(featurepaths.SpecDir(specsDir, slug))
	}
	return changed, warnings, nil
}

func migrateInspectReport(root, slug, canonicalPath, legacyPath string, dryRun bool, actions *[]string) (bool, []string, error) {
	canonicalExists := fileExists(canonicalPath)
	legacyExists := fileExists(legacyPath)

	switch {
	case !canonicalExists && !legacyExists:
		return false, nil, nil
	case !canonicalExists && legacyExists:
		*actions = append(*actions, fmt.Sprintf("move legacy inspect report from %s to %s", displayPath(root, legacyPath), displayPath(root, canonicalPath)))
		if dryRun {
			return true, nil, nil
		}
		if err := os.MkdirAll(filepath.Dir(canonicalPath), 0o755); err != nil {
			return false, nil, err
		}
		if err := os.Rename(legacyPath, canonicalPath); err != nil {
			return false, nil, fmt.Errorf("move legacy inspect report for slug %s: %w", slug, err)
		}
		return true, nil, nil
	case canonicalExists && legacyExists:
		canonicalContent, err := os.ReadFile(canonicalPath)
		if err != nil {
			return false, nil, fmt.Errorf("read canonical inspect report for slug %s: %w", slug, err)
		}
		legacyContent, err := os.ReadFile(legacyPath)
		if err != nil {
			return false, nil, fmt.Errorf("read legacy inspect report for slug %s: %w", slug, err)
		}
		if !bytes.Equal(canonicalContent, legacyContent) {
			return false, []string{fmt.Sprintf("canonical and legacy inspect reports differ for slug %s; resolve manually", slug)}, nil
		}
		*actions = append(*actions, fmt.Sprintf("remove duplicate legacy inspect report %s", displayPath(root, legacyPath)))
		if dryRun {
			return true, nil, nil
		}
		if err := os.Remove(legacyPath); err != nil {
			return false, nil, fmt.Errorf("remove duplicate legacy inspect report for slug %s: %w", slug, err)
		}
		return true, nil, nil
	default:
		return false, nil, nil
	}
}

func removeEmptyDir(path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	if len(entries) > 0 {
		return nil
	}
	return os.Remove(path)
}
