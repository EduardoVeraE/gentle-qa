<!-- Skill: qa-visual-regression · Template: snapshot-matrix -->
<!-- Placeholders: {{project}}, {{release}}, {{owner}}, {{date}}, {{tool}}, {{primary_browser}}, {{secondary_browsers}}, {{viewports}}, {{themes}}, {{locales}}, {{state_coverage}}, {{snapshot_budget}}, {{notes}} -->

# Snapshot Matrix

| Field    | Value          |
| -------- | -------------- |
| Project  | {{project}}    |
| Release  | {{release}}    |
| Owner    | {{owner}}      |
| Date     | {{date}}       |
| Tool     | {{tool}}       <!-- Percy / Chromatic / Playwright -->

## Axes

| Axis | Values | Justification |
| ---- | ------ | ------------- |
| Primary browser | {{primary_browser}} | Default; covers >90% traffic |
| Secondary browsers | {{secondary_browsers}} | Only if analytics justify cost |
| Viewports | {{viewports}} <!-- 375 / 768 / 1280 --> | 80/20 traffic |
| Themes | {{themes}} <!-- light, dark --> | Match shipped themes only |
| Locales | {{locales}} <!-- en + 1 RTL --> | RTL only if supported |
| State coverage | {{state_coverage}} <!-- empty, loading, error, populated --> | Per critical component |

## Snapshot Budget

Estimated snapshots per build: **{{snapshot_budget}}**.

Formula: `pages_or_stories × browsers × viewports × themes × locales × states`.

If exceeding tool quota → reduce axes (drop a viewport, drop a theme, scope state coverage).

## Coverage Tiers

| Tier | Pages / Stories | Browsers | Viewports | Themes | Total |
| ---- | --------------- | -------- | --------- | ------ | ----- |
| Tier 1 critical | | | | | |
| Tier 2 standard | | | | | |
| Tier 3 long-tail | | | | | |

## Notes

{{notes}}
