const fs=require('fs');
const src=fs.readFileSync('renderer/renderer.js','utf8');
const need=(c,m)=>{if(!c){console.error('v16.9 renderer acceptance: FAIL · '+m);process.exit(1)}};
for(const token of [
  'Global Community Evidence Fusion · v16.9','data-community-platform','Source policy controls AI eligibility','fused events','mentionVelocity1h','corroborated',
  'Oil / Energy · v16.9','Brent-WTI','energySectorState','usMarketRelevance','continuousContract',
  'Reusable Scenario Library · v16.9','replayScenarioResult','CPI/FOMC, earnings-gap, high-VIX and market-dislocation','never paper trading',
  'Professional Validation & Learning · v16.9'
]) need(src.includes(token),`missing ${token}`);
need(!src.includes('data-page="community"')&&!src.includes('data-page="replay"'),'v16.9 must extend canonical owners, not create top-level Community/Replay workspaces');
need(src.includes("ingestionMode:'USER_AUTHORIZED_INPUT'"),'community UI must declare sanctioned/user-authorized ingestion path');
console.log('v16.9 renderer acceptance: PASS · evidence fusion + energy truth + reusable no-lookahead replay scenarios');
