locals {
  prefix = "depulse-${var.environment}"
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

resource "azurerm_kubernetes_cluster" "this" {
  name                = "aks-${local.prefix}"
  location            = azurerm_resource_group.this.location
  resource_group_name = azurerm_resource_group.this.name
  dns_prefix          = local.prefix
  kubernetes_version  = var.kubernetes_version

  private_cluster_enabled = var.private_cluster_enabled
  local_account_disabled  = true
  oidc_issuer_enabled     = true
  workload_identity_enabled = true
  role_based_access_control_enabled = true
  azure_policy_enabled            = true

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.aks.id]
  }

  default_node_pool {
    name           = "system"
    vm_size        = var.node_vm_size
    node_count     = var.node_count
    vnet_subnet_id = azurerm_subnet.aks.id
    os_disk_type   = "Managed"
    type           = "VirtualMachineScaleSets"
    only_critical_addons_enabled = true
    upgrade_settings {
      max_surge = "33%"
    }
  }

  network_profile {
    network_plugin = "azure"
    network_policy = "azure"
    outbound_type  = "loadBalancer"
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

  tags = local.common_tags
}
