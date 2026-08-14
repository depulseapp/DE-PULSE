#!/usr/bin/env python3
"""v17 Major Closure performance/capacity evidence-map gate.
Heavy tests remain independent jobs and are not recursively executed here (RL-027/RL-028).
"""
from pathlib import Path
R=Path(__file__).resolve().parent
errs=[]
report=R/'renderer/qa/v17.5.0-performance-capacity-review.md'
if not report.exists(): errs.append('performance/capacity review missing')
else:
    txt=report.read_text(errors='ignore')
    for token in ['Critical-path acceptance budgets','selected-symbol freshness','Provider queue','Backpressure / load shedding','Persistence / warm start','298801 ns/op','Verdict']:
        if token not in txt: errs.append('performance report missing '+token)
for f in ['runtime_slo.go','runtime_observability.go','workload_controller.go','v17_3_performance_test.go','v16_11_performance_stability_gate.py','bounded_race_gate.py','randomized_order_gate.py']:
    if not (R/f).exists(): errs.append('required independent performance/reliability evidence owner missing: '+f)
if errs:
    print('v17.5 Performance / Capacity Closure Gate: FAIL')
    for e in errs: print(' -',e)
    raise SystemExit(2)
print('v17.5 Performance / Capacity Closure Gate: PASS · SLO/backpressure/persistence/long-run evidence mapped; heavy evidence independently blocking')
