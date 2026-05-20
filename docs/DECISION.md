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

## PR #2 additions

- real HTTP crawler (`--crawler http`)
- MCP stub client in `pkg/mcp`
- richer `explain` / `inspect` / `simulate-failure` CLI output
- markdown inspect format

## Future (not current scope)

- MCP integration
- real browser automation backend
- plugin system for providers
- OpenTelemetry export
