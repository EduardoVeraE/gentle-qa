---
name: qa-visual-regression
description: Visual regression testing toolkit covering Percy (BrowserStack), Chromatic (Storybook), and Playwright's built-in screenshot diffing (`toHaveScreenshot`). Use when asked to set up pixel/perceptual diff baselines, snapshot UI components or full pages, gate PRs on visual changes, manage baseline updates across branches, configure browser/viewport matrices for visual snapshots, mask dynamic regions, or tune diff thresholds and anti-aliasing. Trigger keywords - visual regression, visual testing, pixel diff, screenshot diff, snapshot test (UI), Percy, Chromatic, toHaveScreenshot, baseline image, visual baseline, visual approval, Storybook visual, Argos, Applitools (alternative). NOT for functional E2E testing — use `playwright-e2e-testing` or `selenium-e2e-testing`. NOT for accessibility/contrast audits — use `a11y-playwright-testing` or `a11y-selenium-testing`. NOT for performance/load — use `k6-load-test`. NOT for live browser inspection — use `playwright-mcp-inspect`. NOT for regression suite tiering strategy — use `playwright-regression-strategy`.
license: MIT
metadata:
  author: gentleman-programming
  version: "1.0"
---

# Visual Regression Testing Toolkit

Pixel and perceptual visual regression for web UIs across three first-class tools: **Percy**, **Chromatic**, and **Playwright's built-in `toHaveScreenshot`**. Covers baseline workflow, CI gating, browser/viewport matrices, dynamic-content masking, threshold tuning, and review/approval flows.

**Core principle**: Visual regression catches the bugs functional assertions cannot — CSS regressions, layout shifts, font swaps, z-index breakage, dark-mode bleed, RTL flips. But snapshot tests are the highest-flake test type if you do not control rendering determinism. **A flaky visual suite trains teams to ignore real regressions** — invest in determinism BEFORE coverage.

## When to Use This Skill

- Setting up **visual baselines** for a component library, marketing site, or full app
- Configuring **PR gating** on visual approval (block merge until reviewed)
- Choosing between **Percy / Chromatic / Playwright built-in** for a given engagement
- Defining a **browser × viewport × theme matrix** for visual snapshots
- **Masking dynamic regions** (timestamps, ads, animations, video) to eliminate noise
- Tuning **diff thresholds** (pixel ratio, anti-aliasing, color delta) to balance signal vs flake
- Establishing a **baseline update workflow** across feature branches and main
- Integrating visual testing into **CI** (GitHub Actions, GitLab, CircleCI) with quota awareness
- Triaging **visual diffs**: real regression vs rendering noise vs intentional design change
- Migrating from **Applitools / BackstopJS** to a supported tool

## ISTQB Layer

Layer 4 — **Non-functional testing by type → Visual / UI consistency**.

This skill complements (does not replace) other layers:

| Layer | Coverage |
| ----- | -------- |
| 1. Foundation | `qa-manual-istqb` |
| 2. Strategy | `qa-manual-istqb`, `playwright-regression-strategy` |
| 3. Functional by level | `api-testing`, `playwright-e2e-testing`, `selenium-e2e-testing`, `qa-mobile-testing` |
| 4. Non-functional by type | `qa-owasp-security`, `k6-load-test`, `a11y-playwright-testing`, **`qa-visual-regression`** (this skill) |
| 5. Tooling | `playwright-cli`, `playwright-mcp-inspect` |

## Tool Comparison

Three first-class options. Pick ONE primary tool per project; mixing creates duplicate baselines and approval fatigue.

| Tool                  | Hosting / Cost                    | Best for                                                | Diff engine        | Storybook native | CI cost model      | Reference                           |
| --------------------- | --------------------------------- | ------------------------------------------------------- | ------------------ | ---------------- | ------------------ | ----------------------------------- |
| **Percy**             | SaaS (BrowserStack), paid + free  | Full-page web app snapshots, multi-browser/viewport     | Perceptual + pixel | Via integration  | Per-snapshot quota | `references/percy.md`               |
| **Chromatic**         | SaaS (Chromatic.com), paid + free | Component-library / Storybook-first projects            | Pixel + threshold  | First-class      | Per-snapshot quota | `references/chromatic.md`           |
| **Playwright `toHaveScreenshot`** | Self-hosted, free      | Repos that already use Playwright; no SaaS budget       | Pixel + threshold  | No (E2E only)    | Free (CI minutes)  | `references/playwright-screenshots.md` |

### Headline tradeoffs

- **Percy** — Best multi-browser coverage, generous free tier for OSS, perceptual diff reduces noise. Costs money beyond free tier; vendor lock-in to BrowserStack ecosystem; baselines live off-repo.
- **Chromatic** — Native Storybook integration, **TurboSnap** only re-snapshots changed components (huge cost saver), best DX for design systems. Tied to Storybook; weaker for full-page app flows; quota-driven.
- **Playwright `toHaveScreenshot`** — Free, baselines committed to repo (PR diff is visible), zero vendor lock-in, full control. Higher flake risk (no perceptual diff, OS font rendering differences), needs Docker or pinned runners for determinism, no review UI — diffs reviewed in PR.

