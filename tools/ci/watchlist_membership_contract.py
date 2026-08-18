#!/usr/bin/env python3
"""Regression contract for canonical Day/Swing/Long membership behavior."""

from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
JS = (ROOT / "renderer" / "watchlist-v18.5.1.js").read_text(encoding="utf-8")
CSS = (ROOT / "renderer" / "watchlist-v18.5.1.css").read_text(encoding="utf-8")

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

print("watchlist membership contract PASS")
