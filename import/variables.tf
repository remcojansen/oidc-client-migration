variable "provider_pingfederate_url" {
  description = "The url (https_host) for the pingfederate provider. If left unset (null), the underlying PingFederate Terraform provider falls back to the PINGFEDERATE_PROVIDER_HTTPS_HOST environment variable."
  type        = string
  default     = null
}

variable "provider_pingfederate_product_version" {
  description = "The product version for the pingfederate provider. If left unset (null), the underlying PingFederate Terraform provider falls back to the PINGFEDERATE_PROVIDER_PRODUCT_VERSION environment variable."
  type        = string
  default     = null
}

variable "provider_pingfederate_tls_skip_verify" {
  description = "Whether to trust all TLS certificates. If left unset (null), the underlying PingFederate Terraform provider falls back to the PINGFEDERATE_PROVIDER_INSECURE_TRUST_ALL_TLS environment variable, or its own default if that is unset too."
  type        = bool
  default     = null
}

variable "provider_pingfederate_bypass_external_validation_header" {
  description = "Whether to bypass external validation header. If left unset (null), the underlying PingFederate Terraform provider falls back to the PINGFEDERATE_PROVIDER_X_BYPASS_EXTERNAL_VALIDATION_HEADER environment variable, or its own default if that is unset too."
  type        = bool
  default     = null
}

variable "provider_keycloak_url" {
  description = "The url for the keycloak provider. If left unset (null), the underlying Keycloak Terraform provider falls back to the KEYCLOAK_URL environment variable."
  type        = string
  default     = null
}

variable "provider_keycloak_client_id" {
  description = "The client id for the keycloak provider. If left unset (null), the underlying Keycloak Terraform provider falls back to the KEYCLOAK_CLIENT_ID environment variable."
  type        = string
  default     = null
}

variable "provider_keycloak_tls_skip_verify" {
  description = "Whether to skip TLS verification for the Keycloak provider. The Keycloak Terraform provider does not support configuring this via an environment variable, so this must be set via this Terraform variable."
  type        = bool
  default     = false
}

variable "provider_keycloak_username" {
  description = "The username for the keycloak provider. If left unset (null), the underlying Keycloak Terraform provider falls back to the KEYCLOAK_USER environment variable."
  type        = string
  default     = null
}

variable "provider_keycloak_password" {
  description = "The password for the keycloak provider. If left unset (null), the underlying Keycloak Terraform provider falls back to the KEYCLOAK_PASSWORD environment variable."
  type        = string
  default     = null
  sensitive   = true
}

variable "config_base" {
  description = "The base path to the client config files"
  type        = string
  default     = "../client-configurations"
}

variable "keycloak_realm_id" {
  description = "The ID of the realm where the clients will be created"
  type        = string
}

variable "pingfederate_enabled" {
  description = "Whether PingFederate client provisioning is enabled"
  type        = bool
  default     = false
}

variable "keycloak_enabled" {
  description = "Whether Keycloak client provisioning is enabled"
  type        = bool
  default     = false
}
