#!/usr/bin/env python3
"""Offline positive/adverse evidence for HOST-020 deployment admission."""
from __future__ import annotations

import importlib.util
import json
import subprocess
import sys
import tempfile
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
GENERATOR = ROOT / "tools" / "release" / "hosted_deployment_admission.py"
LOADER = ROOT / "tools" / "hosted" / "deployment_admission.py"


def load_contract():
    spec = importlib.util.spec_from_file_location("depulse_deployment_admission", LOADER)
    if spec is None or spec.loader is None:
        raise SystemExit("cannot import deployment admission contract")
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def main() -> None:
    source = "b" * 40
    image_digest = "a" * 64
    image = "ghcr.io/depulseapp/de-pulse@sha256:" + image_digest
    with tempfile.TemporaryDirectory() as tmp_raw:
        tmp = Path(tmp_raw)
        sbom = tmp / "sbom.json"
        metadata = tmp / "govuln-metadata.txt"
        report = tmp / "govuln.txt"
        provenance = tmp / "provenance.json"
        admission = tmp / "admission.json"
        sbom.write_text(json.dumps({
            "spdxVersion": "SPDX-2.3",
            "packages": [
                {"name": "DE.PULSE", "sourceInfo": "Exact Git source commit: " + source},
                {"name": image, "checksums": [{"algorithm": "SHA256", "checksumValue": image_digest}]},
            ],
        }), encoding="utf-8")
        metadata.write_text(f"source_sha={source}\nscanner=golang.org/x/vuln/cmd/govulncheck@v1.1.4\ndatabase=https://vuln.go.dev\n", encoding="utf-8")
        report.write_text("No vulnerabilities found.\n", encoding="utf-8")
        provenance.write_text(json.dumps({"schema": "DE.PULSE-HOSTED-OCI-PROVENANCE-VERIFY-1", "status": "PASS", "sourceSha": source, "image": image, "imageSha256": image_digest}), encoding="utf-8")
        command = [sys.executable, str(GENERATOR), "--source-sha", source, "--environment", "dev", "--image", image, "--sbom", str(sbom), "--advisory-metadata", str(metadata), "--advisory-report", str(report), "--provenance", str(provenance), "--out", str(admission)]
        result = subprocess.run(command, cwd=ROOT, text=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT, check=False)
        if result.returncode != 0 or "HOST-020 hosted deployment admission: PASS" not in result.stdout:
            raise SystemExit("positive deployment admission failed:\n" + result.stdout)
        contract = load_contract()
        loaded, evidence_digest = contract.load_deployment_admission(admission, "dev", image)
        if loaded.get("sourceSha") != source or len(evidence_digest) != 64:
            raise SystemExit("deployment admission loader lost exact evidence identity")
        adverse = json.loads(admission.read_text(encoding="utf-8"))
        adverse["targetEnvironment"] = "stage"
        admission.write_text(json.dumps(adverse), encoding="utf-8")
        try:
            contract.load_deployment_admission(admission, "dev", image)
        except SystemExit:
            pass
        else:
            raise SystemExit("deployment admission accepted cross-environment reuse")
        report.write_text("scanner unavailable\n", encoding="utf-8")
        failed = subprocess.run(command, cwd=ROOT, text=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT, check=False)
        if failed.returncode == 0:
            raise SystemExit("deployment admission accepted missing advisory PASS")
    print("HOST-020 hosted deployment admission self-test: PASS")


if __name__ == "__main__":
    main()
