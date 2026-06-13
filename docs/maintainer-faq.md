# Maintainer FAQ

## Is this proof that a test is flaky?

No. `gh-flakefinder` reports evidence-backed suspects. A `failure -> success` pattern is a useful clue, not a root-cause verdict.

## Does it modify my repository?

`scan` is read-only. It does not edit workflows, rerun jobs, create issues, post comments, or collect telemetry.

The only write path is `gh flakefinder issue --create`, and that flag must be passed explicitly.

## What permissions does it need?

For public repositories, normal GitHub API read access is enough.

For private repositories, your `gh` authentication must have permission to read Actions metadata for the repository.

Creating an issue with `--create` requires issue write permission.

## Why are some results workflow-level and some job-level?

GitHub Actions exposes both workflow runs and jobs. If a specific job changes from a failure-like conclusion to success, `gh-flakefinder` reports a job-level suspect. If the workflow changes result but no matching job-level signal is available, it can report a workflow-level suspect.

## Why did it report no suspects?

Common reasons:

- no recent rerun-like `failure -> success` patterns
- scan window too small
- `--limit` too low
- the failure never had a later success
- repository permissions prevented reading Actions history

Try:

```bash
gh flakefinder scan --repo owner/repo --days 90 --limit 1000
```

## Why does it include `timed_out -> success`?

Timeouts that later pass can point to runner capacity, external services, network calls, or dependency downloads. They are reported with lower confidence than same-job `failure -> success` signals.

## Does v0.1 parse logs?

No. v0.1 uses GitHub Actions run and job metadata only.
