#!/usr/bin/env python3
from pathlib import Path
import json,subprocess,sys
R=Path(__file__).resolve().parent
go='\n'.join(p.read_text(errors='ignore') for p in R.glob('*.go') if not p.name.endswith('_test.go'))
js=(R/'renderer/renderer.js').read_text(); html=(R/'renderer/index.html').read_text(); tests=(R/'v16_2_test.go').read_text(); m=json.loads((R/'renderer/qa/v16.2.0-master-scope.json').read_text()); research=(R/'research_providers.go').read_text(); macro=(R/'macro_events.go').read_text(); e=[]
req={'V16.2-01','V16.2-02','V16.2-03','V16.2-04','V16.2-21','V16.2-28','V16.2-29'}
if {x['id'] for x in m.get('scope',[])}!=req:e.append('v16.2 scope identity/count mismatch')
for tok in ['buildEventNewsIntelligence','buildEconomicCalendar','buildFedIntelligence','buildSmartNotifications','buildReactionIntelligence','buildEventDecisionCorrelation','buildEventIntelligenceSnapshot']:
    if tok not in go:e.append('implementation missing '+tok)
for tok in ['marketEventIntelligenceMarkup','researchEventIntelligenceMarkup','eventDecisionRiskForSymbol','eventOperationalStatusMarkup','eventIntelligence:runtime?.eventIntelligence']:
    if tok not in js:e.append('integration missing '+tok)
if 'Reserved for v16.2' in js:e.append('Market Intelligence EVENTS placeholder was not activated')
if 'data-page="alerts"' in html.lower() or 'data-page="alerts"' in js.lower():e.append('retired standalone Alerts surface returned')
# v16.2 derived Event Intelligence must propagate immediately when its canonical
# News/Earnings/Macro stores commit.  This is a presentation/invalidation
# invariant only; it must never become a second fetch or scheduler path.
for owner,src,marker in [
    ('News',research,'func (e *Engine) refreshNews'),
    ('Earnings',research,'func (e *Engine) refreshEarnings'),
    ('Macro',macro,'func (e *Engine) refreshMacroEvents'),
]:
    start=src.find(marker)
    if start<0:
        e.append(owner+' canonical refresh owner missing')
        continue
    next_func=src.find('\nfunc ', start+len(marker))
    body=src[start: next_func if next_func>=0 else len(src)]
    if 'map[string]any{"type": "runtime", "runtime": e.Snapshot()}' not in body:
        e.append(owner+' refresh does not publish derived Event Intelligence runtime snapshot')
for tok in ['TestV162NewsClustersDuplicateSourcesAndRanksMateriality','TestV162MacroSurpriseRequiresActualAndForecast','TestV162FedIntelligenceUsesSourcedTimeline','TestV162SmartNotificationsAreEventTriggeredNotConditionCards','TestV162ReactionIntelligenceReusesCatalystAndMacroEvidence','TestV162EventDecisionCorrelatesWithoutScoreMutation']:
    if tok not in tests:e.append('regression missing '+tok)
q=subprocess.run(['go','test','-count=1','-run','TestV162','./...'],cwd=R,text=True,capture_output=True)
if q.returncode:e.append('v16.2 regressions failed: '+(q.stdout+q.stderr)[-900:])
if e:
    print('v16.2 Professional Event Intelligence Scope Gate: FAIL');print('\n'.join('- '+x for x in e));sys.exit(1)
print('v16.2 Professional Event Intelligence Scope Gate: PASS · 7/7')
