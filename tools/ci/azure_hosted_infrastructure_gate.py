#!/usr/bin/env python3
"""Fail-closed repository gate for the Azure AKS HOST-013..014 adapter."""
from pathlib import Path
import re

ROOT = Path(__file__).resolve().parents[2]
AZ = ROOT / "governance" / "hosted-infrastructure" / "azure"
REQUIRED = ["README.md", "versions.tf", "variables.tf", "main.tf", "outputs.tf"]


def fail(msg: str) -> None:
    print("FAIL:", msg)
    raise SystemExit(1)


def require(text: str, pattern: str, label: str) -> None:
    if not re.search(pattern, text, re.MULTILINE):
        fail("Azure hosted infrastructure missing " + label)


def main() -> int:
    for name in REQUIRED:
        if not (AZ / name).is_file():
            fail("missing Azure hosted infrastructure file: " + name)

    main_tf = (AZ / "main.tf").read_text(encoding="utf-8")
    vars_tf = (AZ / "variables.tf").read_text(encoding="utf-8")
    versions_tf = (AZ / "versions.tf").read_text(encoding="utf-8")
    readme = (AZ / "README.md").read_text(encoding="utf-8")

    require(versions_tf, r'source\s*=\s*"hashicorp/azurerm"', "AzureRM provider")
    require(main_tf, r'private_cluster_enabled\s*=\s*var\.private_cluster_enabled', "private AKS control-plane binding")
    require(vars_tf, r'condition\s*=\s*var\.private_cluster_enabled', "fail-closed private-cluster validation")
    require(main_tf, r'local_account_disabled\s*=\s*true', "disabled local AKS accounts")
    require(main_tf, r'oidc_issuer_enabled\s*=\s*true', "OIDC issuer")
    require(main_tf, r'workload_identity_enabled\s*=\s*true', "Azure Workload Identity")
    require(main_tf, r'role_based_access_control_enabled\s*=\s*true', "Kubernetes RBAC")
    require(main_tf, r'azure_policy_enabled\s*=\s*true', "Azure Policy")
    require(main_tf, r'network_policy\s*=\s*"azure"', "network policy")
    require(main_tf, r'service_mesh_profile\s*\{[\s\S]*?mode\s*=\s*"Istio"[\s\S]*?revisions\s*=\s*var\.istio_revisions', "managed AKS Istio profile")
    require(vars_tf, r'variable\s+"istio_revisions"', "governed Istio revision input")
    require(vars_tf, r'asm-\[0-9\]\+\-\[0-9\]\+', "Azure managed Istio revision format validation")
    require(main_tf, r'key_vault_secrets_provider\s*\{', "Key Vault CSI integration hook")
    require(main_tf, r'role_definition_name\s*=\s*"Network Contributor"', "BYO-network role assignment")
    require(main_tf, r'scope\s*=\s*azurerm_virtual_network\.this\.id', "network role least-privilege scope")
    require(main_tf, r'host_gate\s*=\s*"HOST-013-HOST-014"', "HOST ownership tag")
    require(readme, r'HOST-013\.\.014 may move to VERIFIED only after a real Azure deployment', "truthful live-verification boundary")

    forbidden = ["client_secret", "password", "api_key", "MARKETDATA_TOKEN", "FINNHUB"]
    joined = "\n".join((AZ / name).read_text(encoding="utf-8") for name in REQUIRED)
    for token in forbidden:
        if token.lower() in joined.lower():
            fail("Azure IaC contains forbidden secret/token material marker: " + token)

    print("PASS: Azure AKS HOST-013..014 adapter is fail-closed, managed-mesh aware and secret-free")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
