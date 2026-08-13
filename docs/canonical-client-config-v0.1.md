# Canonical OIDC Client Configuration v0.1

This document describes the canonical client configuration model implemented in [oidcconfig/client_config.go](../export/oidcconfig/client_config.go).

The goal of the model is to capture enough information to:

1. export client registrations from PingFederate and Keycloak,
2. store them in a system-agnostic format, and
3. provision them back through Terraform/OpenTofu.

The model deliberately combines two layers:

- Standard OAuth 2.0 / OpenID Connect registration metadata.
- Project-specific deployment and policy extensions that are not part of the RFCs.

## Design Boundaries

The model is intended to be a canonical migration shape, not a strict RFC-only client registration document.

That means:

- `client` contains the interoperable registration data.
- `extensions` contains authorization-server policy knobs and implementation-specific settings.
- `metadata` contains ownership and operational tracking data.
  In v0.1, metadata is out-of-band and is not provisioned to target authorization servers.
- `secrets` contains sensitive material that should not be treated as ordinary configuration.

## Top-Level Structure

| Field | Type | Required | Purpose |
| --- | --- | --- | --- |
| `metadata` | object | no | Ownership and source-tracking data. Out-of-band only in v0.1 (not provisioned). |
| `client` | object | yes | Canonical client registration data. |
| `extensions` | object | no | Server-specific policy and deployment settings. |
| `secrets` | object | no | Sensitive credentials or encrypted secret material. |

## Field Reference

### metadata

Metadata is out-of-band in v0.1. It is retained in canonical files for ownership and governance, but current Terraform modules do not apply metadata to provider resources.

| Field | Type | Notes |
| --- | --- | --- |
| `owner_email` | string | Human owner contact address. |
| `owner_team` | string | Human-readable name of the owning team or stakeholder. |
| `owner_slack_channel` | string | Operational contact channel. |
| `owner_team_ad_group` | string | Owning team directory group. |

### client

| Field | Type | Notes |
| --- | --- | --- |
| `client_id` | string | Stable client identifier. Required. |
| `client_name` | string | Human-readable client name. Required. |
| `description` | string | Optional description. |
| `contacts` | string[] | Registration contacts, ideally email addresses. |
| `client_uri` | string | Client home or information URI. |
| `logo_uri` | string | Logo URI. |
| `tos_uri` | string | Terms-of-service URI. |
| `policy_uri` | string | Privacy or security policy URI. |
| `jwks_uri` | string | Client JWKS endpoint. |
| `redirect_uris` | string[] | Allowed redirect URIs. |
| `response_types` | string[] | OAuth/OIDC response types such as `code` or `token id_token`. |
| `grant_types` | string[] | Supported grant types. |
| `token_endpoint_auth_method` | string | Client auth method such as `client_secret_basic`, `client_secret_jwt`, `private_key_jwt`, or `none`. |
| `token_endpoint_auth_signing_alg` | string | Signing algorithm for assertion-based client auth. |
| `id_token_signed_response_alg` | string | ID token signing algorithm. |
| `request_object_signing_alg` | string | Request object signing algorithm. |
| `scopes` | string[] | Scopes associated with the client. |
| `consent_required` | bool | Whether consent is required. Defaults to `true` when omitted. |
| `application_type` | string | `web` or `native`. In this model, `web` is assumed to correspond to a confidential client and `native` is assumed to correspond to a public client. |
| `subject_type` | string | `public` or `pairwise`. |
| `sector_identifier_uri` | string | Sector identifier for pairwise subject generation. |
| `backchannel_logout_uri` | string | Back-channel logout URI. |
| `frontchannel_logout_uri` | string | Front-channel logout URI. |
| `post_logout_redirect_uris` | string[] | Allowed post-logout redirect URIs. |
| `default_acr_values` | string[] | Default ACR values when none are requested. |
| `initiate_login_uri` | string | Third-party initiated login URI. |
| `request_uris` | string[] | Allowed request object URIs. |

### extensions

