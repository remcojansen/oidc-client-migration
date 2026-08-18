package pingfederate

import (
	"ocm/export/authserver"
	"reflect"
	"testing"
)

func TestGetCanonicalClientConfig(t *testing.T) {
	cfg := &PingFederateClientConfig{
		ClientID:     "my-client",
		Name:         "My Client",
		Enabled:      true,
		RedirectUris: []string{"https://app.example.com/callback"},
		GrantTypes:   []string{"AUTHORIZATION_CODE", "REFRESH_TOKEN", "CLIENT_CREDENTIALS", "DEVICE_CODE", "ACCESS_TOKEN_VALIDATION"},
		ClientAuth: AuthConfig{
			Type:            "SECRET",
			EncryptedSecret: "enc-s3cr3t",
		},
		DefaultAccessTokenManagerRef:      AccessTokenManager{Id: "jwt-long"},
		RequireProofKeyForCodeExchange:    true,
		PersistentGrantExpirationType:     "OVERRIDE_SERVER_DEFAULT",
		PersistentGrantExpirationTime:     2,
		PersistentGrantExpirationTimeUnit: "HOURS",
		accessTokenManagerMapping: map[string]AccessTokenManagerInfo{
			"jwt-long": {Format: authserver.AccessTokenFormatJwt, LifetimeSeconds: authserver.AccessTokenLifetimeLong},
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
		authserver.GrantTypeDeviceCode,
	}
	if !reflect.DeepEqual(canonical.Client.GrantTypes, wantGrantTypes) {
		t.Errorf("GrantTypes = %v, want %v", canonical.Client.GrantTypes, wantGrantTypes)
	}

	if !canonical.Extensions.IntrospectionEnabled {
		t.Error("expected IntrospectionEnabled to be true when ACCESS_TOKEN_VALIDATION is present")
	}

	if canonical.Client.TokenEndpointAuthMethod != authserver.AuthMethodClientSecretBasic {
		t.Errorf("TokenEndpointAuthMethod = %q, want %q", canonical.Client.TokenEndpointAuthMethod, authserver.AuthMethodClientSecretBasic)
	}

	if canonical.Extensions.AccessTokenFormat != authserver.AccessTokenFormatJwt {
		t.Errorf("AccessTokenFormat = %q, want %q", canonical.Extensions.AccessTokenFormat, authserver.AccessTokenFormatJwt)
	}

	if !canonical.Extensions.PKCERequired {
		t.Error("expected PKCERequired to be true")
	}

	if canonical.Extensions.OfflineSessionMaxLifetimeSeconds != 7200 {
		t.Errorf("OfflineSessionMaxLifetimeSeconds = %d, want 7200 (2 hours)", canonical.Extensions.OfflineSessionMaxLifetimeSeconds)
	}

	if canonical.Secrets.EncryptedSecret != "enc-s3cr3t" {
		t.Errorf("EncryptedSecret = %q, want %q", canonical.Secrets.EncryptedSecret, "enc-s3cr3t")
	}
}

func TestMapIntrospectionEnabledFalseWhenAbsent(t *testing.T) {
	cfg := &PingFederateClientConfig{
		GrantTypes: []string{"AUTHORIZATION_CODE"},
	}
	if cfg.mapIntrospectionEnabled() {
		t.Error("expected IntrospectionEnabled to be false when ACCESS_TOKEN_VALIDATION is absent")
	}
}

func TestCalculateSeconds(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		unit    string
		want    int
		wantErr bool
	}{
		{"seconds", 30, "SECONDS", 30, false},
		{"minutes", 5, "MINUTES", 300, false},
		{"hours", 2, "HOURS", 7200, false},
		{"days", 1, "DAYS", 86400, false},
		{"unknown unit", 1, "WEEKS", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculateSeconds(tt.value, tt.unit)
			if (err != nil) != tt.wantErr {
				t.Fatalf("calculateSeconds(%d, %q) error = %v, wantErr %v", tt.value, tt.unit, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("calculateSeconds(%d, %q) = %d, want %d", tt.value, tt.unit, got, tt.want)
			}
		})
	}
}

func TestContains(t *testing.T) {
	if !contains([]string{"a", "b", "c"}, "b") {
		t.Error("expected contains to find \"b\" in the slice")
	}
	if contains([]string{"a", "b", "c"}, "z") {
		t.Error("expected contains to not find \"z\" in the slice")
	}
}
