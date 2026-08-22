# DE.PULSE — Current Adaptive CI / Versioning Convergence

**Program:** #70 / `ADAPT-CI-CONVERGENCE-001`  
**Status:** APPROVED / AUDITED / EXECUTABLE IMPLEMENTATION REQUIRED / NOT STARTED  
**Scope:** Process, CI, versioning, repository hygiene and release-evidence hardening only; no product behavior change.  
**Ordering:** finish #64 / `v18.9.1` truthfully first, then execute this convergence packet before accumulating further version-specific CI/release machinery.

**Closure rule:** documentation, audit notes, plans, or comments alone cannot satisfy #70. Completion requires actual repository/code/workflow/tooling changes plus executable evidence.

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
20. Execute repository-root hygiene as real implementation: add safe ignore policy, make source-health package-aware, eliminate contradictory legacy certification ownership, migrate historical/version-scoped root machinery, migrate active version-named tests safely, move stable assets/policies to canonical owners, and enforce a root allowlist.

## 3. Mandatory executable repository-hygiene packet

The repository-hygiene work discovered during the root audit is part of #70 implementation, not an optional follow-up.

Current audit facts:
- root mixes production `package main` Go source, package-local tests, active CI/release entrypoints, historical version-scoped tests/gates/contracts, policies/registries, renderer/browser tests, release metadata and retained assets;
- the existing executable inventory reported `160` version-stacked root executable tests/gates: `80 ACTIVE_REQUIRED`, `80 UNREFERENCED_USEFUL`;
- `source_health_architecture_gate.py` currently uses flat-root Go discovery and must become recursive/package-aware before production package moves;
- the root legacy cluster `certification_plan.json`, `certification_runner.py`, `ci_pipeline.py`, `ci_pipeline_plan.json` represents older orchestration and must not remain a competing current CI/release owner;
- package-local Go tests must not be cosmetically moved into a different directory/package without preserving their access/coverage semantics;
- retained branding/source assets remain intentional until migrated with all consumers atomically updated.

Implementation waves:

1. **Guardrails before moves** — add a safe repository `.gitignore`; make source-health/package discovery recursive; extend the existing legacy inventory into a complete root-layout inventory/allowlist with consumer/reason classification.
2. **Canonical non-Go tooling ownership** — move reusable CI/dev/release scripts from root into `tools/ci`, `tools/release`, or `tools/dev`; update all active consumers atomically.
3. **Eliminate contradictory legacy certification orchestration** — consolidate useful resume/checkpoint behavior into the canonical version-neutral release owner, then retire stale root orchestration/metadata once no active consumer remains.
4. **Historical/version-scoped non-Go cleanup** — reclassify after stale-control removal; move immutable historical contracts/evidence/config into governed release/history ownership; delete only after assertion/evidence equivalence is proven.
5. **Active version-named Go test migration** — rename to capability-oriented names while preserving package-local access and `go test ./...`, race, randomized and focused coverage; prohibit new version-prefixed active tests after migration.
6. **Production package decomposition** — only after source-health is package-aware, extract cohesive implementation + tests together into `internal/<capability>` packages with canonical-owner/equivalence proof.
7. **Policies, registries and retained assets** — converge them into stable governance/config/asset owners with atomic path-consumer updates.
8. **Permanent root allowlist** — CI fails on newly introduced arbitrary/version-prefixed root gates/tests/contracts/release scripts unless explicitly transitional with a migration owner and expiry.

## 4. Versioning / build identity target

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

## 5. Adaptive CI target

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

## 6. Release target

G11 becomes the cheap fail-fast release checkpoint: exact source/fingerprint, required-platform matrix, version/build/predecessor compatibility, current Stable expectation, target tag/release state and publication feasibility.

G12 uses one canonical executor. G13/G14 retain actual supported artifact/deployment proof. G15 consumes exact same-run certified artifacts. Publication does not rebuild and cannot mutate differing Stable bytes. G16 records evidence, release provenance, CI efficiency and cleanup/debt.

## 7. Cleanup discipline

Historical release evidence remains immutable. Cleanup is incremental:

`inventory -> classify -> prove equivalent owner/coverage -> switch active consumer -> verify -> retire obsolete active path`

Never delete version-named tests/gates merely because their names are old. Never create a parallel replacement owner before the existing canonical owner is consolidated.

## 8. Implementation ordering

1. #64 / `v18.9.1` remains the active product corrective and must obtain exact-head Fast + full Qualified + required native proof.
2. After truthful v18.9.1 closure, execute #70 / `ADAPT-CI-CONVERGENCE-001` as a dedicated process-hardening packet before further version-specific CI/release machinery accumulates.
3. The repository-hygiene packet above is part of #70 implementation and must be executed, not merely documented.
4. Reconcile future v18.9.x product reservations against the approved product-version/work-slice separation before starting the next product implementation.
5. No #70 work may be used to weaken v18.9.1 evidence or retroactively redefine shipped Stable history.

## 9. Acceptance

#70 is complete only when all issue acceptance criteria pass in executable evidence and the four Adaptive overlays, handoff, workflow policy and actual repository settings agree.

In addition, #70 cannot close unless:
- actual root/code/workflow/tooling changes have landed;
- source-health still covers every production Go package after any relocation;
- only one current CI/release orchestration owner remains;
- package-local Go regression coverage is preserved;
- historical Stable tags/releases/evidence remain unchanged;
- a CI-enforced root allowlist prevents recurrence;
- G16 records before/after root inventory and proves no quality requirement was removed.

Documentation-only completion is prohibited.

## Exactly one next action

Restore/confirm GitHub Actions hosted-runner execution and finish exact-head #64 / `v18.9.1` qualification. #70 is approved and must execute immediately after truthful v18.9.1 closure, including the mandatory executable repository-hygiene packet, before the next product implementation creates additional version-specific CI/release machinery.
