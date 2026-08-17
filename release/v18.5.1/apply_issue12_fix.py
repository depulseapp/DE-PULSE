#!/usr/bin/env python3
"""Materialize the v18.5.1 Issue #12 watchlist semantics fix.

This is deliberately idempotent and assertion-heavy because renderer.js/styles.css
are large generated-era files. The script changes only the frozen Issue #12
contracts:
- desk row × => global tracked-symbol removal with exact Undo snapshot
- membership pills retain one-desk/final-membership semantics
- current desk receives visible + accessible non-color-only semantics
"""
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
RENDERER = ROOT / "renderer" / "renderer.js"
STYLES = ROOT / "renderer" / "styles.css"


def replace_once(text: str, old: str, new: str, label: str) -> str:
    count = text.count(old)
    if count == 0 and new in text:
        return text
    if count != 1:
        raise SystemExit(f"{label}: expected exactly one source match, found {count}")
    return text.replace(old, new, 1)


def replace_between(text: str, start: str, end: str, replacement: str, label: str) -> str:
    if replacement in text:
        return text
    start_i = text.find(start)
    if start_i < 0:
        raise SystemExit(f"{label}: start anchor not found")
    end_i = text.find(end, start_i)
    if end_i < 0:
        raise SystemExit(f"{label}: end anchor not found")
    if text.find(start, start_i + 1) >= 0 and text.find(start, start_i + 1) < end_i:
        raise SystemExit(f"{label}: ambiguous repeated start anchor")
    return text[:start_i] + replacement + text[end_i:]


