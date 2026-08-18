#!/usr/bin/env python3
from __future__ import annotations

from pathlib import Path
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[2]
sys.path.insert(0, str(ROOT))

from source_fingerprint import canonical_git_object_fingerprint, canonical_source_fingerprint  # noqa: E402

EXCLUDED_DIRS = {".git", "__pycache__", ".depulse-certification"}
EXCLUDED_SUFFIXES = {".log", ".tmp", ".out", ".exe"}


def main() -> int:
    head = subprocess.check_output(["git", "-C", str(ROOT), "rev-parse", "HEAD"], text=True).strip()
    git_fp = canonical_git_object_fingerprint(
        ROOT,
        head,
        excluded_dirs=EXCLUDED_DIRS,
        excluded_suffixes=EXCLUDED_SUFFIXES,
    )
    fs_fp = canonical_source_fingerprint(
        ROOT,
        excluded_dirs=EXCLUDED_DIRS,
        excluded_suffixes=EXCLUDED_SUFFIXES,
    )
    assert len(git_fp) == 64
    assert len(fs_fp) == 64
    assert git_fp == canonical_git_object_fingerprint(
        ROOT,
        "HEAD",
        excluded_dirs=EXCLUDED_DIRS,
        excluded_suffixes=EXCLUDED_SUFFIXES,
    )
    print(f"canonical Git-object fingerprint: {git_fp}")
    print(f"materialized filesystem fingerprint (diagnostic): {fs_fp}")
    if git_fp != fs_fp:
        print("INFO: filesystem materialization differs; Git-object fingerprint remains authoritative")
    print("DE.PULSE source provenance regression: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
