#!/usr/bin/env bash
# run-validation.sh — top-level harness for the qa-mobile-testing skill.
# Runs every spec under specs/ and tallies pass/fail.
#
# Unlike the qa-owasp-security harness this does NOT boot a containerized
# target — mobile testing requires real devices/simulators that cannot be
# bootstrapped deterministically in CI. Instead the harness validates the
# DETERMINISTIC surfaces the skill ships:
#   - mobile_artifacts.mjs CLI (list / help / create / --self-test / errors)
#   - example scaffolds parse as JS (detox-rn) and TypeScript (appium-wdio-ts)
#   - example package.json + tsconfig + .detoxrc.js shape
set -euo pipefail

HARNESS_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
# shellcheck disable=SC1091
source "$HARNESS_DIR/lib.sh"

# ---------- pre-flight ----------------------------------------------------

REQUIRED=(node jq npx)
MISSING=()
for tool in "${REQUIRED[@]}"; do
  command -v "$tool" >/dev/null 2>&1 || MISSING+=("$tool")
done
if [[ "${#MISSING[@]}" -gt 0 ]]; then
  log_fail "preflight" "missing required tools: ${MISSING[*]}"
  cat <<HINT >&2
Install hints (macOS / Homebrew):
  brew install node jq
  # npx ships with node
HINT
  exit 2
fi

# Node version sanity — the CLI uses Node 18+ ESM features. 18, 20, 22, 24+ all OK.
NODE_MAJOR=$(node -e 'console.log(process.versions.node.split(".")[0])')
if [[ "$NODE_MAJOR" -lt 18 ]]; then
  log_fail "preflight" "node 18+ required, found $NODE_MAJOR"
  exit 2
fi

# Pre-warm the TypeScript install so the per-spec timing is honest. Skips if
# the cache directory already has typescript installed.
if [[ ! -f "$HARNESS_DIR/.cache/node_modules/typescript/package.json" ]]; then
  log_info "Pre-warming TypeScript cache (one-time, ~10s)"
  node "$HARNESS_DIR/check-ts-syntax.mjs" "$EXAMPLES_DIR/appium-wdio-ts/wdio.conf.ts" >/dev/null
fi

# ---------- run specs -----------------------------------------------------

tally_init

shopt -s nullglob
SPECS=("$HARNESS_DIR"/specs/test-*.sh)
shopt -u nullglob
IFS=$'\n' SPECS=($(printf "%s\n" "${SPECS[@]}" | sort))
unset IFS

if [[ "${#SPECS[@]}" -eq 0 ]]; then
  log_fail "harness" "no specs found under $HARNESS_DIR/specs"
  exit 1
fi

for spec in "${SPECS[@]}"; do
  name=$(basename "$spec" .sh)
  printf "\n%s>>>>%s running %s\n" "$C_CYAN" "$C_RESET" "$name"
  set +e
  bash "$spec"
  rc=$?
  set -e
  if [[ "$rc" -eq 0 ]]; then
    tally_record pass
  else
    tally_record fail
  fi
done

# ---------- report --------------------------------------------------------

if tally_report; then
  exit 0
fi
exit 1
