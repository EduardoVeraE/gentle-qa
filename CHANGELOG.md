# Changelog

All notable changes to **Gentle-QA** are documented in this file.

The format is loosely based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project follows semantic versioning where applicable.

> Gentle-QA is a fork of [Gentleman-Programming/gentle-ai](https://github.com/Gentleman-Programming/gentle-ai) (originally **gentle-ai**) tailored for Quality Engineering / SDET workflows. This changelog tracks both fork-specific work and the upstream features incorporated through periodic syncs.

---

## [Unreleased] — 2026-04-29

### Upstream sync — 171 commits from `Gentleman-Programming/gentle-ai`

Major refresh from upstream `gentle-ai`, integrated without losing the Gentle-QA branding, QE skill set, or SDET tooling.

#### Added (from upstream)

- **OpenCode Community Plugins** — install curated community plugins for the OpenCode TUI alongside the SDD foundation:
  - `sub-agent-statusline` — surfaces the active sub-agent in the statusline
  - `sdd-engram-plugin` — manage SDD profiles and browse Engram memories from inside OpenCode, with runtime profile activation (no restart)
  - New TUI screen for plugin selection and a `OpenCodePlugins` field on `Selection`
- **Per-phase Claude sub-agent model assignments** — every SDD phase (`sdd-explore`, `sdd-propose`, `sdd-spec`, `sdd-design`, `sdd-tasks`, `sdd-apply`, `sdd-verify`, `sdd-archive`) can be pinned to `opus`, `sonnet`, or `haiku` via Claude's native `~/.claude/agents/` install
- **Native Claude sub-agents** — 8 SDD sub-agent files with per-phase model frontmatter; tool scope tightened per phase
- **Claude SDD slash commands** — `/sdd-*` now installs as native slash commands
- **OpenCode external profile sync strategy** (`external-single-active`) — compatibility mode for community profile managers
- **Antigravity agent integration improvements** — settings bootstrap moves into the engram flow; backups exclude `antigravity_tmp`
- **Qwen Code idempotency hardening** — engram skips unsupported qwen setup; e2e cleanup extended
- **Kimi agent** — `uv` preflight aligned with install flow

#### Fixed (from upstream)

- `engram` — migrate stale Homebrew paths; retry latest release lookup anonymously
- `sdd` — replace wildcard `task` permissions with explicit allowlists; complete Claude native sub-agent integration
- `claude` — align SDD agent tools with artifact backends
- `opencode` — orchestrator prompt drift; community plugin upgrade
- `persona` — keep English sessions in English; add response length discipline
- `e2e` — Claude SDD tool assertions, pacman keyring init in Arch image, qwen idempotency

#### Changed

- Engram MCP goldens regenerated for stable absolute paths (`/opt/homebrew/bin/engram`)
- `skills-presets.json` golden refreshed to include Gentle-QA QE skills under `full-gentleman`
- Embedded skill directory count: 17 → 30 (10 SDD + judgment-day + 5 foundation + `_shared` + 13 QA/SDET)

### Preserved (Gentle-QA fork)

Everything that makes this a QA-focused fork is intact:

- **Module path**: `github.com/EduardoVeraE/Gentle-QA`
- **Binary**: `gentle-qa`
- **Environment variables**: `GENTLE_QA_*` (e.g. `GENTLE_QA_ENGRAM_SETUP_MODE`, `GENTLE_QA_ENGRAM_SETUP_STRICT`)
- **QE skills** still installed as part of the catalog:
  - `playwright-bdd`
  - `playwright-cli`
  - `k6-load-test`
  - `karate-dsl`
- **QE presets**: `qe-front`, `qe-perf`, `qe-api`, `qe-sdet`
- **PersonaSDET** persona
- **Cosmonaut PNG renderer** (`golang.org/x/image v0.39.0`)
- **Beads** workspace integration (issue tracking under `.beads/`)
- All third-party `Gentleman-Programming/homebrew-tap` references for `engram` and `gga` binaries (these are upstream tooling, unrelated to the gentle-ai fork itself, and must remain pointing at the original tap)

### Notes

- Backup branch preserved at `backup/pre-upstream-merge-2026-04-29` before the merge
- Conflict resolution strategy: `ours-canonical` (build manifests, branding-critical files), `additive-merge` (model/selection/skills code that gained both sides), `theirs-bulk` (most upstream-evolved files) — followed by automated rebrand sweeps for module paths, env vars, and casing variants
- Build (`go build ./...`) is clean
- Tests (`go test ./...`) are all green

---

## Fork history (pre-sync)

Earlier work that turned upstream `gentle-ai` into the QE-focused Gentle-QA fork. Highlights only — see git history for the full record.

### Branding & module rename

- Module path renamed from `github.com/Gentleman-Programming/gentle-ai` to `github.com/EduardoVeraE/Gentle-QA`
- CLI binary renamed from `gentle-ai` to `gentle-qa`
- Env var prefix renamed from `GENTLE_AI_*` to `GENTLE_QA_*`
- All TUI strings, banners, docs, and goldens updated to **Gentle-QA**
- Homebrew tap and Scoop bucket switched to `EduardoVeraE/tap` and `EduardoVeraE/scoop-bucket`
- CI release workflow updated (with `workflow_dispatch` for manual releases)

### QE / SDET extensions

- Added 4 QA-focused skills: `playwright-bdd`, `playwright-cli`, `k6-load-test`, `karate-dsl`
- Added 4 presets aimed at QA workflows: `qe-front`, `qe-perf`, `qe-api`, `qe-sdet`
- Introduced the `PersonaSDET` persona
- Catalog entries for QE skills under `qe-e2e`, `qe-performance`, and `qe-api` categories

### Tooling & infrastructure

- Beads workspace bootstrapped under `.beads/` for issue tracking
- `.gitignore` extended for `.dolt/`, `*.db`, `.beads-credential-key`, `.engram/`, `.claude/`, `.beads`, `.windsurf`, etc.
- Workflow alignment with the personal `~/.claude/CLAUDE.md` rules

---

## Upstream credit

Gentle-QA stands on the work done in [Gentleman-Programming/gentle-ai](https://github.com/Gentleman-Programming/gentle-ai). Maintainer: **Alan Buscaglia** ([@Alan-TheGentleman](https://github.com/Alan-TheGentleman)). See [CONTRIBUTORS.md](CONTRIBUTORS.md) for the full list of contributors whose work flows into this fork.
