package workflow

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"draftspec/src/internal/config"
)

type CheckResult struct {
	Lines    []string
	Errors   int
	Warnings int
	Failed   bool
}

type TaskStateSummary struct {
	Total         int
	Completed     int
	Open          int
	TaskIDs       int
	CoverageLines int
}

type docSections struct {
	Goal         string
	Context      string
	Requirements string
	Acceptance   string
	Questions    string
	Coverage     string
}

var (
	taskIDPattern        = regexp.MustCompile(`T[0-9]+\.[0-9]+`)
	coverageLinePattern  = regexp.MustCompile(`AC-[0-9][0-9][0-9].*->.*T[0-9]+\.[0-9]+`)
	placeholderPattern   = regexp.MustCompile(`\[[A-Z0-9_][A-Z0-9_ -]*\]`)
	acceptanceIDPattern  = regexp.MustCompile(`AC-[0-9][0-9][0-9]`)
	requirementIDPattern = regexp.MustCompile(`RQ-[0-9][0-9][0-9]`)
	decisionIDPattern    = regexp.MustCompile(`DEC-[0-9][0-9][0-9]`)
)

func (r *CheckResult) AddOK(message string) {
	r.Lines = append(r.Lines, "OK: "+message)
}

func (r *CheckResult) AddWarn(message string) {
	r.Warnings++
	r.Lines = append(r.Lines, "WARN: "+message)
}

func (r *CheckResult) AddError(message string) {
	r.Errors++
	r.Failed = true
	r.Lines = append(r.Lines, "ERROR: "+message)
}

func (r *CheckResult) AddRaw(line string) {
	r.Lines = append(r.Lines, line)
}

func (r *CheckResult) Merge(other CheckResult) {
	r.Lines = append(r.Lines, other.Lines...)
	r.Errors += other.Errors
	r.Warnings += other.Warnings
	if other.Failed {
		r.Failed = true
	}
}

func CheckConstitution(root, constitutionPath string) (CheckResult, error) {
	root, cfg, err := loadCheckConfig(root)
	if err != nil {
		return CheckResult{}, err
	}

	displayPath := constitutionPath
	if strings.TrimSpace(displayPath) == "" {
		displayPath = cfg.Project.ConstitutionFile
	}
	displayPath, absolutePath := resolveUserPath(root, displayPath)

	result := CheckResult{}
	if !fileExists(absolutePath) {
		result.AddError(fmt.Sprintf("constitution file not found: %s", displayPath))
		result.Failed = false
		return result, nil
	}

	content, err := os.ReadFile(absolutePath)
	if err != nil {
		return CheckResult{}, fmt.Errorf("read constitution %s: %w", displayPath, err)
	}

	sections := constitutionSections(cfg.Language.Docs)
	for _, section := range sections {
		if hasHeading(string(content), section) {
			result.AddOK(section)
		} else {
			result.AddError(fmt.Sprintf("missing section: %s", section))
		}
	}

	principlesCount := countMatchingLines(string(content), `(?m)^### `)
	if principlesCount >= 5 {
		result.AddOK(fmt.Sprintf("principles count is %d", principlesCount))
	} else {
		result.AddError(fmt.Sprintf("expected at least 5 principles, found %d", principlesCount))
	}

	if placeholderPattern.Match(content) {
		result.AddWarn("placeholder tokens remain in constitution")
	} else {
		result.AddOK("no placeholder tokens detected")
	}

	// Preserve the current script contract: report issues, but do not fail the wrapper itself.
	result.Failed = false
	return result, nil
}

