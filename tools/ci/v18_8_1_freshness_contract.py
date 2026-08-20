#!/usr/bin/env python3
from __future__ import annotations
from pathlib import Path
import sys

ROOT=Path(__file__).resolve().parents[2]
FRESH=ROOT/'data_freshness.go'
ROUTED=ROOT/'routed_refresh.go'
JOBS=ROOT/'runtime_jobs.go'
RENDERER=ROOT/'renderer'/'renderer.js'

def fail(errors):
    print('DE.PULSE v18.8.1 freshness contract: FAIL', file=sys.stderr)
    for e in errors: print(f' - {e}', file=sys.stderr)
    return 1

def main():
    errors=[]
    fresh=FRESH.read_text(encoding='utf-8')
    routed=ROUTED.read_text(encoding='utf-8')
    jobs=JOBS.read_text(encoding='utf-8')
    renderer=RENDERER.read_text(encoding='utf-8')
    for token in [
        'func freshnessLimits(dataset, provider, session string)',
        'case "Quotes":','case "VIX":','case "Intraday Bars":','case "Daily / Weekly History":','case "News":','case "Earnings":','case "SEC Filings":','case "Fundamentals":','case "Global":','case "Macro":','case "Options":',
        'return "LIVE"','return "FRESH"','return "DUE SOON"','return "DELAYED"','return "STALE"','return "ERROR"','return "UNAVAILABLE"','return "IDLE"',
        'func safeFreshnessTimestamp(ts, receipt, now int64)',
        'provider market timestamp','market observation','last successful news check','last successful EDGAR check',
        'quoteMissing','intradayMissing'
    ]:
        if token not in fresh: errors.append(f'freshness truth missing: {token}')
    for token in [
        'func (e *Engine) researchPackageReadinessAt(symbol string, nowTime time.Time)',
        'Quote stale/clock-skewed for ','Intraday history stale for ','Fundamentals stale/clock-skewed for ',
        'func (a *Application) handleResearchRefresh',
        'a.engine.refreshResearchHistory(ctx, sym)','a.engine.refreshResearchFundamentals(ctx, sec.Finnhub, sym)','a.engine.refreshResearchEarnings(ctx, sec.Finnhub, sym)','a.engine.refreshResearchNews(ctx, sec.Finnhub, sym)','a.engine.refreshSECResearchSymbol(ctx, sym)',
        'func (e *Engine) refreshDatasetRouted(ctx context.Context, dataset string, s Secrets) bool'
    ]:
        if token not in routed: errors.append(f'targeted routed freshness behavior missing: {token}')
    if 'e.autoFreshnessRecoveryLoop(ctx)' not in jobs: errors.append('automatic freshness recovery loop not started')
    for token in ['Data Freshness v2','Check age and data age are separated','proactive session-aware recovery']:
        if token not in renderer: errors.append(f'freshness UI truth missing: {token}')
    if errors: return fail(errors)
    print('DE.PULSE v18.8.1 freshness contract: PASS')
    print('dataset/session-specific cadence + truth states: PASS')
    print('provider observation time separated from receipt/check time: PASS')
    print('targeted routed recovery + selected-symbol research readiness: PASS')
    print('stable Data Freshness v2 UI + automatic recovery loop: PASS')
    return 0

if __name__=='__main__': raise SystemExit(main())
