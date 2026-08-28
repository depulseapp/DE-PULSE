# DE.PULSE — Adaptive Build Process

**Status:** PERMANENT EXECUTION CONTRACT / AUDIT-REBASELINED  
**Authority:** `governance/ADAPTIVE-OPERATING-CONTRACT.md`  
**Plan:** `adaptive-governance/ADAPTIVE_BUILD_PLAN.md`  
**Live state:** `governance/current-state.json` + active closure ledger + `handoff/CURRENT.md`

Execution loop:

`RESUME -> LOOKUP -> COMPARE -> CLASSIFY -> CHARACTERIZE -> BUILD/SHADOW -> TEST -> RECONCILE -> Fast -> Qualified -> G11-G16 when releasing -> AIPLC`

## 1. Resume and exact-baseline process

Before planning or editing:

1. read `AGENTS.md`, portability/CI contracts, `governance/README.md` and `handoff/CURRENT.md`;
2. inspect actual GitHub default branch, active branch/PR/head, checks, artifacts and latest Stable;
3. read `release_identity.json`, machine current state, active scope/closure ledgers and resume checkpoints;
4. run or inspect `tools/ci/adaptive_resume_gate.py` and `tools/ci/workflow_policy.py`;
5. reconcile any disagreement and resume from the last trustworthy PASS / earliest open G0–G16 responsibility.

Assistant memory, a prior handoff copy or an old source snapshot cannot outrank live GitHub evidence.

## 2. Source overlap and classification

For each requested change or audit finding:

- search source, tests, schemas, providers, migrations, workflows and historical evidence;
- trace `UI -> state -> service -> intelligence -> provider/data -> persistence -> consumer/output`;
- classify capability as `FULL / INCOMPLETE / WEAK / UI_ONLY / BACKEND_ONLY / PROTOTYPE / STUB / FLAG_ONLY / PLANNED / DOCUMENTED_ONLY / LEGACY / DEAD / DUPLICATE / BROKEN / UNVERIFIED`;
- classify work as `INHERITED / EXTEND_EXISTING_OWNER / REPLACE_CONSOLIDATE / NEW_RESIDUAL / EXTERNAL_BLOCKED / N_A`;
- identify the canonical owner before design;
- mark missing evidence UNVERIFIED rather than inferring implementation from prose.

Apply `REUSE -> CONSOLIDATE -> REFACTOR -> DELETE/REPLACE -> ADD` before creating a new abstraction.

## 3. Scope and design gates

### G0 — Exact Baseline

Freeze predecessor, Stable identity, active head, source fingerprint, open defects, current provider/data/platform conditions and last valid evidence.

### G1 — Immutable Scope

Freeze the coherent version scope and explicit exclusions. Future audit work cannot bypass an unresolved active dependency. New material product scope requires an approved decision.

### G2 — Architecture & Data Utility

For every affected datum/capability prove:

`source -> canonical owner -> consumers -> purpose -> freshness -> retention -> cost -> failure/recovery`

Reject duplicate engines, orphan data, unbounded polling, needless persistence and user-visible Data Engine plumbing.

### G3 — Design & Dependency Readiness

Freeze contracts, migrations, compatibility/shadow plan, provider/rights dependencies, temporal semantics, security boundary, load budget, platform matrix, test oracle and rollback before expensive implementation.

## 4. Compatibility-first authority migration

Any move of renderer logic, `RuntimeSnapshot`, Opportunity state, Watchlist/workspace state, provider paths, full-snapshot SSE, auth/session or persistence authority uses:

`CHARACTERIZE -> NEW OWNER -> DUAL/SHADOW -> COMPARE -> MIGRATE -> PROVE -> RETIRE -> REGRESSION`

Required behavior:

- freeze golden vectors and known defects separately;
- version new contracts;
- preserve old consumers through adapters during migration;
- compare equivalent inputs at the same as-of cutoff;
- explicitly approve intentional behavior corrections;
- keep rollback until consumer and package evidence proves the new owner;
- remove the old authority only after all consumers migrate.

## 5. Canonical intelligence process

Implement in this order:

`observations/events -> rights/quality -> deterministic evidence/features -> SymbolIntelligenceSnapshot -> Opportunity Lifecycle -> projections/Decision Brief -> outcomes -> governed adaptation`

Guardrails:

- deterministic calculations and lifecycle invariants are server-owned;
- Watchlist changes universe/user intent only;
- Discovery is broad universe; Radar/Rapid Move submit evidence;
- Research uses the exact frozen snapshot/transition that caused promotion or alert;
- Alerts deduplicate material transitions/incidents and never rescore;
- LLMs synthesize rights-filtered evidence IDs after deterministic/statistical layers;
- no consumer may bypass canonical state with a page-specific provider fetch or duplicate score.

## 6. Temporal and historical process

For each material fact/event define source, observed, ingested, effective, as-of, expiry/half-life, session/timezone, provider/dataset, rights, quality, revision/supersedes and instrument identity.

- retrieval/cache time is not evidence time;
- unknown time remains unknown;
- macro/fundamental vintages are retained;
- raw/adjusted basis is explicit;
- OI as-of differs from live option quote/IV time;
- late/corrected events trigger revisions and dependent re-evaluation;
- replay uses strict point-in-time joins;
- incomplete/halts/missing/delisted windows may produce CENSORED, not false PASS/FAIL.

## 7. Adaptive and AI process

Keep deterministic rules, registered statistical/adaptive models and LLM synthesis as separate authorities.

Controlled sequence:

`freeze feature/evidence/regime/rights/policy -> decision -> outcome -> point-in-time evaluation -> time-split challenger -> shadow comparison -> drift/subgroup/sample gates -> human approval -> bounded rollout -> rollback`

Evaluate calibration, ranking lift, precision at attention capacity, false-promotion cost, stability, uncertainty and sample counts. Include regime/liquidity/catalyst/coverage cohorts, censored outcomes, deterministic baselines and selection-bias controls.

Prompt/schema success is not factual correctness. Test grounding, unsupported claims, prompt injection, source isolation, rights filtering and deterministic non-AI fallback separately.

## 8. Provider and Data Health process

All provider/data changes conserve #80/#81/#82/#83/#78/#84 and the canonical matrix/SLO/fetch-path artifacts.

- Smart Provider Router v2 remains sole general routing/admission authority;
- direct authorities such as SEC/EDGAR remain explicit;
- Adaptive Provider Registry only projects adapters/capabilities;
- new paths are classified before production use and missing required classification must fail closed;
- canonical freshness uses provider observation/event/publication/filing time;
- lawful cache/fallback is reused before degradation;
- unresolved required evidence stays `PARTIAL COVERAGE` / `DATA DEGRADED` at the smallest truthful scope;
- recovery uses hysteresis and re-evaluates dependents;
- upstream provider correlation and revision lineage are explicit;
- lifecycle remains `SHADOW -> VALIDATED -> APPROVED -> PRODUCTION` without automatic authority or rights promotion.

Provider adoption sequence:

`overlap audit -> adapter/manifest -> secret metadata -> fixtures -> capability/entitlement probe -> rights/Data Health classification -> Registry -> Router proof -> canonical-state cross-integration -> SHADOW -> adverse/recovery -> governed promotion`

## 9. Security, persistence and platform process

For hosted/client changes prove applicable:

- tenant/account isolation independent of role, product entitlement and provider rights;
- managed/OS secrets, rotation/revocation and zero leakage;
- session/refresh/reauth/CSRF/device/deep-link controls;
- tenant Postgres keys/RLS disposition, migrations, outbox, conflicts, tombstones and retention;
- real backup/PITR/restore/failover evidence when claimed;
- privacy deletion cannot be silently reversed by restore;
- old-client API/event compatibility and forced-upgrade policy;
- encrypted last-known desktop cache is visibly stale/offline;
- Mac Apple Silicon, Windows x64 and Web semantic parity;
- source/build/artifact/environment provenance and supply-chain evidence.

No Execution remains permanent throughout hosted, provider, Watchlist, AI and platform work.

## 10. Development and qualification loop

### G4 — Development Exit

Run focused formatting, static, unit, schema, migration and adverse tests alongside source. Update regression ownership and machine traceability with implementation.

### G5 — FAST Qualification

Use one coherent exact-head Fast checkpoint after cheap local gates pass. Fast never dispatches Release.

### G6 — Integration & MEDIUM Qualification

Exercise cross-owner/data/provider/persistence/client integration selected by impact and risk.

### G7 — Data, Security & Adaptive Intelligence

Prove provenance/freshness/rights, point-in-time truth, tenant/security boundaries, failure states, adaptive separation and non-AI fallback.

### G8 — Performance, Capacity & Stability

Measure provider calls/credits, CPU/memory/goroutines/queues/DB I/O, cache value, fanout payloads, latency, recovery and long-running behavior. Protect critical evidence and shed optional work first.

### G9 — Cross-Module UI/UX

Prove one product model across Dashboard, Market Intelligence, Desks, Discovery, Watchlist, Radar, Alerts and Research; correct responsive/accessibility behavior; concise explanation/freshness/contradiction; role-gated internals.

### G10 — Pre-Freeze Qualification

Reconcile the whole version: requirements, audit register, 5/5 rows, cross-integration, compatibility migration, platforms, adverse behavior, professional acceptance and exact source identity. Qualified must pass on the candidate head.

### G11–G16

Freeze, certify, package, runtime-audit, promote only when authorized, then perform AIPLC/G16 handoff. Source-changing repairs require a new fingerprint and affected requalification.

## 11. Testing and CI efficiency

- use deterministic fixtures, clocks and point-in-time replay for normal provider/intelligence CI;
- use bounded live tests only when real transport, entitlement, hosted recovery or platform mechanics are the evidence subject;
- batch coherent changes and use impact routing;
- preserve independent PASS only when its fingerprint/input/environment contract is unchanged;
- classify product, test, harness and infrastructure failures distinctly;
- fix recurrence in canonical code/tests/gates instead of creating retry branches or workflows;
- never weaken tests to obtain PASS.

## 12. AIPLC execution

After a meaningful checkpoint, record affected surfaces, findings, root causes, changes, prevention, ten-dimension scores/residuals, metrics and next disposition.

Required defect chain:

`symptom -> root cause -> canonical owner -> immediate fix -> recurrence prevention -> cross-product scan -> regression/gate -> follow-up metric`

AIPLC may change the next Build Plan but cannot change G1, promote an adaptive model, authorize Commercial/Public use or silently modify permanent scope.

## 13. Process exit and handoff

Before ending meaningful work:

- commit one coherent change to the active branch/PR;
- update machine state/closure evidence when status actually changes;
- update `handoff/CURRENT.md` with exactly one next action;
- record current head/check status truthfully;
- preserve immutable Stable/provenance history;
- leave enough GitHub-backed evidence for a fresh Codex, ChatGPT, Claude or human contributor to resume without conversation memory.
