# DE.PULSE — Adaptive Delivery Process

**Status:** ACTIVE / GOVERNED  
**Authority:** `governance/ADAPTIVE-OPERATING-CONTRACT.md` and the certified immutable RC. This file defines efficient trustworthy delivery inside G13–G16.

Delivery principle:

**certified source → independent package lanes → actual artifact audit → assurance → promotion → user delivery → G16 learning.**

G0–G16 remains permanent. No delivery workflow may create G17+.

---

## 1. Delivery Resume Reconciliation

When delivery resumes after interruption:

**read immutable RC identity → inspect actual workflow/artifact/release state → verify package hashes/provenance → preserve unchanged PASS evidence → resume only incomplete/affected lanes.**

Never assume a failed/stopped conversation or wrapper means nothing was produced. Inspect durable GitHub objects first.

Never duplicate tags/releases/assets without checking existing state.

---

## 2. Exact-artifact evidence inheritance

Delivery evidence is reusable only when its relevant:
- immutable RC/source identity;
- package/artifact identity/hash;
- packaging/toolchain contract;
- runtime audit definition;
- platform assumptions

remain equivalent.

If one platform is PASS and the other fails/interruption occurs, preserve the PASS platform when the exact RC/package identity is unchanged.

Metadata/checkpoint/documentation commits that do not alter immutable RC/package/test contract do not justify rebuilding native artifacts or rerunning full G12.

---

## 3. G13 — Native Packaging / Provenance

Build required native artifacts from the certified immutable RC.

Required current platform truth:
- macOS Apple Silicon — independent package/provenance lane;
- Windows x64 — independent package/provenance lane.

Share safe immutable setup/evidence, but do not couple platform outcomes unnecessarily.

Each artifact records source/RC identity, package identity/hash and provenance needed by G14/G15.

---

## 4. G14 — Actual Artifact Runtime Audit

Audit the **actual packaged artifact**, not merely source tests.

Verify applicable:
- launch/run behavior;
- runtime/release identity;
- persistence/migration/restart behavior;
- configuration;
- critical user flows;
- role/capability behavior;
- data/freshness/degradation truth;
- packaging integrity;
- platform-specific behavior;
- no execution capability.

Platform runtime evidence is independent. One platform cannot inherit the other's platform-specific PASS.

---

## 5. G15 — Release Assurance / Promotion

G15 consumes the complete authoritative evidence graph:
- immutable RC/source provenance;
- G12 certification;
- required G13/G14 platform results;
- security/data-rights/reliability/performance closure where applicable;
- rollback/reproducibility truth;
- final release identity.

Promotion is allowed only when required evidence is complete and trustworthy.

Generate final artifact/hash manifest only after final contents are immutable.

---

## 6. User Delivery completion invariant

A Stable native release is `READY` when required permanent certified assets exist, but `DELIVERED` only when the user-facing handoff surfaces the required runnable assets and provenance/hash status.

Concise release status:

**Code | G0–G12 | macOS G13/G14 | Windows G13/G14 | G15 | G16 | User Delivery**

Routine recovery/delivery should require no manual user work except approved external blockers.

---

## 7. Delivery AIPLC

A meaningful packaging/runtime/delivery event runs a concise Delta AIPLC for the affected delivery responsibility.

Ask:
- what actually failed or cost unnecessary work?
- was it product, packaging, runtime, CI harness, infra, publishing, provenance or user-delivery logic?
- what is the root cause?
- what can be reused?
- what prevention makes recurrence less likely?
- can the fix remain platform/lane-scoped?

Every material delivery challenge produces:
1. immediate fix/truthful disposition;
2. reusable prevention.

Do not rerun unrelated platforms/lanes merely to generate another report.

---

## 8. Role-aware delivery invariant

Before G15 for affected scope, actual certified packages must prove:
- authorized OWNER/SUPER_OWNER composition;
- explicitly delegated ADMIN capability boundaries;
- limited ADMIN restrictions;
- USER/DEMO no implementation machinery/privileged payloads;
- frontend/backend authorization agreement;
- direct unauthorized route/API denial;
- natural layout reflow after suppressed content;
- role-invariant protected market logic.

