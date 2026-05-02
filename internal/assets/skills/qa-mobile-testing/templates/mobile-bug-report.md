<!-- Skill: qa-mobile-testing · Template: mobile-bug-report -->
<!-- Placeholders: {{title}}, {{severity}}, {{priority}}, {{device_make}}, {{device_model}}, {{os_name}}, {{os_version}}, {{os_build}}, {{locale}}, {{network}}, {{orientation}}, {{dark_mode}}, {{app_version}}, {{app_build_number}}, {{environment}}, {{steps}}, {{expected}}, {{actual}}, {{frequency}}, {{workaround}}, {{logs_link}}, {{screenshot_link}}, {{video_link}}, {{crash_trace}}, {{env_only}}, {{regression_target}}, {{reporter}}, {{date}} -->

# Mobile Bug Report

## Summary

| Field        | Value             |
| ------------ | ----------------- |
| Title        | {{title}}         |
| Severity     | {{severity}}      |
| Priority     | {{priority}}      |
| Reporter     | {{reporter}}      |
| Date         | {{date}}          |
| Environment  | {{environment}}   |
| Env-only?    | {{env_only}}      |

> Severity scale: **S1** Crash / data loss · **S2** Major feature broken · **S3** Minor / workaround exists · **S4** Cosmetic.
> Priority scale: **P0** Fix now · **P1** Fix this sprint · **P2** Backlog · **P3** Won't fix unless escalated.

---

## 1. Device Fingerprint

| Field         | Value               |
| ------------- | ------------------- |
| Make          | {{device_make}}     |
| Model         | {{device_model}}    |
| OS            | {{os_name}}         |
| OS Version    | {{os_version}}      |
| OS Build      | {{os_build}}        |
| Locale        | {{locale}}          |
| Network       | {{network}}         |
| Orientation   | {{orientation}}     |
| Dark mode     | {{dark_mode}}       |

## 2. App Build

| Field          | Value                  |
| -------------- | ---------------------- |
| App version    | {{app_version}}        |
| Build number   | {{app_build_number}}   |
| Distribution   | TestFlight / Play Internal / Ad-hoc |

---

## 3. Steps to Reproduce

{{steps}}

> Use a numbered list. Be explicit about taps, swipes, dialog choices. Reference any preconditions (logged-in user, feature flag, seeded data).

## 4. Expected Result

{{expected}}

## 5. Actual Result

{{actual}}

## 6. Frequency

| Value      | Meaning                                        |
| ---------- | ---------------------------------------------- |
| **always** | 100% reproduction on the listed device         |
| **often**  | More than half the attempts                    |
| **rare**   | Less than half; intermittent                   |

Observed: **{{frequency}}**

## 7. Workaround

{{workaround}}

> If `none`, explicitly write "none". Do not leave blank.

---

## 8. Evidence

| Asset       | Link                  |
| ----------- | --------------------- |
| Logs        | {{logs_link}}         |
| Screenshot  | {{screenshot_link}}   |
| Video       | {{video_link}}        |

> Logs MUST contain the timestamp window of the repro. Strip PII before attaching.

## 9. Crash Trace

```
{{crash_trace}}
```

> Paste the symbolicated trace if the bug is a crash. Otherwise write `n/a`.

---

## 10. Triage Metadata

| Field                | Value                |
| -------------------- | -------------------- |
| Environment-only     | {{env_only}}         |
| Regression target    | {{regression_target}} |
| Suspected component  |                      |
| Linked test case     |                      |
| Linked telemetry     |                      |

> **Environment-only** means the bug only reproduces on a specific device, OS version, locale, or network configuration — not a universal regression. If `true`, the device fingerprint above is the root cause hypothesis.
> **Regression target** is the earliest known good version (e.g. `1.42.0`). Used to bisect the breaking commit.

---

## 11. Reporter Checklist

- [ ] Title is action-oriented and specific
- [ ] Severity and Priority are set (not both `medium`)
- [ ] Steps reproduce on a clean install
- [ ] Logs / video / screenshot attached
- [ ] Device fingerprint is complete
- [ ] PII redacted from all evidence
- [ ] Linked to feature ticket and test case
