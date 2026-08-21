# DE.PULSE — Current Adaptive Roadmap

**Certified Stable:** `v18.9.0-stable`  
**Certified candidate:** `9ea81cddae4875ae15d3719ca028519a36c597b6`  
**Certified fingerprint:** `a8719090c341c874dbd1279cc31ad98e84075d5701c46a800bf951340780ecb9`  
**Build ID:** `v18.9.0-stable-20260821`  
**Active corrective program:** issue #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Immediate blocker / next patch:** issue #64 / `ADAPT-RUNTIME-CRASH-001` -> `v18.9.1`.

## v18.9.0 — COMPLETE / IMMUTABLE STABLE

Issue #61 / `ADAPT-TRADEINSIGHT-001` is closed completed. Exact source head `9e86b5e731f7a585cc77c1521f3639fc7a208efc` passed Fast #481 and Qualified #153. Merged candidate `9ea81cddae4875ae15d3719ca028519a36c597b6` passed canonical Release #32 through G11–G16. Durable release evidence is `release/v18.9.0/stable-evidence-manifest.json`.

The post-release audit found that the architecture is sound but v18.9.0 did not fully realize the intended adaptive multi-provider/product UX. Those findings are not retroactively added to the immutable Stable artifact; they are governed by #65 as small v18.9.x patches.

## Permanent Small-Patch Operating Rule

DE.PULSE prefers **many small, dependency-ordered, complete patches** over heavy multi-domain builds. This applies permanently to v18.9.x, v19, v20 and later releases.

- One primary responsibility per patch; only tightly coupled support work may accompany it.
- No stability + routing + provider-expansion + UX bundles.
- Each patch starts from an exact G0 baseline and immutable G1 scope.
- Before the next patch starts, the current patch gets implementation-miss review, focused regression proof, runtime/browser proof where applicable, open-issue reconciliation and durable handoff.
- A known implementation miss must be fixed in-scope or explicitly registered against a later patch before closure; it may not disappear into chat memory.
- One development branch + one PR per patch; no retry/certification branch families and no duplicate CI runs.
- Exact future patch numbers are assigned/frozen only at that patch's G1; the roadmap sequence describes responsibility order, not a reason to bundle work.
- G0–G16 remains the only release model.

## v18.9.x ordered patch roadmap

1. **v18.9.1 — Runtime crash corrective ONLY** — #64. Diagnose/fix the real macOS Apple Silicon SIGABRT from evidence/reproduction; preserve user state/API keys; add lifecycle regression and actual packaged macOS proof.
2. **v18.9.2 — TradeInsight Settings/API-key UX ONLY.** Existing Data Provider Settings/secret owner; masked Save/Test/Clear; truthful status; environment override only as developer/runtime fallback.
3. **v18.9.3 — Coverage-aware Smart Provider Router core ONLY.** Upgrade first-success behavior to requirement/coverage-aware fulfillment; persistence/cache first; compute residual gaps; merge/provenance/re-evaluate; validation lifecycle separated from serving role.
4. **v18.9.4 — Canonical company identity + all-desk presentation ONLY.** Shared symbol/company identity; `APP - AppLovin : In Entry Zone` with symbol-only fallback; reused by desks/Research/Discovery/Add Symbol.
5. **v18.9.5 — Market Data Modes + capability diagnostics ONLY.** Behavior-oriented Adaptive modes rather than provider-brand modes; capability-level source/freshness/coverage diagnostics; no separate TradeInsight mode.
6. **v18.9.6 — TradeInsight SEC Form 4 enrichment ONLY.** Contract-validated SHADOW-first enrichment/corroboration; direct SEC/EDGAR authoritative; source-family de-duplication.
7. **v18.9.7 — TradeInsight ticker/company search ONLY.** Contract-validated fallback/corroboration through canonical symbol validation/company identity; U.S.-equity boundary final.
8. **v18.9.8 — TradeInsight movers/ranking evidence ONLY.** Contract-validated candidate evidence into Opportunity Radar; existing scanner/ranker remains canonical; SHADOW-first usefulness proof.
9. **v18.9.9 — Remaining useful TradeInsight capability sweep ONLY.** Every useful entitlement gets explicit disposition and consumer; retest Congress/history/corporate actions under coverage-aware routing; no invented endpoints or Python/MCP production dependency.
10. **v18.9.10 — Provider efficiency + Adaptive Intelligence telemetry ONLY.** Coverage filled, residual gaps, DB/cache hits, calls avoided, provider usefulness, provider-capacity reserve, latency/rate-limit/freshness/conflict telemetry, bounded fan-out and runtime-load proof.
11. **v18.9.11 — Session-Aware Data Readiness Maintenance ONLY.** One canonical coordinator: light overnight gap filling/readiness preparation plus heavier bounded weekend reconciliation/backfill. Existing U.S. market calendar owns session truth. Pre-market, regular market and after-hours are protected Tier-0 sessions with first claim on provider quota/headroom, CPU, memory, DB, network and worker capacity. Maintenance must drain/preempt/checkpoint before protected work and never create self-inflicted degradation.
12. **v18.9.12 — Whole v18.9.x professional closure audit ONLY.** End-to-end implementation-miss audit, #57/#64 regression, deterministic Day/Swing/Long equivalence, DB-first reuse/residual-gap acquisition, overnight/weekend maintenance, protected-session capacity reservation/preemption, actual macOS/Windows packaged proof, Adaptive Intelligence Scorecard and zero unexplained carry-forward/orphan useful capability/duplicate owner.

