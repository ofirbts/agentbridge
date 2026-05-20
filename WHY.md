# WHY — Phase 4: stability and performance

## Problem

No baseline for latency/allocs; step logs had inconsistent provider metadata on plan/mcp steps; contributors lacked a single command for benchmarks and race checks.

## Solution

- `make benchmark`, `make profile`, `make race`
- Benchmarks for engine, search, run store, classify
- `PERFORMANCE.md` with measured results and tuning notes
- Step log metadata: all engine steps use `AddStepWithMeta` with explicit providers

## Alternatives rejected

| Alternative | Why rejected |
|-------------|--------------|
| pprof HTTP server | Out of scope for CLI tool |
| Optimizing away disk persist | Hides real CLI cost; document instead |
| Benchmarks in CI gates | Flaky across machines |

## Tradeoffs

| Choice | Benefit | Cost |
|--------|---------|------|
| Benchmark includes disk persist | Realistic | Higher ns/op vs memory-only |
| Document don't optimize yet | Honest portfolio signal | Not "blazing fast" marketing |

## Expected impact

- Shows production awareness without over-engineering
- Gives Phase 5 packagers concrete numbers
