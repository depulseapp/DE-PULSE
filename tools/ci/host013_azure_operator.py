#!/usr/bin/env python3
"""Governed HOST-013..014 Azure AKS operator.

This is the single repository owner for the manual Azure infrastructure run.
It assumes GitHub/Azure OIDC authentication is already established and Terraform is
installed. It never accepts a client secret, storage key, SAS token or kubeconfig.
"""
from __future__ import annotations

import argparse
import hashlib
import json
import os
from pathlib import Path
import re
import subprocess
import sys
import time

from azure_oidc_cli import refresh_azure_cli_oidc

ROOT = Path(__file__).resolve().parents[2]
AZURE_DIR = ROOT / "governance" / "hosted-infrastructure" / "azure"
RENDERER = ROOT / "tools" / "hosted" / "render_kubernetes_trust.py"
LIVE_EVIDENCE = ROOT / "tools" / "ci" / "host013_azure_live_evidence.py"
TRAFFIC_PROBE = ROOT / "tools" / "ci" / "host013_azure_traffic_probe.py"
CONFIRMATION = "HOST013_AZURE_AKS_OPERATOR_DRILL"
ISTIO_RE = re.compile(r"^asm-(\d+)-(\d+)$")
AKS_VERIFY_ROLE = "Azure Kubernetes Service RBAC Cluster Admin"
AKS_RBAC_PROPAGATION_TIMEOUT_SECONDS = 300
AKS_RBAC_POLL_SECONDS = 10
ISTIO_READY_TIMEOUT_SECONDS = 240
ISTIO_READY_POLL_SECONDS = 10
ISTIO_READY_STABLE_PASSES = 3
INGRESS_READY_TIMEOUT_SECONDS = 600
INGRESS_READY_POLL_SECONDS = 10
INGRESS_READY_STABLE_PASSES = 3
INGRESS_RECONCILE_AFTER_SECONDS = 30


def fail(message: str) -> None:
    raise SystemExit("HOST-013/014 Azure operator: " + message)


def run(
    args: list[str],
    *,
    env: dict[str, str] | None = None,
    capture: bool = False,
    check: bool = True,
) -> subprocess.CompletedProcess[str]:
    if args and args[0] == "az":
        try:
            refresh_azure_cli_oidc()
        except RuntimeError as exc:
            fail("Azure CLI OIDC refresh failed: " + str(exc))
    result = subprocess.run(
        args,
        cwd=ROOT,
        env=env,
        text=True,
        stdout=subprocess.PIPE if capture else None,
        stderr=subprocess.STDOUT if capture else None,
        check=False,
    )
    if check and result.returncode != 0:
        detail = (result.stdout or "").strip() if capture else ""
        fail(f"command failed ({result.returncode}): {' '.join(args)}" + ("\n" + detail if detail else ""))
    return result


def run_json(args: list[str]) -> object:
    result = run(args, capture=True)
    try:
        return json.loads(result.stdout or "null")
    except json.JSONDecodeError as exc:
        fail(f"non-JSON command response from {' '.join(args)}: {exc}")


def aks_remote_succeeded(payload: object) -> bool:
    if not isinstance(payload, dict):
        return False
    exit_code = payload.get("exitCode")
    return exit_code == 0 or str(exit_code).strip() == "0"


def require_aks_command_success(payload: object, label: str) -> dict[str, object]:
    if not isinstance(payload, dict):
        fail(f"{label} returned a non-object AKS command result")
    if not aks_remote_succeeded(payload):
        exit_code = payload.get("exitCode")
        logs = str(payload.get("logs") or "").strip()
        diagnostic = logs[:8000] if logs else "<no remote logs>"
        fail(f"{label} remote command failed: exitCode={exit_code!r}; logs={diagnostic}")
    return payload


def revision_names(value: object) -> list[str]:
    names: set[str] = set()

    def walk(node: object) -> None:
        if isinstance(node, str) and ISTIO_RE.fullmatch(node):
            names.add(node)
        elif isinstance(node, list):
            for item in node:
                walk(item)
        elif isinstance(node, dict):
            for item in node.values():
                walk(item)

    walk(value)
    return sorted(names, key=lambda item: tuple(int(v) for v in ISTIO_RE.fullmatch(item).groups()))


def choose_revision(value: object) -> str:
    names = revision_names(value)
    if not names:
        fail("Azure returned no supported managed Istio asm-X-Y revision")
    return names[-1]


