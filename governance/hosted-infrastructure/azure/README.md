# DE.PULSE Azure AKS hosted trust baseline

This directory is the Azure-specific deployment adapter for the canonical `internal/hostedenv` desired state. It does not replace the portable Kubernetes/Istio owner under `tools/hosted`; it provides the managed Azure substrate needed to prove HOST-013..014 in a real non-production environment.

## Design

- Azure Resource Group + VNet + dedicated AKS subnet.
- Private AKS control plane with local accounts disabled.
- User-assigned AKS control-plane identity with only the required Network Contributor authority on the environment VNet.
- OIDC issuer and Azure Workload Identity enabled for workload identity without long-lived client secrets.
- Kubernetes RBAC and Azure Policy enabled.
- Azure network policy enabled.
- Microsoft-managed Istio service-mesh add-on; supported `asm-X-Y` revisions are discovered from Azure before provisioning and may use a two-revision canary window during governed upgrades.
- Azure Key Vault Secrets Provider enabled as the managed-secret integration hook for HOST-017..020; application/provider secret material is never declared by this module.
- Separate environment naming and tags for `dev`, `test`, `stage`, and `prod`.
- No public/commercial readiness claim is created by applying this module.

## Terraform state and authentication

Terraform state uses the AzureRM backend with `use_oidc=true` and `use_azuread_auth=true`. GitHub Actions must authenticate by Microsoft Entra workload-identity federation/OIDC; no Azure client secret, storage-account key or SAS token is an approved deployment path.

The remote-state storage account/container is a bootstrap prerequisite and is intentionally separate from the application environment. `terraform init` supplies the non-secret backend coordinates (`resource_group_name`, `storage_account_name`, `container_name`, and an environment-specific state key). The federated deployment principal should receive `Storage Blob Data Contributor` at the container scope plus only the management-plane permissions required by the chosen backend endpoint mode.

GitHub/Azure identifiers required by the future operator run are non-secret identifiers:

- `AZURE_CLIENT_ID`
- `AZURE_TENANT_ID`
- `AZURE_SUBSCRIPTION_ID`
- Terraform-state resource group/storage account/container names

Do not store a client secret, subscription credential, kubeconfig, storage key, provider API token, or application secret in repository files or retained CI evidence.

## Verification boundary

Repository CI validates structure, policy, the OIDC state contract and the live-evidence collector's positive/adverse self-tests. HOST-013..014 may move to VERIFIED only after a real Azure deployment proves environment isolation, workload identity, TLS, strict internal mTLS, ingress authorization, denied unregistered egress, private-cluster operability and drift detection.

For the private cluster, live inspection can use `az aks command invoke` through the Azure API so the Kubernetes API server does not need public exposure. `tools/ci/host013_azure_live_evidence.py` retains only sanitized evidence and explicitly rejects a missing/failed run-command status.

## Typical non-production flow

1. Establish the GitHub-to-Azure federated identity and least-privilege role assignments.
2. Create/identify the governed AzureRM remote-state container and grant the federated identity least-privilege state access.
3. Discover a Canada Central AKS Kubernetes version and compatible managed Istio revision with Azure's supported-version/revision APIs.
4. Run `terraform init` with the remote-state backend coordinates.
5. Run `terraform plan` and review the exact non-production resource changes.
6. Run `terraform apply` only through the explicitly confirmed non-production operator path.
7. Render/apply the portable Kubernetes/Istio trust bundle from `tools/hosted/render_kubernetes_trust.py`.
8. Use private-cluster `az aks command invoke` checks plus adverse traffic/network tests.
9. Run `tools/ci/host013_azure_live_evidence.py` and retain its secret-free exact-head evidence artifact.
10. Re-run `terraform plan -detailed-exitcode` after verification; unexpected drift prevents HOST-013..014 closure.

The Azure module intentionally does not create application secrets or provider API tokens. Later HOST-017..020 work will bind managed secrets through the existing secret owners rather than embedding secret values in Terraform state.
