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
AZURE_GATE = ROOT / "tools" / "ci" / "azure_hosted_infrastructure_gate.py"
EGRESS_GATE = ROOT / "tools" / "ci" / "hosted_external_egress_gate.py"


def fail(message: str) -> None:
    raise SystemExit(f"HOST-013/014 infrastructure gate: {message}")


def load_renderer():
    # The renderer intentionally composes sibling hosted renderers. Dynamic
    # import by absolute file path does not automatically expose that sibling
    # directory on sys.path, so make the renderer directory explicit only for
    # module initialization and then restore the caller's import path.
    renderer_dir = str(RENDERER.parent)
    added_path = renderer_dir not in sys.path
    if added_path:
        sys.path.insert(0, renderer_dir)
    try:
        spec = importlib.util.spec_from_file_location("depulse_hosted_renderer", RENDERER)
        if spec is None or spec.loader is None:
            fail("cannot load renderer")
        module = importlib.util.module_from_spec(spec)
        spec.loader.exec_module(module)
        return module
    finally:
        if added_path:
            sys.path.remove(renderer_dir)


def run_gate(path: Path, pass_marker: str, label: str) -> None:
    if not path.is_file():
        fail(label + " gate missing")
    result = subprocess.run(
        [sys.executable, str(path)],
        cwd=ROOT,
        text=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.STDOUT,
        check=False,
    )
    if result.returncode != 0:
        fail(label + " validation failed:\n" + result.stdout)
    if pass_marker not in result.stdout:
        fail(label + " gate did not emit canonical PASS marker")


def verify_common(rendered: str, environment: str, state: dict, module, manifest: dict, test_hosts: list[str]) -> None:
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
    for environment in ("dev", "test", "stage", "prod"):
        state = environments[environment]
        portable = module.render(environment, test_hosts)
        verify_common(portable, environment, state, module, manifest, test_hosts)
        if "istio-injection: enabled" not in portable or "kubernetes.io/metadata.name: istio-system" not in portable:
            fail(f"{environment} portable mesh projection drift")

        managed = module.render(
            environment,
            test_hosts,
            mesh_profile="aks-managed",
            istio_revision="asm-1-27",
            workload_identity_client_id="11111111-1111-1111-1111-111111111111",
        )
        verify_common(managed, environment, state, module, manifest, test_hosts)
        required_managed = (
            'istio.io/rev: "asm-1-27"',
            "kubernetes.io/metadata.name: aks-istio-ingress",
            "istio: aks-istio-ingressgateway-external",
            'namespaces: ["aks-istio-ingress"]',
            'azure.workload.identity/client-id: "11111111-1111-1111-1111-111111111111"',
        )
        for token in required_managed:
            if token not in managed:
                fail(f"{environment} AKS managed render missing {token!r}")
        if "istio-injection: enabled" in managed or "kubernetes.io/metadata.name: istio-system" in managed:
            fail(f"{environment} AKS managed render leaked self-managed Istio conventions")

    try:
        module.render(
            "dev",
            test_hosts,
            mesh_profile="aks-managed",
            istio_revision="invalid",
            workload_identity_client_id="11111111-1111-1111-1111-111111111111",
        )
    except SystemExit:
        pass
    else:
        fail("AKS managed renderer accepted invalid Istio revision")

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

    run_gate(EGRESS_GATE, "HOST-014 external egress conservation: PASS", "HOST-014 external egress")
    run_gate(AZURE_GATE, "PASS: Azure AKS HOST-013..014 adapter", "Azure AKS hosted infrastructure")
    print("HOST-013/014 infrastructure contract: PASS")


if __name__ == "__main__":
    main()
