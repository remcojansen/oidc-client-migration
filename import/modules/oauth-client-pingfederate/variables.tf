variable "config" {
  description = "Canonical client configuration"
  type = object({
    metadata = optional(object({
      owner_email         = optional(string)
      owner_slack_channel = optional(string)
      owner_team_ad_group = optional(string)
    }))

    client = object({
      client_id                       = string
      client_name                     = string
      description                     = optional(string)
      contacts                        = optional(list(string))
      client_uri                      = optional(string)
      logo_uri                        = optional(string)
      tos_uri                         = optional(string)
      policy_uri                      = optional(string)
      jwks_uri                        = optional(string)
      redirect_uris                   = optional(list(string))
      response_types                  = optional(list(string))
      grant_types                     = optional(list(string))
      token_endpoint_auth_method      = optional(string)
      token_endpoint_auth_signing_alg = optional(string)
      id_token_signed_response_alg    = optional(string)
      request_object_signing_alg      = optional(string)
      scopes                          = optional(list(string))
      consent_required                = optional(bool)
      application_type                = optional(string)
      subject_type                    = optional(string)
      sector_identifier_uri           = optional(string)
      backchannel_logout_uri          = optional(string)
      frontchannel_logout_uri         = optional(string)
      post_logout_redirect_uris       = optional(list(string))
      default_acr_values              = optional(list(string))
      initiate_login_uri              = optional(string)
      request_uris                    = optional(list(string))
    })

    extensions = optional(object({
      enabled                              = optional(bool)
      pkce_required                        = optional(bool)
      dpop_required                        = optional(bool)
      par_required                         = optional(bool)
      access_token_format                  = optional(string)
      access_token_lifetime_seconds        = optional(number)
      offline_session_max_lifetime_seconds = optional(number)
      offline_session_idle_timeout_seconds = optional(number)
      session_max_lifetime_seconds         = optional(number)
      session_idle_timeout_seconds         = optional(number)
      rotate_refresh_tokens                = optional(bool)
      minimum_acr_value                    = optional(string)
      terms_and_conditions_required        = optional(bool)
    }))

    secrets = optional(object({
      encrypted_secret = optional(string)
    }))
  })
}