# DE.PULSE Current Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS**

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Certified Stable candidate:** `e55d8d25b15cec2ffb0f5411bc358bc40b359cf9`  
**Certified Stable qualified source:** `d7276c3421dd2b4529ac2a987466be3cffa05678`  
**Certified Stable source fingerprint:** `0062f46dea5690d0b3fcd8a9ed3b1f71ebe1522c7dee2cb218e9d36b9e0076ff`  
**Certified Stable build ID:** `v18.9.1-stable-20260821`  
**Current closure branch:** `adapt-v18-final-closure-10-10-001`  
**Mandatory next release:** `v18.10.0` — **10/10 Future-Proof Final v18 Closure**  
**Candidate identity:** `v18.10.0` / `v18.10.0-stable-20260825` / platform build `181000`  
**Parent:** #113 / `ADAPT-V18-FINAL-CLOSURE-10-10-001`  
**Parent closure ledger:** `governance/work-slices/ADAPT-V18-FINAL-CLOSURE-10-10-001/closure.json`  
**Active child:** #122 / T9 — cross-platform packaged runtime / release / provenance — **IN_PROGRESS**  
**Next child:** #123 / T10 — future-proof zero-gap final certification — **NOT_STARTED / BLOCKED UNTIL T9 VERIFIED**  
**Future hosted umbrella:** #66 / `ADAPT-HOSTED-SYNC-001` — **BLOCKED**.

## Completed closure tracks

- T1 #114 / PR #127 / merge `b434886751fde1804a4906c9cac41dcce4584834`
- T2 #115 / PR #128 / merge `bd26f91964235a8fdb390184684f4fdd216eb1e6`
- T3 #116 / PR #129 / merge `cf0cc9aa877188775d5c36cb805e6c8935d776bb`
- T4 #117 / PR #130 / merge `fadd9247f63779195138f8d79cbaa0f304b16e61`
- T5 #118 / PR #132 / merge `81137e4481cdf4a351e45c599efa0777ef93460a`
- T6 #119 / PR #133 / merge `ca051884c733118c11e13b0b5ef169c810c39714`
- T7 #120 / PR #134 / merge `80dcbf483378dbf8886bf6f3e421b17fce01b7d5`
- T8 #121 / PR #135 / merge `a3a82aee21fe8e6822d11ba678afefdd6361ef97`

T8 final source `e4542890eddf0d6c94e59e991d9d17ae8cc3ca35` passed exact-head Fast #1080/run `32899149506` and identical-head Qualified #216/run `32899321362`. Qualified full Go, race detector and randomized package order passed. `T8_PERFORMANCE_LOAD_SOAK_CONCURRENCY_RESOURCE_ASSURANCE.json` is `COMPLETE` with zero gaps.

## T9 current authority

T9/#122 is the only active closure track. It certifies the actual distributable v18.10.0 candidate rather than source-only behavior. T9 does **not** publish Stable and does **not** start T10.

Canonical T9 evidence owners:
- `governance/programs/ADAPT-V18-FINAL-CLOSURE-10-10-001/T9_PACKAGED_CROSS_PLATFORM_RELEASE_ASSURANCE.json`
- `release/v18.10.0/release_contract.json`
- `release/v18.10.0/certification-manifest.json`
- `tools/ci/v18_t9_packaged_cross_platform_release_assurance_gate.py`
- `tools/release/native_macos.sh`
- `tools/release/native_windows.ps1`
- `tools/release/g15_assurance.py`
- `tools/ci/impact_plan.py`
- `.github/workflows/ci-fast.yml`
- `.github/workflows/ci-qualified.yml`
- `.github/workflows/release.yml`

The candidate identity is synchronized across `release_identity.json`, `app_bootstrap.go`, `VERSION.txt`, `renderer/release-identity.js` and `renderer/index.html`:
- product version `18.10.0`
- build ID `v18.10.0-stable-20260825`
- platform build `181000`
- previous Stable `v18.9.1`
- runtime/config continuity `PersonalMarketTerminal`

