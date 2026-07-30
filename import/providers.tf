terraform {
  required_version = ">=1.6"
  required_providers {
    pingfederate = {
      source  = "pingidentity/pingfederate"
      version = ">= 1.8, < 2.0"
    }
    keycloak = {
      source  = "keycloak/keycloak"
      version = ">= 5.8.0"
    }
  }
}

provider "pingfederate" {
  https_host                          = var.provider_pingfederate_url
  product_version                     = var.provider_pingfederate_product_version
  insecure_trust_all_tls              = var.provider_pingfederate_tls_skip_verify
  x_bypass_external_validation_header = var.provider_pingfederate_bypass_external_validation_header
}

provider "keycloak" {
  url                      = var.provider_keycloak_url
  client_id                = var.provider_keycloak_client_id
  tls_insecure_skip_verify = var.provider_keycloak_tls_skip_verify
  username                 = var.provider_keycloak_username
  password                 = var.provider_keycloak_password
}