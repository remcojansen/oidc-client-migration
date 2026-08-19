# Configuration reference

This is the single reference for every environment variable and Terraform variable used across
the export tool and the Terraform import. See [exporting.md](exporting.md) and
[importing.md](importing.md) for full usage instructions — this page only consolidates the
settings themselves.

## Export tool (`export/`, Go)

The export tool is configured entirely through environment variables (credentials shouldn't be
passed as command-line flags or committed to a config file).

| Variable | Applies to | Required | Description |
| --- | --- | --- | --- |
| `AUTH_SERVER_BASE_URL` | Keycloak, PingFederate | Yes | Base URL of the authorization server's admin API. |
| `AUTH_SERVER_ACCESS_TOKEN` | Keycloak | Yes (Keycloak only) | A valid access token for the Keycloak admin API. See [exporting.md](exporting.md#keycloak) for how to fetch one. |
| `AUTH_SERVER_USERNAME` | PingFederate | Yes (PingFederate only) | Username for the PingFederate admin API. |
| `AUTH_SERVER_PASSWORD` | PingFederate | Yes (PingFederate only) | Password for the PingFederate admin API. |
| `AUTH_SERVER_ACCESS_TOKEN_MANAGER_MAP` | PingFederate | No (default empty), required to populate `extensions.access_token_format`/`extensions.access_token_lifetime_seconds` | JSON object mapping PingFederate access token manager IDs to their format/lifetime, e.g. `{"my-jwt-manager-short": {"format": "jwt", "lifetime_seconds": 300}}`. Access token managers are user-created, deployment-specific resources with no universal naming convention, so this must be supplied per-deployment. |
| `AUTH_SERVER_MINIMUM_ACR_VALUE_PARAM_NAME` | PingFederate | No (default `minimum_acr_value`) | Name of the PingFederate Extended Parameter storing `extensions.minimum_acr_value`. PingFederate has no native concept for this; it's an arbitrary, deployment-chosen Extended Parameter name. |

## Import tool (`import/`, Terraform)

The Terraform import is configured through a mix of Terraform variables (settable via
`terraform.tfvars`, `-var`, or `TF_VAR_<name>` environment variables) and, for credentials, native
provider environment variables. Credentials should be supplied via environment variables rather
than committed to `terraform.tfvars`. Copy [`import/terraform.tfvars.example`](../import/terraform.tfvars.example)
to `import/terraform.tfvars` (gitignored) as a starting point.

### PingFederate

| Setting | Terraform variable | Environment variable | Required | Description |
| --- | --- | --- | --- | --- |
| Admin API URL | `provider_pingfederate_url` | `PINGFEDERATE_PROVIDER_HTTPS_HOST` | Yes | HTTPS URL for the PingFederate admin API. Leave the Terraform variable unset (default `null`) to use the environment variable instead. |
| Product version | `provider_pingfederate_product_version` | `PINGFEDERATE_PROVIDER_PRODUCT_VERSION` | Yes | PingFederate server version being configured. Leave the Terraform variable unset (default `null`) to use the environment variable instead. |
| Skip TLS verification | `provider_pingfederate_tls_skip_verify` | `PINGFEDERATE_PROVIDER_INSECURE_TRUST_ALL_TLS` | No (default `false`) | Trust any TLS certificate when connecting. |
| Bypass external validation header | `provider_pingfederate_bypass_external_validation_header` | `PINGFEDERATE_PROVIDER_X_BYPASS_EXTERNAL_VALIDATION_HEADER` | No (default `true`; falls back to the provider's own default if both the variable and environment variable are unset) | Sets the `X-BypassExternalValidation` admin API header. |
| Username | — | `PINGFEDERATE_PROVIDER_USERNAME` | Yes | Username for the PingFederate admin API. Not wired to a Terraform variable — set directly as an environment variable. |
| Password | — | `PINGFEDERATE_PROVIDER_PASSWORD` | Yes | Password for the PingFederate admin API. Not wired to a Terraform variable — set directly as an environment variable. |

### Keycloak

| Setting | Terraform variable | Environment variable | Required | Description |
| --- | --- | --- | --- | --- |
| Admin API URL | `provider_keycloak_url` | `KEYCLOAK_URL` | Yes | Base URL of the Keycloak instance. Leave the Terraform variable unset (default `null`) to use the environment variable instead. |
| Client ID | `provider_keycloak_client_id` | `KEYCLOAK_CLIENT_ID` | No | Client used for admin API authentication (password grant). Leave the Terraform variable unset (default `null`) to use the environment variable instead. |
| Skip TLS verification | `provider_keycloak_tls_skip_verify` | — | No (default `false`) | Trust any TLS certificate when connecting. Not supported as an environment variable by the Keycloak Terraform provider — must be set via this Terraform variable. |
| Username | `provider_keycloak_username` | `KEYCLOAK_USER` | Yes | Admin username. Leave the Terraform variable unset (default `null`) to use the `KEYCLOAK_USER` environment variable instead. |
| Password | `provider_keycloak_password` | `KEYCLOAK_PASSWORD` | Yes | Admin password. Leave the Terraform variable unset (default `null`) to use the `KEYCLOAK_PASSWORD` environment variable instead. |

### Shared / general

| Setting | Terraform variable | Required | Description |
| --- | --- | --- | --- |
| Canonical config path | `config_base` | No (default `../client-configurations`) | Path to the directory of canonical client configuration YAML files. |
| Keycloak realm | `keycloak_realm_id` | Yes | ID of the realm where Keycloak clients will be created. |
| PingFederate provisioning toggle | `pingfederate_enabled` | No (default `false`) | Enables/disables PingFederate client provisioning. |
| Keycloak provisioning toggle | `keycloak_enabled` | No (default `false`) | Enables/disables Keycloak client provisioning. |
| PingFederate access token manager mapping | `pingfederate_access_token_manager_mapping` | No (default `{}`), but required when `pingfederate_enabled` is `true` | Maps `extensions.access_token_format` (`jwt`/`opaque`) and `extensions.access_token_lifetime_seconds` to this PingFederate deployment's access token manager IDs. Access token managers are user-created resources with no universal naming convention, so this must be supplied per-deployment (see `terraform.tfvars.example`). |
| PingFederate minimum ACR value param name | `pingfederate_minimum_acr_value_param_name` | No (default `minimum_acr_value`) | Name of the PingFederate Extended Parameter storing `extensions.minimum_acr_value`. PingFederate has no native concept for this; it's an arbitrary, deployment-chosen Extended Parameter name. |
