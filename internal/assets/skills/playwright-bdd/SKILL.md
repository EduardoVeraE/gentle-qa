---
name: playwright-bdd
description: >
  BDD/ATDD with Playwright and Cucumber for SDET/QE Engineers. ISTQB System and
  Acceptance test level. Trigger: When writing Gherkin feature files, step definitions,
  Cucumber scenarios, or integrating Playwright with playwright-bdd or @cucumber/cucumber.
license: Apache-2.0
metadata:
  author: gentle-qa
  version: "1.1"
---

## ISTQB Mapping

| Aspect | Value |
|--------|-------|
| Test level | System Testing, Acceptance Testing (UAT / ATDD) |
| Test type | Functional, Black-box |
| Techniques | Use Case Testing, Decision Table, Equivalence Partitioning |
| Approach | ATDD — tests are written BEFORE implementation, driven by requirements |
| Test oracle | Acceptance criteria in the feature file = expected behavior |

**Core principle**: A failing scenario means a user story is not yet done. A passing scenario means the requirement is provably met — nothing more.

---

## When to Use

- System-level behavior validation (user-visible flows)
- Acceptance criteria formalized as executable specifications
- Cross-functional scenarios where business stakeholders read the tests
- Regression suite for critical user journeys

**Not here**: unit logic, component rendering, API contract checks — use the right layer.

---

## Critical Patterns

### Pattern 1: Feature File — Gherkin Rules

```gherkin
# features/checkout.feature
Feature: Checkout flow
  As a registered user
  I want to complete a purchase
  So that my order is placed and confirmed

  # Background = precondition shared by ALL scenarios in this feature
  Background:
    Given I am authenticated as a registered user

  # Scenario: one acceptance criterion = one scenario
  Scenario: Successful purchase with credit card
    Given I have 1 item in my cart
    When I complete checkout with a valid credit card
    Then my order is confirmed
    And I receive a confirmation email

  # Scenario Outline: same flow, different input classes (Equivalence Partitioning)
  Scenario Outline: Payment method renders correct fields
    Given I am on the checkout page
    When I select "<payment_method>"
    Then I see "<expected_fields>"

    Examples:
      | payment_method | expected_fields          |
      | credit_card    | card number, CVV, expiry |
      | paypal         | email, password          |
      | bank_transfer  | IBAN, account holder     |
```

**Gherkin rules**:
- One `When` per scenario — multiple `When`s = multiple responsibilities
- Steps describe BEHAVIOR, not UI actions (`I complete checkout`, not `I click the blue button`)
- `Background` only for true preconditions — not setup that belongs in hooks
- Use `Scenario Outline` when applying Equivalence Partitioning across input classes

### Pattern 2: Step Definitions — Thin Layer

```typescript
// steps/checkout.steps.ts
import { Given, When, Then } from '@cucumber/cucumber';
import { expect } from '@playwright/test';
import { ICustomWorld } from '../support/world';

// Steps are thin — all logic lives in Page Objects
Given('I am authenticated as a registered user', async function (this: ICustomWorld) {
  await this.authPage.loginAs(this.testData.users.registered);
  await expect(this.page).toHaveURL('/dashboard');
});

When('I complete checkout with a valid credit card', async function (this: ICustomWorld) {
  await this.checkoutPage.fillPayment(this.testData.payment.validCard);
  await this.checkoutPage.submit();
});

Then('my order is confirmed', async function (this: ICustomWorld) {
  // Test oracle: what EXACTLY proves this passed?
  await expect(this.page).toHaveURL(/\/orders\/\d+\/confirmation/);
  await expect(this.page.getByTestId('order-status')).toHaveText('Confirmed');
});

Then('I receive a confirmation email', async function (this: ICustomWorld) {
  // Test oracle: check via API or mail sink, never skip this assertion
  const email = await this.mailSink.getLastEmailFor(this.testData.users.registered.email);
  expect(email.subject).toContain('Order confirmed');
});
```

### Pattern 3: Custom World — Shared Context

```typescript
// support/world.ts
import { World, IWorldOptions, setWorldConstructor } from '@cucumber/cucumber';
import { BrowserContext, Page } from '@playwright/test';
import { CheckoutPage } from '../pages/CheckoutPage';
import { AuthPage } from '../pages/AuthPage';
import { MailSink } from '../support/mail-sink';

export interface ICustomWorld extends World {
  context: BrowserContext;
  page: Page;
  checkoutPage: CheckoutPage;
  authPage: AuthPage;
  mailSink: MailSink;
  testData: TestData;
}

class CustomWorld extends World implements ICustomWorld {
  context!: BrowserContext;
  page!: Page;
  checkoutPage!: CheckoutPage;
  authPage!: AuthPage;
  mailSink!: MailSink;
  testData: TestData = loadTestData();

  constructor(options: IWorldOptions) {
    super(options);
  }
}

setWorldConstructor(CustomWorld);
```

### Pattern 4: Hooks — Lifecycle Control

