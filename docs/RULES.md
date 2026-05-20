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

### Crawler

- `Fetch(url string) ([]byte, error)`
- must surface fetch errors to the engine retry loop

## Non-goals (current version)

- no guarantee of real web fidelity (mocks only)
- no cross-run memory beyond stored run records
- no automatic commit or deploy behavior
