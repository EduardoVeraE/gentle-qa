<!-- Skill: qa-contract-pact · Template: pact-mismatch-report -->
<!-- Placeholders: {{project}}, {{title}}, {{date}}, {{author}}, {{owner}}, {{consumer}}, {{provider}}, {{consumer_version}}, {{provider_version}}, {{environment}}, {{interaction}}, {{state}}, {{expected}}, {{actual}}, {{drift_kind}}, {{owning_repo}}, {{fix_summary}}, {{verification_url}} -->

# Pact Mismatch Triage Report

> One report per failed verification interaction.

## 1. Summary

| Field             | Value                       |
| ----------------- | --------------------------- |
| Title             | {{title}}                   |
| Project           | {{project}}                 |
| Date              | {{date}}                    |
| Author            | {{author}}                  |
| Owner             | {{owner}}                   |
| Consumer          | {{consumer}}                |
| Provider          | {{provider}}                |
| Consumer version  | {{consumer_version}}        |
| Provider version  | {{provider_version}}        |
| Environment       | {{environment}}             |
| Verification URL  | {{verification_url}}        |

---

## 2. Failed Interaction

- Description: {{interaction}} <!-- e.g., "a request for user 1" -->
- Provider state: {{state}}    <!-- e.g., "user 1 exists" -->

### Expected (consumer pact)

```
{{expected}}
```

### Actual (provider response)

```
{{actual}}
```

---

## 3. Classification

Drift kind: **{{drift_kind}}**

| Kind          | Meaning                                                            | Owning side          |
| ------------- | ------------------------------------------------------------------ | -------------------- |
| SHAPE drift   | Provider changed response shape (renamed/typed/removed field)      | Provider             |
| STATE drift   | Provider state handler did not seed expected data                  | Provider             |
| MATCHER drift | Consumer over-specified an exact value that provider varies        | Consumer             |
| AUTH drift    | Provider auth required, request filter missing in verification     | Provider verification |

---

## 4. Root Cause

<!-- Describe what changed and why the contract surfaced it. -->

---

## 5. Fix

Owning repo: `{{owning_repo}}`

Fix summary: {{fix_summary}}

Concrete change (link to PR / commit when available):

- [ ] Code change applied
- [ ] Unit tests updated
- [ ] Pact re-published (if consumer-side)
- [ ] Verification re-run and green
- [ ] `can-i-deploy` re-run after fix

---

## 6. Prevention

- [ ] Add lint / matcher in consumer to prevent re-introducing exact value.
- [ ] Add idempotent reset to state handler (if STATE drift).
- [ ] Add CI step to publish verification result (if AUTH drift on result publishing).
- [ ] Document in contract test charter under "Provider State Catalog" or "Risks".

---

## 7. Timeline

| Event                       | When     |
| --------------------------- | -------- |
| Verification failed         | {{date}} |
| Triage started              |          |
| Root cause identified       |          |
| Fix merged                  |          |
| Verification green          |          |
| can-i-deploy unblocked      |          |
| Deploy resumed              |          |
