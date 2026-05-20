---
name: review-gate
description: Mandatory strict review before merge or task completion. Use after implementation, before PR, or when user says review or בודק.
---

# Review Gate

1. Ensure working tree has the changes to review (git status, git diff)
2. Invoke strict-reviewer subagent
3. If BLOCKED: developer fixes each critical item
4. Re-run strict-reviewer until not BLOCKED
5. Summarize verdict for the user in Hebrew if they wrote in Hebrew

Do not mark a task complete while strict-reviewer is BLOCKED unless user explicitly overrides.
