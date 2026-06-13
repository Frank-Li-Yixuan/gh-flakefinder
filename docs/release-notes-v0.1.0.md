# gh-flakefinder v0.1.0

First public release.

`gh-flakefinder` is a GitHub CLI extension that finds GitHub Actions jobs that fail, get rerun, and mysteriously pass.

## Features

- Scan recent GitHub Actions workflow/job history for likely flaky CI suspects.
- Detect same-run/job rerun patterns where attempt 1 fails and a later attempt passes.
- Detect same SHA/workflow/job `failure -> success` and `timed_out -> success` patterns.
- Output table, JSON, and issue-ready Markdown reports.
- Use fixture mode for deterministic demos and tests.
- Render issue bodies in dry-run mode by default.
- Stay read-only by default.

## Install

```bash
gh extension install Frank-Li-Yixuan/gh-flakefinder
```

## Try It

```bash
gh flakefinder scan --repo owner/repo --days 30
```

## Safety

`scan` does not modify repositories, workflows, issues, comments, or Actions runs.

Issue creation requires:

```bash
gh flakefinder issue --repo owner/repo --create
```

## Limitations

- Reports heuristic suspects, not proof of test-case flakiness.
- Does not parse logs or test reports in v0.1.
- Can miss failures that were never rerun or never followed by success.
- Live scans depend on GitHub API permissions and rate limits.
