# DE.PULSE — Current Adaptive Delivery Process

**Certified Stable:** `v18.9.0-stable`  
**Stable candidate:** `9ea81cddae4875ae15d3719ca028519a36c597b6`  
**Stable fingerprint:** `a8719090c341c874dbd1279cc31ad98e84075d5701c46a800bf951340780ecb9`  
**Master corrective program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Immediate patch:** `v18.9.1` / #64 / `ADAPT-RUNTIME-CRASH-001`.

## v18.9.0 delivery — COMPLETE / IMMUTABLE

Fast #481, Qualified #153 and Release #32 are the release authority for v18.9.0. G11-G16 passed, including actual macOS Apple Silicon + Windows x64 packaged-runtime audits, G15 assurance, no-rebuild Stable publication and G16 handoff evidence. Later corrective work does not rewrite that evidence.

## Small-patch delivery policy

The v18.9.x program ships **small, independently auditable patches** rather than heavy releases.

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

No duplicate release workflows, retry branches, certification branches or avoidable duplicate CI runs. Failure classification is required before rerun.

## Ordered delivery train

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

The train is dependency-ordered, but a later patch may be split further if its G0/G1 proves it is still too large. Patch numbers are not a reason to bundle scope. Smaller is preferred when it improves correctness and auditability.

## Delivery acceptance per patch

A patch is not complete merely because code compiles or CI is green. Delivery evidence must show:
- frozen G1 scope fully implemented;
- explicit non-goals untouched;
- deterministic regression tests for the changed responsibility;
- truthful degraded/partial states;
- no parallel owner introduced;
- no newly discovered implementation miss left without a durable target;
- packaged-runtime proof where the patch affects native/runtime behavior;
- current issue state and handoff agree with executable evidence.

## Final v18.9.x closure

v18.9.11 must prove zero unexplained carry-forward, zero orphan useful provider capability, zero duplicate routing/freshness/SEC/symbol/Market-Mode owner, deterministic Day/Swing/Long equivalence, #57 and #64 regressions, Adaptive Intelligence Scorecard results and actual macOS Apple Silicon + Windows x64 packaged-runtime evidence.

## Permanent boundaries

U.S. Equities Processing; GLD/SLV/USO actionable exceptions; No Execution; Smart Provider Router v2 sole routing authority; canonical freshness/recovery sole freshness owner; existing multi-feed allocator sole subscription owner; BroadSnapshotBroker canonical reuse owner; direct SEC/EDGAR authoritative; existing canonical cache/persistence/telemetry/symbol/state owners; G0-G16 only.

## Exactly one next action

Obtain/reproduce the #64 macOS crash and freeze the narrow v18.9.1 scope. Do not create the v18.9.2 branch until v18.9.1 is truthfully closed or the crash is proven external/non-product.
