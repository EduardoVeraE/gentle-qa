#!/usr/bin/env bash
# xss-scan.sh — XSS scanning via dalfox (primary) with ZAP fallback.
# OWASP: A03 Injection / A05. REQUIRES AUTHORIZATION.
set -euo pipefail

SCRIPT_NAME="xss-scan"
TIMESTAMP="$(date -u +%Y%m%dT%H%M%SZ)"
DEFAULT_OUT="./security-out/${SCRIPT_NAME}/${TIMESTAMP}"

TARGET=""
OUT_DIR=""
SEVERITY_THRESHOLD="high"
CRAWL=0
USE_ZAP=0

usage() {
  cat <<EOF
Usage: $0 --target <url> [options]

XSS scanner. Default tool: dalfox. Fallback: OWASP ZAP active scan.

Options:
  --target <url>             Target URL (single endpoint by default)
  --out <dir>                Output directory (default: ${DEFAULT_OUT})
  --severity-threshold <s>   low|medium|high|critical (default: high)
  --crawl                    Crawl from the target before fuzzing (off by default)
  --zap                      Force ZAP fallback even if dalfox is installed
  -h, --help                 Show help and exit

Examples:
  $0 --target 'https://staging.example.com/search?q=foo'
  $0 --target 'https://staging.example.com/' --crawl
  $0 --target 'https://staging.example.com/' --zap
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    -h|--help) usage; exit 0 ;;
    --target) TARGET="$2"; shift 2 ;;
    --out) OUT_DIR="$2"; shift 2 ;;
    --severity-threshold) SEVERITY_THRESHOLD="$2"; shift 2 ;;
    --crawl) CRAWL=1; shift ;;
    --zap) USE_ZAP=1; shift ;;
    *) echo "Unknown option: $1" >&2; usage; exit 64 ;;
  esac
done

cat <<'BANNER' >&2
[!] AUTHORIZATION REQUIRED
[!] XSS scanning sends active payloads (script/img/svg) to the target.
[!] Run only against systems you own or have written permission to test.
BANNER

if [[ -z "$TARGET" ]]; then
  echo "Error: --target is required." >&2
  usage
  exit 64
fi

OUT_DIR="${OUT_DIR:-$DEFAULT_OUT}"
mkdir -p "$OUT_DIR"
SUMMARY="$OUT_DIR/summary.txt"
JSON_OUT="$OUT_DIR/dalfox.json"

run_dalfox() {
  local mode="url"
  [[ "$CRAWL" -eq 1 ]] && mode="url --deep-domain-xss --skip-bav"
  echo "Running: dalfox $mode $TARGET --format json -o $JSON_OUT" | tee "$SUMMARY"
  set +e
  # shellcheck disable=SC2086
  dalfox $mode "$TARGET" --format json -o "$JSON_OUT" 2>&1 | tee -a "$SUMMARY"
  set -e
}

run_zap() {
  if ! command -v docker >/dev/null 2>&1; then
    cat <<'HINT' >&2
Error: ZAP fallback requires Docker.
Install Docker Desktop: https://docs.docker.com/get-docker/
Or install dalfox: go install github.com/hahwul/dalfox/v2@latest
HINT
    exit 2
  fi
  echo "Running ZAP baseline (passive) — fallback mode" | tee "$SUMMARY"
  docker run --rm -v "$OUT_DIR":/zap/wrk/:rw zaproxy/zap-stable \
    zap-baseline.py -t "$TARGET" -J zap.json -r zap.html 2>&1 | tee -a "$SUMMARY" || true
  JSON_OUT="$OUT_DIR/zap.json"
}

if [[ "$USE_ZAP" -eq 1 ]] || ! command -v dalfox >/dev/null 2>&1; then
  if [[ "$USE_ZAP" -ne 1 ]]; then
    cat <<'HINT' >&2
Warning: dalfox not installed — falling back to ZAP.
Install dalfox: go install github.com/hahwul/dalfox/v2@latest
            or: brew install dalfox
HINT
  fi
  run_zap
else
  run_dalfox
fi

FINDINGS=0
if [[ -f "$JSON_OUT" ]] && grep -Eq '"severity"\s*:\s*"(High|Medium|Low|Critical)"|"type"\s*:\s*"V"' "$JSON_OUT" 2>/dev/null; then
  FINDINGS=$(grep -Eo '"severity"\s*:\s*"[^"]+"' "$JSON_OUT" | wc -l | tr -d ' ')
fi

{
  echo "---"
  echo "Target: $TARGET"
  echo "Crawl: $CRAWL"
  echo "Findings: $FINDINGS"
  echo "Severity threshold: $SEVERITY_THRESHOLD"
  echo "Report: $JSON_OUT"
} | tee -a "$SUMMARY"

case "$SEVERITY_THRESHOLD" in
  low|medium|high|critical) ;;
  *) echo "Invalid --severity-threshold: $SEVERITY_THRESHOLD" >&2; exit 64 ;;
esac

if [[ "$FINDINGS" -gt 0 ]]; then
  echo "XSS findings detected." >&2
  exit 1
fi
echo "No XSS findings."
exit 0
