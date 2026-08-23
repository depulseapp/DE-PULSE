const fs=require('fs'),vm=require('vm'),assert=require('assert');
const releaseIdentity=JSON.parse(fs.readFileSync('release_identity.json','utf8'));
const c={console,globalThis:{__DEPULSE_TEST__:true},document:{querySelector(){return null},querySelectorAll(){return []}},window:{addEventListener(){}},fetch:async()=>({ok:false,status:404,text:async()=>'',json:async()=>({})}),EventSource:function(){},ResizeObserver:function(){this.observe=()=>{}},setInterval(){},setTimeout,clearTimeout,requestAnimationFrame:f=>f()};c.globalThis=c;c.__DEPULSE_TEST__=true;vm.createContext(c);vm.runInContext(fs.readFileSync('renderer/renderer.js','utf8'),c);const L=c.DePulseLogic;
const now=Date.now(),future=d=>new Date(now+d*86400000).toISOString().slice(0,10);
function rising(n=120,start=120){return Array.from({length:n},(_,i)=>{const p=start+i*.7;return {o:p-.3,h:p+1,l:p-1,c:p,v:1000000+i*25000}})}
const b=rising(),last=b[b.length-1].c;
const quotes={NVDA:{symbol:'NVDA',price:last,previousClose:last-3,changePercent:2,updatedAt:now,providerTimestamp:now,source:'finnhub-websocket',dataState:'live'},SPY:{price:600,changePercent:1.2,updatedAt:now},QQQ:{price:520,changePercent:1.5,updatedAt:now},DIA:{price:450,changePercent:.8,updatedAt:now},IWM:{price:225,changePercent:1.1,updatedAt:now},XLK:{price:250,changePercent:1.4,updatedAt:now},TLT:{price:90,changePercent:-.5,updatedAt:now},GLD:{price:300,changePercent:-.2,updatedAt:now},SLV:{price:35,changePercent:0,updatedAt:now},USO:{price:80,changePercent:.2,updatedAt:now},VIX:{price:15,changePercent:-6,updatedAt:now,source:'vix-index',dataState:'live'}};
const bars={};for(const s of Object.keys(quotes))bars[s]={daily:b,intraday:b,weekly:b.filter((_,i)=>i%5===0)};
const state={version:releaseIdentity.version,settings:{dataMode:'live',overnightDataMode:'auto',earningsPenalty:10,marketContext:15,signalProfile:'balanced',swingWatchlistId:'swing',dayWatchlistId:'day',longWatchlistId:'long',discoveryWatchlistId:'discovery',dayEnabled:true,swingEnabled:true,longEnabled:true},watchlists:[{id:'day',name:'Day',symbols:['NVDA']},{id:'swing',name:'Swing',symbols:['NVDA']},{id:'long',name:'Long',symbols:['NVDA']},{id:'discovery',name:'Discovery',symbols:[]}],ui:{selectedTicker:'NVDA',scopeType:'watchlist',watchlistId:'day'}};
const runtime={status:'running',mode:'live',quotes,bars,history:{},fundamentals:{NVDA:{revenueGrowth:25,epsGrowth:30,roe:35,netMargin:28,debtToEquity:20,pe:30,forwardPe:25,freeCashFlow:20000000000}},news:[],earnings:[],filings:[],secIntelligence:{},scanner:{mode:'day',status:'ready',results:[]},lastUpdated:{quotes:now,history:now,fundamentals:now,news:now,earnings:now,filings:now,vix:now},health:{},feed:{marketSession:'regular',feedState:'streaming',webSocketConnected:true,subscribedSymbols:['NVDA']},global:{tone:'RISK-OFF',confidence:85,drivers:{semiconductors:{state:'SUPPORTIVE',label:'Global Semiconductor Tone',detail:'SMH real proxy'},taiwan:{state:'SUPPORTIVE',label:'Taiwan',detail:'EWT real proxy'},rates_10y:{state:'HEADWIND',label:'10Y',detail:'4.4%'}}},macroMetrics:{DGS10:{value:4.4},DFII10:{value:2.0},BAMLH0A0HYM2:{value:3.4}},macroEvents:[{id:'cpi',name:'CPI',impact:'HIGH',timeKnown:true,startsAt:now+10*60000,date:future(0)}],eventMode:{active:true,eventId:'cpi',name:'CPI',startsAt:now+10*60000,countdownS:600,phase:'PREP'},eventReactions:[],options:{NVDA:{symbol:'NVDA',provider:'Alpaca Options',feed:'OPRA',state:'CURRENT',bias:'BEARISH',callVolume:400,putVolume:1400,putCallVolume:3.5,averageIv:.55,expectedMove:12,nearestExpiration:future(20),updatedAt:now,provenance:'REAL OPTIONS SNAPSHOT'}},capabilities:[],signalValidation:{snapshots:[]}};
L.setState(state,runtime);
const dayPlan=L.planFor('NVDA','day');assert(dayPlan.direction.includes('BULLISH'),`fixture must be bullish, got ${dayPlan.direction}`);
const dayR=L.tradeReadiness('NVDA','day');assert.equal(dayR.status,'CAUTION');assert(dayR.risks.some(x=>x.includes('macro')));assert(dayR.risks.some(x=>x.includes('Options')));assert(dayR.risks.some(x=>x.includes('Global')));
// A macro CAUTION cannot be downgraded to CONDITIONAL by a later contradiction.
assert.equal(dayR.status,'CAUTION');
// Day gets immediate options/event context.
assert(L.optionsContextLine('NVDA','day').includes('P/C'));assert(L.globalContextLine('NVDA','day').includes('HIGH event'));assert(L.contextLens('NVDA','day').includes('Macro Event Reaction')&&L.contextLens('NVDA','day').includes('waiting for measured snapshot'));
// Swing gets persistent-positioning style context without Day-only event countdown.
assert(L.optionsContextLine('NVDA','swing').includes('positioning'));assert(!L.globalContextLine('NVDA','swing').includes('HIGH event'));
// Long excludes short-dated options and current-session global noise from readiness.
const longR=L.tradeReadiness('NVDA','long');assert(longR.reasons.some(x=>x.includes('Short-dated options intentionally excluded')));assert(L.optionsContextLine('NVDA','long').includes('Intentionally hidden'));assert(L.globalContextLine('NVDA','long').includes('session noise secondary'));
// If long-dated options exist, they can be shown selectively.
runtime.options.NVDA.nearestExpiration=future(70);L.setState(state,runtime);assert(L.optionsContextLine('NVDA','long').includes('long-dated context'));
// AI evidence carries correlated context but remains separate from deterministic plan.
const ctx=L.buildAIClientContext('NVDA','day');assert(ctx.options&&ctx.eventMode&&ctx.globalContext&&ctx.readiness);assert.equal(L.planFor('NVDA','day').score,dayPlan.score);
assert.notEqual(L.aiResultKey('NVDA','day','ticker'),L.aiResultKey('NVDA','day','risk'));assert.notEqual(L.aiResultKey('NVDA','day','risk'),L.aiResultKey('NVDA','day','news'));