def terraform_env(args: argparse.Namespace, istio_revision: str) -> dict[str, str]:
    env = os.environ.copy()
    env.update({
        "ARM_USE_OIDC": "true",
        "ARM_USE_AZUREAD": "true",
        "ARM_CLIENT_ID": args.client_id,
        "ARM_TENANT_ID": args.tenant_id,
        "ARM_SUBSCRIPTION_ID": args.subscription_id,
        "TF_VAR_subscription_id": args.subscription_id,
        "TF_VAR_tenant_id": args.tenant_id,
        "TF_VAR_environment": args.environment,
        "TF_VAR_location": args.location,
        "TF_VAR_istio_revisions": json.dumps([istio_revision]),
        "TF_IN_AUTOMATION": "true",
    })
    return env


def terraform(args: argparse.Namespace, env: dict[str, str], *extra: str, capture: bool = False) -> subprocess.CompletedProcess[str]:
    return run(["terraform", f"-chdir={AZURE_DIR}", *extra], env=env, capture=capture)


def terraform_outputs(args: argparse.Namespace, env: dict[str, str]) -> dict[str, object]:
    result = terraform(args, env, "output", "-json", capture=True)
    try:
        payload = json.loads(result.stdout or "{}")
    except json.JSONDecodeError as exc:
        fail("Terraform output was not valid JSON: " + str(exc))
    if not isinstance(payload, dict):
        fail("Terraform output JSON must be an object")
    return payload


def output_value(outputs: dict[str, object], name: str) -> str:
    row = outputs.get(name)
    if not isinstance(row, dict):
        fail("missing Terraform output: " + name)
    value = str(row.get("value") or "").strip()
    if not value:
        fail("empty Terraform output: " + name)
    return value


def fingerprint(value: str) -> str:
    return hashlib.sha256(value.encode("utf-8")).hexdigest()[:16]


def require_evidence(path: Path, expected_schema: str, expected_status: str) -> dict[str, object]:
    try:
        payload = json.loads(path.read_text(encoding="utf-8"))
    except Exception as exc:
        fail(f"evidence unreadable {path.name}: {exc}")
    if not isinstance(payload, dict):
        fail(f"evidence must be an object: {path.name}")
    if payload.get("schema") != expected_schema or payload.get("status") != expected_status:
        fail(f"evidence did not reach required state: {path.name}")
    if payload.get("containsSecrets") is not False:
        fail(f"evidence must explicitly be secret-free: {path.name}")
    return payload


def role_assignment_present(rows: object, assignment_id: str) -> bool:
    if not isinstance(rows, list):
        return False
    target = assignment_id.rstrip("/").lower()
    for row in rows:
        if isinstance(row, dict) and str(row.get("id") or "").rstrip("/").lower() == target:
            return True
    return False


def create_temporary_aks_verification_role(principal_object_id: str, cluster_id: str) -> str:
    payload = run_json([
        "az", "role", "assignment", "create",
        "--assignee-object-id", principal_object_id,
        "--assignee-principal-type", "ServicePrincipal",
        "--role", AKS_VERIFY_ROLE,
        "--scope", cluster_id,
        "-o", "json",
    ])
    if not isinstance(payload, dict):
        fail("temporary AKS verification RBAC assignment returned a non-object response")
    assignment_id = str(payload.get("id") or "").strip()
    if not assignment_id:
        fail("temporary AKS verification RBAC assignment returned no assignment id")
    return assignment_id


def wait_for_aks_verification_access(rg: str, cluster: str) -> None:
    deadline = time.monotonic() + AKS_RBAC_PROPAGATION_TIMEOUT_SECONDS
    command = (
        "kubectl auth can-i create namespaces && "
        "kubectl auth can-i get namespaces && "
        "kubectl auth can-i create peerauthentications.security.istio.io -n depulse-dev"
    )
    last_logs = ""
    while time.monotonic() < deadline:
        payload = run_json([
            "az", "aks", "command", "invoke",
            "--resource-group", rg,
            "--name", cluster,
            "--command", command,
            "-o", "json",
        ])
        if aks_remote_succeeded(payload):
            logs = str((payload or {}).get("logs") or "").lower() if isinstance(payload, dict) else ""
            if logs.count("yes") >= 3:
                return
        if isinstance(payload, dict):
            last_logs = str(payload.get("logs") or "").strip()[:1000]
        time.sleep(AKS_RBAC_POLL_SECONDS)
    fail("temporary AKS verification RBAC did not propagate within 300s" + (f"; last logs={last_logs}" if last_logs else ""))


