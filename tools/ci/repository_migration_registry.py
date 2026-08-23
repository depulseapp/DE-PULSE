#!/usr/bin/env python3
"""Compose permanent and work-slice repository migration ledgers.

The permanent registry remains the durable baseline. Individual architecture/process
work slices may retain their exact old->new relocation evidence beside their own
scope/closure metadata. The canonical migration gates consume the composed view, so
splitting evidence by work slice never weakens stale-reference, test-identity, mode,
or root-ownership validation.
"""
from __future__ import annotations

from copy import deepcopy
import json
from pathlib import Path
from typing import Any

ROOT = Path(__file__).resolve().parents[2]
GLOBAL_REGISTRY = ROOT / "governance" / "repository-migrations.json"
WORK_SLICE_ROOT = ROOT / "governance" / "work-slices"
SUPPORTED_WORK_SLICE_SCHEMA = "DE.PULSE-WORK-SLICE-REPOSITORY-MIGRATIONS-1"


class MigrationRegistryError(RuntimeError):
    pass


def _load_object(path: Path) -> dict[str, Any]:
    try:
        value = json.loads(path.read_text(encoding="utf-8"))
    except Exception as exc:
        raise MigrationRegistryError(f"cannot load migration registry {path.relative_to(ROOT)}: {exc}") from exc
    if not isinstance(value, dict):
        raise MigrationRegistryError(f"migration registry must be an object: {path.relative_to(ROOT)}")
    return value


def work_slice_ledgers() -> list[Path]:
    if not WORK_SLICE_ROOT.is_dir():
        return []
    return sorted(path for path in WORK_SLICE_ROOT.glob("*/repository-migrations.json") if path.is_file())


def _expanded_work_slice_moves(path: Path, ledger: dict[str, Any]) -> list[dict[str, Any]]:
    if ledger.get("schema") != SUPPORTED_WORK_SLICE_SCHEMA:
        raise MigrationRegistryError(
            f"unsupported work-slice migration schema in {path.relative_to(ROOT)}: {ledger.get('schema')!r}"
        )
    work_slice_id = str(ledger.get("workSliceId", "")).strip()
    if not work_slice_id or path.parent.name != work_slice_id:
        raise MigrationRegistryError(f"work-slice migration identity/path mismatch: {path.relative_to(ROOT)}")

    defaults = ledger.get("moveDefaults", {})
    if defaults is None:
        defaults = {}
    if not isinstance(defaults, dict):
        raise MigrationRegistryError(f"moveDefaults must be an object: {path.relative_to(ROOT)}")

    rows: list[dict[str, Any]] = []
    direct = ledger.get("moves", [])
    if direct is not None and not isinstance(direct, list):
        raise MigrationRegistryError(f"moves must be a list: {path.relative_to(ROOT)}")
    for row in direct or []:
        if not isinstance(row, dict):
            raise MigrationRegistryError(f"move row must be an object: {path.relative_to(ROOT)}")
        merged = dict(defaults)
        merged.update(row)
        rows.append(merged)

    historical = ledger.get("historicalRootRelocations", [])
    if historical is not None and not isinstance(historical, list):
        raise MigrationRegistryError(f"historicalRootRelocations must be a list: {path.relative_to(ROOT)}")
    for pair in historical or []:
        if not isinstance(pair, dict):
            raise MigrationRegistryError(f"historical relocation must be an object: {path.relative_to(ROOT)}")
        old = str(pair.get("oldPath", "")).strip()
        new = str(pair.get("newPath", "")).strip()
        if not old or not new:
            raise MigrationRegistryError(f"historical relocation requires oldPath/newPath: {path.relative_to(ROOT)}")
        merged = dict(defaults)
        merged.update(
            {
                "oldPath": old,
                "newPath": new,
                "owner": str(pair.get("owner") or defaults.get("owner") or f"{work_slice_id} / historical relocation"),
                "reason": str(pair.get("reason") or defaults.get("reason") or "Move historical current-root evidence to its governed history owner."),
                "removalCondition": str(
                    pair.get("removalCondition")
                    or defaults.get("removalCondition")
                    or "Old root path absent; exact relocated path remains tracked and historical evidence is not a current control owner."
                ),
                "referenceScope": str(pair.get("referenceScope") or defaults.get("referenceScope") or "ACTIVE_RUNTIME_AND_CONTROL"),
                "allowedReferenceFiles": list(pair.get("allowedReferenceFiles") or defaults.get("allowedReferenceFiles") or []),
            }
        )
        rows.append(merged)
    return rows


