# PACT-JVM — Consumer + Provider (Java / Kotlin / Scala)

Reference for `au.com.dius.pact` v4.6+ on JDK 21+. Examples target Java + JUnit 5 with Maven and Gradle. Kotlin/Scala bindings follow the same API surface.

## 1. Install

### Maven

```xml
<dependencies>
  <!-- Consumer -->
  <dependency>
    <groupId>au.com.dius.pact.consumer</groupId>
    <artifactId>junit5</artifactId>
    <version>4.6.14</version>
    <scope>test</scope>
  </dependency>
  <!-- Provider -->
  <dependency>
    <groupId>au.com.dius.pact.provider</groupId>
    <artifactId>junit5</artifactId>
    <version>4.6.14</version>
    <scope>test</scope>
  </dependency>
  <!-- Provider Spring Boot integration (optional) -->
  <dependency>
    <groupId>au.com.dius.pact.provider</groupId>
    <artifactId>spring</artifactId>
    <version>4.6.14</version>
    <scope>test</scope>
  </dependency>
</dependencies>
```

### Gradle (Kotlin DSL)

```kotlin
plugins {
  id("au.com.dius.pact") version "4.6.14"
}

dependencies {
  testImplementation("au.com.dius.pact.consumer:junit5:4.6.14")
  testImplementation("au.com.dius.pact.provider:junit5:4.6.14")
  testImplementation("au.com.dius.pact.provider:spring:4.6.14")
}

pact {
  publish {
    pactBrokerUrl = System.getenv("PACT_BROKER_BASE_URL") ?: ""
    pactBrokerToken = System.getenv("PACT_BROKER_TOKEN") ?: ""
    consumerVersion = System.getenv("GIT_SHA") ?: "local"
    consumerBranch = System.getenv("GIT_BRANCH") ?: "local"
  }
}
```

## 2. HTTP Consumer Test (JUnit 5)

```java
package com.example.pact.consumer;

import au.com.dius.pact.consumer.MockServer;
import au.com.dius.pact.consumer.dsl.LambdaDsl;
import au.com.dius.pact.consumer.dsl.PactDslWithProvider;
import au.com.dius.pact.consumer.junit5.PactConsumerTestExt;
import au.com.dius.pact.consumer.junit5.PactTestFor;
import au.com.dius.pact.core.model.V4Pact;
import au.com.dius.pact.core.model.annotations.Pact;
import com.example.client.UserClient;
import com.example.client.User;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;

import java.util.Map;

import static org.assertj.core.api.Assertions.assertThat;

@ExtendWith(PactConsumerTestExt.class)
@PactTestFor(providerName = "UserService", pactVersion = au.com.dius.pact.core.model.PactSpecVersion.V4)
class UserClientPactTest {

    @Pact(consumer = "WebApp")
    public V4Pact userExistsPact(PactDslWithProvider builder) {
        return builder
            .given("user 1 exists")
            .uponReceiving("a request for user 1")
                .method("GET")
                .path("/api/users/1")
                .headers(Map.of("Accept", "application/json"))
            .willRespondWith()
                .status(200)
                .headers(Map.of("Content-Type", "application/json"))
                .body(LambdaDsl.newJsonBody(o -> {
                    o.numberType("id", 1);
                    o.stringType("name", "Alice");
                    o.stringMatcher(".+@.+\\..+", "email", "alice@example.com");
                    o.datetime("createdAt", "yyyy-MM-dd'T'HH:mm:ssXXX", java.time.Instant.parse("2026-01-01T00:00:00Z").atZone(java.time.ZoneOffset.UTC));
                    o.minArrayLike("roles", 1, PactDslJsonRootValue -> {}, 1);
                }).build())
            .toPact(V4Pact.class);
    }

    @Test
    @PactTestFor(pactMethod = "userExistsPact")
    void fetchesUserById(MockServer mockServer) {
        UserClient client = new UserClient(mockServer.getUrl());
        User user = client.fetchUser(1L);
        assertThat(user.id()).isEqualTo(1L);
        assertThat(user.email()).matches(".+@.+\\..+");
    }
}
```

### Matcher cheat-sheet (LambdaDsl)

