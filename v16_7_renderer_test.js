const fs=require('fs');
const src=fs.readFileSync('renderer/renderer.js','utf8');
const need=(c,m)=>{if(!c){console.error('v16.7 renderer acceptance: FAIL · '+m);process.exit(1)}};
for(const token of [
  "marketCalendarFilters={impact:'ALL',scope:'US',category:'ALL',date:'',sort:'TIME_ASC'}",
  'data-calendar-filter="impact"','data-calendar-filter="scope"','data-calendar-filter="category"','data-calendar-filter="date"','data-calendar-filter="sort"',
  'miBreadthInternalsMarkup','Above Key MAs','20-Session High / Low','Sector Participation',
  'miTradeabilityComponentsMarkup','Breadth / Internals','Macro / Event','Setup Availability','Options Context','Market Tradeability ·',
  "String(x.horizon||'').toUpperCase()",'Member breadth',
]) need(src.includes(token),`missing ${token}`);
need(src.includes('A ${x.actual==null')&&src.includes('/ F ${x.forecast==null')&&src.includes('/ P ${x.previous==null'),'calendar actual/forecast/previous UI missing');
need(src.includes('No calendar events match the selected filters.'),'calendar empty-filter state missing');
need(src.includes('tradeReadiness')&&src.includes('relativeStrength'),'readiness/relative-strength integration missing');
need(src.includes("marketTradeabilityState==='DATA DEGRADED'")&&src.includes("marketTradeabilityState==='WAIT'")&&src.includes("marketTradeabilityState==='REDUCE SIZE'"),'Tradeability not consumed by desk readiness');
need(src.includes('queuePriorityFor')&&src.includes('tradeReadiness'),'Decision Queue not fed by readiness/Tradeability');
console.log('v16.7 renderer acceptance: PASS · calendar filters + breadth internals + tradeability components + Day/Swing RS surface');
