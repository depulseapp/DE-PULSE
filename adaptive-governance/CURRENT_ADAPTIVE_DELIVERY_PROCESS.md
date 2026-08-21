# DE.PULSE — Current Adaptive Delivery Process

**Certified Stable:** `v18.9.0-stable`  
**Stable candidate:** `9ea81cddae4875ae15d3719ca028519a36c597b6`  
**Stable fingerprint:** `a8719090c341c874dbd1279cc31ad98e84075d5701c46a800bf951340780ecb9`  
**Master corrective program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Immediate patch:** `v18.9.1` / #64 / `ADAPT-RUNTIME-CRASH-001`.

## v18.9.0 delivery — COMPLETE / IMMUTABLE

Fast #481, Qualified #153 and Release #32 are the release authority for v18.9.0. G11-G16 passed, including actual macOS Apple Silicon + Windows x64 packaged-runtime audits, G15 assurance, no-rebuild Stable publication and G16 handoff evidence. Later corrective work does not rewrite that evidence.

## Permanent small-patch delivery policy

v18.9.x, v19 and v20 ship **small, independently auditable patches** rather than heavy releases.

For each patch:
1. establish exact predecessor/baseline at G0;
2. freeze one primary responsibility at G1 with explicit non-goals;
3. complete G2/G3 ownership/contract readiness before source edits;
4. use one development branch and one PR;
5. batch a coherent code+test candidate before opening/advancing CI;
6. exact-head Fast -> same PR Ready -> Qualified;
7. merge only the qualified exact head;
8. run the single canonical G11-G16 release path when release-capable;
9. publish only certified same-run artifacts without rebuild;
10. perform a post-release implementation-miss/open-issue audit and reconcile handoff/checkpoints before the next patch begins.

No duplicate release workflows, retry branches, certification branches or avoidable duplicate CI runs. Failure classification is required before rerun. If a planned patch is too large, split it before G1 rather than weakening auditability.

## v18.9.x ordered delivery train

- **v18.9.1** — runtime crash corrective only (#64).
- **v18.9.2** — TradeInsight Settings/API-key UX only.
- **v18.9.3** — coverage-aware Smart Provider Router v2 core only.
- **v18.9.4** — canonical company identity/all-desk presentation only.
- **v18.9.5** — Market Data Modes + capability diagnostics only.
- **v18.9.6** — TradeInsight SEC Form 4 enrichment only.
- **v18.9.7** — TradeInsight ticker/company search only.
- **v18.9.8** — TradeInsight movers/ranking evidence only.
- **v18.9.9** — remaining useful TradeInsight capability sweep only.
- **v18.9.10** — provider efficiency + Adaptive Intelligence telemetry only.
- **v18.9.11** — professional closure audit only; no new feature scope.

The train is dependency-ordered, but a later patch may be split further if its G0/G1 proves it is still too large. Patch numbers are not a reason to bundle scope.

## v19 delivery train — Professional Data Infrastructure

v19 begins only after v18 final closure. Deliver as small packets for: provider capability/entitlement/rights; provider quality/cost/coverage/SLOs; source disagreement/history quality; 13F infrastructure; two-sided evidence substrate; AODR point-in-time outcome lineage; ADR-GDI professional reliability/capacity; measured specialized/paid-provider gap evaluation; v20 research-readiness evidence audit; then mandatory v19 Major Closure.

No v19 packet may recreate Smart Provider Router v2, canonical identity, freshness or other owners established in v18. v19 measures and hardens them.

## v20 delivery train — Adaptive Intelligence & Decision Research

v20 begins only after v19 Major Closure. Deliver separately: adaptive research/experiment ledger; historical analogue/regime outcomes; calibration/FP-FN/miss/contradiction/drift; ASBI fingerprints/state transitions; ASBI scenarios/probability momentum; adaptive 13F; TDTI competing Long/Short/No Reliable Edge; TDTI two-sided trade-plan validation; AODR adaptive shared ranking; AODR diversity/personalized relevance after shared truth; ADR-GDI adaptive optimization; model/prompt + Champion/Challenger governance; then v20 Professional Closure.

Adaptive production influence remains `SHADOW -> VALIDATED -> APPROVED -> PRODUCTION`. No silent self-promotion, no silent deterministic formula change and No Execution.

## Delivery acceptance per patch

A patch is not complete merely because code compiles or CI is green. Delivery evidence must show:
- frozen G1 scope fully implemented;
- explicit non-goals untouched;
- deterministic regression tests for the changed responsibility;
- truthful degraded/partial/UNKNOWN/ABSTAIN states;
- no parallel owner introduced;
- no newly discovered implementation miss left without a durable target;
- packaged-runtime proof where the patch affects native/runtime behavior;
- current issue state and handoff agree with executable evidence.

For v19, acceptance additionally requires point-in-time provenance/data-rights/reliability truth where applicable. For v20, acceptance additionally requires leakage-safe evaluation, calibration/utility evidence, reproducibility, SHADOW/Champion-Challenger evidence and explicit rollback before promotion.

## Major-version handoff contract

- **v18 -> v19:** acquisition, canonical identity, adaptive Market Modes, useful provider capabilities and provider telemetry are trustworthy and zero-gap enough for professional measurement.
- **v19 -> v20:** provider/data quality, rights, provenance, point-in-time feature/outcome lineage and reliability history are sufficient for governed adaptive research.
- A major-version closure cannot be bypassed by starting the next major's implementation early.

## Final v18.9.x closure

v18.9.11 must prove zero unexplained carry-forward, zero orphan useful provider capability, zero duplicate routing/freshness/SEC/symbol/Market-Mode owner, deterministic Day/Swing/Long equivalence, #57 and #64 regressions, Adaptive Intelligence Scorecard results and actual macOS Apple Silicon + Windows x64 packaged-runtime evidence.

## Permanent boundaries

U.S. Equities Processing; GLD/SLV/USO actionable exceptions; No Execution; Smart Provider Router v2 sole routing authority; canonical freshness/recovery sole freshness owner; existing multi-feed allocator sole subscription owner; BroadSnapshotBroker canonical reuse owner; direct SEC/EDGAR authoritative; existing canonical cache/persistence/telemetry/symbol/state owners; G0-G16 only.

## Exactly one next action

Obtain/reproduce the #64 macOS crash and freeze the narrow v18.9.1 scope. Do not create the v18.9.2 branch until v18.9.1 is truthfully closed or the crash is proven external/non-product.
