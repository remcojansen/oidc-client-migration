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
# - client.response_types --> Unsupported (client type is derived from grant_types instead)
# - client.application_type --> Unsupported (client type is derived from token_endpoint_auth_method instead)
# - client.session_max_lifetime_seconds --> Unsupported (no PingFederate equivalent; PingFederate's
#   persistent grant lifetime applies unconditionally, so it is only represented by
#   client.offline_session_max_lifetime_seconds)
# - client.session_idle_timeout_seconds --> Unsupported (see above)
#
# client.access_token_format / client.access_token_lifetime_seconds are resolved to a
# PingFederate access token manager ID via var.access_token_manager_mapping, since access token
# managers are user-created, deployment-specific resources with no universal naming convention.

locals {
  # Define which scopes are "common" and hence available to all clients by default;
  # these scopes will be excluded from the scopes that are assigned to the client
  common_scopes = ["acr", "basic", "profile", "openid", "email", "address", "phone"]

  client  = var.config.client
  secrets = try(var.config.secrets, {})

  grant_type_mapping = {
    "authorization_code"                              = "AUTHORIZATION_CODE"
    "implicit"                                        = "IMPLICIT"
    "password"                                        = "RESOURCE_OWNER_PASSWORD_CREDENTIALS"
    "client_credentials"                              = "CLIENT_CREDENTIALS"
    "refresh_token"                                   = "REFRESH_TOKEN"
    "urn:ietf:params:oauth:grant-type:device_code"    = "DEVICE_CODE"
    "urn:ietf:params:oauth:grant-type:token-exchange" = "TOKEN_EXCHANGE"
  }

  grant_types = concat(
    [
      for gt in try(local.client.grant_types, []) : lookup(local.grant_type_mapping, gt, local.grant_type_invalid_sentinel)
    ],
    try(local.client.introspection_enabled, true) ? ["ACCESS_TOKEN_VALIDATION"] : []
  )

  # Sentinel value used to detect a grant type with no PingFederate equivalent; checked by the
  # resource's lifecycle.precondition below so the failure surfaces with a clear message instead
  # of silently provisioning a client with a bogus grant type.
  grant_type_invalid_sentinel = "invalid grant type"

  unsupported_grant_types = [
    for gt in try(local.client.grant_types, []) : gt if !contains(keys(local.grant_type_mapping), gt)
  ]

  access_token_lifetime_seconds = try(local.client.access_token_lifetime_seconds, 300)

  refresh_token_lifetime_minutes     = max(1, ceil(try(local.client.offline_session_max_lifetime_seconds, 0) / 60))
  refresh_token_idle_timeout_minutes = max(1, ceil(try(local.client.offline_session_idle_timeout_seconds, 0) / 60))

  atm_id = lookup(
    lookup(var.access_token_manager_mapping, try(local.client.access_token_format, "jwt"), {}),
    local.access_token_lifetime_seconds,
    local.atm_id_invalid_sentinel
  )

  # Sentinel value used to detect an unsupported access_token_format/access_token_lifetime_seconds
  # combination; checked by the resource's lifecycle.precondition below so the failure surfaces
  # with a clear message instead of silently provisioning a client with a bogus manager id.
  atm_id_invalid_sentinel = "invalid access token format or lifetime"

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

  # Normalize nullable booleans for use in conditions/provider boolean fields.
  enabled          = try(local.client.enabled, null) != false
  consent_required = try(local.client.consent_required, null) != false
  pkce_required    = try(local.client.pkce_required, null) != false
  dpop_required    = try(local.client.dpop_required, null) == true
  par_required     = try(local.client.par_required, null) == true
}

resource "pingfederate_oauth_client" "this" {
  enabled     = local.enabled
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

  bypass_approval_page = !local.consent_required
  redirect_uris        = try(local.client.redirect_uris, [])
  grant_types          = local.grant_types

  require_proof_key_for_code_exchange   = local.pkce_required
  require_dpop                          = local.dpop_required
  require_pushed_authorization_requests = local.par_required

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
    (var.minimum_acr_value_param_name) = { values = [try(local.client.minimum_acr_value, "")] }
  }


  # defaults
  allow_authentication_api_init           = false
  enable_cookieless_authentication_api    = false
  lockout_max_malicious_actions_type      = "SERVER_DEFAULT"
  persistent_grant_expiration_time        = local.refresh_token_lifetime_minutes
  persistent_grant_expiration_time_unit   = "MINUTES"
  persistent_grant_expiration_type        = try(local.client.offline_session_max_lifetime_seconds, 0) > 0 ? "OVERRIDE_SERVER_DEFAULT" : "SERVER_DEFAULT"
  persistent_grant_idle_timeout           = local.refresh_token_idle_timeout_minutes
  persistent_grant_idle_timeout_time_unit = "MINUTES"
  persistent_grant_idle_timeout_type      = try(local.client.offline_session_idle_timeout_seconds, 0) > 0 ? "OVERRIDE_SERVER_DEFAULT" : "SERVER_DEFAULT"
  refresh_rolling                         = try(local.client.rotate_refresh_tokens, true) ? "ROLL" : "DONT_ROLL"
  # refresh_token_rolling_grace_period       = 0
  refresh_token_rolling_grace_period_type = "SERVER_DEFAULT"
  # refresh_token_rolling_interval           = 0
  # refresh_token_rolling_interval_time_unit = "MINUTES"
  refresh_token_rolling_interval_type      = "SERVER_DEFAULT" # or "OVERRIDE_SERVER_DEFAULT"
  require_signed_requests                  = false
  restrict_to_default_access_token_manager = false
  token_exchange_processor_policy_ref      = null
  validate_using_all_eligible_atms         = true

  lifecycle {
    precondition {
      condition     = local.atm_id != local.atm_id_invalid_sentinel
      error_message = "Unsupported access_token_format/access_token_lifetime_seconds combination: ${try(local.client.access_token_format, "jwt")}/${local.access_token_lifetime_seconds}. Supported combinations: ${jsonencode(var.access_token_manager_mapping)}."
    }
    precondition {
      condition     = !contains(local.grant_types, local.grant_type_invalid_sentinel)
      error_message = "Unsupported grant_types: ${jsonencode(local.unsupported_grant_types)}. Supported grant types: ${jsonencode(keys(local.grant_type_mapping))}."
    }
  }
}
