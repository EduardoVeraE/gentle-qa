<div align="center">

  <img alt="Gentle-QA Banner" src="docs/banner.png" />

  <h1>Gentle-QA</h1>

  <p><strong>One command. Any test. Any framework. Your QA agents — configured and ready.</strong></p>

  <p>
    <a href="https://github.com/EduardoVeraE/Gentle-QA/releases"><img src="https://img.shields.io/github/v/release/EduardoVeraE/Gentle-QA" alt="Release"></a>
    <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License: MIT"></a>
    <img src="https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go&logoColor=white" alt="Go 1.24+" />
    <img src="https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey" alt="Platform" />
  </p>
</div>

---

## What It Does

This is NOT just an AI agent installer. This is a **QA ecosystem configurator** — it takes whatever AI coding agent(s) you use and supercharges them with a complete quality engineering stack: persistent memory, Spec-Driven Development workflow, curated QA skills (Playwright, Karate DSL, k6, BDD), an SDET-oriented persona, MCP servers, and per-phase model assignment so each SDD step runs on the right model.

**Before**: "I have Claude Code / OpenCode / Cursor, but it's just a chatbot that writes code."

**After**: Your agent has memory, testing skills, SDD workflow, and an SDET persona that helps you ship quality software — not just code.

### 10 Supported Agents

| Agent | Delegation Model | Key Feature |
|-------|:---:|---|
| **Claude Code** | Full (Task tool) | Sub-agents, output styles |
| **OpenCode** | Full (multi-mode overlay) | Per-phase model routing |
| **Gemini CLI** | Full (experimental) | Custom agents in `~/.gemini/agents/` |
| **Cursor** | Full (native subagents) | 9 SDD agents in `~/.cursor/agents/` |
| **VS Code Copilot** | Full (runSubagent) | Parallel execution |
| **Codex** | Solo-agent | CLI-native, TOML config |
| **Windsurf** | Solo-agent | Plan Mode, Code Mode, native workflows |
| **Antigravity** | Solo-agent + Mission Control | Built-in Browser/Terminal sub-agents |
| **Kiro IDE** | Full (native subagents) | Native `~/.kiro/agents/` + steering orchestration |
| **Qwen Code** | Full (native sub-agents) | Slash commands, `~/.qwen/commands/`, `auto_edit` mode |

---

## Quick Start

### macOS / Linux

```bash
curl -fsSL https://raw.githubusercontent.com/EduardoVeraE/Gentle-QA/main/scripts/install.sh | bash
```

### Windows

```powershell
scoop bucket add gentle-qa https://github.com/EduardoVeraE/scoop-bucket
scoop install gentle-qa
```

Or via PowerShell script: `irm https://raw.githubusercontent.com/EduardoVeraE/Gentle-QA/main/scripts/install.ps1 | iex`

### After install: project-level setup

Once your agents are configured, open your AI agent in a project and run these two commands to register the project context:

| Command | What it does | When to re-run |
|---------|-------------|----------------|
| `/sdd-init` | Detects stack, testing capabilities, activates Strict TDD Mode if available | When your project adds/removes test frameworks, or first time in a new project |
| `skill-registry` | Scans installed skills and project conventions, builds the registry | After installing/removing skills, or first time in a new project |

These are **not required** for basic usage. The SDD orchestrator runs `/sdd-init` automatically if it detects no context. But if something changed in your project (new test runner, new dependencies), re-running them manually ensures the agents have up-to-date context.

---

## Install

### Recommended

```bash
# macOS / Linux
brew tap EduardoVeraE/tap
brew install gentle-qa

# Windows
scoop bucket add gentle-qa https://github.com/EduardoVeraE/scoop-bucket
scoop install gentle-qa
```

<details>
<summary><strong>Other install methods</strong> (Go install, PowerShell script, binary download)</summary>

#### Go install (any platform with Go 1.24+)

```bash
go install github.com/EduardoVeraE/Gentle-QA/cmd/gentle-qa@latest
```

#### Windows (PowerShell script)

```powershell
irm https://raw.githubusercontent.com/EduardoVeraE/Gentle-QA/main/scripts/install.ps1 | iex
```

**Migrating from PowerShell installer to Scoop?** Remove the old binary first:

```powershell
Remove-Item "$env:LOCALAPPDATA\gentle-qa" -Recurse -Force
# Then install via Scoop as shown above
```

#### From releases

