# AgentBridge API Rules

## Deterministic mode

When `--mode deterministic`:

- same task string → same `run_id` (hash-based)
- mock search results sorted by stable URL order
- normalized output uses stable field set
- failure simulation accepts explicit `--seed`

## CLI guarantees

- `run` always prints JSON to stdout (even on partial failure)
- `inspect` reads from the run store (default: `.agentbridge/runs`)
- `explain` never executes providers; plan only

## Provider contracts

### SearchProvider

- `Search(query string) ([]Result, error)`
- must be mockable for tests
- `mock` — stable fixtures (always used in `--mode deterministic`)
- `http` — DuckDuckGo Instant Answer API (no API key)

### Crawler

- `Fetch(url string) ([]byte, error)`
- must surface fetch errors to the engine retry loop
- `mock` — local HTML fixture
- `http` — real HTTP fetch with timeout and retry-friendly errors

### MCP

- `--mcp-endpoint` optional
- `--mcp-transport` `http` (JSON-RPC) or `stub` (offline)
- default transport for `http(s)` endpoints is `http`

## Non-goals (current version)

- no guarantee of real search fidelity (mock search by default)
- no cross-run memory beyond stored run records
- no automatic commit or deploy behavior
