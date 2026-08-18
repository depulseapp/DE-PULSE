# DE.PULSE — Current Authoritative Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS.** GitHub is authoritative; chat history and temporary workspaces are advisory only.

**Release:** `v18.6.0`  
**Active branch:** `v18.6-development`  
**Stable predecessor:** v18.5.2 STABLE / G0–G16 CLOSED  
**Current candidate state:** v18.6.0 NON-BLOCKING RELEASE TRACKING FIX / FRESH G10 REQUIRED / NOT PROMOTED  
**Repository:** `depulseapp/DE-PULSE`  
**Main release PR:** `#16`  
**Stable predecessor tag:** `v18.5.2-stable`  
**Stable predecessor promotion commit:** `d30e54db4908ca57c52ae298cc4ada3416fab46b`  
**Stable predecessor certified source fingerprint:** `807de082d43e83d1d3548bca9350d13b72ef4dc71a848940b73b63e4b4d215b0`  
**v18.6 runtime build ID:** `v18.6.0-stable-20260818`  
**Last updated:** 2026-08-18 America/Vancouver

## Resume rule

Read `AGENTS.md` or `CLAUDE.md`, then reconcile this handoff, `release_identity.json`, both `.depulse-certification/resume/` checkpoints, actual branch/PR/check state and the immutable Stable predecessor tag/release. Never resume from model memory alone. `v18.5.2-stable` remains the certified product authority until v18.6.0 completes G11–G16 and is promoted.

`source_fingerprint.py` excludes `.depulse-certification` but **does not exclude `handoff/CURRENT.md` or workflow definitions**. Machine resume/evidence checkpoints may advance as fingerprint-excluded operational metadata after a source candidate. This handoff and `.github/workflows/*.yml` are candidate-bound source; any change to them requires fresh G10 before G11.

## v18.6.0 implementation state

All eight assigned v18.6 implementation/audit slices remain code-complete:

- watchlist membership/add remediation with regression coverage;
- shared Scanner / Opportunity Radar broad-snapshot acquisition with bounded cache, freshness reuse, partial-miss fetch and coalescing;
- one serialized Session Intelligence Coordinator for Pre-Market Prep and Market Open Prep;
- Market Activity demotion to supporting drill-down plus legacy evidence-route consolidation;
- role-aware documentation with server-authoritative privileged access and Documentation Impact evidence;
- external dependency/provider readiness with durable User Action Required tracking and fail-closed entitlement/rights governance;
- bounded AI context, full cache identity + TTL, strict structured-output/citation validation, safe abstention and continuous-evaluation telemetry;
- provider×dataset rights-aware AI egress with unknown/denied evidence withheld before external model calls.

Protected deterministic Day/Swing/Long formulas remain unchanged. Smart Provider Router remains the sole provider-routing owner; provider count does not change Market Mode. GLD/SLV/USO tradable exceptions, desktop SQLite / hosted PostgreSQL architecture, the U.S. Equities Processing Boundary and the permanent No Execution Boundary remain preserved.

## Last exact-source G10 PASS and final orchestration finding

Source `eafc1d56aee2597b35f06c45f333d1275e801de0` passed the REST-tracking connector-safe qualification:

- CI Fast #224 / run `32195517663` / attempt 1 — PASS;
- CI Qualified #95 / run `32195517662` / attempt 1 — PASS;
- coverage included REST tracking/connector-safe dispatcher policy, Ubuntu/macOS/Windows portability, workflow/provenance contracts, browser behavior, renderer contracts, full Go suite, race detector, randomized package order and final evidence summary.

The release-certification branch was then reconciled to base candidate `5086a2156c15fa25df1d869e3f2875e6eb397bca`, preserving the final G10 source and excluded checkpoint metadata. Trigger PR #23 proved the remaining path precisely:

- Fast preflight passed;
- the release-dispatch job ran;
- release identity passed for `18.6.0` / `v18.6.0-stable-20260818`;
- immutable release ref resolved to `v18.6-release-certification`;
- immutable candidate resolved to base SHA `5086a2156c15fa25df1d869e3f2875e6eb397bca`;
- canonical source fingerprint resolved to `2837843eb9ec67f8edb88052efcb1defe954c0b0419d4f29d1a040d722401231`;
- `publish=false` was enforced;
- the workflow token still received HTTP 403 `Resource not accessible by integration` when posting the REST issue comment to PR #16, despite the job exposing `Issues: write`.

The remaining defect is therefore **tracking transport permission only**, not release identity, source provenance, product code, runner allocation, or certification authorization.

## Final release-tracking contract

PR comments are observability evidence and must not be allowed to block otherwise-valid non-publishing certification. The dispatcher therefore now:

1. attempts both PR #16 tracking comments through the REST issue-comments endpoint;
2. treats an integration-level comment denial as a warning rather than a certification failure;
3. continues to dispatch/reuse canonical `release.yml` only after exact release ref, candidate SHA, fingerprint and `publish=false` have been resolved;
4. writes the canonical release run ID/URL, release ref, candidate SHA, fingerprint and publish state into the Actions job summary/log regardless of PR-comment permission;
5. allows the connected GitHub app to bind that already-proven run identity to PR #16 after the dispatcher run, using the connector's own comment capability.

