# Social Posts

## Short Launch Post

I released gh-flakefinder, a GitHub CLI extension that finds GitHub Actions jobs that fail, get rerun, and mysteriously pass.

Install:

```bash
gh extension install Frank-Li-Yixuan/gh-flakefinder
```

Run:

```bash
gh flakefinder scan --repo owner/repo --days 30
```

No server, no workflow changes, read-only by default.

## Maintainer-Focused Post

CI failed, someone clicked rerun, and now the job is green. Useful, but the signal disappeared.

`gh-flakefinder` scans recent GitHub Actions history and reports likely workflow/job-level flaky suspects with evidence links.

It does not claim perfect root cause. It gives maintainers a short list worth investigating.

```bash
gh extension install Frank-Li-Yixuan/gh-flakefinder
gh flakefinder scan --repo owner/repo --days 30 --format markdown
```

## Longer Launch Note

I shipped `gh-flakefinder` v0.1.0.

It is a small GitHub CLI extension for a common CI maintenance problem: a GitHub Actions job fails, gets rerun, and passes. That can be a flaky test, runner instability, network failure, dependency issue, or manual cancellation noise. Either way, the pattern is easy to lose once the PR turns green.

`gh-flakefinder` scans workflow/job history and reports likely suspects with evidence links. It is read-only by default, has deterministic fixture demos, and can output table, JSON, or issue-ready Markdown.

Repo: https://github.com/Frank-Li-Yixuan/gh-flakefinder
