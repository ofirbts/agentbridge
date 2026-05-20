# Contributing to AgentBridge

Thanks for your interest. This project optimizes for **small, reviewable PRs** and a releasable `main` at all times.

## Before you start

1. Read [README.md](README.md) and [docs/DECISION.md](docs/DECISION.md).
2. Check [CHANGELOG.md](CHANGELOG.md) and open issues for large ideas first.
3. Out of scope: web UI, SaaS, databases, auth, browser automation, distributed systems.

## Development setup

```bash
git clone https://github.com/ofirbts/agentbridge
cd agentbridge
make build
make test
```

Go 1.23+ required. Optional: `bash scripts/setup-dev.sh` for full smoke (build + test + e2e).

## Pull requests

1. Branch from `main` (`feat/…`, `fix/…`, `chore/…`).
2. Keep PRs focused — one concern per PR.
3. Run before opening:

```bash
make test
go test -race ./...
go vet ./...
```

4. Add or update tests for behavior changes.
5. No secrets in commits (use `.env.example` patterns only).

## Code guidelines

- Match existing package layout: `cmd/`, `internal/`, `pkg/mcp/`.
- Prefer extending provider interfaces over new frameworks.
- Prefer deletion over addition when simplifying.
- No comments unless they explain non-obvious invariants.

## Commit messages

Use imperative, concise subjects:

```
feat: add search provider timeout
fix: surface crawl errors in inspect
chore: update CI go version
```

## Review expectations

PRs should state:

- **Problem** — what user pain this solves
- **Approach** — what changed and why
- **Test plan** — commands run locally

## Questions

Open a GitHub issue with the `feature request` or `bug report` template.
