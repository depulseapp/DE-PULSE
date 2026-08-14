#!/usr/bin/env python3
from pathlib import Path
import json,sys
R=Path(__file__).resolve().parent; errs=[]
try:
    scope=json.loads((R/'renderer/qa/v16.10.0-master-scope.json').read_text())
    rows=scope.get('scope_lock',[])
    if len(rows)!=10 or {int(x['id']) for x in rows}!=set(range(1,11)): errs.append('v16.10 immutable scope must contain exactly 10 clauses')
    audit=(R/'v16_10_information_architecture_data_efficiency_audit.md').read_text(errors='ignore')
    for token in ['Surface decisions','Data / processing decisions','Rapid Price Dislocation Watch','Short Trade Plan','SPY / QQQ drawdown from ATH','Discovery/Scanner remains the only scanner owner']:
        if token not in audit: errs.append('IA/Data Efficiency audit missing '+token)
    radar=(R/'opportunity_radar.go').read_text(errors='ignore')
    for token in ['opportunityMaxPromotions','sessionRelativeVolumeFromSnapshot','rangeExpansionFromSnapshot','livePriorityHints','requestHistoryHydration','buildCommunityEvidenceFusion','marketActivity']:
        if token not in radar: errs.append('Radar implementation missing '+token)
    if 'order execution' in radar.lower() or 'paper trading' in radar.lower(): errs.append('Radar source crossed execution boundary')
    adaptive=(R/'adaptive_data_policy.go').read_text(errors='ignore')
    for token in ['ShadowControlState','SHADOW','CanMutateProduction']:
        if token not in adaptive: errs.append('Adaptive Shadow implementation missing '+token)
    renderer=(R/'renderer/renderer.js').read_text(errors='ignore')
    if 'data-page="opportunity"' in renderer or 'data-page="radar"' in renderer: errs.append('Opportunity Radar created a duplicate top-level workspace')
    dh=json.loads((R/'data_health_policy.json').read_text())
    if dh.get('version')!='16.10.0' or dh.get('policy',{}).get('shadow_policy_can_mutate_production') is not False: errs.append('Adaptive Data Health v16.10 policy incomplete')
    du=json.loads((R/'data_utility_registry.json').read_text()); names={x.get('dataset') for x in du.get('datasets',[])}
    if not {'Opportunity Radar','Market Activity'}.issubset(names): errs.append('Data Utility registry missing Radar/Market Activity utility')
    status169=json.loads((R/'renderer/qa/v16.9.0-original-roadmap-status.json').read_text())
    if status169.get('summary')!='30 FULL / 0 PARTIAL / 0 MISSING': errs.append('v16.9 inherited roadmap status drift')
    shadow=json.loads((R/'shadow_experiments.json').read_text())
    if shadow.get('promotion_path')!='SHADOW -> VALIDATED -> APPROVED -> PRODUCTION': errs.append('Shadow promotion path drift')
    if any(x.get('can_mutate_production') for x in shadow.get('experiments',[])): errs.append('Shadow experiment can mutate production')
    evidence=json.loads((R/'renderer/qa/v16.10.0-acceptance-evidence.json').read_text())
    if evidence.get('status')!='10/10 PASS' or len(evidence.get('clauses',[]))!=10: errs.append('v16.10 acceptance evidence incomplete')
except Exception as ex:
    errs.append('v16.10 scope evidence unreadable: '+str(ex))
if errs:
    print('v16.10 Scope Gate: FAIL'); [print(' -',e) for e in errs]; sys.exit(2)
print('v16.10 Scope Gate: PASS · 10/10 static/traceability evidence · executable regressions owned by independent checkpoints (RL-027)')