Use fresh + equivalence-bound inherited coverage from G9/G10/G12. Do not brute-force unchanged role×surface×viewport combinations again without a dependency reason.

---

## 9. Functionality Utility delivery invariant

A release is not delivery-complete if it ships working but unnecessary/unintegrated machinery.

Certified evidence must show applicable changed functionality has:
- purpose/consumer;
- canonical owner;
- reused/consolidated acquisition/computation where practical;
- relevant correlation;
- appropriate freshness/materiality/retention/governance;
- one deep-evidence home with concise contextual reuse elsewhere;
- justified user-facing prominence/tab placement;
- obsolete/superseded paths retired or explicitly carried forward.

G10 full coverage reconciliation remains the pre-freeze completeness boundary.

---

## 10. Shared Symbol Intelligence delivery invariant

For releases affecting multi-user symbol demand, provider acquisition, Scanner/Radar, preparation/event processing, canonical state, AI synthesis or hosted scale, delivery evidence must prove:
- equivalent canonical work is shared across compatible consumers;
- per-user workspaces do not create duplicate market pipelines;
- in-flight work coalesces where practical;
- dynamic attention/backpressure protects high-value work;
- private/tenant/rights/entitlement partitions remain isolated;
- overlapping user cost scales primarily with unique canonical demand;
- no unauthorized cross-user leakage.

Use the material subset of the shared-efficiency scorecard appropriate to the release. Full realistic overlapping-demand closure is mandatory at v18.5.

---

## 11. v18.2 delivery closure

v18.2 cannot promote based on governance documents alone.

Before promotion, prove applicable frozen-scope and current process-hardening obligations including:
- capability-based ADMIN authorization;
- Administration composition;
- role/session/presence lifecycle behavior;
- authoritative current Build State Ledger/checkpoint;
- trusted Build Coordinator/evidence graph;
- canonical naming/evidence state;
- first v18.2 AIPLC and G10 coverage reconciliation.

Do not pull v18.3 PostgreSQL/web or unrelated source-changing consolidation work into frozen v18.2.

---

## 12. G16 — Deep Adaptive Retrospective / Handoff

G16 must:
- verify Stable `main` / tag / release identity agreement;
- record source/native artifact provenance;
- archive checkpoints/evidence;
- clean obsolete release branches/workflows/artifacts when safe;
- seed the next approved release from exact Stable;
- aggregate AIPLC findings;
- record calls/reruns/provider/model work avoided where measurable;
- record role/functionality/shared-processing closure;
- convert recurring failures into reusable prevention;
- remove redundant delivery lanes/checkpoints that added ceremony without assurance;
- verify required user delivery was surfaced.

No issue silently disappears.

---

## 12A. Adaptive CI delivery control

Delivery uses the same evidence planner, but release guardrails are deterministic:

- G10 and G12 must be green for the exact immutable candidate before native packaging starts.
- macOS Apple Silicon and Windows x64 package/audit lanes run in parallel only after shared cheaper prerequisites pass.
- Native evidence is never waived to meet a dollar target. The $5 budget is a soft alert and planning input, not a release gate.
- A failed platform reruns only that platform and any evidence genuinely invalidated by the fix; the other platform PASS is preserved when fingerprints remain equivalent.
- G15 binds hashes, provenance, source commit, requirement ledger and native audit outputs.
- Publication uploads the already-certified artifacts and performs no rebuild.
- Development diagnostics use short retention; RC evidence uses bounded retention; certified artifacts, checksums and provenance live in the GitHub Release.
- Workflow permissions are read-only except for the isolated publication job.
- Every release run records selected/skipped lanes, reason, runtime/cost estimate, cache result, failure class, evidence reuse and invalidation for G16 analysis.
- G16 must distinguish quality-producing spend from avoidable orchestration, record the prevention action and hand the next release an updated—but not silently self-modified—planner proposal.

Delivery cannot pass while `CI-ADAPTIVE-18.5.1-001` is unresolved for the applicable slice or while required native/exact-artifact evidence is missing.

