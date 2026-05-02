#!/usr/bin/env node
// check-ts-syntax.mjs — syntax-only TypeScript validation.
//
// Why not `tsc --noEmit`? The example scaffolds reference @wdio/types,
// @wdio/globals, @types/node, @types/mocha — installing them would balloon
// CI time and fail offline. The harness only needs to guarantee the example
// files PARSE as valid TypeScript; full type checking is the user's job once
// they `npm install` in the example.
//
// Strategy: install typescript ONCE into `<harness>/.cache/node_modules` and
// reuse it. We use `transpileModule`, which is parser-only — it does NOT
// resolve modules or check types — so it never fails on missing @wdio/types.
//
// Usage:
//   node check-ts-syntax.mjs <file.ts> [<file2.ts> ...]
import { readFileSync, existsSync, mkdirSync, writeFileSync } from "node:fs";
import { execSync } from "node:child_process";
import { createRequire } from "node:module";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import process from "node:process";

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);
const CACHE_DIR = resolve(__dirname, ".cache");
const TS_VERSION = "5.4.5";

const files = process.argv.slice(2);
if (files.length === 0) {
  console.error("[check-ts-syntax] usage: node check-ts-syntax.mjs <file.ts> [...]");
  process.exit(2);
}

function ensureTypescript() {
  const tsPkgJson = join(CACHE_DIR, "node_modules", "typescript", "package.json");
  if (existsSync(tsPkgJson)) return;
  if (!existsSync(CACHE_DIR)) mkdirSync(CACHE_DIR, { recursive: true });
  // Bootstrap a stub package.json so npm can install into CACHE_DIR.
  const stub = join(CACHE_DIR, "package.json");
  if (!existsSync(stub)) {
    writeFileSync(stub, JSON.stringify({ name: "qa-mobile-ts-cache", private: true }), "utf8");
  }
  console.error(`[check-ts-syntax] installing typescript@${TS_VERSION} into ${CACHE_DIR} (one-time)`);
  execSync(
    `npm install --silent --no-audit --no-fund --no-save typescript@${TS_VERSION}`,
    { cwd: CACHE_DIR, stdio: ["ignore", "ignore", "inherit"] },
  );
}

ensureTypescript();
const tsRoot = join(CACHE_DIR, "node_modules", "typescript", "package.json");
const requireFromCache = createRequire(tsRoot);
const ts = requireFromCache("typescript");

let failed = 0;
for (const file of files) {
  const source = readFileSync(file, "utf8");
  // transpileModule is parser-only — no module resolution, no type checking.
  // Diagnostics here surface ONLY syntax errors, which is exactly what the
  // harness needs to validate.
  const result = ts.transpileModule(source, {
    fileName: file,
    reportDiagnostics: true,
    compilerOptions: {
      target: ts.ScriptTarget.ES2022,
      module: ts.ModuleKind.CommonJS,
      jsx: ts.JsxEmit.Preserve,
      noEmit: true,
      allowJs: true,
    },
  });
  const diags = (result.diagnostics ?? []).filter(
    (d) => d.category === ts.DiagnosticCategory.Error,
  );
  if (diags.length === 0) {
    console.log(`OK   ${file}`);
  } else {
    failed += 1;
    console.error(`FAIL ${file}`);
    for (const d of diags) {
      const msg = ts.flattenDiagnosticMessageText(d.messageText, "\n");
      const pos = d.file && d.start !== undefined
        ? d.file.getLineAndCharacterOfPosition(d.start)
        : null;
      const where = pos ? `:${pos.line + 1}:${pos.character + 1}` : "";
      console.error(`  ${file}${where}: ${msg}`);
    }
  }
}

process.exit(failed === 0 ? 0 : 1);
