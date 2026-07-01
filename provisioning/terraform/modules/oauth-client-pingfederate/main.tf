# Maps the canonical client configuration to PingFederate provider resources for OAuth clients.
#
# Note: This module is intended to be used by the provisioning/terraform/main.tf file, which 
# reads canonical client configuration from YAML files and passes them to this module.
#
# The following fields lack a direct PingFederate equivalent and are therefore not currently
# mapped to PingFederate provider resources:
# - client.contacts --> Unsupported
# - client.default_acr_values --> Unsupported
# - client.request_uris --> Unsupported
# - client.initiate_login_uri --> Unsupported
# - client.sector_identifier_uri --> Unsupported
# - client.token_endpoint_auth_signing_alg --> Unsupported
# - client.request_object_signing_alg --> Unsupported
# - client.client_uri --> Unsupported
# - client.tos_uri --> Unsupported
# - client.policy_uri --> Unsupported

locals {
  # Define which scopes are "common" and hence available to all clients by default;
  # these scopes will be excluded from the scopes that are assigned to the client
  common_scopes = ["acr", "basic", "profile", "openid", "email", "address", "phone"]

  client     = var.config.client
  extensions = try(var.config.extensions, {})
  secrets    = try(var.config.secrets, {})

  grant_type_mapping = {
    "authorization_code"                              = "AUTHORIZATION_CODE"
    "implicit"                                        = "IMPLICIT"
    "password"                                        = "RESOURCE_OWNER_PASSWORD_CREDENTIALS"
    "client_credentials"                              = "CLIENT_CREDENTIALS"
    "refresh_token"                                   = "REFRESH_TOKEN"
    "device_code"                                     = "DEVICE_CODE"
    "introspect"                                      = "ACCESS_TOKEN_VALIDATION"
    "urn:ietf:params:oauth:grant-type:token-exchange" = "TOKEN_EXCHANGE"
  }

  grant_types = [
    for gt in try(local.client.grant_types, []) : lookup(local.grant_type_mapping, gt, "invalid grant type")
  ]

  access_token_manager_mapping = {
    "opaque" = {
      300 = "refstandardshort"
      # TODO: verify opaque manager id for non-default access token lifetimes.
    }
    "jwt" = {
      300  = "jwtstandardshort"
      1800 = "jwtstandardlong"
      # TODO: verify JWT manager ids for non-default access token lifetimes.
    }
  }

  access_token_lifetime_seconds = try(local.extensions.access_token_lifetime_seconds, 300)

  refresh_token_lifetime_minutes     = max(1, ceil(try(local.extensions.refresh_token_lifetime_seconds, 0) / 60))
  refresh_token_idle_timeout_minutes = max(1, ceil(try(local.extensions.refresh_token_idle_timeout_seconds, 0) / 60))

  atm_id = lookup(
    lookup(local.access_token_manager_mapping, try(local.extensions.access_token_format, "jwt"), {}),
    local.access_token_lifetime_seconds,
    "invalid access token format or lifetime"
  )

  logout_mode = (
    try(local.client.frontchannel_logout_uri, "") != "" ? "OIDC_FRONT_CHANNEL" :
    try(local.client.backchannel_logout_uri, "") != "" ? "OIDC_BACK_CHANNEL" :
    "NONE"
  )

  client_auth_type = (
    try(local.client.token_endpoint_auth_method, "") == "none" ? "NONE" :
    try(local.client.token_endpoint_auth_method, "") == "client_secret_jwt" ? "CLIENT_SECRET_JWT" :
    try(local.client.token_endpoint_auth_method, "") == "private_key_jwt" ? "PRIVATE_KEY_JWT" :
    "SECRET"
  )

  scopes = try(local.client.scopes, [])
}

resource "pingfederate_oauth_client" "this" {
  enabled     = try(local.extensions.enabled, true)
  client_id   = local.client.client_id
  name        = local.client.client_name
  description = try(local.client.description, null)
  logo_url    = try(local.client.logo_uri, null)

  client_auth = {
    type   = local.client_auth_type
    secret = try(local.secrets.encrypted_secret, null)
  }

  jwks_settings = {
    jwks_url = try(local.client.jwks_uri, null)
  }

  bypass_approval_page = !try(local.client.consent_required, true)
  redirect_uris        = try(local.client.redirect_uris, [])
  grant_types          = local.grant_types

  require_proof_key_for_code_exchange   = try(local.extensions.pkce_required, false)
  require_dpop                          = try(local.extensions.dpop_required, false)
  require_pushed_authorization_requests = try(local.extensions.par_required, false)

  default_access_token_manager_ref = {
    id = local.atm_id
  }

  oidc_policy = {
    logout_mode                   = local.logout_mode
    logout_uris                   = try(local.client.frontchannel_logout_uri, "") != "" ? [local.client.frontchannel_logout_uri] : []
    back_channel_logout_uri       = try(local.client.backchannel_logout_uri, "")
    post_logout_redirect_uris     = try(local.client.post_logout_redirect_uris, [])
    pairwise_identifier_user_type = try(local.client.subject_type, "") == "pairwise"
    id_token_signing_algorithm    = try(local.client.id_token_signed_response_alg, "")

    policy_group = {
      id = local.atm_id
    }
  }

  restrict_scopes  = false
  exclusive_scopes = local.scopes

  extended_parameters = {
    exclude_tnc = { values = [try(local.extensions.require_terms_and_conditions_approval, false) ? "false" : "true"] }
    enforce_2sv = { values = [try(local.extensions.minimum_acr_value, "")] }
  }

  # defaults
  allow_authentication_api_init           = false
  enable_cookieless_authentication_api    = false
  lockout_max_malicious_actions_type      = "SERVER_DEFAULT"
  persistent_grant_expiration_time        = local.refresh_token_lifetime_minutes
  persistent_grant_expiration_time_unit   = "MINUTES"
  persistent_grant_expiration_type        = try(local.extensions.refresh_token_lifetime_seconds, 0) > 0 ? "OVERRIDE_SERVER_DEFAULT" : "SERVER_DEFAULT"
  persistent_grant_idle_timeout           = local.refresh_token_idle_timeout_minutes
  persistent_grant_idle_timeout_time_unit = "MINUTES"
  persistent_grant_idle_timeout_type      = try(local.extensions.refresh_token_idle_timeout_seconds, 0) > 0 ? "OVERRIDE_SERVER_DEFAULT" : "SERVER_DEFAULT"
  refresh_rolling                         = try(local.extensions.rotate_refresh_tokens, true) ? "ROLL" : "DONT_ROLL"
  # refresh_token_rolling_grace_period       = 0
  refresh_token_rolling_grace_period_type = "SERVER_DEFAULT"
  # refresh_token_rolling_interval           = 0
  # refresh_token_rolling_interval_time_unit = "MINUTES"
  refresh_token_rolling_interval_type      = "SERVER_DEFAULT" # or "OVERRIDE_SERVER_DEFAULT"
  require_signed_requests                  = false
  restrict_to_default_access_token_manager = false
  token_exchange_processor_policy_ref      = null
  validate_using_all_eligible_atms         = true
}