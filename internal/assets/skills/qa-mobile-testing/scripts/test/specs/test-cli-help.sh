#!/usr/bin/env bash
# test-cli-help.sh — `mobile_artifacts.mjs help <template>` must surface the
# placeholder manifest. Without this signal users cannot discover required
# flags for `create`.
set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
# shellcheck disable=SC1091
source "$SCRIPT_DIR/../lib.sh"

NAME="cli-help"
LOG="/tmp/qa-mobile-validation-${NAME}.log"

log_info "[$NAME] running: node $CLI help mobile-test-plan"
set +e
node "$CLI" help mobile-test-plan > "$LOG" 2>&1
ACTUAL_EXIT=$?
set -e

FAILED=0
assert_exit_code 0 "$ACTUAL_EXIT" "$NAME: exits 0"                          || FAILED=1
assert_grep '^Template: mobile-test-plan' "$LOG" "$NAME: shows template name" || FAILED=1
assert_grep '^Placeholders \([0-9]+\):'   "$LOG" "$NAME: shows placeholder header" || FAILED=1
# A few well-known placeholders from templates/mobile-test-plan.md.
assert_grep '\-\-project'   "$LOG" "$NAME: lists --project flag"   || FAILED=1
assert_grep '\-\-release'   "$LOG" "$NAME: lists --release flag"   || FAILED=1
assert_grep '\-\-platforms' "$LOG" "$NAME: lists --platforms flag" || FAILED=1

if [[ "$FAILED" -ne 0 ]]; then
  log_warn "[$NAME] log tail follows:"
  tail -n 30 "$LOG" >&2 || true
  exit 1
fi
exit 0
