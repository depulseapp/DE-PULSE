# DE.PULSE — Current Adaptive Build Plan

**Immutable Stable:** `v18.8.0-stable`  
**Stable candidate:** `3a32d57dd4c74c6f812cc942a9d8049a7b517718`  
**Next:** `v18.8.1 — Renderer Modularization II + 10/10 Audit Hardening`.

## Entry condition

First close `ADAPT-REL-001`: reconcile the actual v18.8.0 Stable tag/Release into build checkpoint, release-evidence checkpoint, Stable evidence manifest and `handoff/CURRENT.md`. This is continuity/source-of-truth closure only; unchanged certified v18.8.0 binaries do not rebuild.

## v18.8.1 mandatory work packets

1. `ADAPT-CI-001` Release State Coherence validator covering release identity, VERSION/app/renderer/cache-bust identity, checkpoints, Stable manifest, handoff, predecessor/baseline and target tag state in one result.
2. `ADAPT-CI-002` G11 target-tag/version/build/predecessor/release-scaffold preflight before G12/native work, retaining the final publication collision guard.
3. `ADAPT-CI-003` cheap-first Fast ordering; expensive Go/Node/browser setup starts only after cheap governance/coherence/identity/provenance checks are green.
4. `ADAPT-CI-004` safe manual dispatch defaults; full qualification must be explicit.
5. `ADAPT-DATA-001` explicit discovery-universe eligibility and intentional treatment of provider filters such as `has_options`.
6. `ADAPT-DATA-002` provider evidence time separated from retrieval time; missing evidence time is UNKNOWN/degraded/ABSTAIN rather than wall-clock `now`.
7. `ADAPT-ARCH-001` neutral shared-universe diagnostics, panic/cancellation-safe single-flight cleanup and explicit eligibility/concurrency tests.
8. `ADAPT-UI-001` Renderer Modularization II via bounded capability-owner extraction, not a broad redesign.
9. `ADAPT-QA-001` inventory and gradual capability-based consolidation of version-stacked tests/gates; preserve unique historical evidence.
10. `ADAPT-GOV-001` separate historical reconciliation baseline identity from current release identity.
11. `ADAPT-COST-001` measure avoided runs/setup/minutes and cost per trustworthy evidence without weakening required gates.

## Required execution order

`G0 continuity closure → G2/G3 owner/semantic design → CI hardening → data-truth/shared-universe hardening → bounded renderer/test modularization → G10 reconciliation → exact-head Fast/Qualified → one G11–G16 Release when release-capable`.

Use one `v18.8.1-development` branch, coherent batches, one Draft PR, automatic Fast, same PR Ready, Qualified, exact-head merge and one Release. Never create retry/certification/promotion branches or duplicate workflows.

Permanent boundaries remain: U.S. Equities only, No Execution, G0–G16 only, Smart Provider Router v2 sole provider-routing owner, BroadSnapshotBroker sole broad snapshot reuse/coalescing owner, deterministic Day/Swing/Long truth protected.