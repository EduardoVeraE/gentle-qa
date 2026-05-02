<!-- Skill: qa-visual-regression · Template: visual-test-plan -->
<!-- Placeholders: {{project}}, {{release}}, {{owner}}, {{date}}, {{approvers}}, {{tool_choice}}, {{tool_rationale}}, {{coverage_targets}}, {{browsers}}, {{viewports}}, {{themes}}, {{locales}}, {{ci_environment}}, {{quota_or_runner}}, {{baseline_branch}}, {{schedule_start}}, {{schedule_end}}, {{report_cadence}}, {{stakeholders}}, {{exit_criteria}} -->

# Visual Regression Test Plan

> Fill in remaining fields; placeholders pre-populate when generated from context.

## 1. Document Control

| Field      | Value             |
| ---------- | ----------------- |
| Project    | {{project}}       <!-- e.g., Storefront web app -->
| Release    | {{release}}       <!-- e.g., R1 / 2026.05 -->
| Owner      | {{owner}}         |
| Date       | {{date}}          |
| Status     | Draft             |
| Approvers  | {{approvers}}     |

---

## 2. Tool Choice

**Selected tool**: {{tool_choice}}  <!-- e.g., Percy / Chromatic / Playwright toHaveScreenshot -->

**Rationale**: {{tool_rationale}}  <!-- e.g., Storybook-first design system, TurboSnap saves cost; OR no SaaS budget, Playwright already in use -->

Alternatives considered (and why rejected):

- ...
- ...

---

## 3. Coverage Targets

{{coverage_targets}}  <!-- e.g., 20 critical pages OR 200 Storybook stories OR 50 component states -->

| Tier        | Pages / Stories / Components | Snapshot count (approx) |
| ----------- | ---------------------------- | ------------------------ |
| Tier 1 — must pass |                       |                          |
| Tier 2 — should pass |                     |                          |
| Tier 3 — nice to have |                    |                          |

---

## 4. Snapshot Matrix

| Axis        | Values                |
| ----------- | --------------------- |
| Browsers    | {{browsers}}          <!-- e.g., Chromium (mandatory), WebKit -->
| Viewports   | {{viewports}}         <!-- e.g., 375 / 768 / 1280 -->
| Themes      | {{themes}}            <!-- e.g., light, dark -->
| Locales     | {{locales}}           <!-- e.g., en, ar (RTL) -->

Total snapshots per build: `pages × browsers × viewports × themes × locales`.

---

## 5. CI / Runner

- CI environment: {{ci_environment}}  <!-- e.g., GitHub Actions; Docker mcr.microsoft.com/playwright:v1.49.0-jammy -->
- Quota or runner constraints: {{quota_or_runner}}  <!-- e.g., 5,000 snapshots/month free tier -->
- Baseline branch: {{baseline_branch}}  <!-- e.g., main -->

## 6. Determinism Controls

- Animations disabled globally
- Caret hidden on inputs
- Fonts: `document.fonts.ready` before capture
- Dynamic regions masked: timestamps, user names, ads, A/B variants
- Test data: deterministic fixtures via API mocking

## 7. Schedule

- Start: {{schedule_start}}
- End: {{schedule_end}}
- Reporting cadence: {{report_cadence}}  <!-- e.g., per PR + weekly summary -->

## 8. Stakeholders

{{stakeholders}}  <!-- e.g., Design system team, frontend leads, QA -->

## 9. Exit Criteria

{{exit_criteria}}  <!-- e.g., zero unapproved diffs on main, <2% PR flake rate over 2 weeks -->

## 10. Out of Scope

- Functional E2E correctness (use playwright-e2e-testing / selenium-e2e-testing)
- Accessibility audits (use a11y-playwright-testing / a11y-selenium-testing)
- Performance / load (use k6-load-test)
- Mobile native UI (out of scope; web only)
