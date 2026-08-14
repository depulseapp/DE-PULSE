const fs=require('fs');
const src=fs.readFileSync('renderer/renderer.js','utf8');
const need=(c,m)=>{if(!c){console.error('v16.10 renderer acceptance: FAIL · '+m);process.exit(1)}};
for(const token of [
  'Always-On Opportunity Radar','Unusual volume, volatility & rapid dislocation','sessionRelativeVolume','rangeExpansion','opportunityScore',
  'Adaptive Data Policy & Shadow Control','SHADOW → VALIDATED → APPROVED → PRODUCTION','observes alternatives only',
  'discovery-market-scanner-v1610','Opportunity Radar starts automatically','Targeted hot-symbol refresh'
]) need(src.includes(token),`missing ${token}`);
need(!src.includes('data-page="opportunity"')&&!src.includes('data-page="radar"'),'Opportunity Radar must extend Discovery, not create a top-level workspace');
need(src.includes("runtime?.scanner?.radar?.candidates"),'Discovery quote-impact path must include radar candidates');
console.log('v16.10 renderer acceptance: PASS · Opportunity Radar + adaptive data policy + Shadow controls');
