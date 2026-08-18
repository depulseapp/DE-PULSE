#!/usr/bin/env bash
# DE.PULSE v18.5.2 no-cost exact-source certification.
# Run with: bash release/v18.5.2/run_free_certification.sh
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$repo_root"

for tool in git go gofmt node python3; do
  command -v "$tool" >/dev/null 2>&1 || {
    echo "ERROR: required tool is missing: $tool" >&2
    exit 2
  }
done

source_sha="$(git rev-parse HEAD)"
source_branch="$(git branch --show-current)"
expected_sha="${DEPULSE_EXPECTED_SHA:-}"

if [[ -n "$expected_sha" && "$source_sha" != "$expected_sha" ]]; then
  echo "ERROR: expected $expected_sha but checkout is $source_sha" >&2
  exit 2
fi
if [[ -n "$(git status --porcelain --untracked-files=normal)" ]]; then
  echo "ERROR: certification requires a clean exact-source checkout." >&2
  git status --short >&2
  exit 2
fi

evidence_root="${DEPULSE_EVIDENCE_DIR:-$repo_root/.depulse-certification/manual-v18.5.2/$source_sha}"
mkdir -p "$evidence_root"
log_file="$evidence_root/certification.log"
manifest_file="$evidence_root/certification-result.json"
started_at="$(date -u +%Y-%m-%dT%H:%M:%SZ)"

exec > >(tee "$log_file") 2>&1

echo "DE.PULSE v18.5.2 no-cost certification"
echo "Source SHA: $source_sha"
echo "Branch: $source_branch"
echo "Started: $started_at"

echo
echo "[G0/G1] Canonical identity and version consistency"
python3 release_identity.py --verify
python3 version_consistency_test.py

echo
echo "[G2/G10] Adaptive functionality, provider and Market Mode integration"
python3 functionality_utility_checkpoint_gate.py

echo
echo "[G10/G12] Go formatting, vet, deterministic suite, race and randomized order"
unformatted="$(gofmt -l identity_model.go http_auth.go http_api.go v18_5_2_display_name_test.go)"
if [[ -n "$unformatted" ]]; then
  echo "ERROR: gofmt drift:"
  echo "$unformatted"
  exit 1
fi
go vet ./...
go test -count=1 ./...
go test -race -count=1 ./...
go test -shuffle=on -count=2 ./...

echo
echo "[G12] Renderer syntax and protected deterministic contracts"
node --check renderer/renderer.js
node --check renderer/live-dom-reconcile.js
node --check renderer/watchlist-v18.5.1.js
node --check renderer/header-v18.5.1.js
node deterministic_equivalence_test.js
node renderer_logic_test.js
node v18_0_5_renderer_test.js

browser_tests=(
  release/v18.5.1/browser_live_render_test.py
  release/v18.5.1/browser_watchlist_membership_test.py
  release/v18.5.1/browser_auth_copy_test.py
  release/v18.5.1/browser_ui_hierarchy_test.py
  release/v18.5.2/browser_master_symbol_input_test.py
  release/v18.5.2/browser_profile_display_name_test.py
  release/v18.5.2/browser_settings_save_bar_test.py
)
python3 -m py_compile "${browser_tests[@]}"
python3 -c 'import playwright' >/dev/null 2>&1 || {
  echo "ERROR: Python Playwright is missing." >&2
  echo "Install once with: python3 -m pip install --user playwright" >&2
  exit 2
}

if [[ -z "${CHROME_BIN:-}" ]]; then
  if [[ -x "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" ]]; then
    export CHROME_BIN="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
  elif command -v google-chrome >/dev/null 2>&1; then
    export CHROME_BIN="$(command -v google-chrome)"
  elif command -v chromium >/dev/null 2>&1; then
    export CHROME_BIN="$(command -v chromium)"
  else
    echo "ERROR: set CHROME_BIN to an installed Chrome/Chromium executable." >&2
    exit 2
  fi
fi
echo "Chrome: $CHROME_BIN"

for test_file in "${browser_tests[@]}"; do
  python3 "$test_file"
done

completed_at="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
python3 - "$manifest_file" "$source_sha" "$source_branch" "$started_at" "$completed_at" "$log_file" <<'PY'
import hashlib
import json
import pathlib
import sys

manifest, sha, branch, started, completed, log = sys.argv[1:]
log_path = pathlib.Path(log)
payload = {
    "schema": "DE.PULSE-MANUAL-CERTIFICATION-1",
    "release": "v18.5.2",
    "source_sha": sha,
    "source_branch": branch,
    "started_at_utc": started,
    "completed_at_utc": completed,
    "result": "PASS",
    "lane": "NO_COST_LOCAL_EXACT_SOURCE",
    "log_sha256": hashlib.sha256(log_path.read_bytes()).hexdigest(),
    "native_packaging_status": "REQUIRED_BEFORE_G15_PROMOTION",
}
pathlib.Path(manifest).write_text(json.dumps(payload, indent=2) + "\n", encoding="utf-8")
PY

echo
echo "PASS: exact-source v18.5.2 source and browser certification completed."
echo "Evidence: $manifest_file"
echo "Next: G13/G14 native macOS and Windows packaging/runtime audit before G15 promotion."
