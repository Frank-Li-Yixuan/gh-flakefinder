# gh-flakefinder

Find flaky GitHub Actions runs before they waste another hour.

`gh-flakefinder` is a GitHub CLI extension that scans recent GitHub Actions workflow and job history for likely flaky CI signals, such as the same commit failing and then passing after a rerun.

It is read-only by default, requires no workflow changes, and can generate table, JSON, or issue-ready Markdown reports.

## Install

From a published repository:

```bash
gh extension install Frank-Li-Yixuan/gh-flakefinder
```

For local development from this checkout:

```bash
go build -o gh-flakefinder ./cmd/gh-flakefinder
gh extension install .
```

After installation, run the extension as:

```bash
gh flakefinder --help
```

## Quickstart

Scan the current repository, inferred from the local git remote or `gh repo view`:

```bash
gh flakefinder scan
```

Scan an explicit repository:

```bash
gh flakefinder scan --repo owner/repo --days 30 --limit 300
```

Use fixture mode for a deterministic demo without GitHub network access:

```bash
gh flakefinder scan --fixture testdata/workflow_runs.json --format table
gh flakefinder scan --fixture testdata/workflow_runs.json --format json
gh flakefinder scan --fixture testdata/workflow_runs.json --format markdown
```

Create an issue-ready report without writing to GitHub:

```bash
gh flakefinder issue --repo owner/repo --days 30 --dry-run
```

Actually creating an issue requires the explicit write flag:

```bash
gh flakefinder issue --repo owner/repo --days 30 --create --label flaky-ci
```

## Commands

### `scan`

```bash
gh flakefinder scan [flags]
```

Flags:

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
gh flakefinder issue [flags]
```

`issue` reuses the scan flags and renders the Markdown report as an issue body. It defaults to dry-run output. GitHub writes only happen when `--create` is present.

## Example Output

```text
Suspected flaky GitHub Actions jobs in owner/repo, last 30 days

CONF  WORKFLOW  JOB                 SHA      SIGNAL                CATEGORY
0.90  CI        test-node-20-linux  abc1234  failure -> success    flaky-test
0.85  CI        (workflow)          abc1234  failure -> success    flaky-test
0.70  Build     integration         def9876  timed_out -> success  infra/network
0.70  Build     (workflow)          def9876  timed_out -> success  infra/network

Run with --format markdown to create a report for an issue.
```

## Detection Model

The MVP detects workflow/job-level suspects, not exact flaky test cases. It groups recent history by:

- Workflow level: `head_sha + workflow_id + event + branch`
- Job level: `head_sha + workflow_id + job_name`

It then reports failure-like conclusions followed by success:

- `failure -> success`
- `timed_out -> success`
- `cancelled -> success`, below the default confidence threshold

Each suspect includes evidence links to the failed and later successful run or job.

## Permissions

`scan` is read-only. For public repositories, normal GitHub API read access is enough. For private repositories, authenticate `gh` with a token that can read Actions metadata for that repository.

`issue --create` writes to GitHub and requires issue creation permission. Without `--create`, `issue` only prints Markdown locally.

## Limitations

- Heuristic suspects are not proof of a flaky test.
- The tool does not modify workflow files or rerun jobs.
- The MVP does not parse logs or JUnit reports.
- Runs that were never rerun may be missed.
- GitHub API rate limits and repository permissions can limit live scans.

See [docs/algorithm.md](docs/algorithm.md) and [docs/limitations.md](docs/limitations.md) for details.

## Development

```bash
go test ./...
go vet ./...
go build -o gh-flakefinder ./cmd/gh-flakefinder
./gh-flakefinder scan --fixture testdata/workflow_runs.json --format json
./gh-flakefinder scan --fixture testdata/workflow_runs.json --format markdown
```

If `gh` is installed:

```bash
gh extension install .
gh flakefinder scan --fixture testdata/workflow_runs.json --format table
```

## Release

Push a tag to build release binaries:

```bash
git tag v0.1.0
git push origin v0.1.0
```

The release workflow uses `cli/gh-extension-precompile` to publish precompiled binaries for supported platforms.

## License

MIT
