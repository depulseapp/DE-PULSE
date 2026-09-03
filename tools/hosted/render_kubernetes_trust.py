#!/usr/bin/env python3
"""Render DE.PULSE hosted Kubernetes + Istio trust resources.

The canonical environment policy authority is internal/hostedenv/desired_state_v1.json.
The canonical external egress authority is governance/hosted-infrastructure/external-egress-v1.json.
HOST-017/018 workload rendering composes the version-pinned managed-secret reference
contract without accepting or serializing secret values.
"""
from __future__ import annotations

import argparse
import hashlib
import ipaddress
import json
import re
from pathlib import Path
from urllib.parse import urlsplit

from render_managed_secrets import (
    load_reference_manifest,
    reference_generation,
    required_aliases,
    secret_provider_class_name,
)

ROOT = Path(__file__).resolve().parents[2]
DESIRED_STATE = ROOT / "internal" / "hostedenv" / "desired_state_v1.json"
EGRESS_INVENTORY = ROOT / "governance" / "hosted-infrastructure" / "external-egress-v1.json"
CANONICAL_ENVIRONMENTS = ("dev", "test", "stage", "prod")
HOST_RE = re.compile(r"^(?=.{1,253}$)(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$", re.I)
ISTIO_REVISION_RE = re.compile(r"^asm-[0-9]+-[0-9]+$")
AZURE_CLIENT_ID_RE = re.compile(r"^[0-9a-fA-F-]{36}$")
IMMUTABLE_IMAGE_RE = re.compile(r"^[^\s]+@sha256:[0-9a-fA-F]{64}$")


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
    encoded = json.dumps(payload, separators=(",", ":"), ensure_ascii=False).encode("utf-8")
    return hashlib.sha256(encoded).hexdigest()


def validate_host(host: str) -> str:
    host = host.strip().lower().rstrip(".")
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
    return host


def parse_hosts(path: Path) -> list[str]:
    hosts: list[str] = []
    seen: set[str] = set()
    for raw in path.read_text(encoding="utf-8").splitlines():
        host = raw.strip().lower()
        if not host or host.startswith("#"):
            continue
        host = validate_host(host)
        if host not in seen:
            seen.add(host)
            hosts.append(host)
    if not hosts:
        raise SystemExit("canonical external-host inventory is empty; refusing broad egress")
    return sorted(hosts)


def canonical_egress_hosts(path: Path = EGRESS_INVENTORY) -> list[str]:
    payload = json.loads(path.read_text(encoding="utf-8"))
    if payload.get("schema") != "DE.PULSE-HOSTED-EXTERNAL-EGRESS-1" or payload.get("version") != "v1":
        raise SystemExit("unsupported hosted external-egress inventory")
    if str(payload.get("defaultPolicy", "")).upper() != "DENY":
        raise SystemExit("external egress inventory must default DENY")
    allowed_protocols = payload.get("allowedProtocols")
    if allowed_protocols != ["https", "wss"]:
        raise SystemExit("external egress allowedProtocols must remain exactly https,wss")
    hosts: set[str] = set()
    for row in payload.get("rules", []):
        if not isinstance(row, dict) or row.get("runtimeAllowed") is not True:
            continue
        protocols = [str(value).lower() for value in row.get("protocols", [])]
        if not protocols or any(protocol not in allowed_protocols for protocol in protocols):
            raise SystemExit("runtime-authorized egress entry uses a forbidden protocol")
        hosts.add(validate_host(str(row.get("host", ""))))
    if not hosts:
        raise SystemExit("canonical external-host inventory is empty; refusing broad egress")
    return sorted(hosts)


def q(value: str) -> str:
    return json.dumps(value)


def validate_public_origin(value: str) -> str:
    value = value.strip()
    parts = urlsplit(value)
    if (
        parts.scheme != "https"
        or not parts.netloc
        or parts.username is not None
        or parts.password is not None
        or parts.path not in ("", "/")
        or parts.query
        or parts.fragment
        or "*" in parts.netloc
    ):
        raise SystemExit("hosted workload public origin must be one explicit HTTPS origin")
    return value.rstrip("/")


