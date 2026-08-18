#!/usr/bin/env python3
"""Browser proof that the Settings save bar stays fully visible at short desktop heights."""

import os
from pathlib import Path

from playwright.sync_api import sync_playwright

ROOT = Path(__file__).resolve().parents[2]
BASE_CSS = (ROOT / "renderer" / "styles.css").read_text(encoding="utf-8")
UI_CSS = (ROOT / "renderer" / "ui-v18.5.1.css").read_text(encoding="utf-8")

FIXTURE = """
<header class="topbar"><strong>DE.PULSE</strong></header>
<section class="market-status-bar"><div class="market-status-content">Market Pulse</div></section>
<div class="ticker-tape">Market Instruments</div>
<div class="workspace">
  <aside class="sidebar"></aside>
  <main id="main" class="main" data-page="settings">
    <div class="settings-shell">
      <section class="settings" data-scope-marker="settings-v1421">
        <article class="settings-section" style="height:360px"><h2>Account</h2></article>
        <article class="settings-section" style="height:360px"><h2>Data</h2></article>
        <article class="settings-section" style="height:360px"><h2>Application</h2></article>
      </section>
      <div class="save-bar settings-persistent-save">
        <span class="save-security-note">API keys stay local and are hidden after saving.</span>
        <div class="save-actions">
          <span class="last-saved-inline">Last saved: just now</span>
          <button class="btn primary" data-save-settings>Save Settings</button>
        </div>
      </div>
    </div>
  </main>
</div>
"""


def assert_visible_contract(page, height: int) -> None:
    page.set_viewport_size({"width": 1280, "height": height})
    page.set_content(FIXTURE)
    page.add_style_tag(content=BASE_CSS)
    page.add_style_tag(content=UI_CSS)
    page.wait_for_timeout(50)

    before = page.evaluate(
        """() => {
          const bar=document.querySelector('.settings-persistent-save');
          const button=document.querySelector('[data-save-settings]');
          const pane=document.querySelector('.settings-shell>.settings');
          const main=document.querySelector('main');
          const br=bar.getBoundingClientRect();
          const rr=button.getBoundingClientRect();
          return {
            barTop:br.top,barBottom:br.bottom,buttonTop:rr.top,buttonBottom:rr.bottom,
            viewport:innerHeight,documentHeight:document.documentElement.scrollHeight,
            paneClient:pane.clientHeight,paneScroll:pane.scrollHeight,
            mainHeight:main.getBoundingClientRect().height,windowY:scrollY
          };
        }"""
    )
    assert before["barTop"] >= 0, (height, before)
    assert before["barBottom"] <= before["viewport"] - 8, (height, before)
    assert before["buttonTop"] >= 0 and before["buttonBottom"] <= before["viewport"] - 8, (height, before)
    assert before["paneScroll"] > before["paneClient"], (height, before)
    assert before["documentHeight"] <= before["viewport"] + 2, (height, before)
    assert before["windowY"] == 0, (height, before)

    page.locator(".settings-shell>.settings").evaluate("(node)=>{node.scrollTop=node.scrollHeight}")
    after = page.evaluate(
        """() => {
          const bar=document.querySelector('.settings-persistent-save').getBoundingClientRect();
          const button=document.querySelector('[data-save-settings]').getBoundingClientRect();
          const pane=document.querySelector('.settings-shell>.settings');
          return {
            barBottom:bar.bottom,buttonBottom:button.bottom,viewport:innerHeight,
            paneAtBottom:Math.ceil(pane.scrollTop+pane.clientHeight)>=pane.scrollHeight,
            windowY:scrollY
          };
        }"""
    )
    assert after["paneAtBottom"], (height, after)
    assert after["barBottom"] <= after["viewport"] - 8, (height, after)
    assert after["buttonBottom"] <= after["viewport"] - 8, (height, after)
    assert after["windowY"] == 0, (height, after)


def main() -> None:
    assert "v18.5.2 Settings viewport contract" in UI_CSS
    assert "height:calc(100dvh - var(--v1851-chrome-h))!important" in UI_CSS
    with sync_playwright() as p:
        kwargs = {"headless": True}
        chrome = os.environ.get("CHROME_BIN", "").strip()
        if chrome:
            assert Path(chrome).is_file(), f"CHROME_BIN missing: {chrome}"
            kwargs["executable_path"] = chrome
        browser = p.chromium.launch(**kwargs)
        page = browser.new_page()
        for height in (900, 520, 402):
            assert_visible_contract(page, height)
        browser.close()
    print("PASS: Settings content scrolls independently and the complete Save Settings bar remains visible at desktop window heights.")


if __name__ == "__main__":
    main()
