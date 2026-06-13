package detect_test

import (
	"testing"
	"time"

	"github.com/Frank-Li-Yixuan/gh-flakefinder/internal/detect"
	"github.com/Frank-Li-Yixuan/gh-flakefinder/internal/githubapi"
)

func TestDetectFindsJobAndWorkflowFlakeSignals(t *testing.T) {
	data := githubapi.DataSet{
		Repository: "owner/repo",
		Runs: []githubapi.WorkflowRun{
			run(1001, 1, "CI", 10, "push", "abc123456789", "main", "failure", "https://github.com/owner/repo/actions/runs/1001"),
			run(1002, 1, "CI", 10, "push", "abc123456789", "main", "success", "https://github.com/owner/repo/actions/runs/1002"),
			run(2001, 1, "Build", 20, "pull_request", "def987654321", "feature", "timed_out", "https://github.com/owner/repo/actions/runs/2001"),
			run(2002, 1, "Build", 20, "pull_request", "def987654321", "feature", "success", "https://github.com/owner/repo/actions/runs/2002"),
		},
		Jobs: []githubapi.WorkflowJob{
			job(5001, 1001, 1, "test-node-20-linux", "failure", "https://github.com/owner/repo/actions/runs/1001/job/5001"),
			job(5002, 1002, 1, "test-node-20-linux", "success", "https://github.com/owner/repo/actions/runs/1002/job/5002"),
			job(6001, 2001, 1, "integration", "timed_out", "https://github.com/owner/repo/actions/runs/2001/job/6001"),
			job(6002, 2002, 1, "integration", "success", "https://github.com/owner/repo/actions/runs/2002/job/6002"),
		},
	}

	suspects := detect.FindSuspects(data, detect.Options{MinConfidence: 0.60})

	if len(suspects) != 4 {
		t.Fatalf("expected 4 suspects, got %d: %#v", len(suspects), suspects)
	}
	assertSuspect(t, suspects[0], "CI", "test-node-20-linux", "failure -> success", 0.90)
	assertSuspect(t, suspects[1], "CI", "", "failure -> success", 0.85)
	assertSuspect(t, suspects[2], "Build", "integration", "timed_out -> success", 0.70)
	assertSuspect(t, suspects[3], "Build", "", "timed_out -> success", 0.70)
}

func TestDetectFindsSameRunJobRerunAcrossAttempts(t *testing.T) {
	data := githubapi.DataSet{
		Repository: "owner/repo",
		Runs: []githubapi.WorkflowRun{
			run(9001, 1, "CI", 10, "pull_request", "feedface1234567890", "feature", "failure", "https://github.com/owner/repo/actions/runs/9001/attempts/1"),
			run(9001, 2, "CI", 10, "pull_request", "feedface1234567890", "feature", "success", "https://github.com/owner/repo/actions/runs/9001/attempts/2"),
		},
		Jobs: []githubapi.WorkflowJob{
			job(9101, 9001, 1, "test-rerun", "failure", "https://github.com/owner/repo/actions/runs/9001/job/9101"),
			job(9102, 9001, 2, "test-rerun", "success", "https://github.com/owner/repo/actions/runs/9001/job/9102"),
		},
	}

	suspects := detect.FindSuspects(data, detect.Options{MinConfidence: 0.91})

	if len(suspects) != 1 {
		t.Fatalf("expected exactly one high-confidence same-run job suspect, got %d: %#v", len(suspects), suspects)
	}
	assertSuspect(t, suspects[0], "CI", "test-rerun", "failure -> success", 0.95)
	if suspects[0].Evidence[0].RunID != suspects[0].Evidence[1].RunID {
		t.Fatalf("expected same workflow run ID evidence, got %#v", suspects[0].Evidence)
	}
}