func CheckSpecReady(root string) (CheckResult, error) {
	root, cfg, err := loadCheckConfig(root)
	if err != nil {
		return CheckResult{}, err
	}

	result := CheckResult{}
	constitutionDisplay := cfg.Project.ConstitutionFile
	templateDisplay := joinDisplay(cfg.Paths.TemplatesDir, cfg.Templates.Spec)
	promptDisplay := joinDisplay(cfg.Paths.TemplatesDir, cfg.Templates.SpecPrompt)

	checkFile(&result, constitutionDisplay, absFromRoot(root, constitutionDisplay))
	checkFile(&result, templateDisplay, absFromRoot(root, templateDisplay))
	checkFile(&result, promptDisplay, absFromRoot(root, promptDisplay))

	templateAbs := absFromRoot(root, templateDisplay)
	if fileExists(templateAbs) {
		content, err := os.ReadFile(templateAbs)
		if err != nil {
			return CheckResult{}, fmt.Errorf("read spec template %s: %w", templateDisplay, err)
		}
		checkPattern(&result, string(content), `(?m)^## (Требования|Requirements)$`, "spec template has requirements section")
		checkPattern(&result, string(content), requirementIDPattern.String(), "spec template includes requirement IDs")
		checkPattern(&result, string(content), acceptanceIDPattern.String(), "spec template includes acceptance IDs")
		checkPattern(&result, string(content), `Given`, "spec template includes Given marker")
		checkPattern(&result, string(content), `When`, "spec template includes When marker")
		checkPattern(&result, string(content), `Then`, "spec template includes Then marker")
	}

	return result, nil
}

func CheckInspectReady(root, slug string) (CheckResult, error) {
	root, cfg, err := loadCheckConfig(root)
	if err != nil {
		return CheckResult{}, err
	}

	result := CheckResult{}
	specDisplay, specAbs := resolveSpecDisplayPath(root, cfg.Paths.SpecsDir, slug)
	reportTemplateDisplay := joinDisplay(cfg.Paths.TemplatesDir, cfg.Templates.InspectReport)
	promptDisplay := joinDisplay(cfg.Paths.TemplatesDir, cfg.Templates.InspectPrompt)

	checkFile(&result, cfg.Project.ConstitutionFile, absFromRoot(root, cfg.Project.ConstitutionFile))
	checkFile(&result, specDisplay, specAbs)
	checkFile(&result, reportTemplateDisplay, absFromRoot(root, reportTemplateDisplay))
	checkFile(&result, promptDisplay, absFromRoot(root, promptDisplay))

	if !result.Failed {
		inspectResult, err := InspectSpec(root, specDisplay, "")
		if err != nil {
			return CheckResult{}, err
		}
		result.Merge(inspectResult)
		if !inspectResult.Failed {
			result.AddOK(fmt.Sprintf("inspect is ready for slug '%s'", slug))
		}
	}

	return result, nil
}

func CheckPlanReady(root, slug string) (CheckResult, error) {
	root, cfg, err := loadCheckConfig(root)
	if err != nil {
		return CheckResult{}, err
	}

	result := CheckResult{}
	specDisplay, specAbs := resolveSpecDisplayPath(root, cfg.Paths.SpecsDir, slug)
	inspectDisplay, inspectAbs := resolveInspectDisplayPath(root, cfg.Paths.SpecsDir, slug)
	legacyInspectDisplay := joinDisplay(cfg.Paths.PlansDir, slug, "inspect.md")
	templateDisplay := joinDisplay(cfg.Paths.TemplatesDir, cfg.Templates.Plan)
	promptDisplay := joinDisplay(cfg.Paths.TemplatesDir, cfg.Templates.PlanPrompt)

	checkFile(&result, cfg.Project.ConstitutionFile, absFromRoot(root, cfg.Project.ConstitutionFile))
	checkFile(&result, specDisplay, specAbs)
	checkFile(&result, templateDisplay, absFromRoot(root, templateDisplay))
	checkFile(&result, promptDisplay, absFromRoot(root, promptDisplay))
	if fileExists(specAbs) {
		content, err := os.ReadFile(specAbs)
		if err != nil {
			return CheckResult{}, fmt.Errorf("read spec %s: %w", specDisplay, err)
		}
		checkPattern(&result, string(content), `(?m)^## (Критерии приемки|Acceptance Criteria)$`, "spec has acceptance criteria section")
		checkPattern(&result, string(content), acceptanceIDPattern.String(), "spec has stable acceptance IDs")
	}

	inspectDisplayPath := inspectDisplay
	if !fileExists(inspectAbs) {
		legacyAbs := absFromRoot(root, legacyInspectDisplay)
		if fileExists(legacyAbs) {
			inspectDisplayPath = legacyInspectDisplay
			inspectAbs = legacyAbs
		} else {
			result.AddError(fmt.Sprintf("missing inspect report %s", inspectDisplay))
		}
	}

	if fileExists(inspectAbs) {
		report, err := ParseReport(inspectAbs)
		if err != nil {
			result.AddError(err.Error())
		} else {
			if ValidStatus(report.Status) {
				result.AddOK("inspect report has a valid status")
			} else {
				result.AddError("inspect report has a valid status")
			}
			if report.Status == StatusBlocked {
				result.AddError("inspect report is blocked")
			}
		}
		_ = inspectDisplayPath
	}

	if fileExists(specAbs) {
		inspectResult, err := InspectSpec(root, specDisplay, "")
		if err != nil {
			return CheckResult{}, err
		}
		result.Merge(inspectResult)
	}

	return result, nil
}

