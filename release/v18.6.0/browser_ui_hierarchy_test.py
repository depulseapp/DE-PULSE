#!/usr/bin/env python3
"""Chromium proof for the v18.6 header and Research hierarchy contract.

Uses the actual accumulated styles.css, retained ui-v18.5.1.css/header-v18.5.1.js
implementation layers, and the canonical v18.6 release cache-buster from
release_identity.json. No product network calls are required.
"""
import json
import os
from pathlib import Path

from playwright.sync_api import sync_playwright

ROOT = Path(__file__).resolve().parents[2]
BASE_CSS = (ROOT / "renderer" / "styles.css").read_text(encoding="utf-8")
UI_CSS = (ROOT / "renderer" / "ui-v18.5.1.css").read_text(encoding="utf-8")
HEADER_JS = (ROOT / "renderer" / "header-v18.5.1.js").read_text(encoding="utf-8")
INDEX = (ROOT / "renderer" / "index.html").read_text(encoding="utf-8")
RELEASE_VERSION = json.loads((ROOT / "release_identity.json").read_text(encoding="utf-8"))["version"]

HEADER_HTML = r"""
<div id="app-shell">
  <header class="topbar">
    <div class="brand"><img alt="DE.PULSE"><div class="brand-copy"><strong class="brand-name"><span>DE</span><i class="brand-dot"></i><span>PULSE</span><svg class="brand-pulse"></svg></strong><small>Personal Market Intelligence</small></div></div>
    <div id="header-notification" class="header-notification" aria-hidden="true"></div>
    <div class="runtime">
      <div id="market-session-context" class="market-session-context open" aria-label="US market session"><b id="market-session-label">MARKET OPEN</b><small id="market-session-detail">Closes in 3h 15m · Mon 4 PM ET · 1 PM PT</small></div>
      <div class="market-clocks" aria-label="Eastern and Pacific market times">
        <span class="market-clock" title="US Eastern Time"><b id="zone-et">ET</b><span class="market-clock-value"><small id="date-et">AUG 17, 2026</small><time id="clock-et">6:07:42 PM</time></span></span>
        <span class="market-clock" title="US Pacific Time"><b id="zone-pt">PT</b><span class="market-clock-value"><small id="date-pt">AUG 17, 2026</small><time id="clock-pt">3:07:42 PM</time></span></span>
      </div>
      <span id="runtime-status" class="runtime-pill running" aria-label="Market data status">DATA HEALTHY</span>
      <button id="runtime-toggle" class="btn danger">STOP DATA</button>
      <span id="identity-principal" class="identity-principal">OWNER</span>
      <button id="identity-signout" class="btn ghost identity-signout">SIGN OUT</button>
    </div>
  </header>
  <div id="macro-event-banner" class="macro-event-banner" aria-hidden="true"></div>
  <div id="ticker-tape" class="ticker-tape" aria-label="Market Instruments"><strong class="ticker-tape-label">Market Instruments</strong><div class="ticker-viewport"><span class="ticker-item"><b>SPY</b><span>$500.00</span></span></div></div>
  <div class="workspace"><aside class="sidebar"></aside><main class="main">
    <div class="page research-page research-v2">
      <section class="card research-command research-command-v2">
        <div class="research-command-heading">
          <div><span class="eyebrow">Research Target</span><p>Select one symbol, confirm the complete Research evidence state, then review sourced analysis below.</p></div>
          <button id="research-back" class="btn ghost" data-research-back>Back to Dashboard</button>
        </div>
        <div class="research-command-symbol">
          <div class="research-target-grid">
            <label id="research-primary"><small>Choose from List</small><select class="ticker-control" data-research-symbol><option>NVDA</option></select></label>
            <label id="research-add"><small>Add Symbol</small><span class="research-add-symbol"><input class="ticker-input" placeholder="Ticker"><button class="btn">Load</button></span></label>
            <div id="research-freshness" class="research-target-status"><small>Research freshness</small><span class="data-badge live">CURRENT</span></div>
          </div>
          <p class="research-origin-context">Opened from Dashboard · return context preserved.</p>
        </div>
      </section>
      <nav class="research-tabs research-tabs-v2" aria-label="Research views">
        <button>Overview</button><button>Earnings</button><button>Fundamentals</button><button>SEC &amp; Ownership</button><button>Catalysts</button><button>Technical Context</button>
      </nav>
    </div>
  </main></div>
</div>
"""

STUB_JS = r"""
let __baseUpdateCount=0;
function updateChrome(){__baseUpdateCount+=1}
function headerDataHealth(){return {label:'DATA HEALTHY',tone:'running',detail:'Alpaca primary snapshots current · full actionable coverage'} }
"""


def visible_box(page, selector):
    loc = page.locator(selector)
    assert loc.count() == 1, selector
    assert loc.is_visible(), selector
    box = loc.bounding_box()
    assert box and box["width"] > 0 and box["height"] > 0, (selector, box)
    return box


def bottom(box):
    return box["y"] + box["height"]


