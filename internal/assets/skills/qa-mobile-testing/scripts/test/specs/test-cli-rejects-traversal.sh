#!/usr/bin/env bash
# test-cli-rejects-traversal.sh — `--out` must reject ".." segments, otherwise
# users can be tricked into writing artifacts outside their workspace.
set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
# shellcheck disable=SC1091
source "$SCRIPT_DIR/../lib.sh"

NAME="cli-rejects-traversal"
LOG="/tmp/qa-mobile-validation-${NAME}.log"

log_info "[$NAME] running: node $CLI create mobile-test-plan --out ../escape"
set +e
node "$CLI" create mobile-test-plan --out "../escape" --title "x" > "$LOG" 2>&1
ACTUAL_EXIT=$?
set -e

FAILED=0
assert_exit_code 1 "$ACTUAL_EXIT" "$NAME: exits 1 on traversal" || FAILED=1
assert_grep 'must not contain "\.\."' "$LOG" "$NAME: error message present" || FAILED=1

if [[ "$FAILED" -ne 0 ]]; then
  log_warn "[$NAME] log tail follows:"
  tail -n 30 "$LOG" >&2 || true
  exit 1
fi
exit 0
