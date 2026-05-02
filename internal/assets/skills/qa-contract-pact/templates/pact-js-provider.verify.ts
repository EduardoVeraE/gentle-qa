// pact-js-provider.verify.ts — runnable starter
//
// Starter Pact provider verification using @pact-foundation/pact v12+.
// Replace provider name, server bootstrap, and state handlers with your own.
//
// Run:
//   PACT_BROKER_BASE_URL=... PACT_BROKER_TOKEN=... GIT_SHA=... GIT_BRANCH=main \
//     npx jest pact-js-provider.verify.ts
//
// Required env:
//   PACT_BROKER_BASE_URL  Broker / PactFlow URL
//   PACT_BROKER_TOKEN     Broker token with verify-publish permission
//   GIT_SHA               Provider version (immutable git SHA)
//   GIT_BRANCH            Provider branch
//   PROVIDER_TEST_TOKEN   (optional) bearer token injected via requestFilter

import { Verifier } from "@pact-foundation/pact";

// Replace with your real server bootstrap. Must return a base URL and a close fn.
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
            await seedUser({
              id: 1,
              name: "Alice",
              email: "alice@example.com",
              roles: ["admin"],
            });
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
        // Inject a valid bearer if the provider requires auth that the consumer
        // does not own (typical for cross-team setups).
        if (process.env.PROVIDER_TEST_TOKEN) {
          req.headers["authorization"] = `Bearer ${process.env.PROVIDER_TEST_TOKEN}`;
        }
      },
      timeout: 30000,
    }).verifyProvider();
  }, 60000);
});
