# DE.PULSE — Current Adaptive Build Plan

**Operational overlay date:** 2026-08-19  
**Immutable Stable baseline:** `v18.6.1-stable`  
**Current engineering slice:** Phase 0 Packet E — renderer modularization foundation / Documentation owner  
**Authority:** current execution overlay. Permanent contracts and historical release evidence remain intact.

## 1. Normal target lifecycle

`1 development branch → batch coherent branch work → 1 Draft PR → Fast → same PR Ready → Qualified → merge → Release only when release identity/release workflow requires it → Stable`

Rules:

- Build a coherent branch before opening the PR whenever practical so preparatory file writes do not manufacture `synchronize` events.
- Never create trigger/retry/certification/promotion branches.
- Never create a second PR just to retrigger CI.
- Fix source/test/gate defects on the same branch and same PR; return to Draft when appropriate, then obtain new exact-head Fast and Qualified evidence.
- Same-SHA infrastructure failure reruns only affected work when possible.
- Main push performs hygiene only.
- Publication uses exact same-run certified native artifacts; no post-certification rebuild.

## 2. G0–G16 execution map

- **G0 Exact Stable Intake:** immutable Stable identity, source SHA/fingerprint, open defects/issues, CI state and dependencies.
- **G1 Immutable Scope:** every committed scope item has stable traceability; no silent additions/removals.
- **G2 Architecture / Data Utility:** owner, consumer, provider, entitlement/rights, source of truth, reuse, freshness, retention and duplication.
- **G3 Design / Dependency / Impact Readiness:** Impact Planner classifies affected surfaces, tests, portability, browser risk and expected CI cost.
- **G4 Development Exit:** one version-development branch, one Draft PR, clean source and scope traceability.
- **G5 Fast Qualification:** cheap exact-head syntax/format/unit/contract checks for affected risk.
- **G6 Integration / Medium Qualification:** affected integration and cross-module evidence.
- **G7 Data / Security / Adaptive Intelligence:** provider/data-rights/security/adaptive evidence when applicable.
- **G8 Performance / Capacity / Stability:** load/runtime/backpressure/stability evidence when applicable.
- **G9 Cross-Module / UI / UX:** affected renderer/browser/interaction evidence.
- **G10 Pre-Freeze Qualified Candidate:** exact-head Qualified success; Release Rehearsal for CI/release changes; Chrome + WebKit primary browser evidence whenever selected risk requires browser qualification; compact CI telemetry retained.
- **G11 Immutable Release Candidate:** merged candidate bound to exact Fast + Qualified source head and equal source fingerprint.
- **G12 Full Certification:** authoritative full certification from immutable candidate.
- **G13 Native Packaging / Provenance:** required native packages from candidate.
- **G14 Actual Artifact Runtime Audit:** packaged macOS Apple Silicon and Windows x64 behavior/provenance.
- **G15 Release Assurance / Promotion:** native evidence graphs and exact artifact hashes verified.
- **Publish:** exact certified artifacts only; no rebuild.
- **G16 Adaptive Retrospective / Handoff:** current source of truth, durable release evidence, defects, CI performance and next intake.

No new top-level gates beyond G0–G16.

## 3. Impact Planner v2

Change classes:

- `CI_HARNESS`
- `RELEASE_TOOLING`
- `BACKEND`
- `RENDERER_UI`
- `AUTH_SECURITY`
- `PROVIDER_ROUTER`
- `DATA_RIGHTS`
- `PERSISTENCE`
- `RELIABILITY_PERFORMANCE`
- `CERTIFICATION_GOVERNANCE`

Unknown non-process content fails closed to full qualification.

### Lane selection

- Process/governance/CI-only: `ci-harness` + portability.
- Product/mixed/uncertain: `full` Qualified.
- Normal `full` and explicit `browser` qualification: both primary browser engines must pass.
- `RENDERER_UI`: WebKit evidence is mandatory even if the lane is narrowed.
- WebKit harness/routing changes: real WebKit evidence is mandatory while remaining process-only.
- Backend/provider-only narrowed work: no unnecessary browser-engine runtime.
- `CI_HARNESS` or `RELEASE_TOOLING`: Release Rehearsal required through workflow policy.
- A file matching `release/v*/stable-evidence-manifest.json` is a retrospective evidence index and is process-only, but remains `RELEASE_TOOLING` governed. Other release scripts/artifacts remain full-qualification scope.

### Browser policy

- **Chrome and WebKit are co-primary browser engines.**
- Chrome carries the broad behavioral regression suite.
- WebKit carries the primary cross-engine compatibility suite for core UI/interaction contracts.
- `full` and `browser` candidates require both.
- Renderer/UI changes require WebKit through `webkit_required`.
- Primary WebKit executes on `macos-15` with exact pinned Playwright plus `playwright install webkit`; Linux `--with-deps webkit` is prohibited.
- Other engines, including Firefox if introduced later, remain secondary/risk-directed unless evidence justifies promotion.

## 4. Failure taxonomy

- `PRODUCT_FAIL`
- `GATE_TEST_FAIL`
- `CI_HARNESS_FAIL`
- `INFRA_FAIL`
- `EXPECTED_NOOP`
- `SUPERSEDED`

Failure classification never permits bypassing required quality gates.

## 5. Phase 0 status

### Packet A — COMPLETE

Merged PR #46. Delivered Impact Planner v2, Release Rehearsal, current governance overlays, honest v18.6.1 reconciliation and portable handoff.

### Packet B — COMPLETE

Merged PR #47. Delivered immutable Fast/Qualified Action pins, dependency lock, Playwright `1.62.0` pin, deterministic pip caching and reproducibility/permission gate. Release-workflow Action pinning remains deferred to the next genuine release-capable product slice.

