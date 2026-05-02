---
name: qa-mobile-testing
description: Mobile functional testing toolkit for native iOS/Android and React Native apps using Appium, Detox, XCUITest, and Espresso. Covers device matrix strategy, gesture handling (tap, swipe, long-press, pinch), real-device vs simulator/emulator tradeoffs, and cloud device farms (BrowserStack, Sauce Labs, AWS Device Farm, Firebase Test Lab). Use when asked to plan a mobile test campaign, write Appium or Detox scripts, automate native iOS (XCUITest) or Android (Espresso) tests, build a device matrix, design a mobile test plan, file a mobile bug report with device fingerprint, or run tests on a cloud device farm. Trigger keywords - mobile testing, Appium, Detox, XCUITest, Espresso, native app test, React Native test, gesture, swipe, tap, mobile UI test, device farm, BrowserStack, Sauce Labs, AWS Device Farm, Firebase Test Lab, mobile bug report, device matrix, mobile test plan. NOT for mobile SECURITY testing - use `qa-owasp-security`. NOT for mobile accessibility - use `a11y-playwright-testing` (web a11y patterns apply) or escalate. NOT for mobile performance - use `k6-load-test` for backend, dedicated mobile perf is out of scope. NOT for general E2E web - use `playwright-e2e-testing` or `selenium-e2e-testing`.
license: MIT
metadata:
  author: gentleman-programming
  version: "1.0"
---

# Mobile Functional Testing Toolkit

ISTQB-aligned functional testing for mobile applications across native iOS, native Android, and cross-platform React Native stacks. Covers tooling (Appium, Detox, XCUITest, Espresso), device strategy (real devices, simulators, emulators, cloud farms), and the gesture/UI primitives that distinguish mobile testing from web testing.

**Core principle**: Mobile bugs are device-specific. Every mobile test artifact MUST include a device fingerprint (model, OS version, locale, network, orientation). A passing test on a single emulator is not coverage — it is a smoke check.

## When to Use This Skill

- Planning a **mobile test campaign** with a documented device matrix and OS coverage strategy
- Writing **Appium** scripts for cross-platform native iOS/Android tests
- Writing **Detox** tests for a React Native app (grey-box, in-process)
- Automating **iOS native** flows with **XCUITest**
- Automating **Android native** flows with **Espresso**
- Validating **gesture-driven UI** (tap, double-tap, long-press, swipe, pinch, drag)
- Choosing between **real devices, simulators, emulators, and cloud farms** for a given engagement
- Running tests on **BrowserStack**, **Sauce Labs**, **AWS Device Farm**, or **Firebase Test Lab**
- Filing **mobile bug reports** with full device fingerprint (model, OS, locale, network)
- Triaging **device-specific** failures that reproduce on one model/OS combination but not others

## ISTQB Layer

Layer 3 — Functional testing by level → **System Test (mobile)**.

This skill complements (does not replace) other layers:

| Layer | Coverage |
| ----- | -------- |
| 1. Foundation | `qa-manual-istqb` |
| 2. Strategy | `qa-manual-istqb`, `playwright-regression-strategy` |
| 3. Functional by level | `api-testing`, `playwright-e2e-testing`, `selenium-e2e-testing`, **`qa-mobile-testing`** (this skill) |
| 4. Non-functional by type | `qa-owasp-security`, `k6-load-test`, `a11y-playwright-testing` |
| 5. Tooling | `playwright-cli`, `playwright-mcp-inspect` |

## Frameworks Coverage

| Framework  | Best for                                    | Language       | Reference file                              |
| ---------- | ------------------------------------------- | -------------- | ------------------------------------------- |
| Appium     | Cross-platform native iOS/Android (black-box) | TS / Java / Python | `references/appium.md`                      |
| Detox      | React Native (grey-box, in-process sync)    | JS / TS        | `references/detox.md`                       |
| XCUITest   | iOS native (Apple-supported)                | Swift          | `references/native-and-device-strategy.md`  |
| Espresso   | Android native (Google-supported)           | Kotlin / Java  | `references/native-and-device-strategy.md`  |

