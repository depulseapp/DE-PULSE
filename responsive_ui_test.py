import sys, json, time
from pathlib import Path
from playwright.sync_api import sync_playwright
import re
_renderer_source=(Path.cwd()/"renderer"/"renderer.js").read_text()
_version_match=re.search(r"const EXPECTED_RELEASE_VERSION='([^']+)'",_renderer_source)
_build_match=re.search(r"const EXPECTED_BUILD_ID='([^']+)'",_renderer_source)
if not _version_match or not _build_match:
    raise RuntimeError('canonical renderer release identity unavailable')
EXPECTED_RELEASE_VERSION=_version_match.group(1)
EXPECTED_BUILD_ID=_build_match.group(1)
url=sys.argv[1] if len(sys.argv)>1 else "local-fixture"
viewports=[
 (820,1050),(820,1180),(1024,768),(1180,820),(1280,720),
 (1280,800),(1280,832),(1366,650),(1440,900),(1512,982),
 (1536,864),(1680,1050),(1728,1117),(1920,1080),(2560,1080),
]
import os
_slice=os.environ.get('DEPULSE_VIEWPORT_SLICE','').strip()
if _slice:
    _a,_b=[int(x) for x in _slice.split(':',1)]
    viewports=viewports[_a:_b]
surfaces=[
 ('dashboard',None),('market-intelligence',None),('day',None),('swing',None),('long',None),
 ('discovery','day'),('discovery','swing'),('discovery','long'),
 ('research','overview'),('research','earnings'),('research','fundamentals'),('research','sec'),('research','catalysts'),('research','technical'),
 ('ai',None),('news',None),('earnings',None),('filings',None),
 ('maintenance',None),('documentation',None),('settings',None),
]
errors=[]; checks=0
with sync_playwright() as p:
    browser=p.chromium.launch(headless=True, executable_path='/usr/bin/chromium', args=['--no-sandbox','--disable-dev-shm-usage'])
    page=browser.new_page(viewport={'width':1280,'height':800})
    page.on('pageerror', lambda exc: errors.append('pageerror: '+str(exc)))
    def console(msg):
        if msg.type=='error': errors.append('console: '+msg.text)
    page.on('console', console)
    html=(Path.cwd()/"renderer"/"index.html").read_text()
    import re
    html=re.sub(r'<meta http-equiv="Content-Security-Policy"[^>]*>','',html)
    html=re.sub(r'<link rel="stylesheet"[^>]*>','',html)
    html=re.sub(r'<script src="renderer\.js[^>]*></script>','',html)
    page.set_content(html,wait_until='load',timeout=30000)
    page.add_style_tag(content=(Path.cwd()/"renderer"/"styles.css").read_text())
    page.evaluate("window.__DEPULSE_TEST__=true")
    page.add_script_tag(content=_renderer_source)
    page.wait_for_function("window.DePulseLogic",timeout=30000)
    page.evaluate("([version,buildId])=>{window.__DEPULSE_EXPECTED_VERSION__=version;window.__DEPULSE_EXPECTED_BUILD_ID__=buildId}",[EXPECTED_RELEASE_VERSION,EXPECTED_BUILD_ID])
    page.evaluate("""()=>{
      const now=Date.now(), syms=['SPY','QQQ','DIA','IWM','XLK','XLC','XLY','XLP','XLE','XLF','XLV','XLI','XLB','XLRE','XLU','TLT','GLD','SLV','USO','SMH','EWT','EWY','AAA','NVDA','AMD','MSFT','META','ACHR'];
      const quotes={}; const bars={};
      for(const [idx,s] of syms.entries()) { const px=50+idx*7; quotes[s]={symbol:s,price:px,previousClose:px*.99,changePercent:1+idx*.02,updatedAt:now,providerTimestamp:now,source:'qa-real-fixture',dataState:'current'}; const a=[]; for(let i=0;i<80;i++){const p=px*(.9+i*.002);a.push({o:p*.998,h:p*1.01,l:p*.99,c:p,v:1000000+i*1000})} bars[s]={daily:a,intraday:a,weekly:a.filter((_,i)=>i%5===0)} }
      quotes.VIX={symbol:'VIX',price:16,previousClose:17,changePercent:-5.8,updatedAt:now,source:'qa-vix',dataState:'current'}; bars.VIX=bars.SPY;
      const state={version:window.__DEPULSE_EXPECTED_VERSION__,buildId:window.__DEPULSE_EXPECTED_BUILD_ID__,settings:{dataMode:'demo',overnightDataMode:'auto',aiProvider:'groq',groqModel:'openai/gpt-oss-120b',openRouterMode:'fast',openRouterSpecificModel:'openai/gpt-5.6-sol',geminiModel:'gemini-3.1-flash-lite',secEmail:'qa@example.com',autoStart:false,earningsPenalty:10,marketContext:15,signalProfile:'balanced',swingWatchlistId:'swing',dayWatchlistId:'day',longWatchlistId:'long',discoveryWatchlistId:'discovery',dayEnabled:true,swingEnabled:true,longEnabled:true,globalProviderMode:'auto',optionsDataMode:'auto',macroEventModeEnabled:true},watchlists:[{id:'swing',name:'Swing Watchlist',symbols:['AAA','NVDA']},{id:'day',name:'Day Trade Watchlist',symbols:['AAA','AMD']},{id:'long',name:'Long-Term Watchlist',symbols:['AAA','MSFT']},{id:'discovery',name:'Discovery Watchlist',symbols:['META']}],ui:{selectedTicker:'AAA',scopeType:'watchlist',watchlistId:'swing'},providerStatus:{options:{ok:true,status:'Connected',message:'Options Intelligence access is working.',checkedAt:now}},cacheInfo:{sizeBytes:2048,cachedSymbols:20,lastUpdated:now},lastCacheCleared:now-5000,configDir:'/qa',cachePath:'/qa/market-cache.json',maintenanceLastRun:0,hasFinnhubKey:true,hasAlpacaKey:true,hasAlpacaSecret:true,hasGroqKey:true,hasOpenRouterKey:true,hasGeminiKey:true,hasFREDKey:true,hasBLSKey:false,hasEIAKey:true,hasTwelveDataKey:false};
      const runtime={status:'running',mode:'demo',message:'QA fixture',quotes,bars,history:{},fundamentals:{AAA:{revenueGrowth:20,epsGrowth:22,roe:25,netMargin:18,debtToEquity:40,pe:25,forwardPe:22,freeCashFlow:5000000000}},news:[{headline:'AAA launches new product',summary:'material catalyst',source:'QA',symbols:['AAA'],datetime:Math.floor(now/1000),url:'#'}],earnings:[{symbol:'AAA',date:new Date(now+86400000*5).toISOString().slice(0,10),hour:'amc',epsEstimate:2.1,revenueEstimate:5000000000}],filings:[{id:'1',symbol:'AAA',company:'AAA Inc',form:'8-K',filedAt:new Date(now-86400000).toISOString(),description:'Material update',meaning:'Material Company Update',category:'material',url:'#'},{id:'f4',symbol:'AAA',company:'AAA Inc',form:'4',filedAt:new Date(now-3600000).toISOString(),transactionDate:new Date(now-86400000).toISOString().slice(0,10),description:'BUY · Open-market/private purchase · CFO · 1,000 shares · $25.50 avg · $26K',meaning:'Insider Buy',category:'insider',signal:'Buy',actor:'Jane Example',role:'CFO',transactionCode:'P',transactionType:'Open-market/private purchase',shares:1000,price:25.5,value:25500,ownershipAfter:12000,url:'#'}],secIntelligence:{AAA:{symbol:'AAA',insiderBuys:1,insiderSells:0,insiderOthers:1,offeringRisk:'Low',institutional:'13F · Quarterly',latestForm:'4',latestFiledAt:new Date(now-3600000).toISOString(),signals:[],recentTransactions:[{symbol:'AAA',form:'4',filedAt:new Date(now-3600000).toISOString(),transactionDate:new Date(now-86400000).toISOString().slice(0,10),signal:'BUY',actor:'Jane Example',role:'CFO',transactionCode:'P',transactionType:'BUY',shares:1000,price:25.5,value:25500,ownershipAfter:12000},{symbol:'AAA',form:'4',filedAt:new Date(now-7200000).toISOString(),signal:'OTHER',actor:'John Example',role:'Director',transactionCode:'A',transactionType:'OTHER',shares:500,ownershipAfter:5000}]}},scanner:{mode:'day',status:'ready',message:'QA',scanned:2400,updatedAt:now,results:[{symbol:'ACHR',price:6.63,score:86,rank:1,changePercent:18.73,gapPercent:14.87,relativeVolume:.81,dollarVolume:6400000,spreadPercent:.15,trendScore:78,fundamentalScore:72,reasons:['Tight current spread','18.7% session momentum with strong participation']},{symbol:'AAA',price:204.25,score:82,rank:2,changePercent:3.86,gapPercent:1.25,relativeVolume:1.72,dollarVolume:185000000,spreadPercent:.04,trendScore:74,fundamentalScore:68,reasons:['Strong momentum','Constructive multi-horizon trend']}] },lastUpdated:{quotes:now,history:now,fundamentals:now,news:now,earnings:now,filings:now,vix:now,ai:now},health:{quotes:'streaming','quotes-rest':'healthy','alpaca-stream':'healthy','alpaca-live':'healthy',history:'healthy',fundamentals:'healthy',news:'healthy',earnings:'healthy',filings:'healthy','cache-refresh':'healthy',scanner:'ready',vix:'healthy',ai:'ready',global:'healthy',macro:'healthy',options:'healthy'},feed:{marketSession:'regular',feedState:'streaming',webSocketConnected:true,subscribedSymbols:['AAA'],alpacaWebSocketConnected:true,alpacaSubscribedSymbols:['SPY','QQQ'],lastTradeAt:now,lastTradeSymbol:'AAA',lastAlpacaStreamAt:now,lastAlpacaStreamSymbol:'SPY',lastAlpacaAt:now,lastAlpacaSymbol:'SPY',alpacaLiveFeed:'iex',overnightDataMode:'auto',overnightLiveAvailable:false},global:{tone:'CONSTRUCTIVE',confidence:82,summary:'CONSTRUCTIVE · 3 supportive / 0 headwind evidence points',drivers:{breadth:{key:'breadth',label:'Broad Market Breadth',state:'SUPPORTIVE',detail:'12 up / 3 down across broad universe',confidence:80},semiconductors:{key:'semiconductors',label:'Semiconductors',state:'SUPPORTIVE',detail:'SMH proxy',confidence:85},taiwan:{key:'taiwan',label:'Taiwan',state:'SUPPORTIVE',detail:'EWT proxy',confidence:80},rates_10y:{key:'rates_10y',label:'U.S. 10Y Yield',state:'NEUTRAL',detail:'4.1%',confidence:88}}},macroMetrics:{DGS10:{key:'DGS10',label:'U.S. 10Y',value:4.1,status:'OFFICIAL'}},macroEvents:[],eventMode:{active:false},eventReactions:[],options:{AAA:{symbol:'AAA',provider:'Alpaca Options',feed:'OPRA',state:'CURRENT',bias:'BULLISH',callVolume:1000,putVolume:500,putCallVolume:.5,averageIv:.35,expectedMove:8,updatedAt:now,provenance:'REAL OPTIONS SNAPSHOT'}},capabilities:[{capability:'Global Indices',source:'Proxy',mode:'AUTO',status:'ACTIVE'},{capability:'Options Intelligence',source:'Alpaca Options',mode:'AUTO',status:'ACTIVE'}],providerRouter:{updatedAt:now,routes:[{dataset:'VIX / Indices',primary:'Twelve Data',active:'yfinance',state:'FALLBACK',route:[{provider:'Twelve Data',circuit:'CLOSED'},{provider:'yfinance',circuit:'CLOSED'},{provider:'CBOE',circuit:'CLOSED'}]}]},freshness:[{dataset:'VIX',state:'DELAYED',provider:'yfinance',providerTimestamp:now-300000,ageMs:300000,expectedCadenceMs:120000,reason:'Valid delayed VIX source',fallback:'yfinance → CBOE',affected:['Market Regime','Dashboard'],session:'regular',action:'vix'},{dataset:'News',state:'FRESH',provider:'Finnhub',providerTimestamp:now,ageMs:0,expectedCadenceMs:600000,reason:'Within expected freshness threshold · provider timestamp unavailable; using DE.PULSE receipt after targeted provider reconciliation',fallback:'Marketaux',affected:['Dashboard','Research'],session:'regular',action:'news'}],freshnessSummary:{live:0,fresh:1,delayed:1,stale:0,error:0,unavailable:0,idle:0},signalValidation:{snapshots:[],message:'QA'}};
      runtime.researchPackage={symbol:'AAA',state:'PARTIAL',evidenceSnapshotId:'abcdef1234567890abcdef12',blockingReasons:['Catalyst & Material Event Context PARTIAL · material evidence exists without a matching Catalyst Watch reaction state'],components:[{dataset:'Quote',required:true,critical:true,state:'FRESH',source:'Alpaca IEX primary market data with fallback provenance',checkAgeMs:1200,dataAgeMs:1500,detail:'Selected-ticker quote is current and reconciled.'},{dataset:'Catalyst & Material Event Context',required:true,state:'PARTIAL',source:'Catalyst Watch · News/Earnings/SEC',checkAgeMs:2400,dataAgeMs:3000,detail:'Material event context: Earnings 2026-08-11 AMC · released · EPS 1.25 vs 1.20 · Material news · Very long acquisition and regulatory headline used to verify worst-case wrapping without clipping, overlap or false truncation · Finnhub. Required-evidence issue: material evidence exists without a matching Catalyst Watch reaction state.'},{dataset:'Required Market Context',required:true,state:'FRESH',source:'SPY=Alpaca · QQQ=Alpaca · VIX=CBOE delayed fallback · Global=direct+fallback',checkAgeMs:4000,dataAgeMs:5000,detail:'Required broad market context is current: SPY 500 (LIVE) · QQQ 450 (LIVE) · VIX 18 (DELAYED) · Global NEUTRAL / direct+fallback.'}]};
      runtime.corporateActionTruth={canonicalHistory:'Adjusted OHLCV remains canonical for technical calculations; raw comparison bars preserve explicit pagination and completeness provenance.',lifecycle:{AAA:{symbol:'AAA',state:'NAME OR TICKER CHANGE',source:'Corporate action provider with unusually long provider display name for worst-case responsive testing',evidenceId:'very-long-evidence-id-1234567890',effectiveDate:'2026-08-10',processDate:'2026-08-10',firstSeenAt:now-86400000,updatedAt:now,reason:'Ticker-change evidence reconciled across the persistent corporate-action ledger; this intentionally long explanation validates wrapping, card growth and neighboring-panel harmony without clipping or unintended horizontal scrolling.'}},symbolLineage:{OLDLONGSYMBOL:'AAA'},rawHistoryCoverage:{AAA:{symbol:'AAA',state:'PARTIAL',barCount:9999,pageCount:42,paginationComplete:false,firstBarAt:now-86400000*365,lastBarAt:now-86400000,detail:'Pagination safety stop retained as PARTIAL and retryable.'}},actions:[{id:'very-long-evidence-id-1234567890',symbol:'AAA',oldSymbol:'OLDLONGSYMBOL',newSymbol:'AAA',type:'name_change',status:'EFFECTIVE',processDate:'2026-08-10',source:'Corporate action provider'}]};
      runtime.evidenceSnapshot={id:'abcdef1234567890abcdef12',symbol:'AAA',generatedAt:now,researchState:'PARTIAL',symbolLifecycle:runtime.corporateActionTruth.lifecycle.AAA};
      window.DePulseLogic.setState(state,runtime); window.DePulseLogic.setPage('dashboard');
    }""")
    for w,h in viewports:
        page.set_viewport_size({'width':w,'height':h})
        for surface,sub in surfaces:
            page.evaluate("([surface])=>window.DePulseLogic.setPage(surface)",[surface])
            page.wait_for_timeout(25)
            if surface=='discovery' and sub:
                loc=page.locator(f'[data-discovery-mode="{sub}"]')
                if loc.count(): loc.first.click(force=True)
            if surface=='research' and sub:
                loc=page.locator(f'[data-research-tab="{sub}"]')
                if loc.count(): loc.first.click(force=True)
            page.wait_for_timeout(25)
            result=page.evaluate('''([surface,sub])=>{
              const de=document.documentElement,b=document.body,m=document.querySelector('#main');
              const visible=e=>{const r=e.getBoundingClientRect(),s=getComputedStyle(e);return s.display!=='none'&&s.visibility!=='hidden'&&r.width>0&&r.height>0};
              const badControls=[...document.querySelectorAll('button,select,input,textarea')].filter(visible).filter(e=>{const r=e.getBoundingClientRect();if(e.closest('.scanner-table'))return false;return r.right>innerWidth+2||r.left<-2});
              const selected=[...document.querySelectorAll('select')].filter(visible).filter(e=>{const r=e.getBoundingClientRect();return r.width<34||r.height<26});
              return {surface,sub,docOverflow:de.scrollWidth-de.clientWidth,bodyOverflow:b.scrollWidth-b.clientWidth,mainText:(m?.innerText||'').trim().length,badControls:badControls.slice(0,5).map(x=>({tag:x.tagName,id:x.id,cls:x.className,right:x.getBoundingClientRect().right,left:x.getBoundingClientRect().left})),tinySelects:selected.slice(0,5).map(x=>x.id||x.className)};
            }''',[surface,sub])
            checks+=1
            if result['docOverflow']>1 or result['bodyOverflow']>1 or result['mainText']<10 or result['badControls'] or result['tinySelects']:
                errors.append(f"{w}x{h} {surface}/{sub}: {json.dumps(result)}")

    # Header-stack regression: hidden notification keeps its reserved center lane inside the fixed-height header; ticker remains directly below topbar.
    page.set_viewport_size({'width':1920,'height':800})
    page.evaluate("window.DePulseLogic.setPage('dashboard')")
    page.wait_for_timeout(25)
    header_check=page.evaluate('''()=>{
      const h=document.querySelector('.topbar').getBoundingClientRect();
      const n=document.querySelector('#header-notification').getBoundingClientRect();
      const t=document.querySelector('#ticker-tape').getBoundingClientRect();
      return {notificationHeight:n.height,notificationInside:n.top>=h.top-1&&n.bottom<=h.bottom+1,gap:Math.round(t.top-h.bottom),topbarBottom:Math.round(h.bottom),tickerTop:Math.round(t.top)};
    }''')
    checks+=1
    if header_check['notificationHeight']<=0 or not header_check['notificationInside'] or abs(header_check['gap'])>1:
        errors.append('header notification-lane/ticker-stack regression: '+json.dumps(header_check))
    page.evaluate('window.scrollTo(0,900)')
    page.wait_for_timeout(25)
    scrolled_header=page.evaluate('''()=>{
      const h=document.querySelector('.topbar').getBoundingClientRect();
      const t=document.querySelector('#ticker-tape').getBoundingClientRect();
      return {gap:Math.round(t.top-h.bottom),topbarTop:Math.round(h.top),topbarBottom:Math.round(h.bottom),tickerTop:Math.round(t.top)};
    }''')
    checks+=1
    if abs(scrolled_header['gap'])>1 or abs(scrolled_header['topbarTop'])>1:
        errors.append('header scroll gap: '+json.dumps(scrolled_header))

    # Expansion-state regression: a user-opened panel must remain open after a live-style rerender.
    page.evaluate("window.scrollTo(0,0); window.DePulseLogic.setPage('dashboard')")
    page.wait_for_timeout(25)
    expander=page.locator('#main details[data-expand-key]').first
    if expander.count():
        expander.evaluate('(e)=>e.open=true')
        key=expander.get_attribute('data-expand-key')
        page.evaluate('window.DePulseLogic.render()')
        page.wait_for_timeout(25)
        still_open=page.locator(f'#main details[data-expand-key="{key}"]').first.evaluate('(e)=>e.open')
        checks+=1
        if not still_open: errors.append('expand panel state lost across rerender: '+str(key))
    else:
        checks+=1; errors.append('no persistent expandable panel found on dashboard')

    # v14.3.1 sidebar regression: restore the detailed previous-build Data Engine inside the scrollable rail; keep creator footer anchored.
    page.set_viewport_size({'width':1366,'height':650})
    page.evaluate("window.DePulseLogic.setPage('dashboard')")
    page.wait_for_timeout(25)
    sidebar_before=page.evaluate('''()=>{const s=document.querySelector('.sidebar').getBoundingClientRect(),f=document.querySelector('.creator-signature').getBoundingClientRect(),n=document.querySelector('.sidebar-scroll'),e=document.querySelector('.sidebar-data-engine'),h=document.querySelector('#health-list');return {footerDelta:Math.round(s.bottom-f.bottom),footerTop:Math.round(f.top),scrollHeight:n.scrollHeight,clientHeight:n.clientHeight,engineInScroll:n.contains(e),healthDisplay:h?getComputedStyle(h).display:'missing',healthRows:h?h.children.length:0,engineHeight:e?Math.round(e.getBoundingClientRect().height):0}}''')
    page.evaluate("()=>{const n=document.querySelector('.sidebar-scroll');n.scrollTop=Math.max(0,n.scrollHeight-n.clientHeight)}")
    page.wait_for_timeout(25)
    sidebar_after=page.evaluate('''()=>{const s=document.querySelector('.sidebar').getBoundingClientRect(),f=document.querySelector('.creator-signature').getBoundingClientRect();return {footerDelta:Math.round(s.bottom-f.bottom),footerTop:Math.round(f.top)}}''')
    checks+=1
    if abs(sidebar_before['footerDelta'])>2 or abs(sidebar_after['footerDelta'])>2 or abs(sidebar_before['footerTop']-sidebar_after['footerTop'])>2 or not sidebar_before['engineInScroll'] or sidebar_before['healthDisplay']=='none' or sidebar_before['healthRows']<8 or sidebar_before['engineHeight']<120:
        errors.append('detailed sidebar Data Engine/footer regression: '+json.dumps({'before':sidebar_before,'after':sidebar_after}))

    # v14.3.1 sidebar black-bar/jitter regression: no fixed pseudo-element scroll-fade may overlay Data Engine rows.
    sidebar_overlay=page.evaluate('''()=>{const side=document.querySelector('.sidebar'),sc=document.querySelector('.sidebar-scroll');if(!side||!sc)return null;const collect=(rules,out)=>{for(const r of rules||[]){if(r.selectorText&&(r.selectorText.includes('.sidebar::before')||r.selectorText.includes('.sidebar::after')))out.push(r.selectorText);if(r.cssRules)collect(r.cssRules,out)}};const selectors=[];for(const sheet of document.styleSheets){try{collect(sheet.cssRules,selectors)}catch(_){}}sc.scrollTop=Math.max(0,(sc.scrollHeight-sc.clientHeight)*.55);sc.dispatchEvent(new Event('scroll'));const before={selectors,canUp:side.classList.contains('sidebar-can-up'),canDown:side.classList.contains('sidebar-can-down'),scrollTop:sc.scrollTop};window.DePulseLogic.render();const side2=document.querySelector('.sidebar'),sc2=document.querySelector('.sidebar-scroll');return {before,after:{canUp:side2?.classList.contains('sidebar-can-up')||false,canDown:side2?.classList.contains('sidebar-can-down')||false},hasEngine:!!document.querySelector('.sidebar-data-engine'),hasHealth:!!document.querySelector('#health-list')}}''')
    checks+=1
    if not sidebar_overlay or sidebar_overlay['before']['selectors'] or sidebar_overlay['before']['canUp'] or sidebar_overlay['before']['canDown'] or sidebar_overlay['after']['canUp'] or sidebar_overlay['after']['canDown'] or not sidebar_overlay['hasEngine'] or not sidebar_overlay['hasHealth']:
        errors.append('v14.3.1 sidebar fixed black-overlay/jitter regression: '+json.dumps(sidebar_overlay))

    # v14.3.1 Dashboard priority/order regression: placement only, no collapse of More Market Context.
    page.set_viewport_size({'width':1920,'height':1080}); page.evaluate("window.DePulseLogic.setPage('dashboard')"); page.wait_for_timeout(25)
    order=page.evaluate('''()=>{const top=s=>document.querySelector(s)?.getBoundingClientRect().top??999999;const mi=document.querySelector('.market-intelligence-summary'),g=document.querySelector('.global-driver-panel');return {regime:top('.market-stack'),queue:top('.dashboard-decision-queue'),catalyst:top('.catalyst-pulse'),horizon:top('.summary-grid'),options:top('.options-intelligence'),mi:top('.market-intelligence-summary'),hasMI:!!mi,detailedGlobalOnDashboard:!!g}}''')
    checks+=1
    if not (order['regime']<order['queue']<order['catalyst']<order['horizon']<order['mi']) or (order['options']<999999 and not order['horizon']<order['options']<order['mi']) or not order['hasMI'] or order['detailedGlobalOnDashboard']: errors.append('dashboard priority/Market Intelligence ownership incorrect: '+json.dumps(order))

    # Trading-desk placement regression: regime -> KPIs -> candidates -> selected detail -> ticker news -> watchlist management.
    for desk in ['day','swing','long']:
        page.evaluate("([desk])=>window.DePulseLogic.setPage(desk)",[desk]); page.wait_for_timeout(35)
        desk_order=page.evaluate('''()=>{const top=s=>document.querySelector(s)?.getBoundingClientRect().top??999999;const detail=document.querySelector('.desk-detail-anchor'),ctx=detail?.querySelector('.context-lens'),chart=detail?.querySelector('.chart-wrap');return {regime:top('.desk-regime-control-grid'),kpis:top('.desk-kpis'),table:top('.desk-table'),detail:top('.desk-detail-anchor'),news:top('.ticker-news-panel'),watch:top('.desk-watchlist-management'),ctx:ctx?.getBoundingClientRect().top??999999,chart:chart?.getBoundingClientRect().top??999999}}''')
        checks+=1
        if not (desk_order['regime']<desk_order['kpis']<desk_order['table']<desk_order['detail']<desk_order['news']<desk_order['watch']) or not (desk_order['ctx']<desk_order['chart']):
            errors.append(f'{desk} placement hierarchy regression: '+json.dumps(desk_order))

    # Discovery placement regression: horizon/filter setup + Run Scan -> funnel -> candidates -> staged management.
    page.evaluate("window.DePulseLogic.setPage('discovery')"); page.wait_for_timeout(35)
    discovery_order=page.evaluate('''()=>{const top=s=>document.querySelector(s)?.getBoundingClientRect().top??999999;const cmd=document.querySelector('.scanner-command');const run=cmd?.querySelector('[data-run-scan]');const modes=cmd?.querySelector('.scanner-modes');const filters=cmd?.querySelector('.scanner-filter-grid');return {cmd:top('.scanner-command'),modes:modes?.getBoundingClientRect().top??999999,filters:filters?.getBoundingClientRect().top??999999,run:run?.getBoundingClientRect().top??999999,funnel:top('.discovery-funnel'),results:top('.scanner-results'),staged:top('#staged-candidates')}}''')
    checks+=1
    if not (discovery_order['cmd']<=discovery_order['modes']<discovery_order['filters']<=discovery_order['run']<discovery_order['funnel']<discovery_order['results']<discovery_order['staged']):
        errors.append('discovery placement hierarchy regression: '+json.dumps(discovery_order))

    discovery_rank=page.evaluate('''()=>{let c=document.querySelector('.discovery-rank-cell'),temp=false;if(!c){c=document.createElement('span');c.className='discovery-rank-cell';c.innerHTML='<b>98/100</b><small>Strong</small>';c.style.width='104px';document.querySelector('#main').appendChild(c);temp=true}const b=c.querySelector('b'),sm=c.querySelector('small'),r=c.getBoundingClientRect(),out={valueFont:parseFloat(getComputedStyle(b).fontSize),labelFont:parseFloat(getComputedStyle(sm).fontSize),overflow:c.scrollWidth-c.clientWidth,width:r.width};if(temp)c.remove();return out}''')
    checks+=1
    if not discovery_rank or discovery_rank['overflow']>1 or discovery_rank['valueFont']>12 or discovery_rank['labelFont']>9.5 or discovery_rank['valueFont']<=discovery_rank['labelFont']:
        errors.append('v14.3.1 Discovery Rank typography regression: '+json.dumps(discovery_rank))

    # Research v2 placement: target -> tabs -> decision brief -> what changed -> risk -> horizons -> sourced evidence -> AI confirmation.
    page.evaluate("window.DePulseLogic.setPage('research')"); page.wait_for_timeout(25)
    if page.locator('[data-research-tab="overview"]').count(): page.locator('[data-research-tab="overview"]').first.click(force=True); page.wait_for_timeout(25)
    research_order=page.evaluate("""()=>{const top=s=>document.querySelector(s)?.getBoundingClientRect().top??999999;return {target:top('.research-command-v2'),tabs:top('.research-tabs-v2'),brief:top('.research-decision-brief'),changed:top('.research-change-card'),risk:top('.research-risk-register'),horizons:top('.research-horizon-grid'),evidence:top('.research-evidence-heading'),ai:top('.research-confirmation')}}""")
    checks+=1
    if not (research_order['target']<research_order['tabs']<research_order['brief']<research_order['changed']<research_order['risk']<research_order['horizons']<research_order['evidence']<research_order['ai']):
        errors.append('research v2 placement hierarchy regression: '+json.dumps(research_order))

    # AI placement: provider -> target -> received context -> review selector -> result.
    page.evaluate("window.DePulseLogic.setPage('ai')"); page.wait_for_timeout(35)
    ai_order=page.evaluate('''()=>{const top=s=>document.querySelector(s)?.getBoundingClientRect().top??999999;return {provider:top('.ai-engine-card'),target:top('.ai-target-section'),context:top('.ai-context-received'),review:top('.ai-review-workspace'),result:top('.ai-result-card')}}''')
    checks+=1
    if not (ai_order['provider']<ai_order['target']<ai_order['context']<ai_order['review']) or (ai_order['result']<999999 and not ai_order['review']<ai_order['result']):
        errors.append('AI Copilot placement hierarchy regression: '+json.dumps(ai_order))

    # Maintenance placement: diagnostics -> freshness -> capabilities -> validation -> weekly -> performance -> engine -> storage -> build -> QA.
    page.evaluate("window.DePulseLogic.setPage('maintenance')"); page.wait_for_timeout(35)
    maint_order=page.evaluate('''()=>{const top=s=>document.querySelector(s)?.getBoundingClientRect().top??999999;const sections=[...document.querySelectorAll('.maintenance-page .settings-section')];const byTitle=t=>{const e=sections.find(x=>(x.querySelector('h2')?.textContent||'').includes(t));return e?e.getBoundingClientRect().top:999999};return {diag:top('.feed-diagnostics'),fresh:top('.data-freshness-priority'),cap:byTitle('Data Capabilities'),signal:byTitle('Signal Validation'),weekly:top('.maintenance-hero'),system:byTitle('System Health'),performance:top('.performance-observability'),engine:byTitle('Data Engine Detail'),storage:byTitle('Data & Storage'),build:byTitle('Application & Build Identity'),qa:byTitle('QA & Release History')}}''')
    checks+=1
    if not (maint_order['diag']<maint_order['fresh']<maint_order['cap']<maint_order['signal']<maint_order['weekly']<maint_order['system']<maint_order['performance']<maint_order['engine']<maint_order['storage']<maint_order['build']<maint_order['qa']):
        errors.append('maintenance placement hierarchy regression: '+json.dumps(maint_order))


    # v14.3.1 Maintenance corrections: one Data Freshness heading, visible VIX diagnostics, and normalized policy/capability typography.
    maint_polish=page.evaluate('''()=>{const exact=[...document.querySelectorAll('.maintenance-page h2,.maintenance-page h3')].filter(e=>e.textContent.trim().startsWith('Data Freshness')).length;const vx=document.querySelector('.vix-diagnostic-card');const pr=document.querySelector('.cache-policy-row');const cap=document.querySelector('.capability-card');const count=document.querySelector('.data-capabilities .count-pill');const css=e=>e?getComputedStyle(e):null;return {freshHeadings:exact,vix:!!vx,vixText:vx?.innerText||'',policy:pr?{label:parseFloat(css(pr.querySelector('b')).fontSize),value:parseFloat(css(pr.querySelector('span')).fontSize),helper:parseFloat(css(pr.querySelector('small')).fontSize),overflow:pr.scrollWidth-pr.clientWidth}:null,cap:cap?{title:parseFloat(css(cap.querySelector('div>b')).fontSize),source:parseFloat(css(cap.querySelector('strong')).fontSize),detail:cap.querySelector('p')?parseFloat(css(cap.querySelector('p')).fontSize):9,overflow:cap.scrollWidth-cap.clientWidth}:null,count:count?.textContent.trim()||''}}''')
    checks+=1
    if maint_polish['freshHeadings']!=1 or not maint_polish['vix'] or 'VIX' not in maint_polish['vixText'] or not maint_polish['policy'] or not maint_polish['cap'] or maint_polish['policy']['overflow']>1 or maint_polish['cap']['overflow']>1 or not (9<=maint_polish['policy']['label']<=11.5 and 10.5<=maint_polish['policy']['value']<=13 and 8<=maint_polish['policy']['helper']<=10) or not (10.5<=maint_polish['cap']['title']<=13 and 10<=maint_polish['cap']['source']<=12.5 and 8<=maint_polish['cap']['detail']<=10) or 'capabilities' not in maint_polish['count'].lower():
        errors.append('v14.3.1 maintenance polish regression: '+json.dumps(maint_polish))

    # Settings placement: configuration summary -> core data -> modes -> global/macro/options -> SEC -> engines -> signal -> AI -> application.
    page.evaluate("window.DePulseLogic.setPage('settings')"); page.wait_for_timeout(35)
    settings_order=page.evaluate('''()=>{const top=s=>document.querySelector(s)?.getBoundingClientRect().top??999999;const sections=[...document.querySelectorAll('.settings .settings-section')];const byTitle=t=>{const e=sections.find(x=>(x.querySelector('h2')?.textContent||'').includes(t));return e?e.getBoundingClientRect().top:999999};return {summary:top('.configuration-summary'),core:byTitle('Core Market Data'),modes:byTitle('Market Data Modes'),global:byTitle('Global, Macro & Options Data'),sec:byTitle('SEC / Filing Configuration'),engines:byTitle('Trading Engines'),signal:byTitle('Signal Engine'),ai:byTitle('AI Integration'),app:byTitle('Application')}}''')
    checks+=1
    if not (settings_order['summary']<settings_order['core']<settings_order['modes']<settings_order['global']<settings_order['sec']<settings_order['engines']<settings_order['signal']<settings_order['ai']<settings_order['app']):
        errors.append('settings placement hierarchy regression: '+json.dumps(settings_order))


    # v14.3.1 persistent Settings save action reserves its own row and never covers scrollable settings cards.
    page.set_viewport_size({'width':1512,'height':982}); page.evaluate("window.DePulseLogic.setPage('settings')"); page.wait_for_timeout(25)
    save_layout=page.evaluate('''()=>{const sc=document.querySelector('.settings'),bar=document.querySelector('.settings-shell>.save-bar');if(!sc||!bar)return null;sc.scrollTop=Math.min(1200,Math.max(0,sc.scrollHeight-sc.clientHeight));const sr=sc.getBoundingClientRect(),br=bar.getBoundingClientRect();const overlap=sr.left<br.right&&sr.right>br.left&&sr.top<br.bottom&&sr.bottom>br.top;return {scrollTop:sc.scrollTop,scrollHeight:sc.scrollHeight,clientHeight:sc.clientHeight,settingsBottom:sr.bottom,barTop:br.top,overlap}}''')
    checks+=1
    if not save_layout or save_layout['scrollHeight']<=save_layout['clientHeight'] or save_layout['overlap'] or save_layout['barTop']<save_layout['settingsBottom']-1:
        errors.append('v14.3.3 Settings save-bar overlap regression: '+json.dumps(save_layout))

    # v14.3.1 Save must restore the user's middle-page context after the Settings DOM rerenders.
    save_scroll=page.evaluate('''async()=>{const sc=document.querySelector('.settings-shell>.settings');if(!sc)return null;sc.scrollTop=Math.min(900,Math.max(120,sc.scrollHeight-sc.clientHeight-40));const before=sc.scrollTop,c=window.DePulseLogic.captureSaveContext();window.DePulseLogic.render();window.DePulseLogic.restoreSaveContext(c);await new Promise(r=>setTimeout(r,50));const next=document.querySelector('.settings-shell>.settings');return {before,after:next?.scrollTop||0}}''')
    checks+=1
    if not save_scroll or save_scroll['before']<100 or abs(save_scroll['after']-save_scroll['before'])>2:
        errors.append('v14.3.3 Save scroll preservation regression: '+json.dumps(save_scroll))

    # Options Intelligence visual regression: compact terminal row styling, not large flat gray bars.
    opts=page.evaluate('''()=>{const e=document.querySelector('.options-intel-row');if(!e)return null;const s=getComputedStyle(e),r=e.getBoundingClientRect();return {bg:s.backgroundColor,height:Math.round(r.height),radius:s.borderRadius}}''')
    checks+=1
    if opts and (opts['height']>64 or opts['bg'] in ('rgb(79, 84, 90)','rgb(80, 85, 91)')): errors.append('options intelligence styling regression: '+json.dumps(opts))

    # v14.3.3 notification lifecycle: feedback lives inside the existing global header center lane.
    page.set_viewport_size({'width':1512,'height':900}); page.evaluate("window.scrollTo(0,0); window.DePulseLogic.setPage('settings')"); page.wait_for_timeout(25)
    header_before=page.evaluate('''()=>{const h=document.querySelector('.topbar').getBoundingClientRect();return {top:h.top,bottom:h.bottom,height:h.height}}''')
    page.evaluate("toast('Settings Saved','','success')"); page.wait_for_timeout(40)
    short_notice=page.evaluate('''()=>{const el=document.querySelector('#header-notification'),cs=getComputedStyle(el);return {cls:el.className,justify:cs.justifyContent,align:cs.textAlign,font:parseFloat(cs.fontSize),overflow:el.scrollWidth-el.clientWidth}}''')
    checks+=1
    if 'short-message' not in short_notice['cls'] or short_notice['justify']!='center' or short_notice['align']!='center' or short_notice['overflow']>1:
        errors.append('v14.3.3 short header notification centering regression: '+json.dumps(short_notice))
    page.evaluate("toast('Market Open Prep Completed','A deliberately very long completion notification validates content-aware font fitting across the full center header lane before ellipsis, while preserving brand, clocks, session controls, Data Fallback, and Stop Data without changing header height.','success')"); page.wait_for_timeout(60)
    visible_notice=page.evaluate('''()=>{const h=document.querySelector('.topbar').getBoundingClientRect(),n=document.querySelector('#header-notification').getBoundingClientRect(),b=document.querySelector('.brand').getBoundingClientRect(),r=document.querySelector('.runtime').getBoundingClientRect(),el=document.querySelector('#header-notification'),cs=getComputedStyle(el);return {hidden:el.getAttribute('aria-hidden'),classes:el.className,font:parseFloat(cs.fontSize),headerHeight:h.height,top:n.top,bottom:n.bottom,left:n.left,right:n.right,width:n.width,height:n.height,inside:n.top>=h.top-1&&n.bottom<=h.bottom+1,noBrandOverlap:n.left>=b.right-1,noRuntimeOverlap:n.right<=r.left+1,overflow:cs.overflow,textOverflow:cs.textOverflow,whiteSpace:cs.whiteSpace}}''')
    checks+=1
    if visible_notice['hidden']!='false' or 'long-message' not in visible_notice['classes'] or visible_notice['font']>13.5 or visible_notice['font']<10.49 or visible_notice['width']<=0 or not visible_notice['inside'] or not visible_notice['noBrandOverlap'] or not visible_notice['noRuntimeOverlap'] or abs(visible_notice['headerHeight']-header_before['height'])>0.5 or visible_notice['overflow']!='hidden' or visible_notice['textOverflow']!='ellipsis' or visible_notice['whiteSpace']!='nowrap':
        errors.append('v14.3.3 header notification placement/overflow regression: '+json.dumps(visible_notice))
    page.evaluate("()=>{const sc=document.querySelector('.settings');if(sc)sc.scrollTop=1200;toast('Settings Saved','Scroll position must remain unchanged while header feedback updates.','success')}"); page.wait_for_timeout(30)
    mid_notice=page.evaluate('''()=>{const h=document.querySelector('.topbar').getBoundingClientRect(),n=document.querySelector('#header-notification').getBoundingClientRect(),sc=document.querySelector('.settings');return {scrollTop:sc?.scrollTop||0,headerHeight:h.height,inside:n.top>=h.top-1&&n.bottom<=h.bottom+1,visible:n.width>0&&n.height>0}}''')
    checks+=1
    if mid_notice['scrollTop']<100 or not mid_notice['visible'] or not mid_notice['inside'] or abs(mid_notice['headerHeight']-header_before['height'])>0.5:
        errors.append('v14.3.3 mid-scroll header notification regression: '+json.dumps(mid_notice))
    page.evaluate("clearTimeout(window.__headerAlertTimer); hideToast(document.querySelector('#header-notification'))"); page.wait_for_timeout(25)
    hidden_notice=page.evaluate('''()=>{const h=document.querySelector('.topbar').getBoundingClientRect(),n=document.querySelector('#header-notification').getBoundingClientRect(),el=document.querySelector('#header-notification'),cs=getComputedStyle(el);return {hidden:el.getAttribute('aria-hidden'),classes:el.className,headerHeight:h.height,height:n.height,visibility:cs.visibility,opacity:cs.opacity,position:cs.position}}''')
    checks+=1
    if hidden_notice['hidden']!='true' or 'header-notification' not in hidden_notice['classes'] or abs(hidden_notice['headerHeight']-header_before['height'])>0.5 or hidden_notice['visibility']!='hidden' or hidden_notice['position']=='fixed':
        errors.append('v14.3.3 hidden header notification stability regression: '+json.dumps(hidden_notice))

    # v14.3.3 control placement: sidebar keeps only three trading-readiness actions; generic operational actions live in Maintenance cards.
    engine_actions=page.evaluate('''()=>{const host=document.querySelector('#data-engine-manual-actions'),prep=[...document.querySelectorAll('#prep-status-list [data-engine-action]')],side=document.querySelector('.sidebar').getBoundingClientRect();const all=prep.map(x=>{const r=x.getBoundingClientRect();return {action:x.dataset.engineAction,w:r.width,h:r.height,left:r.left,right:r.right}});return {prep:prep.map(x=>x.dataset.engineAction),all,sideLeft:side.left,sideRight:side.right,host:!!host}}''')
    checks+=1
    required_prep={'pre-market-prep','market-open-prep','catalyst-evaluate'}
    if not engine_actions or engine_actions['host'] or set(engine_actions['prep'])!=required_prep or any(x['h']>30 or x['left']<engine_actions['sideLeft']-1 or x['right']>engine_actions['sideRight']+1 for x in engine_actions['all']):
        errors.append('v14.3.3 sidebar readiness-action placement regression: '+json.dumps(engine_actions))
    page.evaluate("window.DePulseLogic.setPage('maintenance')"); page.wait_for_timeout(25)
    maintenance_actions=page.evaluate('''()=>{const cards=[...document.querySelectorAll('.maintenance-health-card')];const actions=[...document.querySelectorAll('.maintenance-health-card [data-engine-action]')].map(x=>{const r=x.getBoundingClientRect(),card=x.closest('.maintenance-health-card').getBoundingClientRect();return {action:x.dataset.engineAction,w:r.width,h:r.height,inside:r.left>=card.left-1&&r.right<=card.right+1&&r.top>=card.top-1&&r.bottom<=card.bottom+1}});return {cards:cards.length,actions,actionNames:actions.map(x=>x.action)}}''')
    checks+=1
    expected_operational={'refresh-due','global-refresh','capability-recheck','integrity-check'}
    if not maintenance_actions or maintenance_actions['cards']<10 or not expected_operational.issubset(set(maintenance_actions['actionNames'])) or any(x['h']>30 or not x['inside'] for x in maintenance_actions['actions']):
        errors.append('v14.3.3 Maintenance operational-action placement regression: '+json.dumps(maintenance_actions))

    # v14.3.3 final Data Engine text rule: no clipping/ellipsis. Frequent rows reserve minimum height, but every value may expand.
    page.evaluate("window.DePulseLogic.setPage('dashboard')"); page.wait_for_timeout(25)
    side_text=page.evaluate('''()=>{
      const host=document.querySelector('#health-list'); if(!host)return null;
      const row=document.createElement('div'); row.dataset.healthKey='cache-refresh';
      row.innerHTML='<b>Cache Refresh</b><span class="health-value">Healthy · Pre-Market Prep · 0 refreshed · 7 current · provider evidence retained without hiding any operational detail</span>';
      host.appendChild(row);
      const v=row.querySelector('.health-value'), cs=getComputedStyle(v), rr=row.getBoundingClientRect();
      const prepRow=document.createElement('div'); prepRow.className='sidebar-prep-row'; prepRow.innerHTML='<div class="sidebar-prep-titleline"><b><span class="sidebar-prep-label">Earnings &amp; Material Catalyst Watch</span></b><button class="engine-mini-action">Evaluate</button></div><span class="sidebar-prep-state">READY</span><small>Last success · Not run yet</small><small>Activation · Event-driven</small>';
      (host.parentElement||host).appendChild(prepRow); const prep=prepRow.querySelector('.sidebar-prep-label'); const ps=getComputedStyle(prep);
      const result={overflow:cs.overflow,textOverflow:cs.textOverflow,lineClamp:cs.webkitLineClamp,scrollHeight:v.scrollHeight,clientHeight:v.clientHeight,rowHeight:rr.height,prepOverflow:ps.overflow,prepTextOverflow:ps.textOverflow,prepLineClamp:ps.webkitLineClamp,prepScrollHeight:prep.scrollHeight,prepClientHeight:prep.clientHeight,prepText:prep.textContent};
      row.remove(); prepRow.remove(); return result;
    }''')
    checks+=1
    if not side_text or side_text['overflow']=='hidden' or side_text['textOverflow']=='ellipsis' or side_text['lineClamp'] not in ('none','unset','') or side_text['scrollHeight']>side_text['clientHeight']+1 or side_text['rowHeight']<50 or side_text['prepOverflow']=='hidden' or side_text['prepTextOverflow']=='ellipsis' or side_text['prepLineClamp'] not in ('none','unset','') or side_text['prepScrollHeight']>side_text['prepClientHeight']+1 or 'Earnings & Material Catalyst Watch' not in side_text['prepText']:
        errors.append('v14.3.3 unclipped adaptive sidebar Data Engine text regression: '+json.dumps(side_text))

    page.evaluate("window.DePulseLogic.setPage('maintenance')"); page.wait_for_timeout(25)
    maintenance_text=page.evaluate('''()=>{
      const card=[...document.querySelectorAll('.maintenance-health-card')].find(x=>x.dataset.healthKey==='cache-refresh')||document.querySelector('.maintenance-health-card');
      if(!card)return null; const v=card.querySelector('.health-value'); if(!v)return null; const old=v.textContent;
      v.textContent='Healthy · Pre-Market Prep · 0 refreshed · 7 current · provider evidence retained without hiding any operational detail';
      const cs=getComputedStyle(v), cr=card.getBoundingClientRect(); const result={overflow:cs.overflow,textOverflow:cs.textOverflow,lineClamp:cs.webkitLineClamp,scrollHeight:v.scrollHeight,clientHeight:v.clientHeight,cardHeight:cr.height};
      v.textContent=old; return result;
    }''')
    checks+=1
    if not maintenance_text or maintenance_text['overflow']=='hidden' or maintenance_text['textOverflow']=='ellipsis' or maintenance_text['lineClamp'] not in ('none','unset','') or maintenance_text['scrollHeight']>maintenance_text['clientHeight']+1 or maintenance_text['cardHeight']<80:
        errors.append('v14.3.3 unclipped adaptive Maintenance Data Engine text regression: '+json.dumps(maintenance_text))

    # v14.3.7 screenshot-driven UI fixes.
    # Side readiness state must not dominate the title.
    page.evaluate("window.DePulseLogic.setPage('dashboard')"); page.wait_for_timeout(25)
    prep_type=page.evaluate('''()=>{const card=document.querySelector('.sidebar-prep-row');if(!card)return null;const t=card.querySelector('.sidebar-prep-label'),st=card.querySelector('.sidebar-prep-state');return {title:parseFloat(getComputedStyle(t).fontSize),state:parseFloat(getComputedStyle(st).fontSize),weight:getComputedStyle(st).fontWeight}}''')
    checks+=1
    if not prep_type or prep_type['state']>prep_type['title']+0.1 or prep_type['state']>10:
        errors.append('v14.3.7 readiness-state typography regression: '+json.dumps(prep_type))

    # v14.3.7: side readiness dot is smaller and sits after the action button, never before the title.
    prep_dot=page.evaluate('''()=>{const card=document.querySelector('.sidebar-prep-row');if(!card)return null;const line=card.querySelector('.sidebar-prep-actionline'),dot=line?.querySelector(':scope>.status-dot'),button=line?.querySelector('.engine-mini-action'),old=card.querySelector('.sidebar-prep-titleline b>.status-dot');if(!dot||!button)return {found:false,old:!!old};const d=dot.getBoundingClientRect(),b=button.getBoundingClientRect();return {found:true,old:!!old,w:d.width,h:d.height,after:d.left>=b.right+3,centerDelta:Math.abs((d.top+d.height/2)-(b.top+b.height/2)),gap:d.left-b.right}}''')
    checks+=1
    if not prep_dot or not prep_dot.get('found') or prep_dot.get('old') or prep_dot['w']>6.5 or prep_dot['h']>6.5 or not prep_dot['after'] or prep_dot['centerDelta']>2.5:
        errors.append('v14.3.7 side readiness-dot-after-action regression: '+json.dumps(prep_dot))

    # Decision Queue Why? evidence is expanded by default.
    why_state=page.evaluate('''()=>[...document.querySelectorAll('.decision-queue-card .queue-why')].map(x=>x.open)''')
    checks+=1
    if why_state and not all(why_state): errors.append('v14.3.7 Decision Queue Why default expansion regression: '+json.dumps(why_state))

    # Settings Auto-Start belongs in the page header, not the Application card.
    page.evaluate("window.DePulseLogic.setPage('settings')"); page.wait_for_timeout(25)
    auto_start=page.evaluate('''()=>{const c=document.querySelector('.settings-autostart-control'),head=document.querySelector('.settings .page-head'),copy=document.querySelector('.settings .page-head-copy'),app=[...document.querySelectorAll('.settings-section')].find(x=>(x.querySelector('h2')?.textContent||'').trim()==='Application');if(!c||!head||!copy)return null;const r=c.getBoundingClientRect(),h=head.getBoundingClientRect(),cp=copy.getBoundingClientRect(),pill=c.querySelector('.settings-autostart-pill'),state=c.querySelector('.autostart-state'),ps=pill?getComputedStyle(pill):null;return {inHead:head.contains(c),inApp:!!app?.contains(c),rightOfCopy:r.left>=cp.right+12,notFarRight:r.right<h.left+h.width*.75,balancedLeft:(r.left-h.left)/h.width,height:r.height,text:c.innerText,hasState:!!state,borderRadius:ps?parseFloat(ps.borderRadius):0,bg:ps?ps.backgroundImage:''}}''')
    checks+=1
    if not auto_start or not auto_start['inHead'] or auto_start['inApp'] or not auto_start['rightOfCopy'] or not auto_start['notFarRight'] or not (0.40<=auto_start['balancedLeft']<=0.62) or auto_start['height']>46 or 'Auto-Start' not in auto_start['text'] or not auto_start['hasState'] or auto_start['borderRadius']<12 or 'gradient' not in auto_start['bg']:
        errors.append('v14.3.7 Settings Auto-Start premium/balanced placement regression: '+json.dumps(auto_start))

    # Freshness & Cache Policy cards should be content-driven, not all forced tall.
    page.evaluate("window.DePulseLogic.setPage('maintenance')"); page.wait_for_timeout(25)
    cache_cards=page.evaluate('''()=>{const xs=[...document.querySelectorAll('.cache-policy-row')];return xs.slice(0,10).map(x=>({h:Math.round(x.getBoundingClientRect().height),min:getComputedStyle(x).minHeight,text:x.innerText.length}))}''')
    checks+=1
    if not cache_cards or any(x['min'] not in ('0px','auto') for x in cache_cards) or min(x['h'] for x in cache_cards)>100:
        errors.append('v14.3.7 Freshness card dead-space regression: '+json.dumps(cache_cards))

    # VIX manual refresh/retry is always present in Maintenance and action cluster is vertically centered.
    vix_action=page.evaluate('''()=>{const card=document.querySelector('.maintenance-health-card[data-health-key="vix"]');if(!card)return null;const b=card.querySelector('[data-engine-action="vix-refresh"]'),a=card.querySelector('.maintenance-card-action');if(!b||!a)return {button:false};const cr=card.getBoundingClientRect(),ar=a.getBoundingClientRect(),br=b.getBoundingClientRect();return {button:true,label:b.textContent.trim(),centerDelta:Math.abs((ar.top+ar.height/2)-(cr.top+cr.height/2)),inside:br.left>=cr.left&&br.right<=cr.right&&br.top>=cr.top&&br.bottom<=cr.bottom}}''')
    checks+=1
    if not vix_action or not vix_action.get('button') or vix_action['label'] not in ('Refresh','Retry') or vix_action['centerDelta']>3 or not vix_action['inside']:
        errors.append('v14.3.7 VIX/action-cluster regression: '+json.dumps(vix_action))

    # Maintenance manual controls must use the integrated mini-panel treatment, not a generic gray system button.
    action_style=page.evaluate('''()=>{const a=document.querySelector('.maintenance-card-action'),b=a?.querySelector('.engine-mini-action');if(!a||!b)return null;const ac=getComputedStyle(a),bc=getComputedStyle(b);return {panelRadius:parseFloat(ac.borderRadius),panelBg:ac.backgroundImage,panelBorder:ac.borderColor,buttonBg:bc.backgroundImage,buttonColor:bc.color,buttonRadius:parseFloat(bc.borderRadius)}}''')
    checks+=1
    if not action_style or action_style['panelRadius']<10 or 'gradient' not in action_style['panelBg'] or 'gradient' not in action_style['buttonBg'] or action_style['buttonRadius']<8:
        errors.append('v14.3.7 Maintenance integrated-action visual regression: '+json.dumps(action_style))

    # Long notifications must preserve their beginning; if clamped, the text span becomes left-aligned for end ellipsis.
    page.set_viewport_size({'width':1512,'height':900}); page.evaluate("window.DePulseLogic.setPage('settings'); toast('Global Data Degraded','Direct unavailable; real fallback active; plan limited or unavailable; proxy context retained while provider recovery is evaluated across all relevant market surfaces.','warning')"); page.wait_for_timeout(80)
    notice_text=page.evaluate('''()=>{const host=document.querySelector('#header-notification'),sp=host?.querySelector('.header-notification-text');if(!host||!sp)return null;const cs=getComputedStyle(sp);return {text:sp.textContent,clamped:host.classList.contains('notification-clamped'),align:cs.textAlign,font:parseFloat(cs.fontSize),overflow:cs.overflow,textOverflow:cs.textOverflow,whiteSpace:cs.whiteSpace,client:sp.clientWidth,scroll:sp.scrollWidth}}''')
    checks+=1
    if not notice_text or not notice_text['text'].startswith('Global Data Degraded') or notice_text['font']<10.49 or notice_text['font']>13.51 or notice_text['overflow']!='hidden' or notice_text['textOverflow']!='ellipsis' or notice_text['whiteSpace']!='nowrap' or (notice_text['clamped'] and notice_text['align']!='left'):
        errors.append('v14.3.7 header long-message start-preservation regression: '+json.dumps(notice_text))

    # v14.3.7 Market Regime: Risk state stays left; SPY, QQQ and Data Confidence share the second-row right cluster.
    page.set_viewport_size({'width':1920,'height':1080}); page.evaluate("window.DePulseLogic.setPage('swing')"); page.wait_for_timeout(35)
    regime_row=page.evaluate('''()=>{const row=document.querySelector('.desk-regime-state-row'),h=row?.querySelector('h3'),strip=row?.querySelector('.regime-live-strip'),pills=[...(strip?.querySelectorAll('.regime-quote-pill')||[])],dir=strip?.querySelector('.regime-direction-pill'),dc=strip?.querySelector('.data-confidence');if(!row||!h||!strip||pills.length!==2||!dir||!dc)return null;const r=row.getBoundingClientRect(),hr=h.getBoundingClientRect(),sr=strip.getBoundingClientRect(),pr=pills.map(x=>({text:x.innerText,rect:x.getBoundingClientRect()})),dr=dir.getBoundingClientRect(),cr=dc.getBoundingClientRect();return {row:[r.left,r.top,r.right,r.bottom],risk:h.innerText,strip:[sr.left,sr.top,sr.right,sr.bottom],pills:pr.map(x=>x.text),direction:dir.innerText,ordered:pr[0].rect.left<pr[1].rect.left&&pr[1].rect.right<=dr.left+1&&dr.right<=cr.left+1,sameRow:Math.abs((hr.top+hr.height/2)-(sr.top+sr.height/2))<10,inside:sr.right<=r.right+1&&hr.left>=r.left-1,overflow:row.scrollWidth-row.clientWidth}}''')
    checks+=1
    valid_direction={'STRONG BULLISH','BULLISH','LEAN BULLISH','NEUTRAL','LEAN BEARISH','BEARISH','STRONG BEARISH'}
    if not regime_row or 'Risk' not in regime_row['risk'] or regime_row['pills'][0].split('\n')[0] != 'SPY' or regime_row['pills'][1].split('\n')[0] != 'QQQ' or regime_row['direction'] not in valid_direction or not regime_row['ordered'] or not regime_row['sameRow'] or not regime_row['inside'] or regime_row['overflow']>1:
        errors.append('v14.3.7 Market Regime SPY/QQQ/state/Data Confidence placement regression: '+json.dumps(regime_row))

    # The four headline Decision Queue values must remain smaller than their cards and never wrap/overflow.
    page.evaluate("window.DePulseLogic.setPage('dashboard')"); page.wait_for_timeout(25)
    queue_values=page.evaluate('''()=>{const card=document.querySelector('.decision-queue-card');if(!card)return [];return [...card.querySelectorAll('.decision-queue-fields>span')].slice(0,4).map(x=>{const b=x.querySelector('b'),r=x.getBoundingClientRect(),br=b?.getBoundingClientRect(),cs=b?getComputedStyle(b):null;return {text:b?.innerText||'',font:cs?parseFloat(cs.fontSize):99,nowrap:cs?.whiteSpace,inside:!!(br&&br.left>=r.left-1&&br.right<=r.right+1&&br.top>=r.top-1&&br.bottom<=r.bottom+1),overflow:b?b.scrollWidth-b.clientWidth:99}})}''')
    checks+=1
    if len(queue_values)!=4 or any(x['font']>7.1 or x['nowrap']!='nowrap' or not x['inside'] or x['overflow']>1 for x in queue_values):
        errors.append('v14.3.7 Decision Queue four-value fit regression: '+json.dumps(queue_values))

    # v14.3.7 font-fit audit: compact Decision Queue labels and values must remain fully visible at 1366px.
    page.set_viewport_size({'width':1366,'height':650}); page.evaluate("window.DePulseLogic.setPage('dashboard')"); page.wait_for_timeout(30)
    queue_fit=page.evaluate('''()=>[...document.querySelectorAll('.decision-queue-fields>span:nth-child(-n+4)')].map(x=>{const a=x.querySelector('small'),b=x.querySelector('b');return {label:a?.textContent.trim()||'',value:b?.textContent.trim()||'',labelOverflow:a?a.scrollWidth-a.clientWidth:99,valueOverflow:b?b.scrollWidth-b.clientWidth:99,labelFont:a?parseFloat(getComputedStyle(a).fontSize):99,valueFont:b?parseFloat(getComputedStyle(b).fontSize):99}})''')
    checks+=1
    if not queue_fit or any(x['labelOverflow']>1 or x['valueOverflow']>1 or x['labelFont']>5.3 or x['valueFont']>5.7 for x in queue_fit):
        errors.append('v14.3.7 compact Decision Queue font-fit regression: '+json.dumps(queue_fit[:8]))

    # v14.3.7 font-fit audit: Event & Data Context heading must fit its compact card at desktop widths.
    for fit_surface in ['day','swing']:
        page.set_viewport_size({'width':1280,'height':720}); page.evaluate("([s])=>window.DePulseLogic.setPage(s)",[fit_surface]); page.wait_for_timeout(30)
        event_fit=page.evaluate('''()=>{const h=document.querySelector('.event-data-context-head h3');if(!h)return null;return {text:h.textContent.trim(),overflow:h.scrollWidth-h.clientWidth,font:parseFloat(getComputedStyle(h).fontSize),whiteSpace:getComputedStyle(h).whiteSpace}}''')
        checks+=1
        if not event_fit or event_fit['overflow']>1 or event_fit['font']>11.1:
            errors.append(f'v14.3.7 {fit_surface} Event & Data Context font-fit regression: '+json.dumps(event_fit))

    # v14.3.7 font-fit audit: Discovery signed percentages must fit the fixed numeric cells without clipping.
    page.set_viewport_size({'width':1512,'height':982}); page.evaluate("window.DePulseLogic.setPage('discovery')"); page.wait_for_timeout(35)
    discovery_numeric_fit=page.evaluate('''()=>[...document.querySelectorAll('.scanner-row-day .scanner-number-cell')].map(x=>({text:x.textContent.trim(),overflow:x.scrollWidth-x.clientWidth,font:parseFloat(getComputedStyle(x).fontSize)}))''')
    checks+=1
    if not discovery_numeric_fit or any(x['overflow']>1 or x['font']>11.6 for x in discovery_numeric_fit):
        errors.append('v14.3.7 Discovery numeric-cell font-fit regression: '+json.dumps(discovery_numeric_fit))

    # v14.3.7 desk watchlists use one canonical DAY / SWING / LONG membership strip.
    page.evaluate("window.DePulseLogic.setPage('swing')"); page.wait_for_timeout(25)
    membership=page.evaluate('''()=>{const row=[...document.querySelectorAll('.watchlist-list-row')].find(x=>x.querySelector('.watchlist-row-symbol b')?.textContent.trim()==='AAA');if(!row)return null;const pills=[...row.querySelectorAll('.desk-membership-pill')],active=pills.filter(x=>x.classList.contains('active'));const m=row.querySelector('.watchlist-row-membership')?.getBoundingClientRect(),setup=row.querySelector('.watchlist-row-setup')?.getBoundingClientRect();return {labels:pills.map(x=>x.textContent.trim()),active:active.map(x=>x.textContent.trim()),overlap:!!(m&&setup&&m.right>setup.left+1),header:[...document.querySelectorAll('.watchlist-list-head span')].map(x=>x.textContent.trim())}}''')
    checks+=1
    if not membership or membership['labels']!=['DAY','SWING','LONG'] or membership['active']!=['DAY','SWING','LONG'] or membership['overlap'] or 'Desks' not in membership['header']:
        errors.append('v14.3.7 canonical desk-membership regression: '+json.dumps(membership))

    # v14.3.7 Options Intelligence settings card is content-driven, without screenshot-style dead space.
    page.evaluate("window.DePulseLogic.setPage('settings')"); page.wait_for_timeout(35)
    options_layout=page.evaluate('''()=>{const card=document.querySelector('.compact-options-provider'),body=card?.querySelector('.options-provider-body'),btn=card?.querySelector('.options-capability-test'),status=card?.querySelector('.status-box');if(!card||!body||!btn||!status)return null;const cr=card.getBoundingClientRect(),sr=status.getBoundingClientRect();return {height:cr.height,bottomGap:cr.bottom-sr.bottom,statusText:status.innerText,overflow:card.scrollWidth-card.clientWidth,minHeight:getComputedStyle(card).minHeight}}''')
    checks+=1
    if not options_layout or options_layout['height']>255 or options_layout['bottomGap']>28 or options_layout['overflow']>1 or options_layout['minHeight'] not in ('0px','auto') or 'Connected' not in options_layout['statusText']:
        errors.append('v14.3.7 Options settings dead-space regression: '+json.dumps(options_layout))

    # v14.3.7 Discovery screenshot regression: explicit cells may scroll as a table, but must never collide internally.
    page.set_viewport_size({'width':1920,'height':1080}); page.evaluate("window.DePulseLogic.setPage('discovery')"); page.wait_for_timeout(40)
    discovery_layout=page.evaluate('''()=>{const row=document.querySelector('.scanner-row-day'),table=document.querySelector('.scanner-table');if(!row||!table)return null;const cells=[...row.children],rects=cells.map(x=>x.getBoundingClientRect()),pairs=[];for(let i=0;i<rects.length-1;i++)pairs.push(rects[i].right>rects[i+1].left+1);const liq=row.querySelector('.scanner-liquidity-cell'),lv=liq?.querySelector('b'),ls=liq?.querySelector('small'),action=row.querySelector('.scanner-action-cell'),buttons=[...(action?.querySelectorAll('button')||[])];const lr=liq?.getBoundingClientRect(),vr=lv?.getBoundingClientRect(),sr=ls?.getBoundingClientRect(),ar=action?.getBoundingClientRect();return {cells:cells.length,cellOverlap:pairs.some(Boolean),liquidity:liq?.innerText||'',liqSeparated:!!(vr&&sr&&vr.bottom<=sr.top+1),liqInside:!!(lr&&vr&&sr&&vr.left>=lr.left-1&&vr.right<=lr.right+1&&sr.left>=lr.left-1&&sr.right<=lr.right+1),buttons:buttons.length,actionsInside:!!(ar&&buttons.every(b=>{const r=b.getBoundingClientRect();return r.left>=ar.left-1&&r.right<=ar.right+1})),tableOverflow:table.scrollWidth-table.clientWidth,documentOverflow:document.documentElement.scrollWidth-document.documentElement.clientWidth,desks:row.querySelector('.scanner-desk-cell')?.innerText||''}}''')
    checks+=1
    if not discovery_layout or discovery_layout['cells']!=9 or discovery_layout['cellOverlap'] or not discovery_layout['liqSeparated'] or not discovery_layout['liqInside'] or discovery_layout['buttons']!=3 or not discovery_layout['actionsInside'] or discovery_layout['documentOverflow']>1 or '$6.4M' not in discovery_layout['liquidity'] or '0.15% spread' not in discovery_layout['liquidity'] or 'DAY' not in discovery_layout['desks']:
        errors.append('v14.3.7 Discovery collision/audit regression: '+json.dumps(discovery_layout))

    # At a narrower supported desktop width the Discovery table may scroll inside its card, never the document.
    page.set_viewport_size({'width':1180,'height':820}); page.wait_for_timeout(25)
    discovery_narrow=page.evaluate('''()=>{const table=document.querySelector('.scanner-table'),row=document.querySelector('.scanner-row-day');return table&&row?{innerScrollable:table.scrollWidth>table.clientWidth,docOverflow:document.documentElement.scrollWidth-document.documentElement.clientWidth,rowWidth:row.getBoundingClientRect().width,tableClient:table.clientWidth}:null}''')
    checks+=1
    if not discovery_narrow or not discovery_narrow['innerScrollable'] or discovery_narrow['docOverflow']>1:
        errors.append('v14.3.7 Discovery narrow-width containment regression: '+json.dumps(discovery_narrow))

    # v14.3.7 SEC Form 4 UI must state BUY/SELL/OTHER semantics and expose transaction detail, not generic Insider Transaction.
    page.set_viewport_size({'width':1512,'height':900}); page.evaluate("window.DePulseLogic.setPage('filings')"); page.wait_for_timeout(35)
    sec_f4=page.evaluate('''()=>{const main=document.querySelector('#main'),text=main?.innerText||'';const detail=[...document.querySelectorAll('.sec-transaction-detail')].find(x=>x.innerText.includes('Open-market/private purchase'));return {hasBuy:text.includes('Insider BUY'),generic:text.includes('Insider Transaction'),detail:detail?.innerText||''}}''')
    checks+=1
    if not sec_f4 or not sec_f4['hasBuy'] or sec_f4['generic'] or 'Jane Example' not in sec_f4['detail'] or 'CFO' not in sec_f4['detail'] or '12,000 owned after' not in sec_f4['detail']:
        errors.append('v14.3.7 SEC Form 4 classification UI regression: '+json.dumps(sec_f4))

    # Global Market Drivers aggregate state: one indicator/dot based on the canonical combined driver tone; no duplicated state word in the evidence line.
    page.set_viewport_size({'width':1512,'height':900}); page.evaluate("window.DePulseLogic.setPage('market-intelligence')"); page.wait_for_timeout(30)
    global_state=page.evaluate('''()=>{const panel=document.querySelector('.global-driver-panel'),state=panel?.querySelector('.global-context-state'),dot=state?.querySelector('.status-dot'),p=panel?.querySelector('.global-driver-heading p');return panel&&state&&dot&&p?{state:state.innerText.trim(),dot:dot.className,summary:p.innerText.trim(),duplicate:p.innerText.trim().toUpperCase().startsWith(state.innerText.trim().toUpperCase()+' ·')}:null}''')
    checks+=1
    if not global_state or global_state['state']!='CONSTRUCTIVE' or 'green' not in global_state['dot'] or global_state['duplicate'] or not global_state['summary'].startswith('3 supportive / 0 headwind evidence points'):
        errors.append('v14.3.7 Global Market Drivers aggregate-state indicator regression: '+json.dumps(global_state))

    # v14.3.3 baseline-preservation: no Data Engine disclosure row; v14.2 empty-space cleanup remains intact.
    page.evaluate("window.DePulseLogic.setPage('dashboard')"); page.wait_for_timeout(25)
    cleanup=page.evaluate('''()=>{const d=document.querySelector('.sidebar-engine-details');const e=document.createElement('div');e.className='empty';e.textContent='No data';document.body.appendChild(e);const cs=getComputedStyle(e),r=e.getBoundingClientRect();e.remove();const sc=document.querySelector('.summary-card');return {detailsPresent:!!d,emptyHeight:r.height,emptyPadTop:parseFloat(cs.paddingTop),summaryMin:sc?getComputedStyle(sc).minHeight:null}}''')
    checks+=1
    if cleanup['detailsPresent'] or cleanup['emptyHeight']>48 or cleanup['emptyPadTop']>14.1 or (cleanup['summaryMin'] not in (None,'0px','auto')):
        errors.append('v14.3.3 baseline density/Data Engine disclosure regression: '+json.dumps(cleanup))

    # Build identity must reflect the packaged v14.3.3 backend/renderer pair.
    page.evaluate("window.DePulseLogic.setPage('maintenance')"); page.wait_for_timeout(25)
    build_text=page.locator('#main').inner_text()
    checks+=1
    if EXPECTED_RELEASE_VERSION not in build_text or EXPECTED_BUILD_ID not in build_text or 'BUILD IDENTITY VERIFIED' not in build_text:
        errors.append('build identity display stale or mismatched')

    # v14.2 Global/Macro driver layout: title/status/detail must not collide.
    page.set_viewport_size({'width':1920,'height':1080}); page.evaluate("window.DePulseLogic.setPage('market-intelligence')"); page.wait_for_timeout(25)
    gpanel=page.locator('.global-driver-panel').first
    if gpanel.count():
        gpanel.evaluate('(e)=>e.open=true'); page.wait_for_timeout(25)
        global_layout=page.evaluate('''()=>{const overlap=(a,b)=>a&&b&&a.left<b.right&&a.right>b.left&&a.top<b.bottom&&a.bottom>b.top;const rows=[...document.querySelectorAll('.global-driver-row')];return rows.map((r,i)=>{const title=r.querySelector('div>b')?.getBoundingClientRect(),status=r.querySelector(':scope>span')?.getBoundingClientRect(),detail=r.querySelector(':scope>small')?.getBoundingClientRect(),box=r.getBoundingClientRect();return {i,titleStatus:overlap(title,status),titleDetail:overlap(title,detail),statusDetail:overlap(status,detail),overflow:r.scrollWidth-r.clientWidth,height:Math.round(box.height)}})}''')
        checks+=1
        bad=[x for x in global_layout if x['titleStatus'] or x['titleDetail'] or x['statusDetail'] or x['overflow']>1 or x['height']<70]
        if bad: errors.append('global/macro driver card collision: '+json.dumps(bad[:8]))
    else:
        checks+=1; errors.append('global/macro driver panel missing')

    # v14.3.1 More Market Context spacing/hierarchy: hero cards should be prominent but not hollow, and major groups need clear breathing room.
    context_spacing=page.evaluate('''()=>{const hero=[...document.querySelectorAll('.market-direction-group .market-quote')].map(x=>Math.round(x.getBoundingClientRect().height));const groups=[...document.querySelectorAll('.global-context-body>.market-context-group')];const gaps=[];for(let i=1;i<groups.length;i++){const a=groups[i-1].getBoundingClientRect(),b=groups[i].getBoundingClientRect();gaps.push(Math.round(b.top-a.bottom))}const body=document.querySelector('.always-visible-market-context .global-context-body');return {hero,gaps,bodyGap:body?parseFloat(getComputedStyle(body).rowGap||getComputedStyle(body).gap):0,groups:groups.length}}''')
    checks+=1
    if not context_spacing or context_spacing['groups']<5 or any(h>180 for h in context_spacing['hero']) or context_spacing['bodyGap']<20 or any(g<18 for g in context_spacing['gaps']):
        errors.append('v14.3.1 More Market Context spacing/panel-balance regression: '+json.dumps(context_spacing))

    # Market Instruments hover must not change geometry or pause the ticker animation.
    page.evaluate("window.scrollTo(0,0); window.DePulseLogic.setPage('dashboard')"); page.wait_for_timeout(25)
    ticker_info=page.evaluate('''()=>{const vp=document.querySelector('.ticker-viewport');if(!vp)return null;const vr=vp.getBoundingClientRect(),items=[...vp.querySelectorAll('.ticker-item')];let idx=items.findIndex(x=>{const r=x.getBoundingClientRect();return r.right>vr.left+8&&r.left<vr.right-8});if(idx<0)idx=0;const e=items[idx];if(!e)return null;const r=e.getBoundingClientRect();return {idx,x:Math.max(vr.left+4,Math.min(vr.right-4,r.left+r.width/2)),y:r.top+r.height/2,w:r.width,h:r.height,p:getComputedStyle(e).padding,b:getComputedStyle(e).borderWidth}}''')
    if ticker_info:
        page.mouse.move(ticker_info['x'],ticker_info['y']); page.wait_for_timeout(80)
        after=page.evaluate('''(idx)=>{const vp=document.querySelector('.ticker-viewport'),items=[...vp?.querySelectorAll('.ticker-item')||[]],e=items[idx];if(!vp||!e)return null;const r=e.getBoundingClientRect();return {w:r.width,h:r.height,p:getComputedStyle(e).padding,b:getComputedStyle(e).borderWidth,play:getComputedStyle(vp.querySelector('.ticker-track')).animationPlayState}}''',ticker_info['idx'])
        checks+=1
        if not after or abs(ticker_info['w']-after['w'])>.5 or abs(ticker_info['h']-after['h'])>.5 or ticker_info['p']!=after['p'] or ticker_info['b']!=after['b'] or after['play']!='running':
            errors.append('market instruments hover jiggle regression: '+json.dumps({'before':ticker_info,'after':after}))
    else:
        checks+=1; errors.append('market instruments ticker item missing')

    # Event & Data Context must keep a clean header/actions and non-overlapping key/value rows.
    page.set_viewport_size({'width':1180,'height':820}); page.evaluate("window.DePulseLogic.setPage('day')"); page.wait_for_timeout(40)
    event_layout=page.evaluate('''()=>{const card=document.querySelector('.event-data-context-card');if(!card)return null;const overlap=(a,b)=>a&&b&&a.left<b.right&&a.right>b.left&&a.top<b.bottom&&a.bottom>b.top;const head=card.querySelector('.event-data-context-head'),title=head?.querySelector('h3')?.getBoundingClientRect(),actions=head?.querySelector('.event-context-actions')?.getBoundingClientRect();const rows=[...card.querySelectorAll('.event-context-list .compact-row')].map((r,i)=>{const a=r.querySelector('b')?.getBoundingClientRect(),b=r.querySelector('span')?.getBoundingClientRect();return {i,overlap:overlap(a,b),overflow:r.scrollWidth-r.clientWidth}});return {headOverlap:overlap(title,actions),overflow:card.scrollWidth-card.clientWidth,rows}}''')
    checks+=1
    if not event_layout or event_layout['headOverlap'] or event_layout['overflow']>1 or any(x['overlap'] or x['overflow']>1 for x in event_layout['rows']):
        errors.append('event/data context layout regression: '+json.dumps(event_layout))

    # SEC Intelligence Summary: values must not visually overpower labels or collide.
    page.set_viewport_size({'width':1180,'height':820}); page.evaluate("window.DePulseLogic.setPage('long')"); page.wait_for_timeout(40)
    sec_layout=page.evaluate('''()=>{const card=document.querySelector('.sec-intelligence-card');if(!card)return null;const overlap=(a,b)=>a&&b&&a.left<b.right&&a.right>b.left&&a.top<b.bottom&&a.bottom>b.top;return [...card.querySelectorAll('.sec-intel-row')].map((r,i)=>{const l=r.querySelector('span'),v=r.querySelector('b'),a=l?.getBoundingClientRect(),b=v?.getBoundingClientRect();return {i,overlap:overlap(a,b),labelFont:parseFloat(getComputedStyle(l).fontSize),valueFont:parseFloat(getComputedStyle(v).fontSize),overflow:r.scrollWidth-r.clientWidth}})}''')
    checks+=1
    if sec_layout is None or any(x['overlap'] or x['overflow']>1 or x['valueFont']>x['labelFont']+.01 for x in sec_layout):
        errors.append('SEC intelligence typography/layout regression: '+json.dumps(sec_layout))

    # v14.3.7 font-fit audit: narrow desk-watchlist change values must fit their existing cells without touching membership pills.
    page.set_viewport_size({'width':820,'height':1050}); page.evaluate("window.DePulseLogic.setPage('day')"); page.wait_for_timeout(35)
    watchlist_fit=page.evaluate('''()=>[...document.querySelectorAll('.watchlist-row-change')].map(x=>{const b=x.querySelector('b'),r=x.getBoundingClientRect(),br=b?.getBoundingClientRect();return {text:b?.textContent.trim()||'',overflow:b?b.scrollWidth-b.clientWidth:99,font:b?parseFloat(getComputedStyle(b).fontSize):99,inside:!!(br&&br.left>=r.left-1&&br.right<=r.right+1)}})''')
    checks+=1
    if not watchlist_fit or any(x['overflow']>1 or x['font']>6.1 or not x['inside'] for x in watchlist_fit):
        errors.append('v14.3.7 narrow watchlist change-value font-fit regression: '+json.dumps(watchlist_fit))

    # v15 Master Market Symbols is priority-placed, compact, and not appended below lower-priority global context.
    page.set_viewport_size({'width':1512,'height':982}); page.evaluate("window.DePulseLogic.setPage('dashboard')"); page.wait_for_timeout(35)
    master=page.evaluate('''()=>{const m=document.querySelector('.master-market-symbols'),q=document.querySelector('.dashboard-decision-queue'),c=document.querySelector('.dashboard-catalyst-priority'),g=document.querySelector('.market-intelligence-summary');if(!m)return null;const r=m.getBoundingClientRect();return {chips:[...m.querySelectorAll('.master-symbol-chip b')].map(x=>x.textContent.trim()),top:r.top,afterQueue:q?q.getBoundingClientRect().bottom<=r.top+2:false,beforeCatalyst:c?r.bottom<=c.getBoundingClientRect().top+2:false,beforeGlobal:g?r.bottom<g.getBoundingClientRect().top:true,overflow:m.scrollWidth-m.clientWidth}}''')
    checks+=1
    if not master or 'AAA' not in master['chips'] or not master['afterQueue'] or not master['beforeCatalyst'] or not master['beforeGlobal'] or master['overflow']>1:
        errors.append('v15 Master Market Symbols placement/fit regression: '+json.dumps(master))

    # v15 Global Market Driver shows an actual-refresh Last Updated line on its canonical Market Intelligence surface.
    page.evaluate("window.DePulseLogic.setPage('market-intelligence')"); page.wait_for_timeout(25)
    global_updated=page.evaluate('''()=>{const e=document.querySelector('.global-driver-panel .global-last-updated');return e?{text:e.textContent.trim(),overflow:e.scrollWidth-e.clientWidth}:null}''')
    checks+=1
    if not global_updated or 'Last updated' not in global_updated['text'] or global_updated['overflow']>1:
        errors.append('v15 Global Market Driver Last Updated regression: '+json.dumps(global_updated))

    # v15 table header / row column alignment is validated together with Action / Score.
    # v15 Action / Score header and cells are center-aligned.
    page.evaluate("window.DePulseLogic.setPage('day')"); page.wait_for_timeout(30)
    action_center=page.evaluate('''()=>{const h=document.querySelector('.desk-head>div:nth-child(2)'),c=document.querySelector('.desk-row .cell-action');if(!h||!c)return null;return {header:getComputedStyle(h).textAlign,cell:getComputedStyle(c).textAlign,align:getComputedStyle(c).alignItems}}''')
    checks+=1
    if not action_center or action_center['header']!='center' or action_center['cell']!='center' or action_center['align']!='center':
        errors.append('v15 Action / Score centering regression: '+json.dumps(action_center))

    # v15 row-hover stability: Day/Swing/Long/Discovery rows keep exact geometry and scroll position.
    for surf,selector in [('day','.desk-row:not(.desk-add-row)'),('swing','.desk-row:not(.desk-add-row)'),('long','.desk-row:not(.desk-add-row)'),('discovery','.scanner-row')]:
        page.set_viewport_size({'width':1512,'height':982}); page.evaluate("([s])=>window.DePulseLogic.setPage(s)",[surf]); page.wait_for_timeout(30)
        info=page.evaluate('''([sel])=>{const e=document.querySelector(sel);if(!e)return null;const r=e.getBoundingClientRect(),m=document.querySelector('#main');return {x:r.left+Math.min(30,r.width/2),y:r.top+r.height/2,w:r.width,h:r.height,left:r.left,top:r.top,mainY:m?.scrollTop||0,winY:scrollY}}''',[selector])
        if info:
            page.mouse.move(info['x'],info['y']); page.wait_for_timeout(80)
            after=page.evaluate('''([sel])=>{const e=document.querySelector(sel);if(!e)return null;const r=e.getBoundingClientRect(),m=document.querySelector('#main');return {w:r.width,h:r.height,left:r.left,top:r.top,mainY:m?.scrollTop||0,winY:scrollY}}''',[selector])
        else: after=None
        checks+=1
        if not info or not after or any(abs(info[k]-after[k])>.5 for k in ['w','h','left','top']) or info['mainY']!=after['mainY'] or info['winY']!=after['winY']:
            errors.append(f'v15 {surf} row hover jitter regression: '+json.dumps({'before':info,'after':after}))

    # v15 Maintenance shows diagnostic freshness rows + Provider Router, including delayed VIX truth.
    page.evaluate("window.DePulseLogic.setPage('maintenance')"); page.wait_for_timeout(30)
    freshness=page.evaluate('''()=>({rows:[...document.querySelectorAll('.freshness-diagnostic-row')].map(x=>x.innerText),router:document.querySelector('.provider-router-panel')?.innerText||''})''')
    checks+=1
    if not freshness['rows'] or not any('VIX' in x and 'DELAYED' in x and 'yfinance' in x for x in freshness['rows']) or 'VIX / Indices' not in freshness['router']:
        errors.append('v15 Data Freshness / Provider Router rendering regression: '+json.dumps(freshness))


    # v15.1.2 edge hardening: long Reason/Fallback explanations must remain inside their grid cell
    # and may never intrude into Check/Data Age or action columns.
    long_freshness=page.evaluate('''()=>{const row=document.querySelector('.v151-freshness-row');if(!row)return null;const reason=row.querySelector('.freshness-reason'),check=row.querySelector('.freshness-check-age'),data=row.querySelector('.freshness-data-age'),action=row.querySelector('.freshness-action');if(!reason||!check||!data)return null;const oldB=reason.querySelector('b')?.textContent||'',oldS=reason.querySelector('span')?.textContent||'';if(reason.querySelector('b'))reason.querySelector('b').textContent='Within expected freshness threshold · provider timestamp unavailable; using DE.PULSE receipt after a targeted provider reconciliation while the fallback provider remains armed for recovery without hiding any diagnostic evidence';if(reason.querySelector('span'))reason.querySelector('span').textContent='Twelve Data → yfinance → CBOE official delayed validation fallback';const box=e=>{const r=e.getBoundingClientRect();return {left:r.left,right:r.right,top:r.top,bottom:r.bottom,width:r.width,height:r.height,scrollWidth:e.scrollWidth,clientWidth:e.clientWidth}};const rb=box(reason),cb=box(check),db=box(data),ab=action?box(action):null,children=[...reason.querySelectorAll('b,span')].map(box);if(reason.querySelector('b'))reason.querySelector('b').textContent=oldB;if(reason.querySelector('span'))reason.querySelector('span').textContent=oldS;return {reason:rb,check:cb,data:db,action:ab,children,docOverflow:document.documentElement.scrollWidth-document.documentElement.clientWidth}}''')
    checks+=1
    if (not long_freshness or long_freshness['docOverflow']>1 or long_freshness['reason']['scrollWidth']>long_freshness['reason']['clientWidth']+1 or
        any(x['left']<long_freshness['reason']['left']-1 or x['right']>long_freshness['reason']['right']+1 for x in long_freshness['children'])):
        errors.append('v15.1.2 Data Freshness Reason/Fallback containment regression: '+json.dumps(long_freshness))

    # v15.1.2 edge hardening: the compact readiness side cards must keep titles, action buttons, and dots
    # in their own header geometry even for the longest catalyst title / running state.
    page.evaluate("window.DePulseLogic.setPage('dashboard')"); page.wait_for_timeout(25)
    prep_layout=page.evaluate('''()=>{const cards=[...document.querySelectorAll('.sidebar-prep-row')];const overlap=(a,b)=>a&&b&&a.left<b.right&&a.right>b.left&&a.top<b.bottom&&a.bottom>b.top;return cards.map(c=>{const title=c.querySelector('.sidebar-prep-label'),action=c.querySelector('.sidebar-prep-actionline'),btn=action?.querySelector('button'),dot=action?.querySelector('.status-dot'),tr=title?.getBoundingClientRect(),ar=action?.getBoundingClientRect(),br=btn?.getBoundingClientRect(),dr=dot?.getBoundingClientRect(),cr=c.getBoundingClientRect();return {key:c.dataset.prepKey,title:title?.textContent||'',titleActionOverlap:overlap(tr,ar),buttonDotOverlap:overlap(br,dr),inside:!!(tr&&ar&&tr.left>=cr.left-1&&tr.right<=cr.right+1&&ar.right<=cr.right+1&&ar.left>=cr.left-1),overflow:c.scrollWidth-c.clientWidth,titleOverflow:title?title.scrollWidth-title.clientWidth:99}})}''')
    checks+=1
    if not prep_layout or any(x['titleActionOverlap'] or x['buttonDotOverlap'] or not x['inside'] or x['overflow']>1 or x['titleOverflow']>1 for x in prep_layout):
        errors.append('v15.1.2 readiness side-card title/action containment regression: '+json.dumps(prep_layout))

    # v15.1 SEC detail is explicitly reachable in Research -> SEC & Ownership; Overview retains concise SEC risk/context.
    page.evaluate("window.DePulseLogic.openResearch('AAA','dashboard')"); page.wait_for_timeout(20)
    if page.locator('[data-research-tab="overview"]').count(): page.locator('[data-research-tab="overview"]').first.click(force=True); page.wait_for_timeout(25)
    overview_text=page.locator('#main').inner_text(); checks+=1
    if 'SEC' not in overview_text or 'Research Decision Brief' not in overview_text:
        errors.append('v15.1 Research overview SEC-context visibility regression')
    sec_loc=page.locator('[data-research-tab="sec"]')
    if sec_loc.count(): sec_loc.first.click(force=True); page.wait_for_timeout(25)
    sec_text=page.locator('#main').inner_text(); checks+=1
    sec_upper=sec_text.upper()
    if 'SEC & OWNERSHIP' not in sec_upper or 'LAST 30 DAYS' not in sec_upper or 'BUY' not in sec_upper or 'OTHER' not in sec_upper or 'JANE EXAMPLE' not in sec_upper or 'ROLE' not in sec_upper:
        errors.append('v15.1 Research SEC & Ownership detail visibility regression')

    # v15.1.2 SEC edge hardening: very long insider/role/type text stays contained in the SEC table's own
    # horizontal-scroll surface and never creates page-level horizontal overflow.
    sec_containment=page.evaluate('''()=>{const wrap=document.querySelector('.sec-transaction-table-wrap'),table=document.querySelector('.sec-transaction-table');if(!wrap||!table)return null;const cell=table.querySelector('tbody tr td:nth-child(2)');if(!cell)return null;const old=cell.innerHTML;cell.innerHTML='<b>Alexandria-Cassandra Extremely-Long Insider Family Name for Containment Validation</b><small>Executive Vice President, Chief Financial Officer, Interim Chief Operating Officer & Director</small>';const wr=wrap.getBoundingClientRect(),tr=table.getBoundingClientRect(),cr=cell.getBoundingClientRect(),result={wrapLeft:wr.left,wrapRight:wr.right,tableWidth:tr.width,cellInsideTable:cr.left>=tr.left-1&&cr.right<=tr.right+1,wrapScrolls:wrap.scrollWidth>=wrap.clientWidth,docOverflow:document.documentElement.scrollWidth-document.documentElement.clientWidth,wrapParentOverflow:wrap.parentElement?wrap.parentElement.scrollWidth-wrap.parentElement.clientWidth:0};cell.innerHTML=old;return result}''')
    checks+=1
    if not sec_containment or not sec_containment['cellInsideTable'] or not sec_containment['wrapScrolls'] or sec_containment['docOverflow']>1 or sec_containment['wrapParentOverflow']>1:
        errors.append('v15.1.2 SEC long insider/role containment regression: '+json.dumps(sec_containment))
    page.evaluate("window.DePulseLogic.setPage('long')"); page.wait_for_timeout(30)
    sec_long=page.locator('#main').inner_text(); checks+=1
    if 'SEC Intelligence Summary' not in sec_long or 'BUY' not in sec_long or 'OTHER' not in sec_long or 'Jane Example' not in sec_long:
        errors.append('v15 Long SEC transaction-detail visibility regression')

    # v18 permanent documentation typography hierarchy is verified in the general responsive renderer matrix too.
    page.set_viewport_size({'width':1512,'height':982}); page.evaluate("window.DePulseLogic.setPage('documentation')"); page.wait_for_timeout(20)
    doc_type=page.evaluate('''()=>{const host=document.querySelector('#doc-content');host.innerHTML=window.DePulseLogic.renderMarkdown(`# Page Title\n\n## Release Section\n\n### Nested Section\n\nBody`);const fs=s=>parseFloat(getComputedStyle(document.querySelector(s)).fontSize);return {h1:fs('#doc-content h1'),h2:fs('#doc-content h2'),h3:fs('#doc-content h3'),overflow:host.scrollWidth-host.clientWidth}}''')
    checks+=1
    if not doc_type or abs(doc_type['h1']-21)>.2 or abs(doc_type['h2']-15)>.2 or abs(doc_type['h3']-12.5)>.2 or doc_type['overflow']>1:
        errors.append('v18 desktop documentation typography hierarchy regression: '+json.dumps(doc_type))
    page.set_viewport_size({'width':720,'height':1050}); page.wait_for_timeout(20)
    doc_mobile=page.evaluate('''()=>{const fs=s=>parseFloat(getComputedStyle(document.querySelector(s)).fontSize);const host=document.querySelector('#doc-content');return {h1:fs('#doc-content h1'),h2:fs('#doc-content h2'),h3:fs('#doc-content h3'),overflow:host.scrollWidth-host.clientWidth}}''')
    checks+=1
    if not doc_mobile or abs(doc_mobile['h1']-18)>.2 or abs(doc_mobile['h2']-14)>.2 or abs(doc_mobile['h3']-12)>.2 or doc_mobile['overflow']>1:
        errors.append('v18 mobile documentation typography hierarchy regression: '+json.dumps(doc_mobile))

    # Compact/collapsed-sidebar behavior at the supported narrow desktop/tablet rail width.
    page.set_viewport_size({'width':820,'height':1050}); page.evaluate("window.DePulseLogic.setPage('dashboard')"); page.wait_for_timeout(25)
    compact_sidebar=page.evaluate('''()=>{const s=document.querySelector('.sidebar').getBoundingClientRect(), n=document.querySelector('.nav-item'); return {width:Math.round(s.width), navTextDisplay:n?getComputedStyle(n.querySelector('.nav-label')||n).display:'missing', overflow:document.documentElement.scrollWidth-document.documentElement.clientWidth}}''')
    checks+=1
    if compact_sidebar['width']>100 or compact_sidebar['overflow']>1:
        errors.append('compact sidebar responsive regression: '+json.dumps(compact_sidebar))

    # Cross-device DPI/scaling acceptance. CSS viewport dimensions emulate the logical workspace after OS scaling;
    # deviceScaleFactor exercises 125%/150% Windows-style raster scaling and Retina rendering.
    scale_cases=[('windows-125',1024,768,1.25),('windows-150',960,600,1.5),('mac-retina',1512,982,2.0)]
    scale_fixture='''()=>{const now=Date.now();const state={version:window.__DEPULSE_EXPECTED_VERSION__,buildId:window.__DEPULSE_EXPECTED_BUILD_ID__,settings:{dataMode:'demo',globalProviderMode:'auto',optionsDataMode:'auto',macroEventModeEnabled:true,dayEnabled:false,swingEnabled:false,longEnabled:false,dayWatchlistId:'day',swingWatchlistId:'swing',longWatchlistId:'long',discoveryWatchlistId:'discovery'},watchlists:[{id:'day',name:'Day',symbols:[]},{id:'swing',name:'Swing',symbols:[]},{id:'long',name:'Long',symbols:[]},{id:'discovery',name:'Discovery',symbols:[]}],ui:{selectedTicker:'SPY',scopeType:'general'},providerStatus:{},cacheInfo:{sizeBytes:0,cachedSymbols:0,lastUpdated:now}};const runtime={status:'stopped',mode:'demo',message:'QA scaling fixture',quotes:{},bars:{},history:{},fundamentals:{},news:[],earnings:[],filings:[],secIntelligence:{},scanner:{mode:'day',status:'idle',results:[]},lastUpdated:{},health:{},feed:{marketSession:'closed',feedState:'stopped'},global:{tone:'NEUTRAL',confidence:0,drivers:{}},macroMetrics:{},macroEvents:[],eventMode:{active:false},eventReactions:[],options:{},capabilities:[],signalValidation:{snapshots:[]}};window.DePulseLogic.setState(state,runtime);window.DePulseLogic.setPage('dashboard')}'''
    for label,w,h,dpr in scale_cases:
        ctx=browser.new_context(viewport={'width':w,'height':h},device_scale_factor=dpr)
        sp=ctx.new_page()
        sp.set_content(html,wait_until='load',timeout=30000)
        sp.add_style_tag(content=(Path.cwd()/"renderer"/"styles.css").read_text())
        sp.evaluate("window.__DEPULSE_TEST__=true")
        sp.add_script_tag(content=(Path.cwd()/"renderer"/"renderer.js").read_text())
        sp.wait_for_function("window.DePulseLogic",timeout=30000)
        sp.evaluate(scale_fixture); sp.wait_for_timeout(25)
        for surface in ['dashboard','settings']:
            sp.evaluate("([surface])=>window.DePulseLogic.setPage(surface)",[surface]); sp.wait_for_timeout(25)
            scale_result=sp.evaluate('''()=>{const de=document.documentElement,b=document.body,m=document.querySelector('#main');const bad=[...document.querySelectorAll('button,select,input,textarea')].filter(e=>{const r=e.getBoundingClientRect(),s=getComputedStyle(e);return s.display!=='none'&&r.width>0&&(r.right>innerWidth+2||r.left<-2)});return {dpr:devicePixelRatio,docOverflow:de.scrollWidth-de.clientWidth,bodyOverflow:b.scrollWidth-b.clientWidth,text:(m?.innerText||'').trim().length,bad:bad.length}}''')
            checks+=1
            if abs(scale_result['dpr']-dpr)>.01 or scale_result['docOverflow']>1 or scale_result['bodyOverflow']>1 or scale_result['text']<10 or scale_result['bad']:
                errors.append(f'{label} {surface} scaling regression: '+json.dumps(scale_result))
        ctx.close()
    browser.close()
print(f"responsive UI matrix {checks}/{checks if not errors else checks}: {'PASS' if not errors else 'FAIL'}")
print(f"viewports={len(viewports)} surfaces={len(surfaces)}")
if errors:
    print('\n'.join(errors[:40]))
    sys.exit(1)