func CheckTasksReady(root, slug string) (CheckResult, error) {
	root, cfg, err := loadCheckConfig(root)
	if err != nil {
		return CheckResult{}, err
	}

	result := CheckResult{}
	specDisplay, specAbs := resolveSpecDisplayPath(root, cfg.Paths.SpecsDir, slug)
	planDisplay := joinDisplay(cfg.Paths.PlansDir, slug, "plan.md")
	dataModelDisplay := joinDisplay(cfg.Paths.PlansDir, slug, "data-model.md")
	tasksTemplateDisplay := joinDisplay(cfg.Paths.TemplatesDir, cfg.Templates.Tasks)
	promptDisplay := joinDisplay(cfg.Paths.TemplatesDir, cfg.Templates.TasksPrompt)
	contractsDisplay := joinDisplay(cfg.Paths.PlansDir, slug, "contracts")

	checkFile(&result, cfg.Project.ConstitutionFile, absFromRoot(root, cfg.Project.ConstitutionFile))
	checkFile(&result, specDisplay, specAbs)
	checkFile(&result, planDisplay, absFromRoot(root, planDisplay))
	checkFile(&result, dataModelDisplay, absFromRoot(root, dataModelDisplay))
	checkFile(&result, tasksTemplateDisplay, absFromRoot(root, tasksTemplateDisplay))
	checkFile(&result, promptDisplay, absFromRoot(root, promptDisplay))

	if isDir(absFromRoot(root, contractsDisplay)) {
		result.AddOK(contractsDisplay)
	} else {
		result.AddOK("optional contracts directory not present")
	}

	if fileExists(specAbs) {
		content, err := os.ReadFile(specAbs)
		if err != nil {
			return CheckResult{}, fmt.Errorf("read spec %s: %w", specDisplay, err)
		}
		checkPattern(&result, string(content), acceptanceIDPattern.String(), "spec has stable acceptance IDs")

		inspectResult, err := InspectSpec(root, specDisplay, "")
		if err != nil {
			return CheckResult{}, err
		}
		result.Merge(inspectResult)
	}

	planAbs := absFromRoot(root, planDisplay)
	if fileExists(planAbs) {
		content, err := os.ReadFile(planAbs)
		if err != nil {
			return CheckResult{}, fmt.Errorf("read plan %s: %w", planDisplay, err)
		}
		checkPattern(&result, string(content), decisionIDPattern.String(), "plan has stable decision IDs")
	}

	return result, nil
}