Pick **Appium** when one suite must cover both platforms. Pick **Detox** when the app is React Native and you want fast, deterministic in-process tests. Pick **XCUITest / Espresso** when you need native fidelity, deep platform integration, or CI parity with platform engineers.

## Prerequisites

| Requirement                  | Notes                                                                                  |
| ---------------------------- | -------------------------------------------------------------------------------------- |
| Node.js 18+                  | Required for the artifact CLI (`scripts/mobile_artifacts.mjs`) and for Appium / Detox  |
| Xcode + iOS Simulator (Mac)  | Required for any iOS work — XCUITest, Detox iOS, Appium iOS                            |
| Android SDK + emulator       | Required for any Android work — Espresso, Detox Android, Appium Android                |
| Real device + developer mode | Strongly recommended for serious testing; emulators miss hardware/sensor/network bugs  |
| Cloud device farm account    | Optional — BrowserStack / Sauce Labs / AWS Device Farm / Firebase Test Lab for scale   |
| App build artifacts          | `.ipa` (iOS) and/or `.apk` / `.aab` (Android), debug-signed for instrumentation        |

## Quick Start

Generate mobile test artifacts from templates (CLI implemented in `scripts/mobile_artifacts.mjs`):

```bash
# List available templates
node scripts/mobile_artifacts.mjs list

# Create a mobile test plan
node scripts/mobile_artifacts.mjs create mobile-test-plan --out specs --project "MyApp"

# Create a device matrix
node scripts/mobile_artifacts.mjs create device-matrix --out specs --release "R1"

# Create a mobile bug report (with device fingerprint section)
node scripts/mobile_artifacts.mjs create mobile-bug-report --out specs/bugs --title "Swipe-to-delete fails on iOS 17.4"
```

Working starter projects live under `examples/`:

- `examples/appium-wdio-ts/` — Appium + WebdriverIO + TypeScript, runs against iOS Simulator and Android emulator
- `examples/detox-rn/` — Detox configured for a React Native app, with iOS and Android targets

## Workflows

### 1) Plan a mobile test campaign (device matrix, OS coverage)

1. Inventory the **target user base**: top device models, top OS versions, top locales, network conditions.
2. Build a **device matrix**: tier devices into Primary (must pass), Secondary (best-effort), Excluded (out of scope) — with explicit justification.
3. Cover **OS spread**: at minimum N, N-1 for each platform; extend if analytics show a long tail.
4. Define **environment variables**: locale, timezone, network (wifi/4G/offline), orientation, dark mode, accessibility settings.
5. Map test cases to the matrix; not every case runs on every device — risk-based selection.

Use: `templates/device-matrix.md` and `templates/mobile-test-plan.md`.

### 2) Set up Appium for cross-platform tests

1. Install Appium server and the relevant drivers (`xcuitest`, `uiautomator2`).
2. Configure capabilities per platform (`platformName`, `platformVersion`, `deviceName`, `app`, `automationName`).
3. Pick a client (WebdriverIO recommended for TS) and structure tests with Page Objects.
4. Use **accessibility IDs** as the primary locator strategy — stable across platforms.
5. Wire up artifact capture (screenshots, video, device logs, Appium logs) for triage.

See `references/appium.md`.

### 3) Set up Detox for a React Native app

1. Install Detox and its platform dependencies (`applesimutils` on Mac, Android SDK build-tools).
2. Configure `.detoxrc.js` with iOS and Android device configs.
3. Build the app with the Detox-instrumented configuration (`detox build`).
4. Use `testID` props on RN components — Detox's primary locator.
5. Leverage Detox's built-in synchronization — avoid manual waits, let Detox idle the runtime.

See `references/detox.md`.

### 4) Choose between native (XCUITest / Espresso) and Appium

