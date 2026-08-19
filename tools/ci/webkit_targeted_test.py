#!/usr/bin/env python3
"""Focused WebKit compatibility proof for renderer/UI-sensitive DE.PULSE changes.

Chrome and WebKit are the primary browser qualification engines. Chrome keeps
the broad regression suite; this WebKit proof is the primary compatibility
counterpart for core cross-engine UI behavior. Other engines remain secondary.
"""
from __future__ import annotations

from pathlib import Path
from playwright.sync_api import sync_playwright

ROOT = Path(__file__).resolve().parents[2]
CONTRACT = (ROOT / "renderer" / "watchlist-desk-contract-v18.6.1.js").read_text(encoding="utf-8")
EXT = (ROOT / "renderer" / "watchlist-v18.5.1.js").read_text(encoding="utf-8")
BASE_CSS = (ROOT / "renderer" / "styles.css").read_text(encoding="utf-8")
SETTINGS_CSS = (ROOT / "renderer" / "ui-v18.5.1.css").read_text(encoding="utf-8")
PATCH_CSS = (ROOT / "renderer" / "ui-v18.6.1.css").read_text(encoding="utf-8")

WATCHLIST_HARNESS = r"""
let state={},runtime={quotes:{}},selected={day:'NVDA',swing:'NVDA',long:'NVDA'},watchlistDraft={day:'',swing:'',long:''};
const deskCfg={day:{title:'Day Trade Desk'},swing:{title:'Swing Desk'},long:{title:'Long-Term Desk'}};
window.__members={day:[],swing:[],long:[]};window.__calls=[];window.__toasts=[];window.__failRemove=false;
const $=(s,r=document)=>r.querySelector(s), $$=(s,r=document)=>[...r.querySelectorAll(s)];
const num=v=>Number(v)||0; const esc=v=>String(v??'');
function deskWL(k){return {id:'wl-'+k,symbols:[...(window.__members[k]||[])]}}
function deskMembershipStrip(){return ''}
function captureSaveContext(){return {y:window.scrollY}} function restoreSaveContext(){} function updateChrome(){}
function toast(title,msg='',tone=''){window.__toasts.push({title,msg,tone});}
function render(){const row=$('[data-desk-remove]');if(row){const s=row.dataset.deskRemove.split(':')[2];if(!Object.values(window.__members).some(xs=>xs.includes(s)))row.remove()}}
async function api(path,payload){
 window.__calls.push({path,payload:JSON.parse(JSON.stringify(payload||{}))});
 if(path==='/api/master-symbol/remove'){
   if(window.__failRemove) throw new Error('simulated remove failure');
   const removed={};for(const k of DESKS){removed[k]=(window.__members[k]||[]).includes(payload.symbol);window.__members[k]=(window.__members[k]||[]).filter(x=>x!==payload.symbol)}return {removed};
 }
 if(path==='/api/master-symbol/restore'){for(const k of DESKS){if(payload.membership?.[k]&&!window.__members[k].includes(payload.symbol))window.__members[k].push(payload.symbol)}return {restored:true}}
 if(path==='/api/desk/membership'){
   const {desk,symbol,active}=payload,total=DESKS.filter(k=>(window.__members[k]||[]).includes(symbol)).length;
   if(!active&&total<=1)return {protected:true};const set=new Set(window.__members[desk]||[]);active?set.add(symbol):set.delete(symbol);window.__members[desk]=[...set];return {protected:false};
 }
 if(path==='/api/bootstrap')return {state:{},runtime:{quotes:{}}};
 throw new Error('unexpected '+path);
}
function bindDynamic(){}
function mount(kind='day'){
 document.body.innerHTML='<div id="membership"></div><button data-desk-remove="'+kind+':wl-'+kind+':NVDA">×</button>';
 document.getElementById('membership').innerHTML=deskMembershipStrip('NVDA');bindDynamic();
}
"""