def main() -> None:
    assert RELEASE_VERSION.startswith("18.6."), RELEASE_VERSION
    assert f"ui-v18.5.1.css?v={RELEASE_VERSION}" in INDEX
    assert f"header-v18.5.1.js?v={RELEASE_VERSION}" in INDEX
    assert "grid-template-columns:minmax(260px,1.3fr)" in UI_CSS
    assert "align-items:end!important" in UI_CSS
    assert "top:calc(var(--v1851-topbar-h) + var(--v1851-statusbar-h))" in UI_CSS
    assert "ensureSecondaryMarketStatus" in HEADER_JS
    assert "content.appendChild(clocks)" in HEADER_JS
    assert "research-command-heading" in UI_CSS

    with sync_playwright() as p:
        kwargs = {"headless": True}
        chrome = os.environ.get("CHROME_BIN", "").strip()
        if chrome:
            assert Path(chrome).is_file(), f"CHROME_BIN missing: {chrome}"
            kwargs["executable_path"] = chrome
        browser = p.chromium.launch(**kwargs)
        page = browser.new_page(viewport={"width": 1440, "height": 900})
        page.set_content(f"<style>{BASE_CSS}\n{UI_CSS}</style>{HEADER_HTML}")
        page.add_script_tag(content=STUB_JS)
        page.add_script_tag(content=HEADER_JS)

        for width in (1440, 900, 560):
            page.set_viewport_size({"width": width, "height": 900})
            page.wait_for_timeout(30)

            overflow = page.evaluate("document.documentElement.scrollWidth - window.innerWidth")
            assert overflow <= 1, (width, "page horizontal overflow", overflow)

            top = visible_box(page, ".topbar")
            status = visible_box(page, "#market-status-bar")
            ticker = visible_box(page, "#ticker-tape")
            assert status["y"] >= top["y"] + top["height"] - 1, (width, top, status)
            assert ticker["y"] >= status["y"] + status["height"] - 1, (width, status, ticker)

            assert page.locator(".topbar .runtime #market-session-context").count() == 0, width
            assert page.locator(".topbar .runtime .market-clocks").count() == 0, width
            assert page.locator(".topbar .runtime #runtime-status").count() == 0, width
            assert page.locator(".topbar .runtime #runtime-toggle").count() == 0, width
            assert page.locator("#market-status-bar #market-session-context").count() == 1, width
            assert page.locator("#market-status-bar .market-status-content > .market-clocks").count() == 1, width
            assert page.locator("#market-status-bar #runtime-status").count() == 1, width
            assert page.locator("#market-status-bar #runtime-toggle").count() == 1, width

            ribbon_order = page.locator("#market-status-bar .market-status-content > *").evaluate_all(
                "(nodes)=>nodes.map(x=>x.id||x.className)"
            )
            assert ribbon_order[:4] == [
                "market-session-context",
                "market-data-summary",
                "market-clocks",
                "runtime-toggle",
            ], (width, ribbon_order)
            if width == 1440:
                ribbon_content = visible_box(page, ".market-status-content")
                assert ribbon_content["width"] < status["width"] - 80, (status, ribbon_content)

            for zone in ("et", "pt"):
                visible_box(page, f"#zone-{zone}")
                visible_box(page, f"#date-{zone}")
                visible_box(page, f"#clock-{zone}")
                clock_box = visible_box(page, f".market-clock:has(#clock-{zone})")
                assert clock_box["width"] >= 84, (width, zone, clock_box)
                assert page.locator(f"#date-{zone}").inner_text() == "AUG 17, 2026"
                assert ":" in page.locator(f"#clock-{zone}").inner_text()

            visible_box(page, "#market-session-label")
            visible_box(page, "#runtime-status")
            visible_box(page, "#runtime-toggle")
            visible_box(page, "#market-data-detail")

            visible_box(page, ".research-command-heading")
            visible_box(page, "#research-back")
            assert page.locator(".research-command-v2 .research-command-actions").count() == 0, width
            assert page.locator(".research-command-v2 [data-research-ai]").count() == 0, width
            assert page.evaluate(
                "document.querySelector('.research-command-v2').nextElementSibling.matches('.research-tabs-v2')"
            ), width
            assert page.evaluate(
                "document.querySelector('.research-command-heading').nextElementSibling.matches('.research-command-symbol')"
            ), width
            primary = visible_box(page, "#research-primary")
            add = visible_box(page, "#research-add")
            fresh = visible_box(page, "#research-freshness")
            tabs = page.locator(".research-tabs-v2 button")
            assert tabs.count() == 6
            for i in range(6):
                assert tabs.nth(i).is_visible(), (width, i)
                box = tabs.nth(i).bounding_box()
                assert box and box["x"] >= -1 and box["x"] + box["width"] <= width + 1, (width, i, box)

            if width == 1440:
                assert primary["width"] > add["width"] > fresh["width"] * 0.95, (primary, add, fresh)
                assert abs(bottom(primary) - bottom(add)) <= 2, (primary, add)
                assert abs(bottom(add) - bottom(fresh)) <= 2, (add, fresh)
            elif width == 900:
                assert primary["width"] > add["width"], (primary, add)
                assert primary["y"] < add["y"], (primary, add)
                assert abs(add["y"] - fresh["y"]) <= 2, (add, fresh)
            else:
                assert primary["y"] < add["y"] < fresh["y"], (primary, add, fresh)
                assert abs(primary["width"] - add["width"]) <= 3, (primary, add)

        page.evaluate("updateChrome('SPY')")
        assert page.locator("#market-data-detail").inner_text() == "Alpaca primary snapshots current · full actionable coverage"
        assert page.evaluate("__baseUpdateCount") == 1

        browser.close()

    print("PASS: v18.6 header hierarchy and Research Target responsive hierarchy are behaviorally contained with complete ET/PT clocks and canonical release cache-busters.")


if __name__ == "__main__":
    main()
