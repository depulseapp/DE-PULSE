# DE.PULSE — Current Authoritative Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS. GitHub is authoritative; model memory is advisory only.**

## Current identity

**Certified Stable:** `v18.7.0-stable`  
**Certified candidate / tag target:** `75e494fb92441439c73c8ace41a40118e4518c1c`  
**Certified source fingerprint:** `350e1f87f2046410ae52623de9dacba8fca2a16fda9b116232107fa8f8cac963`  
**Certified Build ID:** `v18.7.0-stable-20260819`  
**Engineering branch:** `v18.8.0-development`  
**Candidate package identity:** `18.8.0` / `v18.8.0-stable-20260819`  
**Engineering baseline:** `6ab2094cb7aefdb3e1e21862cbd17b64ed850c28` (PR #54 merged product source)  
**Release recovery:** Release #27 (`32335527140`) proved G11, G12, macOS/Windows G13-G14, G15 and exact-artifact verification, then correctly blocked publication because the merged candidate still identified itself as v18.7.0.  
**Stable checkpoints:** checkpoints intentionally remain anchored to immutable v18.7.0 Stable while v18.8.0 is an in-flight candidate.

## v18.8.0 frozen scope — Shared Intelligence Consolidation

G0-G3 source-owner audit proved one meaningful consolidation gap and no justification for a broader rewrite: Discovery Scanner and Opportunity Radar had separate lifecycle/cache ownership around the Alpaca U.S.-equity universe. v18.8.0 consolidates that into `canonicalUSSymbolUniverse` while preserving the existing 12-hour TTL, coalescing concurrent refreshes, preserving provider-backed timestamps after failed refresh, using bounded retry suppression plus identifiable deterministic seed fallback, and preserving Scanner ranking plus Radar sampling/promotion semantics.

Smart Provider Router v2 remains the sole provider-routing owner. BroadSnapshotBroker remains the canonical snapshot reuse/coalescing owner. Session Intelligence, Event Intelligence, Research evidence hydration, deterministic Day/Swing/Long market truth, GLD/SLV/USO actionable tradable exceptions, U.S. Equities Processing boundary and No Execution boundary remain unchanged.

## Release recovery rule

Do not rerun Release #27 unchanged and do not modify, move or overwrite `v18.7.0-stable`. The bounded recovery is one recreated `v18.8.0-development` branch from merged candidate `6ab2094c`, one release-closure commit containing the correct v18.8 identity/scaffolding, one Draft recovery PR to `main`, exact-head Fast, the same PR Ready, full exact-head Qualified, exact-head merge, then one automatic Release G11-G16 run. No retry/certification/promotion branches and no duplicate CI/release workflows.

## v18.7.0 immutable Stable

v18.7.0 Runtime Reliability & Data Truth remains the last certified Stable until v18.8.0 completes G11-G16. Its immutable tag, candidate, fingerprint, artifacts and prior evidence are not rewritten to make the v18.8 candidate look Stable.

## Resume rule

Read `AGENTS.md` or `CLAUDE.md`, `governance/AI-ASSISTANT-PORTABILITY-CONTRACT.md`, this handoff, `release_identity.json`, both `.depulse-certification/resume/` Stable checkpoints, `release/v18.7.0/stable-evidence-manifest.json`, `v18_8_0_scope.json`, `v18_8_0_g0_g3_contract.json`, `release/v18.8.0/release_contract.json`, live GitHub state, the four current Adaptive overlays, and the conserved v17/v18 reconciliation ledger. Never resume from model memory alone.

## Exactly one next action

Complete the single bounded v18.8.0 release-recovery PR lifecycle: exact-head Fast → same PR Ready → full Qualified → exact-head merge → one automatic G11-G16 Release. If any gate fails, diagnose that exact failure and do not rerun unchanged evidence.