### Packet C — COMPLETE

Merged PR #48 at `23ecb71f60e1658d68bcef6248044ce53b6dd851`.

- Chrome + WebKit co-primary policy implemented.
- Primary WebKit runs on `macos-15` without Linux apt amplification.
- Core watchlist/global-remove, membership semantics, short-height Settings save bar and centered header compatibility are covered.
- Final Fast #393 PASS.
- Final Qualified #138 PASS with real WebKit + Ubuntu/macOS/Windows portability.

### Packet D — COMPLETE

Merged PR #49 at `2885de409c86f771d582f09f54e0f6c564f6c59d`.

Delivered durable Stable evidence indexing, release-evidence drift validation, Qualified queue/runtime/platform telemetry, workflow-amplification warnings, browser setup/cache observability, 30-day compact telemetry retention, and zero-network workflow structural lint. Final Fast #396 and Qualified #139 passed through the intended process-only path.

### Packet E — ACTIVE: Documentation capability owner

Fresh renderer inventory shows `renderer.js` at about 425 KB and `styles.css` about 316 KB. Runtime composition still begins with the classic monolith and then applies specialized layers. The correct first move is ownership extraction, not a big-bang rewrite.

Bound Packet E scope:

1. **Active owner:** `renderer/documentation-ui.js` owns Documentation Markdown, hydration and view rendering.
2. **Load order:** `renderer.js` → `documentation-ui.js` → later compatibility layers → `documentation-access-v18.6.js`.
3. **Access decorator:** existing v18.6 role policy wraps the active owner and registers itself as a decorator; it does not create another authorization API.
4. **Owner registry:** `__DE_PULSE_RENDERER_OWNERS__.documentation` records owner, responsibilities, dependencies, decorator and deletion gate.
5. **Truthful strangler state:** `ACTIVE_OWNER_WITH_LEGACY_FALLBACK`. Old monolith Documentation functions remain present but are inactive after the capability owner loads.
6. **Explicit remaining coupling:** Documentation Markdown still delegates architecture diagrams to legacy `architectureDiagram`; this is recorded, not hidden.
7. **New naming rule:** long-lived owner is capability-oriented; do not introduce `documentation-ui-v18.x.js`.
8. **Static owner gate:** load order, duplicate loads, naming, fallback truth, access wrapping and primary-engine test wiring fail closed.
9. **Fast proof:** existing v18.6 Documentation role/access regression is owner-aware.
10. **Qualified renderer proof:** `documentation_ui_owner_test.js` verifies owner replacement, Markdown, hydration and role decorator integration.
11. **Primary engines:** the same focused owner behavior is executed in Chrome and WebKit through `documentation_owner_browser_test.py`.
12. **Deterministic safety:** market math/scoring is untouched; existing deterministic equivalence remains mandatory.

### Packet E qualification requirement

Because Packet E changes renderer/product source, Impact Planner must choose `full`. Required evidence:

- Fast exact-head PASS, including owner/static/access checks;
- backend full/race/randomized PASS even though no backend semantics should change, because mixed product + CI wiring fails closed;
- renderer syntax + deterministic equivalence + renderer logic + owner regression PASS;
- Chrome broad suite + direct Documentation owner Chrome PASS;
- WebKit core compatibility + direct Documentation owner WebKit PASS;
- Qualified exact-head PASS;
- telemetry retained.

No Stable Release is expected unless release identity or the canonical Release workflow is deliberately changed; Packet E does neither.

### Packet E physical-deletion gate

This packet establishes active ownership but does **not** physically delete the old monolith Documentation definitions. Physical deletion is a later controlled step only after no consumer needs the fallback, equivalence evidence is direct, both primary engines pass, and the owner state can truthfully move to a no-fallback designation. This avoids a risky monolith rewrite while still removing runtime ownership ambiguity now.

## 6. CI telemetry contract

Qualified emits compact operational evidence with:

- exact candidate SHA and selected lane;
- per-job queue seconds and execution seconds;
- consumed runner minutes split into Linux, macOS, Windows and unknown;
- Chrome/WebKit dependency setup seconds and pip cache-hit state when applicable;
- current-PR Fast/Qualified/Release run counts;
- amplification warnings above conservative thresholds;
- explicit `actualCurrencyCost: null` because GitHub billing remains the authority for financial cost.

Telemetry is diagnostic/operational evidence. It cannot replace functional qualification, exact-head statuses, native release evidence or the immutable Stable evidence manifest.

## 7. Source and repository hygiene

A file being several days old is not a defect. Do not touch unchanged files merely to refresh GitHub dates. Remove a file only when references/consumers/evidence needs are proven absent and protected history is unaffected.

Historical certification, governance, approved reference assets and actively loaded compatibility layers stay until explicit cleanup proof says otherwise.

## 8. Quality floor

Efficiency may reduce duplicate work, runner choice or irrelevant lanes; it may not reduce:

- exact-source provenance;
- deterministic tests;
- Chrome + WebKit primary evidence when browser qualification is required;
- data/security/rights controls;
- macOS Apple Silicon + Windows x64 Stable certification;
- same-artifact publication;
- conserved requirement traceability;
- No Execution and other permanent product boundaries.

## 9. Product intake after Phase 0

After Packet E reaches its governed exit, run fresh G0–G3 against current reconciliation and select the highest-value coherent v18.x product slice. Provisional priority order remains: user-trust defects → runtime/ADR-GDI reliability → shared intelligence utility consolidation → renderer maintainability → controlled TradeInsight SHADOW integration.