func TestDetectFiltersCancelledByDefaultConfidence(t *testing.T) {
	data := githubapi.DataSet{
		Repository: "owner/repo",
		Runs: []githubapi.WorkflowRun{
			run(1001, 1, "CI", 10, "push", "abc123456789", "main", "cancelled", "https://github.com/owner/repo/actions/runs/1001"),
			run(1002, 1, "CI", 10, "push", "abc123456789", "main", "success", "https://github.com/owner/repo/actions/runs/1002"),
		},
	}

	suspects := detect.FindSuspects(data, detect.Options{MinConfidence: 0.60})

	if len(suspects) != 0 {
		t.Fatalf("expected cancelled signal to be filtered, got %#v", suspects)
	}
}

func TestDetectWorkflowNameFilter(t *testing.T) {
	data := githubapi.DataSet{
		Repository: "owner/repo",
		Runs: []githubapi.WorkflowRun{
			run(1001, 1, "CI", 10, "push", "abc123456789", "main", "failure", "https://github.com/owner/repo/actions/runs/1001"),
			run(1002, 1, "CI", 10, "push", "abc123456789", "main", "success", "https://github.com/owner/repo/actions/runs/1002"),
			run(2001, 1, "Deploy", 20, "push", "abc123456789", "main", "failure", "https://github.com/owner/repo/actions/runs/2001"),
			run(2002, 1, "Deploy", 20, "push", "abc123456789", "main", "success", "https://github.com/owner/repo/actions/runs/2002"),
		},
	}

	suspects := detect.FindSuspects(data, detect.Options{MinConfidence: 0.60, Workflow: "Deploy"})

	if len(suspects) != 1 {
		t.Fatalf("expected 1 filtered suspect, got %d: %#v", len(suspects), suspects)
	}
	if suspects[0].WorkflowName != "Deploy" {
		t.Fatalf("expected Deploy suspect, got %q", suspects[0].WorkflowName)
	}
}

func TestDetectDoesNotFlagOnlyFailuresOrOnlySuccesses(t *testing.T) {
	data := githubapi.DataSet{
		Repository: "owner/repo",
		Runs: []githubapi.WorkflowRun{
			run(1001, 1, "CI", 10, "push", "abc123456789", "main", "failure", "https://github.com/owner/repo/actions/runs/1001"),
			run(1002, 1, "CI", 10, "push", "abc123456789", "main", "failure", "https://github.com/owner/repo/actions/runs/1002"),
			run(2001, 1, "Build", 20, "push", "def987654321", "main", "success", "https://github.com/owner/repo/actions/runs/2001"),
			run(2002, 1, "Build", 20, "push", "def987654321", "main", "success", "https://github.com/owner/repo/actions/runs/2002"),
		},
	}

	suspects := detect.FindSuspects(data, detect.Options{MinConfidence: 0.60})

	if len(suspects) != 0 {
		t.Fatalf("expected no suspects, got %#v", suspects)
	}
}

func assertSuspect(t *testing.T, got detect.Suspect, workflow, jobName, signal string, confidence float64) {
	t.Helper()
	if got.WorkflowName != workflow || got.JobName != jobName || got.Signal != signal || got.Confidence != confidence {
		t.Fatalf("unexpected suspect: %#v", got)
	}
	if len(got.Evidence) != 2 {
		t.Fatalf("expected 2 evidence links, got %#v", got.Evidence)
	}
}

func run(id int64, attempt int, name string, workflowID int64, event, sha, branch, conclusion, url string) githubapi.WorkflowRun {
	return githubapi.WorkflowRun{
		ID:         id,
		Attempt:    attempt,
		Name:       name,
		WorkflowID: workflowID,
		Event:      event,
		HeadSHA:    sha,
		HeadBranch: branch,
		Status:     "completed",
		Conclusion: conclusion,
		CreatedAt:  time.Unix(id, 0).UTC(),
		UpdatedAt:  time.Unix(id+60, 0).UTC(),
		HTMLURL:    url,
	}
}

func job(id, runID int64, attempt int, name, conclusion, url string) githubapi.WorkflowJob {
	return githubapi.WorkflowJob{
		ID:          id,
		RunID:       runID,
		Attempt:     attempt,
		Name:        name,
		Status:      "completed",
		Conclusion:  conclusion,
		StartedAt:   time.Unix(id, 0).UTC(),
		CompletedAt: time.Unix(id+60, 0).UTC(),
		HTMLURL:     url,
	}
}