### Decision tree

```
Is this a Storybook-driven design system?
├── Yes → Chromatic (TurboSnap + component-first review)
└── No → Is the team OK with SaaS + paying once you exceed free tier?
        ├── Yes → Percy (full-page, multi-browser, perceptual diff)
        └── No → Playwright toHaveScreenshot
                 (commit baselines to repo, run in pinned Docker on CI)
```

## Prerequisites

| Requirement                       | Notes                                                                              |
| --------------------------------- | ---------------------------------------------------------------------------------- |
| Node.js 18+                       | Required for the artifact CLI (`scripts/visual_artifacts.mjs`) and all three tools |
| Playwright (for built-in)         | `@playwright/test` 1.40+ for stable `toHaveScreenshot` API                         |
| Storybook 7+ (for Chromatic)      | Chromatic builds Storybook; without it, Chromatic loses its main advantage         |
| Percy / Chromatic project token   | Required SaaS option; stored as CI secret (`PERCY_TOKEN`, `CHROMATIC_PROJECT_TOKEN`) |
| Pinned CI runner / Docker image   | Mandatory for Playwright built-in to avoid OS-font rendering drift                 |
| Stable test data + feature flags  | Visual tests fail loudly on data drift; freeze fixtures and flag states            |

## Quick Start

Generate visual regression artifacts from templates (CLI implemented in `scripts/visual_artifacts.mjs`):

```bash
# List available templates
node scripts/visual_artifacts.mjs list

# Create a visual test plan
node scripts/visual_artifacts.mjs create visual-test-plan --out specs --project "Storefront"

# Create a snapshot matrix (browsers x viewports x themes)
node scripts/visual_artifacts.mjs create snapshot-matrix --out specs --release "R1"

# Create a baseline-update runbook
node scripts/visual_artifacts.mjs create baseline-runbook --out docs --tool "Playwright"

# Create a visual-diff triage report
node scripts/visual_artifacts.mjs create visual-diff-report --out reports --title "Header redesign"
```

Starter templates per tool live under `templates/`:

- `templates/percy-config.yml` — `.percy.yml` with project-level config
- `templates/chromatic-config.json` — `chromatic.config.json` + GitHub Action snippet
- `templates/playwright-screenshot.spec.ts` — Playwright spec with `toHaveScreenshot` patterns
- `templates/visual-test-plan.md`, `templates/snapshot-matrix.md`, `templates/baseline-runbook.md`, `templates/visual-diff-report.md`

## Workflows

### 1) Choose the right tool

1. Inventory the project: Storybook present? Existing Playwright suite? OSS or commercial? Budget for SaaS?
2. Apply the decision tree above.
3. Document the choice in the visual test plan with an explicit rationale (cost, lock-in, DX, coverage).
4. Resist the urge to run two tools "for coverage" — duplicate baselines double review cost without doubling signal.

### 2) Set up Percy

1. Create a Percy project; copy `PERCY_TOKEN` into CI secrets.
2. Install the Percy SDK for your test runner (`@percy/playwright`, `@percy/cypress`, `@percy/cli` for static sites).
3. Add `percySnapshot(page, "Name")` calls inside existing E2E tests OR run Percy against a static build.
4. Define snapshot widths in `.percy.yml` (mobile 375, tablet 768, desktop 1280, wide 1920).
5. Wire CI: `npx percy exec -- npm test`. First run establishes baselines automatically.
6. Configure the GitHub integration so PR checks block merge until visual review is approved.

See `references/percy.md`.

### 3) Set up Chromatic

1. Create a Chromatic project linked to your repo; copy `CHROMATIC_PROJECT_TOKEN`.
2. Ensure Storybook builds clean (`npm run build-storybook`).
3. Run `npx chromatic --project-token=$CHROMATIC_PROJECT_TOKEN` once locally to seed baselines on `main`.
4. Add the GitHub Action so every PR triggers a Chromatic build.
5. Enable **TurboSnap** (`--only-changed`) — only re-snapshots stories whose dependencies changed, cutting cost 70-90%.
6. Configure UI Review and UI Tests gates; require both to pass before merge.

See `references/chromatic.md`.

### 4) Set up Playwright `toHaveScreenshot`

1. Pin browser version in `playwright.config.ts` AND run CI inside `mcr.microsoft.com/playwright:v1.x.x-jammy` Docker image — OS font rendering varies between Ubuntu releases.
2. Disable animations globally: `page.addStyleTag({ content: '*{animation:none!important;transition:none!important}' })` in a fixture.
3. Mask dynamic regions: `await expect(page).toHaveScreenshot({ mask: [page.getByTestId('timestamp')] })`.
4. Set tolerant thresholds initially: `maxDiffPixelRatio: 0.01`, `threshold: 0.2` — tighten once stable.
5. Commit baselines under `tests/__snapshots__/`; review them in PR diffs visually.
6. Update with `npx playwright test --update-snapshots` ONLY on a dedicated branch; never auto-update in CI.

