# Homebrew Distribution

## Requirements

### Requirement: Homebrew Formula Generation

GoReleaser MUST generate a valid Homebrew formula in the `homebrew-tap` repository.

The system SHALL create `Formula/gentle-qa.rb` with binary URL pointing to GitHub Releases.

#### Scenario: Formula created on release

- GIVEN a new GitHub release is published with version tag
- WHEN GoReleaser runs with brews configuration
- THEN it MUST commit `gentle-qa.rb` to `EduardoVeraE/homebrew-tap/Formula/`

### Requirement: Homebrew Installation

The Homebrew formula MUST allow installation via `brew install gentle-qa`.

The formula SHALL point to the pre-built binary assets from GitHub Releases (no build from source).

#### Scenario: Install via Homebrew on macOS

- GIVEN the formula exists in homebrew-tap
- WHEN user runs `brew install gentle-qa` on macOS
- THEN it MUST download and install the darwin binary
- AND the binary MUST be executable

#### Scenario: Install via Homebrew on Linux

- GIVEN the formula exists in homebrew-tap
- WHEN user runs `brew install gentle-qa` on Linux
- THEN it MUST download and install the linux binary
- AND the binary MUST be executable

### Requirement: Homebrew Upgrade

The Homebrew formula MUST support upgrades via `brew upgrade gentle-qa`.

The formula SHALL use dynamic version resolution to fetch the latest release URL.

#### Scenario: Upgrade to new version

- GIVEN gentle-qa v0.1.0 is installed via Homebrew
- WHEN a new release v0.1.1 is published
- AND user runs `brew upgrade gentle-qa`
- THEN it MUST download and replace with v0.1.1 binary
