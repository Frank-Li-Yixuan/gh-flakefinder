package githubapi

import "time"

type DataSet struct {
	Repository string        `json:"repository"`
	Runs       []WorkflowRun `json:"workflow_runs"`
	Jobs       []WorkflowJob `json:"workflow_jobs"`
}

type WorkflowRun struct {
	ID         int64     `json:"id"`
	Attempt    int       `json:"run_attempt"`
	Name       string    `json:"name"`
	WorkflowID int64     `json:"workflow_id"`
	Event      string    `json:"event"`
	HeadSHA    string    `json:"head_sha"`
	HeadBranch string    `json:"head_branch"`
	Conclusion string    `json:"conclusion"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	HTMLURL    string    `json:"html_url"`
}

type WorkflowJob struct {
	ID          int64     `json:"id"`
	RunID       int64     `json:"run_id"`
	Attempt     int       `json:"run_attempt"`
	Name        string    `json:"name"`
	Conclusion  string    `json:"conclusion"`
	Status      string    `json:"status"`
	StartedAt   time.Time `json:"started_at"`
	CompletedAt time.Time `json:"completed_at"`
	HTMLURL     string    `json:"html_url"`
}

func (d DataSet) WithLimit(limit int) DataSet {
	if limit <= 0 || len(d.Runs) <= limit {
		return d
	}
	limited := DataSet{
		Repository: d.Repository,
		Runs:       append([]WorkflowRun(nil), d.Runs[:limit]...),
	}
	allowedRuns := make(map[int64]bool, len(limited.Runs))
	for _, run := range limited.Runs {
		allowedRuns[run.ID] = true
	}
	for _, job := range d.Jobs {
		if allowedRuns[job.RunID] {
			limited.Jobs = append(limited.Jobs, job)
		}
	}
	return limited
}
