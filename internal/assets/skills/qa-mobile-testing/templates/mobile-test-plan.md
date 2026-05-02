<!-- Skill: qa-mobile-testing · Template: mobile-test-plan -->
<!-- Placeholders: {{project}}, {{release}}, {{owner}}, {{date}}, {{approvers}}, {{platforms}}, {{ios_min_os}}, {{ios_max_os}}, {{android_min_api}}, {{android_max_api}}, {{features_in_scope}}, {{features_out_of_scope}}, {{ios_build_id}}, {{android_build_id}}, {{tooling_choice}}, {{tooling_reasoning}}, {{cloud_farm}}, {{schedule_start}}, {{schedule_end}}, {{report_cadence}}, {{stakeholders}}, {{test_data_source}} -->

# Mobile Test Plan

> Fill in remaining fields; placeholders pre-populate when generated from context.

## 1. Document Control

| Field      | Value             |
| ---------- | ----------------- |
| Project    | {{project}}       |
| Release    | {{release}}       |
| Owner      | {{owner}}         |
| Date       | {{date}}          |
| Status     | Draft             |
| Approvers  | {{approvers}}     |
| Platforms  | {{platforms}}     |

---

## 2. Scope

### 2.1 In Scope

- Features: {{features_in_scope}}
- Platforms: {{platforms}} (iOS {{ios_min_os}}–{{ios_max_os}}, Android API {{android_min_api}}–{{android_max_api}})
- Form factors: phone, tablet (where applicable)
- Build artifacts: iOS `{{ios_build_id}}`, Android `{{android_build_id}}`

### 2.2 Out of Scope

- {{features_out_of_scope}}
- Wearables / TV / Auto surfaces (unless explicitly listed)
- Carrier-specific provisioning

---

## 3. Device Matrix

The full device matrix lives in [`device-matrix.md`](./device-matrix.md). Tier 1 devices block release; Tier 2 are best-effort; Tier 3 are spot-checked.

| Tier | Coverage Goal             | Block Release |
| ---- | ------------------------- | :-----------: |
| T1   | 100% smoke + regression   | Yes           |
| T2   | Smoke + targeted features | No            |
| T3   | Spot-check only           | No            |

---

## 4. Test Types

| Type                  | Goal                                                | Approach                          |
| --------------------- | --------------------------------------------------- | --------------------------------- |
| Functional            | Verify feature behavior                             | Scripted + exploratory            |
| Install / Uninstall   | Fresh install, reinstall, app data wipe             | Manual + automated on T1          |
| Upgrade               | Migration from N-1 and N-2 versions                 | Automated build hop matrix        |
| Network               | Wi-Fi, 4G/5G, offline, flaky, captive portal        | Charles / Network Link Conditioner |
| Locale & RTL          | en, {{platforms}} required locales, RTL (ar/he)     | Locale switch + screenshot diff   |
| Gesture               | Tap, swipe, pinch, long-press, drag (see catalog)   | See `gesture-catalog.md`          |
| Accessibility-pointer | VoiceOver / TalkBack focus order, Dynamic Type, CT  | Manual + axe-mobile               |
| Orientation           | Portrait / landscape transitions                    | Automated rotate steps            |
| Background / Foreground | Suspend, resume, OS-killed restore                | Lifecycle scripts                 |
| Performance           | Cold start, scroll FPS, memory, battery             | Sampled on T1                     |

---

## 5. Test Environment

| Env             | Purpose                | Notes                                              |
| --------------- | ---------------------- | -------------------------------------------------- |
| Simulator (iOS) | Fast feedback dev loop | Xcode-managed, not for perf or biometrics          |
| Emulator (And.) | CI + dev               | AVD with Google APIs; Play services as needed      |
| Real devices    | Release-blocking       | Local rack + cloud farm `{{cloud_farm}}`           |
| Build source    | CI artifacts           | iOS `{{ios_build_id}}`, Android `{{android_build_id}}` |

---

## 6. Tooling

Selected: **{{tooling_choice}}**.

Reasoning: {{tooling_reasoning}}

| Concern             | Tool                                  |
| ------------------- | ------------------------------------- |
| UI automation       | {{tooling_choice}}                    |
| Reporting           | Allure / JUnit XML                    |
| Cloud devices       | {{cloud_farm}}                        |
| Network shaping     | Charles, Network Link Conditioner     |
| Crash & telemetry   | Crashlytics / Sentry                  |
| Accessibility audit | axe-mobile, Accessibility Inspector   |

---

## 7. Test Data

- Source: {{test_data_source}}
- Test accounts seeded via fixtures, isolated per environment.
- Sensitive data: synthetic only; never use prod PII on devices.

---

## 8. Risk Matrix

| Risk                                           | Likelihood | Impact | Mitigation                                  |
| ---------------------------------------------- | :--------: | :----: | ------------------------------------------- |
| Fragmented Android OEM behavior                | High       | Med    | Enforce T1 OEM coverage; cloud farm rotation |
| iOS biometric features fail on simulator       | High       | Low    | Real-device-only suite for FaceID/TouchID   |
| Flaky tests block CI                           | Med        | Med    | Quarantine + retry policy (max 2)           |
| Build provisioning expiry mid-cycle            | Low        | High   | Automated profile renewal alert             |
| Accessibility regression unnoticed             | Med        | High   | a11y suite gated on T1                      |

---

## 9. Schedule

| Milestone           | Date              |
| ------------------- | ----------------- |
| Test design freeze  | {{schedule_start}} |
| Smoke pass green    | {{schedule_start}} + 3d |
| Full regression     | {{schedule_end}} - 5d |
| Release sign-off    | {{schedule_end}}  |

---

## 10. Entry & Exit Criteria

### Entry

- Build deployed to TestFlight / internal track
- Smoke suite passes on at least one T1 iOS and one T1 Android device
- No P0 / P1 bugs open from prior cycle

### Exit

- 100% of T1 regression executed; ≥ 95% pass rate
- 0 open P0, ≤ 2 open P1 with documented workarounds
- Crash-free sessions ≥ 99.5% on staging
- Accessibility checks green on T1
- Sign-off from {{stakeholders}}

---

## 11. Reporting

- Cadence: {{report_cadence}}
- Channel: dashboard + Slack digest
- Final test summary uses the report template at exit.

---

## 12. Stakeholders

| Role               | Name              |
| ------------------ | ----------------- |
| Test Lead          | {{owner}}         |
| Product Owner      | {{stakeholders}}  |
| Engineering Lead   | {{stakeholders}}  |
| Release Manager    | {{stakeholders}}  |
