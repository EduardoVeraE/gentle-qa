#!/usr/bin/env bash
# test-cli-list.sh — `mobile_artifacts.mjs list` must succeed and report
# exactly the four templates the skill ships.
set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
# shellcheck disable=SC1091
source "$SCRIPT_DIR/../lib.sh"

NAME="cli-list"
LOG="/tmp/qa-mobile-validation-${NAME}.log"

log_info "[$NAME] running: node $CLI list"
set +e
node "$CLI" list > "$LOG" 2>&1
ACTUAL_EXIT=$?
set -e

FAILED=0
assert_exit_code 0 "$ACTUAL_EXIT" "$NAME: exits 0" || FAILED=1
assert_grep '^mobile-test-plan ' "$LOG"  "$NAME: lists mobile-test-plan"  || FAILED=1
assert_grep '^device-matrix '    "$LOG"  "$NAME: lists device-matrix"     || FAILED=1
assert_grep '^gesture-catalog '  "$LOG"  "$NAME: lists gesture-catalog"   || FAILED=1
assert_grep '^mobile-bug-report ' "$LOG" "$NAME: lists mobile-bug-report" || FAILED=1

if [[ "$FAILED" -ne 0 ]]; then
  log_warn "[$NAME] log tail follows:"
  tail -n 30 "$LOG" >&2 || true
  exit 1
fi
exit 0
