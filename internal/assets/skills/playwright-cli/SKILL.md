---
name: playwright-cli
description: >
  Browser automation CLI for SDET/QE Engineers. 40+ commands for navigation,
  interaction, network mocking, and exploratory testing sessions.
  Source: microsoft/playwright-cli (skills.sh, adapted for QE workflow).
  Trigger: When using playwright-cli for browser automation, scripted exploratory
  testing, visual regression capture, accessibility audits, or building CLI test flows.
license: Apache-2.0
metadata:
  author: microsoft (adapted for gentle-qa)
  version: "1.1"
  source: https://skills.sh/microsoft/playwright-cli/playwright-cli
  weekly_installs: "20.2K"
---

## ISTQB Mapping

| Aspect | Value |
|--------|-------|
| Test level | System Testing |
| Test type | Functional (exploratory), Non-functional (accessibility) |
| Techniques | Exploratory Testing, Error Guessing, Checklist-based Testing |
| Test oracle | Session notes + expected vs actual screenshots / state |

**Core principle**: playwright-cli is your scripted exploratory testing tool. Use it when you need to probe behavior quickly, capture evidence, or automate a session without building a full test suite. Every exploratory session MUST produce a session report.

---

## When to Use

- Rapid exploratory testing with recorded evidence (screenshots, traces, video)
- Reproducing a reported bug step-by-step before writing a formal test
- Accessibility audits on live pages
- API mocking during manual testing to isolate front-end behavior
- Multi-user/multi-role scenarios via named sessions

**Not here**: regression suites, BDD scenarios — use playwright-bdd for those.

---

## Critical Patterns

### Pattern 1: Exploratory Testing Session

```bash
# Start a timed exploratory session with evidence capture
playwright-cli tracing-start --screenshots --snapshots
playwright-cli video-start ./sessions/session-$(date +%Y%m%d-%H%M%S).webm

# Explore the feature
playwright-cli open https://staging.example.com/checkout
playwright-cli goto https://staging.example.com/checkout/payment

# Probe edge cases (Error Guessing technique)
playwright-cli fill "[name='card-number']" "0000 0000 0000 0000"  # Invalid card
playwright-cli click "[data-testid='submit']"
playwright-cli screenshot --full-page ./evidence/invalid-card-response.png

playwright-cli fill "[name='card-number']" "4111 1111 1111 1111"  # Valid Visa
playwright-cli fill "[name='cvv']" "999"                           # BVA: max CVV
playwright-cli click "[data-testid='submit']"

# Stop and save session evidence
playwright-cli tracing-stop ./sessions/trace-$(date +%Y%m%d).zip
```

### Pattern 2: Bug Reproduction Script

```bash
# Reproduce bug #1234: checkout fails when coupon applied + new card
playwright-cli -s bug1234 open https://staging.example.com
playwright-cli -s bug1234 state-load ./auth/registered-user.json  # Reuse auth state

playwright-cli -s bug1234 goto https://staging.example.com/cart
playwright-cli -s bug1234 fill "[data-testid='coupon-input']" "SAVE10"
playwright-cli -s bug1234 click "[data-testid='apply-coupon']"

# Network interception to observe the API call
playwright-cli -s bug1234 network

playwright-cli -s bug1234 goto https://staging.example.com/checkout
playwright-cli -s bug1234 fill "[name='card-number']" "4111 1111 1111 1111"
playwright-cli -s bug1234 click "[data-testid='submit-payment']"

playwright-cli -s bug1234 screenshot --full-page ./evidence/bug1234-reproduced.png
```

### Pattern 3: Auth State Management (reuse across sessions)

```bash
# One-time: capture authenticated state
playwright-cli open https://example.com/login
playwright-cli fill "[name='email']" "user@example.com"
playwright-cli fill "[name='password']" "password123"
playwright-cli click "[type='submit']"
playwright-cli state-save ./auth/registered-user.json

# Every session after: load state instead of logging in
playwright-cli state-load ./auth/registered-user.json
playwright-cli goto https://example.com/dashboard  # Already authenticated
```

### Pattern 4: Multi-Role Testing (named sessions)

```bash
# Parallel sessions: admin reviews, user submits
playwright-cli -s admin state-load ./auth/admin.json
playwright-cli -s user  state-load ./auth/user.json

playwright-cli -s user  goto https://example.com/orders/new
playwright-cli -s user  fill "[name='description']" "Test order"
playwright-cli -s user  click "[data-testid='submit-order']"

playwright-cli -s admin goto https://example.com/admin/orders
playwright-cli -s admin screenshot ./evidence/admin-sees-new-order.png

playwright-cli -s admin click "[data-testid='approve-order']"
playwright-cli -s user  reload  # User refreshes
playwright-cli -s user  screenshot ./evidence/user-sees-approved.png
```

