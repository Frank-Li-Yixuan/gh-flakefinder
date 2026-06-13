package detect

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Frank-Li-Yixuan/gh-flakefinder/internal/githubapi"
)

type Options struct {
	MinConfidence float64
	Workflow      string
}

type EvidenceLink struct {
	Label      string `json:"label"`
	URL        string `json:"url"`
	Conclusion string `json:"conclusion,omitempty"`
	RunID      int64  `json:"run_id,omitempty"`
	JobID      int64  `json:"job_id,omitempty"`
}

type Suspect struct {
	WorkflowName string         `json:"workflow_name"`
	WorkflowID   int64          `json:"workflow_id"`
	JobName      string         `json:"job_name,omitempty"`
	HeadSHA      string         `json:"head_sha"`
	Branch       string         `json:"branch"`
	Event        string         `json:"event"`
	Signal       string         `json:"signal"`
	Category     string         `json:"category"`
	Confidence   float64        `json:"confidence"`
	Attempts     int            `json:"attempts_count"`
	Evidence     []EvidenceLink `json:"evidence"`
	Suggestion   string         `json:"suggestion"`
}

func FindSuspects(data githubapi.DataSet, opts Options) []Suspect {
	if opts.MinConfidence <= 0 {
		opts.MinConfidence = 0.60
	}
	runByID := make(map[int64]githubapi.WorkflowRun, len(data.Runs))
	for _, run := range data.Runs {
		if isCompleted(run.Status) && workflowMatches(run, opts.Workflow) {
			runByID[run.ID] = run
		}
	}

	var suspects []Suspect
	suspects = append(suspects, jobSuspects(data.Jobs, runByID, opts.MinConfidence)...)
	suspects = append(suspects, workflowSuspects(runByID, opts.MinConfidence)...)

	sort.SliceStable(suspects, func(i, j int) bool {
		iRun := firstRunID(suspects[i])
		jRun := firstRunID(suspects[j])
		if iRun != jRun {
			return iRun < jRun
		}
		iJob := suspects[i].JobName != ""
		jJob := suspects[j].JobName != ""
		if iJob != jJob {
			return iJob
		}
		if suspects[i].WorkflowName != suspects[j].WorkflowName {
			return suspects[i].WorkflowName < suspects[j].WorkflowName
		}
		return suspects[i].JobName < suspects[j].JobName
	})
	return suspects
}

func jobSuspects(jobs []githubapi.WorkflowJob, runByID map[int64]githubapi.WorkflowRun, minConfidence float64) []Suspect {
	type jobWithRun struct {
		job githubapi.WorkflowJob
		run githubapi.WorkflowRun
	}
	groups := make(map[string][]jobWithRun)
	for _, job := range jobs {
		run, ok := runByID[job.RunID]
		if !ok || !isCompleted(job.Status) {
			continue
		}
		key := strings.Join([]string{
			run.HeadSHA,
			fmt.Sprintf("%d", run.WorkflowID),
			strings.ToLower(job.Name),
		}, "\x00")
		groups[key] = append(groups[key], jobWithRun{job: job, run: run})
	}

	var suspects []Suspect
	for _, group := range groups {
		sort.Slice(group, func(i, j int) bool {
			if !group[i].job.StartedAt.Equal(group[j].job.StartedAt) {
				return group[i].job.StartedAt.Before(group[j].job.StartedAt)
			}
			return group[i].job.ID < group[j].job.ID
		})
		failed, passed, ok := firstFailureThenSuccess(group, func(item jobWithRun) string {
			return item.job.Conclusion
		})
		if !ok {
			continue
		}
		confidence := confidenceFor(failed.job.Conclusion)
		if normalizeConclusion(failed.job.Conclusion) == "failure" {
			confidence = 0.90
		}
		if failed.job.RunID == passed.job.RunID {
			confidence = 0.95
		}
		if confidence < minConfidence {
			continue
		}
		suspects = append(suspects, Suspect{
			WorkflowName: failed.run.Name,
			WorkflowID:   failed.run.WorkflowID,
			JobName:      failed.job.Name,
			HeadSHA:      failed.run.HeadSHA,
			Branch:       failed.run.HeadBranch,
			Event:        failed.run.Event,
			Signal:       signal(failed.job.Conclusion, passed.job.Conclusion),
			Category:     categoryFor(failed.job.Conclusion),
			Confidence:   confidence,
			Attempts:     len(group),
			Evidence: []EvidenceLink{
				{Label: "Failed", URL: failed.job.HTMLURL, Conclusion: failed.job.Conclusion, RunID: failed.job.RunID, JobID: failed.job.ID},
				{Label: "Passed after rerun", URL: passed.job.HTMLURL, Conclusion: passed.job.Conclusion, RunID: passed.job.RunID, JobID: passed.job.ID},
			},
			Suggestion: suggestionFor(failed.job.Conclusion),
		})
	}
	return suspects
}

