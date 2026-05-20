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

## MCP (PR #2 stub)

`pkg/mcp` provides:

- JSON-RPC types (`Request`, `Response`, `ToolCallParams`)
- `Client` interface and `StubClient`
- CLI wiring via `--mcp-endpoint`

Replace the stub with SSE/stdio transport when connecting to a live MCP server.

## Real HTTP crawler (PR #2)

`internal/crawl/http.go`:

- timeout + user-agent
- 429 / 5xx mapped to retry-friendly errors
- enabled with `--crawler http`

## MCP (future full integration)
