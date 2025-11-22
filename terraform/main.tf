terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~>4.3.0"
    }
  }
  backend "azurerm" {
    resource_group_name   = "KingRG"
    storage_account_name  = "kingst"
    container_name        = "tfstate"
    key                   = "devops-portfolio.tfstate"
  }
}


provider "azurerm" {
  features {}
  subscription_id = var.subscription_id
}

# Create Resource Group
resource "azurerm_resource_group" "KingRG" {
  name     = "KingRG"
  location = "South Africa North"
}

resource "azurerm_kubernetes_cluster" "aks" {
  name                = "DevOpsCluster"
  location            = "South Africa North"
  resource_group_name = "KingRG"
  dns_prefix          = "devops"
  default_node_pool {
    name       = "default"
    node_count = 2
    vm_size    = "Standard_DS2_v2"
  }
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_container_registry" "acr" {
  name                = "kingacr" 
  resource_group_name = azurerm_resource_group.KingRG.name
  location            = azurerm_resource_group.KingRG.location
  sku                 = "Basic"
  admin_enabled       = false # Best practice is to disable admin access for security
}

