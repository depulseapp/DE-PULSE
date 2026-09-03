# DE.PULSE Product Audit Rebaseline — 2026-08-27

**Status:** AUTHORITATIVE AUDIT ADDENDUM TO `governance/V19_V20_REBASELINE.md`  
**Audit reference:** `87baad3a90357b9af8e040a5ff34203a269fd861` / 27 Aug 2026 / stable comparison `v18.10.0`  
**Live implementation authority:** current source + `governance/current-state.json` + active work-slice/PR evidence  
**Current lifecycle:** `DEVELOPMENT`  
**Commercial activation:** `OWNER_EXPLICIT_ONLY` — never inferred from technical maturity, CI, hosted readiness, provider configuration, signing, release packaging, or a 5/5 score.

This addendum accepts the architecture/product audit direction without reopening the immutable `v18.10.0` certification and without discarding any existing capability. The audit is an architecture and planning rebaseline, not permission to delete working behavior.

## 1. Permanent zero-capability-loss rule

Rearchitecture must preserve capability while authority moves.

The following are cumulative conservation baselines and none may be silently dropped, weakened, bypassed, hidden beyond usefulness, or replaced without explicit equivalence evidence:

1. all **180 certified v18 responsibilities** and their durable regression ownership;
2. all `HOST-001..HOST-072` hosted requirements;
3. all mapped v19/v20 backlog requirements and source-discovered responsibilities;
4. all conserved legacy future-roadmap/build-plan commitments;
5. the current whole-product surface map and every useful data/intelligence consumer;
6. Smart Provider Router v2, Data Health/freshness/recovery, canonical cache/persistence/subscription/reconciliation/telemetry/lifecycle, identity/session and direct-authority boundaries;
7. approved platform requirements: macOS Apple Silicon, Windows x64 and future Web parity where applicable;
8. permanent product boundaries: U.S. Equities Processing, approved ETF exceptions, **No Execution**, no parallel canonical subsystems, no automatic provider lifecycle promotion;
9. every audit finding classified as `CURRENT`, `GAP`, `TARGET` or `UNVERIFIED` until explicitly dispositioned.

A UI consolidation may move or collapse presentation, but useful canonical evidence and downstream decision utility must remain available. A new canonical owner may replace an old implementation only after compatibility/equivalence evidence proves no capability loss.

## 2. Product architecture north star

DE.PULSE evolves toward one market-intelligence operating system, not a collection of parallel scorers.

Canonical flow:

`Providers / events -> canonical observations/evidence -> SymbolIntelligenceSnapshot -> Opportunity Lifecycle -> projections -> Decision Brief -> outcomes -> governed adaptation`

Required architecture rules:

- Preserve the Go core and evolve it into a **modular monolith** before considering service extraction.
- Introduce a versioned **`SymbolIntelligenceSnapshot`** as the authoritative per-symbol/as-of domain contract.
- Introduce one shared **Opportunity Lifecycle** used by Discovery, Opportunity Radar, Rapid Move, Watchlist, Alerts, Research and Desk handoffs.
- **Do not create a Watchlist scoring engine.** Watchlist changes the eligible user-selected universe and presentation; evaluation stays shared.
- Move authoritative technical features, deterministic horizon policies, side/setup geometry and scoring from renderer-owned logic into server-owned versioned Go domain packages.
- Keep compatibility adapters around `RuntimeSnapshot` and existing clients until equivalence and telemetry prove migration complete.
- Replace full high-churn snapshot fan-out with versioned user-scoped delta/event contracts gradually, not as a big-bang rewrite.
- Use Postgres v2 as the hosted normalized tenant-aware persistence target; keep SQLite as the local/offline development/client bridge where useful.
- Use one typed provider capability gateway/Registry + Smart Provider Router v2 path; preserve explicit direct-source authority such as SEC/EDGAR where governed.
- Make Research the authoritative frozen-as-of **Decision Brief** destination. AI may synthesize/explain bounded evidence but does not become a second source of truth.
- Web, macOS and Windows must consume the same domain truth and authorization semantics; platform differences are shell/mechanics, not intelligence forks.

## 3. Development-only lifecycle rule

DE.PULSE is in **DEVELOPMENT** until the Owner explicitly changes that state.

The following do **not** activate commercialization on their own:
- reaching `DEVELOPMENT_PRODUCTION_READY`;
- completing hosted infrastructure;
- obtaining provider rights;
- signing/notarizing installers;
- enabling Web;
- completing billing/support code;
- all domains reaching target 5/5;
- a release tag or G16 closure.

