# Limitations

`gh-flakefinder` intentionally keeps v0.1 small and evidence-based.

## What It Can Do

- Detect workflow/job-level history where the same commit and workflow failed and later passed.
- Highlight timed-out jobs that later passed.
- Generate table, JSON, and Markdown reports with GitHub evidence links.
- Run fully offline against fixtures for demos and tests.

## What It Does Not Prove

- It does not prove a specific test case is flaky.
- It does not prove root cause.
- It does not parse logs, JUnit XML, or framework-specific output in v0.1.
- It does not find failures that were never rerun or followed by a success.

## Operational Boundaries

- Live scans depend on GitHub API availability, rate limits, and token permissions.
- Private repositories require authenticated `gh` access with permission to read Actions metadata.
- `issue --create` requires issue write permission and is never used unless explicitly requested.
- The tool never modifies workflow files and never reruns Actions jobs.
