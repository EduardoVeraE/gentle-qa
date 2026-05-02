# Appium 2.x Reference

Comprehensive reference for mobile test automation with Appium 2.x, focused on
TypeScript + WebdriverIO v9 as the primary stack and Java + Selenium-Appium as
an alternative.

---

## 1. Introduction

Appium is a cross-platform automation framework for native, hybrid, and mobile
web apps. Appium 2.x is a significant evolution of the project:

- **Modular driver architecture**: drivers (XCUITest, UiAutomator2, Espresso,
  Mac2, Windows, Chromium, Gecko) are installed as plugins, not bundled.
- **W3C WebDriver protocol only**: the legacy MJSONWP/JSON Wire protocol is
  removed. All capabilities use the W3C `appium:` vendor prefix.
- **Plugin system**: cross-cutting features (relaxed caps, image comparison,
  device farm orchestration) ship as plugins via `appium plugin install`.
- **Independent server**: `appium` is the server; clients (WebdriverIO,
  Selenium, appium-python-client) speak HTTP to it on port 4723 by default.

### Architecture at a glance

```
Test Process  --HTTP/W3C-->  Appium Server  -->  Driver  -->  Device
(WDIO/Selenium)              (Node.js)           (XCUITest /  (Real or Sim/
                                                  UiAutomator2) Emulator)
```

### When Appium beats native frameworks (XCUITest / Espresso directly)

| Need | Choose Appium | Choose Native (XCUITest/Espresso) |
|------|---------------|-----------------------------------|
| One test suite, both iOS + Android | Yes | No |
| No access to app source code (testing third-party builds) | Yes | No |
| Mix native + webview + mobile web flows | Yes | Painful |
| Use existing Selenium/WebdriverIO skills + tooling | Yes | No |
| Need lowest possible latency (<10ms per action) | No | Yes |
| Deep platform-specific gesture recorder integration | Maybe | Yes |
| Tight CI integration with Xcode/Gradle build pipeline | Maybe | Yes |

In short: **Appium is the right choice when cross-platform reach, language
freedom, or app-source independence matters more than raw speed**.

---

## 2. Setup

### 2.1 Install Appium 2.x server

```bash
# Node 18+ required
node --version

# Install the server globally
npm install -g appium

# Verify
appium --version  # should print 2.x

# Run with default port 4723
appium
```

### 2.2 Install drivers

Drivers are no longer bundled. Install them per platform:

```bash
# iOS
appium driver install xcuitest

# Android
appium driver install uiautomator2

# (Optional) Espresso for Android (faster, requires app rebuild)
appium driver install espresso

# List installed drivers
appium driver list --installed
```

### 2.3 iOS prerequisites (XCUITest driver)

- macOS host (XCUITest only runs on macOS).
- **Xcode 15+** with Command Line Tools: `xcode-select --install`.
- iOS Simulator runtimes installed via Xcode -> Settings -> Platforms.
- For real devices: an Apple Developer account, valid provisioning profile,
  device UDID (`xcrun xctrace list devices`).
- WebDriverAgent is built and signed automatically by the driver on first run;
  for real devices configure `appium:xcodeOrgId` and `appium:xcodeSigningId`.

Verify the host is ready:

```bash
appium driver doctor xcuitest
```

### 2.4 Android prerequisites (UiAutomator2 driver)

- **Android SDK** (cmdline-tools + platform-tools).
- `ANDROID_HOME` (or `ANDROID_SDK_ROOT`) pointing at the SDK.
- `JAVA_HOME` pointing at JDK 17+.
- An emulator AVD (`avdmanager create avd ...`) or a real device with USB
  debugging enabled (`adb devices`).

Useful env block (zsh / bash):

```bash
export ANDROID_HOME="$HOME/Library/Android/sdk"
export PATH="$ANDROID_HOME/platform-tools:$ANDROID_HOME/cmdline-tools/latest/bin:$PATH"
export JAVA_HOME="$(/usr/libexec/java_home -v 17)"
```

Verify:

```bash
appium driver doctor uiautomator2
adb devices
emulator -list-avds
```

### 2.5 Real devices vs simulators / emulators

See section 9 for tradeoffs. For setup:

- **iOS real**: device must be unlocked, paired with Xcode once, developer
  mode enabled (iOS 16+). Provisioning profile must include
  `WebDriverAgentRunner-Runner` bundle id.
- **Android real**: USB debugging on, vendor USB driver installed (Windows),
  computer authorized via the on-device dialog.

---

## 3. Capabilities

Appium 2.x is **strict W3C**. The session payload uses
`alwaysMatch`/`firstMatch`. All non-standard caps must be prefixed `appium:`.

### 3.1 W3C standard caps (no prefix)

| Cap | Meaning |
|-----|---------|
| `platformName` | `iOS` or `Android` |
| `browserName` | only for mobile web sessions (`Safari`, `Chrome`) |

Everything else is vendor-specific and needs `appium:`.

### 3.2 iOS minimal capabilities (XCUITest)

```ts
// wdio.conf.ts (excerpt) — iOS Simulator
export const iosCaps = {
  platformName: 'iOS',
  'appium:automationName': 'XCUITest',
  'appium:platformVersion': '17.4',
  'appium:deviceName': 'iPhone 15',
  'appium:app': '/absolute/path/to/MyApp.app', // OR appium:bundleId for installed
  'appium:newCommandTimeout': 240,
  'appium:noReset': false,
  'appium:wdaLaunchTimeout': 120000,
};
```

For an already-installed app, use:

```ts
'appium:bundleId': 'com.acme.myapp',
```

Real device additions:

```ts
'appium:udid': '00008110-0011223344556677',
'appium:xcodeOrgId': 'ABCDE12345',
'appium:xcodeSigningId': 'iPhone Developer',
'appium:updatedWDABundleId': 'com.acme.WebDriverAgentRunner',
```

### 3.3 Android minimal capabilities (UiAutomator2)

```ts
// wdio.conf.ts (excerpt) — Android Emulator
export const androidCaps = {
  platformName: 'Android',
  'appium:automationName': 'UiAutomator2',
  'appium:platformVersion': '14',
  'appium:deviceName': 'Pixel_7_API_34',
  'appium:app': '/absolute/path/to/app-debug.apk',
  // OR launch an installed app:
  // 'appium:appPackage': 'com.acme.myapp',
  // 'appium:appActivity': '.MainActivity',
  'appium:autoGrantPermissions': true,
  'appium:newCommandTimeout': 240,
};
```

Real device: add `'appium:udid'` from `adb devices`.

### 3.4 Common caps you almost always want

| Cap | Why |
|-----|-----|
| `appium:newCommandTimeout` | how long the server waits for the next command before killing the session (default 60s — too short for debugging) |
| `appium:noReset` / `appium:fullReset` | controls app data between tests (see section 10) |
| `appium:autoAcceptAlerts` (iOS) | dismiss native permission dialogs automatically |
| `appium:autoGrantPermissions` (Android) | grant runtime permissions at install |
| `appium:printPageSourceOnFindFailure` | helpful for debugging xpath fails |

---

## 4. Locator strategies

Order of preference (fastest + most stable first):

1. **Accessibility ID** — works on both platforms, single locator string.
2. **iOS predicate string** / **iOS class chain** for iOS.
3. **UiSelector** for Android.
4. **xpath** — last resort. Slow on big trees, brittle to layout changes.

### 4.1 Accessibility ID (preferred)

iOS: matches `accessibilityIdentifier` on the view.
Android: matches `content-desc` attribute.

```ts
// WebdriverIO uses ~ prefix for accessibility id
const loginButton = await $('~login-button');
await loginButton.click();
```

In your app code, set stable accessibility IDs in development:

```swift
// SwiftUI
Button("Log in") { ... }.accessibilityIdentifier("login-button")
```

```kotlin
// Jetpack Compose
Button(onClick = { ... }, modifier = Modifier.semantics { contentDescription = "login-button" }) { ... }
```

### 4.2 iOS predicate string

Powerful, fast, supports compound conditions:

```ts
// type and label
const cell = await $('-ios predicate string:type == "XCUIElementTypeCell" AND label CONTAINS "Inbox"');

// visible only
const btn = await $('-ios predicate string:name == "Submit" AND visible == 1');
```

### 4.3 iOS class chain

XPath-like but native and faster:

