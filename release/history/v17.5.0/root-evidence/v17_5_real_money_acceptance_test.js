const fs=require('fs'),vm=require('vm'),assert=require('assert');
const context={console,globalThis:{__DEPULSE_TEST__:true},document:{querySelector(){return null},querySelectorAll(){return []}},window:{addEventListener(){}},fetch:async()=>({ok:false,status:404,text:async()=>'',json:async()=>({})}),EventSource:function(){},ResizeObserver:function(){this.observe=()=>{}},setInterval(){},setTimeout,clearTimeout,requestAnimationFrame:f=>f()};
context.globalThis=context;context.__DEPULSE_TEST__=true;vm.createContext(context);vm.runInContext(fs.readFileSync('renderer/renderer.js','utf8'),context);
const L=context.DePulseLogic,now=Date.now();

// v17.4/v17 Major Closure: one shared freshness problem must not masquerade as many independent liquidity problems.
const exceptions=Array.from({length:9},(_,i)=>({symbol:`RISK${i+1}`,reason:'LIQUIDITY RISK',severity:'HIGH',target:'research',source:'Market Open Prep',updatedAt:now}));
exceptions.push({reason:'Quotes · STALE · upstream live quote coverage delayed',severity:'HIGH',target:'maintenance',source:'Live equity router',updatedAt:now});
const groups=L.v174PreparationExceptionGroups(exceptions);
const liquidity=groups.find(x=>x.reason==='LIQUIDITY RISK');
assert(liquidity&&liquidity.members.length===9,'genuine symbol liquidity risks must be grouped without losing members');
const exceptionMarkup=L.v174PreparationExceptionsMarkup(exceptions);
assert((exceptionMarkup.match(/LIQUIDITY RISK/g)||[]).length===1,'repeated genuine cause must render once at group level');
assert(exceptionMarkup.includes('9 symbols'),'group count must remain explicit');
assert(exceptionMarkup.includes('Quotes · STALE'),'shared freshness root cause must remain explicit');
assert((exceptionMarkup.match(/data-readiness-drill=/g)||[]).length===9,'every affected symbol must retain Review drill-down');

// Operational truth must expose whether critical evidence is usable and the v17 SLO/warm-start/recovery diagnostics.
const state={version:'17.5.1',settings:{},watchlists:[],ui:{page:'maintenance',selectedTicker:'SPY'}};
const runtime={runtimeLoad:{cpuUtilizationPct:17,goroutines:21,gomaxprocs:4,heapAllocBytes:1024,numGC:2,sampledAt:now,http:{interactiveP95Ms:42,requests:10,inFlight:0,slowInteractive:0,interactiveMaxMs:65},persistence:{backend:'sqlite',ready:true,queueDepth:0,writeBatchesLastMinute:2,rowsWrittenLastMinute:3,materialWritesSuppressed:10,store:{schemaVersion:3,activeSymbolCount:1,symbolCount:1,canonicalQuotes:1,quoteHistoryRows:1,evidenceRows:1,decisionRows:1,storageBytes:4096}},providerCallsAvoided:12,canonicalReuseHitRatePct:80,storageGrowthBytes:0,startup:{warmStartCoveragePct:100,bootstrapDurationMs:15,warmStartQuotes:1,warmStartTargetQuotes:1,cacheQuotesLoaded:1,persistedQuotesApplied:1},recovery:{currentlyStaleDatasets:0,staleToCurrentEvents:2,lastStaleToCurrentMs:900,degradationRecoveryEvents:1},workload:[],providerRequests:[]},degradation:{code:'PARTIAL COVERAGE',criticalUsable:true,detail:'Broad Discovery delayed · Selected/Watchlist LIVE'},runtimeSlo:{status:'PASS',checks:[]},feed:{},quotes:{},bars:{},lastUpdated:{},health:{},scanner:{},marketIntelligence:{},freshness:[]};
L.setState(state,runtime);const maintenance=L.renderMaintenance();
for(const token of ['Performance & Runtime Load','Warm Start','Recovery','Critical Decision Data','USABLE','PARTIAL COVERAGE']) assert(maintenance.includes(token),`Maintenance must expose ${token}`);

// Executable HTTP surface must remain decision-support only: no order/broker/portfolio/journal endpoints entered in v17.
const api=fs.readFileSync('http_api.go','utf8');
const routes=[...api.matchAll(/HandleFunc\("([^"]+)"/g)].map(m=>m[1]);
assert(routes.length>10,'expected executable HTTP route inventory');
const forbidden=routes.filter(r=>/(^|\/)(orders?|broker|execute|execution|portfolio|journal|blotter)(\/|$)/i.test(r));
assert.deepEqual(forbidden,[],`No Execution Boundary violated by routes: ${forbidden.join(', ')}`);

console.log(`v17.5 real-money closure acceptance: PASS · grouped risk truth · critical-data usability/SLO diagnostics · No Execution HTTP surface (${routes.length} routes)`);
