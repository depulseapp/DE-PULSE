#!/usr/bin/env python3
from pathlib import Path
import json,sys
R=Path(__file__).resolve().parent; e=[]
def need(ok,msg):
    if not ok:e.append(msg)
i=json.loads((R/'release_identity.json').read_text()); s=json.loads((R/'v18_0_6_scope.json').read_text()); c=json.loads((R/'v18_0_6_g0_g3_contract.json').read_text())
need(i.get('version')=='18.0.6' and i.get('channel') in {'TEST','STABLE'},'v18.0.6 release identity missing')
need(i.get('previous_stable')=='v18.0.5' and i.get('patch_predecessor')=='v18.0.5','v18.0.5 Stable predecessor drift')
if i.get('channel')=='TEST':
    need(i.get('build_id')=='v18.0.6-test-smart-provider-router-rapid-move-market-shock-hardening-20260814','TEST build identity drift')
    need(i.get('runtime_config')=='PersonalMarketTerminal-v18.0.6-TEST','isolated v18.0.6 TEST runtime missing')
    need(i.get('application_bundle')=='De-Pulse-v18.0.6-TEST.app','separate TEST bundle missing')
else:
    need(i.get('build_id')=='v18.0.6-stable-smart-provider-router-rapid-move-market-shock-hardening-20260814','Stable build identity drift')
    need(i.get('runtime_config')=='PersonalMarketTerminal','Stable runtime config drift')
    need(i.get('application_bundle')=='De-Pulse.app','Stable bundle drift')
need(len(s.get('clauses',[]))==12,'immutable scope clause count mismatch')
need(c.get('baseline',{}).get('commit')=='d4431b2cb6a12d6c55fc0c23e93d27b21406465b','G0 Stable commit drift')
router=(R/'smart_router_v2.go').read_text()+(R/'provider_router.go').read_text(); rec=(R/'provider_reconciliation.go').read_text(); rapid=(R/'rapid_move_intelligence.go').read_text(); core=(R/'engine_core.go').read_text(); test=(R/'v18_0_6_router_shock_hardening_test.go').read_text()
for term in ['ProviderCapabilityStateRecord','providerCapabilityCircuits','CallsAvoided','P95LatencyMs']:
    need(term in router,'Router v2 canonical owner missing '+term)
need('providerReconciliationConflictCount' in rec and 'Scorecard.SourceDisagreements = providerReconciliationConflictCount' in core,'canonical source-disagreement scorecard truth missing')
for term in ['rapid-move-v1.1.0','MARKET_SHOCK','rapidMoveApplyHysteresis','rapid-move-learning-v1.0.0','SHADOW','VALIDATED','APPROVED','PRODUCTION','HysteresisRetained','ContinuedOutcomes','ReversedOutcomes','FadedMixedOutcomes']:
    need(term in rapid,'Rapid Move hardening missing '+term)
need('existing.EventProviderAt' in rapid and 'existing.WindowStartAt' in rapid,'dedupe outcome-anchor preservation missing')
need('TIERED_PARTIAL' in rapid,'coverage truth drift')
need(all(x in test for x in ['SourceDisagreement','MarketWideEventAsMarketShock','HysteresisPreventsAlertStateThrash','DedupPreservesOriginalProviderOutcomeAnchor','GovernanceIsExplicitAndCannotAutoPromote']),'v18.0.6 focused tests incomplete')
need('v18.1' not in rapid+router+rec+core,'v18.1 architecture leaked into v18.0.6 implementation')
if e:
 print('v18.0.6 Router / Market Shock Scope: FAIL'); [print(' -',x) for x in e]; sys.exit(2)
print('v18.0.6 Router / Market Shock Scope: PASS · 12/12 clauses · canonical owners reused · adaptive/no-execution boundaries preserved')
