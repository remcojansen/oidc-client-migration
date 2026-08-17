package keycloak

import (
	"ocm/export/authserver"
	"reflect"
	"testing"
)

func TestGetCanonicalClientConfig(t *testing.T) {
	cfg := &KeycloakClientConfig{
		ClientID:                  "my-client",
		Name:                      "My Client",
		Enabled:                   true,
		PublicClient:              false,
		ClientAuthenticatorType:   "client-secret",
		Secret:                    "s3cr3t",
		RedirectUris:              []string{"https://app.example.com/callback"},
		StandardFlowEnabled:       true,
		DirectAccessGrantsEnabled: true,
		ServiceAccountsEnabled:    true,
		DefaultClientScopes:       []string{"profile"},
		OptionalClientScopes:      []string{"offline_access"},
		Attributes: map[string]string{
			"use.refresh.tokens":         "true",
			"access.token.lifespan":      "600",
			"pkce.code.challenge.method": "S256",
		},
	}

	canonical := cfg.GetCanonicalClientConfig()

	if canonical.Client.ClientID != "my-client" {
		t.Errorf("ClientID = %q, want %q", canonical.Client.ClientID, "my-client")
	}

	wantGrantTypes := []string{
		authserver.GrantTypeAuthorizationCode,
		authserver.GrantTypeRefreshToken,
		authserver.GrantTypeClientCredentials,
		authserver.GrantTypeResourceOwnerPassword,
	}
	if !reflect.DeepEqual(canonical.Client.GrantTypes, wantGrantTypes) {
		t.Errorf("GrantTypes = %v, want %v", canonical.Client.GrantTypes, wantGrantTypes)
	}

	if canonical.Client.TokenEndpointAuthMethod != authserver.AuthMethodClientSecretBasic {
		t.Errorf("TokenEndpointAuthMethod = %q, want client_secret_basic", canonical.Client.TokenEndpointAuthMethod)
	}

	if !canonical.Extensions.PKCERequired {
		t.Error("expected PKCERequired to be true when pkce.code.challenge.method is set to S256")
	}

	if canonical.Extensions.AccessTokenLifetimeSeconds != 600 {
		t.Errorf("AccessTokenLifetimeSeconds = %d, want 600", canonical.Extensions.AccessTokenLifetimeSeconds)
	}

	if canonical.Secrets.PlainSecret != "s3cr3t" {
		t.Errorf("PlainSecret = %q, want %q", canonical.Secrets.PlainSecret, "s3cr3t")
	}
}

func TestMapTokenEndpointAuthMethod_PublicClient(t *testing.T) {
	cfg := &KeycloakClientConfig{PublicClient: true}
	if got := cfg.mapTokenEndpointAuthMethod(); got != authserver.AuthMethodNone {
		t.Errorf("mapTokenEndpointAuthMethod() = %q, want %q", got, authserver.AuthMethodNone)
	}
}

func TestMapTokenEndpointAuthMethod_PublicClientWithDefaultClientAuthenticatorType(t *testing.T) {
	// Keycloak sets clientAuthenticatorType to "client-secret" by default on every client,
	// including public ones, where it is ignored. A public client must still map to "none".
	cfg := &KeycloakClientConfig{PublicClient: true, ClientAuthenticatorType: "client-secret"}
	if got := cfg.mapTokenEndpointAuthMethod(); got != authserver.AuthMethodNone {
		t.Errorf("mapTokenEndpointAuthMethod() = %q, want %q", got, authserver.AuthMethodNone)
	}
}

func TestMapAccessTokenFormat_AlwaysJwt(t *testing.T) {
	cfg := &KeycloakClientConfig{}
	if got := cfg.mapAccessTokenFormat(); got != authserver.AccessTokenFormatJwt {
		t.Errorf("mapAccessTokenFormat() = %q, want %q (Keycloak always issues JWTs)", got, authserver.AccessTokenFormatJwt)
	}
}

func TestMapGrantTypes_DeviceCode(t *testing.T) {
	cfg := &KeycloakClientConfig{
		Attributes: map[string]string{
			"oauth2.device.authorization.grant.enabled": "true",
		},
	}
	got := cfg.mapGrantTypes()
	want := []string{authserver.GrantTypeDeviceCode}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("mapGrantTypes() = %v, want %v", got, want)
	}
}

func TestMapGrantTypes_DeviceCodeNotEnabled(t *testing.T) {
	cfg := &KeycloakClientConfig{
		Attributes: map[string]string{
			"oauth2.device.authorization.grant.enabled": "false",
		},
	}
	got := cfg.mapGrantTypes()
	if len(got) != 0 {
		t.Errorf("mapGrantTypes() = %v, want empty", got)
	}
}

func TestMapIntrospectionEnabled_AlwaysTrue(t *testing.T) {
	cfg := &KeycloakClientConfig{}
	if !cfg.mapIntrospectionEnabled() {
		t.Error("expected mapIntrospectionEnabled() to always be true for Keycloak")
	}
}

func TestMapOfflineSessionMaxLifetimeSeconds(t *testing.T) {
	cfg := &KeycloakClientConfig{
		Attributes: map[string]string{
			"client.offline.session.max.lifespan": "259200",
		},
	}
	if got := cfg.mapOfflineSessionMaxLifetimeSeconds(); got != 259200 {
		t.Errorf("mapOfflineSessionMaxLifetimeSeconds() = %d, want 259200", got)
	}
}

func TestMapOfflineSessionIdleTimeoutSeconds(t *testing.T) {
	cfg := &KeycloakClientConfig{
		Attributes: map[string]string{
			"client.offline.session.idle.timeout": "86400",
		},
	}
	if got := cfg.mapOfflineSessionIdleTimeoutSeconds(); got != 86400 {
		t.Errorf("mapOfflineSessionIdleTimeoutSeconds() = %d, want 86400", got)
	}
}

func TestMapSessionMaxLifetimeSeconds(t *testing.T) {
	cfg := &KeycloakClientConfig{
		Attributes: map[string]string{
			"client.session.max.lifespan": "3600",
		},
	}
	if got := cfg.mapSessionMaxLifetimeSeconds(); got != 3600 {
		t.Errorf("mapSessionMaxLifetimeSeconds() = %d, want 3600", got)
	}
}

func TestMapSessionIdleTimeoutSeconds(t *testing.T) {
	cfg := &KeycloakClientConfig{
		Attributes: map[string]string{
			"client.session.idle.timeout": "1800",
		},
	}
	if got := cfg.mapSessionIdleTimeoutSeconds(); got != 1800 {
		t.Errorf("mapSessionIdleTimeoutSeconds() = %d, want 1800", got)
	}
}
