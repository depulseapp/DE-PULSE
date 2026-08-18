# DE.PULSE — Current Authoritative Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS.** GitHub is authoritative; chat history and temporary workspaces are advisory only.

**Release:** `v18.6.0`  
**Active branch:** `v18.6-development`  
**Stable predecessor:** `v18.5.2-stable` / G0–G16 CLOSED  
**Current candidate state:** v18.6.0 G12 BROWSER-HARNESS REMEDIATION / FRESH G10 REQUIRED / NOT PROMOTED  
**Repository:** `depulseapp/DE-PULSE`  
**Main release PR:** `#16`  
**v18.6 runtime build ID:** `v18.6.0-stable-20260818`  
**Last updated:** 2026-08-18 America/Vancouver

## Resume rule

Read `AGENTS.md` or `CLAUDE.md`, then reconcile this handoff, `release_identity.json`, both `.depulse-certification/resume/` checkpoints, PR #16, current checks, the release-certification branch, and the immutable Stable predecessor. Never resume from model memory alone.

`source_fingerprint.py` excludes `.depulse-certification` only. Workflow definitions, certification harnesses, browser proofs/contracts, and this handoff are source-fingerprinted. Any change to them requires fresh G10 before G11.

## Product / architecture state

All eight assigned v18.6 implementation/audit slices remain code-complete: watchlist remediation; shared Scanner/Radar broad-snapshot acquisition; serialized Session Intelligence Coordinator; Market Activity/legacy-route consolidation; role-aware documentation; external dependency/provider readiness; bounded AI context/cache/schema/evaluation hardening; and provider×dataset rights-aware fail-closed AI egress.

Protected invariants remain unchanged: deterministic Day/Swing/Long formulas; Smart Provider Router sole routing ownership; provider count never changes Market Mode; GLD/SLV/USO tradable live exceptions; desktop SQLite / hosted PostgreSQL; U.S. Equities Processing Boundary; permanent No Execution Boundary; exactly three canonical workflows; G0–G16 only.

## Trustworthy qualification before this remediation

Source `a84b7028f8723f47f11a59e2225c10ddf1a38e3b`, exercised at metadata head `66fd480c3fcf91f1ce56a6077dff323625aaf0b3`, passed fresh G10 after the first G12 watchlist-harness correction:
- CI Fast #242 / run `32197304582` — PASS;
- CI Qualified #102 / run `32197304782` — PASS;
- harness contract, Ubuntu/macOS/Windows portability, browser behavior, renderer contracts, full Go suite, race detector, randomized package order and final evidence summary all passed.

Release-certification base `a6b58016f4921e5927579fc4eae4ea5e81f026ad` was dispatched through unmerged trigger PR #25. Dispatcher CI Fast #246 / run `32198077508` passed and resolved:
- release ref `v18.6-release-certification`;
- candidate `a6b58016f4921e5927579fc4eae4ea5e81f026ad`;
- canonical source fingerprint `f38e3815307382ec111c54558b36b503a06ab143fb1a5b84a279e0f48f1092ba`;
- canonical release run `32198058592`;
- `publish=false`.

The connected GitHub app bound that exact run identity to PR #16 because the workflow token still cannot post PR comments. Trigger PR #25 was then closed unmerged after G12 failed. No Stable publication occurred.

## G12 progress and second failure classification

Canonical release run `32198058592` proved the first harness remediation worked: G12 passed the hardened watchlist contract and did **not** execute the obsolete v18.5.1 CURRENT/aria-current membership proof. It then passed Go full suite, race detector, randomized order, deterministic equivalence 2403/2403, renderer logic, v18.0.5 role/responsive acceptance, v18.6 surface consolidation, documentation access, live DOM reconciliation, and first-run auth-copy browser proof.

The next failure was in historical `release/v18.5.1/browser_ui_hierarchy_test.py`, before its behavioral browser assertions:

