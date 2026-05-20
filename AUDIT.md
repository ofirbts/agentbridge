# AgentBridge — Phase 0 Repository Audit

**Date:** 2026-05-20  
**Auditor role:** Lead engineer / release manager  
**Scope:** Signal improvement without new features  
**Assumption:** PRs #1–#3, tests, CI, and MCP HTTP are complete on `main`

---

## Executive summary

AgentBridge is a **coherent, test-backed Go CLI** with a clear thesis (execution reliability for web-aware agents). The codebase is small (~2.5k LOC Go), readable, and already demonstrates AI-infra thinking.

The main gap for **portfolio / hiring signal** is not missing features—it is **packaging, narrative clarity, and a few internal rough edges** (dead code, duplicate result types, silent error paths) that dilute the “polished infrastructure tool” story.

**Verdict:** Strong foundation. Phase 0 should be a **cleanup + docs-signal PR** (no behavior change) before v0.1.0 release work.

---

## Repository structure

```
agentbridge/
├── cmd/                 # Cobra CLI (run, explain, inspect, simulate-failure)
├── internal/
│   ├── workflow/        # Engine orchestration + Plan (central hub)
│   ├── search/          # SearchProvider + mock
│   ├── crawl/           # Crawler + mock + http
│   ├── extract/         # HTML extract + normalize
│   ├── failure/         # Retry policy + simulator
│   ├── observability/   # Tracer + step logs (underused in UX)
│   ├── run/             # Run model + file store
│   └── output/          # JSON/markdown formatting
├── pkg/mcp/             # MCP Client (stub + HTTP JSON-RPC)
├── tests/               # CLI e2e (go run)
├── examples/            # Shell demos
├── docs/                # DECISION, RULES, INFRA-NOTES
└── .github/workflows/   # CI (build, test, vet)
```

**Dependency hub:** `internal/workflow` imports all provider packages + `pkg/mcp`. Acceptable at current size; watch growth.

---

## Strengths (recruiter signal)

| Area | Evidence |
|------|----------|
| **Clear product thesis** | Execution layer, not framework; documented in `docs/DECISION.md` |
| **Go maturity** | Interfaces for providers, table-driven tests, `go vet`, race-safe tests run locally |
| **CLI design** | Cobra, consistent JSON stdout, `explain` vs `run` separation |
| **Reliability thinking** | Retry policy, failure simulator, deterministic mode with stable run IDs |
| **Test pyramid** | Unit per package + `tests/e2e_test.go` CLI flow |
| **Incremental delivery** | Mock-first, HTTP crawler, MCP stub → HTTP transport |
| **CI** | `.github/workflows/ci.yml` on push/PR |

---

## Risks

| Risk | Impact | Severity |
|------|--------|----------|
| **No LICENSE / CHANGELOG / CONTRIBUTING** | OSS credibility; recruiters check these first | High |
| **README not “60-second”** | Strong content but long; hero/demo buried | High |
| **Dead code & silent paths** | Suggests unfinished polish | Medium |
| **`workflow.Result` vs `run.Run` duplication** | Confusing for contributors; `inspect` adapter hack | Medium |
| **Observability collected but hidden** | `Tracer` logs stored on `Run.Logs` but not in `inspect --format markdown` | Medium |
| **Version drift** | `project.config.json` says `0.2.0`, no git tag, no CHANGELOG | Medium |
| **MCP client errors swallowed** | `NewEngine` ignores `mcp.NewClient` error | Medium |
| **Naming inconsistency** | `mock` vs `mock_http`, undocumented `real` crawler alias | Low |
| **Docs fragmentation** | 4 doc locations (README + 3 under `docs/`) with overlap | Low |
| **CI gap** | No `-race` in CI (only local discipline) | Low |
| **`.cursor/` in repo** | Useful for you; neutral/negative for external readers unless framed | Low |

---

## Dead code

| Item | Location | Notes |
|------|----------|-------|
| `BuildPlan()` | `internal/workflow/plan.go` | Never referenced |
| `FormatStepLogs()` | `internal/output/log_formatter.go` | Never referenced |
| `Recorder.AddError()` | `internal/observability/recorder.go` | Never called |
| `ConfigFile` field | `workflow.Config` + `cmd/run.go` | Flag exists, never loaded |

**Estimated fix:** Small (S) — delete or wire up intentionally.

---

## Duplicate abstractions

| Duplication | Files | Recommendation |
|-------------|-------|----------------|
| **Run record vs CLI result** | `run.Run` ≈ `workflow.Result` | Single domain type in `internal/run`, engine returns it; output formats accept `run.Run` |
| **Inspect markdown path** | `cmd/inspect.go` rebuilds `workflow.Result` | Format from `run.Run` directly |
| **JSON formatting** | `output.FormatRunResult` vs raw `inspect` encode | One formatter for run records |

**Estimated fix:** Medium (M) — one PR, behavior-preserving tests.

---

## Naming inconsistencies

