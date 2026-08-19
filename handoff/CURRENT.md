# DE.PULSE — Current Authoritative Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS.** GitHub is authoritative; chat history and temporary workspaces are advisory only.

**Release:** `v18.6.0`  
**Active branch:** `v18.6-development`  
**Certified Stable:** `v18.6.0-stable`  
**Stable promotion commit:** `2abc4a4a3fbbe623aff57948ec875f45e7ef0a1c`  
**Certified source:** `d375852d846f8c9f0045ac929da1830b85ad629e`  
**Certified source fingerprint:** `e8c009c16eedb448ed5b9731d8dd24026a7ea0b5a2b5c82e26490a2941b7b4c8`  
**Canonical certification run:** `32225064225`  
**Historical v18.6 promotion run:** `32225910416`  
**Next product branch:** `v18.6.1-development`  
**Repository:** `depulseapp/DE-PULSE`  
**Last updated:** 2026-08-19 America/Vancouver

`Active branch` above intentionally preserves the certified v18.6 checkpoint identity required by the permanent resume contract; v18.6 product work is closed and the next product branch is `v18.6.1-development`.

## Resume rule

Read `AGENTS.md` or `CLAUDE.md`, then `governance/CI-EFFICIENCY-CONTRACT.md`, this handoff, `release_identity.json`, both `.depulse-certification/resume/` checkpoints, current GitHub PR/check state and the immutable Stable tag. Never resume from model memory alone. Resume from the last trustworthy PASS recorded in GitHub evidence.

## Current Stable truth

v18.6.0 is fully implemented, certified and published. G0–G16 are closed for that immutable Stable. The certified native macOS Apple Silicon and Windows x64 artifacts and G15/G16 evidence remain authoritative through the Stable release/checkpoints. No v18.6.0 release blocker remains.

The product invariants remain unchanged: U.S. Equities Processing Boundary, permanent No Execution Boundary, deterministic Day/Swing/Long formulas, Smart Provider Router ownership, provider-count/Market-Mode rule, GLD/SLV/USO tradable exceptions, desktop SQLite / hosted PostgreSQL architecture, and governed SHADOW → VALIDATED → APPROVED → PRODUCTION influence.

## CI/process correction after v18.6

The v18.6 release exposed CI event amplification: temporary certification/promotion/retry/dispatch branches and PRs were used to manufacture GitHub events; CI Fast listened to overlapping PR and version-branch events; PR close could rerun validated source; Fast also dispatched Release; and macOS/Windows portability ran more often than necessary.

The permanent replacement is `governance/CI-EFFICIENCY-CONTRACT.md` with exactly three canonical workflows:

- **CI Fast:** PR opened/synchronize/reopened only for validation. It explicitly validates the PR source head and records `DE.PULSE/fast-head` on that exact SHA. The workflow also listens to `main` pushes solely for branch hygiene; Fast validation is skipped on those pushes. No Release dispatcher or per-success impact artifact exists.
- **CI Qualified / G10:** normal automatic trigger is `Ready for review`. It validates the exact candidate source head and records `DE.PULSE/qualified-head`. Product candidates use Linux backend/renderer/browser qualification; macOS/Windows portability is reserved for CI/harness changes.
- **Release G11–G16:** one merged `v*-development → main` PR that changes `release_identity.json` starts one Release workflow. G11 fails closed unless both exact-head G10 statuses are successful, the pull-request source ref still identifies that exact head, and its canonical fingerprint equals the merged candidate. The same run then performs G12, native macOS/Windows G13/G14, G15 assurance, exact same-run no-rebuild publication and G16.
- **Retries:** same SHA infrastructure failures rerun only failed jobs; source/test fixes stay on the same development branch and PR. Never create retry/trigger/fallback branches or PRs.
- **Branch hygiene:** main pushes perform cleanup without rerunning Fast tests. Open-PR branches are preserved; merged PR heads (including squash merges), fully merged branches, already-Stable release-line branches and closed/orphaned release-temp branches are removed conservatively.

Exactly three active workflow files remain permitted: `ci-fast.yml`, `ci-qualified.yml`, and `release.yml`. G0–G16 remains the only top-level gate model.

## Cost/quality rule

Optimize event count, paid runner use and artifact retention before removing quality checks. Linux handles routine Fast/Qualified work; macOS/Windows remain mandatory where they establish real cross-platform CI-harness portability or final native G13/G14 runtime truth. Stable publication is impossible unless exact-head G10 evidence and G11–G15 all pass for the immutable merged candidate.

## Provider-neutral continuation instruction

> Connect to `depulseapp/DE-PULSE`, read `AGENTS.md` or `CLAUDE.md`, `governance/CI-EFFICIENCY-CONTRACT.md`, `handoff/CURRENT.md`, `release_identity.json`, both `.depulse-certification/resume/` checkpoints, current PR/check state and latest Stable release. Treat `v18.6.0-stable` as immutable. Use one development branch + one PR for the next release, never create CI-trigger/retry/certification/promotion branches, and preserve all permanent product boundaries.

## Exactly one next action

After the CI-efficiency hardening PR is qualified and merged, start normal v18.6.1 intake on `v18.6.1-development` using one draft PR to `main`; no additional v18.6.0 release work is required.
