# Algorithm

`gh-flakefinder` reports likely flaky GitHub Actions workflow and job signals from recent history. It is a heuristic detector, not a test-case attribution engine.

## Inputs

The detector consumes normalized workflow runs and workflow jobs:

- Workflow runs from `GET /repos/{owner}/{repo}/actions/runs`
- Workflow jobs from `GET /repos/{owner}/{repo}/actions/runs/{run_id}/jobs`
- Local fixture JSON in the same normalized shape for tests and demos

## Grouping

Workflow-level grouping:

```text
head_sha + workflow_id + event + branch
```

Job-level grouping:

```text
head_sha + workflow_id + job_name
```

The job-level signal is preferred when a specific job changed from a failure-like conclusion to success.

## Signals

The detector sorts each group by observed time and looks for a failure-like conclusion followed by `success`.

Default confidence scores:

```text
failure -> success:   0.85 workflow, 0.90 job
timed_out -> success: 0.70
cancelled -> success: 0.45
same run/job attempt changing result: 0.95
```

The default `--min-confidence 0.60` filters out low-confidence cancellation signals.

## Output Contract

Each suspect includes:

- workflow name and ID
- job name when available
- full head SHA plus short SHA in human reports
- branch and event
- signal and category
- confidence score
- failed and passed evidence URLs
- suggested next action

JSON output is deterministic and intentionally does not include generation timestamps.