```ts
const firstCell = await $('-ios class chain:**/XCUIElementTypeTable/XCUIElementTypeCell[1]');
```

### 4.4 Android UiSelector

```ts
const tile = await $(
  'android=new UiSelector().className("android.widget.TextView").textContains("Inbox")'
);
```

Scrollable variant (auto-scrolls to find):

```ts
const item = await $(
  'android=new UiScrollable(new UiSelector().scrollable(true))' +
    '.scrollIntoView(new UiSelector().text("Settings"))'
);
```

### 4.5 xpath (last resort)

```ts
// Slow — avoid when possible
const banner = await $('//XCUIElementTypeStaticText[@name="Welcome"]');
const banner2 = await $('//android.widget.TextView[@text="Welcome"]');
```

Why xpath is slow: Appium serializes the full UI tree for each query. On
complex screens (lists, tabs) this can take seconds.

---

## 5. Page Object Model

Goal: encapsulate selectors and actions so tests read like business intent.

### 5.1 Base page

```ts
// test/pageobjects/BasePage.ts
import type { ChainablePromiseElement } from 'webdriverio';

export abstract class BasePage {
  /**
   * Wait for the page to be ready. Each subclass must define what "ready"
   * means — typically a unique anchor element being displayed.
   */
  abstract isReady(): Promise<void>;

  protected async waitFor(
    el: ChainablePromiseElement,
    timeout = 10_000
  ): Promise<void> {
    await el.waitForDisplayed({ timeout });
  }

  protected async tap(el: ChainablePromiseElement): Promise<void> {
    await this.waitFor(el);
    await el.click();
  }

  protected async type(
    el: ChainablePromiseElement,
    text: string
  ): Promise<void> {
    await this.waitFor(el);
    await el.setValue(text);
  }
}
```

### 5.2 Concrete page

```ts
// test/pageobjects/LoginPage.ts
import { BasePage } from './BasePage';

export class LoginPage extends BasePage {
  private get emailField() {
    return $('~login-email');
  }
  private get passwordField() {
    return $('~login-password');
  }
  private get submitButton() {
    return $('~login-submit');
  }
  private get errorBanner() {
    return $('~login-error');
  }

  async isReady(): Promise<void> {
    await this.waitFor(this.submitButton);
  }

  async loginWith(email: string, password: string): Promise<void> {
    await this.type(this.emailField, email);
    await this.type(this.passwordField, password);
    await this.tap(this.submitButton);
  }

  async getErrorText(): Promise<string> {
    await this.waitFor(this.errorBanner);
    return this.errorBanner.getText();
  }
}
```

### 5.3 Init pattern from a spec

```ts
// test/specs/login.spec.ts
import { LoginPage } from '../pageobjects/LoginPage';

describe('Login', () => {
  const loginPage = new LoginPage();

  beforeEach(async () => {
    await loginPage.isReady();
  });

  it('rejects bad credentials', async () => {
    await loginPage.loginWith('user@acme.com', 'wrong');
    await expect(await loginPage.getErrorText()).toContain('Invalid');
  });
});
```

Rules of thumb:

- Selectors live in the page, never in the spec.
- One assertion target per page method when possible.
- Pages return primitive data (strings, booleans, DTOs), not WDIO elements.
- Pages do not call `expect` — that's the spec's job.

---

## 6. Gestures

### 6.1 Simple gestures via WebdriverIO

```ts
const btn = await $('~login-submit');
await btn.click();             // tap
await btn.doubleClick();       // double tap (WDIO v9)
await driver.hideKeyboard();
await driver.setOrientation('LANDSCAPE');
await driver.setOrientation('PORTRAIT');
```

### 6.2 W3C Actions API (preferred for complex gestures)

```ts
// Swipe up by 60% of screen height
async function swipeUp() {
  const { width, height } = await driver.getWindowSize();
  const startX = Math.floor(width / 2);
  const startY = Math.floor(height * 0.8);
  const endY = Math.floor(height * 0.2);

  await driver.performActions([
    {
      type: 'pointer',
      id: 'finger1',
      parameters: { pointerType: 'touch' },
      actions: [
        { type: 'pointerMove', duration: 0, x: startX, y: startY },
        { type: 'pointerDown', button: 0 },
        { type: 'pause', duration: 100 },
        { type: 'pointerMove', duration: 600, x: startX, y: endY },
        { type: 'pointerUp', button: 0 },
      ],
    },
  ]);
  await driver.releaseActions();
}
```

