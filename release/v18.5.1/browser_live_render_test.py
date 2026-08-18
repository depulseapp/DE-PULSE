#!/usr/bin/env python3
"""Behavior-first browser proof for v18.5.2 incremental live rendering.

This is intentionally a real Chromium test. It verifies the user-visible contract
that quote ticks update data without replacing/reordering the row under active
interaction, while non-quote structural events retain the full-render path.
"""
import os
from pathlib import Path
from playwright.sync_api import sync_playwright

ROOT = Path(__file__).resolve().parents[2]
LIVE_RECONCILER = ROOT / "renderer" / "live-dom-reconcile.js"

INITIAL_HTML = """
<style>
  #main { height: 160px; overflow: auto; width: 700px; }
  .spacer { height: 120px; }
  article { min-height: 70px; border: 1px solid transparent; }
</style>
<main id="main">
  <div class="page">
    <div class="spacer">top</div>
    <section data-live-key="master-market-symbols">
      <div data-live-key="master-symbols-heading">
        <div>Tracked Symbols</div>
        <div data-live-key="master-symbol-controls">
          <input id="master-symbol-input" data-master-add-input data-live-key="master-symbol-input" value="" aria-label="Add tracked ticker">
          <button data-master-add>Add Symbol</button><button data-master-remove-all>Remove All</button>
        </div>
      </div>
      <div data-live-key="master-symbol-list"><span data-master-symbol="AAPL">AAPL</span></div>
    </section>
    <section data-live-key="quote-list">
      <article data-live-key="row:AAPL" class="quote neutral">
        <button data-symbol="AAPL">AAPL</button>
        <span data-live-key="price:AAPL">$100.00</span>
        <input id="note" value="keep-me" aria-label="interaction sentinel">
      </article>
      <article data-live-key="row:MSFT" class="quote neutral">
        <button data-symbol="MSFT">MSFT</button>
        <span data-live-key="price:MSFT">$200.00</span>
      </article>
    </section>
    <div class="spacer">bottom</div>
  </div>
</main>
"""

NEXT_HTML = """
<div class="page">
  <div class="spacer">top</div>
  <section data-live-key="master-market-symbols">
    <div data-live-key="master-symbols-heading">
      <div>Tracked Symbols</div>
      <div data-live-key="master-symbol-controls">
        <input id="master-symbol-input" data-master-add-input data-live-key="master-symbol-input" value="server-draft" aria-label="Add tracked ticker">
        <button data-master-add>Add Symbol</button><button data-master-remove-all>Remove All</button>
      </div>
    </div>
    <div data-live-key="master-symbol-list"><span data-master-symbol="AAPL">AAPL</span></div>
  </section>
  <section data-live-key="quote-list">
    <article data-live-key="row:MSFT" class="quote negative">
      <button data-symbol="MSFT">MSFT</button>
      <span data-live-key="price:MSFT">$199.50</span>
    </article>
    <article data-live-key="row:AAPL" class="quote positive">
      <button data-symbol="AAPL">AAPL</button>
      <span data-live-key="price:AAPL">$101.25</span>
      <input id="note" value="server-value" aria-label="interaction sentinel">
    </article>
  </section>
  <div class="spacer">bottom</div>
</div>
"""

HARNESS = r"""
var state = {settings:{dataMode:'live'}};
var runtime = {status:'running'};
var page = 'dashboard';
var perfMetrics = {renderRequests:0, renderExecuted:0, renderSkipped:0};
window.__fullRenderCalls = 0;
window.__fullScheduleCalls = 0;
window.__chromeUpdates = 0;
window.__prioritySyncs = 0;
window.__nextHtml = '';
function quoteAffectsPage(){ return true; }
function watchlistEditInProgress(){ return false; }
function pageRenderer(){ return function(){ return window.__nextHtml; }; }
function updateChrome(){ window.__chromeUpdates += 1; }
function syncLivePriorityHints(){ window.__prioritySyncs += 1; }
function render(){
  window.__fullRenderCalls += 1;
  document.getElementById('main').innerHTML = window.__nextHtml;
}
function scheduleLiveRender(){
  window.__fullScheduleCalls += 1;
  clearTimeout(window.__renderTimer);
  window.__renderTimer = setTimeout(render, 50);
}
"""


