# Appium 2 + WebdriverIO v9 + TypeScript

Minimal cross-platform mobile UI automation starter. Targets Node 20+, iOS Simulator and Android Emulator out of the box; the same suite runs on real devices via cloud farms by swapping capabilities.

## Stack

- **Appium 2.5.x** with `appium-xcuitest-driver` (iOS) and `appium-uiautomator2-driver` (Android)
- **WebdriverIO v9** with `@wdio/local-runner` and `@wdio/appium-service`
- **Mocha** with retries
- **TypeScript 5** (transpile-only via ts-node)

## Setup

```bash
nvm use 20
npm install
npx appium driver install xcuitest
npx appium driver install uiautomator2
```

Place built artifacts under `apps/`:

```
apps/MyApp.app          # iOS Simulator build (.app bundle)
apps/app-debug.apk      # Android debug APK
```

## Run

```bash
# iOS Simulator
PLATFORM=ios npm test

# Android Emulator
PLATFORM=android npm test
```

Capability overrides via env vars: `IOS_VERSION`, `IOS_DEVICE`, `ANDROID_VERSION`, `ANDROID_DEVICE`.

## Layout

```
test/
  pages/      PageObjects (one class per screen)
  specs/      Mocha specs that drive PageObjects
wdio.conf.ts  Capabilities + framework + services
```

## Conventions

- Use accessibility ids (`~id`) — set via `testID` (RN), `accessibilityIdentifier` (iOS), `contentDescription` (Android).
- Never hard-code waits. Use `waitForDisplayed`, `waitFor`, or WDIO's auto-wait.
- Specs only call PageObject methods — no direct selectors in specs.
- Retries are set to 1 in `mochaOpts.retries`. Quarantine flaky specs, do not raise the retry budget.

## Troubleshooting

- iOS WDA build fails → ensure Xcode CLI tools and a valid signing identity are present.
- Android `ANDROID_HOME` errors → export `ANDROID_HOME` and add `platform-tools` to `PATH`.
- Stale Appium session → `pkill -f appium` before re-running.