### 6.3 Driver-specific mobile commands

The XCUITest and UiAutomator2 drivers expose dedicated `mobile:` commands —
faster and less flaky than W3C Actions:

```ts
// iOS
await driver.execute('mobile: swipe', { direction: 'up' });
await driver.execute('mobile: scroll', { direction: 'down' });
await driver.execute('mobile: pinch', { scale: 0.5, velocity: 1.0 });
await driver.execute('mobile: tap', { x: 100, y: 200 });
await driver.execute('mobile: doubleTap', { elementId: (await $('~zoomable')).elementId });

// Android
await driver.execute('mobile: swipeGesture', {
  left: 100,
  top: 800,
  width: 500,
  height: 600,
  direction: 'up',
  percent: 0.75,
});
await driver.execute('mobile: longClickGesture', {
  elementId: (await $('~item')).elementId,
  duration: 1500,
});
await driver.execute('mobile: scrollGesture', {
  left: 0,
  top: 200,
  width: 1080,
  height: 1500,
  direction: 'down',
  percent: 1.0,
});
```

These commands are documented in the driver READMEs:
- https://github.com/appium/appium-xcuitest-driver
- https://github.com/appium/appium-uiautomator2-driver

---

## 7. Waits

### 7.1 Implicit wait — avoid as primary strategy

Implicit waits apply to every command and silently slow your suite. They also
hide flakiness instead of surfacing it. Use them only as a small safety net
(0–1s) or not at all.

```ts
// If you must, set once, low value
await driver.setTimeout({ implicit: 0 });
```

### 7.2 Explicit waits (preferred)

WebdriverIO ships predicates on every element:

```ts
const btn = await $('~login-submit');

await btn.waitForExist({ timeout: 10_000 });
await btn.waitForDisplayed({ timeout: 10_000 });
await btn.waitForClickable({ timeout: 10_000 });
await btn.waitForEnabled({ timeout: 10_000 });
```

### 7.3 `waitUntil` with custom predicate

For application-state conditions that don't map to a single element:

```ts
await browser.waitUntil(
  async () => {
    const items = await $$('~todo-item');
    return items.length >= 3;
  },
  {
    timeout: 8_000,
    interval: 250,
    timeoutMsg: 'Expected 3 todo items to appear',
  }
);
```

### 7.4 Retry strategy

Wrap genuinely flaky interactions in a retry, but ONLY after checking that the
flake isn't caused by missing waits:

```ts
async function retry<T>(fn: () => Promise<T>, attempts = 3, delayMs = 500): Promise<T> {
  let lastErr: unknown;
  for (let i = 0; i < attempts; i++) {
    try { return await fn(); } catch (err) { lastErr = err; await browser.pause(delayMs); }
  }
  throw lastErr;
}
await retry(() => $('~submit').click());
```

For runner-level retries use Mocha's `this.retries(2)` or WDIO's
`mochaOpts.retries`. Always log retried tests so you can hunt the root cause.

---

## 8. Parallel execution

WebdriverIO is the easiest path to parallel mobile runs.

### 8.1 Multi-capability config

```ts
// wdio.conf.ts (excerpt)
export const config: WebdriverIO.Config = {
  port: 4723, // ignored when each cap defines its own server, see below

  // Run up to 2 sessions in parallel
  maxInstances: 2,

  capabilities: [
    {
      // Device 1 — iPhone Simulator
      protocol: 'http',
      hostname: '127.0.0.1',
      port: 4723,
      path: '/',
      platformName: 'iOS',
      'appium:automationName': 'XCUITest',
      'appium:deviceName': 'iPhone 15',
      'appium:platformVersion': '17.4',
      'appium:app': process.env.IOS_APP_PATH,
    },
    {
      // Device 2 — Android Emulator on a separate Appium server
      protocol: 'http',
      hostname: '127.0.0.1',
      port: 4724,
      path: '/',
      platformName: 'Android',
      'appium:automationName': 'UiAutomator2',
      'appium:deviceName': 'Pixel_7_API_34',
      'appium:platformVersion': '14',
      'appium:app': process.env.ANDROID_APP_PATH,
      'appium:systemPort': 8201, // unique per parallel Android session
    },
  ],

  framework: 'mocha',
  specs: ['./test/specs/**/*.spec.ts'],
  // ...
};
```

