# AgentBridge Use Cases

Short scenarios showing where the execution layer fits. Each maps to commands in `examples/`.

## 1. Reliable RAG context fetch

**Problem:** Your agent needs web context but retries and failures are ad-hoc.

**Flow:** search → crawl → extract → normalize

```bash
./agentbridge run "summarize latest AI infra tooling" --search http --crawler http
./agentbridge inspect <run-id> --format markdown
```

**Example script:** `examples/ai_infra_search.sh`

## 2. Deterministic CI smoke test

**Problem:** Pipelines need stable outputs for regression checks.

```bash
./agentbridge run "ci smoke task" --mode deterministic
```

Same task string → same `run_id` and mock payloads. See `docs/RULES.md`.

**Example script:** `examples/rag_task.sh`

## 3. Failure injection before production

**Problem:** You want to validate retry behavior without hitting production APIs.

```bash
./agentbridge simulate-failure --seed 42 --count 5
./agentbridge run "task" --mode normal
```

**Example script:** `examples/simulate_then_fix.sh`

## 4. Plan without executing

**Problem:** Operators want to see providers and observability fields before a live run.

```bash
./agentbridge explain --task "multi-step web task" --search http --crawler http
```

Returns steps, providers, and `observability` metadata (inspect fields, error classes).

## 5. MCP-augmented run (offline stub)

**Problem:** Prototype tool-backed steps before wiring a live MCP server.

```bash
./agentbridge run "task" --mcp-endpoint http://localhost:8080/mcp --mcp-transport stub
```

**Example script:** `examples/search_http.sh` (search only; combine flags as needed)

## What AgentBridge is not

- Not a browser automation tool
- Not a hosted SaaS or dashboard
- Not an agent framework with personas or memory

See [DECISION.md](DECISION.md).
