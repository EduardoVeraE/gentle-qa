# Installer Script

## Requirements

### Requirement: Platform Detection

The system MUST support macOS, Linux, AND Windows.

- GIVEN any supported platform (macOS, Linux, Windows)
- WHEN platform detection runs
- THEN it MUST set OS to one of: darwin, linux, windows
- AND set ARCH to one of: amd64, arm64

#### Scenario: Detect macOS Apple Silicon

- GIVEN a user runs install.sh on macOS with Apple Silicon
- WHEN the script executes platform detection
- THEN it MUST set OS="darwin", OS_LABEL="macOS", ARCH="arm64"

#### Scenario: Detect Linux x86_64

- GIVEN a user runs install.sh on Linux on x86_64
- WHEN the script executes platform detection
- THEN it MUST set OS="linux", OS_LABEL="Linux", ARCH="amd64"

### Requirement: Windows Platform Detection

The installer script MUST detect Windows environments (MINGW, MSYS, CYGWIN) and select appropriate binary format.

The system SHALL recognize Windows by checking `OSTYPE` environment variable or `uname -s` output containing "MINGW", "MSYS", or "CYGWIN".

#### Scenario: Detect Windows on MINGW64

- GIVEN a user runs install.sh on Windows with MINGW64 terminal
- WHEN the script executes platform detection
- THEN it MUST set OS="windows" and OS_LABEL="Windows"

#### Scenario: Detect Windows on MSYS2

- GIVEN a user runs install.sh on Windows with MSYS2 terminal
- WHEN the script executes platform detection
- THEN it MUST set OS="windows" and OS_LABEL="Windows"

### Requirement: Windows Binary Format Selection

The installer script MUST download the correct binary format for Windows (.zip with .exe extension).

The system SHALL use `.zip` archive format instead of `.tar.gz` when OS is Windows.

The system SHALL reference the `.exe` binary file in the download and extraction logic.

#### Scenario: Download Windows binary

- GIVEN the platform is detected as Windows
- WHEN the installer fetches the release archive
- THEN it MUST download `gentle-qa_{version}_windows_{arch}.zip`
- AND the binary path MUST reference `gentle-qa.exe`

### Requirement: Windows Installation Path

The installer script MUST install to a valid Windows directory.

The system SHALL use `%LOCALAPPDATA%\Programs\gentle-qa` or `%USERPROFILE%\.local\bin` as fallback.

#### Scenario: Install to user directory on Windows

- GIVEN the user has write permissions to `%LOCALAPPDATA%\Programs`
- WHEN the installer copies the binary
- THEN it MUST copy to `%LOCALAPPDATA%\Programs\gentle-qa\gentle-qa.exe`
