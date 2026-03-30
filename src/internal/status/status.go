package status

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"draftspec/src/internal/config"
)

type Result struct {
	Slug           string `json:"slug"`
	Phase          string `json:"phase"`
	SpecExists     bool   `json:"spec_exists"`
	PlanExists     bool   `json:"plan_exists"`
	TasksExists    bool   `json:"tasks_exists"`
	Archived       bool   `json:"archived"`
	TasksTotal     int    `json:"tasks_total"`
	TasksCompleted int    `json:"tasks_completed"`
	TasksOpen      int    `json:"tasks_open"`
	ReadyFor       string `json:"ready_for,omitempty"`
	Blocked        bool   `json:"blocked"`
}

var (
	checkboxPattern = regexp.MustCompile(`^- \[([ x])\]`)
)

func Check(root, slug string) (Result, error) {
	if slug == "" {
		return Result{}, fmt.Errorf("slug cannot be empty")
	}

	cfg, err := config.Load(root)
	if err != nil {
		return Result{}, err
	}

	specsDir, err := cfg.SpecsDir(root)
	if err != nil {
		return Result{}, err
	}
	plansDir, err := cfg.PlansDir(root)
	if err != nil {
		return Result{}, err
	}
	archiveDir, err := cfg.ArchiveDir(root)
	if err != nil {
		return Result{}, err
	}

	specPath := filepath.Join(specsDir, slug+".md")
	planPath := filepath.Join(plansDir, slug, "plan.md")
	tasksPath := filepath.Join(plansDir, slug, "tasks.md")
	archiveSlugDir := filepath.Join(archiveDir, slug)

	result := Result{
		Slug:       slug,
		SpecExists: fileExists(specPath),
		PlanExists: fileExists(planPath),
		TasksExists: fileExists(tasksPath),
	}

	if result.TasksExists {
		total, completed, open, err := taskCounts(tasksPath)
		if err != nil {
			return Result{}, err
		}
		result.TasksTotal = total
		result.TasksCompleted = completed
		result.TasksOpen = open
	}

	result.Archived = archiveExists(archiveSlugDir)

	switch {
	case result.Archived:
		result.Phase = "archive"
	case !result.SpecExists:
		result.Phase = "constitution"
		result.ReadyFor = "spec"
		result.Blocked = true
	case !result.PlanExists:
		result.Phase = "spec"
		result.ReadyFor = "plan"
	case !result.TasksExists:
		result.Phase = "plan"
		result.ReadyFor = "tasks"
	case result.TasksOpen > 0:
		result.Phase = "implement"
		result.ReadyFor = "implement"
	default:
		result.Phase = "verify"
		result.ReadyFor = "verify"
	}

	return result, nil
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func archiveExists(path string) bool {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}
	return len(entries) > 0
}

func taskCounts(path string) (total int, completed int, open int, err error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("read tasks file: %w", err)
	}

	lines := regexp.MustCompile(`\r?\n`).Split(string(content), -1)
	for _, line := range lines {
		match := checkboxPattern.FindStringSubmatch(line)
		if len(match) == 0 {
			continue
		}
		total++
		if match[1] == "x" {
			completed++
		}
		if match[1] == " " {
			open++
		}
	}

	return total, completed, open, nil
}
