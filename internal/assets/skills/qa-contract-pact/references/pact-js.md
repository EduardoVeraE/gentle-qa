# PACT-JS — Consumer + Provider (Node / TypeScript)

Reference for `@pact-foundation/pact` v12+ in Node 18+. Covers HTTP consumer tests (Jest / Mocha / Vitest), HTTP provider verification, and message contracts for async stacks (Kafka, RabbitMQ, SNS/SQS).

## 1. Install

```bash
npm i -D @pact-foundation/pact
# Optional: broker CLI
npm i -D @pact-foundation/pact-cli
```

`package.json` scripts:

```json
{
  "scripts": {
    "test:pact:consumer": "jest --testPathPattern='pact.test.ts$'",
    "test:pact:provider": "jest --testPathPattern='verify.test.ts$'",
    "pact:publish": "pact-broker publish ./pacts --consumer-app-version=$GIT_SHA --branch=$GIT_BRANCH --broker-base-url=$PACT_BROKER_BASE_URL --broker-token=$PACT_BROKER_TOKEN",
    "pact:can-i-deploy": "pact-broker can-i-deploy --pacticipant WebApp --version=$GIT_SHA --to-environment production --broker-base-url=$PACT_BROKER_BASE_URL --broker-token=$PACT_BROKER_TOKEN"
  }
}
```

## 2. HTTP Consumer Test (PactV3)

`PactV3` produces V3 pact files. The mock server starts on a random port; the consumer's HTTP client points at `mockServer.url`.

```typescript
// consumer/user-service.pact.test.ts
import path from "node:path";
import { PactV3, MatchersV3 } from "@pact-foundation/pact";
import { fetchUser } from "../src/clients/user-client";

const { like, integer, string, eachLike, regex, iso8601DateTime } = MatchersV3;

const provider = new PactV3({
  consumer: "WebApp",
  provider: "UserService",
  dir: path.resolve(process.cwd(), "pacts"),
  logLevel: "warn",
});

describe("UserService consumer pact", () => {
  it("fetches a user by id (happy path)", async () => {
    provider
      .given("user 1 exists")
      .uponReceiving("a request for user 1")
      .withRequest({
        method: "GET",
        path: "/api/users/1",
        headers: { Accept: "application/json" },
      })
      .willRespondWith({
        status: 200,
        headers: { "Content-Type": regex(/application\/json.*/, "application/json; charset=utf-8") },
        body: like({
          id: integer(1),
          name: string("Alice"),
          email: regex(/.+@.+\..+/, "alice@example.com"),
          createdAt: iso8601DateTime("2026-01-01T00:00:00Z"),
          roles: eachLike("admin", { min: 1 }),
        }),
      });

    await provider.executeTest(async (mockServer) => {
      const user = await fetchUser(mockServer.url, 1);
      expect(user.id).toBe(1);
      expect(user.email).toMatch(/@/);
    });
  });

  it("returns 404 when user does not exist", async () => {
    provider
      .given("user 999 does not exist")
      .uponReceiving("a request for missing user 999")
      .withRequest({ method: "GET", path: "/api/users/999" })
      .willRespondWith({
        status: 404,
        headers: { "Content-Type": "application/problem+json" },
        body: like({ type: string("about:blank"), title: string("Not Found"), status: integer(404) }),
      });

    await provider.executeTest(async (mockServer) => {
      await expect(fetchUser(mockServer.url, 999)).rejects.toThrow(/not found/i);
    });
  });
});
```

### Matcher cheat-sheet

| Matcher | What it asserts | Example |
| ------- | --------------- | ------- |
| `like(x)` | Same TYPE as `x` | `like({ id: 1 })` |
| `integer(n)` | Any integer | `integer(42)` |
| `string(s)` | Any string | `string("Alice")` |
| `boolean(b)` | Any boolean | `boolean(true)` |
| `eachLike(item, { min: 1 })` | Array of items shaped like `item` | `eachLike({ id: integer(1) }, { min: 2 })` |
| `regex(pattern, sample)` | String matching regex | `regex(/^\d{4}$/, "1234")` |
| `iso8601DateTime(sample)` | RFC 3339 timestamp | `iso8601DateTime("2026-01-01T00:00:00Z")` |
| `uuid(sample)` | UUID string | `uuid("3f29c6...")` |
| `equal(x)` | EXACT value | `equal("ACTIVE")` — for enums |

Rule: wrap every dynamic field in a matcher. Reserve `equal` for enums and discriminator fields.

## 3. HTTP Provider Verification

