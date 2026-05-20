# AgentBridge Use Cases

## AI research agent

**Input**

```txt
find AI infra startups
```

**Pipeline**

```txt
search → crawl → extract → normalize
```

**Command**

```bash
./agentbridge run "find AI infra startups" --mode deterministic
./agentbridge inspect <run-id> --format markdown
```

**Output**

Structured JSON with `run_id`, `step_logs`, and `result.normalized` (text, tokens, word_count) ready for a downstream LLM or RAG store.

**Why AgentBridge**

- **Inspectable** — every run persisted; `inspect` shows per-step providers, duration, retries.
- **Retryable** — search/crawl steps use retry policy; errors get `classified_errors`.
- **Deterministic option** — same task → same output for CI and demos.

Walkthrough: [examples/quickstart.md](../examples/quickstart.md)

---

## Other scenarios

| Scenario | Doc |
|----------|-----|
| First-time flow (explain → run → inspect) | [quickstart.md](../examples/quickstart.md) |
| CI / reproducible runs | [deterministic.md](../examples/deterministic.md) |
| Real HTTP fetch | [http.md](../examples/http.md) |
| Retry testing | [failure.md](../examples/failure.md) |

## What AgentBridge is not

- Not a browser automation tool
- Not a hosted SaaS or dashboard
- Not an agent framework with personas or memory

See [DECISION.md](DECISION.md).
