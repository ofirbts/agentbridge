#!/usr/bin/env bash
set -euo pipefail
input=$(cat)
status=$(echo "$input" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('status',''))" 2>/dev/null || echo "")
if [ "$status" = "completed" ] && git rev-parse --git-dir >/dev/null 2>&1; then
  if ! git diff --quiet || ! git diff --cached --quiet; then
    echo '{"followup_message":"יש שינויים לא committed. להפעיל strict-reviewer לפני commit?"}'
  fi
fi
