# Detox Reference (v20.x)

Grey-box end-to-end testing for React Native applications. This reference targets Detox 20.x, the current major line maintained by Wix.

---

## 1. Introduction

Detox is a grey-box end-to-end test runner built specifically for React Native. Unlike pure black-box tools (Appium, XCUITest, Espresso used directly), Detox is compiled INTO the app under test as a native dependency and synchronizes with the JavaScript bridge, the Animated API, timers, and pending network requests.

### Black-box vs Grey-box

| Approach  | Tool        | Visibility into app internals | Sync model                    |
| --------- | ----------- | ----------------------------- | ----------------------------- |
| Black-box | Appium      | None (UI only)                | Polling + arbitrary `sleep()` |
| Grey-box  | Detox       | JS bridge, timers, network    | Automatic idle detection      |
| White-box | Jest + RNTL | Component tree, props, state  | Synchronous (no real device)  |

Detox knows when the app is "idle" — no pending fetches, no running animations, no scheduled timers — and waits before each action. This eliminates the single biggest source of flake in mobile E2E: arbitrary waits.

### When Detox beats Appium for RN apps

- **Synchronization**: Detox does not need `sleep(2000)` or `waitForElementToBeVisible(timeout)`. The test framework knows the app's state.
- **Speed**: a typical Detox spec runs 2-4x faster than the equivalent Appium spec because it does not poll the screen for elements.
- **Stability**: flake rates of <1% are achievable on CI; Appium typically sits at 5-15% on the same suite.
- **Dev ergonomics**: tests are JavaScript/TypeScript, the same language as the app. No Java/Python bridge.

### When Appium is still the right choice

- The app is **not** React Native (native iOS/Android, Flutter, hybrid Cordova/Ionic).
- You need **cross-platform** tests covering both an RN app and a native one with the same suite.
- You do not have access to the native build pipeline (e.g. testing a third-party APK without source).
- You need **real device clouds** with broad device coverage (BrowserStack, Sauce Labs) — Detox supports this through `genycloud` and EAS, but Appium has wider device coverage.

### Honest limitations

- **React Native only**. There is no Detox-for-Flutter or Detox-for-native-iOS. If your app is not RN, stop reading.
- **Native build access required**. You must be able to `pod install`, run Gradle, and ship a custom build of the app with the Detox dependency linked. Closed-source apps cannot be tested.
- **Slower iteration than Jest unit tests**. A Detox suite still has to launch a simulator, install an app, and run tests over a real bridge. Expect ~30s-2min per spec file. Reserve Detox for critical user journeys; use Jest + React Native Testing Library for components.
- **iOS Simulator and Android Emulator only by default**. Real-device support is improving but still rough; most teams run on simulators in CI.
- **Hermes quirks**. Some sync hooks behave differently on Hermes vs JSC (more on this below).

---

## 2. Setup

Detox requires a working React Native project. Bare RN works out of the box; Expo requires the bare workflow (`expo prebuild`) or Expo's dev-client.

### Install

```bash
npm install detox --save-dev
npm install jest @types/jest --save-dev   # if not already present
npm install --save-dev @config-plugins/detox  # Expo only
```

Globally install the CLI for ergonomics:

```bash
npm install -g detox-cli
```

### Initialize

```bash
npx detox init -r jest
```

This creates:

```
.detoxrc.js          # Detox configuration (apps + devices + configs)
e2e/
├── jest.config.js   # Jest config for e2e (separate from unit tests)
├── starter.test.js  # Example spec
└── tsconfig.json    # If TypeScript was detected
```

### iOS native setup

Add to `ios/Podfile`:

```ruby
target 'YourApp' do
  # ...existing config...

  # Detox is automatically linked through autolinking.
  # No manual entry required for v20.x.
end
```

Then:

```bash
cd ios && pod install && cd ..
```

If you use Hermes, ensure `:hermes_enabled => true` in your Podfile is consistent across debug and release configurations — mismatches cause launch hangs.

