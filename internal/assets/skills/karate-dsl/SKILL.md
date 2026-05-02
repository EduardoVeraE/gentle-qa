---
name: karate-dsl
description: >
  API testing with Karate DSL for SDET/QE Engineers. ISTQB Integration and System test level.
  Trigger: When writing Karate feature files, API contract tests, schema validations,
  security tests, data-driven API tests, or setting up Karate mocks.
license: Apache-2.0
metadata:
  author: gentle-qa
  version: "1.1"
---

## ISTQB Mapping

| Aspect | Value |
|--------|-------|
| Test level | Integration Testing, System Testing |
| Test type | Functional, Security, Structural (contract) |
| Techniques | Equivalence Partitioning, Boundary Value Analysis, Decision Tables, State Transition |
| Test oracle | OpenAPI spec / contract = expected API behavior |

**Core principle**: Every Karate scenario must have a clear oracle — what EXACTLY defines a correct response? Status code alone is not enough. Schema + values + business rules = complete oracle.

---

## When to Use

- API endpoint behavior (CRUD, auth, error handling)
- Strict-schema contract checks within a single Karate suite (`match response == { ... }` for breaking-change detection)
- Schema validation against OpenAPI spec
- Business rule assertions (calculations, state machines)
- Security surface testing (auth, injection, exposure)
- API mocking for downstream service isolation

**Not here**: UI flows (use playwright-bdd), performance (use k6), unit logic, consumer-driven contract testing across teams (use `qa-contract-pact` for Pact / broker / can-i-deploy / message contracts).

---

## Critical Patterns

### Pattern 1: ISTQB Techniques Applied to API Testing

```gherkin
# Equivalence Partitioning — group inputs into valid/invalid classes
# Valid class: price > 0 → test ONE representative (e.g., 50.00)
# Invalid classes: price = 0, price < 0, price = null, price = string
Scenario Outline: Price validation — equivalence partitions
  Given url baseUrl
  And path '/api/products'
  And request { name: 'Test', price: <price> }
  When method POST
  Then status <expected_status>

  Examples:
    | price  | expected_status | partition        |
    | 50.00  | 201             | valid — normal   |
    | 0.01   | 201             | valid — min BVA  |
    | 0      | 400             | invalid — zero   |
    | -1     | 400             | invalid — below  |
    | 999999 | 201             | valid — max BVA  |
    | 1000000| 400             | invalid — above  |
```

### Pattern 2: Basic API Test — Complete Oracle

```gherkin
# src/test/java/api/products/get-product.feature
Feature: Get Product API

  Background:
    * url baseUrl
    * def token = call read('classpath:auth/get-token.feature')
    * header Authorization = 'Bearer ' + token.value
    * header Content-Type = 'application/json'

  Scenario: Get existing product — full oracle
    Given path '/api/products/1'
    When method GET
    Then status 200
    # Oracle layer 1: schema (types and structure)
    And match response.id == '#number'
    And match response.name == '#string'
    And match response.price == '#? _ > 0'
    And match response.stock == '#? _ >= 0'
    And match response.isActive == '#boolean'
    # Oracle layer 2: business rules (values)
    And assert response.price > 0
    And assert response.currency.length == 3
    # Oracle layer 3: relationship integrity
    And match response.category.id == '#number'
    And match response.category.slug == '#regex [a-z-]+'

  Scenario: Get non-existent product — error contract
    Given path '/api/products/99999'
    When method GET
    Then status 404
    And match response.error == '#string'
    And match response.code == 'PRODUCT_NOT_FOUND'
    And match response !contains { stack: '#present' }  # Security: no stack trace
```

### Pattern 3: Contract Testing — Strict Schema

