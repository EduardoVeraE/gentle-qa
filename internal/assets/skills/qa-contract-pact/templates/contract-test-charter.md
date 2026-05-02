<!-- Skill: qa-contract-pact · Template: contract-test-charter -->
<!-- Placeholders: {{project}}, {{consumer}}, {{provider}}, {{owner}}, {{date}}, {{author}}, {{stack}}, {{broker_url}}, {{environments}}, {{branch_strategy}}, {{state_catalog}}, {{message_channels}}, {{interactions_in_scope}}, {{interactions_out_of_scope}}, {{can_i_deploy_pacticipant}}, {{verification_schedule}}, {{stakeholders}} -->

# Contract Test Charter

> Charter for a single consumer/provider pair. One charter per pact relationship.

## 1. Document Control

| Field        | Value             |
| ------------ | ----------------- |
| Project      | {{project}}       |
| Consumer     | {{consumer}}      |
| Provider     | {{provider}}      |
| Owner        | {{owner}}         |
| Author       | {{author}}        |
| Date         | {{date}}          |
| Status       | Draft             |
| Stakeholders | {{stakeholders}}  |

---

## 2. Stack

- Consumer stack: {{stack}} <!-- e.g., Node 20 + TypeScript + @pact-foundation/pact 12 -->
- Provider stack: {{stack}}
- Broker: {{broker_url}}
- Pact spec version: V4

---

## 3. Scope

### 3.1 In Scope

- Interactions: {{interactions_in_scope}} <!-- e.g., GET /api/users/{id}, POST /api/orders -->
- Message channels (if any): {{message_channels}} <!-- e.g., kafka topic orders.events -->

### 3.2 Out of Scope

- Interactions: {{interactions_out_of_scope}}
- Business-logic correctness (covered by integration / E2E)
- Performance and security (covered by `k6-load-test` and `qa-owasp-security`)

---

## 4. Branch + Versioning Strategy

- Pacticipant version: immutable `$GIT_SHA`.
- Branch strategy: {{branch_strategy}} <!-- e.g., main + feature/* with matchingBranch selector -->
- Environments tracked: {{environments}} <!-- e.g., staging, production -->

---

## 5. Provider State Catalog

Each consumer interaction declares a `given(...)` state. The provider must implement a state handler with the same name. Names are case-sensitive strings.

| State Name                      | Setup Action                                           | Owner    |
| ------------------------------- | ------------------------------------------------------ | -------- |
| {{state_catalog}}               | (reset DB; seed expected fixture)                      | provider |

---

## 6. CI Integration

- Consumer pipeline:
  1. Run consumer pact tests.
  2. Publish pact to broker with `--consumer-app-version=$GIT_SHA --branch=$GIT_BRANCH`.
  3. Run `can-i-deploy --pacticipant {{can_i_deploy_pacticipant}} --version=$GIT_SHA --to-environment <env>`.
  4. Deploy only if can-i-deploy returns 0.
  5. Record deployment via `pact-broker record-deployment`.

- Provider pipeline:
  1. Run provider verification with `consumerVersionSelectors: [mainBranch, deployedOrReleased, matchingBranch]`.
  2. Publish verification result (`publishVerificationResult: true`).
  3. Verification cadence: {{verification_schedule}} <!-- e.g., every PR + main push + nightly -->

- Broker webhooks:
  - On `contract_content_changed` → trigger provider verification CI.
  - On `provider_verification_published` → trigger consumer can-i-deploy.

---

## 7. Definition of Done

- [ ] Consumer pact tests committed under `pacts.test.*` naming convention.
- [ ] Provider verification tests committed and registered with broker.
- [ ] Broker webhook configured for both events.
- [ ] `can-i-deploy` gate inserted before each deploy step.
- [ ] State catalog reviewed and signed off by both teams.
- [ ] Charter linked from both repos' README.

---

## 8. Risks

| Risk                                                | Mitigation                                                |
| --------------------------------------------------- | --------------------------------------------------------- |
| State handler non-determinism causes flaky verifies | Reset before seed; isolated fixtures                      |
| Over-specification breaks on unrelated provider tweaks | Use matchers; assert only fields consumer reads          |
| Broker downtime blocks deploys                      | Cache last successful can-i-deploy outcome with TTL       |
| Token leakage                                       | Per-pacticipant scoped tokens; rotate quarterly           |
