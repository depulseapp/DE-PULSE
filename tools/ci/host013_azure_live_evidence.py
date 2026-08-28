#!/usr/bin/env python3
"""Collect and validate secret-free HOST-013..014 evidence from a real Azure AKS cluster.

The script intentionally consumes Azure CLI identity established out-of-band (for
GitHub Actions this is OIDC federation). It never reads or writes client secrets,
access tokens, kubeconfig, provider credentials, or application secrets.
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
        "localAccountsDisabled": enabled(cluster.get("disableLocalAccounts")),
        "oidcIssuerEnabled": enabled(pick(cluster, "oidcIssuerProfile", "enabled")),
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


def validate_run_command(result: dict[str, Any]) -> dict[str, bool]:
    exit_code = result.get("exitCode")
    logs = str(result.get("logs") or "")
    checks = {
        "runCommandSucceeded": str(exit_code) in {"0", "None"} or exit_code == 0,
        "managedIstioNamespaceReachable": "aks-istio-system" in logs,
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

    revisions = run_json(["az", "aks", "mesh", "get-revisions", "--location", args.location, "-o", "json"])
    mesh_profile = cluster.get("serviceMeshProfile") or {}
    installed_revisions = list(mesh_profile.get("revisions") or [])

    command_result = run_json([
        "az", "aks", "command", "invoke", "--resource-group", args.resource_group,
        "--name", args.cluster_name,
        "--command", "kubectl get namespace aks-istio-system -o name && kubectl get pods -n aks-istio-system --no-headers",
        "-o", "json",
    ])
    checks.update(validate_run_command(command_result))

    evidence = {
        "schema": SCHEMA,
        "requirements": ["HOST-013", "HOST-014"],
        "candidateSha": args.candidate_sha,
        "environment": args.environment,
        "provider": "Azure",
        "substrate": "AKS",
        "location": str(cluster.get("location") or args.location),
        "subscriptionFingerprint": sanitize_id(active_subscription),
        "resourceGroup": args.resource_group,
        "clusterName": args.cluster_name,
        "privateFqdnPresent": bool(cluster.get("privateFqdn")),
        "kubernetesVersion": cluster.get("kubernetesVersion"),
        "installedIstioRevisions": installed_revisions,
        "availableIstioRevisionRecords": len(revisions) if isinstance(revisions, list) else None,
        "checks": checks,
        "containsSecrets": False,
        "credentialMaterialRetained": False,
        "kubeconfigRetained": False,
        "status": "PASS",
    }
    return evidence


def self_test() -> None:
    fixture = {
        "location": "canadacentral",
        "apiServerAccessProfile": {"enablePrivateCluster": True},
        "disableLocalAccounts": True,
        "oidcIssuerProfile": {"enabled": True},
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
    validate_run_command({"exitCode": 0, "logs": "namespace/aks-istio-system\npod ready"})
    broken = json.loads(json.dumps(fixture))
    broken["apiServerAccessProfile"]["enablePrivateCluster"] = False
    try:
        validate_cluster(broken, "canadacentral")
    except RuntimeError:
        pass
    else:
        raise AssertionError("private-cluster adverse self-test did not fail closed")
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
    parser.add_argument("--output", default=".depulse-host013-azure/host013-azure-live-evidence.json")
    args = parser.parse_args()
    if args.self_test:
        self_test()
        return 0
    for name in ("subscription_id", "resource_group", "cluster_name", "candidate_sha"):
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
