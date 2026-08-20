#!/usr/bin/env bash
# DE.PULSE v18.8.1 exact-source G12 certification.
# Adaptive Trust, Data Truth & Qualification Closure.
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

evidence_root="${DEPULSE_EVIDENCE_DIR:-$repo_root/.depulse-certification/v18.8.1/$source_sha}"
mkdir -p "$evidence_root"
log_file="$evidence_root/certification.log"
manifest_file="$evidence_root/certification-result.json"
started_at="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
exec > >(tee "$log_file") 2>&1

echo "DE.PULSE v18.8.1 G12 exact-source certification"
echo "Source SHA: $source_sha"
echo "Branch: $source_branch"
echo "Started: $started_at"

echo
echo "[G0/G1] Canonical release identity and v18.8.1 zero-miss closure"
python3 release_identity.py --verify
python3 version_consistency_test.py
python3 - <<'PY'
import json
identity=json.load(open('release_identity.json'))
contract=json.load(open('release/v18.8.1/release_contract.json'))
recon=json.load(open('adaptive-governance/V18.8.1-ZERO-MISS-RECONCILIATION.json'))
assert identity['version']=='18.8.1'
assert identity['previous_stable']=='v18.8.0'
assert identity['stable_baseline']=='v18.8.0'
assert contract['release']=='18.8.1'
assert contract['certifiedStableBaseline']=='v18.8.0-stable'
assert contract['identity_asset']=='release-identity-v18.8.1.js'
assert recon['requireAllClosed'] is True
assert recon['currentPackets'] and all(row['state']=='CLOSED' for row in recon['currentPackets'])
print('PASS: v18.8.1 identity / release contract / zero-miss closure')
PY

echo
echo "[G2/G7/G10] CI, reproducibility, portability, data utility and rights"
python3 tools/ci/workflow_policy.py
python3 tools/ci/reproducibility_gate.py
python3 tools/ci/source_provenance_test.py
python3 dependency_readiness_gate.py
python3 adaptive_resume_gate.py
python3 functionality_utility_checkpoint_gate.py
python3 ai_continuous_eval_gate.py
python3 tools/ci/watchlist_membership_contract.py
python3 v18_5_1_v17_v18_reconciliation_gate.py
python3 tools/ci/v18_8_1_zero_miss_reconciliation_gate.py
python3 tools/ci/v18_8_1_readiness_contract.py
python3 tools/ci/v18_8_1_freshness_contract.py

echo
echo "[G2/G7/G8/G10/G12] v18.8.1 data-truth and universe proof"
go test -count=1 -run 'TestV1881|TestV188' ./...

echo
echo "[G7/G8/G10/G12] Go format, vet, full suite, race and randomized order"
unformatted="$(gofmt -l $(git ls-files '*.go'))"
if [[ -n "$unformatted" ]]; then
  echo "ERROR: gofmt drift:"
  echo "$unformatted"
  exit 1
fi
go vet ./...
go test -count=1 ./...
go test -race -count=1 ./...
go test -shuffle=on -count=1 ./...

echo
echo "[G9/G10/G12] Renderer syntax, deterministic and trust-closure regressions"
while IFS= read -r -d '' file; do node --check "$file"; done < <(find renderer -type f -name '*.js' -print0)
node deterministic_equivalence_test.js
node --require ./tools/ci/release_identity_renderer_preload.js renderer_logic_test.js
node v18_0_5_renderer_test.js
node documentation_ui_owner_test.js
node tests/renderer/v18_8_1_trust_closure_test.js
node tests/renderer/surface_consolidation_test.js
node tests/renderer/documentation_access_test.js

echo
echo "[G9/G10/G12] Current primary Chrome behavior"
python3 -m py_compile \
  release/v18.5.1/browser_live_render_test.py \
  release/v18.6.0/browser_watchlist_membership_test.py \
  release/v18.6.1/browser_watchlist_global_remove_test.py \
  release/v18.5.2/browser_master_symbol_input_test.py \
  release/v18.5.2/browser_profile_display_name_test.py \
  release/v18.5.2/browser_settings_save_bar_test.py \
  tools/ci/documentation_owner_browser_test.py
python3 -c 'import playwright' >/dev/null 2>&1 || {
  echo "ERROR: Python Playwright is missing." >&2
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
python3 release/v18.5.1/browser_live_render_test.py
python3 release/v18.6.0/browser_watchlist_membership_test.py
python3 release/v18.6.1/browser_watchlist_global_remove_test.py
python3 release/v18.5.2/browser_master_symbol_input_test.py
python3 release/v18.5.2/browser_profile_display_name_test.py
python3 release/v18.5.2/browser_settings_save_bar_test.py
python3 tools/ci/documentation_owner_browser_test.py --engine chrome

echo
echo "[G10/G12] Pre-merge qualification dependency"
echo "G11 independently requires exact-head DE.PULSE/fast-head + DE.PULSE/qualified-head success."
echo "Full Qualified is the mandatory Chrome + WebKit + race/randomized proof before this merged candidate can enter G12."

completed_at="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
python3 - "$manifest_file" "$source_sha" "$source_branch" "$started_at" "$completed_at" "$log_file" <<'PY'
import hashlib
import json
import pathlib
import sys
manifest, sha, branch, started, completed, log = sys.argv[1:]
log_path = pathlib.Path(log)
payload = {
    'schema': 'DE.PULSE-G12-CERTIFICATION-1',
    'release': 'v18.8.1',
    'source_sha': sha,
    'source_branch': branch,
    'started_at_utc': started,
    'completed_at_utc': completed,
    'result': 'PASS',
    'lane': 'EXACT_SOURCE_FULL_CERTIFICATION',
    'zero_miss_reconciliation': 'PASS',
    'active_renderer_owner_release_identity': 'PASS',
    'evidence_time_and_freshness_truth': 'PASS',
    'research_symbol_readiness_closure': 'PASS',
    'cost_per_trustworthy_evidence_contract': 'PASS',
    'go_full': 'PASS',
    'race': 'PASS',
    'randomized': 'PASS',
    'deterministic_renderer': 'PASS',
    'chrome': 'PASS',
    'webkit': 'REQUIRED_PREMERGE_EXACT_HEAD_QUALIFIED',
    'conserved_scope_reconciliation': 'PASS',
    'ci_reproducibility': 'PASS',
    'log_sha256': hashlib.sha256(log_path.read_bytes()).hexdigest(),
    'native_packaging_status': 'REQUIRED_BEFORE_G15_PROMOTION',
    'protected_boundaries': [
        'US_EQUITIES_PROCESSING',
        'NO_EXECUTION',
        'DETERMINISTIC_DAY_SWING_LONG',
        'SMART_PROVIDER_ROUTER_V2_SOLE_ROUTING_OWNER',
        'BROAD_SNAPSHOT_BROKER_CANONICAL_REUSE_OWNER'
    ]
}
pathlib.Path(manifest).write_text(json.dumps(payload, indent=2) + '\n', encoding='utf-8')
PY

echo
echo "PASS: exact-source v18.8.1 G12 certification completed."
echo "Evidence: $manifest_file"
echo "Next: G13/G14 native macOS Apple Silicon and Windows x64 package/runtime audit."
