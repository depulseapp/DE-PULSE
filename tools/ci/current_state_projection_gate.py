#!/usr/bin/env python3
"""Executable convergence guard for DE.PULSE current-state projections."""
from __future__ import annotations

import json
from pathlib import Path
import sys

ROOT = Path(__file__).resolve().parents[2]
STATE = ROOT / "governance" / "current-state.json"
SURFACES = (
    "handoff/CURRENT.md",
    "adaptive-governance/CURRENT_ADAPTIVE_ROADMAP.md",
    "adaptive-governance/CURRENT_ADAPTIVE_BUILD_PLAN.md",
    "adaptive-governance/CURRENT_ADAPTIVE_BUILD_PROCESS.md",
    "adaptive-governance/CURRENT_ADAPTIVE_DELIVERY_PROCESS.md",
    "adaptive-governance/CURRENT_ADAPTIVE_CI_CONVERGENCE.md",
    "adaptive-governance/CURRENT_ADAPTIVE_GAP_CLOSURE.md",
)
WAIVER_PROJECTIONS = (
    "handoff/CURRENT.md",
    "adaptive-governance/CURRENT_ADAPTIVE_CI_CONVERGENCE.md",
    "adaptive-governance/CURRENT_ADAPTIVE_GAP_CLOSURE.md",
)
CANONICAL_ADAPTIVE_DOCS = {
    "governance/ADAPTIVE-OPERATING-CONTRACT.md": (
        "SymbolIntelligenceSnapshot",
        "Opportunity Lifecycle",
        "SHADOW -> VALIDATED -> APPROVED -> PRODUCTION",
        "CURRENT_ADAPTIVE_*",
    ),
    "governance/ROADMAP.md": (
        "Rebaselined:** 2026-08-28",
        "Watchlist",
        "v19.4.0",
        "v20.0.0",
        "governance/current-state.json",
    ),
    "adaptive-governance/ADAPTIVE_BUILD_PLAN.md": (
        "RuntimeSnapshot",
        "SymbolIntelligenceSnapshot",
        "Opportunity Lifecycle",
        "Watchlist",
        "point-in-time",
    ),
    "adaptive-governance/ADAPTIVE_BUILD_PROCESS.md": (
        "CHARACTERIZE -> NEW OWNER -> DUAL/SHADOW",
        "Smart Provider Router v2",
        "PARTIAL COVERAGE",
        "CENSORED",
    ),
    "adaptive-governance/ADAPTIVE_DELIVERY_PROCESS.md": (
        "actual artifact/runtime audit",
        "Commercial/Public",
        "macOS Apple Silicon",
        "Windows x64",
        "Web",
    ),
    "adaptive-governance/README.md": (
        "DEC-2026-08-28-001",
        "Canonical narrative authorities",
        "Current projections",
        "Historical evidence",
    ),
}


def projected_active_state(state: dict) -> tuple[dict, str]:
    product = state.get("productCapabilityGate", {})
    if isinstance(product, dict) and str(product.get("reservationStatus", "")).strip() == "IN_PROGRESS":
        return {
            "workSliceId": product.get("reservedWorkSliceId"),
            "issue": product.get("reservedIssue"),
            "branch": product.get("reservedBranch"),
            "closureLedger": product.get("closureLedger"),
        }, "PRODUCT_CAPABILITY"
    active = state.get("activeWorkSlice", {})
    return (active if isinstance(active, dict) else {}), "WORK_SLICE"


