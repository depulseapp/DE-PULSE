#!/usr/bin/env python3
"""Canonical version-neutral DE.PULSE G12 certification executor.

Release-specific requirements live in a declarative certification manifest under
release/v<productVersion>/. Historical shipped executors remain immutable in
Git history/tags, but future releases use this single orchestration owner.
"""
from __future__ import annotations

import argparse
import hashlib
import importlib.metadata
import importlib.util
import json
import os
from pathlib import Path
import shutil
import subprocess
import sys
from datetime import datetime, timezone
from typing import Any

ROOT = Path(__file__).resolve().parents[2]
REQUIRED_TOOLS = ("git", "go", "gofmt", "node", "python3")
TOOLCHAIN_MANIFEST = ROOT / "governance" / "toolchain-manifest.json"


def utc_now() -> str:
    return datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")


def git(*args: str) -> str:
    return subprocess.check_output(("git", *args), cwd=ROOT, text=True).strip()


def load_json(path: Path) -> dict[str, Any]:
    return json.loads(path.read_text(encoding="utf-8"))


def resolve_json_path(payload: Any, dotted: str) -> Any:
    current = payload
    for part in dotted.split("."):
        if isinstance(current, dict) and part in current:
            current = current[part]
        else:
            raise KeyError(dotted)
    return current


def write_line(log, text: str = "") -> None:
    print(text, flush=True)
    log.write(text + "\n")
    log.flush()


def run(log, argv: list[str], *, env: dict[str, str] | None = None) -> None:
    if not argv or not all(isinstance(item, str) and item for item in argv):
        raise ValueError(f"invalid manifest command: {argv!r}")
    write_line(log, "+ " + " ".join(argv))
    proc = subprocess.Popen(
        argv,
        cwd=ROOT,
        env=env,
        stdout=subprocess.PIPE,
        stderr=subprocess.STDOUT,
        text=True,
        bufsize=1,
    )
    assert proc.stdout is not None
    for line in proc.stdout:
        print(line, end="", flush=True)
        log.write(line)
    rc = proc.wait()
    log.flush()
    if rc != 0:
        raise subprocess.CalledProcessError(rc, argv)


def resolved_toolchain() -> dict[str, Any]:
    manifest = load_json(TOOLCHAIN_MANIFEST)
    resolved = {
        "go version": subprocess.check_output(["go", "version"], cwd=ROOT, text=True).strip(),
        "node --version": subprocess.check_output(["node", "--version"], cwd=ROOT, text=True).strip(),
        "python3 --version": subprocess.check_output(["python3", "--version"], cwd=ROOT, text=True).strip(),
        "playwright package version": importlib.metadata.version("playwright"),
    }
    runner = {
        "RUNNER_OS": os.environ.get("RUNNER_OS", "local"),
        "ImageOS": os.environ.get("ImageOS", "local"),
        "ImageVersion": os.environ.get("ImageVersion", "local"),
        "RUNNER_ARCH": os.environ.get("RUNNER_ARCH", "local"),
    }
    return {
        "manifest": "governance/toolchain-manifest.json",
        "requested": manifest,
        "resolved": resolved,
        "runner": runner,
    }


def validate_resolved_toolchain(toolchain: dict[str, Any]) -> None:
    requested = toolchain["requested"]
    resolved = toolchain["resolved"]
    expected = {
        "go version": str(requested["go"]["version"]),
        "node --version": str(requested["node"]["version"]),
        "python3 --version": str(requested["python"]["version"]),
        "playwright package version": str(requested["playwright"]["version"]),
    }
    checks = {
        "go version": f"go{expected['go version']}",
        "node --version": f"v{expected['node --version']}",
        "python3 --version": f"Python {expected['python3 --version']}",
        "playwright package version": expected["playwright package version"],
    }
    for key, token in checks.items():
        value = str(resolved[key])
        if key == "playwright package version":
            ok = value == token
        else:
            ok = token in value
        if not ok:
            raise AssertionError(f"resolved toolchain mismatch for {key}: actual={value!r} expected token={token!r}")


def validate_manifest(identity: dict[str, Any], manifest: dict[str, Any], manifest_path: Path) -> None:
    if manifest.get("schema") != "DE.PULSE-G12-EVIDENCE-MANIFEST-1":
        raise ValueError(f"unsupported G12 manifest schema: {manifest.get('schema')!r}")
    version = str(identity.get("version", ""))
    if manifest.get("productVersion") != version:
        raise ValueError(
            f"manifest productVersion {manifest.get('productVersion')!r} does not match release identity {version!r}"
        )
    if not manifest.get("workSliceId"):
        raise ValueError("manifest workSliceId is required")
    if not isinstance(manifest.get("evidenceSchemaVersion"), int):
        raise ValueError("manifest evidenceSchemaVersion must be an integer")
    expected = ROOT / "release" / f"v{version}" / "certification-manifest.json"
    if manifest_path.resolve() != expected.resolve():
        raise ValueError(f"canonical manifest path must be {expected.relative_to(ROOT)}")
    contract = ROOT / str(manifest.get("releaseContract", ""))
    if not contract.is_file():
        raise ValueError(f"release contract missing: {contract.relative_to(ROOT)}")


