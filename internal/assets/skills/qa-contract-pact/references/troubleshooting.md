# Pact Troubleshooting

Reference for diagnosing common contract-test failures: flaky verifications, mismatch reports, broker auth, state setup, and matcher pitfalls. Use this when verifications fail intermittently or when a pact passes locally but fails in CI.

## 1. Reading a Verification Report

Every verification produces a structured report. The key sections:

```
Failure/s in pact between WebApp and UserService:

  a request for user 1
    given user 1 exists
      returns a response which
        has status code 200          (OK)
        has a matching body          (FAILED)
          $.email                    (Expected 'string' but got 'null')
```

Classify the failure:

| Symptom | Cause | Fix |
| ------- | ----- | --- |
| `Expected '<type>' but got 'null'` | Provider returned `null` for a typed field | Provider bug OR consumer over-specified |
| `Expected '<value>' but got '<other>'` | Exact match mismatch | Wrap consumer in matcher (`like`, `string`, `integer`) |
| `Provider state '<name>' not found` | State handler not registered | Add `@State("<name>")` (JVM) or `stateHandlers` entry (JS) |
| `Connection refused` | Provider not booted | Ensure verifier waits for ready |
| `401 Unauthorized` | Provider auth required | Use `requestFilter` to inject test token |

## 2. Flaky Verifications — Root Causes

### 2.1 Non-idempotent state handlers

WRONG:

```typescript
"user 1 exists": async () => {
  await seedUser({ id: 1, name: "Alice" }); // Fails on second run with PK collision
}
```

RIGHT:

```typescript
"user 1 exists": async () => {
  await resetDb();
  await seedUser({ id: 1, name: "Alice", email: "alice@example.com" });
}
```

Every state handler MUST reset the relevant data first.

### 2.2 Shared mutable state across interactions

If interaction A creates user 1 and interaction B expects user 1 to NOT exist, the order matters and tests are flaky. Fix:

- Each handler resets fully and seeds only what THAT state needs.
- Use scoped DBs (Testcontainers) where possible.

### 2.3 Time-dependent assertions

`createdAt: "2026-01-01T00:00:00Z"` will eventually drift. Use:

- JS: `iso8601DateTime("2026-01-01T00:00:00Z")` — asserts shape.
- JVM: `o.datetime("createdAt", "yyyy-MM-dd'T'HH:mm:ssXXX", sample)`.

Never assert exact timestamps in pacts.

### 2.4 Concurrency / DB races

Verifier runs interactions sequentially by default; if you parallelize, state handlers can race. Keep verification sequential unless every handler operates on isolated data.

## 3. Mismatch Triage Workflow

```
Verification FAILED
├── Read report: which interaction, which path, expected vs actual
├── Reproduce locally: download pact from broker, run verifier with same selectors
├── Classify:
│   ├── SHAPE drift → coordinate with consumer; provider adds backward-compat field
│   ├── STATE drift → fix state handler (idempotent + correct seed)
│   └── MATCHER drift → loosen consumer pact (replace exact value with matcher)
├── Apply fix in the OWNING repo (consumer for matcher drift; provider for shape drift)
├── Re-run verification
└── Publish result (if green) → unblocks can-i-deploy
```

## 4. Broker Auth Issues

| Symptom | Cause | Fix |
| ------- | ----- | --- |
| `401 Unauthorized` from broker | Token missing or expired | Set `PACT_BROKER_TOKEN`; rotate token |
| `403 Forbidden` from broker | Token lacks pacticipant scope | Use a token with write scope on the pacticipant (PactFlow RBAC) |
| `404 Not Found` for pact | Pact never published | Run `pact-broker publish` from consumer CI |
| `422 Unprocessable` on publish | Invalid version (mutable, missing) | Use `$GIT_SHA`; never `latest` |

## 5. CI-Only Failures (passes locally)

Common causes:

- **Different Node / JDK version** in CI — pin via `setup-node` / `setup-java`.
- **Missing env vars** — broker URL/token only set on `main` workflow but verification runs on PR; gate accordingly.
- **Network egress** — corporate runners may need broker URL allow-listed.
- **Time zone drift** — assert UTC explicitly in datetime matchers.
- **Provider boot timing** — verifier starts before provider is ready; add a readiness probe.

## 6. Over-Specification → Brittleness

If pact verifications break every time the provider tweaks an unrelated field, the pact is over-specified. Fix:

1. Identify fields the consumer NEVER reads.
2. Remove them from the consumer pact body.
3. Replace exact values with matchers for fields you do read.

```typescript
// BEFORE — brittle: any change to email format, role list, or createdAt fails
body: {
  id: 1,
  name: "Alice",
  email: "alice@example.com",
  createdAt: "2026-01-01T00:00:00Z",
  roles: ["admin"],
  internalDebugFlag: false, // consumer never reads this
}

// AFTER — robust
body: like({
  id: integer(1),
  name: string("Alice"),
  email: regex(/.+@.+\..+/, "alice@example.com"),
  // createdAt removed — consumer doesn't read it
  // roles removed — consumer doesn't read it
})
```

## 7. Schema Validation vs Pact

A common confusion: "we already have OpenAPI / Zod / JSON Schema, why Pact?"

| Tool | Catches |
| ---- | ------- |
| OpenAPI lint (Spectral) | STATIC drift between spec and itself |
| `oasdiff` | STATIC breaking changes between two spec versions |
| Zod / JSON Schema runtime validation | Runtime mismatch between PROVIDER response and SPEC |
| Pact verification | Runtime mismatch between PROVIDER response and CONSUMER expectation |

Pact catches the case where the spec says one thing, the consumer reads something compatible-with-spec-but-different, and the provider implements yet a third compatible-with-spec-but-different shape. OpenAPI alone misses this — both sides conform to the spec but disagree on usage.

## 8. When NOT to Use Pact

- Single-team monorepo with same release train → OpenAPI as contract is enough; Pact is overhead.
- Internal API with no team boundaries → karate-dsl strict-match contracts handle this in a single suite.
- Business-logic correctness ("placing an order sends a welcome email") → integration / E2E test, not Pact.
- Performance assertions ("response under 100 ms") → `k6-load-test`.
- Security assertions ("403 on cross-tenant access") → `qa-owasp-security` BOLA tests.

If the team has Pact and is using it for any of the above, they will hit walls. Redirect to the appropriate skill.
