package githubapi

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestFetchBuildsRunQueryAndRespectsWorkflowFilter(t *testing.T) {
	t.Setenv("GH_FLAKEFINDER_FAKE_GH", "1")
	t.Setenv("GH_FLAKEFINDER_EXPECT_CREATED", "created=>=2026-06-06")
	client := GHClient{
		Command: os.Args[0],
		Now: func() time.Time {
			return time.Date(2026, 6, 13, 12, 0, 0, 0, time.UTC)
		},
	}

	data, err := client.Fetch(context.Background(), "owner/repo", FetchOptions{Days: 7, Limit: 1, Workflow: "CI"})

	if err != nil {
		t.Fatalf("Fetch returned error: %v", err)
	}
	if len(data.Runs) != 1 {
		t.Fatalf("expected one filtered run, got %d", len(data.Runs))
	}
	if data.Runs[0].Name != "CI" || data.Runs[0].WorkflowID != 10 {
		t.Fatalf("unexpected run: %#v", data.Runs[0])
	}
}

func TestFetchJobsPaginatesRunJobsEndpoint(t *testing.T) {
	t.Setenv("GH_FLAKEFINDER_FAKE_GH", "1")
	client := GHClient{Command: os.Args[0]}

	jobs, err := client.fetchJobs(context.Background(), "owner/repo", WorkflowRun{ID: 1234, Attempt: 1})

	if err != nil {
		t.Fatalf("fetchJobs returned error: %v", err)
	}
	if len(jobs) != 101 {
		t.Fatalf("expected 101 jobs across two pages, got %d", len(jobs))
	}
}

func TestFetchJobsPaginatesAttemptJobsEndpoint(t *testing.T) {
	t.Setenv("GH_FLAKEFINDER_FAKE_GH", "1")
	client := GHClient{Command: os.Args[0]}

	jobs, err := client.fetchJobs(context.Background(), "owner/repo", WorkflowRun{ID: 5678, Attempt: 2})

	if err != nil {
		t.Fatalf("fetchJobs returned error: %v", err)
	}
	if len(jobs) != 202 {
		t.Fatalf("expected 202 jobs across two pages for two attempts, got %d", len(jobs))
	}
	attempts := map[int]int{}
	for _, job := range jobs {
		attempts[job.Attempt]++
	}
	if attempts[1] != 101 || attempts[2] != 101 {
		t.Fatalf("expected 101 jobs per attempt, got %#v", attempts)
	}
}

func TestMain(m *testing.M) {
	if os.Getenv("GH_FLAKEFINDER_FAKE_GH") == "1" {
		fakeGHAPI()
		return
	}
	os.Exit(m.Run())
}

func fakeGHAPI() {
	endpoint := ""
	page := 1
	for i, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "/repos/") {
			endpoint = arg
		}
		if arg == "-f" && i+2 <= len(os.Args[1:]) {
			keyValue := os.Args[1:][i+1]
			if strings.HasPrefix(keyValue, "page=") {
				if parsed, err := strconv.Atoi(strings.TrimPrefix(keyValue, "page=")); err == nil {
					page = parsed
				}
			}
		}
	}
	if strings.Contains(endpoint, "/jobs") {
		writeFakeJobs(endpoint, page)
		return
	}
	if strings.Contains(endpoint, "/actions/runs") {
		writeFakeRuns(page)
		return
	}
	fmt.Fprintln(os.Stderr, "unexpected fake gh endpoint:", endpoint)
	os.Exit(1)
}

func writeFakeRuns(page int) {
	if page != 1 {
		writeJSON(struct {
			WorkflowRuns []WorkflowRun `json:"workflow_runs"`
		}{})
		return
	}
	expectedCreated := os.Getenv("GH_FLAKEFINDER_EXPECT_CREATED")
	for _, arg := range os.Args[1:] {
		if expectedCreated != "" && arg == expectedCreated {
			expectedCreated = ""
			break
		}
	}
	if expectedCreated != "" {
		fmt.Fprintln(os.Stderr, "missing expected created query:", expectedCreated)
		os.Exit(1)
	}
	writeJSON(struct {
		WorkflowRuns []WorkflowRun `json:"workflow_runs"`
	}{
		WorkflowRuns: []WorkflowRun{
			{ID: 7777, Attempt: 1, Name: "CI", WorkflowID: 10, Event: "push", HeadSHA: "abc", HeadBranch: "main", Status: "completed", Conclusion: "success", HTMLURL: "https://github.com/owner/repo/actions/runs/7777"},
		},
	})
}

func writeFakeJobs(endpoint string, page int) {
	runID := int64(1234)
	attempt := 1
	if strings.Contains(endpoint, "/runs/5678/") {
		runID = 5678
	}
	if strings.Contains(endpoint, "/attempts/2/") {
		attempt = 2
	}
	count := 100
	if page == 2 {
		count = 1
	}
	if page > 2 {
		count = 0
	}
	jobs := make([]WorkflowJob, 0, count)
	for i := 0; i < count; i++ {
		jobs = append(jobs, WorkflowJob{
			ID:         int64(attempt*100000 + page*1000 + i),
			RunID:      runID,
			Attempt:    attempt,
			Name:       fmt.Sprintf("job-%d-%d-%d", attempt, page, i),
			Status:     "completed",
			Conclusion: "success",
			HTMLURL:    fmt.Sprintf("https://github.com/owner/repo/actions/runs/%d/job/%d", runID, i),
		})
	}
	writeJSON(struct {
		Jobs []WorkflowJob `json:"jobs"`
	}{Jobs: jobs})
}

func writeJSON(value any) {
	if err := json.NewEncoder(os.Stdout).Encode(value); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
