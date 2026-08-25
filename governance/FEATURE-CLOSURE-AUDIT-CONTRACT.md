# DE.PULSE Feature Closure Audit Contract

**Status:** PERMANENT / MANDATORY FOR MAJOR OR FINAL CLOSURE BUILDS  
**Owner:** Adaptive Delivery / closure governance  
**Purpose:** Make every closure build perform a source-driven audit of every shipped feature, not merely a test rerun or documentation inventory.

## 1. Closure-audit invariant

Every major/final closure build must begin with a **Feature Reality & Quality Audit** over the exact current source-of-truth head.

The audit is exhaustive across shipped capabilities, surfaces, data paths, background jobs, roles/security paths, persistence/state, runtime/platform behavior, release responsibilities, and materially active compatibility paths.

A closure may not claim `10/10`, future-proof, production-ready, or equivalent merely because:
- an issue is closed;
- a test exists;
- CI is green;
- a feature appears in the UI;
- a historical audit said it was good;
- documentation says it is implemented.

The current source, executable behavior, current tests, current runtime evidence, and current GitHub history must reconcile.

## 2. What must be audited for every feature

Each feature row must answer all applicable questions below.

### A. Reality / traceability
- What exact shipped responsibility exists?
- What requirement, issue, PR, decision, or inherited Stable behavior created it?
- What is the canonical source owner?
- Which consumers depend on it?
- Which tests and executable evidence own its regression?
- Is any apparently related code actually a separate canonical responsibility?

### B. Product utility
- What user, operator, research, reliability, security, or delivery problem does it solve?
- Is that utility still real in the current product?
- Is the feature redundant, superseded, too prominent, misplaced, or no longer worth its complexity?
- Does another feature already provide the same outcome better?

### C. Correctness and data truth
- Does the implementation match the current requirement and authority model?
- Are missing, stale, partial, contradictory, unavailable, delayed, cached, degraded, or future-dated states represented truthfully?
- Can fallback/recovery accidentally fabricate health, zero values, certainty, or authority?
- Are direct authorities, secondary evidence, proxies, alternative evidence, and provider lifecycle states kept distinct?

### D. Architecture and ownership quality
- Is there exactly one canonical owner for the responsibility?
- Does the feature reuse Smart Provider Router v2, canonical freshness/degradation/cache/persistence/subscription/telemetry/reconciliation/lifecycle/identity/session owners where applicable?
- Is there duplicate acquisition, state, routing, recovery, persistence, telemetry, scheduler, or UI logic?
- Is a compatibility layer still justified, or has it become a second architecture?

### E. Intelligence / adaptive fit
Audit the **right intelligence level**, not the maximum possible intelligence level.

Allowed expectations:
- `DETERMINISTIC_REQUIRED` — security, authority, persistence, release/provenance, validation and other fail-closed rules should remain explicit/deterministic.
- `CONTEXT_AWARE` — behavior should react to session/state/material context without learning from outcomes.
- `ADAPTIVE_POLICY` — bounded policy may adjust cadence, priority, provider usage, load shedding, or similar behavior from current evidence.
- `EVIDENCE_LEARNING` — durable outcomes/quality history should improve future prioritization, synthesis, provider usefulness, ranking, or explanations.
- `NOT_APPLICABLE` — no intelligent behavior is useful for this responsibility.

For each row determine:
- current intelligence level;
- expected intelligence level;
- whether it is under-intelligent;
- whether it is over-engineered or unnecessarily adaptive;
- whether adaptive behavior is bounded, explainable, reversible and evidence-driven;
- whether learning uses durable outcome evidence instead of hidden or unbounded self-modification.

A feature is not improved by adding AI/LLM/adaptation where deterministic behavior is safer or clearer.

### F. Maintainability / stale-code audit
Inspect implementation and supporting code for:
- unreachable or unreferenced production paths;
- stale version-specific branches or compatibility code;
- obsolete feature flags;
- duplicate helpers/normalizers/state owners;
- dead provider adapters or dormant routes presented as current;
- historical tests/gates still acting as active owners;
- misleading comments/names that no longer match runtime truth;
- oversized files/functions or coupling that materially obstruct safe change;
- old UI routes/components retained after canonical replacement;
- code whose only justification is historical but has no current consumer.

Source disposition must be explicit:
`KEEP`, `CLEANUP`, `REFACTOR`, `MERGE`, `REMOVE`, or `REWRITE`.

Do not remove code only because it looks old. Removal requires consumer/reachability/equivalence evidence. Do not rewrite working code merely for style; `REWRITE` requires a concrete correctness, ownership, safety, maintainability, performance, or future-evolution reason plus regression/equivalence protection.

### G. UX / information architecture
For visible behavior determine:
- whether the feature belongs where it is;
- whether terminology is current and understandable;
- whether duplicate surfaces should be merged or redirected;
- whether loading/empty/degraded/error/permission states are truthful;
- whether control placement, hierarchy, responsive behavior and accessibility are appropriate.

UI disposition remains explicit:
`KEEP`, `MOVE`, `MERGE`, `REMOVE`, `RENAME`, `REDESIGN`, or `NOT_USER_VISIBLE`.

### H. Security / rights / privacy
Where applicable verify:
- role/capability and direct-route/API parity;
- session/reauth/revocation behavior;
- no plaintext or accidental secret exposure;
- provider entitlement and data-rights boundaries;
- authority and alternative-evidence boundaries;
- per-user/workspace isolation;
- export/import and diagnostics do not leak secrets or another user's state.

### I. Persistence / restart / lifecycle
Where applicable verify:
- save/load/atomicity;
- restart continuity;
- migration/upgrade behavior;
- cache versus durable truth separation;
- interrupted-write/corrupt-state behavior;
- packaged shutdown/restart lifecycle;
- no duplicate persistence owner.

