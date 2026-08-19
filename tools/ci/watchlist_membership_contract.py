#!/usr/bin/env python3
"""Regression contract for canonical Day/Swing/Long membership behavior."""
from pathlib import Path
import json

ROOT = Path(__file__).resolve().parents[2]
JS = (ROOT / "renderer" / "watchlist-v18.5.1.js").read_text(encoding="utf-8")
CSS = (ROOT / "renderer" / "watchlist-v18.5.1.css").read_text(encoding="utf-8")
INDEX = (ROOT / "renderer" / "index.html").read_text(encoding="utf-8")
IDENTITY = json.loads((ROOT / "release_identity.json").read_text(encoding="utf-8"))
VERSION = IDENTITY["version"]
CONTRACT_PATH = ROOT / "renderer" / f"watchlist-desk-contract-v{VERSION}.js"
PATCH_PROOF = ROOT / "release" / f"v{VERSION}" / "browser_watchlist_global_remove_test.py"
PATCH_CERT = ROOT / "release" / f"v{VERSION}" / "run_full_certification.sh"
BASE_PROOF = ROOT / "release" / "v18.6.0" / "browser_watchlist_membership_test.py"

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

# v18.6.1 defect root-cause guard: the compatibility extension references DESKS
# as a classic-script global. Production must define it explicitly before the
# extension, otherwise Safari/WebKit reports `Can't find variable: DESKS`.
if not CONTRACT_PATH.is_file():
    raise SystemExit(f"watchlist contract FAIL: missing runtime desk contract {CONTRACT_PATH}")
contract = CONTRACT_PATH.read_text(encoding="utf-8")
required_contract = "var DESKS = Object.freeze(['day', 'swing', 'long'])"
if required_contract not in contract:
    raise SystemExit("watchlist contract FAIL: canonical DESKS runtime binding missing")
contract_asset = f"watchlist-desk-contract-v{VERSION}.js?v={VERSION}"
extension_asset = f"watchlist-v18.5.1.js?v={VERSION}"
if contract_asset not in INDEX or extension_asset not in INDEX or INDEX.index(contract_asset) > INDEX.index(extension_asset):
    raise SystemExit("watchlist contract FAIL: DESKS runtime contract must load before watchlist extension")
if not PATCH_PROOF.is_file():
    raise SystemExit("watchlist contract FAIL: patch-specific global-remove browser proof missing")
if not PATCH_CERT.is_file() or str(PATCH_PROOF.relative_to(ROOT)) not in PATCH_CERT.read_text(encoding="utf-8"):
    raise SystemExit("watchlist contract FAIL: G12 patch certification is not bound to global-remove browser proof")

print(f"watchlist membership contract PASS · v{VERSION} DESKS binding, toggle semantics and G12 edge proof wired")
