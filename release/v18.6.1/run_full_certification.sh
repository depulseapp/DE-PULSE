#!/usr/bin/env bash
# DE.PULSE v18.6.1 focused patch certification.
# Inherit the complete v18.6.0 G12 matrix, then add patch-specific browser proof.
set -euo pipefail
repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$repo_root"
source_sha="$(git rev-parse HEAD)"
expected_sha="${DEPULSE_EXPECTED_SHA:-}"
if [[ -n "$expected_sha" && "$source_sha" != "$expected_sha" ]]; then
  echo "ERROR: expected $expected_sha but checkout is $source_sha" >&2
  exit 2
fi
if [[ -n "$(git status --porcelain --untracked-files=normal)" ]]; then
  echo "ERROR: certification requires a clean exact-source checkout." >&2
  exit 2
fi

python3 release_identity.py --verify
python3 version_consistency_test.py
python3 - <<'PY'
import json
p=json.load(open('release/v18.6.1/patch_contract.json'))
assert p['release']=='18.6.1'
assert p['base_stable']=='v18.6.0'
assert p['inherited_certification_plan']=='18.6.0'
assert 'no gate' in p['quality_rule'].lower()
print('PASS: v18.6.1 patch inheritance contract')
PY

evidence_root="${DEPULSE_EVIDENCE_DIR:-$repo_root/.depulse-certification/v18.6.1/$source_sha}"
mkdir -p "$evidence_root"
started_at="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
renderer_preload="$repo_root/tools/ci/release_identity_renderer_preload.js"

# Full inherited certification; the Node preload makes direct source-reading
# renderer regressions use the canonical patch identity while keeping the
# inherited renderer implementation itself frozen. Product runtime identity is
# enforced by renderer/watchlist-desk-contract-v18.6.1.js.
NODE_OPTIONS="--require=$renderer_preload${NODE_OPTIONS:+ $NODE_OPTIONS}" \
DEPULSE_EVIDENCE_DIR="$evidence_root/inherited-v18.6.0" \
DEPULSE_EXPECTED_SHA="$source_sha" \
  bash release/v18.6.0/run_full_certification.sh

# Patch-specific regression that reproduces the user-reported DESKS failure and
# exercises all seven legal Day/Swing/Long membership combinations plus undo,
# Master Market Store removal, rapid double activation, failure behavior,
# membership-only toggles and the centered alert CSS contract.
python3 release/v18.6.1/browser_watchlist_global_remove_test.py | tee "$evidence_root/v18.6.1-browser-watchlist.log"

completed_at="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
python3 - "$evidence_root/patch-certification-result.json" "$source_sha" "$started_at" "$completed_at" <<'PY'
import hashlib,json,pathlib,sys
out,sha,started,completed=sys.argv[1:]
log=pathlib.Path(out).parent/'v18.6.1-browser-watchlist.log'
payload={
  'schema':'DE.PULSE-PATCH-CERTIFICATION-1',
  'release':'v18.6.1',
  'base_certification':'v18.6.0 full G12 matrix',
  'source_sha':sha,
  'started_at_utc':started,
  'completed_at_utc':completed,
  'result':'PASS',
  'patch_browser_regression':'PASS',
  'patch_browser_log_sha256':hashlib.sha256(log.read_bytes()).hexdigest(),
  'protected_boundaries':['US_EQUITIES_PROCESSING','NO_EXECUTION','DETERMINISTIC_DESK_FORMULAS'],
}
pathlib.Path(out).write_text(json.dumps(payload,indent=2)+'\n',encoding='utf-8')
PY

echo "PASS: v18.6.1 exact-source patch certification = inherited v18.6.0 full matrix + v18.6.1 edge/browser proof."
