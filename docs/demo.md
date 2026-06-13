# 30-Second Demo

This demo uses the checked-in fixture, so it works without GitHub network access and always produces the same suspects.

## Install Locally

```bash
git clone https://github.com/Frank-Li-Yixuan/gh-flakefinder.git
cd gh-flakefinder
gh extension install .
```

## Run The Fixture Scan

```bash
gh flakefinder scan --fixture testdata/workflow_runs.json --format table
```

You should see four suspects:

- `CI / test-node-20-linux`: `failure -> success`
- `CI`: workflow-level `failure -> success`
- `Build / integration`: `timed_out -> success`
- `Build`: workflow-level `timed_out -> success`

## Generate Issue-Ready Markdown

```bash
gh flakefinder scan --fixture testdata/workflow_runs.json --format markdown > flaky-ci.md
```

The Markdown report includes direct evidence links for the failed and later successful runs/jobs.

## Scan A Real Repository

```bash
gh flakefinder scan --repo owner/repo --days 30 --limit 300
```

If there are no likely suspects, the tool exits successfully and says no matching signals were found.
