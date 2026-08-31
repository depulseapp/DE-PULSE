#!/usr/bin/env python3
"""Ephemeral real-AKS HOST-013/014 positive/adverse trust probe.

The probe is dev-only and intentionally destructive only to temporary resources.
It retains secret-free JSON evidence only. A one-day self-signed TLS key/certificate
is generated for the probe and deleted locally and from the cluster in cleanup.
"""
from __future__ import annotations

import argparse
import base64
import hashlib
import json
from pathlib import Path
import socket
import subprocess
import tempfile
import time

PROBE_IMAGE = "python:3.13-alpine@sha256:540c7d91f98ff6880174c40e99067bf5941eb54d818a7a5e094d188b196a934d"
PROBE_HOST = "depulse-probe.invalid"
ALLOWED_EGRESS = "https://api.twelvedata.com/"
DENIED_EGRESS = "https://example.com/"
CLEANUP_MARKER = "PROBE_CLEANUP_VERIFIED"


def fail(message: str) -> None:
    raise RuntimeError("HOST-013/014 Azure traffic probe: " + message)


def local_run(args: list[str], *, capture: bool = False, check: bool = True) -> subprocess.CompletedProcess[str]:
    result = subprocess.run(
        args,
        text=True,
        stdout=subprocess.PIPE if capture else None,
        stderr=subprocess.STDOUT if capture else None,
        check=False,
    )
    if check and result.returncode != 0:
        fail(f"command failed ({result.returncode}): {' '.join(args)}\n{(result.stdout or '').strip()}")
    return result


def aks_command(rg: str, cluster: str, command: str, *, file: Path | None = None, expect_success: bool = True) -> dict:
    args = [
        "az", "aks", "command", "invoke",
        "--resource-group", rg,
        "--name", cluster,
        "--command", command,
    ]
    if file is not None:
        args.extend(["--file", str(file)])
    args.extend(["-o", "json"])
    result = local_run(args, capture=True, check=True)
    try:
        payload = json.loads(result.stdout or "{}")
    except json.JSONDecodeError as exc:
        fail("AKS command returned non-JSON: " + str(exc))
    exit_code = payload.get("exitCode")
    succeeded = exit_code == 0 or str(exit_code).strip() == "0"
    if expect_success and not succeeded:
        fail("AKS command failed: " + str(payload.get("logs") or ""))
    if not expect_success and succeeded:
        fail("AKS adverse command unexpectedly succeeded")
    return payload


def generate_cert(directory: Path) -> tuple[Path, Path]:
    cert = directory / "probe.crt"
    key = directory / "probe.key"
    local_run([
        "openssl", "req", "-x509", "-newkey", "rsa:2048", "-sha256", "-nodes",
        "-days", "1", "-subj", f"/CN={PROBE_HOST}",
        "-addext", f"subjectAltName=DNS:{PROBE_HOST}",
        "-keyout", str(key), "-out", str(cert),
    ], capture=True)
    return cert, key


