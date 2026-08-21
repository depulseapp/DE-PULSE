# DE.PULSE — Current Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS**

**GitHub source of truth:** `depulseapp/DE-PULSE`  
**Certified Stable:** `v18.8.1-stable`  
**Stable candidate:** `410679ba0d6459f66a44db15a0a55f30741a7c53`  
**Certified fingerprint:** `bfefa3605ab29b4678275936a3e60e45133d0b592b91298551731f6d629a9d92`  
**Certified Stable build ID:** `v18.8.1-stable-20260820`  
**Engineering branch:** `v18.8.2-development`  
**Candidate package identity:** `18.8.2` / `v18.8.2-stable-20260820`  
**Current issue:** #57 / `ADAPT-FRESHNESS-001 REOPENED`

The v18.8.1 build/release-evidence checkpoints intentionally remain anchored to the immutable certified Stable while v18.8.2 is an in-flight candidate/recovery line. Candidate identity, branch and qualification progress live in GitHub/current handoff; prior Stable PASS evidence is never rewritten merely to match candidate version strings.

## Product scope and qualified implementation

v18.8.2 remains the bounded Market Intelligence reliability repair for issue #57 only. The implementation makes the existing canonical breadth universe participate in quote freshness/recovery accountability and renders genuinely unavailable/degraded Market Tradeability evidence as `UNAVAILABLE` instead of a misleading numeric zero. Smart Provider Router v2 remains sole routing authority; canonical freshness/recovery and routed refresh remain sole recovery owners; existing multi-feed allocation remains sole subscription owner; deterministic Day/Swing/Long and No Execution are preserved.

The release-capable source head `186dd18bcd33a2d891b3df738478ba88cf7b98b6` passed Fast #435 / `32434635563` and full Qualified #150 / `32434742951`. PR #59 then squash-merged to `main` as `d855607426bc56372656d3b0baad67611aae7a96` with G11 source-head → merged-candidate fingerprint equivalence PASS.

## Release #30 failure classification

Automatic Release #30 / run `32435511692` passed G11 and failed immediately in G12 before native packaging. `release_identity.py --verify`, workflow policy, Stable evidence, pre-merge rehearsal and candidate provenance all passed. The sole failure was `version_consistency_test.py` requiring the obsolete literal README heading `## Resume with any AI assistant or account`.

The current version-neutral README already carries the durable portability semantics through `AGENTS.md`, `CLAUDE.md`, `governance/AI-ASSISTANT-PORTABILITY-CONTRACT.md`, `handoff/CURRENT.md`, and the explicit statement that GitHub is the durable source of truth. Therefore Release #30 is classified as **CI/release-harness failure**, not product/runtime failure.

## Approved recovery path

Repository precedent is v18.7.0 recovery PR #53. Following that contract, branch hygiene deleted the merged development branch and the same canonical `v18.8.2-development` line has been recreated from failed merged candidate `d855607426bc56372656d3b0baad67611aae7a96`.

Recovery scope is release-harness/governance only:
- `version_consistency_test.py` validates durable README portability semantics instead of one presentation heading;
- `release/v18.8.2/release_contract.json` records Release #30 and the fail-closed recovery rule;
- `release_identity.json` carries a non-runtime recovery note so the existing Release workflow path filter re-enters the same canonical G11–G16 workflow after a fresh exact-head merge;
- this handoff records portable continuation state.

No Market Intelligence runtime code, provider routing, freshness behavior, deterministic scoring, execution boundary, application build ID or package identity is changed by the recovery.

## Exactly one next action

Ensure exactly one Draft recovery PR exists from `v18.8.2-development` to `main`. If it does not yet exist, create it; if it exists, inspect its single automatic Fast run. Only after exact-head Fast passes may that same PR be marked Ready once for fresh full Qualified. Do not rerun Release #30 unchanged, do not create a retry/certification/promotion branch, and do not merge the recovery PR until the new exact head has both `DE.PULSE/fast-head` and `DE.PULSE/qualified-head` success.

After the recovery PR merges, allow only the automatic canonical G11–G16 Release run. Stable may be claimed only after G12, macOS Apple Silicon + Windows x64 G13/G14 actual packaged-runtime audits, G15 assurance, same-run no-rebuild publication and G16 complete.

After v18.8.2 Stable and continuity reconciliation, resume v18.9.0 `ADAPT-TRADEINSIGHT-001` SHADOW integration through Smart Provider Router v2.

## Resume rule

Any ChatGPT account, Claude or other assistant must read `AGENTS.md` / `CLAUDE.md`, `governance/AI-ASSISTANT-PORTABILITY-CONTRACT.md`, this file, CURRENT Adaptive overlays, `release_identity.json`, `release/v18.8.2/release_contract.json`, both `.depulse-certification/resume/` checkpoints, the v18.8.1 Stable evidence manifest, issue #57, merged PR #59, the single live recovery PR (if opened), and live workflow/branch state before changing source. Never resume from model memory alone.
