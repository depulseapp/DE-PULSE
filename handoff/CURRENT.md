# DE.PULSE — Current Authoritative Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS.** GitHub is authoritative; chat history and temporary workspaces are advisory only.

**Release:** `v18.6.0`  
**Active branch:** `v18.6-development`  
**Stable predecessor:** v18.5.2 STABLE / G0–G16 CLOSED  
**Current candidate state:** v18.6.0 FINAL DISPATCH AUTHORIZATION FIX / FRESH G10 REQUIRED / NOT PROMOTED  
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

## Latest trustworthy qualification and why it must be repeated

Exact source `7b8e3ab40a55edce8e4b16acbd3781a531f70585` eventually passed its blocked G10 retries after GitHub Actions allocation recovered:

- CI Fast #181 / run `32190500889` / attempt 5 — PASS;
- CI Qualified #78 / run `32190500878` / attempt 4 — PASS;
- coverage included Ubuntu/macOS/Windows portability, workflow/provenance contracts, browser behavior, renderer contracts, full Go suite, race detector, randomized package order and qualified evidence summary.

That evidence is trustworthy for `7b8e3ab...` but is now historical because the release dispatcher workflow was subsequently corrected. It **must not** be reused as G10 evidence for the new source.

## Release-dispatch defect discovered after G10

The release-certification branch was reconciled without force-pushing and a metadata-only GitHub merge was used to produce a normal `web-flow` merge commit on `v18.6-release-certification`. The expected PR #16 dispatcher/run-ID comments still did not appear.

The remaining orchestration defect was the authorization expression in `ci-fast.yml`: it depended on `github.event.head_commit.author.username` as a fallback when `github.actor` was not the repository owner. That field is not a reliable authorization identity for a push event and should not control the release dispatcher.

The corrected contract is deliberately asymmetric:

1. **Release-certification branch:** an exact `v<release-line>-release-certification` push may enter the dispatcher without a user-identity gate. The dispatcher itself validates the release identity and exact branch, and the branch maps to `publish=false`. Therefore this path can certify but cannot publish Stable.
2. **Stable-promotion branch:** publication remains owner-gated. The exact `v<release-line>-stable-promotion` push requires repository-owner authorization via `github.actor` or the push payload `sender.login`, and then maps to `publish=true`.
3. `tools/ci/workflow_policy.py` now enforces this contract, rejects the unreliable `head_commit.author.username` authorization pattern, and requires the exact release-certification/promotion branch-to-publish mapping.
4. The canonical workflow set remains exactly `ci-fast.yml`, `ci-qualified.yml`, and `release.yml`. No fourth workflow and no G17+ gate exists.

This correction changes source-fingerprinted files, so a fresh exact-source Fast + Qualified G10 is mandatory before any new G11 run.

## Canonical G11–G16 contract after fresh G10

After fresh G10 passes, advance `v18.6-release-certification` to the qualified candidate plus fingerprint-excluded evidence metadata. A normal push on that branch must let the Fast dispatcher:

- resolve release version/build/certification script and canonical source fingerprint from the pushed SHA;
- prove the pushed ref matches the release identity;
- resolve the open release PR from `v18.6-development`;
- post a dispatcher-active comment to PR #16 with the Fast/dispatcher run ID, exact SHA, fingerprint and `publish=false`;
- reuse an existing `release.yml` run only when both `headBranch` and `headSha` match the exact candidate, otherwise dispatch a new exact-candidate run;
- post the canonical `release.yml` run ID/URL to PR #16.

`release.yml` remains the sole G11–G16 owner. It performs exact-SHA G11 provenance, G12 full certification via `release/v18.6.0/run_full_certification.sh`, macOS Apple Silicon + Windows x64 G13/G14 native package/runtime audit, G15 evidence binding, and G16 durable adaptive handoff. Pre-merge certification uses `publish=false`.

Stable publication is allowed only later through the exact Stable-promotion path and must remain no-rebuild, owner-gated, and evidence-bound.

## Release sequence still required

Fresh G10 on the final dispatcher source → G11 immutable candidate/provenance → G12 full certification → G13 native packaging/provenance → G14 actual packaged runtime audit on macOS Apple Silicon and Windows x64 → G15 release assurance → G16 adaptive release handoff → only then no-rebuild Stable promotion.

No v18.6 Stable tag, package, native artifact hash or publication claim is valid before those gates pass.

## Known residuals / User Action Required

- TradeInsight remains governed future/shadow implementation work until separately validated and approved; it has no automatic production influence merely because other providers exist.
- Deployment-specific provider keys, entitlements and commercial/redistribution/AI-use rights remain User Action Required where not already configured or approved and fail closed when required.
- v18.6 native macOS Apple Silicon and Windows x64 evidence remains pending until G13/G14 passes on the final G11 candidate.

## Exactly one next action

**Obtain fresh canonical CI Fast + CI Qualified PASS on the final dispatcher authorization source; after both pass, bind only fingerprint-excluded checkpoints, reconcile `v18.6-release-certification` to that qualified state without rewriting source, and require the PR-posted canonical `release.yml` run ID before accepting G11.**

## Provider-neutral continuation instruction

> Connect to `depulseapp/DE-PULSE`, read `AGENTS.md` or `CLAUDE.md`, inspect `release_identity.json`, `handoff/CURRENT.md`, both `.depulse-certification/resume/` checkpoints, PR #16 checks/comments and release branch state. Treat `v18.5.2-stable` as the immutable certified predecessor until v18.6.0 G11–G16 promotion is complete. Resume from the last trustworthy exact-source PASS and the single next action above. Preserve G0–G16, assistant/account independence, Smart Provider Router ownership, deterministic Day/Swing/Long formulas, the U.S. Equities Processing Boundary and permanent No Execution Boundary.