### Pattern 5: Network Mocking (front-end isolation)

```bash
# Mock API to test front-end error handling in isolation
playwright-cli route "https://api.example.com/orders" \
  --method POST \
  --status 500 \
  --body '{"error": "Internal Server Error"}'

playwright-cli goto https://example.com/checkout
playwright-cli click "[data-testid='place-order']"
playwright-cli screenshot ./evidence/500-error-handling.png  # What does the UI show?

# Mock slow network (performance degradation)
playwright-cli route "https://api.example.com/products" \
  --delay 5000 \
  --status 200 \
  --body '{"items": [], "total": 0}'

playwright-cli goto https://example.com/products
playwright-cli screenshot ./evidence/slow-api-loading-state.png

playwright-cli unroute "https://api.example.com/orders"
playwright-cli unroute "https://api.example.com/products"
```

### Pattern 6: Accessibility Audit

```bash
# Capture page for accessibility review
playwright-cli open https://example.com/checkout
playwright-cli snapshot  # Captures ARIA tree + DOM state

# Tab through the page (keyboard navigation test)
playwright-cli press "Tab"
playwright-cli press "Tab"
playwright-cli press "Tab"
playwright-cli screenshot ./evidence/focus-state-3.png  # Is focus visible?
playwright-cli press "Enter"  # Does Enter activate focused element?

# Check color contrast manually on screenshots
playwright-cli screenshot --full-page ./evidence/full-page-a11y-review.png
```

### Pattern 7: Cross-Browser Verification

```bash
# Same flow, 3 browsers — visual diff evidence
playwright-cli --browser chromium open https://example.com/checkout
playwright-cli screenshot ./evidence/checkout-chromium.png

playwright-cli --browser firefox open https://example.com/checkout
playwright-cli screenshot ./evidence/checkout-firefox.png

playwright-cli --browser webkit open https://example.com/checkout
playwright-cli screenshot ./evidence/checkout-webkit.png

# Compare screenshots to detect rendering differences
```

---

## Anti-patterns — Never Do This

| Anti-pattern | Why it fails | Fix |
|---|---|---|
| Exploratory session with no evidence | Nothing to attach to bug report | Always run `tracing-start` + `screenshot` |
| Log in manually every session | Slow, inconsistent | Use `state-save` / `state-load` |
| Hardcoded `sleep 3` in scripts | Flaky on slow envs | Use `waitForSelector` patterns instead |
| Mocking without `unroute` | Leaks mock to next test | Always `unroute` after the scenario |
| No session notes | Evidence without context is useless | Write a short session charter before starting |

---

## Exploratory Session Template

Before every session, define:

```markdown
## Session Charter
**Mission**: Explore [feature] to find defects in [area]
**Scope**: [what's in], [what's out]
**Time box**: 45 minutes
**Technique**: Error Guessing / Boundary Value Analysis / Checklist

## Findings
- [bug found, severity, screenshot path]

## Session Coverage
- [what was tested]
- [what was NOT tested — for next session]
```

---

## Commands Reference

```bash
playwright-cli open <url>                           # Open browser
playwright-cli goto <url>                           # Navigate
playwright-cli click "<selector>"                   # Click element
playwright-cli fill "<selector>" "<value>"          # Fill input
playwright-cli screenshot --full-page <path>        # Capture evidence
playwright-cli tracing-start --screenshots          # Start trace
playwright-cli tracing-stop <path>                  # Save trace
playwright-cli state-save <path>                    # Save auth state
playwright-cli state-load <path>                    # Restore auth state
playwright-cli -s <name> <command>                  # Named session
playwright-cli route <url> --status <n> --body '…' # Mock response
playwright-cli unroute <url>                        # Remove mock
playwright-cli network                              # Inspect requests
playwright-cli console                              # Monitor console
playwright-cli --browser <firefox|webkit> <cmd>    # Select browser
```

---

## Resources

- [playwright-cli GitHub](https://github.com/microsoft/playwright-cli)
- [skills.sh](https://skills.sh/microsoft/playwright-cli/playwright-cli)
- [ISTQB Exploratory Testing](https://www.istqb.org)
- [Session-Based Test Management](https://www.satisfice.com/blog/archives/1100)