## Permanent adaptive provider + persistence architecture

DE.PULSE operates as:

`consumer requirement -> in-memory canonical cache -> persisted canonical DB/state -> validate coverage/freshness/schema/provenance/rights -> exact residual gap -> eligible-provider ranking -> targeted acquisition -> canonical merge/provenance -> coverage re-evaluation -> next provider only if still needed -> persist -> synthesized consumer state`

A successful provider response is not enough to stop if required coverage/freshness/fields/quality remain incomplete. No fixed global chain such as `Alpaca -> TradeInsight -> Twelve Data -> yfinance` is the decision model; static ordering is at most a prior/tiebreaker. Smart Provider Router v2 remains sole executable routing authority.

**Persistence-first / reuse-first is permanent:** do not refetch or recompute trustworthy evidence already available for the consumer requirement. Fetch only missing/expired/revised/materially insufficient evidence where rights and provider contracts permit.

**Protected-session priority is permanent:** pre-market, regular market and after-hours decision-support workloads outrank maintenance. Maintenance uses only bounded surplus capacity after provider/runtime reserves. Light overnight maintenance prepares the next session; heavy weekend maintenance performs deeper useful backfill/reconciliation. Both must yield to protected or market-shock workloads. Machine details live in `adaptive-governance/PERSISTENCE_REUSE_AND_OFF_HOURS_DATA_READINESS_CONTRACT.md`.

## v19 — Professional Data Infrastructure

**Entry condition:** v18.9.x final closure must be zero-gap enough to serve as the trusted acquisition/identity/telemetry/persistence-readiness foundation. v19 does not redo the router, maintenance coordinator or create new provider-specific owners; it measures and professionalizes the data plane created in v18.

Canonical purpose: make provider/data quality, rights, cost, reliability and suitability measurable rather than assumption-driven, and create sufficient point-in-time evidence/provenance/outcome history for v20.

### Mandatory hosted account/sync/provider-gateway program — issue #66 / `ADAPT-HOSTED-SYNC-001`

v19 MUST include the approved single-account macOS/Windows/web architecture as executable roadmap scope, not as a side document. The program preserves native SQLite as the offline edge/warm working set, introduces PostgreSQL only as shared hosted authority for sync-eligible account/device state and lawful hosted evidence, and uses typed authenticated incremental application-level synchronization rather than raw database replication.

Production commercial users follow the **end-user zero-key** model: clients authenticate only to DE.PULSE; platform-owned provider credentials remain server-side in the canonical managed-secret/KMS owner; all REST/snapshot/live-stream provider access crosses the authenticated DE.PULSE hosted boundary and reuses Smart Provider Router v2, canonical freshness/cache/persistence, the existing multi-feed subscription owner, rights/entitlement gates and protected-session resource controls.

The program is not complete unless it covers canonical `SUPER_OWNER/OWNER/ADMIN/USER/DEMO` capability truth; hosted identity/device/session lifecycle; privileged re-authentication; account/device revocation; PostgreSQL HA/PITR/backup/restore/migration safety; secret rotation/compromise recovery; entitlement-aware cache/live fan-out; bootstrap for new or stale devices; SQLite outbox/idempotency/checkpoints; tombstone/change-log retention and compaction; local account isolation/lost-device behavior; mixed-version sync protocol compatibility; rights-aware research synchronization; per-account usage/cost/abuse telemetry; and multi-user load/security/DR assurance.

