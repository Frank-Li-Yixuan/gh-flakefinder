---
name: Bug report
about: Report a problem with gh-flakefinder
title: "Bug: "
labels: bug
assignees: ""
---

## What happened?

Describe the behavior you saw.

## What did you expect?

Describe what you expected instead.

## Command

```bash
gh flakefinder scan --repo owner/repo --days 30
```

## Output

Paste the relevant output. Remove private repository names, tokens, and internal URLs if needed.

## Environment

- OS:
- `gh --version`:
- `gh flakefinder --help` works: yes/no
- Repository type: public/private

## Notes

If this involves a live repository, include whether fixture mode works:

```bash
gh flakefinder scan --fixture testdata/workflow_runs.json --format table
```
