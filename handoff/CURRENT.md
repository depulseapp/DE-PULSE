# DE.PULSE — Current Authoritative Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS.** GitHub is authoritative; chat history and temporary workspaces are advisory only.

**Release:** `v18.6.0`  
**Active branch:** `v18.6-development`  
**Stable predecessor:** v18.5.2 STABLE / G0–G16 CLOSED  
**Current candidate state:** v18.6.0 G10 REQUALIFICATION AFTER RELEASE-OBSERVABILITY HARDENING / NOT PROMOTED  
**Repository:** `depulseapp/DE-PULSE`  
**Main release PR:** `#16`  
**Stable predecessor tag:** `v18.5.2-stable`  
**Stable predecessor promotion commit:** `d30e54db4908ca57c52ae298cc4ada3416fab46b`  
**Stable predecessor certified source fingerprint:** `807de082d43e83d1d3548bca9350d13b72ef4dc71a848940b73b63e4b4d215b0`  
**v18.6 runtime build ID:** `v18.6.0-stable-20260818`  
**Last updated:** 2026-08-18 America/Vancouver

## Resume rule

Read `AGENTS.md` or `CLAUDE.md`, then reconcile this handoff, `release_identity.json`, both `.depulse-certification/resume/` checkpoints, actual branch/PR/check state and the immutable Stable predecessor tag/release. The last trustworthy PASS is the highest exact-source gate evidence visible in GitHub. Never resume from model memory alone. `v18.5.2-stable` remains the current certified product authority until v18.6.0 completes G11–G16 and is promoted.

`source_fingerprint.py` excludes `.depulse-certification` but **does not exclude `handoff/CURRENT.md` or workflow definitions**. Machine resume/evidence checkpoints may advance as fingerprint-excluded operational metadata after a candidate commit, while this handoff and `.github/workflows/*.yml` are candidate-bound source. Any source-bound change after qualification requires a new fingerprint and fresh G10 before G11.

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

The hardened implementation head `84efadd58b4b279e9055a1b5db4a78f3b41b693c` passed CI Fast `32182124289` and CI Qualified `32182124317`.

The canonical v18.6 release candidate source `aa72c8b2c7571d5633eb2e6dce622749704b1637` then passed real-branch G10:

- CI Fast #142 / run `32185509389` — PASS;
- CI Qualified #63 / run `32185509377` — PASS, including Ubuntu/macOS/Windows portability, workflow/provenance contracts, browser behavior, renderer contracts, full Go suite, race detector, randomized package order and evidence summary.

The canonical product fingerprint observed after that qualification was `f7ae32d088563019f406abc70f4781a85b061fdd237ccabd9d4539ad8d1a3622`. Subsequent `.depulse-certification` checkpoint-only commits did not alter that fingerprint.

## Release delivery hardening after G10

The first branch-gated certification attempt exposed an assistant/account portability gap: connected assistants could trigger `workflow_dispatch` but could not reliably enumerate that run later. The release machinery is therefore being hardened before G11 is frozen; no product behavior or protected trading logic changed.

The canonical workflow set remains exactly `ci-fast.yml`, `ci-qualified.yml`, and `release.yml`. No new top-level gate or workflow was introduced.

The permanent release path now requires:

- owner-only push dispatch from exact `v<release-line>-release-certification` or `v<release-line>-stable-promotion` branches;
- release workflow definition loaded from the candidate branch itself rather than implicitly from `main`;
- idempotent reuse only when both release-run `headBranch` and `headSha` match the current candidate;
- explicit capture of the release workflow run ID/URL;
- durable tracking on the open v18.6 release PR;
- G11 PASS and G16 PASS status comments from `release.yml` itself;
- `publish=false` during pre-merge certification and `publish=true` only on the Stable-promotion path;
- no rebuild during publication.

Because `ci-fast.yml`, `release.yml`, and this handoff are source-fingerprinted, this delivery hardening intentionally invalidates reuse of the earlier G10 evidence for the new source candidate. Fresh Fast + Qualified qualification is required before the observable G11 run becomes authoritative.

## Release sequence still required

Fresh G10 on the final observable/idempotent release source → G11 immutable candidate/provenance → G12 full certification via `release/v18.6.0/run_full_certification.sh` → G13 native packaging/provenance → G14 actual packaged runtime audit on macOS Apple Silicon and Windows x64 → G15 release assurance → G16 adaptive release handoff. Pre-merge certification uses `publish=false`; Stable publication remains prohibited until the certified candidate is promoted through the no-rebuild Stable-promotion path.

No v18.6 Stable tag, package, native artifact hash or publication claim is valid before those gates pass.

## Known residuals / User Action Required

- TradeInsight remains governed future/shadow implementation work until separately validated and approved; it has no automatic production influence merely because other providers exist.
- Deployment-specific provider keys, entitlements and commercial/redistribution/AI-use rights remain User Action Required where not already configured or approved and fail closed when required.
- v18.6 native macOS Apple Silicon and Windows x64 evidence remains pending until G13/G14 passes on the final G11 candidate.

## Exactly one next action

**Obtain fresh canonical CI Fast + CI Qualified PASS on the final source containing observable/idempotent release tracking, then update only fingerprint-excluded checkpoints, fast-forward `v18.6-release-certification` to that exact qualified metadata head, and follow the PR-recorded canonical release run through G11–G16 with `publish=false`.**

## Provider-neutral continuation instruction

> Connect to `depulseapp/DE-PULSE`, read `AGENTS.md` or `CLAUDE.md`, inspect `release_identity.json`, `handoff/CURRENT.md`, both `.depulse-certification/resume/` checkpoints, PR #16 and exact GitHub checks/comments. Treat `v18.5.2-stable` as the immutable certified predecessor until v18.6.0 G11–G16 promotion is complete. Resume from the last trustworthy PASS and the single next action above. Preserve G0–G16, assistant/account independence, Smart Provider Router ownership, deterministic Day/Swing/Long formulas, the U.S. Equities Processing Boundary and permanent No Execution Boundary.