def manifest(environment: str, istio_revision: str, cert: bytes, key: bytes) -> str:
    namespace = f"depulse-{environment}"
    service_account = f"depulse-web-{environment}"
    cert_b64 = base64.b64encode(cert).decode("ascii")
    key_b64 = base64.b64encode(key).decode("ascii")
    return f"""apiVersion: v1
kind: Secret
metadata:
  name: depulse-probe-tls
  namespace: aks-istio-ingress
type: kubernetes.io/tls
data:
  tls.crt: {cert_b64}
  tls.key: {key_b64}
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: depulse-trust-probe
  namespace: {namespace}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: depulse-trust-probe
  template:
    metadata:
      labels:
        app: depulse-trust-probe
        app.kubernetes.io/name: depulse
        azure.workload.identity/use: "true"
    spec:
      serviceAccountName: {service_account}
      containers:
        - name: probe
          image: {PROBE_IMAGE}
          imagePullPolicy: IfNotPresent
          command: ["python", "-m", "http.server", "8080", "--bind", "0.0.0.0"]
          ports:
            - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: depulse-trust-probe
  namespace: {namespace}
spec:
  selector:
    app: depulse-trust-probe
  ports:
    - name: http
      port: 8080
      targetPort: 8080
---
apiVersion: networking.istio.io/v1
kind: Gateway
metadata:
  name: depulse-probe-gateway
  namespace: {namespace}
spec:
  selector:
    istio: aks-istio-ingressgateway-external
  servers:
    - port:
        number: 443
        name: https
        protocol: HTTPS
      tls:
        mode: SIMPLE
        credentialName: depulse-probe-tls
        minProtocolVersion: TLSV1_2
      hosts:
        - {PROBE_HOST}
---
apiVersion: networking.istio.io/v1
kind: VirtualService
metadata:
  name: depulse-probe-vs
  namespace: {namespace}
spec:
  hosts:
    - {PROBE_HOST}
  gateways:
    - depulse-probe-gateway
  http:
    - route:
        - destination:
            host: depulse-trust-probe
            port:
              number: 8080
---
apiVersion: v1
kind: Namespace
metadata:
  name: depulse-crossenv-probe
---
apiVersion: v1
kind: Pod
metadata:
  name: crossenv-client
  namespace: depulse-crossenv-probe
spec:
  restartPolicy: Never
  containers:
    - name: client
      image: {PROBE_IMAGE}
      command: ["python", "-c", "import time; time.sleep(1800)"]
"""


def wait_ready(rg: str, cluster: str, namespace: str) -> None:
    command = (
        f"kubectl rollout status deployment/depulse-trust-probe -n {namespace} --timeout=180s && "
        "kubectl wait --for=condition=Ready pod/crossenv-client -n depulse-crossenv-probe --timeout=180s && "
        f"kubectl get pod -n {namespace} -l app=depulse-trust-probe -o jsonpath='{{.items[0].spec.containers[*].name}}'"
    )
    payload = aks_command(rg, cluster, command)
    logs = str(payload.get("logs") or "")
    if "istio-proxy" not in logs or "probe" not in logs:
        fail("managed Istio sidecar was not injected into the probe workload")


def ingress_ip(rg: str, cluster: str) -> str:
    for _ in range(24):
        payload = aks_command(
            rg,
            cluster,
            "printf 'IP=' && kubectl get service aks-istio-ingressgateway-external -n aks-istio-ingress -o jsonpath='{.status.loadBalancer.ingress[0].ip}' && echo",
        )
        logs = str(payload.get("logs") or "")
        marker = "IP="
        if marker in logs:
            value = logs.split(marker, 1)[1].splitlines()[0].strip()
            if value:
                return value
        time.sleep(5)
    fail("managed external Istio ingress did not receive a public IP")
    return ""


def openssl_tls11_client(target: str, cert: Path) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        [
            "openssl", "s_client",
            "-connect", target,
            "-servername", PROBE_HOST,
            "-tls1_1",
            "-cipher", "DEFAULT:@SECLEVEL=0",
            "-CAfile", str(cert),
            "-verify_return_error",
            "-verify_hostname", PROBE_HOST,
            "-brief",
        ],
        input="",
        text=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.STDOUT,
        check=False,
        timeout=15,
    )


def prove_tls11_client_capability(cert: Path, key: Path) -> None:
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as reservation:
        reservation.bind(("127.0.0.1", 0))
        port = reservation.getsockname()[1]
    server = subprocess.Popen(
        [
            "openssl", "s_server",
            "-accept", str(port),
            "-cert", str(cert),
            "-key", str(key),
            "-tls1_1",
            "-cipher", "DEFAULT:@SECLEVEL=0",
            "-naccept", "1",
        ],
        stdin=subprocess.DEVNULL,
        stdout=subprocess.PIPE,
        stderr=subprocess.STDOUT,
        text=True,
    )
    try:
        time.sleep(0.5)
        if server.poll() is not None:
            output = server.communicate(timeout=1)[0] if server.stdout is not None else ""
            fail("local TLS 1.1 capability server could not start: " + output)
        client = openssl_tls11_client(f"127.0.0.1:{port}", cert)
        output = client.stdout or ""
        if client.returncode != 0 or "TLSv1.1" not in output:
            fail("runner cannot prove TLS 1.1 client capability before edge rejection test: " + output)
    finally:
        if server.poll() is None:
            server.terminate()
            try:
                server.wait(timeout=3)
            except subprocess.TimeoutExpired:
                server.kill()
                server.wait(timeout=3)


