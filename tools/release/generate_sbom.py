#!/usr/bin/env python3
from __future__ import annotations

import argparse
import hashlib
import json
import re
import subprocess
from datetime import datetime, timezone
from pathlib import Path
from typing import Any

ROOT = Path(__file__).resolve().parents[2]
LOCK_PATH = ROOT / "tools" / "ci" / "ci_dependency_lock.json"


def sha256(path: Path) -> str:
    h = hashlib.sha256()
    with path.open("rb") as fh:
        for chunk in iter(lambda: fh.read(1024 * 1024), b""):
            h.update(chunk)
    return h.hexdigest()


def git_source_sha(explicit: str) -> str:
    value = explicit.strip()
    if not value:
        proc = subprocess.run(
            ["git", "rev-parse", "HEAD"],
            cwd=ROOT,
            check=True,
            text=True,
            capture_output=True,
        )
        value = proc.stdout.strip()
    if not re.fullmatch(r"[0-9a-f]{40}", value):
        raise SystemExit(f"invalid exact source SHA: {value!r}")
    return value


def go_modules() -> list[dict[str, Any]]:
    proc = subprocess.run(
        ["go", "list", "-m", "-json", "all"],
        cwd=ROOT,
        check=True,
        text=True,
        capture_output=True,
    )
    decoder = json.JSONDecoder()
    text = proc.stdout
    pos = 0
    modules: list[dict[str, Any]] = []
    while pos < len(text):
        while pos < len(text) and text[pos].isspace():
            pos += 1
        if pos >= len(text):
            break
        row, pos = decoder.raw_decode(text, pos)
        modules.append(row)
    return modules


def supply_chain_policy() -> dict[str, Any]:
    lock = json.loads(LOCK_PATH.read_text(encoding="utf-8"))
    policy = lock.get("go_supply_chain")
    if not isinstance(policy, dict) or policy.get("schema") != "DE.PULSE-GO-SUPPLY-CHAIN-1":
        raise SystemExit("go supply-chain policy missing or schema invalid")
    if not isinstance(policy.get("components"), dict) or not policy["components"]:
        raise SystemExit("go supply-chain component inventory missing")
    return policy


def spdx_id(value: str) -> str:
    cleaned = "".join(ch if ch.isalnum() or ch in ".-" else "-" for ch in value)
    return "SPDXRef-" + cleaned[:180]


def main() -> int:
    parser = argparse.ArgumentParser(description="Generate DE.PULSE release SPDX 2.3 SBOM")
    parser.add_argument("--version", required=True)
    parser.add_argument("--source-sha", default="", help="Exact 40-hex Git source SHA; defaults to checked-out HEAD")
    parser.add_argument("--artifact", action="append", default=[])
    parser.add_argument("--out", required=True)
    args = parser.parse_args()

    source_sha = git_source_sha(args.source_sha)
    policy = supply_chain_policy()
    governed = policy["components"]
    approved_licenses = set(policy.get("approved_licenses", []))

    artifacts = [Path(x) for x in args.artifact]
    for path in artifacts:
        if not path.is_file():
            raise SystemExit(f"artifact missing: {path}")

    packages: list[dict[str, Any]] = []
    relationships: list[dict[str, str]] = []
    root_id = "SPDXRef-DE-PULSE"
    packages.append({
        "name": "DE.PULSE",
        "SPDXID": root_id,
        "versionInfo": args.version,
        "downloadLocation": f"https://github.com/depulseapp/DE-PULSE/commit/{source_sha}",
        "filesAnalyzed": False,
        "licenseConcluded": "NOASSERTION",
        "licenseDeclared": "NOASSERTION",
        "supplier": "Organization: DE.PULSE",
        "sourceInfo": f"Exact Git source commit: {source_sha}",
    })

    seen: set[str] = set()
    observed: set[str] = set()
    for mod in go_modules():
        path = str(mod.get("Path", "")).strip()
        version = str(mod.get("Version", "")).strip() or "workspace"
        if not path or path == "depulse":
            continue
        key = f"{path}@{version}"
        if key in seen:
            continue
        seen.add(key)
        observed.add(path)
        meta = governed.get(path)
        if not isinstance(meta, dict):
            raise SystemExit(f"ungoverned Go module in SBOM: {key}")
        expected_version = str(meta.get("version", "")).strip()
        if version != expected_version:
            raise SystemExit(f"Go module version drift for {path}: observed={version} governed={expected_version}")
        if meta.get("revoked") is not False:
            raise SystemExit(f"revoked or unapproved Go module blocks SBOM: {key}")
        license_id = str(meta.get("license", "")).strip()
        if not license_id or license_id not in approved_licenses:
            raise SystemExit(f"missing/unapproved Go module license for {key}: {license_id!r}")
        source = str(meta.get("source", "")).strip()
        if not source.startswith("https://"):
            raise SystemExit(f"missing governed HTTPS source for {key}")
        pid = spdx_id(key)
        packages.append({
            "name": path,
            "SPDXID": pid,
            "versionInfo": version,
            "downloadLocation": source,
            "filesAnalyzed": False,
            "licenseConcluded": license_id,
            "licenseDeclared": license_id,
            "externalRefs": [{
                "referenceCategory": "PACKAGE-MANAGER",
                "referenceType": "purl",
                "referenceLocator": f"pkg:golang/{path}@{version}",
            }],
        })
        relationships.append({"spdxElementId": root_id, "relationshipType": "DEPENDS_ON", "relatedSpdxElement": pid})

    missing = sorted(set(governed) - observed)
    if missing:
        raise SystemExit("governed Go modules absent from resolved module graph: " + ", ".join(missing))

    for artifact in artifacts:
        aid = spdx_id("artifact-" + artifact.name)
        packages.append({
            "name": artifact.name,
            "SPDXID": aid,
            "versionInfo": args.version,
            "downloadLocation": "NOASSERTION",
            "filesAnalyzed": False,
            "licenseConcluded": "NOASSERTION",
            "licenseDeclared": "NOASSERTION",
            "checksums": [{"algorithm": "SHA256", "checksumValue": sha256(artifact)}],
            "primaryPackagePurpose": "APPLICATION",
            "sourceInfo": f"Built from exact Git source commit: {source_sha}",
        })
        relationships.append({"spdxElementId": aid, "relationshipType": "GENERATED_FROM", "relatedSpdxElement": root_id})

    namespace_seed = f"{args.version}:{source_sha}"
    out = {
        "spdxVersion": "SPDX-2.3",
        "dataLicense": "CC0-1.0",
        "SPDXID": "SPDXRef-DOCUMENT",
        "name": f"DE.PULSE-v{args.version}-release-sbom",
        "documentNamespace": f"https://github.com/depulseapp/DE-PULSE/sbom/v{args.version}/{source_sha}/{hashlib.sha256(namespace_seed.encode()).hexdigest()[:16]}",
        "documentDescribes": [root_id],
        "creationInfo": {
            "created": datetime.now(timezone.utc).replace(microsecond=0).isoformat().replace("+00:00", "Z"),
            "creators": ["Tool: DE.PULSE tools/release/generate_sbom.py"],
        },
        "packages": packages,
        "relationships": relationships,
    }
    target = Path(args.out)
    target.parent.mkdir(parents=True, exist_ok=True)
    target.write_text(json.dumps(out, indent=2, sort_keys=True) + "\n", encoding="utf-8")
    print(
        f"DE.PULSE SPDX SBOM: PASS source={source_sha} packages={len(packages)} "
        f"artifacts={len(artifacts)} out={target}"
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
