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
)


def main() -> int:
    errors: list[str] = []
    if not STATE.is_file():
        print("DE.PULSE current-state projection convergence: FAIL", file=sys.stderr)
        print(" - governance/current-state.json missing", file=sys.stderr)
        return 1

    state = json.loads(STATE.read_text(encoding="utf-8"))
    stable = state.get("stable", {})
    active = state.get("activeWorkSlice", {})
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

    handoff = texts.get("handoff/CURRENT.md", "")
    if stable_sha and stable_sha not in handoff:
        errors.append("handoff/CURRENT.md does not project certified Stable candidate SHA")
    if closure and closure not in handoff:
        errors.append("handoff/CURRENT.md does not name the canonical closure ledger")

    ci_projection = texts.get("adaptive-governance/CURRENT_ADAPTIVE_CI_CONVERGENCE.md", "")
    if closure and closure not in ci_projection:
        errors.append("CURRENT_ADAPTIVE_CI_CONVERGENCE.md does not name canonical closure ledger")

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
    print(f"active work slice: #{issue} / {work_slice} / {branch}")
    print(f"projected current surfaces: {len(texts)}/{len(SURFACES)}")
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
