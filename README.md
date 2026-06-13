# gh-flakefinder

Find GitHub Actions jobs that fail, get rerun, and mysteriously pass.

`gh-flakefinder` is a GitHub CLI extension for maintainers who want a quick, evidence-backed view of likely flaky workflow or job runs. It scans recent GitHub Actions history and reports cases where the same commit, workflow, or job failed and later passed.

```bash
gh extension install Frank-Li-Yixuan/gh-flakefinder
```

```bash
gh flakefinder scan --repo owner/repo --days 30
```

> Safety note: `scan` is read-only by default. It does not modify workflows, rerun jobs, create issues, post comments, or send telemetry. The only write path is `gh flakefinder issue --create`, which must be passed explicitly.

## Before / After

Before:

- CI fails on a PR that did not touch the failing area.
- Someone clicks rerun.
- The job turns green.
- The failure disappears from the review flow, and nobody records the pattern.

After:

- Run `gh flakefinder scan`.
- Get a compact list of workflow/job suspects.
- Copy Markdown evidence into an issue or incident note.
- Decide whether to investigate tests, runner capacity, network calls, or dependencies.

## 30-Second Demo

Try the deterministic fixture demo:

```bash
git clone https://github.com/Frank-Li-Yixuan/gh-flakefinder.git
cd gh-flakefinder
gh extension install .
gh flakefinder scan --fixture testdata/workflow_runs.json --format table
```

Example output:

```text
Suspected flaky GitHub Actions jobs in owner/repo, last 30 days

CONF  WORKFLOW  JOB                 SHA      SIGNAL                CATEGORY
0.90  CI        test-node-20-linux  abc1234  failure -> success    flaky-test
0.85  CI        (workflow)          abc1234  failure -> success    flaky-test
0.70  Build     integration         def9876  timed_out -> success  infra/network
0.70  Build     (workflow)          def9876  timed_out -> success  infra/network

Run with --format markdown to create a report for an issue.
```

For a real repository:

```bash
gh flakefinder scan --repo owner/repo --days 30 --limit 300
```

## What It Detects

`gh-flakefinder` reports likely workflow/job-level suspects when recent Actions history shows:

- same run/job rerun where a failed job later passed
- same `head_sha + workflow + job_name` with `failure -> success`
- same `head_sha + workflow` with `failure -> success`
- `timed_out -> success` patterns that may point to runner, network, or dependency instability
- low-confidence `cancelled -> success` patterns only when the confidence threshold allows them

Every suspect includes direct GitHub evidence links.

## What It Does Not Detect

This is a heuristic detector, not a flaky-test oracle. It does not:

- prove a specific test case is flaky
- parse logs, JUnit XML, pytest output, Jest output, or Go test output in v0.1
- detect failures that were never rerun or never followed by success
- modify workflow files
- rerun jobs
- create issues unless `issue --create` is passed explicitly

## Commands

### `scan`

```bash
gh flakefinder scan [flags]
```

Common flags:

```text
--repo owner/repo            Repository to scan; default is inferred from the current directory.
--days 30                    Number of recent days to scan.
--limit 300                  Maximum workflow runs to fetch.
--workflow "CI"              Only scan one workflow name or workflow ID.
--format table|json|markdown Output format. Default: table.
--min-confidence 0.60        Minimum confidence threshold.
--fixture path.json          Load local fixture data instead of calling GitHub.
--verbose                    Print API diagnostics.
```

### `issue`

```bash
gh flakefinder issue --repo owner/repo --days 30 --dry-run
```

`issue` reuses the scan flags and renders the Markdown report as an issue body. It defaults to dry-run output. GitHub writes only happen when `--create` is present:

```bash
gh flakefinder issue --repo owner/repo --days 30 --create --label flaky-ci
```

## Output Formats

Table is best for terminal scanning:

```bash
gh flakefinder scan --repo owner/repo --format table
```

JSON is best for scripts:

```bash
gh flakefinder scan --repo owner/repo --format json
```

Markdown is best for issues, PR comments, or incident notes:

```bash
gh flakefinder scan --repo owner/repo --format markdown > flaky-ci.md
```

See [docs/demo-output.md](docs/demo-output.md) for a fuller example.

## Permissions

For public repositories, normal GitHub API read access is enough.

For private repositories, authenticate `gh` with a token that can read Actions metadata for that repository.

`issue --create` requires issue creation permission. Without `--create`, `issue` only prints Markdown locally.

## Documentation

- [30-second demo](docs/demo.md)
- [Demo output](docs/demo-output.md)
- [Maintainer FAQ](docs/maintainer-faq.md)
- [Social posts](docs/social-posts.md)
- [Repository settings](docs/repo-settings.md)
- [Algorithm](docs/algorithm.md)
- [Limitations](docs/limitations.md)
- [Release notes](docs/release-notes-v0.1.0.md)

## Development

Run the release gate:

```bash
scripts/verify.sh
```

On Windows:

```powershell
.\scripts\verify.ps1
```

Core commands:

```bash
go test ./...
go vet ./...
go build -o gh-flakefinder ./cmd/gh-flakefinder
./gh-flakefinder scan --fixture testdata/workflow_runs.json --format json
./gh-flakefinder scan --fixture testdata/workflow_runs.json --format markdown
```

## Release

`v0.1.0` is published at:

https://github.com/Frank-Li-Yixuan/gh-flakefinder/releases/tag/v0.1.0

Future releases are built by pushing a tag:

```bash
git tag v0.1.1
git push origin v0.1.1
```

The release workflow uses `cli/gh-extension-precompile` to publish precompiled binaries for supported platforms.

## License

MIT
