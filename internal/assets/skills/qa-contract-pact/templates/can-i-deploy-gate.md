<!-- Skill: qa-contract-pact · Template: can-i-deploy-gate -->
<!-- Placeholders: {{project}}, {{owner}}, {{date}}, {{author}}, {{pacticipant}}, {{environment}}, {{ci_platform}}, {{deploy_command}}, {{record_deployment_command}}, {{rollback_strategy}}, {{stakeholders}} -->

# can-i-deploy CI Gate Spec

> Defines the contract-driven deploy gate for a single pacticipant + environment.

## 1. Document Control

| Field         | Value              |
| ------------- | ------------------ |
| Project       | {{project}}        |
| Pacticipant   | {{pacticipant}}    |
| Environment   | {{environment}}    |
| Owner         | {{owner}}          |
| Author        | {{author}}         |
| Date          | {{date}}           |
| Stakeholders  | {{stakeholders}}   |

---

## 2. CI Platform

Selected: **{{ci_platform}}** <!-- e.g., GitHub Actions, GitLab CI, CircleCI, Buildkite -->

---

## 3. Pipeline Order

```
build → test → can-i-deploy → deploy → record-deployment
```

- `build` → produces immutable artifact tagged with `$GIT_SHA`.
- `test` → unit + contract tests; consumer pact published.
- `can-i-deploy` → BLOCKS on incompatibility.
- `deploy` → runs only if can-i-deploy returns exit 0.
- `record-deployment` → updates broker so future can-i-deploy queries succeed.

---

## 4. can-i-deploy Step

```bash
pact-broker can-i-deploy \
  --pacticipant {{pacticipant}} \
  --version $GIT_SHA \
  --to-environment {{environment}} \
  --broker-base-url $PACT_BROKER_BASE_URL \
  --broker-token $PACT_BROKER_TOKEN
```

Exit codes:

| Code | Meaning                            | CI Action               |
| ---- | ---------------------------------- | ----------------------- |
| 0    | Compatible — safe to deploy        | Continue to deploy step |
| 1    | Incompatible / unverified          | Fail pipeline           |
| 2+   | CLI error (network, auth)          | Fail pipeline           |

---

## 5. Deploy Step

Executes only if can-i-deploy passes:

```bash
{{deploy_command}}
```

---

## 6. record-deployment Step

After successful deploy:

```bash
{{record_deployment_command}}
```

Default form:

```bash
pact-broker record-deployment \
  --pacticipant {{pacticipant}} \
  --version $GIT_SHA \
  --environment {{environment}}
```

---

## 7. Rollback Strategy

If a regression is detected post-deploy:

{{rollback_strategy}} <!-- e.g., `record-undeployment` for previous version, redeploy last green SHA -->

```bash
pact-broker record-undeployment \
  --pacticipant {{pacticipant}} \
  --environment {{environment}}
```

---

## 8. Failure Modes

| Failure                              | Diagnosis                                      | Fix                                            |
| ------------------------------------ | ---------------------------------------------- | ---------------------------------------------- |
| can-i-deploy returns 1               | Missing verification or known mismatch         | Re-run provider verification; investigate mismatch |
| can-i-deploy returns 2+              | Broker auth / network                          | Check token + URL; retry                       |
| Provider verified locally, not in CI | Result not published (`publishVerificationResult` off) | Enable publish in provider CI            |
| New consumer not gated               | record-deployment never ran for prior version  | Backfill record-deployment                     |

---

## 9. Definition of Done

- [ ] can-i-deploy step inserted before every deploy job to {{environment}}.
- [ ] record-deployment step inserted after every successful deploy.
- [ ] Pipeline fails closed (no override / continue-on-error).
- [ ] Rollback runbook linked.
- [ ] One full green run executed end-to-end.
