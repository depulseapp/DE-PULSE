#!/usr/bin/env python3
from pathlib import Path
import json,sys
R=Path(__file__).resolve().parent; go='\n'.join(p.read_text(errors='ignore') for p in R.glob('*.go') if not p.name.endswith('_test.go')); js=(R/'renderer/renderer.js').read_text(); html=(R/'renderer/index.html').read_text(); tests=(R/'v16_1_test.go').read_text(); m=json.loads((R/'renderer/qa/v16.1.0-master-scope.json').read_text()); e=[]
req={'V16.1-07','V16.1-12','V16.1-13','V16.1-14','V16.1-15','V16.1-27','V16.1-UI-01'}
if m.get('count')!=7 or {x['id'] for x in m.get('requirements',[])}!=req:e.append('v16.1 scope identity/count mismatch')
for tok in ['marketStructureFor','marketTradeability','marketBreadthState','relativeStrengthFor','canonicalSymbolClassifications','benchmarkRegime','liquiditySlippageState','buildMarketIntelligenceSnapshot']:
    if tok not in go:e.append('implementation missing '+tok)
if html.count('data-page="market-intelligence"')!=1:e.append('Market Intelligence primary nav must exist exactly once')
for tok in ['renderMarketIntelligence','marketIntelligenceDashboardSummary','MARKET','Leadership','Cross-Asset / Alternative']:
    if tok not in js:e.append('surface missing '+tok)
for tok in ['TestV161MarketStructureUsesExactHorizons','TestV161TradeabilityRequiresSPYQQQVIX','TestV161BreadthUsesExplicitDenominatorAndStaleCannotVote','TestV161RelativeStrengthExactHorizonNoShortening','TestV161CanonicalClassificationAndIndustryRegime','TestV161LiquidityMappingUnknownIsNeverSafe']:
    if tok not in tests:e.append('regression missing '+tok)
if e: print('v16.1 Market Intelligence Scope Gate: FAIL'); print('\n'.join('- '+x for x in e)); sys.exit(1)
print('v16.1 Market Intelligence Scope Gate: PASS · 7/7')
