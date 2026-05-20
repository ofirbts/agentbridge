# WHY — Phase 3: feat/observability

## Problem

Tracer data was persisted as opaque `logs` but **never surfaced** in `inspect` or `explain`. Hiring managers could not see per-step duration, retries, or error types without reading raw JSON files.

## Solution

- Structured `step_logs` on every run (duration, retries, provider, status, error_class)
- `classified_errors` with stable taxonomy (rate_limit, timeout, server_error, …)
- Per-step retry counts (delta per step, not cumulative)
- Rich `inspect` JSON/markdown and `explain` observability section
- Legacy `logs` field still loads into `step_logs`

## Alternatives rejected

| Alternative | Why rejected |
|-------------|--------------|
| Separate `agentbridge logs` command | Extra surface; inspect is the observability command |
| OpenTelemetry export | Phase 4+ scope; violates “no dashboards” |
| Drop plain `errors` array | Breaks existing scripts; keep both |

## Tradeoffs

| Choice | Benefit | Cost |
|--------|---------|------|
| Error classification by message heuristics | No new deps | Imperfect taxonomy |
| Dual errors + classified_errors | Backward compatible | Slight redundancy |
| Markdown table for steps | Readable inspect | Wide terminal output |

## Expected impact

- `inspect` demonstrates reliability thinking clearly
- `explain` sets expectations before execution
- Foundation for Phase 4 benchmarks without new commands