def managed_istio_diagnostics(rg: str, cluster: str, istio_revision: str, environment: str = "dev") -> str:
    namespace = f"depulse-{environment}"
    profile_text = "<service mesh profile unavailable>"
    try:
        profile = run_json([
            "az", "aks", "show", "--resource-group", rg, "--name", cluster,
            "--query", "serviceMeshProfile", "-o", "json",
        ])
        profile_text = json.dumps(profile, sort_keys=True)
    except BaseException as exc:
        profile_text = "profile collection failed: " + str(exc)[:1200]

    service = f"istiod-{istio_revision}"
    command = (
        "printf '%s\\n' '=== MANAGED ISTIO SYSTEM PODS ==='; "
        "kubectl get pods -n aks-istio-system -o wide || true; "
        "printf '%s\\n' '=== MANAGED ISTIO SYSTEM DEPLOYMENTS/HPA/DAEMONSETS ==='; "
        "kubectl get deployments,hpa,daemonsets -n aks-istio-system -o wide || true; "
        "printf '%s\\n' '=== MANAGED ISTIO SYSTEM REPLICASETS ==='; "
        "kubectl get replicasets -n aks-istio-system -l app=istiod -o wide || true; "
        f"printf '%s\\n' '=== {service} SERVICE/ENDPOINTS ==='; "
        f"kubectl get service,endpoints {service} -n aks-istio-system -o wide || true; "
        "printf '%s\\n' '=== MANAGED ISTIO INGRESS WORKLOADS ==='; "
        "kubectl get pods,deployments,hpa,services,endpoints -n aks-istio-ingress -o wide || true; "
        "printf '%s\\n' '=== DE.PULSE WORKLOAD TRUST OBJECTS ==='; "
        f"kubectl get namespace {namespace} --show-labels || true; "
        f"kubectl get serviceaccount depulse-web-{environment} -n {namespace} "
        "-o custom-columns='NAME:.metadata.name,WI_CLIENT:.metadata.annotations.azure\\.workload\\.identity/client-id' || true; "
        f"kubectl get networkpolicies -n {namespace} -o name || true; "
        f"kubectl get peerauthentications.security.istio.io,authorizationpolicies.security.istio.io,sidecars.networking.istio.io,serviceentries.networking.istio.io -n {namespace} -o name || true; "
        "printf '%s\\n' '=== AKS NODES ==='; "
        "kubectl get nodes -o custom-columns='NAME:.metadata.name,READY:.status.conditions[?(@.type==\"Ready\")].status,CPU:.status.allocatable.cpu,MEMORY:.status.allocatable.memory' || true; "
        "printf '%s\\n' '=== RECENT ISTIO SYSTEM EVENTS ==='; "
        "kubectl get events -n aks-istio-system --sort-by=.lastTimestamp | tail -n 40 || true; "
        "printf '%s\\n' '=== RECENT ISTIO INGRESS EVENTS ==='; "
        "kubectl get events -n aks-istio-ingress --sort-by=.lastTimestamp | tail -n 40 || true; "
        "printf '%s\\n' '=== RECENT WORKLOAD EVENTS ==='; "
        f"kubectl get events -n {namespace} --sort-by=.lastTimestamp | tail -n 40 || true"
    )
    try:
        payload = run_json([
            "az", "aks", "command", "invoke",
            "--resource-group", rg,
            "--name", cluster,
            "--command", command,
            "-o", "json",
        ])
        logs = str(payload.get("logs") or "").strip() if isinstance(payload, dict) else ""
    except BaseException as exc:
        logs = "Kubernetes diagnostics collection failed: " + str(exc)[:2000]
    return ("=== AZURE SERVICE MESH PROFILE ===\n" + profile_text + "\n" + logs)[:16000]


