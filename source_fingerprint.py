#!/usr/bin/env python3
"""Canonical cross-platform DE.PULSE source fingerprinting.

Relative paths are normalized to POSIX separators and sorted by the same
canonical path key on every host. This makes one immutable source tree produce
the same fingerprint on Linux, macOS and Windows.
"""
from __future__ import annotations
from pathlib import Path, PurePath
import hashlib
from typing import Callable

ExcludeFile = Callable[[Path, PurePath], bool]

def canonical_rel_key(path: PurePath) -> str:
    return "/".join(path.parts)

def canonical_source_fingerprint(root: Path, *, excluded_dirs=frozenset(), excluded_suffixes=frozenset(), exclude_file: ExcludeFile | None = None) -> str:
    root = Path(root).resolve()
    candidates = [p for p in root.rglob("*") if p.is_file()]
    candidates.sort(key=lambda p: canonical_rel_key(p.relative_to(root)))
    h = hashlib.sha256()
    suffixes = {str(x).lower() for x in excluded_suffixes}
    for p in candidates:
        rel = p.relative_to(root)
        if any(part in excluded_dirs for part in rel.parts):
            continue
        if p.suffix.lower() in suffixes:
            continue
        if exclude_file is not None and exclude_file(p, rel):
            continue
        h.update(canonical_rel_key(rel).encode("utf-8")); h.update(b"\0")
        h.update(hashlib.sha256(p.read_bytes()).digest()); h.update(b"\0")
    return h.hexdigest()
