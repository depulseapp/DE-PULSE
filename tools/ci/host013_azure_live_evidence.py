#!/usr/bin/env python3
"""Collect and validate secret-free HOST-013..014 evidence from a real Azure AKS cluster.

The script consumes Azure CLI identity established out-of-band (GitHub Actions uses
OIDC federation). It never reads or writes client secrets, access tokens, kubeconfig,
provider credentials, or application secrets.
"""
from __future__ import annotations

import argparse
import hashlib
import json
from pathlib import Path
import subprocess
from typing import Any

SCHEMA = "DE.PULSE-HOST013-AZURE-LIVE-EVIDENCE-1"


def run_json(args: list[str]) -> Any:
    result = subprocess.run(args, text=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, check=False)
    if result.returncode != 0:
        raise RuntimeError(f"command failed ({result.returncode}): {' '.join(args)}\n{result.stderr.strip()}")
    return json.loads(result.stdout)


def pick(value: Any, *path: str) -> Any:
    current = value
    for key in path:
        if not isinstance(current, dict):
            return None
        current = current.get(key)
    return current


def enabled(value: Any) -> bool:
    return value is True or str(value).strip().lower() == "true"


def sanitize_id(value: str) -> str:
    return hashlib.sha256(value.encode("utf-8")).hexdigest()[:16]


def validate_cluster(cluster: dict[str, Any], expected_location: str) -> dict[str, bool]:
    addon = cluster.get("addonProfiles") or {}
    mesh = cluster.get("serviceMeshProfile") or {}
    identity = cluster.get("identity") or {}
    checks = {
        "privateControlPlane": enabled(pick(cluster, "apiServerAccessProfile", "enablePrivateCluster")),
        "privateFqdnPresent": bool(cluster.get("privateFqdn")),
        "localAccountsDisabled": enabled(cluster.get("disableLocalAccounts")),
        "oidcIssuerEnabled": enabled(pick(cluster, "oidcIssuerProfile", "enabled")),
        "oidcIssuerUrlPresent": bool(pick(cluster, "oidcIssuerProfile", "issuerUrl")),
        "workloadIdentityEnabled": enabled(pick(cluster, "securityProfile", "workloadIdentity", "enabled")),
        "azurePolicyEnabled": enabled(pick(addon, "azurepolicy", "enabled")),
        "keyVaultSecretsProviderEnabled": enabled(pick(addon, "azureKeyvaultSecretsProvider", "enabled")),
        "azureNetworkPolicy": str(pick(cluster, "networkProfile", "networkPolicy") or "").lower() == "azure",
        "managedIstioEnabled": str(mesh.get("mode") or "").lower() == "istio" and bool(mesh.get("revisions")),
        "userAssignedControlPlaneIdentity": "userassigned" in str(identity.get("type") or "").lower(),
        "locationMatches": str(cluster.get("location") or "").replace(" ", "").lower() == expected_location.replace(" ", "").lower(),
    }
    failed = [name for name, ok in checks.items() if not ok]
    if failed:
        raise RuntimeError("AKS hosted-trust checks failed: " + ", ".join(failed))
    return checks


def validate_federated_credential(
    rows: Any,
    *,
    expected_subject: str,
    expected_issuer: str,
) -> dict[str, bool]:
    records = rows if isinstance(rows, list) else []
    matched = None
    for row in records:
        if not isinstance(row, dict):
            continue
        if str(row.get("subject") or "") == expected_subject:
            matched = row
            break
    audiences = list((matched or {}).get("audiences") or []) if isinstance(matched, dict) else []
    checks = {
        "federatedCredentialPresent": matched is not None,
        "federatedCredentialIssuerMatchesCluster": str((matched or {}).get("issuer") or "") == expected_issuer,
        "federatedCredentialAudienceBound": "api://AzureADTokenExchange" in audiences,
    }
    failed = [name for name, ok in checks.items() if not ok]
    if failed:
        raise RuntimeError("Azure workload federated-credential checks failed: " + ", ".join(failed))
    return checks