1. Default to **Appium** if the test suite must cover both platforms with shared logic.
2. Choose **XCUITest** (Swift) for iOS-only deep integration, in-platform CI, or platform-team ownership.
3. Choose **Espresso** (Kotlin/Java) for Android-only deep integration, fast in-process tests, or platform-team ownership.
4. Mixed strategy is valid: Appium for E2E user journeys, native frameworks for component-level UI tests.

See `references/native-and-device-strategy.md`.

### 5) Run tests on cloud farms (BrowserStack, Sauce, AWS Device Farm)

1. Pick a farm based on device coverage, region, integration with your CI, and budget.
2. Upload the app build per run (or use the farm's app storage to avoid re-uploads).
3. Configure capabilities for the farm (each has its own caps namespace — check vendor docs).
4. Parallelize across the device matrix; cap parallelism to your concurrency quota.
5. Pull artifacts (video, logs, screenshots) for failed sessions; archive per release.

See `references/cloud-farms.md` (created by sibling task).

### 6) Triage mobile bugs (device-specific reproduction)

1. Capture the **full device fingerprint**: model, OS version, build number, locale, timezone, network, orientation, accessibility settings.
2. Reproduce on a second device of the same model/OS to confirm device specificity.
3. If reproducible only on one model: pull device logs (`adb logcat`, Xcode console) for native errors.
4. If reproducible across models on one OS version: likely an OS-API regression — file with OS version as the discriminator.
5. Score severity using customer impact × device tier (a P0 on a Primary-tier device may be a P2 on Excluded).

Use: `templates/mobile-bug-report.md`.

## Inputs to Collect

- **App build**: `.ipa` (iOS, debug-signed for Simulator or ad-hoc for real devices) and/or `.apk` / `.aab` (Android).
- **Target devices and OS versions**: ranked by user analytics; document Primary / Secondary / Excluded tiers.
- **Test data**: accounts, payment methods, content fixtures, feature flags, backend environment.
- **Network conditions**: wifi, 4G/5G, throttled, offline — mobile apps must degrade gracefully.
- **Accessibility requirements**: VoiceOver (iOS), TalkBack (Android), Dynamic Type, large fonts, reduce-motion (escalate detailed a11y to `a11y-playwright-testing` or a dedicated mobile a11y skill).
- **Languages and locales**: layout direction (LTR/RTL), text expansion, date/number/currency formats.
- **Backend integration**: API base URL per environment, mock vs real services, contract test references.

## Outputs

| Artifact                                              | When produced                              | Template                               |
| ----------------------------------------------------- | ------------------------------------------ | -------------------------------------- |
| Mobile test plan                                      | Per release / per major feature            | `templates/mobile-test-plan.md`        |
| Device matrix                                         | Per release; reviewed quarterly            | `templates/device-matrix.md`           |
| Mobile bug report (with device fingerprint)           | Per confirmed defect                       | `templates/mobile-bug-report.md`       |
| Test artifacts (screenshots, video, device logs)      | Per failed session                         | Captured by framework / farm           |
| Appium / Detox / XCUITest / Espresso scripts          | Per automated test case                    | `examples/appium-wdio-ts/`, `examples/detox-rn/` |

## Exclusions

This skill is deliberately scoped. Do NOT use it for:

- **Mobile SECURITY testing** (OWASP Mobile Top 10, MobSF, reverse engineering) — use `qa-owasp-security`.
- **Mobile accessibility** (VoiceOver, TalkBack, WCAG mapping) — start with `a11y-playwright-testing` (web a11y patterns largely apply) or escalate to a dedicated mobile a11y workflow.
- **Mobile performance / load** — backend load belongs to `k6-load-test`; dedicated mobile client perf (FPS, jank, battery, memory profiling) is out of scope for this skill.
- **General E2E web testing** — use `playwright-e2e-testing` or `selenium-e2e-testing`.
- **Test planning, traceability, ISTQB foundations** — use `qa-manual-istqb`.

If a request mixes concerns (e.g., "test the mobile checkout under load and check for auth bypass"), split it: functional flows here, load with `k6-load-test`, auth bypass with `qa-owasp-security`.
