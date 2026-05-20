#!/usr/bin/env bash
set -euo pipefail
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$DIR"
make build
./agentbridge explain --task "find top AI infra startups using web search" --crawler http --format json
./agentbridge run "find top AI infra startups using web search" --mode normal --crawler http