| Issue | Current | Suggested |
|-------|---------|-----------|
| Mock crawler name | `Name()` → `mock_http` | `mock` (or document why `_http`) |
| Plan provider label | `mock` + `mock_crawler` | Align with `Crawler.Name()` |
| Crawler flag values | `mock`, `http`, hidden `real` | Document in README; remove `real` or alias in help text |
| Run ID field | struct `ID`, JSON `run_id` | OK; document in RULES |
| Module path | `github.com/ofirbts/agentbridge` | Consistent ✓ |

**Estimated fix:** Small (S) — mostly docs + one rename if behavior unchanged.

---

## Package coupling

```
cmd ──► workflow ──► search, crawl, extract, failure, observability, run
                 └──► pkg/mcp

cmd ──► output ──► workflow   ← coupling: output should not depend on workflow

cmd ──► run (inspect)
```

**Concerns:**

1. `internal/output` imports `internal/workflow` only for `Result` type — invert: move shared view model to `internal/run` or `internal/api`.
2. `pkg/mcp` is public but only consumed by `workflow` — acceptable for “integration surface” story; otherwise move to `internal/mcp`.

**Estimated fix:** Medium (M) — optional for v0.1.0; do before v1.

---

## Behavioral gaps (not bugs, signal gaps)

1. **Step-level observability** is recorded but **not shown** in default `inspect` output (Phase 3 aligns).
2. **`--config` flag** advertises future config file — misleading until implemented or removed.
3. **Examples** (`examples/*.sh`) are good but not linked from a single “Demo” entry in README hero.
4. **RULES.md** says “mocks only” for search — outdated after HTTP crawler (docs drift).

---

## Proposed fixes (Phase 0 — no new features)

Prioritized for **signal / clarity** only. Each should be a tiny PR; `main` stays releasable.

| ID | Fix | Complexity | User-visible? |
|----|-----|------------|---------------|
| P0-1 | Remove dead code (`BuildPlan`, `FormatStepLogs`, `AddError` or use) | S | No |
| P0-2 | Remove or stub `ConfigFile` until implemented | S | Yes (flag honesty) |
| P0-3 | Fix silent `mcp.NewClient` error → surface in run errors | S | Yes |
| P0-4 | Unify `run.Run` / `workflow.Result` (single type) | M | No |
| P0-5 | Decouple `output` from `workflow` | M | No |
| P0-6 | Align crawler naming (`mock` vs `mock_http`) + help text | S | Yes |
| P0-7 | Update `docs/RULES.md` drift (HTTP crawler, MCP HTTP) | S | No |
| P0-8 | Add `LICENSE` (MIT recommended) | S | Yes |
| P0-9 | CI: add `go test -race ./...` job step | S | No |

**Defer to Phase 1 (not Phase 0):** README rewrite, CHANGELOG, CONTRIBUTING, diagrams, screenshots.

**Reject for now:** new providers, databases, UI, refactors >300 LOC without signal gain.

---

## 60-second comprehension test (current)

| Question | Can a stranger answer in 60s? |
|----------|-------------------------------|
| What is it? | Partially — README explains but slowly |
| Why not LangChain? | Yes — DECISION.md (if they open it) |
| How do I try it? | Yes — `make build` + `run` |
| What’s unique? | Buried — deterministic mode + execution layer |

**Phase 1 target:** Hero + 3-line demo + architecture diagram above the fold.

---

## Complexity estimates (Phase 0 total)

| Bucket | PRs | Effort |
|--------|-----|--------|
| S | 6–7 | ~1–2 days |
| M | 2 | ~1–2 days |
| **Total Phase 0** | **8–9 tiny PRs** | **~2–4 days** |

---

## Recommended PR sequence (after approval)

1. `chore/remove-dead-code` (P0-1)
2. `fix/mcp-client-error-surface` (P0-3)
3. `chore/remove-config-flag` or `feat/config-load` — prefer **remove** per “prefer deletion” (P0-2)
4. `chore/crawler-naming` (P0-6)
5. `docs/rules-drift` (P0-7)
6. `chore/license-mit` (P0-8)
7. `ci/race-detector` (P0-9)
8. `refactor/unify-run-type` (P0-4 + P0-5) — largest; split if needed

Each PR requires `WHY.md` + `REVIEW.md` per project process.

---

## STATUS

| Item | State |
|------|--------|
| Phase 0 audit | **Complete** |
| Implementation | **Not started** (waiting approval) |
| `main` | Releasable; tests pass |
| OSS packaging | Incomplete (no LICENSE/CHANGELOG/tags) |

## NEXT STEP

**Await approval** to execute Phase 0 PRs (starting with P0-1 + P0-3), **or** skip to Phase 1 `release/v0.1.0` if you prefer docs-first.

**Recommendation:** Approve Phase 0 quick wins (P0-1, P0-2, P0-3, P0-8) in one sitting, then Phase 1 README/LICENSE for v0.1.0 tag.

## RISKS

- **Scope creep:** Phase 0 must not become Phase 2 (real search).
- **Refactor PR (P0-4):** Highest regression risk; keep tests strict.
- **Recruiter skim:** Without LICENSE + sharp README, strong code is under-valued.

## RECOMMENDATION

**Approve Phase 0 as defined.** Do not start Phase 1 until dead code and misleading `--config` flag are resolved—these are negative signal if a hiring manager clones and greps the repo.

---

*End of Phase 0. Do not proceed automatically.*
