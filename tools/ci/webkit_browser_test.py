#!/usr/bin/env python3
"""Primary WebKit browser-only compatibility evidence for DE.PULSE.

This owner intentionally does not package or launch the native macOS app.
Native JXA/Cocoa/WKWebView rehearsal is selected independently by Planner v3
and reuses tools/release/native_macos.sh.
"""
from __future__ import annotations

from playwright.sync_api import sync_playwright

from webkit_targeted_test import EXT, INDEX, PATCH_CSS, settings_layout_contract, watchlist_contract


def main() -> None:
    assert "CURRENT" not in EXT
    assert "aria-pressed" in EXT
    assert 'id="header-notification" class="header-notification"' in INDEX
    assert "justify-self:stretch!important" in PATCH_CSS
    assert "text-align:center!important" in PATCH_CSS

    with sync_playwright() as p:
        browser = p.webkit.launch(headless=True)
        page = browser.new_page(viewport={"width": 1440, "height": 800})
        watchlist_contract(page)
        settings_layout_contract(page)
        browser.close()

    print(
        "PASS: primary WebKit browser compatibility for watchlist/global-remove, "
        "no-CURRENT membership semantics, Settings short-height save bar, and centered header alert."
    )


if __name__ == "__main__":
    main()
