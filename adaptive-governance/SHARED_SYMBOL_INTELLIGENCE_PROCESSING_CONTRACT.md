# DE.PULSE — Shared Symbol Intelligence Processing Contract

Status: PERMANENT
Effective: v18.2 governance; mandatory implementation progression through v18.3–v18.5

## 1. North Star

DE.PULSE must process market intelligence at the **shared canonical symbol layer first** and personalize delivery second.

Canonical flow:

**User/System demand → Global Symbol Registry → Shared Demand Union → Smart Provider Router → Canonical Symbol State → Shared calculations/correlation/intelligence → material-change propagation → capability/user-specific composition.**

The system must not create a separate market-data, indicator, event, research or AI pipeline merely because another user tracks the same symbol.

The target behavior is:

**process once where lawful and equivalent → remember once → correlate once → synthesize once where reusable → fan out many times.**

User-specific state may determine *what the user requests, is authorized to receive, and how intelligence is presented*; it must not automatically duplicate the underlying market-processing pipeline.

## 2. Canonical ownership

### Global Symbol Registry

The Global Symbol Registry is the canonical owner of instrument identity and shared processing membership. It must represent the active union of relevant demand, including where applicable:

- user tracked/watchlist symbols;
- selected/focused symbols;
- SPY/QQQ permanent context;
- explicitly tradable GLD/SLV/USO;
- Opportunity Radar/Discovery promotions;
- Day/Swing/Long candidates;
- Decision Queue demand;
- Rapid Move / Market Shock candidates;
- earnings/material catalyst demand;
- Pre-Market / Market Open temporal checkpoint demand;
- research/manual deep-dive demand;
- required market/sector/context symbols.

A user workspace owns personal membership, preferences and workflow context. It does **not** become a second owner of quotes, provider subscriptions, canonical history, indicators, filings, news, event truth or common intelligence.

### Canonical processing key

Shared work is deduplicated using a canonical processing identity that is sufficiently specific to preserve correctness. At minimum the key considers:

**symbol/instrument + dataset/capability + market session/time window + freshness/materiality requirement + provider/entitlement/data-rights domain + policy/model/version where relevant.**

Work may be shared only when those dimensions are equivalent. Intentional independent-provider reconciliation is not considered duplicate work.

## 3. Shared execution plane

For each canonical processing key, DE.PULSE should have one logical shared execution path for:

1. acquisition/subscription;
2. normalization;
3. timestamp/provenance/freshness classification;
4. validation and provider reconciliation;
5. canonical state update;
6. deterministic calculations;
7. event/catalyst correlation;
8. intelligent/adaptive synthesis;
9. bounded persistence/history/outcome linkage;
10. downstream fan-out.

Scanner, Opportunity Radar, preparation checkpoints, catalyst evaluation, desks, Research and user views consume this shared state rather than independently refetching or recalculating equivalent evidence.

## 4. Demand union and dynamic attention

The processing union is global/shared, but processing intensity is adaptive rather than uniform.

The scheduler/router must continuously rank work by materiality, freshness need, user/system demand, session, event risk, current decision relevance and provider/runtime capacity.

High-attention examples include selected symbols, Day candidates, Decision Queue items, Rapid Move/Market Shock candidates, imminent earnings/material catalysts and actionable tracked symbols. Lower-value/inactive/background symbols may use slower refresh, cached state, event-driven reactivation or bounded historical processing.

No user's oversized list or burst of requests may starve materially higher-priority shared work or other authorized users. Fairness, backpressure, provider budgets and bounded concurrency are mandatory.

## 5. Fetch once / calculate once / single-flight

Equivalent simultaneous requests must be coalesced using shared cache/in-flight ownership (single-flight or equivalent).

Required behavior:

- one fresh canonical acquisition may satisfy multiple consumers;
- one calculation result may satisfy multiple surfaces/users when its inputs/version are identical;
- one active provider subscription may feed all authorized equivalent consumers;
- cache entries are freshness/provenance aware, not blind TTL shortcuts;
- a cache miss for one consumer must not create a request storm when the same work is already in flight;
- duplicate background jobs must not refetch the same dataset merely because they have different feature names.

