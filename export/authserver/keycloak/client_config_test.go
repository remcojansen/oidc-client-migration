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

func TestMapAccessTokenFormat_AlwaysJwt(t *testing.T) {
	cfg := &KeycloakClientConfig{}
	if got := cfg.mapAccessTokenFormat(); got != authserver.AccessTokenFormatJwt {
		t.Errorf("mapAccessTokenFormat() = %q, want %q (Keycloak always issues JWTs)", got, authserver.AccessTokenFormatJwt)
	}
}