Previous-Stable upgrade authority is immutable `v18.9.1-stable` at candidate `e55d8d25b15cec2ffb0f5411bc358bc40b359cf9`, fingerprint `0062f46dea5690d0b3fcd8a9ed3b1f71ebe1522c7dee2cb218e9d36b9e0076ff`, build `v18.9.1-stable-20260821`. Both native harnesses invoke the previous Stable tag's own certified native harness in an isolated worktree and then launch the exact v18.10.0 package against the same profile to prove migrations, SQLite integrity and persisted identity/symbol state are preserved.

The T9 Planner v3 contract requires one unchanged source head to run:
- backend full Go suite
- race detector
- randomized package order
- persistence/DB integration
- renderer contracts
- Chrome primary behavior
- WebKit primary compatibility
- security/data-rights contracts
- macOS Apple Silicon actual packaged lifecycle
- Windows x64 actual packaged lifecycle

macOS must prove clean extraction, arm64/SQLite/code-sign identity, packaged backend fresh/warm relaunch, actual Cocoa/WKWebView fresh/warm lifecycle, and v18.9.1 profile upgrade. Windows must prove clean extraction, PE x64 identity, packaged fresh/warm relaunch, and v18.9.1 profile upgrade. Both native evidence files bind exact source SHA, Git-object fingerprint, build ID and package SHA-256.

Current T9 open evidence gaps are intentionally limited to exact-head Fast, identical-head full Qualified, and native evidence binding from that Qualified run. T9 remains `IMPLEMENTED_UNVERIFIED` until those executable results pass.

## Exactly one next action

Fetch fresh `main` and `adapt-v18-final-closure-10-10-001`, confirm #122 remains active and no newer source has superseded the branch, then finish T9 only: run exact-head Fast, run identical-head Qualified with backend/race/randomized + DB + renderer + Chrome + WebKit + security/data-rights + macOS + Windows selected, inspect the native package evidence, fix any real failures, and only then mark T9 VERIFIED/merge/close. Do not start T10, publish v18.10.0 Stable or begin v19 in the same step before T9 is durably verified.

## Retained architecture

Smart Provider Router v2 remains sole routing/admission authority. Direct SEC/EDGAR remains filing/Form 4 authority. Canonical Data Health/freshness/degradation/cache/persistence/subscription/telemetry/reconciliation/lifecycle and identity/session owners remain authoritative. Provider usefulness remains advisory only. U.S. equities processing, GLD/SLV/USO actionable exceptions, governed SHADOW → VALIDATED → APPROVED → PRODUCTION provider maturity, and No Execution remain permanent. Linux is CI/test only; required packaged release targets are macOS Apple Silicon and Windows x64. Hosted Web GA remains v19 scope. No parallel subsystem may be created merely to satisfy closure testing.

## Resume rule

1. Fetch live `main`, `adapt-v18-final-closure-10-10-001`, #113, #122 and current PR/workflow state before modifying anything.
2. Read this file and `governance/current-state.json` first, then the T9 assurance artifact, v18.10 release contract/manifest and parent closure ledger.
3. GitHub objects and executable evidence outrank this prose and all chat memory.
4. Treat T1-T8 as complete only while their durable evidence remains intact; T8 authority is source `e4542890...`, Fast #1080, Qualified #216 and merge `a3a82aee...`.
5. T9/#122 is IN_PROGRESS; T10/#123 is NOT_STARTED and v19/#66 is blocked.
6. Never infer assurance from file existence or issue state. Require exact-head executable positive, negative/failure and platform evidence.
7. Preserve Smart Provider Router v2, canonical Data Health/freshness/cache/persistence/subscription/telemetry/reconciliation/lifecycle/identity/session owners, direct SEC/EDGAR, U.S. equities, GLD/SLV/USO and No Execution.
8. Use G0-G16 only. T9 must earn exact-head Fast and identical-head Qualified; final public v18.10.0 publication is T10-only after T1-T10 are VERIFIED.
9. A new ChatGPT account, Codex or Claude must be able to resume from GitHub alone.
