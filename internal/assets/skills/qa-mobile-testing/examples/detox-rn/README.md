# Detox 20 + React Native ~0.74

End-to-end testing scaffold for a React Native ~0.74 app using Detox 20.x and Jest, on Node 20+.

## Stack

- **Detox 20.20.x** — gray-box E2E for native React Native apps
- **Jest 29** — test runner with Detox's environment + reporter
- **iOS Simulator** (Xcode-managed) and **Android Emulator** (AVD)

## Prereqs

```bash
nvm use 20
npm install
brew tap wix/brew && brew install applesimutils   # iOS only
```

Make sure you have:

- Xcode 15+ with Command Line Tools (iOS)
- Android SDK + an AVD called `Pixel_8_API_34` (or change `.detoxrc.js`)
- A React Native ~0.74 project at the repo root (`ios/`, `android/`, `metro.config.js`)

## Build

```bash
# iOS Simulator debug build
npm run build:ios

# Android Emulator debug build
npm run build:android
```

The build commands wrap `xcodebuild` and `gradlew assembleDebug assembleAndroidTest`.

## Run

```bash
npm run test:ios
npm run test:android
```

## Layout

```
.detoxrc.js          apps + devices + configurations
e2e/
  jest.config.js     Jest config wired to Detox runners
  starter.test.js    Smoke: app launch + first navigation
  login.test.js      Login flow: success and failure paths
```

## Conventions

- Identify elements by `testID` prop in RN components — never by visible text.
- Use `device.reloadReactNative()` in `beforeEach` for fast, deterministic state reset.
- Use `waitFor(...).toBeVisible().withTimeout(N)` instead of `device.pause()`.
- Detox is gray-box and synchronous — no manual sleeps required.

## Troubleshooting

- `Could not find a connected device` (Android) → start the emulator first or set `--device-name` to a running AVD.
- `applesimutils` not found → `brew install applesimutils`.
- App build runs Metro automatically; if you see stale JS, delete `ios/build` and `android/app/build` and rebuild.
