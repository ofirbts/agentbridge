# Deterministic mode

Same task string → same `run_id` and stable mock payloads. Use in CI and demos.

## Command

```bash
./agentbridge run "find AI infra startups" --mode deterministic
```

## Output (excerpt)

```json
{
  "status": "success",
  "run_id": "run_det_23434e37",
  "mode": "deterministic",
  "search": "mock",
  "crawler": "mock",
  "retries": 0,
  "result": {
    "query": "find AI infra startups",
    "url": "https://example.com",
    "normalized": {
      "text": "AgentBridge mock page",
      "tokens": ["AgentBridge", "mock", "page"],
      "word_count": 3
    }
  },
  "step_logs": [
    { "step": "search", "status": "ok", "provider": "mock", "duration_ms": 0, "retries": 0 },
    { "step": "crawl", "status": "ok", "provider": "mock", "duration_ms": 0, "retries": 0 }
  ]
}
```

Run again with the same task — `run_id` and `result` stay the same.

## What this means

- **deterministic** forces mock search (even if you pass `--search http`).
- Outputs are reproducible — important for tests and portfolio demos.
- `step_logs` still records per-step providers and timing (here near-zero with mocks).
