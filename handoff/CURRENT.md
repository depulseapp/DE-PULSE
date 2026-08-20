# DE.PULSE — Current Authoritative Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS. GitHub is authoritative; model memory is advisory only.**

## Current identity

**Certified Stable:** `v18.8.0-stable`  
**Certified candidate / tag target:** `3a32d57dd4c74c6f812cc942a9d8049a7b517718`  
**Qualified source head:** `a8cf5f4e818609d191f977da846be31203d76f06`  
**Certified source fingerprint:** `fa7b49ec9001d5ef95b829834f6268100e2eaf7c3da6bcc1f1a0b9bcba208d46`  
**Certified Build ID:** `v18.8.0-stable-20260819`  
**Fast:** #419 / `32336519003`  
**Qualified:** #145 / `32336619446`  
**Release:** #28 / `32336898662`  
**Engineering branch:** `v18.8.1-development`  
**Latest completed implementation commit:** `c01bd4d683be6236ec884aaea0cc14f7adf4d0b8` (`ADAPT-DATA-002`)  
**Current engineering line:** `v18.8.1 — Renderer Modularization II + 10/10 Audit Hardening`.

## v18.8.0 Stable closure

v18.8.0 Shared Intelligence Consolidation is complete and immutable. `v18.8.0-stable` resolves exactly to candidate `3a32d57dd4c74c6f812cc942a9d8049a7b517718`. Release #28 passed G11, authoritative G12, macOS Apple Silicon and Windows x64 G13-G14 actual packaged-runtime audits, G15 Release Assurance, same-artifact no-rebuild publication and G16. The certified source fingerprint is `fa7b49ec9001d5ef95b829834f6268100e2eaf7c3da6bcc1f1a0b9bcba208d46`.

Durable v18.8.0 continuity evidence is bound in:
- `.depulse-certification/resume/build-checkpoint.json`;
- `.depulse-certification/resume/release-evidence-checkpoint.json`;
- `release/v18.8.0/stable-evidence-manifest.json`.

`ADAPT-REL-001` is CLOSED. This was continuity/source-of-truth repair only; certified v18.8.0 binaries were not rebuilt.

## v18.8.1 completed engineering packets

The following packets are implemented on `v18.8.1-development` and are now reflected in `adaptive-governance/CURRENT_ADAPTIVE_BUILD_PLAN.md`:

- `ADAPT-CI-001` Release State Coherence — implemented (`35531cd98465565c60cb2a26a3e066692d0f2168`).
- `ADAPT-CI-002` Early G11 Target Preflight — implemented with release-state coherence and invariant coverage.
- `ADAPT-CI-003` Cheap-First Fast Ordering — implemented (`89e78f0da57c4819c2ad818c73541fcc7713f269`).
- `ADAPT-CI-004` Safe Manual Dispatch Defaults — implemented (`02da4e1a560e16671188abdc69037116de4994c2`).
- CI hardening invariants are locked by `fd58ce8949013f5e6ffd70b69a3a991f8d4453f1`.
- `ADAPT-DATA-001` Universe Eligibility Contract — implemented by `3337386c5492a7af49e6e4dc49ef25dd23f94a44` plus follow-up `d63d0753d4d321362934e0ad7dabc52b7dca9b32`.
- `ADAPT-DATA-002` Evidence-Time Truth — implemented by `c01bd4d683be6236ec884aaea0cc14f7adf4d0b8`. Provider observation/evidence time drives Scanner/Opportunity Radar freshness; retrieval time is bookkeeping; old evidence fetched now remains old; unknown provider evidence remains unknown.
- `ADAPT-ARCH-001` Shared-Universe Robustness — CLOSED BY FRESH EVIDENCE. `symbol_universe.go` has a neutral shared U.S.-equity universe owner, cancellation/panic-safe deferred single-flight cleanup and retrieval/evidence-time separation; `v18_8_1_universe_hardening_test.go` carries focused eligibility and recovery evidence. Smart Provider Router v2 and BroadSnapshotBroker remain canonical owners for routing and broad snapshot reuse respectively.

No v18.8.1 GitHub Actions qualification has been triggered for this engineering line yet. Do not claim v18.8.1 Fast/Qualified PASS until those exact-head workflows actually run later in qualification.

## Remaining frozen v18.8.1 sequence

Continue on the **same** `v18.8.1-development` branch, without retry/certification/promotion branches or duplicate workflows, in this exact order:

1. `ADAPT-UI-001` Renderer Modularization II
2. `ADAPT-QA-001` Test/Gate Consolidation
3. `ADAPT-GOV-001` Historical Reconciliation Identity
4. `ADAPT-COST-001` Cost per Trustworthy Evidence
5. `ADAPT-RECON-001` Zero-Miss Reconciliation
6. `ADAPT-UX-RESEARCH-001` Research Information Architecture
7. `ADAPT-SYMBOL-001` Symbol/Desk Correctness
8. `ADAPT-READINESS-001` Prep/Readiness Semantics
9. `ADAPT-FRESHNESS-001` Freshness/Data Engine Correctness
10. `ADAPT-RESEARCH-002` Research Correctness Closure

A historical carry-forward that is already correct closes by fresh evidence. A reproducible gap keeps its original requirement identity and is fixed in its assigned coherent slice. Do not defer a known reproducible gap to a generic v18.10 catch-all.

## Provider continuity

Smart Provider Router v2 remains the sole provider-routing authority. BroadSnapshotBroker remains the canonical broad snapshot reuse/coalescing owner. Deterministic Day/Swing/Long truth remains protected. GLD/SLV/USO remain actionable tradable exceptions. U.S. Equities Processing and No Execution remain permanent boundaries.

v18.9.0 retains `ADAPT-TRADEINSIGHT-001` Full Capability Discovery, Utility Mapping & SHADOW Integration as the next provider-intelligence release after v18.8.1.

## Resume rule

Read `AGENTS.md` or `CLAUDE.md`, `governance/AI-ASSISTANT-PORTABILITY-CONTRACT.md`, this handoff, `release_identity.json`, both `.depulse-certification/resume/` checkpoints, `release/v18.8.0/stable-evidence-manifest.json`, the four CURRENT Adaptive overlays, live GitHub branch/PR/workflow state, and the conserved reconciliation authority before changing source. Never resume from model memory alone.

## Exactly one next action

Implement `ADAPT-UI-001` on `v18.8.1-development`: perform bounded capability-owner extraction for Renderer Modularization II, preserving existing visual design, deterministic Day/Swing/Long behavior, Research capability contracts, scroll/focus behavior, provider/snapshot ownership and U.S.-equities/no-execution boundaries. Add focused regression evidence and update this handoff/build-plan state after completion. Do not create a new branch or trigger qualification CI merely for packet development.
