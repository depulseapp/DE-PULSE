variable "subscription_id" {
  description = "Azure subscription hosting the non-production DE.PULSE environment."
  type        = string
  sensitive   = true
}

variable "tenant_id" {
  description = "Microsoft Entra tenant id."
  type        = string
  sensitive   = true
}

variable "environment" {
  description = "Canonical hosted environment."
  type        = string
  validation {
    condition     = contains(["dev", "test", "stage", "prod"], var.environment)
    error_message = "environment must be one of dev, test, stage, prod"
  }
}

variable "location" {
  description = "Azure region."
  type        = string
  default     = "canadacentral"
}

variable "vnet_cidr" {
  description = "Environment-isolated VNet CIDR."
  type        = string
  default     = "10.40.0.0/16"
}

variable "aks_subnet_cidr" {
  description = "AKS subnet CIDR."
  type        = string
  default     = "10.40.0.0/20"
}

variable "kubernetes_version" {
  description = "AKS Kubernetes version. Null lets Azure choose the current supported default."
  type        = string
  default     = null
}

variable "node_vm_size" {
  description = "Non-production node size."
  type        = string
  default     = "Standard_D2as_v5"
}

variable "node_count" {
  description = "Initial system pool size."
  type        = number
  default     = 1
  validation {
    condition     = var.node_count >= 1 && var.node_count <= 3
    error_message = "non-production node_count must remain between 1 and 3"
  }
}

variable "private_cluster_enabled" {
  description = "Keep AKS API server private. Must remain true for governed deployment."
  type        = bool
  default     = true
  validation {
    condition     = var.private_cluster_enabled
    error_message = "DE.PULSE hosted trust baseline requires private_cluster_enabled=true"
  }
}

variable "tags" {
  description = "Additional resource tags."
  type        = map(string)
  default     = {}
}
