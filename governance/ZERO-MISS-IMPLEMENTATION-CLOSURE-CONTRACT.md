# DE.PULSE — Strict Zero-Miss Implementation & Closure Contract

**Status:** PERMANENT / GOVERNING / FAIL-CLOSED  
**Applies to:** Adaptive Roadmap, Adaptive Build Plan, Adaptive Build Process, Adaptive Delivery Process, every work slice, dependency band, closure ledger, handoff and release.  
**Core invariant:** `Framework ≠ Implementation ≠ Production Integration ≠ Verification ≠ Release Qualification`.

## 1. Non-negotiable completion rule

DE.PULSE may not treat a capability, requirement group, dependency band, version, closure item or handoff item as complete merely because code, types, helpers, contracts, mocks, configuration, documentation or unit tests exist.

The only allowed forward lifecycle is:

`PLANNED → IMPLEMENTED → PRODUCTION_INTEGRATED → VERIFIED → RELEASE_QUALIFIED`

Intermediate/open states may exist only while work is actively incomplete. They are **never** a completion state, advancement state, handoff-as-done state, merge authorization or release qualification.

A dependency band may advance only when every mandatory responsibility owned by that band is `VERIFIED` or explicitly `NOT_APPLICABLE` from G1 with evidence. Mandatory unresolved work may not be silently carried forward.

## 2. Zero-Miss Closure Chain

Before a requirement can become `VERIFIED`, all applicable links below must be proven on current source:

`requirement → canonical owner → production integration → consumer reachability → positive behavior → adverse/fail-closed behavior → persistence/restart/lifecycle behavior → security/rights/privacy behavior → observability/diagnosability → executable regression evidence → real integration/infrastructure/external evidence when required → exact-head CI → closure ledger`

Missing one mandatory link means the requirement is not VERIFIED.

## 3. Strict evidence semantics

The following never prove completion by themselves:
- a model/type/schema exists;
- a helper/evaluator exists;
- documentation says implemented;
- an issue/PR is closed;
- a unit test passes;
- Fast CI is green when the requirement needs real integration/infrastructure evidence;
- a runtime environment variable declares readiness;
- an API key/connection succeeds;
- a provider is configured;
- a mock/stub proves a production boundary;
- historical evidence passed on a different source fingerprint.

Evidence must prove the actual requirement at the actual production boundary it governs.

## 4. Adaptive Roadmap rule

Roadmap/version placement conserves scope but never proves implementation. Every planned capability carries its lifecycle state separately. A version may not be declared complete while any mandatory owned responsibility is below `VERIFIED`.

Small coherent versions/chunks remain preferred, but smaller scope does not weaken closure depth.

## 5. Adaptive Build Plan rule

Before implementation, every applicable responsibility must identify:
- canonical source/state owner;
- production consumers/integration points;
- positive and adverse/fail-closed behavior;
- persistence/restart/lifecycle impact;
- security/rights/privacy impact;
- observability requirements;
- executable regression owner;
- real external/infrastructure evidence required;
- exact-head CI/qualification requirement.

If any of these cannot be resolved, the build plan must expose the blocker before the item is called implemented.

## 6. Adaptive Build Process rule

Implementation is not complete until production reachability is proven. Build work must follow:

`REUSE → CONSOLIDATE → REFACTOR → DELETE/REPLACE → ADD → INTEGRATE → PROVE`

For each changed requirement, perform a source-driven reverse scan from runtime/consumer back to requirement so orphan helpers, unconsumed evaluators, duplicate owners and declaration-only controls are detected before closure.

No dependency band advances merely because its happy-path tests are green.

## 7. Adaptive Delivery Process rule

Delivery must fail closed when closure evidence is incomplete. Before merge/release/handoff-as-complete:
- every mandatory requirement is VERIFIED;
- exact-head Fast evidence is current;
- impact-selected Qualified evidence is current where required;
- real infrastructure/database/provider/security/platform evidence exists where required;
- closure ledger and handoff match executable truth;
- no known material implementation miss is hidden as deferred cleanup.

Only a fully verified release candidate may become `RELEASE_QUALIFIED`.

## 8. No partial-completion masquerade

`OPEN`, `PARTIAL`, `IMPLEMENTED_UNVERIFIED`, `CONTRACT_ONLY`, `FRAMEWORK_ONLY`, `DECLARED_READY`, or equivalent states are truthful diagnostic states only. They must never be summarized as "done", "closed", "complete", "ready", "10/10", "production-ready" or used to authorize the next dependency band.

If an audit discovers a previously missed mandatory integration/evidence link, the item automatically reopens and downstream completion claims are invalid until the miss is corrected and requalified.

## 9. External-blocker rule

A mandatory external dependency is not a reason to mark a requirement VERIFIED. It remains blocking until one of these occurs **before release advancement**:
1. the required real evidence is obtained and verified; or
2. G1 is explicitly rebaselined with durable rationale and user-approved scope disposition showing the responsibility is genuinely not applicable to the current release.

No post-hoc scope shrinking is allowed merely to obtain closure.

## 10. Handoff / AI portability rule

Every assistant, account, model or developer must independently reconstruct completion from GitHub/source/executable evidence. Chat memory is never required and cannot upgrade a lifecycle state.

Every meaningful handoff must state:
- exact live branch/head;
- earliest not-yet-VERIFIED dependency band;
- exact missing closure-chain links;
- latest executable evidence;
- whether the current candidate is allowed to advance.

## 11. Current v19 corrective application

The 2026-08-26 v19 zero-miss audit is the motivating incident. Existing truthful OPEN/PARTIAL rows remain truthful until actually fixed; this contract does **not** relabel missing work as VERIFIED. Instead, it prevents any such row from being treated as completed or skipped.

For v19.0.0, execution resumes from the earliest reopened dependency band and may advance only after full `VERIFIED` evidence under this contract.

## 12. Enforcement priority

When prose, historical status, issue state or prior CI conflicts with current executable evidence, use this priority:

`current source/runtime truth → current executable/infrastructure evidence → closure ledger → current handoff/governance → issue/PR prose → chat history`

A false positive completion is considered a governance defect and must be corrected immediately.