def validate_run_command(
    result: dict[str, Any],
    *,
    expected_revision: str,
    expected_client_id: str,
) -> dict[str, bool]:
    exit_code = result.get("exitCode")
    logs = str(result.get("logs") or "")
    checks = {
        "runCommandExitCodePresent": exit_code is not None,
        "runCommandSucceeded": exit_code == 0 or str(exit_code).strip() == "0",
        "managedIstioSystemNamespaceReachable": "namespace/aks-istio-system" in logs,
        "managedIstioIngressNamespaceReachable": "namespace/aks-istio-ingress" in logs,
        "managedExternalIngressGatewayPresent": "service/aks-istio-ingressgateway-external" in logs,
        "workloadNamespaceRevisionMatches": f"REV={expected_revision}" in logs,
        "serviceAccountWorkloadIdentityClientMatches": f"CLIENT={expected_client_id}" in logs,
    }
    failed = [name for name, ok in checks.items() if not ok]
    if failed:
        raise RuntimeError("AKS private-cluster command evidence failed: " + ", ".join(failed))
    return checks


def collect(args: argparse.Namespace) -> dict[str, Any]:
    account = run_json(["az", "account", "show", "-o", "json"])
    active_subscription = str(account.get("id") or "")
    if args.subscription_id and active_subscription != args.subscription_id:
        raise RuntimeError("Azure CLI subscription does not match requested subscription")

    cluster = run_json([
        "az", "aks", "show", "--resource-group", args.resource_group,
        "--name", args.cluster_name, "-o", "json",
    ])
    checks = validate_cluster(cluster, args.location)
    mesh_profile = cluster.get("serviceMeshProfile") or {}
    installed_revisions = list(mesh_profile.get("revisions") or [])
    if args.expected_istio_revision not in installed_revisions:
        raise RuntimeError("expected managed Istio revision is not installed")

    oidc_issuer = str(pick(cluster, "oidcIssuerProfile", "issuerUrl") or "")
    expected_subject = f"system:serviceaccount:depulse-{args.environment}:depulse-web-{args.environment}"
    federated = run_json([
        "az", "identity", "federated-credential", "list",
        "--resource-group", args.resource_group,
        "--identity-name", f"id-depulse-{args.environment}-workload",
        "-o", "json",
    ])
    checks.update(validate_federated_credential(
        federated,
        expected_subject=expected_subject,
        expected_issuer=oidc_issuer,
    ))

    revisions = run_json(["az", "aks", "mesh", "get-revisions", "--location", args.location, "-o", "json"])
    namespace = f"depulse-{args.environment}"
    service_account = f"depulse-web-{args.environment}"
    command = (
        "kubectl get namespace aks-istio-system -o name && "
        "kubectl get namespace aks-istio-ingress -o name && "
        "kubectl get service aks-istio-ingressgateway-external -n aks-istio-ingress -o name && "
        f"printf 'REV=' && kubectl get namespace {namespace} -o jsonpath='{{.metadata.labels.istio\\.io/rev}}' && echo && "
        f"printf 'CLIENT=' && kubectl get serviceaccount {service_account} -n {namespace} "
        "-o jsonpath='{.metadata.annotations.azure\\.workload\\.identity/client-id}' && echo"
    )
    command_result = run_json([
        "az", "aks", "command", "invoke", "--resource-group", args.resource_group,
        "--name", args.cluster_name,
        "--command", command,
        "-o", "json",
    ])
    checks.update(validate_run_command(
        command_result,
        expected_revision=args.expected_istio_revision,
        expected_client_id=args.workload_identity_client_id,
    ))

    return {
        "schema": SCHEMA,
        "requirements": ["HOST-013", "HOST-014"],
        "candidateSha": args.candidate_sha,
        "environment": args.environment,
        "provider": "Azure",
        "substrate": "AKS",
        "location": str(cluster.get("location") or args.location),
        "subscriptionFingerprint": sanitize_id(active_subscription),
        "workloadIdentityClientFingerprint": sanitize_id(args.workload_identity_client_id),
        "workloadIdentitySubject": expected_subject,
        "resourceGroup": args.resource_group,
        "clusterName": args.cluster_name,
        "kubernetesVersion": cluster.get("kubernetesVersion"),
        "installedIstioRevisions": installed_revisions,
        "availableIstioRevisionRecords": len(revisions) if isinstance(revisions, list) else None,
        "checks": checks,
        "containsSecrets": False,
        "credentialMaterialRetained": False,
        "kubeconfigRetained": False,
        "trafficProbeStatus": "NOT_YET_EXECUTED",
        "status": "PASS_CONFIGURATION_AND_IDENTITY",
    }


