# import

Terraform/OpenTofu definitions that read canonical client configurations from
[`client-configurations/`](../client-configurations/) and provision them as clients in an
authorization server.

Supports Keycloak and PingFederate as target systems, each behind its own module in `modules/`.

See [docs/importing.md](../docs/importing.md) for configuration and usage instructions.
