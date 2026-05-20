#!/usr/bin/env bash
set -euo pipefail
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$DIR"
make build
./agentbridge simulate-failure --seed 42 --count 3 || true
./agentbridge run "retry search after failure simulation" --mode normal --crawler mock
./agentbridge explain --task "retry search after failure simulation" --format text
