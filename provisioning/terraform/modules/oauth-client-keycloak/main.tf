# Maps the canonical client configuration to Keycloak provider resources for OAuth clients.
#
# Note: This module is intended to be used by the provisioning/terraform/main.tf file, which 
# reads canonical client configuration from YAML files and passes them to this module.
#
# The following fields lack a direct Keycloak equivalent and are therefore not currently
# mapped to Keycloak provider resources:
# - client.contacts
# - client.response_types
# - client.initiate_login_uri
# - client.sector_identifier_uri
# - client.subject_type
# - extensions.refresh_token_lifetime_seconds
# - extensions.refresh_token_idle_timeout_seconds
# - extensions.access_token_format

locals {
  # Define which scopes should be available to all clients by default;
  # these scopes will be added to the scopes that are assigned to the client
  default_scopes = ["acr", "basic", "profile", "openid", "email", "address", "phone"]

  # Define the login theme to use for all clients
  login_theme = "keycloak"

  client     = var.config.client
  extensions = try(var.config.extensions, {})
  secrets    = try(var.config.secrets, {})

  # Canonical application_type convention: web -> public, native -> confidential
  access_type = local.client.application_type == "web" ? "PUBLIC" : "CONFIDENTIAL"

  client_authenticator_type = (
    try(local.client.token_endpoint_auth_method, "") == "client_secret_jwt" ? "client-secret-jwt" :
    try(local.client.token_endpoint_auth_method, "") == "private_key_jwt" ? "client-jwt" :
    try(local.client.token_endpoint_auth_method, "") == "none" ? null :
    local.access_type == "PUBLIC" ? null : "client-secret"
  )

  # Keycloak splits scopes into default and optional scopes, while the canonical model has one 
  # scopes list. This mapping writes canonical scopes as optional scopes.
  optional_scopes = try(local.client.scopes, [])

  extra_config = merge(
    {
      "par.required"                    = tostring(try(local.extensions.par_required, false))
      "dpop.required"                   = tostring(try(local.extensions.dpop_required, false))
      "require-tnc"                     = tostring(try(local.extensions.require_terms_and_conditions_approval, false))
      "post.logout.redirect.uris"       = join(",", try(local.client.post_logout_redirect_uris, []))
      "backchannel.logout.url"          = try(local.client.backchannel_logout_uri, "")
      "frontchannel.logout.url"         = try(local.client.frontchannel_logout_uri, "")
      "id.token.signed.response.alg"    = try(local.client.id_token_signed_response_alg, "")
      "token.endpoint.auth.signing.alg" = try(local.client.token_endpoint_auth_signing_alg, "")
      "request.object.signature.alg"    = try(local.client.request_object_signing_alg, "")
      "minimum.acr.value"               = try(local.client.minimum_acr_value, "")
      "default.acr.values"              = join(",", try(local.client.default_acr_values, []))
      "logoUri"                         = try(local.client.logo_uri, "")
      "tosUri"                          = try(local.client.tos_uri, "")
      "policyUri"                       = try(local.client.policy_uri, "")
      "jwks.url"                        = try(local.client.jwks_uri, "")
      "use.jwks.url"                    = try(local.client.jwks_uri, "") != ""
      "request.uris"                    = join(",", try(local.client.request_uris, []))
    }
  )
}

resource "keycloak_openid_client" "this" {
  enabled     = try(local.extensions.enabled, true)
  realm_id    = var.realm_id
  client_id   = local.client.client_id
  name        = local.client.client_name
  description = try(local.client.description, null)

  client_secret             = try(local.secrets.plain_secret, null)
  access_type               = local.access_type
  client_authenticator_type = local.client_authenticator_type

  consent_required    = try(local.client.consent_required, false)
  valid_redirect_uris = try(local.client.redirect_uris, [])

  standard_flow_enabled                     = contains(try(local.client.grant_types, []), "authorization_code")
  implicit_flow_enabled                     = contains(try(local.client.grant_types, []), "implicit")
  direct_access_grants_enabled              = contains(try(local.client.grant_types, []), "password")
  service_accounts_enabled                  = contains(try(local.client.grant_types, []), "client_credentials")
  oauth2_device_authorization_grant_enabled = contains(try(local.client.grant_types, []), "device_code")
  standard_token_exchange_enabled           = contains(try(local.client.grant_types, []), "urn:ietf:params:oauth:grant-type:token-exchange")
  use_refresh_tokens                        = contains(try(local.client.grant_types, []), "refresh_token")

  pkce_code_challenge_method = try(local.extensions.pkce_required, false) ? "S256" : null

  extra_config = local.extra_config

  access_token_lifespan = try(local.extensions.access_token_lifetime_seconds, null)

  frontchannel_logout_enabled     = try(local.client.frontchannel_logout_uri, "") == true
  frontchannel_logout_url         = try(local.client.frontchannel_logout_uri, null)
  valid_post_logout_redirect_uris = try(local.client.post_logout_redirect_uris, [])

  # defaults
  admin_url                           = null
  base_url                            = null
  always_display_in_console           = true
  backchannel_logout_session_required = false
  consent_screen_text                 = ""
  login_theme                         = local.login_theme
  root_url                            = try(local.client.client_uri, null)
  web_origins                         = []
}

resource "keycloak_openid_client_optional_scopes" "this" {
  realm_id  = var.realm_id
  client_id = keycloak_openid_client.this.id

  optional_scopes = merge(local.default_scopes, local.optional_scopes)
}

# Create an audience mapping to set the client ID as audience
resource "keycloak_openid_audience_protocol_mapper" "audience_mapper" {
  realm_id                 = var.realm_id
  client_id                = keycloak_openid_client.this.id
  name                     = "audience-mapper"
  included_custom_audience = local.client.client_id
}