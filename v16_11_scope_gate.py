#!/usr/bin/env python3
"""v16.11 Major Closure scope/traceability gate.
Intentionally static/bounded. Fresh v16-family regressions, expert reviews,
performance and inherited 30/30 are separate fingerprint-scoped evidence jobs.
"""
from pathlib import Path
import json,sys
R=Path(__file__).resolve().parent; errs=[]
required=[
 'v16_11_major_closure_scope_matrix.json','v16_11_major_scope_matrix_gate.py',
 'v16_11_senior_engineer_gate.py','v16_11_trader_investor_gate.py',
 'v16_11_performance_stability_gate.py','v16_11_data_utility_source_hygiene_gate.py',
 'v16_11_adaptive_retrospective_gate.py','v16_11_source_hygiene_gate.py',
 'v16_11_real_money_acceptance_test.js','v16_11_major_closure_test.go',
 'renderer/qa/v16.11.0-master-scope.json','renderer/qa/v16.11.0-acceptance-evidence.json',
 'renderer/qa/v16.11.0-traceability.md','renderer/qa/v16.11.0-senior-engineer-review.md',
 'renderer/qa/v16.11.0-trader-investor-review.md','renderer/qa/v16.11.0-adaptive-retrospective.md',
 'renderer/qa/v16.11.0-data-utility-source-hygiene-audit.md']
for f in required:
 if not (R/f).exists(): errs.append('missing '+f)
try:
 s=json.loads((R/'renderer/qa/v16.11.0-master-scope.json').read_text())
 if len(s.get('scope_lock',[]))!=10: errs.append('major closure scope must be 10 clauses')
 a=json.loads((R/'renderer/qa/v16.11.0-acceptance-evidence.json').read_text())
 if len(a.get('clauses',[]))!=10 or any(x.get('status')!='PASS' for x in a.get('clauses',[])): errs.append('acceptance evidence is not 10/10 PASS')
 if a.get('original_professional_roadmap')!='30 FULL / 0 PARTIAL / 0 MISSING': errs.append('original roadmap closure drift')
 if a.get('v16_10_scope')!='10/10 PASS': errs.append('v16.10 inherited closure drift')
 m=json.loads((R/'v16_11_major_closure_scope_matrix.json').read_text())
 if len(m.get('original_professional_roadmap',[]))!=30: errs.append('major matrix original roadmap !=30')
 if len(m.get('major_family_milestones',[]))!=12: errs.append('major matrix v16 milestones !=12')
 ident=json.loads((R/'release_identity.json').read_text())
 version=str(ident.get('version',''))
 if version=='16.11.0':
  if ident.get('previous_stable')!='v16.10.0': errs.append('v16.11 historical previous-Stable identity drift')
 elif version.startswith('17.'):
  # RL-034: later v17 TEST/RC may be current while v16.11 remains the authoritative Stable baseline.
  if ident.get('stable_baseline')!='v16.11.0' or ident.get('previous_stable')!='v16.11.0': errs.append('v16.11 Stable-baseline preservation drift')
 elif version.startswith('18.'):
  # v18 inherits the v16.11 closure through the certified v17.5.1 Stable major-closure baseline.
  if ident.get('stable_baseline')!='v17.5.1' or ident.get('previous_stable')!='v17.5.1': errs.append('v16.11 closure inheritance through v17.5.1 drift')
 else: errs.append('release identity is outside v16.11/v17/v18 preservation contract')
 reg=json.loads((R/'release_learning_registry.json').read_text())
 for rid in ('RL-029','RL-030','RL-031'):
  if not any(x.get('id')==rid and x.get('status')=='ACTIVE' for x in reg.get('entries',[])): errs.append(rid+' not ACTIVE')
except Exception as ex: errs.append('closure metadata unreadable: '+str(ex))
# Strong static safety truth for the real-money fix; behavior is independently tested.
rend=(R/'renderer/renderer.js').read_text(errors='ignore')
for token in ['STALE','CACHED','HISTORY ONLY','INCOMPLETE','CONDITIONAL']:
 if token not in rend: errs.append('readiness/freshness safety token missing: '+token)
if errs:
 print('v16.11 Major Closure Scope Gate: FAIL'); [print(' -',e) for e in errs]; raise SystemExit(2)
print('v16.11 Major Closure Scope Gate: PASS · 10/10 closure obligations structurally evidenced · heavy inherited evidence remains independently checkpointed')