Commercial/public activation requires a separate explicit Owner declaration and a dedicated commercial-activation audit. Until then, commercial-only approvals must not be misclassified as blockers for ordinary development when they are genuinely commercial-only. Technical security, privacy, persistence, recovery, data-truth, provider-health, rights enforcement and supply-chain controls remain Development requirements where they protect the product itself.

## 4. Zero-implementation-miss contract

Every version/build must maintain a complete implementation matrix. For each applicable responsibility/finding/capability record:

`requirement -> baseline capability -> canonical owner -> source path -> upstream evidence -> consumers -> user decision purpose -> data/freshness/point-in-time semantics -> persistence/migration -> security/RBAC/rights -> UI/platform applicability -> positive test -> adverse test -> recovery/re-evaluation -> performance/backpressure -> regression owner -> compatibility/rollback -> closure evidence`

Mandatory rules:

1. No row is considered implemented merely because code exists.
2. No row is considered verified merely because a happy-path test passes.
3. No architecture extraction may delete the old owner before equivalence is proven on golden vectors and affected clients.
4. Every changed intelligence capability receives explicit `REQUIRED / CONDITIONAL / NOT_USEFUL` dispositions for Market Intelligence/Regime, Tradeability, Discovery/Radar, Watchlist, Research, Day/Swing/Long Desks, Prep, Alerts/Rapid Move/Catalyst, Outcome/Pattern Learning, Data Health/Maintenance and processing priority.
5. Every visible surface receives `KEEP / IMPROVE / CONSOLIDATE / COLLAPSE / MOVE / REMOVE`, with useful evidence conserved even if presentation changes.
6. G10 requires **zero unexplained applicable row** across all conservation families, including this audit addendum.
7. G16 may claim a maturity improvement only from executable/current evidence; roadmap completion text is not evidence.
8. `UNVERIFIED` remains `UNVERIFIED` until evidence exists; it cannot be converted to PASS by assumption.

## 5. 5/5 maturity objective

The audit scorecard is a **current-state baseline**, not a permanent rating. The engineering target is **5/5 in every domain**. Scores are increased only when the closure criteria below are evidenced.

| Domain | Audit baseline | Target | 5/5 closure meaning |
|---|---:|---:|---|
| Local product usefulness | 4/5 | 5/5 | coherent decision workflow; Dashboard/MI/Desks/Discovery/Watchlist/Research roles are distinct; no material UX fragmentation; useful deep evidence remains reachable |
| Deterministic intelligence | 4/5 | 5/5 | authoritative server-owned, versioned, replayable technical/horizon/setup logic; two-sided geometry; golden-vector parity; no client truth forks |
| Provider architecture | 3.5/5 | 5/5 | all provider capabilities classified; generic Registry/gateway + Router v2; direct-authority exceptions explicit; rights/freshness/cost/health/recovery/coalescing complete; no feature-owned bypass |
| Canonical state | 3/5 | 5/5 | versioned SymbolIntelligenceSnapshot + typed evidence/events/deltas; field-level freshness/provenance; shared Opportunity identity/lifecycle; compatibility migration complete |
| Adaptive intelligence | 2.5/5 | 5/5 | point-in-time outcome corpus, calibration/drift/model registry, shadow/champion-challenger, bounded reversible approved influence; deterministic fallback remains |
| Persistence | 2.5/5 | 5/5 | normalized tenant-aware Postgres v2, RLS/isolation, revisions/outbox/conflict/tombstones, retention/partitioning, backup/PITR/restore/DR evidence; local SQLite path remains correct |
| Security foundation | 3/5 | 5/5 | hosted identity/OIDC/PKCE/device/session/RBAC/entitlement isolation, managed/OS secrets, tenant/cache/stream negative tests, durable audit, supply-chain/signing controls |
| Testing / release assurance | 4/5 | 5/5 | current exact-head Fast/Qualified/G0-G16 plus contract/golden/chaos/load/hosted/client/adverse coverage, measurable coverage/flake/SLO evidence and zero unexplained responsibilities |
| Commercial distribution/readiness | 1.5/5 | 5/5 readiness, activation OFF | rights evidence, signed/notarized installers/update/rollback, support/ops/privacy/billing readiness and release channels are technically complete; **commercial/public activation still requires explicit Owner declaration** |
| Web + multi-platform readiness | 2/5 | 5/5 | production-capable versioned API/client boundary; Web/macOS/Windows parity; auth/deep links/sync/offline/degraded behavior; same domain truth across clients |
| Documentation coherence | 2/5 | 5/5 | one canonical current roadmap/architecture/build/release truth, generated current status, machine-state precedence, historical evidence archived/immutable, duplicate CURRENT/base narratives removed or demoted |