def validate_assertions(log, manifest: dict[str, Any]) -> None:
    cache: dict[Path, dict[str, Any]] = {}
    for item in manifest.get("assertions", []):
        path = ROOT / item["file"]
        payload = cache.setdefault(path, load_json(path))
        value = resolve_json_path(payload, item["path"])
        if "equals" in item and value != item["equals"]:
            raise AssertionError(f"{item['file']}:{item['path']}={value!r} != {item['equals']!r}")
        if "contains" in item and str(item["contains"]) not in str(value):
            raise AssertionError(f"{item['file']}:{item['path']} does not contain {item['contains']!r}")
    write_line(log, "PASS: declarative release/identity/scope assertions")

    for item in manifest.get("staticSourceAssertions", []):
        path = ROOT / item["file"]
        text = path.read_text(encoding="utf-8")
        for token in item.get("present", []):
            if token not in text:
                raise AssertionError(f"{item['file']} missing required source token {token!r}")
        for token in item.get("absent", []):
            if token in text:
                raise AssertionError(f"{item['file']} contains forbidden source token {token!r}")
    write_line(log, "PASS: declarative static source assertions")


def chrome_available() -> bool:
    explicit = os.environ.get("CHROME_BIN", "").strip()
    if explicit and Path(explicit).exists():
        return True
    for candidate in (
        "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
        shutil.which("google-chrome") or "",
        shutil.which("chromium") or "",
    ):
        if candidate and Path(candidate).exists():
            os.environ["CHROME_BIN"] = candidate
            return True
    return False


def ensure_python_module(name: str) -> None:
    if importlib.util.find_spec(name) is None:
        raise RuntimeError(f"required Python module is missing: {name}")