See `references/playwright-screenshots.md`.

### 5) Manage baselines across branches

Universal rule: **`main` baselines are the source of truth**. Feature branches diff against `main`'s baselines.

| Step               | Percy                                | Chromatic                                  | Playwright built-in                          |
| ------------------ | ------------------------------------ | ------------------------------------------ | -------------------------------------------- |
| Establish baseline | First build on `main`                | First build on `main` (auto-accept)        | `git checkout main && --update-snapshots` + commit |
| Approve PR change  | Review UI in Percy                   | Review UI in Chromatic                     | Re-run with `--update-snapshots` on PR branch + commit |
| Reject change      | Reject in Percy UI                   | Deny in Chromatic UI                       | Do not regenerate; fix the code              |
| Cross-branch sync  | Auto (cloud)                         | Auto (cloud)                               | Rebase on main; regenerate if main changed   |

### 6) Define the snapshot matrix

A matrix of `browsers × viewports × themes × locales` explodes quickly. Use the 80/20 rule.

1. **Browsers**: Chromium (mandatory). Add WebKit if iOS Safari is a top-3 traffic source. Firefox only if explicit requirement.
2. **Viewports**: 375 (mobile), 768 (tablet), 1280 (desktop). Skip 1920 unless wide layouts differ.
3. **Themes**: Light + Dark if dark mode shipped. Skip if not.
4. **Locales**: Default + 1 RTL locale (`ar` or `he`) if RTL is supported. Skip otherwise.
5. **States**: empty, loading, error, populated — at minimum one snapshot per state per critical component.

Use: `templates/snapshot-matrix.md`.

### 7) Triage visual diffs

For every diff, classify in this order:

1. **Real regression** — bug. File defect, link build URL, fix code, do not update baseline.
2. **Intentional design change** — accept new baseline. Document in PR description.
3. **Rendering noise** — flake. Investigate root cause: animation, font load, dynamic content, OS rendering. Fix determinism BEFORE updating baseline.

Never auto-accept all diffs. Never update baselines on a "make it pass" branch.

Use: `templates/visual-diff-report.md`.

## Inputs to Collect

- **App URL or build artifact**: stable staging URL or static build directory.
- **Critical pages / components**: top 10-20 by user traffic; do not snapshot every page.
- **Browser + viewport coverage targets**: derived from analytics, not aspiration.
- **Theme / locale / state coverage**: dark mode, RTL, empty, loading, error.
- **Dynamic regions**: timestamps, user names, ads, A/B test variants, animations — must be masked or stubbed.
- **CI runner pinning**: Docker image tag (Playwright) or known SaaS environment (Percy/Chromatic).
- **Feature flag fixtures**: every snapshot test must run with deterministic flag state.

## Outputs

| Artifact                                        | When produced                          | Template                              |
| ----------------------------------------------- | -------------------------------------- | ------------------------------------- |
| Visual test plan                                | Per project / per major release        | `templates/visual-test-plan.md`       |
| Snapshot matrix                                 | Per release; reviewed quarterly        | `templates/snapshot-matrix.md`        |
| Baseline runbook                                | Once per repo; per tool                | `templates/baseline-runbook.md`       |
| Visual diff triage report                       | Per PR with diffs                      | `templates/visual-diff-report.md`     |
| Tool config files                               | Once per repo                          | `templates/percy-config.yml`, `templates/chromatic-config.json`, `templates/playwright-screenshot.spec.ts` |

## Cross-reference

Visual regression was previously buried inside `playwright-e2e-testing/SKILL.md`. That skill now points here for any visual concern. If you are writing functional E2E tests with Playwright and only need a single screenshot assertion, you can stay there; if you need a full visual suite, use this skill.

## Exclusions

This skill is deliberately scoped. Do NOT use it for:

- **Functional E2E testing** (UI flows, form validation, navigation correctness) — use `playwright-e2e-testing` or `selenium-e2e-testing`.
- **Accessibility audits** (contrast ratios, ARIA, keyboard nav) — use `a11y-playwright-testing` or `a11y-selenium-testing`. Visual diff catches a CSS color regression but cannot tell you if it broke WCAG.
- **Performance / load testing** — use `k6-load-test`. Visual tests are not perf tests.
- **Live browser inspection / debugging** — use `playwright-mcp-inspect`.
- **Regression suite tiering** (smoke / sanity / full) — use `playwright-regression-strategy`. This skill defines the visual layer; tier strategy is owned elsewhere.
- **Mobile native visual testing** (iOS / Android native UI) — escalate; Percy/Chromatic/Playwright cover web only. Appium has limited image-comparison support — out of scope for this skill.
- **Email / PDF rendering** — out of scope; use dedicated tools (Litmus, percy-pdf, etc.).

If a request mixes concerns (e.g., "snapshot the dashboard and check it loads under 200ms"), split it: visual layer here, perf with `k6-load-test`.
