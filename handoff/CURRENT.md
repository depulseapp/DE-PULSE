# DE.PULSE — Current Authoritative Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS.** GitHub is authoritative; chat history and temporary workspaces are advisory only.

**Certified Stable:** `v18.6.1-stable`  
**Stable target / certified merged candidate:** `42e8432f7530ae39cbfd6ceb0b0bd5f6311dc5cc`  
**Qualified v18.6.1 source head:** `5c3fae486f3e4b4a39a0b1d549916aea9e1295fd`  
**Certified source fingerprint:** `b01c14e1d54b736785eab6c03407801c527edd7769ff6f3d41fd4b20dabebd75`  
**Build ID:** `v18.6.1-stable-20260819`  
**Canonical v18.6.1 certification/publication run:** `32279232665` (Release #24)  
**Repository:** `depulseapp/DE-PULSE`  
**Active engineering branch:** `v18.6.2-development`  
**Current slice:** post-v18.6.1 engineering hardening; process/governance only, not a product Stable build  
**Last updated:** 2026-08-19 America/Vancouver

## Resume rule

Read `AGENTS.md` or `CLAUDE.md`, `governance/CI-EFFICIENCY-CONTRACT.md`, this handoff, `release_identity.json`, both `.depulse-certification/resume/` checkpoints, current GitHub PR/issue state and immutable `v18.6.1-stable`. Never resume from model memory alone.

For current operating state also read:

- `adaptive-governance/CURRENT_ADAPTIVE_ROADMAP.md`
- `adaptive-governance/CURRENT_ADAPTIVE_BUILD_PLAN.md`
- `adaptive-governance/CURRENT_ADAPTIVE_BUILD_PROCESS.md`
- `adaptive-governance/CURRENT_ADAPTIVE_DELIVERY_PROCESS.md`
- `adaptive-governance/V18_6_1_CURRENT_RECONCILIATION.md`

Older adaptive files remain permanent/historical evidence, but older statements that v18.2 or v18.5.1 is the active release are historical and must not override these current overlays.

## Immutable v18.6.1 Stable evidence

- Final exact-head Fast: run `32279055139` · PASS.
- Final exact-head Qualified: run `32279113032` · PASS using the intended process-only CI-harness portability lane.
- Product candidate full Qualified: run `32276304863` · backend/race/randomized, renderer, deterministic 2403/2403 and browser/global-remove edge coverage PASS.
- Canonical Release #24: run `32279232665` · G11, G12, macOS Apple Silicon G13/G14, Windows x64 G13/G14, G15, same-run no-rebuild publication and G16 PASS.
- Stable tag: `v18.6.1-stable` → `42e8432f7530ae39cbfd6ceb0b0bd5f6311dc5cc`.
- Measured v18.6.1 clock: first Draft PR to final Stable = **1h 56m 20s**; final canonical G11–G16 run = **6m 43s**.

Permanent boundaries remain unchanged: U.S. Equities Processing, No Execution, deterministic Day/Swing/Long truth, Smart Provider Router ownership, provider/Market-Mode rules, GLD/SLV/USO tradable exceptions and Adaptive Intelligence governance.

## Current post-Stable hardening slice

Purpose: turn the v18.6.1 lessons into permanent process safeguards before broad feature work.

Implemented on `v18.6.2-development`:

1. `tools/ci/impact_plan.py` schema 2 adds explicit change classes, conservative full-lane fallback, Release Rehearsal signal, targeted-WebKit planning signal and canonical failure taxonomy while retaining existing workflow outputs.
2. `tools/ci/impact_plan_self_test.py` proves process-only, renderer, provider/data-rights and mixed fail-closed classifications.
3. `tools/ci/release_rehearsal.py` performs side-effect-free pre-merge checks of G11–G16 topology, exact-head/fingerprint requirements, Stable-tag absent/same/mismatch behavior and no-rebuild publication.
4. `tools/ci/workflow_policy.py` invokes both new gates, so Fast and the process-only Qualified harness automatically exercise them without adding another workflow.
5. Current adaptive roadmap/build-plan/build-process/delivery-process overlays replace stale current-state claims without rewriting historical/permanent contracts.
6. `V18_6_1_CURRENT_RECONCILIATION.md` records an honest baseline. It explicitly does not invent the historical 296 ledger rows; the original conserved IDs/history must be preserved when the source ledger is located.

This slice intentionally does **not** change `release_identity.json` or `.github/workflows/release.yml`. Therefore merging it should not trigger a new Stable Release. The expected validation path is Fast → same PR Ready → Qualified `ci-harness` + portability → merge → main-push hygiene.

## Next hardening packets after this slice

1. Verify and pin third-party GitHub Actions to immutable SHAs; pin Playwright/browser; add safe caches and generic workflow linting; tighten least-privilege permissions.
2. Connect `webkit_required` to a focused WebKit Qualified lane for renderer/UI-sensitive changes.
3. Add durable compact Stable evidence manifest and CI runtime/queue/cache/failure telemetry.
4. Start incremental renderer strangler modularization only after the process hardening is green.

## Repository/source-age rule

Do not modify or delete files because GitHub says they are several days old. Age is not evidence of obsolescence. Historical release evidence, governance and actively loaded compatibility assets stay until runtime references/consumers and tests prove safe removal.

## External/open governance item

Repository `main` branch protection/ruleset was previously observed as not protected and has not been claimed fixed by this slice. Do not report it as changed unless GitHub confirms a ruleset/branch-protection update through an authorized mechanism.

## Exactly one next action

Open one Draft PR from `v18.6.2-development` to `main`, let Fast validate the exact head, mark that same PR Ready only after Fast passes, then require Qualified `ci-harness` + portability. Fix any legitimate source/gate defect on the same branch/PR; do not create retry/certification/promotion branches.
