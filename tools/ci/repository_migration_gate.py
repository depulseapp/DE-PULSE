#!/usr/bin/env python3
"""Canonical #70 repository migration/root/current-state convergence gate."""
from __future__ import annotations

import json
from pathlib import Path
import tempfile

import current_state_projection_gate
import repository_migration_gate_impl
import root_ownership_gate

ROOT = Path(__file__).resolve().parents[2]
POLICY_PATH = ROOT / "governance" / "root-layout-policy.json"
MIGRATIONS_PATH = ROOT / "governance" / "repository-migrations.json"


def synthesized_migration_policy() -> dict[str, object]:
    policy = json.loads(POLICY_PATH.read_text(encoding="utf-8"))
    migrations = json.loads(MIGRATIONS_PATH.read_text(encoding="utf-8"))
    transitional: dict[str, dict[str, str]] = {}
    moves = migrations.get("moves", [])
    for row in moves if isinstance(moves, list) else []:
        if not isinstance(row, dict):
            continue
        new_path = str(row.get("newPath", "")).strip()
        if not new_path or "/" in new_path:
            continue
        owner = str(row.get("owner", "")).strip()
        reason = str(row.get("reason", "")).strip()
        removal = str(row.get("removalCondition", "")).strip()
        if not all((owner, reason, removal)):
            continue
        transitional[new_path] = {
            "owner": owner,
            "reason": reason,
            "expiry": "REGISTERED_MIGRATION_TARGET",
            "removalCondition": removal,
        }
    compatibility = dict(policy)
    compatibility["transitionalRootFiles"] = transitional
    return compatibility


def main() -> int:
    if root_ownership_gate.main() != 0:
        return 1
    if current_state_projection_gate.main() != 0:
        return 1

    compatibility = synthesized_migration_policy()
    original_policy_path = repository_migration_gate_impl.POLICY_PATH
    with tempfile.NamedTemporaryFile(mode="w", encoding="utf-8", suffix=".json", delete=True) as tmp:
        json.dump(compatibility, tmp, indent=2, sort_keys=True)
        tmp.write("\n")
        tmp.flush()
        repository_migration_gate_impl.POLICY_PATH = Path(tmp.name)
        try:
            result = repository_migration_gate_impl.main()
        finally:
            repository_migration_gate_impl.POLICY_PATH = original_policy_path
    if result != 0:
        return result
    print("DE.PULSE canonical repository/root/current-state convergence: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
