# Components, Skills & Presets

← [Back to README](../README.md)

---

## Components

| Component | ID | Description |
|-----------|-----|-------------|
| Engram | `engram` | Persistent cross-session memory via MCP — auto-detection of project name, full-text search, git sync, project consolidation. See [engram repo](https://github.com/Gentleman-Programming/engram) |
| SDD | `sdd` | Spec-Driven Development workflow (9 phases) — the agent handles SDD organically when the task warrants it, or when you ask; you don't need to learn the commands |
| Skills | `skills` | Curated coding skill library |
| Context7 | `context7` | MCP server for live framework/library documentation |
| Persona | `persona` | Managed Gentleman/neutral persona injection, or unmanaged custom persona mode |
| Permissions | `permissions` | Security-first defaults and guardrails |
| GGA | `gga` | Gentleman Guardian Angel — AI provider switcher |
| Theme | `theme` | Gentleman Kanagawa theme overlay |

## GGA Behavior

`gentle-qa --component gga` installs/provisions the `gga` binary globally on your machine.

It does **not** run project-level hook setup automatically (`gga init` / `gga install`) because that should be an explicit decision per repository.

After global install, enable GGA per project with:

```bash
gga init
gga install
```

---

## Skills

### Included Skills (installed by gentle-qa)

34 skill directories organized by category, embedded in the binary and injected into your agent's configuration. Skills follow the **5-layer ISTQB taxonomy** (Foundation → Strategy → Functional-by-level → Non-functional-by-type → Tooling) with disjoint triggers and explicit exclusion clauses (`NOT for X — use Y`) to keep orchestrator routing unambiguous.

#### SDD (Spec-Driven Development)

| Skill | ID | Description |
|-------|-----|-------------|
| SDD Init | `sdd-init` | Bootstrap SDD context in a project |
| SDD Explore | `sdd-explore` | Investigate codebase before committing to a change |
| SDD Propose | `sdd-propose` | Create change proposal with intent, scope, approach |
| SDD Spec | `sdd-spec` | Write specifications with requirements and scenarios |
| SDD Design | `sdd-design` | Technical design with architecture decisions |
| SDD Tasks | `sdd-tasks` | Break down a change into implementation tasks |
| SDD Apply | `sdd-apply` | Implement tasks following specs and design |
| SDD Verify | `sdd-verify` | Validate implementation matches specs |
| SDD Archive | `sdd-archive` | Sync delta specs to main specs and archive |
| Judgment Day | `judgment-day` | Parallel adversarial review — two independent judges review the same target |

#### Foundation

| Skill | ID | Description |
|-------|-----|-------------|
| Go Testing | `go-testing` | Go testing patterns including Bubbletea TUI testing |
| Skill Creator | `skill-creator` | Create new AI agent skills following the Agent Skills spec |
| Branch & PR | `branch-pr` | PR creation workflow with conventional commits, branch naming, and issue-first enforcement |
| Issue Creation | `issue-creation` | Issue filing workflow with bug report and feature request templates |
| Upstream Sync | `upstream-sync` | Workflow for merging upstream `gentle-ai` into the Gentle-QA fork while preserving rebranding |

These foundation skills are installed by default with both `full-gentleman` and `ecosystem-only` presets.

#### QA / SDET (Gentle-QA fork)

These are the QA-focused skills shipped with the Gentle-QA fork, organized by ISTQB layer.

##### Foundation & Strategy

| Skill | ID | Description |
|-------|-----|-------------|
| Manual ISTQB | `qa-manual-istqb` | ISTQB-aligned test planning, analysis, design, execution, completion (Layer 1) |
| Regression Strategy | `playwright-regression-strategy` | Tier model, risk-based selection, sharding, flake quarantine (Layer 2) |

##### Functional — by level

| Skill | ID | Description |
|-------|-----|-------------|
| Playwright E2E | `playwright-e2e-testing` | E2E SUITES with Page Object Model, fixtures, network interception |
| Playwright BDD | `playwright-bdd` | BDD/ATDD with Cucumber/Gherkin (acceptance level) |
| Selenium E2E | `selenium-e2e-testing` | Selenium WebDriver + Java/JUnit 5 with explicit waits, AssertJ |
| API Testing | `api-testing` | REST/GraphQL with schema validation, OpenAPI-first, mandatory headers — Playwright TS + REST Assured (Java 21+) |
| Karate DSL | `karate-dsl` | BDD-style API tests, mocking, GraphQL, performance scenarios |
| Contract — PACT | `qa-contract-pact` | Consumer-driven contract testing with PACT-JS, PACT-JVM, broker, `can-i-deploy` (integration level) |

##### Non-functional — by type

| Skill | ID | Description |
|-------|-----|-------------|
| OWASP Security | `qa-owasp-security` | OWASP Web/API/Mobile testing, ZAP scans, deps scanning, validated against DVWA |
| A11y (Playwright) | `a11y-playwright-testing` | WCAG 2.1/2.2 AA with axe-core, ARIA patterns, POUR principles |
| A11y (Selenium) | `a11y-selenium-testing` | WCAG 2.1/2.2 AA with Selenium + axe-core |
| Visual Regression | `qa-visual-regression` | Percy / Chromatic / Playwright `toHaveScreenshot` with baseline workflow + CI gating |
| K6 Load Test | `k6-load-test` | Performance and load testing with virtual users, thresholds, scenarios |

##### Cross-cutting platform

| Skill | ID | Description |
|-------|-----|-------------|
| Mobile Testing | `qa-mobile-testing` | Appium (iOS/Android) + Detox (React Native), device strategy, gesture/wait patterns |

##### Tooling

| Skill | ID | Description |
|-------|-----|-------------|
| Playwright MCP Inspect | `playwright-mcp-inspect` | Live browser inspection in a running session via MCP |
| Playwright CLI | `playwright-cli` | Scripted exploratory testing with session evidence |

### Coding Skills (separate repository)

For framework-specific skills (React 19, Angular, TypeScript, Tailwind 4, Zod 4, Playwright, etc.), see [Gentleman-Programming/Gentleman-Skills](https://github.com/Gentleman-Programming/Gentleman-Skills). These are maintained by the community and installed separately by cloning the repo and copying skills to your agent's skills directory.

---

## Presets

| Preset | ID | What's Included |
|--------|-----|-----------------|
| Full Gentleman | `full-gentleman` | All components (Engram + SDD + Skills + Context7 + GGA + Persona + Permissions + Theme) + all skills + gentleman persona |
| Ecosystem Only | `ecosystem-only` | Core components (Engram + SDD + Skills + Context7 + GGA) + all skills + gentleman persona |
| Minimal | `minimal` | Engram + SDD skills only |
| Custom | `custom` | You choose components and skills manually while keeping any existing persona/settings unmanaged |
