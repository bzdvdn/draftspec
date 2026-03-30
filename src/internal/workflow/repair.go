package workflow

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"draftspec/src/internal/config"
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

	canonicalInspectPath := filepath.Join(specsDir, slug+".inspect.md")
	legacyInspectPath := filepath.Join(plansDir, slug, "inspect.md")

	result := RepairResult{Slug: slug, DryRun: dryRun}
	canonicalExists := fileExists(canonicalInspectPath)
	legacyExists := fileExists(legacyInspectPath)

	switch {
	case !canonicalExists && !legacyExists:
		result.Warnings = append(result.Warnings, fmt.Sprintf("no inspect report found for slug %s", slug))
		return result, nil
	case !canonicalExists && legacyExists:
		action := fmt.Sprintf("move legacy inspect report from %s to %s", displayPath(root, legacyInspectPath), displayPath(root, canonicalInspectPath))
		result.Actions = append(result.Actions, action)
		result.Changed = true
		if dryRun {
			return result, nil
		}
		if err := os.MkdirAll(filepath.Dir(canonicalInspectPath), 0o755); err != nil {
			return RepairResult{}, err
		}
		if err := os.Rename(legacyInspectPath, canonicalInspectPath); err != nil {
			return RepairResult{}, fmt.Errorf("move legacy inspect report for slug %s: %w", slug, err)
		}
		return result, nil
	case canonicalExists && legacyExists:
		canonicalContent, err := os.ReadFile(canonicalInspectPath)
		if err != nil {
			return RepairResult{}, fmt.Errorf("read canonical inspect report for slug %s: %w", slug, err)
		}
		legacyContent, err := os.ReadFile(legacyInspectPath)
		if err != nil {
			return RepairResult{}, fmt.Errorf("read legacy inspect report for slug %s: %w", slug, err)
		}
		if !bytes.Equal(canonicalContent, legacyContent) {
			result.Warnings = append(result.Warnings, fmt.Sprintf("canonical and legacy inspect reports differ for slug %s; resolve manually", slug))
			return result, nil
		}
		action := fmt.Sprintf("remove duplicate legacy inspect report %s", displayPath(root, legacyInspectPath))
		result.Actions = append(result.Actions, action)
		result.Changed = true
		if dryRun {
			return result, nil
		}
		if err := os.Remove(legacyInspectPath); err != nil {
			return RepairResult{}, fmt.Errorf("remove duplicate legacy inspect report for slug %s: %w", slug, err)
		}
		return result, nil
	default:
		return result, nil
	}
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
