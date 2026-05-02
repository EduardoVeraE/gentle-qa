# Playwright `toHaveScreenshot` Reference

Reference for visual regression with **Playwright's built-in `toHaveScreenshot` API** (`@playwright/test` 1.40+). Self-hosted, free, baselines committed to the repo. Highest control, highest determinism burden.

---

## 1. Introduction

Playwright Test ships a first-class visual comparison API. Each test calls `expect(page).toHaveScreenshot()`; on first run it writes a baseline; on subsequent runs it diffs against the baseline using a pixel-diff algorithm (pixelmatch under the hood).

Key differences from Percy / Chromatic:

- **Baselines live in the repo** — `tests/__snapshots__/<test-name>-<browser>-<platform>.png`.
- **PR diff is visible in GitHub** — reviewers see the binary diff in the file changes tab.
- **No SaaS, no quota** — runs on CI minutes only.
- **No review UI** — accepting a baseline = regenerating + committing the file.
- **Determinism is your job** — OS font rendering, browser version, GPU all affect output.

### When to choose this over Percy/Chromatic

| Need | Playwright built-in | Percy | Chromatic |
|------|---------------------|-------|-----------|
| Zero SaaS cost | Yes | No | No |
| Baselines in PR diff | Yes | No | No |
| No vendor lock-in | Yes | No | No |
| Review UI (web app) | No | Yes | Yes |
| Perceptual diff (anti-aliasing tolerant) | No | Yes | Partial |
| Cross-platform determinism without effort | No | Yes | Yes |

---

## 2. Setup

### 2.1 Pin browser AND OS

The single most common cause of flake: snapshots taken on a developer's macOS render differently from Linux CI. **You MUST run snapshots in a pinned environment.**

The recommended setup:

```yaml
# .github/workflows/visual.yml
jobs:
  visual:
    runs-on: ubuntu-latest
    container: mcr.microsoft.com/playwright:v1.49.0-jammy
    steps:
      - uses: actions/checkout@v4
      - run: npm ci
      - run: npx playwright test --project=visual
```

The `mcr.microsoft.com/playwright:v1.x.x-jammy` Docker image pins:

- Ubuntu 22.04 (jammy) base
- Playwright browsers at exact version
- System fonts at exact version

Without this pinning, snapshots taken locally fail in CI and vice versa — endless flake.

### 2.2 `playwright.config.ts`

```ts
import { defineConfig, devices } from "@playwright/test";

export default defineConfig({
  testDir: "./tests",
  expect: {
    toHaveScreenshot: {
      // Default thresholds for ALL screenshot assertions
      maxDiffPixelRatio: 0.01,    // 1% of pixels may differ
      threshold: 0.2,              // per-pixel YIQ color delta tolerance
      animations: "disabled",
      caret: "hide",
      scale: "css",
    },
  },
  projects: [
    {
      name: "visual-chromium",
      use: { ...devices["Desktop Chrome"] },
      testMatch: /.*\.visual\.spec\.ts/,
    },
  ],
});
```

Key options:

- `maxDiffPixelRatio: 0.01` — tolerate up to 1% pixel variance. Tighten to 0.001 once stable.
- `threshold: 0.2` — color sensitivity per pixel (0 = exact, 1 = anything). 0.2 ignores anti-aliasing noise.
- `animations: "disabled"` — Playwright stops animations and clamps them to their final frame.
- `caret: "hide"` — hide the text caret in input fields (it blinks).
- `scale: "css"` — use CSS pixels, not device pixels — reproducible across DPI.

---

## 3. Authoring screenshot tests

### 3.1 Full page

```ts
import { test, expect } from "@playwright/test";

test("homepage hero", async ({ page }) => {
  await page.goto("/");
  await expect(page.locator("[data-testid='hero']")).toBeVisible();
  await expect(page).toHaveScreenshot("homepage-hero.png", {
    fullPage: true,
  });
});
```

