#!/usr/bin/env python3
"""Behavior-first Chromium proof for v18.5.1 Issue #12.

The harness loads the same watchlist-v18.5.1.js extension as renderer/index.html,
while retaining the actual renderer.js membership-pill handler. Only network,
bootstrap and render dependencies are stubbed. It proves:
- desk row × uses canonical global tracked-symbol removal;
- all 7 legal desk-membership combinations remove to zero;
- Undo restores exact memberships and prior desk selections;
- final membership-pill removal remains protected and separate;
- the current desk is visible and exposed with aria-current.
"""
import os
from pathlib import Path

from playwright.sync_api import sync_playwright

ROOT = Path(__file__).resolve().parents[2]
RENDERER = ROOT / "renderer" / "renderer.js"
EXTENSION = ROOT / "renderer" / "watchlist-v18.5.1.js"
INDEX = ROOT / "renderer" / "index.html"
CSS = ROOT / "renderer" / "watchlist-v18.5.1.css"
DESKS = ("day", "swing", "long")


def between(text: str, start: str, end: str, label: str) -> str:
    s = text.find(start)
    if s < 0:
        raise AssertionError(f"{label}: start anchor missing")
    e = text.find(end, s)
    if e < 0:
        raise AssertionError(f"{label}: end anchor missing")
    return text[s:e]


