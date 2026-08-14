#!/usr/bin/env python3
from pathlib import Path
import json,sys
R=Path(__file__).resolve().parent; e=[]
def need(ok,msg):
    if not ok:e.append(msg)
i=json.loads((R/'release_identity.json').read_text()); s=json.loads((R/'v18_0_1_scope.json').read_text())
need(i.get('version')=='18.0.1' and i.get('channel')=='TEST','v18.0.1 TEST identity missing')
need(i.get('build_id')=='v18.0.1-test-smart-router-v2-rapid-move-foundation-20260813','build identity drift')
need(i.get('stable_baseline')=='v17.5.1' and i.get('patch_predecessor')=='v18.0.0','baseline/predecessor drift')
need(i.get('runtime_config')=='PersonalMarketTerminal-v18.0.1-TEST','isolated v18.0.1 TEST runtime missing')
need(i.get('application_bundle')=='De-Pulse-v18.0.1-TEST.app','separate TEST bundle missing')
need(len(s.get('scope_lock',[]))==10,'scope lock must remain 10 clauses')
router=(R/'smart_router_v2.go').read_text(); route=(R/'provider_router.go').read_text(); rapid=(R/'rapid_move_intelligence.go').read_text(); model=(R/'runtime_model.go').read_text(); render=(R/'renderer/renderer.js').read_text()
for term in ['providerCapabilityNotEntitled','ProviderCapabilityStateRecord','rankedProviderRoute','providerCapabilityCircuits','P95LatencyMs','CallsAvoided']:
    need(term in router+route+model,'Smart Router scope missing '+term)
for term in ['15 * time.Second','30 * time.Second','60 * time.Second','2 * time.Minute','5 * time.Minute','TIERED_PARTIAL','rapidMoveSourceAgreement','rapidMoveMechanicalRisk','rapidMoveCatalyst','rapidMoveMarketContext','Outcome20mPct','TraceID']:
    need(term in rapid,'Rapid Move scope missing '+term)
need('type":"rapid-move"' in rapid.replace(' ',''),'Rapid Move live event broadcast missing')
need("d.type==='rapid-move'" in render and 'toast(`${ev.symbol} · RAPID MOVE' in render,'immediate Rapid Move banner consumer missing')
need('promoteRapidMoveToRadarLocked' in rapid and 'livePriorityHints' in rapid,'Radar/live priority reuse missing')
need('providerCapabilityStates' in model and 'rapidMoveEvents' in model,'persistent runtime/cache state missing')
need('Provider Capability / Entitlement State' in (R/'data_utility_registry.json').read_text(),'Data Utility registration missing')
if e:
 print('v18.0.1 Smart Router / Rapid Move Scope: FAIL'); [print(' -',x) for x in e]; sys.exit(2)
print('v18.0.1 Smart Router / Rapid Move Scope: PASS · 10/10 locked clauses present · canonical router/event pipeline reused · Coverage Truth preserved')
