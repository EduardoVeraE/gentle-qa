---
name: qa-contract-pact
description: Consumer-driven contract testing with Pact for HTTP and message-based APIs across PACT-JVM (Java/Kotlin/Scala) and PACT-JS (Node/TypeScript). Covers consumer expectation tests, pact-file generation, provider verification, Pact Broker / PactFlow integration, can-i-deploy CI gating, consumer-version selectors, and main-branch / environment tagging. Use when asked to write Pact consumer tests, write Pact provider verification, set up a Pact Broker, gate deploys with can-i-deploy, debug Pact verification mismatches, choose between Pact and Spring Cloud Contract, or design a consumer-driven contract workflow across teams. Trigger keywords - Pact, PACT-JVM, PACT-JS, consumer-driven contract, pact-broker, PactFlow, can-i-deploy, provider verification, consumer version selectors, message pact, pact matchers, like/eachLike, provider state, stateHandlers, broker webhooks, pact tags. NOT for OpenAPI schema validation or OpenAPI-as-contract — use `api-testing` (`references/openapi-driven-testing.md`). NOT for end-to-end functional API tests, request/response status checks, or generic schema assertions — use `api-testing` or `karate-dsl`. NOT for Spring Cloud Contract (provider-driven, Groovy DSL) — that is a different paradigm; if the team needs producer-driven contracts use Spring Cloud Contract directly. NOT for API security testing — use `qa-owasp-security`. NOT for API performance / load — use `k6-load-test`. NOT for E2E browser flows — use `playwright-e2e-testing` / `selenium-e2e-testing`.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

# Consumer-Driven Contract Testing with Pact

Pact is the de-facto consumer-driven contract testing framework. It captures the actual HTTP (or message) interactions a consumer needs from a provider, persists them as a **pact file**, and replays them against the real provider in a verification step. The Pact Broker (or PactFlow) is the registry that connects the two sides and gates deploys via `can-i-deploy`.

**Core principle**: the CONSUMER drives the contract. If the consumer never reads a field, it does not belong in the pact. Over-specified pacts are brittle; under-specified ones miss drift. Use matchers (`like`, `integer`, `eachLike`) to assert SHAPE, not exact values.

## When to Use This Skill

- Writing **PACT-JS** consumer tests (Node / TypeScript) with `@pact-foundation/pact` v12+
- Writing **PACT-JVM** consumer tests (Java / Kotlin / Scala) with `au.com.dius.pact` v4.6+
- Writing **provider verification** — replaying pact files against the real provider with state handlers
- Setting up a **Pact Broker** (self-hosted) or **PactFlow** (managed)
- **Publishing pacts** from CI with version + branch + environment tags
- Implementing **`can-i-deploy`** as a deploy gate
- Designing **consumer version selectors** (`mainBranch`, `deployedOrReleased`, `matchingBranch`) so providers do not verify every pact ever published
- **Message contracts** for async messaging (Kafka, RabbitMQ, SNS/SQS)
- **Webhooks** between Pact Broker and CI for cross-team verification triggers
- Debugging **flaky provider verifications** caused by non-deterministic state handlers
- Choosing **Pact vs OpenAPI-as-contract vs Spring Cloud Contract** for a given team topology

## ISTQB Position

This skill sits at **Layer 3 — Functional by level**, specifically **Integration level**, with two distinct phases:

| Phase | Level | What runs | Mocks |
| ----- | ----- | --------- | ----- |
| Consumer test | Component (with mock provider) | Real consumer code against Pact mock server | Pact mock server replaces provider |
| Provider verification | Integration | Real provider against pact file | State handlers seed real or test DB |

5-layer taxonomy alignment:

| Layer | Coverage |
| ----- | -------- |
| 1. Foundation | `qa-manual-istqb` |
| 2. Strategy | `qa-manual-istqb`, `playwright-regression-strategy` |
| 3. Functional by level | `api-testing`, **`qa-contract-pact`** (this skill — integration), `playwright-e2e-testing`, `selenium-e2e-testing`, `qa-mobile-testing`, `karate-dsl` |
| 4. Non-functional by type | `qa-owasp-security`, `k6-load-test`, `a11y-playwright-testing`, `a11y-selenium-testing` |
| 5. Tooling | `playwright-cli`, `playwright-mcp-inspect` |

## Pact vs Other Contract Approaches

| Approach | Driver | Best for | Reference |
| -------- | ------ | -------- | --------- |
| **Pact (consumer-driven)** | Consumer | External / many consumers, polyglot, cross-team negotiation | THIS skill |
| **OpenAPI as contract** | Provider | Internal monorepo, single-language, one provider | `api-testing/references/openapi-driven-testing.md` |
| **Spring Cloud Contract** | Provider | Java/Spring shops, producer owns contracts in Groovy DSL | Out of scope — use SCC docs directly |
| **Karate strict match** | Test author | Black-box request/response shape lock-in within a Karate suite | `karate-dsl` — `match response == { ... }` |

Decision tree:

```
External or third-party consumers?
├── Yes → Pact
└── No → Same monorepo, same release train, same language?
        ├── Yes → OpenAPI as contract (lower setup)
        └── No → Polyglot or separate cadences?
                ├── Yes → Pact (or both)
                └── Java/Spring with provider-owned contracts?
                        └── Spring Cloud Contract
```

Heuristic:

- 1 provider, 1–3 internal consumers → **OpenAPI**.
- 1 provider, N external consumers → **Pact**.
- N providers, M consumers, separate teams / cadences → **Pact + OpenAPI**.

## Stacks Covered

| Stack | Library | Versions | Reference |
| ----- | ------- | -------- | --------- |
| Node / TypeScript | `@pact-foundation/pact` | 12+ | `references/pact-js.md` |
| Java / Kotlin / Scala | `au.com.dius.pact:consumer:junit5` / `au.com.dius.pact:provider:junit5` | 4.6+ | `references/pact-jvm.md` |
| Pact Broker / PactFlow | `pact-broker` CLI, GitHub Actions | broker 2.107+ | `references/broker-and-cicd.md` |
| Troubleshooting | matchers, state setup, broker auth, webhooks | — | `references/troubleshooting.md` |

## Prerequisites

| Requirement | Notes |
| ----------- | ----- |
| Node.js 18+ | For PACT-JS consumer + provider verification, broker CLI |
| JDK 21+ + Gradle 8 / Maven 3.9 | For PACT-JVM consumer + provider verification |
| Pact Broker URL + token | Self-hosted (Docker `pactfoundation/pact-broker`) or PactFlow account |
| CI variables | `PACT_BROKER_BASE_URL`, `PACT_BROKER_TOKEN`, `GIT_SHA`, `GIT_BRANCH` |
| Test framework | Jest / Mocha / Vitest (JS); JUnit 5 (JVM) |
| Provider state seeding | Idempotent helpers that reset + seed per state name |

## Quick Start

Generate a contract artifact from templates (CLI implemented in `scripts/contract_artifacts.mjs`):

```bash
# List available templates
node scripts/contract_artifacts.mjs list

# Create a contract test charter
node scripts/contract_artifacts.mjs create contract-test-charter --out specs --project "PaymentsAPI"

# Create a broker setup checklist
node scripts/contract_artifacts.mjs create broker-setup-checklist --out specs --release "R1"

# Create a can-i-deploy CI gate spec
node scripts/contract_artifacts.mjs create can-i-deploy-gate --out specs/ci --title "WebApp deploy gate"

# Create a Pact mismatch triage report
node scripts/contract_artifacts.mjs create pact-mismatch-report --out specs/bugs --title "OrderService field rename break"
```

Runnable starters live under `templates/`:

- `templates/pact-js-consumer.test.ts` — PACT-JS consumer test (Jest, TypeScript, V3 matchers)
- `templates/pact-js-provider.verify.ts` — PACT-JS provider verification with state handlers
- `templates/pact-jvm-consumer.java` — PACT-JVM consumer test (JUnit 5, Java)
- `templates/pact-jvm-provider.java` — PACT-JVM provider verification (JUnit 5 + Spring Boot)

## Workflows

### 1) Bootstrap a consumer-driven contract suite (PACT-JS)

1. Add `@pact-foundation/pact` v12+ as a dev dependency.
2. Create a consumer test that **does not call the real provider** — use `PactV3` mock server.
3. Use matchers (`like`, `integer`, `string`, `eachLike`) for every dynamic field. Exact values are for enums and discriminators only.
4. Run the consumer test → it generates a pact JSON under `pacts/`.
5. Publish the pact to the broker with consumer version + branch (`pact-broker publish`).
6. Coordinate with the provider team to enable verification.

See `references/pact-js.md`.

### 2) Bootstrap a consumer-driven contract suite (PACT-JVM)

1. Add `au.com.dius.pact.consumer:junit5` to `build.gradle` / `pom.xml`.
2. Annotate the test class with `@ExtendWith(PactConsumerTestExt.class)` and `@PactTestFor(providerName = "...")`.
3. Define interactions in a `@Pact` method using `PactDslWithProvider` and `LambdaDsl` for body.
4. Run with `mvn test` / `gradle test` → pact JSON appears under `target/pacts` / `build/pacts`.
5. Publish via `pact-broker publish` or the Gradle plugin (`pactPublish` task).

See `references/pact-jvm.md`.

### 3) Implement provider verification

1. Add `@pact-foundation/pact` (JS) or `au.com.dius.pact.provider:junit5` (JVM) to the provider repo.
2. Boot the provider in a known clean state (random port, in-memory or test DB).
3. Configure `Verifier` with `pactBrokerUrl`, `providerVersion` (= `GIT_SHA`), `providerVersionBranch` (= `GIT_BRANCH`), `publishVerificationResult: true`.
4. Use `consumerVersionSelectors` to limit verification scope: `{ mainBranch: true }`, `{ deployedOrReleased: true }`, `{ matchingBranch: true }`.
5. Wire `stateHandlers` — each handler resets DB and seeds the named state idempotently. Non-determinism is the #1 cause of flake.
6. Run verification in CI on every provider PR + main branch push.