// v14.3 contextual intelligence can strengthen readiness/context without changing the frozen deterministic plan.
const frozenDayScore=L.planFor('NVDA','day').score;
runtime.liquidity={NVDA:{symbol:'NVDA',state:'RISK',spreadPercent:1.25,quoteAgeMs:250,detail:'wide spread'}};
runtime.marketOpenFlags={NVDA:['LIQUIDITY RISK','EVENT RISK']};runtime.marketOpenCheckpoint={state:'COMPLETE',runAt:now,premarket:{NVDA:{gapPercent:4.5,rangePercent:6.2,volume:2500000}},optionsContexts:1,dayCandidates:1,swingContextChanges:['NVDA'],longContextChanges:[]};
runtime.catalystReactions={NVDA:{symbol:'NVDA',triggerType:'EARNINGS',state:'EXTREME VOLATILITY',movePercent:11.4,flags:['LIQUIDITY RISK'],detail:'premarket reaction; confirmation pending'}};
runtime.symbolIntelligence={NVDA:{symbol:'NVDA',consecutiveBeats:3,peers:['AMD','AVGO'],recommendationTrend:'BULLISH',priceTarget:250}};
runtime.intelligence={rates:{label:'Rates State',state:'RISING',detail:'10Y/2Y pressure'},credit:{label:'Credit State',state:'HEALTHY',detail:'HY spread contained'},'financial-conditions':{label:'Financial Conditions',state:'NEUTRAL'},inflation:{label:'Inflation State',state:'MIXED'},labor:{label:'Labor State',state:'BALANCED'},energy:{label:'Energy State',state:'BALANCED'}};
runtime.corporateActions=[{symbol:'NVDA',type:'split',exDate:future(5)}];
L.setState(state,runtime);
const v143R=L.tradeReadiness('NVDA','day');assert.equal(v143R.status,'CAUTION');assert(v143R.risks.some(x=>x.includes('Liquidity')));assert(v143R.risks.some(x=>x.includes('EVENT RISK')));assert(v143R.risks.some(x=>x.includes('OUTSIDE ENTRY ZONE')));assert(v143R.reasons.some(x=>x.includes('EARNINGS reaction')));
assert.equal(L.planFor('NVDA','day').score,frozenDayScore,'v14.3 context must not alter deterministic day score');
const v143Ctx=L.buildAIClientContext('NVDA','day');assert(v143Ctx.liquidity&&v143Ctx.marketOpenFlags.length===2&&v143Ctx.catalystReaction&&v143Ctx.providerIntelligence&&v143Ctx.corporateActions.length===1);
assert(L.contextLens('NVDA','day').includes('Material Catalyst Reaction')&&L.contextLens('NVDA','day').includes('Provider Intelligence'));
assert(L.globalContextLine('NVDA','swing').includes('Rates'));

// Numeric guidance is evidence-gated: no structured source values means no invented range.
assert.equal(L.numericGuidanceMarkup({symbol:'NVDA'}),'');const gm=L.numericGuidanceMarkup({guidanceRevenuePrevLow:10e9,guidanceRevenuePrevHigh:11e9,guidanceRevenueNewLow:11e9,guidanceRevenueNewHigh:12e9,guidanceSource:'8-K'});assert(gm.includes('8-K')&&gm.includes('Revenue guidance'));
// Cross-module UI: Dashboard, desk, Research, Queue/AI all surface shared context without creating an Options desk.
assert(L.renderDashboard().includes('Options Intelligence')&&L.renderDashboard().includes('Market Intelligence')&&L.renderDashboard().includes('Decision Queue'));assert(L.renderMarketIntelligence().includes('Global Market Drivers'));
assert(L.renderDesk('day').includes('Trade Readiness')&&L.renderDesk('day').includes('Options Context')&&L.renderDesk('day').includes('Global Context'));
assert(L.renderResearch().includes('Cross-Module Intelligence'));
assert(!fs.readFileSync('renderer/index.html','utf8').includes('Options Trading Desk'));
console.log('professional trader acceptance scenarios: PASS');
