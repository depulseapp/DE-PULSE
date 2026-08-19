# DE.PULSE — Current Authoritative Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS.** GitHub is authoritative; chat history and temporary workspaces are advisory only.

**Release:** `v18.6.0`  
**Active branch:** `v18.6-development`  
**Stable predecessor:** `v18.5.2-stable` / G0–G16 CLOSED  
**Current candidate state:** v18.6.0 PROMOTION-ONLY HARDENING / FRESH G10 REQUIRED / NOT PROMOTED  
**Repository:** `depulseapp/DE-PULSE`  
**Main release PR:** `#16`  
**Runtime build ID:** `v18.6.0-stable-20260818`  
**Last updated:** 2026-08-18 America/Vancouver

## Resume rule

Read `AGENTS.md` or `CLAUDE.md`, then reconcile this handoff, `release_identity.json`, both `.depulse-certification/resume/` checkpoints, PR #16, current checks, the release-certification branch, and immutable Stable predecessor. Never resume from model memory alone.

`source_fingerprint.py` excludes `.depulse-certification` only. Workflow definitions, release tooling, certification harnesses, browser proofs/contracts, and this handoff are source-fingerprinted. Any change to them requires fresh G10 before G11.

## Product / architecture state

All eight v18.6 implementation/audit slices remain code-complete: watchlist remediation; shared Scanner/Radar broad-snapshot acquisition; serialized Session Intelligence Coordinator; Market Activity/legacy-route consolidation; role-aware documentation; external dependency/provider readiness; bounded AI context/cache/schema/evaluation hardening; and provider×dataset rights-aware fail-closed AI egress.

Protected invariants remain unchanged: deterministic Day/Swing/Long formulas; Smart Provider Router sole routing ownership; provider count never changes Market Mode; GLD/SLV/USO tradable live exceptions; desktop SQLite / hosted PostgreSQL; U.S. Equities Processing Boundary; permanent No Execution Boundary; exactly three canonical workflows; G0–G16 only.

## Last trustworthy pre-change certification

Source `f706b205fcb36bc74ca31e113be6a6add3c2afdb`, qualified by CI Fast #250 / run `32198719049` PASS and CI Qualified #105 / run `32198719053` PASS, was reconciled to release-certification candidate `10a89f2b94629aa83d588f5a893c0dd7e83334d6`.

Canonical non-publishing release run `32199369265` completed successfully:
- G11 immutable candidate/provenance — PASS;
- G12 authoritative full certification — PASS;
- G13/G14 macOS Apple Silicon native package + actual packaged runtime audit — PASS;
- G13/G14 Windows x64 native package + actual packaged runtime audit — PASS;
- G15 Release Assurance — PASS;
- G16 adaptive handoff evidence — PASS;
- publication — SKIPPED as required for `publish=false`.

That run produced exact certified native binaries and G15 evidence bound to source fingerprint `e1460c36e60ade65a91fccb5ff0c24c769f5b7b95d3ab8829c4402f74eaf6b27`, build ID `v18.6.0-stable-20260818`, and No Execution Boundary `PRESERVED`.

The artifacts were independently inspected before promotion hardening:
- macOS binary SHA-256 `4e6607610965589aceea826c4d642480ad316906db926dfc9f1b6a58b5bd3bd0`;
- Windows binary SHA-256 `900d8255bdded40dd7f8433fb993ddfbb4c93fb42ebcd0bb407eaffd29f90332`;
- both native evidence graphs and `G15-Release-Assurance.json` agree on release, build ID, certified source, fingerprint, PASS state, and promotion authorization.

## Promotion-path audit finding

The prior `release.yml` publish path was named no-rebuild, but a separate `publish=true` workflow run would still execute G11/G12/G13/G14/G15 before its publish job. That would rebuild native artifacts instead of publishing the exact artifacts that had already passed G13/G14/G15.

This is a release-orchestration contract defect, not a product defect. Stable publication has therefore remained blocked and `v18.5.2-stable` remains the immutable Stable predecessor.

## Promotion-only hardening in this source change

The canonical workflow set remains exactly `ci-fast.yml`, `ci-qualified.yml`, and `release.yml`.

The hardening changes are:
1. `release.yml` now separates certification from publication. `publish=false` executes G11→G16 certification. `publish=true` executes a promotion-only path and does **not** run G12/G13/G14/G15 builds again.
2. Promotion-only requires a successful canonical non-publishing `release.yml` run ID, downloads the macOS, Windows, and G15 artifacts from that exact run using GitHub Actions cross-run artifact retrieval, and verifies the certification run identity before publication.
3. `tools/release/verify_promotion_evidence.py` validates release, build ID, certified source SHA, source fingerprint, G15 promotion authorization, No Execution Boundary, every native PASS check, and exact SHA-256 of both native binary ZIPs.
4. Stable promotion target must preserve the certified canonical source fingerprint. The Stable tag/release uploads the exact downloaded certified assets; no native build command runs in promotion-only mode.
5. `ci-fast.yml` supports owner-gated PR fallback for both release-certification and stable-promotion because connector-originated pushes may not emit Actions runs. PR fallback is additionally restricted to fingerprint-excluded `.depulse-certification/resume/` changes only and uses the immutable PR base SHA as the target.
6. Stable promotion reads the canonical certification run ID from the release evidence checkpoint. Workflow-token PR comments remain best-effort observability; Actions run identity + checkpoint/PR binding remain authoritative.
7. `tools/ci/workflow_policy.py` rejects regression to same-run rebuild publication and requires the exact-artifact cross-run promotion contract.

Because `release.yml`, `ci-fast.yml`, workflow policy, release verification tooling, and this handoff are source-fingerprinted, the successful historical run `32199369265` cannot certify this new promotion orchestration. It is retained as historical proof and as validation input for the verifier, but the final source requires fresh G10 and one new canonical non-publishing G11–G16 certification before Stable promotion.

## Exactly one next action

**Qualify the consolidated promotion-only hardening as one fresh Fast + Qualified G10 candidate. If both pass, bind only fingerprint-excluded checkpoints, reconcile `v18.6-release-certification`, run one canonical `publish=false` G11–G16 certification, bind its exact run ID and native artifact identities into the checkpoint, merge PR #16 only after that run is fully green, create the exact-fingerprint `v18.6-stable-promotion` target, and invoke the owner-gated metadata-only promotion fallback so `publish=true` reuses that certified run’s exact artifacts without rebuilding.**

## Known residuals / User Action Required

- TradeInsight remains governed future/shadow work until separately validated and approved; it gains no production influence automatically.
- Deployment-specific provider keys, entitlements and commercial/redistribution/AI-use rights remain User Action Required where absent or unapproved and fail closed.
- No v18.6 Stable claim is valid until the final source passes fresh G10, new G11–G16 certification succeeds, exact-artifact promotion succeeds, and `v18.6.0-stable` is verified published.

## Provider-neutral continuation instruction

> Connect to `depulseapp/DE-PULSE`, read `AGENTS.md` or `CLAUDE.md`, then inspect `release_identity.json`, `handoff/CURRENT.md`, both `.depulse-certification/resume/` checkpoints, PR #16, active checks, and the release-certification branch. Treat `v18.5.2-stable` as immutable Stable until v18.6 exact-artifact no-rebuild promotion completes. Resume from the last trustworthy PASS and source fingerprint, not chat memory. Preserve G0–G16, GitHub source-of-truth hierarchy, assistant/account independence, Smart Provider Router ownership, deterministic desk formulas, U.S. Equities Processing Boundary, and permanent No Execution Boundary.
