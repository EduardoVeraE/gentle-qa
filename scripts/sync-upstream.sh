#!/usr/bin/env bash
set -euo pipefail

# ============================================================================
# gentle-qa — Upstream Sync Script
# Orchestrates a safe merge from upstream/main (Gentleman-Programming/gentle-ai)
# into the gentle-qa fork, with automatic token rewriting and branding safety.
#
# Usage:
#   scripts/sync-upstream.sh
#   SYNC_UPSTREAM_YES=1 scripts/sync-upstream.sh   # skip confirmation prompt
#
# Exit codes:
#   0 — sync completed successfully (or already up to date)
#   1 — branding leak detected after auto-rewrite (manual intervention required)
#   2 — pre-flight failure (dirty tree, wrong branch, missing remote, etc.)
#   3 — merge conflict (resolve manually, then re-run verify and commit)
# ============================================================================

# ============================================================================
# Allowlist — paths that must NOT be rewritten during token substitution.
# These match the same set used in verify-branding.sh.
# ============================================================================
REWRITE_SKIP=(
    "CHANGELOG.md"
    "CONTRIBUTORS.md"
    "internal/components/filemerge/section_test.go"
    "scripts/verify-branding.sh"
    "scripts/sync-upstream.sh"
    "internal/assets/skills/upstream-sync/SKILL.md"
)

# ============================================================================
# Color support
# ============================================================================
setup_colors() {
    if [ -t 1 ] && [ "${TERM:-}" != "dumb" ]; then
        RED='\033[0;31m'
        GREEN='\033[0;32m'
        YELLOW='\033[1;33m'
        CYAN='\033[0;36m'
        BOLD='\033[1m'
        DIM='\033[2m'
        NC='\033[0m'
    else
        RED='' GREEN='' YELLOW='' CYAN='' BOLD='' DIM='' NC=''
    fi
}

# ============================================================================
# Logging helpers
# ============================================================================
info()  { echo -e "${DIM}[sync-upstream]${NC} $*"; }
ok()    { echo -e "${GREEN}✓${NC} $*"; }
fail()  { echo -e "${RED}✗${NC} $*" >&2; }
warn()  { echo -e "${YELLOW}[warn]${NC} $*"; }
step()  { echo -e "\n${CYAN}${BOLD}==>${NC}${BOLD} $*${NC}"; }

die() {
    local exit_code="${1}"; shift
    fail "$*"
    exit "$exit_code"
}

# ============================================================================
# Step 1 — Pre-flight checks
# ============================================================================
preflight() {
    step "Pre-flight checks"

    # Working tree must be clean
    if ! git diff --quiet || ! git diff --cached --quiet; then
        die 2 "Working tree is dirty. Commit or stash your changes before syncing."
    fi
    ok "Working tree is clean"

    # Must be on main
    local current_branch
    current_branch="$(git rev-parse --abbrev-ref HEAD)"
    if [ "$current_branch" != "main" ]; then
        die 2 "Must be on 'main' branch (currently on '${current_branch}'). Run: git checkout main"
    fi
    ok "On branch main"

    # upstream remote must exist
    if ! git remote get-url upstream &>/dev/null; then
        die 2 "Remote 'upstream' not found. Run: git remote add upstream https://github.com/Gentleman-Programming/gentle-ai"
    fi

    local upstream_url
    upstream_url="$(git remote get-url upstream)"
    if ! echo "$upstream_url" | grep -qi "gentle-ai"; then
        warn "upstream remote URL does not look like gentle-ai: ${upstream_url}"
        warn "Continuing anyway — verify this is correct."
    fi
    ok "Remote 'upstream' exists: ${upstream_url}"

    # Fetch upstream
    info "Fetching upstream..."
    git fetch upstream
    ok "Fetched upstream"

    # Compute how far behind we are
    BEHIND_COUNT="$(git rev-list --count HEAD..upstream/main)"

    if [ "$BEHIND_COUNT" -eq 0 ]; then
        ok "Already up to date with upstream/main"
        exit 0
    fi

    ok "Behind upstream/main by ${BEHIND_COUNT} commit(s)"
}

