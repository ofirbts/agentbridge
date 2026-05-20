#!/usr/bin/env bash
set -euo pipefail
if ! command -v go >/dev/null 2>&1; then
  echo "Go not found. Install: sudo apt-get update && sudo apt-get install -y golang-go"
  exit 1
fi
cd "$(dirname "$0")/.."
go mod tidy
make build
make test
make e2e
echo "agentbridge dev setup OK"