def main() -> None:
    renderer = RENDERER.read_text(encoding="utf-8")
    extension = EXTENSION.read_text(encoding="utf-8")
    index = INDEX.read_text(encoding="utf-8")
    css = CSS.read_text(encoding="utf-8")

    membership_binding = between(
        renderer,
        "$$('[data-desk-membership]').forEach",
        "$$('[data-master-remove]').forEach",
        "actual membership-pill binding",
    )

    assert "watchlist-v18.5.1.js?v=18.5.1" in index
    assert "watchlist-v18.5.1.css?v=18.5.1" in index
    assert "/api/master-symbol/remove" in extension
    assert "/api/master-symbol/restore" in extension
    assert "bindGlobalTrackedSymbolRemoval" in extension
    assert "aria-current=\"true\"" in extension
    assert "CURRENT" in extension
    assert "Remove ${symbol} from Tracked Symbols and all desks" in extension
    assert "/api/desk/membership" not in extension, "extension must not duplicate one-desk membership semantics"
    assert "/api/desk/membership" in membership_binding
    assert ".desk-membership-pill.current-desk" in css
    assert ".desk-membership-current" in css

    harness = r"""
const DESKS=['day','swing','long'];
let page='day';
let state={},runtime={},selected={day:'NVDA',swing:'NVDA',long:'NVDA'};
window.__members={day:[],swing:[],long:[]};
window.__calls=[];
window.__renderCount=0;
window.__toast=[];
function normalizeSymbol(v){return String(v||'').trim().toUpperCase()}
function esc(v){return String(v??'').replace(/[&<>\"']/g,c=>({'&':'&amp;','<':'&lt;','>':'&gt;','\"':'&quot;',"'":'&#39;'}[c]))}
function titleCaseText(v){return String(v||'').replace(/(^|[-_ ])([a-z])/g,(m,a,b)=>a+b.toUpperCase())}
function deskWL(k){return {id:'wl-'+k,symbols:[...(window.__members[k]||[])]}}
const deskCfg={day:{title:'Day Trade Desk'},swing:{title:'Swing Desk'},long:{title:'Long-Term Desk'}};
const $=(s,r=document)=>r.querySelector(s), $$=(s,r=document)=>[...r.querySelectorAll(s)];
function captureSaveContext(){return {scrollY:window.scrollY}}
function restoreSaveContext(){}
function toast(title,msg='',tone=''){
 window.__toast.push({title,msg,tone});
 const host=$('#header-notification');
 if(host)host.innerHTML='<span class="toast-title">'+title+'</span>';
}
function render(){
 window.__renderCount+=1;
 const host=$('#state-output');
 if(host)host.textContent=JSON.stringify(window.__members);
 const row=$('[data-desk-remove]');
 if(row){const parts=row.dataset.deskRemove.split(':');const k=parts[0],sym=parts[2];if(!(window.__members[k]||[]).includes(sym))row.remove()}
}
async function api(path,payload){
 window.__calls.push({path,payload:JSON.parse(JSON.stringify(payload||{}))});
 if(path==='/api/master-symbol/remove'){
  const sym=payload.symbol,removed={};
  for(const k of DESKS){removed[k]=(window.__members[k]||[]).includes(sym);window.__members[k]=(window.__members[k]||[]).filter(x=>x!==sym)}
  return {removed};
 }
 if(path==='/api/master-symbol/restore'){
  const sym=payload.symbol;
  for(const k of DESKS){window.__members[k]=(payload.membership&&payload.membership[k])?[sym]:[]}
  return {restored:true};
 }
 if(path==='/api/desk/membership'){
  const {desk,symbol,active}=payload;
  const total=DESKS.filter(k=>(window.__members[k]||[]).includes(symbol)).length;
  if(!active && total<=1)return {protected:true,state};
  const set=new Set(window.__members[desk]||[]);active?set.add(symbol):set.delete(symbol);window.__members[desk]=[...set];
  return {protected:false,state};
 }
 if(path==='/api/bootstrap')return {state:{members:JSON.parse(JSON.stringify(window.__members))},runtime:{status:'running'}};
 throw new Error('unexpected API '+path);
}
function mountDeskRemove(kind){
 const old=$('[data-desk-remove]');if(old)old.remove();
 const b=document.createElement('button');b.type='button';b.textContent='×';b.dataset.deskRemove=kind+':wl-'+kind+':NVDA';b.setAttribute('aria-label','legacy one-desk label');document.body.appendChild(b);return b;
}
function mountMembership(kind){
 let host=$('#membership-host');if(!host){host=document.createElement('div');host.id='membership-host';document.body.appendChild(host)}
 host.innerHTML=deskMembershipStrip('NVDA');
 return host.querySelector('[data-desk-membership="'+kind+':NVDA"]');
}
function deskMembershipStrip(){return ''}
"""

    # Preserve the actual renderer membership-pill contract inside the base bind.
    base_bind = "function bindDynamic(){\n" + membership_binding + "\n}"

    combos = [
        ("D", {"day"}),
        ("S", {"swing"}),
        ("L", {"long"}),
        ("DS", {"day", "swing"}),
        ("DL", {"day", "long"}),
        ("SL", {"swing", "long"}),
        ("DSL", {"day", "swing", "long"}),
    ]

    with sync_playwright() as p:
        kwargs = {"headless": True}
        chrome = os.environ.get("CHROME_BIN", "").strip()
        if chrome:
            assert Path(chrome).is_file(), f"CHROME_BIN missing: {chrome}"
            kwargs["executable_path"] = chrome
        browser = p.chromium.launch(**kwargs)
        pg = browser.new_page(viewport={"width": 1000, "height": 700})
        pg.set_content('<div id="header-notification"></div><div id="membership-host"></div><pre id="state-output"></pre>')
        pg.add_script_tag(content=harness)
        pg.add_script_tag(content=base_bind)
        pg.add_script_tag(content=extension)

        for name, active in combos:
            current = next(iter(active))
            initial = {k: (["NVDA"] if k in active else []) for k in DESKS}
            selected_before = {
                "day": "NVDA" if "day" in active else "SPY",
                "swing": "NVDA" if "swing" in active else "QQQ",
                "long": "NVDA" if "long" in active else "IWM",
            }
            pg.evaluate(
                """cfg=>{
                  window.__members=JSON.parse(JSON.stringify(cfg.members));
                  window.__calls=[];window.__toast=[];window.__masterUndo=null;
                  selected={...cfg.selected};page=cfg.current;
                  document.getElementById('membership-host').innerHTML=deskMembershipStrip('NVDA');
                  mountDeskRemove(cfg.current);
                  bindDynamic();
                }""",
                {"members": initial, "selected": selected_before, "current": current},
            )

            row = pg.locator('[data-desk-remove]')
            assert row.get_attribute("aria-label") == "Remove NVDA from Tracked Symbols and all desks", (name, row.get_attribute("aria-label"))
            row.click()
            pg.wait_for_function("window.__calls.some(x=>x.path==='/api/bootstrap')")
            result = pg.evaluate(
                """() => ({
                  calls:window.__calls,
                  members:window.__members,
                  rowExists:Boolean(document.querySelector('[data-desk-remove]')),
                  undoExists:Boolean(document.querySelector('[data-master-undo]')),
                  selected:{...selected}
                })"""
            )
            remove_calls = [x for x in result["calls"] if x["path"] == "/api/master-symbol/remove"]
            assert remove_calls == [{"path": "/api/master-symbol/remove", "payload": {"symbol": "NVDA"}}], (name, result)
            assert all(not result["members"][k] for k in DESKS), (name, result)
            assert not result["rowExists"], (name, result)
            assert result["undoExists"], (name, result)

            pg.locator('[data-master-undo]').click()
            pg.wait_for_function("window.__calls.some(x=>x.path==='/api/master-symbol/restore')")
            restored = pg.evaluate("({members:window.__members,calls:window.__calls,selected:{...selected}})")
            assert restored["members"] == initial, (name, restored, initial)
            assert restored["selected"] == selected_before, (name, restored, selected_before)
            restore = [x for x in restored["calls"] if x["path"] == "/api/master-symbol/restore"][-1]
            assert restore["payload"]["symbol"] == "NVDA", (name, restore)
            assert restore["payload"]["membership"] == {k: (k in active) for k in DESKS}, (name, restore)

            # Current desk is inferred from the actual active page because the
            # existing deskWatchlistCard calls deskMembershipStrip(sym) with one arg.
            pg.evaluate(
                """cfg=>{window.__members=JSON.parse(JSON.stringify(cfg.members));page=cfg.current;document.getElementById('membership-host').innerHTML=deskMembershipStrip('NVDA')}""",
                {"members": initial, "current": current},
            )
            current_button = pg.locator(f'#membership-host [data-desk-membership="{current}:NVDA"]')
            assert current_button.get_attribute("aria-current") == "true", (name, current)
            assert "CURRENT" in current_button.inner_text(), (name, current)
            assert "current-desk" in (current_button.get_attribute("class") or ""), (name, current)
            for kind in DESKS:
                btn = pg.locator(f'#membership-host [data-desk-membership="{kind}:NVDA"]')
                assert btn.get_attribute("aria-pressed") == ("true" if kind in active else "false"), (name, kind)
                if kind != current:
                    assert btn.get_attribute("aria-current") is None, (name, kind)

        # Final membership pill remains protected. This uses the actual original
        # renderer.js one-desk handler through the extension's wrapped bindDynamic.
        pg.evaluate(
            """()=>{
              window.__members={day:['NVDA'],swing:[],long:[]};window.__calls=[];window.__toast=[];page='day';
              document.getElementById('membership-host').innerHTML=deskMembershipStrip('NVDA');bindDynamic();
            }"""
        )
        pg.locator('[data-desk-membership="day:NVDA"]').click()
        pg.wait_for_function("window.__calls.some(x=>x.path==='/api/desk/membership')")
        protected = pg.evaluate("({members:window.__members,calls:window.__calls,toast:window.__toast})")
        assert protected["members"] == {"day": ["NVDA"], "swing": [], "long": []}, protected
        desk_calls = [x for x in protected["calls"] if x["path"] == "/api/desk/membership"]
        assert desk_calls[-1]["payload"] == {"desk": "day", "symbol": "NVDA", "active": False}, protected
        assert not any(x["path"] == "/api/master-symbol/remove" for x in protected["calls"]), protected
        assert any("Kept" in t["title"] for t in protected["toast"]), protected

        browser.close()

    print("PASS: production watchlist extension globally removes all 7 membership combinations with exact Undo/selection restoration; one-desk protection and current-desk accessibility remain correct.")


if __name__ == "__main__":
    main()
