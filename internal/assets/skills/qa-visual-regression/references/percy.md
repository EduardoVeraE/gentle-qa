# Percy Reference

Reference for visual regression with **Percy** (BrowserStack), focused on the Playwright + TypeScript SDK as the primary integration and the `@percy/cli` for static-site snapshots as an alternative.

---

## 1. Introduction

Percy is a SaaS visual review platform. It captures DOM snapshots from your tests, renders them server-side across multiple browsers and widths, performs **perceptual diffing**, and presents an approval UI gating PR merges.

Key concepts:

- **Snapshot**: a single captured page-state. Percy bills per snapshot × browser × width.
- **Build**: a collection of snapshots associated with a Git commit.
- **Baseline build**: the build a current build is compared against. Defaults to the parent commit on the target branch.
- **Approval**: a human review step in the Percy UI; PR check passes only when approved or no diff exists.

### Why Percy over alternatives

| Need | Choose Percy | Choose Chromatic | Choose Playwright built-in |
|------|--------------|------------------|----------------------------|
| Multi-browser coverage out of the box (Chrome + Firefox + Safari + Edge) | Yes | No (Chrome only on free) | Manual (one project per browser) |
| Full-page web app snapshots, not Storybook | Yes | Painful | Yes |
| Perceptual diff that ignores anti-aliasing | Yes | No (pixel) | No (pixel) |
| Free for OSS | Yes (generous) | Yes (smaller) | N/A (self-hosted) |
| No vendor lock-in | No | No | Yes |

---

## 2. Setup

### 2.1 Install the CLI and SDK

```bash
# Pick ONE SDK based on your test runner
npm install --save-dev @percy/cli @percy/playwright   # Playwright
npm install --save-dev @percy/cli @percy/cypress      # Cypress
npm install --save-dev @percy/cli @percy/selenium-webdriver  # Selenium

# Static site snapshots (no test runner)
npm install --save-dev @percy/cli
```

### 2.2 Project token

Create a project at percy.io, copy the `PERCY_TOKEN`, and add it to CI secrets:

```bash
# Local development
export PERCY_TOKEN="..."

# GitHub Actions
# Settings → Secrets → Actions → PERCY_TOKEN
```

NEVER commit the token to the repo. It grants write access to the project.

### 2.3 `.percy.yml` configuration

Place at repo root. Minimum useful config:

```yaml
version: 2
snapshot:
  widths: [375, 768, 1280]
  min-height: 1024
  percy-css: |
    /* Hide elements that are inherently dynamic */
    .timestamp, .user-avatar-initials, [data-testid="ad-slot"] {
      visibility: hidden !important;
    }
    /* Disable animations and transitions */
    *, *::before, *::after {
      animation-duration: 0s !important;
      transition-duration: 0s !important;
    }
discovery:
  network-idle-timeout: 750
  allowed-hostnames: []
```

Key options:

- `widths`: viewport widths in pixels. Each adds a billed snapshot. Start with 3.
- `percy-css`: CSS injected only for Percy snapshots — use to hide dynamic regions globally.
- `discovery.network-idle-timeout`: how long Percy waits for assets after page settle. Increase to 1500-2000 for asset-heavy SPAs.

---

## 3. Capturing snapshots (Playwright + TypeScript)

### 3.1 Inside an existing E2E test

```ts
import { test, expect } from "@playwright/test";
import percySnapshot from "@percy/playwright";

test("checkout — empty cart", async ({ page }) => {
  await page.goto("/cart");
  await expect(page.getByText("Your cart is empty")).toBeVisible();
  await percySnapshot(page, "Cart - empty state");
});

test("checkout — populated cart", async ({ page }) => {
  await page.goto("/cart?seed=happy-path");
  await expect(page.getByRole("heading", { name: /order summary/i })).toBeVisible();
  await percySnapshot(page, "Cart - populated", {
    widths: [375, 1280],   // override per-snapshot
    minHeight: 1500,
  });
});
```

### 3.2 Running

```bash
# Wraps your test runner; Percy uploads snapshots to the cloud
npx percy exec -- npx playwright test
```

