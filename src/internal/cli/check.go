package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"draftspec/src/internal/workflow"
	"github.com/spf13/cobra"
)

type checkResult struct {
	Slug        string         `json:"slug"`
	Phase       string         `json:"phase"`
	Verdict     string         `json:"verdict"`
	NextCommand string         `json:"next_command,omitempty"`
	Blocked     bool           `json:"blocked"`
	Artifacts   checkArtifacts `json:"artifacts"`
}

type checkArtifacts struct {
	Spec    checkArtifact `json:"spec"`
	Inspect checkArtifact `json:"inspect"`
	Plan    checkArtifact `json:"plan"`
	Tasks   checkArtifact `json:"tasks"`
	Verify  checkArtifact `json:"verify"`
}

type checkArtifact struct {
	Present bool   `json:"present"`
	Detail  string `json:"detail,omitempty"`
}

type checkAllResult struct {
	Blocked  bool          `json:"blocked"`
	Features []checkResult `json:"features"`
}

func newCheckCmd() *cobra.Command {
	var jsonOutput bool
	var allFeatures bool

	cmd := &cobra.Command{
		Use:   "check <slug> [path]",
		Short: "Check feature readiness and show next action",
		Args: func(cmd *cobra.Command, args []string) error {
			all, _ := cmd.Flags().GetBool("all")
			if all {
				if len(args) > 1 {
					return fmt.Errorf("accepts at most 1 arg (path) when --all is set, received %d", len(args))
				}
				return nil
			}
			if len(args) < 1 || len(args) > 2 {
				return fmt.Errorf("accepts 1 or 2 args (slug [path]), received %d", len(args))
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if allFeatures {
				root := "."
				if len(args) == 1 {
					root = args[0]
				}
				return runCheckAll(cmd, root, jsonOutput)
			}

			root := "."
			if len(args) == 2 {
				root = args[1]
			}

			state, err := workflow.State(root, args[0])
			if err != nil {
				return err
			}

			result := buildCheckResult(state)

			if jsonOutput {
				payload, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(payload))
				if result.Blocked {
					return fmt.Errorf("blocked")
				}
				return nil
			}

			printCheck(cmd, state, result)
			if result.Blocked {
				return fmt.Errorf("blocked")
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "output as JSON; exits with code 1 when blocked")
	cmd.Flags().BoolVar(&allFeatures, "all", false, "check all features; exits with code 1 if any are blocked")
	return cmd
}

func runCheckAll(cmd *cobra.Command, root string, jsonOutput bool) error {
	states, err := workflow.States(root)
	if err != nil {
		return err
	}

	results := make([]checkResult, 0, len(states))
	anyBlocked := false
	for _, state := range states {
		r := buildCheckResult(state)
		results = append(results, r)
		if r.Blocked {
			anyBlocked = true
		}
	}

	if jsonOutput {
		payload, err := json.MarshalIndent(checkAllResult{Blocked: anyBlocked, Features: results}, "", "  ")
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), string(payload))
		if anyBlocked {
			return fmt.Errorf("blocked")
		}
		return nil
	}

	printCheckAll(cmd, states, results)
	if anyBlocked {
		return fmt.Errorf("blocked")
	}
	return nil
}

func buildCheckResult(state workflow.FeatureState) checkResult {
	result := checkResult{
		Slug:    state.Slug,
		Phase:   state.Phase,
		Blocked: state.Blocked,
		Artifacts: checkArtifacts{
			Spec:    checkArtifact{Present: state.SpecExists},
			Inspect: checkArtifact{Present: state.InspectExists},
			Plan:    checkArtifact{Present: state.PlanExists},
			Tasks:   checkArtifact{Present: state.TasksExists},
			Verify:  checkArtifact{Present: state.VerifyExists},
		},
	}

	if state.InspectExists && state.InspectStatus != "" {
		result.Artifacts.Inspect.Detail = state.InspectStatus
	}

	if state.TasksExists {
		result.Artifacts.Tasks.Detail = fmt.Sprintf("%d/%d done", state.TasksCompleted, state.TasksTotal)
	}

	if state.VerifyExists && state.VerifyStatus != "" {
		result.Artifacts.Verify.Detail = state.VerifyStatus
	}

	if state.Blocked {
		result.Verdict = "blocked"
	} else {
		result.Verdict = "ready"
	}

	result.NextCommand = nextCommand(state)
	return result
}

