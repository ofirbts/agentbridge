#!/usr/bin/env bash
set -euo pipefail
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$DIR"
make build
./agentbridge run --mode deterministic --crawler mock "get latest AI infra notes from GitHub"
RUN_ID=$(./agentbridge run --mode deterministic --crawler mock "get latest AI infra notes from GitHub" | python3 -c "import sys,json; print(json.load(sys.stdin)['run_id'])")
./agentbridge inspect "$RUN_ID" --format json
