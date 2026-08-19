variable "config" {
  description = "Canonical client configuration"
  type = object({
    annotations = optional(object({
      owner_email          = optional(string)
      owner_team           = optional(string)
      owner_channel        = optional(string)
      owner_security_group = optional(string)
    }))

    client = object({
      client_id        = string
      client_name      = string
      description      = optional(string)
      enabled          = optional(bool)
      application_type = optional(string)

      token_endpoint_auth_method = optional(string)
      grant_types                = optional(list(string))
      response_types             = optional(list(string))
      redirect_uris              = optional(list(string))
      post_logout_redirect_uris  = optional(list(string))
      scopes                     = optional(list(string))
      consent_required           = optional(bool)

      pkce_required       = optional(bool)
      dpop_required       = optional(bool)
      par_required        = optional(bool)
      access_token_format = optional(string)

      access_token_lifetime_seconds        = optional(number)
      rotate_refresh_tokens                = optional(bool)
      offline_session_max_lifetime_seconds = optional(number)
      offline_session_idle_timeout_seconds = optional(number)
      session_max_lifetime_seconds         = optional(number)
      session_idle_timeout_seconds         = optional(number)
      introspection_enabled                = optional(bool)

      minimum_acr_value     = optional(string)
      default_acr_values    = optional(list(string))
      subject_type          = optional(string)
      sector_identifier_uri = optional(string)

      id_token_signed_response_alg    = optional(string)
      token_endpoint_auth_signing_alg = optional(string)
      request_object_signing_alg      = optional(string)
      jwks_uri                        = optional(string)

      contacts                = optional(list(string))
      client_uri              = optional(string)
      logo_uri                = optional(string)
      tos_uri                 = optional(string)
      policy_uri              = optional(string)
      backchannel_logout_uri  = optional(string)
      frontchannel_logout_uri = optional(string)
      initiate_login_uri      = optional(string)
      request_uris            = optional(list(string))
    })

    secrets = optional(object({
      encrypted_secret = optional(string)
    }))
  })
}

variable "access_token_manager_mapping" {
  description = "Maps client.access_token_format (\"jwt\" or \"opaque\") and client.access_token_lifetime_seconds to this PingFederate deployment's access token manager IDs. PingFederate access token managers are user-created resources with no universal naming convention, so this must be supplied per-deployment. Example: { jwt = { 300 = \"my-jwt-manager-short\", 1800 = \"my-jwt-manager-long\" }, opaque = { 300 = \"my-opaque-manager-short\" } }."
  type        = map(map(string))
  default     = {}
}

variable "minimum_acr_value_param_name" {
  description = "PingFederate Extended Parameter name storing client.minimum_acr_value. PingFederate has no native concept for this; it is an arbitrary, deployment-chosen Extended Parameter."
  type        = string
  default     = "minimum_acr_value"
}