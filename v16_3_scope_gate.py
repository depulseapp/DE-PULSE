#!/usr/bin/env python3
from pathlib import Path
import json, subprocess, sys

R = Path(__file__).resolve().parent
GO = '\n'.join(p.read_text(errors='ignore') for p in R.glob('*.go') if not p.name.endswith('_test.go'))
JS = (R/'renderer/renderer.js').read_text(errors='ignore')
HTML = (R/'renderer/index.html').read_text(errors='ignore')
T1 = (R/'v16_3_test.go').read_text(errors='ignore')
TP = (R/'v16_3_professional_test.go').read_text(errors='ignore')
TR = (R/'v16_3_renderer_test.js').read_text(errors='ignore')
MASTER = json.loads((R/'renderer/qa/v16.3.0-master-scope.json').read_text())
HISTORY = (R/'history_providers.go').read_text(errors='ignore')
ERR=[]

expected={8,18,19,20,26}
actual={int(x['id']) for x in MASTER.get('scope_lock',[])}
if actual != expected:
    ERR.append(f'v16.3 scope identity/count mismatch: expected {sorted(expected)}, got {sorted(actual)}')
if MASTER.get('baseline',{}).get('source_fingerprint') != 'da718a33cab7fd177644cb0a84a427d84856a373fe956816b0c945f3cb103b00':
    ERR.append('authoritative v16.2 baseline fingerprint changed')

# Canonical implementation: analytics layer reuses existing signal validation, bars,
# corporate action truth and the same deterministic renderer core.
for tok in [
    'evaluateSignalSnapshotsWithActions',
    'postSnapshotSplitAdjustment',
    'buildSeasonalitySnapshot',
    'buildCalibrationSnapshot',
    'buildCorrelationConcentrationSnapshot',
    'buildValidationLearningSnapshot',
    'ValidationLearningSnapshot',
]:
    if tok not in GO: ERR.append('implementation missing '+tok)
for tok in [
    'deterministicScoreFromFamilies',
    'replayRuntimeAt',
    'computePlanAt',
    'replaySignalSnapshot',
    'researchValidationLearningMarkup',
    'marketValidationLearningMarkup',
    'validationLearningMaintenanceMarkup',
    'concentrationAttentionForSymbol',
    'opportunity quality unchanged',
]:
    if tok not in JS: ERR.append('renderer/integration missing '+tok)

# #8 must get deeper SPY/QQQ history from the existing canonical history owner,
# not a second history store/provider path.
for tok in ['SPY', 'QQQ', '10']:
    if tok not in HISTORY: ERR.append('seasonality history ownership missing '+tok)
if 'daily-raw' not in JS or '_v163ReplayTruth' not in JS:
    ERR.append('corporate-action cutoff-safe replay path missing')

# No forbidden surface/scope leakage.
low_html=HTML.lower(); low_js=JS.lower()
if 'data-page="alerts"' in low_html or 'data-page="alerts"' in low_js:
    ERR.append('retired standalone Alerts surface returned')
for forbidden in ['portfolio', 'positions/p&l', 'brokerage execution']:
    # Those words may be present in explicit explanatory exclusions. We only block new page nav.
    if f'data-page="{forbidden}"' in low_html:
        ERR.append('forbidden top-level surface returned: '+forbidden)

# Professional blocking regression inventory.
regressions=[
    'TestV163OutcomeTargetBeforeEntryNeverCountsAsSuccess',
    'TestV163OutcomeInvalidationBeforeTarget',
    'TestV163OutcomeEntryThenTarget',
    'TestV163OutcomeEntryNeverTouched',
    'TestV163OutcomeMissingMaturedHistoryUnavailable',
    'TestV163OutcomeCoarseBarOrderingAmbiguous',
    'TestV163FrozenReplayFieldsSurviveAPIRoundTrip',
    'TestV163SeasonalityInsufficientSampleDoesNotFabricateSignal',
    'TestV163CalibrationTinySampleAndScoreNotProbability',
    'TestV163CorrelationRequiresAlignedReturnWindow',
    'TestV163SemiconductorConcentrationIsAttentionContextOnly',
    'TestV163ProfessionalSplitAdjustmentPreservesFrozenOutcomeEconomics',
    'TestV163ProfessionalWeekendGapDoesNotDelayCompletedDailyEvidence',
]
alltests=T1+'\n'+TP
for tok in regressions:
    if tok not in alltests: ERR.append('professional regression missing '+tok)
for tok in ['date-only SEC filing', 'RAW_HISTORY', 'PARTIAL', 'restor']:
    if tok.lower() not in TR.lower(): ERR.append('renderer replay regression evidence missing '+tok)

# Execute the focused affected-gate loops. Deterministic equivalence is a hard
# invariant, so it is part of this scope gate even though the formulas are unchanged.
commands=[
    (['go','test','-count=1','-run','TestV163','./...'], 'Go v16.3 regressions'),
    (['node','v16_3_renderer_test.js'], 'renderer replay/isolation'),
    (['node','deterministic_equivalence_test.js'], 'deterministic 2403/2403'),
]
for cmd,label in commands:
    q=subprocess.run(cmd,cwd=R,text=True,capture_output=True)
    if q.returncode:
        ERR.append(label+' failed: '+(q.stdout+q.stderr)[-1200:])

if ERR:
    print('v16.3 Professional Validation & Learning Scope Gate: FAIL')
    for x in ERR: print('- '+x)
    sys.exit(1)
print('v16.3 Professional Validation & Learning Scope Gate: PASS · 5/5 (#8/#18/#19/#20/#26)')
