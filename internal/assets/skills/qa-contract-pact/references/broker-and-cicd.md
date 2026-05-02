# Pact Broker + CI/CD

The Pact Broker is the central registry for pacts and verification results. Without it, contract testing is just unit tests with extra steps. This reference covers self-hosted and PactFlow setups, publishing strategies, `can-i-deploy` gating, consumer version selectors, branch tagging, environments, webhooks, and the matrix view.

## 1. Self-Hosted vs PactFlow

| Feature | Self-hosted (`pactfoundation/pact-broker` Docker) | PactFlow (managed) |
| ------- | ------------------------------------------------- | ------------------ |
| Cost | Free (you pay infra) | SaaS subscription |
| Ops burden | DB backups, TLS, upgrades on you | Managed |
| RBAC | Basic auth + tokens | RBAC roles, SSO, audit log |
| Secrets | Env vars only | Per-pacticipant secrets store |
| UI matrix | Basic | Enhanced with deployment timeline |
| Webhooks | Yes | Yes + retry policy UI |
| Choose when | Strict data residency, on-prem requirement | Cross-team org, want it to "just work" |

For teams just starting out: spin up self-hosted in a sandbox, migrate to PactFlow once you have > 5 pacticipants and cross-team friction.

## 2. Self-Hosted Broker (Docker Compose)

```yaml
# docker-compose.broker.yml
services:
  postgres:
    image: postgres:16
    environment:
      POSTGRES_USER: pact
      POSTGRES_PASSWORD: pact
      POSTGRES_DB: pact_broker
    volumes:
      - pgdata:/var/lib/postgresql/data
  broker:
    image: pactfoundation/pact-broker:2.116.0
    depends_on: [postgres]
    ports:
      - "9292:9292"
    environment:
      PACT_BROKER_DATABASE_URL: "postgres://pact:pact@postgres/pact_broker"
      PACT_BROKER_BASIC_AUTH_USERNAME: admin
      PACT_BROKER_BASIC_AUTH_PASSWORD: ${BROKER_ADMIN_PASS}
      PACT_BROKER_PUBLIC_HEARTBEAT: "true"
volumes:
  pgdata:
```

Front it with TLS (Caddy / Nginx) before exposing publicly. Generate read/write tokens via the broker UI; rotate quarterly.

## 3. Versioning + Tagging Convention

The single most important convention: **the consumer / provider version is the immutable git SHA**, with branches and environments as mutable tags.

| Identifier | Source | Mutable? |
| ---------- | ------ | :------: |
| `pacticipantVersion` | `$GIT_SHA` | No |
| `branch` | `$GIT_BRANCH` (`main`, `feature/x`) | Yes |
| `environment` tag | `staging`, `production` | Yes |

```bash
# Publish with branch
pact-broker publish ./pacts \
  --consumer-app-version=$GIT_SHA \
  --branch=$GIT_BRANCH \
  --broker-base-url=$PACT_BROKER_BASE_URL \
  --broker-token=$PACT_BROKER_TOKEN

# After successful prod deploy
pact-broker record-deployment \
  --pacticipant WebApp \
  --version=$GIT_SHA \
  --environment=production
```

## 4. Consumer Version Selectors

Without selectors, the provider verifies EVERY pact ever published. With selectors, it verifies only what matters.

Recommended baseline:

```typescript
consumerVersionSelectors: [
  { mainBranch: true },        // latest pact from each consumer's main branch
  { deployedOrReleased: true },// latest pact deployed to any environment
  { matchingBranch: true },    // pact from a consumer branch matching provider branch
],
```

| Selector | Use when |
| -------- | -------- |
| `{ mainBranch: true }` | Always — verify trunk |
| `{ deployedOrReleased: true }` | Always — never break what is live |
| `{ matchingBranch: true }` | Coordinated feature branches across consumer + provider |
| `{ branch: "release/2.0" }` | Verify a specific release line |
| `{ tag: "production" }` | Legacy tag-based selectors (pre-environments) |

## 5. `can-i-deploy` Gate