| Method | Asserts | Example |
| ------ | ------- | ------- |
| `numberType("id", 1)` | Any number | `o.numberType("id", 42)` |
| `stringType("name", "X")` | Any string | `o.stringType("name", "Alice")` |
| `stringMatcher(regex, key, sample)` | String matching regex | `o.stringMatcher(".+@.+", "email", "a@b.c")` |
| `booleanType("flag", true)` | Any boolean | — |
| `eachLike(key, sample)` / `minArrayLike(key, min, sample, count)` | Array of items | — |
| `datetime("when", pattern, sample)` | Timestamp | — |
| `uuid("id", sample)` | UUID | — |
| `stringValue("status", "ACTIVE")` | EXACT value | for enums only |

Same rule as PACT-JS: matchers for dynamic fields, exact values for enums and discriminators.

## 3. HTTP Provider Verification (JUnit 5 + Spring Boot)

```java
package com.example.pact.provider;

import au.com.dius.pact.provider.junit5.HttpTestTarget;
import au.com.dius.pact.provider.junit5.PactVerificationContext;
import au.com.dius.pact.provider.junitsupport.Provider;
import au.com.dius.pact.provider.junitsupport.State;
import au.com.dius.pact.provider.junitsupport.loader.PactBroker;
import au.com.dius.pact.provider.junitsupport.loader.VersionSelector;
import au.com.dius.pact.provider.junit5.PactVerificationInvocationContextProvider;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.TestTemplate;
import org.junit.jupiter.api.extension.ExtendWith;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.context.SpringBootTest.WebEnvironment;
import org.springframework.boot.web.server.LocalServerPort;

@SpringBootTest(webEnvironment = WebEnvironment.RANDOM_PORT)
@Provider("UserService")
@PactBroker(
    url = "${PACT_BROKER_BASE_URL}",
    authentication = @au.com.dius.pact.provider.junitsupport.loader.PactBrokerAuth(token = "${PACT_BROKER_TOKEN}"),
    consumerVersionSelectors = {
        @VersionSelector(mainBranch = "true"),
        @VersionSelector(deployedOrReleased = "true"),
        @VersionSelector(matchingBranch = "true")
    }
)
@ExtendWith(PactVerificationInvocationContextProvider.class)
class UserServiceProviderPactTest {

    @LocalServerPort int port;
    @Autowired UserFixtures fixtures;

    @BeforeEach
    void before(PactVerificationContext context) {
        context.setTarget(new HttpTestTarget("localhost", port));
    }

    @TestTemplate
    void pactVerificationTestTemplate(PactVerificationContext context) {
        context.verifyInteraction();
    }

    // State handlers — names MUST match consumer's `given(...)`
    @State("user 1 exists")
    void userOneExists() {
        fixtures.resetDb();
        fixtures.seedUser(1L, "Alice", "alice@example.com", java.util.List.of("admin"));
    }

    @State("user 999 does not exist")
    void userNineNineNineMissing() {
        fixtures.resetDb();
        fixtures.deleteUser(999L);
    }
}
```

System properties picked up by the verifier:

```bash
mvn test \
  -Dpact.provider.version=$GIT_SHA \
  -Dpact.provider.branch=$GIT_BRANCH \
  -Dpact.verifier.publishResults=true
```

## 4. Async Message Contract (PACT-JVM)