Issue #66 and its architecture-audit addenda are the durable detailed scope. At v19 G0/G1 this program is split into small dependency-ordered patches; it must not be bundled into one heavy build. It must be operationally mature before v19 Major Closure and before v20 relies on synchronized evidence/outcomes.

### Provisional v19 small-patch train
Exact semantic patch numbers are frozen only at each G1 and any item may split further.

1. **Provider capability + entitlement + rights registry.** One machine-readable capability matrix covering entitlement, serving role, validation lifecycle, redistribution/persistent-storage/AI-use/commercial rights and U.S.-equity suitability.
2. **Provider quality / cost / coverage / SLO scorecards.** Measured freshness, completeness, latency, reliability, rate-limit pressure, cost/value, contribution/usefulness, calls avoided, maintenance value and fallback quality using v18.9 telemetry.
3. **Data reconciliation + disagreement + historical-quality hardening.** Source independence, conflict/reconciliation policy, corporate-action/adjustment correctness, historical depth/gaps, revision preservation and point-in-time provenance.
4. **Institutional / 13F evidence infrastructure hardening.** Direct SEC truth, manager identity, CIK, CUSIP/FIGI/security mapping, amendments/restatements, filing-lag truth, point-in-time holdings, storage/indexing and outcome lineage.
5. **Two-sided thesis evidence substrate.** Point-in-time Long/Short plan snapshots, target/invalidation ordering, side-aware MFE/MAE and reliable short-interest/crowding/borrow/SSR context only where lawful/trustworthy; explicit UNKNOWN otherwise.
6. **AODR opportunity evidence/outcome infrastructure.** My Market vs Global truth, point-in-time candidate/rank/reason lineage, NOW/WATCH/PASS/ABSTAIN transitions, shared-ranking efficiency, diversity/correlation metadata, recommendation usefulness and missed-opportunity outcomes.
7. **ADR-GDI professional reliability hardening.** Capability SLOs, degradation history, provider/DB/runtime reliability scorecards, restart/warm-start, query/index/pool/capacity tuning, load shedding, bounded operating limits, protected-session reserves and maintenance/preemption economics.
8. **`ADAPT-HOSTED-SYNC-001` hosted account/sync/zero-key provider-gateway program.** Deliver issue #66 through small sub-packets covering identity/device/session and account lifecycle; PostgreSQL HA/PITR/schema/pool foundation; managed secrets/KMS and zero-key Provider Gateway; entitlement-aware REST/live-stream serving; sync bootstrap/outbox/idempotency/change-log/checkpoint/compaction protocol; macOS pilot; desks/workspaces; Windows parity; hosted-web parity; rights-aware research/state sync; and final multi-user cost/security/licensing/DR assurance. Reuse all canonical owners; no client provider secrets and no raw SQLite↔PostgreSQL replication.
9. **Specialized/paid-provider gap evaluation.** Only consider replacement/additional paid data where measured v19 evidence proves a material capability/quality/rights gap; integrate through the same Smart Provider Router/persistence/session-priority contracts, never by special path.
10. **v20 research-readiness dataset/lineage audit.** Prove sufficient point-in-time evidence, feature history, outcomes, provenance, rights, independence and reliability history for ASBI, 13F Intelligence, TDTI, AODR and adaptive reliability optimization, including synchronized evidence only after #66 rights/security/sync assurance is proven.
11. **v19 Major Closure — mandatory before v20.** Whole-system data-quality/data-rights/performance/security/utility audit; zero unexplained provider role, zero unowned dataset, truthful commercial/data-rights posture and executable package/runtime evidence. Closure additionally requires #66 Mac/Windows/web account parity, zero-key provider-gateway security, sync correctness/offline recovery, entitlement isolation, secret rotation, PostgreSQL recovery and multi-user load/cost assurance to PASS.

## v20 — Adaptive Intelligence & Decision Research

**Entry condition:** v19 Major Closure PASS. v20 consumes the trusted point-in-time evidence/outcome history built in v18/v19; it must not compensate for unreliable acquisition by inventing confidence.

Canonical purpose: improve decision support from historical outcomes while preventing a silent self-modifying trading system. Production influence remains `SHADOW -> VALIDATED -> APPROVED -> PRODUCTION`, with explicit rollback and no execution.

### Provisional v20 small-patch train
Exact patch numbers are frozen only at G1 and complex items are split further instead of bundled.

