// pact-jvm-provider.java — runnable starter
//
// Starter Pact provider verification using au.com.dius.pact:provider:junit5 4.6+
// with Spring Boot. Maven deps:
//   <dependency>
//     <groupId>au.com.dius.pact.provider</groupId>
//     <artifactId>junit5</artifactId>
//     <version>4.6.14</version>
//     <scope>test</scope>
//   </dependency>
//   <dependency>
//     <groupId>au.com.dius.pact.provider</groupId>
//     <artifactId>spring</artifactId>
//     <version>4.6.14</version>
//     <scope>test</scope>
//   </dependency>
//
// Run:
//   mvn -B test \
//     -Dpact.provider.version=$GIT_SHA \
//     -Dpact.provider.branch=$GIT_BRANCH \
//     -Dpact.verifier.publishResults=true

package com.example.pact.provider;

import au.com.dius.pact.provider.junit5.HttpTestTarget;
import au.com.dius.pact.provider.junit5.PactVerificationContext;
import au.com.dius.pact.provider.junit5.PactVerificationInvocationContextProvider;
import au.com.dius.pact.provider.junitsupport.Provider;
import au.com.dius.pact.provider.junitsupport.State;
import au.com.dius.pact.provider.junitsupport.loader.PactBroker;
import au.com.dius.pact.provider.junitsupport.loader.PactBrokerAuth;
import au.com.dius.pact.provider.junitsupport.loader.VersionSelector;
import com.example.fixtures.UserFixtures;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.TestTemplate;
import org.junit.jupiter.api.extension.ExtendWith;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.context.SpringBootTest.WebEnvironment;
import org.springframework.boot.web.server.LocalServerPort;

import java.util.List;

@SpringBootTest(webEnvironment = WebEnvironment.RANDOM_PORT)
@Provider("UserService")
@PactBroker(
    url = "${PACT_BROKER_BASE_URL}",
    authentication = @PactBrokerAuth(token = "${PACT_BROKER_TOKEN}"),
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

    @State("user 1 exists")
    void userOneExists() {
        fixtures.resetDb();
        fixtures.seedUser(1L, "Alice", "alice@example.com", List.of("admin"));
    }

    @State("user 999 does not exist")
    void userNineNineNineMissing() {
        fixtures.resetDb();
        fixtures.deleteUser(999L);
    }
}