### Android native setup

In `android/build.gradle`:

```gradle
buildscript {
  ext {
    minSdkVersion = 24       // Detox requires API 24+
    kotlinVersion = '1.8.0'
  }
}

allprojects {
  repositories {
    google()
    mavenCentral()
  }
}
```

In `android/app/build.gradle`:

```gradle
android {
  defaultConfig {
    testBuildType System.getProperty('testBuildType', 'debug')
    testInstrumentationRunner 'androidx.test.runner.AndroidJUnitRunner'
  }

  testOptions {
    animationsDisabled = true
  }
}

dependencies {
  androidTestImplementation('com.wix:detox:+')
}
```

Create `android/app/src/androidTest/java/com/yourapp/DetoxTest.java`:

```java
package com.yourapp;

import com.wix.detox.Detox;
import com.wix.detox.config.DetoxConfig;

import org.junit.Rule;
import org.junit.Test;
import org.junit.runner.RunWith;

import androidx.test.ext.junit.runners.AndroidJUnit4;
import androidx.test.filters.LargeTest;
import androidx.test.rule.ActivityTestRule;

@RunWith(AndroidJUnit4.class)
@LargeTest
public class DetoxTest {
    @Rule
    public ActivityTestRule<MainActivity> mActivityRule = new ActivityTestRule<>(MainActivity.class, false, false);

    @Test
    public void runDetoxTests() {
        DetoxConfig detoxConfig = new DetoxConfig();
        detoxConfig.idlePolicyConfig.masterTimeoutSec = 90;
        detoxConfig.idlePolicyConfig.idleResourceTimeoutSec = 60;
        detoxConfig.rnContextLoadTimeoutSec = (BuildConfig.DEBUG ? 180 : 60);
        Detox.runTests(mActivityRule, detoxConfig);
    }
}
```

---

## 3. Configuration (`.detoxrc.js`)

The config file has three sections that compose into named configurations: **apps**, **devices**, and **configurations**.

```js
/** @type {Detox.DetoxConfig} */
module.exports = {
  testRunner: {
    args: {
      $0: 'jest',
      config: 'e2e/jest.config.js',
    },
    jest: {
      setupTimeout: 120000,
    },
  },

  apps: {
    'ios.debug': {
      type: 'ios.app',
      binaryPath: 'ios/build/Build/Products/Debug-iphonesimulator/YourApp.app',
      build:
        'xcodebuild -workspace ios/YourApp.xcworkspace -scheme YourApp -configuration Debug -sdk iphonesimulator -derivedDataPath ios/build',
    },
    'ios.release': {
      type: 'ios.app',
      binaryPath: 'ios/build/Build/Products/Release-iphonesimulator/YourApp.app',
      build:
        'xcodebuild -workspace ios/YourApp.xcworkspace -scheme YourApp -configuration Release -sdk iphonesimulator -derivedDataPath ios/build',
    },
    'android.debug': {
      type: 'android.apk',
      binaryPath: 'android/app/build/outputs/apk/debug/app-debug.apk',
      testBinaryPath:
        'android/app/build/outputs/apk/androidTest/debug/app-debug-androidTest.apk',
      build:
        'cd android && ./gradlew assembleDebug assembleAndroidTest -DtestBuildType=debug',
      reversePorts: [8081],
    },
    'android.release': {
      type: 'android.apk',
      binaryPath: 'android/app/build/outputs/apk/release/app-release.apk',
      testBinaryPath:
        'android/app/build/outputs/apk/androidTest/release/app-release-androidTest.apk',
      build:
        'cd android && ./gradlew assembleRelease assembleAndroidTest -DtestBuildType=release',
    },
  },

  devices: {
    simulator: {
      type: 'ios.simulator',
      device: { type: 'iPhone 15 Pro' },
    },
    attached: {
      type: 'android.attached',
      device: { adbName: '.*' },
    },
    emulator: {
      type: 'android.emulator',
      device: { avdName: 'Pixel_7_API_34' },
    },
    genycloud: {
      type: 'android.genycloud',
      device: { recipeUUID: '<recipe-uuid>' },
    },
  },

  configurations: {
    'ios.sim.debug':     { device: 'simulator', app: 'ios.debug' },
    'ios.sim.release':   { device: 'simulator', app: 'ios.release' },
    'android.emu.debug': { device: 'emulator',  app: 'android.debug' },
    'android.emu.release': { device: 'emulator', app: 'android.release' },
    'android.att.release': { device: 'attached', app: 'android.release' },
  },
};
```

