#!/usr/bin/env python3
import json
from pathlib import Path
R=Path(__file__).resolve().parents[2]
p=json.loads((R/'governance'/'policies'/'data_health_policy.json').read_text())
errs=[]
policy=p.get('policy',{})
for k in ('session_aware','provider_vs_cache_timestamps','stale_while_revalidate','fallback_recovery','material_change_priority','selected_symbol_priority','adaptive_policy_changes_require_shadow_validation'):
    if policy.get(k) is not True: errs.append(k)
if policy.get('market_critical_priority')[:2] != ['SPY','QQQ']: errs.append('SPY/QQQ market-critical priority')
src=(R/'data_freshness.go').read_text()
for token in ('ProviderTimestamp','CacheAt','CheckAgeMs','DataAgeMs','FreshLimitMs','StaleLimitMs','Fallback','Reason'):
    if token not in src: errs.append('freshness diagnostic '+token)
live=(R/'live_subscription_manager.go').read_text()
if 'marketCriticalLiveSymbols = []string{"SPY", "QQQ"}' not in live: errs.append('live priority source')
if errs:
    print('Adaptive Data Health: FAIL · '+', '.join(errs)); raise SystemExit(2)
print('Adaptive Data Health: PASS · session-aware freshness/cache/fallback observability + SPY/QQQ critical priority')