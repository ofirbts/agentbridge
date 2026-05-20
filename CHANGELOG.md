# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

Nothing yet.

## [0.2.0] - 2026-05-20

### Added

- HTTP search provider (`--search http`, DuckDuckGo Instant Answer API)
- `step_logs` and `classified_errors` on runs; richer `inspect` / `explain`
- `make benchmark`, `make profile`, `make race`; `PERFORMANCE.md`
- `docs/use-cases.md`, good first issue template
- Examples: `search_http.sh`

### Changed

- Deterministic mode forces mock search
- Per-step retry counts in step logs
- Step log provider metadata on all engine steps
- Documentation refresh (README v0.1 structure, DECISION/INFRA-NOTES drift fixes)

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

[Unreleased]: https://github.com/ofirbts/agentbridge/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/ofirbts/agentbridge/releases/tag/v0.2.0
[0.1.0]: https://github.com/ofirbts/agentbridge/releases/tag/v0.1.0
