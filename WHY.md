# WHY — Phase 1: release/v0.1.0

## Problem

The codebase was functionally complete but **hard to evaluate in under 60 seconds**. Missing OSS staples (CHANGELOG, CONTRIBUTING), README buried the value prop, and version story was unclear for recruiters and hiring managers.

## Solution

Ship documentation-only release packaging:

- Restructured README (hero → problem → demo → quick start → architecture → commands → decisions → examples → roadmap)
- CHANGELOG for v0.1.0
- CONTRIBUTING for contribution norms
- Mermaid architecture diagram and real terminal output examples

No runtime behavior changes.

## Alternatives rejected

| Alternative | Why rejected |
|-------------|--------------|
| Wait until v1 for docs | Delays portfolio signal |
| Single wiki / GitHub Pages | Extra infra; README is enough for v0.1.0 |
| Auto-generated README from code | Worse narrative for AI-infra positioning |
| PNG screenshots in repo | Heavy maintenance; fenced real output is enough |

## Tradeoffs

| Choice | Benefit | Cost |
|--------|---------|------|
| Docs-only PR | Zero regression risk | No new features |
| v0.1.0 tag name | Clear “first release” story | `project.config` was 0.2.0 — aligned to 0.1.0 |
| Mermaid in README | Renders on GitHub | Not visible in plain-text clones |

## Expected impact

- Faster comprehension for new readers
- Stronger open-source credibility (LICENSE already present)
- Clear baseline for Phase 2+ PRs