### 8.2 Run each Appium server on its own port

For two real devices in parallel:

```bash
# Terminal 1
appium --port 4723 --base-path /

# Terminal 2
appium --port 4724 --base-path /
```

Then point each capability at the matching port, plus a unique
`appium:systemPort` (Android) or `appium:wdaLocalPort` and
`appium:mjpegServerPort` (iOS) so the per-device helpers don't collide.

| Cap | Allocate per session |
|-----|----------------------|
| `appium:systemPort` (Android) | yes |
| `appium:chromedriverPort` (Android webview) | yes |
| `appium:wdaLocalPort` (iOS) | yes |
| `appium:mjpegServerPort` (iOS) | yes |

A small helper:

```ts
let portCounter = 0;
function nextPort(base: number) {
  return base + portCounter++;
}
```

For larger fleets, use a device farm (BrowserStack, SauceLabs, AWS Device Farm)
or `appium-device-farm` plugin instead of hand-allocating.

---

## 9. Real device vs Simulator / Emulator

| Concern | Simulator / Emulator | Real device |
|---------|----------------------|-------------|
| Speed of feedback | Fast (no signing, no provisioning) | Slow first run, fast after |
| Cost | Free | Hardware + lab |
| Push notifications | Limited (iOS Simulator: APNs simulator only; Android emulator: yes) | Full fidelity |
| Biometrics (FaceID, fingerprint) | Mockable (simulator commands) | Requires UI helpers |
| Camera, GPS, Bluetooth | Stubbed/limited | Full |
| Performance characteristics | NOT representative | Representative |
| Signing constraints | None | Provisioning profile, USB cable, dev mode |
| Debugging access | `xcrun simctl`, `adb` against emulator | `idevicesyslog`, `adb logcat` |
| Timing | Often faster than real | Source of truth |
| Flake risk from device state | Low | Higher (battery, thermal, network) |

Practical rule:

- **Dev loop and PR CI**: simulators / emulators. Cheap and parallel.
- **Nightly / pre-release**: real devices, ideally in a device farm, on the
  top 3–5 OS + form-factor combos that match your analytics.
- **Performance, sensors, push**: always real device.

---

## 10. Common pitfalls

### 10.1 Keyboard occlusion

Soft keyboard hides the field you just typed in, so a follow-up tap on a
button below misses.

```ts
// Always dismiss before tapping next control
await driver.hideKeyboard();
// iOS: tap "return" key on the keyboard if hideKeyboard is unreliable
await driver.execute('mobile: keys', { keys: ['\n'] });
```

### 10.2 Scroll-to-element

Don't blindly swipe. Use the platform helpers:

```ts
// iOS
await driver.execute('mobile: scroll', {
  elementId: (await $('~settings-row')).elementId,
  toVisible: true,
});

// Android — UiScrollable does the work
const row = await $(
  'android=new UiScrollable(new UiSelector().scrollable(true))' +
    '.scrollIntoView(new UiSelector().description("settings-row"))'
);
```

### 10.3 Native alerts

Permission dialogs steal focus and break locators. Two strategies:

```ts
// 1. Auto-dismiss at session level (preferred)
'appium:autoAcceptAlerts': true,        // iOS
'appium:autoGrantPermissions': true,    // Android

// 2. Handle on demand
await driver.acceptAlert();
await driver.dismissAlert();
const text = await driver.getAlertText();
```

Be aware: system-level alerts (App Tracking Transparency on iOS, Google Play
Protect on Android) sometimes need vendor-specific commands or a paired
WebDriverAgent build.

### 10.4 Deep links

Trigger them via the driver instead of building UI flows:

```ts
// iOS — uses xcrun simctl on simulators
await driver.execute('mobile: deepLink', {
  url: 'myapp://orders/42',
  bundleId: 'com.acme.myapp',
});

// Android
await driver.execute('mobile: deepLink', {
  url: 'myapp://orders/42',
  package: 'com.acme.myapp',
});
```

