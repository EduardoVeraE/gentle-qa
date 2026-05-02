# Changelog

All notable changes to **Gentle-QA** are documented in this file.

The format is loosely based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project follows semantic versioning where applicable.

> Gentle-QA is a fork of [Gentleman-Programming/gentle-ai](https://github.com/Gentleman-Programming/gentle-ai) (originally **gentle-ai**) tailored for Quality Engineering / SDET workflows. This changelog tracks both fork-specific work and the upstream features incorporated through periodic syncs.

---

## [1.25.1] — 2026-05-02

### Fixed

- **Wire qa-* skills end-to-end** — the 4 QA skills shipped in `v1.25.0`
  (`qa-owasp-security`, `qa-mobile-testing`, `qa-visual-regression`,
  `qa-contract-pact`) were embedded as assets but unreachable to end users
  because they were not registered in `internal/model/types.go` and not
  included in any preset. Now wired into `model`, `catalog`, and the
  `full-gentleman` and `qe-sdet` presets so the installer actually
  delivers them.
- **`e2e_test.sh` SKILL.md count** — `full-gentleman` now expects 25 files
  (11 SDD + 5 foundation + 4 QE + 4 QA + `_shared/SKILL.md`); the previous
  17 was stale since 04-25 and caused CI to fail on every push to `main`
  since 04-29.

### Convention established

- **Asset → distribution wiring is mandatory.** Adding a skill, preset, or
  configuration to `internal/assets/` is NOT enough. Every addition must
  also: (1) register a constant in `internal/model/types.go`, (2) appear
  in the appropriate preset arrays in
  `internal/components/skills/presets.go`, (3) update
  `testdata/golden/skills-presets.json`, (4) update `e2e_test.sh`
  expected counts, (5) update catalog entries. Asset-only changes ship in
  the binary but stay orphaned to end users — that is a bug, not a
  feature.

---

## [1.25.0] — 2026-05-02

This release combines two unreleased streams of work since `v1.24.3`:

1. **QA skill expansion** (2026-05-02) — 4 new skills + 1 enriched skill under the new ISTQB taxonomy, with deterministic validation harnesses.
2. **Upstream sync** (2026-04-29) — 172 commits from `Gentleman-Programming/gentle-ai` integrated without losing Gentle-QA branding or QE skills.

### Added

#### New QA skills (5-layer ISTQB taxonomy)

- **`qa-owasp-security`** (Layer 4 — Non-functional / Security) — OWASP Web Top 10 + API Top 10 + MASVS coverage, ZAP active scans, dependency scanning (trivy), XSS/SQLi attack scripts. Ships a deterministic validation harness (docker-compose + DVWA fixtures + `lib.sh` with `assert_exit_code` to avoid the pipe-to-tail trap) — 4/4 PASS.
- **`qa-mobile-testing`** (cross-cutting platform skill) — Appium (iOS + Android) and Detox (React Native grey-box), device strategy guidance (sims/emulators vs real devices vs cloud farms), gesture/wait patterns, flake mitigation. Validation harness validates the deterministic surfaces only (CLI behavior, scaffold syntax via `ts.transpileModule` parser-only) — 8/8 PASS.
- **`qa-visual-regression`** (Layer 4 — Non-functional / Tooling) — Percy, Chromatic, and Playwright `toHaveScreenshot` covered as separate tools with explicit tradeoffs (paid SaaS vs free in-repo baselines, perceptual vs pixel diff). Includes baseline workflow, masking, threshold tuning, CI gating, and Docker pinning guidance for OS-font determinism. CLI self-test 9/9 PASS.
- **`qa-contract-pact`** (Layer 3 — Functional / Integration) — Consumer-driven contract testing with PACT-JS v12+ and PACT-JVM v4.6+, Pact Broker (self-hosted Docker + Postgres) and PactFlow, `can-i-deploy` gates, `record-deployment`, broker webhooks, HTTP and async (Kafka/RabbitMQ/SNS/SQS) message contracts. CLI self-test 4/4 PASS.

#### From upstream `gentle-ai`

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

### Enriched

- **`api-testing`** — promoted to ISTQB-aligned skill: mandatory headers, OpenAPI-first workflow, contract-testing bridge to `qa-contract-pact`, scope explicitly bounded (NOT for security → `qa-owasp-security`, NOT for performance → `k6-load-test`). Playwright TypeScript + REST Assured (Java 21+) examples.

### Changed

- **5-layer ISTQB taxonomy** adopted for skill organization: Foundation → Strategy → Functional-by-level → Non-functional-by-type → Tooling.
- **Disjoint triggers + exclusion clauses** required in every skill description (`NOT for X — use Y`) — the orchestrator matches by trigger text; overlap caused both-fire or neither-fire problems before this convention.
- **CLI artifact-generation pattern** standardized across `security_artifacts.mjs`, `mobile_artifacts.mjs`, `api_artifacts.mjs`, `visual_artifacts.mjs`, `contract_artifacts.mjs` — byte-near-identical except DESCRIPTIONS map + name strings + template count.
- `playwright-e2e-testing` → points to `qa-visual-regression` (no longer buries visual regression inside the E2E skill).
- `karate-dsl` → clarifies strict-match scope and points to `qa-contract-pact` for consumer-driven testing.
- `api-testing/references/headers-and-contracts.md` → linked to `qa-contract-pact`.
- Engram MCP goldens regenerated for stable absolute paths (`/opt/homebrew/bin/engram`).
- `skills-presets.json` golden refreshed to include Gentle-QA QE skills under `full-gentleman`.
- Embedded skill directory count: 17 → 30 (10 SDD + judgment-day + 5 foundation + `_shared` + 13 QA/SDET).

### Fixed (from upstream)

- `engram` — migrate stale Homebrew paths; retry latest release lookup anonymously
- `sdd` — replace wildcard `task` permissions with explicit allowlists; complete Claude native sub-agent integration
- `claude` — align SDD agent tools with artifact backends
- `opencode` — orchestrator prompt drift; community plugin upgrade
- `persona` — keep English sessions in English; add response length discipline
- `e2e` — Claude SDD tool assertions, pacman keyring init in Arch image, qwen idempotency

### Tech debt opened (P3)

- iOS Simulator-based Detox/Appium runner, Android Emulator-based runner, cloud device-farm smoke runner — tracked under beads as dependents of the closed mobile harness issue.

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
