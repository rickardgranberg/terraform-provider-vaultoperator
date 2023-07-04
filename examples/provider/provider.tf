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
    # optional exec:
    exec {
      api_version = "client.authentication.k8s.io/v1beta1"
      args        = ["eks", "get-token", "--cluster-name", var.cluster_name]
      command     = "aws"
    }
  }
}
