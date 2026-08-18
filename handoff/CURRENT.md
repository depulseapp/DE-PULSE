# DE.PULSE — Current Authoritative Handoff

**Status:** v18.5.2 STABLE / G0–G16 CLOSED  
**Repository:** `depulseapp/DE-PULSE`  
**Stable tag:** `v18.5.2-stable`  
**Stable promotion commit:** `d30e54db4908ca57c52ae298cc4ada3416fab46b`  
**Certified source checkout:** `e9ca615cf7ab97ac476128c81ee9ae2d7340c0d9`  
**Certified source fingerprint:** `807de082d43e83d1d3548bca9350d13b72ef4dc71a848940b73b63e4b4d215b0`  
**Runtime build ID:** `v18.5.2-stable-hotfix-20260817`  
**PR #13:** merged  
**Machine state:** `.depulse-certification/resume/build-checkpoint.json`  
**Last updated:** 2026-08-17 America/Vancouver

## Authority and resume rule

GitHub is authoritative. Start with `AGENTS.md` or `CLAUDE.md`, then reconcile this handoff and the machine checkpoints against the immutable Stable tag/release and current repository state.

`v18.5.2-stable` is the immutable product/release authority. `main` may contain later **post-release G16 operational metadata or handoff commits**; those do not redefine the certified Stable product source, native artifact hashes, or Stable tag. Do not infer a new product release from post-release documentation/checkpoint commits.

## v18.5.2 final result

The desktop recovery is complete and promoted:

- Day, Swing and Long-Term desk rendering recovered;
- complete ET/PT clocks retained in the compact Market Pulse Ribbon;
- Research hierarchy/reading width recovered;
- tracked-symbol draft and first-attempt add behavior recovered;
- symbol controls aligned with responsive stacking;
- configurable display/sign-in names added while OWNER remains a separate role;
- Settings viewport/save-row behavior corrected;
- deterministic Day/Swing/Long formulas remain protected;
- U.S. Equities Processing and permanent No Execution boundaries remain intact.

Provider → Market Mode governance remains in force with 22 capability assessments across 12 providers. Provider count alone cannot change a mode; deterministic/statistical code owns numeric Market Mode truth. TradeInsight remains `NOT_IMPLEMENTED` with `NONE` production influence until its adapter and lifecycle evidence close.

## Certification and native delivery

### G0–G12 — PASS

Canonical exact-source certification passed at checkout `e9ca615cf7ab97ac476128c81ee9ae2d7340c0d9` with fingerprint `807de082d43e83d1d3548bca9350d13b72ef4dc71a848940b73b63e4b4d215b0`.

Key protected result: deterministic trading equivalence **2403/2403 PASS**.

Four G12 defects/test-contract defects were found and closed before the final full pass: bundled v18.5.2 QA identity, refactor-safe role/identity regression coverage, Playwright scroll baseline ordering, and release-bound watchlist cache-busting coverage.

### G13/G14 — PASS on both required native targets

**macOS Apple Silicon** — run `32102570376`

- actual arm64 packaged artifact runtime audit PASS;
- native `libsqlite3` linkage PASS;
- clean extraction/code-sign verification PASS;
- health/release identity PASS;
- SQLite ready/integrity/migrations 1–4 PASS;
- package: `De-Pulse-v18.5.2-Stable-macOS-Apple-Silicon.zip`;
- SHA-256: `91f7d64a433474c4efbed0bd5c7d065508b2ce6eeae0544b1c217849fe62a4ae`.

**Windows x64** — run `32103078336`

- actual AMD64 packaged artifact runtime audit PASS;
- PE x64 and real `winsqlite3` runtime PASS;
- clean extraction/launch PASS;
- health/release identity PASS;
- SQLite ready/integrity/migrations 1–4 PASS;
- package: `De-Pulse-v18.5.2-Stable-Windows-x64.zip`;
- SHA-256: `6431562a67bcd55db6ebab2e6a09724006119910ea8adbc358afb8644a752326`.

### G15 — PASS

Run `32103078336` consumed the preserved macOS PASS and final Windows PASS, verified both exact package hashes/evidence graphs, preserved the No Execution Boundary, and returned `promotionAuthorized=true`.

PR #13 was then promoted with exact-head protection. Merge/promotion commit: `d30e54db4908ca57c52ae298cc4ada3416fab46b`.

Stable publication run `32103499711` re-downloaded and re-verified the already-certified artifacts, created `v18.5.2-stable` at the promotion commit, uploaded the certified packages/evidence without rebuilding, and verified the release is neither draft nor prerelease.

Published release assets include both native packages, both G13/G14 evidence JSON files, G15 assurance JSON, and the G15 SHA-256 manifest.

## G16 retrospective / reusable learning

1. **Preserve independent platform PASS.** macOS passed on the first native run and was never rerun merely because Windows failed. G15 reused the preserved exact macOS artifact.
2. **Classify infrastructure/test-harness failures truthfully.** Earlier Actions pre-step failures were infrastructure failures, not product failures. Windows native attempts that stopped before build on fingerprint materialization were harness/provenance failures, not runtime failures.
3. **Cross-platform source provenance must use canonical Git bytes.** Windows filesystem materialization produced a different byte fingerprint even for the exact Git commit. The successful Windows lane recomputed the certified fingerprint directly from raw Git object bytes and treated the materialized filesystem fingerprint as diagnostic only. Future process hardening should make this platform-neutral provenance method canonical without modifying the already-certified v18.5.2 source.
4. **Do not rebuild during publication.** Stable publication uploaded the exact artifacts already certified at G13/G14/G15.
5. **Minimize unnecessary spend/reruns.** G0–G12 was not rerun after it became authoritative; macOS was not rerun for Windows-only harness fixes; publication performed verification/upload only.
6. **Release identity and post-release metadata are separate.** The immutable Stable tag/release and artifact hashes define v18.5.2. Later G16 documentation/checkpoint commits on `main` are operational metadata and cannot silently redefine Stable.

## Truthful residuals after v18.5.2 closure

- Owner/Admin/User role-specific tab/navigation composition remains future work; v18.5.2 does not claim it.
- TradeInsight secure key configuration, adapter/normalization, Smart Router SHADOW routing, validation and promotion remain future work; no API key is committed.
- Cross-platform fingerprint implementation should be hardened around canonical Git-object provenance in a future authorized release.
- Native-certification/recovery orchestration branches are retained while they remain useful audit references; cleanup may occur later when evidence retention no longer depends on them.

None of these residuals blocks the certified v18.5.2 Stable release.

## Exactly one next action

For any future authorized DE.PULSE work, **start from the immutable `v18.5.2-stable` baseline and reconcile the remaining approved backlog/audit items before selecting the next release scope.** Do not infer or start a new major release merely from chat history.

## Provider-neutral continuation instruction

> Connect to `depulseapp/DE-PULSE`, read `AGENTS.md` or `CLAUDE.md`, resolve Stable through `v18.5.2-stable`, verify it maps to promotion commit `d30e54db4908ca57c52ae298cc4ada3416fab46b`, and read `.depulse-certification/resume/build-checkpoint.json` plus this handoff. Treat later `main` commits that only close G16 documentation/checkpoints as post-release operational metadata, not a new product baseline. Preserve G0–G16, the U.S. Equities Processing Boundary, No Execution Boundary, deterministic Market Mode ownership, TradeInsight `NOT_IMPLEMENTED/NONE` truth, and the native package hashes recorded above.
