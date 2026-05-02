# qa-mobile-testing — Validation Harness

Automated, deterministic validation of the qa-mobile-testing skill assets:

- The `mobile_artifacts.mjs` CLI — `list`, `help`, `create`, `--self-test`,
  and the two error paths (path traversal, unknown template).
- The example scaffolds the skill ships under `examples/`:
  - `appium-wdio-ts/` — TypeScript files parse cleanly; `tsconfig.json`,
    `package.json`, and `wdio.conf.ts` keep their canonical shape.
  - `detox-rn/` — JS files parse under `node --check`; `.detoxrc.js`
    exports the keys Detox 20 requires; `package.json` declares the
    canonical scripts.

The harness deliberately **does NOT** boot a real iOS Simulator or Android
emulator. Mobile testing requires real devices/sims that cannot be
bootstrapped deterministically across CI runners. The harness validates
what *is* deterministic and leaves device-bound runs to the user once they
`npm install` inside an example.

## Why this exists

Mirrors the `qa-owasp-security` harness (commit `e40581d`). Without a
deterministic gate every change to `mobile_artifacts.mjs` or the example
scaffolds is a roll of the dice — broken templates ship to users with no
warning. The CLI ships its own `--self-test`; the harness extends that
with shell-level coverage of error paths and example-scaffold integrity.

The lesson encoded throughout `lib.sh`: **never** capture an exit code
from the right side of a pipe. `cmd | tail; echo $?` returns `tail`'s
status, not `cmd`'s. Use `cmd > log 2>&1; echo $?` directly.

## Prerequisites

| Tool | Version | Why                                                           |
| ---- | ------- | ------------------------------------------------------------- |
| node | 18+     | The CLI uses Node ESM; `node --check` validates Detox JS files |
| jq   | any     | Asserts shape of `package.json` and `tsconfig.json`           |
| npx  | any     | One-time fetch of `typescript@5.4` into a local cache for TS parsing |

The first run downloads `typescript@5.4` into
`scripts/test/.cache/node_modules/`. Subsequent runs reuse it offline.

## Usage

From the repo root:

```sh
bash internal/assets/skills/qa-mobile-testing/scripts/test/run-validation.sh
```

The runner:

1. Verifies prereqs (exits 2 if any required tool is missing).
2. Pre-warms the TypeScript cache on first run.
3. Runs every `specs/test-*.sh` in lexicographic order.
4. Exits 0 if every spec passed; 1 otherwise.

Expected wall-clock: ~10s after the TypeScript cache is warm; ~30-60s on
the first run.

## What each spec asserts

| Spec                                  | Exit | What it guards                                                          |
| ------------------------------------- | ---- | ----------------------------------------------------------------------- |
| `specs/test-appium-syntax.sh`         | `0`  | `wdio.conf.ts`, `login.spec.ts`, `login.page.ts` parse; tsconfig + package.json shape; `export const config` present |
| `specs/test-cli-create.sh`            | `0`  | `create <template>` produces a file with NO `{{placeholder}}` left, manifests + hints stripped |
| `specs/test-cli-help.sh`              | `0`  | `help mobile-test-plan` lists Placeholders header + `--project`/`--release`/`--platforms` |
| `specs/test-cli-list.sh`              | `0`  | `list` prints all four templates                                       |
| `specs/test-cli-rejects-traversal.sh` | `1`  | `create ... --out ../escape` rejects path traversal                    |
| `specs/test-cli-self-test.sh`         | `0`  | `--self-test` reports `self-test: OK` and no `FAIL` lines              |
| `specs/test-cli-unknown-template.sh`  | `1`  | `create totally-bogus` exits 1 and lists known templates               |
| `specs/test-detox-syntax.sh`          | `0`  | All `.detoxrc.js` + `e2e/*.js` parse; required keys present; package.json declares detox+jest+scripts |

## Adding a new spec

1. Drop `specs/test-<name>.sh` next to the existing files; copy any of
   them as a starting point.
2. `set -euo pipefail`, source `../lib.sh` via the standard `SCRIPT_DIR`
   idiom, run the script under test with output redirected to
   `/tmp/qa-mobile-validation-<name>.log` and capture `$?` directly into
   `ACTUAL_EXIT`. **Never pipe.**
3. Use `assert_exit_code`, `assert_grep`, `assert_no_grep`,
   `assert_file_exists`. Tail the log on failure for debuggability.
4. The runner picks it up automatically — no registry to update.

## CI snippet (GitHub Actions)

```yaml
jobs:
  qa-mobile-validation:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with: { node-version: "20" }
      - run: sudo apt-get update && sudo apt-get install -y jq
      - run: bash internal/assets/skills/qa-mobile-testing/scripts/test/run-validation.sh
```

The harness is platform-agnostic: it runs on macOS, Linux, and CI without
needing a simulator/emulator.

## What this harness DOES NOT cover (technical debt)

These are deliberate omissions. They require infrastructure that does not
fit the deterministic-CI bar:

- Booting an iOS Simulator and running Detox/Appium against `MyApp.app`
- Booting an Android Emulator and running Espresso/Detox against an APK
- Cloud device-farm runs (BrowserStack, Sauce Labs, AWS Device Farm,
  Firebase Test Lab)

Track candidates as separate beads issues. Once a deterministic emulator
strategy lands (e.g. headless `emulator -no-window` on Linux), extend the
runner with a `--with-emulator` flag rather than rewriting it.