**Naming convention**: `<platform>.<device-type>.<build-type>`. This is enforced by community convention, not Detox itself, but it pays off in CI scripts.

---

## 4. Building the app for tests

```bash
# iOS, debug
detox build --configuration ios.sim.debug

# Android, release
detox build --configuration android.emu.release
```

### Debug vs release

| Build   | When to use                              | Trade-off                                      |
| ------- | ---------------------------------------- | ---------------------------------------------- |
| Debug   | Local dev, fast rebuilds, Metro attached | Slower at runtime; flakier on CI               |
| Release | CI, smoke tests on PRs                   | Slower to build; closer to production behavior |

Always run **release** builds on CI. Debug builds depend on Metro running in the background; that race condition will bite you.

### iOS signing

Simulator builds do NOT need code signing. For release configs targeting the simulator, ensure `CODE_SIGNING_ALLOWED=NO` in the build command:

```bash
xcodebuild ... CODE_SIGNING_ALLOWED=NO
```

Real-device runs require a development team and provisioning profile. Use `fastlane match` to keep the team's certs in sync.

### Hermes considerations

- Hermes is the default JS engine on RN 0.70+. Detox 20.x supports Hermes natively, but the **first launch is slower** (Hermes precompiles bytecode).
- If you see `RN_BRIDGE_TIMEOUT` errors on the first test, increase `rnContextLoadTimeoutSec` in `DetoxTest.java` (Android) or `detoxIPCSetUp` config (iOS).
- Source maps are different in Hermes; if you need readable JS stack traces in the Detox logs, enable `hermesEnabled: false` for debug builds only.

---

## 5. Matchers

Matchers locate elements. Always combine with `element(...)` or `expect(element(...))`.

### `by.id` — preferred

```ts
element(by.id('login-button'));
```

Add `testID="login-button"` to the React Native component. `testID` is stable across translations and design changes; it is the ONLY matcher that should appear in production specs.

### `by.text`

```ts
element(by.text('Sign in'));
```

Brittle: breaks on i18n changes, A/B tests, copy edits. Acceptable only for assertions, never for actions.

### `by.label`

```ts
element(by.label('Submit form'));
```

Maps to `accessibilityLabel` on iOS and `contentDescription` on Android. Useful when accessibility labels are mandated, but still locale-sensitive if labels are translated.

### `by.type`

```ts
element(by.type('RCTImageView'));   // iOS native class name
element(by.type('android.widget.ImageView'));   // Android
```

Returns ALL views of that type. Combine with index (`.atIndex(0)`) or descendants. Avoid in app code; useful only for debugging.

### `by.traits` (iOS only)

```ts
element(by.traits(['button', 'selected']));
```

Maps to UIAccessibilityTraits. iOS-only; do not use in cross-platform suites.

### Composing matchers

```ts
// AND
element(by.id('cell').and(by.label('Item 3')));

// Ancestor / descendant
element(by.id('child').withAncestor(by.id('parent-list')));
element(by.id('parent-list').withDescendant(by.id('child')));

// Index disambiguation
element(by.text('Delete')).atIndex(1);   // second match
```

### Why `testID` is preferred