A 5/5 claim is not aspirational. Each domain must have objective evidence and no unresolved critical/high gap contradicting the score.

## 6. Rebaselined execution order

The existing active `v19.0.0` work is **not restarted**. Audit improvements are integrated around it using compatibility-first sequencing.

### A. Rebaseline + conservation controls — immediate governance overlay
- adopt this audit addendum;
- add the 5/5 machine target/closure checklist;
- bind audit findings to existing v19/v20 versions or explicit future residuals;
- preserve the active PR/work slice and exact current source truth;
- freeze golden vectors/characterization coverage before moving renderer-owned intelligence.

### B. Finish active v19.0 technical trust foundation
Continue existing HOST-001..023 work to Development Production Ready evidence. Commercial/public-only activation remains OFF.

### C. v19.1 canonical intelligence foundation
In addition to existing canonical runtime/global symbol/provider work:
- define Observation/Evidence/Snapshot/Transition contracts;
- introduce `SymbolIntelligenceSnapshot` behind compatibility adapters;
- establish server-owned technical feature and horizon policy package boundaries;
- create golden equivalence vectors for current `computePlan`, Radar, Rapid Move and Research behavior;
- no user-visible capability loss.

### D. v19.2 hosted serving/sync + normalized persistence
- tenant-aware Postgres v2, outbox/revision/conflict/tombstone;
- versioned user-scoped APIs/deltas;
- hosted gateway/cache/fan-out with rights/entitlement isolation;
- compatibility clients continue until migrated.

### E. v19.3 shared Opportunity Lifecycle + cross-platform product contract
- establish shared Opportunity aggregate/states/hysteresis/transitions;
- adapt Rapid Move, Opportunity Radar and Discovery as evidence producers/projections rather than parallel lifecycle owners;
- preserve Day/Swing/Long workflows while moving authoritative setup math server-side;
- establish cross-platform auth/role/IA parity.

### F. v19.4 / v19.4.1 Watchlist + Research + Discovery convergence
- dedicated first-class Watchlist tab as the user-selected-universe projection of shared lifecycle;
- ranked compact table, promotion/demotion transitions, contradictions, freshness/trust and Research handoff;
- frozen-as-of Decision Brief identity;
- durable transition alerts; no second Watchlist scorer;
- Discovery remains broad-universe projection; Radar remains a detector/evidence producer.

### G. v19.5.x intelligence enrichment
Price/volume/event and options/GEX capabilities feed the same canonical snapshot/lifecycle/Decision Brief and obey quality/rights gates. No isolated competing product state.

### H. v19.6.x point-in-time outcomes, reliability and 5/5 readiness reconciliation
- point-in-time/bitemporal data and outcome definitions;
- MFE/MAE/horizon return/false-positive/miss/confidence usefulness by regime/freshness/contradiction;
- observability, SLO, cost/capacity and recovery;
- full maturity matrix evidence review;
- any domain below 5/5 creates explicit residual work before major closure.

### I. v19.7 major deterministic/hosted closure
No new feature scope. Zero-gap reconciliation across conservation families and 5/5 deterministic/platform foundations. This closure still does **not** activate commercialization.

### J. v20 governed adaptation
Only after point-in-time outcomes are trustworthy: shadow/champion-challenger, calibration, pattern similarity, regime-conditioned bounded learning, adaptive provider usefulness and controlled AI/agent orchestration. Learned behavior remains reversible, measurable and policy-governed.

## 7. Mandatory compatibility strategy

For architecture migrations use:

`CHARACTERIZE -> ADD NEW OWNER -> DUAL WRITE/DUAL READ OR SHADOW -> COMPARE -> MIGRATE CONSUMERS -> PROVE EQUIVALENCE/IMPROVEMENT -> REMOVE OLD AUTHORITY -> KEEP REGRESSION`

Do not use:

`REWRITE -> DELETE OLD PATH -> DISCOVER MISSED CAPABILITY LATER`.

This applies especially to renderer scoring, RuntimeSnapshot, Discovery/Radar promotion semantics, watchlists/workspaces, provider direct paths, full-snapshot SSE, authentication/session state and persistence migrations.

## 8. Completion definition

The rebaseline is successful when:
- the product retains all currently useful capabilities;
- every planned capability has an assigned implementation/verification home;
- no parallel scorer/state/router/persistence/auth truth is introduced;
- the canonical intelligence/lifecycle architecture becomes authoritative gradually with measured parity;
- every maturity domain reaches evidence-backed 5/5;
- the product remains `DEVELOPMENT` until explicit Owner commercial activation.
