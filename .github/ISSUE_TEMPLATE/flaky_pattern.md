---
name: Flaky pattern report
about: Share a GitHub Actions pattern that gh-flakefinder should or should not flag
title: "Pattern: "
labels: flaky-pattern
assignees: ""
---

## Pattern

Describe the workflow/job history pattern.

## Expected classification

Should this be reported as a likely flaky suspect?

- [ ] Yes
- [ ] No
- [ ] Unsure

## Evidence shape

Do not paste private URLs unless they are safe to share. A redacted shape is enough:

```text
workflow: CI
job: test
head_sha: abc1234
attempt 1: failure
attempt 2: success
```

## Why it matters

Explain how this pattern affects maintainers.

## Current gh-flakefinder output

Paste the relevant table, JSON, or Markdown output if available.
