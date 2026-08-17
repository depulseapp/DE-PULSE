# DE.PULSE GitHub Actions Efficiency Contract

Status: PERMANENT BUILD / DELIVERY GOVERNANCE

## Objective
Use GitHub Actions only when hosted execution materially adds evidence that cannot be obtained more efficiently elsewhere. GitHub is the durable source/release archive, not the default compute engine for every development change.

## Default execution model
1. Ordinary source edits, repository cleanup, documentation, inventory, dependency analysis, governance updates and release preparation MUST NOT trigger hosted GitHub Actions automatically.
2. Development qualification should prefer local/tooling execution and the least-expensive suitable environment.
3. Linux-hosted GitHub jobs are used only when remote reproducibility, service containers, immutable-artifact binding or independent CI evidence materially adds value.
4. macOS Apple Silicon and Windows x64 GitHub-hosted runners are reserved for the exact **final release candidate / Stable deliverables** that require native packaging and actual-artifact runtime certification, or for a proven platform-specific defect that cannot be qualified elsewhere.
5. A routine TEST-native pass MUST NOT be followed by an equivalent Stable-native pass when the Stable transition is behavior-preserving/identity-only and the final Stable source will receive fresh full certification. In that case, run native G13/G14/G15 once on the final Stable deliverables.
6. A pre-Stable native pass is justified only when it resolves a material platform uncertainty that cannot safely wait for final Stable certification; it is not a default release step.
7. Native runners MUST NOT run for documentation-only, governance-only, release-tooling-only or unrelated platform-neutral changes.

## Trigger policy
DE.PULSE release tooling commonly lives on non-default branches, so expensive release workflows use **dedicated intent-trigger files** rather than broad push triggers or default-branch-dependent manual dispatch.

Required pattern:
- each expensive workflow watches exactly one `.depulse-certification/triggers/<workflow>.json` path on its tooling branch;
- ordinary commits cannot match that path and therefore cannot start the workflow;
- a run is requested only by intentionally changing that trigger file with `run=true`, an incremented nonce, and any exact prerequisite run/commit/fingerprint values required by the workflow;
- the workflow must validate the trigger contents before doing expensive work;
- one trigger file controls one workflow only;
- updating workflow/source/docs files alone must not start release certification.

If an automatic development workflow is justified later, it MUST use narrow path filters, cheap-first prerequisite jobs, and cancellation/concurrency controls. Expensive jobs require explicit dependencies so they cannot start before cheaper prerequisite gates pass.

## Evidence reuse
Evidence may be reused only when all of the following are unchanged:
- exact source commit / fingerprint relevant to that evidence;
- product semantics relevant to that evidence;
- gate implementation / test semantics relevant to the evidence;
- platform/runtime requirement relevant to the evidence.

A release identity-only transition does not invalidate already-proven product semantics, but the final release identity/source still requires the gate families that materially depend on that identity/runtime. Do not rerun already-passed expensive gates solely because documentation, orchestration metadata or unrelated tooling changed. If product/source semantics changed, invalidate and rerun the affected gate family.

## Progressive qualification
Use the cheapest blocking order:
1. static/source/governance checks;
2. focused tests;
3. full Linux/unit/integration qualification where required;
4. real PostgreSQL / browser / heavy performance qualification only when applicable;
5. macOS Apple Silicon + Windows x64 packaging/runtime only on the final closure candidate unless a specific earlier platform uncertainty requires otherwise.

A failed cheap prerequisite must prevent all downstream expensive jobs.

## Native closure
Required release targets remain:
- macOS Apple Silicon;
- Windows x64.

For a Stable release, native G13/G14 evidence must be produced from the exact immutable final candidate and G15 must verify artifact hashes/source binding before promotion/publication. No native gate may be waived to save Actions minutes.

Prefer one parameterized native-audit implementation per OS rather than duplicated TEST/STABLE scripts. Release identity, runtime profile, bundle/executable names and artifact names are parameters; audit semantics remain shared. The normal path is **one native pass on the final Stable artifacts**.

## Release publication
GitHub Releases remain the canonical durable location for final runnable Stable artifacts, source archive, checksums, provenance and certification evidence. Upload/publish operations should reuse already-certified artifacts; they must not rebuild binaries merely to publish them.

Publication should run as a small metadata/upload job only after Stable G15 passes. It may fast-forward `main`, create the immutable Stable tag, generate G16/provenance/checksum metadata and upload the G15-certified bundle, but it must perform no compilation or packaging rebuild.

The immutable Stable tag is exact product/release truth. Post-tag G16 governance or process metadata may be synchronized to `main` without moving the Stable tag, provided it does not alter the certified product artifact/source identity.

## Budget-aware behavior
- A zero/limited GitHub Actions budget is not a reason to lower certification standards.
- If hosted minutes are unavailable, freeze the exact qualified source and continue all runner-independent work.
- Resume only the blocked hosted/native gate when capacity returns; do not restart unrelated gates unless evidence validity requires it.
- Prefer self-hosted runners when they provide trustworthy equivalent platform evidence and the runner itself is controlled and qualified.

## Cost-efficiency acceptance criteria
A DE.PULSE release process is considered efficient only when:
- no expensive job runs without a blocking reason;
- ordinary commits cannot trigger release certification;
- macOS/Windows normally run once on the exact final closure candidate unless a proven platform-specific uncertainty or native failure requires another pass;
- macOS and Windows run in parallel once their shared cheaper prerequisite is green;
- unchanged certified evidence is reused safely;
- failed prerequisites stop downstream compute;
- Stable identity promotion does not rerun a full suite that the immediately following Stable G12 will rerun authoritatively;
- only one fresh Stable G12 is run for the exact promoted Stable candidate;
- final publication does not rebuild certified artifacts;
- repository changes cannot silently trigger legacy historical workflows;
- post-tag G16 governance sync does not mutate the immutable Stable tag or falsely broaden certified product scope.

This contract supplements the permanent G0-G16 model; it does not add a new top-level release gate.
