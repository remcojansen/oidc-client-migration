package authserver

import (
	"encoding/json"
	"ocm/oidcconfig"
	"strconv"
	"strings"
)

func CreateKeycloakClientConfig() *KeycloakClientConfig {
	return &KeycloakClientConfig{}
}

type KeycloakClientConfig struct {
	// Basic fields
	ID                  string `json:"id"`
	ClientID            string `json:"clientId"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	Enabled             bool   `json:"enabled"`
	Protocol            string `json:"protocol"`
	PublicClient        bool   `json:"publicClient"`
	ClientAuthenticator string `json:"clientAuthenticatorType"`
	Secret              string `json:"secret"`

	// OAuth/OIDC flows and URIs
	RedirectUris           []string `json:"redirectUris"`
	WebOrigins             []string `json:"webOrigins"`
	StandardFlowEnabled    bool     `json:"standardFlowEnabled"`
	ImplicitFlowEnabled    bool     `json:"implicitFlowEnabled"`
	DirectAccessGrants     bool     `json:"directAccessGrantsEnabled"`
	ServiceAccountsEnabled bool     `json:"serviceAccountsEnabled"`

	// OIDC specific
	UseRefreshTokens   bool `json:"useRefreshTokens"`
	ConsentRequired    bool `json:"consentRequired"`
	FrontchannelLogout bool `json:"frontchannelLogoutSessionRequired"`

	// Security
	FullScopeAllowed     bool     `json:"fullScopeAllowed"`
	DefaultClientScopes  []string `json:"defaultClientScopes"`
	OptionalClientScopes []string `json:"optionalClientScopes"`

	// Custom attributes for PKCE, DPoP, PAR, etc.
	Attributes map[string]string `json:"attributes"`

	// Other fields
	Access map[string]bool `json:"access"`
}

func (cc *KeycloakClientConfig) ToJSON() ([]byte, error) {
	jsonData, err := json.MarshalIndent(cc, "", "  ")
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func (c *KeycloakClientConfig) GetCanonicalClientConfig() *oidcconfig.CanonicalClientConfig {
	return &oidcconfig.CanonicalClientConfig{
		Metadata: oidcconfig.CanonicalClientConfigMetadata{
			// Not provided by Keycloak API
		},
		Client: oidcconfig.CanonicalClientConfigClient{
			ClientID:                c.mapClientID(),
			Name:                    c.mapName(),
			Description:             c.mapDescription(),
			RedirectURIs:            c.mapRedirectURIs(),
			ResponseTypes:           c.mapResponseTypes(),
			GrantTypes:              c.mapGrantTypes(),
			TokenEndpointAuthMethod: c.mapTokenEndpointAuthMethod(),
			Scope:                   c.mapScopes(),
			BackchannelLogoutURIs:   c.mapBackchannelLogoutURIs(),
			FrontchannelLogoutURIs:  c.mapFrontchannelLogoutURIs(),
			PostLogoutRedirectURIs:  c.mapPostLogoutRedirectURIs(),
			Enabled:                 c.Enabled,
			ClientType:              c.mapClientType(),
			PKCERequired:            c.mapPKCERequired(),
			DPoPRequired:            c.mapDPoPRequired(),
			PARRequired:             c.mapPARRequired(),
			PromptScopeApproval:     c.mapPromptScopeApproval(),
			AccessTokenFormat:       c.mapAccessTokenFormat(),
			AccessTokenLifetime:     c.mapAccessTokenLifetime(),
			RotateRefreshTokens:     c.mapRotateRefreshTokens(),
		},
		Extensions: oidcconfig.CanonicalClientConfigExtensions{
			SSO:                      c.mapSSO(),
			Enforce2SV:               c.mapEnforce2SV(),
			PromptTermsAndConditions: c.mapPromptTermsAndConditions(),
		},
		Secrets: oidcconfig.CanonicalClientConfigSecrets{
			PlainSecret: c.mapClientSecret(),
		},
	}
}

func (c *KeycloakClientConfig) GetClientID() string {
	return c.ClientID
}

func (c *KeycloakClientConfig) mapClientID() string {
	return c.ClientID
}

func (c *KeycloakClientConfig) mapName() string {
	return c.Name
}

func (c *KeycloakClientConfig) mapDescription() string {
	return c.Description
}

func (c *KeycloakClientConfig) mapRedirectURIs() []string {
	return c.RedirectUris
}

func (c *KeycloakClientConfig) mapResponseTypes() []string {
	// Derive response types from enabled flows
	var responseTypes []string
	if c.StandardFlowEnabled {
		responseTypes = append(responseTypes, "code")
	}
	if c.ImplicitFlowEnabled {
		responseTypes = append(responseTypes, "token", "id_token")
	}
	return responseTypes
}

func (c *KeycloakClientConfig) mapGrantTypes() []string {
	var g []string

	if c.StandardFlowEnabled {
		g = append(g, "authorization_code")
	}
	if c.ImplicitFlowEnabled {
		g = append(g, "implicit")
	}
	if c.UseRefreshTokens {
		g = append(g, "refresh_token")
	}
	if c.ServiceAccountsEnabled {
		g = append(g, "client_credentials")
	}
	if c.DirectAccessGrants {
		g = append(g, "password")
	}

	// Check attributes for additional grant types
	if val, ok := c.Attributes["token.response.type"]; ok && val == "device_code" {
		g = append(g, "device_code")
	}

	return g
}

func (c *KeycloakClientConfig) mapTokenEndpointAuthMethod() string {
	// Map client authenticator type to standard method
	switch c.ClientAuthenticator {
	case "client-secret":
		return AuthMethodClientSecretBasic
	case "client-secret-jwt":
		return AuthMethodClientSecretJwt
	case "client-jwt":
		return AuthMethodPrivateKeyJwt
	default:
		if c.PublicClient {
			return AuthMethodNone
		}
		return AuthMethodClientSecretBasic
	}
}

func (c *KeycloakClientConfig) mapScopes() string {
	var scopes []string
	scopes = append(scopes, c.DefaultClientScopes...)
	scopes = append(scopes, c.OptionalClientScopes...)
	return strings.Join(scopes, ",")
}

func (c *KeycloakClientConfig) mapBackchannelLogoutURIs() []string {
	if val, ok := c.Attributes["backchannel.logout.url"]; ok && val != "" {
		return strings.Split(val, ",")
	}
	return nil
}

func (c *KeycloakClientConfig) mapFrontchannelLogoutURIs() []string {
	if c.FrontchannelLogout {
		if val, ok := c.Attributes["frontchannel.logout.url"]; ok && val != "" {
			return strings.Split(val, ",")
		}
	}
	return nil
}

func (c *KeycloakClientConfig) mapPostLogoutRedirectURIs() []string {
	if val, ok := c.Attributes["post.logout.redirect.uris"]; ok && val != "" {
		return strings.Split(val, ",")
	}
	return nil
}

func (c *KeycloakClientConfig) mapClientType() string {
	if c.PublicClient {
		return ClientTypePublic
	}
	return ClientTypeConfidential
}

func (c *KeycloakClientConfig) mapPKCERequired() bool {
	if val, ok := c.Attributes["pkce.code.challenge.method"]; ok {
		return val != "none" && val != ""
	}
	return false
}

func (c *KeycloakClientConfig) mapDPoPRequired() bool {
	if val, ok := c.Attributes["dpop.required"]; ok {
		return val == "true"
	}
	return false
}

func (c *KeycloakClientConfig) mapPARRequired() bool {
	if val, ok := c.Attributes["par.required"]; ok {
		return val == "true"
	}
	return false
}

func (c *KeycloakClientConfig) mapPromptScopeApproval() bool {
	return c.ConsentRequired
}

func (c *KeycloakClientConfig) mapAccessTokenFormat() string {
	// Keycloak issues JWTs by default
	return TokenFormatJwt
}

func (c *KeycloakClientConfig) mapAccessTokenLifetime() int {
	// Check attributes for access token lifespan
	if val, ok := strconv.Atoi(c.Attributes["access.token.lifespan"]); ok == nil {
		return val
	}
	return TokenLifetimeDefault
}

func (c *KeycloakClientConfig) mapRotateRefreshTokens() bool {
	// Keycloak rotates refresh tokens by default when use_refresh_tokens is enabled
	return c.UseRefreshTokens
}

func (c *KeycloakClientConfig) mapSSO() bool {
	// Keycloak doesn't have a direct SSO flag; check if standard flow is enabled
	return c.StandardFlowEnabled
}

func (c *KeycloakClientConfig) mapEnforce2SV() bool {
	// Check attributes for 2FA requirement
	if val, ok := c.Attributes["enforce2fa"]; ok {
		return val == "true"
	}
	return false
}

func (c *KeycloakClientConfig) mapPromptTermsAndConditions() bool {
	// Check attributes for terms and conditions requirement
	if val, ok := c.Attributes["require-tnc"]; ok {
		return val == "true"
	}
	return false
}

func (c *KeycloakClientConfig) mapClientSecret() string {
	return c.Secret
}
