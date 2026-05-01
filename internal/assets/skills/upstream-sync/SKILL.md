---
name: upstream-sync
description: >
  Upstream synchronization workflow for merging Gentleman-Programming/gentle-ai into the Gentle-QA fork while preserving rebranding.
  Use when synchronizing upstream gentle-ai changes into the fork, pulling in upstream commits, or resolving upstream drift.
  NOT for routine git merges or branch syncs — use only for upstream/main → main fork synchronization.
---

# Upstream Sync — gentle-ai → Gentle-QA Fork

Workflow for pulling upstream `Gentleman-Programming/gentle-ai` commits into this fork without leaking prohibited branding tokens into the codebase.

## When to Use This Skill

- Pulling new commits from `upstream/main` into the fork's `main`
- Checking how far behind the fork has drifted from upstream
- Recovering from a partial or failed upstream merge
- Auditing the codebase for accidental branding leaks after any upstream integration

## Prerequisites

| Requirement | Notes |
|-------------|-------|
| Clean working tree | `git status` must report nothing to commit |
| Remotes configured | `origin` → `git@github.com:EduardoVeraE/gentle-qa.git`; `upstream` → `https://github.com/Gentleman-Programming/gentle-ai.git` |
| Go toolchain | `go build ./...` and `go test ./...` must be runnable |
| `scripts/sync-upstream.sh` | Executor — performs merge + token rewrites |
| `scripts/verify-branding.sh` | Integrity checker — exits non-zero on any prohibited token |

## Workflow Overview

1. **Pre-flight** — verify the working tree is clean, then fetch upstream:
   ```bash
   git status            # must be clean
   git fetch upstream
   git rev-list HEAD..upstream/main --count   # commits behind
   ```

2. **Pre-merge scan** — preview which upstream files contain branding tokens (informational only, not blocking):
   ```bash
   scripts/verify-branding.sh --upstream
   ```

3. **Merge** — run the sync script; it performs the merge, applies known token rewrites, and runs post-merge verification internally:
   ```bash
   scripts/sync-upstream.sh
   ```

4. **Build and test** — both must pass before committing:
   ```bash
   go build ./...
   go test ./...
   ```

5. **Verify branding integrity** — must exit 0; any leak aborts and leaves the merge open for manual resolution:
   ```bash
   scripts/verify-branding.sh
   ```

6. **Commit** — use this exact conventional commit format, filling in the actual count:
   ```
   chore(sync): merge upstream gentle-ai (<N> commits) into gentle-qa fork
   ```

7. **Do NOT push automatically** — Eduardo controls when to push to `origin`.

## Branding Allowlist

These files contain prohibited-looking tokens that are legitimate and MUST NOT be flagged:

| File | Token(s) | Why it is allowed |
|------|----------|-------------------|
| `CHANGELOG.md` | Any upstream name references | Historical changelog entries document what was merged from upstream — removing them would destroy the audit trail |
| `CONTRIBUTORS.md` | Any upstream name references | Credit to upstream contributors must be preserved intact; removing attribution would be both incorrect and disrespectful |
| `internal/components/filemerge/section_test.go` | `gentleAiMarkerSection` | This is a test fixture variable name for a persona-marker section; it is a legacy compatibility identifier for the merge-marker format, not a live reference to the upstream product |

**Judgment rule for new cases:** if a token appears inside a test fixture, a historical document (changelog, attribution), or an archived spec, it is almost certainly safe. If it appears in production Go source, a help string, or a UI label, it is a leak and must be rewritten.

## Expected Conflict Pattern: Branding-Rewritten Files

### Why these conflicts happen

Every upstream commit that touches a line we previously rewrote (import paths, renamed function identifiers, replaced string literals) will cause a merge conflict on that line. Git cannot auto-merge because both sides changed the same line — upstream changed the logic, we changed the identifier. This is **structural to the rebrand-fork model**, not a script bug. Exit 3 from `sync-upstream.sh` is correct behavior. Expect this pattern on every upstream commit that touches a file containing our rebranded identifiers (e.g. anything using `gentleQaHint`, `EduardoVeraE/Gentle-QA`, or `GENTLE_QA`).

