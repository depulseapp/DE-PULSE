#!/usr/bin/env python3
"""Behavior-first Chromium proof for v18.5.1 Issue #12.

The test executes the actual product helper/markup/bindings extracted from
renderer.js, with only network/bootstrap/render dependencies stubbed. It proves:
- desk row × uses the canonical global-remove endpoint, not one-desk demotion;
- all 7 legal membership combinations remove to zero and Undo restores exactly;
- membership-pill final-desk protection remains a separate one-desk contract;
- the active current desk has visible + aria-current semantics.
"""
import os
from pathlib import Path

from playwright.sync_api import sync_playwright

ROOT = Path(__file__).resolve().parents[2]
RENDERER = ROOT / "renderer" / "renderer.js"


def between(text: str, start: str, end: str, label: str) -> str:
    s = text.find(start)
    if s < 0:
        raise AssertionError(f"{label}: start anchor missing")
    e = text.find(end, s)
    if e < 0:
        raise AssertionError(f"{label}: end anchor missing")
    return text[s:e]


def main() -> None:
    source = RENDERER.read_text(encoding="utf-8")
    helper = between(
        source,
        "async function removeTrackedSymbolEverywhere(symbol,ctx){",
        "function bindDynamic(){",
        "global-remove helper",
    )
    membership_fn = between(
        source,
        "function deskMembershipStrip(sym,currentDesk='')",
        "function discoveryMembership",
        "membership markup",
    )
    desk_remove_binding = between(
        source,
        "$$('[data-desk-remove]').forEach",
        "\n $$('[data-add-desk]')",
        "desk remove binding",
    )
    membership_binding = between(
        source,
        "$$('[data-desk-membership]').forEach",
        "$$('[data-master-remove]').forEach",
        "membership binding",
    )

    assert "/api/master-symbol/remove" in helper
    assert "/api/master-symbol/restore" in helper
    assert "removeTrackedSymbolEverywhere(symbol,ctx)" in desk_remove_binding
    assert "/api/desk/membership" not in desk_remove_binding
    assert "/api/desk/membership" in membership_binding
    assert "aria-current=\"true\"" in membership_fn
    assert "desk-membership-current" in membership_fn
    assert "${deskMembershipStrip(sym,kind)}" in source

    harness = r"""
const DESKS=['day','swing','long'];
let state={},runtime={},selected={day:'NVDA',swing:'NVDA',long:'NVDA'};
window.__members={day:[],swing:[],long:[]};
window.__calls=[];
window.__renderCount=0;
window.__toast=[];
window.__currentDesk='day';
function normalizeSymbol(v){return String(v||'').trim().toUpperCase()}
function esc(v){return String(v??'').replace(/[&<>\"']/g,c=>({'&':'&amp;','<':'&lt;','>':'&gt;','\"':'&quot;',"'":'&#39;'}[c]))}
function titleCaseText(v){return String(v||'').replace(/(^|[-_ ])([a-z])/g,(m,a,b)=>a+b.toUpperCase())}
function deskWL(k){return {id:'wl-'+k,symbols:[...(window.__members[k]||[])]}}
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
 const b=document.createElement('button');b.type='button';b.textContent='×';b.dataset.deskRemove=kind+':wl-'+kind+':NVDA';document.body.appendChild(b);return b;
}
function mountMembership(kind){
 const old=$('[data-desk-membership]');if(old)old.remove();
 const b=document.createElement('button');b.type='button';b.dataset.deskMembership=kind+':NVDA';b.setAttribute('aria-pressed',(window.__members[kind]||[]).includes('NVDA')?'true':'false');document.body.appendChild(b);return b;
}
"""

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
        pg.set_content('<div id="header-notification"></div><pre id="state-output"></pre>')
        pg.add_script_tag(content=harness)
        pg.add_script_tag(content=helper)
        pg.add_script_tag(content=membership_fn)

        for name, active in combos:
            current = next(iter(active))
            initial = {k: (["NVDA"] if k in active else []) for k in ("day", "swing", "long")}
            pg.evaluate(
                """cfg=>{
                  window.__members=JSON.parse(JSON.stringify(cfg.members));
                  window.__calls=[];window.__toast=[];window.__masterUndo=null;
                  selected={day:'NVDA',swing:'NVDA',long:'NVDA'};
                  window.__currentDesk=cfg.current;
                  mountDeskRemove(cfg.current);
                }""",
                {"members": initial, "current": current},
            )
            # Bind using the actual renderer.js desk-row binding.
            pg.add_script_tag(content=desk_remove_binding)
            pg.locator('[data-desk-remove]').click()
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
            assert all(not result["members"][k] for k in ("day", "swing", "long")), (name, result)
            assert not result["rowExists"], (name, result)
            assert result["undoExists"], (name, result)

            pg.locator('[data-master-undo]').click()
            pg.wait_for_function("window.__calls.some(x=>x.path==='/api/master-symbol/restore')")
            restored = pg.evaluate("({members:window.__members,calls:window.__calls,selected:{...selected}})")
            expected = initial
            assert restored["members"] == expected, (name, restored, expected)
            restore = [x for x in restored["calls"] if x["path"] == "/api/master-symbol/restore"][-1]
            assert restore["payload"]["symbol"] == "NVDA", (name, restore)
            assert restore["payload"]["membership"] == {k: (k in active) for k in ("day", "swing", "long")}, (name, restore)

            # Execute actual membership markup for the same legal combination.
            markup = pg.evaluate(
                """cfg=>{
                  window.__members=JSON.parse(JSON.stringify(cfg.members));
                  return deskMembershipStrip('NVDA',cfg.current);
                }""",
                {"members": initial, "current": current},
            )
            pg.evaluate("html=>{let h=document.getElementById('membership-host');if(!h){h=document.createElement('div');h.id='membership-host';document.body.appendChild(h)}h.innerHTML=html}", markup)
            current_button = pg.locator(f'#membership-host [data-desk-membership="{current}:NVDA"]')
            assert current_button.get_attribute("aria-current") == "true", (name, current, markup)
            assert "CURRENT" in current_button.inner_text(), (name, current, markup)
            assert "current-desk" in (current_button.get_attribute("class") or ""), (name, current, markup)
            for k in ("day", "swing", "long"):
                btn = pg.locator(f'#membership-host [data-desk-membership="{k}:NVDA"]')
                assert btn.get_attribute("aria-pressed") == ("true" if k in active else "false"), (name, k, markup)
                if k != current:
                    assert btn.get_attribute("aria-current") is None, (name, k, markup)

        # Final membership pill remains protected; it must not silently become
        # the global-removal contract.
        pg.evaluate("""()=>{window.__members={day:['NVDA'],swing:[],long:[]};window.__calls=[];window.__toast=[];mountMembership('day')}""")
        pg.add_script_tag(content=membership_binding)
        pg.locator('[data-desk-membership="day:NVDA"]').click()
        pg.wait_for_function("window.__calls.some(x=>x.path==='/api/desk/membership')")
        protected = pg.evaluate("({members:window.__members,calls:window.__calls,toast:window.__toast})")
        assert protected["members"] == {"day": ["NVDA"], "swing": [], "long": []}, protected
        desk_calls = [x for x in protected["calls"] if x["path"] == "/api/desk/membership"]
        assert desk_calls[-1]["payload"] == {"desk": "day", "symbol": "NVDA", "active": False}, protected
        assert not any(x["path"] == "/api/master-symbol/remove" for x in protected["calls"]), protected
        assert any("Kept" in t["title"] for t in protected["toast"]), protected

        browser.close()

    print("PASS: Issue #12 row × globally removes all 7 membership combinations with exact Undo; membership-pill final protection and current-desk accessibility remain correct.")


if __name__ == "__main__":
    main()