1. **Stable**: changes only when the developer explicitly changes it.
2. **Locale-independent**: not affected by i18n.
3. **Fast**: native code performs ID lookup in O(1); text/label lookup walks the view tree.
4. **Self-documenting**: a `testID="cart.checkout.confirm"` namespacing convention doubles as analytics IDs.

---

## 6. Actions

```ts
// Basic taps
await element(by.id('submit')).tap();
await element(by.id('item')).multiTap(2);
await element(by.id('item')).longPress();
await element(by.id('item')).longPress(2000);  // 2 seconds

// Custom-position tap (relative to element bounds)
await element(by.id('canvas')).tap({ x: 100, y: 50 });

// Swipe (direction, speed, percentage)
await element(by.id('list')).swipe('up', 'fast');
await element(by.id('list')).swipe('left', 'slow', 0.75);

// Scroll (offset-based)
await element(by.id('scrollview')).scroll(200, 'down');
await element(by.id('scrollview')).scroll(100, 'up', NaN, 0.85);  // start from 85% of view

// Scroll until element is visible
await waitFor(element(by.id('row-99')))
  .toBeVisible()
  .whileElement(by.id('list'))
  .scroll(300, 'down');

// Text input
await element(by.id('email')).typeText('user@example.com');
await element(by.id('email')).replaceText('new@example.com');
await element(by.id('email')).clearText();
await element(by.id('email')).tapReturnKey();
await element(by.id('search')).tapBackspaceKey();

// Pinch (iOS only)
await element(by.id('image')).pinch(0.5);   // zoom out
await element(by.id('image')).pinch(2.0);   // zoom in

// Screenshot
const path = await element(by.id('hero')).takeScreenshot('hero-state');
```

### Action timing

Detox waits for the app to be idle before executing each action. You typically do NOT need explicit `sleep`. The exceptions:

- **Animations driven by `requestAnimationFrame` outside React's scheduler** — Detox cannot detect these. Use `waitFor(...)` instead.
- **WebViews** — Detox does not see inside WebViews on the JS bridge. Use `device.disableSynchronization()` around WebView interactions, then re-enable.

---

## 7. Expectations

```ts
await expect(element(by.id('home'))).toBeVisible();
await expect(element(by.id('hidden'))).toBeNotVisible();
await expect(element(by.id('label'))).toExist();
await expect(element(by.id('label'))).toNotExist();
await expect(element(by.id('label'))).toHaveText('Welcome');
await expect(element(by.id('label'))).toHaveLabel('greeting');
await expect(element(by.id('input'))).toHaveValue('initial');
await expect(element(by.id('toggle'))).toHaveToggleValue(true);
await expect(element(by.id('cell'))).toBeFocused();
```

### `waitFor`

When idle-sync is not enough (animations, network responses arriving via WebSocket, splash screens):

```ts
// Timeout-based
await waitFor(element(by.id('toast')))
  .toBeVisible()
  .withTimeout(5000);

// Whilst performing an action
await waitFor(element(by.id('row-99')))
  .toBeVisible()
  .whileElement(by.id('list'))
  .scroll(300, 'down');
```

`waitFor` polls every ~100ms. Keep timeouts small (3-10s); a 30s timeout hides a real bug.

---

## 8. Device API

```ts
// Launch with options
await device.launchApp({
  newInstance: true,
  permissions: { camera: 'YES', location: 'inuse', notifications: 'YES' },
  url: 'myapp://product/42',
  launchArgs: { detoxPrintBusyIdleResources: 'YES' },
});

// Reload only the JS bundle (keeps native state)
await device.reloadReactNative();

// Send a deep link to a running app
await device.openURL({ url: 'myapp://settings' });

// Notifications
await device.sendUserNotification({
  trigger: { type: 'push' },
  title: 'New message',
  body: 'Hello',
});

// Network mocking — block list of URL patterns
await device.setURLBlacklist(['.*analytics\\.com.*']);
await device.clearURLBlacklist();

// Location
await device.setLocation(37.7749, -122.4194);

// Orientation
await device.setOrientation('landscape');
await device.setOrientation('portrait');

// Background / foreground
await device.sendToHome();
await device.launchApp({ newInstance: false });

// Shake gesture (Android only triggers dev menu)
await device.shake();

// System bars (iOS 13+)
await device.setStatusBar({ time: '12:00', batteryState: 'charged' });

// Disable / enable synchronization (use sparingly)
await device.disableSynchronization();
// ...interact with WebView or long animation...
await device.enableSynchronization();
```

