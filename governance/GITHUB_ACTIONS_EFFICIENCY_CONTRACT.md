# DE.PULSE GitHub Actions Efficiency Contract

Status: PERMANENT BUILD / DELIVERY GOVERNANCE

## Objective
Optimize **cost per trustworthy, non-duplicative evidence**. GitHub Actions runs whenever hosted execution is required for quality, reproducibility, platform truth, security or immutable release evidence. Spend is optimized through selection and reuse; it is never capped by lowering assurance.

Control ID: `CI-ADAPTIVE-18.5.1-001`.

## Default execution model
1. Every change receives a machine-readable impact decision. Affected source/config/workflow changes run automatic cheap preflight; documentation/governance-only changes run only the relevant schema/link/gate checks.
2. Development qualification may use local/tooling execution for fast feedback, but local success does not replace required independent CI evidence.
3. Linux-hosted jobs provide remote reproducibility, service-container, browser, security, integration and immutable-artifact evidence when selected by impact or mandatory policy.
4. macOS Apple Silicon and Windows x64 hosted runners are reserved for exact final-candidate packaging/runtime certification or an earlier proven platform-specific uncertainty.
5. A routine TEST-native pass must not be duplicated by an equivalent Stable-native pass. Prefer one authoritative native pass on the exact final candidate unless source or relevant evidence changed.
6. Native and other expensive jobs start only after their cheaper shared prerequisites pass.
7. The planner may skip only fingerprint-proven unaffected lanes. Deterministic guardrails own all mandatory gates and may not be overridden by budget.

## Trigger and workflow policy
Keep a small durable workflow surface:

- `ci-fast.yml`: automatic cheap preflight with precise path/requirement classification;
- `ci-qualified.yml`: reusable affected-area qualification invoked by the planner or `workflow_dispatch`;
- `release.yml`: exact-candidate G10/G12/G13/G14/G15/G16 orchestration and no-rebuild publication.

Requirements:

- use reusable workflows/composite actions and version/config inputs rather than copying TEST/Stable implementations;
- use `workflow_dispatch` inputs for diagnostics; do not commit temporary/one-shot workflows or trigger-file churn merely to request a run;
- CI must not edit, delete or push product/workflow source;
- default token permissions are read-only; only isolated publication receives the narrow write permissions it requires;
- use path filters, cheap-first dependencies, concurrency cancellation and bounded matrices;
- workflow/config changes must trigger workflow lint and relevant harness tests;
- stale historical/version-specific workflows must be retired from the active branch after inventory.

## Adaptive evidence planner

Planner input:

`diff → canonical owners → dependency blast radius → invalidated evidence → historical failures → required lanes → estimated runtime/cost`

Planner output must record:

- requirement IDs and source/dependency fingerprints;
- selected and skipped lanes with explicit reasons;
- mandatory-policy overrides;
- expected and actual runtime by runner OS;
- estimated cost, cache hit/miss and artifact retention;
- failure class, retry relationship and reused/invalidated evidence.

The planner may add risk-responsive lanes and skip only unaffected lanes. It cannot waive release reconciliation, security, native, artifact-provenance or other mandatory evidence. Planner policy changes follow `SHADOW → VALIDATED → APPROVED → PRODUCTION`; no CI job silently self-modifies it.

## Failure classification and retry

Every non-PASS is classified before retry:

- `PRODUCT_FAIL` — product behavior/source defect;
- `GATE_TEST_FAIL` — incorrect or stale assertion/test contract;
- `CI_HARNESS_FAIL` — orchestration, portability, readiness or artifact-handling defect;
- `INFRA_FAIL` — runner/provider/network/service failure outside product truth;
- `EXPECTED_NOOP` — correctly skipped/no-change result;
- `SUPERSEDED` — obsolete evidence replaced by a newer authoritative run.

Retry only the failed lane and dependent evidence invalidated by the correction. Preserve independent matching-fingerprint PASS evidence.

## Portability, cache and retention controls

- All text reads are explicitly UTF-8.
- Permission assertions are OS-aware; POSIX mode expectations do not create false Windows failures.
- Native readiness uses bounded process/port/health probes and captures diagnostics, not fragile timing against one log line.
- Enable supported dependency caches and record cache effectiveness.
- Development artifacts default to 3–7 days, failure diagnostics 7–14 days and RC evidence 30 days. Certified binaries, source, checksums and provenance remain in the immutable GitHub Release.
- Optional controlled self-hosted runners may accelerate iterative native debugging; clean exact-candidate native release proof remains independently reproducible.

## Budget policy

The user budget—currently $5—is a **soft alert and learning signal**, not a hard quality ceiling. A separate higher emergency cap may stop runaway loops. A legitimate required CI lane must never be skipped, weakened or falsely marked PASS because of cost.

When an alert is crossed, continue any already-authorized required evidence, diagnose duplicated work/cache misses/retries, replan the smallest remaining lane and record the decision for G16.

## 2026-08-17 harness lessons

Run `32009262146` demonstrated three CI-harness defects rather than product proof: default Windows text decoding rejected UTF-8 documents, a POSIX archive-permission assertion produced a Windows false failure, and native readiness timed out after the application had already logged a live local terminal. The later stable run `32049647088` completed G11, G12, macOS, Windows, G15 and G16.

Permanent prevention:

- test the harness on both target OS families before expensive certification;
- distinguish harness/infra failure from product failure;
- never discard successful native evidence merely because orchestration was inefficient;
- treat the native lanes as necessary quality evidence and remove only duplicated orchestration.


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
