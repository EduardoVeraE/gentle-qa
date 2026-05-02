<!-- Skill: qa-contract-pact · Template: broker-setup-checklist -->
<!-- Placeholders: {{project}}, {{owner}}, {{date}}, {{author}}, {{release}}, {{broker_kind}}, {{broker_url}}, {{environments}}, {{pacticipants}}, {{token_rotation_cadence}}, {{webhook_targets}}, {{retention_policy}} -->

# Pact Broker Setup Checklist

> Use during initial broker rollout or migration (self-hosted ↔ PactFlow).

## 1. Document Control

| Field    | Value          |
| -------- | -------------- |
| Project  | {{project}}    |
| Owner    | {{owner}}      |
| Author   | {{author}}     |
| Date     | {{date}}       |
| Release  | {{release}}    |
| Status   | Draft          |

---

## 2. Decision

Selected broker: **{{broker_kind}}** <!-- e.g., self-hosted pactfoundation/pact-broker 2.116 OR PactFlow managed -->

Broker URL: `{{broker_url}}`

---

## 3. Environments

Configure first-class environments in the broker:

- {{environments}} <!-- e.g., staging, production -->

```bash
pact-broker create-environment --name staging
pact-broker create-environment --name production
```

---

## 4. Pacticipants

Initial pacticipants to register:

- {{pacticipants}} <!-- e.g., WebApp (consumer), UserService (provider), OrderService (provider) -->

---

## 5. Tokens + Auth

- [ ] Admin credential created (admin user / org admin).
- [ ] Read token created (CI read-only access).
- [ ] Write tokens created per pacticipant (least privilege).
- [ ] Tokens stored in CI secret manager (no plaintext in repo).
- [ ] Token rotation cadence: {{token_rotation_cadence}}

---

## 6. Webhooks

| Event                              | Target                            | Purpose                                |
| ---------------------------------- | --------------------------------- | -------------------------------------- |
| `contract_content_changed`         | {{webhook_targets}}               | Trigger provider verification CI       |
| `provider_verification_published`  | {{webhook_targets}}               | Trigger consumer can-i-deploy         |

- [ ] Webhook secrets stored securely.
- [ ] Webhook retry policy configured (default 5 retries with backoff).
- [ ] Webhooks tested with a synthetic event.

---

## 7. CI Integration

- [ ] Consumer pipelines publish pacts with `$GIT_SHA` + `$GIT_BRANCH`.
- [ ] Provider pipelines verify with `consumerVersionSelectors: [mainBranch, deployedOrReleased, matchingBranch]`.
- [ ] Provider pipelines publish verification result (`publishVerificationResult: true`).
- [ ] `can-i-deploy --to-environment <env>` runs before every deploy.
- [ ] `record-deployment --environment <env>` runs after every successful deploy.

---

## 8. Retention

Pact retention policy: {{retention_policy}} <!-- e.g., keep main + last 30d feature branches + all deployed versions -->

```bash
pact-broker clean --keep-version-selector '{"branch":"main"}' \
                  --keep-version-selector '{"deployedOrReleased":true}' \
                  --keep-version-selector '{"maxAgeInDays":30}'
```

Schedule cleanup: weekly cron.

---

## 9. Observability

- [ ] Broker uptime monitor (HTTP heartbeat: `/diagnostic/status/heartbeat`).
- [ ] Postgres backup scheduled (self-hosted only); restore drill quarterly.
- [ ] Audit log retention defined.

---

## 10. Definition of Done

- [ ] All pacticipants registered.
- [ ] All webhooks fire on synthetic events.
- [ ] At least one consumer pact published end-to-end and verified.
- [ ] `can-i-deploy` gates one production deploy successfully.
- [ ] Runbook linked from team handbook.
