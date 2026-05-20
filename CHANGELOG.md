# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-05-20

### Added

- CLI: `run`, `explain`, `inspect`, `simulate-failure`
- Workflow engine with retry/backoff and deterministic mode
- Mock search, mock crawl, HTTP crawler (`--crawler http`)
- MCP stub and MCP HTTP JSON-RPC client (`--mcp-transport`)
- File-backed run store (`.agentbridge/runs`)
- JSON and markdown output for `run` and `inspect`
- Examples under `examples/`
- CI: build, test, vet, race detector
- MIT LICENSE

### Changed

- Module path: `github.com/ofirbts/agentbridge`
- Unified run model (`internal/run.Run`) for engine and inspect
- Mock crawler reported as `mock` (was `mock_http`)

### Removed

- Unused `--config` flag
- Dead code: `BuildPlan`, `FormatStepLogs`, undocumented `real` crawler alias

### Fixed

- MCP client initialization errors surface via `NewEngine`
- Mock search URL uses `https://example.com` for successful HTTP demos

## [Unreleased]

### Added

- `make benchmark`, `make profile`, `make race`
- Package benchmarks (engine, search, run store, classify)
- `PERFORMANCE.md` with measured results

### Changed

- Step logs: plan/mcp/normalize use explicit provider metadata

### Added (prior)

- `step_logs` on runs: per-step duration, retries, provider, status, error_class
- `classified_errors` with taxonomy (rate_limit, timeout, server_error, …)
- `explain` observability section; `inspect` providers map in JSON
- Legacy run files with `logs` field load into `step_logs`
- `--search` flag: `mock` (default) or `http` (DuckDuckGo Instant Answer API)
- `internal/search` HTTP provider; `examples/search_http.sh`

### Changed

- Deterministic mode always uses mock search
- Per-step retry counts (not cumulative) in step logs
- Docs: DECISION.md and INFRA-NOTES.md reflect MCP HTTP and search providers

### Planned

- Benchmarks and performance notes (Phase 4)

[0.1.0]: https://github.com/ofirbts/agentbridge/releases/tag/v0.1.0
