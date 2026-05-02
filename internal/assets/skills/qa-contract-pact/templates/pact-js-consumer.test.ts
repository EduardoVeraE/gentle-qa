// pact-js-consumer.test.ts — runnable starter
//
// Starter Pact consumer test using @pact-foundation/pact v12+ (PactV3 spec).
// Replace consumer/provider names, paths, headers, and matchers with your own.
//
// Run:
//   npm i -D @pact-foundation/pact jest ts-jest @types/jest
//   npx jest pact-js-consumer.test.ts
// Output:
//   pacts/<consumer>-<provider>.json
//
// Then publish from CI:
//   npx pact-broker publish ./pacts \
//     --consumer-app-version=$GIT_SHA \
//     --branch=$GIT_BRANCH \
//     --broker-base-url=$PACT_BROKER_BASE_URL \
//     --broker-token=$PACT_BROKER_TOKEN

import path from "node:path";
import { PactV3, MatchersV3 } from "@pact-foundation/pact";

// Replace with your real client. The client MUST accept a baseURL so the test can
// point it at the Pact mock server URL.
import { fetchUser } from "../src/clients/user-client";

const { like, integer, string, regex, eachLike, iso8601DateTime } = MatchersV3;

const provider = new PactV3({
  consumer: "WebApp",
  provider: "UserService",
  dir: path.resolve(process.cwd(), "pacts"),
  logLevel: "warn",
});

describe("UserService consumer pact", () => {
  it("fetches a user by id", async () => {
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
        headers: {
          "Content-Type": regex(/application\/json.*/, "application/json; charset=utf-8"),
        },
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
        body: like({
          type: string("about:blank"),
          title: string("Not Found"),
          status: integer(404),
        }),
      });

    await provider.executeTest(async (mockServer) => {
      await expect(fetchUser(mockServer.url, 999)).rejects.toThrow(/not found/i);
    });
  });
});