def test_edge_tls(ip: str, cert: Path, key: Path) -> dict[str, bool]:
    base = ["curl", "--silent", "--show-error", "--fail", "--cacert", str(cert), "--resolve", f"{PROBE_HOST}:443:{ip}"]
    positive = local_run(base + ["--tlsv1.2", "--tls-max", "1.2", f"https://{PROBE_HOST}/"], capture=True, check=False)
    if positive.returncode != 0:
        fail("TLS 1.2 managed ingress probe failed: " + (positive.stdout or ""))

    prove_tls11_client_capability(cert, key)
    legacy = openssl_tls11_client(f"{ip}:443", cert)
    if legacy.returncode == 0:
        fail("TLS 1.1 unexpectedly negotiated at managed ingress: " + (legacy.stdout or ""))
    return {
        "edgeTls12Succeeded": True,
        "tls11ClientCapabilityProved": True,
        "edgeTls11Rejected": True,
    }


def test_direct_ingress_denial(rg: str, cluster: str, environment: str) -> bool:
    namespace = f"depulse-{environment}"
    command = (
        "if kubectl exec -n depulse-crossenv-probe crossenv-client -- python -c \""
        "import urllib.request; urllib.request.urlopen('http://depulse-trust-probe."
        f"{namespace}.svc.cluster.local:8080/', timeout=5).read(1)\"; "
        "then echo UNAUTHORIZED_DIRECT_INGRESS_SUCCEEDED; exit 9; else echo DIRECT_INGRESS_DENIED; fi"
    )
    payload = aks_command(rg, cluster, command)
    if "DIRECT_INGRESS_DENIED" not in str(payload.get("logs") or ""):
        fail("cross-environment direct ingress denial marker missing")
    return True


def test_egress(rg: str, cluster: str, environment: str) -> dict[str, bool]:
    namespace = f"depulse-{environment}"
    allowed_py = (
        "import urllib.request,urllib.error; u='" + ALLOWED_EGRESS + "'; "
        "\ntry:\n urllib.request.urlopen(u,timeout=10).read(1)\nexcept urllib.error.HTTPError:\n pass"
    )
    allowed = (
        f"kubectl exec -n {namespace} deploy/depulse-trust-probe -c probe -- python -c \"{allowed_py}\""
    )
    aks_command(rg, cluster, allowed)

    denied_py = "import urllib.request; urllib.request.urlopen('" + DENIED_EGRESS + "',timeout=5).read(1)"
    denied = (
        f"if kubectl exec -n {namespace} deploy/depulse-trust-probe -c probe -- python -c \"{denied_py}\"; "
        "then echo UNREGISTERED_EGRESS_SUCCEEDED; exit 9; else echo UNREGISTERED_EGRESS_DENIED; fi"
    )
    payload = aks_command(rg, cluster, denied)
    if "UNREGISTERED_EGRESS_DENIED" not in str(payload.get("logs") or ""):
        fail("unregistered egress denial marker missing")
    return {"canonicalEgressSucceeded": True, "unregisteredEgressRejected": True}


