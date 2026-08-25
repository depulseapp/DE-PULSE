#!/usr/bin/env python3
"""Regression contract for canonical Day/Swing/Long membership behavior."""
from pathlib import Path
import hashlib
import json
import re

ROOT = Path(__file__).resolve().parents[2]
JS_PATH = ROOT / "renderer" / "watchlist-ui.js"
JS = JS_PATH.read_text(encoding="utf-8")
CSS = (ROOT / "renderer" / "watchlist-desk.css").read_text(encoding="utf-8")
INDEX = (ROOT / "renderer" / "index.html").read_text(encoding="utf-8")
IDENTITY = json.loads((ROOT / "release_identity.json").read_text(encoding="utf-8"))
VERSION = IDENTITY["version"]
CURRENT_MANIFEST = ROOT / "release" / f"v{VERSION}" / "certification-manifest.json"
CANONICAL_EXECUTOR = ROOT / "tools" / "release" / "run_full_certification.py"
BASE_PROOF = ROOT / "release" / "v18.6.0" / "browser_watchlist_membership_test.py"


def git_blob_token(path: Path) -> str:
    data = path.read_bytes()
    return hashlib.sha1(f"blob {len(data)}\0".encode() + data).hexdigest()[:16]


required_js = (
    "deskMembershipStripV186",
    "aria-pressed",
    "addSymbolToDesk",
    "'/api/desk/membership'",
    "[data-add-desk]",
    "[data-add-desk-table]",
    "[data-ai-add-desk]",
    "'/api/master-symbol/remove'",
    "bindGlobalTrackedSymbolRemoval",
)
for token in required_js:
    if token not in JS:
        raise SystemExit(f"watchlist contract FAIL: missing {token!r}")
for forbidden in ("CURRENT", "current-desk", "/api/watchlists/add-symbol"):
    if forbidden in JS:
        raise SystemExit(f"watchlist contract FAIL: override contains deprecated {forbidden!r}")
if '[aria-pressed="true"]' not in CSS:
    raise SystemExit("watchlist contract FAIL: selected membership CSS must follow aria-pressed=true")
if ".current-desk" in CSS:
    raise SystemExit("watchlist contract FAIL: deprecated current-desk styling remains")
if not BASE_PROOF.is_file():
    raise SystemExit("watchlist contract FAIL: inherited v18.6 browser membership proof is missing")

# Runtime behavior ownership is determined by the stable capability asset actually
# loaded in production. The desk contract retains its compatibility patch version
# inside the file and in its cache key without embedding that version in the path.
contract_matches = re.findall(r'(watchlist-desk-contract\.js\?v=([0-9.]+))', INDEX)
if len(contract_matches) != 1:
    raise SystemExit(f"watchlist contract FAIL: expected exactly one runtime desk contract asset, found {len(contract_matches)}")
contract_asset, contract_cache_version = contract_matches[0]
contract_filename = contract_asset.split('?', 1)[0]
contract_path = ROOT / "renderer" / contract_filename
if not contract_path.is_file():
    raise SystemExit(f"watchlist contract FAIL: missing runtime desk contract {contract_path}")
contract = contract_path.read_text(encoding="utf-8")
required_contract = "var DESKS = Object.freeze(['day', 'swing', 'long'])"
if required_contract not in contract:
    raise SystemExit("watchlist contract FAIL: canonical DESKS runtime binding missing")
patch_version_match = re.search(r"DEPULSE_PATCH_VERSION = '([^']+)'", contract)
if not patch_version_match or patch_version_match.group(1) != contract_cache_version:
    raise SystemExit("watchlist contract FAIL: runtime desk contract patch/cache version mismatch")
extension_asset = f"watchlist-ui.js?v={git_blob_token(JS_PATH)}"
if contract_asset not in INDEX or extension_asset not in INDEX or INDEX.index(contract_asset) > INDEX.index(extension_asset):
    raise SystemExit("watchlist contract FAIL: DESKS runtime contract must load before watchlist extension")

# Current releases use the version-neutral G12 executor with a declarative
# per-release manifest. The global-remove proof may be inherited from the release
# where the behavior was introduced; do not reintroduce retired per-release shell
# orchestrators or cosmetic copies merely to satisfy this contract.
if not CANONICAL_EXECUTOR.is_file():
    raise SystemExit("watchlist contract FAIL: canonical version-neutral G12 executor is missing")
if not CURRENT_MANIFEST.is_file():
    raise SystemExit("watchlist contract FAIL: current G12 certification manifest is missing")
manifest = json.loads(CURRENT_MANIFEST.read_text(encoding="utf-8"))
if manifest.get("schema") != "DE.PULSE-G12-EVIDENCE-MANIFEST-1" or manifest.get("productVersion") != VERSION:
    raise SystemExit("watchlist contract FAIL: current G12 manifest identity/schema mismatch")
chrome_tests = manifest.get("chromeTests", [])
proof_matches = sorted({
    str(command[1])
    for command in chrome_tests
    if isinstance(command, list)
    and len(command) >= 2
    and str(command[1]).startswith("release/")
    and re.fullmatch(r"release/v[0-9.]+/browser_watchlist_global_remove_test\.py", str(command[1]))
})
if len(proof_matches) != 1:
    raise SystemExit(f"watchlist contract FAIL: expected exactly one G12 global-remove proof binding, found {proof_matches}")
proof_path = ROOT / proof_matches[0]
if not proof_path.is_file():
    raise SystemExit(f"watchlist contract FAIL: G12 global-remove browser proof is missing: {proof_matches[0]}")
executor = CANONICAL_EXECUTOR.read_text(encoding="utf-8")
for marker in ("DE.PULSE-G12-EVIDENCE-MANIFEST-1", 'manifest.get("chromeTests"', "run(log, list(command), env=browser_env)"):
    if marker not in executor:
        raise SystemExit(f"watchlist contract FAIL: canonical G12 executor lost manifest Chrome execution marker {marker!r}")

print(
    f"watchlist membership contract PASS · app v{VERSION} uses {contract_filename} "
    f"with toggle semantics, content-derived extension cache identity, and manifest-bound G12 edge proof {proof_matches[0]}"
)