func CheckImplementReady(root, slug string) (CheckResult, error) {
	root, cfg, err := loadCheckConfig(root)
	if err != nil {
		return CheckResult{}, err
	}

	result := CheckResult{}
	specDisplay, specAbs := resolveSpecDisplayPath(root, cfg.Paths.SpecsDir, slug)
	planDisplay := joinDisplay(cfg.Paths.PlansDir, slug, "plan.md")
	tasksDisplay := joinDisplay(cfg.Paths.PlansDir, slug, "tasks.md")
	dataModelDisplay := joinDisplay(cfg.Paths.PlansDir, slug, "data-model.md")
	promptDisplay := joinDisplay(cfg.Paths.TemplatesDir, cfg.Templates.ImplementPrompt)
	contractsDisplay := joinDisplay(cfg.Paths.PlansDir, slug, "contracts")

	checkFile(&result, cfg.Project.ConstitutionFile, absFromRoot(root, cfg.Project.ConstitutionFile))
	checkFile(&result, specDisplay, specAbs)
	checkFile(&result, planDisplay, absFromRoot(root, planDisplay))
	checkFile(&result, tasksDisplay, absFromRoot(root, tasksDisplay))
	checkFile(&result, dataModelDisplay, absFromRoot(root, dataModelDisplay))
	checkFile(&result, promptDisplay, absFromRoot(root, promptDisplay))

	if isDir(absFromRoot(root, contractsDisplay)) {
		result.AddOK(contractsDisplay)
	} else {
		result.AddOK("optional contracts directory not present")
	}

	tasksAbs := absFromRoot(root, tasksDisplay)
	if fileExists(tasksAbs) {
		content, err := os.ReadFile(tasksAbs)
		if err != nil {
			return CheckResult{}, fmt.Errorf("read tasks %s: %w", tasksDisplay, err)
		}
		checkPattern(&result, string(content), `(?m)^## (Покрытие критериев приемки|Acceptance Coverage)$`, "tasks include acceptance coverage section")
		checkPattern(&result, string(content), taskIDPattern.String(), "tasks include phase-scoped task IDs")
		checkPattern(&result, string(content), coverageLinePattern.String(), "tasks include AC-to-task coverage lines")
	}

	if fileExists(specAbs) && fileExists(tasksAbs) {
		inspectResult, err := InspectSpec(root, specDisplay, tasksDisplay)
		if err != nil {
			return CheckResult{}, err
		}
		result.Merge(inspectResult)
	}

	constitutionResult, err := CheckConstitution(root, cfg.Project.ConstitutionFile)
	if err != nil {
		return CheckResult{}, err
	}
	result.Lines = append(result.Lines, constitutionResult.Lines...)
	result.Warnings += constitutionResult.Warnings

	return result, nil
}

func CheckVerifyReady(root, slug string) (CheckResult, error) {
	root, cfg, err := loadCheckConfig(root)
	if err != nil {
		return CheckResult{}, err
	}

	result := CheckResult{}
	specDisplay, specAbs := resolveSpecDisplayPath(root, cfg.Paths.SpecsDir, slug)
	tasksDisplay := joinDisplay(cfg.Paths.PlansDir, slug, "tasks.md")
	reportTemplateDisplay := joinDisplay(cfg.Paths.TemplatesDir, cfg.Templates.VerifyReport)
	promptDisplay := joinDisplay(cfg.Paths.TemplatesDir, cfg.Templates.VerifyPrompt)

	checkFile(&result, cfg.Project.ConstitutionFile, absFromRoot(root, cfg.Project.ConstitutionFile))
	checkFile(&result, specDisplay, specAbs)
	checkFile(&result, tasksDisplay, absFromRoot(root, tasksDisplay))
	checkFile(&result, reportTemplateDisplay, absFromRoot(root, reportTemplateDisplay))
	checkFile(&result, promptDisplay, absFromRoot(root, promptDisplay))

	tasksAbs := absFromRoot(root, tasksDisplay)
	if fileExists(specAbs) && fileExists(tasksAbs) {
		inspectResult, err := InspectSpec(root, specDisplay, tasksDisplay)
		if err != nil {
			return CheckResult{}, err
		}
		result.Merge(inspectResult)
	}
	if fileExists(tasksAbs) {
		taskStateResult, _, err := VerifyTaskState(root, slug)
		if err != nil {
			return CheckResult{}, err
		}
		result.Merge(taskStateResult)
	}

	return result, nil
}

func CheckArchiveReady(root, slug, status, reason string) (CheckResult, error) {
	root, cfg, err := loadCheckConfig(root)
	if err != nil {
		return CheckResult{}, err
	}

	result := CheckResult{}
	if strings.TrimSpace(status) == "" {
		result.AddError("archive status is required")
		return result, nil
	}
	if strings.TrimSpace(reason) == "" {
		result.AddError("archive reason is required")
		return result, nil
	}
	switch status {
	case "completed", "superseded", "abandoned", "rejected", "deferred":
	default:
		result.AddError(fmt.Sprintf("invalid archive status: %s", status))
		return result, nil
	}

	specDisplay, specAbs := resolveSpecDisplayPath(root, cfg.Paths.SpecsDir, slug)
	if !fileExists(specAbs) {
		result.AddError(fmt.Sprintf("missing required file: %s", specDisplay))
		return result, nil
	}

	tasksDisplay := joinDisplay(cfg.Paths.PlansDir, slug, "tasks.md")
	tasksAbs := absFromRoot(root, tasksDisplay)
	if status == "completed" && fileExists(tasksAbs) {
		taskStateResult, summary, err := VerifyTaskState(root, slug)
		if err != nil {
			return CheckResult{}, err
		}
		result.Merge(taskStateResult)
		if summary.Open > 0 {
			result.AddError("completed archive requested while open tasks remain")
		}
	}

	if result.Failed {
		return result, nil
	}

	archiveDisplay := joinDisplay(cfg.Paths.ArchiveDir, slug)
	if err := os.MkdirAll(absFromRoot(root, archiveDisplay), 0o755); err != nil {
		return CheckResult{}, fmt.Errorf("create archive directory %s: %w", archiveDisplay, err)
	}
	result.AddOK(fmt.Sprintf("archive is ready for slug '%s' with status '%s'", slug, status))
	return result, nil
}

