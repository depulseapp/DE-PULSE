#!/usr/bin/env python3
"""Fail-closed hosted workload deployment-admission contract."""
from __future__ import annotations

import hashlib
import json
import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
DEPENDENCY_POLICY = ROOT / "tools" / "ci" / "ci_dependency_lock.json"
SCHEMA = "DE.PULSE-HOSTED-DEPLOYMENT-ADMISSION-1"
SHA256_RE = re.compile(r"^[0-9a-f]{64}$")
SOURCE_RE = re.compile(r"^[0-9a-f]{40}$")
IMAGE_RE = re.compile(r"^[^\s]+@sha256:([0-9a-f]{64})$")


def file_sha256(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def validate_digest(value: object, label: str) -> str:
    text = str(value or "").strip().lower()
    if not SHA256_RE.fullmatch(text):
        raise SystemExit(f"hosted deployment admission {label} must be a 64-hex SHA-256")
    return text


def load_deployment_admission(path: Path, environment: str, image: str) -> tuple[dict, str]:
    raw = path.read_bytes()
    try:
        admission = json.loads(raw)
    except (UnicodeDecodeError, json.JSONDecodeError) as exc:
        raise SystemExit(f"hosted deployment admission is invalid JSON: {exc}") from exc
    if not isinstance(admission, dict) or admission.get("schema") != SCHEMA:
        raise SystemExit("unsupported hosted deployment-admission schema")
    if admission.get("status") != "PASS":
        raise SystemExit("hosted deployment admission status is not PASS")
    source_sha = str(admission.get("sourceSha", "")).strip().lower()
    if not SOURCE_RE.fullmatch(source_sha):
        raise SystemExit("hosted deployment admission sourceSha must be exact 40-hex Git identity")
    if admission.get("targetEnvironment") != environment:
        raise SystemExit("hosted deployment admission target environment mismatch")
    image_match = IMAGE_RE.fullmatch(image.strip().lower())
    if image_match is None or admission.get("image") != image:
        raise SystemExit("hosted deployment admission immutable image mismatch")
    if validate_digest(admission.get("imageSha256"), "image digest") != image_match.group(1):
        raise SystemExit("hosted deployment admission image digest does not match image reference")
    for field in ("sbomSha256", "advisoryEvidenceSha256", "provenanceEvidenceSha256"):
        validate_digest(admission.get(field), field)
    if admission.get("sbomStatus") != "PASS":
        raise SystemExit("hosted deployment admission SBOM status is not PASS")
    if admission.get("advisoryStatus") != "PASS":
        raise SystemExit("hosted deployment admission advisory status is not PASS")
    if admission.get("provenanceStatus") != "PASS":
        raise SystemExit("hosted deployment admission provenance status is not PASS")
    if admission.get("noRebuild") is not True:
        raise SystemExit("hosted deployment admission must bind the exact no-rebuild artifact")
    if admission.get("revokedComponents") not in ([], None):
        raise SystemExit("hosted deployment admission contains revoked components")
    expected_policy = file_sha256(DEPENDENCY_POLICY)
    if validate_digest(admission.get("dependencyPolicySha256"), "dependency policy digest") != expected_policy:
        raise SystemExit("hosted deployment admission dependency policy drift")
    return admission, hashlib.sha256(raw).hexdigest()
