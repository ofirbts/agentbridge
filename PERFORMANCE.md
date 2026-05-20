# AgentBridge Performance Notes

Measured on Linux amd64, Go 1.23, Intel i5-7400 @ 3.00GHz (`make benchmark`, `-count=3`).

## Summary

AgentBridge is **I/O-bound** when using real HTTP providers. The mock deterministic pipeline is ~**5–7 ms/op** end-to-end including run-store persistence. Hot paths are suitable for CLI batch use; not tuned for high-QPS server deployment (non-goal).

## Benchmarks

| Benchmark | ns/op | B/op | allocs/op |
|-----------|------:|-----:|----------:|
| `BenchmarkEngineRunTaskDeterministic` | ~6.5M | ~2.2M | ~77 |
| `BenchmarkMockSearch` | ~410 | 312 | 6 |
| `BenchmarkRunStoreInsertGet` | ~120 | 224 | 1 |
| `BenchmarkClassifyError` | ~69 | 0 | 0 |

Run locally:

```bash
make benchmark
make profile   # writes cpu.prof / mem.prof under internal/workflow/
make race
```

## Findings

### End-to-end run (mock, deterministic)

- Dominated by **JSON marshal + disk write** to `.agentbridge/runs` each `RunTask`.
- In-memory workflow (search, crawl, extract, normalize) is sub-millisecond with mocks.
- Recommendation: for high-frequency runs, add optional `--store memory` (future) or batch persistence.

### Search / classify

- Mock search: sub-microsecond per op excluding allocator overhead.
- Error classification: **zero allocations** per call — safe on hot error paths.

### Run store

- `Get` after warm memory cache: ~120 ns/op, 1 alloc (defensive copy).
- Disk cold-read not benchmarked separately; expect file I/O on first `inspect`.

## Memory

- Typical deterministic run allocates ~2 MB/op (mostly JSON + step log slices).
- No unbounded growth: step logs scale with fixed pipeline steps (~6–8 entries per run).
- Tracer stores all step logs in memory until persist; bounded by step count.

## Race detector

CI and `make race` run `go test -race ./...`. No races reported on `main` at Phase 4.

## Profiling

```bash
make profile
go tool pprof -top cpu.prof
go tool pprof -top mem.prof
```

Expect CPU time in `encoding/json`, `os.WriteFile`, and extract/normalize on larger HTML payloads.

## Non-goals

- Sub-millisecond SLA for full pipeline with disk persistence
- Parallel step execution
- Connection pooling for HTTP (single client per provider instance today)

## Phase 5+ ideas (not implemented)

- Optional in-memory store for tests/CI
- Reuse `json.Encoder` buffer for run persistence
- HTTP client transport tuning for real search/crawl only