def managed_istio_ready_command(istio_revision: str) -> str:
    deployment = f"istiod-{istio_revision}"
    return (
        f"rollout=\"$(kubectl rollout status deployment/{deployment} -n aks-istio-system --timeout=5s 2>&1)\" || "
        "{ printf '%s\\n' \"$rollout\"; exit 1; }; "
        f"endpoint=\"$(kubectl get endpoints {deployment} -n aks-istio-system "
        "-o jsonpath='{.subsets[0].addresses[0].ip}' 2>/dev/null || true)\"; "
        f"desired=\"$(kubectl get deployment {deployment} -n aks-istio-system -o jsonpath='{{.status.replicas}}' 2>/dev/null || true)\"; "
        f"updated=\"$(kubectl get deployment {deployment} -n aks-istio-system -o jsonpath='{{.status.updatedReplicas}}' 2>/dev/null || true)\"; "
        f"ready=\"$(kubectl get deployment {deployment} -n aks-istio-system -o jsonpath='{{.status.readyReplicas}}' 2>/dev/null || true)\"; "
        f"available=\"$(kubectl get deployment {deployment} -n aks-istio-system -o jsonpath='{{.status.availableReplicas}}' 2>/dev/null || true)\"; "
        f"unavailable=\"$(kubectl get deployment {deployment} -n aks-istio-system -o jsonpath='{{.status.unavailableReplicas}}' 2>/dev/null || true)\"; "
        "test -n \"$endpoint\" && test \"${desired:-0}\" -ge 1 && "
        "test \"${updated:-0}\" = \"$desired\" && test \"${ready:-0}\" = \"$desired\" && "
        "test \"${available:-0}\" = \"$desired\" && test \"${unavailable:-0}\" = \"0\" && "
        "printf 'ISTIO_READY deployment=%s endpoint=%s desired=%s updated=%s ready=%s available=%s unavailable=%s\\n' "
        f"'{deployment}' \"$endpoint\" \"$desired\" \"$updated\" \"$ready\" \"$available\" \"${{unavailable:-0}}\""
    )


def wait_for_managed_istio_ready(rg: str, cluster: str, istio_revision: str) -> None:
    command = managed_istio_ready_command(istio_revision)
    deadline = time.monotonic() + ISTIO_READY_TIMEOUT_SECONDS
    last_logs = ""
    stable_passes = 0
    while time.monotonic() < deadline:
        payload = run_json([
            "az", "aks", "command", "invoke",
            "--resource-group", rg,
            "--name", cluster,
            "--command", command,
            "-o", "json",
        ])
        if aks_remote_succeeded(payload):
            logs = str((payload or {}).get("logs") or "").strip() if isinstance(payload, dict) else ""
            if "ISTIO_READY" in logs:
                stable_passes += 1
                last_logs = logs[:1200]
                if stable_passes >= ISTIO_READY_STABLE_PASSES:
                    return
            else:
                stable_passes = 0
                last_logs = logs[:1200]
        else:
            stable_passes = 0
            if isinstance(payload, dict):
                last_logs = str(payload.get("logs") or "").strip()[:1200]
        time.sleep(ISTIO_READY_POLL_SECONDS)
    diagnostic = managed_istio_diagnostics(rg, cluster, istio_revision)
    fail(
        f"managed Istio revision {istio_revision} did not complete its rollout and remain ready for "
        f"{ISTIO_READY_STABLE_PASSES} consecutive checks within {ISTIO_READY_TIMEOUT_SECONDS}s; "
        f"last readiness logs={last_logs or '<none>'}; diagnostics={diagnostic}"
    )


def external_ingress_ready_command() -> str:
    service = "aks-istio-ingressgateway-external"
    return (
        f"svc=\"$(kubectl get service {service} -n aks-istio-ingress -o name 2>/dev/null || true)\"; "
        f"port=\"$(kubectl get service {service} -n aks-istio-ingress -o jsonpath='{{.spec.ports[?(@.port==443)].port}}' 2>/dev/null || true)\"; "
        f"ip=\"$(kubectl get service {service} -n aks-istio-ingress -o jsonpath='{{.status.loadBalancer.ingress[0].ip}}' 2>/dev/null || true)\"; "
        f"endpoint=\"$(kubectl get endpoints {service} -n aks-istio-ingress -o jsonpath='{{.subsets[0].addresses[0].ip}}' 2>/dev/null || true)\"; "
        f"if test -z \"$endpoint\"; then endpoint=\"$(kubectl get endpointslices.discovery.k8s.io -n aks-istio-ingress -l kubernetes.io/service-name={service} -o jsonpath='{{.items[0].endpoints[0].addresses[0]}}' 2>/dev/null || true)\"; fi; "
        "ready=\"$(kubectl get pods -n aks-istio-ingress -l istio=aks-istio-ingressgateway-external "
        "-o jsonpath='{range .items[*]}{.status.containerStatuses[0].ready}{\"\\n\"}{end}' 2>/dev/null | grep -c true || true)\"; "
        f"test \"$svc\" = \"service/{service}\" && test \"$port\" = \"443\" && test -n \"$ip\" && test -n \"$endpoint\" && test \"${{ready:-0}}\" -ge 1 && "
        "printf 'INGRESS_READY service=%s httpsPort=%s ip=%s endpoint=%s readyPods=%s\\n' \"$svc\" \"$port\" \"$ip\" \"$endpoint\" \"$ready\""
    )