`can-i-deploy` asks the broker: "given pacticipant X at version V, can it be safely deployed to environment E?" Answers yes only if every required pact has been verified between the consumer and the provider versions in that environment.

```bash
pact-broker can-i-deploy \
  --pacticipant WebApp \
  --version $GIT_SHA \
  --to-environment production \
  --broker-base-url $PACT_BROKER_BASE_URL \
  --broker-token $PACT_BROKER_TOKEN
```

Exit codes:

- `0` → safe to deploy.
- `1` → not safe (missing verification or known incompatibility).
- `2+` → CLI error (network, auth).

CI integration: place `can-i-deploy` BEFORE the deploy step. On failure, the deploy step never runs.

```yaml
- name: can-i-deploy production
  run: pact-broker can-i-deploy --pacticipant WebApp --version $GIT_SHA --to-environment production --broker-base-url $PACT_BROKER_BASE_URL --broker-token $PACT_BROKER_TOKEN

- name: deploy production
  if: success()
  run: ./deploy.sh production

- name: record deployment
  if: success()
  run: pact-broker record-deployment --pacticipant WebApp --version $GIT_SHA --environment production
```

## 6. Webhooks

Webhooks let the broker trigger CI in another repo when relevant events happen. Two events matter most:

| Event | Trigger CI for | Why |
| ----- | -------------- | --- |
| `contract_content_changed` | Provider | New pact published → re-verify |
| `provider_verification_published` | Consumer | Verification result available → re-run can-i-deploy |

GitHub Actions webhook example (broker side):

```bash
pact-broker create-webhook \
  https://api.github.com/repos/<org>/<provider-repo>/dispatches \
  --request POST \
  --header "Authorization: Bearer ${GH_TOKEN}" \
  --header "Accept: application/vnd.github+json" \
  --data '{"event_type":"pact_changed","client_payload":{"pact_url":"${pactbroker.pactUrl}"}}' \
  --consumer WebApp \
  --provider UserService \
  --contract-content-changed
```

The provider repo listens on `repository_dispatch` for `pact_changed` and runs verification.

## 7. Matrix View

The broker's matrix view answers: "for pacticipant X, which versions of consumers/providers are compatible?" Use it to:

- Diagnose `can-i-deploy` failures (which pact / version is missing verification).
- Audit which consumer versions are deployed to production.
- Confirm a new provider version is verified by all main-branch consumers before promote.

URL: `${PACT_BROKER_BASE_URL}/matrix?q[][pacticipant]=WebApp&q[][version]=$GIT_SHA`.

## 8. Branch Strategy

Recommended branch model for contract verification:

```
main (provider) ──────── verifies main + production + matchingBranch ─────────►
                                                                     │
feature/x (provider) ─── verifies main + matchingBranch ─────────────┤
                                                                     │
main (consumer)    ────► publishes pact with branch=main ────────────┤
feature/x (consumer)──► publishes pact with branch=feature/x ────────┘
```

- Every consumer push publishes the pact tagged with its branch.
- Every provider push runs verification with `mainBranch + deployedOrReleased + matchingBranch`.
- Coordinated breaking changes use the SAME branch name across consumer + provider; `matchingBranch` selector finds the matching pact.

## 9. Environments

Environments are first-class in modern broker (replaces tag-based promotion). Define them once:

```bash
pact-broker create-environment --name staging
pact-broker create-environment --name production
```

Then `record-deployment` with `--environment` to promote a version.

## 10. Pitfalls

- **Mutable consumer versions** (e.g., `version=latest`) — `can-i-deploy` cannot gate; always use `$GIT_SHA`.
- **No `record-deployment`** — broker thinks nothing is deployed; `deployedOrReleased` selector returns empty.
- **No webhooks** — provider does not know a new pact landed; verifications run only on schedule, drift goes undetected.
- **Skipping `--to-environment`** in `can-i-deploy` — falls back to "is this version verified at all" which is a weaker check.
- **Single shared token** with full write permissions — lost token leaks all pacts. Use per-pacticipant scoped tokens (PactFlow) or rotate often.
- **Unbounded pact retention** — broker grows forever. Configure clean-up: keep main + last N feature branches + deployed versions.