## 6. Material-change propagation

DE.PULSE should prefer event/material-change propagation over repeated full recomputation.

A new canonical event/state change should invalidate or recompute only affected derived intelligence and downstream consumers. Unchanged evidence should retain reusable calculations/synthesis where safe.

Examples:

- a new quote updates only price-dependent state that needs the new point;
- a new filing/catalyst invalidates relevant catalyst/readiness/research conclusions;
- a role/UI change does not trigger market-data reacquisition;
- a user adding an already-active symbol normally adds a new consumer, not a new market pipeline;
- a user removing a symbol removes personal demand but shared processing continues while another legitimate consumer still requires it.

## 7. AI/LLM/adaptive intelligence reuse

AI/LLM-style processing follows the same shared-first architecture.

The preferred flow is:

**canonical evidence package → evidence fingerprint → correlation/synthesis → reusable intelligence result → user-specific composition/question layer.**

A reusable AI/adaptive result should be keyed by the evidence fingerprint plus relevant model/prompt-policy/version/context identity. Identical evidence must not be repeatedly re-summarized simply because multiple users open the same symbol.

User-specific LLM work is justified only when private user context, an explicit user question, entitlement, workflow state or presentation preference materially changes the required reasoning/output.

Private prompts, private user context, restricted data and authorization-scoped outputs must never be shared through a cross-user cache. Shared intelligence must remain within compatible data-rights, tenant/security and entitlement domains.

Production adaptive-policy changes continue to follow **SHADOW → VALIDATED → APPROVED → PRODUCTION** and may not silently alter protected deterministic Day/Swing/Long logic.

## 8. Personalization boundary

Per-user processing should primarily consist of:

- tracked-symbol membership and personal workflow state;
- role/capability/authorization filtering;
- user-specific ranking/preferences when approved;
- notification eligibility;
- UI composition and hierarchy;
- explicit user research questions/private context.

Per-user processing should **not** normally duplicate:

- provider quote/history/news/SEC/earnings acquisition;
- canonical normalization/freshness/provenance;
- common technical calculations;
- shared event correlation;
- common Rapid Move/Market Shock validation;
- reusable evidence packages or identical AI synthesis.

Authorization is always applied before delivery. Shared processing is not permission to expose shared data to an unauthorized user.

## 9. Persistence and warm state

Live/current canonical symbol state is memory-first. Durable persistence stores meaningful last-known state, evidence/events, decisions, outcomes, provider usefulness, learning state and recovery information needed for restart/warm-start truth.

Do not persist raw high-frequency data merely because it exists. Retention follows Data Utility, point-in-time truth, rights and learning value.

Restart recovery should restore reusable canonical state once, then revalidate freshness/materiality before fan-out. It must not rebuild equivalent state separately for each user.

## 10. 10/10 efficiency scorecard

Affected releases must measure, where applicable:

- unique active symbols and canonical processing keys;
- total consumer/user demand vs unique demand union;
- provider calls/subscriptions per unique processing key;
- duplicate acquisition rate;
- in-flight coalescing/single-flight hit rate;
- cache reuse and freshness-safe cache hit rate;
- duplicate calculation rate;
- reusable AI/evidence-synthesis hit rate;
- fan-out ratio: consumers served per canonical result;
- marginal cost of adding an overlapping user/symbol consumer;
- CPU/memory/storage cost per active symbol/key;
- provider rate-limit/backpressure pressure;
- p50/p95 acquisition-to-canonical latency;
- p50/p95 material-change-to-consumer latency;
- stale/degraded fan-out incidents;
- fairness/starvation incidents;
- data-rights/authorization leakage incidents;
- provider calls avoided and useful-provider ratio from the Adaptive Intelligence Scorecard.

### 10/10 architectural acceptance

For equivalent overlapping demand:

- duplicate acquisition/calculation should be **zero by design** except documented independent-validation/reconciliation or technically unavoidable cases;
- simultaneous equivalent misses must collapse to one in-flight owner where practical;
- adding another user who tracks an already-active symbol should have marginal cost dominated by authorization/composition/fan-out, not another full market pipeline;
- processing cost should scale primarily with **unique canonical demand**, not `users × symbols`;
- no sharing optimization may weaken freshness, provenance, point-in-time truth, user isolation, data rights or authorization;
- any justified duplicate work must have an explicit owner, reason, bound and scorecard visibility.

A regression to per-user duplicate market pipelines, uncontrolled provider fan-out, repeated identical AI synthesis or unbounded duplicate computation is an architecture/performance defect, not an acceptable scaling strategy.

## 11. Gate/checkpoint integration

This contract is enforced inside the existing canonical gate model; it does not require a new gate by itself.

- **G1 — Immutable Scope:** identify affected symbol/data/intelligence consumers and any new demand sources.
- **G2 — Architecture & Data Utility:** prove canonical owner, processing key, demand-union behavior, rights/entitlement boundary, reuse/coalescing and correlation design.
- **G3 — Design & Dependency Readiness:** map producers/consumers, invalidation/material-change graph, concurrency/backpressure and scorecard/load-test plan.
- **G4 — Development Exit:** no new parallel per-user or per-feature equivalent market pipeline may remain unintentionally.
- **G6 — Integration & MEDIUM Qualification:** prove multiple consumers share canonical state and one consumer's lifecycle does not corrupt others.
- **G7 — Data, Security & Adaptive Intelligence:** prove provenance, point-in-time truth, AI reuse boundaries, data-rights isolation and no cross-user private-context leakage.
- **G8 — Performance, Capacity & Stability:** measure union scaling, provider calls, coalescing, duplicate computation, memory/CPU, fairness and long-run stability.
- **G9 — Cross-Module UI/UX:** ensure shared intelligence is composed appropriately without duplicating deep evidence or exposing implementation machinery.
- **G10 — Pre-Freeze Qualification:** unresolved unjustified duplicate acquisition/computation/shared-intelligence ownership blocks freeze for affected scope.
- **G12 — Full Certification:** replay architecture/integration/performance/security evidence on the immutable RC where applicable.
- **G16 — Adaptive Retrospective & Handoff:** record reuse ratio, calls avoided, remaining justified duplication, regressions and generalized preventative learning.

## 12. Release progression

### v18.2

Preserve the invariant that `UserWorkspace` and role/capability work do not create independent per-user market-data/provider/intelligence pipelines. Governance adoption and regression protection are current-release process requirements; unrelated large source refactors remain outside immutable v18.2 product scope.

### v18.3 — mandatory shared execution entry

The existing mandatory functionality-remediation work becomes the first concrete implementation stage of this contract:

- shared/coalesced Scanner + Opportunity Radar acquisition;
- Session Intelligence Coordinator over canonical state for Pre-Market / Market Open checkpoints;
- Event Intelligence ownership for earnings/material catalyst lifecycle;
- hosted/shared-server demand-union behavior;
- no per-user provider engine duplication;
- canonical shared state compatible with PostgreSQL persistence/recovery boundaries.

### v18.4 — rights/security isolation closure

Prove that shared processing respects role/capability boundaries, provider/data rights, tenant/user isolation, secrets/security policy and commercial/governance requirements. Shared caches/synthesis must be partitioned where rights or private context differ.

### v18.5 — 10/10 major closure

Run the full shared-processing efficiency/capacity/security scorecard under realistic multi-user overlapping-demand scenarios. v18.5 major closure must demonstrate that processing scales with unique canonical demand, that duplicate work is bounded/justified, and that material-change propagation, warm state, provider budgets and adaptive intelligence reuse remain stable under load.

Any material shortfall is retained as an explicit blocker/closure item rather than being hidden by average UI responsiveness.

## 13. Permanent invariant

**Global/shared market intelligence first; authorized personal composition second.**

DE.PULSE should become more efficient as users overlap in symbols/evidence—not proportionally more expensive because the same market truth is requested more than once.
