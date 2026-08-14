const fs=require('fs');
const src=fs.readFileSync('renderer/renderer.js','utf8');
const need=(c,m)=>{if(!c){console.error('v16.8 renderer acceptance: FAIL · '+m);process.exit(1)}};
for(const token of [
  'v168HeatModes',"['sector','watchlist','broad']",'coverageBasis',
  'v168SeasonalityRows','10-Year Seasonality','new Date().getFullYear()','bestReturnPct','worstReturnPct',
  'Major strikes','Concentration','Expirations','never measured dealer positioning',
  'relativeSpreadMultiple','dollarVolume','openingLiquidity','slippageRisk',
  'Smart Notifications','No new material state-change notification.'
]) need(src.includes(token),`missing ${token}`);
need(!src.includes('data-jump=\"alerts\"')&&!src.includes('renderAlerts'), 'renderer must not introduce standalone Alerts workspace');
need(src.includes('EXPECTED_RELEASE_VERSION='),'canonical current-release identity hook missing');
console.log('v16.8 renderer acceptance: PASS · heat modes + 10-year seasonality + GEX structure + liquidity/slippage + selective notification surfaces');
