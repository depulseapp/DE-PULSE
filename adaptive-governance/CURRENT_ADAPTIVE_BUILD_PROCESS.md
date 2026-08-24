# CURRENT Adaptive Build Process

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Retained completed process authority:** #73 / `ADAPT-ROOT-CONVERGENCE-001` / `adapt-root-convergence-001` / closure `governance/work-slices/ADAPT-ROOT-CONVERGENCE-001/closure.json`  
**Completed foundation:** #80 / `ADAPT-DATAHEALTH-BASELINE-001`  
**Completed Router adoption:** #81 / PR #87 / merge `1870dd3881dbe7f6463f242e35fdc19e70d9ae15`  
**Completed runtime Data Health:** #82 / Fast #894 / Qualified #187 / PR #88 / merge `4882b6d53c0c34463239faae752b86de377fb19a`  
**Active product work:** #83 / `ADAPT-PROVIDER-LIFECYCLE-001` / `adapt-provider-lifecycle-001`  
**Parent program:** #79 / `ADAPT-PROVIDER-PRODUCTION-001`

For #83 and remaining Data Health slices, implementation is executable-first and canonical-owner-first:
- reuse Smart Provider Router v2 as the sole general market-data routing/admission owner;
- reuse canonical freshness, cache, persistence, ProviderTelemetry, Router v2 capability state and validation owners before adding anything;
- REUSE/CONSOLIDATE/REFACTOR before ADD, and fail closed on a parallel provider-specific lifecycle/health subsystem;
- lifecycle promotion is explicit governed-only; runtime may compute readiness but must never auto-promote;
- automatic circuit suppression, rate-limit cooldown, probing, fallback and recovery remain runtime reliability behavior and do not mutate governed lifecycle;
- insufficient observations, unstable auth/errors, stale evidence, schema/semantic uncertainty, poor latency/quota behavior, unproven fallback, excessive independent disagreement, weak provenance, unproven utility or truth-boundary violations fail readiness closed;
- direct SEC/EDGAR and other explicit first-party authorities remain authority rules, not rank-promotion candidates;
- TradeInsight uses the shared lifecycle/readiness evaluator and remains SHADOW until #78 evidence and explicit promotion;
- use focused regressions before canonical Fast exact-head PASS, then deliberate impact-selected Qualified exact-head PASS;
- do not add workflow families or weaken G0–G16/source-health/architecture gates.

The governed program sequence token remains `#81/#82/#83/#78/#84`; #81 and #82 are complete, so current execution is **#83 → #78 → #84**. Documentation alone never verifies closure.
