# Example client configurations

Sample canonical client configuration files illustrating realistic, distinct client
patterns supported by the [canonical model](canonical-client-config-v0.1.md). Together,
these examples exercise every field defined in the
[JSON Schema](canonical-client-config-v0.1.schema.json).

These are illustrative only — not consumed by `export/` or `import/` at runtime. Any
secret values are placeholders and must be replaced before real use.

| File | Scenario |
| --- | --- |
| `public-native-mobile-pkce.yaml` | Native iOS/Android app, public client, PKCE required. |
| `public-spa-with-dpop.yaml` | Browser SPA, public client, PKCE + DPoP-bound tokens. |
| `confidential-web-app.yaml` | Server-rendered web app, `client_secret_basic`, full standard metadata. |
| `confidential-private-key-jwt.yaml` | Confidential client authenticating via `private_key_jwt` and a published JWKS. |
| `confidential-client-credentials-service.yaml` | Machine-to-machine service account, `client_credentials` grant only. |
| `public-device-code-smarttv.yaml` | Smart-TV app using the device authorization grant. |
| `confidential-high-assurance-par-dpop.yaml` | High-assurance client requiring PAR, PKCE, DPoP, pairwise subjects, and step-up ACR. |
| `confidential-token-exchange-gateway.yaml` | Internal API gateway using RFC 8693 token exchange with opaque access tokens. |
| `confidential-disabled-legacy-client.yaml` | Deprecated client disabled at the authorization server (`client.enabled: false`). |
