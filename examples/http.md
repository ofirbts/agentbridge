# HTTP crawl

Fetches a real page over the network. Search stays mock unless you add `--search http`.

## Command

```bash
./agentbridge run "demo http task" --crawler http
```

## Output (excerpt)

```json
{
  "status": "success",
  "run_id": "run_5a205ced",
  "crawler": "http",
  "search": "mock",
  "duration_ms": 123,
  "result": {
    "query": "demo http task",
    "url": "https://example.com",
    "normalized": {
      "text": "Example Domain ... This domain is for use in documentation examples ...",
      "word_count": 23
    }
  },
  "step_logs": [
    { "step": "crawl", "status": "ok", "provider": "http", "duration_ms": 123, "retries": 0 }
  ]
}
```

## What this means

- **search** (mock) picks `https://example.com` as the first hit.
- **crawl** (http) downloads real HTML — `duration_ms` reflects network + parse.
- **extract/normalize** turn HTML into structured text for downstream agents.
- Failures (429, 5xx) surface in `errors` and `classified_errors` with retry attempts.

For live search, add `--search http` (depends on DuckDuckGo Instant Answer API).
