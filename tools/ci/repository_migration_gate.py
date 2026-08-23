#!/usr/bin/env python3
"""Canonical repository migration/root/current-state convergence gate."""
from __future__ import annotations

import json
from pathlib import Path
import tempfile

import current_state_projection_gate
import repository_migration_gate_impl
import repository_migration_registry
import root_convergence_gate
import root_ownership_gate

ROOT = Path(__file__).resolve().parents[2]
POLICY_PATH = ROOT / "governance" / "root-layout-policy.json"
EPHEMERAL_DIR = ROOT / ".depulse-certification"


def synthesized_migration_policy(migrations: dict[str, object]) -> dict[str, object]:
    policy = json.loads(POLICY_PATH.read_text(encoding="utf-8"))
    transitional: dict[str, dict[str, str]] = {}
    for row in migrations.get("moves", []) if isinstance(migrations.get("moves", []), list) else []:
        if not isinstance(row, dict):
            continue
        new = str(row.get("newPath", "")).strip()
        owner = str(row.get("owner", "")).strip()
        reason = str(row.get("reason", "")).strip()
        removal = str(row.get("removalCondition", "")).strip()
        if new and "/" not in new and all((owner, reason, removal)):
            transitional[new] = {
                "owner": owner,
                "reason": reason,
                "expiry": "REGISTERED_MIGRATION_TARGET",
                "removalCondition": removal,
            }

    final = policy.get("finalRootEvidenceFiles", {})
    if isinstance(final, dict):
        for path, meta in final.items():
            if not isinstance(meta, dict) or "/" in str(path):
                continue
            owner = str(meta.get("owner", "")).strip()
            reason = str(meta.get("reason", "")).strip()
            if owner and reason:
                transitional[str(path)] = {
                    "owner": owner,
                    "reason": reason,
                    "expiry": "FINAL_PACKAGE_MAIN_EVIDENCE",
                    "removalCondition": "MOVE_WITH_CAPABILITY_ONLY_WHEN_PRIVATE_PACKAGE_MAIN_ACCESS_CAN_BE_PRESERVED_WITHOUT_TEST_ONLY_PRODUCTION_EXPORTS",
                }
    compatibility = dict(policy)
    compatibility["transitionalRootFiles"] = transitional
    return compatibility


def main() -> int:
    try:
        migrations = repository_migration_registry.load_repository_migrations()
    except Exception as exc:
        print("DE.PULSE canonical repository migration registry: FAIL")
        print(f" - {exc}")
        return 1

    compatibility = synthesized_migration_policy(migrations)
    original_root_migrations = root_ownership_gate.MIGRATIONS
    original_impl_policy = repository_migration_gate_impl.POLICY_PATH
    original_impl_migrations = repository_migration_gate_impl.MIGRATIONS_PATH
    EPHEMERAL_DIR.mkdir(parents=True, exist_ok=True)

    with tempfile.NamedTemporaryFile(mode="w", encoding="utf-8", suffix=".json", prefix="composed-migrations-", dir=EPHEMERAL_DIR, delete=True) as migration_tmp, tempfile.NamedTemporaryFile(mode="w", encoding="utf-8", suffix=".json", prefix="composed-root-policy-", dir=EPHEMERAL_DIR, delete=True) as policy_tmp:
        json.dump(migrations, migration_tmp, indent=2, sort_keys=True)
        migration_tmp.write("\n")
        migration_tmp.flush()
        json.dump(compatibility, policy_tmp, indent=2, sort_keys=True)
        policy_tmp.write("\n")
        policy_tmp.flush()

        root_ownership_gate.MIGRATIONS = Path(migration_tmp.name)
        repository_migration_gate_impl.MIGRATIONS_PATH = Path(migration_tmp.name)
        repository_migration_gate_impl.POLICY_PATH = Path(policy_tmp.name)
        try:
            if root_ownership_gate.main() != 0:
                return 1
            if root_convergence_gate.main() != 0:
                return 1
            if current_state_projection_gate.main() != 0:
                return 1
            result = repository_migration_gate_impl.main()
        finally:
            root_ownership_gate.MIGRATIONS = original_root_migrations
            repository_migration_gate_impl.POLICY_PATH = original_impl_policy
            repository_migration_gate_impl.MIGRATIONS_PATH = original_impl_migrations

    if result != 0:
        return result
    ledgers = migrations.get("composedWorkSliceLedgers", [])
    print(f"composed work-slice migration ledgers: {len(ledgers) if isinstance(ledgers, list) else 0}")
    print("DE.PULSE canonical repository/root/current-state convergence: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
