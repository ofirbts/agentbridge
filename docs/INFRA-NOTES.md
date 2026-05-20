# AgentBridge in the AI infra world

AgentBridge is not a web data platform. It is a **light execution layer** that can sit on top of web data infra:

- plug your own `SearchProvider` that uses Bright-style proxy networks, browser APIs, or MCP-style servers
- plug your own `Crawler` that uses Puppeteer / Playwright / Browser API / proxy fetchers

## Goal

Decouple **agent logic** from **execution environment**.

## Suitable for

- RAG pipelines
- agent orchestration
- AI infra tools that need reliable web steps

## Avoid being

- yet another scraper
- yet another agent framework

Think of AgentBridge as the **execution reliability layer** for web-aware AI agents.

## Integration pattern

```
Your Agent Logic
      ↓
AgentBridge Engine (retries, tracing, deterministic mode)
      ↓
Your Providers (Bright Data / MCP / browser API)
      ↓
The Web
```

## MCP (future)

MCP servers can back `SearchProvider` and `Crawler` implementations without changing the CLI surface.
