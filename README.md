# DE.PULSE

DE.PULSE is a U.S.-equity research and decision-support application with a permanent **No Execution** boundary. Deterministic Day/Swing/Long market truth remains protected while adaptive intelligence, provider selection, reliability and synthesis evolve under governed evidence.

## Current project truth

Do **not** infer the current Stable version from this README. The authoritative live project state is intentionally version-neutral here so the root README cannot become a stale release ledger again.

Use, in order:
1. GitHub immutable Stable tag + GitHub Release and certified artifacts;
2. `release_identity.json` and `VERSION.txt`;
3. `.depulse-certification/resume/build-checkpoint.json` and `release-evidence-checkpoint.json`;
4. `handoff/CURRENT.md`;
5. `adaptive-governance/CURRENT_ADAPTIVE_ROADMAP.md`, `CURRENT_ADAPTIVE_BUILD_PLAN.md`, `CURRENT_ADAPTIVE_BUILD_PROCESS.md`, and `CURRENT_ADAPTIVE_DELIVERY_PROCESS.md`.

A fresh ChatGPT/Codex account starts with `AGENTS.md`; Claude starts with `CLAUDE.md`. Both point to `governance/AI-ASSISTANT-PORTABILITY-CONTRACT.md`. GitHub is the durable source of truth; chat memory is never required.

## Core permanent boundaries

- Actionable processing: U.S.-listed equities and approved U.S.-listed tradable exceptions/context under canonical eligibility rules.
- No execution, brokerage order routing, autonomous trading or simulated execution product scope.
- Smart Provider Router v2 is the sole provider-routing authority.
- BroadSnapshotBroker is the canonical broad snapshot reuse/coalescing owner.
- Canonical freshness/evidence-time truth must never fabricate `now` or hide genuine stale/missing evidence.
- GLD, SLV and USO remain actionable tradable live-priority exceptions.
- Adaptive production influence follows `SHADOW → VALIDATED → APPROVED → PRODUCTION`.
- G0–G16 is the only top-level release model.

## Build / release lifecycle

Normal release flow:

`reconcile Stable → G0 exact baseline → G1 immutable scope → G2/G3 architecture/design → implementation + targeted evidence → Fast → Qualified → exact-head merge → G11-G16 certify/package/audit/publish → post-Stable continuity reconciliation`.

Use one `v<version>-development` branch and one PR for a normal release. Do not create retry/certification/promotion/dispatch branches or duplicate workflow families merely to trigger CI.

The only canonical GitHub Actions workflows are documented in `governance/ACTIVE_WORKFLOW_MANIFEST.md`.

## Runnable artifacts

Primary location: repository **Releases** with immutable Stable tags. Forward release evidence lives under `release/<version>/` and `.depulse-certification/`.

Required native Stable targets are macOS Apple Silicon and Windows x64. Desktop runtime/config continuity uses `PersonalMarketTerminal` unless a future approved migration changes that contract.

## Reliability north star

`DATA DEGRADED` must mean decision-relevant evidence is genuinely degraded—not that DE.PULSE overloaded itself, duplicated work, rejected valid evidence incorrectly or allowed optional context to contaminate unrelated consumers.

Canonical reliability contract: `governance/ADAPTIVE-DATA-RELIABILITY-CONTRACT.md`.

## Historical release detail

Historical release notes, packages, certification evidence and implementation history remain preserved in Git history, immutable tags/Releases, `release/<version>/`, governance ledgers and Stable evidence manifests. The root README is deliberately not an accumulating version-by-version release ledger.
