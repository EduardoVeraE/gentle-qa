<!-- Skill: qa-mobile-testing · Template: device-matrix -->
<!-- Placeholders: {{project}}, {{date}}, {{cloud_farm}}, {{primary_market}}, {{secondary_markets}} -->

# Device Matrix

> Project: {{project}} · Date: {{date}} · Cloud farm: {{cloud_farm}}

Tiering rule of thumb:

- **Tier 1 (T1)**: release-blocking. Covers ≥ 70% of installed base in {{primary_market}}.
- **Tier 2 (T2)**: best-effort, run weekly. Covers long-tail OEMs and older OS.
- **Tier 3 (T3)**: spot-check, run on regression milestones.

---

## Tier 1 — Release Blocking

| Device                | OS Version       | RAM   | Screen     | Network    | Coverage Rationale                          | Cloud Farm Name        |
| --------------------- | ---------------- | ----- | ---------- | ---------- | ------------------------------------------- | ---------------------- |
| iPhone 15 Pro         | iOS 17.x         | 8 GB  | 6.1" 460ppi | Wi-Fi/5G  | Latest flagship, dynamic island, A17        | {{cloud_farm}}-ios-15p |
| iPhone 13             | iOS 17.x         | 4 GB  | 6.1" 460ppi | Wi-Fi/5G  | High install base, no dynamic island        | {{cloud_farm}}-ios-13  |
| iPhone SE (3rd gen)   | iOS 16.x         | 4 GB  | 4.7" 326ppi | Wi-Fi/4G  | Smallest supported screen, TouchID          | {{cloud_farm}}-ios-se3 |
| Pixel 8               | Android 14       | 8 GB  | 6.2" 428ppi | Wi-Fi/5G  | AOSP reference, latest Android              | {{cloud_farm}}-and-px8 |
| Samsung Galaxy S23    | Android 14       | 8 GB  | 6.1" 425ppi | Wi-Fi/5G  | OneUI flagship, Knox, Samsung keyboard      | {{cloud_farm}}-and-s23 |
| Samsung Galaxy A54    | Android 13       | 6 GB  | 6.4" 403ppi | Wi-Fi/4G  | Mid-range, dominant in {{secondary_markets}} | {{cloud_farm}}-and-a54 |

---

## Tier 2 — Best Effort

| Device              | OS Version | RAM  | Screen     | Network   | Coverage Rationale                  | Cloud Farm Name        |
| ------------------- | ---------- | ---- | ---------- | --------- | ----------------------------------- | ---------------------- |
| iPhone 12 mini      | iOS 16.x   | 4 GB | 5.4" 476ppi | Wi-Fi/5G | Smallest modern iPhone              | {{cloud_farm}}-ios-12m |
| iPad (10th gen)     | iPadOS 17  | 4 GB | 10.9"      | Wi-Fi    | Tablet layout regressions           | {{cloud_farm}}-ios-ip10 |
| Xiaomi Redmi Note 12 | Android 13 | 4 GB | 6.67" 395ppi | Wi-Fi/4G | MIUI quirks, large APAC/LatAm share | {{cloud_farm}}-and-rn12 |
| OnePlus 11          | Android 13 | 8 GB | 6.7" 525ppi | Wi-Fi/5G | OxygenOS, gaming-tier perf          | {{cloud_farm}}-and-op11 |
| Motorola G54        | Android 13 | 4 GB | 6.5" 405ppi | Wi-Fi/4G | Stock-ish Android, low-mid market   | {{cloud_farm}}-and-g54 |

---

## Tier 3 — Spot Check

| Device              | OS Version | RAM  | Screen | Network  | Coverage Rationale                  | Cloud Farm Name         |
| ------------------- | ---------- | ---- | ------ | -------- | ----------------------------------- | ----------------------- |
| iPhone 11           | iOS 15.x   | 4 GB | 6.1"   | Wi-Fi/4G | Min supported iOS                   | {{cloud_farm}}-ios-11   |
| Pixel 5a            | Android 13 | 6 GB | 6.34"  | Wi-Fi/5G | Older AOSP                          | {{cloud_farm}}-and-px5a |
| Galaxy Tab S8       | Android 13 | 8 GB | 11"    | Wi-Fi    | Android tablet split layout         | {{cloud_farm}}-and-tabs8 |
| Foldable (Z Fold 5) | Android 14 | 12 GB | 7.6" inner | Wi-Fi/5G | Foldable continuity / multi-window  | {{cloud_farm}}-and-fold5 |

---

## Coverage Axes (orthogonal to device tier)

These axes must each be exercised on at least one T1 device.

| Axis           | Variants                              | Notes                                           |
| -------------- | ------------------------------------- | ----------------------------------------------- |
| Orientation    | Portrait, Landscape, Auto-rotate off  | Verify navigation + modals on rotation          |
| Locale         | en-US, {{primary_market}}, RTL (ar)   | Strings, date/number format, layout mirroring   |
| Dark mode      | Light, Dark, Follow system            | Contrast, asset variants, status bar tinting    |
| Dynamic Type   | Default, Largest accessibility size   | iOS only; truncation and overflow checks        |
| Font scale     | 100%, 130%, 200% (Android)            | Same as above for Android                       |
| Network        | Wi-Fi, 4G, 5G, Offline, Flaky 1% loss | Use Network Link Conditioner / Charles          |
| Battery        | 100%, < 20% low-power mode            | Background sync throttling                      |
| Storage        | Plenty, near-full (< 500 MB free)     | Install and cache eviction paths                |
