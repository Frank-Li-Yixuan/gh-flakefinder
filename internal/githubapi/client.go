package githubapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

type FetchOptions struct {
	Days     int
	Limit    int
	Workflow string
	Verbose  bool
}

type Fetcher interface {
	Fetch(ctx context.Context, repo string, opts FetchOptions) (DataSet, error)
}

type GHClient struct {
	Command string
	Stderr  io.Writer
	Now     func() time.Time
}

func (c GHClient) Fetch(ctx context.Context, repo string, opts FetchOptions) (DataSet, error) {
	if err := validateRepo(repo); err != nil {
		return DataSet{}, err
	}
	if opts.Limit <= 0 {
		opts.Limit = 300
	}
	if opts.Days <= 0 {
		opts.Days = 30
	}
	runs, err := c.fetchRuns(ctx, repo, opts)
	if err != nil {
		return DataSet{}, err
	}
	data := DataSet{Repository: repo, Runs: runs}
	for _, run := range runs {
		if opts.Workflow != "" && !workflowMatches(run, opts.Workflow) {
			continue
		}
		jobs, err := c.fetchJobs(ctx, repo, run)
		if err != nil {
			return DataSet{}, err
		}
		data.Jobs = append(data.Jobs, jobs...)
	}
	normalizeDataSet(&data)
	return data, nil
}

func (c GHClient) fetchRuns(ctx context.Context, repo string, opts FetchOptions) ([]WorkflowRun, error) {
	perPage := 100
	var runs []WorkflowRun
	now := time.Now
	if c.Now != nil {
		now = c.Now
	}
	since := now().UTC().AddDate(0, 0, -opts.Days).Format("2006-01-02")
	for page := 1; len(runs) < opts.Limit; page++ {
		args := []string{
			"--method", "GET",
			fmt.Sprintf("/repos/%s/actions/runs", repo),
			"-f", fmt.Sprintf("per_page=%d", perPage),
			"-f", fmt.Sprintf("page=%d", page),
			"-f", fmt.Sprintf("created=>=%s", since),
		}
		if opts.Workflow != "" {
			debug(c.Stderr, opts.Verbose, "fetch workflow runs page %d for %s", page, repo)
		}
		out, err := c.ghAPI(ctx, args...)
		if err != nil {
			return nil, err
		}
		var response struct {
			WorkflowRuns []WorkflowRun `json:"workflow_runs"`
		}
		if err := json.Unmarshal(out, &response); err != nil {
			return nil, fmt.Errorf("parse workflow runs response: %w", err)
		}
		if len(response.WorkflowRuns) == 0 {
			break
		}
		for _, run := range response.WorkflowRuns {
			if opts.Workflow == "" || workflowMatches(run, opts.Workflow) {
				runs = append(runs, run)
			}
			if len(runs) >= opts.Limit {
				break
			}
		}
		if len(response.WorkflowRuns) < perPage {
			break
		}
	}
	return runs, nil
}

func (c GHClient) fetchJobs(ctx context.Context, repo string, run WorkflowRun) ([]WorkflowJob, error) {
	attempts := run.Attempt
	if attempts < 1 {
		attempts = 1
	}
	var jobs []WorkflowJob
	for attempt := 1; attempt <= attempts; attempt++ {
		endpoint := fmt.Sprintf("/repos/%s/actions/runs/%d/attempts/%d/jobs", repo, run.ID, attempt)
		if attempts == 1 {
			endpoint = fmt.Sprintf("/repos/%s/actions/runs/%d/jobs", repo, run.ID)
		}
		for page := 1; ; page++ {
			out, err := c.ghAPI(ctx, "--method", "GET", endpoint, "-f", "per_page=100", "-f", fmt.Sprintf("page=%d", page))
			if err != nil {
				return nil, err
			}
			var response struct {
				Jobs []WorkflowJob `json:"jobs"`
			}
			if err := json.Unmarshal(out, &response); err != nil {
				return nil, fmt.Errorf("parse workflow jobs response for run %d: %w", run.ID, err)
			}
			for _, job := range response.Jobs {
				if job.RunID == 0 {
					job.RunID = run.ID
				}
				if job.Attempt == 0 {
					job.Attempt = attempt
				}
				jobs = append(jobs, job)
			}
			if len(response.Jobs) < 100 {
				break
			}
		}
	}
	return jobs, nil
}

func (c GHClient) ghAPI(ctx context.Context, args ...string) ([]byte, error) {
	command := c.Command
	if command == "" {
		command = "gh"
	}
	cmd := exec.CommandContext(ctx, command, append([]string{"api"}, args...)...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		detail := strings.TrimSpace(stderr.String())
		if detail == "" {
			detail = err.Error()
		}
		return nil, fmt.Errorf("gh api failed: %s", detail)
	}
	return out, nil
}

func validateRepo(repo string) error {
	parts := strings.Split(repo, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("invalid repo %q; expected owner/repo", repo)
	}
	return nil
}

func workflowMatches(run WorkflowRun, filter string) bool {
	if filter == "" {
		return true
	}
	if strings.EqualFold(run.Name, filter) {
		return true
	}
	return fmt.Sprintf("%d", run.WorkflowID) == filter
}

func debug(w io.Writer, verbose bool, format string, args ...any) {
	if !verbose || w == nil {
		return
	}
	fmt.Fprintf(w, format+"\n", args...)
}
