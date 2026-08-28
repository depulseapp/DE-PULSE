#!/usr/bin/env python3
"""Render DE.PULSE HOST-013/014 Kubernetes + Istio trust resources.

The canonical policy authority is internal/hostedenv/desired_state_v1.json.
This renderer is intentionally fail-closed and contains no cloud credentials.
"""
from __future__ import annotations

import argparse
import hashlib
import ipaddress
import json
import re
from pathlib import Path
from urllib.parse import urlparse

ROOT = Path(__file__).resolve().parents[2]
DESIRED_STATE = ROOT / "internal" / "hostedenv" / "desired_state_v1.json"
CANONICAL_ENVIRONMENTS = ("dev", "test", "stage", "prod")
HOST_RE = re.compile(r"^(?=.{1,253}$)(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$", re.I)


def load_manifest() -> dict:
    manifest = json.loads(DESIRED_STATE.read_text(encoding="utf-8"))
    if manifest.get("schema") != "DE.PULSE-HOSTED-DESIRED-STATE-1":
        raise SystemExit("unsupported hosted desired-state schema")
    environments = manifest.get("environments")
    if not isinstance(environments, dict) or tuple(environments.keys()) != CANONICAL_ENVIRONMENTS:
        raise SystemExit("desired state must define ordered dev/test/stage/prod environments")
    return manifest


def canonical_digest(manifest: dict, environment: str) -> str:
    payload = {
        "schema": manifest["schema"],
        "version": manifest["version"],
        "environment": environment,
        "state": manifest["environments"][environment],
    }
    # Match Go encoding/json for this structure: compact JSON, insertion order, no ASCII escaping changes needed here.
    encoded = json.dumps(payload, separators=(",", ":"), ensure_ascii=False).encode("utf-8")
    return hashlib.sha256(encoded).hexdigest()


def parse_hosts(path: Path) -> list[str]:
    hosts: list[str] = []
    seen: set[str] = set()
    for raw in path.read_text(encoding="utf-8").splitlines():
        host = raw.strip().lower()
        if not host or host.startswith("#"):
            continue
        if any(token in host for token in ("*", "/", ":")):
            raise SystemExit(f"invalid egress host {host!r}: wildcards, URLs and ports are forbidden")
        try:
            ipaddress.ip_address(host)
        except ValueError:
            pass
        else:
            raise SystemExit(f"invalid egress host {host!r}: IP literals are forbidden")
        if not HOST_RE.fullmatch(host):
            raise SystemExit(f"invalid egress host {host!r}: expected concrete DNS hostname")
        if host not in seen:
            seen.add(host)
            hosts.append(host)
    if not hosts:
        raise SystemExit("canonical external-host inventory is empty; refusing broad egress")
    return sorted(hosts)


def q(value: str) -> str:
    return json.dumps(value)


def render(environment: str, hosts: list[str]) -> str:
    manifest = load_manifest()
    if environment not in CANONICAL_ENVIRONMENTS:
        raise SystemExit(f"unsupported environment {environment!r}")
    state = manifest["environments"][environment]
    required = {
        "ingressPolicy": "managed-edge:https:443->mesh-mtls:8080",
        "egressPolicy": "default-deny:canonical-external-host-allowlist:https-wss:443",
        "networkPolicy": "isolated-namespace:default-deny",
        "tlsPolicy": "edge-min-tls1.2:internal-mtls-required",
    }
    for field, expected in required.items():
        if state.get(field) != expected:
            raise SystemExit(f"canonical desired state drift for {environment}: {field}")
    if state.get("internalMTLS") is not True:
        raise SystemExit(f"canonical desired state drift for {environment}: internalMTLS")

    namespace = state["isolationId"]
    service_account = state["serviceIdentity"]
    digest = canonical_digest(manifest, environment)
    annotations = (
        f"    depulse.io/desired-state-version: {q(manifest['version'])}\n"
        f"    depulse.io/desired-state-sha256: {q(digest)}"
    )
    host_yaml = "\n".join(f"    - {q(host)}" for host in hosts)
    host_csv = ",".join(hosts)
    ingress_principal = "cluster.local/ns/istio-system/sa/istio-ingressgateway-service-account"

    docs = [
        f"""apiVersion: v1
kind: Namespace
metadata:
  name: {namespace}
  labels:
    depulse.io/environment: {environment}
    istio-injection: enabled
  annotations:
{annotations}
""",
        f"""apiVersion: v1
kind: ServiceAccount
metadata:
  name: {service_account}
  namespace: {namespace}
  annotations:
{annotations}
automountServiceAccountToken: false
""",
        f"""apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny
  namespace: {namespace}
  annotations:
{annotations}
spec:
  podSelector: {{}}
  policyTypes: [Ingress, Egress]
""",
        f"""apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-managed-ingress
  namespace: {namespace}
  annotations:
{annotations}
spec:
  podSelector:
    matchLabels:
      app.kubernetes.io/name: depulse
  policyTypes: [Ingress]
  ingress:
    - from:
        - namespaceSelector:
            matchLabels:
              kubernetes.io/metadata.name: istio-system
          podSelector:
            matchLabels:
              app: istio-ingressgateway
      ports:
        - protocol: TCP
          port: 8080
""",
        f"""apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-governed-egress
  namespace: {namespace}
  annotations:
{annotations}
spec:
  podSelector:
    matchLabels:
      app.kubernetes.io/name: depulse
  policyTypes: [Egress]
  egress:
    - to:
        - namespaceSelector:
            matchLabels:
              kubernetes.io/metadata.name: kube-system
      ports:
        - protocol: UDP
          port: 53
        - protocol: TCP
          port: 53
    - ports:
        - protocol: TCP
          port: 443
""",
        f"""apiVersion: security.istio.io/v1
kind: PeerAuthentication
metadata:
  name: strict-mtls
  namespace: {namespace}
  annotations:
{annotations}
spec:
  mtls:
    mode: STRICT
""",
        f"""apiVersion: security.istio.io/v1
kind: AuthorizationPolicy
metadata:
  name: managed-ingress-only
  namespace: {namespace}
  annotations:
{annotations}
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: depulse
  action: ALLOW
  rules:
    - from:
        - source:
            principals: [{q(ingress_principal)}]
      to:
        - operation:
            ports: ["8080"]
""",
        f"""apiVersion: networking.istio.io/v1
kind: Sidecar
metadata:
  name: governed-outbound
  namespace: {namespace}
  annotations:
{annotations}
spec:
  workloadSelector:
    labels:
      app.kubernetes.io/name: depulse
  outboundTrafficPolicy:
    mode: REGISTRY_ONLY
""",
        f"""apiVersion: networking.istio.io/v1
kind: ServiceEntry
metadata:
  name: canonical-external-egress
  namespace: {namespace}
  annotations:
{annotations}
    depulse.io/egress-hosts-sha256: {q(hashlib.sha256(host_csv.encode()).hexdigest())}
spec:
  hosts:
{host_yaml}
  location: MESH_EXTERNAL
  resolution: DNS
  ports:
    - number: 443
      name: tls
      protocol: TLS
""",
    ]
    return "---\n".join(doc.rstrip() + "\n" for doc in docs)


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--environment", required=True, choices=CANONICAL_ENVIRONMENTS)
    parser.add_argument("--egress-hosts-file", required=True, type=Path)
    parser.add_argument("--output", type=Path)
    args = parser.parse_args()
    hosts = parse_hosts(args.egress_hosts_file)
    rendered = render(args.environment, hosts)
    if args.output:
        args.output.write_text(rendered, encoding="utf-8")
    else:
        print(rendered, end="")


if __name__ == "__main__":
    main()
