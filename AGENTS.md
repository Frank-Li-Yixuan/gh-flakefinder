# AGENTS.md

## Project

This repository is a public GitHub CLI extension named gh-flakefinder. It detects likely flaky GitHub Actions workflow/job runs from repository history and generates terminal, JSON, and Markdown reports.

## Product Constraints

- Read-only by default.
- No server, database, telemetry, paid API, or LLM dependency.
- Do not write to GitHub unless the user passes an explicit write flag such as `--create`.
- Prefer deterministic fixture tests for detector and renderer logic.
- Keep README commands aligned with the implementation.

## Build And Test

Run before declaring work complete:

```bash
go test ./...
go vet ./...
go build -o gh-flakefinder ./cmd/gh-flakefinder
./gh-flakefinder scan --fixture testdata/workflow_runs.json --format json
./gh-flakefinder scan --fixture testdata/workflow_runs.json --format markdown
```

If `gh` is installed, also run:

```bash
gh extension install .
gh flakefinder scan --fixture testdata/workflow_runs.json
```

## Code Style

- Prefer small internal packages.
- Core detection logic must not depend on live GitHub API calls.
- GitHub API access should be behind an interface so tests can use fixtures.
- Return useful errors with actionable remediation.
- Keep JSON output stable and documented.