### 3.2 Element-scoped (preferred)

Element snapshots are 10× more stable than full-page — tighter scope, fewer dynamic regions.

```ts
test("nav bar — logged out", async ({ page }) => {
  await page.goto("/");
  const nav = page.getByRole("navigation");
  await expect(nav).toHaveScreenshot("nav-logged-out.png");
});
```

### 3.3 Masking dynamic regions

```ts
test("dashboard with timestamp", async ({ page }) => {
  await page.goto("/dashboard");
  await expect(page).toHaveScreenshot("dashboard.png", {
    mask: [
      page.getByTestId("last-updated"),
      page.getByTestId("user-avatar"),
      page.locator(".ad-slot"),
    ],
    maskColor: "#FF00FF",   // hot pink — easy to spot in diffs
  });
});
```

Masked regions are filled with `maskColor`. They are not compared. Use generously for any region that changes per-run (timestamps, ads, A/B variants, animated regions).

### 3.4 Disable animations and freeze fonts

```ts
import { test as base, expect } from "@playwright/test";

export const test = base.extend({
  page: async ({ page }, use) => {
    // Wait for fonts before any test runs
    await page.addInitScript(() => {
      // @ts-expect-error - inject CSS to kill animations
      const style = document.createElement("style");
      style.textContent = `
        *, *::before, *::after {
          animation-duration: 0s !important;
          animation-delay: 0s !important;
          transition-duration: 0s !important;
          transition-delay: 0s !important;
        }
      `;
      document.head.appendChild(style);
    });
    await use(page);
    await page.evaluate(() => document.fonts.ready);
  },
});
```

Use this fixture in every visual spec.

---

## 4. Baseline workflow

### 4.1 Generating baselines

```bash
# In the SAME environment that CI uses (Docker recommended)
docker run --rm -v "$PWD:/work" -w /work mcr.microsoft.com/playwright:v1.49.0-jammy \
  npx playwright test --update-snapshots --project=visual-chromium
```

Commit the resulting `__snapshots__/` directory. Review the PNGs in your PR diff before merging — git renders binary image diffs.

NEVER run `--update-snapshots` on a developer machine if CI uses Docker. The PNGs will diff every time.

### 4.2 Updating baselines (intentional change)

1. Make your code change.
2. Push PR. CI fails with diff artifacts.
3. Inspect the diff: `playwright-report/data/<id>-diff.png`.
4. If the change is intended, on a fresh branch:

```bash
docker run --rm -v "$PWD:/work" -w /work mcr.microsoft.com/playwright:v1.49.0-jammy \
  npx playwright test --update-snapshots
```

5. Commit and push. CI passes. Reviewer sees binary diff in PR.

### 4.3 Rejecting a diff

If the diff is a real regression: do nothing to the snapshots. Fix the code. Re-run.

NEVER `--update-snapshots` to "make CI pass". This silently swallows regressions and is the #1 misuse of this API.

---

## 5. CI integration (GitHub Actions)

```yaml
name: Visual

on: [pull_request]

jobs:
  visual:
    runs-on: ubuntu-latest
    container:
      image: mcr.microsoft.com/playwright:v1.49.0-jammy
    steps:
      - uses: actions/checkout@v4
      - run: npm ci
      - run: npx playwright test --project=visual-chromium
      - uses: actions/upload-artifact@v4
        if: failure()
        with:
          name: playwright-visual-report
          path: |
            playwright-report/
            test-results/
          retention-days: 7
```

The `test-results/` directory contains, per failed test:

- `*-actual.png` — what was rendered this run
- `*-expected.png` — the committed baseline
- `*-diff.png` — pixelmatch output highlighting differences

Reviewers download the artifact and inspect.

---

## 6. Browser × viewport matrix

