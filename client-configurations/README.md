# client-configurations

Canonical client configuration files (YAML or JSON), one per client.

This directory is the hand-off point between [`export/`](../export/) and [`import/`](../import/):

- `export/` writes canonical configuration files here.
- `import/` reads canonical configuration files from here to provision clients.

The configuration files themselves are generated data and are not committed to this repository
(see `.gitignore`); only this README is tracked so the directory exists by default.

See the [canonical client configuration model](../docs/canonical-client-config-v0.1.md) for the
file format.
