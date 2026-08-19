# Maps the canonical client configuration to Keycloak provider resources for OAuth clients.
#
# Note: This module is intended to be used by the provisioning/terraform/main.tf file, which 
# reads canonical client configuration from YAML files and passes them to this module.
#
# The following fields lack a direct Keycloak equivalent and are therefore not currently
# mapped to Keycloak provider resources:
# - client.contacts --> Unsupported
# - client.response_types --> Unsupported
# - client.initiate_login_uri --> Unsupported
# - client.sector_identifier_uri --> Unsupported
# - client.subject_type --> Unsupported
# - extensions.rotate_refresh_tokens --> Unsupported
# - extensions.access_token_format --> Support for lightweight access tokens intended for later
# - extensions.introspection_enabled --> Unsupported; no per-client toggle in Keycloak

locals {
  # Define which scopes should be available to all clients by default;
  # these scopes will be added to the scopes that are assigned to the client
  default_scopes = ["profile", "openid", "email", "address", "phone"]

  # Define the login theme to use for all clients
  login_theme = "keycloak"

  client     = var.config.client
  extensions = try(var.config.extensions, {})
  secrets    = try(var.config.secrets, {})

  # Normalize nullable list fields so downstream functions like join() never receive null.
  post_logout_redirect_uris = try(local.client.post_logout_redirect_uris, null) == null ? [] : local.client.post_logout_redirect_uris
  default_acr_values        = try(local.client.default_acr_values, null) == null ? [] : local.client.default_acr_values
  request_uris              = try(local.client.request_uris, null) == null ? [] : local.client.request_uris
  redirect_uris             = try(local.client.redirect_uris, null) == null ? [] : local.client.redirect_uris
  grant_types               = try(local.client.grant_types, null) == null ? [] : local.client.grant_types
  optional_scopes           = try(local.client.scopes, null) == null ? [] : local.client.scopes

  # Normalize nullable booleans for use in conditions/provider boolean fields.
  enabled          = try(local.extensions.enabled, null) != false
  consent_required = try(local.client.consent_required, null) != false
  par_required     = try(local.extensions.par_required, null) == true
  dpop_required    = try(local.extensions.dpop_required, null) == true
  pkce_required    = try(local.extensions.pkce_required, null) != false

  # Normalize nullable strings used in extra_config and URL toggles.
  frontchannel_logout_uri = try(local.client.frontchannel_logout_uri, null) == null ? "" : local.client.frontchannel_logout_uri
  backchannel_logout_uri  = try(local.client.backchannel_logout_uri, null) == null ? "" : local.client.backchannel_logout_uri
  jwks_uri                = try(local.client.jwks_uri, null) == null ? "" : local.client.jwks_uri

  # Offline session lifetime/idle timeout: the persistent-refresh-token attributes, only
  # honored by Keycloak for clients that request the offline_access scope, so we add that
  # scope below whenever either of these is configured.
  offline_session_max_lifetime_seconds = try(local.extensions.offline_session_max_lifetime_seconds, null)
  offline_session_idle_timeout_seconds = try(local.extensions.offline_session_idle_timeout_seconds, null)
  # Purely an internal computed flag (not a Keycloak provider setting) used below to decide
  # whether offline_access needs to be added to the client's optional scopes.
  needs_offline_access_scope = local.offline_session_max_lifetime_seconds != null || local.offline_session_idle_timeout_seconds != null

  # Session lifetime/idle timeout: governs refresh tokens tied to the browser SSO session
  # (i.e. clients not using offline_access).
  session_max_lifetime_seconds = try(local.extensions.session_max_lifetime_seconds, null)
  session_idle_timeout_seconds = try(local.extensions.session_idle_timeout_seconds, null)

  # Canonical application_type convention: native -> public, anything else (including unset,
  # which defaults to web) -> confidential. Defaulting to confidential is the safer choice,
  # since application_type is optional and a missing value should not silently downgrade an
  # intended confidential client to public.
  access_type = try(local.client.application_type, "") == "native" ? "PUBLIC" : "CONFIDENTIAL"

  client_authenticator_type = (
    try(local.client.token_endpoint_auth_method, "") == "client_secret_jwt" ? "client-secret-jwt" :
    try(local.client.token_endpoint_auth_method, "") == "private_key_jwt" ? "client-jwt" :
    try(local.client.token_endpoint_auth_method, "") == "none" ? null :
    local.access_type == "PUBLIC" ? null : "client-secret"
  )

  extra_config = merge(
    {
      "par.required"                    = tostring(local.par_required)
      "dpop.required"                   = tostring(local.dpop_required)
      "id.token.signed.response.alg"    = try(local.client.id_token_signed_response_alg, "")
      "token.endpoint.auth.signing.alg" = try(local.client.token_endpoint_auth_signing_alg, "")
      "request.object.signature.alg"    = try(local.client.request_object_signing_alg, "")
      "minimum.acr.value"               = try(local.extensions.minimum_acr_value, "")
      "default.acr.values"              = join(",", local.default_acr_values)
      "logoUri"                         = try(local.client.logo_uri, "")
      "tosUri"                          = try(local.client.tos_uri, "")
      "policyUri"                       = try(local.client.policy_uri, "")
      "jwks.url"                        = local.jwks_uri
      "use.jwks.url"                    = tostring(local.jwks_uri != "")
      "request.uris"                    = join(",", local.request_uris)
    }
  )
}

