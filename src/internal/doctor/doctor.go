package doctor

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"draftspec/src/internal/agents"
	"draftspec/src/internal/config"
)

type Finding struct {
	Level   string
	Message string
}

type Result struct {
	Findings []Finding
}

func Check(root string) (Result, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return Result{}, err
	}

	cfg, err := config.Load(root)
	if err != nil {
		return Result{}, err
	}

	var findings []Finding

	draftspecDir, err := cfg.DraftspecDir(root)
	if err != nil {
		return Result{}, err
	}
	configPath, err := cfg.ConfigPath(root)
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
	templatesDir, err := cfg.TemplatesDir(root)
	if err != nil {
		return Result{}, err
	}
	scriptsDir, err := cfg.ScriptsDir(root)
	if err != nil {
		return Result{}, err
	}

	for _, path := range []string{draftspecDir, specsDir, plansDir, archiveDir, templatesDir, scriptsDir} {
		checkPath(&findings, path, true)
	}
	for _, path := range []string{
		configPath,
		filepath.Join(root, cfg.Project.ConstitutionFile),
		filepath.Join(templatesDir, cfg.Templates.Spec),
		filepath.Join(templatesDir, cfg.Templates.Plan),
		filepath.Join(templatesDir, cfg.Templates.Tasks),
		filepath.Join(templatesDir, cfg.Templates.InspectReport),
		filepath.Join(templatesDir, cfg.Templates.VerifyReport),
		filepath.Join(templatesDir, cfg.Templates.ConstitutionPrompt),
		filepath.Join(templatesDir, cfg.Templates.SpecPrompt),
		filepath.Join(templatesDir, cfg.Templates.InspectPrompt),
		filepath.Join(templatesDir, cfg.Templates.PlanPrompt),
		filepath.Join(templatesDir, cfg.Templates.TasksPrompt),
		filepath.Join(templatesDir, cfg.Templates.ImplementPrompt),
		filepath.Join(templatesDir, cfg.Templates.ArchivePrompt),
		filepath.Join(templatesDir, cfg.Templates.VerifyPrompt),
		filepath.Join(scriptsDir, cfg.Scripts.CheckConstitution),
		filepath.Join(scriptsDir, cfg.Scripts.CheckSpecReady),
		filepath.Join(scriptsDir, cfg.Scripts.CheckInspectReady),
		filepath.Join(scriptsDir, cfg.Scripts.CheckPlanReady),
		filepath.Join(scriptsDir, cfg.Scripts.CheckTasksReady),
		filepath.Join(scriptsDir, cfg.Scripts.CheckImplementReady),
		filepath.Join(scriptsDir, cfg.Scripts.CheckArchiveReady),
		filepath.Join(scriptsDir, cfg.Scripts.CheckVerifyReady),
		filepath.Join(scriptsDir, cfg.Scripts.VerifyTaskState),
	} {
		checkPath(&findings, path, false)
	}

	if cfg.Language.Default != "en" && cfg.Language.Default != "ru" {
		findings = append(findings, Finding{Level: "error", Message: fmt.Sprintf("unsupported default language: %s", cfg.Language.Default)})
	}
	for _, value := range []string{cfg.Language.Docs, cfg.Language.Agent, cfg.Language.Comments} {
		if value != "en" && value != "ru" {
			findings = append(findings, Finding{Level: "error", Message: fmt.Sprintf("unsupported configured language: %s", value)})
		}
	}
	if _, err := config.NormalizeShell(cfg.Runtime.Shell); err != nil {
		findings = append(findings, Finding{Level: "error", Message: err.Error()})
	}

	enabledTargets := map[string]struct{}{}
	for _, target := range cfg.Agents.Targets {
		enabledTargets[target] = struct{}{}
		paths, err := agents.PathsForTarget(target)
		if err != nil {
			findings = append(findings, Finding{Level: "error", Message: err.Error()})
			continue
		}
		for _, relPath := range paths {
			checkPath(&findings, filepath.Join(root, filepath.FromSlash(relPath)), false)
		}
	}

	for _, target := range agents.SupportedTargets() {
		if _, ok := enabledTargets[target]; ok {
			continue
		}
		paths, err := agents.PathsForTarget(target)
		if err != nil {
			continue
		}
		for _, relPath := range paths {
			fullPath := filepath.Join(root, filepath.FromSlash(relPath))
			if _, err := os.Stat(fullPath); err == nil {
				findings = append(findings, Finding{Level: "warning", Message: fmt.Sprintf("orphaned agent artifact for disabled target %s: %s", target, fullPath)})
			}
		}
	}

	hasErrors := false
	for _, finding := range findings {
		if finding.Level == "error" {
			hasErrors = true
			break
		}
	}
	if !hasErrors {
		findings = append(findings, Finding{Level: "ok", Message: "draftspec workspace looks healthy"})
	}

	sort.Slice(findings, func(i, j int) bool {
		ri := severityRank(findings[i].Level)
		rj := severityRank(findings[j].Level)
		if ri != rj {
			return ri < rj
		}
		if findings[i].Message != findings[j].Message {
			return findings[i].Message < findings[j].Message
		}
		return findings[i].Level < findings[j].Level
	})
	return Result{Findings: findings}, nil
}

func checkPath(findings *[]Finding, path string, expectDir bool) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			*findings = append(*findings, Finding{Level: "error", Message: fmt.Sprintf("missing %s", path)})
			return
		}
		*findings = append(*findings, Finding{Level: "error", Message: fmt.Sprintf("failed to stat %s: %v", path, err)})
		return
	}
	if expectDir && !info.IsDir() {
		*findings = append(*findings, Finding{Level: "error", Message: fmt.Sprintf("expected directory: %s", path)})
		return
	}
	if !expectDir && info.IsDir() {
		*findings = append(*findings, Finding{Level: "error", Message: fmt.Sprintf("expected file: %s", path)})
	}
}

func severityRank(level string) int {
	switch level {
	case "error":
		return 0
	case "warning":
		return 1
	case "ok":
		return 2
	default:
		return 3
	}
}