Download the binary for your platform from [GitHub Releases](https://github.com/EduardoVeraE/Gentle-QA/releases).

</details>

---

## QA Skills Included

Gentle-QA ships with **20+ curated skills** across every layer of the quality engineering stack. Your AI agent gets context-aware assistance for all of them out of the box — skills auto-activate based on the files you touch and the task at hand.

> The **SDET persona** orients your agent toward quality-first thinking: test pyramid strategy, shift-left practices, risk-based coverage, and CI/CD gate design.

---

### Playwright Testing

| Skill | When to use it | What it covers |
|-------|---------------|----------------|
| **playwright-e2e-testing** | Writing E2E test SUITES for UI flows, auth, forms | Page Object Model, fixtures, network interception, screenshot comparison, parallel execution |
| **playwright-bdd** | BDD/ATDD tests (Given/When/Then) | Cucumber/Gherkin integration, fixture composition, ISTQB acceptance level |
| **playwright-regression-strategy** | Planning/organizing regression suites, debugging flaky tests | Tier model, risk- and change-based selection, sharding, CI/CD, flake quarantine |
| **a11y-playwright-testing** | Accessibility testing (WCAG 2.1/2.2 AA) | ARIA patterns, axe-core integration, WCAG checklist, POUR principles |
| **playwright-mcp-inspect** | Live browser inspection in a running session | MCP-driven navigation, interactive debugging, screenshot capture, console log inspection |
| **playwright-cli** | Scripted exploratory testing from CLI | Session evidence (trace/screenshot/video), state save/load, network mocking, multi-role |

---

### API Testing

| Skill | When to use it | What it covers |
|-------|---------------|----------------|
| **api-testing** | Creating or debugging API tests (REST, GraphQL) | Schema validation (Zod/JSON Schema), auth flows (OAuth2, JWT), error state coverage, contract testing, idempotency — Playwright TypeScript and REST Assured (Java 21+) |

---

### Selenium Testing

| Skill | When to use it | What it covers |
|-------|---------------|----------------|
| **selenium-e2e-testing** | Writing E2E test SUITES with Selenium WebDriver + Java/JUnit 5 | Locator strategies, Page Object Model, explicit waits, AssertJ, Maven scaffolding |
| **a11y-selenium-testing** | Accessibility testing with Selenium + axe-core | WCAG 2.1/2.2 AA checklist, keyboard nav, ARIA semantics, POUR principles |

---

### Load & Performance Testing

| Skill | When to use it | What it covers |
|-------|---------------|----------------|
| **k6-load-test** | Performance and load testing | Virtual users, thresholds, scenarios, cloud execution |

---

### Manual QA & Planning

| Skill | When to use it | What it covers |
|-------|---------------|----------------|
| **qa-manual-istqb** | ISTQB-aligned test planning and design (canonical) | Full lifecycle: test planning → analysis → design → execution → completion. Techniques: EP, BVA, decision tables, state transitions. Templates: test plan, bug report, regression suite, exploratory charter, traceability matrix |

---

### Other Frameworks

| Skill | When to use it | What it covers |
|-------|---------------|----------------|
| **karate-dsl** | API testing with Karate | BDD-style API tests, mocking, GraphQL, performance scenarios |

---

### Workflow & Tooling

| Skill | When to use it | What it covers |
|-------|---------------|----------------|
| **go-testing** | Writing Go tests (unit, integration, Bubbletea TUI) | teatest, mocking patterns, TUI test harness |
| **skill-creator** | Creating a new skill for the project | Skill scaffolding, metadata, compact rules |
| **skill-registry** | Updating the skill registry after adding/removing skills | Scans all skills, writes `.atl/skill-registry.md`, saves to Engram |
| **branch-pr** | Creating GitHub PRs | Issue-first enforcement, PR template, structured description |
| **issue-creation** | Reporting bugs or requesting features on GitHub | Issue template, structured format |
| **judgment-day** | Adversarial code review | Launches two independent blind judge agents in parallel, synthesizes findings, auto-fixes, re-judges |

---

### Recommended Usage by Scenario

| Scenario | Skills to activate |
|----------|-------------------|
| Writing E2E tests for a web app | `playwright-e2e-testing` + `qa-manual-istqb` |
| Testing REST/GraphQL APIs | `api-testing` + `qa-manual-istqb` |
| Accessibility audit (Playwright) | `a11y-playwright-testing` |
| Accessibility audit (Selenium) | `a11y-selenium-testing` |
| Performance / load testing | `k6-load-test` |
| ISTQB-aligned test plan & design | `qa-manual-istqb` |
| Debugging flaky tests / regression suite design | `playwright-regression-strategy` |
| Live browser debugging (no test file) | `playwright-mcp-inspect` |
| Exploratory testing with evidence | `playwright-cli` |
| Structured change (new feature, big refactor) | `/sdd-init` → `/sdd-explore` → `/sdd-propose` → `/sdd-apply` |
| Code review with quality bar | `judgment-day` |

Skills activate automatically based on context. You can also invoke any skill explicitly: `/playwright-e2e-testing`, `/qa-manual-istqb`, `/api-testing`, etc.

---

## Backups

Every install, sync, and upgrade automatically snapshots your config files. Backups are **compressed** (tar.gz), **deduplicated** (identical configs are not re-backed up), and **auto-pruned** (keeps the 5 most recent). Pin important backups via the TUI (`p` key) to protect them from pruning.

See [Backup & Rollback Guide](docs/rollback.md) for details.

---

## Key Features You Should Know About

### OpenCode SDD Profiles

Assign different AI models to different SDD phases — a powerful model for test design, a fast one for implementation, a cheap one for exploration. Create multiple profiles and switch between them with Tab in OpenCode.

```bash
# Via CLI
gentle-qa sync --profile cheap:openrouter/qwen/qwen3-30b-a3b:free
gentle-qa sync --profile-phase cheap:sdd-design:anthropic/claude-sonnet-4-20250514

# Or via TUI: gentle-qa → "OpenCode SDD Profiles" → Create
```

After creating a profile, open OpenCode and press **Tab** to switch between `sdd-orchestrator` (default) and your custom profiles.

**Full guide**: [OpenCode SDD Profiles](docs/opencode-profiles.md)

### OpenCode Community Plugins

Gentle-QA can install curated community plugins for the OpenCode TUI alongside the SDD foundation:

- **sub-agent-statusline** — see which sub-agent is active in the OpenCode statusline
- **sdd-engram-plugin** — manage SDD profiles and browse Engram memories from inside OpenCode (runtime profile activation, no restart)

Pick them in the TUI under **OpenCode Community Plugins**, or pass them through `gentle-qa sync` for unattended setups.

### Per-Phase Claude Models

When Claude Code is your active agent, every SDD phase (`sdd-explore`, `sdd-propose`, `sdd-spec`, `sdd-design`, `sdd-tasks`, `sdd-apply`, `sdd-verify`, `sdd-archive`) can be pinned to a specific Claude alias (`opus`, `sonnet`, `haiku`). The orchestrator reads the assignments table and forwards the right `model` parameter on every delegation — no manual switching during a session. Configure it from the TUI's **Configure Models** screen or via `--claude-model`/`--kiro-model` flags.

### Engram (Persistent Memory)

Your AI agent automatically remembers test strategy decisions, bug patterns, and project context across sessions. You don't need to do anything — but when you do:

```bash
engram projects list          # See all projects with memory counts
engram projects consolidate   # Fix name drift ("my-app" vs "My-App")
engram search "flaky test"    # Find a past fix from the terminal
engram tui                    # Visual memory browser
```

**Full reference**: [Engram Commands](docs/engram.md)

---

## Documentation

| Topic | Description |
|-------|-------------|
| [Intended Usage](docs/intended-usage.md) | How gentle-qa is meant to be used — the mental model |
| [OpenCode SDD Profiles](docs/opencode-profiles.md) | Create and manage per-phase model profiles for OpenCode |
| [Engram Commands](docs/engram.md) | CLI commands, MCP tools, project management, team sharing |
| [Agents](docs/agents.md) | Supported agents, feature matrix, config paths, and per-agent notes |
| [Components, Skills & Presets](docs/components.md) | All components, GGA behavior, skill catalog, and preset definitions |
| [Usage](docs/usage.md) | Persona modes, interactive TUI, CLI flags, and dependency management |
| [Backup & Rollback](docs/rollback.md) | Backup retention, compression, dedup, pinning, and restore |
| [Kiro IDE](docs/kiro.md) | Kiro-specific setup, config paths, native subagents, and SDD behavior |
| [Platforms](docs/platforms.md) | Supported platforms, Windows notes, security verification, config paths |
| [Architecture & Development](docs/architecture.md) | Codebase layout, testing, and relationship to Gentleman.Dots |

---

## Community Highlights

This project gets better when the community builds on top of it.

### Community Integrations

- [sdd-engram-plugin](https://github.com/j0k3r-dev-rgl/sdd-engram-plugin) — manage OpenCode SDD profiles and browse Engram memories directly from OpenCode, with runtime profile activation and no restart required.

Using a tool like this? Gentle-QA now supports a safer OpenCode sync compatibility path for external single-active profile managers.

## Contributors

This project exists because of the community. See [CONTRIBUTORS.md](CONTRIBUTORS.md) for the full list.

<a href="https://github.com/EduardoVeraE/Gentle-QA/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=EduardoVeraE/Gentle-QA" />
</a>

---

## Next Steps

- **Just installed?** Read [Intended Usage](docs/intended-usage.md) — the one page that explains the mental model.
- **Using OpenCode?** Set up [SDD Profiles](docs/opencode-profiles.md) to assign different models per phase.
- **Want to share memory across machines?** Learn `engram sync` in the [Engram reference](docs/engram.md).
- **Ready to contribute?** Check [CONTRIBUTING.md](CONTRIBUTING.md) and the [open issues](https://github.com/EduardoVeraE/Gentle-QA/issues?q=is%3Aissue+is%3Aopen+label%3A%22status%3Aapproved%22).

---

<div align="center">
<a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License: MIT"></a>
</div>
