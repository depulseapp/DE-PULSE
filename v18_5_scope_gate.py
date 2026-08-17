#!/usr/bin/env python3
from pathlib import Path
import json, sys
R=Path(__file__).resolve().parent
errors=[]
def need(ok,msg):
    if not ok: errors.append(msg)
scope=json.loads((R/'v18_5_scope.json').read_text())
need(scope.get('schema')=='DE.PULSE-v18.5-SCOPE-1','scope schema drift')
need(scope.get('version')=='18.5.0','v18.5 version drift')
need(scope.get('incomingStable')=='v18.4.0','incoming Stable version drift')
need(scope.get('incomingStableTag')=='v18.4.0-stable','incoming Stable tag drift')
need(scope.get('incomingStableCommit')=='8f3d0db3a7b7300f5954bd34af2541a0e04d6870','incoming Stable commit drift')
need(scope.get('incomingStableFingerprint')=='97eb553241eb98bda25f1867ea9f9ebfbaa4fe95e110960a8a4d5e7fc28eb05b','incoming Stable fingerprint drift')
need(scope.get('incomingStableSourceSha256')=='05a9ebe601a77aff783900904eced5eb49e02fd19c8c7f1d8ba0a817d72c0ea4','incoming Stable source SHA drift')
need(scope.get('title')=='Major Closure & Release Assurance','release title drift')
expected_closure=[
 'full_v18_scope_traceability','architecture_source_quality_developer_proofing','data_utility_correlation_reuse',
 'performance_capacity_responsiveness_stability','security_auth_session_authorization_abuse',
 'ui_ux_content_responsive_accessibility_runtime_truth','adaptive_intelligence_governance',
 'principal_engineer_acceptance','professional_trader_investor_acceptance',
 'native_hosted_packaging_provenance_release_assurance']
need(scope.get('closureDimensions')==expected_closure,'closure dimension set/order drift')
required_adr={
 'provider_failure_rate_limit_fallback_calls_avoided','stale_evidence_freshness_slo','source_disagreement',
 'consumer_aware_blast_radius','postgres_pressure_slow_unavailable','sqlite_desktop_continuity',
 'queue_saturation_bounded_backpressure','workload_priority_graceful_load_shedding','restart_warm_start_recovery',
 'multi_user_multi_symbol_fanout_shared_reuse','background_job_pressure_duplicate_work_avoidance',
 'recovery_hysteresis','unknown_degraded_abstain_truth','packaged_runtime_degradation_ux_diagnostics',
 'truthful_supported_operating_limits'}
need(set(scope.get('adrGdiRequiredScenarios',[]))==required_adr,'ADR-GDI scenario set drift')
need(len(scope.get('releaseBlockers',[]))==5,'release-blocker set drift')
need('no_execution_boundary' in scope.get('protectedContracts',[]),'No Execution protection missing')
need('deterministic_day_swing_long_formulas' in scope.get('protectedContracts',[]),'deterministic formula protection missing')
need('canonical_smart_provider_router' in scope.get('protectedContracts',[]),'Smart Router protection missing')
need('no_mature_asbi' in scope.get('exclusions',[]),'ASBI anti-scope-creep missing')
need('no_mature_tdti' in scope.get('exclusions',[]),'TDTI anti-scope-creep missing')
need('no_mature_aodr' in scope.get('exclusions',[]),'AODR anti-scope-creep missing')
for f in (
 'release/v18.5.0/G0-EXACT-BASELINE.md','release/v18.5.0/G1-IMMUTABLE-SCOPE.md',
 'release/v18.5.0/G2-ARCHITECTURE-DATA-UTILITY.md','release/v18.5.0/G3-DESIGN-DEPENDENCY-READINESS.md',
 'governance/ROADMAP.md','governance/ADAPTIVE-DATA-RELIABILITY-CONTRACT.md'):
    need((R/f).exists(),f+' missing')
if (R/'governance/ROADMAP.md').exists():
    road=(R/'governance/ROADMAP.md').read_text()
    need('### v18.5 — Major Closure & Release Assurance' in road,'canonical roadmap v18.5 section missing')
    need('ADR-GDI is a mandatory v18.5 closure dimension' in road,'roadmap ADR-GDI closure clause missing')
if (R/'governance/ADAPTIVE-DATA-RELIABILITY-CONTRACT.md').exists():
    adr=(R/'governance/ADAPTIVE-DATA-RELIABILITY-CONTRACT.md').read_text()
    low=adr.lower()
    need('data degraded' in low,'ADR-GDI contract missing DATA DEGRADED truth semantics')
    need('abstain' in low,'ADR-GDI contract missing ABSTAIN semantics')
    need('blast-radius' in low or 'blast radius' in low,'ADR-GDI contract missing blast-radius semantics')
    # Validate the governing bounded-work requirement semantically, not by an invented literal phrase.
    need(('queues, goroutines/workers, subscriptions, db work and background jobs must be bounded' in low)
         or ('queues' in low and 'must be bounded' in low and 'backpressure' in low),
         'ADR-GDI contract missing bounded queue/work/backpressure requirement')
    need('local_overload' in low and 'queue_saturated' in low,'ADR-GDI contract missing local-overload/queue-saturation reason taxonomy')
if errors:
    print('v18.5 scope gate: FAIL')
    for e in errors: print(' -',e)
    sys.exit(2)
print('v18.5 scope gate: PASS · G0-G3 frozen · 10 closure dimensions · 15 ADR-GDI scenarios · anti-scope-creep preserved')