def render(
    environment: str,
    hosts: list[str],
    *,
    mesh_profile: str = "portable",
    istio_revision: str | None = None,
    workload_identity_client_id: str | None = None,
    workload_image: str | None = None,
    public_origin: str | None = None,
    managed_secret_reference_manifest: Path | None = None,
) -> str:
    manifest = load_manifest()
    if environment not in CANONICAL_ENVIRONMENTS:
        raise SystemExit(f"unsupported environment {environment!r}")
    if mesh_profile not in {"portable", "aks-managed"}:
        raise SystemExit(f"unsupported mesh profile {mesh_profile!r}")
    if mesh_profile == "aks-managed":
        if not istio_revision or not ISTIO_REVISION_RE.fullmatch(istio_revision):
            raise SystemExit("AKS managed Istio requires an explicit asm-X-Y revision")
        if not workload_identity_client_id or not AZURE_CLIENT_ID_RE.fullmatch(workload_identity_client_id):
            raise SystemExit("AKS managed workload identity requires a concrete Azure client ID")
    elif istio_revision or workload_identity_client_id:
        raise SystemExit("Istio revision/workload identity client ID are valid only for aks-managed rendering")

    workload_requested = any(value is not None for value in (workload_image, public_origin, managed_secret_reference_manifest))
    if workload_requested:
        if mesh_profile != "aks-managed":
            raise SystemExit("managed hosted workload rendering requires the aks-managed profile")
        if not workload_image or not IMMUTABLE_IMAGE_RE.fullmatch(workload_image):
            raise SystemExit("hosted workload image must be pinned by sha256 digest")
        if not public_origin:
            raise SystemExit("hosted workload rendering requires --public-origin")
        public_origin = validate_public_origin(public_origin)
        if managed_secret_reference_manifest is None:
            raise SystemExit("hosted workload rendering requires --managed-secret-reference-manifest")

    hosts = sorted({validate_host(host) for host in hosts})
    if not hosts:
        raise SystemExit("canonical external-host inventory is empty; refusing broad egress")
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

    if mesh_profile == "aks-managed":
        namespace_mesh_label = f"    istio.io/rev: {q(istio_revision or '')}"
        ingress_namespace = "aks-istio-ingress"
        control_plane_namespace = "aks-istio-system"
        ingress_pod_selector = "              istio: aks-istio-ingressgateway-external"
        auth_from = "            namespaces: [\"aks-istio-ingress\"]"
        service_account_extra = f"\n    azure.workload.identity/client-id: {q(workload_identity_client_id or '')}"
    else:
        namespace_mesh_label = "    istio-injection: enabled"
        ingress_namespace = "istio-system"
        control_plane_namespace = "istio-system"
        ingress_pod_selector = "              app: istio-ingressgateway"
        auth_from = "            principals: [\"cluster.local/ns/istio-system/sa/istio-ingressgateway-service-account\"]"
        service_account_extra = ""

    docs = [
        f"""apiVersion: v1
kind: Namespace
metadata:
  name: {namespace}
  labels:
    depulse.io/environment: {environment}
{namespace_mesh_label}
  annotations:
{annotations}
""",
        f"""apiVersion: v1
kind: ServiceAccount
metadata:
  name: {service_account}
  namespace: {namespace}
  annotations:
{annotations}{service_account_extra}
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
              kubernetes.io/metadata.name: {ingress_namespace}
          podSelector:
            matchLabels:
{ingress_pod_selector}
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
    - to:
        - namespaceSelector:
            matchLabels:
              kubernetes.io/metadata.name: {control_plane_namespace}
      ports:
        - protocol: TCP
          port: 15012
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
{auth_from}
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

    if workload_requested:
        references = load_reference_manifest(managed_secret_reference_manifest, environment)
        secret_generation = reference_generation(environment, references)
        secret_aliases = required_aliases(references)
        provider_class = secret_provider_class_name(secret_generation)
        workload_annotations = (
            f"        depulse.io/desired-state-version: {q(manifest['version'])}\n"
            f"        depulse.io/desired-state-sha256: {q(digest)}\n"
            f"        depulse.io/secret-generation-sha256: {q(secret_generation)}"
        )
        docs.extend([
            f"""apiVersion: apps/v1
kind: Deployment
metadata:
  name: depulse-web
  namespace: {namespace}
  annotations:
{annotations}
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: depulse
      app.kubernetes.io/component: hosted-web
  template:
    metadata:
      labels:
        app.kubernetes.io/name: depulse
        app.kubernetes.io/component: hosted-web
        azure.workload.identity/use: "true"
      annotations:
{workload_annotations}
    spec:
      serviceAccountName: {service_account}
      containers:
        - name: depulse
          image: {q(workload_image or '')}
          imagePullPolicy: IfNotPresent
          env:
            - name: DEPULSE_RUNTIME_MODE
              value: "hosted"
            - name: DEPULSE_LISTEN_ADDR
              value: ":8080"
            - name: DEPULSE_TRUST_PROXY_HEADERS
              value: "true"
            - name: DEPULSE_PUBLIC_ORIGIN
              value: {q(public_origin or '')}
            - name: DEPULSE_HOSTED_ENVIRONMENT
              value: {q(environment)}
            - name: DEPULSE_HOSTED_DESIRED_STATE_VERSION
              value: {q(str(manifest['version']))}
            - name: DEPULSE_HOSTED_DESIRED_STATE_SHA256
              value: {q(digest)}
            - name: DEPULSE_HOSTED_ISOLATION_ID
              value: {q(str(state['isolationId']))}
            - name: DEPULSE_HOSTED_SERVICE_IDENTITY
              value: {q(str(state['serviceIdentity']))}
            - name: DEPULSE_HOSTED_INGRESS_POLICY
              value: {q(str(state['ingressPolicy']))}
            - name: DEPULSE_HOSTED_EGRESS_POLICY
              value: {q(str(state['egressPolicy']))}
            - name: DEPULSE_HOSTED_NETWORK_POLICY
              value: {q(str(state['networkPolicy']))}
            - name: DEPULSE_HOSTED_TLS_POLICY
              value: {q(str(state['tlsPolicy']))}
            - name: DEPULSE_HOSTED_INTERNAL_MTLS
              value: "true"
            - name: DEPULSE_HOSTED_SECRETS_DIR
              value: "/var/run/depulse/secrets"
            - name: DEPULSE_HOSTED_REQUIRED_SECRETS
              value: {q(secret_aliases)}
            - name: DEPULSE_HOSTED_SECRET_GENERATION
              value: {q(secret_generation)}
            - name: XDG_CONFIG_HOME
              value: "/var/lib/depulse"
          ports:
            - name: http
              containerPort: 8080
          readinessProbe:
            httpGet:
              path: /api/ready
              port: http
            periodSeconds: 5
            timeoutSeconds: 2
            failureThreshold: 3
          livenessProbe:
            httpGet:
              path: /api/health
              port: http
            periodSeconds: 15
            timeoutSeconds: 2
            failureThreshold: 3
          volumeMounts:
            - name: managed-secrets
              mountPath: /var/run/depulse/secrets
              readOnly: true
            - name: runtime-config
              mountPath: /var/lib/depulse
      volumes:
        - name: managed-secrets
          csi:
            driver: secrets-store.csi.k8s.io
            readOnly: true
            volumeAttributes:
              secretProviderClass: {provider_class}
        - name: runtime-config
          emptyDir: {{}}
""",
            f"""apiVersion: v1
kind: Service
metadata:
  name: depulse-web
  namespace: {namespace}
  annotations:
{annotations}
spec:
  selector:
    app.kubernetes.io/name: depulse
    app.kubernetes.io/component: hosted-web
  ports:
    - name: http
      port: 8080
      targetPort: http
""",
        ])
    return "---\n".join(doc.rstrip() + "\n" for doc in docs)


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--environment", required=True, choices=CANONICAL_ENVIRONMENTS)
    source = parser.add_mutually_exclusive_group()
    source.add_argument("--egress-hosts-file", type=Path, help="test/compatibility seam for explicit DNS-host fixtures")
    source.add_argument("--egress-inventory", type=Path, default=EGRESS_INVENTORY, help="canonical JSON egress inventory")
    parser.add_argument("--mesh-profile", choices=["portable", "aks-managed"], default="portable")
    parser.add_argument("--istio-revision")
    parser.add_argument("--workload-identity-client-id")
    parser.add_argument("--workload-image")
    parser.add_argument("--public-origin")
    parser.add_argument("--managed-secret-reference-manifest", type=Path)
    parser.add_argument("--output", type=Path)
    args = parser.parse_args()
    hosts = parse_hosts(args.egress_hosts_file) if args.egress_hosts_file else canonical_egress_hosts(args.egress_inventory)
    rendered = render(
        args.environment,
        hosts,
        mesh_profile=args.mesh_profile,
        istio_revision=args.istio_revision,
        workload_identity_client_id=args.workload_identity_client_id,
        workload_image=args.workload_image,
        public_origin=args.public_origin,
        managed_secret_reference_manifest=args.managed_secret_reference_manifest,
    )
    if args.output:
        args.output.write_text(rendered, encoding="utf-8")
    else:
        print(rendered, end="")


if __name__ == "__main__":
    main()
