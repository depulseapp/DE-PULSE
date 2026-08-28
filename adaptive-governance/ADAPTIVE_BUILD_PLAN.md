# DE.PULSE — Adaptive Build Plan

**Status:** PERMANENT PLANNING CONTRACT / AUDIT-REBASELINED  
**Roadmap authority:** `governance/ROADMAP.md`  
**Operating authority:** `governance/ADAPTIVE-OPERATING-CONTRACT.md`  
**Live execution state:** `governance/current-state.json` + active closure ledger + `handoff/CURRENT.md`

This file defines how roadmap scope becomes dependency-correct build work. It does not own a changing branch SHA or next action.

## 1. Planning north star

Plan coherent versions through:

`approved scope -> exact source overlap -> canonical owner -> smallest safe build slices -> evidence -> AIPLC -> next dependency`

Always apply:

`REUSE -> CONSOLIDATE -> REFACTOR -> DELETE/REPLACE -> ADD`

Builds must improve one shared intelligence system, not accumulate a Watchlist engine, Discovery engine, Radar engine, Alert engine and Research engine that recalculate the same truth.

## 2. Mandatory planning inputs

Before G3, reconcile the version against:

- certified v18 responsibility and regression ownership;
- `governance/V19_V20_REBASELINE.md` plus backlog/HOST/legacy/cross-integration/surface maps;
- `governance/PRODUCT_AUDIT_REBASELINE_2026_08_27.md`;
- `governance/PRODUCT_AUDIT_COVERAGE_2026_08_27.md`;
- the full audit finding register and 5/5 maturity target under `governance/programs/V19-V20-REBASELINE/`;
- current source/runtime architecture, open issues and defects;
- provider capability, Data Health, rights, cost and entitlement evidence;
- the current G1 scope, closure ledger and exact prior PASS.

Do not treat documentation as implementation evidence. Unverified source/runtime/package claims stay OPEN or UNVERIFIED.

## 3. Version/build sizing

The planning and release unit is a coherent version/build. Requirements, audit rows and issue bullets remain traceability rows inside it.

- combine related small changes when they share owners, dependencies and rollback;
- split feature-heavy or high-risk work when behavior/evidence would become opaque;
- do not create a version, branch, PR or workflow per requirement;
- use coherent commits and risk-based CI checkpoints within the one active version PR;
- preserve one exactly named next action only in `handoff/CURRENT.md` and machine state.

## 4. Required build matrix

Every assigned requirement, defect, audit row or legacy commitment records:

- source-overlap disposition: `INHERITED / EXTEND_EXISTING_OWNER / REPLACE_CONSOLIDATE / NEW_RESIDUAL / EXTERNAL_BLOCKED / N_A`;
- canonical owner, upstream evidence and downstream consumers;
- user/trader decision purpose and materiality;
- temporal/freshness/provenance/revision/raw-adjusted semantics;
- persistence, retention, restart, migration and rollback applicability;
- positive, negative, failure and recovery evidence;
- RBAC, tenant, product-entitlement, provider-right and privacy applicability;
- load, provider quota/cost, concurrency/backpressure and platform applicability;
- Mac Apple Silicon, Windows x64 and Web requirement or justified N/A;
- compatibility/shadow/dual-read strategy when authority moves;
- durable regression owner;
- mapped `AUDIT-EXEC-*`, `AUDIT-RISK-*`, ADR and 5/5 closure rows.

Intelligence-bearing work additionally records:

- maturity: `DETERMINISTIC_ONLY / ADAPTIVE_CANDIDATE / LEARNING_ENABLED / AI_ASSISTED / NOT_USEFUL`;
- cross-integration `REQUIRED / CONDITIONAL / NOT_USEFUL` for Market Regime, Tradeability, Discovery, Watchlist, Research, Desks, Prep, Alerts, Outcome/Pattern, Data Health and processing priority;
- evidence families, contradictions, confidence, decay/expiry and recovery re-evaluation;
- point-in-time joins, censoring, sample uncertainty and selection-bias controls where outcomes/learning apply;
- personalization boundary versus shared model/policy truth;
- deterministic fallback and explanation lineage.

Visible work records `KEEP / IMPROVE / CONSOLIDATE / COLLAPSE / MOVE / REMOVE / ADD_FIRST_CLASS` and conserves #171 UX/intelligence-density requirements.

## 5. Audit-mandated architecture packets

The following packets are mandatory and must extend existing owners where possible.

### A. Characterization and domain extraction

- golden vectors for renderer `computePlan`, technical/family scoring, side/geometry, Rapid Move, Radar and Research;
- versioned Go Observation, Evidence, `SymbolIntelligenceSnapshot`, Transition and DecisionBrief contracts;
- compatibility adapter behind `RuntimeSnapshot` consumers;
- dual/shadow comparison, consumer migration and retirement of renderer/domain duplication;
- explicit intentional fixes for current long-oriented SHORT geometry defects.

### B. Shared Opportunity Lifecycle

- one symbol/opportunity identity and transition state machine;
- detectors submit evidence; they do not own lifecycle state;
- broad Discovery and selected Watchlist universe projections;
- material transition events for Alerts and frozen Research Brief lineage;
- temporal decay, contradictions, quality/freshness and deterministic explanation factors;
- halt/LULD/pause/resume and causal-deduplication behavior.

### C. Watchlist intelligence

