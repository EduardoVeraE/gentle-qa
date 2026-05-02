#!/usr/bin/env bash
# test-cli-unknown-template.sh — unknown template names must exit 1 and list
# the available templates so users self-correct without reading source.
set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
# shellcheck disable=SC1091
source "$SCRIPT_DIR/../lib.sh"

NAME="cli-unknown-template"
LOG="/tmp/qa-mobile-validation-${NAME}.log"

log_info "[$NAME] running: node $CLI create totally-bogus-template --out /tmp/nope"
set +e
node "$CLI" create totally-bogus-template --out "/tmp/qa-mobile-bogus-out" > "$LOG" 2>&1
ACTUAL_EXIT=$?
set -e

FAILED=0
assert_exit_code 1 "$ACTUAL_EXIT" "$NAME: exits 1 on unknown template" || FAILED=1
assert_grep 'Unknown template'        "$LOG" "$NAME: error mentions unknown template" || FAILED=1
assert_grep 'mobile-bug-report'       "$LOG" "$NAME: error lists available templates" || FAILED=1

if [[ "$FAILED" -ne 0 ]]; then
  log_warn "[$NAME] log tail follows:"
  tail -n 30 "$LOG" >&2 || true
  exit 1
fi
exit 0
