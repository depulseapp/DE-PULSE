#!/usr/bin/env python3
from __future__ import annotations
import json, re, sys
from pathlib import Path
R=Path(__file__).resolve().parent
errors=[]
need=lambda ok,msg: errors.append(msg) if not ok else None
major=json.loads((R/'v17_major_scope.json').read_text())
slices=json.loads((R/'v17_delivery_slices.json').read_text())
need(major.get('scope_frozen') is True,'v17 major scope is not frozen')
need(len(major.get('items',[]))==20,'v17 major scope must retain exactly 20 approved G1 items')
v170=next((x for x in slices.get('slices',[]) if x.get('id')=='v17.0'),None)
need(v170 is not None and (v170.get('status')=='IN DEVELOPMENT' or str(v170.get('status','')).startswith('SOURCE TEST CHECKPOINT')),'v17.0 foundation slice missing/inactive')
for f in ['persistence_repository.go','persistence_backend_sqlite.go','persistence_backend_fallback.go','runtime_load_profile.go','runtime_observability.go','runtime_slo.go','workload_controller.go','runtime_degradation.go','v17_0_persistence_test.go','v17_0_degradation_test.go']:
    need((R/f).exists(),f'missing v17.0 evidence file: {f}')
repo=(R/'persistence_repository.go').read_text()
sql=(R/'persistence_backend_sqlite.go').read_text()
engine=(R/'engine_core.go').read_text()
need('type PersistenceBackend interface' in repo and 'LoadSymbols(' in repo and 'Stats(' in repo,'repository contract incomplete')
need('schema_migrations' in sql and 'migrationApplied' in sql,'SQLite schema-versioned migration contract missing')
need('symbol_registry' in sql and 'first_seen_ms' in sql,'Global Symbol Registry persistence missing')
need('canonical_quotes' in sql and 'quote_history' in sql,'canonical quote/history persistence missing')
need('loadPersistedCanonicalQuotes' in engine and 'candidate.DataState = "persisted"' in engine,'warm-start truth contract missing')
need('quoteMateriallyChanged' in repo and 'MaterialWritesSuppressed' in repo,'material-change write suppression missing')
work=(R/'workload_controller.go').read_text(); obs=(R/'runtime_observability.go').read_text(); deg=(R/'runtime_degradation.go').read_text()
need('provider-rest' in work and 'scanner' in work and 'background' in work,'bounded workload classes missing')
need('ProviderRequestDiagnostics' in obs and 'providerGetJSON' in obs and 'RequestTelemetry' in obs,'provider/API telemetry missing')
slo=(R/'runtime_slo.go').read_text(); need('RuntimeSLOAssessment' in slo and 'Interactive API p95' in slo and 'Critical decision datasets' in slo,'runtime SLO acceptance missing')
for reason in ['LOCAL LOAD','RATE LIMITED','PROVIDER DEGRADED','PARTIAL COVERAGE','NETWORK']:
    need(reason in deg,f'degradation reason missing: {reason}')
need('evidence_records' in sql and 'decision_lineage' in sql and 'outcome_history' in sql and 'derived_features' in sql,'structured evidence/lineage/outcome/feature persistence missing')
load=(R/'runtime_load_profile.go').read_text(); need('ProviderCallsAvoided' in load and 'ProviderRequests' in load and 'Persistence' in load,'runtime load profile incomplete')

registry=json.loads((R/'data_utility_registry.json').read_text())
registered={str(x.get('dataset','')).strip() for x in registry.get('datasets',[])}
for dataset in ['Global Symbol Registry','Persistent Canonical Quotes','Persistent Quote History','Evidence Records','Decision Lineage','Outcome History','Derived Feature Store','Runtime Load / Provider Telemetry']:
    need(dataset in registered,f'v17 data-utility dataset not registered: {dataset}')
# permanent boundaries must remain represented in source/docs
alltext='\n'.join((R/f).read_text(errors='ignore') for f in ['renderer/docs/developer.md','renderer/docs/limitations.md','README.md'])
need('No Execution' in alltext or 'no execution' in alltext.lower(),'No Execution boundary documentation missing')
if errors:
    print('v17.0 scope gate: FAIL')
    for e in errors: print(' -',e)
    sys.exit(1)
print('v17.0 scope gate: PASS · frozen v17 major scope retained · foundation slice evidence present')
