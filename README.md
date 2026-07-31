# Client Migration tool

This project helps you move OIDC/OAuth2 client configurations between authorization servers. It
exports client configurations from an existing authorization server into a canonical,
system-agnostic format, which can then be provisioned into the same or a different authorization
server.

Currently supported authorization servers: Keycloak and PingFederate.

## How it works

```
+------------------------------+
|     Authorization server     |
|           (source)           |
|                              |
|   Keycloak / PingFederate    |
+------------------------------+
              |
              | export (export/)
              v
+------------------------------+
|    client-configurations/    |
|                              |
|   (canonical YAML / JSON)    |
+------------------------------+
              |
              | import (import/)
              v
+------------------------------+
|     Authorization server     |
|           (target)           |
|                              |
|   Keycloak / PingFederate    |
+------------------------------+
```

| Directory | Purpose |
| --- | --- |
| [`export/`](export/) | Go tool (`ocm`) that exports client configurations from an authorization server into the canonical format. |
| [`client-configurations/`](client-configurations/) | Canonical client configuration files: written by `export/`, read by `import/`. |
| [`import/`](import/) | Terraform/OpenTofu definitions that provision canonical client configurations into an authorization server. |

## Documentation

- [docs/exporting.md](docs/exporting.md): building and running the export tool.
- [docs/importing.md](docs/importing.md): configuring and running the Terraform import.
- [docs/configuration.md](docs/configuration.md): single reference for every environment variable
  and Terraform variable used by the export tool and the Terraform import.
- [docs/canonical-client-config-v0.1.md](docs/canonical-client-config-v0.1.md) /
  [docs/canonical-client-config-v0.1.schema.json](docs/canonical-client-config-v0.1.schema.json):
  the canonical client configuration model shared by both.

## Quick start

```bash
# 1. Build the export tool
make build

# 2. Export client configurations from Keycloak (see docs/exporting.md for credentials setup)
bin/ocm -source keycloak -dir client-configurations -format yaml

# 3. Provision the exported configurations into a target authorization server (see docs/importing.md)
cd import
terraform init
terraform apply
```