# ============================================================================
# Step 2 — Preview and confirm
# ============================================================================
preview_and_confirm() {
    step "Preview: ${BEHIND_COUNT} commit(s) to merge from upstream/main"

    git log --oneline HEAD..upstream/main
    echo ""

    # Run branding scan on upstream diff — informational only
    info "Checking upstream diff for prohibited tokens (informational)..."
    if ! scripts/verify-branding.sh --upstream; then
        warn "verify-branding.sh --upstream returned non-zero (unexpected for informational mode)"
    fi
    echo ""

    # Skip prompt if SYNC_UPSTREAM_YES is set
    if [ "${SYNC_UPSTREAM_YES:-}" = "1" ]; then
        info "SYNC_UPSTREAM_YES=1 — skipping confirmation prompt"
        return
    fi

    printf "Proceed with merge? [y/N] "
    read -r answer </dev/tty
    case "$answer" in
        y|Y) : ;;
        *) die 2 "Merge aborted by user." ;;
    esac
}

# ============================================================================
# Step 3 — Merge
# ============================================================================
do_merge() {
    step "Merging upstream/main"

    local merge_msg="chore(sync): merge upstream gentle-ai (${BEHIND_COUNT} commits) into gentle-qa fork"

    if ! git merge --no-ff upstream/main -m "$merge_msg"; then
        # Check for conflicts
        if git diff --name-only --diff-filter=U | grep -q .; then
            local conflicted_files
            conflicted_files="$(git diff --name-only --diff-filter=U)"
            echo "" >&2
            fail "Merge conflicts detected. Conflicted files:"
            echo "$conflicted_files" | while IFS= read -r f; do
                echo -e "  ${YELLOW}${f}${NC}" >&2
            done
            echo "" >&2
            warn "This is expected when upstream modifies lines that our rebrand previously rewrote."
            warn "See internal/assets/skills/upstream-sync/SKILL.md#expected-conflict-pattern-branding-rewritten-files for resolution steps."
            echo "" >&2
            warn "After resolving manually:"
            warn "  1. Scan for branding leaks in non-conflict regions: rg -i 'gentle-ai|gentleAi|GENTLE_AI' <file>"
            warn "  2. Run: scripts/verify-branding.sh  (must exit 0)"
            warn "  3. Run: go build ./... && go test ./..."
            warn "  4. Run: git add <file> && git commit --no-edit"
            warn "Do NOT skip the rewrite — branding leaks can exist outside conflict markers."
            echo "" >&2
            exit 3
        fi
        die 2 "Merge failed for an unexpected reason. Check git output above."
    fi

    ok "Merge succeeded (no conflicts)"
}

# ============================================================================
# Step 4 — Token rewrite
# ============================================================================
is_skipped() {
    local filepath="$1"
    local entry
    for entry in "${REWRITE_SKIP[@]}"; do
        if [ "$filepath" = "$entry" ]; then
            return 0
        fi
    done
    return 1
}

