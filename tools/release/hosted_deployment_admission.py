#!/usr/bin/env python3
"""Create a source-to-environment admission record for one immutable hosted image."""
from __future__ import annotations

import argparse
import hashlib
import json
import re
from datetime import datetime, timezone
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
DEPENDENCY_POLICY = ROOT / "tools" / "ci" / "ci_dependency_lock.json"
SOURCE_RE = re.compile(r"^[0-9a-f]{40}$")
IMAGE_RE = re.compile(r"^[^\s]+@sha256:([0-9a-f]{64})$")
SCHEMA = "DE.PULSE-HOSTED-DEPLOYMENT-ADMISSION-1"


def digest(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def load_json(path: Path) -> dict:
    value = json.loads(path.read_text(encoding="utf-8"))
    if not isinstance(value, dict):
        raise SystemExit(f"{path}: expected JSON object")
    return value


def require_sbom(path: Path, source_sha: str, image: str, image_digest: str) -> None:
    sbom = load_json(path)
    if sbom.get("spdxVersion") != "SPDX-2.3":
        raise SystemExit("hosted deployment SBOM must be SPDX 2.3")
    packages = [row for row in sbom.get("packages", []) if isinstance(row, dict)]
    roots = [row for row in packages if row.get("name") == "DE.PULSE"]
    if len(roots) != 1 or source_sha not in str(roots[0].get("sourceInfo", "")):
        raise SystemExit("hosted deployment SBOM is not bound to exact source SHA")
    subjects = [row for row in packages if row.get("name") == image]
    if len(subjects) != 1:
        raise SystemExit("hosted deployment SBOM does not identify the immutable image")
    checksums = subjects[0].get("checksums", [])
    if not any(row.get("algorithm") == "SHA256" and str(row.get("checksumValue", "")).lower() == image_digest for row in checksums if isinstance(row, dict)):
        raise SystemExit("hosted deployment SBOM image checksum mismatch")


def require_advisory(metadata_path: Path, report_path: Path, source_sha: str, policy: dict) -> str:
    metadata = {}
    for line in metadata_path.read_text(encoding="utf-8").splitlines():
        key, separator, value = line.partition("=")
        if separator:
            metadata[key.strip()] = value.strip()
    advisory = policy.get("go_supply_chain", {}).get("advisory", {})
    expected_scanner = f"{advisory.get('scanner')}@{advisory.get('version')}"
    if metadata.get("source_sha") != source_sha or metadata.get("scanner") != expected_scanner or metadata.get("database") != advisory.get("database"):
        raise SystemExit("hosted deployment advisory evidence is not bound to exact source/scanner/database")
    report = report_path.read_text(encoding="utf-8")
    if "No vulnerabilities found." not in report:
        raise SystemExit("hosted deployment advisory evidence is not a clean PASS")
    combined = metadata_path.read_bytes() + b"\x00" + report_path.read_bytes()
    return hashlib.sha256(combined).hexdigest()


def require_provenance(path: Path, source_sha: str, image: str, image_digest: str) -> None:
    provenance = load_json(path)
    if provenance.get("schema") != "DE.PULSE-HOSTED-OCI-PROVENANCE-VERIFY-1" or provenance.get("status") != "PASS":
        raise SystemExit("hosted deployment OCI provenance verification is not PASS")
    if provenance.get("sourceSha") != source_sha or provenance.get("image") != image or provenance.get("imageSha256") != image_digest:
        raise SystemExit("hosted deployment OCI provenance subject/source mismatch")


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--source-sha", required=True)
    parser.add_argument("--environment", required=True, choices=("dev", "test", "stage"))
    parser.add_argument("--image", required=True)
    parser.add_argument("--sbom", required=True, type=Path)
    parser.add_argument("--advisory-metadata", required=True, type=Path)
    parser.add_argument("--advisory-report", required=True, type=Path)
    parser.add_argument("--provenance", required=True, type=Path)
    parser.add_argument("--out", required=True, type=Path)
    args = parser.parse_args()
    source_sha = args.source_sha.strip().lower()
    if not SOURCE_RE.fullmatch(source_sha):
        raise SystemExit("hosted deployment source SHA must be exact 40-hex Git identity")
    match = IMAGE_RE.fullmatch(args.image.strip().lower())
    if match is None:
        raise SystemExit("hosted deployment image must be immutable sha256 reference")
    image_digest = match.group(1)
    policy = load_json(DEPENDENCY_POLICY)
    components = policy.get("go_supply_chain", {}).get("components", {})
    revoked = sorted(name for name, row in components.items() if not isinstance(row, dict) or row.get("revoked") is not False)
    if revoked:
        raise SystemExit("revoked/unapproved components block deployment: " + ", ".join(revoked))
    require_sbom(args.sbom, source_sha, args.image, image_digest)
    advisory_digest = require_advisory(args.advisory_metadata, args.advisory_report, source_sha, policy)
    require_provenance(args.provenance, source_sha, args.image, image_digest)
    admission = {
        "schema": SCHEMA,
        "status": "PASS",
        "sourceSha": source_sha,
        "targetEnvironment": args.environment,
        "image": args.image,
        "imageSha256": image_digest,
        "dependencyPolicySha256": digest(DEPENDENCY_POLICY),
        "sbomStatus": "PASS",
        "sbomSha256": digest(args.sbom),
        "advisoryStatus": "PASS",
        "advisoryEvidenceSha256": advisory_digest,
        "provenanceStatus": "PASS",
        "provenanceEvidenceSha256": digest(args.provenance),
        "revokedComponents": [],
        "noRebuild": True,
        "createdAt": datetime.now(timezone.utc).replace(microsecond=0).isoformat().replace("+00:00", "Z"),
    }
    args.out.parent.mkdir(parents=True, exist_ok=True)
    args.out.write_text(json.dumps(admission, indent=2, sort_keys=True) + "\n", encoding="utf-8")
    print(f"HOST-020 hosted deployment admission: PASS source={source_sha} environment={args.environment} image=sha256:{image_digest}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
