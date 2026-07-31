# Field support matrix (capabilities per authorization server)

This is the single source of truth for which fields of the
[canonical client configuration model](canonical-client-config-v0.1.md) are actually supported by
each authorization server, on export (`export/`, reading from the source system) and on import
(`import/`, provisioning into the target system). It replaces scattered "unsupported" comments in
individual Terraform modules and Go source files — see those files for the authoritative
implementation if this table and the code ever disagree.

## Legend

| Symbol | Meaning |
| --- | --- |
| ✅ | Fully supported — a direct, faithful mapping exists. |
| ~ | Partial / approximate — supported, but via a heuristic, a lossy mapping, or only a subset of the field's range. See the note. |
| ❌ | Unsupported — no equivalent exists in this system/direction; the field is dropped (export) or ignored (import). |

## `metadata`

All `metadata.*` fields (`owner_email`, `owner_slack_channel`, `owner_team_ad_group`) are
out-of-band in v0.1: neither exported from a source system nor provisioned to a target system.
They only round-trip through canonical files if set by hand.

## `client`

| Field | Keycloak export | Keycloak import | PingFederate export | PingFederate import |
| --- | --- | --- | --- | --- |
| `client_id` | ✅ | ✅ | ✅ | ✅ |
| `client_name` | ✅ | ✅ | ✅ | ✅ |
| `description` | ✅ | ✅ | ✅ | ✅ |
| `contacts` | ❌ | ❌ | ❌ | ❌ |
| `client_uri` | ✅ | ✅ | ❌ | ❌ |
| `logo_uri` | ✅ | ✅ | ✅ | ✅ |
| `tos_uri` | ✅ | ✅ | ❌ | ❌ |
| `policy_uri` | ✅ | ✅ | ❌ | ❌ |
| `jwks_uri` | ✅ | ✅ | ✅ | ✅ |
| `redirect_uris` | ✅ | ✅ | ✅ | ✅ |
| `response_types` | ~ [1] | ❌ [2] | ~ [1] | ❌ [2] |
| `grant_types` | ✅ | ✅ | ✅ | ✅ |
| `token_endpoint_auth_method` | ✅ | ✅ | ✅ | ✅ |
| `token_endpoint_auth_signing_alg` | ✅ | ✅ | ❌ | ❌ |
| `id_token_signed_response_alg` | ✅ | ✅ | ✅ | ✅ |
| `request_object_signing_alg` | ✅ | ✅ | ❌ | ❌ |
| `scopes` | ✅ | ✅ | ✅ | ✅ |
| `consent_required` | ✅ | ✅ | ✅ [3] | ✅ [3] |
| `application_type` | ~ [4] | ✅ [12] | ~ [4] | ❌ [2] |
| `subject_type` | ❌ [5] | ❌ | ✅ | ✅ |
| `sector_identifier_uri` | ❌ | ❌ | ❌ | ❌ |
| `backchannel_logout_uri` | ✅ | ✅ | ✅ | ✅ |
| `frontchannel_logout_uri` | ✅ | ✅ | ✅ | ✅ |
| `post_logout_redirect_uris` | ✅ | ✅ | ✅ | ✅ |
| `default_acr_values` | ✅ | ✅ | ❌ | ❌ |
| `initiate_login_uri` | ❌ | ❌ | ❌ | ❌ |
| `request_uris` | ✅ | ✅ | ❌ | ❌ |

## `extensions`

| Field | Keycloak export | Keycloak import | PingFederate export | PingFederate import |
| --- | --- | --- | --- | --- |
| `enabled` | ✅ | ✅ | ✅ | ✅ |
| `pkce_required` | ✅ | ✅ | ✅ | ✅ |
| `dpop_required` | ✅ | ✅ | ✅ | ✅ |
| `par_required` | ✅ | ✅ | ✅ | ✅ |
| `access_token_format` | ~ [6] | ❌ [7] | ~ [8] | ~ [8] |
| `access_token_lifetime_seconds` | ✅ | ✅ | ~ [8] | ~ [8] |
| `offline_session_max_lifetime_seconds` | ✅ | ✅ [10] | ✅ | ✅ |
| `offline_session_idle_timeout_seconds` | ✅ | ✅ [10] | ✅ | ✅ |
| `session_max_lifetime_seconds` | ✅ | ✅ | ❌ [11] | ❌ [11] |
| `session_idle_timeout_seconds` | ✅ | ✅ | ❌ [11] | ❌ [11] |
| `rotate_refresh_tokens` | ~ [9] | ❌ | ✅ | ✅ |
| `minimum_acr_value` | ✅ | ✅ | ✅ | ✅ |
| `require_terms_and_conditions_approval` | ✅ | ✅ | ✅ | ✅ |

## `secrets`

| Field | Keycloak export | Keycloak import | PingFederate export | PingFederate import |
| --- | --- | --- | --- | --- |
| `plain_secret` | ✅ (only form Keycloak provides) | ✅ | ❌ | ❌ |
| `encrypted_secret` | ❌ | ❌ | ✅ (only form PingFederate provides) | ✅ |

## Notes

1. **`response_types` (export, both systems)**: not read directly — derived from enabled grant
   types/flows (e.g. `authorization_code` flow implies the `code` response type).
2. **`response_types` / `application_type` (PingFederate import)**: not consumed. PingFederate
   client type (public vs. confidential) is derived entirely from `token_endpoint_auth_method`
   instead.
3. **`consent_required` (PingFederate)**: represented inversely as `bypass_approval_page`.
4. **`application_type` (export, both systems)**: derived, not read from a dedicated field —
   Keycloak from `publicClient`, PingFederate from `client_auth.type == "NONE"`.
5. **`subject_type` (Keycloak export)**: Keycloak's client API has no equivalent concept; always
   reported as `public`.
6. **`access_token_format` (Keycloak export)**: always reported as `jwt` — Keycloak issues JWT
   access tokens only in this mapping; no opaque-token detection.
7. **`access_token_format` (Keycloak import)**: reserved for future support of lightweight/opaque
   access tokens; currently not applied to the provisioned client.
8. **`access_token_format` / `access_token_lifetime_seconds` (PingFederate)**: PingFederate has no
   per-client format/lifetime fields — both are approximated by selecting one of a small,
   preconfigured set of access token managers (e.g. `jwtstandardshort` = JWT/300s,
   `jwtstandardlong` = JWT/1800s), so only those discrete combinations round-trip faithfully.
9. **`rotate_refresh_tokens` (Keycloak export)**: Keycloak has no native "rotate refresh tokens"
   toggle; approximated as `true` whenever refresh tokens are in use at all.
10. **`offline_session_max_lifetime_seconds` / `offline_session_idle_timeout_seconds` (Keycloak
    import)**: mapped to Keycloak's per-client `client.offline.session.max.lifespan` /
    `client.offline.session.idle.timeout` attributes, which only govern true persistent refresh
    tokens issued to clients that request the `offline_access` scope. Since setting these
    attributes has no effect otherwise, the import automatically adds `offline_access` to the
    client's optional scopes whenever either value is configured.
11. **`session_max_lifetime_seconds` / `session_idle_timeout_seconds` (PingFederate)**:
    PingFederate has no concept of a refresh token lifetime that is distinct from its persistent
    grant settings, so these fields have no PingFederate equivalent and are not read/applied by
    the PingFederate export/import.
12. **`application_type` (Keycloak import)**: since the field is optional in the canonical model,
    `access_type` defaults to `CONFIDENTIAL` unless `application_type` is explicitly `native`, to
    avoid accidentally provisioning an intended confidential client as public when the field is
    left unset.
