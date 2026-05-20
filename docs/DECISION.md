# Why AgentBridge is an execution layer, not a framework

AgentBridge deliberately avoids becoming another agent framework.

## Decision

Build a **thin execution reliability layer** that:

- orchestrates search → crawl → extract → normalize
- handles retries, backoff, and failure simulation
- records structured run logs for inspection
- exposes pluggable provider interfaces

Do **not** build:

- agent personas or prompt orchestration
- a web UI or dashboard (future roadmap only)
- a full scraping engine (mock providers first)

## Rationale

Production agent failures usually come from the web execution environment:

- rate limits and timeouts
- inconsistent HTML and blocked requests
- non-deterministic tool behavior

Frameworks optimize for agent logic. AgentBridge optimizes for **durable web execution**.

## Trade-offs

| Choice | Benefit | Cost |
|--------|---------|------|
| Mock providers first | fast iteration, testable | no real web data yet |
| CLI-first | simple ops, easy CI | no GUI |
| File-backed run store | inspect works across CLI invocations | local filesystem dependency |

## Shipped capabilities

| Area | CLI | Notes |
|------|-----|-------|
| HTTP crawl | `--crawler http` | timeouts, retry-friendly errors |
| MCP | `--mcp-endpoint`, `--mcp-transport` | stub + HTTP JSON-RPC (`tools/list`, `tools/call`) |
| Search | `--search mock` (default) | HTTP provider in Phase 2 |

## Future (not current scope)

- MCP SSE/stdio transport
- real browser automation backend
- plugin registry for third-party providers
- OpenTelemetry export
