#!/usr/bin/env python3
from pathlib import Path
import subprocess,sys
R=Path(__file__).resolve().parent
GO='\n'.join(p.read_text(errors='ignore') for p in R.glob('*.go') if not p.name.endswith('_test.go'))
JS=(R/'renderer/renderer.js').read_text(errors='ignore')
ERR=[]
for tok in ['SentimentCompositeState','v165Sentiment','ComponentsExpected','Missing']:
    if tok not in GO: ERR.append('#5 Sentiment Composite missing '+tok)
for tok in ['MarketSectorHeatMap','v165HeatMap','v165SectorHeatUniverse','CoveragePct']:
    if tok not in GO: ERR.append('#6 Market/Sector Heat Map missing '+tok)
for tok in ['GEXContextState','GammaOICoveragePct','fetchOptionOpenInterest','GEXState','structural signed gamma']:
    if tok not in GO: ERR.append('#9 GEX missing '+tok)
for tok in ['CommunityEvidenceItem','UNTRUSTED COMMUNITY INTELLIGENCE','handleCommunityEvidence','Untrusted']:
    if tok not in GO: ERR.append('#10 Community Intelligence missing '+tok)
for tok in ['OilEnergyContextState','WTI_OFFICIAL','USO is a tradable futures-based ETF','REFINERY_UTIL']:
    if tok not in GO: ERR.append('#11 Oil/Energy Context missing '+tok)
for tok in ['Sentiment Composite','Market / Sector Heat Map','Gamma Exposure Context','UNTRUSTED COMMUNITY INTELLIGENCE','Oil / Energy']:
    if tok not in JS: ERR.append('v16.5 UI integration missing '+tok)
for cmd,label in [(['go','test','-count=1','-run','TestV165','./...'],'Go v16.5 regressions'),(['node','v16_5_renderer_test.js'],'renderer v16.5'),(['node','deterministic_equivalence_test.js'],'deterministic 2403/2403')]:
    q=subprocess.run(cmd,cwd=R,text=True,capture_output=True)
    if q.returncode: ERR.append(label+' failed: '+(q.stdout+q.stderr)[-1600:])
if ERR:
    print('v16.5 Context & Alternative Intelligence Scope Gate: FAIL'); print('\n'.join('- '+x for x in ERR)); sys.exit(1)
print('v16.5 Context & Alternative Intelligence Scope Gate: PASS · 5/5 (#5/#6/#9/#10/#11)')
