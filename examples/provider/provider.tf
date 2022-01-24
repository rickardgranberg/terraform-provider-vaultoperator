terraform {
  required_providers {
    vaultoperator = {
      version = "0.2.0"
      source  = "rickardgranberg/vaultoperator"
    }
  }
}

provider "vaultoperator" {
  # example configuration here
  vault_url = "http://vault:8200"
  kube_config {
    path       = "~/.kube/config"
    namespace  = "vault"
    service    = "vault"
    localPort  = "8200"
    remotePort = "8200"
  }
}
