#!/usr/bin/env bash
# test-cli-self-test.sh — defer to the CLI's built-in `--self-test`. Renders
# every template with sample values and asserts no leftover placeholders.
# Cheap regression guard for template + CLI co-evolution.
set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
# shellcheck disable=SC1091
source "$SCRIPT_DIR/../lib.sh"

NAME="cli-self-test"
LOG="/tmp/qa-mobile-validation-${NAME}.log"

log_info "[$NAME] running: node $CLI --self-test"
set +e
node "$CLI" --self-test > "$LOG" 2>&1
ACTUAL_EXIT=$?
set -e

FAILED=0
assert_exit_code 0 "$ACTUAL_EXIT" "$NAME: exits 0"          || FAILED=1
assert_grep '^self-test: OK' "$LOG" "$NAME: self-test OK"   || FAILED=1
assert_no_grep '^FAIL'       "$LOG" "$NAME: no FAIL lines"  || FAILED=1

if [[ "$FAILED" -ne 0 ]]; then
  log_warn "[$NAME] log tail follows:"
  tail -n 40 "$LOG" >&2 || true
  exit 1
fi
exit 0
