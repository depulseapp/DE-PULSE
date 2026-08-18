#!/usr/bin/env python3
"""Canonical DE.PULSE source fingerprinting.

Two fingerprint modes are intentionally provided:

- ``canonical_source_fingerprint`` hashes materialized filesystem bytes. It is
  retained for local diagnostics and backwards-compatible callers.
- ``canonical_git_object_fingerprint`` hashes the raw Git blob bytes at an
  immutable commit. This is the authoritative cross-platform provenance method
  for CI/release certification because checkout/materialization rules can alter
  filesystem bytes on Windows even when the Git commit is identical.

Both modes normalize relative paths to POSIX separators and use the same
exclusion semantics. Git blob reads use one persistent ``git cat-file --batch``
process so provenance stays deterministic without per-file process overhead.
"""
from __future__ import annotations

from pathlib import Path, PurePath, PurePosixPath
import hashlib
import subprocess
from typing import Callable, Iterable

ExcludeFile = Callable[[Path, PurePath], bool]


def canonical_rel_key(path: PurePath) -> str:
    return "/".join(path.parts)


def _suffixes(values: Iterable[str]) -> set[str]:
    return {str(x).lower() for x in values}


def canonical_source_fingerprint(
    root: Path,
    *,
    excluded_dirs=frozenset(),
    excluded_suffixes=frozenset(),
    exclude_file: ExcludeFile | None = None,
) -> str:
    """Hash materialized filesystem bytes.

    This remains useful as a diagnostic, but native/release provenance should
    prefer :func:`canonical_git_object_fingerprint`.
    """
    root = Path(root).resolve()
    candidates = [p for p in root.rglob("*") if p.is_file()]
    candidates.sort(key=lambda p: canonical_rel_key(p.relative_to(root)))
    h = hashlib.sha256()
    suffixes = _suffixes(excluded_suffixes)
    for p in candidates:
        rel = p.relative_to(root)
        if any(part in excluded_dirs for part in rel.parts):
            continue
        if p.suffix.lower() in suffixes:
            continue
        if exclude_file is not None and exclude_file(p, rel):
            continue
        h.update(canonical_rel_key(rel).encode("utf-8"))
        h.update(b"\0")
        h.update(hashlib.sha256(p.read_bytes()).digest())
        h.update(b"\0")
    return h.hexdigest()


def _git_blob_hashes(root: Path, items: list[tuple[str, str]]) -> list[tuple[str, bytes]]:
    """Return SHA-256 digests of Git blobs using one persistent Git process."""
    proc = subprocess.Popen(
        ["git", "-C", str(root), "cat-file", "--batch"],
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
    )
    assert proc.stdin is not None
    assert proc.stdout is not None
    digests: list[tuple[str, bytes]] = []
    try:
        for rel, object_id in items:
            proc.stdin.write(object_id.encode("ascii") + b"\n")
            proc.stdin.flush()
            header = proc.stdout.readline()
            if not header:
                raise RuntimeError(f"git cat-file ended while reading {rel}")
            parts = header.rstrip(b"\n").split()
            if len(parts) != 3 or parts[1] != b"blob":
                raise RuntimeError(f"unexpected git cat-file header for {rel}: {header!r}")
            size = int(parts[2])
            data = proc.stdout.read(size)
            if len(data) != size:
                raise RuntimeError(f"short Git blob read for {rel}: {len(data)} != {size}")
            terminator = proc.stdout.read(1)
            if terminator != b"\n":
                raise RuntimeError(f"missing Git blob terminator for {rel}")
            digests.append((rel, hashlib.sha256(data).digest()))
    finally:
        proc.stdin.close()
        return_code = proc.wait()
        if return_code != 0:
            stderr = b"" if proc.stderr is None else proc.stderr.read()
            raise RuntimeError(f"git cat-file --batch failed ({return_code}): {stderr.decode(errors='replace')}")
    return digests


def canonical_git_object_fingerprint(
    root: Path,
    commit: str = "HEAD",
    *,
    excluded_dirs=frozenset(),
    excluded_suffixes=frozenset(),
) -> str:
    """Hash canonical raw Git blob bytes for ``commit``.

    This is the authoritative CI/release source identity. It deliberately does
    not read the checked-out file contents, so line-ending conversion, file mode
    materialization and other OS checkout behavior cannot change provenance.
    """
    root = Path(root).resolve()
    suffixes = _suffixes(excluded_suffixes)
    raw = subprocess.check_output(
        ["git", "-C", str(root), "ls-tree", "-r", "-z", "--full-tree", commit]
    )
    objects: list[tuple[str, str]] = []
    for entry in raw.split(b"\0"):
        if not entry:
            continue
        meta, path_bytes = entry.split(b"\t", 1)
        _mode, kind, object_id = meta.decode("ascii").split()
        if kind != "blob":
            continue
        rel = path_bytes.decode("utf-8")
        pure = PurePosixPath(rel)
        if any(part in excluded_dirs for part in pure.parts):
            continue
        if pure.suffix.lower() in suffixes:
            continue
        objects.append((rel, object_id))

    objects.sort(key=lambda item: item[0])
    h = hashlib.sha256()
    for rel, blob_digest in _git_blob_hashes(root, objects):
        h.update(rel.encode("utf-8"))
        h.update(b"\0")
        h.update(blob_digest)
        h.update(b"\0")
    return h.hexdigest()


if __name__ == "__main__":
    import argparse

    parser = argparse.ArgumentParser(description="Compute DE.PULSE source fingerprints")
    parser.add_argument("--root", default=".")
    parser.add_argument("--commit", default="HEAD")
    parser.add_argument(
        "--mode",
        choices=("git", "filesystem"),
        default="git",
        help="git is authoritative for cross-platform CI/release provenance",
    )
    args = parser.parse_args()
    excluded_dirs = {".git", "__pycache__", ".depulse-certification"}
    excluded_suffixes = {".log", ".tmp", ".out", ".exe"}
    if args.mode == "git":
        print(
            canonical_git_object_fingerprint(
                Path(args.root),
                args.commit,
                excluded_dirs=excluded_dirs,
                excluded_suffixes=excluded_suffixes,
            )
        )
    else:
        print(
            canonical_source_fingerprint(
                Path(args.root),
                excluded_dirs=excluded_dirs,
                excluded_suffixes=excluded_suffixes,
            )
        )