See `references/pact-js.md` and `references/pact-jvm.md`.

### 4) Stand up a Pact Broker (or PactFlow)

1. Decide self-hosted vs PactFlow. PactFlow adds RBAC, secrets management, and a UI matrix view.
2. Self-hosted: run `pactfoundation/pact-broker` Docker image with Postgres backend behind HTTPS.
3. Provision read/write tokens; rotate on a schedule.
4. Configure webhooks: on `contract_content_changed` → trigger provider verification CI; on `provider_verification_published` → trigger consumer can-i-deploy.
5. Document the broker URL, tokens (in secret store), and tagging convention in the team handbook.

See `references/broker-and-cicd.md`.

### 5) Gate deploys with `can-i-deploy`

1. Before any deploy step in CI, run:

   ```bash
   pact-broker can-i-deploy \
     --pacticipant <name> \
     --version $GIT_SHA \
     --to-environment production \
     --broker-base-url $PACT_BROKER_BASE_URL \
     --broker-token $PACT_BROKER_TOKEN
   ```

2. Exit code 0 → safe to deploy. Non-zero → block.
3. After successful deploy, record the version in the broker:

   ```bash
   pact-broker record-deployment \
     --pacticipant <name> \
     --version $GIT_SHA \
     --environment production
   ```

4. Tag versions by environment (`staging`, `production`) and main branch (`main`).

See `references/broker-and-cicd.md`.

### 6) Add a message contract (async)

1. Identify the message channel (Kafka topic, RabbitMQ queue, SNS topic).
2. Write a consumer message pact: assert the message handler can process a message of the given shape.
3. Provider verification: assert the producer code emits messages matching the pact, using a function adapter (no real broker).
4. Use the same broker for message pacts; selectors and `can-i-deploy` work identically.

See `references/pact-js.md` (message section) and `references/pact-jvm.md` (message section).

### 7) Triage a Pact verification mismatch

1. Pull the verification report from the broker (or CI logs) — it lists the FAILED interaction, expected vs actual.
2. Classify: SHAPE drift (field renamed, type changed) vs STATE drift (DB not seeded) vs MATCHER drift (over-specified).
3. SHAPE drift → coordinate change with consumer; bump consumer pact OR add backwards-compatible field.
4. STATE drift → fix `stateHandlers` to be idempotent and reset before seeding.
5. MATCHER drift → loosen consumer pact (use `like` instead of exact value).
6. Re-run verification; if green, publish result.

See `references/troubleshooting.md`.

## Inputs to Collect

- **Consumer / provider names** — stable IDs registered in the broker.
- **Broker URL + token** — self-hosted or PactFlow.
- **Versioning convention** — recommended `GIT_SHA` (immutable) + branch tag.
- **Branch model** — main-branch verification + per-PR consumer pacts.
- **Environments** — `staging`, `production` minimum; map to broker tags.
- **State catalog** — list of provider states each consumer test depends on.
- **CI integration** — webhooks from broker to provider CI; can-i-deploy gate before deploy.
- **Message channels** (if async) — topic / queue names and message envelope schema.

## Outputs

| Artifact | When produced | Template |
| -------- | ------------- | -------- |
| Contract test charter | Per consumer/provider pair | `templates/contract-test-charter.md` (via CLI) |
| Broker setup checklist | Per broker rollout | `templates/broker-setup-checklist.md` (via CLI) |
| can-i-deploy CI gate spec | Per pacticipant | `templates/can-i-deploy-gate.md` (via CLI) |
| Pact mismatch triage report | Per failed verification | `templates/pact-mismatch-report.md` (via CLI) |
| Consumer test code | Per interaction | `templates/pact-js-consumer.test.ts`, `templates/pact-jvm-consumer.java` |
| Provider verification code | Per provider | `templates/pact-js-provider.verify.ts`, `templates/pact-jvm-provider.java` |

## Exclusions

This skill is deliberately scoped. Do NOT use it for:

- **OpenAPI schema validation, OpenAPI-as-contract, runtime spec validation** — use `api-testing` (`references/openapi-driven-testing.md` and `references/headers-and-contracts.md`).
- **General API functional testing, status code coverage, schema assertions** — use `api-testing` (TS / Java) or `karate-dsl`.
- **Spring Cloud Contract (provider-driven Groovy DSL)** — different paradigm; refer the team to SCC documentation directly.
- **API security testing** (OWASP API Top 10, BOLA, JWT attacks) — use `qa-owasp-security`.
- **API performance / load** — use `k6-load-test`.
- **E2E browser flows** — use `playwright-e2e-testing` or `selenium-e2e-testing`.
- **Mobile harness, gestures, device matrix** — use `qa-mobile-testing` (Pact still applies if the mobile app consumes an API; write the Pact in this skill, run it in the mobile project).

If a request mixes concerns (e.g., "verify the contract AND check rate-limit headers AND test SQLi"), split it: contract here, headers in `api-testing`, SQLi in `qa-owasp-security`.