func VerifyTaskState(root, slug string) (CheckResult, TaskStateSummary, error) {
	root, cfg, err := loadCheckConfig(root)
	if err != nil {
		return CheckResult{}, TaskStateSummary{}, err
	}

	tasksDisplay := joinDisplay(cfg.Paths.PlansDir, slug, "tasks.md")
	tasksAbs := absFromRoot(root, tasksDisplay)
	result := CheckResult{}
	if !fileExists(tasksAbs) {
		result.AddError(fmt.Sprintf("missing %s", tasksDisplay))
		return result, TaskStateSummary{}, nil
	}

	summary, err := computeTaskState(tasksAbs)
	if err != nil {
		return CheckResult{}, TaskStateSummary{}, err
	}
	result.AddRaw(fmt.Sprintf("TASKS_TOTAL=%d", summary.Total))
	result.AddRaw(fmt.Sprintf("TASKS_COMPLETED=%d", summary.Completed))
	result.AddRaw(fmt.Sprintf("TASKS_OPEN=%d", summary.Open))
	result.AddRaw(fmt.Sprintf("TASK_IDS=%d", summary.TaskIDs))
	result.AddRaw(fmt.Sprintf("AC_COVERAGE_LINES=%d", summary.CoverageLines))

	if summary.Total == 0 {
		result.AddError(fmt.Sprintf("no task checkboxes found in %s", tasksDisplay))
		return result, summary, nil
	}
	if summary.TaskIDs == 0 {
		result.AddError(fmt.Sprintf("no stable task IDs found in %s", tasksDisplay))
		return result, summary, nil
	}
	if summary.CoverageLines == 0 {
		result.AddWarn(fmt.Sprintf("no AC-to-task coverage lines found in %s", tasksDisplay))
	}
	if summary.Open > 0 {
		result.AddWarn(fmt.Sprintf("open tasks remain in %s", tasksDisplay))
		return result, summary, nil
	}
	result.AddOK(fmt.Sprintf("all tasks are marked complete in %s", tasksDisplay))
	return result, summary, nil
}