1. **Adaptive research control plane + immutable experiment ledger.** Dataset/version lineage, feature/provenance snapshotting, evaluation cohorts, model/prompt versions, reproducibility, leakage controls and promotion/rollback evidence.
2. **Historical analogues + regime-conditioned outcomes.** Point-in-time analogue retrieval, regime/sector/catalyst conditioning and outcome distributions without changing deterministic Day/Swing/Long truth.
3. **Calibration / false-positive / miss / contradiction / drift intelligence.** Confidence calibration, FP/FN/missed-opportunity analysis, evidence contradiction, distribution drift and abstention thresholds.
4. **ASBI I — Behavioral Fingerprints + state transitions.** Canonical behavior features, sequence states, hierarchical symbol/peer/sector/regime context and immutable Behavior Intelligence Ledger.
5. **ASBI II — scenarios / probability momentum / calibration.** Competing paths, multi-horizon outlooks, Behavior Probability Momentum, expected-move distributions, evidence sufficiency/conflict and ABSTAIN/NO RELIABLE EDGE.
6. **Adaptive Institutional / 13F Intelligence.** Manager behavioral fingerprints, persistence/concentration, accumulation/reduction breadth, consensus/crowding, rotation, usefulness by regime/stock type, stale-data penalties and calibrated outcomes.
7. **TDTI I — Competing Long / Short / No Reliable Edge theses.** Same canonical snapshot; separate direction probability, thesis strength, confidence and opportunity quality; cause-aware confirmation/invalidation.
8. **TDTI II — Two-sided trade-plan intelligence + validation.** Long/Short entry/target/invalidation/R:R, side-aware readiness, probability momentum, time-to-resolution, MFE/MAE, risk intelligence and historical calibration; still No Execution.
9. **AODR I — Adaptive shared opportunity ranking.** Cross-candidate ranking using canonical ASBI/TDTI readiness/quality, expected magnitude/time-to-resolution, extension/chase/R:R/degradation penalties and candidate-vs-surfaced-vs-missed outcomes.
10. **AODR II — diversity + personalized relevance after shared truth.** Correlation/theme/catalyst diversity, opportunity cost, user relevance layered after canonical market truth, recommendation utility and ABSTAIN/no-strong-opportunity as a valid result.
11. **ADR-GDI adaptive optimization.** Governed SHADOW/Champion-Challenger learning for provider recovery prediction, cooldown/backoff, workload priority, maintenance value, protected-session reserve sizing, fallback usefulness and capacity policy; cannot self-promote or reduce live-session protection without evidence/approval.
12. **Model/prompt governance + Champion/Challenger system.** Explainability, independent evaluation, model/prompt drift, reproducible comparisons, approval/rollback and evidence-bound promotion across adaptive intelligence.
13. **v20 Professional Closure.** Principal Engineer + Professional Trader/Investor audit, calibration/utility/drift/abstention proof, deterministic-boundary protection, privacy/security/data-rights review, actual supported packages and zero silent self-modification/zero execution.

## Why the versions fit together

- **v18.9.x = trustworthy plumbing and truth:** stable runtime, Settings, coverage-aware acquisition, persistence-first reuse, canonical identity, Market Data Modes, useful TradeInsight evidence, provider-efficiency telemetry and session-aware overnight/weekend readiness without compromising live sessions.
- **v19 = professional measurement, hosted account/sync and evidence infrastructure:** prove which sources are useful/reliable/lawful, reconcile disagreement, harden storage/revision/capacity/maintenance economics, deliver the #66 zero-key Mac/Windows/web account architecture, and preserve 13F/TDTI/AODR/ADR-GDI point-in-time evidence/outcome lineage.
- **v20 = governed learning from that evidence:** ASBI, adaptive 13F, TDTI, AODR and reliability optimization use the v19 dataset/scorecards instead of learning from noisy or ambiguous inputs.

This ordering prevents v20 from learning provider artifacts, stale/partial data, survivorship leakage or undocumented provenance. It also prevents v19 from creating a second router or duplicating the adaptive acquisition/maintenance work completed in v18.9.x.

Permanent constraints: U.S. Equities Processing, No Execution, Smart Provider Router v2 sole routing owner, canonical freshness/recovery sole freshness owner, existing multi-feed allocator sole subscription owner, BroadSnapshotBroker canonical reuse owner, direct SEC/EDGAR authoritative, canonical persistence/cache owners reused, canonical U.S. session calendar reused, GLD/SLV/USO actionable tradable exceptions and deterministic Day/Swing/Long truth protected.

## Exactly one next action

Perform issue #64 / v18.9.1 G0 crash diagnosis from concrete macOS evidence or deterministic reproduction. Do not start v18.9.2 or any v19 implementation branch until v18.9.1 is closed with truthful evidence or the crash is proven external/non-product.