def self_test() -> None:
    fixture = {
        "location": "canadacentral",
        "privateFqdn": "private.example.invalid",
        "apiServerAccessProfile": {"enablePrivateCluster": True},
        "disableLocalAccounts": True,
        "oidcIssuerProfile": {"enabled": True, "issuerUrl": "https://issuer.example.invalid/"},
        "securityProfile": {"workloadIdentity": {"enabled": True}},
        "addonProfiles": {
            "azurepolicy": {"enabled": True},
            "azureKeyvaultSecretsProvider": {"enabled": True},
        },
        "networkProfile": {"networkPolicy": "azure"},
        "serviceMeshProfile": {"mode": "Istio", "revisions": ["asm-1-27"]},
        "identity": {"type": "UserAssigned"},
    }
    checks = validate_cluster(fixture, "canadacentral")
    assert all(checks.values())
    validate_federated_credential(
        [{"subject": "system:serviceaccount:depulse-dev:depulse-web-dev", "issuer": "https://issuer.example.invalid/", "audiences": ["api://AzureADTokenExchange"]}],
        expected_subject="system:serviceaccount:depulse-dev:depulse-web-dev",
        expected_issuer="https://issuer.example.invalid/",
    )
    validate_run_command(
        {"exitCode": 0, "logs": "namespace/aks-istio-system\nnamespace/aks-istio-ingress\nservice/aks-istio-ingressgateway-external\nREV=asm-1-27\nCLIENT=11111111-1111-1111-1111-111111111111\n"},
        expected_revision="asm-1-27",
        expected_client_id="11111111-1111-1111-1111-111111111111",
    )

    adverse = [
        lambda: validate_cluster({**fixture, "privateFqdn": ""}, "canadacentral"),
        lambda: validate_federated_credential([], expected_subject="x", expected_issuer="y"),
        lambda: validate_run_command(
            {"logs": "namespace/aks-istio-system"},
            expected_revision="asm-1-27",
            expected_client_id="11111111-1111-1111-1111-111111111111",
        ),
    ]
    for index, target in enumerate(adverse):
        try:
            target()
        except RuntimeError:
            pass
        else:
            raise AssertionError(f"adverse self-test {index} did not fail closed")
    print("HOST-013/014 Azure live-evidence self-test: PASS")


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--self-test", action="store_true")
    parser.add_argument("--subscription-id", default="")
    parser.add_argument("--resource-group", default="")
    parser.add_argument("--cluster-name", default="")
    parser.add_argument("--environment", default="dev", choices=["dev", "test", "stage", "prod"])
    parser.add_argument("--location", default="canadacentral")
    parser.add_argument("--candidate-sha", default="")
    parser.add_argument("--expected-istio-revision", default="")
    parser.add_argument("--workload-identity-client-id", default="")
    parser.add_argument("--output", default=".depulse-host013-azure/host013-azure-live-evidence.json")
    args = parser.parse_args()
    if args.self_test:
        self_test()
        return 0
    for name in (
        "subscription_id", "resource_group", "cluster_name", "candidate_sha",
        "expected_istio_revision", "workload_identity_client_id",
    ):
        if not getattr(args, name):
            parser.error(f"--{name.replace('_', '-')} is required")
    evidence = collect(args)
    path = Path(args.output)
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(evidence, indent=2, sort_keys=True) + "\n", encoding="utf-8")
    print(f"HOST-013/014 Azure live evidence: PASS -> {path}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
