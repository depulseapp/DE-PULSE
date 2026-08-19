#!/usr/bin/env python3
"""v18.6.1 browser regression for watchlist toggle/global removal and alert layout."""
from pathlib import Path
import os
from playwright.sync_api import sync_playwright

ROOT = Path(__file__).resolve().parents[2]
CONTRACT = (ROOT / 'renderer/watchlist-desk-contract-v18.6.1.js').read_text(encoding='utf-8')
EXT = (ROOT / 'renderer/watchlist-v18.5.1.js').read_text(encoding='utf-8')
CSS = (ROOT / 'renderer/ui-v18.6.1.css').read_text(encoding='utf-8')
INDEX = (ROOT / 'renderer/index.html').read_text(encoding='utf-8')

assert "var DESKS = Object.freeze(['day', 'swing', 'long'])" in CONTRACT
assert 'CURRENT' not in EXT
assert 'aria-current' not in EXT
assert 'aria-pressed' in EXT
assert 'watchlist-desk-contract-v18.6.1.js?v=18.6.1' in INDEX
assert INDEX.index('watchlist-desk-contract-v18.6.1.js') < INDEX.index('watchlist-v18.5.1.js')
assert 'ui-v18.6.1.css?v=18.6.1' in INDEX
assert 'justify-self:stretch!important' in CSS
assert 'width:100%!important' in CSS
assert 'text-align:center!important' in CSS

HARNESS = r"""
let state={},runtime={quotes:{}},selected={day:'NVDA',swing:'NVDA',long:'NVDA'},watchlistDraft={day:'',swing:'',long:''};
const deskCfg={day:{title:'Day Trade Desk'},swing:{title:'Swing Desk'},long:{title:'Long-Term Desk'}};
window.__members={day:[],swing:[],long:[]};window.__calls=[];window.__toasts=[];window.__failRemove=false;
const $=(s,r=document)=>r.querySelector(s), $$=(s,r=document)=>[...r.querySelectorAll(s)];
const num=v=>Number(v)||0; const esc=v=>String(v??'');
function deskWL(k){return {id:'wl-'+k,symbols:[...(window.__members[k]||[])]}}
function deskMembershipStrip(){return ''}
function captureSaveContext(){return {y:window.scrollY}} function restoreSaveContext(){} function updateChrome(){}
function toast(title,msg='',tone=''){window.__toasts.push({title,msg,tone});const h=$('#header-notification');if(h)h.innerHTML='<span class="toast-title">'+title+'</span>'}
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
 document.body.innerHTML='<div id="header-notification"></div><div id="membership"></div><button data-desk-remove="'+kind+':wl-'+kind+':NVDA">×</button><button data-master-remove="NVDA">remove all</button>';
 document.getElementById('membership').innerHTML=deskMembershipStrip('NVDA');bindDynamic();
}
"""

combos = [
    {'day'}, {'swing'}, {'long'}, {'day','swing'}, {'day','long'}, {'swing','long'}, {'day','swing','long'}
]

with sync_playwright() as p:
    kwargs={'headless':True}
    chrome=os.environ.get('CHROME_BIN','').strip()
    if chrome: kwargs['executable_path']=chrome
    browser=p.chromium.launch(**kwargs)
    pg=browser.new_page(viewport={'width':1440,'height':800})
    pg.set_content('<body></body>')
    pg.add_script_tag(content=HARNESS)
    # Production order: canonical desk contract before compatibility extension.
    pg.add_script_tag(content=CONTRACT)
    pg.add_script_tag(content=EXT)

    assert pg.evaluate('DESKS.join(",")') == 'day,swing,long'

    for active in combos:
        members={k:(['NVDA'] if k in active else []) for k in ('day','swing','long')}
        kind=next(iter(active))
        pg.evaluate('(x)=>{window.__members=x.members;window.__calls=[];window.__toasts=[];window.__failRemove=false;selected={day:"NVDA",swing:"NVDA",long:"NVDA"};mount(x.kind)}', {'members':members,'kind':kind})
        # Toggle UI is membership-only; there is never a CURRENT marker.
        assert 'CURRENT' not in pg.locator('#membership').inner_text()
        for desk in ('day','swing','long'):
            expected='true' if desk in active else 'false'
            assert pg.locator(f'[data-desk-membership="{desk}:NVDA"]').get_attribute('aria-pressed') == expected
            assert pg.locator(f'[data-desk-membership="{desk}:NVDA"]').get_attribute('aria-current') is None

        pg.locator('[data-desk-remove]').click()
        pg.wait_for_function("window.__calls.some(x=>x.path==='/api/bootstrap')")
        result=pg.evaluate('({members:window.__members,calls:window.__calls,undo:!!document.querySelector("[data-master-undo]")})')
        assert all(result['members'][k] == [] for k in ('day','swing','long')), (active,result)
        assert len([x for x in result['calls'] if x['path']=='/api/master-symbol/remove']) == 1, result
        assert result['undo'], result

        pg.locator('[data-master-undo]').click()
        pg.wait_for_function("window.__calls.some(x=>x.path==='/api/master-symbol/restore')")
        restored=pg.evaluate('window.__members')
        assert {k for k,v in restored.items() if 'NVDA' in v} == active, (active,restored)

    # Master-store removal uses the same single global-remove owner.
    pg.evaluate('()=>{window.__members={day:["NVDA"],swing:["NVDA"],long:["NVDA"]};window.__calls=[];mount("day")}')
    pg.locator('[data-master-remove="NVDA"]').click()
    pg.wait_for_function("window.__calls.some(x=>x.path==='/api/bootstrap')")
    calls=pg.evaluate('window.__calls')
    assert len([x for x in calls if x['path']=='/api/master-symbol/remove']) == 1

    # Double activation of a desk-row remove is idempotent while the action is busy.
    pg.evaluate('()=>{window.__members={day:["NVDA"],swing:["NVDA"],long:[]};window.__calls=[];mount("day")}')
    pg.locator('[data-desk-remove]').evaluate('(b)=>{b.click();b.click()}')
    pg.wait_for_function("window.__calls.some(x=>x.path==='/api/bootstrap')")
    calls=pg.evaluate('window.__calls')
    assert len([x for x in calls if x['path']=='/api/master-symbol/remove']) == 1, calls

    # Failed backend removal does not mutate memberships and produces an error toast, not ReferenceError/DESKS failure.
    pg.evaluate('()=>{window.__members={day:["NVDA"],swing:[],long:[]};window.__calls=[];window.__toasts=[];window.__failRemove=true;mount("day")}')
    pg.locator('[data-desk-remove]').click()
    pg.wait_for_function('window.__toasts.length>0')
    failed=pg.evaluate('({members:window.__members,toasts:window.__toasts})')
    assert failed['members']['day'] == ['NVDA'], failed
    assert failed['toasts'][-1]['title'] == 'Global Remove Failed', failed
    assert 'DESKS' not in failed['toasts'][-1]['msg'], failed

    browser.close()

print('PASS: v18.6.1 global remove works from desk rows and Master Market Store across all 7 membership combinations; toggle state has no CURRENT marker; failure/double-action edges and centered alert contract pass.')