func nextCommand(state workflow.FeatureState) string {
	switch state.ReadyFor {
	case "spec":
		return "/draftspec.spec " + state.Slug
	case "inspect":
		return "/draftspec.inspect " + state.Slug
	case "plan":
		return "/draftspec.plan " + state.Slug
	case "tasks":
		return "/draftspec.tasks " + state.Slug
	case "implement":
		return "/draftspec.implement " + state.Slug
	case "verify":
		return "/draftspec.verify " + state.Slug
	case "archive":
		return "/draftspec.archive " + state.Slug
	default:
		return ""
	}
}

func printCheck(cmd *cobra.Command, state workflow.FeatureState, result checkResult) {
	w := cmd.OutOrStdout()

	fmt.Fprintf(w, "feature:  %s\n", state.Slug)
	fmt.Fprintf(w, "phase:    %s\n", state.Phase)
	fmt.Fprintln(w)

	fmt.Fprintln(w, "artifacts:")
	fmt.Fprintf(w, "  spec      %s\n", artifactLine(state.SpecExists, ""))
	fmt.Fprintf(w, "  inspect   %s\n", artifactLine(state.InspectExists, state.InspectStatus))

	tasksDetail := ""
	if state.TasksExists {
		tasksDetail = fmt.Sprintf("%d/%d done", state.TasksCompleted, state.TasksTotal)
		if state.TasksOpen > 0 {
			tasksDetail += fmt.Sprintf("  (%d open)", state.TasksOpen)
		}
	}
	fmt.Fprintf(w, "  plan      %s\n", artifactLine(state.PlanExists, ""))
	fmt.Fprintf(w, "  tasks     %s\n", artifactLine(state.TasksExists, tasksDetail))
	fmt.Fprintf(w, "  verify    %s\n", artifactLine(state.VerifyExists, state.VerifyStatus))

	fmt.Fprintln(w)

	if state.Blocked {
		fmt.Fprintf(w, "verdict:  blocked\n")
	} else {
		fmt.Fprintf(w, "verdict:  ready\n")
	}

	if result.NextCommand != "" {
		fmt.Fprintf(w, "next:     %s\n", result.NextCommand)
	}

	if state.InspectLegacy {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "warning:  inspect report is at legacy path — run `draftspec feature repair "+state.Slug+"`")
	}
}

func artifactLine(present bool, detail string) string {
	if !present {
		return "-  missing"
	}
	parts := []string{"✓"}
	if detail != "" {
		parts = append(parts, " ", detail)
	}
	return strings.Join(parts, "")
}

func printCheckAll(cmd *cobra.Command, states []workflow.FeatureState, results []checkResult) {
	w := cmd.OutOrStdout()

	if len(results) == 0 {
		fmt.Fprintln(w, "no features found")
		return
	}

	// Compute column widths.
	slugWidth := len("feature")
	phaseWidth := len("phase")
	for _, r := range results {
		if len(r.Slug) > slugWidth {
			slugWidth = len(r.Slug)
		}
		if len(r.Phase) > phaseWidth {
			phaseWidth = len(r.Phase)
		}
	}

	format := fmt.Sprintf("%%-%ds  %%-%ds  %%-8s  %%s\n", slugWidth, phaseWidth)
	fmt.Fprintf(w, format, "feature", "phase", "verdict", "next")
	fmt.Fprintf(w, format,
		strings.Repeat("-", slugWidth),
		strings.Repeat("-", phaseWidth),
		"-------",
		"----",
	)

	blockedCount := 0
	for i, r := range results {
		_ = states[i]
		verdict := "ready"
		if r.Blocked {
			verdict = "blocked"
			blockedCount++
		}
		fmt.Fprintf(w, format, r.Slug, r.Phase, verdict, r.NextCommand)
	}

	fmt.Fprintln(w)
	if blockedCount > 0 {
		fmt.Fprintf(w, "verdict:  %d of %d features blocked\n", blockedCount, len(results))
	} else {
		fmt.Fprintf(w, "verdict:  all %d features ready\n", len(results))
	}
}
