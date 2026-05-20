---
name: supabase-dev
description: Supabase development workflow via MCP. Use for database schema, migrations, edge functions, logs, and TypeScript type generation.
---

# Supabase Dev

Before schema changes: call list_tables.

## Schema change flow

1. list_tables — understand current structure
2. apply_migration — DDL only through migrations
3. generate_typescript_types — if the app uses generated types
4. get_advisors — check security/performance after migration

## Debugging

1. get_logs — auth, postgres, edge-function, or api as needed
2. execute_sql — read-only SELECT for investigation

Never invent credentials. Use get_project_url and get_publishable_keys for client config.