### Permissions

iOS supports: `calendar`, `camera`, `contacts`, `health`, `homekit`, `location`, `medialibrary`, `microphone`, `motion`, `notifications`, `photos`, `reminders`, `siri`, `speech`, `faceid`. Values: `YES` | `NO` | `unset` | `inuse` | `always`.

Android handles permissions through ADB; Detox auto-grants known runtime permissions if `permissions` is set.

---

## 9. Test isolation

The number-one rule: **every test must start from a known state**. Tests that depend on each other will rot.

### `beforeEach` launch

```ts
describe('Cart flow', () => {
  beforeAll(async () => {
    await device.launchApp({ newInstance: true });
  });

  beforeEach(async () => {
    await device.reloadReactNative();
  });

  // ...specs...
});
```

`reloadReactNative()` is much faster than `launchApp({ newInstance: true })` (1-2s vs 5-15s) but only resets JS state. Native state — keychain, AsyncStorage, push tokens — survives.

### Resetting persistent state

For full isolation:

```ts
beforeEach(async () => {
  await device.launchApp({
    delete: true,           // reinstall the app
    newInstance: true,
    permissions: { notifications: 'YES' },
  });
});
```

`delete: true` reinstalls — slowest option, ~10-30s. Reserve for tests that depend on first-run behavior.

### Fixture data via deep links

Best pattern: expose a `myapp://__test__/seed?fixture=cart-with-3-items` deep link in debug builds only. The test then reduces to:

```ts
beforeEach(async () => {
  await device.launchApp({
    newInstance: true,
    url: 'myapp://__test__/seed?fixture=cart-with-3-items',
  });
});
```

This is faster than driving the UI to set up state, and it lives in the app code (typed, refactor-safe).

### AsyncStorage cleanup

Add a debug-only IPC handler:

```ts
// In your app, behind __DEV__:
import AsyncStorage from '@react-native-async-storage/async-storage';
import { Linking } from 'react-native';

Linking.addEventListener('url', ({ url }) => {
  if (url.startsWith('myapp://__test__/clear')) AsyncStorage.clear();
});
```

Then call `await device.openURL({ url: 'myapp://__test__/clear' })` in `beforeEach`.

---

## 10. Debugging

### Inspect the JS test code

```bash
detox test --configuration ios.sim.debug --inspect-brk
```

Opens a Node inspector on port 9229. Attach Chrome DevTools or VS Code.

### Screenshots on failure

```js
// e2e/jest.config.js
module.exports = {
  testEnvironment: 'detox/runners/jest/testEnvironment',
  reporters: ['detox/runners/jest/reporter'],
  globalSetup: 'detox/runners/jest/globalSetup',
  globalTeardown: 'detox/runners/jest/globalTeardown',
  rootDir: '..',
  testMatch: ['<rootDir>/e2e/**/*.test.{js,ts}'],
  setupFilesAfterEnv: ['./e2e/setup.ts'],
};
```

```bash
detox test -c ios.sim.debug --record-logs all --take-screenshots failing
```

Artifacts land under `artifacts/<configName>.<timestamp>/`.

### Video recording

```bash
detox test -c ios.sim.debug --record-videos failing
```

Videos are heavy; use `failing` (only when a test fails) rather than `all`.

### Logs

