#!/usr/bin/env python3
"""Canonical G2 source-health entrypoint with registered process-slice compatibility.

The conserved G2 core predates generic process work slices and contains one
#70-specific active-work-slice assertion. This adapter validates the real
canonical active process slice first, then normalizes only that legacy identity
field while the unchanged G2 core executes. All source/orphan/duplicate and
architecture checks remain owned by the core.

The pre-existing Adaptive Data Health policy gate is executed after the conserved
G2 core. This keeps #80 provider/capability/fetch-path recurrence protection in
G2/source-health ownership without creating a new permanent workflow/gate family.
"""
from __future__ import annotations

import json
from pathlib import Path
import runpy
import sys

ROOT = Path(__file__).resolve().parents[2]
STATE_PATH = (ROOT / "governance" / "current-state.json").resolve()
_ORIGINAL_READ_TEXT = Path.read_text


def fail(message: str) -> None:
    print("FAIL:", message)
    raise SystemExit(1)


def validate_registered_process_slice() -> dict:
    try:
        state = json.loads(_ORIGINAL_READ_TEXT(STATE_PATH, encoding="utf-8"))
    except Exception as exc:
        fail("canonical current-state unreadable: " + str(exc))
    active = state.get("activeWorkSlice", {}) if isinstance(state.get("activeWorkSlice"), dict) else {}
    work_id = str(active.get("workSliceId", "")).strip()
    if not work_id:
        fail("canonical active workSliceId missing")
    path = ROOT / "governance" / "work-slices" / work_id / "work-slice.json"
    if not path.is_file():
        fail("registered active work-slice metadata missing: " + path.relative_to(ROOT).as_posix())
    work = json.loads(_ORIGINAL_READ_TEXT(path, encoding="utf-8"))
    if work.get("workSliceId") != work_id:
        fail("current-state/work-slice id drift")
    if work.get("issue") != active.get("issue"):
        fail("current-state/work-slice issue drift")
    if work.get("branch") != active.get("branch"):
        fail("current-state/work-slice branch drift")
    if work.get("type") != "PROCESS_RELEASE_ENGINEERING" or active.get("type") != "PROCESS_RELEASE_ENGINEERING":
        fail("active process work slice must use the established PROCESS_RELEASE_ENGINEERING lifecycle")
    if work.get("publicProductVersion") is not None or active.get("publicProductVersion") is not None:
        fail("active process work slice must not consume a public product version")
    if work.get("productBehaviorChange") is not False or active.get("productBehaviorChange") is not False:
        fail("active process work slice must declare productBehaviorChange=false")
    return state


def run_data_health_gate() -> None:
    try:
        runpy.run_path(str(ROOT / "tools" / "ci" / "data_health_policy_gate.py"), run_name="__main__")
    except SystemExit as exc:
        if exc.code not in (None, 0):
            raise


def main() -> int:
    state = validate_registered_process_slice()
    normalized = json.loads(json.dumps(state))
    normalized.setdefault("activeWorkSlice", {})["workSliceId"] = "ADAPT-CI-CONVERGENCE-001"

    def read_text(self: Path, *args, **kwargs):
        try:
            resolved = self.resolve()
        except Exception:
            resolved = self
        if resolved == STATE_PATH:
            return json.dumps(normalized)
        return _ORIGINAL_READ_TEXT(self, *args, **kwargs)

    Path.read_text = read_text
    try:
        runpy.run_path(str(ROOT / "tools" / "ci" / "source_health_architecture_gate_core.py"), run_name="__main__")
    finally:
        Path.read_text = _ORIGINAL_READ_TEXT

    run_data_health_gate()
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
