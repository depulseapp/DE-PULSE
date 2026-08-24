#!/usr/bin/env python3
"""Behavior-first Chromium proof for the v18.5.2 tracked-symbol hotfix.

Loads the production panel renderer and exact production add-symbol binding. It
proves a typed symbol survives input normalization, adds on the first click, and
uses the requested desktop/control layout with responsive stacking.
"""
import os
from pathlib import Path

from playwright.sync_api import sync_playwright

ROOT = Path(__file__).resolve().parents[2]
RENDERER = (ROOT / "renderer" / "renderer.js").read_text(encoding="utf-8")
BASE_CSS = (ROOT / "renderer" / "styles.css").read_text(encoding="utf-8")
UI_CSS = (ROOT / "renderer" / "ui-layout-contracts.css").read_text(encoding="utf-8")


def between(text: str, start: str, end: str, label: str) -> str:
    s = text.find(start)
    if s < 0:
        raise AssertionError(f"{label}: start anchor missing")
    e = text.find(end, s)
    if e < 0:
        raise AssertionError(f"{label}: end anchor missing")
    return text[s:e]


PANEL_SOURCE = between(
    RENDERER,
    "function masterDeskSymbols()",
    "function miTone",
    "tracked-symbol panel",
)
ADD_BINDING = between(
    RENDERER,
    "const masterInput=$('[data-master-add-input]')",
    "$$('[data-mark-research-reviewed]')",
    "tracked-symbol add binding",
)

HARNESS = r"""
let state={},runtime={},masterSymbolDraft='';
window.__members={day:[],swing:[],long:[]};
window.__calls=[];
window.__toast=[];
const $=(s,r=document)=>r.querySelector(s), $$=(s,r=document)=>[...r.querySelectorAll(s)];
function esc(v){return String(v??'').replace(/[&<>"']/g,c=>({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c]))}
function deskWL(k){return {symbols:[...(window.__members[k]||[])]}}
function quoteState(){return {light:'green',label:'CURRENT'}}
function statusDot(tone,label){return '<span class="status-dot '+tone+'" aria-label="'+label+'"></span>'}
function captureSaveContext(){return {}}
function restoreSaveContext(){}
function toast(title,msg='',tone=''){window.__toast.push({title,msg,tone})}
async function api(path,payload){
  window.__calls.push({path,payload:JSON.parse(JSON.stringify(payload||{}))});
  if(path==='/api/master-symbol/add'){
    for(const kind of ['day','swing','long']){
      if(!window.__members[kind].includes(payload.symbol))window.__members[kind].push(payload.symbol);
    }
    return {changed:true,message:'Added to all desks.'};
  }
  if(path==='/api/bootstrap')return {state:{},runtime:{status:'running'}};
  throw new Error('unexpected API '+path);
}
function render(){
  document.getElementById('panel-host').innerHTML=masterMarketSymbolsPanel();
  bindMaster();
}
"""


def box(page, selector):
    loc = page.locator(selector)
    assert loc.count() == 1, selector
    value = loc.bounding_box()
    assert value and value["width"] > 0 and value["height"] > 0, (selector, value)
    return value


def main() -> None:
    assert 'id="master-symbol-input"' in PANEL_SOURCE
    assert 'data-live-key="master-symbol-input"' in PANEL_SOURCE
    assert 'data-live-key="master-symbol-controls"' in PANEL_SOURCE
    assert "masterSymbolDraft=next" in ADD_BINDING
    assert "masterSymbolDraft=''" in ADD_BINDING
    assert "masterSymbolDraft=sym" in ADD_BINDING
    assert "master-symbols-heading-row" in UI_CSS

    with sync_playwright() as p:
        kwargs = {"headless": True}
        chrome = os.environ.get("CHROME_BIN", "").strip()
        if chrome:
            assert Path(chrome).is_file(), f"CHROME_BIN missing: {chrome}"
            kwargs["executable_path"] = chrome
        browser = p.chromium.launch(**kwargs)
        page = browser.new_page(viewport={"width": 1440, "height": 800})
        page.set_content(f"<style>{BASE_CSS}\n{UI_CSS}</style><div id='panel-host'></div>")
        page.add_script_tag(content=HARNESS)
        page.add_script_tag(content=PANEL_SOURCE)
        page.add_script_tag(content="function bindMaster(){\n"+ADD_BINDING+"\n}")
        page.evaluate("render()")

        ticker = page.locator("[data-master-add-input]")
        ticker.fill(" nke ")
        assert ticker.input_value() == "NKE"
        assert page.evaluate("masterSymbolDraft") == "NKE"

        page.locator("[data-master-add]").click()
        page.wait_for_function("window.__calls.some(x=>x.path==='/api/bootstrap')")
        result = page.evaluate(
            """() => ({
              calls:window.__calls,
              members:window.__members,
              draft:masterSymbolDraft,
              input:document.querySelector('[data-master-add-input]').value,
              count:document.querySelector('.count-pill').textContent,
              removeDisabled:document.querySelector('[data-master-remove-all]').disabled,
              chips:[...document.querySelectorAll('[data-master-symbol]')].map(x=>x.dataset.masterSymbol)
            })"""
        )
        add_calls = [x for x in result["calls"] if x["path"] == "/api/master-symbol/add"]
        assert add_calls == [{"path": "/api/master-symbol/add", "payload": {"symbol": "NKE"}}], result
        assert result["members"] == {"day": ["NKE"], "swing": ["NKE"], "long": ["NKE"]}, result
        assert result["draft"] == "", result
        assert result["input"] == "", result
        assert result["count"] == "1 tracked", result
        assert not result["removeDisabled"], result
        assert result["chips"] == ["NKE"], result

        # Desktop: copy and controls share one row; input/Add/Remove are ordered
        # side by side without shrinking the input below its usable width.
        copy = box(page, ".master-symbols-copy")
        controls = box(page, ".master-add-row")
        input_box = box(page, "[data-master-add-input]")
        add_box = box(page, "[data-master-add]")
        remove_box = box(page, "[data-master-remove-all]")
        assert min(copy["y"]+copy["height"], controls["y"]+controls["height"]) > max(copy["y"], controls["y"]), (copy, controls)
        assert input_box["x"] < add_box["x"] < remove_box["x"], (input_box, add_box, remove_box)
        assert abs(input_box["y"]-add_box["y"]) <= 2 and abs(add_box["y"]-remove_box["y"]) <= 2
        assert input_box["width"] >= 220, input_box

        # Narrow: the heading context stays readable and controls stack only
        # after the desktop row no longer fits.
        page.set_viewport_size({"width": 650, "height": 900})
        page.wait_for_timeout(30)
        copy = box(page, ".master-symbols-copy")
        controls = box(page, ".master-add-row")
        input_box = box(page, "[data-master-add-input]")
        add_box = box(page, "[data-master-add]")
        assert copy["y"] + copy["height"] <= controls["y"] + 2, (copy, controls)
        assert input_box["y"] < add_box["y"], (input_box, add_box)
        assert document_overflow(page) <= 1

        browser.close()

    print("PASS: tracked ticker normalizes and adds on the first click; the draft clears only after success and the requested desktop/responsive layout is enforced.")


def document_overflow(page):
    return page.evaluate("document.documentElement.scrollWidth-window.innerWidth")


if __name__ == "__main__":
    main()