### How to recognize the pattern

After `exit 3`, run:
```bash
git diff --name-only --diff-filter=U
```
If the conflicted files are files you know were rebranded (they contain our markers — `gentleQaHint`, `github.com/EduardoVeraE/Gentle-QA`, etc.) you are in this pattern. The conflict is on the renamed line; upstream's new logic is on `>>>>>>> upstream/main` side.

### Manual resolution checklist

1. Open the conflicted file in full. Read both `<<<<<<< HEAD` and `>>>>>>> upstream/main` sides to understand what upstream changed.

2. **Cross-file dependency check** — identify every file that imports or calls into the conflicted file:
   ```bash
   rg -l '<conflicted_basename without extension>'   # e.g. rg -l 'instructions'
   ```
   Read those callers. Confirm which symbols they reference. Any NEW symbol introduced by upstream on the `>>>>>>> upstream/main` side that is called from another file MUST be preserved in your resolution — dropping it breaks the build silently after the merge commit.

3. Resolve the conflict: keep our renamed identifiers and module paths; fold upstream's new logic under our names; rewrite any branding tokens in new code that upstream introduced.

4. **Re-apply branding rewrite before continuing** — the script aborts before its rewrite step when there are conflicts, so non-conflict regions of the same file can contain upstream tokens that were never rewritten. After saving the conflict resolution, scan for leaks:
   ```bash
   rg -i 'gentle-ai|gentleAi|GENTLE_AI' <resolved-file>
   ```
   Rewrite any hits manually (map to `gentle-qa` / `gentleQa` / `GENTLE_QA`).

5. Stage and complete the merge:
   ```bash
   git add <file>
   git commit --no-edit
   ```

6. Run `scripts/verify-branding.sh` standalone — must exit 0.

7. Run `go build ./...` and `go test ./...` — both must pass.

**First-run reference:** this pattern was first encountered with upstream commit `a164adc` (`fix(opencode): handle plugin update layers`, 11 files in `internal/update/**`), resolved manually into commit `f15b945`.

---

## Failure Modes & Recovery

**`verify-branding.sh` exits non-zero after merge**

The merge is left in place (not rolled back). Inspect the reported files:
```bash
scripts/verify-branding.sh   # read the output — it lists file:line pairs
```
Manually rewrite the offending tokens, stage the changes, then re-run verification. Do not commit until the script exits 0.

**Merge conflict in a branding-sensitive file**

See **Expected Conflict Pattern: Branding-Rewritten Files** above.

**`go build ./...` or `go test ./...` fails after merge**

This is an upstream API change or Go module incompatibility. Do not commit. Investigate the failure, apply the necessary fix (update imports, adjust call sites), verify, then commit. The sync commit should include both the upstream merge and the compatibility fix — document this in the commit body.

**Remotes not configured**

```bash
git remote add upstream https://github.com/Gentleman-Programming/gentle-ai.git
git remote -v   # confirm both origin and upstream are present
```

## Post-sync Checklist

- [ ] `go build ./...` passes
- [ ] `go test ./...` passes
- [ ] `scripts/verify-branding.sh` exits 0
- [ ] Commit message follows the exact format: `chore(sync): merge upstream gentle-ai (<N> commits) into gentle-qa fork`
- [ ] If conflict was resolved manually, verify cross-file callers of any modified function still build (see **Expected Conflict Pattern** above)
- [ ] Save a memory entry to engram with `topic_key: sync/upstream-history` noting the upstream commit range merged and the date
- [ ] Push to `origin` only when Eduardo explicitly decides to

## Related Files

| File | Role |
|------|------|
| `scripts/sync-upstream.sh` | Executor — runs the merge, applies token rewrites, calls verify internally |
| `scripts/verify-branding.sh` | Integrity check — scans working tree for prohibited tokens; exits non-zero on any leak |
