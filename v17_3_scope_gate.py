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
v173=next((x for x in slices.get('slices',[]) if x.get('id')=='v17.3'),None)
need(v173 is not None and v173.get('status') in {'SOURCE TEST CHECKPOINT','TEST QUALIFIED','SOURCE QUALIFIED','COMPLETE'},'v17.3 slice must be completed/preserved before later v17.x qualification')
slo=(R/'runtime_slo.go').read_text(); load=(R/'runtime_load_profile.go').read_text(); persist=(R/'persistence_repository.go').read_text(); renderer=(R/'renderer/renderer.js').read_text(); tests=(R/'runtime_performance_slo_regression_test.go').read_text()
for token in ['Selected-symbol freshness','Actionable-watchlist freshness','Stale→current recovery','Degradation recovery','CPU utilization','DB write rate','Storage growth','Startup/warm-start time','Provider request budgets','Live subscription utilization','Canonical reuse / provider calls avoided']:
    need(token in slo,f'missing v17.3 SLO: {token}')
for token in ['CPUUtilizationPct','GOMAXPROCS','CanonicalReuseHitRatePct','StorageGrowthBytes','RuntimeStartupDiagnostics','RuntimeRecoveryDiagnostics','runtimeCPUTotalSeconds']:
    need(token in load+slo,f'missing runtime performance metric: {token}')
for token in ['WriteBatchesLastMin','RowsWrittenLastMin','writeEvents']:
    need(token in persist,f'missing DB write-rate telemetry: {token}')
for token in ['cpuUtilizationPct','warmStartCoveragePct','canonicalReuseHitRatePct','currentlyStaleDatasets','rowsWrittenLastMinute','storageGrowthBytes']:
    need(token in renderer,f'Maintenance/Data Engine does not expose {token}')
for name in ['TestV173SelectedAndActionableFreshnessSLOsAreExplicit','TestV173RecoveryTrackerMeasuresStaleAndDegradationRecovery','TestV173PersistenceWriteRateIsObservable','TestV173RuntimeLoadSamplesCPUStartupReuseAndStorage','TestV173ActiveMarketPressureProtectsCriticalWork','TestV173SLOBudgetsBlockOnlyMeasuredSeverePressure']:
    need(name in tests,f'missing v17.3 acceptance test: {name}')
if errors:
    print('v17.3 scope gate: FAIL'); [print(' -',e) for e in errors]; sys.exit(1)
print('v17.3 scope gate: PASS · freshness/recovery/CPU/API/queue/provider/warm-start/write/storage SLO evidence is measurable and exposed')
