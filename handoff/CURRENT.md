# DE.PULSE — Current Authoritative Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS.** GitHub is authoritative; chat history and temporary workspaces are advisory only.

**Release:** `v18.6.1`  
**Active branch:** `main`  
**Certified Stable:** `v18.6.1-stable`  
**Stable target / certified merged candidate:** `42e8432f7530ae39cbfd6ceb0b0bd5f6311dc5cc`  
**Qualified v18.6.1 source head:** `5c3fae486f3e4b4a39a0b1d549916aea9e1295fd`  
**Certified source fingerprint:** `b01c14e1d54b736785eab6c03407801c527edd7769ff6f3d41fd4b20dabebd75`  
**Build ID:** `v18.6.1-stable-20260819`  
**Canonical v18.6.1 certification/publication run:** `32279232665` (Release #24)  
**Repository:** `depulseapp/DE-PULSE`  
**Engineering branch:** `v18.6.3-development`  
**Current slice:** Phase 0 Packet B — CI reproducibility hardening; process-only, not a product Stable build  
**Last updated:** 2026-08-19 America/Vancouver

`Release` and `Active branch` intentionally mirror the immutable build-resume checkpoint (`v18.6.1` / `main`). `Engineering branch` records in-flight engineering without mutating Stable checkpoint identity.

## Resume rule

Read `AGENTS.md` or `CLAUDE.md`, `governance/CI-EFFICIENCY-CONTRACT.md`, this handoff, `release_identity.json`, both `.depulse-certification/resume/` checkpoints, current GitHub PR/issue state and immutable `v18.6.1-stable`. Never resume from model memory alone.

For current operating state also read:

- `adaptive-governance/CURRENT_ADAPTIVE_ROADMAP.md`
- `adaptive-governance/CURRENT_ADAPTIVE_BUILD_PLAN.md`
- `adaptive-governance/CURRENT_ADAPTIVE_BUILD_PROCESS.md`
- `adaptive-governance/CURRENT_ADAPTIVE_DELIVERY_PROCESS.md`
- `adaptive-governance/V18_6_1_CURRENT_RECONCILIATION.md`

## Immutable v18.6.1 Stable evidence

- Final exact-head Fast: run `32279055139` · PASS.
- Final exact-head Qualified: run `32279113032` · PASS.
- Product candidate full Qualified: run `32276304863` · backend/race/randomized, renderer, deterministic 2403/2403 and browser/global-remove coverage PASS.
- Canonical Release #24: run `32279232665` · G11–G16 PASS including macOS Apple Silicon and Windows x64 native evidence and same-run publication.
- Stable tag: `v18.6.1-stable` → `42e8432f7530ae39cbfd6ceb0b0bd5f6311dc5cc`.
- Measured v18.6.1 clock: first Draft PR to final Stable = **1h 56m 20s**; final canonical G11–G16 run = **6m 43s**.

Permanent boundaries remain unchanged: U.S. Equities Processing, No Execution, deterministic Day/Swing/Long truth, Smart Provider Router ownership, provider/Market-Mode rules, GLD/SLV/USO tradable exceptions and Adaptive Intelligence governance.

## Phase 0 hardening status

### Packet A — COMPLETE

Merged PR #46 to `main` at `a3beb28322c2c53227bac037e546d6863c8d279e`.

Delivered:

1. Impact Planner v2 change classes, conservative fail-closed routing, WebKit planning signal and failure taxonomy.
2. Planner self-test.
3. Side-effect-free Release Rehearsal for G11–G16 topology, exact-head/fingerprint contracts, Stable-tag behavior and no-rebuild publication.
4. Workflow policy integration of the planner/rehearsal safeguards.
5. Current Adaptive Roadmap / Build Plan / Build Process / Delivery Process overlays.
6. Honest v18.6.1 reconciliation baseline without inventing missing historical ledger rows.

Packet A final exact-head Fast #381 and Qualified #133 passed; Qualified correctly used CI-harness + Ubuntu/macOS/Windows portability and skipped irrelevant product suites. No new Stable release was triggered.

### Packet B — ACTIVE

Goal: make routine CI dependencies reproducible and auditable without spending a release-only native certification run.

Current branch scope:

1. Pin third-party GitHub Actions in `ci-fast.yml` and `ci-qualified.yml` to immutable 40-hex commits, retaining readable upstream version comments.
2. Add `tools/ci/ci_dependency_lock.json` as the canonical Action/browser dependency record.
3. Pin Playwright to `1.62.0` through `tools/ci/browser-requirements.txt`.
4. Use setup-python's pip cache keyed by the pinned browser requirements file; no separate cache Action is added.
5. Add `tools/ci/reproducibility_gate.py` to reject movable Action refs, dependency-lock drift, unpinned Playwright and scoped permission expansion.
6. Run the reproducibility gate through canonical `tools/ci/workflow_policy.py`.
7. Keep `release.yml` Action pinning explicitly deferred to the next release-capable product slice so changing it does not create a process-only G11–G16/macOS/Windows spend.

Expected Packet B path: Draft PR → Fast exact-head → same PR Ready → Qualified `ci-harness` + portability → merge → main-push hygiene. No Stable Release is expected because neither `release_identity.json` nor `.github/workflows/release.yml` is changed.

## Remaining Phase 0 packets

1. **Packet C — Browser Risk Routing:** connect Impact Planner `webkit_required` to a focused Qualified WebKit lane for renderer/UI-sensitive changes without duplicating backend-only work.
2. **Packet D — Durable CI Evidence & Telemetry:** compact Stable evidence manifest plus lane runtime/queue/cache/failure taxonomy and cost trend evidence.
3. **Packet E — Renderer Modularization Foundation:** incremental strangler extraction from the large renderer/compatibility stack, preserving equivalence and browser/native evidence before deleting former owners.

## After Phase 0

Adaptive v18.x priority order: user-trust defects → runtime/ADR-GDI reliability → shared intelligence utility consolidation → renderer maintainability/consolidation → controlled TradeInsight SHADOW integration. Then perform v18 Major Closure, followed by v19 Professional Data Infrastructure and v20 Adaptive Intelligence.

## Repository/source-age rule

Do not modify or delete files because GitHub says they are several days old. Historical release evidence, governance and actively loaded compatibility assets stay until references/consumers/evidence requirements prove safe removal.

## External/open governance item

Repository `main` branch protection/ruleset remains an external governance item. Do not report it as fixed unless GitHub confirms an authorized branch-protection/ruleset change.

## Exactly one next action

Open one Draft PR from `v18.6.3-development` to `main` for Packet B, require exact-head Fast, then mark the same PR Ready only after Fast passes and require Qualified `ci-harness` + portability. Fix any legitimate defect on the same branch/PR; do not create retry/certification/promotion branches.
