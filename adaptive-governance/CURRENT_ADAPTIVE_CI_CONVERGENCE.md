# DE.PULSE — Current Adaptive CI / Versioning Convergence

**Program:** #70 / `ADAPT-CI-CONVERGENCE-001`  
**Status:** APPROVED / AUDITED / IMPLEMENTATION NOT STARTED  
**Scope:** Process, CI, versioning and release-evidence hardening only; no product behavior change.  
**Ordering:** finish #64 / `v18.9.1` truthfully first, then execute this convergence packet before accumulating further version-specific CI/release machinery.

## 1. Permanent architecture

Keep exactly three active GitHub Actions workflow families:

1. `ci-fast.yml`
2. `ci-qualified.yml`
3. `release.yml`

Do not create retry/certification/promotion/recovery workflow families. Version, candidate, platform and evidence requirements are parameters of permanent CI.

No CI efficiency optimization may weaken G0-G16, exact-head evidence, macOS Apple Silicon, Windows x64, Chrome/WebKit, security/data/provider truth, Stable immutability, U.S. Equities Processing or No Execution.

## 2. Approved convergence targets

1. Evolve Impact Planner v2 into deterministic Planner v3 job selection with explicit dependency invalidation and fail-closed `full` fallback for unknown/mixed risk.
2. Manual/reusable qualification must bind the full intended delta through explicit base SHA/ref or an authoritative merge-base; do not silently classify only `HEAD^`.
3. Consume native/release-harness impact as real targeted rehearsal inside existing Qualified; no fourth workflow.
4. Separate WebKit browser ownership from macOS native-lifecycle rehearsal ownership while reusing one canonical `tools/release/native_macos.sh` implementation.
5. Replace future per-version G12 executable scripts with one canonical version-neutral certification executor driven by release identity plus declarative release/capability manifests.
6. Reduce release-identity fan-out and version-only source churn; prefer one canonical identity resource/build injection and content-hash cache busting.
7. Migrate active version-named tests incrementally to capability-oriented ownership paths with equivalence proof before old active-path retirement.
8. Move version/predecessor/tag/publication-feasibility checks into G11 so impossible releases fail before expensive G12/G13/G14 work.
9. Stable publication is immutable: existing same-digest asset = idempotent no-op; differing digest = release-integrity failure; never overwrite differing Stable bytes.
10. Serialize Stable G11-G16 publication repository-wide with `cancel-in-progress: false`.
11. Verify repository rulesets/branch protection. If equivalent protection is absent, protect `main` with PR-only changes, trustworthy required checks, no force-push/deletion and controlled bypass.
12. Separate `productVersion`, `workSliceId`, `candidateSha`, `sourceFingerprint`, `buildId` and `evidenceSchemaVersion`.
13. Make one explicit prospective versioning decision: adopt SemVer or formally document a custom DE.PULSE release-train scheme. Do not renumber or rewrite shipped history.
14. If SemVer is adopted, use normal Stable tags `vX.Y.Z`, prerelease tags such as `vX.Y.Z-rc.N`, and keep channel separately. Historical `-stable` tags remain immutable.
15. Enforce a collision-free monotonic native build-number contract independent of display product version; apply the canonical product/build identity model across supported installers/artifacts.
16. Add modern runnable-artifact provenance/attestations and SBOM evidence at the software-supply-chain milestone with narrowly scoped release permissions.
17. Define one canonical toolchain manifest and record resolved Go/Node/Python/Playwright/runner-image identity in release evidence.
18. Converge `handoff/CURRENT.md` and the four `CURRENT_ADAPTIVE_*` overlays on one actual current-state truth; prefer a machine-readable current release state over repeated hand-authored fields.
19. G16 must report runner minutes/reruns avoided, evidence reuse and confirmation that no quality requirement was removed.

## 3. Versioning / build identity target

Do not make every CI run or independently certifiable work packet a public product version.

Preferred model:

`work packet -> commits -> Fast -> deliberate Qualified checkpoint -> exact candidate SHA/fingerprint -> release grouping -> one public product version -> one canonical G11-G16/native publication`

Identity dimensions are independent:

- `productVersion` — user-visible compatibility/release version;
- `workSliceId` — independently governed engineering responsibility;
- `candidateSha` — exact Git source;
- `sourceFingerprint` — exact relevant source identity;
- `buildId` — immutable build/evidence identity;
- `evidenceSchemaVersion` — evidence contract version;
- native bundle/file build number — monotonic platform packaging identity.

Urgent user-impacting reliability/security/bug corrections may ship as PATCH releases. Backward-compatible feature groups normally become MINOR releases if SemVer is adopted. Existing v18.9.x reservations remain planning reservations until the prospective versioning decision is implemented; shipped releases/tags/evidence are immutable.

## 4. Adaptive CI target

Planner v3 should select the smallest trustworthy evidence graph, for example:

- backend/core;
- renderer;
- Chrome;
- WebKit;
- persistence/DB/restart/migration;
- security/rights;
- provider/router/data;
- portability/CI harness;
- macOS native rehearsal;
- Windows native rehearsal;
- full fail-closed fallback.

The planner must emit both selections and reasons. Cost is never a reason to omit required evidence. Shared evidence may be reused only when candidate identity and dependency validity remain exact.

## 5. Release target

G11 becomes the cheap fail-fast release checkpoint: exact source/fingerprint, required-platform matrix, version/build/predecessor compatibility, current Stable expectation, target tag/release state and publication feasibility.

G12 uses one canonical executor. G13/G14 retain actual supported artifact/deployment proof. G15 consumes exact same-run certified artifacts. Publication does not rebuild and cannot mutate differing Stable bytes. G16 records evidence, release provenance, CI efficiency and cleanup/debt.

## 6. Cleanup discipline

Historical release evidence remains immutable. Cleanup is incremental:

`inventory -> classify -> prove equivalent owner/coverage -> switch active consumer -> verify -> retire obsolete active path`

Never delete version-named tests/gates merely because their names are old. Never create a parallel replacement owner before the existing canonical owner is consolidated.

## 7. Implementation ordering

1. #64 / `v18.9.1` remains the active product corrective and must obtain exact-head Fast + full Qualified + required native proof.
2. After truthful v18.9.1 closure, execute #70 / `ADAPT-CI-CONVERGENCE-001` as a dedicated process-hardening packet before further version-specific CI/release machinery accumulates.
3. Reconcile future v18.9.x product reservations against the approved product-version/work-slice separation before starting the next product implementation.
4. No #70 work may be used to weaken v18.9.1 evidence or retroactively redefine shipped Stable history.

## 8. Acceptance

#70 is complete only when all issue acceptance criteria pass in executable evidence and the four Adaptive overlays, handoff, workflow policy and actual repository settings agree.

## Exactly one next action

Restore/confirm GitHub Actions hosted-runner execution and finish exact-head #64 / `v18.9.1` qualification. #70 is approved and must execute immediately after truthful v18.9.1 closure, before the next product implementation creates additional version-specific CI/release machinery.