def main() -> None:
    assert LIVE_RECONCILER.is_file(), f"missing reconciler: {LIVE_RECONCILER}"
    with sync_playwright() as p:
        launch_kwargs = {"headless": True}
        chrome_bin = os.environ.get("CHROME_BIN", "").strip()
        if chrome_bin:
            assert Path(chrome_bin).is_file(), f"CHROME_BIN does not exist: {chrome_bin}"
            launch_kwargs["executable_path"] = chrome_bin
        browser = p.chromium.launch(**launch_kwargs)
        page_obj = browser.new_page(viewport={"width": 1000, "height": 700})
        page_obj.set_content(INITIAL_HTML)
        page_obj.add_script_tag(content=HARNESS)
        page_obj.add_script_tag(path=str(LIVE_RECONCILER))

        assert page_obj.evaluate("Boolean(window.__DEPULSE_LIVE_DOM__)")
        assert page_obj.evaluate("window.__DEPULSE_LIVE_DOM__.version") == "18.5.2"

        page_obj.evaluate(
            """next => {
              window.__nextHtml = next;
              window.__aaplNode = document.querySelector('[data-live-key="row:AAPL"]');
              const main = document.getElementById('main');
              // Stay deliberately away from the maximum-scroll boundary. At the
              // boundary Chromium may clamp by one CSS pixel after style/text
              // reconciliation even though the viewport anchor is unchanged.
              main.scrollTop = 80;
              const note = document.getElementById('note');
              note.focus();
              note.setSelectionRange(1, 4);
            }""",
            NEXT_HTML,
        )
        page_obj.hover('[data-live-key="row:AAPL"]')
        # Playwright may scroll the target into view while establishing hover.
        # Capture the preservation baseline only after that user-interaction
        # positioning is complete, immediately before DE.PULSE reconciles.
        page_obj.evaluate("window.__scrollBefore = document.getElementById('main').scrollTop")
        page_obj.evaluate("scheduleLiveRender('AAPL')")
        page_obj.wait_for_timeout(450)

        result = page_obj.evaluate(
            """() => {
              const main = document.getElementById('main');
              const aapl = document.querySelector('[data-live-key="row:AAPL"]');
              const rows = [...document.querySelectorAll('article[data-live-key]')].map(x => x.dataset.liveKey);
              const note = document.getElementById('note');
              return {
                sameNode: window.__aaplNode === aapl,
                price: aapl.querySelector('[data-live-key="price:AAPL"]').textContent,
                className: aapl.className,
                rows,
                fullRenderCalls: window.__fullRenderCalls,
                fullScheduleCalls: window.__fullScheduleCalls,
                chromeUpdates: window.__chromeUpdates,
                prioritySyncs: window.__prioritySyncs,
                focused: document.activeElement === note,
                noteValue: note.value,
                selectionStart: note.selectionStart,
                selectionEnd: note.selectionEnd,
                scrollBefore: window.__scrollBefore,
                scrollAfter: main.scrollTop,
                renderExecuted: perfMetrics.renderExecuted
              };
            }"""
        )

        assert result["sameNode"], result
        assert result["price"] == "$101.25", result
        assert "positive" in result["className"], result
        assert result["rows"] == ["row:AAPL", "row:MSFT"], result
        assert result["fullRenderCalls"] == 0, result
        assert result["fullScheduleCalls"] == 0, result
        assert result["chromeUpdates"] == 1, result
        assert result["prioritySyncs"] == 1, result
        assert result["focused"], result
        assert result["noteValue"] == "keep-me", result
        assert (result["selectionStart"], result["selectionEnd"]) == (1, 4), result
        assert result["scrollAfter"] == result["scrollBefore"], result
        assert result["renderExecuted"] == 1, result

        # The tracked-symbol draft is a keyed interaction surface. Even if a
        # quote reconciliation runs while it is focused, it must preserve the
        # exact node, typed value, focus and selection.
        typing_next = NEXT_HTML.replace("$101.25", "$102.00").replace(
            'value="server-draft"', 'value="server-replacement"'
        )
        page_obj.evaluate(
            """next => {
              window.__nextHtml = next;
              const input = document.getElementById('master-symbol-input');
              window.__masterInputNode = input;
              input.focus();
              input.value = 'NVDA';
              input.setSelectionRange(2, 4);
            }""",
            typing_next,
        )
        page_obj.evaluate("scheduleLiveRender('AAPL')")
        page_obj.wait_for_timeout(450)
        draft_result = page_obj.evaluate(
            """() => {
              const input = document.getElementById('master-symbol-input');
              return {
                sameNode: window.__masterInputNode === input,
                focused: document.activeElement === input,
                value: input.value,
                selectionStart: input.selectionStart,
                selectionEnd: input.selectionEnd,
                price: document.querySelector('[data-live-key="price:AAPL"]').textContent,
                fullRenderCalls: window.__fullRenderCalls,
                fullScheduleCalls: window.__fullScheduleCalls
              };
            }"""
        )
        assert draft_result["sameNode"], draft_result
        assert draft_result["focused"], draft_result
        assert draft_result["value"] == "NVDA", draft_result
        assert (draft_result["selectionStart"], draft_result["selectionEnd"]) == (2, 4), draft_result
        assert draft_result["price"] == "$102.00", draft_result
        assert draft_result["fullRenderCalls"] == 0, draft_result
        assert draft_result["fullScheduleCalls"] == 0, draft_result

        # A structural event (empty symbol) must still use renderer.js's original
        # full-render contract; the incremental layer is quote-selective.
        page_obj.evaluate("scheduleLiveRender('')")
        page_obj.wait_for_timeout(150)
        structural = page_obj.evaluate(
            "({fullRenderCalls:window.__fullRenderCalls, fullScheduleCalls:window.__fullScheduleCalls})"
        )
        assert structural == {"fullRenderCalls": 1, "fullScheduleCalls": 1}, structural

        browser.close()

    print("PASS: v18.5.2 live quote rendering preserves DOM identity, focus, selection and scroll; structural events retain full render.")


if __name__ == "__main__":
    main()
