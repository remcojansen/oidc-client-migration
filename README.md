# Client Migration tool

This tool helps you export OIDC or oAuth2 client configurations from an existing authorization server and store the configurations in a canonical, system-agnostic format.

The canonical client model is documented in [docs/canonical-client-config-v0.1.md](docs/canonical-client-config-v0.1.md) and [docs/canonical-client-config-v0.1.schema.json](docs/canonical-client-config-v0.1.schema.json).

The canonical format can be used to provision configurations to the same or a different authorization server by reading from it using an available Terraform provider or a custom client.

## Prerequisites

- Go (version 1.16 or later)

## Build

To build the tool, run the following command in the current directory:

```bash
make build
```

## Usage

Run the tool with the following command:

```bash
bin/ocm -source <auth-server> -dir <path-to-configurations> -format <yaml|json>
``` 

- `-source`: Indicate which authorization server to export configuration from.
- `-dir`: Specifies the path to the directory containing the client configurations. This should point to the `provisioning/client-configurations/` directory in this repository.
- `-format`: Specifies the output format for the generated files. It can be either `yaml` or `json`.

## Keycloak

Keycloak requires obtaining an access token to consume the admin API. You can run the following command to fetch a token and set the respective environment variable:

```bash
export AUTH_SERVER_ACCESS_TOKEN=$(curl -d "client_id=admin-cli" \
     -d "username=<username>" \
     -d "password=<password>" \
     -d "grant_type=password" \
     "https://<hostname>/realms/master/protocol/openid-connect/token" | jq -r .access_token)
```

## Provisioning client configurations

One way of using the generated canonical configurations is to provision these to an authorization server using Terraform / OpenTofu. The `provisioning/` directory contains example configurations that show how the canonical configuration can be transformed into valid HCL for the supported authorization servers. 