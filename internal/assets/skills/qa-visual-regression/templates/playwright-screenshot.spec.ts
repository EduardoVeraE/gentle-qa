// Skill: qa-visual-regression · Template: playwright-screenshot
// Starter spec for Playwright `toHaveScreenshot`.
// Run inside Docker `mcr.microsoft.com/playwright:v1.x.x-jammy` for deterministic baselines.

import { test as base, expect } from "@playwright/test";

/**
 * Visual fixture: kills animations, hides caret, waits for fonts.
 * Use this in EVERY visual spec to eliminate the top sources of flake.
 */
const test = base.extend({
  page: async ({ page }, use) => {
    await page.addInitScript(() => {
      const style = document.createElement("style");
      style.textContent = `
        *, *::before, *::after {
          animation-duration: 0s !important;
          animation-delay: 0s !important;
          transition-duration: 0s !important;
          transition-delay: 0s !important;
        }
        input, textarea { caret-color: transparent !important; }
      `;
      document.head.appendChild(style);
    });
    await use(page);
  },
});

test.describe("homepage — visual", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/");
    // Wait for the page to be functionally ready BEFORE waiting for fonts.
    await expect(page.getByRole("heading", { level: 1 })).toBeVisible();
    await page.evaluate(() => document.fonts.ready);
  });

  test("hero section", async ({ page }) => {
    const hero = page.getByTestId("hero");
    await expect(hero).toHaveScreenshot("hero.png");
  });

  test("nav — logged out", async ({ page }) => {
    const nav = page.getByRole("navigation");
    await expect(nav).toHaveScreenshot("nav-logged-out.png");
  });

  test("full page with masked dynamics", async ({ page }) => {
    await expect(page).toHaveScreenshot("home-full.png", {
      fullPage: true,
      mask: [
        page.getByTestId("last-updated"),
        page.getByTestId("user-avatar"),
        page.locator(".ad-slot"),
      ],
      maskColor: "#FF00FF",
    });
  });
});

test.describe("dashboard — visual", () => {
  test("loaded state", async ({ page }) => {
    await page.goto("/dashboard?seed=visual-test");
    await expect(page.getByTestId("loaded-marker")).toBeVisible();
    await page.evaluate(() => document.fonts.ready);

    await expect(page.getByTestId("dashboard-grid")).toHaveScreenshot(
      "dashboard-grid.png",
    );
  });

  test("empty state", async ({ page }) => {
    await page.goto("/dashboard?seed=empty");
    await expect(page.getByText(/no items yet/i)).toBeVisible();

    await expect(page.getByTestId("dashboard-grid")).toHaveScreenshot(
      "dashboard-grid-empty.png",
    );
  });
});
