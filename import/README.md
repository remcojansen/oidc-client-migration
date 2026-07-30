# Terraform definitions for client configurations

This directory contains Terraform definitions for provisioning canonical client configurations.

## Organization

The directory is organized as follows:

- `main.tf`: Loads canonical client configuration files and instantiates provisioning modules for PingFederate and Keycloak.
- `variables.tf`: Defines root-level Terraform variables such as provider settings and `config_base`.
- `providers.tf`: Configures required providers and provider authentication settings.
- `terraform.tfvars`: Environment-specific Terraform variable values.
- `modules/oauth-client-pingfederate`: Reusable module for provisioning a client in PingFederate.
- `modules/oauth-client-keycloak`: Reusable module for provisioning a client in Keycloak.

## Module Toggles

Provisioning for each authorization server can be toggled independently using root variables:

- `pingfederate_enabled` (bool): Enables or disables PingFederate provisioning.
- `keycloak_enabled` (bool): Enables or disables Keycloak provisioning.

Example in `terraform.tfvars`:

```hcl
pingfederate_enabled = false
keycloak_enabled     = true
```

When a toggle is set to `false`, no resources for that authorization server are instantiated.

## Configuration Loading

Client configurations are loaded from YAML files in the path configured by `config_base` (default: `../client-configurations/`).

Each YAML file is decoded and passed to both modules:

- `modules/oauth-client-pingfederate`
- `modules/oauth-client-keycloak`

Both modules consume the canonical model (`client`, `extensions`, `secrets`).

## Usage

To use these Terraform definitions, make sure you have the following environment variables set:
- `PINGFEDERATE_PROVIDER_USERNAME`: The username for the PingFederate admin API.
- `PINGFEDERATE_PROVIDER_PASSWORD`: The password for the PingFederate admin API.
- `KEYCLOAK_USER`: The username for the Keycloak admin API.
- `KEYCLOAK_PASSWORD`: The password for the Keycloak admin API.

Other environment-specific configuration is managed using Terraform variables, which can be set in `terraform.tfvars` or passed via environment variables and/or command-line flags.

When the environment variables are set, run the following commands in this directory:
```bash
terraform init
terraform plan
terraform apply
```