```typescript
// support/hooks.ts
import { Before, After, BeforeAll, AfterAll, Status } from '@cucumber/cucumber';
import { chromium, Browser } from '@playwright/test';
import { ICustomWorld } from './world';

let browser: Browser;

BeforeAll(async () => {
  browser = await chromium.launch({ headless: process.env.CI === 'true' });
});

AfterAll(async () => browser.close());

Before(async function (this: ICustomWorld) {
  this.context = await browser.newContext({ baseURL: process.env.BASE_URL });
  this.page = await this.context.newPage();
  this.checkoutPage = new CheckoutPage(this.page);
  this.authPage = new AuthPage(this.page);
  this.mailSink = new MailSink();
});

After(async function (this: ICustomWorld, scenario) {
  if (scenario.result?.status === Status.FAILED) {
    // Attach evidence — always, on failure
    this.attach(await this.page.screenshot({ fullPage: true }), 'image/png');
  }
  await this.context.close();
});
```

### Pattern 5: Page Object — Behavioral API

```typescript
// pages/CheckoutPage.ts
import { Page, Locator, expect } from '@playwright/test';

export class CheckoutPage {
  // Selectors: data-testid over CSS/XPath — resilient to style changes
  private readonly cardNumber: Locator;
  private readonly cvv: Locator;
  private readonly expiry: Locator;
  private readonly submitBtn: Locator;

  constructor(private readonly page: Page) {
    this.cardNumber = page.getByTestId('card-number');
    this.cvv = page.getByTestId('cvv');
    this.expiry = page.getByTestId('expiry');
    this.submitBtn = page.getByTestId('submit-payment');
  }

  async fillPayment(card: { number: string; cvv: string; expiry: string }) {
    await this.cardNumber.fill(card.number);
    await this.cvv.fill(card.cvv);
    await this.expiry.fill(card.expiry);
  }

  async submit() {
    await this.submitBtn.click();
    await this.page.waitForLoadState('networkidle');
  }
}
```

### Pattern 6: Config — playwright-bdd

```typescript
// cucumber.config.ts
import { defineConfig } from '@playwright/test';
import { defineBddConfig } from 'playwright-bdd';

const testDir = defineBddConfig({
  features: 'features/**/*.feature',
  steps: ['steps/**/*.steps.ts', 'support/hooks.ts'],
});

export default defineConfig({
  testDir,
  reporter: [
    ['html', { outputFolder: 'reports/playwright' }],
    ['json', { outputFile: 'reports/cucumber-report.json' }],
  ],
  use: {
    baseURL: process.env.BASE_URL ?? 'http://localhost:3000',
    headless: true,
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
    trace: 'retain-on-failure',
  },
});
```

---

## Anti-patterns — Never Do This

| Anti-pattern | Why it fails | Fix |
|---|---|---|
| `When I click the blue submit button` | Couples to UI, not behavior | `When I submit my order` |
| Multiple `When` in one scenario | Tests two things at once | Split into two scenarios |
| `And I wait for 3 seconds` | Flaky, non-deterministic | `waitForSelector` / `waitForResponse` |
| Assertions inside `When` steps | `When` is action, `Then` is oracle | Move assertions to `Then` |
| Business logic in step definitions | Breaks reuse, hides complexity | Move to Page Object |
| Skipping `Then` for side effects | No oracle = no test | Always assert the outcome |
| Giant `Background` with 10 steps | Obscures scenario intent | Extract to fixture/hook |

---

## Test Oracle Checklist

Before marking a scenario as complete, verify:
- [ ] `Then` steps assert the OUTCOME visible to the user, not internal state
- [ ] URL, page title, or visible element confirms the flow completed
- [ ] Side effects (emails, DB writes, API calls) are also asserted — not assumed
- [ ] Failure message is readable: "expected '/orders/123/confirmation' but got '/checkout'"

---

## Decision Tree

```
Writing a test?
├── User-visible behavior? → Feature file + step definition
├── Internal logic/function? → Unit test (wrong layer here)
├── API contract? → Karate DSL (wrong layer here)
│
Feature file structure?
├── One acceptance criterion → One Scenario
├── Same flow, different inputs → Scenario Outline (Equivalence Partitioning)
├── Shared precondition → Background (max 3 steps)
│
Step is getting complex?
├── > 3 lines? → Extract to Page Object method
├── Reused across features? → Move to common steps file
├── Network call needed? → Abstract in World or helper
│
Test is flaky?
├── Hard-coded wait? → Replace with explicit wait
├── Shared test data? → Isolate per scenario
├── Race condition? → Use waitForResponse / waitForEvent
└── Env-dependent? → Add to env-specific config
```

---

## Project Structure

```
tests/
├── features/
│   ├── checkout.feature
│   ├── login.feature
│   └── search.feature
├── steps/
│   ├── checkout.steps.ts
│   ├── login.steps.ts
│   └── common.steps.ts
├── pages/
│   ├── CheckoutPage.ts
│   ├── AuthPage.ts
│   └── BasePage.ts
├── support/
│   ├── world.ts
│   ├── hooks.ts
│   ├── mail-sink.ts
│   └── test-data.ts
├── reports/
└── cucumber.config.ts
```

---

## Commands

```bash
npx playwright test                          # Run all BDD tests
npx playwright test --grep "@smoke"          # Run by tag
npx playwright test --grep "Checkout flow"  # Run by feature name
npx playwright show-report                  # Open HTML report
npx playwright test --debug                 # Debug mode
```

---

## Resources

- [playwright-bdd](https://github.com/vitalets/playwright-bdd)
- [@cucumber/cucumber](https://github.com/cucumber/cucumber-js)
- [ISTQB Glossary](https://glossary.istqb.org)