```typescript
// provider/user-service.verify.test.ts
import { Verifier } from "@pact-foundation/pact";
import { startServer, type TestServer } from "../src/server";
import { resetDb, seedUser, deleteUser } from "../test/db-helpers";

describe("UserService provider verification", () => {
  let server: TestServer;

  beforeAll(async () => {
    server = await startServer({ port: 0 }); // random port
  });

  afterAll(async () => {
    await server.close();
  });

  it("verifies all consumer pacts", async () => {
    await new Verifier({
      provider: "UserService",
      providerBaseUrl: server.url,
      providerVersion: process.env.GIT_SHA!,
      providerVersionBranch: process.env.GIT_BRANCH!,
      pactBrokerUrl: process.env.PACT_BROKER_BASE_URL!,
      pactBrokerToken: process.env.PACT_BROKER_TOKEN!,
      publishVerificationResult: process.env.CI === "true",
      consumerVersionSelectors: [
        { mainBranch: true },
        { deployedOrReleased: true },
        { matchingBranch: true },
      ],
      stateHandlers: {
        "user 1 exists": {
          setup: async () => {
            await resetDb();
            await seedUser({ id: 1, name: "Alice", email: "alice@example.com", roles: ["admin"] });
          },
          teardown: async () => {
            await deleteUser(1);
          },
        },
        "user 999 does not exist": async () => {
          await resetDb();
          await deleteUser(999);
        },
      },
      requestFilter: (req) => {
        // Inject a valid bearer token; the broker stores the consumer's request as-is,
        // and the provider may need credentials the consumer does not own.
        req.headers["authorization"] = `Bearer ${process.env.PROVIDER_TEST_TOKEN}`;
      },
      timeout: 30000,
    }).verifyProvider();
  }, 60000);
});
```

### State handler rules

- Each handler MUST be IDEMPOTENT — reset before seed.
- Use `setup` / `teardown` form when teardown is non-trivial; bare `async` form when the handler resets state on each call.
- Names are FREE-FORM strings; consumer + provider must agree on them. Document the catalog in the contract test charter.
- Avoid coupling state to a specific user ID — prefer "an active user exists" with a known fixture.

## 4. Async Message Contract

For Kafka, RabbitMQ, SNS/SQS, etc. The consumer test asserts that the message HANDLER can process a message of the given shape. The provider test asserts that the producer code emits messages matching the pact.

```typescript
// consumer/order-events.pact.test.ts
import { MessageConsumerPact, asynchronousBodyHandler, MatchersV3 } from "@pact-foundation/pact";
import { handleOrderCreated } from "../src/handlers/order-created";

const { like, integer, string, iso8601DateTime } = MatchersV3;

const messagePact = new MessageConsumerPact({
  consumer: "WarehouseService",
  provider: "OrderService",
  dir: "./pacts",
});

describe("OrderCreated message contract", () => {
  it("processes an OrderCreated event", () => {
    return messagePact
      .given("a customer placed an order")
      .expectsToReceive("an OrderCreated event")
      .withContent(
        like({
          eventType: string("OrderCreated"),
          orderId: integer(42),
          customerId: integer(7),
          createdAt: iso8601DateTime("2026-01-01T00:00:00Z"),
          items: [{ sku: string("SKU-1"), qty: integer(2) }],
        }),
      )
      .withMetadata({ "content-type": "application/json", "kafka.topic": "orders.events" })
      .verify(asynchronousBodyHandler(handleOrderCreated));
  });
});
```

Provider side: feed a function that returns the message envelope to `MessageProviderPact.verify`. No real Kafka needed.

## 5. CI Pipeline (GitHub Actions sketch)

```yaml
# .github/workflows/pact-consumer.yml
name: pact-consumer
on: [push, pull_request]
jobs:
  consumer:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with: { node-version: 20 }
      - run: npm ci
      - run: npm run test:pact:consumer
      - name: Publish pact
        if: github.ref == 'refs/heads/main' || github.event_name == 'pull_request'
        env:
          PACT_BROKER_BASE_URL: ${{ secrets.PACT_BROKER_BASE_URL }}
          PACT_BROKER_TOKEN:    ${{ secrets.PACT_BROKER_TOKEN }}
          GIT_SHA:    ${{ github.sha }}
          GIT_BRANCH: ${{ github.head_ref || github.ref_name }}
        run: npm run pact:publish
      - name: can-i-deploy
        if: github.ref == 'refs/heads/main'
        env:
          PACT_BROKER_BASE_URL: ${{ secrets.PACT_BROKER_BASE_URL }}
          PACT_BROKER_TOKEN:    ${{ secrets.PACT_BROKER_TOKEN }}
          GIT_SHA: ${{ github.sha }}
        run: npm run pact:can-i-deploy
```

## 6. Common Pitfalls

- **Hard-coded values without matchers** → `body: { id: 1 }` asserts EXACTLY `1`. Real responses with `id: 2` will fail. Wrap in `integer(1)`.
- **Asserting fields the consumer does not read** → makes the contract brittle. Drop them.
- **Provider not idempotent in state handlers** → flaky verifications. Always `resetDb()` first.
- **Forgetting `consumerVersionSelectors`** → provider verifies every pact ever published. Use `mainBranch + deployedOrReleased + matchingBranch`.
- **Using `PactV2` for new work** → V3 is the modern format with proper matchers. Use `PactV3`.
- **Running real provider against pact mock** in consumer tests → defeats the purpose. The mock server is the boundary.
- **Skipping `publishVerificationResult: true`** → broker never learns the verification passed; `can-i-deploy` returns false.
