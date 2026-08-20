# DE.PULSE — Current Adaptive Build Process

**Certified Stable:** `v18.8.0-stable`  
**Current activity:** `v18.8.1-development` governance hardening and G0 entry preparation.

Process remains:

**GitHub source of truth → reconcile actual Stable → freeze only verified G1 scope → one version-development branch → one Draft PR → exact-head Fast → same PR Ready → Qualified → exact-head merge → one G11–G16 Release when release-capable.**

## v18.8.1 process hardening

- Run `ADAPT-REL-001` first so v18.8.0 Stable tag/Release, checkpoints, Stable evidence manifest, handoff and CURRENT overlays agree before normal product work.
- Add `ADAPT-CI-001` Release State Coherence so one cheap preflight reports every release-state mismatch instead of discovering VERSION/checkpoint/manifest/handoff drift sequentially.
- Add `ADAPT-CI-002` at G11 so a conflicting target Stable tag/version/build/predecessor fails before G12/native certification. Keep the publication-time tag guard too.
- Reorder Fast per `ADAPT-CI-003`: Python/impact/coherence/governance/identity/provenance before expensive Go/Node/browser setup.
- Manual CI defaults follow `ADAPT-CI-004`: smallest safe/adaptive lane by default; `full` requires explicit intent.
- For market evidence, `ADAPT-DATA-002` separates provider evidence time from retrieval time. Missing provider time never becomes `time.Now()` freshness.
- For Scanner/Radar, `ADAPT-DATA-001` makes universe eligibility explicit; acquisition filters such as `has_options` cannot silently redefine the advertised broad U.S.-equity universe.
- `ADAPT-ARCH-001` keeps one neutral shared-universe owner with cancellation/panic-safe refresh cleanup and no duplicate router/broker/freshness engine.
- Renderer/test cleanup uses strangler migration: stable capability owners first, equivalence proof, then retire superseded release-number owners. No mass delete/rename.

Failure handling remains classify-first: `PRODUCT_FAIL`, `GATE_TEST_FAIL`, `CI_HARNESS_FAIL`, `INFRA_FAIL`, `EXPECTED_NOOP`, `SUPERSEDED`. No unchanged reruns or CI event manufacturing.

Historical v18.5.1 reconciliation remains provenance, not current-release identity. Actual GitHub objects and CURRENT overlays identify the current Stable/release line.

No new top-level gates. G0–G16 remains the only release model.