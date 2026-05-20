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

## MCP (`pkg/mcp`)

- JSON-RPC types (`Request`, `Response`, `ToolCallParams`)
- `Client` interface: `StubClient` and `HTTPClient`
- CLI: `--mcp-endpoint`, `--mcp-transport` (`http` | `stub`)

SSE/stdio transport is future work.

## HTTP crawler (`internal/crawl/http.go`)

- timeout + user-agent
- 429 / 5xx mapped to retry-friendly errors
- CLI: `--crawler http`

## Search (`internal/search`)

- `SearchProvider` interface
- `mock` (default, deterministic-safe)
- `http` — DuckDuckGo Instant Answer API (no API key)
- CLI: `--search mock|http`