func InspectSpec(root, specPath, tasksPath string) (CheckResult, error) {
	root, cfg, err := loadCheckConfig(root)
	if err != nil {
		return CheckResult{}, err
	}

	specDisplay, specAbs := resolveUserPath(root, specPath)
	result := CheckResult{}
	if !fileExists(specAbs) {
		result.AddError(fmt.Sprintf("spec file not found: %s", specDisplay))
		return result, nil
	}

	content, err := os.ReadFile(specAbs)
	if err != nil {
		return CheckResult{}, fmt.Errorf("read spec %s: %w", specDisplay, err)
	}
	sections := docsSections(cfg.Language.Docs)
	text := string(content)

	checkRequiredHeading(&result, text, sections.Goal)
	checkOptionalHeading(&result, text, sections.Context)
	checkRequiredHeading(&result, text, sections.Requirements)
	checkRequiredHeading(&result, text, sections.Acceptance)
	checkOptionalHeading(&result, text, sections.Questions)

	acceptanceBody := markdownSection(text, sections.Acceptance)
	if strings.TrimSpace(acceptanceBody) == "" {
		result.AddError("empty acceptance criteria section")
	} else {
		if strings.Contains(acceptanceBody, "Given") {
			result.AddOK("Given marker found")
		} else {
			result.AddError("missing Given marker in acceptance criteria")
		}
		if strings.Contains(acceptanceBody, "When") {
			result.AddOK("When marker found")
		} else {
			result.AddError("missing When marker in acceptance criteria")
		}
		if strings.Contains(acceptanceBody, "Then") {
			result.AddOK("Then marker found")
		} else {
			result.AddError("missing Then marker in acceptance criteria")
		}
	}

	criteriaCount := countMatchingLines(acceptanceBody, `(?m)^### `)
	if criteriaCount > 0 {
		result.AddOK(fmt.Sprintf("acceptance criteria count: %d", criteriaCount))
	} else {
		result.AddWarn("no explicit acceptance criterion headings found")
	}

	acceptanceIDCount := len(acceptanceIDPattern.FindAllString(acceptanceBody, -1))
	if acceptanceIDCount > 0 {
		result.AddOK(fmt.Sprintf("acceptance IDs found: %d", acceptanceIDCount))
	} else {
		result.AddWarn("no stable acceptance IDs found in acceptance criteria")
	}

	if strings.TrimSpace(tasksPath) != "" {
		tasksDisplay, tasksAbs := resolveUserPath(root, tasksPath)
		if fileExists(tasksAbs) {
			tasksContentBytes, err := os.ReadFile(tasksAbs)
			if err != nil {
				return CheckResult{}, fmt.Errorf("read tasks %s: %w", tasksDisplay, err)
			}
			tasksContent := string(tasksContentBytes)
			if hasHeading(tasksContent, sections.Coverage) {
				result.AddOK(sections.Coverage)
				coverageBody := markdownSection(tasksContent, sections.Coverage)
				coverageLines := countMatchingLines(coverageBody, `(?m)->`)
				malformedLines := countMalformedCoverageLines(coverageBody)
				if criteriaCount > 0 && coverageLines < criteriaCount {
					result.AddError(fmt.Sprintf("acceptance coverage entries (%d) are fewer than acceptance criteria (%d)", coverageLines, criteriaCount))
				} else {
					result.AddOK(fmt.Sprintf("acceptance coverage entries: %d", coverageLines))
				}
				if acceptanceIDCount > 0 && coverageLines < acceptanceIDCount {
					result.AddError(fmt.Sprintf("acceptance coverage entries (%d) are fewer than acceptance IDs (%d)", coverageLines, acceptanceIDCount))
				}
				if malformedLines > 0 {
					result.AddError("acceptance coverage contains malformed entries; expected lines like AC-001 -> T1.1")
				} else if coverageLines > 0 {
					result.AddOK("acceptance coverage format uses AC and task IDs")
				}
			} else {
				result.AddError(fmt.Sprintf("tasks file is missing required section: %s", sections.Coverage))
			}
		}
	}

	result.AddRaw(fmt.Sprintf("SUMMARY: errors=%d warnings=%d", result.Errors, result.Warnings))
	return result, nil
}

func loadCheckConfig(root string) (string, config.Config, error) {
	absoluteRoot, err := filepath.Abs(root)
	if err != nil {
		return "", config.Config{}, err
	}
	cfg, err := config.Load(absoluteRoot)
	if err != nil {
		return "", config.Config{}, err
	}
	return absoluteRoot, cfg, nil
}

func docsSections(language string) docSections {
	if strings.EqualFold(strings.TrimSpace(language), "ru") {
		return docSections{
			Goal:         "Цель",
			Context:      "Контекст",
			Requirements: "Требования",
			Acceptance:   "Критерии приемки",
			Questions:    "Открытые вопросы",
			Coverage:     "Покрытие критериев приемки",
		}
	}
	return docSections{
		Goal:         "Goal",
		Context:      "Context",
		Requirements: "Requirements",
		Acceptance:   "Acceptance Criteria",
		Questions:    "Open Questions",
		Coverage:     "Acceptance Coverage",
	}
}

func constitutionSections(language string) []string {
	if strings.EqualFold(strings.TrimSpace(language), "ru") {
		return []string{
			"Назначение",
			"Ключевые принципы",
			"Ограничения",
			"Языковая политика",
			"Процесс разработки",
			"Управление",
			"Последнее обновление",
		}
	}
	return []string{
		"Purpose",
		"Core Principles",
		"Constraints",
		"Language Policy",
		"Development Workflow",
		"Governance",
		"Last Updated",
	}
}

