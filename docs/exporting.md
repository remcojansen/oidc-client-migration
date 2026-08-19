# Exporting client configurations

The `export/` directory contains a Go command-line tool (`ocm` — OIDC Client Migration) that connects to an
authorization server, reads its OIDC/OAuth2 client configurations, and writes them out in the
[canonical client configuration format](canonical-client-config-v0.1.md) as YAML or JSON files.

Supported source systems:

- Keycloak
- PingFederate

## Prerequisites

- Go (version 1.26 or later, matching the `go` directive in [go.mod](../go.mod))

## Build

```bash
make build
```

This produces the binary at `bin/ocm`.

## Usage

```bash
bin/ocm -source <auth-server> -dir <path-to-configurations> -format <yaml|json>
```

| Flag | Description |
| --- | --- |
| `-source` | Which authorization server to export from: `keycloak` or `pingfederate`. Defaults to `keycloak`. |
| `-dir` | Directory to write the generated canonical configuration files to. This should point to the [`client-configurations/`](../client-configurations/) directory in this repository. |
| `-format` | Output format: `yaml` or `json`. Defaults to `yaml`. Note: the Terraform `import/` module currently only reads `.yaml` files from `client-configurations/`, so use `json` only for inspection/tooling purposes outside of this repository's import workflow. |
| `-client-id` | Optional. Export a single client by its client ID instead of all clients. |
| `-verbose` | Optional. Print the URL and response status code for each HTTP request made to the authorization server. |

## Configuration

The tool is configured entirely through environment variables, since these typically hold
credentials that shouldn't be passed as command-line flags or committed to a config file. See the
[configuration reference](configuration.md#export-tool-export-go) for the full list of variables.

Set `AUTH_SERVER_BASE_URL` to the base URL of the authorization server's admin API. This is
required for both Keycloak and PingFederate.

### Keycloak

Keycloak requires obtaining an access token to consume the admin API. You can run the following
command to fetch a token and set the `AUTH_SERVER_ACCESS_TOKEN` environment variable:

```bash
export AUTH_SERVER_ACCESS_TOKEN=$(curl -d "client_id=admin-cli" \
     -d "username=<username>" \
     -d "password=<password>" \
     -d "grant_type=password" \
     "https://<hostname>/realms/master/protocol/openid-connect/token" | jq -r .access_token)
```

### PingFederate

Set `AUTH_SERVER_USERNAME` and `AUTH_SERVER_PASSWORD` for the PingFederate admin API.

PingFederate has no per-client access token format/lifetime fields — these are approximated by
selecting one of the deployment's preconfigured access token managers. Set
`AUTH_SERVER_ACCESS_TOKEN_MANAGER_MAP` to a JSON object mapping your PingFederate access token
manager IDs to their format/lifetime, so the exporter can resolve
`client.access_token_format`/`client.access_token_lifetime_seconds`, e.g.:

```bash
export AUTH_SERVER_ACCESS_TOKEN_MANAGER_MAP='{
  "my-jwt-manager-short": {"format": "jwt", "lifetime_seconds": 300},
  "my-jwt-manager-long": {"format": "jwt", "lifetime_seconds": 1800}
}'
```

`client.minimum_acr_value` also has no native PingFederate equivalent — it's read from a
deployment-chosen Extended Parameter. If your deployment uses a name other than the default
(`minimum_acr_value`), set `AUTH_SERVER_MINIMUM_ACR_VALUE_PARAM_NAME` accordingly.

## Next step

Once configurations are exported to `client-configurations/`, use the
[Terraform definitions in `import/`](importing.md) to provision them into an authorization server.
