#!/usr/bin/env python3
"""Fail-closed repository gate for the Azure AKS HOST-013..014 adapter."""
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[2]
AZ = ROOT / "governance" / "hosted-infrastructure" / "azure"
REQUIRED = ["README.md", "versions.tf", "backend.tf", "variables.tf", "main.tf", "outputs.tf"]
LIVE_EVIDENCE = ROOT / "tools" / "ci" / "host013_azure_live_evidence.py"
TRAFFIC_PROBE = ROOT / "tools" / "ci" / "host013_azure_traffic_probe.py"
OPERATOR = ROOT / "tools" / "ci" / "host013_azure_operator.py"
RENDERER = ROOT / "tools" / "hosted" / "render_kubernetes_trust.py"


def fail(msg: str) -> None:
    print("FAIL:", msg)
    raise SystemExit(1)


def require(text: str, pattern: str, label: str) -> None:
    if not re.search(pattern, text, re.MULTILINE):
        fail("Azure hosted infrastructure missing " + label)


def run_self_test(path: Path, label: str) -> None:
    if not path.is_file():
        fail("missing " + label)
    result = subprocess.run([sys.executable, str(path), "--self-test"], cwd=ROOT, check=False)
    if result.returncode != 0:
        fail(label + " self-test failed")


def main() -> int:
    for name in REQUIRED:
        if not (AZ / name).is_file():
            fail("missing Azure hosted infrastructure file: " + name)
    for path in (LIVE_EVIDENCE, TRAFFIC_PROBE, OPERATOR, RENDERER):
        if not path.is_file():
            fail("missing Azure hosted trust owner: " + path.name)

    main_tf = (AZ / "main.tf").read_text(encoding="utf-8")
    vars_tf = (AZ / "variables.tf").read_text(encoding="utf-8")
    versions_tf = (AZ / "versions.tf").read_text(encoding="utf-8")
    backend_tf = (AZ / "backend.tf").read_text(encoding="utf-8")
    readme = (AZ / "README.md").read_text(encoding="utf-8")
    renderer = RENDERER.read_text(encoding="utf-8")
    probe_text = TRAFFIC_PROBE.read_text(encoding="utf-8")

    require(versions_tf, r'source\s*=\s*"hashicorp/azurerm"', "AzureRM provider")
    require(versions_tf, r'version\s*=\s*"=\s*4\.81\.0"', "exact live-tested AzureRM 4.81.0 pin")
    require(backend_tf, r'backend\s+"azurerm"', "AzureRM remote state backend")
    require(backend_tf, r'use_oidc\s*=\s*true', "OIDC state authentication")
    require(backend_tf, r'use_azuread_auth\s*=\s*true', "Microsoft Entra state data-plane authentication")
    require(main_tf, r'private_cluster_enabled\s*=\s*var\.private_cluster_enabled', "private AKS control-plane binding")
    require(vars_tf, r'condition\s*=\s*var\.private_cluster_enabled', "fail-closed private-cluster validation")
    require(main_tf, r'local_account_disabled\s*=\s*true', "disabled local AKS accounts")
    require(main_tf, r'azure_active_directory_role_based_access_control\s*\{', "AKS-managed Microsoft Entra integration")
    require(main_tf, r'azure_active_directory_role_based_access_control\s*\{[\s\S]*?tenant_id\s*=\s*var\.tenant_id[\s\S]*?azure_rbac_enabled\s*=\s*true', "tenant-bound Azure RBAC")
    require(main_tf, r'oidc_issuer_enabled\s*=\s*true', "OIDC issuer")
    require(main_tf, r'workload_identity_enabled\s*=\s*true', "Azure Workload Identity")
    require(main_tf, r'role_based_access_control_enabled\s*=\s*true', "Kubernetes RBAC")
    require(main_tf, r'azure_policy_enabled\s*=\s*true', "Azure Policy")
    require(main_tf, r'only_critical_addons_enabled\s*=\s*false', "schedulable single-pool dev verification workloads")
    require(main_tf, r'network_policy\s*=\s*"azure"', "network policy")
    require(main_tf, r'service_mesh_profile\s*\{[\s\S]*?mode\s*=\s*"Istio"[\s\S]*?revisions\s*=\s*var\.istio_revisions', "managed AKS Istio profile")
    require(main_tf, r'external_ingress_gateway_enabled\s*=\s*true', "managed external Istio ingress gateway")
    require(main_tf, r'internal_ingress_gateway_enabled\s*=\s*false', "single governed external ingress posture")
    require(vars_tf, r'variable\s+"istio_revisions"', "governed Istio revision input")
    require(vars_tf, r'asm-\[0-9\]\+\-\[0-9\]\+', "Azure managed Istio revision format validation")
    require(main_tf, r'resource\s+"azurerm_federated_identity_credential"\s+"workload"', "workload federated identity credential")
    require(main_tf, r'audience\s*=\s*\["api://AzureADTokenExchange"\]', "workload token-exchange audience")
    require(main_tf, r'subject\s*=\s*local\.workload_identity_subject', "canonical service-account federation subject")
    require(main_tf, r'key_vault_secrets_provider\s*\{', "Key Vault CSI integration hook")
    require(main_tf, r'role_definition_name\s*=\s*"Network Contributor"', "BYO-network role assignment")
    require(main_tf, r'scope\s*=\s*azurerm_virtual_network\.this\.id', "network role least-privilege scope")
    require(main_tf, r'host_gate\s*=\s*"HOST-013-HOST-014"', "HOST ownership tag")
    require(readme, r'HOST-013\.\.014 may move to VERIFIED only after a real Azure deployment', "truthful live-verification boundary")

    require(renderer, r'"aks-managed"', "AKS-managed mesh renderer profile")
    require(renderer, r'istio\.io/rev', "AKS managed revision namespace label")
    require(renderer, r'aks-istio-ingress', "AKS managed ingress namespace")
    require(renderer, r'aks-istio-ingressgateway-external', "AKS managed external ingress selector")
    require(renderer, r'control_plane_namespace\s*=\s*"aks-istio-system"', "AKS managed Istio control-plane namespace")
    require(renderer, r'port:\s*15012', "secure Istio XDS/CA NetworkPolicy egress")
    require(renderer, r'azure\.workload\.identity/client-id', "Kubernetes workload identity service-account annotation")

    require(probe_text, r'python:3\.13-alpine@sha256:[0-9a-f]{64}', "digest-pinned live probe image")
    require(probe_text, r'minProtocolVersion: TLSV1_2', "TLS 1.2 managed-edge minimum probe")
    require(probe_text, r'azure\.workload\.identity/use', "workload identity pod mutation probe")
    require(probe_text, r'WORKLOAD_IDENTITY_TOKEN_OK', "real workload identity token-exchange probe")
    require(probe_text, r'UNREGISTERED_EGRESS_DENIED', "unregistered egress adverse probe")
    require(probe_text, r'DIRECT_INGRESS_DENIED', "cross-environment/direct-ingress adverse probe")
    require(probe_text, r'PROBE_CLEANUP_VERIFIED', "verified ephemeral probe cleanup marker")
    require(probe_text, r'cleanupVerified.*True', "fail-closed cleanup evidence")
    require(probe_text, r'containsSecrets.*False', "secret-free traffic evidence")

    forbidden = ["client_secret", "password", "api_key", "MARKETDATA_TOKEN", "FINNHUB"]
    joined = "\n".join((AZ / name).read_text(encoding="utf-8") for name in REQUIRED)
    for token in forbidden:
        if token.lower() in joined.lower():
            fail("Azure IaC contains forbidden secret/token material marker: " + token)

    run_self_test(LIVE_EVIDENCE, "Azure live-evidence collector")
    run_self_test(TRAFFIC_PROBE, "Azure live traffic probe")
    run_self_test(OPERATOR, "Azure AKS operator")

    operator_text = OPERATOR.read_text(encoding="utf-8")
    require(operator_text, r'HOST013_AZURE_AKS_OPERATOR_DRILL', "explicit destructive/non-production operator acknowledgement")
    require(operator_text, r'choices=\["dev"\]', "dev-only operator scope")
    require(operator_text, r'ARM_USE_OIDC.*true', "operator OIDC Terraform authentication")
    require(operator_text, r'--mesh-profile", "aks-managed"', "operator AKS-managed rendering")
    require(operator_text, r'workload_identity_client_id', "operator workload identity output binding")
    require(operator_text, r'host013_azure_traffic_probe', "operator live traffic probe binding")
    require(operator_text, r'DE.PULSE-HOST013-AZURE-TRAFFIC-EVIDENCE-1', "operator traffic evidence validation")
    require(operator_text, r'plan.*-detailed-exitcode', "post-verification Terraform drift check")
    if "client-secret" in operator_text.lower() or "client_secret" in operator_text.lower():
        fail("Azure operator must not expose a client-secret path")

    print("PASS: Azure AKS HOST-013..014 adapter/operator is fail-closed, Entra-integrated, reproducibly pinned, schedulable, OIDC-state-backed, managed-Istio-correct, XDS-reachable, workload-identity-bound, live-traffic-tested, cleanup-verified and secret-free")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
