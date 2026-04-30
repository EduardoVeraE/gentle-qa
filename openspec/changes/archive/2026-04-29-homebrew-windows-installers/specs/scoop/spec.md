# Delta for Scoop Distribution

## ADDED Requirements

### Requirement: Scoop Bucket JSON Generation

GoReleaser MUST generate a valid Scoop manifest JSON in the `scoop-bucket` repository.

The system SHALL create `bucket/gentle-qa.json` with Windows binary URL and checksum.

#### Scenario: Scoop manifest created on release

- GIVEN a new GitHub release is published with version tag
- WHEN GoReleaser runs with scoops configuration
- THEN it MUST commit `gentle-qa.json` to `EduardoVeraE/scoop-bucket/bucket/`

### Requirement: Scoop Installation

The Scoop manifest MUST allow installation via `scoop install gentle-qa`.

The manifest SHALL point to the pre-built Windows binary from GitHub Releases.

#### Scenario: Install via Scoop on Windows

- GIVEN the manifest exists in scoop-bucket
- WHEN user runs `scoop install gentle-qa` on Windows
- THEN it MUST download and install the windows-amd64 binary
- AND the binary MUST be executable

### Requirement: Scoop Update

The Scoop manifest MUST support updates via `scoop update gentle-qa`.

The manifest SHALL include version and checksum for automatic update detection.

#### Scenario: Update to new version

- GIVEN gentle-qa v0.1.0 is installed via Scoop
- WHEN a new release v0.1.1 is published
- AND user runs `scoop update gentle-qa`
- THEN it MUST download and replace with v0.1.1 binary