```bash
detox test -c ios.sim.debug --record-logs all
# artifacts/<run>/<spec>/<index>/{device.log, app.log}
```

The `device.log` contains the simulator/emulator system log; `app.log` contains your `console.log` output.

### Idle resource debugging

If tests time out at "Detox is busy":

```bash
detox test -c ios.sim.debug --loglevel trace
```

Look for `BusyIdleResources` entries. Common offenders: long-running timers (`setInterval`), unfinished fetch calls, looping animations.

---

## 11. CI integration

### Caching builds

Detox builds are slow (3-15min). Cache them by hashing the inputs that affect the binary:

```yaml
# .github/workflows/e2e.yml
- name: Cache iOS build
  uses: actions/cache@v4
  with:
    path: ios/build
    key: ios-${{ hashFiles('ios/Podfile.lock', 'package-lock.json', 'ios/**/*.swift', 'ios/**/*.m', 'ios/**/*.h') }}

- name: Cache Pods
  uses: actions/cache@v4
  with:
    path: ios/Pods
    key: pods-${{ hashFiles('ios/Podfile.lock') }}
```

Invalidating on JS-only changes is wasteful — those do not require a rebuild. Hash native files only.

### Parallel runs

**A single Detox configuration does NOT parallelize across workers**. The simulator is a single-tenant resource; you cannot run two tests on the same device in parallel.

To parallelize, use **shards**:

```bash
# Worker 1
detox test -c ios.sim.release --maxWorkers 1 \
  --testPathPattern '(login|cart)\.test\.ts'

# Worker 2 (different machine or different simulator)
detox test -c ios.sim.release --maxWorkers 1 \
  --testPathPattern '(checkout|profile)\.test\.ts'
```

Or with Jest's built-in sharding (Detox 20.10+):

```bash
detox test -c ios.sim.release --shard 1/4
detox test -c ios.sim.release --shard 2/4
detox test -c ios.sim.release --shard 3/4
detox test -c ios.sim.release --shard 4/4
```

Each shard runs its own simulator. On GitHub Actions, use a matrix:

```yaml
strategy:
  matrix:
    shard: [1, 2, 3, 4]
steps:
  - run: detox test -c ios.sim.release --shard ${{ matrix.shard }}/4
```

### Retries

Detox supports per-spec retries via Jest:

```js
// e2e/setup.ts
jest.retryTimes(2, { logErrorsBeforeRetry: true });
```

Retries hide flake. Use them only on CI as a safety net AND track flake rate; if a test retries more than once a week, fix the test, do not raise the retry count.

---

## 12. Code examples

A complete minimal setup for a small RN app.

### `.detoxrc.js`

```js
/** @type {Detox.DetoxConfig} */
module.exports = {
  testRunner: {
    args: { $0: 'jest', config: 'e2e/jest.config.js' },
    jest: { setupTimeout: 120000 },
  },
  apps: {
    'ios.release': {
      type: 'ios.app',
      binaryPath:
        'ios/build/Build/Products/Release-iphonesimulator/Sample.app',
      build:
        'xcodebuild -workspace ios/Sample.xcworkspace -scheme Sample -configuration Release -sdk iphonesimulator -derivedDataPath ios/build CODE_SIGNING_ALLOWED=NO',
    },
    'android.release': {
      type: 'android.apk',
      binaryPath: 'android/app/build/outputs/apk/release/app-release.apk',
      testBinaryPath:
        'android/app/build/outputs/apk/androidTest/release/app-release-androidTest.apk',
      build:
        'cd android && ./gradlew assembleRelease assembleAndroidTest -DtestBuildType=release',
    },
  },
  devices: {
    simulator: {
      type: 'ios.simulator',
      device: { type: 'iPhone 15 Pro' },
    },
    emulator: {
      type: 'android.emulator',
      device: { avdName: 'Pixel_7_API_34' },
    },
  },
  configurations: {
    'ios.sim.release':     { device: 'simulator', app: 'ios.release' },
    'android.emu.release': { device: 'emulator',  app: 'android.release' },
  },
};
```