```gherkin
# src/test/java/contracts/product-contract.feature
Feature: Product API Contract

  Scenario: Product response matches published contract
    Given url baseUrl
    And path '/api/products/1'
    When method GET
    Then status 200

    # Strict match — fails on missing OR extra fields (breaking change detection)
    And match response ==
      """
      {
        "id": "#number",
        "name": "#string",
        "description": "#string",
        "price": "#? _ > 0",
        "currency": "#regex [A-Z]{3}",
        "stock": "#? _ >= 0",
        "category": {
          "id": "#number",
          "name": "#string",
          "slug": "#regex [a-z-]+"
        },
        "images": "#[] #string",
        "createdAt": "#? _ != null",
        "isActive": "#boolean"
      }
      """

  Scenario: Product list pagination contract
    Given url baseUrl
    And path '/api/products'
    When method GET
    Then status 200
    And match response.items == '#[] { id: #number, name: #string, price: #number }'
    And match response.total == '#number'
    And match response.page == '#? _ >= 1'
    And match response.pageSize == '#? _ > 0 && _ <= 100'
```

### Pattern 4: Value Assertions — Business Rules

```gherkin
# Decision Table technique: each row = one rule combination
Scenario: Order total calculation — Decision Table
  Given url baseUrl
  And path '/api/orders'
  And request { items: [{ productId: 1, quantity: 2 }], couponCode: 'SAVE10' }
  When method POST
  Then status 201

  # Rule 1: subtotal = price × quantity
  * def itemSubtotal = response.items[0].price * response.items[0].quantity
  And assert response.subtotal == itemSubtotal

  # Rule 2: 10% discount applied to subtotal
  * def expectedDiscount = itemSubtotal * 0.10
  And assert response.discount == expectedDiscount

  # Rule 3: total = subtotal - discount
  And assert response.total == itemSubtotal - expectedDiscount

  # Rule 4: status must be a valid state (State Transition — initial state)
  And match response.status == '#? ["pending","confirmed"].contains(_)'
```

### Pattern 5: Security Testing

```gherkin
# src/test/java/security/api-security.feature
Feature: API Security Surface

  Scenario: Unauthenticated request rejected
    Given url baseUrl
    And path '/api/orders'
    When method GET
    Then status 401
    And match response.error == 'Unauthorized'
    And match response !contains { token: '#present' }

  Scenario: SQL injection sanitized
    Given url baseUrl
    And param search = "' OR '1'='1"
    And path '/api/products'
    When method GET
    Then status 200
    And match response.items == '#[]'  # Empty result, not a DB dump

  Scenario: Mass assignment rejected — role escalation blocked
    Given url baseUrl
    And def userToken = call read('classpath:auth/get-token.feature') { role: 'user' }
    And header Authorization = 'Bearer ' + userToken.value
    And path '/api/users/me'
    And request { name: 'Updated', role: 'admin', isAdmin: true }
    When method PATCH
    Then status 200
    * def profile = call read('classpath:auth/get-profile.feature')
    And assert profile.role == 'user'  # Oracle: role MUST NOT change

  Scenario: Sensitive data not leaked in errors
    Given url baseUrl
    And path '/api/products/not-a-number'
    When method GET
    Then status 400
    And match response !contains { stack: '#present' }
    And match response !contains { query: '#present' }
    And match response !contains { dbError: '#present' }

  Scenario: Rate limiting enforced
    * def results = karate.repeat(30, function(i){ return karate.call('classpath:helpers/single-get.feature') })
    * def blocked = results.filter(function(r){ return r.status == 429 })
    And assert blocked.length > 0
```

### Pattern 6: Data-Driven Testing

```gherkin
# Scenario Outline: Equivalence Partitioning across input classes
Scenario Outline: Product search by category
  Given url baseUrl
  And path '/api/products'
  And param category = '<category>'
  When method GET
  Then status 200
  And match response.items[*].category.slug contains '<category>'
  And assert response.items.length > 0

  Examples:
    | category    |
    | electronics |
    | clothing    |
    | books       |
```

### Pattern 7: Karate Mock Server

```gherkin
# src/test/java/mocks/product-mock.feature
Feature: Product Service Mock

  Background:
    * configure cors = true

  Scenario: pathMatches('/api/products/{id}') && methodIs('get')
    * def response = { id: '#(pathParams.id)', name: 'Mock Product', price: 29.99 }
    * def responseStatus = 200

  Scenario: pathMatches('/api/products') && methodIs('post')
    * def response = { id: 999, createdAt: '2026-01-01T00:00:00Z' }
    * def responseStatus = 201

  Scenario: pathMatches('/api/products/{id}') && methodIs('delete')
    * def responseStatus = 204
```