def reconcile_external_ingress_gateway(rg: str, cluster: str) -> None:
    result = run([
        "az", "aks", "mesh", "enable-ingress-gateway",
        "--resource-group", rg,
        "--name", cluster,
        "--ingress-gateway-type", "external",
        "--only-show-errors",
        "-o", "none",
    ], capture=True, check=False)
    if result.returncode == 0:
        return
    detail = (result.stdout or "").strip()
    lowered = detail.lower()
    if "already" in lowered and ("enabled" in lowered or "exist" in lowered):
        return
    fail("managed external Istio ingress reconciliation failed" + (": " + detail[:4000] if detail else ""))


def wait_for_managed_external_ingress_ready(rg: str, cluster: str, istio_revision: str, environment: str) -> None:
    command = external_ingress_ready_command()
    deadline = time.monotonic() + INGRESS_READY_TIMEOUT_SECONDS
    started = time.monotonic()
    stable_passes = 0
    reconciled = False
    last_logs = ""
    while time.monotonic() < deadline:
        payload = run_json([
            "az", "aks", "command", "invoke",
            "--resource-group", rg,
            "--name", cluster,
            "--command", command,
            "-o", "json",
        ])
        if aks_remote_succeeded(payload):
            logs = str((payload or {}).get("logs") or "").strip() if isinstance(payload, dict) else ""
            if "INGRESS_READY" in logs:
                stable_passes += 1
                last_logs = logs[:1200]
                if stable_passes >= INGRESS_READY_STABLE_PASSES:
                    return
            else:
                stable_passes = 0
                last_logs = logs[:1200]
        else:
            stable_passes = 0
            if isinstance(payload, dict):
                last_logs = str(payload.get("logs") or "").strip()[:1200]
        if not reconciled and time.monotonic() - started >= INGRESS_RECONCILE_AFTER_SECONDS:
            reconcile_external_ingress_gateway(rg, cluster)
            reconciled = True
            stable_passes = 0
        time.sleep(INGRESS_READY_POLL_SECONDS)
    diagnostic = managed_istio_diagnostics(rg, cluster, istio_revision, environment)
    fail(
        "managed external Istio ingress did not remain ready for "
        f"{INGRESS_READY_STABLE_PASSES} consecutive checks within {INGRESS_READY_TIMEOUT_SECONDS}s; "
        f"reconciled={reconciled}; last readiness logs={last_logs or '<none>'}; diagnostics={diagnostic}"
    )


def delete_temporary_aks_verification_role(assignment_id: str, cluster_id: str, rg: str, cluster: str) -> None:
    run([
        "az", "role", "assignment", "delete",
        "--ids", assignment_id,
        "--only-show-errors",
    ], capture=True)

    rows = run_json([
        "az", "role", "assignment", "list",
        "--scope", cluster_id,
        "-o", "json",
    ])
    if role_assignment_present(rows, assignment_id):
        fail("temporary AKS verification RBAC assignment still exists after delete")

    deadline = time.monotonic() + AKS_RBAC_PROPAGATION_TIMEOUT_SECONDS
    command = "kubectl auth can-i create namespaces"
    while time.monotonic() < deadline:
        payload = run_json([
            "az", "aks", "command", "invoke",
            "--resource-group", rg,
            "--name", cluster,
            "--command", command,
            "-o", "json",
        ])
        if not aks_remote_succeeded(payload):
            return
        logs = str((payload or {}).get("logs") or "").strip().lower() if isinstance(payload, dict) else ""
        if logs == "no" or "\nno" in logs:
            return
        time.sleep(AKS_RBAC_POLL_SECONDS)
    fail("temporary AKS verification RBAC remained effective after deletion timeout")