```java
package com.example.pact.consumer.message;

import au.com.dius.pact.consumer.MessagePactBuilder;
import au.com.dius.pact.consumer.dsl.PactDslJsonBody;
import au.com.dius.pact.consumer.junit5.PactConsumerTestExt;
import au.com.dius.pact.consumer.junit5.PactTestFor;
import au.com.dius.pact.consumer.junit5.ProviderType;
import au.com.dius.pact.core.model.annotations.Pact;
import au.com.dius.pact.core.model.messaging.MessagePact;
import com.example.handlers.OrderCreatedHandler;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;

import static org.assertj.core.api.Assertions.assertThat;

@ExtendWith(PactConsumerTestExt.class)
@PactTestFor(providerName = "OrderService", providerType = ProviderType.ASYNCH)
class OrderCreatedConsumerTest {

    @Pact(consumer = "WarehouseService")
    public MessagePact orderCreatedPact(au.com.dius.pact.consumer.MessagePactBuilder builder) {
        PactDslJsonBody body = new PactDslJsonBody()
            .stringValue("eventType", "OrderCreated")
            .numberType("orderId", 42)
            .numberType("customerId", 7)
            .datetime("createdAt", "yyyy-MM-dd'T'HH:mm:ssXXX")
            .minArrayLike("items", 1, new PactDslJsonBody()
                .stringType("sku", "SKU-1")
                .numberType("qty", 2));

        return builder
            .given("a customer placed an order")
            .expectsToReceive("an OrderCreated event")
            .withContent(body)
            .withMetadata(java.util.Map.of("kafka.topic", "orders.events"))
            .toPact();
    }

    @Test
    @PactTestFor(pactMethod = "orderCreatedPact")
    void handlesOrderCreated(java.util.List<au.com.dius.pact.core.model.messaging.Message> messages) {
        OrderCreatedHandler handler = new OrderCreatedHandler();
        for (var msg : messages) {
            assertThat(handler.handle(msg.contentsAsBytes())).isTrue();
        }
    }
}
```

## 5. CI Pipeline (Maven sketch)

```yaml
# .github/workflows/pact-consumer-jvm.yml
jobs:
  consumer:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-java@v4
        with: { distribution: temurin, java-version: 21 }
      - run: mvn -B test
      - name: Publish pacts
        if: github.event_name == 'push' || github.event_name == 'pull_request'
        env:
          PACT_BROKER_BASE_URL: ${{ secrets.PACT_BROKER_BASE_URL }}
          PACT_BROKER_TOKEN:    ${{ secrets.PACT_BROKER_TOKEN }}
          GIT_SHA:    ${{ github.sha }}
          GIT_BRANCH: ${{ github.head_ref || github.ref_name }}
        run: |
          mvn -B au.com.dius.pact.provider:maven-plugin:publish \
            -Dpact.broker.url=$PACT_BROKER_BASE_URL \
            -Dpact.broker.token=$PACT_BROKER_TOKEN \
            -Dpact.consumer.version=$GIT_SHA \
            -Dpact.consumer.tags=$GIT_BRANCH
      - name: can-i-deploy
        if: github.ref == 'refs/heads/main'
        env:
          PACT_BROKER_BASE_URL: ${{ secrets.PACT_BROKER_BASE_URL }}
          PACT_BROKER_TOKEN:    ${{ secrets.PACT_BROKER_TOKEN }}
        run: |
          npx -y @pact-foundation/pact-cli can-i-deploy \
            --pacticipant WebApp \
            --version ${{ github.sha }} \
            --to-environment production \
            --broker-base-url $PACT_BROKER_BASE_URL \
            --broker-token $PACT_BROKER_TOKEN
```

## 6. Spring Cloud Contract — When to Choose Instead

Spring Cloud Contract is **provider-driven** (Groovy DSL on the producer side). Pact is **consumer-driven**. They solve overlapping problems differently.

| Signal | Choose |
| ------ | ------ |
| Provider team owns the contract definition | Spring Cloud Contract |
| Consumers can shape provider responses | Pact |
| External / third-party consumers | Pact (cannot ask SCC of an external team) |
| Java/Spring shop, single org, internal | Either; SCC is lower friction |
| Polyglot consumers | Pact |

If the team is already on Spring Cloud Contract and it works, keep it. If they need consumer-driven (external consumers, polyglot, separate cadence), use Pact.

## 7. Common Pitfalls (JVM-specific)

- **Mixing Pact V3 / V4 specs across consumer and provider** → use `PactSpecVersion.V4` consistently; older V3 code is still supported but matchers differ.
- **Forgetting `@PactBroker` on the provider** → falls back to local pact files; verifications never publish results.
- **State methods named differently than `given(...)`** → broker reports "no provider state" mismatches. Names are case-sensitive strings; document them in the contract test charter.
- **Random ports without `@LocalServerPort`** → verifier hits the wrong URL; tests pass against an empty server.
- **Hard-coded sample values without matchers (`stringValue` instead of `stringType`)** → contract is brittle; any provider response with a different sample value fails.
- **Spring Boot context not isolated per test** → state handlers from one consumer's pact bleed into the next. Use `@DirtiesContext` or scoped fixtures.
