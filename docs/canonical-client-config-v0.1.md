# Canonical OIDC Client Configuration v0.1

This document describes the canonical client configuration model implemented in [oidcconfig/client_config.go](../export/oidcconfig/client_config.go).

The goal of the model is to capture enough information to:

1. export client registrations from PingFederate and Keycloak,
2. store them in a system-agnostic format, and
3. provision them back through Terraform/OpenTofu.

The model deliberately combines two layers within a single `client` block:

- Standard OAuth 2.0 / OpenID Connect registration fields, most of which are defined by an RFC or OpenID Connect specification.
- Project-specific authorization-server policy and provisioning fields that have no RFC or OIDC specification definition.

## Design Boundaries

The model is intended to be a canonical migration shape, not a strict RFC-only client registration document.

See the [Top-Level Structure](#top-level-structure) table below for what each top-level block
contains, and the [Field Reference](#field-reference) section for which `client` fields are
RFC/OIDC-defined versus project-specific.

## Top-Level Structure

| Field | Type | Required | Purpose |
| --- | --- | --- | --- |
| `annotations` | object | no | Ownership and source-tracking data. Out-of-band only (not provisioned). |
| `client` | object | yes | Client registration data combined with authorization-server policy and provisioning settings. |
| `secrets` | object | no | Sensitive credentials or encrypted secret material. |

## Field Reference

### annotations

Annotations are out-of-band. They are retained in canonical files for ownership and governance, but current Terraform modules do not apply annotations to provider resources.

| Field | Type | Notes |
| --- | --- | --- |
| `owner_email` | string | Human owner contact address. |
| `owner_team` | string | Human-readable name of the owning team or stakeholder. |
| `owner_channel` | string | Operational contact channel. |
| `owner_security_group` | string | Owning team directory security group. |

### client

Fields are grouped by purpose (identity, auth flow, token/session policy, advanced OIDC, registration
metadata) rather than by standards-origin. The **Standard** column indicates whether the field is
defined by an OAuth 2.0/OpenID Connect RFC or OpenID Connect specification; project-specific fields
have no RFC or OIDC specification definition and instead capture authorization-server policy or
provisioning behavior needed to migrate clients between systems.

| Field | Type | Standard | Notes |
| --- | --- | --- | --- |
| `client_id` | string | [RFC 6749 §2.2](https://www.rfc-editor.org/rfc/rfc6749#section-2.2) | Stable client identifier. Required. |
| `client_name` | string | [RFC 7591 §2](https://www.rfc-editor.org/rfc/rfc7591#section-2) | Human-readable client name. Required. |
| `description` | string | — (project-specific) | Optional description. |
| `enabled` | bool | — (project-specific) | Whether the client is enabled in the target authorization server. Required when `client` is present. Defaults to `true` when omitted. |
| `application_type` | string | [OIDC Dynamic Client Registration §2](https://openid.net/specs/openid-connect-registration-1_0.html#ClientMetadata) | `web` or `native`. In this model, `web` is assumed to correspond to a confidential client and `native` is assumed to correspond to a public client (see [Application Type Assumption](#application-type-assumption)). |
| `token_endpoint_auth_method` | string | [RFC 7591 §2](https://www.rfc-editor.org/rfc/rfc7591#section-2) | Client auth method such as `client_secret_basic`, `client_secret_jwt`, `private_key_jwt`, or `none`. |
| `grant_types` | string[] | [RFC 6749](https://www.rfc-editor.org/rfc/rfc6749), [RFC 8628](https://www.rfc-editor.org/rfc/rfc8628), [RFC 8693](https://www.rfc-editor.org/rfc/rfc8693) | Supported grant types, restricted to the exact identifiers registered by the relevant OAuth 2.0 RFCs: `authorization_code`, `implicit`, `client_credentials`, `password`, `refresh_token` (RFC 6749); `urn:ietf:params:oauth:grant-type:device_code` (RFC 8628); `urn:ietf:params:oauth:grant-type:token-exchange` (RFC 8693). Note: `implicit` is not an actual `grant_type` wire value (the implicit flow never sends one), but is included here as the conventional label for that RFC 6749 grant type. |
| `response_types` | string[] | [RFC 6749 §3.1.1](https://www.rfc-editor.org/rfc/rfc6749#section-3.1.1) | OAuth/OIDC response types such as `code` or `token id_token`. |
| `redirect_uris` | string[] | [RFC 6749 §3.1.2](https://www.rfc-editor.org/rfc/rfc6749#section-3.1.2) | Allowed redirect URIs. |
| `post_logout_redirect_uris` | string[] | [OIDC RP-Initiated Logout 1.0](https://openid.net/specs/openid-connect-rpinitiated-1_0.html) | Allowed post-logout redirect URIs. |
| `scopes` | string[] | [RFC 6749 §3.3](https://www.rfc-editor.org/rfc/rfc6749#section-3.3) | Scopes associated with the client. |
| `consent_required` | bool | — (project-specific) | Whether user consent is required. Defaults to `true` when omitted. |
| `pkce_required` | bool | — (project-specific; PKCE itself is [RFC 7636](https://www.rfc-editor.org/rfc/rfc7636)) | Whether PKCE is required. Defaults to `true` when omitted. |
| `dpop_required` | bool | — (project-specific; DPoP itself is [RFC 9449](https://www.rfc-editor.org/rfc/rfc9449)) | Whether DPoP is required. Defaults to `false` when omitted. |
| `par_required` | bool | — (project-specific; PAR itself is [RFC 9126](https://www.rfc-editor.org/rfc/rfc9126)) | Whether pushed authorization requests are required. Defaults to `false` when omitted. |
| `access_token_format` | string | — (project-specific; JWT access tokens are profiled by [RFC 9068](https://www.rfc-editor.org/rfc/rfc9068)) | `jwt` or `opaque`. Authorization-server policy, not registration metadata. |
| `access_token_lifetime_seconds` | integer | — (project-specific) | Access token lifetime in seconds. |
| `rotate_refresh_tokens` | bool | — (project-specific) | Whether refresh token rotation is enabled. |
| `offline_session_max_lifetime_seconds` | integer | — (project-specific) | Maximum lifetime, in seconds, of a persistent/offline refresh token — one that must remain usable independently of any browser SSO session (e.g. for native/mobile apps, or backend services refreshing tokens unattended). |
| `offline_session_idle_timeout_seconds` | integer | — (project-specific) | Idle timeout, in seconds, of a persistent/offline refresh token. |
| `session_max_lifetime_seconds` | integer | — (project-specific) | Maximum lifetime, in seconds, of a refresh token tied to the authorization server's browser SSO session. Has no PingFederate equivalent — PingFederate refresh tokens are always persistent grants, not session-bound. |
| `session_idle_timeout_seconds` | integer | — (project-specific) | Idle timeout, in seconds, of a refresh token tied to the authorization server's browser SSO session. Has no PingFederate equivalent. |
| `introspection_enabled` | bool | [RFC 7662](https://www.rfc-editor.org/rfc/rfc7662) | Whether the client is authorized to call the authorization server's token introspection endpoint to validate tokens. Defaults to `true` when omitted. |
| `minimum_acr_value` | string | — (project-specific; `acr_values` itself is [OIDC Core §3.1.2.1](https://openid.net/specs/openid-connect-core-1_0.html#AuthRequest)) | Minimum ACR value enforced by the authorization server, regardless of the requested ACR values in the authorization request. |
| `default_acr_values` | string[] | [OIDC Dynamic Client Registration §2](https://openid.net/specs/openid-connect-registration-1_0.html#ClientMetadata) | Default ACR values to apply when none are requested. |
| `subject_type` | string | [OIDC Core §8](https://openid.net/specs/openid-connect-core-1_0.html#SubjectIDTypes) | `public` or `pairwise`. |
| `sector_identifier_uri` | string | [OIDC Core §8.1](https://openid.net/specs/openid-connect-core-1_0.html#SectorIdentifierValidation) | Sector identifier for pairwise subject generation. |
| `id_token_signed_response_alg` | string | [OIDC Dynamic Client Registration §2](https://openid.net/specs/openid-connect-registration-1_0.html#ClientMetadata) | ID token signing algorithm. |
| `token_endpoint_auth_signing_alg` | string | [RFC 7591 §2](https://www.rfc-editor.org/rfc/rfc7591#section-2) | Signing algorithm for assertion-based client auth. |
| `request_object_signing_alg` | string | [RFC 7591 §2](https://www.rfc-editor.org/rfc/rfc7591#section-2), [RFC 9101](https://www.rfc-editor.org/rfc/rfc9101) | Request object signing algorithm. |
| `jwks_uri` | string | [RFC 7591 §2](https://www.rfc-editor.org/rfc/rfc7591#section-2) | Client JWKS endpoint. |
| `contacts` | string[] | [RFC 7591 §2](https://www.rfc-editor.org/rfc/rfc7591#section-2) | Registration contacts, ideally email addresses. |
| `client_uri` | string | [RFC 7591 §2](https://www.rfc-editor.org/rfc/rfc7591#section-2) | Client home or information URI. |
| `logo_uri` | string | [RFC 7591 §2](https://www.rfc-editor.org/rfc/rfc7591#section-2) | Logo URI. |
| `tos_uri` | string | [RFC 7591 §2](https://www.rfc-editor.org/rfc/rfc7591#section-2) | Terms-of-service URI. |
| `policy_uri` | string | [RFC 7591 §2](https://www.rfc-editor.org/rfc/rfc7591#section-2) | Privacy or security policy URI. |
| `backchannel_logout_uri` | string | [OIDC Back-Channel Logout 1.0](https://openid.net/specs/openid-connect-backchannel-1_0.html) | Back-channel logout URI. |
| `frontchannel_logout_uri` | string | [OIDC Front-Channel Logout 1.0](https://openid.net/specs/openid-connect-frontchannel-1_0.html) | Front-channel logout URI. |
| `initiate_login_uri` | string | [OIDC Dynamic Client Registration §2](https://openid.net/specs/openid-connect-registration-1_0.html#ClientMetadata) | Third-party initiated login URI. |
| `request_uris` | string[] | [RFC 7591 §2](https://www.rfc-editor.org/rfc/rfc7591#section-2), [RFC 9101](https://www.rfc-editor.org/rfc/rfc9101) | Allowed request object URIs. |

`consent_required`, `pkce_required`, `dpop_required`, `par_required`, and `rotate_refresh_tokens`
are always written out explicitly by `export/` (as `true` or
`false`), even when `false` — the underlying value is always determinable from the source system,
so it is never simply left out. The defaults noted above apply when these fields are omitted from
a hand-authored configuration; the `import/` Terraform modules apply them consistently regardless
of whether the enclosing `client` block is entirely absent or merely missing that field.

### secrets

| Field | Type | Notes |
| --- | --- | --- |
| `plain_secret` | string | Plain-text secret. INSECURE: Do not use for production clients. |
| `encrypted_secret` | string | Encrypted client secret or secret reference. |

## Application Type Assumption

For this model, the `application_type` field is interpreted as follows:

- `web` is assumed to correspond to a confidential client.
- `native` is assumed to correspond to a public client.

This is a project-level convention used to normalize source authorization-server data into a single canonical representation.

See the [Field Reference](#field-reference) section above for which `client` fields are
project-specific (no RFC/OIDC basis) versus RFC/OIDC-defined; `annotations.*` and
`secrets.encrypted_secret` are project-specific in their entirety.

## Example

```json
{
  "annotations": {
    "owner_email": "team@example.com",
    "owner_security_group": "ad-group-example"
  },
  "client": {
    "client_id": "example-client",
    "client_name": "Example Client",
    "enabled": true,
    "redirect_uris": ["https://app.example.com/callback"],
    "grant_types": ["authorization_code", "refresh_token"],
    "token_endpoint_auth_method": "client_secret_basic",
    "scopes": ["openid", "profile"],
    "subject_type": "public",
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