func resolveUserPath(root, value string) (string, string) {
	display := filepath.ToSlash(strings.TrimSpace(value))
	if filepath.IsAbs(value) {
		return display, value
	}
	return display, filepath.Join(root, filepath.FromSlash(value))
}

func joinDisplay(parts ...string) string {
	return filepath.ToSlash(filepath.Join(parts...))
}

func absFromRoot(root, rel string) string {
	if filepath.IsAbs(rel) {
		return rel
	}
	return filepath.Join(root, filepath.FromSlash(rel))
}

func resolveSpecDisplayPath(root, specsDir, slug string) (string, string) {
	display := joinDisplay(specsDir, slug, "spec.md")
	abs := absFromRoot(root, display)
	if fileExists(abs) {
		return display, abs
	}

	legacyDisplay := joinDisplay(specsDir, slug+".md")
	legacyAbs := absFromRoot(root, legacyDisplay)
	if fileExists(legacyAbs) {
		return legacyDisplay, legacyAbs
	}
	return display, abs
}

func resolveInspectDisplayPath(root, specsDir, slug string) (string, string) {
	display := joinDisplay(specsDir, slug, "inspect.md")
	abs := absFromRoot(root, display)
	if fileExists(abs) {
		return display, abs
	}

	legacyDisplay := joinDisplay(specsDir, slug+".inspect.md")
	legacyAbs := absFromRoot(root, legacyDisplay)
	if fileExists(legacyAbs) {
		return legacyDisplay, legacyAbs
	}
	return display, abs
}

func checkFile(result *CheckResult, displayPath, absolutePath string) {
	if fileExists(absolutePath) {
		result.AddOK(displayPath)
		return
	}
	result.AddError(fmt.Sprintf("missing %s", displayPath))
}

func checkPattern(result *CheckResult, content, pattern, label string) {
	if regexp.MustCompile(pattern).FindStringIndex(content) != nil {
		result.AddOK(label)
		return
	}
	result.AddError(label)
}

func hasHeading(content, section string) bool {
	return strings.Contains(content, "\n## "+section+"\n") || strings.HasPrefix(content, "## "+section+"\n") || strings.HasSuffix(content, "\n## "+section)
}

func checkRequiredHeading(result *CheckResult, content, section string) {
	if hasHeading(content, section) {
		result.AddOK(section)
	} else {
		result.AddError(fmt.Sprintf("missing required section: %s", section))
	}
}

func checkOptionalHeading(result *CheckResult, content, section string) {
	if hasHeading(content, section) {
		result.AddOK(section)
	} else {
		result.AddWarn(fmt.Sprintf("missing section: %s", section))
	}
}

func markdownSection(content, section string) string {
	lines := strings.Split(content, "\n")
	var captured []string
	inSection := false
	target := "## " + section
	for _, line := range lines {
		if line == target {
			inSection = true
			continue
		}
		if inSection && strings.HasPrefix(line, "## ") {
			break
		}
		if inSection {
			captured = append(captured, line)
		}
	}
	return strings.Join(captured, "\n")
}

func countMatchingLines(content, pattern string) int {
	return len(regexp.MustCompile(pattern).FindAllString(content, -1))
}

func countMalformedCoverageLines(content string) int {
	lines := strings.Split(content, "\n")
	count := 0
	for _, line := range lines {
		if !strings.Contains(line, "->") {
			continue
		}
		if !coverageLinePattern.MatchString(line) {
			count++
		}
	}
	return count
}

func computeTaskState(tasksPath string) (TaskStateSummary, error) {
	content, err := os.ReadFile(tasksPath)
	if err != nil {
		return TaskStateSummary{}, fmt.Errorf("read tasks file: %w", err)
	}

	text := string(content)
	total, completed, open, err := taskCounts(tasksPath)
	if err != nil {
		return TaskStateSummary{}, err
	}
	return TaskStateSummary{
		Total:         total,
		Completed:     completed,
		Open:          open,
		TaskIDs:       len(taskIDPattern.FindAllString(text, -1)),
		CoverageLines: len(coverageLinePattern.FindAllString(text, -1)),
	}, nil
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