def load_repository_migrations() -> dict[str, Any]:
    if not GLOBAL_REGISTRY.is_file():
        raise MigrationRegistryError("permanent migration registry missing: governance/repository-migrations.json")
    combined = deepcopy(_load_object(GLOBAL_REGISTRY))
    base = str(combined.get("baselineCommit", "")).strip()
    moves = list(combined.get("moves", [])) if isinstance(combined.get("moves", []), list) else []
    aliases = list(combined.get("temporaryAliases", [])) if isinstance(combined.get("temporaryAliases", []), list) else []
    removed_tests = list(combined.get("removedGoTestIdentities", [])) if isinstance(combined.get("removedGoTestIdentities", []), list) else []
    retired_exec = list(combined.get("retiredExecutablePaths", [])) if isinstance(combined.get("retiredExecutablePaths", []), list) else []
    rename_map = dict(combined.get("testIdentityRenames", {})) if isinstance(combined.get("testIdentityRenames", {}), dict) else {}
    included: list[str] = []

    seen_moves: dict[str, str] = {}
    for row in moves:
        if isinstance(row, dict):
            old = str(row.get("oldPath", "")).strip()
            new = str(row.get("newPath", "")).strip()
            if old and new:
                seen_moves[old] = new

    for path in work_slice_ledgers():
        ledger = _load_object(path)
        ledger_base = str(ledger.get("repositoryMigrationBaselineCommit", base)).strip()
        if base and ledger_base and ledger_base != base:
            raise MigrationRegistryError(
                f"migration baseline drift in {path.relative_to(ROOT)}: {ledger_base} != {base}"
            )
        for row in _expanded_work_slice_moves(path, ledger):
            old = str(row.get("oldPath", "")).strip()
            new = str(row.get("newPath", "")).strip()
            if not old or not new:
                raise MigrationRegistryError(f"composed migration move requires oldPath/newPath: {path.relative_to(ROOT)}")
            prior = seen_moves.get(old)
            if prior and prior != new:
                raise MigrationRegistryError(f"conflicting migration target for {old}: {prior} != {new}")
            if not prior:
                moves.append(row)
                seen_moves[old] = new

        for alias in ledger.get("temporaryAliases", []) if isinstance(ledger.get("temporaryAliases", []), list) else []:
            if alias not in aliases:
                aliases.append(alias)
        for identity in ledger.get("removedGoTestIdentities", []) if isinstance(ledger.get("removedGoTestIdentities", []), list) else []:
            if identity not in removed_tests:
                removed_tests.append(identity)
        for path_value in ledger.get("retiredExecutablePaths", []) if isinstance(ledger.get("retiredExecutablePaths", []), list) else []:
            if path_value not in retired_exec:
                retired_exec.append(path_value)
        local_renames = ledger.get("testIdentityRenames", {})
        if isinstance(local_renames, dict):
            for old, new in local_renames.items():
                if old in rename_map and rename_map[old] != new:
                    raise MigrationRegistryError(f"conflicting test identity rename for {old}")
                rename_map[old] = new
        included.append(path.relative_to(ROOT).as_posix())

    combined["moves"] = moves
    combined["temporaryAliases"] = aliases
    combined["removedGoTestIdentities"] = removed_tests
    combined["retiredExecutablePaths"] = retired_exec
    combined["testIdentityRenames"] = rename_map
    combined["composedWorkSliceLedgers"] = included
    return combined
