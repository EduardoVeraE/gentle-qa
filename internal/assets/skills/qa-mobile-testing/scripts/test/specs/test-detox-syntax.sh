#!/usr/bin/env bash
# test-detox-syntax.sh — every JS file the detox-rn example ships must parse
# under `node --check`. Running an actual Detox suite needs an emulator; the
# harness deliberately stops at "the file is syntactically valid JavaScript".
#
# Also asserts the .detoxrc.js evaluates as a CommonJS module exporting an
# object with the keys Detox 20 requires (testRunner, apps, devices,
# configurations).
set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
# shellcheck disable=SC1091
source "$SCRIPT_DIR/../lib.sh"

NAME="detox-syntax"
LOG="/tmp/qa-mobile-validation-${NAME}.log"
EX_DIR="$EXAMPLES_DIR/detox-rn"
: > "$LOG"

JS_FILES=(
  "$EX_DIR/.detoxrc.js"
  "$EX_DIR/e2e/jest.config.js"
  "$EX_DIR/e2e/login.test.js"
  "$EX_DIR/e2e/starter.test.js"
)

FAILED=0
for f in "${JS_FILES[@]}"; do
  if [[ ! -f "$f" ]]; then
    log_fail "$NAME" "missing example file: $f"
    FAILED=1
    continue
  fi
  log_info "[$NAME] node --check $f"
  set +e
  node --check "$f" >> "$LOG" 2>&1
  ACTUAL_EXIT=$?
  set -e
  assert_exit_code 0 "$ACTUAL_EXIT" "$NAME: parses $(basename "$f")" || FAILED=1
done

# Shape check on .detoxrc.js — load it as a module and inspect exports. Detox
# refuses to run when these keys are missing, so this catches dropped keys
# from a future template edit.
log_info "[$NAME] loading .detoxrc.js and inspecting exports"
set +e
node --input-type=commonjs -e "
  const cfg = require('$EX_DIR/.detoxrc.js');
  const required = ['testRunner', 'apps', 'devices', 'configurations'];
  const missing = required.filter((k) => !(k in cfg));
  if (missing.length) {
    console.error('missing keys: ' + missing.join(', '));
    process.exit(1);
  }
  if (!cfg.configurations['ios.sim.debug']) {
    console.error('missing configurations[\"ios.sim.debug\"]');
    process.exit(1);
  }
  if (!cfg.configurations['android.emu.debug']) {
    console.error('missing configurations[\"android.emu.debug\"]');
    process.exit(1);
  }
  console.log('detoxrc shape OK');
" >> "$LOG" 2>&1
ACTUAL_EXIT=$?
set -e
assert_exit_code 0 "$ACTUAL_EXIT" "$NAME: .detoxrc.js shape valid" || FAILED=1

# package.json must declare detox + jest under devDependencies and the four
# canonical scripts.
log_info "[$NAME] inspecting detox-rn/package.json"
set +e
jq -e '
  (.devDependencies.detox != null) and
  (.devDependencies.jest != null) and
  (.scripts["test:ios"] != null) and
  (.scripts["test:android"] != null) and
  (.scripts["build:ios"] != null) and
  (.scripts["build:android"] != null)
' "$EX_DIR/package.json" >> "$LOG" 2>&1
ACTUAL_EXIT=$?
set -e
assert_exit_code 0 "$ACTUAL_EXIT" "$NAME: package.json declares detox+jest+scripts" || FAILED=1

if [[ "$FAILED" -ne 0 ]]; then
  log_warn "[$NAME] log tail follows:"
  tail -n 40 "$LOG" >&2 || true
  exit 1
fi
exit 0
