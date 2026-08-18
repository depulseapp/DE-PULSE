#!/usr/bin/env python3
"""Behavior-first Chromium proof for v18.6 watchlist membership and add recovery.

The harness loads the production watchlist-v18.5.1.js compatibility extension
used by renderer/index.html while retaining the actual renderer.js membership
button handler. It proves:
- Day/Swing/Long pressed state is the only desk-membership indicator;
- the deprecated CURRENT/current-desk concept is absent;
- desk row × uses canonical global tracked-symbol removal;
- all 7 legal desk-membership combinations remove to zero;
- Undo restores exact memberships and prior desk selections;
- final membership-button removal remains protected and separate;
- desk, table and AI add surfaces all use canonical /api/desk/membership.
"""
import json
import os
from pathlib import Path

from playwright.sync_api import sync_playwright

ROOT = Path(__file__).resolve().parents[2]
RENDERER = ROOT / "renderer" / "renderer.js"
EXTENSION = ROOT / "renderer" / "watchlist-v18.5.1.js"
INDEX = ROOT / "renderer" / "index.html"
CSS = ROOT / "renderer" / "watchlist-v18.5.1.css"
RELEASE_IDENTITY = ROOT / "release_identity.json"
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
    release_version = json.loads(RELEASE_IDENTITY.read_text(encoding="utf-8"))["version"]

    membership_binding = between(
        renderer,
        "$$('[data-desk-membership]').forEach",
        "$$('[data-master-remove]').forEach",
        "actual membership-button binding",
    )

    assert f"watchlist-v18.5.1.js?v={release_version}" in index
    assert f"watchlist-v18.5.1.css?v={release_version}" in index
    assert "/api/master-symbol/remove" in extension
    assert "/api/master-symbol/restore" in extension
    assert "/api/desk/membership" in extension
    assert "/api/watchlists/add-symbol" not in extension
    assert "bindCanonicalDeskAdds" in extension
    assert "bindGlobalTrackedSymbolRemoval" in extension
    assert "deskMembershipStripV186" in extension
    assert "aria-pressed" in extension
    assert "aria-current" not in extension
    assert "CURRENT" not in extension
    assert "current-desk" not in extension
    assert "Remove ${symbol} from Tracked Symbols and all desks" in extension
    assert "function normalizeTrackedSymbol" in extension
    assert "normalizeSymbol" not in extension, "extension must remain self-contained"
    assert "/api/desk/membership" in membership_binding
    assert '[aria-pressed="true"]' in css
    assert ".current-desk" not in css
    assert ".desk-membership-current" not in css

    harness = r"""
const DESKS=['day','swing','long'];
let page='day';
let state={},runtime={quotes:{}},selected={day:'NVDA',swing:'NVDA',long:'NVDA'},watchlistDraft={day:'',swing:'',long:''};
window.__members={day:[],swing:[],long:[]};
window.__calls=[];
window.__renderCount=0;
window.__toast=[];
function esc(v){return String(v??'').replace(/[&<>\"']/g,c=>({'&':'&amp;','<':'&lt;','>':'&gt;','\"':'&quot;',"'":'&#39;'}[c]))}
function titleCaseText(v){return String(v||'').replace(/(^|[-_ ])([a-z])/g,(m,a,b)=>a+b.toUpperCase())}
function deskWL(k){return {id:'wl-'+k,symbols:[...(window.__members[k]||[])]}}
const deskCfg={day:{title:'Day Trade Desk'},swing:{title:'Swing Desk'},long:{title:'Long-Term Desk'}};
const $=(s,r=document)=>r.querySelector(s), $$=(s,r=document)=>[...r.querySelectorAll(s)];
const num=v=>Number(v)||0;
function captureSaveContext(){return {scrollY:window.scrollY}}
function restoreSaveContext(){}
function updateChrome(){}
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
 if(path==='/api/bootstrap')return {state:{members:JSON.parse(JSON.stringify(window.__members))},runtime:{status:'running',quotes:{}}};
 throw new Error('unexpected API '+path);
}
function mountDeskRemove(kind){
 const old=$('[data-desk-remove]');if(old)old.remove();
 const b=document.createElement('button');b.type='button';b.textContent='×';b.dataset.deskRemove=kind+':wl-'+kind+':NVDA';b.setAttribute('aria-label','legacy one-desk label');document.body.appendChild(b);return b;
}
function deskMembershipStrip(){return ''}
"""

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

            # Membership, not the current page, is the sole source of selected state.
            pg.evaluate(
                """cfg=>{window.__members=JSON.parse(JSON.stringify(cfg.members));page=cfg.current;document.getElementById('membership-host').innerHTML=deskMembershipStrip('NVDA')}""",
                {"members": initial, "current": current},
            )
            host_text = pg.locator('#membership-host').inner_text()
            assert "CURRENT" not in host_text, (name, host_text)
            for kind in DESKS:
                btn = pg.locator(f'#membership-host [data-desk-membership="{kind}:NVDA"]')
                expected = "true" if kind in active else "false"
                assert btn.get_attribute("aria-pressed") == expected, (name, kind, expected)
                assert btn.get_attribute("aria-current") is None, (name, kind)
                classes = btn.get_attribute("class") or ""
                assert ("active" in classes) == (kind in active), (name, kind, classes)

        # Final membership button remains protected. This uses the actual renderer
        # one-desk handler through the extension's wrapped bindDynamic.
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

        # All desk-add surfaces must now target canonical per-user membership.
        pg.evaluate(
            """()=>{
              window.__members={day:[],swing:[],long:[]};window.__calls=[];window.__toast=[];
              watchlistDraft={day:'',swing:'',long:''};
              document.body.insertAdjacentHTML('beforeend',
                '<input data-add-input="swing" value="NKE"><button data-add-desk="swing">Add</button>'+ 
                '<input data-add-table-input="day" value="MSFT"><button data-add-desk-table="day">Add table</button>'+ 
                '<button data-ai-add-desk="long:AAPL">AI Add</button>');
              bindDynamic();
            }"""
        )
        pg.locator('[data-add-desk="swing"]').click()
        pg.wait_for_function("window.__calls.filter(x=>x.path==='/api/desk/membership').length>=1")
        pg.locator('[data-add-desk-table="day"]').click()
        pg.wait_for_function("window.__calls.filter(x=>x.path==='/api/desk/membership').length>=2")
        pg.locator('[data-ai-add-desk="long:AAPL"]').click()
        pg.wait_for_function("window.__calls.filter(x=>x.path==='/api/desk/membership').length>=3")
        added = pg.evaluate("({members:window.__members,calls:window.__calls})")
        assert not any(x["path"] == "/api/watchlists/add-symbol" for x in added["calls"]), added
        add_calls = [x for x in added["calls"] if x["path"] == "/api/desk/membership"]
        assert add_calls[0]["payload"] == {"symbol": "NKE", "desk": "swing", "active": True}, added
        assert add_calls[1]["payload"] == {"symbol": "MSFT", "desk": "day", "active": True}, added
        assert add_calls[2]["payload"] == {"symbol": "AAPL", "desk": "long", "active": True}, added
        assert added["members"] == {"day": ["MSFT"], "swing": ["NKE"], "long": ["AAPL"]}, added

        pg.evaluate("document.getElementById('membership-host').innerHTML=deskMembershipStrip('NKE')")
        assert pg.locator('[data-desk-membership="swing:NKE"]').get_attribute("aria-pressed") == "true"
        assert pg.locator('[data-desk-membership="day:NKE"]').get_attribute("aria-pressed") == "false"
        assert pg.locator('[data-desk-membership="long:NKE"]').get_attribute("aria-pressed") == "false"
        assert "CURRENT" not in pg.locator('#membership-host').inner_text()

        browser.close()

    print("PASS: v18.6 membership buttons are canonical pressed-state controls with no CURRENT marker; global remove/Undo/protection remain intact and desk/table/AI adds use canonical membership state.")


if __name__ == "__main__":
    main()