---

## 12B. Independent-audit delivery contract

Applies to: `AUDIT-18-UI-001`, `AUDIT-18-AI-001`, `AUDIT-18-AI-RIGHTS-001`, `AUDIT-18-PROVIDER-001`, `AUDIT-18-ARCH-001`, `AUDIT-18-CI-001`, `AUDIT-18-SECURITY-001`, `AUDIT-18-PROVENANCE-001`, `AUDIT-18-TRADER-001` and `AUDIT-18-QA-001`.

Every v18.x delivery manifest must include the ten `AUDIT-18-*` records, their subfindings, release placement, source/test fingerprints and final/open status. The authoritative inventory is **295 rows** after audit integration.

Delivery requirements:

1. **UI/live runtime:** actual packaged macOS and Windows apps must dwell on Dashboard, Research and every desk while live/replayed SSE updates arrive; capture hover, focus, selected state, semantic scroll anchor, layout-shift/repaint observations and symbol mutation behavior.
2. **AI:** ship a prompt/safety/schema/model/provider manifest, bounded-context proof, cache isolation/expiry proof, golden/injection eval summary, token/cost/latency telemetry and safe abstention behavior.
3. **AI rights:** ship the provider×dataset AI-egress decision summary with secrets and contractual evidence redacted; blocked/unknown rights must remain blocked in the package.
4. **Provider learning:** SHADOW comparisons must state sample size, dataset/session/regime, cost/calls, disagreement/truth limits, useful/missed/false evidence and non-production status. No paid-provider promotion occurs as an incidental release action.
5. **Trader learning:** probability outputs, if any, require cutoff-safe calibration evidence and explicit distinction from setup scores; evidence-thesis logs cannot introduce positions, P&L, orders or execution.
6. **Architecture:** package provenance identifies canonical owners and retired duplicates; refactors must preserve formulas, persisted profile compatibility, migration/rollback and cross-platform startup.
7. **Security/supply chain:** include vulnerability/dependency results, SBOM, credential-storage migration evidence, HTTP/server fault tests, secret-scan result and any accepted risk with owner/expiry.
8. **CI:** release evidence states which durable workflow selected each lane, why it ran/skipped, its fingerprint, runtime/cost/cache result and failure/retry class. A branch-trigger gap blocks delivery.
9. **Provenance:** bind source commit, verified tag/attestation, artifact hashes, SBOM and G13/G14 results. Publication uploads certified artifacts without rebuilding.
10. **Handoff:** G16 records every audit ID as final, explicitly future-placed or still blocking; it includes the maturity stage, next action, owner and invalidation conditions so a new chat never restarts from memory.

No v18 final closure delivery may retain a reopened/not-implemented/revalidation audit row, an unresolved audit subfinding, unsigned/unattested future Stable provenance, or missing exact-package proof for affected macOS/Windows surfaces.

---

## 13. 10/10 Delivery acceptance

Delivery Process is 10/10 only when:
1. exact immutable source/artifact identity controls evidence;
2. independent platform PASS is preserved safely;
3. failures rerun the smallest affected lane;
4. metadata changes do not trigger unnecessary native/full certification;
5. actual packaged runtime is audited;
6. role/security/data truth survives packaging;
7. user delivery is explicit;
8. AIPLC produces prevention, not just retries;
9. G16 leaves reproducible provenance and a clean next-release seed;
10. no promotion status is inferred or inflated without durable evidence.

## 14. Portable G16 delivery handoff

G16 delivery is incomplete until GitHub alone can onboard a newly authorized ChatGPT/Codex, Claude or human maintainer. Update `handoff/CURRENT.md` with delivered artifacts/hashes, actual residuals and exactly one next action; update the Build State Ledger after the candidate commit; and verify both root assistant adapters still point to the same vendor-neutral portability contract.

Do not make ChatGPT Library, Claude Projects, a local Mac, email, chat history or a temporary AI workspace the only location of required source, evidence or continuation instructions. Those may be mirrors, never delivery authority.
