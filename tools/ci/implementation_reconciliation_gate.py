#!/usr/bin/env python3
"""Archive-aware entrypoint for the conserved v17/v18 reconciliation gate.

#73 moved historical version-scoped non-Go root evidence to governed release
history without rewriting the bytes. The conserved reconciliation core still
uses the historical root-relative identities embedded in its immutable ledger,
so this entrypoint resolves only missing root-level v17/v18 evidence paths to
their byte-preserved archive owners while the core executes.
"""
from __future__ import annotations

from pathlib import Path
import importlib
import sys

ROOT = Path(__file__).resolve().parents[2]
_ORIG_READ_TEXT = Path.read_text
_ORIG_EXISTS = Path.exists
_ORIG_IS_FILE = Path.is_file

_PREFIX_DIRS = (
    ("v17_5_1_", "release/history/v17.5.1/root-evidence"),
    ("v17_5_", "release/history/v17.5.0/root-evidence"),
    ("v17_4_", "release/history/v17.4.0/root-evidence"),
    ("v17_3_", "release/history/v17.3.0/root-evidence"),
    ("v17_2_", "release/history/v17.2.0/root-evidence"),
    ("v17_1_", "release/history/v17.1.0/root-evidence"),
    ("v17_0_", "release/history/v17.0.0/root-evidence"),
    ("v18_0_6_", "release/history/v18.0.6/root-evidence"),
    ("v18_0_5_", "release/history/v18.0.5/root-evidence"),
    ("v18_0_4_", "release/history/v18.0.4/root-evidence"),
    ("v18_0_3_", "release/history/v18.0.3/root-evidence"),
    ("v18_0_2_", "release/history/v18.0.2/root-evidence"),
    ("v18_0_1_", "release/history/v18.0.1/root-evidence"),
    ("v18_0_", "release/history/v18.0.0/root-evidence"),
    ("v18_1_", "release/history/v18.1.0/root-evidence"),
    ("v18_2_", "release/history/v18.2.0/root-evidence"),
    ("v18_3_", "release/history/v18.3.0/root-evidence"),
    ("v18_4_", "release/history/v18.4.0/root-evidence"),
    ("v18_5_", "release/history/v18.5.0/root-evidence"),
    ("v18_7_0_", "release/history/v18.7.0/root-evidence"),
    ("v18_8_0_", "release/history/v18.8.0/root-evidence"),
    ("v18_9_1_", "release/v18.9.1/legacy-root"),
)
_MAJOR_FILES = {
    "v17_baseline_contract.json", "v17_baseline_gate.py", "v17_delivery_slices.json", "v17_major_scope.json",
    "v18_baseline_contract.json", "v18_baseline_gate.py", "v18_delivery_slices.json", "v18_documentation_typography_gate.py", "v18_major_scope.json",
}


def archived(path: Path) -> Path | None:
    try:
        if path.parent != ROOT:
            return None
    except Exception:
        return None
    name = path.name
    if name.startswith("v17_") and name in _MAJOR_FILES:
        candidate = ROOT / "release/history/v17-major/root-evidence" / name
        return candidate if _ORIG_IS_FILE(candidate) else None
    if name.startswith("v18_") and name in _MAJOR_FILES:
        candidate = ROOT / "release/history/v18-major/root-evidence" / name
        return candidate if _ORIG_IS_FILE(candidate) else None
    for prefix, directory in _PREFIX_DIRS:
        if name.startswith(prefix):
            candidate = ROOT / directory / name
            return candidate if _ORIG_IS_FILE(candidate) else None
    return None


def read_text(self: Path, *args, **kwargs):
    if _ORIG_IS_FILE(self):
        return _ORIG_READ_TEXT(self, *args, **kwargs)
    alternate = archived(self)
    if alternate is not None:
        return _ORIG_READ_TEXT(alternate, *args, **kwargs)
    return _ORIG_READ_TEXT(self, *args, **kwargs)


def exists(self: Path) -> bool:
    return _ORIG_EXISTS(self) or archived(self) is not None


def is_file(self: Path) -> bool:
    return _ORIG_IS_FILE(self) or archived(self) is not None


def main() -> int:
    Path.read_text = read_text
    Path.exists = exists
    Path.is_file = is_file
    try:
        core = importlib.import_module("implementation_reconciliation_gate_core")
        return int(core.main())
    finally:
        Path.read_text = _ORIG_READ_TEXT
        Path.exists = _ORIG_EXISTS
        Path.is_file = _ORIG_IS_FILE


if __name__ == "__main__":
    raise SystemExit(main())
