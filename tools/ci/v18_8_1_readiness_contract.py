#!/usr/bin/env python3
from __future__ import annotations
from pathlib import Path
import sys

ROOT=Path(__file__).resolve().parents[2]
PREP=ROOT/'preparation_catalyst.go'
TYPES=ROOT/'preparation_types_liquidity.go'
JOBS=ROOT/'runtime_jobs.go'

def fail(errors):
    print('DE.PULSE v18.8.1 readiness contract: FAIL', file=sys.stderr)
    for e in errors: print(f' - {e}', file=sys.stderr)
    return 1

def main():
    errors=[]
    prep=PREP.read_text(encoding='utf-8')
    types=TYPES.read_text(encoding='utf-8')
    jobs=JOBS.read_text(encoding='utf-8')
    required_prep=[
        'func (e *Engine) runMarketOpenPrep(reason string)',
        '"LIQUIDITY RISK"','"EXTENDED"','"EVENT RISK"',
        'materialSECFilingForTradingRisk(fi)',
        'readinessFreshnessGate(freshSnap.Freshness, []string{"Quotes", "VIX", "Intraday Bars", "News", "Earnings", "SEC Filings"}, now)',
        'e.setPreparationRich("market-open-prep", attention',
        'e.evaluateCatalystWatch(now)',
        'if marketOpenPrepWindow(now)',
        'e.runMarketOpenPrep("scheduled")',
        'mins > 9*60+25 && mins <= 10*60+15',
        'e.runMarketOpenPrep("missed-window catch-up")',
        'return "PREMARKET REACTION"','return "OPENING REACTION"','return "5m"','return "15m"','return "30m"','return "60m"','return "SESSION REACTION"','return "COMPLETE"'
    ]
    for token in required_prep:
        if token not in prep: errors.append(f'preparation semantic missing: {token}')
    required_types=[
        'type PreparationJobStatus struct','LastAttempt  int64','LastSuccess  int64','AttemptCount int','NextWindow   int64','TradingDay   string','Late         bool','Attention    string','Changed      []string','Exceptions   []CheckpointException',
        'type CatalystReactionState struct','GapPercent','RelativeVolume','SpreadPct','Liquidity','VWAPState','OpeningRangeState','HoldFadeState','VolatilityState','ReactionPercent','CompletedAt',
        'Window: "9:20–9:25 AM ET"','Window: "Event-driven only"'
    ]
    for token in required_types:
        if token not in types: errors.append(f'persistent readiness state missing: {token}')
    if 'e.sessionIntelligenceCoordinatorLoop(ctx, key, alpacaKey, alpacaSecret)' not in jobs:
        errors.append('session intelligence coordinator is not started')
    if 'e.autoFreshnessRecoveryLoop(ctx)' not in jobs:
        errors.append('freshness recovery support is not started')

    start=prep.find('func (e *Engine) runMarketOpenPrep(reason string)')
    end=prep.find('\nfunc (e *Engine) marketOpenPrepLoop', start)
    body=prep[start:end] if start>=0 and end>start else ''
    for forbidden in ('refreshQuotesRouted(', 'refreshHistoryRouted', 'refreshNews(', 'refreshEarnings', 'refreshFilings', 'refreshFundamentals', 'refreshMacroRouted', 'refreshOptions'):
        if forbidden in body: errors.append(f'Market Open Prep must reconcile current evidence, not broad-refetch: {forbidden}')
    for cached in ('quotes := clone(e.quotes)','bars := clone(e.bars)','earnings := clone(e.earnings)','filings := clone(e.filings)','news := clone(e.news)','options := clone(e.options)'):
        if cached not in body: errors.append(f'Market Open Prep cached-evidence reuse missing: {cached}')
    if errors: return fail(errors)
    print('DE.PULSE v18.8.1 readiness contract: PASS')
    print('Market Open Prep cached-evidence reconciliation / no broad refetch: PASS')
    print('scheduled + missed-window catch-up semantics: PASS')
    print('liquidity/EXTENDED/SEC-event-risk/freshness exceptions: PASS')
    print('persistent preparation + catalyst lifecycle measurements: PASS')
    return 0

if __name__=='__main__': raise SystemExit(main())
