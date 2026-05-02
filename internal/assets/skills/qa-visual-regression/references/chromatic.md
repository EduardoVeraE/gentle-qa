# Chromatic Reference

Reference for visual regression with **Chromatic** (chromatic.com), focused on the Storybook-first workflow that is Chromatic's primary value proposition. Chromatic also supports Playwright/Cypress capture, but if you are not using Storybook, Percy or Playwright built-in are usually better fits.

---

## 1. Introduction

Chromatic is a SaaS visual review platform built by the Storybook maintainers. It builds your Storybook in the cloud, captures snapshots of every story across configured viewports/browsers, performs **pixel diffing** with configurable thresholds, and presents a UI Tests + UI Review workflow.

Key concepts:

- **Story**: a Storybook component instance — Chromatic's unit of capture.
- **Build**: snapshots captured for a single Git commit.
- **Baseline**: the build a current build is compared against; tracked per-branch.
- **TurboSnap**: only re-snapshots stories whose component dependency graph changed since the baseline. Cuts cost 70-90% on large libraries.
- **UI Tests**: pixel-diff gate. Fail = something changed visually.
- **UI Review**: human review gate. Required for merge.

### Why Chromatic over alternatives

| Need | Choose Chromatic | Choose Percy | Choose Playwright built-in |
|------|------------------|--------------|----------------------------|
| Storybook-first design system | Yes | Painful | Painful |
| TurboSnap to reduce cost on large libraries | Yes | No | N/A |
| Per-component visual review (not per-page) | Yes | No | No |
| Full-page web app flows | No | Yes | Yes |
| Multi-browser on free tier | No (Chrome only) | Limited | Free, manual setup |
| No vendor lock-in | No | No | Yes |

**Bottom line**: Chromatic's killer feature is Storybook integration + TurboSnap. If you do not use Storybook, this skill recommends Percy or Playwright built-in instead.

---

## 2. Setup

### 2.1 Prerequisites

- Storybook 7+ (Storybook 6 supported but legacy).
- A Storybook build that runs clean: `npm run build-storybook` produces a `storybook-static/` directory.

### 2.2 Install and authenticate

```bash
npm install --save-dev chromatic

# Project token from chromatic.com → Settings → Manage
export CHROMATIC_PROJECT_TOKEN="..."
```

Add token to CI secrets — same security posture as Percy. Never commit.

### 2.3 First build

```bash
# From repo root, on main
npx chromatic --project-token=$CHROMATIC_PROJECT_TOKEN
```

The first build on `main` becomes the baseline for all stories. Subsequent builds diff against the latest baseline on the target branch.

### 2.4 `chromatic.config.json`

```json
{
  "projectId": "Project:abc123",
  "buildScriptName": "build-storybook",
  "onlyChanged": true,
  "autoAcceptChanges": "main",
  "exitZeroOnChanges": false,
  "exitOnceUploaded": false
}
```

Key options:

- `onlyChanged: true` — enables TurboSnap. Mandatory for any non-trivial library.
- `autoAcceptChanges: "main"` — auto-accept new baselines on `main` after merge. Without this, every merge requires manual re-approval in Chromatic.
- `exitZeroOnChanges: false` — fail the CI step on any diff (correct for PR gating).
- `exitOnceUploaded: true` — option for fire-and-forget; default keeps CI green and runs the gate via GitHub check.

---

## 3. Capturing stories

Chromatic snapshots **every story** by default. Control coverage via story-level parameters:

```ts
// Button.stories.tsx
import type { Meta, StoryObj } from "@storybook/react";
import { Button } from "./Button";

const meta: Meta<typeof Button> = {
  component: Button,
  parameters: {
    chromatic: {
      viewports: [375, 768, 1280],
      modes: {
        light: { theme: "light" },
        dark: { theme: "dark" },
      },
      delay: 300,            // wait 300ms before capture
      pauseAnimationAtEnd: true,
    },
  },
};
export default meta;

export const Primary: StoryObj<typeof Button> = {
  args: { variant: "primary", children: "Submit" },
};

// Skip a story from visual capture (still rendered in Storybook UI)
export const Animated: StoryObj<typeof Button> = {
  args: { variant: "primary", children: "Loading", loading: true },
  parameters: { chromatic: { disableSnapshot: true } },
};
```

### 3.1 Modes (themes / locales)

Modes multiply snapshots: `viewports × modes × browsers`. Be deliberate.

```ts
parameters: {
  chromatic: {
    modes: {
      "light-en": { theme: "light", locale: "en" },
      "dark-en":  { theme: "dark",  locale: "en" },
      "light-ar": { theme: "light", locale: "ar" },  // RTL
    },
  },
},
```

---

## 4. CI integration (GitHub Actions)

