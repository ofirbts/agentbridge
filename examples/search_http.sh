#!/usr/bin/env bash
set -euo pipefail
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$DIR"
make build
./agentbridge explain --task "open source AI infrastructure tools" --search http
./agentbridge run "open source AI infrastructure tools" --search http --crawler mock
