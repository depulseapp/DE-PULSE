#!/usr/bin/env python3
from __future__ import annotations
import json, sys
from pathlib import Path
R=Path(__file__).resolve().parent
errors=[]
def need(ok,msg):
    if not ok: errors.append(msg)
identity=json.loads((R/'release_identity.json').read_text())
slices=json.loads((R/'v17_delivery_slices.json').read_text())
version=str(identity.get('version','')).split('.')
try:
    vmajor, vminor = int(version[0]), int(version[1])
except (ValueError, IndexError):
    vmajor, vminor = -1, -1
is_native_v17=vmajor==17 and vminor>=1 and identity.get('channel') in {'TEST','RC','STABLE'}
is_inherited_v18=vmajor>=18 and identity.get('channel') in {'TEST','RC','STABLE'} and identity.get('stable_baseline')=='v17.5.1'
need(is_native_v17 or is_inherited_v18,'v17.1 inherited scope release identity missing')
if is_native_v17: need(identity.get('stable_baseline')=='v16.11.0','v16.11 Stable baseline drifted')
v171=next((x for x in slices.get('slices',[]) if x.get('id')=='v17.1'),None)
need(v171 is not None and v171.get('status') in {'IN DEVELOPMENT','SOURCE TEST CHECKPOINT','TEST QUALIFIED'},'v17.1 delivery slice missing/inactive')
work=(R/'workload_controller.go').read_text()
for token in ['WorkTierMarketCritical','WorkTierUserActionable','WorkTierRadarPromoted','WorkTierBroadDiscovery','WorkTierBackground','maxQueue','reservedCritical','AcquireTier','TryAcquireTier','ShouldShed','QueuedByTier','InFlightByTier','RejectedByTier']:
    need(token in work,f'workload/backpressure contract missing: {token}')
need('newWorkClass("provider-rest", 6, 24, 2)' in work,'provider REST capacity/queue/critical reserve contract drifted')
obs=(R/'runtime_observability.go').read_text()
need('providerGetJSONTier' in obs and 'workTierFromContext' in obs,'tier-aware provider request owner missing')
for token in ['BudgetPerMinute','BudgetRemaining','BudgetUtilizationPct','BudgetState','CooldownUntil','BudgetShed','providerBudgetPerMinute','Allow(provider string, tier WorkTier)']:
    need(token in obs,f'provider request-budget contract missing: {token}')
scanner=(R/'scanner.go').read_text()
radar=(R/'opportunity_radar.go').read_text()
jobs=(R/'runtime_jobs.go').read_text()
need('WorkTierBroadDiscovery' in scanner,'Broad Discovery is not routed through Tier 3')
need('WorkTierRadarPromoted' in radar,'Opportunity Radar is not routed through Tier 2')
need('WorkTierBackground' in jobs and 'withWorkTier(ctx, WorkTierBackground)' in jobs,'background isolation/context propagation missing')
load=(R/'runtime_load_profile.go').read_text()
for token in ['LiveSubscriptionBudgetDiagnostics','ReservedCapacity','ReserveUsed','UtilizationPct','Saturated','LiveSubscriptions']:
    need(token in load,f'live subscription budget telemetry missing: {token}')
deg=(R/'runtime_degradation.go').read_text()
need('LIVE CAPACITY SATURATED' in deg,'live subscription saturation reason code missing')
tests=(R/'v17_1_backpressure_test.go').read_text()
for name in ['TestV171ReservedCapacityProtectsTierZeroAndOne','TestV171ProviderQueueIsHardBounded','TestV171LowPriorityShedsBeforeActionable','TestV171LiveSubscriptionBudgetsExposeReservedHeadroom','TestV171ProviderRateLimitCooldownShedsLowTierButProtectsActionable','TestV171FinnhubBudgetComesFromLocalPacingNotGuessedEntitlement']:
    need(name in tests,f'v17.1 acceptance test missing: {name}')
if errors:
    print('v17.1 scope gate: FAIL')
    for e in errors: print(' -',e)
    sys.exit(1)
print('v17.1 scope gate: PASS · Tier 0-4 priority, bounded queues, critical reserve, load shedding and live-capacity telemetry present')
