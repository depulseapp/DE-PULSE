#!/usr/bin/env python3
from __future__ import annotations
import json, sys
from pathlib import Path
R=Path(__file__).resolve().parent; errors=[]
def need(ok,msg):
    if not ok: errors.append(msg)
identity=json.loads((R/'release_identity.json').read_text()); slices=json.loads((R/'v17_delivery_slices.json').read_text())
version=str(identity.get('version',''))
is_native_v17=version.startswith('17.') and identity.get('channel') in {'TEST','RC','STABLE'}
is_inherited_v18=version.startswith('18.') and identity.get('channel') in {'TEST','RC','STABLE'} and identity.get('stable_baseline')=='v17.5.1'
need(is_native_v17 or is_inherited_v18,'inherited v17 scope release identity missing')
if is_native_v17: need(identity.get('stable_baseline')=='v16.11.0','v16.11 Stable baseline drifted')
v172=next((x for x in slices.get('slices',[]) if x.get('id')=='v17.2'),None)
need(v172 is not None and v172.get('status') in {'SOURCE TEST CHECKPOINT','TEST QUALIFIED','SOURCE QUALIFIED','COMPLETE'},'v17.2 slice must be completed/preserved before later v17.x qualification')
pipe=(R/'canonical_pipeline.go').read_text()
for token in ['CanonicalPipelineDiagnostics','propagateCanonicalQuoteChange','quoteMateriallyChanged','SuppressedDownstream','PersistenceEnqueues','CatalystEvaluations','catalystQuoteReactionActiveForSymbol','updateCanonicalSessionClose']:
    need(token in pipe,f'canonical propagation missing: {token}')
need('CanonicalPipeline' in (R/'runtime_load_profile.go').read_text(),'canonical diagnostics not exposed')
live=(R/'runtime_live_loop.go').read_text()
need('e.propagateCanonicalQuoteChange(symbol, prev, q)' in live,'live quote path bypasses canonical propagation')
need('alloc := e.multiFeedAllocation()\n\te.mu.Lock()' in live,'IEX allocation must occur before engine write lock')
need('mergeAlpacaIEXStreamAt' in live,'IEX deterministic deadlock test seam missing')
need((R/'quote_providers.go').read_text().count('propagateCanonicalQuoteChange')>=2,'snapshot providers bypass canonical propagation')
hist=(R/'history_providers.go').read_text()+(R/'fallback_providers.go').read_text()
need(hist.count('updateCanonicalSessionClose')>=3,'history session-close enrichment bypasses canonical propagation')
intel=(R/'persistence_intelligence.go').read_text()
need('decisionPayload' in intel and 'decisionRaw, err := json.Marshal(s)' not in intel,'decision lineage still carries mutable outcomes')
for src,label in [((R/'persistence_backend_sqlite.go').read_text(),'native SQLite'),((R/'persistence_backend_windows.go').read_text(),'Windows SQLite')]:
    need('ON CONFLICT(evidence_id) DO NOTHING' in src,f'{label} evidence immutability missing')
    need('ON CONFLICT(decision_id) DO NOTHING' in src,f'{label} decision immutability missing')
    need('derived_features.source_hash<>excluded.source_hash' in src and 'derived_features.as_of_ms<excluded.as_of_ms' not in src,f'{label} feature store not source-hash incremental')
fb=(R/'persistence_backend_fallback.go').read_text()
need('if _, exists := b.data.Evidence[r.ID]; !exists' in fb,'fallback evidence immutability missing')
need('if _, exists := b.data.Decisions[r.ID]; !exists' in fb,'fallback decision immutability missing')
need('old.SourceHash != r.SourceHash' in fb and 'old.AsOf < r.AsOf' not in fb,'fallback feature store not source-hash incremental')
tests=(R/'v17_2_canonical_pipeline_test.go').read_text()
for name in ['TestV172ImmaterialTicksUpdateMemoryButSuppressHeavyDownstream','TestV172MaterialPriceAndTruthStateChangesPropagate','TestV172CatalystQuotePropagationIsSymbolScoped','TestV172SnapshotProvidersUseCanonicalPropagationOwner','TestV172DecisionLineagePayloadIsFrozenFromLaterOutcomes','TestV172DerivedFeatureHashIgnoresOutcomeOnlyChanges','TestV172AlpacaIEXStreamDoesNotSelfDeadlockOnAllocationRead']:
    need(name in tests,f'acceptance test missing: {name}')
if errors:
    print('v17.2 scope gate: FAIL'); [print(' -',e) for e in errors]; sys.exit(1)
print('v17.2 scope gate: PASS · canonical propagation, immutable lineage, incremental features and IEX deadlock regression present')