`tools/ci/workflow_policy.py` continues to require the REST endpoint and forbid GraphQL `gh pr comment`. This is deliberately not a relaxation of G11 identity or publication controls: failure to comment may not block dispatch, but G11 is accepted only after the exact canonical release run is identified and its identity is durably bound to PR #16 by either the workflow or connected GitHub app.

Because `ci-fast.yml` and this handoff are source-fingerprinted, the `eafc1d56...` G10 evidence becomes historical after this non-blocking tracking correction. One fresh exact-source G10 is required; no other source changes are planned before G11.

## Final release-dispatch contract

The canonical workflow set remains exactly `ci-fast.yml`, `ci-qualified.yml`, and `release.yml`. No fourth workflow and no G17+ gate exists.

The dispatcher supports two certification entry mechanisms plus one publication mechanism:

1. **Normal release-certification push (primary):** exact `v<release-line>-release-certification` push may dispatch. Release resolution validates the exact branch and forces `publish=false`.
2. **Owner-gated PR certification fallback:** an owner-triggered `pull_request` whose **base** is exact `v<release-line>-release-certification` may dispatch certification using the immutable PR base ref/SHA, not the PR merge SHA or head SHA. It resolves `release_ref` from `pull_request.base.ref`, `candidate_sha` from `pull_request.base.sha`, and explicitly prohibits publication.
3. **Stable promotion:** exact `v<release-line>-stable-promotion` remains push-only and owner-gated via `github.actor` or push `sender.login`, and the resolver maps it to `publish=true`. There is deliberately **no pull-request fallback for Stable promotion**.

The PR fallback cannot publish Stable, cannot choose a non-release-certification base, and dispatches canonical `release.yml` against the existing release-certification branch/base SHA.

## Canonical G11–G16 contract after fresh G10

After fresh G10 passes, bind only fingerprint-excluded checkpoints and reconcile `v18.6-release-certification` to that qualified state. Then create one owner-controlled fingerprint-excluded trigger PR **targeting** `v18.6-release-certification` without merging it. Its Fast lane must pass and the owner-gated PR fallback must:

- resolve the exact release-certification base ref and base SHA;
- compute the canonical source fingerprint from that immutable base SHA;
- force `publish=false`;
- dispatch or reuse a canonical `release.yml` run only when both `headBranch` and `headSha` match that exact release-certification candidate;
- retain the canonical release run ID/URL in the dispatcher job summary/log;
- attempt REST tracking comments without treating integration-level comment permission failure as a release failure;
- if the workflow token cannot comment, bind the dispatcher/run identity to PR #16 through the connected GitHub app before accepting G11.

`release.yml` remains the sole G11–G16 owner. It performs exact-SHA G11 provenance, G12 full certification via `release/v18.6.0/run_full_certification.sh`, macOS Apple Silicon + Windows x64 G13/G14 native package/runtime audit, G15 evidence binding, and G16 durable adaptive handoff. Pre-merge certification remains `publish=false`.

Stable publication is allowed only later through the exact Stable-promotion push path and must remain no-rebuild, owner-gated, and evidence-bound.

## Release sequence still required

Fresh G10 on this final non-blocking tracking source → bind excluded checkpoints → reconcile release-certification → one owner-gated trigger PR → canonical release run ID → bind run identity to PR #16 if workflow token cannot → G11 immutable candidate/provenance → G12 full certification → G13 native packaging/provenance → G14 actual packaged runtime audit on macOS Apple Silicon and Windows x64 → G15 release assurance → G16 adaptive release handoff → only then no-rebuild Stable promotion.

No v18.6 Stable tag, package, native artifact hash or publication claim is valid before those gates pass.

## Known residuals / User Action Required

- TradeInsight remains governed future/shadow implementation work until separately validated and approved; it has no automatic production influence merely because other providers exist.
- Deployment-specific provider keys, entitlements and commercial/redistribution/AI-use rights remain User Action Required where not already configured or approved and fail closed when required.
- v18.6 native macOS Apple Silicon and Windows x64 evidence remains pending until G13/G14 passes on the final G11 candidate.

## Exactly one next action

**Obtain one fresh canonical CI Fast + CI Qualified PASS on this final non-blocking tracking source; then bind only fingerprint-excluded checkpoints, reconcile `v18.6-release-certification`, open one owner-controlled fingerprint-excluded trigger PR, capture the canonical `release.yml` run ID from workflow output, bind that identity to PR #16 through the connected GitHub app if necessary, and continue G11–G16 with `publish=false`.**

## Provider-neutral continuation instruction

> Connect to `depulseapp/DE-PULSE`, read `AGENTS.md` or `CLAUDE.md`, inspect `release_identity.json`, `handoff/CURRENT.md`, both `.depulse-certification/resume/` checkpoints, PR #16 checks/comments and release branch state. Treat `v18.5.2-stable` as the immutable certified predecessor until v18.6.0 G11–G16 promotion is complete. Resume from the last trustworthy exact-source PASS and the single next action above. Preserve G0–G16, assistant/account independence, Smart Provider Router ownership, deterministic Day/Swing/Long formulas, the U.S. Equities Processing Boundary and permanent No Execution Boundary.
