# WHY — Phase 5: open-source packaging

## Problem

v0.1.0 tag existed but most portfolio value shipped after it. No use-case index, no contributor onboarding issue template, CHANGELOG `[Unreleased]` was cluttered, and GitHub had no v0.2.0 release artifact.

## Solution

- Consolidate changelog into **v0.2.0** (v1 candidate baseline)
- `docs/use-cases.md` + examples table in README
- Good first issue template + issue chooser config
- Release badge; tag `v0.2.0` with release notes from CHANGELOG

## Alternatives rejected

| Alternative | Why rejected |
|-------------|--------------|
| Jump to v1.0.0 | Scope still CLI-only; 0.2.0 honest semver |
| GitHub Pages site | Extra maintenance; README is primary |
| Auto-release Action | YAGNI for solo OSS portfolio |

## Tradeoffs

| Choice | Benefit | Cost |
|--------|---------|------|
| v0.2.0 not v1.0 | Signals maturity without overclaiming | Another version to track |
| docs/use-cases vs wiki | Lives in repo | Another doc file |

## Expected impact

- Recruiters see releases, use cases, and contribution path in one glance
- Clear “first issue” lowers friction for drive-by contributors
