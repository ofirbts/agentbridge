# Failure simulation

Test how retries behave without touching production web APIs.

## Command

```bash
./agentbridge simulate-failure --seed 42 --count 3
```

## Output

```json
{
  "count": 3,
  "failures": 1,
  "kind": "random",
  "results": [
    { "attempt": 1, "kind": "random", "success": true },
    { "attempt": 2, "kind": "random", "success": true },
    { "attempt": 3, "error": "simulated failure", "kind": "random", "success": false }
  ]
}
```

## What this means

- Injects random failures with a fixed **seed** for reproducibility.
- Shows which attempt failed and the error string.
- Use before wiring AgentBridge into a pipeline to reason about retry policy behavior.
- Pair with `run` on real tasks to see `classified_errors` (e.g. `rate_limit`, `timeout`) when providers fail.
