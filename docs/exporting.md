# Exporting client configurations

The `export/` directory contains a Go command-line tool (`ocm`) that connects to an
authorization server, reads its OIDC/OAuth2 client configurations, and writes them out in the
[canonical client configuration format](canonical-client-config-v0.1.md) as YAML or JSON files.

Supported source systems:

- Keycloak
- PingFederate

## Prerequisites

- Go (version 1.16 or later)

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
| `-format` | Output format: `yaml` or `json`. Defaults to `yaml`. |
| `-client-id` | Optional. Export a single client by its client ID instead of all clients. |

## Configuration

The tool is configured entirely through environment variables, since these typically hold
credentials that shouldn't be passed as command-line flags or committed to a config file.

Both source systems require:

- `AUTH_SERVER_BASE_URL`: Base URL of the authorization server's admin API.

### Keycloak

- `AUTH_SERVER_ACCESS_TOKEN`: A valid access token for the Keycloak admin API.

Keycloak requires obtaining an access token to consume the admin API. You can run the following
command to fetch a token and set the environment variable:

```bash
export AUTH_SERVER_ACCESS_TOKEN=$(curl -d "client_id=admin-cli" \
     -d "username=<username>" \
     -d "password=<password>" \
     -d "grant_type=password" \
     "https://<hostname>/realms/master/protocol/openid-connect/token" | jq -r .access_token)
```

### PingFederate

- `AUTH_SERVER_USERNAME`: Username for the PingFederate admin API.
- `AUTH_SERVER_PASSWORD`: Password for the PingFederate admin API.

## Next step

Once configurations are exported to `client-configurations/`, use the
[Terraform definitions in `import/`](importing.md) to provision them into an authorization server.
