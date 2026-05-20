#!/usr/bin/env bash
set -euo pipefail
if ! command -v go >/dev/null 2>&1; then
  if [ -x "${HOME}/sdk/go/bin/go" ]; then
    export PATH="${HOME}/sdk/go/bin:${PATH}"
  else
    echo "Go not found. Install tarball: https://go.dev/dl/ or: sudo apt-get install -y golang-go"
    exit 1
  fi
fi
cd "$(dirname "$0")/.."
go mod tidy
make build
make test
make e2e
echo "agentbridge dev setup OK"
