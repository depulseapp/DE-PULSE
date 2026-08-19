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
**Engineering branch:** `v18.6.5-development`  
**Current PR:** Packet D Draft PR from `v18.6.5-development` to `main`; resolve the live PR number from GitHub rather than model memory  
**Current slice:** Phase 0 Packet D — durable CI evidence and telemetry; process-only, not a product Stable build  
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
- `release/v18.6.1/stable-evidence-manifest.json`

## Immutable v18.6.1 Stable evidence

- Final exact-head Fast: run `32279055139` · PASS.
- Final exact-head Qualified: run `32279113032` · PASS.
- Product candidate full Qualified: run `32276304863` · backend/race/randomized, renderer, deterministic 2403/2403 and browser/global-remove coverage PASS.
- Canonical Release #24: run `32279232665` · G11–G16 PASS including macOS Apple Silicon and Windows x64 native evidence and same-run publication.
- Stable tag: `v18.6.1-stable` → `42e8432f7530ae39cbfd6ceb0b0bd5f6311dc5cc`.
- Measured v18.6.1 clock: first Draft PR to final Stable = **1h 56m 20s**; final canonical G11–G16 run = **6m 43s**.
- Packet D adds a retrospective compact Stable evidence manifest bound to the authoritative release-evidence checkpoint. It does not redefine the immutable Stable tag/binaries.

Permanent boundaries remain unchanged: U.S. Equities Processing, No Execution, deterministic Day/Swing/Long truth, Smart Provider Router ownership, provider/Market-Mode rules, GLD/SLV/USO tradable exceptions and Adaptive Intelligence governance.

## Phase 0 hardening status

### Packet A — COMPLETE

Merged PR #46 at `a3beb28322c2c53227bac037e546d6863c8d279e`.

Delivered Impact Planner v2, planner self-test, side-effect-free Release Rehearsal, workflow-policy integration, current adaptive operating overlays and an honest v18.6.1 reconciliation baseline. Final Fast #381 and Qualified #133 passed.

### Packet B — COMPLETE

Merged PR #47 at `31b697d175317f0132c0a3fff7283beb1b79662d`.

Delivered immutable Action pins for Fast/Qualified, canonical CI dependency lock, Playwright `1.62.0` pin, safe pip caching and reproducibility/permission gate. Final Fast #384 and Qualified #134 passed. `release.yml` Action pinning remains intentionally deferred to the next genuine release-capable product slice.

### Packet C — COMPLETE

Merged PR #48 at `23ecb71f60e1658d68bcef6248044ce53b6dd851`.

Browser policy:

- Chrome and WebKit are the two primary engines.
- Chrome carries broad behavioral regression.
- WebKit carries co-primary cross-engine compatibility.
- `full`/`browser` qualification requires both; renderer/UI and WebKit-harness risk requires WebKit.
- backend/provider/process-only work avoids browser runtime when unaffected.
- other engines remain secondary/risk-directed.

Final Packet C exact-head evidence:

- Fast #393 / run `32296746701`: PASS.
- Qualified #138 / run `32296793338`: PASS.
- CI-harness: PASS.
- Ubuntu/macOS/Windows portability: PASS.
- real macOS WebKit core compatibility: PASS.
- unrelated backend/renderer/Chrome product lanes: correctly SKIPPED for the process-only Packet C candidate.

Packet C also proved the failure taxonomy in practice: initial Linux WebKit dependency amplification was `CI_HARNESS_FAIL`; two subsequent WebKit assertion failures were production-fixture mismatches in the new test harness and were fixed without changing product runtime or weakening assertions.

### Packet D — ACTIVE

Branch `v18.6.5-development` is assembled before PR opening to avoid preparatory `synchronize` amplification.

Scope:

1. Durable `release/v18.6.1/stable-evidence-manifest.json` bound to the authoritative checkpoint.
2. `stable_evidence_gate.py` verifies Stable candidate, source fingerprint, run IDs and native/G15/G16 artifact digests.
3. `ci_telemetry.py` + self-test records per-job queue/runtime/platform consumption.
4. Qualified records Linux/macOS/Windows runner minutes, Chrome/WebKit setup duration and pip cache-hit state when applicable.
5. Qualified counts current-PR Fast/Qualified/Release runs and warns on abnormal amplification without blocking legitimate fixes.
6. Compact telemetry JSON retained 30 days plus job summary.
7. Currency cost is intentionally not invented; GitHub billing remains authoritative.
8. `workflow_structural_lint.py` provides zero-network generic structural checks in addition to semantic workflow policy.
9. Impact Planner treats only the exact retrospective Stable manifest path as process-only; executable release files remain full-qualification scope.
10. Current Roadmap / Build Plan / Build Process / Delivery Process / CI Efficiency Contract are synchronized to this operating model.

Expected Packet D qualification: process-only `ci-harness` + Ubuntu/macOS/Windows portability. WebKit should remain skipped because Packet D observes browser metrics but does not change renderer/WebKit behavior. No Stable Release is expected because neither `release_identity.json` nor `.github/workflows/release.yml` changes.

## Remaining Phase 0

**Packet E — Renderer Modularization Foundation:** incremental strangler extraction from the large renderer/compatibility stack, preserving deterministic equivalence and required Chrome/WebKit/native evidence before deleting former owners.

## After Phase 0

Adaptive v18.x priority order: user-trust defects → runtime/ADR-GDI reliability → shared intelligence utility consolidation → renderer maintainability/consolidation → controlled TradeInsight SHADOW integration. Then v18 Major Closure, followed by v19 Professional Data Infrastructure and v20 Adaptive Intelligence.

## Repository/source-age rule

Do not modify or delete files because GitHub says they are several days old. Historical release evidence, governance and actively loaded compatibility assets stay until references/consumers/evidence requirements prove safe removal.

## External/open governance item

Repository `main` branch protection/ruleset is still disabled as of the post-Packet-C live check. Do not report it as fixed unless GitHub confirms an authorized branch-protection/ruleset change.

## Exactly one next action

Open exactly one Draft PR from fully assembled `v18.6.5-development` to `main`. Require one exact-head Fast candidate. Only after Fast passes, mark the same PR Ready and require Qualified `ci-harness` + portability + telemetry evidence. Fix any legitimate defect on the same branch/PR; do not create retry/certification/promotion branches.