do_rewrite() {
    step "Rewriting prohibited tokens"

    local files_rewritten=0

    # Enumerate tracked files via git ls-files
    local tracked_files
    # Read into array portably (no mapfile)
    local IFS=$'\n'
    local all_files=()
    while IFS= read -r f; do
        all_files+=("$f")
    done < <(git ls-files)
    unset IFS

    for filepath in "${all_files[@]}"; do
        # Skip allowlisted paths
        if is_skipped "$filepath"; then
            continue
        fi

        # Skip non-regular or non-text files
        if [ ! -f "$filepath" ]; then
            continue
        fi

        # Check if file contains any prohibited token before running perl
        if ! grep -qE \
            'github\.com/gentleman-programming/gentle-ai|Gentleman-Programming/gentle-ai|gentleman-programming/gentle-ai|gentle-ai|GENTLE_AI|gentleAI|gentleAi' \
            "$filepath" 2>/dev/null; then
            continue
        fi

        # Apply all substitutions in a single perl invocation.
        # Order matters: most-specific patterns first, then the shorter ones.
        # NOTE: BSD sed (macOS) inconsistencies make perl the safe choice.
        perl -i -pe '
            s{github\.com/gentleman-programming/gentle-ai}{github.com/EduardoVeraE/Gentle-QA}gi;
            s{Gentleman-Programming/gentle-ai}{EduardoVeraE/Gentle-QA}g;
            s{gentleman-programming/gentle-ai}{EduardoVeraE/Gentle-QA}g;
            s{GENTLE_AI}{GENTLE_QA}g;
            s{gentleAI}{gentleQA}g;
            s{gentleAi}{gentleQa}g;
            s{gentle-ai}{gentle-qa}g;
        ' "$filepath"

        files_rewritten=$((files_rewritten + 1))
    done

    REWRITTEN_COUNT="$files_rewritten"

    if [ "$files_rewritten" -gt 0 ]; then
        ok "Rewrote tokens in ${files_rewritten} file(s)"
        # Stage all tracked modifications and amend the merge commit
        git add -u
        git commit --amend --no-edit
        ok "Amended merge commit with rewritten tokens"
    else
        ok "No files needed token rewriting"
    fi
}

# ============================================================================
# Step 5 — Verify branding
# ============================================================================
do_verify() {
    step "Verifying branding"

    if ! scripts/verify-branding.sh; then
        fail "Branding leak detected after auto-rewrite. Manual intervention required."
        fail "The merge commit is in place — inspect the leaks above, fix them,"
        fail "then run: scripts/verify-branding.sh && git commit --amend --no-edit"
        exit 1
    fi

    ok "Branding verification passed"
}

# ============================================================================
# Step 6 — Build + test gate
# ============================================================================
do_build_test() {
    step "Build + test gate"

    info "Running: go build ./..."
    if ! go build ./...; then
        die 2 "Build failed after merge. Fix compilation errors before pushing."
    fi
    ok "Build passed"

    info "Running: go test ./..."
    if ! go test ./...; then
        die 2 "Tests failed after merge. Fix test failures before pushing."
    fi
    ok "All tests passed"
}

# ============================================================================
# Step 7 — Done
# ============================================================================
print_summary() {
    step "Sync complete"

    echo ""
    echo -e "${GREEN}${BOLD}Summary${NC}"
    echo -e "  Commits merged : ${BOLD}${BEHIND_COUNT}${NC}"
    echo -e "  Files rewritten: ${BOLD}${REWRITTEN_COUNT:-0}${NC}"
    echo -e "  Build          : ${GREEN}passing${NC}"
    echo -e "  Tests          : ${GREEN}passing${NC}"
    echo ""
    echo -e "${DIM}Run 'git log --oneline -5' to inspect.${NC}"
    echo -e "${DIM}Push when ready: git push origin main${NC}"
    echo ""
}

# ============================================================================
# Main
# ============================================================================
main() {
    setup_colors

    # Global counters (set by the step functions)
    BEHIND_COUNT=0
    REWRITTEN_COUNT=0

    while [ $# -gt 0 ]; do
        case "$1" in
            -h|--help)
                echo "Usage: $0 [--help]"
                echo "  Merge upstream gentle-ai into gentle-qa fork with branding safety."
                echo "  Set SYNC_UPSTREAM_YES=1 to skip the confirmation prompt."
                exit 0
                ;;
            *) die 2 "Unknown option: $1. Use --help for usage." ;;
        esac
    done

    preflight
    preview_and_confirm
    do_merge
    do_rewrite
    do_verify
    do_build_test
    print_summary
}

main "$@"