def write_failure_evidence(
    path: Path,
    *,
    candidate_sha: str,
    environment: str,
    istio_revision: str,
    verification_stage: str,
    verification_error: BaseException | None,
    cleanup_error: BaseException | None,
    istio_diagnostics: str | None,
) -> None:
    payload = {
        "schema": "DE.PULSE-HOST013-AZURE-FAILURE-1",
        "candidateSha": candidate_sha,
        "environment": environment,
        "managedIstioRevision": istio_revision,
        "verificationStage": verification_stage,
        "verificationFailure": str(verification_error)[:12000] if verification_error is not None else None,
        "cleanupFailure": str(cleanup_error)[:4000] if cleanup_error is not None else None,
        "managedIstioDiagnostics": istio_diagnostics[:16000] if istio_diagnostics else None,
        "temporaryKubernetesAdminRemoved": cleanup_error is None,
        "containsSecrets": False,
        "status": "FAIL",
    }
    path.write_text(json.dumps(payload, indent=2, sort_keys=True) + "\n", encoding="utf-8")


def self_test() -> None:
    fixture = {
        "revisions": [
            {"name": "asm-1-26"},
            {"revision": "asm-1-27"},
            {"note": "not-a-revision"},
        ]
    }
    if choose_revision(fixture) != "asm-1-27":
        fail("managed Istio revision selection self-test failed")
    try:
        choose_revision({"revisions": []})
    except SystemExit:
        pass
    else:
        fail("empty managed Istio revision set did not fail closed")
    if output_value({"x": {"value": "abc"}}, "x") != "abc":
        fail("Terraform output extraction self-test failed")
    try:
        output_value({}, "missing")
    except SystemExit:
        pass
    else:
        fail("missing Terraform output did not fail closed")
    if require_aks_command_success({"exitCode": 0, "logs": "ok"}, "self-test").get("exitCode") != 0:
        fail("AKS remote command success self-test failed")
    try:
        require_aks_command_success({"exitCode": 1, "logs": "synthetic kubectl failure"}, "self-test")
    except SystemExit as exc:
        if "synthetic kubectl failure" not in str(exc):
            fail("AKS remote command failure diagnostics were not preserved")
    else:
        fail("AKS remote command nonzero exit did not fail closed")
    assignment = "/subscriptions/x/providers/Microsoft.Authorization/roleAssignments/abc"
    if not role_assignment_present([{"id": assignment}], assignment + "/"):
        fail("temporary AKS RBAC assignment presence self-test failed")
    if role_assignment_present([], assignment):
        fail("temporary AKS RBAC assignment absence self-test failed")
    if ISTIO_READY_STABLE_PASSES < 2 or INGRESS_READY_STABLE_PASSES < 2:
        fail("managed Istio and ingress readiness must require consecutive passes")
    istio_command = managed_istio_ready_command("asm-1-30")
    for token in ("kubectl rollout status deployment/istiod-asm-1-30", "updatedReplicas", "readyReplicas", "availableReplicas", "unavailableReplicas", "ISTIO_READY"):
        if token not in istio_command:
            fail("managed Istio rollout-readiness self-test missing " + token)
    ingress_command = external_ingress_ready_command()
    for token in ("aks-istio-ingressgateway-external", "httpsPort=%s", "INGRESS_READY", "endpointslices.discovery.k8s.io"):
        if token not in ingress_command:
            fail("managed ingress readiness self-test missing " + token)
    print("HOST-013/014 Azure operator self-test: PASS")


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--self-test", action="store_true")
    parser.add_argument("--confirmation", default="")
    parser.add_argument("--candidate-sha", default="")
    parser.add_argument("--client-id", default="")
    parser.add_argument("--tenant-id", default="")
    parser.add_argument("--subscription-id", default="")
    parser.add_argument("--state-resource-group", default="")
    parser.add_argument("--state-storage-account", default="")
    parser.add_argument("--state-container", default="")
    parser.add_argument("--environment", default="dev", choices=["dev"])
    parser.add_argument("--location", default="canadacentral")
    parser.add_argument("--evidence-dir", default=".depulse-host013-azure")
    args = parser.parse_args()
    if args.self_test:
        self_test()
        return 0
    if args.confirmation != CONFIRMATION:
        fail(f"explicit confirmation must equal {CONFIRMATION}")
    required = (
        "candidate_sha", "client_id", "tenant_id", "subscription_id",
        "state_resource_group", "state_storage_account", "state_container",
    )
    for name in required:
        if not str(getattr(args, name)).strip():
            fail("missing required non-secret identifier --" + name.replace("_", "-"))
    head = run(["git", "rev-parse", "HEAD"], capture=True).stdout.strip()
    if head != args.candidate_sha:
        fail(f"exact-head mismatch: repository={head} requested={args.candidate_sha}")

    account = run_json(["az", "account", "show", "-o", "json"])
    if not isinstance(account, dict) or str(account.get("id") or "") != args.subscription_id:
        fail("active Azure subscription does not match requested subscription")

    revision_payload = run_json(["az", "aks", "mesh", "get-revisions", "--location", args.location, "-o", "json"])
    istio_revision = choose_revision(revision_payload)
    env = terraform_env(args, istio_revision)
    state_key = f"host013/{args.environment}/terraform.tfstate"

    terraform(args, env, "fmt", "-check", "-diff")
    terraform(
        args, env, "init", "-input=false", "-reconfigure",
        f"-backend-config=resource_group_name={args.state_resource_group}",
        f"-backend-config=storage_account_name={args.state_storage_account}",
        f"-backend-config=container_name={args.state_container}",
        f"-backend-config=key={state_key}",
    )
    terraform(args, env, "validate")

    evidence_dir = ROOT / args.evidence_dir
    evidence_dir.mkdir(parents=True, exist_ok=True)
    plan_path = evidence_dir / "host013-azure.tfplan"
    terraform(args, env, "plan", "-input=false", "-out=" + str(plan_path))
    terraform(args, env, "apply", "-input=false", "-auto-approve", str(plan_path))

    outputs = terraform_outputs(args, env)
    workload_client_id = output_value(outputs, "workload_identity_client_id")
    workload_subject = output_value(outputs, "workload_identity_subject")
    operator_client_id = output_value(outputs, "operator_identity_client_id")
    operator_object_id = output_value(outputs, "operator_identity_object_id")
    cluster_id = output_value(outputs, "aks_cluster_id")
    if operator_client_id.lower() != args.client_id.lower():
        fail("Terraform authenticated operator client ID does not match requested GitHub OIDC client ID")
    expected_subject = f"system:serviceaccount:depulse-{args.environment}:depulse-web-{args.environment}"
    if workload_subject != expected_subject:
        fail(f"workload identity subject drift: expected={expected_subject} actual={workload_subject}")

    rg = f"rg-depulse-{args.environment}"
    cluster = f"aks-depulse-{args.environment}"
    manifest_path = evidence_dir / "host013-aks-managed-trust.yaml"
    run([
        sys.executable, str(RENDERER),
        "--environment", args.environment,
        "--mesh-profile", "aks-managed",
        "--istio-revision", istio_revision,
        "--workload-identity-client-id", workload_client_id,
        "--output", str(manifest_path),
    ])

    temporary_assignment_id = create_temporary_aks_verification_role(operator_object_id, cluster_id)
    verification_error: BaseException | None = None
    verification_stage = "temporary-aks-rbac-propagation"
    evidence_path = evidence_dir / "host013-azure-live-evidence.json"
    traffic_path = evidence_dir / "host013-azure-traffic-evidence.json"
    istio_diagnostics: str | None = None
    try:
        wait_for_aks_verification_access(rg, cluster)

        verification_stage = "managed-istio-readiness"
        wait_for_managed_istio_ready(rg, cluster, istio_revision)

        verification_stage = "managed-external-ingress-readiness"
        wait_for_managed_external_ingress_ready(rg, cluster, istio_revision, args.environment)

        verification_stage = "canonical-trust-manifest-apply"
        apply_result = run_json([
            "az", "aks", "command", "invoke", "--resource-group", rg,
            "--name", cluster, "--command", "kubectl apply -f host013-aks-managed-trust.yaml",
            "--file", str(manifest_path), "-o", "json",
        ])
        require_aks_command_success(apply_result, "canonical trust manifest apply")

        verification_stage = "configuration-and-identity-evidence"
        try:
            refresh_azure_cli_oidc()
        except RuntimeError as exc:
            fail("Azure CLI OIDC refresh before live evidence failed: " + str(exc))
        run([
            sys.executable, str(LIVE_EVIDENCE),
            "--subscription-id", args.subscription_id,
            "--resource-group", rg,
            "--cluster-name", cluster,
            "--environment", args.environment,
            "--location", args.location,
            "--candidate-sha", args.candidate_sha,
            "--expected-istio-revision", istio_revision,
            "--workload-identity-client-id", workload_client_id,
            "--output", str(evidence_path),
        ], capture=True)
        require_evidence(evidence_path, "DE.PULSE-HOST013-AZURE-LIVE-EVIDENCE-1", "PASS_CONFIGURATION_AND_IDENTITY")

        verification_stage = "live-traffic-and-adverse-evidence"
        try:
            refresh_azure_cli_oidc()
        except RuntimeError as exc:
            fail("Azure CLI OIDC refresh before traffic evidence failed: " + str(exc))
        run([
            sys.executable, str(TRAFFIC_PROBE),
            "--resource-group", rg,
            "--cluster-name", cluster,
            "--environment", args.environment,
            "--candidate-sha", args.candidate_sha,
            "--istio-revision", istio_revision,
            "--output", str(traffic_path),
        ], capture=True)
        traffic = require_evidence(traffic_path, "DE.PULSE-HOST013-AZURE-TRAFFIC-EVIDENCE-1", "PASS")
        traffic_checks = traffic.get("checks")
        if not isinstance(traffic_checks, dict) or not traffic_checks or not all(value is True for value in traffic_checks.values()):
            fail("traffic probe evidence is incomplete")
    except BaseException as exc:
        verification_error = exc
        try:
            istio_diagnostics = managed_istio_diagnostics(rg, cluster, istio_revision, args.environment)
        except BaseException as diagnostic_exc:
            istio_diagnostics = "managed Istio diagnostics collection failed: " + str(diagnostic_exc)[:2000]

    cleanup_error: BaseException | None = None
    try:
        delete_temporary_aks_verification_role(temporary_assignment_id, cluster_id, rg, cluster)
    except BaseException as exc:
        cleanup_error = exc

    if verification_error is not None or cleanup_error is not None:
        write_failure_evidence(
            evidence_dir / "host013-azure-failure-evidence.json",
            candidate_sha=args.candidate_sha,
            environment=args.environment,
            istio_revision=istio_revision,
            verification_stage=verification_stage,
            verification_error=verification_error,
            cleanup_error=cleanup_error,
            istio_diagnostics=istio_diagnostics,
        )
    if cleanup_error is not None:
        detail = str(cleanup_error)
        if verification_error is not None:
            detail += "; original verification failure=" + str(verification_error)
        fail("temporary AKS verification RBAC cleanup failed: " + detail)
    if verification_error is not None:
        raise verification_error

    drift = subprocess.run(
        ["terraform", f"-chdir={AZURE_DIR}", "plan", "-input=false", "-detailed-exitcode"],
        cwd=ROOT, env=env, text=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT, check=False,
    )
    if drift.returncode != 0:
        fail(f"post-verification Terraform drift check returned {drift.returncode}; expected 0")

    summary = {
        "schema": "DE.PULSE-HOST013-AZURE-OPERATOR-1",
        "candidateSha": args.candidate_sha,
        "environment": args.environment,
        "location": args.location,
        "managedIstioRevision": istio_revision,
        "managedIstioReadinessProved": True,
        "managedIstioReadinessStablePasses": ISTIO_READY_STABLE_PASSES,
        "managedExternalIngressReadinessProved": True,
        "managedExternalIngressReadinessStablePasses": INGRESS_READY_STABLE_PASSES,
        "operatorIdentityFingerprint": fingerprint(operator_object_id),
        "temporaryKubernetesAdminRole": AKS_VERIFY_ROLE,
        "temporaryKubernetesAdminRemoved": True,
        "workloadIdentityClientFingerprint": fingerprint(workload_client_id),
        "workloadIdentitySubject": workload_subject,
        "stateKey": state_key,
        "terraformPostApplyDriftExitCode": drift.returncode,
        "configurationIdentityEvidence": evidence_path.name,
        "trafficEvidence": traffic_path.name,
        "containsSecrets": False,
        "status": "PASS",
    }
    (evidence_dir / "host013-azure-operator-summary.json").write_text(
        json.dumps(summary, indent=2, sort_keys=True) + "\n", encoding="utf-8"
    )
    plan_path.unlink(missing_ok=True)
    manifest_path.unlink(missing_ok=True)
    print("HOST-013/014 Azure operator: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
