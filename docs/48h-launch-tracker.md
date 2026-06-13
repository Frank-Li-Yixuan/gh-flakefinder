# 48-Hour Launch Tracker

Checklist start: 2026-06-13T12:24:29+08:00

Repository: https://github.com/Frank-Li-Yixuan/gh-flakefinder

Release: https://github.com/Frank-Li-Yixuan/gh-flakefinder/releases/tag/v0.1.0

Tag: v0.1.0

Starting main commit: 4114587199a4ee6219f4a72979c80f0fc8aa372e

## Metrics

| Metric | Target | Status | Evidence | Owner | Next action |
|---|---:|---|---|---|---|
| GitHub description/topics set | 100% | done | `gh repo view` confirms target description and 10 topics | Codex | None |
| Social preview uploaded | 100% | manual | `assets/social-preview.svg` exists; upload instructions in `docs/social-preview-upload.md` | Human | Export/upload preview image in GitHub settings |
| CI warning cleanup | 100% | done | CI run https://github.com/Frank-Li-Yixuan/gh-flakefinder/actions/runs/27456198494 succeeded after workflow update | Codex | Watch CI after this checklist commit |
| Remote install retest | 2 environments | partial | Local Windows retest succeeded; post-release smoke workflow pending | Codex | Push and run `post-release-smoke` |
| Launch post prepared/sent | 2-3 channels | prepared | Drafts in `docs/48h-launch-posts.md` | Human | Manually post and capture URLs |
| Real repo scans | 3 repos | done | 3 public repos scanned; see `docs/real-scan-log.md` | Codex | None |

## Verified Repository Metadata

Description:

```text
Find GitHub Actions jobs that fail, get rerun, and mysteriously pass.
```

Topics:

```text
ci, cli, developer-tools, devops, flaky-tests, gh-extension, github-actions, golang, open-source, testing-tools
```

Verification command:

```bash
gh repo view Frank-Li-Yixuan/gh-flakefinder --json description,repositoryTopics,url
```

## CI Warning Cleanup

Current workflow state:

- `.github/workflows/ci.yml` uses `actions/checkout@v6` and `actions/setup-go@v6`.
- `actions/setup-go` has `cache: false`, avoiding cache warnings for this `go.sum`-free repository.
- `.github/workflows/release.yml` uses `actions/checkout@v6` and preserves `cli/gh-extension-precompile@v2`.
- No release tag was changed.

Latest verified CI run before this checklist work:

- https://github.com/Frank-Li-Yixuan/gh-flakefinder/actions/runs/27456198494
- conclusion: success
- head commit: 4114587199a4ee6219f4a72979c80f0fc8aa372e

## Post-Release Smoke

Workflow file: `.github/workflows/post-release-smoke.yml`

Status: pending push and workflow run.

Expected environments:

- `ubuntu-latest`
- `windows-latest`
- `macos-latest`

## Local Remote Install Retest

Status: local environment complete; GitHub-hosted smoke pending.

Commands:

```bash
gh extension remove flakefinder
gh extension install Frank-Li-Yixuan/gh-flakefinder
gh flakefinder --help
gh flakefinder scan --repo Frank-Li-Yixuan/gh-flakefinder --days 30 --limit 100 --format table
```

Result:

- environment: Windows 11
- `gh` version: 2.92.0
- extension source: `Frank-Li-Yixuan/gh-flakefinder`
- installed extension version: `v0.1.0`
- remove exit: 0
- install exit: 0
- help exit: 0
- self-scan exit: 0
- self-scan output: no likely flaky workflow or job signals matched the selected scan window

## Real Repo Scans

Completed public repositories:

- `Frank-Li-Yixuan/gh-flakefinder`
- `cli/safeexec`
- `cli/gh-webhook`

All scans used:

```bash
gh flakefinder scan --repo OWNER/REPO --days 30 --limit 100 --format table
gh flakefinder scan --repo OWNER/REPO --days 30 --limit 100 --format json
```

All completed scans exited 0 and returned `suspect_count: 0`.

## Metrics Snapshot

Snapshot file: `docs/metrics-snapshot.md`

Collected:

- stars/forks/watchers
- release asset download counts
- issues
- recent workflow runs
- traffic views/clones/referrers/paths

## Human Manual Actions

- Export `assets/social-preview.svg` to PNG at approximately 1280x640 and under 1 MB, or use an equivalent image exported from the SVG source.
- Upload the image through GitHub repository settings: Settings -> Social preview -> Edit -> Upload image.
- Capture a screenshot or note the upload date in this tracker.
- Post the prepared launch drafts manually and update the posted URL fields.
- Watch early feedback for installation, permission, false-positive, and missing-pattern reports.
