locals {
  prefix                    = "depulse-${var.environment}"
  workload_namespace        = "depulse-${var.environment}"
  workload_service_account  = "depulse-web-${var.environment}"
  workload_identity_subject = "system:serviceaccount:${local.workload_namespace}:${local.workload_service_account}"
  common_tags = merge({
    application = "DE.PULSE"
    environment = var.environment
    lifecycle   = "DEVELOPMENT"
    owner       = "hostedenv"
    managed_by  = "terraform"
    host_gate   = "HOST-013-HOST-014"
  }, var.tags)
}

resource "azurerm_resource_group" "this" {
  name     = "rg-${local.prefix}"
  location = var.location
  tags     = local.common_tags
}

resource "azurerm_virtual_network" "this" {
  name                = "vnet-${local.prefix}"
  location            = azurerm_resource_group.this.location
  resource_group_name = azurerm_resource_group.this.name
  address_space       = [var.vnet_cidr]
  tags                = local.common_tags
}

resource "azurerm_subnet" "aks" {
  name                 = "snet-aks"
  resource_group_name  = azurerm_resource_group.this.name
  virtual_network_name = azurerm_virtual_network.this.name
  address_prefixes     = [var.aks_subnet_cidr]
}

resource "azurerm_user_assigned_identity" "aks" {
  name                = "id-${local.prefix}-aks"
  location            = azurerm_resource_group.this.location
  resource_group_name = azurerm_resource_group.this.name
  tags                = local.common_tags
}

resource "azurerm_user_assigned_identity" "workload" {
  name                = "id-${local.prefix}-workload"
  location            = azurerm_resource_group.this.location
  resource_group_name = azurerm_resource_group.this.name
  tags                = merge(local.common_tags, { identity_role = "depulse-workload" })
}

resource "azurerm_role_assignment" "aks_network" {
  scope                = azurerm_virtual_network.this.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.aks.principal_id
}

resource "azurerm_kubernetes_cluster" "this" {
  name                = "aks-${local.prefix}"
  location            = azurerm_resource_group.this.location
  resource_group_name = azurerm_resource_group.this.name
  dns_prefix          = local.prefix
  kubernetes_version  = var.kubernetes_version

  private_cluster_enabled           = var.private_cluster_enabled
  local_account_disabled            = true
  oidc_issuer_enabled               = true
  workload_identity_enabled         = true
  role_based_access_control_enabled = true
  azure_policy_enabled              = true

  azure_active_directory_role_based_access_control {
    tenant_id          = var.tenant_id
    azure_rbac_enabled = true
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.aks.id]
  }

  default_node_pool {
    name                         = "system"
    vm_size                      = var.node_vm_size
    node_count                   = var.node_count
    vnet_subnet_id               = azurerm_subnet.aks.id
    os_disk_type                 = "Managed"
    type                         = "VirtualMachineScaleSets"
    only_critical_addons_enabled = false
    upgrade_settings {
      max_surge = "33%"
    }
  }

  network_profile {
    network_plugin = "azure"
    network_policy = "azure"
    outbound_type  = "loadBalancer"
  }

  service_mesh_profile {
    mode                             = "Istio"
    revisions                        = var.istio_revisions
    external_ingress_gateway_enabled = true
    internal_ingress_gateway_enabled = false
  }

  key_vault_secrets_provider {
    secret_rotation_enabled  = true
    secret_rotation_interval = "2m"
  }

  lifecycle {
    precondition {
      condition     = var.private_cluster_enabled
      error_message = "HOST-013..014 requires a private AKS control plane"
    }
  }

  depends_on = [azurerm_role_assignment.aks_network]
  tags       = local.common_tags
}

resource "azurerm_federated_identity_credential" "workload" {
  name      = "fic-${local.prefix}-workload"
  parent_id = azurerm_user_assigned_identity.workload.id
  audience  = ["api://AzureADTokenExchange"]
  issuer    = azurerm_kubernetes_cluster.this.oidc_issuer_url
  subject   = local.workload_identity_subject
}
