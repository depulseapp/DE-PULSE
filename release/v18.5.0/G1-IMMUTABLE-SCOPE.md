# DE.PULSE v18.5.0 — G1 Immutable Scope

Status: PASS

Canonical release title: **Major Closure & Release Assurance**.

v18.5 is mandatory before v19 and is a closure/certification release, not a feature-expansion release.

## Frozen closure dimensions

1. Reconstruct and trace the full approved v18.0–v18.4 scope against the immutable v18.4.0 Stable baseline.
2. Fresh architecture and developer-proof source-quality review, including ownership clarity, duplicate-work avoidance, dead/obsolete-code hygiene and maintainability without functional drift.
3. Fresh Data Utility / Correlation review: every material dataset, computation and prominent UI consumer must retain purpose, provenance, freshness/materiality and reuse.
4. Fresh performance, capacity, responsiveness and long-running stability qualification under realistic supported activity.
5. Fresh security/auth/session/authorization/CSRF/cookie/quota/abuse regression and hosted-state isolation review.
6. Fresh UI/UX/content/responsive/accessibility/runtime-truth review on actual supported surfaces.
7. Fresh adaptive-intelligence governance review, including point-in-time evidence/outcome lineage and SHADOW → VALIDATED → APPROVED → PRODUCTION boundaries where applicable.
8. Principal Engineer architecture acceptance and Professional Trader/Investor decision-support acceptance.
9. Fresh macOS Apple Silicon, Windows x64, desktop SQLite and hosted PostgreSQL 17 actual-runtime/package/provenance certification.
10. Fresh exact-source G11–G15 release assurance, Stable promotion, branch hygiene, retrospective and G16 handoff.

## Mandatory ADR-GDI closure dimension

v18.5 MUST prove, under realistic supported load and failure injection, that DE.PULSE itself does not materially create broad, unexplained or misleading `DATA DEGRADED` states. Closure must exercise and evidence:

- provider failure, provider rate limits, retry/circuit/fallback behavior and calls avoided;
- stale evidence, freshness-SLO breach and source disagreement;
- capability/consumer-aware degradation blast radius rather than unnecessary whole-app degradation;
- PostgreSQL pressure/slow/unavailable behavior and desktop SQLite continuity where applicable;
- queue saturation, bounded backpressure, workload priority and graceful load shedding;
- restart/warm-start recovery from canonical persisted state;
- multi-user and multi-symbol fan-out with shared canonical computation/fetch reuse;
- background-job pressure and duplicate-fetch/calculation avoidance;
- recovery hysteresis so state does not flap between healthy/degraded;
- truthful `UNKNOWN`, degraded and `ABSTAIN` semantics when required evidence is insufficient;
- actual packaged-runtime degradation UX plus deeper Maintenance diagnostics;
- explicit, evidence-backed supported operating limits when a bounded limit is necessary.

**Release blocker:** if local/self-inflicted overload can materially delay, stale, misstate or hide decision-critical current evidence, or if a narrow dependency failure unnecessarily degrades unrelated capabilities, v18.5 cannot promote until corrected or explicitly constrained with truthful operating limits.

## Preserved protected contracts

- Permanent No Execution Boundary: no paper trading, order routing, portfolio/P&L or execution features.
- Protected deterministic Day/Swing/Long formulas must not silently change.
- Smart Provider Router remains canonical; no duplicate router or rights-driven routing mutation.
- GLD, SLV and USO remain approved tradable live-priority exceptions.
- Shared Symbol Intelligence / Global Symbol Registry reuse remains mandatory; per-user state must not cause duplicate market-wide work.
- Desktop remains SQLite/local by default; hosted shared state remains PostgreSQL where configured.
- Security/commercial/data-rights readiness remains fail closed; v18.5 makes no provider licensing/legal determination.
- Adaptive intelligence may not silently self-modify production behavior.

## Explicit exclusions / anti-scope-creep

- Do not force mature ASBI into v18.5.
- Do not force mature Two-Sided Directional Thesis & Trade Plan Intelligence into v18.5.
- Do not force mature Adaptive Opportunity Discovery & Recommendations into v18.5.
- Do not force mature adaptive Institutional Holdings / 13F Intelligence into v18.5.
- Do not introduce v19 professional-data-infrastructure scope early unless a minimal dependency-compatible correction is required to close a proven v18 defect.
- No visual redesign unrelated to closure defects.

Dependency-compatible foundations already present may be validated and hardened, but mature adaptive engines remain governed later-roadmap work.
