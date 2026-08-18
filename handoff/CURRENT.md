# DE.PULSE — Current Authoritative Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS.** GitHub is authoritative; chat history and temporary workspaces are advisory only.

**Release:** `v18.6.0`  
**Active branch:** `v18.6-development`  
**Stable predecessor:** v18.5.2 STABLE / G0–G16 CLOSED  
**Current candidate state:** v18.6.0 G10 PRE-FREEZE / NOT PROMOTED  
**Repository:** `depulseapp/DE-PULSE`  
**Stable predecessor tag:** `v18.5.2-stable`  
**Stable predecessor promotion commit:** `d30e54db4908ca57c52ae298cc4ada3416fab46b`  
**Stable predecessor certified source fingerprint:** `807de082d43e83d1d3548bca9350d13b72ef4dc71a848940b73b63e4b4d215b0`  
**v18.6 runtime build ID:** `v18.6.0-stable-20260818`  
**Last updated:** 2026-08-18 America/Vancouver

## Resume rule

Read `AGENTS.md` or `CLAUDE.md`, then reconcile this handoff, `release_identity.json`, both `.depulse-certification/resume/` checkpoints, actual branch/PR/check state and the immutable Stable predecessor tag/release. The last trustworthy PASS is the highest exact-source gate evidence visible in GitHub. Never resume from model memory alone. `v18.5.2-stable` remains the current certified product authority until v18.6.0 completes G11–G16 and is promoted.

`source_fingerprint.py` excludes `.depulse-certification` but **does not exclude `handoff/CURRENT.md`**. Therefore the machine resume/evidence checkpoints may advance as fingerprint-excluded operational metadata after a product candidate, while this handoff is candidate-bound source and must be frozen with the product before G11. Changing this handoff after G11 requires a new candidate/fingerprint and requalification.

## v18.6.0 implementation state

All eight assigned v18.6 implementation/audit slices are code-complete:

- watchlist membership/add remediation with regression coverage;
- shared Scanner / Opportunity Radar broad-snapshot acquisition with bounded cache, freshness reuse, partial-miss fetch and coalescing;
- one serialized Session Intelligence Coordinator for Pre-Market Prep and Market Open Prep;
- Market Activity demotion to supporting drill-down plus legacy evidence-route consolidation;
- role-aware documentation with server-authoritative privileged access and Documentation Impact evidence;
- external dependency/provider readiness with durable User Action Required tracking and fail-closed entitlement/rights governance;
- bounded AI context, full cache identity + TTL, strict structured-output/citation validation, safe abstention and continuous-evaluation telemetry;
- provider×dataset rights-aware AI egress with unknown/denied evidence withheld before external model calls.

Protected deterministic Day/Swing/Long formulas remain unchanged. Smart Provider Router remains the sole provider-routing owner; provider count does not change Market Mode. GLD/SLV/USO tradable exceptions, desktop SQLite / hosted PostgreSQL architecture, the U.S. Equities Processing Boundary and the permanent No Execution Boundary are preserved.

## Qualification evidence

The hardened implementation head `84efadd58b4b279e9055a1b5db4a78f3b41b693c` passed:

- CI Fast run `32182124289` — PASS
- CI Qualified run `32182124317` — PASS, including Ubuntu/Windows/macOS portability, full Go suite, race detector, randomized package order, browser behavior, renderer contracts and evidence summary.

The canonical release identity was then generated in isolation. Candidate identity values are:

- version `18.6.0`;
- build `v18.6.0-stable-20260818`;
- previous Stable `v18.5.2`;
- major v18 provenance anchor `v17.5.1`;
- runtime/config continuity `PersonalMarketTerminal`.

Pre-freeze qualification correctly surfaced and closed stale resume identity, release-coupled asset cache versions, QA-history identity, and version-hardcoded regression assumptions. Those failures were not bypassed: the canonical identity tooling and regression contracts were hardened so future release freezes derive from `release_identity.json`.

The canonical Fast workflow now also contains an **owner-only release dispatcher inside the existing G0–G16 workflow model**. An exact `v18.6-release-certification` branch runs Fast and then dispatches the existing `release.yml` G11–G16 workflow with `publish=false`. After pre-promotion G11–G16 closure and merge, an exact `v18.6-stable-promotion` branch can run the same full release workflow with `publish=true`, preserving the existing no-rebuild publication contract. This adds no fourth workflow and no G17+ gate; `ci-fast.yml`, `ci-qualified.yml`, and `release.yml` remain the canonical workflow set.

The remaining G10 requirement is a fresh Fast + Qualified PASS on this final candidate-bound handoff and exact source, including the release-dispatch governance change.

## Release sequence still required

G10 pre-freeze exact-identity qualification → G11 immutable RC/provenance → G12 full certification via `release/v18.6.0/run_full_certification.sh` → G13 native packaging/provenance → G14 actual packaged runtime audit on macOS Apple Silicon and Windows x64 → G15 release assurance → G16 adaptive release handoff. The pre-merge certification branch uses `publish=false`; Stable publication occurs only from the post-merge `v18.6-stable-promotion` path with `publish=true` and already-certified artifacts within that run.

No v18.6 Stable tag, package, native artifact hash or publication claim is valid before those gates pass.

## Known residuals / User Action Required

- TradeInsight remains governed future/shadow implementation work until separately validated and approved; it has no automatic production influence merely because other providers exist.
- Deployment-specific provider keys, entitlements and commercial/redistribution/AI-use rights remain User Action Required where not already configured or approved and fail closed when required.
- v18.6 native macOS Apple Silicon and Windows x64 evidence remains pending until G13/G14.

## Exactly one next action

**Reconcile the fingerprint-excluded `.depulse-certification/resume` checkpoints to this final candidate-bound source, obtain fresh canonical CI Fast + CI Qualified PASS, fast-forward the qualified result to `v18.6-development` / PR #16, then create `v18.6-release-certification` at that exact candidate to execute G11–G16 with publication disabled.**

## Provider-neutral continuation instruction

> Connect to `depulseapp/DE-PULSE`, read `AGENTS.md` or `CLAUDE.md`, inspect `release_identity.json`, `handoff/CURRENT.md`, both `.depulse-certification/resume/` checkpoints, PR #16 and exact GitHub checks. Treat `v18.5.2-stable` as the immutable certified predecessor until v18.6.0 G11–G16 promotion is complete. Resume from the last trustworthy PASS and the single next action above. Preserve G0–G16, assistant/account independence, Smart Provider Router ownership, deterministic Day/Swing/Long formulas, the U.S. Equities Processing Boundary and permanent No Execution Boundary.