def main() -> int:
    errors: list[str] = []
    if not STATE.is_file():
        print("DE.PULSE current-state projection convergence: FAIL", file=sys.stderr)
        print(" - governance/current-state.json missing", file=sys.stderr)
        return 1

    state = json.loads(STATE.read_text(encoding="utf-8"))
    stable = state.get("stable", {})
    process_active = state.get("activeWorkSlice", {})
    if not isinstance(process_active, dict):
        process_active = {}
    active, projection_authority = projected_active_state(state)
    stable_tag = str(stable.get("tag", "")).strip()
    stable_sha = str(stable.get("candidateSha", "")).strip()
    work_slice = str(active.get("workSliceId", "")).strip()
    issue = active.get("issue")
    branch = str(active.get("branch", "")).strip()
    closure = str(active.get("closureLedger", "")).strip()
    if not all((stable_tag, stable_sha, work_slice, branch, closure)) or issue is None:
        errors.append("canonical current-state is missing stable/work-slice projection fields")

    required_common = (
        "governance/current-state.json",
        stable_tag,
        work_slice,
        branch,
    )
    texts: dict[str, str] = {}
    for rel in SURFACES:
        path = ROOT / rel
        if not path.is_file():
            errors.append(f"missing current projection: {rel}")
            continue
        text = path.read_text(encoding="utf-8", errors="replace")
        texts[rel] = text
        for token in required_common:
            if token and token not in text:
                errors.append(f"{rel} does not project canonical token: {token}")
        if f"#{issue}" not in text:
            errors.append(f"{rel} does not project active issue #{issue}")

        if rel.startswith("adaptive-governance/CURRENT_ADAPTIVE_"):
            if "COMPATIBILITY PROJECTION" not in text:
                errors.append(f"{rel} is not marked as a compatibility projection")
            if "## Exactly one next action" in text:
                errors.append(f"{rel} duplicates the handoff next action")
            if len(text.splitlines()) > 40:
                errors.append(f"{rel} is no longer a thin current-state projection")

    for rel, tokens in CANONICAL_ADAPTIVE_DOCS.items():
        path = ROOT / rel
        if not path.is_file():
            errors.append(f"missing canonical adaptive document: {rel}")
            continue
        text = path.read_text(encoding="utf-8", errors="replace")
        for token in tokens:
            if token not in text:
                errors.append(f"{rel} is missing audit-rebaseline token: {token}")

    compatibility_roadmap = ROOT / "adaptive-governance" / "ADAPTIVE_ROADMAP.md"
    if not compatibility_roadmap.is_file():
        errors.append("missing adaptive-governance/ADAPTIVE_ROADMAP.md compatibility path")
    else:
        text = compatibility_roadmap.read_text(encoding="utf-8", errors="replace")
        if "COMPATIBILITY POINTER" not in text or "governance/ROADMAP.md" not in text:
            errors.append("ADAPTIVE_ROADMAP.md must remain a pointer to the canonical roadmap")
        if "## Exactly one next action" in text or len(text.splitlines()) > 30:
            errors.append("ADAPTIVE_ROADMAP.md contains duplicate current roadmap/next-action narrative")

    handoff = texts.get("handoff/CURRENT.md", "")
    if stable_sha and stable_sha not in handoff:
        errors.append("handoff/CURRENT.md does not project certified Stable candidate SHA")
    if closure and closure not in handoff:
        errors.append("handoff/CURRENT.md does not name the canonical closure ledger")

    next_band_id = str(process_active.get("nextDependencyBandId", "")).strip()
    if next_band_id:
        next_action = handoff.split("## Exactly one next action", 1)
        next_action_text = next_action[1].split("\n## ", 1)[0] if len(next_action) == 2 else ""
        if next_band_id not in next_action_text:
            errors.append(
                "handoff/CURRENT.md next action does not match "
                f"activeWorkSlice.nextDependencyBandId: {next_band_id}"
            )

    ci_projection = texts.get("adaptive-governance/CURRENT_ADAPTIVE_CI_CONVERGENCE.md", "")
    if closure and closure not in ci_projection:
        errors.append("CURRENT_ADAPTIVE_CI_CONVERGENCE.md does not name canonical closure ledger")

    waivers = process_active.get("externalControlWaivers", [])
    if waivers is not None and not isinstance(waivers, list):
        errors.append("activeWorkSlice.externalControlWaivers must be an array")
        waivers = []
    for waiver in waivers or []:
        if not isinstance(waiver, dict):
            errors.append("externalControlWaivers entries must be objects")
            continue
        wid = str(waiver.get("id", "")).strip()
        gap = str(waiver.get("gapId", "")).strip()
        path = str(waiver.get("path", "")).strip()
        status = str(waiver.get("status", "")).strip()
        if not all((wid, gap, path, status)):
            errors.append("external control waiver projection fields missing")
            continue
        if not (ROOT / path).is_file():
            errors.append(f"external control waiver file missing: {path}")
        for rel in WAIVER_PROJECTIONS:
            text = texts.get(rel, "")
            for token in (wid, gap, path):
                if token not in text:
                    errors.append(f"{rel} does not project external waiver token: {token}")

    stale_claims = (
        "Status:** NOT STARTED",
        "Active work slice:** #64",
        "Active work slice:** #69",
        "Certified Stable:** `v18.9.0-stable`",
    )
    for rel, text in texts.items():
        for stale in stale_claims:
            if stale in text:
                errors.append(f"stale current-state claim re-emerged in {rel}: {stale}")

    print("DE.PULSE current-state projection convergence")
    print(f"canonical Stable: {stable_tag} @ {stable_sha}")
    print(f"projection authority: {projection_authority}")
    print(f"active projected work: #{issue} / {work_slice} / {branch}")
    print(f"projected current surfaces: {len(texts)}/{len(SURFACES)}")
    print(f"projected external-control waivers: {len(waivers or [])}")
    if errors:
        print("DE.PULSE current-state projection convergence: FAIL", file=sys.stderr)
        for error in errors:
            print(" - " + error, file=sys.stderr)
        return 1
    print("machine current-state authority: PASS")
    print("human projection drift guard: PASS")
    print("DE.PULSE current-state projection convergence: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
