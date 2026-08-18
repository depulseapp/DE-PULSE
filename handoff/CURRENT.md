# DE.PULSE — Current Authoritative Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS.** GitHub is authoritative; chat history and temporary workspaces are advisory only.

**Release:** `v18.6.0`  
**Active branch:** `v18.6-development`  
**Stable predecessor:** v18.5.2 STABLE / G0–G16 CLOSED  
**Current candidate state:** v18.6.0 G10 REQUALIFICATION AFTER FINAL RELEASE-RUN OBSERVABILITY HARDENING / NOT PROMOTED  
**Repository:** `depulseapp/DE-PULSE`  
**Main release PR:** `#16`  
**Stable predecessor tag:** `v18.5.2-stable`  
**Stable predecessor promotion commit:** `d30e54db4908ca57c52ae298cc4ada3416fab46b`  
**Stable predecessor certified source fingerprint:** `807de082d43e83d1d3548bca9350d13b72ef4dc71a848940b73b63e4b4d215b0`  
**v18.6 runtime build ID:** `v18.6.0-stable-20260818`  
**Last updated:** 2026-08-18 America/Vancouver

## Resume rule

Read `AGENTS.md` or `CLAUDE.md`, then reconcile this handoff, `release_identity.json`, both `.depulse-certification/resume/` checkpoints, actual branch/PR/check state and the immutable Stable predecessor tag/release. The last trustworthy PASS is the highest exact-source gate evidence visible in GitHub. Never resume from model memory alone. `v18.5.2-stable` remains the current certified product authority until v18.6.0 completes G11–G16 and is promoted.

`source_fingerprint.py` excludes `.depulse-certification` but **does not exclude `handoff/CURRENT.md` or workflow definitions**. Machine resume/evidence checkpoints may advance as fingerprint-excluded operational metadata after a candidate commit. This handoff and `.github/workflows/*.yml` are candidate-bound source; any change to them requires a new candidate fingerprint and fresh G10 before G11.

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

## Qualification history

The hardened implementation head `84efadd58b4b279e9055a1b5db4a78f3b41b693c` passed CI Fast `32182124289` and CI Qualified `32182124317`.

The canonical v18.6 release-candidate source `aa72c8b2c7571d5633eb2e6dce622749704b1637` passed real-branch G10 via CI Fast #142 / `32185509389` and CI Qualified #63 / `32185509377`.

After release-run observability hardening, source `9558ee44df55dd10f931f7435f4a8615b1b976ba` passed CI Fast #158 / `32187439582` and CI Qualified #69 / `32187439404`, including Ubuntu/macOS/Windows portability, workflow/provenance contracts, browser behavior, renderer contracts, full Go suite, race detector, randomized package order and final evidence summary. The qualification merge-tree fingerprint observed there was `6125035b40fbee5c1571e26a1ff3fdd3677000f962961b34f7bdd7b3161e717d`.

That qualification then exposed one final pre-G11 compatibility concern: a candidate-only `workflow_dispatch` input is unnecessarily fragile while the workflow must also exist on the default branch. The final source therefore removes the candidate-only tracking input and places all run observability in the already-authorized Fast dispatcher. This final source requires one fresh Fast + Qualified G10 before G11.

## Final release-run observability contract

The canonical workflow set remains exactly `ci-fast.yml`, `ci-qualified.yml`, and `release.yml`. No fourth workflow and no G17+ gate exists.

For exact `v<release-line>-release-certification` and `v<release-line>-stable-promotion` pushes, the owner-gated Fast dispatcher now:

- resolves version/build/certification script and canonical candidate fingerprint from the pushed SHA;
- resolves the open release PR from the canonical development branch;
- posts a **dispatcher-active comment first**, including the Fast/dispatcher run ID, ref, exact SHA, fingerprint and publish mode;
- searches existing `release.yml` workflow-dispatch runs and reuses one only when both `headBranch` and `headSha` match the exact candidate;
- otherwise dispatches `release.yml` from the candidate branch itself using the stable/default-compatible input schema;
- captures the canonical release workflow run ID/URL and posts a second PR comment with that durable lookup key.

This means another ChatGPT account, Claude, a human developer, or another authorized assistant can recover the exact G11–G16 run from GitHub PR #16 without needing the previous chat or an opaque Actions list operation.

`release.yml` remains the sole G11–G16 owner. It performs exact-SHA G11 provenance, G12 full certification, macOS Apple Silicon + Windows x64 G13/G14 native package/runtime audit, G15 evidence binding, optional no-rebuild publication, and G16 durable workflow artifact. Pre-merge certification uses `publish=false`; Stable publication is allowed only from the post-merge Stable-promotion path.

## Release sequence still required

Fresh G10 on this final dispatcher/source contract → G11 immutable candidate/provenance → G12 full certification via `release/v18.6.0/run_full_certification.sh` → G13 native packaging/provenance → G14 actual packaged runtime audit on macOS Apple Silicon and Windows x64 → G15 release assurance → G16 adaptive release handoff. Stable publication remains prohibited until the certified candidate reaches the no-rebuild Stable-promotion path.

No v18.6 Stable tag, package, native artifact hash or publication claim is valid before those gates pass.

## Known residuals / User Action Required

- TradeInsight remains governed future/shadow implementation work until separately validated and approved; it has no automatic production influence merely because other providers exist.
- Deployment-specific provider keys, entitlements and commercial/redistribution/AI-use rights remain User Action Required where not already configured or approved and fail closed when required.
- v18.6 native macOS Apple Silicon and Windows x64 evidence remains pending until G13/G14 passes on the final G11 candidate.

## Exactly one next action

**Obtain fresh canonical CI Fast + CI Qualified PASS on the final default-compatible, PR-audited dispatcher source; update only fingerprint-excluded checkpoints; fast-forward `v18.6-release-certification` to that qualified metadata head; then use the PR-posted canonical release run ID to execute and inspect G11–G16 with `publish=false`.**

## Provider-neutral continuation instruction

> Connect to `depulseapp/DE-PULSE`, read `AGENTS.md` or `CLAUDE.md`, inspect `release_identity.json`, `handoff/CURRENT.md`, both `.depulse-certification/resume/` checkpoints, PR #16 checks and PR comments. Treat `v18.5.2-stable` as the immutable certified predecessor until v18.6.0 G11–G16 promotion is complete. Resume from the last trustworthy PASS and the single next action above. Preserve G0–G16, assistant/account independence, Smart Provider Router ownership, deterministic Day/Swing/Long formulas, the U.S. Equities Processing Boundary and permanent No Execution Boundary.