### 10.5 App reset between tests

Three knobs interact, and getting them wrong is a top source of flake:

| Cap | iOS behavior | Android behavior |
|-----|--------------|------------------|
| `appium:noReset: true` | Keep app data between sessions | Keep app data between sessions |
| `appium:noReset: false` (default) | Clear app state but keep app installed | Clear app state but keep app installed |
| `appium:fullReset: true` | Uninstall + reinstall app | Uninstall + reinstall app |

Recommended:

- Per-suite: `noReset: true` + an explicit `driver.terminateApp` /
  `activateApp` per test for speed.
- Per-test isolation: `noReset: false`. Slower but bulletproof.
- `fullReset` only for the very first run of CI to guarantee a clean install.

```ts
// Programmatic reset between tests
afterEach(async () => {
  await driver.terminateApp('com.acme.myapp');
  await driver.activateApp('com.acme.myapp');
});
```

### 10.6 Flake patterns and how to fix them

| Pattern | Likely cause | Fix |
|---------|--------------|-----|
| "Element not found" only on CI | CI machine slower; default 60s timeout too short for cold start | Bump `appium:newCommandTimeout`, `appium:wdaLaunchTimeout`; add `waitForDisplayed` |
| Test passes locally, fails in parallel | Shared device / port collision | Allocate `systemPort`, `wdaLocalPort` per session |
| Stale element after navigation | Element handle from previous screen | Re-query inside the page object getter (don't cache) |
| Random taps miss the target | Animation still running | Wait for `waitForClickable` or animation-end event, not just `Displayed` |
| "Could not connect to WDA" (iOS) | Stale WDA process from previous run | `xcrun simctl shutdown all && killall -9 xcodebuild`, or `appium:useNewWDA: true` |
| Android keyboard blocks input | Soft keyboard | `driver.hideKeyboard()` after every input |

---

## 11. Code examples

### 11.1 TypeScript + WebdriverIO v9 — full small spec

Project layout: `wdio.conf.ts` + `test/pageobjects/{BasePage,LoginPage}.ts` +
`test/specs/login.spec.ts`. `package.json` deps: `@wdio/cli@^9`,
`@wdio/local-runner@^9`, `@wdio/mocha-framework@^9`, `@wdio/spec-reporter@^9`,
`appium@^2.5`, `typescript@^5.4`. `tsconfig.json` should include
`@wdio/globals/types` in `compilerOptions.types`.

`wdio.conf.ts`:

```ts
import type { Options } from '@wdio/types';

export const config: Options.Testrunner = {
  runner: 'local',
  tsConfigPath: './tsconfig.json',

  specs: ['./test/specs/**/*.spec.ts'],
  maxInstances: 1,

  capabilities: [
    {
      platformName: 'iOS',
      'appium:automationName': 'XCUITest',
      'appium:deviceName': 'iPhone 15',
      'appium:platformVersion': '17.4',
      'appium:app': process.env.IOS_APP_PATH,
      'appium:newCommandTimeout': 240,
      'appium:autoAcceptAlerts': true,
    },
  ],

  logLevel: 'info',
  framework: 'mocha',
  reporters: ['spec'],
  mochaOpts: { ui: 'bdd', timeout: 120_000, retries: 1 },

  services: [
    [
      'appium',
      {
        args: { address: '127.0.0.1', port: 4723, basePath: '/' },
      },
    ],
  ],
};
```

`test/pageobjects/BasePage.ts` and `LoginPage.ts` follow the same shape shown
in section 5 (BasePage abstract + tap/type/waitFor helpers). The spec-specific
LoginPage adds an extra `homeAnchor` for post-login verification:

```ts
import { BasePage } from './BasePage';

export class LoginPage extends BasePage {
  private get emailField() {
    return $('~login-email');
  }
  private get passwordField() {
    return $('~login-password');
  }
  private get submitButton() {
    return $('~login-submit');
  }
  private get errorBanner() {
    return $('~login-error');
  }
  private get homeAnchor() {
    return $('~home-tab');
  }

  async isReady(): Promise<void> {
    await this.waitFor(this.submitButton);
  }

  async loginWith(email: string, password: string): Promise<void> {
    await this.type(this.emailField, email);
    await this.type(this.passwordField, password);
    await driver.hideKeyboard();
    await this.tap(this.submitButton);
  }

  async waitForHome(): Promise<void> {
    await this.waitFor(this.homeAnchor, 15_000);
  }

  async getErrorText(): Promise<string> {
    await this.waitFor(this.errorBanner);
    return this.errorBanner.getText();
  }
}
```

`test/specs/login.spec.ts`:

```ts
import { LoginPage } from '../pageobjects/LoginPage';

describe('Login flow', () => {
  const loginPage = new LoginPage();

  beforeEach(async () => {
    await driver.terminateApp('com.acme.myapp');
    await driver.activateApp('com.acme.myapp');
    await loginPage.isReady();
  });

  it('lets a valid user in', async () => {
    await loginPage.loginWith('user@acme.com', 'CorrectHorse9!');
    await loginPage.waitForHome();
  });

  it('shows an inline error for bad credentials', async () => {
    await loginPage.loginWith('user@acme.com', 'wrong');
    const msg = await loginPage.getErrorText();
    expect(msg).toContain('Invalid');
  });
});
```

Run:

```bash
export IOS_APP_PATH="$PWD/build/MyApp.app"
npm test
```

### 11.2 Java + Selenium-Appium — alternative

For teams already invested in Java/JUnit, Appium ships an official Java
client built on Selenium 4.

`pom.xml` (excerpt):

```xml
<dependencies>
  <dependency>
    <groupId>io.appium</groupId>
    <artifactId>java-client</artifactId>
    <version>9.2.2</version>
  </dependency>
  <dependency>
    <groupId>org.junit.jupiter</groupId>
    <artifactId>junit-jupiter</artifactId>
    <version>5.10.2</version>
    <scope>test</scope>
  </dependency>
</dependencies>
```

`AndroidLoginTest.java`:

```java
package com.acme.tests;

import io.appium.java_client.AppiumBy;
import io.appium.java_client.android.AndroidDriver;
import io.appium.java_client.android.options.UiAutomator2Options;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.openqa.selenium.WebElement;
import org.openqa.selenium.support.ui.WebDriverWait;
import org.openqa.selenium.support.ui.ExpectedConditions;

import java.net.URL;
import java.time.Duration;

import static org.junit.jupiter.api.Assertions.assertTrue;

class AndroidLoginTest {

  private AndroidDriver driver;
  private WebDriverWait wait;

  @BeforeEach
  void setUp() throws Exception {
    UiAutomator2Options options = new UiAutomator2Options()
        .setDeviceName("Pixel_7_API_34")
        .setPlatformVersion("14")
        .setApp(System.getenv("ANDROID_APP_PATH"))
        .setAutoGrantPermissions(true)
        .setNewCommandTimeout(Duration.ofMinutes(4));

    driver = new AndroidDriver(new URL("http://127.0.0.1:4723/"), options);
    wait = new WebDriverWait(driver, Duration.ofSeconds(10));
  }

  @AfterEach
  void tearDown() {
    if (driver != null) driver.quit();
  }

  @Test
  void rejectsBadCredentials() {
    WebElement email = wait.until(ExpectedConditions.visibilityOfElementLocated(
        AppiumBy.accessibilityId("login-email")));
    email.sendKeys("user@acme.com");

    driver.findElement(AppiumBy.accessibilityId("login-password"))
        .sendKeys("wrong");
    driver.hideKeyboard();
    driver.findElement(AppiumBy.accessibilityId("login-submit")).click();

    WebElement banner = wait.until(ExpectedConditions.visibilityOfElementLocated(
        AppiumBy.accessibilityId("login-error")));
    assertTrue(banner.getText().contains("Invalid"));
  }
}
```

Run with Maven:

```bash
export ANDROID_APP_PATH="$PWD/build/app-debug.apk"
mvn -q test
```

---

## Version pinning recap

Appium 2.5.x, XCUITest driver 7.x, UiAutomator2 driver 3.x, WebdriverIO v9
(`@wdio/cli@^9`), appium java-client 9.2.x (Selenium 4.21+), Node 18/20 LTS,
JDK 17. Pin these in `package.json` / `pom.xml` and treat upgrades as their
own PR with a green CI run as the gate.