def main() -> None:
    js = RENDERER.read_text(encoding="utf-8")
    css = STYLES.read_text(encoding="utf-8")

    old_membership = """function deskMembershipStrip(sym){const target=String(sym||'').toUpperCase(),items=[['day','DAY'],['swing','SWING'],['long','LONG']];return `<span class=\"desk-membership-strip\" aria-label=\"Desk membership for ${esc(target)}\">${items.map(([kind,label])=>{const active=(deskWL(kind).symbols||[]).includes(target);return `<button type=\"button\" class=\"desk-membership-pill ${active?'active':''}\" aria-pressed=\"${active?'true':'false'}\" data-desk-membership=\"${kind}:${esc(target)}\" title=\"${active?`${target} is in ${kind==='long'?'Long-Term':kind==='day'?'Day Trade':'Swing'} · click to remove when another desk remains`:`Add ${target} to ${kind==='long'?'Long-Term':kind==='day'?'Day Trade':'Swing'}`}\">${label}</button>`}).join('')}</span>`}"""
    new_membership = """function deskMembershipStrip(sym,currentDesk=''){const target=String(sym||'').toUpperCase(),current=String(currentDesk||'').toLowerCase(),items=[['day','DAY'],['swing','SWING'],['long','LONG']];return `<span class=\"desk-membership-strip\" aria-label=\"Desk membership for ${esc(target)}\">${items.map(([kind,label])=>{const active=(deskWL(kind).symbols||[]).includes(target),isCurrent=active&&kind===current;return `<button type=\"button\" class=\"desk-membership-pill ${active?'active':''}${isCurrent?' current-desk':''}\" aria-pressed=\"${active?'true':'false'}\"${isCurrent?' aria-current=\"true\"':''} data-desk-membership=\"${kind}:${esc(target)}\" title=\"${isCurrent?`${label} is the current desk for ${target}`:active?`${target} is in ${kind==='long'?'Long-Term':kind==='day'?'Day Trade':'Swing'} · click to remove when another desk remains`:`Add ${target} to ${kind==='long'?'Long-Term':kind==='day'?'Day Trade':'Swing'}`}\">${label}${isCurrent?'<span class=\"desk-membership-current\">CURRENT</span>':''}</button>`}).join('')}</span>`}"""
    js = replace_once(js, old_membership, new_membership, "current-desk membership markup")
    js = replace_once(js, "${deskMembershipStrip(sym)}", "${deskMembershipStrip(sym,kind)}", "desk card current-desk context")

    helper = r'''async function removeTrackedSymbolEverywhere(symbol,ctx){
 const sym=normalizeSymbol(symbol);if(!sym)return null;
 const selectedBefore={...selected};
 const res=await api('/api/master-symbol/remove',{symbol:sym});
 const removed=res.removed||{};
 const boot=await api('/api/bootstrap');state=boot.state;runtime=boot.runtime;
 for(const k of DESKS){if(selected[k]===sym)selected[k]=deskWL(k).symbols[0]||''}
 render();restoreSaveContext(ctx);
 toast(`${sym} Removed from Tracked Symbols`,'Removed from every desk. Undo restores the exact previous desk memberships.','warning');
 window.__masterUndo={sym,membership:removed,selected:selectedBefore,expires:Date.now()+9000};
 const host=$('#header-notification');
 if(host){
  host.innerHTML+=` <button class="toast-undo" data-master-undo="${esc(sym)}">UNDO</button>`;
  host.querySelector('[data-master-undo]')?.addEventListener('click',async ev=>{
   ev.stopPropagation();const u=window.__masterUndo;if(!u||u.sym!==sym||Date.now()>u.expires)return;
   const undoCtx=captureSaveContext();await api('/api/master-symbol/restore',{symbol:sym,membership:u.membership});
   const boot2=await api('/api/bootstrap');state=boot2.state;runtime=boot2.runtime;
   selected={...selected,...u.selected};window.__masterUndo=null;render();restoreSaveContext(undoCtx);
   toast(`${sym} Restored`,'Previous desk memberships restored.','success');
  });
 }
 return res;
}
'''
    if "async function removeTrackedSymbolEverywhere(" not in js:
        js = replace_once(js, "function bindDynamic(){", helper + "function bindDynamic(){", "shared global-remove helper")

    master_start = "$$('[data-master-remove]').forEach"
    master_end = "\n $('[data-master-remove-all]')"
    master_replacement = "$$('[data-master-remove]').forEach(b=>b.onclick=async e=>{e.preventDefault();e.stopPropagation();const ctx=captureSaveContext(),sym=b.dataset.masterRemove;try{await removeTrackedSymbolEverywhere(sym,ctx)}catch(err){toast('Global Remove Failed',err.message,'error');restoreSaveContext(ctx)}});"
    js = replace_between(js, master_start, master_end, master_replacement, "Dashboard global remove binding")

    desk_start = "$$('[data-desk-remove]').forEach"
    desk_end = "\n $$('[data-add-desk]')"
    desk_replacement = "$$('[data-desk-remove]').forEach(b=>b.onclick=async e=>{e.preventDefault();e.stopPropagation();if(b.dataset.busy==='1')return;b.dataset.busy='1';b.disabled=true;const ctx=captureSaveContext(),[k,id,symbol]=b.dataset.deskRemove.split(':');try{await removeTrackedSymbolEverywhere(symbol,ctx)}catch(err){b.disabled=false;b.dataset.busy='0';toast('Global Remove Failed',err.message,'error');restoreSaveContext(ctx)}});"
    js = replace_between(js, desk_start, desk_end, desk_replacement, "desk-row global remove binding")

    css_marker = "/* v18.5.1 Issue #12 current-desk semantics */"
    css_block = r'''

/* v18.5.1 Issue #12 current-desk semantics */
.desk-membership-pill.current-desk{
  font-weight:800;
  outline:2px solid currentColor;
  outline-offset:1px;
}
.desk-membership-current{
  display:inline-block;
  margin-left:.35rem;
  padding:.05rem .28rem;
  border:1px solid currentColor;
  border-radius:999px;
  font-size:.58rem;
  font-weight:800;
  letter-spacing:.04em;
  line-height:1.15;
}
'''
    if css_marker not in css:
        css += css_block

    # Closure assertions: these are implementation invariants, not substitutes for
    # the browser/API regression tests that run after materialization.
    required = [
        "async function removeTrackedSymbolEverywhere(symbol,ctx)",
        "await api('/api/master-symbol/remove',{symbol:sym})",
        "await api('/api/master-symbol/restore',{symbol:sym,membership:u.membership})",
        "${deskMembershipStrip(sym,kind)}",
        "aria-current=\"true\"",
        "desk-membership-current",
        "await removeTrackedSymbolEverywhere(symbol,ctx)",
    ]
    missing = [token for token in required if token not in js]
    if missing:
        raise SystemExit(f"Issue #12 materialization invariant missing: {missing}")
    if "const res=await api('/api/desk/membership',{desk:k,symbol,active:false})" in js:
        raise SystemExit("legacy desk-row one-desk removal handler still present")

    RENDERER.write_text(js, encoding="utf-8")
    STYLES.write_text(css, encoding="utf-8")
    print("PASS: materialized Issue #12 global row removal, exact Undo, and current-desk semantics")


if __name__ == "__main__":
    main()
