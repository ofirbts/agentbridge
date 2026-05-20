# AgentBridge

**A production-grade CLI for reliable AI agent web workflows in Go**

> AI agents are easy to build. Making them reliable in the real web is not.

AgentBridge is a thin execution layer that makes web agents predictable, observable, and production-aware.
It is not an agent framework, but a **reliability + execution layer** between your agent logic and the web.

## Why this exists

The web is:

- unstable
- rate-limited
- blocked
- inconsistent
- non-deterministic

Most agent frameworks ignore this reality.

AgentBridge sits between your agent and the web, handling:

- retries + backoff
- observability + structured logs
- deterministic execution mode
- failure simulation
- mockable providers

**Core idea:** Agents should not care about the web. The execution layer should.

## Features

- CLI-first architecture (Go + Cobra)
- Deterministic execution mode (`--mode deterministic`)
- Retry + fallback system
- Structured logs + tracing
- Pluggable providers (Search / Crawl / Extract)
- Failure simulation mode (`agentbridge simulate-failure`)
- Run-based observability (`agentbridge inspect <run-id>`)

## Architecture

```
CLI (Cobra)
└→ Workflow Engine (execution planner)
   ├→ Search Layer (mock/API)
   ├→ Crawl Layer (fetch)
   ├→ Extract + Normalize
   └→ Output Layer (JSON / Markdown / Logs)
      └→ Observability Layer (retries, tracing)
```

## Installation

```bash
git clone https://github.com/ofirbts/agentbridge
cd agentbridge
make build
./agentbridge --help
```

## Usage

### Run a task

```bash
agentbridge run --mode normal "find top AI infra startups"
```

Example output:

```json
{
  "status": "success",
  "steps": ["search", "crawl", "extract", "normalize"],
  "retries": 0,
  "duration_ms": 3200,
  "run_id": "run_abc12345",
  "result": {}
}
```

### Explain execution plan

```bash
agentbridge explain --task "find top AI infra startups"
```

### Inspect execution

```bash
agentbridge inspect run_abc12345
```

### Simulate failure conditions

```bash
agentbridge simulate-failure --seed 42 --count 3
```

### HTTP crawler + MCP stub

```bash
agentbridge run "find top AI infra startups" --crawler http
agentbridge run "mcp backed task" --mcp-endpoint http://localhost:8080/mcp
agentbridge explain --task "plan with http" --crawler http
agentbridge inspect run_abc12345 --format markdown
```

## Design Principles

1. **Execution over abstraction** — predictable execution over abstract frameworks
2. **Observability first** — every action is traceable
3. **Failure is default** — the web is unreliable; assume failure
4. **Minimal surface area** — small interfaces, replaceable components

## Use Cases

- AI agents with web browsing capability
- RAG pipelines with external context
- Data extraction workflows
- AI orchestration systems

## Docs

- [DECISION.md](docs/DECISION.md) — why execution layer, not framework
- [RULES.md](docs/RULES.md) — deterministic mode guarantees
- [INFRA-NOTES.md](docs/INFRA-NOTES.md) — AI infra / MCP fit

## Development

```bash
make build
make test
make e2e
make dev
```

## Contributing

PRs welcome for:

- new providers (web-search APIs, browser-automation backends)
- metrics exporters
- better retries / observability
- docs

## Cursor Workflow

This repo includes `.cursor/` agents, rules, and skills from `project-starter`:

1. `architect` — plan large changes
2. `developer` — implement + tests
3. `strict-reviewer` — review before commit
4. `debugger` — root-cause failures

See `project.config.json` for the workflow map.
