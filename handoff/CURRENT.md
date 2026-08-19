# DE.PULSE — Current Authoritative Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS.** GitHub is authoritative; chat history and temporary workspaces are advisory only.

**Release:** `v18.6.0`  
**Active branch:** `v18.6-development`  
**Stable predecessor:** `v18.5.2-stable` / G0–G16 CLOSED  
**Current candidate state:** v18.6.0 PROMOTION-PATH HARDENING / FRESH G10 REQUIRED / NOT PROMOTED  
**Repository:** `depulseapp/DE-PULSE`  
**Main release PR:** `#16`  
**v18.6 runtime build ID:** `v18.6.0-stable-20260818`  
**Last updated:** 2026-08-18 America/Vancouver

## Resume rule

Read `AGENTS.md` or `CLAUDE.md`, then reconcile this handoff, `release_identity.json`, both `.depulse-certification/resume/` checkpoints, PR #16, current checks, the release-certification branch, and the immutable Stable predecessor. Never resume from model memory alone.

`source_fingerprint.py` excludes `.depulse-certification` only. Workflow definitions, certification harnesses, browser proofs/contracts, and this handoff are source-fingerprinted. Any change to them requires fresh G10 before G11.

## Product state

All eight v18.6 product slices remain code-complete. This change is release-process-only: it does not change the renderer, backend product behavior, deterministic Day/Swing/Long formulas, Smart Provider Router ownership, provider-count/Market-Mode rule, GLD/SLV/USO tradable exceptions, desktop SQLite / hosted PostgreSQL architecture, U.S. Equities Processing Boundary, or permanent No Execution Boundary.

The canonical workflow set remains exactly `ci-fast.yml`, `ci-qualified.yml`, and `release.yml`. G0–G16 remains the only top-level gate model.

## Last trustworthy PASS before promotion-path source change

Canonical `release.yml` run `32199369265` completed SUCCESS on release-certification candidate `10a89f2b94629aa83d588f5a893c0dd7e83334d6`, source fingerprint `e1460c36e60ade65a91fccb5ff0c24c769f5b7b95d3ab8829c4402f74eaf6b27`:
- G11 immutable candidate/provenance: PASS;
- G12 authoritative full certification: PASS;
- G13/G14 macOS Apple Silicon package + actual packaged runtime audit: PASS;
- G13/G14 Windows x64 package + actual packaged runtime audit: PASS;
- G15 Release Assurance: PASS;
- G16 adaptive handoff: PASS;
- publication: deliberately SKIPPED because `publish=false`.

That run remains valid historical proof for the pre-hardening source. It cannot certify the workflow/handoff changes below because those files are source-fingerprinted.

## Promotion-path defect found before Stable publication

The pre-hardening `release.yml` labelled publication “no-rebuild” but `publish=true` still depended on the same-run G12/G13/G14/G15 jobs, so a promotion run would recertify and rebuild native artifacts before publishing. That violated the stronger DE.PULSE contract that Stable promotion must publish the exact already-certified artifacts.

A second portability problem existed for connected assistants: connector-originated branch writes do not reliably emit usable push-triggered Actions runs. Certification already has an owner-gated PR-event fallback. Stable promotion previously had no connector-safe equivalent.

A third issue was caught during promotion review: the release-certification reconciliation commit is not guaranteed to be an ancestor of the eventual `main` merge commit even when both have the identical canonical source fingerprint. Stable identity therefore must be enforced by exact source fingerprint + release identity + certified evidence graph, not by a false release-branch ancestry requirement.

## Hardened no-rebuild promotion contract

The canonical workflows now define two distinct modes inside the same `release.yml`:

1. **Certification mode (`publish=false`)**
   - G11 → G12 → native macOS/Windows G13/G14 → G15 → G16 execute normally.
   - Native artifacts and G15 assurance are retained by the successful certification run.

2. **Promotion-reuse mode (`publish=true`)**
   - G11 verifies the immutable certified candidate/fingerprint.
   - G12, macOS G13/G14, Windows G13/G14, and G15 assurance jobs are explicitly skipped; they are never rebuilt or rerun.
   - promotion requires `certification_run_id` and `promotion_sha` from the durable release-evidence checkpoint;
   - `actions/download-artifact` downloads the exact macOS, Windows, and G15 artifacts from that successful certification run using `run-id` + `github-token`;
   - promotion verifies the certification workflow run completed successfully and binds the same candidate SHA, source fingerprint, release version and build ID;
   - native evidence JSON and G15 assurance JSON are revalidated;
   - actual macOS/Windows ZIP SHA-256 values are recomputed and must equal their certified evidence and G15 graph;
   - the promotion target must have the same canonical source fingerprint as the certified candidate;
   - only then may `v18.6.0-stable` and its release assets be published.

## Connector-safe Stable promotion event

Stable promotion stays owner-controlled and exact-branch-gated. Normal push to exact `v<release-line>-stable-promotion` remains supported. In addition, a PR event may trigger Stable promotion **only after the PR is actually merged** into exact `v<release-line>-stable-promotion`, only when `pull_request.merged == true`, only for the repository owner, and only using the resulting merge commit as `promotion_sha`. There is no unmerged/open-PR Stable publication fallback.

Release-certification keeps its existing owner-gated unmerged PR fallback with `publish=false` only.

`tools/ci/workflow_policy.py` owns regression prevention for this separation, cross-run artifact reuse, publish-mode job suppression, evidence validation, and the merged-PR-only Stable promotion rule.

## Why fresh qualification is mandatory

`.github/workflows/ci-fast.yml`, `.github/workflows/release.yml`, `tools/ci/workflow_policy.py`, and this handoff are source-fingerprinted. Therefore run `32199369265` cannot be reused as certification for this hardened source. Fresh Fast + Qualified G10 must pass on the exact new source, followed by one fresh `publish=false` G11–G16 certification run. After that, Stable promotion must reuse that new run's exact artifacts without recertification or rebuild.

## Known residuals / User Action Required

- TradeInsight remains governed future/shadow work until separately validated and approved; it gains no production influence automatically.
- Deployment-specific provider keys, entitlements and commercial/redistribution/AI-use rights remain User Action Required where absent or unapproved and fail closed.
- Current certified Stable remains `v18.5.2-stable` until v18.6 promotion-reuse mode succeeds and the tag/release/artifacts are verified.

## Exactly one next action

**Run one fresh canonical Fast + Qualified G10 qualification on the exact promotion-hardened v18.6 source. If both pass, bind fingerprint-excluded checkpoints, execute exactly one fresh `publish=false` G11–G16 certification run, then merge PR #16 and execute evidence-bound promotion-reuse from that successful certification run without rerunning G12/G13/G14/G15.**

## Provider-neutral continuation instruction

> Connect to `depulseapp/DE-PULSE`, read `AGENTS.md` or `CLAUDE.md`, then inspect `release_identity.json`, `handoff/CURRENT.md`, both `.depulse-certification/resume/` checkpoints, PR #16, current checks, and release branches. Treat `v18.5.2-stable` as immutable Stable until v18.6 promotion-reuse completes. Resume from the last trustworthy PASS and current source fingerprint, not chat memory. Preserve GitHub source-of-truth portability across ChatGPT/Codex/Claude, G0–G16, deterministic desk formulas, Smart Provider Router ownership, U.S. Equities Processing Boundary and permanent No Execution Boundary.