| Field | Type | Notes |
| --- | --- | --- |
| `enabled` | bool | Whether the client is enabled in the target authorization server. Required when `extensions` is present. Defaults to `true` when omitted. |
| `pkce_required` | bool | Whether PKCE is required. Defaults to `true` when omitted. |
| `dpop_required` | bool | Whether DPoP is required. Defaults to `false` when omitted. |
| `par_required` | bool | Whether pushed authorization requests are required. Defaults to `false` when omitted. |
| `access_token_format` | string | `jwt` or `opaque`. Authorization-server policy, not registration metadata. |
| `access_token_lifetime_seconds` | integer | Access token lifetime in seconds. |
| `offline_session_max_lifetime_seconds` | integer | Maximum lifetime, in seconds, of a persistent/offline refresh token — one that must remain usable independently of any browser SSO session (e.g. for native/mobile apps, or backend services refreshing tokens unattended). |
| `offline_session_idle_timeout_seconds` | integer | Idle timeout, in seconds, of a persistent/offline refresh token. |
| `session_max_lifetime_seconds` | integer | Maximum lifetime, in seconds, of a refresh token tied to the authorization server's browser SSO session. Has no PingFederate equivalent — PingFederate refresh tokens are always persistent grants, not session-bound. |
| `session_idle_timeout_seconds` | integer | Idle timeout, in seconds, of a refresh token tied to the authorization server's browser SSO session. Has no PingFederate equivalent. |
| `rotate_refresh_tokens` | bool | Whether refresh token rotation is enabled. |
| `minimum_acr_value` | string | Minimum ACR value required by the authorization server. |
| `terms_and_conditions_required` | bool | Whether terms-and-conditions approval is required. Defaults to `true` when omitted. |

`consent_required`, `pkce_required`, `dpop_required`, `par_required`, `rotate_refresh_tokens`,
and `terms_and_conditions_required` are always written out explicitly by `export/` (as `true` or
`false`), even when `false` — the underlying value is always determinable from the source system,
so it is never simply left out. The defaults noted above apply when these fields are omitted from
a hand-authored configuration; the `import/` Terraform modules apply them consistently regardless
of whether the enclosing `client`/`extensions` block is entirely absent or merely missing that field.

### secrets

| Field | Type | Notes |
| --- | --- | --- |
| `plain_secret` | string | Plain-text secret. INSECURE: Do not use for production clients. |
| `encrypted_secret` | string | Encrypted client secret or secret reference. |

## Standards Notes

The following fields are directly aligned with OAuth 2.0 / OpenID Connect registration concepts:

- `client_id`
- `client_name`
- `client_uri`
- `logo_uri`
- `tos_uri`
- `policy_uri`
- `jwks_uri`
- `redirect_uris`
- `response_types`
- `grant_types`
- `token_endpoint_auth_method`
- `token_endpoint_auth_signing_alg`
- `id_token_signed_response_alg`
- `request_object_signing_alg`
- `scopes`
- `consent_required`
- `application_type`
- `subject_type`
- `sector_identifier_uri`
- `backchannel_logout_uri`
- `frontchannel_logout_uri`
- `post_logout_redirect_uris`
- `default_acr_values`
- `initiate_login_uri`
- `request_uris`

## Application Type Assumption

For this model, the `application_type` field is interpreted as follows:

- `web` is assumed to correspond to a confidential client.
- `native` is assumed to correspond to a public client.

This is a project-level convention used to normalize source authorization-server data into a single canonical representation.

The following fields are project-specific and should be treated as extensions:

- `metadata.*`
- `extensions.*`
- `secrets.encrypted_secret`

## Validation Guidance

Recommended validation rules for the canonical model:

- `client.client_id` and `client.client_name` should always be present.
- URI fields should contain absolute URIs.
- Array fields should be deduplicated and normalized when exported.
- `application_type` should remain limited to the model's documented values.
- `extensions.enabled` should be explicitly set whenever an `extensions` block is present.

## Example

```json
{
  "metadata": {
    "owner_email": "team@example.com",
    "owner_team_ad_group": "ad-group-example"
  },
  "client": {
    "client_id": "example-client",
    "client_name": "Example Client",
    "redirect_uris": ["https://app.example.com/callback"],
    "grant_types": ["authorization_code", "refresh_token"],
    "token_endpoint_auth_method": "client_secret_basic",
    "scopes": ["openid", "profile"],
    "subject_type": "public"
  },
  "extensions": {
    "enabled": true,
    "pkce_required": true,
    "access_token_format": "jwt",
    "access_token_lifetime_seconds": 300,
    "rotate_refresh_tokens": true
  },
  "secrets": {
    "encrypted_secret": "enc:v1:..."
  }
}
```

## Versioning Note

This document describes v0.1 of the canonical model. If the structure changes in a way that affects exported JSON/YAML compatibility, the spec should be versioned separately and the JSON Schema should be updated alongside the code.