#!/usr/bin/env python3
from pathlib import Path
import json, subprocess, sys
R=Path(__file__).resolve().parent
GO='\n'.join(p.read_text(errors='ignore') for p in R.glob('*.go') if not p.name.endswith('_test.go'))
JS=(R/'renderer/renderer.js').read_text(errors='ignore')
ERR=[]
for tok in ['AIResearchPackage','buildAIResearchPackage','aiEvidenceArchitectureVersion','EvidenceSnapshotID','BullEvidenceIDs','BearEvidenceIDs']:
    if tok not in GO: ERR.append('#22 AI Evidence Architecture missing '+tok)
for tok in ['UNTRUSTED DATA','untrustedExternalContent','externalInstructionWarning','sanitizeAIResponse','Setup Score is not win probability']:
    if tok not in GO: ERR.append('#23 external-content safety missing '+tok)
for tok in ['resolveAIRouting','aiCacheKey','loadAICache','storeAICache','MaxOutputTokens','AIRoutingMode']:
    if tok not in GO: ERR.append('#24 routing/cost control missing '+tok)
for tok in ['Bull case','Base case','Bear case','Routing & Cost Policy','Not win probability','Evidence package']:
    if tok not in JS: ERR.append('v16.4 UI integration missing '+tok)
for cmd,label in [(['go','test','-count=1','-run','TestV164','./...'],'Go v16.4 professional regressions'),(['node','v16_4_renderer_test.js'],'renderer v16.4'),(['node','deterministic_equivalence_test.js'],'deterministic 2403/2403')]:
    q=subprocess.run(cmd,cwd=R,text=True,capture_output=True)
    if q.returncode: ERR.append(label+' failed: '+(q.stdout+q.stderr)[-1400:])
if ERR:
    print('v16.4 Professional Research AI Scope Gate: FAIL'); print('\n'.join('- '+x for x in ERR)); sys.exit(1)
print('v16.4 Professional Research AI Scope Gate: PASS · 3/3 (#22/#23/#24)')