### `e2e/jest.config.js`

```js
/** @type {import('jest').Config} */
module.exports = {
  rootDir: '..',
  testMatch: ['<rootDir>/e2e/**/*.test.ts'],
  testTimeout: 120000,
  maxWorkers: 1,
  globalSetup: 'detox/runners/jest/globalSetup',
  globalTeardown: 'detox/runners/jest/globalTeardown',
  reporters: ['detox/runners/jest/reporter'],
  testEnvironment: 'detox/runners/jest/testEnvironment',
  verbose: true,
  preset: 'ts-jest',
};
```

### `e2e/setup.ts`

```ts
import { device } from 'detox';

beforeAll(async () => {
  await device.launchApp({
    newInstance: true,
    permissions: { notifications: 'YES' },
  });
});

afterAll(async () => {
  await device.terminateApp();
});
```

### `e2e/login.test.ts`

```ts
import { by, device, element, expect, waitFor } from 'detox';

describe('Login flow', () => {
  beforeEach(async () => {
    await device.reloadReactNative();
  });

  it('rejects empty credentials', async () => {
    await element(by.id('login.submit')).tap();
    await expect(element(by.id('login.error'))).toBeVisible();
    await expect(element(by.id('login.error'))).toHaveText(
      'Email and password are required',
    );
  });

  it('signs in with valid credentials', async () => {
    await element(by.id('login.email')).typeText('user@example.com');
    await element(by.id('login.password')).typeText('correct-horse');
    await element(by.id('login.submit')).tap();

    await waitFor(element(by.id('home.greeting')))
      .toBeVisible()
      .withTimeout(5000);

    await expect(element(by.id('home.greeting'))).toHaveText('Hi, user');
  });

  it('shows backend error on 401', async () => {
    await device.launchApp({
      newInstance: true,
      url: 'sample://__test__/seed?fixture=auth-401',
    });

    await element(by.id('login.email')).typeText('user@example.com');
    await element(by.id('login.password')).typeText('whatever');
    await element(by.id('login.submit')).tap();

    await expect(element(by.id('login.error'))).toHaveText(
      'Invalid credentials',
    );
  });
});
```

### `e2e/cart.test.ts`

```ts
import { by, device, element, expect, waitFor } from 'detox';

describe('Cart', () => {
  beforeEach(async () => {
    await device.launchApp({
      newInstance: true,
      url: 'sample://__test__/seed?fixture=cart-with-3-items',
    });
  });

  it('removes an item by swiping left', async () => {
    await expect(element(by.id('cart.row.0'))).toBeVisible();
    await element(by.id('cart.row.0')).swipe('left', 'fast');
    await element(by.id('cart.row.0.delete')).tap();
    await expect(element(by.id('cart.row.0'))).toBeNotVisible();
  });

  it('scrolls to the last row and proceeds to checkout', async () => {
    await waitFor(element(by.id('cart.row.10')))
      .toBeVisible()
      .whileElement(by.id('cart.list'))
      .scroll(300, 'down');

    await element(by.id('cart.checkout')).tap();
    await expect(element(by.id('checkout.title'))).toBeVisible();
  });
});
```

### Running

```bash
# Build once
detox build --configuration ios.sim.release

# Run all specs
detox test --configuration ios.sim.release

# Run a single spec
detox test --configuration ios.sim.release e2e/login.test.ts

# Run a single test
detox test --configuration ios.sim.release \
  --testNamePattern 'signs in with valid credentials'
```

---

## Further reading

- Detox docs: https://wix.github.io/Detox/
- Migration guide v19 → v20: https://wix.github.io/Detox/docs/guide/migration
- Detox + EAS (Expo): https://docs.expo.dev/build-reference/e2e-tests/
- Sample apps: https://github.com/wix/Detox/tree/master/examples
