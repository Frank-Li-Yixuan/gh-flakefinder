# 48-Hour Launch Posts

Core message:

> A GitHub Actions job fails. Someone reruns it. It passes. Everyone moves on, and the flaky signal disappears.

## X / Twitter

```text
A GitHub Actions job fails. Someone reruns it. It passes. Everyone moves on, and the flaky signal disappears.

I released gh-flakefinder, a GitHub CLI extension that scans recent Actions history for likely flaky workflow/job runs.

Read-only by default. No server. No telemetry. No issue creation unless you explicitly use issue --create.

Install:
gh extension install Frank-Li-Yixuan/gh-flakefinder

Scan:
gh flakefinder scan --repo owner/repo --days 30

I am looking for feedback on false positives, missed rerun patterns, and what output format maintainers want next.
```

## LinkedIn

```text
A GitHub Actions job fails. Someone reruns it. It passes. Everyone moves on, and the flaky signal disappears.

I released gh-flakefinder v0.1.0, a small GitHub CLI extension for maintainers who want a quick read on likely flaky CI history.

It scans recent GitHub Actions workflow/job runs and reports cases where the same commit, workflow, or job failed and later passed. The goal is not to prove root cause. The goal is to surface evidence worth investigating before it gets buried by a green rerun.

Safety notes:
- read-only by default
- no server
- no telemetry
- no workflow changes
- issue creation only happens with an explicit issue --create command

Install:
gh extension install Frank-Li-Yixuan/gh-flakefinder

Scan:
gh flakefinder scan --repo owner/repo --days 30 --format markdown

Feedback wanted: false positives, missed flaky patterns, whether the report is useful in an issue or incident note, and which CI metadata would make this more actionable.
```

## Hacker News Show HN

```text
Title: Show HN: gh-flakefinder - find GitHub Actions jobs that fail, rerun, and pass

A GitHub Actions job fails. Someone reruns it. It passes. Everyone moves on, and the flaky signal disappears.

I built gh-flakefinder, a GitHub CLI extension that scans recent GitHub Actions history and reports likely flaky workflow/job runs. It looks for patterns like failure -> success or timed_out -> success for the same commit/workflow/job.

It is intentionally small and read-only by default:
- no server
- no telemetry
- no workflow modification
- no issue creation unless issue --create is passed explicitly

Install:
gh extension install Frank-Li-Yixuan/gh-flakefinder

Scan:
gh flakefinder scan --repo owner/repo --days 30

I would appreciate feedback from maintainers with noisy CI: which detections are useful, which are obvious false positives, and what evidence should be included in the Markdown report.
```

## Chinese Developer Communities

Use for V2EX, Juejin, or a private developer group where posting is appropriate.

```text
A GitHub Actions job fails. Someone reruns it. It passes. Everyone moves on, and the flaky signal disappears.

I released gh-flakefinder, a GitHub CLI extension for finding likely flaky GitHub Actions workflow/job runs from repository history.

It is read-only by default: no server, no telemetry, no workflow changes, and no issue creation unless issue --create is passed explicitly.

Install:
gh extension install Frank-Li-Yixuan/gh-flakefinder

Scan:
gh flakefinder scan --repo owner/repo --days 30

Feedback wanted: false positives, missed flaky patterns, and whether the Markdown report is useful for maintainers.

Repo:
https://github.com/Frank-Li-Yixuan/gh-flakefinder
```

## Posting Checklist

| Channel | Posted? | URL | Early replies | Useful feedback | Follow-up action |
|---|---|---|---|---|---|
| X / Twitter | no | pending | pending | pending | Post manually |
| LinkedIn | no | pending | pending | pending | Post manually |
| Hacker News Show HN | no | pending | pending | pending | Post manually |
| Chinese developer community | no | pending | pending | pending | Post only where appropriate |