def test_workload_identity_token(rg: str, cluster: str, environment: str) -> bool:
    namespace = f"depulse-{environment}"
    script = r'''import json, os, urllib.parse, urllib.request
client=os.environ["AZURE_CLIENT_ID"]
tenant=os.environ["AZURE_TENANT_ID"]
path=os.environ["AZURE_FEDERATED_TOKEN_FILE"]
assert client and tenant and path
assert os.path.exists(path)
assert os.path.getsize(path) > 100
assert os.environ.get("AZURE_AUTHORITY_HOST", "").startswith("https://login.microsoftonline.com")
assertion=open(path, encoding="utf-8").read().strip()
data=urllib.parse.urlencode({
 "client_id": client,
 "scope": "https://management.azure.com/.default",
 "grant_type": "client_credentials",
 "client_assertion_type": "urn:ietf:params:oauth:client-assertion-type:jwt-bearer",
 "client_assertion": assertion,
}).encode()
req=urllib.request.Request(f"https://login.microsoftonline.com/{tenant}/oauth2/v2.0/token", data=data, method="POST")
with urllib.request.urlopen(req, timeout=15) as resp:
 payload=json.load(resp)
assert payload.get("access_token") and payload.get("token_type") == "Bearer"
print("WORKLOAD_IDENTITY_TOKEN_OK")'''
    encoded = base64.b64encode(script.encode()).decode()
    command = (
        f"kubectl exec -n {namespace} deploy/depulse-trust-probe -c probe -- python -c "
        f"\"import base64; exec(base64.b64decode('{encoded}'))\""
    )
    payload = aks_command(rg, cluster, command)
    if "WORKLOAD_IDENTITY_TOKEN_OK" not in str(payload.get("logs") or ""):
        fail("workload identity token exchange marker missing")
    return True


def cleanup_command(environment: str) -> str:
    namespace = f"depulse-{environment}"
    return (
        f"kubectl delete gateway depulse-probe-gateway -n {namespace} --ignore-not-found --wait=true --timeout=120s && "
        f"kubectl delete virtualservice depulse-probe-vs -n {namespace} --ignore-not-found --wait=true --timeout=120s && "
        f"kubectl delete service depulse-trust-probe -n {namespace} --ignore-not-found --wait=true --timeout=120s && "
        f"kubectl delete deployment depulse-trust-probe -n {namespace} --ignore-not-found --wait=true --timeout=120s && "
        "kubectl delete secret depulse-probe-tls -n aks-istio-ingress --ignore-not-found --wait=true --timeout=120s && "
        "kubectl delete namespace depulse-crossenv-probe --ignore-not-found --wait=true --timeout=120s && "
        f"test -z \"$(kubectl get gateway depulse-probe-gateway -n {namespace} --ignore-not-found -o name)\" && "
        f"test -z \"$(kubectl get virtualservice depulse-probe-vs -n {namespace} --ignore-not-found -o name)\" && "
        f"test -z \"$(kubectl get service depulse-trust-probe -n {namespace} --ignore-not-found -o name)\" && "
        f"test -z \"$(kubectl get deployment depulse-trust-probe -n {namespace} --ignore-not-found -o name)\" && "
        "test -z \"$(kubectl get secret depulse-probe-tls -n aks-istio-ingress --ignore-not-found -o name)\" && "
        "test -z \"$(kubectl get namespace depulse-crossenv-probe --ignore-not-found -o name)\" && "
        f"echo {CLEANUP_MARKER}"
    )


def cleanup(rg: str, cluster: str, environment: str) -> bool:
    try:
        payload = aks_command(rg, cluster, cleanup_command(environment))
    except Exception as exc:
        print("WARNING: HOST-013/014 probe cleanup verification failed:", exc)
        return False
    if CLEANUP_MARKER not in str(payload.get("logs") or ""):
        print("WARNING: HOST-013/014 probe cleanup verification marker missing")
        return False
    return True


