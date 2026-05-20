# Quickstart

From repo root after `make build`.

## Command

```bash
./agentbridge explain --task "find AI infra startups"
./agentbridge run "find AI infra startups" --mode deterministic
./agentbridge inspect run_det_23434e37 --format markdown
```

Use the `run_id` from your `run` output (not a placeholder).

## Output — explain

```json
{
  "task": "find AI infra startups",
  "mode": "normal",
  "search": "mock",
  "crawler": "mock",
  "steps": [
    { "name": "search", "provider": "mock", "description": "resolve candidate URLs for the task" },
    { "name": "crawl", "provider": "mock", "description": "fetch page content with retries" },
    { "name": "extract", "provider": "extractor", "description": "strip HTML and extract text" },
    { "name": "normalize", "provider": "normalizer", "description": "structure output for downstream agents" }
  ],
  "observability": {
    "inspect_fields": ["duration_ms", "retries", "step_logs", "classified_errors", "search", "crawler", "mode"]
  }
}
```

## Output — run (excerpt)

```json
{
  "status": "success",
  "run_id": "run_det_23434e37",
  "mode": "deterministic",
  "steps": ["search", "crawl", "extract", "normalize"],
  "result": {
    "query": "find AI infra startups",
    "url": "https://example.com",
    "normalized": { "text": "AgentBridge mock page", "word_count": 3 }
  }
}
```

## Output — inspect

```markdown
# Run run_det_23434e37

- Status: success
- Mode: deterministic
...
## Step timeline

| Step | Status | Provider | Duration | Retries | Error class |
| plan | ok | workflow | 0ms | 0 | - |
| search | ok | mock | 0ms | 0 | - |
...
```

## What this means

1. **explain** — plan only; no network.
2. **run** — executes the pipeline; saves a run under `.agentbridge/runs/`.
3. **inspect** — reads that run later; step timeline shows what happened per stage.

This is the core loop: execute → persist → inspect.