func workflowSuspects(runByID map[int64]githubapi.WorkflowRun, minConfidence float64) []Suspect {
	groups := make(map[string][]githubapi.WorkflowRun)
	for _, run := range runByID {
		key := strings.Join([]string{
			run.HeadSHA,
			fmt.Sprintf("%d", run.WorkflowID),
			run.Event,
			run.HeadBranch,
		}, "\x00")
		groups[key] = append(groups[key], run)
	}

	var suspects []Suspect
	for _, group := range groups {
		sort.Slice(group, func(i, j int) bool {
			if !group[i].CreatedAt.Equal(group[j].CreatedAt) {
				return group[i].CreatedAt.Before(group[j].CreatedAt)
			}
			return group[i].ID < group[j].ID
		})
		failed, passed, ok := firstFailureThenSuccess(group, func(item githubapi.WorkflowRun) string {
			return item.Conclusion
		})
		if !ok {
			continue
		}
		confidence := confidenceFor(failed.Conclusion)
		if failed.ID == passed.ID {
			confidence = 0.95
		}
		if confidence < minConfidence {
			continue
		}
		suspects = append(suspects, Suspect{
			WorkflowName: failed.Name,
			WorkflowID:   failed.WorkflowID,
			HeadSHA:      failed.HeadSHA,
			Branch:       failed.HeadBranch,
			Event:        failed.Event,
			Signal:       signal(failed.Conclusion, passed.Conclusion),
			Category:     categoryFor(failed.Conclusion),
			Confidence:   confidence,
			Attempts:     len(group),
			Evidence: []EvidenceLink{
				{Label: "Failed", URL: failed.HTMLURL, Conclusion: failed.Conclusion, RunID: failed.ID},
				{Label: "Passed after rerun", URL: passed.HTMLURL, Conclusion: passed.Conclusion, RunID: passed.ID},
			},
			Suggestion: suggestionFor(failed.Conclusion),
		})
	}
	return suspects
}

func firstFailureThenSuccess[T any](items []T, conclusion func(T) string) (T, T, bool) {
	var zero T
	var failed T
	foundFailure := false
	for _, item := range items {
		c := normalizeConclusion(conclusion(item))
		if !foundFailure && isFailureSignal(c) {
			failed = item
			foundFailure = true
			continue
		}
		if foundFailure && c == "success" {
			return failed, item, true
		}
	}
	return zero, zero, false
}

func confidenceFor(conclusion string) float64 {
	switch normalizeConclusion(conclusion) {
	case "failure":
		return 0.85
	case "timed_out":
		return 0.70
	case "cancelled":
		return 0.45
	default:
		return 0.60
	}
}

func categoryFor(conclusion string) string {
	switch normalizeConclusion(conclusion) {
	case "timed_out":
		return "infra/network"
	case "cancelled":
		return "manual/cancelled"
	default:
		return "flaky-test"
	}
}

func suggestionFor(conclusion string) string {
	switch normalizeConclusion(conclusion) {
	case "timed_out":
		return "Check runner capacity, network calls, dependency downloads, and external service timeouts."
	case "cancelled":
		return "Confirm whether cancellation was manual before treating this as a flaky signal."
	default:
		return "Inspect test isolation, shared state, and network or dependency calls in this job."
	}
}

func signal(from, to string) string {
	return normalizeConclusion(from) + " -> " + normalizeConclusion(to)
}

func isFailureSignal(conclusion string) bool {
	switch normalizeConclusion(conclusion) {
	case "failure", "timed_out", "cancelled":
		return true
	default:
		return false
	}
}

func normalizeConclusion(conclusion string) string {
	return strings.ToLower(strings.TrimSpace(conclusion))
}

func isCompleted(status string) bool {
	status = strings.ToLower(strings.TrimSpace(status))
	return status == "" || status == "completed"
}

func workflowMatches(run githubapi.WorkflowRun, filter string) bool {
	if filter == "" {
		return true
	}
	return strings.EqualFold(run.Name, filter) || fmt.Sprintf("%d", run.WorkflowID) == filter
}

func firstRunID(s Suspect) int64 {
	if len(s.Evidence) == 0 {
		return 0
	}
	return s.Evidence[0].RunID
}
