# WHY — Phase 2: feat/search-provider

## Problem

Search was hard-coded to mock. The README roadmap promised a real provider, but `run` could not demonstrate live URL discovery — weakening portfolio signal for “web-aware execution.”

## Solution

- `SearchProvider` factory: `mock` (default) and `http` (DuckDuckGo Instant Answer API)
- Timeouts and retry-friendly HTTP errors (429, 5xx) aligned with crawl layer
- `--search` flag on `run` and `explain`
- Deterministic mode forces `mock` search for reproducibility

## Alternatives rejected

| Alternative | Why rejected |
|-------------|--------------|
| Google/Bing APIs | Require API keys; bad OSS onboarding |
| HTML scraping SERPs | Fragile, high maintenance |
| New command `search` | Breaks “run unchanged” contract |
| Real search in deterministic mode | Breaks stable run IDs / outputs |

## Tradeoffs

| Choice | Benefit | Cost |
|--------|---------|------|
| DuckDuckGo API | No key, legal lightweight API | Result quality varies |
| Force mock in deterministic | Tests stay stable | `--search http` ignored in deterministic mode |
| Cap 10 results | Bounded crawl load | May drop relevant hits |

## Expected impact

- Demonstrates provider abstraction end-to-end
- `run` UX unchanged (new optional flag only)
- Clear path to swap in other APIs later