```yaml
name: Visual

on:
  push:
    branches: [main]
  pull_request:

jobs:
  chromatic:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0   # Required for TurboSnap dependency graph
      - uses: actions/setup-node@v4
        with:
          node-version: 20
      - run: npm ci
      - uses: chromaui/action@latest
        with:
          projectToken: ${{ secrets.CHROMATIC_PROJECT_TOKEN }}
          onlyChanged: true
          autoAcceptChanges: main
          exitZeroOnChanges: false
```

Configure required GitHub checks:

- `UI Tests` — pixel-diff gate.
- `UI Review` — human approval gate.

Both must pass before merge. UI Review is the human gate — without it, Chromatic becomes a notification system, not a quality gate.

---

## 5. Baseline workflow

| Action | Step |
|--------|------|
| Establish initial baselines | First green build on `main` after Chromatic wired up |
| Diff a PR | Push PR → Chromatic builds → check on PR |
| Accept diffs | Open Chromatic UI → review each story diff → Approve |
| Reject diffs | Deny in UI → fix code → push → re-run |
| Auto-accept on merge | `autoAcceptChanges: main` in config; Chromatic accepts the new baseline once PR merges |

The standard PR flow:

1. Push PR.
2. Chromatic builds; UI Tests reports diffs.
3. Open Chromatic UI; review each diff.
4. Approve diffs that are intended; deny diffs that are bugs.
5. Once all denied diffs are fixed and re-approved, UI Review gate passes.
6. Merge. New baselines auto-promote on `main`.

---

## 6. TurboSnap

TurboSnap is the cost-saver. It builds a dependency graph from your Storybook's webpack/vite stats and re-snapshots only stories whose dependency closure changed.

To enable:

1. `onlyChanged: true` in config.
2. Ensure your bundler emits stats. For Vite + Storybook 7+, this is automatic.
3. `fetch-depth: 0` in CI checkout — TurboSnap needs the full Git history.

Verify it works: in the Chromatic UI, the build report shows "TurboSnap saved N snapshots".

Without TurboSnap, a 200-story library × 3 viewports × 2 modes = 1,200 snapshots **every PR**. Quota exhaustion within a sprint.

---

## 7. Browser × viewport matrix

- **Free / OSS**: Chrome only.
- **Team / Pro**: Chrome + Firefox + Safari + Edge.

Each browser multiplies cost. For most teams Chrome-only is sufficient — visual regressions are 95% CSS-engine-agnostic in 2025.

Viewports: keep to 3 (375 / 768 / 1280) unless you have specific wide-layout components.

---

## 8. Flake mitigation

The 4 sources of Chromatic flake:

1. **Animations / transitions** — `pauseAnimationAtEnd: true` in story params, or end the story in a final state.
2. **Async data loading** — Stories should be deterministic; mock data with MSW or Storybook decorators.
3. **Font loading** — Storybook 7+ waits for `document.fonts.ready` automatically; older versions need `delay: 500`.
4. **Random data** — Faker, Math.random in stories causes diffs every build. Seed all randomness or use fixtures.

```ts
// Decorator that forces deterministic data
parameters: {
  msw: {
    handlers: [
      rest.get("/api/users", (req, res, ctx) =>
        res(ctx.json([{ id: 1, name: "Ada Lovelace" }])),
      ),
    ],
  },
},
```

---

## 9. Threshold tuning

Chromatic uses pixel diffing with configurable threshold:

```ts
parameters: {
  chromatic: {
    diffThreshold: 0.063,     // 0-1, default 0.063
    diffIncludeAntiAliasing: false,
  },
},
```

- `diffThreshold`: lower = more sensitive. Below 0.05 → expect flake on font rendering. Above 0.1 → may miss real regressions.
- `diffIncludeAntiAliasing: false` — recommended; ignores sub-pixel font rendering differences.

Tune at the project level first (config file), then per-story for known-tricky cases.

---

## 10. Common gotchas

- **`autoAcceptChanges: main` is critical** — without it, every merge to main creates an unapproved build that blocks the next PR.
- TurboSnap requires `fetch-depth: 0` in CI checkout. With `fetch-depth: 1` (default), TurboSnap silently disables and you snapshot everything.
- A story that uses `Math.random()` or `new Date()` will diff on every build. Use Storybook decorators to freeze time/randomness.
- Modes multiply cost: 3 viewports × 4 modes × 4 browsers = 48 snapshots per story. Audit modes regularly.
- Free tier: 5,000 snapshots/month for any account. Easy to exhaust on a 100-story library without TurboSnap.

---

## 11. References

- Chromatic docs: chromatic.com/docs (consult inside SaaS).
- `chromaui/action` GitHub Action for current option surface.
- Storybook's `addon-themes` and MSW addon for deterministic story setup.
