// pact-jvm-consumer.java — runnable starter
//
// Starter Pact consumer test using au.com.dius.pact:consumer:junit5 4.6+.
// Maven dependency:
//   <dependency>
//     <groupId>au.com.dius.pact.consumer</groupId>
//     <artifactId>junit5</artifactId>
//     <version>4.6.14</version>
//     <scope>test</scope>
//   </dependency>
//
// Run:
//   mvn -B test
// Output:
//   target/pacts/<consumer>-<provider>.json

package com.example.pact.consumer;

import au.com.dius.pact.consumer.MockServer;
import au.com.dius.pact.consumer.dsl.LambdaDsl;
import au.com.dius.pact.consumer.dsl.PactDslWithProvider;
import au.com.dius.pact.consumer.junit5.PactConsumerTestExt;
import au.com.dius.pact.consumer.junit5.PactTestFor;
import au.com.dius.pact.core.model.PactSpecVersion;
import au.com.dius.pact.core.model.V4Pact;
import au.com.dius.pact.core.model.annotations.Pact;
import com.example.client.User;
import com.example.client.UserClient;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;

import java.util.Map;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

@ExtendWith(PactConsumerTestExt.class)
@PactTestFor(providerName = "UserService", pactVersion = PactSpecVersion.V4)
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
                    o.minArrayLike("roles", 1, item -> item.stringType("admin"), 1);
                }).build())
            .toPact(V4Pact.class);
    }

    @Pact(consumer = "WebApp")
    public V4Pact userMissingPact(PactDslWithProvider builder) {
        return builder
            .given("user 999 does not exist")
            .uponReceiving("a request for missing user 999")
                .method("GET")
                .path("/api/users/999")
            .willRespondWith()
                .status(404)
                .headers(Map.of("Content-Type", "application/problem+json"))
                .body(LambdaDsl.newJsonBody(o -> {
                    o.stringType("type", "about:blank");
                    o.stringType("title", "Not Found");
                    o.numberType("status", 404);
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

    @Test
    @PactTestFor(pactMethod = "userMissingPact")
    void throwsWhenUserMissing(MockServer mockServer) {
        UserClient client = new UserClient(mockServer.getUrl());
        assertThatThrownBy(() -> client.fetchUser(999L))
            .hasMessageContaining("not found");
    }
}