def main() -> int:
    parser = argparse.ArgumentParser(description="Run canonical DE.PULSE exact-source G12 certification")
    parser.add_argument("--manifest", default="")
    args = parser.parse_args()

    os.chdir(ROOT)
    for tool in REQUIRED_TOOLS:
        if shutil.which(tool) is None:
            raise SystemExit(f"ERROR: required tool is missing: {tool}")

    source_sha = git("rev-parse", "HEAD")
    source_branch = git("branch", "--show-current")
    expected_sha = os.environ.get("DEPULSE_EXPECTED_SHA", "").strip()
    if expected_sha and source_sha != expected_sha:
        raise SystemExit(f"ERROR: expected {expected_sha} but checkout is {source_sha}")
    dirty = git("status", "--porcelain", "--untracked-files=normal")
    if dirty:
        print("ERROR: certification requires a clean exact-source checkout.", file=sys.stderr)
        print(dirty, file=sys.stderr)
        return 2

    identity = load_json(ROOT / "release_identity.json")
    version = str(identity.get("version", ""))
    manifest_path = Path(args.manifest) if args.manifest else Path(f"release/v{version}/certification-manifest.json")
    if not manifest_path.is_absolute():
        manifest_path = ROOT / manifest_path
    if not manifest_path.is_file():
        raise SystemExit(f"ERROR: canonical G12 manifest missing: {manifest_path}")
    manifest = load_json(manifest_path)
    validate_manifest(identity, manifest, manifest_path)

    ensure_python_module("playwright")
    toolchain = resolved_toolchain()
    validate_resolved_toolchain(toolchain)

    evidence_root = Path(
        os.environ.get(
            "DEPULSE_EVIDENCE_DIR",
            str(ROOT / ".depulse-certification" / f"v{version}" / source_sha),
        )
    )
    evidence_root.mkdir(parents=True, exist_ok=True)
    log_file = evidence_root / "certification.log"
    result_file = evidence_root / "certification-result.json"
    started_at = utc_now()

    with log_file.open("w", encoding="utf-8") as log:
        write_line(log, f"DE.PULSE v{version} canonical G12 exact-source certification")
        write_line(log, f"Source SHA: {source_sha}")
        write_line(log, f"Branch: {source_branch}")
        write_line(log, f"Work slice: {manifest['workSliceId']}")
        write_line(log, f"Evidence schema: {manifest['evidenceSchemaVersion']}")
        write_line(log, f"Started: {started_at}")
        write_line(log, "Resolved toolchain: " + json.dumps(toolchain["resolved"], sort_keys=True))
        write_line(log, "Runner identity: " + json.dumps(toolchain["runner"], sort_keys=True))

        write_line(log, "\n[G0/G1] Canonical release identity + declarative release evidence")
        run(log, ["python3", "release_identity.py", "--verify"])
        run(log, ["python3", "tools/release/release_identity_contract.py", "--verify"])
        run(log, ["python3", "version_consistency_test.py"])
        validate_assertions(log, manifest)

        write_line(log, "\n[G2/G7/G10] Governance, reproducibility, portability and inherited trust contracts")
        for command in manifest.get("pythonGates", []):
            run(log, list(command))

        focused = str(manifest.get("focusedGoPattern", "")).strip()
        if focused:
            write_line(log, "\n[G2/G7/G8/G10/G12] Release-focused Go trust coverage")
            run(log, ["go", "test", "-count=1", "-run", focused, "./..."])

        write_line(log, "\n[G7/G8/G10/G12] Go format, vet, full suite, race and randomized order")
        tracked_go = git("ls-files", "*.go").splitlines()
        if tracked_go:
            result = subprocess.run(["gofmt", "-l", *tracked_go], cwd=ROOT, text=True, capture_output=True, check=True)
            if result.stdout.strip():
                raise AssertionError("gofmt drift:\n" + result.stdout)
        write_line(log, "PASS: gofmt")
        run(log, ["go", "vet", "./..."])
        run(log, ["go", "test", "-count=1", "./..."])
        run(log, ["go", "test", "-race", "-count=1", "./..."])
        run(log, ["go", "test", "-shuffle=on", "-count=1", "./..."])

        write_line(log, "\n[G9/G10/G12] Renderer syntax and declarative Node regressions")
        for js in sorted((ROOT / "renderer").rglob("*.js")):
            run(log, ["node", "--check", str(js.relative_to(ROOT))])
        for command in manifest.get("nodeTests", []):
            run(log, list(command))

        chrome_tests = manifest.get("chromeTests", [])
        if chrome_tests:
            write_line(log, "\n[G9/G10/G12] Current primary Chrome behavior")
            if not chrome_available():
                raise RuntimeError("set CHROME_BIN to an installed Chrome/Chromium executable")
            for command in chrome_tests:
                run(log, list(command), env=os.environ.copy())

        write_line(log, "\n[G10/G12] Pre-merge evidence ownership")
        premerge = manifest.get("premergeEvidence", {})
        for key, value in sorted(premerge.items()):
            write_line(log, f"{key}: {value}")
        write_line(log, "G11 independently requires exact-head DE.PULSE/fast-head + DE.PULSE/qualified-head success.")

        completed_at = utc_now()
        log.flush()
        log_sha = hashlib.sha256(log_file.read_bytes()).hexdigest()
        source_fingerprint = subprocess.check_output(
            ["python3", "source_fingerprint.py", "--mode", "git", "--commit", source_sha],
            cwd=ROOT,
            text=True,
        ).strip()
        result = {
            "schema": "DE.PULSE-G12-CERTIFICATION-2",
            "evidenceSchemaVersion": manifest["evidenceSchemaVersion"],
            "productVersion": version,
            "workSliceId": manifest["workSliceId"],
            "buildId": identity.get("build_id"),
            "platformBuildNumber": identity.get("bundle_version"),
            "sourceSha": source_sha,
            "sourceFingerprint": source_fingerprint,
            "sourceBranch": source_branch,
            "startedAtUtc": started_at,
            "completedAtUtc": completed_at,
            "result": "PASS",
            "lane": "CANONICAL_VERSION_NEUTRAL_EXACT_SOURCE_FULL_CERTIFICATION",
            "manifest": str(manifest_path.relative_to(ROOT)),
            "releaseContract": manifest["releaseContract"],
            "logSha256": log_sha,
            "resolvedToolchain": toolchain,
            "goFull": "PASS",
            "race": "PASS",
            "randomized": "PASS",
            "renderer": "PASS",
            "chrome": "PASS" if chrome_tests else "NOT_SELECTED",
            "premergeEvidence": premerge,
            "nativePackagingStatus": "REQUIRED_BEFORE_G15_PROMOTION",
            "protectedBoundaries": manifest.get("protectedBoundaries", []),
        }
        result_file.write_text(json.dumps(result, indent=2, sort_keys=True) + "\n", encoding="utf-8")
        write_line(log, f"\nPASS: canonical v{version} G12 certification completed.")
        write_line(log, f"Evidence: {result_file.relative_to(ROOT)}")
        write_line(log, "Next: G13/G14 native macOS Apple Silicon and Windows x64 package/runtime audit.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
