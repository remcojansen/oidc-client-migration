variable "provider_pingfederate_url" {
  description = "The url for the pingfederate provider"
  type        = string
  default     = ""
}

variable "provider_pingfederate_product_version" {
  description = "The product version for the pingfederate provider"
  type        = string
  default     = ""
}

variable "provider_pingfederate_tls_skip_verify" {
  description = "Whether to trust all TLS certificates"
  type        = bool
  default     = true
}

variable "provider_pingfederate_bypass_external_validation_header" {
  description = "Whether to bypass external validation header"
  type        = bool
  default     = true
}

variable "provider_keycloak_url" {
  description = "The url for the keycloak provider"
  type        = string
  default     = "http://localhost:8088"
}

variable "provider_keycloak_client_id" {
  description = "The client id for the keycloak provider"
  type        = string
  default     = "admin-cli"
}

variable "provider_keycloak_tls_skip_verify" {
  description = "Whether to skip TLS verification for the Keycloak provider"
  type        = bool
  default     = true
}

variable "provider_keycloak_username" {
  description = "The username for the keycloak provider"
  type        = string
  default     = "admin"
}

variable "provider_keycloak_password" {
  description = "The password for the keycloak provider"
  type        = string
  default     = "admin"
}

variable "config_base" {
  description = "The base path to the client config files"
  type        = string
}

variable "keycloak_realm_id" {
  description = "The ID of the realm where the clients will be created"
  type        = string
}

variable "pingfederate_enabled" {
  description = "Whether PingFederate client provisioning is enabled"
  type        = bool
  default     = true
}

variable "keycloak_enabled" {
  description = "Whether Keycloak client provisioning is enabled"
  type        = bool
  default     = true
}
