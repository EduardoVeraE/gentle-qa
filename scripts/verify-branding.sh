#!/usr/bin/env bash
set -euo pipefail

# ============================================================================
# gentle-qa — Branding Verification Script
# Scans the repo for prohibited upstream branding tokens.
#
# Usage:
#   scripts/verify-branding.sh             # scan working tree (default)
#   scripts/verify-branding.sh --upstream  # scan upstream/main diff
#   scripts/verify-branding.sh --verbose   # show every match with context
#
# Exit codes:
#   0 — verification passed (no leaks)
#   1 — verification failed (leaks found)
#   2 — pre-flight failure (missing tool, bad args, etc.)
# ============================================================================

# ============================================================================
# Allowlist — paths that are legitimate exceptions (grep -v patterns)
# Add new entries here if a file needs to reference upstream tokens by design.
# ============================================================================
ALLOWLIST=(
    "CHANGELOG.md"
    "CONTRIBUTORS.md"
    "internal/components/filemerge/section_test.go"
    "internal/components/sdd/inject_test.go"
    "scripts/verify-branding.sh"
    "scripts/sync-upstream.sh"
    "internal/assets/skills/upstream-sync/SKILL.md"
)

# ============================================================================
# Prohibited tokens — combined into a single alternation regex for rg
# ============================================================================
PROHIBITED_PATTERN='gentle-ai|gentleAi|gentleAI|GENTLE_AI|gentleai|Gentle-AI|Gentle AI|gentleman-programming/gentle-ai|Gentleman-Programming/gentle-ai'

# ============================================================================
# Color support
# ============================================================================
setup_colors() {
    if [ -t 1 ] && [ "${TERM:-}" != "dumb" ]; then
        RED='\033[0;31m'
        GREEN='\033[0;32m'
        YELLOW='\033[1;33m'
        BOLD='\033[1m'
        DIM='\033[2m'
        NC='\033[0m'
    else
        RED='' GREEN='' YELLOW='' BOLD='' DIM='' NC=''
    fi
}

# ============================================================================
# Helpers
# ============================================================================
info()  { echo -e "${DIM}[verify-branding]${NC} $*"; }
ok()    { echo -e "${GREEN}✓${NC} $*"; }
fail()  { echo -e "${RED}✗${NC} $*" >&2; }
warn()  { echo -e "${YELLOW}[warn]${NC} $*"; }

die() { fail "$*"; exit 2; }

# ============================================================================
# Build rg glob exclusions from ALLOWLIST
# ============================================================================
build_rg_globs() {
    local args=()
    for entry in "${ALLOWLIST[@]}"; do
        args+=("--glob=!${entry}")
    done
    # Also exclude VCS / generated dirs unconditionally
    args+=(
        "--glob=!.git/**"
        "--glob=!node_modules/**"
        "--glob=!vendor/**"
        "--glob=!dist/**"
        "--glob=!build/**"
    )
    printf '%s\n' "${args[@]}"
}

# Read globs into an array portably (no mapfile — macOS bash 3.x compat)
read_globs_array() {
    local IFS=$'\n'
    local line
    RG_GLOB_ARGS=()
    while IFS= read -r line; do
        RG_GLOB_ARGS+=("$line")
    done < <(build_rg_globs)
}

# ============================================================================
# Scan working tree
# ============================================================================
scan_working_tree() {
    local verbose="$1"

    read_globs_array

    # Collect matches; --verbose adds 2 lines of context around each hit
    local match_output
    if [ "$verbose" = "1" ]; then
        match_output="$(rg --line-number --with-filename -C 2 \
            "${RG_GLOB_ARGS[@]}" \
            -e "$PROHIBITED_PATTERN" . 2>/dev/null || true)"
    else
        match_output="$(rg --line-number --with-filename \
            "${RG_GLOB_ARGS[@]}" \
            -e "$PROHIBITED_PATTERN" . 2>/dev/null || true)"
    fi

    if [ -z "$match_output" ]; then
        # Count files scanned (best-effort — rg doesn't expose this directly)
        local file_count
        file_count="$(rg --files "${RG_GLOB_ARGS[@]}" . 2>/dev/null | wc -l | tr -d '[:space:]' || echo "?")"
        ok "branding verification passed (scanned ${file_count} files, 0 leaks)"
        return 0
    fi

    # Print leaks
    echo "$match_output"
    echo ""

    local leak_count unique_files
    leak_count="$(echo "$match_output" | wc -l | tr -d '[:space:]')"
    unique_files="$(echo "$match_output" | cut -d: -f1 | sort -u | wc -l | tr -d '[:space:]')"

    fail "branding verification failed: ${leak_count} leaks in ${unique_files} files"
    echo -e "${YELLOW}Hint: legitimate exceptions go in the allowlist inside scripts/verify-branding.sh${NC}" >&2
    return 1
}

# ============================================================================
# Scan upstream/main via diff
# ============================================================================
scan_upstream() {
    local verbose="$1"

    info "Scanning diff HEAD..upstream/main for prohibited tokens (informational)"

    # Check upstream remote exists
    if ! git remote get-url upstream &>/dev/null; then
        die "Remote 'upstream' not found. Run: git remote add upstream https://github.com/Gentleman-Programming/gentle-ai"
    fi

    # Get the diff and grep it
    local diff_output
    diff_output="$(git diff HEAD..upstream/main 2>/dev/null || true)"

    if [ -z "$diff_output" ]; then
        ok "No diff between HEAD and upstream/main (or upstream/main not fetched yet)"
        return 0
    fi

    # Filter only added lines (+) that contain prohibited tokens
    local match_output
    match_output="$(echo "$diff_output" | grep '^+' | grep -v '^+++' \
        | grep -E "$PROHIBITED_PATTERN" || true)"

    if [ -z "$match_output" ]; then
        ok "No prohibited tokens found in upstream/main diff"
        return 0
    fi

    warn "Upstream diff introduces the following prohibited tokens:"
    echo "$match_output"
    echo ""

    local leak_count
    leak_count="$(echo "$match_output" | wc -l | tr -d '[:space:]')"
    warn "Found ${leak_count} line(s) with prohibited tokens in upstream diff (these will need rewriting after merge)"

    # --upstream is informational: always exit 0 (caller decides)
    return 0
}

# ============================================================================
# Main
# ============================================================================
main() {
    setup_colors

    local mode="working-tree"
    local verbose="0"

    while [ $# -gt 0 ]; do
        case "$1" in
            --upstream) mode="upstream"; shift ;;
            --verbose)  verbose="1"; shift ;;
            -h|--help)
                echo "Usage: $0 [--upstream] [--verbose]"
                echo "  (no flags)   Scan working tree for prohibited branding tokens"
                echo "  --upstream   Scan diff HEAD..upstream/main (informational)"
                echo "  --verbose    Show every match with file:line:content"
                exit 0
                ;;
            *) die "Unknown option: $1. Use --help for usage." ;;
        esac
    done

    # Pre-flight: rg must be available
    if ! command -v rg &>/dev/null; then
        die "ripgrep (rg) is required but not found. Install it: brew install ripgrep"
    fi

    if [ "$mode" = "upstream" ]; then
        scan_upstream "$verbose"
    else
        scan_working_tree "$verbose"
    fi
}

main "$@"
