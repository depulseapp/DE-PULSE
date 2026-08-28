# DE.PULSE Azure AKS hosted trust baseline

This directory is the Azure-specific deployment adapter for the canonical `internal/hostedenv` desired state. It does not replace the portable Kubernetes/Istio owner under `tools/hosted`; it provides the managed Azure substrate needed to prove HOST-013..014 in a real non-production environment.

## Design

- Azure Resource Group + VNet + dedicated AKS subnet.
- AKS with OIDC issuer and Azure Workload Identity enabled.
- Azure RBAC for Kubernetes authorization.
- Private-cluster-ready networking posture; API server public exposure is disabled by default in variables.
- Network policy enabled.
- Azure Key Vault Secrets Provider enabled as the managed-secret integration hook for HOST-017..020; secret material is never stored in Terraform state by this module.
- Separate environment naming and tags for `dev`, `test`, `stage`, and `prod`.
- No public/commercial readiness claim is created by applying this module.

## Verification boundary

Repository CI may validate syntax, structure, and fail-closed policy. HOST-013..014 may move to VERIFIED only after a real Azure deployment proves environment isolation, workload identity, TLS, strict internal mTLS, ingress authorization, denied unregistered egress, and drift detection.

## Typical non-production flow

1. Authenticate to Azure using a dedicated deployment identity or GitHub OIDC federation.
2. Supply `subscription_id`, `tenant_id`, `environment`, `location`, and approved address ranges.
3. `terraform init`
4. `terraform plan`
5. `terraform apply`
6. Render/apply the portable Kubernetes/Istio trust bundle from `tools/hosted/render_kubernetes_trust.py`.
7. Run live adverse tests and retain secret-free evidence.

The module intentionally does not create application secrets or provider API tokens.