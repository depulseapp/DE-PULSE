# DE.PULSE — Current Authoritative Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS.** GitHub is authoritative; chat history and temporary workspaces are advisory only.

**Release:** `v18.6.0`  
**Active branch:** `v18.6-development`  
**Stable predecessor:** v18.5.2 STABLE / G0–G16 CLOSED  
**Current candidate state:** v18.6.0 CONNECTOR-SAFE CERTIFICATION FALLBACK / FRESH G10 REQUIRED / NOT PROMOTED  
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

## Last exact-source G10 PASS and why it is now historical

Source `8b516d8384f21b865d8e7ea318c1c2c48e08e44c` passed the final dispatcher-authorization qualification before the connector event behavior was fully understood:

- CI Fast #198 / run `32193789742` / attempt 1 — PASS;
- CI Qualified #85 / run `32193789745` / attempt 1 — PASS;
- coverage included the dispatcher authorization/publish boundary, Ubuntu/macOS/Windows portability, workflow/provenance contracts, browser behavior, renderer contracts, full Go suite, race detector, randomized package order and final evidence summary.

After that PASS, release-certification was reconciled and multiple fingerprint-excluded commits/merges were attempted to generate a normal branch `push` event. Connector-originated ref moves, PR merges and direct contents commits did not produce the expected dispatcher comments or discoverable push workflow run. PR-created/synchronized events, however, reliably execute Actions in this connected environment.

Because the workflow source is now being extended with a connector-safe certification fallback, the `8b516d...` G10 evidence is historical and cannot be reused for the new source.

## Final release-dispatch contract

The canonical workflow set remains exactly `ci-fast.yml`, `ci-qualified.yml`, and `release.yml`. No fourth workflow and no G17+ gate exists.

The dispatcher now supports two certification entry mechanisms plus one publication mechanism:

1. **Normal release-certification push (primary):** exact `v<release-line>-release-certification` push may dispatch. Release resolution validates the exact branch and forces `publish=false`.
2. **Owner-gated PR certification fallback:** an owner-triggered `pull_request` whose **base** is exact `v<release-line>-release-certification` may dispatch certification using the immutable PR base ref/SHA, not the PR merge SHA or head SHA. It resolves `release_ref` from `pull_request.base.ref`, `candidate_sha` from `pull_request.base.sha`, and explicitly prohibits publication. This path exists because connected-app writes may suppress normal push-triggered Actions while pull-request Actions execute reliably.
3. **Stable promotion:** exact `v<release-line>-stable-promotion` remains push-only and owner-gated via `github.actor` or push `sender.login`, and the resolver maps it to `publish=true`. There is deliberately **no pull-request fallback for Stable promotion**.

`tools/ci/workflow_policy.py` enforces all of these invariants and rejects both the unreliable `head_commit.author.username` authorization pattern and any pull-request fallback targeting `*-stable-promotion`.

The PR fallback changes only orchestration. It cannot publish Stable, cannot choose a non-release-certification base, and dispatches canonical `release.yml` against the existing release-certification branch/base SHA.

## Canonical G11–G16 contract after fresh G10

After fresh G10 passes, bind only fingerprint-excluded checkpoints and reconcile `v18.6-release-certification` to that qualified state. Then create an owner-controlled fingerprint-excluded trigger PR **targeting** `v18.6-release-certification` without merging it first. Its Fast lane must pass and the owner-gated PR fallback must:

- resolve the exact release-certification base ref and base SHA;
- compute the canonical source fingerprint from that immutable base SHA;
- force `publish=false`;
- resolve release PR #16;
- post a dispatcher-active comment to PR #16 with dispatcher run ID, release ref, exact candidate/base SHA, fingerprint and `publish=false`;
- reuse a canonical release run only when both `headBranch` and `headSha` equal that release-certification base candidate, otherwise dispatch a new `release.yml` run;
- post the canonical `release.yml` run ID/URL to PR #16.

Only after that durable run ID appears is G11 considered started.

`release.yml` remains the sole G11–G16 owner. It performs exact-SHA G11 provenance, G12 full certification via `release/v18.6.0/run_full_certification.sh`, macOS Apple Silicon + Windows x64 G13/G14 native package/runtime audit, G15 evidence binding, and G16 durable adaptive handoff. Pre-merge certification remains `publish=false`.

Stable publication is allowed only later through the exact Stable-promotion push path and must remain no-rebuild, owner-gated, and evidence-bound.

## Release sequence still required

Fresh G10 on this connector-safe dispatcher source → bind excluded checkpoints → reconcile release-certification → owner-gated trigger PR into release-certification → require PR #16 canonical release run ID → G11 immutable candidate/provenance → G12 full certification → G13 native packaging/provenance → G14 actual packaged runtime audit on macOS Apple Silicon and Windows x64 → G15 release assurance → G16 adaptive release handoff → only then no-rebuild Stable promotion.

No v18.6 Stable tag, package, native artifact hash or publication claim is valid before those gates pass.

## Known residuals / User Action Required

- TradeInsight remains governed future/shadow implementation work until separately validated and approved; it has no automatic production influence merely because other providers exist.
- Deployment-specific provider keys, entitlements and commercial/redistribution/AI-use rights remain User Action Required where not already configured or approved and fail closed when required.
- v18.6 native macOS Apple Silicon and Windows x64 evidence remains pending until G13/G14 passes on the final G11 candidate.

## Exactly one next action

**Obtain fresh canonical CI Fast + CI Qualified PASS on the connector-safe PR-certification fallback source; then bind only fingerprint-excluded checkpoints, reconcile `v18.6-release-certification`, open an owner-controlled fingerprint-excluded trigger PR into that branch, and require the PR #16 dispatcher plus canonical `release.yml` run-ID comments before accepting G11.**

## Provider-neutral continuation instruction

> Connect to `depulseapp/DE-PULSE`, read `AGENTS.md` or `CLAUDE.md`, inspect `release_identity.json`, `handoff/CURRENT.md`, both `.depulse-certification/resume/` checkpoints, PR #16 checks/comments and release branch state. Treat `v18.5.2-stable` as the immutable certified predecessor until v18.6.0 G11–G16 promotion is complete. Resume from the last trustworthy exact-source PASS and the single next action above. Preserve G0–G16, assistant/account independence, Smart Provider Router ownership, deterministic Day/Swing/Long formulas, the U.S. Equities Processing Boundary and permanent No Execution Boundary.
