# Importing client configurations

The `import/` directory contains Terraform/OpenTofu definitions that read canonical client
configurations and provision them as clients in an authorization server.

Supported target systems:

- Keycloak
- PingFederate

## Organization

- `main.tf`: Loads canonical client configuration files and instantiates provisioning modules for
  PingFederate and Keycloak.
- `variables.tf`: Defines root-level Terraform variables such as provider settings and
  `config_base`.
- `providers.tf`: Configures required providers and provider authentication settings.
- `terraform.tfvars`: Environment-specific Terraform variable values.
- `modules/oauth-client-pingfederate`: Reusable module for provisioning a client in PingFederate.
- `modules/oauth-client-keycloak`: Reusable module for provisioning a client in Keycloak.

## Configuration loading

Client configurations are loaded from YAML files in the path configured by `config_base` (default:
`../client-configurations/`, i.e. the [`client-configurations/`](../client-configurations/)
directory at the repository root — the same directory the export tool writes to).

Each YAML file is decoded and passed to both modules, which consume the canonical model (`client`,
`extensions`, `secrets`).

## Module toggles

Provisioning for each authorization server can be toggled independently using root variables:

- `pingfederate_enabled` (bool): Enables or disables PingFederate provisioning.
- `keycloak_enabled` (bool): Enables or disables Keycloak provisioning.

Example in `terraform.tfvars`:

```hcl
pingfederate_enabled = false
keycloak_enabled     = true
```

When a toggle is set to `false`, no resources for that authorization server are instantiated.

## Configuration

Set the following environment variables before running Terraform. See the
[configuration reference](configuration.md#import-tool-import-terraform) for the full list of
Terraform variables, environment variables, and defaults.

- `PINGFEDERATE_PROVIDER_USERNAME`: The username for the PingFederate admin API.
- `PINGFEDERATE_PROVIDER_PASSWORD`: The password for the PingFederate admin API.
- `KEYCLOAK_USER`: The username for the Keycloak admin API.
- `KEYCLOAK_PASSWORD`: The password for the Keycloak admin API.

Other environment-specific configuration is managed using Terraform variables, which can be set in
`terraform.tfvars` or passed via environment variables and/or command-line flags.

## Usage

Run the following commands from the `import/` directory:

```bash
terraform init
terraform plan
terraform apply
```

## Feature support

Not every canonical field has a direct equivalent on every target system. See
[capabilities.md](capabilities.md) for the full field support matrix across both export and
import, for both Keycloak and PingFederate.
