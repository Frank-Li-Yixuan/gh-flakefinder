# Real Scan Log

Scope: low-volume scans of public repositories only. Commands use `--limit 100` and do not create issues.

Checklist start: 2026-06-13T12:24:29+08:00

## Summary

| Repository | Days | Limit | Table exit | JSON exit | Suspect count | Actionable? | Notes |
|---|---:|---:|---:|---:|---:|---|---|
| Frank-Li-Yixuan/gh-flakefinder | 30 | 100 | 0 | 0 | 0 | no suspects | Required self-scan; command completed cleanly |
| cli/safeexec | 30 | 100 | 0 | 0 | 0 | no suspects | Small public repo with recent Actions history |
| cli/gh-webhook | 30 | 100 | 0 | 0 | 0 | no suspects | Public repo with Actions history; no suspect signals in the selected window |

## Commands

For each repository:

```bash
gh flakefinder scan --repo OWNER/REPO --days 30 --limit 100 --format table
gh flakefinder scan --repo OWNER/REPO --days 30 --limit 100 --format json
```

## Detailed Notes

### Frank-Li-Yixuan/gh-flakefinder

Commands:

```bash
gh flakefinder scan --repo Frank-Li-Yixuan/gh-flakefinder --days 30 --limit 100 --format table
gh flakefinder scan --repo Frank-Li-Yixuan/gh-flakefinder --days 30 --limit 100 --format json
```

Results:

- table exit: 0
- JSON exit: 0
- suspect count: 0
- output summary: no likely flaky workflow or job signals matched the selected scan window
- API/auth/rate-limit issues: none observed
- obvious false positives: none
- obvious missed pattern: none visible from the scan output

### cli/safeexec

Commands:

```bash
gh flakefinder scan --repo cli/safeexec --days 30 --limit 100 --format table
gh flakefinder scan --repo cli/safeexec --days 30 --limit 100 --format json
```

Results:

- table exit: 0
- JSON exit: 0
- suspect count: 0
- output summary: no likely flaky workflow or job signals matched the selected scan window
- API/auth/rate-limit issues: none observed
- obvious false positives: none
- obvious missed pattern: none visible from the scan output

### cli/gh-webhook

Commands:

```bash
gh flakefinder scan --repo cli/gh-webhook --days 30 --limit 100 --format table
gh flakefinder scan --repo cli/gh-webhook --days 30 --limit 100 --format json
```

Results:

- table exit: 0
- JSON exit: 0
- suspect count: 0
- output summary: no likely flaky workflow or job signals matched the selected scan window
- API/auth/rate-limit issues: none observed
- obvious false positives: none
- obvious missed pattern: none visible from the scan output

### Excluded Candidate

`cli/browser` was tried as a candidate scan target but the live scan did not complete before a 300-second command timeout. It is not counted toward the 3 completed public repository scans.
