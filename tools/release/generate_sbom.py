#!/usr/bin/env python3
from __future__ import annotations

import argparse
import hashlib
import json
import subprocess
from datetime import datetime, timezone
from pathlib import Path
from typing import Any

ROOT = Path(__file__).resolve().parents[2]


def sha256(path: Path) -> str:
    h = hashlib.sha256()
    with path.open("rb") as fh:
        for chunk in iter(lambda: fh.read(1024 * 1024), b""):
            h.update(chunk)
    return h.hexdigest()


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


def spdx_id(value: str) -> str:
    cleaned = "".join(ch if ch.isalnum() or ch in ".-" else "-" for ch in value)
    return "SPDXRef-" + cleaned[:180]


def main() -> int:
    parser = argparse.ArgumentParser(description="Generate DE.PULSE release SPDX 2.3 SBOM")
    parser.add_argument("--version", required=True)
    parser.add_argument("--artifact", action="append", default=[])
    parser.add_argument("--out", required=True)
    args = parser.parse_args()

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
        "downloadLocation": "NOASSERTION",
        "filesAnalyzed": False,
        "licenseConcluded": "NOASSERTION",
        "licenseDeclared": "NOASSERTION",
        "supplier": "Organization: DE.PULSE",
    })

    seen: set[str] = set()
    for mod in go_modules():
        path = str(mod.get("Path", "")).strip()
        version = str(mod.get("Version", "")).strip() or "workspace"
        if not path or path == "depulse":
            continue
        key = f"{path}@{version}"
        if key in seen:
            continue
        seen.add(key)
        pid = spdx_id(key)
        packages.append({
            "name": path,
            "SPDXID": pid,
            "versionInfo": version,
            "downloadLocation": "NOASSERTION",
            "filesAnalyzed": False,
            "licenseConcluded": "NOASSERTION",
            "licenseDeclared": "NOASSERTION",
            "externalRefs": [{
                "referenceCategory": "PACKAGE-MANAGER",
                "referenceType": "purl",
                "referenceLocator": f"pkg:golang/{path}@{version}",
            }],
        })
        relationships.append({"spdxElementId": root_id, "relationshipType": "DEPENDS_ON", "relatedSpdxElement": pid})

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
        })
        relationships.append({"spdxElementId": aid, "relationshipType": "GENERATED_FROM", "relatedSpdxElement": root_id})

    out = {
        "spdxVersion": "SPDX-2.3",
        "dataLicense": "CC0-1.0",
        "SPDXID": "SPDXRef-DOCUMENT",
        "name": f"DE.PULSE-v{args.version}-release-sbom",
        "documentNamespace": f"https://github.com/depulseapp/DE-PULSE/sbom/v{args.version}/{hashlib.sha256(args.version.encode()).hexdigest()[:16]}",
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
    print(f"DE.PULSE SPDX SBOM: PASS packages={len(packages)} artifacts={len(artifacts)} out={target}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