Percy auto-detects Git context (`PERCY_BRANCH`, `PERCY_COMMIT`, `PERCY_TARGET_BRANCH`) on common CI providers; override manually if needed.

---

## 4. Static-site snapshots (no test runner)

For marketing sites or static builds:

```bash
# Snapshot a list of URLs from a sitemap or manual list
npx percy snapshot ./public/sitemap.xml
# or
npx percy snapshot snapshots.yml
```

`snapshots.yml`:

```yaml
- name: Home
  url: http://localhost:8000/
- name: Pricing
  url: http://localhost:8000/pricing/
  waitForSelector: "[data-testid='pricing-table']"
```

Run a local server first (`npx serve ./public`), then `percy snapshot`.

---

## 5. CI integration (GitHub Actions)

```yaml
name: Visual

on: [pull_request]

jobs:
  percy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0   # Required for accurate baseline detection
      - uses: actions/setup-node@v4
        with:
          node-version: 20
      - run: npm ci
      - run: npx playwright install --with-deps chromium
      - run: npx percy exec -- npx playwright test
        env:
          PERCY_TOKEN: ${{ secrets.PERCY_TOKEN }}
```

Configure the Percy GitHub App so the `percy/<project>` check is required to merge. PR cannot merge until either: no diffs, or all diffs are approved in the Percy UI.

---

## 6. Baseline workflow

| Action | Command / Step |
|--------|----------------|
| Establish initial baselines | First green build on `main` after Percy is wired up — auto-becomes baseline |
| Diff a PR | Push PR → Percy build runs → comment on PR with diff link |
| Accept diffs | Open Percy UI → review each snapshot → click Approve |
| Reject diffs | Reject in UI → fix code → push → re-run |
| Reset all baselines (rare) | Re-baseline on `main` by pushing an empty commit |

Never approve diffs you did not introduce. Never approve in bulk without reviewing each.

---

## 7. Flake mitigation

The 4 sources of Percy flake, in order of frequency:

1. **Animations / transitions** — kill via `percy-css` global override.
2. **Dynamic content** — hide via `percy-css` selectors OR mock data at the API layer for visual tests.
3. **Async asset loading** — increase `discovery.network-idle-timeout`; wait for specific `data-testid` to be visible before `percySnapshot`.
4. **Font loading** — wait for `document.fonts.ready` before snapshotting:

```ts
await page.evaluate(() => document.fonts.ready);
await percySnapshot(page, "Page");
```

---

## 8. Browser × viewport matrix

Percy renders each snapshot across the browsers your project plan includes:

- **Free / OSS**: Chrome (latest)
- **Paid**: Chrome, Firefox, Safari, Edge (latest)

Combined cost: `snapshots × widths × browsers`. Watch the quota:

```
50 snapshots × 3 widths × 4 browsers = 600 billed snapshots per build
```

Keep snapshot count under 50 for most projects; expand only with quota headroom.

---

## 9. Threshold tuning

Percy uses perceptual diffing by default — no per-test threshold to tune. The two knobs you have:

- **`percy-css`** to suppress noise globally.
- **`enable-javascript: false`** at the snapshot level if a page has uncontrollable JS animation. Use sparingly — disabling JS may break SPA rendering.

```ts
await percySnapshot(page, "Static page", { enableJavaScript: false });
```

---

## 10. Common gotchas

- Percy snapshots the DOM at capture time, then renders server-side. Custom fonts loaded from third-party CDNs may fail to load on Percy's side — bundle fonts with the app or use `data:` URLs.
- Snapshots inside iframes are not captured by default. Set `enable-javascript: true` and ensure the iframe origin is in `allowed-hostnames`.
- Percy does not support `<canvas>` or WebGL content meaningfully — the rendered image is what the browser renders, but server-side rendering may differ from client. Mask canvas regions.
- Free tier limits: 5,000 screenshots/month for OSS, 100/month for private. Plan accordingly.

---

## 11. References

- Percy docs: percy.io/docs (do not link; consult inside the SaaS).
- `@percy/playwright` SDK source for current API surface.
- BrowserStack Percy pricing page for current quotas.
