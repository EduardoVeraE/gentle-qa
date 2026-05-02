#!/usr/bin/env bash
# test-appium-syntax.sh — every TypeScript file the appium-wdio-ts example
# ships must parse via TypeScript's `transpileModule` (parser-only). We do
# NOT run full type checking because that would require installing
# @wdio/types, @wdio/globals, @types/node, @types/mocha — those belong in
# the user's `npm install`, not the harness.
#
# Also validates wdio.conf.ts contains the expected exported `config` and
# package.json declares the canonical scripts.
set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
# shellcheck disable=SC1091
source "$SCRIPT_DIR/../lib.sh"

NAME="appium-syntax"
LOG="/tmp/qa-mobile-validation-${NAME}.log"
EX_DIR="$EXAMPLES_DIR/appium-wdio-ts"
: > "$LOG"

TS_FILES=(
  "$EX_DIR/wdio.conf.ts"
  "$EX_DIR/test/specs/login.spec.ts"
  "$EX_DIR/test/pages/login.page.ts"
)

# Verify each file exists before invoking the checker.
FAILED=0
for f in "${TS_FILES[@]}"; do
  if [[ ! -f "$f" ]]; then
    log_fail "$NAME" "missing example file: $f"
    FAILED=1
  fi
done

if [[ "$FAILED" -eq 0 ]]; then
  log_info "[$NAME] running TypeScript parser on ${#TS_FILES[@]} files"
  set +e
  node "$SCRIPT_DIR/../check-ts-syntax.mjs" "${TS_FILES[@]}" >> "$LOG" 2>&1
  ACTUAL_EXIT=$?
  set -e
  assert_exit_code 0 "$ACTUAL_EXIT" "$NAME: all TS files parse" || FAILED=1
fi

# Sanity-check tsconfig.json shape — the example must keep strict + ES2022.
log_info "[$NAME] inspecting tsconfig.json"
set +e
jq -e '
  (.compilerOptions.strict == true) and
  (.compilerOptions.target == "ES2022") and
  (.include | type == "array")
' "$EX_DIR/tsconfig.json" >> "$LOG" 2>&1
ACTUAL_EXIT=$?
set -e
assert_exit_code 0 "$ACTUAL_EXIT" "$NAME: tsconfig.json shape valid" || FAILED=1

# package.json must declare the WebdriverIO toolchain and `test` scripts.
log_info "[$NAME] inspecting appium-wdio-ts/package.json"
set +e
jq -e '
  (.devDependencies["@wdio/cli"] != null) and
  (.devDependencies["appium"] != null) and
  (.devDependencies["typescript"] != null) and
  (.scripts.test != null) and
  (.scripts["test:ios"] != null) and
  (.scripts["test:android"] != null)
' "$EX_DIR/package.json" >> "$LOG" 2>&1
ACTUAL_EXIT=$?
set -e
assert_exit_code 0 "$ACTUAL_EXIT" "$NAME: package.json declares wdio+appium+scripts" || FAILED=1

# wdio.conf.ts must expose a `config` export — wdio CLI imports it by name.
assert_grep 'export const config' "$EX_DIR/wdio.conf.ts" "$NAME: wdio.conf.ts exports config" || FAILED=1

if [[ "$FAILED" -ne 0 ]]; then
  log_warn "[$NAME] log tail follows:"
  tail -n 40 "$LOG" >&2 || true
  exit 1
fi
exit 0
