#!/usr/bin/env python3
"""Fail-closed repository gate for HOST-013..014 infrastructure projection."""
from __future__ import annotations

import importlib.util
import json
import subprocess
import sys
import tempfile
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
RENDERER = ROOT / "tools" / "hosted" / "render_kubernetes_trust.py"
DESIRED = ROOT / "internal" / "hostedenv" / "desired_state_v1.json"


def fail(message: str) -> None:
    raise SystemExit(f"HOST-013/014 infrastructure gate: {message}")


def load_renderer():
    spec = importlib.util.spec_from_file_location("depulse_hosted_renderer", RENDERER)
    if spec is None or spec.loader is None:
        fail("cannot load renderer")
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def main() -> None:
    if not RENDERER.is_file():
        fail("canonical renderer missing")
    manifest = json.loads(DESIRED.read_text(encoding="utf-8"))
    environments = manifest.get("environments", {})
    if set(environments) != {"dev", "test", "stage", "prod"}:
        fail("canonical desired state must contain exactly dev/test/stage/prod")
    isolation_ids = [environments[e]["isolationId"] for e in ("dev", "test", "stage", "prod")]
    service_ids = [environments[e]["serviceIdentity"] for e in ("dev", "test", "stage", "prod")]
    if len(set(isolation_ids)) != 4 or len(set(service_ids)) != 4:
        fail("environment isolation IDs and service identities must be unique")

    module = load_renderer()
    test_hosts = ["example.invalid", "second.example.invalid"]
    required_tokens = (
        "kind: Namespace",
        "kind: ServiceAccount",
        "kind: NetworkPolicy",
        "name: default-deny",
        "name: allow-managed-ingress",
        "name: allow-governed-egress",
        "kind: PeerAuthentication",
        "mode: STRICT",
        "kind: AuthorizationPolicy",
        "kind: Sidecar",
        "mode: REGISTRY_ONLY",
        "kind: ServiceEntry",
        "location: MESH_EXTERNAL",
        "port: 8080",
        "port: 443",
    )
    for environment in ("dev", "test", "stage", "prod"):
        rendered = module.render(environment, test_hosts)
        state = environments[environment]
        for token in required_tokens:
            if token not in rendered:
                fail(f"{environment} render missing {token!r}")
        if f"name: {state['isolationId']}" not in rendered:
            fail(f"{environment} namespace not bound to canonical isolationId")
        if f"name: {state['serviceIdentity']}" not in rendered:
            fail(f"{environment} service account not bound to canonical serviceIdentity")
        digest = module.canonical_digest(manifest, environment)
        if f'depulse.io/desired-state-sha256: "{digest}"' not in rendered:
            fail(f"{environment} render missing canonical desired-state digest")
        for host in test_hosts:
            if f'    - "{host}"' not in rendered:
                fail(f"{environment} render missing explicit egress host")
        if "0.0.0.0/0" in rendered or "::/0" in rendered or "hosts:\n    - \"*\"" in rendered:
            fail(f"{environment} render contains broad egress authority")

    with tempfile.TemporaryDirectory() as tmp:
        empty = Path(tmp) / "empty-hosts.txt"
        empty.write_text("# intentionally empty\n", encoding="utf-8")
        result = subprocess.run(
            [sys.executable, str(RENDERER), "--environment", "dev", "--egress-hosts-file", str(empty)],
            cwd=ROOT,
            text=True,
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,
            check=False,
        )
        if result.returncode == 0 or "refusing broad egress" not in result.stdout:
            fail("renderer does not fail closed for empty egress inventory")

    print("HOST-013/014 infrastructure contract: PASS")


if __name__ == "__main__":
    main()