- reuse existing Watchlist/workspace membership and canonical symbol intelligence;
- no second fetch, scanner, scorer, cache, persistence or lifecycle owner;
- rank only selected symbols and expose promotion/demotion, reason, contradiction, confidence, freshness and Research handoff;
- synchronize membership/settings through tenant-aware user state;
- define alert preferences/quiet hours/acknowledgement without rescore;
- defer Long King/Short King and Call/Put Wall until their formal contracts exist.

### D. Point-in-time outcomes and governed adaptation

- stable instrument identity, bitemporal/vintage facts and raw/adjusted basis;
- frozen feature/evidence/decision snapshots and explicit outcome windows;
- revision/supersedes lineage and censored outcomes;
- unbiased/control sampling against policy-induced selection bias;
- model/policy registry, time-split evaluation, shadow/champion-challenger, drift, approval and rollback;
- LLM explanation last, evidence-bounded and separately evaluated.

### E. Hosted/platform foundation

- tenant-normalized Postgres v2, isolation/RLS disposition, outbox, conflict/tombstone/retention and recovery;
- managed secrets, service/environment trust, supply-chain provenance and durable audit;
- versioned user-scoped REST/events/deltas and minimum-client compatibility;
- macOS/Windows secure storage, signed/notarized packages, installers, updates/channels and rollback;
- thin Web/macOS/Windows clients that do not own market calculations or provider credentials.

## 6. Provider and Data Health planning

Smart Provider Router v2 remains the sole general router. The Adaptive Provider Registry supplies adapter/capability metadata; it does not own routing, health, cache, persistence, subscriptions, lifecycle or canonical state.

Every provider change reconciles:

- `governance/data-health/provider-capability-matrix.json`;
- `governance/data-health/data-health-slo.json`;
- `governance/data-health/provider-fetch-paths.json`;
- provider rights for display/cache/persistence/derived/AI/multi-user/redistribution;
- capability/entitlement/freshness/history/quota/cost and upstream-correlation truth;
- Settings secret preserve/replace/clear/test/redaction behavior;
- outage, auth, rate-limit, malformed, stale, downgrade and recovery scenarios;
- cross-integration through canonical state rather than page-specific fetches.

Market Data remains the first generic Registry adopter in v19.1. Its current Bearer auth, `MARKETDATA_TOKEN`, HTTP 200/203 and delayed trial semantics are provider-contract observations to test—not eternal product assumptions.

## 7. Rebaselined version build sequence

`governance/ROADMAP.md` owns detailed placement. Build ownership is:

| Version | Primary build outcome |
|---|---|
| v19.0.0 | Hosted technical trust/identity/persistence/security/data-truth foundation; Commercial/Public OFF |
| v19.1.0 | Characterized server domain, versioned symbol/evidence contracts and generic provider foundation |
| v19.2.0 | Tenant Postgres v2/outbox/versioned serving/sync |
| v19.3.0 | Shared Opportunity Lifecycle, server-owned two-sided policy and cross-platform contract |
| v19.4.0 | Frozen Research Brief and first-class Watchlist projection |
| v19.4.1 | Discovery/Watchlist/Radar/Alert convergence |
| v19.5.0 | Price/volume and event-anchored intelligence |
| v19.5.1 | Formal options structure/GEX semantics |
| v19.6.0 | Point-in-time evidence and outcome-ready storage |
| v19.6.1 | Reliability/economics/observability/distribution and maturity residual closure |
| v19.7.0 | No-feature v19 major closure |
| v20.0–v20.6 | Governed outcomes, patterns, synthesis, institutional/thesis, agents, operations and adaptive closure |

No later row bypasses an unresolved earlier dependency unless an approved governance decision explicitly reclassifies it.

## 8. Three-depth AIPLC plan

- **Delta AIPLC:** affected surfaces/capabilities plus dependency sentinels after a meaningful checkpoint.
- **G10 reconciliation:** whole-version requirements, audit/5×5 rows, compatibility and cross-integration.
- **G16/Major Closure:** outcomes, provider/data utility, incidents, UI/decision value, packaged behavior and prevention fed into the next plan.

AIPLC produces both an immediate disposition and reusable prevention. It may reprioritize future work but does not bypass frozen scope or promote adaptive behavior.

## 9. CI and evidence plan

- focused local/unit/static checks while editing;
- coherent exact-head Fast at a development checkpoint;
- impact-selected Qualified at material risk boundaries and G10;
- Release G11–G16 only for an actual candidate;
- deterministic provider fixtures/replay for normal CI; bounded live smoke only when transport/auth/entitlement is the evidence subject;
- actual package/runtime evidence for platform, migration, installer/update, security and operational claims;
- no stale PASS reuse across a changed fingerprint.

CI cost is reduced through batching, impact routing and evidence equivalence—not by deleting assurance.

## 10. Zero-miss exit criteria

A version does not close until every applicable certified responsibility, backlog/HOST row, legacy commitment, audit finding/risk, surface disposition, ADR, temporal rule, role/right/platform case, compatibility migration and regression owner is implemented/evidenced or has a named truthful future disposition.

A 5/5 claim additionally requires all applicable machine target rows and no contradictory Critical/High gap. Planned scope and documentation do not earn maturity.

## 11. Resume rule

At every resume, verify GitHub, Stable identity, active branch/PR/head/checks, machine state, closure ledger and the last trustworthy PASS. The exact active dependency and one next action are intentionally not copied into this permanent plan; read them from `handoff/CURRENT.md`.
