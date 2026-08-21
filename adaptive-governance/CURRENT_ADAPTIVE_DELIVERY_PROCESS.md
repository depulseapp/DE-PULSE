# DE.PULSE — Current Adaptive Delivery Process

**Certified Stable:** `v18.9.0-stable`  
**Stable candidate:** `9ea81cddae4875ae15d3719ca028519a36c597b6`  
**Stable fingerprint:** `a8719090c341c874dbd1279cc31ad98e84075d5701c46a800bf951340780ecb9`  
**Final Fast:** #481 / `32525637987`  
**Final Qualified:** #153 / `32525738828`  
**Final Release:** #32 / `32526121817`  
**Immediate corrective requirement:** issue #64 / `ADAPT-RUNTIME-CRASH-001`.

## v18.9.0 delivery — COMPLETE

G11–G16 completed in Release #32. G12 full certification passed. macOS Apple Silicon and Windows x64 G13/G14 built and audited the actual packaged runtimes. G15 bound both native evidence graphs to the same candidate/fingerprint/build identity. Publication verified the evidence graph and exact binaries, then published `v18.9.0-stable` from those same-run certified artifacts **without rebuilding**. G16 durable handoff evidence passed.

The immutable Stable tag identifies candidate `9ea81cddae4875ae15d3719ca028519a36c597b6`. Durable release evidence is `release/v18.9.0/stable-evidence-manifest.json` plus both resume checkpoints. Issue #61 is closed completed and PR #62 is merged.

## Post-Stable escape handling

The user subsequently reported a native macOS v18.9.0 crash (`EXC_CRASH (SIGABRT)` / `abort() called`). Issue #63 is closed as superseded rather than fixed; issue #64 is the corrective owner. Stable evidence remains immutable and truthful: the release jobs passed, while real-world use exposed an additional failure not caught by the existing audit.

The corrective delivery must use concrete root-cause evidence/reproduction, regression coverage, actual macOS Apple Silicon packaged-runtime proof and the same one-branch/one-PR exact-head Fast → Qualified → canonical G11–G16 path. Do not silently reset user state, create retry/certification branches, duplicate release workflows, or bundle unrelated feature expansion ahead of this blocker.

Normal delivery remains one development branch → one Draft PR → Fast → same PR Ready → Qualified → exact-head merge → one canonical G11–G16 run → exact certified artifacts published without rebuild → repository continuity reconciliation.

No duplicate release workflows or manual duplicate CI runs. G0–G16 only.

## Exactly one next action

Obtain/reproduce the issue #64 macOS crash, freeze the bounded corrective scope, and only then begin the next release implementation.
