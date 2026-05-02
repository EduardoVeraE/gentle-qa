<!-- Skill: qa-visual-regression · Template: baseline-runbook -->
<!-- Placeholders: {{project}}, {{tool}}, {{owner}}, {{date}}, {{baseline_branch}}, {{ci_environment}}, {{accept_command}}, {{reject_workflow}}, {{escalation}}, {{notes}} -->

# Visual Baseline Runbook

| Field    | Value          |
| -------- | -------------- |
| Project  | {{project}}    |
| Tool     | {{tool}}       <!-- Percy / Chromatic / Playwright -->
| Owner    | {{owner}}      |
| Date     | {{date}}       |
| Baseline branch | {{baseline_branch}} <!-- e.g., main -->
| CI environment | {{ci_environment}} |

## 1. Establishing Baselines (first time)

1. Verify CI environment is pinned (Docker image tag for Playwright; project token for SaaS).
2. From `{{baseline_branch}}` after the first green build, baselines are established.
3. Confirm baseline count matches expected snapshot budget.

## 2. Accepting an Intentional Change

```
{{accept_command}}
```

<!-- Examples per tool -->
<!-- Percy:      Open Percy UI → Approve each diff → PR check turns green -->
<!-- Chromatic:  Open Chromatic UI → Approve each story → UI Review gate passes -->
<!-- Playwright: docker run ... npx playwright test --update-snapshots && git add tests/__snapshots__ && git commit -->

## 3. Rejecting a Change (real regression)

{{reject_workflow}}

Standard reject workflow:

1. Inspect the diff in the tool UI (or `*-diff.png` artifact for Playwright).
2. File a defect linking the build URL.
3. Fix the code; do NOT regenerate baselines.
4. Push fix; CI re-runs; diff resolves.

## 4. Cross-Branch Synchronization

| Scenario | Action |
| -------- | ------ |
| PR is behind main | Rebase on main; re-run visual job |
| Main updated baselines while PR open | Rebase; if PR also touches same baseline, regenerate ON TOP of main |
| Hotfix on release branch | Establish independent baseline on release branch |

## 5. Quota / Cost Hygiene (Percy / Chromatic)

- Audit snapshot count weekly; trim non-critical pages.
- For Chromatic: ensure TurboSnap is active (`onlyChanged: true` + `fetch-depth: 0`).
- For Percy: review widths × browsers multiplier before adding axes.

## 6. Escalation

{{escalation}}  <!-- e.g., Design system team owns Storybook baselines; Web platform team owns app baselines -->

## 7. Notes

{{notes}}
