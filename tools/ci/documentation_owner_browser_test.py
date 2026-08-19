#!/usr/bin/env python3
from __future__ import annotations

import argparse
import os
from pathlib import Path
from playwright.sync_api import sync_playwright

ROOT=Path(__file__).resolve().parents[2]
OWNER=(ROOT/'renderer'/'documentation-ui.js').read_text(encoding='utf-8')
ACCESS=(ROOT/'renderer'/'documentation-access-v18.6.js').read_text(encoding='utf-8')

HARNESS=r"""
let authPrincipal={role:'USER'};
let documentationTab='user';
let docCache={user:'# User\nHello **DE.PULSE**',developer:'# Developer\nRestricted',limitations:null};
let page='documentation';
let renderCalls=0;
const DOC_BRAND='<span class="brand">DE.PULSE</span>';
const esc=v=>String(v??'').replace(/[&<>"']/g,c=>({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c]));
function docInline(v){return esc(v).replace(/\*\*([^*]+)\*\*/g,'<strong>$1</strong>')}
function architectureDiagram(kind){return '<section class="legacy-diagram" data-kind="'+esc(kind)+'"></section>'}
function render(){renderCalls++}
function renderMarkdown(){return 'LEGACY_MARKDOWN'}
async function hydrateDocumentation(){return 'LEGACY_HYDRATE'}
function renderDocumentation(){return 'LEGACY_RENDER'}
window.__fetchCalls=[];
window.fetch=async url=>{window.__fetchCalls.push(url);return {ok:true,status:200,text:async()=> '# Loaded\nBrowser proof'}};
"""


def run(engine:str)->None:
    with sync_playwright() as p:
        if engine=='webkit':
            browser=p.webkit.launch(headless=True)
        else:
            chrome=os.environ.get('CHROME_BIN','/usr/bin/google-chrome')
            browser=p.chromium.launch(headless=True,executable_path=chrome)
        page=browser.new_page(viewport={"width":1280,"height":720})
        page.set_content('<main id="host"></main>')
        page.add_script_tag(content=HARNESS)
        page.add_script_tag(content=OWNER)
        owner=page.evaluate("__DE_PULSE_DOCUMENTATION_UI__.registry()")
        assert owner['owner']=='renderer/documentation-ui.js',owner
        assert owner['state']=='ACTIVE_OWNER_WITH_LEGACY_FALLBACK',owner
        assert 'legacy-architecture-diagram' in owner['dependencies'],owner

        user=page.evaluate("renderDocumentation()")
        assert 'data-render-owner="documentation-ui"' in user,user
        assert '<strong>DE.PULSE</strong>' in user,user
        markdown=page.evaluate("renderMarkdown('## Heading\\n- Item\\n[[diagram:overall]]')")
        assert '<h2>Heading</h2>' in markdown,markdown
        assert '<li>Item</li>' in markdown,markdown
        assert 'data-kind="overall"' in markdown,markdown

        page.add_script_tag(content=ACCESS)
        restricted=page.evaluate("()=>{authPrincipal={role:'USER'};documentationTab='developer';return renderDocumentation()}")
        assert 'data-doc-tab="developer"' not in restricted,restricted
        assert page.evaluate('documentationTab')=='user'
        admin=page.evaluate("()=>{authPrincipal={role:'ADMIN'};documentationTab='developer';return renderDocumentation()}")
        assert 'data-doc-tab="developer"' in admin,admin
        assert 'data-render-owner="documentation-ui"' in admin,admin

        hydration=page.evaluate("""async()=>{authPrincipal={role:'USER'};documentationTab='limitations';docCache.limitations=null;await hydrateDocumentation();return {tab:documentationTab,cache:docCache.limitations,calls:window.__fetchCalls,renderCalls}}""")
        assert hydration['tab']=='limitations',hydration
        assert hydration['calls']==['/docs/limitations.md'],hydration
        assert hydration['cache'].startswith('# Loaded'),hydration
        assert hydration['renderCalls']==1,hydration
        browser.close()
    print(f'PASS: Documentation capability owner + role decorator on {engine}')


def main()->None:
    parser=argparse.ArgumentParser()
    parser.add_argument('--engine',choices=('chrome','webkit'),required=True)
    args=parser.parse_args()
    run(args.engine)


if __name__=='__main__':
    main()
