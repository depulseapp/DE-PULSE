#!/usr/bin/env python3
"""Regression contract for canonical Day/Swing/Long membership behavior."""

from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
JS = (ROOT / "renderer" / "watchlist-v18.5.1.js").read_text(encoding="utf-8")
CSS = (ROOT / "renderer" / "watchlist-v18.5.1.css").read_text(encoding="utf-8")
QUALIFIED = (ROOT / ".github" / "workflows" / "ci-qualified.yml").read_text(encoding="utf-8")
BROWSER_PROOF = ROOT / "release" / "v18.6.0" / "browser_watchlist_membership_test.py"

required_js = (
    "deskMembershipStripV186",
    "aria-pressed",
    "addSymbolToDesk",
    "'/api/desk/membership'",
    "[data-add-desk]",
    "[data-add-desk-table]",
    "[data-ai-add-desk]",
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

qualified_proof = "python3 release/v18.6.0/browser_watchlist_membership_test.py"
if qualified_proof not in QUALIFIED:
    raise SystemExit("watchlist contract FAIL: qualified browser lane is not bound to the v18.6 membership proof")

if not BROWSER_PROOF.is_file():
    raise SystemExit("watchlist contract FAIL: v18.6 browser membership proof is missing")

print("watchlist membership contract PASS")
