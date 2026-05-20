---
name: feature-build
description: End-to-end feature workflow — plan if needed, implement with tests, strict review. Use when user asks to build a feature, add functionality, or implement a task.
---

# Feature Build

## Decide path

Clear task (goal + files + done criteria): skip to developer.
Unclear or multi-file: invoke architect first.

## Developer phase

1. Read surrounding code
2. Implement minimal change
3. Add/update tests
4. Run `make test`

## Review gate

Invoke strict-reviewer on the full diff.
If BLOCKED: fix and re-review.
If APPROVED WITH WARNINGS: fix warnings or get explicit user acceptance.

## Done criteria

- `make test` passes
- strict-reviewer verdict is not BLOCKED
- User has not asked to hold the commit
