locals {
  client_configs = {
    for file in fileset("${var.config_base}/", "*.yaml") :
    trimsuffix(file, ".yaml") => yamldecode(file("${var.config_base}/${file}"))
  }
}

module "pingfederate-client" {
  source   = "./modules/oauth-client-pingfederate"
  for_each = var.pingfederate_enabled ? local.client_configs : {}

  config = each.value
}

module "keycloak-client" {
  source   = "./modules/oauth-client-keycloak"
  for_each = var.keycloak_enabled ? local.client_configs : {}

  realm_id = var.keycloak_realm_id
  config   = each.value
}