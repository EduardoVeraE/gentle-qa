<!-- Skill: qa-visual-regression · Template: visual-diff-report -->
<!-- Placeholders: {{title}}, {{tool}}, {{build_url}}, {{pr_url}}, {{author}}, {{date}}, {{diff_count}}, {{classification}}, {{root_cause}}, {{action}}, {{evidence_link}}, {{notes}} -->

# Visual Diff Triage Report

| Field    | Value          |
| -------- | -------------- |
| Title    | {{title}}      <!-- e.g., Header redesign — diff triage -->
| Tool     | {{tool}}       |
| Build    | {{build_url}}  |
| PR       | {{pr_url}}     |
| Author   | {{author}}     |
| Date     | {{date}}       |
| Diff count | {{diff_count}} |

## Classification

{{classification}}  <!-- Real regression / Intentional design change / Rendering noise -->

> Decision tree: 1) Real regression → file defect, fix code, do NOT update baseline. 2) Intentional → accept new baseline; document in PR. 3) Noise → fix determinism BEFORE updating baseline.

## Root Cause (if regression or noise)

{{root_cause}}  <!-- e.g., CSS specificity bug introduced by Tailwind upgrade; OR animation landed mid-capture -->

## Action Taken

{{action}}  <!-- e.g., Code fix in PR #123; OR baseline regenerated and committed; OR animations disabled in fixture -->

## Evidence

- Build URL: {{build_url}}
- Diff artifact: {{evidence_link}}  <!-- e.g., playwright-report/data/<id>-diff.png -->
- Affected snapshots:
  - ...
  - ...

## Notes

{{notes}}