```ts
// playwright.config.ts
projects: [
  {
    name: "visual-chromium-mobile",
    use: { ...devices["Pixel 5"] },
    testMatch: /.*\.visual\.spec\.ts/,
  },
  {
    name: "visual-chromium-desktop",
    use: { ...devices["Desktop Chrome"], viewport: { width: 1280, height: 720 } },
    testMatch: /.*\.visual\.spec\.ts/,
  },
  {
    name: "visual-webkit-mobile",
    use: { ...devices["iPhone 14"] },
    testMatch: /.*\.visual\.spec\.ts/,
  },
],
```

Each project creates a SEPARATE baseline file (filename suffix includes the project name). Cost is in CI minutes, not SaaS quota.

---

## 7. Flake mitigation

In order of frequency:

1. **OS font rendering** — solved by Docker pinning. No other solution.
2. **Animations** — solved by `animations: "disabled"` + the CSS fixture above.
3. **Caret blink** — `caret: "hide"`.
4. **Async data** — wait for specific elements / responses BEFORE `toHaveScreenshot`.
5. **Custom fonts** — wait for `document.fonts.ready`.
6. **GPU rendering (canvas, WebGL)** — mask the canvas region; pixel diff cannot survive GPU variance.

```ts
// Robust pre-snapshot wait pattern
await page.goto("/dashboard");
await expect(page.getByTestId("loaded-marker")).toBeVisible();
await page.evaluate(() => document.fonts.ready);
await page.waitForLoadState("networkidle"); // last resort; prefer specific waits
await expect(page).toHaveScreenshot("dashboard.png");
```

---

## 8. Threshold tuning strategy

Start permissive, tighten over time. NEVER start strict.

| Stage | `maxDiffPixelRatio` | `threshold` | When |
|-------|---------------------|-------------|------|
| Bootstrap | 0.05 | 0.3 | First 2 weeks; learning what flakes |
| Stable | 0.01 | 0.2 | After flake root-causes are fixed |
| Strict | 0.001 | 0.1 | For design-system primitives only |

Per-test override when needed:

```ts
await expect(page).toHaveScreenshot("noisy-canvas.png", {
  maxDiffPixelRatio: 0.05,
});
```

Document every per-test override with a comment explaining why.

---

## 9. Storage and Git considerations

Baselines are PNGs in the repo. Watch for:

- **Repo size growth**: 50 tests × 3 projects × ~50KB ≈ 7.5 MB. Fine. But 500 tests × 5 projects × 200 KB = 500 MB → use Git LFS.
- **Diff visibility**: GitHub renders binary image diffs. Reviewers can see baseline changes inline. THIS IS A FEATURE — use it.
- **Merge conflicts**: two PRs both update the same baseline → binary conflict. Pick one, regenerate the other on top.

For repos with hundreds of baselines, configure Git LFS:

```gitattributes
tests/__snapshots__/**/*.png filter=lfs diff=lfs merge=lfs -text
```

---

## 10. Common gotchas

- Running `--update-snapshots` on macOS when CI is Linux Docker → baselines fail every build. Always update inside Docker.
- Forgetting `animations: "disabled"` → tests pass locally, fail randomly in CI when an animation lands mid-capture.
- Snapshotting before fonts load → text rendering differs across runs. Always `await document.fonts.ready`.
- Snapshotting `<canvas>` content directly → GPU-dependent. Mask the canvas, snapshot the surrounding UI.
- Updating a snapshot to "make CI green" without reviewing the diff → silent regression. Treat every baseline change as a code review item.
- `fullPage: true` on long pages → tiny scroll-position differences cause large diffs. Prefer element-scoped snapshots.
- Baselines auto-generate on first run if missing — a brand-new test ALWAYS passes. Reviewer must inspect the baseline PNG in PR diff.

---

## 11. References

- Playwright `toHaveScreenshot` API docs (within `@playwright/test` types).
- pixelmatch (the underlying diff algorithm) for threshold semantics.
- `mcr.microsoft.com/playwright` Docker images on Microsoft Container Registry.