### Pattern 8: Global Config

```javascript
// src/test/java/karate-config.js
function fn() {
  var env = karate.env || 'dev';
  var config = {
    baseUrl: 'http://localhost:8080',
    adminCreds: { email: 'admin@example.com', password: 'admin123' },
    userCreds:  { email: 'user@example.com',  password: 'user123'  },
  };

  if (env === 'staging') config.baseUrl = 'https://staging.example.com';

  // Never run mutation tests in prod
  if (env === 'prod') karate.configure('readTimeout', 10000);

  karate.configure('retry', { count: 2, interval: 1000 });
  karate.configure('ssl', true);
  return config;
}
```

---

## Anti-patterns — Never Do This

| Anti-pattern | Why it fails | Fix |
|---|---|---|
| `Then status 200` only | No schema, no values — not a test | Add schema + business rule assertions |
| `match response contains { id: '#present' }` | Partial match misses breaking changes | Use strict `match response == { ... }` for contracts |
| Hardcoded IDs (`productId: 1`) | Breaks on DB reset | Use dynamic setup or known fixture IDs |
| No auth header in security test | Tests the wrong path | Always test both authenticated AND unauthenticated |
| `* def token = 'hardcoded'` | Expires, wrong env | Use `call read('classpath:auth/get-token.feature')` |
| Testing only happy path | Miss EP invalid classes | Apply Equivalence Partitioning — test at least one invalid class per input |
| Ignoring error response body | Leaks internal details | Assert `!contains { stack: '#present' }` always |

---

## Test Oracle Checklist

Before marking an API test complete:
- [ ] Status code asserted
- [ ] Response schema validated (types, structure)
- [ ] Business rules asserted (calculated values, state transitions)
- [ ] Error responses asserted (correct code, no stack trace)
- [ ] Security: unauthenticated path tested for protected endpoints
- [ ] At least one invalid equivalence class tested per input field

---

## Decision Tree

```
What type of test?
├── Endpoint behavior → Basic feature (status + full oracle)
├── Service contract → Contract test (strict schema match)
├── Field types → Schema validation (inline or JSON schema)
├── Business rules → Value assertions (Decision Table technique)
├── Security posture → Security feature (auth, injection, exposure)
└── Multiple input classes → Scenario Outline (Equivalence Partitioning)

Contract fails?
├── Missing field → Producer must add (consumer-driven contract)
├── Wrong type → Align schema between services
├── Extra field → Warn only unless strict mode
└── Breaking change → Version the API endpoint

Security test fails?
├── 401 not returned → Add/fix auth middleware
├── Injection not sanitized → Add input validation
├── Stack trace exposed → Sanitize error responses
└── Rate limit missing → Add rate limiting middleware
```

---

## Project Structure

```
src/test/java/
├── karate-config.js
├── api/
│   ├── products/
│   └── orders/
├── contracts/
├── security/
├── mocks/
├── auth/
│   ├── get-token.feature
│   └── get-profile.feature
├── schemas/
├── data/
└── helpers/
    └── single-get.feature
```

---

## Commands

```bash
mvn test                                                              # All tests
mvn test -Dkarate.options="classpath:security"                        # Security only
mvn test -Dkarate.env=staging                                         # Staging env
mvn test -Dkarate.options="--tags @contract"                          # Contracts only
java -jar karate.jar -m src/test/java/mocks/product-mock.feature -p 8081  # Mock server
```

---

## Resources

- [Karate docs](https://karatelabs.github.io/karate/)
- [Consumer-Driven Contracts](https://martinfowler.com/articles/consumerDrivenContracts.html) — for true consumer-driven workflows use the `qa-contract-pact` skill (Pact, broker, can-i-deploy)
- [OWASP API Security Top 10](https://owasp.org/www-project-api-security/)
- [ISTQB Glossary](https://glossary.istqb.org)