### J. Performance / efficiency
Where applicable verify:
- bounded concurrency, queues and fan-out;
- rate-limit awareness and provider budget behavior;
- coalescing/cache reuse;
- protected market-session priority;
- no broad refresh where targeted work suffices;
- no hidden background job that duplicates another owner;
- UI remains responsive under realistic state volume.

### K. Observability / diagnosability
Determine whether operators and downstream logic can truthfully distinguish:
- healthy versus degraded/unavailable/stale;
- provider failure versus local load versus entitlement versus rights;
- cached/recovered versus live/direct;
- expected idle versus broken scheduler;
- feature failure without creating app-global false degradation.

### L. Improvement decision
Every row receives an explicit closure decision:
- `READY_10_10` — no material feature-specific improvement is required for this closure and all applicable later evidence tracks are owned.
- `IMPROVE_BEFORE_CLOSURE` — implementation/test/UX/architecture/data-truth/security/performance cleanup is required before final closure.
- `REFACTOR_BEFORE_CLOSURE` — behavior may be correct but ownership/maintainability creates unacceptable closure risk.
- `REWRITE_BEFORE_CLOSURE` — current design cannot safely meet the closure contract through bounded correction.
- `MERGE_OR_REMOVE_BEFORE_CLOSURE` — duplicate/stale responsibility should not survive the closure unchanged.
- `DEFER_NON_BLOCKING_IMPROVEMENT` — worthwhile improvement exists but is not required for truthful safe current-version closure; rationale and future owner are mandatory.
- `EXTERNAL_BLOCKED_WITH_EVIDENCE` — a required proof depends on an evidenced external constraint; no false `10/10` claim is allowed while blocked.

## 3. 10/10 rule

`10/10` is **not an arithmetic average** and is not a claim of perfection.

A feature may be declared `READY_10_10` only when every applicable closure dimension is either:
1. evidenced as satisfactory for the current responsibility; or
2. explicitly `NOT_APPLICABLE` with a defensible reason.

A serious weakness in one dimension cannot be averaged away by strengths elsewhere.

Any unexplained P0/P1 correctness, data-truth, security, authority, persistence, platform, duplicate-owner, stale-code, or maintainability gap blocks the feature and therefore blocks final closure.

## 4. Required machine fields per feature

The release-specific feature assurance ledger must be able to resolve at least:
- `id`, `name`, `category`;
- `requirementProvenance`;
- `canonicalSourceOwners`;
- `consumers`;
- `existingRegressionOwners`;
- `positiveFunctionalEvidenceExpectation`;
- downstream T2-T9 assurance ownership/profile;
- `productUtilityAssessment`;
- `correctnessAssessment`;
- `architectureAssessment`;
- `intelligenceAssessment` (`currentLevel`, `expectedLevel`, `fit`, `gaps`);
- `dataTruthAssessment` where applicable;
- `maintainabilityAssessment`;
- `staleCodeAssessment`;
- `sourceDisposition`;
- `uiDisposition`;
- `securityRightsAssessment` where applicable;
- `persistenceLifecycleAssessment` where applicable;
- `performanceEfficiencyAssessment` where applicable;
- `observabilityAssessment` where applicable;
- `improvementOpportunities`;
- `closureDecision`;
- `correctiveIssue` when closure-blocking remediation is required;
- `deferredImprovementOwner` when a non-blocking improvement is intentionally deferred;
- `durableRegressionOwner`;
- `currentAssuranceState` and explicit blocking states.

## 5. Required audit method

A closure Feature Reality & Quality Audit must use multiple independent discovery directions:

1. **Feature/requirement -> source/tests**: start from registry, roadmap, issues, PRs, Stable evidence and UI surfaces; find implementation and evidence.
2. **Source/runtime -> feature ledger**: traverse current production source, routes, jobs, provider registrations, state/persistence, renderer and platform/release code; every shipped responsibility must map back to a row.
3. **Tests/evidence -> owner**: active regression tests, acceptance suites and release gates must map to a current responsibility; orphaned active tests are debt signals.
4. **Duplicate/stale scan**: independently inspect parallel owners, legacy routes, compatibility paths, version-stacked active code, comments/names and unused responsibilities.
5. **Improvement scan**: assess whether each feature is merely present or is actually fit, efficient, truthful, understandable, maintainable and at the correct intelligence level.

At least one independent omission/duplicate/stale scan is mandatory after the initial inventory is populated.

## 6. Corrective rule

If the audit finds a genuine shipped implementation miss or a closure-blocking defect, create a named corrective under the active closure parent. Do not hide the miss in the ledger and do not mark it `VERIFIED` because a historical issue was closed.

If the audit finds stale/duplicate/unnecessary code, cleanup may occur inside the closure only when current behavior and required compatibility are protected by executable evidence. Material corrections must be requalified through the normal delivery process.

## 7. Relationship to later assurance tracks

The feature audit does not replace deep assurance tracks. It defines what they must prove.

Typical closure order:
`Feature Reality & Quality Audit -> unit/contract -> functional/E2E -> adverse/data-truth -> persistence/restart -> security/rights -> UI/UX -> performance/load -> packaged platform/release -> durable regression/final closure`.

Later tracks must feed discoveries back into the same feature rows. A row may regress from `READY_10_10` if later evidence exposes a miss.

## 8. Future closure reuse

This contract is version-neutral and remains active for future major/final closure programs. A future closure must not create a weaker one-off interpretation of T1. It may extend the audit dimensions, but it may not silently drop reality, utility, correctness, intelligence-fit, stale-code, cleanup/rewrite, or improvement assessment.

GitHub source, executable evidence and current program objects remain authoritative. Chat history is never required to reconstruct the audit.