resource "keycloak_openid_client" "this" {
  enabled     = local.enabled
  realm_id    = var.realm_id
  client_id   = local.client.client_id
  name        = local.client.client_name
  description = try(local.client.description, null)

  client_secret             = try(local.secrets.plain_secret, null)
  access_type               = local.access_type
  client_authenticator_type = local.client_authenticator_type

  consent_required    = local.consent_required
  valid_redirect_uris = local.redirect_uris

  standard_flow_enabled                     = contains(local.grant_types, "authorization_code")
  implicit_flow_enabled                     = contains(local.grant_types, "implicit")
  direct_access_grants_enabled              = contains(local.grant_types, "password")
  service_accounts_enabled                  = contains(local.grant_types, "client_credentials")
  oauth2_device_authorization_grant_enabled = contains(local.grant_types, "urn:ietf:params:oauth:grant-type:device_code")
  standard_token_exchange_enabled           = contains(local.grant_types, "urn:ietf:params:oauth:grant-type:token-exchange")
  use_refresh_tokens                        = contains(local.grant_types, "refresh_token")

  pkce_code_challenge_method = local.pkce_required ? "S256" : null

  extra_config = local.extra_config

  access_token_lifespan = try(local.extensions.access_token_lifetime_seconds, null)

  client_session_max_lifespan         = local.session_max_lifetime_seconds != null ? tostring(local.session_max_lifetime_seconds) : null
  client_session_idle_timeout         = local.session_idle_timeout_seconds != null ? tostring(local.session_idle_timeout_seconds) : null
  client_offline_session_max_lifespan = local.offline_session_max_lifetime_seconds != null ? tostring(local.offline_session_max_lifetime_seconds) : null
  client_offline_session_idle_timeout = local.offline_session_idle_timeout_seconds != null ? tostring(local.offline_session_idle_timeout_seconds) : null

  frontchannel_logout_enabled     = local.frontchannel_logout_uri != ""
  frontchannel_logout_url         = local.frontchannel_logout_uri != "" ? local.frontchannel_logout_uri : null
  backchannel_logout_url          = local.backchannel_logout_uri != "" ? local.backchannel_logout_uri : null
  valid_post_logout_redirect_uris = local.post_logout_redirect_uris

  # defaults
  admin_url                           = null
  root_url                            = null
  base_url                            = try(local.client.client_uri, null)
  always_display_in_console           = true
  backchannel_logout_session_required = false
  consent_screen_text                 = ""
  login_theme                         = local.login_theme
  web_origins                         = []
}

resource "keycloak_openid_client_optional_scopes" "this" {
  realm_id  = var.realm_id
  client_id = keycloak_openid_client.this.id

  optional_scopes = distinct(concat(
    local.default_scopes,
    local.optional_scopes,
    local.needs_offline_access_scope ? ["offline_access"] : []
  ))
}

# Create an audience mapping to set the client ID as audience
resource "keycloak_openid_audience_protocol_mapper" "audience_mapper" {
  realm_id                 = var.realm_id
  client_id                = keycloak_openid_client.this.id
  name                     = "audience-mapper"
  included_custom_audience = local.client.client_id
}