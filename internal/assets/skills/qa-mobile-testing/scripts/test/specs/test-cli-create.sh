#!/usr/bin/env bash
# test-cli-create.sh — `mobile_artifacts.mjs create` must produce a rendered
# artifact with NO leftover {{placeholder}} tokens and the manifest comments
# stripped. Validates the most-used CLI surface.
#
# Strategy: introspect placeholders for the chosen template via `help`, then
# pass a sample value for every one. This keeps the spec resilient to
# template-content changes without losing the "no leftover placeholders"
# invariant.
set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
# shellcheck disable=SC1091
source "$SCRIPT_DIR/../lib.sh"

NAME="cli-create"
LOG="/tmp/qa-mobile-validation-${NAME}.log"
OUT_DIR="/tmp/qa-mobile-validation-${NAME}-out"
TEMPLATE="mobile-bug-report"
rm -rf "$OUT_DIR"
mkdir -p "$OUT_DIR"

# Discover placeholders, build flag pairs.
PLACEHOLDERS=$(node "$CLI" help "$TEMPLATE" \
  | awk '/^Placeholders/{found=1; next} found && /^  --/{sub(/^  --/, ""); sub(/[ ].*$/, ""); print}')

ARGS=(create "$TEMPLATE" --out "$OUT_DIR" --strip-hints --title "harness sample bug")
for p in $PLACEHOLDERS; do
  # title is already provided above (it slugs the filename); skip duplicates.
  if [[ "$p" == "title" ]]; then continue; fi
  ARGS+=("--$p" "SAMPLE_${p}")
done

log_info "[$NAME] running: node $CLI ${ARGS[*]}"
set +e
node "$CLI" "${ARGS[@]}" > "$LOG" 2>&1
ACTUAL_EXIT=$?
set -e

FAILED=0
assert_exit_code 0 "$ACTUAL_EXIT" "$NAME: exits 0" || FAILED=1

# Find the rendered file (CLI prints `Created: <path>` on stdout).
RENDERED=$(grep -E '^Created: ' "$LOG" | head -n1 | sed 's/^Created: //')
if [[ -z "$RENDERED" ]]; then
  log_fail "$NAME" "no 'Created: <path>' line in log"
  FAILED=1
else
  assert_file_exists "$RENDERED" "$NAME: rendered file present" || FAILED=1
  # No {{placeholder}} tokens should remain after a fully-flagged create.
  if grep -qE '\{\{[a-z_0-9]+\}\}' "$RENDERED"; then
    LEFT=$(grep -oE '\{\{[a-z_0-9]+\}\}' "$RENDERED" | sort -u | tr '\n' ' ')
    log_fail "$NAME" "leftover placeholders in rendered file: $LEFT"
    FAILED=1
  else
    log_pass "$NAME: no leftover placeholders in rendered file"
  fi
  # Manifest comments must be stripped.
  if grep -qE '^<!--\s*(Skill|Placeholders):' "$RENDERED"; then
    log_fail "$NAME" "manifest comments not stripped from rendered file"
    FAILED=1
  else
    log_pass "$NAME: manifest comments stripped"
  fi
  # --strip-hints must remove "<!-- e.g., ... -->" hints.
  if grep -qE '<!--\s*e\.g\.,' "$RENDERED"; then
    log_fail "$NAME" "--strip-hints did not remove inline hints"
    FAILED=1
  else
    log_pass "$NAME: --strip-hints removed inline hints"
  fi
fi

if [[ "$FAILED" -ne 0 ]]; then
  log_warn "[$NAME] log tail follows:"
  tail -n 40 "$LOG" >&2 || true
  exit 1
fi
exit 0
