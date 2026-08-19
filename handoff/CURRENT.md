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
**Engineering branch:** `v18.6.4-development`  
**Current slice:** Phase 0 Packet C — Chrome-primary targeted-WebKit browser-risk routing; process-only, not a product Stable build  
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

Merged PR #46 at `a3beb28322c2c53227bac037e546d6863c8d279e`.

Delivered Impact Planner v2, planner self-test, side-effect-free Release Rehearsal, workflow-policy integration, current adaptive operating overlays and an honest v18.6.1 reconciliation baseline. Final Fast #381 and Qualified #133 passed through the intended process-only path.

### Packet B — COMPLETE

Merged PR #47 at `31b697d175317f0132c0a3fff7283beb1b79662d`.

Delivered:

1. Immutable 40-hex SHA pins for third-party Actions in Fast/Qualified with readable version comments.
2. `tools/ci/ci_dependency_lock.json` as canonical CI dependency record.
3. Playwright `1.62.0` pin through `tools/ci/browser-requirements.txt`.
4. Safe setup-python pip cache keyed by the browser requirements lock.
5. `tools/ci/reproducibility_gate.py` for movable Action refs, lock drift, unpinned Playwright and scoped permission expansion.
6. Workflow-policy integration of the reproducibility gate.
7. `release.yml` Action pinning remains intentionally deferred to the next genuine release-capable product slice so native G11–G16 spend is useful.

Packet B had one legitimate gate-parser defect on its first Fast attempt; it was fixed on the same branch/PR. Final Fast #384 and Qualified #134 passed; Qualified used CI-harness + Ubuntu/macOS/Windows portability and skipped product suites. No Stable release was triggered.

### Packet C — ACTIVE

Browser strategy is explicitly usage/risk aligned:

- **Chrome is the primary/default browser behavioral qualification target** and keeps the broad browser regression suite.
- **WebKit is a selective secondary compatibility guardrail**, not a duplicate full matrix.
- Impact Planner `webkit_required=true` only when `RENDERER_UI` changes are present.
- Backend/provider/process-only changes must not incur WebKit runtime.
- Targeted WebKit scope focuses on browser-engine-sensitive contracts: watchlist/global-remove and DESKS/no-CURRENT semantics, failure handling, short-height Settings save-bar behavior and centered header alert layout.
- Exact-head Qualified evidence requires WebKit success whenever the impact signal says it is required; otherwise the WebKit job must be skipped.

Current Packet C files:

1. `.github/workflows/ci-qualified.yml` — propagates `webkit_required`, keeps Chrome as the full browser lane, adds a conditional targeted WebKit job and binds its result into Qualified evidence.
2. `tools/ci/webkit_targeted_test.py` — focused Safari/WebKit compatibility proof.
3. `tools/ci/browser_risk_routing_gate.py` — static routing contract ensuring no full Chrome+WebKit matrix and no Fast/backend-only WebKit coupling.
4. `tools/ci/workflow_policy.py` — permanently enforces browser-risk routing.

Expected Packet C path: Draft PR → Fast exact-head → same PR Ready → Qualified `ci-harness` + portability → merge → main-push hygiene. Because Packet C changes CI tooling only, `webkit_required` should be false for Packet C itself; the routing gate proves the conditional wiring and the first future renderer/UI candidate will execute the WebKit lane. No Stable Release is expected because neither `release_identity.json` nor `.github/workflows/release.yml` changes.

## Remaining Phase 0 packets

1. **Packet D — Durable CI Evidence & Telemetry:** compact Stable evidence manifest, lane runtime/queue/cache/failure taxonomy, CI cost trend evidence, and carry forward generic workflow linting if not already incorporated by the end of Packet C.
2. **Packet E — Renderer Modularization Foundation:** incremental strangler extraction from the large renderer/compatibility stack, preserving equivalence and browser/native evidence before deleting former owners.

## After Phase 0

Adaptive v18.x priority order: user-trust defects → runtime/ADR-GDI reliability → shared intelligence utility consolidation → renderer maintainability/consolidation → controlled TradeInsight SHADOW integration. Then perform v18 Major Closure, followed by v19 Professional Data Infrastructure and v20 Adaptive Intelligence.

## Repository/source-age rule

Do not modify or delete files because GitHub says they are several days old. Historical release evidence, governance and actively loaded compatibility assets stay until references/consumers/evidence requirements prove safe removal.

## External/open governance item

Repository `main` branch protection/ruleset remains an external governance item. Do not report it as fixed unless GitHub confirms an authorized branch-protection/ruleset change.

## Exactly one next action

Open one Draft PR from `v18.6.4-development` to `main` for Packet C, require exact-head Fast, then mark the same PR Ready only after Fast passes and require Qualified `ci-harness` + portability. Fix any legitimate defect on the same branch/PR; do not create retry/certification/promotion branches.