def self_test() -> None:
    if "@sha256:" not in PROBE_IMAGE:
        fail("probe image must remain digest pinned")
    with tempfile.TemporaryDirectory(prefix="depulse-host013-tls11-selftest-") as tmp:
        cert, key = generate_cert(Path(tmp))
        prove_tls11_client_capability(cert, key)
    rendered = manifest("dev", "asm-1-27", b"certificate", b"private-key")
    required = (
        "istio: aks-istio-ingressgateway-external",
        "minProtocolVersion: TLSV1_2",
        'azure.workload.identity/use: "true"',
        "serviceAccountName: depulse-web-dev",
        "namespace: aks-istio-ingress",
        PROBE_IMAGE,
    )
    for token in required:
        if token not in rendered:
            fail("probe manifest self-test missing " + token)
    cleanup_text = cleanup_command("dev")
    cleanup_required = (
        "--wait=true --timeout=120s",
        "kubectl get gateway depulse-probe-gateway",
        "kubectl get namespace depulse-crossenv-probe",
        CLEANUP_MARKER,
    )
    for token in cleanup_required:
        if token not in cleanup_text:
            fail("probe cleanup self-test missing " + token)
    print("HOST-013/014 Azure traffic probe self-test: PASS")


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--self-test", action="store_true")
    parser.add_argument("--resource-group", default="")
    parser.add_argument("--cluster-name", default="")
    parser.add_argument("--environment", default="dev", choices=["dev"])
    parser.add_argument("--candidate-sha", default="")
    parser.add_argument("--istio-revision", default="")
    parser.add_argument("--output", default=".depulse-host013-azure/host013-azure-traffic-evidence.json")
    args = parser.parse_args()
    if args.self_test:
        self_test()
        return 0
    for name in ("resource_group", "cluster_name", "candidate_sha", "istio_revision"):
        if not getattr(args, name):
            parser.error(f"--{name.replace('_', '-')} is required")

    namespace = f"depulse-{args.environment}"
    checks: dict[str, bool] = {}
    with tempfile.TemporaryDirectory(prefix="depulse-host013-probe-") as tmp:
        directory = Path(tmp)
        cert, key = generate_cert(directory)
        probe_manifest = directory / "host013-traffic-probe.yaml"
        probe_manifest.write_text(
            manifest(args.environment, args.istio_revision, cert.read_bytes(), key.read_bytes()),
            encoding="utf-8",
        )
        try:
            aks_command(
                args.resource_group,
                args.cluster_name,
                "kubectl apply -f host013-traffic-probe.yaml",
                file=probe_manifest,
            )
            wait_ready(args.resource_group, args.cluster_name, namespace)
            checks["managedIstioSidecarInjected"] = True
            ip = ingress_ip(args.resource_group, args.cluster_name)
            checks.update(test_edge_tls(ip, cert, key))
            checks["managedIngressToStrictMtlsWorkloadSucceeded"] = True
            checks["crossEnvironmentDirectIngressRejected"] = test_direct_ingress_denial(
                args.resource_group, args.cluster_name, args.environment
            )
            checks.update(test_egress(args.resource_group, args.cluster_name, args.environment))
            checks["workloadIdentityTokenExchangeSucceeded"] = test_workload_identity_token(
                args.resource_group, args.cluster_name, args.environment
            )
        finally:
            checks["probeCleanupVerified"] = cleanup(args.resource_group, args.cluster_name, args.environment)

    evidence = {
        "schema": "DE.PULSE-HOST013-AZURE-TRAFFIC-EVIDENCE-1",
        "requirements": ["HOST-013", "HOST-014"],
        "candidateSha": args.candidate_sha,
        "environment": args.environment,
        "managedIstioRevision": args.istio_revision,
        "probeImage": PROBE_IMAGE,
        "probeHost": PROBE_HOST,
        "checks": checks,
        "ephemeralTlsCredentialRetained": False,
        "containsSecrets": False,
        "cleanupAttempted": True,
        "cleanupVerified": checks.get("probeCleanupVerified") is True,
        "status": "PASS",
    }
    if not checks or not all(checks.values()) or evidence["cleanupVerified"] is not True:
        fail("not all live trust probes and cleanup checks passed")
    output = Path(args.output)
    output.parent.mkdir(parents=True, exist_ok=True)
    output.write_text(json.dumps(evidence, indent=2, sort_keys=True) + "\n", encoding="utf-8")
    print("HOST-013/014 Azure traffic evidence: PASS ->", output)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
