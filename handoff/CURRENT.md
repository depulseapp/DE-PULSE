# DE.PULSE — Current Authoritative Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS.** GitHub is authoritative; chat history and temporary workspaces are advisory only.

**Release:** `v18.6.0`  
**Active branch:** `v18.6-development`  
**Stable predecessor:** `v18.5.2-stable` / G0–G16 CLOSED  
**Current candidate state:** v18.6.0 G12 HARNESS REMEDIATION / FRESH G10 REQUIRED / NOT PROMOTED  
**Repository:** `depulseapp/DE-PULSE`  
**Main release PR:** `#16`  
**v18.6 runtime build ID:** `v18.6.0-stable-20260818`  
**Last updated:** 2026-08-18 America/Vancouver

## Resume rule

Read `AGENTS.md` or `CLAUDE.md`, then reconcile this handoff, `release_identity.json`, both `.depulse-certification/resume/` checkpoints, PR #16, current checks, the release-certification branch, and the immutable Stable predecessor. Never resume from model memory alone.

`source_fingerprint.py` excludes `.depulse-certification` only. Workflow, certification-harness, test-contract, and this handoff are source-fingerprinted. Any change to them requires fresh G10 before G11.

## Product / architecture state

All eight assigned v18.6 implementation/audit slices remain code-complete: watchlist remediation; shared Scanner/Radar broad-snapshot acquisition; serialized Session Intelligence Coordinator; Market Activity/legacy-route consolidation; role-aware documentation; external dependency/provider readiness; bounded AI context/cache/schema/evaluation hardening; and provider×dataset rights-aware fail-closed AI egress.

Protected invariants remain unchanged:
- deterministic Day/Swing/Long formulas;
- Smart Provider Router is the sole provider-routing owner and provider count never changes Market Mode;
- GLD/SLV/USO remain tradable live exceptions;
- desktop SQLite / hosted PostgreSQL architecture;
- U.S. Equities Processing Boundary;
- permanent No Execution Boundary;
- exactly three canonical workflows: `ci-fast.yml`, `ci-qualified.yml`, `release.yml`;
- G0–G16 is the only top-level gate model.

## Last trustworthy qualification and release evidence

Source `51be6afc444d8337821dd3c6ecaaecb1297dc84b`, exercised at metadata head `836322345886c9236e1da6ee774d73e6d27463b8`, passed final G10:
- CI Fast #234 / run `32196295742` — PASS;
- CI Qualified #99 / run `32196295782` — PASS;
- Ubuntu/macOS/Windows portability, workflow/provenance, browser behavior, renderer contracts, full Go suite, race detector, randomized package order and evidence summary all passed.

Release-certification base `659f4e59286ede7c8627d5cefa7deceab792c8b3` was dispatched through trigger PR #24 without merging it. Dispatcher CI Fast #238 / run `32196732001` passed. It resolved:
- release ref `v18.6-release-certification`;
- candidate `659f4e59286ede7c8627d5cefa7deceab792c8b3`;
- canonical Git-object source fingerprint `22b14764656a08d3677eb05187232afde3029b5f35ea7f5c21b3f85b7d41f242`;
- `publish=false`;
- canonical `release.yml` run `32196711609`.

Workflow-token PR comments remain unavailable with HTTP 403, but this is non-blocking observability only. Exact run identity is retained in Actions evidence and was durably bound to PR #16 through the connected GitHub app.

Canonical release run `32196711609`:
- G11 immutable candidate/provenance — PASS;
- G12 full certification — FAIL;
- G13/G14 macOS, G13/G14 Windows, G15 and G16 — correctly skipped;
- no Stable publication occurred.

## G12 failure classification

The G12 failure is a **certification-harness contradiction, not a demonstrated product regression**.

`release/v18.5.1/browser_watchlist_membership_test.py` asserts the historical `CURRENT` / `aria-current="true"` / `current-desk` watchlist model. The v18.6 contract intentionally replaced that model with Day/Swing/Long `aria-pressed` membership state and explicitly requires `CURRENT`, `aria-current`, and `current-desk` to be absent. `release/v18.6.0/browser_watchlist_membership_test.py` proves the current semantics, and `tools/ci/watchlist_membership_contract.py` already binds CI Qualified to that v18.6 proof.

`release/v18.6.0/run_full_certification.sh` incorrectly executed **both** the obsolete v18.5.1 membership proof and the v18.6 replacement against the same current production extension. The old test failed first at `assert "aria-current=\"true\"" in extension`, making the G12 list internally contradictory.

## Remediation

The v18.5.1 historical test file remains untouched for audit/history. The v18.6 G12 execution list now excludes only `release/v18.5.1/browser_watchlist_membership_test.py` and continues to run the v18.6 browser membership proof plus all other retained browser regressions.

`tools/ci/watchlist_membership_contract.py` is hardened to require:
- CI Qualified uses `release/v18.6.0/browser_watchlist_membership_test.py`;
- G12 full certification also includes that v18.6 proof;
- G12 must not invoke the contradictory v18.5.1 CURRENT/aria-current membership proof.

This preserves regression coverage while preventing the stale semantic contract from silently re-entering v18.6 certification.

Because the G12 script, watchlist contract, and this handoff are source-fingerprinted, the prior G10/G11 evidence cannot certify the remediated source. Fresh G10 is mandatory before a new G11/G12 run.

## Release orchestration contract

Certification remains non-publishing. The primary path is exact `v<release-line>-release-certification` push with `publish=false`; the connector-safe fallback is an owner-triggered pull request whose base is exact release-certification, using immutable base ref/SHA and forcing `publish=false`. Stable promotion remains exact `*-stable-promotion`, push-only, owner-gated, evidence-bound, and no-rebuild. There is no PR fallback for Stable publication.

PR tracking through the workflow token is best-effort. A comment denial cannot block certification after exact release ref/SHA/fingerprint/publish state is proven. The dispatcher retains canonical run identity in Actions logs/summary; the connected GitHub app binds that exact identity to PR #16 when required.

## Exactly one next action

**Qualify the consolidated G12-harness remediation as one fresh G10 candidate (Fast + Qualified). If both pass, bind only fingerprint-excluded checkpoints, reconcile `v18.6-release-certification`, create one unmerged fingerprint-excluded trigger PR, dispatch one canonical `release.yml` run with `publish=false`, bind its exact identity to PR #16, and continue G11 → G12 → G13/G14 macOS + Windows → G15 → G16. Do not promote Stable unless every gate passes.**

## Known residuals / User Action Required

- TradeInsight remains governed future/shadow implementation work until separately validated and approved; it gains no production influence automatically.
- Deployment-specific provider keys, entitlements and commercial/redistribution/AI-use rights remain User Action Required where absent or unapproved and fail closed.
- v18.6 native macOS Apple Silicon and Windows x64 evidence remains pending until the remediated candidate passes G12 and native G13/G14.

## Provider-neutral continuation instruction

> Connect to `depulseapp/DE-PULSE`, read `AGENTS.md` or `CLAUDE.md`, then inspect `release_identity.json`, `handoff/CURRENT.md`, both `.depulse-certification/resume/` checkpoints, PR #16, active certification PR/checks, and the release-certification branch. Treat `v18.5.2-stable` as immutable Stable until v18.6 G11–G16 and no-rebuild promotion complete. Resume from the exact trustworthy evidence above, not chat memory. Preserve G0–G16, assistant/account independence, Smart Provider Router ownership, deterministic desk formulas, U.S. Equities Processing Boundary and permanent No Execution Boundary.
