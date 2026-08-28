output "resource_group_name" {
  value = azurerm_resource_group.this.name
}

output "aks_cluster_name" {
  value = azurerm_kubernetes_cluster.this.name
}

output "oidc_issuer_url" {
  value = azurerm_kubernetes_cluster.this.oidc_issuer_url
}

output "kubelet_identity_object_id" {
  value = azurerm_kubernetes_cluster.this.kubelet_identity[0].object_id
}

output "deployment_identity_client_id" {
  value = azurerm_user_assigned_identity.aks.client_id
}

output "workload_identity_client_id" {
  value       = azurerm_user_assigned_identity.workload.client_id
  description = "Non-secret client ID bound to the canonical DE.PULSE Kubernetes service account."
}

output "workload_identity_principal_id" {
  value       = azurerm_user_assigned_identity.workload.principal_id
  description = "Object/principal ID of the canonical DE.PULSE workload managed identity."
}

output "workload_identity_subject" {
  value       = local.workload_identity_subject
  description = "Exact Kubernetes service-account subject authorized by the federated identity credential."
}