`assert "ui-v18.5.1.css?v=18.5.2" in INDEX`

Current `renderer/index.html` still deliberately loads the retained implementation layers `ui-v18.5.1.css` and `header-v18.5.1.js`, but their canonical cache-busters are now `?v=18.6.0`. Therefore this is another **release-harness/version-binding defect, not a demonstrated product hierarchy regression**.

The remaining v18.5.2 browser tests for master-symbol input, profile/display-name, and Settings save-bar are behavior-first and were reviewed for this issue; they do not contain the obsolete `?v=18.5.2` index-asset assertion and remain in G12. The already-passing v18.5.1 live-render and auth-copy proofs also remain.

## Consolidated second remediation

Historical `release/v18.5.1/browser_ui_hierarchy_test.py` remains untouched for audit/history.

A dedicated `release/v18.6.0/browser_ui_hierarchy_test.py` now carries the same substantive responsive/header/Research behavior proof while deriving the expected cache-buster from canonical `release_identity.json`. It requires the actual current assets:
- `ui-v18.5.1.css?v=<canonical v18.6 release version>`;
- `header-v18.5.1.js?v=<canonical v18.6 release version>`.

`release/v18.6.0/run_full_certification.sh` now executes the v18.6 hierarchy proof instead of the historical v18.5.1 hierarchy proof, alongside the existing v18.6 watchlist proof and retained compatible legacy behavior tests.

`tools/ci/workflow_policy.py` now prevents regression by requiring both v18.6 browser proofs in G12, forbidding both superseded v18.5.1 hierarchy/watchlist proofs from the v18.6 G12 list, and requiring the v18.6 hierarchy proof to derive its asset cache-buster from canonical release identity rather than hard-code `?v=18.5.2`.

Because these files are source-fingerprinted, the prior G10/G11 evidence cannot certify this second harness remediation. Fresh G10 is mandatory before another canonical release run.

## Release orchestration contract

Certification remains non-publishing. Exact release-certification push or the owner-gated PR fallback may certify only with `publish=false`; the PR fallback uses immutable base ref/SHA. Stable promotion remains exact `*-stable-promotion`, push-only, owner-gated, evidence-bound and no-rebuild. There is no PR fallback for Stable publication.

Workflow-token PR comments are best-effort observability. Exact dispatcher/release run identity remains in Actions evidence and is bound to PR #16 through the connected GitHub app when needed.

## Exactly one next action

**Qualify this consolidated v18.6 hierarchy-harness remediation as one fresh Fast + Qualified G10 candidate. If both pass, bind only fingerprint-excluded checkpoints, reconcile `v18.6-release-certification`, create one unmerged fingerprint-excluded trigger PR, dispatch one canonical `release.yml` run with `publish=false`, bind its exact identity to PR #16, then continue G11 → G12 → G13/G14 macOS + Windows → G15 → G16. Do not promote Stable unless every gate passes.**

## Known residuals / User Action Required

- TradeInsight remains governed future/shadow implementation work until separately validated and approved; it gains no production influence automatically.
- Deployment-specific provider keys, entitlements and commercial/redistribution/AI-use rights remain User Action Required where absent or unapproved and fail closed.
- v18.6 native macOS Apple Silicon and Windows x64 evidence remains pending until G12 and native G13/G14 pass.

## Provider-neutral continuation instruction

> Connect to `depulseapp/DE-PULSE`, read `AGENTS.md` or `CLAUDE.md`, then inspect `release_identity.json`, `handoff/CURRENT.md`, both `.depulse-certification/resume/` checkpoints, PR #16, active checks, and the release-certification branch. Treat `v18.5.2-stable` as immutable Stable until v18.6 G11–G16 and no-rebuild promotion complete. Resume from the exact trustworthy evidence above, not chat memory. Preserve G0–G16, assistant/account independence, Smart Provider Router ownership, deterministic desk formulas, U.S. Equities Processing Boundary and permanent No Execution Boundary.
