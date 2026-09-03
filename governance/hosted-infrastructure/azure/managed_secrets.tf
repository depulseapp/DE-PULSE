locals {
  managed_secret_vault_name = "kvdp${var.environment}${substr(replace(data.azurerm_client_config.current.subscription_id, "-", ""), 0, 8)}"
}

# HOST-017/018: provider credentials remain outside ordinary DE.PULSE product
# state and outside Terraform state. Operators provision secret values directly
# into this environment-scoped vault; Terraform owns only custody boundaries.
resource "azurerm_key_vault" "managed_secrets" {
  name                       = local.managed_secret_vault_name
  location                   = azurerm_resource_group.this.location
  resource_group_name        = azurerm_resource_group.this.name
  tenant_id                  = var.tenant_id
  sku_name                   = "standard"
  enable_rbac_authorization  = true
  purge_protection_enabled   = true
  soft_delete_retention_days = 7
  public_network_access_enabled = false

  tags = merge(local.common_tags, {
    host_gate = "HOST-017-HOST-018"
    custody   = "managed-secrets"
  })
}

resource "azurerm_role_assignment" "workload_managed_secret_reader" {
  scope                = azurerm_key_vault.managed_secrets.id
  role_definition_name = "Key Vault Secrets User"
  principal_id         = azurerm_user_assigned_identity.workload.principal_id
}

resource "azurerm_private_dns_zone" "key_vault" {
  name                = "privatelink.vaultcore.azure.net"
  resource_group_name = azurerm_resource_group.this.name
  tags                = local.common_tags
}

resource "azurerm_private_dns_zone_virtual_network_link" "key_vault" {
  name                  = "link-${local.prefix}-key-vault"
  resource_group_name   = azurerm_resource_group.this.name
  private_dns_zone_name = azurerm_private_dns_zone.key_vault.name
  virtual_network_id    = azurerm_virtual_network.this.id
  registration_enabled  = false
  tags                  = local.common_tags
}

resource "azurerm_private_endpoint" "key_vault" {
  name                = "pe-${local.prefix}-key-vault"
  location            = azurerm_resource_group.this.location
  resource_group_name = azurerm_resource_group.this.name
  subnet_id           = azurerm_subnet.aks.id
  tags                = local.common_tags

  private_service_connection {
    name                           = "psc-${local.prefix}-key-vault"
    private_connection_resource_id = azurerm_key_vault.managed_secrets.id
    subresource_names              = ["vault"]
    is_manual_connection           = false
  }

  private_dns_zone_group {
    name                 = "key-vault"
    private_dns_zone_ids = [azurerm_private_dns_zone.key_vault.id]
  }
}

output "managed_secret_vault_name" {
  value       = azurerm_key_vault.managed_secrets.name
  description = "Environment-scoped Key Vault name used by the HOST-017/018 CSI secret provider contract."
}

output "managed_secret_workload_client_id" {
  value       = azurerm_user_assigned_identity.workload.client_id
  description = "Workload identity client ID authorized for read-only managed-secret resolution."
}
