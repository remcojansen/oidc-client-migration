provider_pingfederate_url                               = "https://localhost:9999"
provider_pingfederate_product_version                   = "12.2.2"
provider_pingfederate_tls_skip_verify                   = true
provider_pingfederate_bypass_external_validation_header = true

provider_keycloak_url             = "http://localhost:8088"
provider_keycloak_client_id       = "admin-cli"
provider_keycloak_tls_skip_verify = true

config_base = "../client-configurations"

keycloak_realm_id = "example"

pingfederate_enabled = false
keycloak_enabled     = false