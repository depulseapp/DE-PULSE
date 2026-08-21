# DE.PULSE — Current Adaptive Build Process

**Certified Stable:** `v18.9.0-stable`  
**Current activity:** post-Stable corrective diagnosis for issue #64 / `ADAPT-RUNTIME-CRASH-001`.  
**Active development branch:** none.

Process remains:

**GitHub source of truth → reconcile actual Stable → perform G0 exact-baseline → freeze verified G1 scope → one version-development branch → one Draft PR → exact-head Fast → same PR Ready → Qualified → exact-head merge → one G11–G16 Release when release-capable → post-Stable continuity reconciliation before the next product line.**

## v18.9.0 process closure

G0–G16 are PASS. Final Fast #481 / `32525637987`, Qualified #153 / `32525738828`, and Release #32 / `32526121817` are the release authority. macOS Apple Silicon and Windows x64 packaged-runtime audits, G15 evidence binding, no-rebuild publication and G16 all passed.

The release nevertheless has a real-world post-certification escape: the user reported a native macOS v18.9.0 `EXC_CRASH (SIGABRT)` / `abort() called`. This does not rewrite the recorded release workflow result; it creates the next corrective learning loop. Version-scoped issue #63 is superseded by issue #64.

## Corrective process rule

Before code changes for #64:
- obtain the complete crash report/backtrace or deterministic reproduction;
- classify the owning lifecycle surface from evidence rather than guesswork;
- preserve user state/API-key continuity as a safety constraint;
- add regression/lifecycle proof tied to the root cause;
- keep the correction isolated from optional provider/Market-Mode expansion unless technically required;
- run the normal exact-head Fast → Qualified → one G11–G16 release path.

Repair order remains REUSE → CONSOLIDATE → REFACTOR → DELETE/REPLACE → ADD. Failure handling remains classify-first: `PRODUCT_FAIL`, `GATE_TEST_FAIL`, `CI_HARNESS_FAIL`, `INFRA_FAIL`, `EXPECTED_NOOP`, `SUPERSEDED`. No new top-level gates; G0–G16 remains the only release model.

Permanent owners remain unchanged: Smart Provider Router v2 sole routing authority, canonical freshness/recovery sole freshness owner, existing multi-feed allocator sole subscription owner, deterministic Day/Swing/Long truth, U.S. Equities Processing, GLD/SLV/USO actionable exceptions and No Execution.

## Exactly one next action

Run issue #64 G0 crash diagnosis from concrete evidence/reproduction and freeze the bounded corrective G1 before implementation.
