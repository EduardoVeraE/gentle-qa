# lib.sh — shared helpers for the qa-mobile-testing validation harness.
# Source it; do not execute it. Callers must `set -euo pipefail` themselves.
#
# CRITICAL DESIGN RULE
# --------------------
# NEVER pipe a script's output into another command and capture $? — the
# pipe captures the LAST command's exit status, not the script under test.
# This bug burned the qa-owasp-security harness once already (the entire
# reason this pattern exists).
#
# Always do:
#     set +e
#     <script_under_test> > "$LOG" 2>&1
#     ACTUAL_EXIT=$?
#     set -e
#
# `assert_exit_code` below assumes the caller already captured $? this way.

# Resolve harness paths once, regardless of caller cwd.
HARNESS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SCRIPTS_DIR="$(cd "$HARNESS_DIR/.." && pwd)"
SKILL_DIR="$(cd "$SCRIPTS_DIR/.." && pwd)"
EXAMPLES_DIR="$SKILL_DIR/examples"
TEMPLATES_DIR="$SKILL_DIR/templates"
CLI="$SCRIPTS_DIR/mobile_artifacts.mjs"

# Colours when stdout is a TTY; plain otherwise.
if [[ -t 1 ]]; then
  C_RED=$'\033[0;31m'; C_GREEN=$'\033[0;32m'; C_CYAN=$'\033[0;36m'
  C_YELLOW=$'\033[0;33m'; C_RESET=$'\033[0m'
else
  C_RED=""; C_GREEN=""; C_CYAN=""; C_YELLOW=""; C_RESET=""
fi

log_info() { printf "%s[i]%s %s\n" "$C_CYAN" "$C_RESET" "$*"; }
log_warn() { printf "%s[!]%s %s\n" "$C_YELLOW" "$C_RESET" "$*" >&2; }
log_pass() { printf "%s[PASS]%s %s\n" "$C_GREEN" "$C_RESET" "$*"; }
log_fail() {
  local name="$1"; shift
  printf "%s[FAIL]%s %s — %s\n" "$C_RED" "$C_RESET" "$name" "$*" >&2
}

# Assertions ----------------------------------------------------------------

# assert_exit_code <expected> <actual> <test_name>
assert_exit_code() {
  local expected="$1" actual="$2" name="$3"
  if [[ "$expected" == "$actual" ]]; then
    log_pass "$name (exit=$actual)"
    return 0
  fi
  log_fail "$name" "expected exit $expected, got $actual"
  return 1
}

# assert_grep <pattern> <file> <test_name>
assert_grep() {
  local pattern="$1" file="$2" name="$3"
  if [[ ! -f "$file" ]]; then
    log_fail "$name" "log file missing: $file"
    return 1
  fi
  if grep -Eq "$pattern" "$file"; then
    log_pass "$name (matched /$pattern/)"
    return 0
  fi
  log_fail "$name" "pattern /$pattern/ not found in $file"
  return 1
}

# assert_no_grep <pattern> <file> <test_name>
assert_no_grep() {
  local pattern="$1" file="$2" name="$3"
  if [[ ! -f "$file" ]]; then
    log_fail "$name" "log file missing: $file"
    return 1
  fi
  if grep -Eq "$pattern" "$file"; then
    log_fail "$name" "pattern /$pattern/ unexpectedly matched in $file"
    return 1
  fi
  log_pass "$name (no match for /$pattern/)"
  return 0
}

# assert_file_exists <path> <test_name>
assert_file_exists() {
  local path="$1" name="$2"
  if [[ -f "$path" ]]; then
    log_pass "$name (exists: $path)"
    return 0
  fi
  log_fail "$name" "file not found: $path"
  return 1
}

# Tally ---------------------------------------------------------------------

tally_init() { PASS=0; FAIL=0; }

tally_record() {
  case "$1" in
    pass) PASS=$((PASS + 1)) ;;
    fail) FAIL=$((FAIL + 1)) ;;
    *)    log_warn "tally_record: unknown result '$1'" ;;
  esac
}

tally_report() {
  printf "\n==================== VALIDATION SUMMARY ====================\n"
  printf "PASS: %s%d%s   FAIL: %s%d%s\n" \
    "$C_GREEN" "$PASS" "$C_RESET" "$C_RED" "$FAIL" "$C_RESET"
  printf "============================================================\n"
  [[ "$FAIL" -eq 0 ]]
}
