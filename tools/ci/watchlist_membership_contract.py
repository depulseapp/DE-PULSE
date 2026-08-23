#!/usr/bin/env python3
"""Regression contract for canonical Day/Swing/Long membership behavior."""
from pathlib import Path
import hashlib
import json
import re

ROOT = Path(__file__).resolve().parents[2]
JS_PATH = ROOT / "renderer" / "watchlist-v18.5.1.js"
JS = JS_PATH.read_text(encoding="utf-8")
CSS = (ROOT / "renderer" / "watchlist-v18.5.1.css").read_text(encoding="utf-8")
INDEX = (ROOT / "renderer" / "index.html").read_text(encoding="utf-8")
IDENTITY = json.loads((ROOT / "release_identity.json").read_text(encoding="utf-8"))
VERSION = IDENTITY["version"]
CURRENT_CERT = ROOT / "release" / f"v{VERSION}" / "run_full_certification.sh"
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

# Runtime behavior ownership is determined by the asset actually loaded in
# production, not by assuming every later app release must duplicate the same
# compatibility contract under a new versioned filename. The desk contract
# itself still carries its compatibility-version cache key, while the shared
# watchlist extension follows the canonical content-derived cache identity.
contract_matches = re.findall(r'(watchlist-desk-contract-v([0-9.]+)\.js\?v=([0-9.]+))', INDEX)
if len(contract_matches) != 1:
    raise SystemExit(f"watchlist contract FAIL: expected exactly one runtime desk contract asset, found {len(contract_matches)}")
contract_asset, contract_version, contract_cache_version = contract_matches[0]
if contract_version != contract_cache_version:
    raise SystemExit("watchlist contract FAIL: runtime desk contract filename/cache version mismatch")
contract_filename = contract_asset.split('?', 1)[0]
contract_path = ROOT / "renderer" / contract_filename
if not contract_path.is_file():
    raise SystemExit(f"watchlist contract FAIL: missing runtime desk contract {contract_path}")
contract = contract_path.read_text(encoding="utf-8")
required_contract = "var DESKS = Object.freeze(['day', 'swing', 'long'])"
if required_contract not in contract:
    raise SystemExit("watchlist contract FAIL: canonical DESKS runtime binding missing")
extension_asset = f"watchlist-v18.5.1.js?v={git_blob_token(JS_PATH)}"
if contract_asset not in INDEX or extension_asset not in INDEX or INDEX.index(contract_asset) > INDEX.index(extension_asset):
    raise SystemExit("watchlist contract FAIL: DESKS runtime contract must load before watchlist extension")

# The current release certification must consume a real global-remove browser
# proof, but that proof may be inherited from the release where the behavior was
# introduced. Do not force cosmetic copies into every later release directory.
if not CURRENT_CERT.is_file():
    raise SystemExit("watchlist contract FAIL: current G12 certification script is missing")
cert_text = CURRENT_CERT.read_text(encoding="utf-8")
proof_matches = sorted(set(re.findall(r'release/v[0-9.]+/browser_watchlist_global_remove_test\.py', cert_text)))
if len(proof_matches) != 1:
    raise SystemExit(f"watchlist contract FAIL: expected exactly one G12 global-remove proof binding, found {proof_matches}")
proof_path = ROOT / proof_matches[0]
if not proof_path.is_file():
    raise SystemExit(f"watchlist contract FAIL: G12 global-remove browser proof is missing: {proof_matches[0]}")

print(
    f"watchlist membership contract PASS · app v{VERSION} uses {contract_filename} "
    f"with toggle semantics, content-derived extension cache identity, and G12 edge proof {proof_matches[0]}"
)