SETTINGS_FIXTURE = """
<header class="topbar"><strong>DE.PULSE</strong><div id="header-notification">Compatibility notice</div></header>
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


def watchlist_contract(page) -> None:
    page.set_content("<body></body>")
    page.add_script_tag(content=WATCHLIST_HARNESS)
    page.add_script_tag(content=CONTRACT)
    page.add_script_tag(content=EXT)
    assert page.evaluate('DESKS.join(",")') == "day,swing,long"

    combos = [
        {"day"}, {"swing"}, {"long"},
        {"day", "swing"}, {"day", "long"}, {"swing", "long"},
        {"day", "swing", "long"},
    ]
    for active in combos:
        members = {k: (["NVDA"] if k in active else []) for k in ("day", "swing", "long")}
        kind = next(iter(active))
        page.evaluate(
            "(x)=>{window.__members=x.members;window.__calls=[];window.__toasts=[];window.__failRemove=false;mount(x.kind)}",
            {"members": members, "kind": kind},
        )
        assert "CURRENT" not in page.locator("#membership").inner_text()
        for desk in ("day", "swing", "long"):
            expected = "true" if desk in active else "false"
            toggle = page.locator(f'[data-desk-membership="{desk}:NVDA"]')
            assert toggle.get_attribute("aria-pressed") == expected
            assert toggle.get_attribute("aria-current") is None

        page.locator("[data-desk-remove]").click()
        page.wait_for_function("window.__calls.some(x=>x.path==='/api/bootstrap')")
        result = page.evaluate("({members:window.__members,calls:window.__calls,undo:!!document.querySelector('[data-master-undo]')})")
        assert all(result["members"][k] == [] for k in ("day", "swing", "long")), (active, result)
        assert len([x for x in result["calls"] if x["path"] == "/api/master-symbol/remove"]) == 1, result
        assert result["undo"], result

        page.locator("[data-master-undo]").click()
        page.wait_for_function("window.__calls.some(x=>x.path==='/api/master-symbol/restore')")
        restored = page.evaluate("window.__members")
        assert {k for k, v in restored.items() if "NVDA" in v} == active, (active, restored)

    page.evaluate("()=>{window.__members={day:['NVDA'],swing:[],long:[]};window.__calls=[];window.__toasts=[];window.__failRemove=true;mount('day')}")
    page.locator("[data-desk-remove]").click()
    page.wait_for_function("window.__toasts.length>0")
    failed = page.evaluate("({members:window.__members,toasts:window.__toasts})")
    assert failed["members"]["day"] == ["NVDA"], failed
    assert failed["toasts"][-1]["title"] == "Global Remove Failed", failed
    assert "DESKS" not in failed["toasts"][-1]["msg"], failed


def settings_layout_contract(page) -> None:
    for height in (520, 402):
        page.set_viewport_size({"width": 1280, "height": height})
        page.set_content(SETTINGS_FIXTURE)
        page.add_style_tag(content=BASE_CSS)
        page.add_style_tag(content=SETTINGS_CSS)
        page.add_style_tag(content=PATCH_CSS)
        page.wait_for_timeout(50)
        state = page.evaluate(
            """() => {
              const bar=document.querySelector('.settings-persistent-save').getBoundingClientRect();
              const button=document.querySelector('[data-save-settings]').getBoundingClientRect();
              const pane=document.querySelector('.settings-shell>.settings');
              const notice=document.querySelector('#header-notification');
              const ns=getComputedStyle(notice);
              return {
                barBottom:bar.bottom, buttonBottom:button.bottom, viewport:innerHeight,
                paneClient:pane.clientHeight, paneScroll:pane.scrollHeight, windowY:scrollY,
                noticeAlign:ns.textAlign, noticeWidth:notice.getBoundingClientRect().width
              };
            }"""
        )
        assert state["barBottom"] <= state["viewport"] - 8, (height, state)
        assert state["buttonBottom"] <= state["viewport"] - 8, (height, state)
        assert state["paneScroll"] > state["paneClient"], (height, state)
        assert state["windowY"] == 0, (height, state)
        assert state["noticeAlign"] == "center", (height, state)
        assert state["noticeWidth"] > 0, (height, state)

        page.locator(".settings-shell>.settings").evaluate("node=>{node.scrollTop=node.scrollHeight}")
        after = page.evaluate(
            """() => {
              const pane=document.querySelector('.settings-shell>.settings');
              const button=document.querySelector('[data-save-settings]').getBoundingClientRect();
              return {atBottom:Math.ceil(pane.scrollTop+pane.clientHeight)>=pane.scrollHeight,
                      buttonBottom:button.bottom,viewport:innerHeight,windowY:scrollY};
            }"""
        )
        assert after["atBottom"], (height, after)
        assert after["buttonBottom"] <= after["viewport"] - 8, (height, after)
        assert after["windowY"] == 0, (height, after)


def main() -> None:
    assert "CURRENT" not in EXT
    assert "aria-pressed" in EXT
    assert "justify-self:stretch!important" in PATCH_CSS
    assert "text-align:center!important" in PATCH_CSS

    with sync_playwright() as p:
        browser = p.webkit.launch(headless=True)
        page = browser.new_page(viewport={"width": 1440, "height": 800})
        watchlist_contract(page)
        settings_layout_contract(page)
        browser.close()

    print("PASS: primary WebKit compatibility for watchlist/global-remove, no-CURRENT membership semantics, Settings short-height save bar, and centered header alert.")


if __name__ == "__main__":
    main()
