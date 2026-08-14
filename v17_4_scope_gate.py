#!/usr/bin/env python3
from pathlib import Path
import json,sys
R=Path(__file__).resolve().parent; errors=[]
def need(ok,msg):
    if not ok: errors.append(msg)
identity=json.loads((R/'release_identity.json').read_text()); slices=json.loads((R/'v17_delivery_slices.json').read_text())
version=str(identity.get('version',''))
is_native_v17=version.startswith('17.') and identity.get('channel') in {'TEST','RC','STABLE'}
is_inherited_v18=version.startswith('18.') and identity.get('channel') in {'TEST','RC','STABLE'} and identity.get('stable_baseline')=='v17.5.1'
need(is_native_v17 or is_inherited_v18,'inherited v17 scope release identity missing')
if is_native_v17: need(identity.get('stable_baseline')=='v16.11.0','v16.11 Stable baseline drifted')
v174=next((x for x in slices.get('slices',[]) if x.get('id')=='v17.4'),None)
need(v174 is not None and v174.get('status') in {'SOURCE TEST CHECKPOINT','TEST QUALIFIED','SOURCE QUALIFIED','COMPLETE'},'v17.4 slice must be completed/preserved before v17 Major Closure')
renderer=(R/'renderer/renderer.js').read_text(); css=(R/'renderer/styles.css').read_text(); prep=(R/'preparation_catalyst.go').read_text(); liq=(R/'preparation_types_liquidity.go').read_text(); gotests=(R/'v17_4_operational_hardening_test.go').read_text(); jstests=(R/'v17_4_renderer_test.js').read_text()
need(renderer.count('function masterMarketSymbolsPanel(')==1 and 'masterMarketSymbolsPanel=function' not in renderer,'OPEN-001 still has duplicate Master Market Symbols renderer owner')
for token in ['master-add-primary','master-remove-all','data-master-add','data-master-remove-all','master-symbol-chip-list']:
    need(token in renderer,f'OPEN-001 behavior/hierarchy missing: {token}')
need('.master-add-row{display:flex;flex-wrap:wrap' in css and '@media(max-width:520px)' in css,'OPEN-001 responsive control layout missing')
for token in ['v174PreparationExceptionGroups','v174PreparationExceptionsMarkup','View all ${rest.length} additional risk group','readiness-exception-group','data-readiness-drill']:
    need(token in renderer+css,f'OPEN-002 grouping/bounded review missing: {token}')
need('slice(0,5)' in renderer,'OPEN-002 visible root-cause list is not bounded')
need('currentLiquidityMarketRisk' in liq,'stale-vs-liquidity classification helper missing')
need(prep.count('currentLiquidityMarketRisk')>=2,'Market Open/Catalyst paths do not both use current liquidity truth')
need(prep.count('deriveLiquidityStatesWithContext(quotes, bars, liquidityBaselines, now)')>=2,'prep/catalyst paths must use canonical liquidity context and learned baselines')
for name in ['TestV174LiquidityRiskRequiresCurrentMarketEvidence','TestV174DerivedStaleLiquidityDoesNotMasqueradeAsMarketRisk']:
    need(name in gotests,f'missing v17.4 backend acceptance: {name}')
need('repeated prep risks grouped/bounded' in jstests,'missing v17.4 renderer acceptance contract')
# v17.3 runtime/Data Engine observability must remain present while v17.4 changes presentation.
for token in ['Performance & Runtime Load','Critical Decision Data','Provider Calls Avoided','Warm Start','Recovery']:
    need(token in renderer,f'Data Engine diagnostics regressed: {token}')
if errors:
    print('v17.4 scope gate: FAIL'); [print(' -',e) for e in errors]; sys.exit(1)
print('v17.4 scope gate: PASS · OPEN-001 balanced controls · OPEN-002 root-cause truth/grouping · operational diagnostics preserved')
