package keycloak

import (
	"encoding/json"
	"ocm/export/authserver"
	"ocm/export/oidcconfig"
	"strconv"
	"strings"
)

// KeycloakRefreshRollingDefault is Keycloak's default refresh-token rotation behavior.
const KeycloakRefreshRollingDefault = true

func CreateKeycloakClientConfig() *KeycloakClientConfig {
	return &KeycloakClientConfig{}
}

type KeycloakClientConfig struct {
	// Basic fields
	ID                      string `json:"id"`
	ClientID                string `json:"clientId"`
	Name                    string `json:"name"`
	Description             string `json:"description"`
	Enabled                 bool   `json:"enabled"`
	Protocol                string `json:"protocol"`
	PublicClient            bool   `json:"publicClient"`
	ClientAuthenticatorType string `json:"clientAuthenticatorType"`
	Secret                  string `json:"secret"`
	RootURL                 string `json:"rootUrl"`
	BaseURL                 string `json:"baseUrl"`

	// OAuth/OIDC flows and URIs
	RedirectUris              []string `json:"redirectUris"`
	WebOrigins                []string `json:"webOrigins"`
	StandardFlowEnabled       bool     `json:"standardFlowEnabled"`
	ImplicitFlowEnabled       bool     `json:"implicitFlowEnabled"`
	DirectAccessGrantsEnabled bool     `json:"directAccessGrantsEnabled"`
	ServiceAccountsEnabled    bool     `json:"serviceAccountsEnabled"`

	// OIDC specific
	ConsentRequired    bool `json:"consentRequired"`
	FrontchannelLogout bool `json:"frontchannelLogout"`

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
			ClientID:                    c.mapClientID(),
			Name:                        c.mapName(),
			Description:                 c.mapDescription(),
			Contacts:                    c.mapContacts(),
			ClientURI:                   c.mapClientURI(),
			LogoURI:                     c.mapLogoURI(),
			TosURI:                      c.mapTosURI(),
			PolicyURI:                   c.mapPolicyURI(),
			JWKSURI:                     c.mapJWKSURI(),
			RedirectURIs:                c.mapRedirectURIs(),
			ResponseTypes:               c.mapResponseTypes(),
			GrantTypes:                  c.mapGrantTypes(),
			TokenEndpointAuthMethod:     c.mapTokenEndpointAuthMethod(),
			TokenEndpointAuthSigningAlg: c.mapTokenEndpointAuthSigningAlg(),
			IdTokenSignedResponseAlg:    c.mapIdTokenSignedResponseAlg(),
			RequestObjectSigningAlg:     c.mapRequestObjectSigningAlg(),
			Scopes:                      c.mapScopes(),
			ConsentRequired:             c.mapConsentRequired(),
			ApplicationType:             c.mapApplicationType(),
			SubjectType:                 c.mapSubjectType(),
			SectorIdentifierURI:         c.mapSectorIdentifierURI(),
			BackchannelLogoutURI:        c.mapBackchannelLogoutURI(),
			FrontchannelLogoutURI:       c.mapFrontchannelLogoutURI(),
			PostLogoutRedirectURIs:      c.mapPostLogoutRedirectURIs(),
			DefaultACRValues:            c.mapDefaultACRValues(),
			InitiateLoginURI:            c.mapInitiateLoginURI(),
			RequestURIs:                 c.mapRequestURIs(),
		},
		Extensions: oidcconfig.CanonicalClientConfigExtensions{
			Enabled:                           c.mapEnabled(),
			PKCERequired:                      c.mapPKCERequired(),
			DPoPRequired:                      c.mapDPoPRequired(),
			PARRequired:                       c.mapPARRequired(),
			AccessTokenFormat:                 c.mapAccessTokenFormat(),
			AccessTokenLifetimeSeconds:        c.mapAccessTokenLifetimeSeconds(),
			RefreshTokenLifetimeSeconds:       c.mapRefreshTokenLifetimeSeconds(),
			RefreshTokenIdleTimeoutSeconds:    c.mapRefreshTokenIdleTimeoutSeconds(),
			RotateRefreshTokens:               c.mapRotateRefreshTokens(),
			MinimumACRValue:                   c.mapMinimumACRValue(),
			RequireTermsAndConditionsApproval: c.mapRequireTermsAndConditionsApproval(),
		},
		Secrets: oidcconfig.CanonicalClientConfigSecrets{
			PlainSecret: c.mapPlainSecret(),
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

func (c *KeycloakClientConfig) mapContacts() []string {
	// Keycloak does not store contacts, returning nil
	return nil
}

func (c *KeycloakClientConfig) mapClientURI() string {
	return c.BaseURL
}

func (c *KeycloakClientConfig) mapLogoURI() string {
	if val, ok := c.Attributes["logoUri"]; ok && val != "" {
		return val
	}
	return ""
}

func (c *KeycloakClientConfig) mapTosURI() string {
	if val, ok := c.Attributes["tosUri"]; ok && val != "" {
		return val
	}
	return ""
}

func (c *KeycloakClientConfig) mapPolicyURI() string {
	if val, ok := c.Attributes["policyUri"]; ok && val != "" {
		return val
	}
	return ""
}

func (c *KeycloakClientConfig) mapJWKSURI() string {
	if c.Attributes["use.jwks.url"] == "true" {
		return c.Attributes["jwks.url"]
	}
	return ""
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
		g = append(g, authserver.GrantTypeAuthorizationCode)
	}
	if c.ImplicitFlowEnabled {
		g = append(g, authserver.GrantTypeImplicit)
	}
	if c.Attributes["use.refresh.tokens"] == "true" {
		g = append(g, authserver.GrantTypeRefreshToken)
	}
	if c.ServiceAccountsEnabled {
		g = append(g, authserver.GrantTypeClientCredentials)
	}
	if c.DirectAccessGrantsEnabled {
		g = append(g, authserver.GrantTypeResourceOwnerPassword)
	}

	// Check attributes for additional grant types
	if c.Attributes["oauth2.device.authorization.grant.enabled"] == "true" {
		g = append(g, authserver.GrantTypeDeviceCode)
	}

	return g
}

func (c *KeycloakClientConfig) mapTokenEndpointAuthSigningAlg() string {
	if val, ok := c.Attributes["token.endpoint.auth.signing.alg"]; ok {
		return val
	}
	return ""
}

func (c *KeycloakClientConfig) mapIdTokenSignedResponseAlg() string {
	if val, ok := c.Attributes["id.token.signed.response.alg"]; ok {
		return val
	}
	return ""
}

func (c *KeycloakClientConfig) mapRequestObjectSigningAlg() string {
	if val, ok := c.Attributes["request.object.signature.alg"]; ok {
		return val
	}
	return ""
}

func (c *KeycloakClientConfig) mapTokenEndpointAuthMethod() string {
	// Map client authenticator type to standard method
	switch c.ClientAuthenticatorType {
	case "client-secret":
		return authserver.AuthMethodClientSecretBasic
	case "client-secret-jwt":
		return authserver.AuthMethodClientSecretJwt
	case "client-jwt":
		return authserver.AuthMethodPrivateKeyJwt
	default:
		if c.PublicClient {
			return authserver.AuthMethodNone
		}
		return authserver.AuthMethodClientSecretBasic
	}
}

func (c *KeycloakClientConfig) mapScopes() []string {
	scopes := append([]string{}, c.DefaultClientScopes...)
	scopes = append(scopes, c.OptionalClientScopes...)
	return scopes
}

func (c *KeycloakClientConfig) mapConsentRequired() bool {
	return c.ConsentRequired
}

func (c *KeycloakClientConfig) mapApplicationType() string {
	if c.PublicClient {
		return authserver.ApplicationTypeNative
	}
	return authserver.ApplicationTypeWeb
}

func (c *KeycloakClientConfig) mapSubjectType() string {
	// Not available in Keycloak client config, defaulting to "public"
	return authserver.SubjectTypePublic
}

func (c *KeycloakClientConfig) mapSectorIdentifierURI() string {
	// Not available in Keycloak client config, returning empty string
	return ""
}

func (c *KeycloakClientConfig) mapBackchannelLogoutURI() string {
	if val, ok := c.Attributes["backchannel.logout.url"]; ok && val != "" {
		return strings.Split(val, ",")[0]
	}
	return ""
}

func (c *KeycloakClientConfig) mapFrontchannelLogoutURI() string {
	if c.FrontchannelLogout {
		if val, ok := c.Attributes["frontchannel.logout.url"]; ok && val != "" {
			return strings.Split(val, ",")[0]
		}
	}
	return ""
}

func (c *KeycloakClientConfig) mapPostLogoutRedirectURIs() []string {
	if val, ok := c.Attributes["post.logout.redirect.uris"]; ok && val != "" {
		return strings.Split(val, ",")
	}
	return nil
}

func (c *KeycloakClientConfig) mapDefaultACRValues() []string {
	if val, ok := c.Attributes["default.acr.values"]; ok && val != "" {
		return strings.Split(val, ",")
	}
	return nil
}

func (c *KeycloakClientConfig) mapInitiateLoginURI() string {
	// Not available in Keycloak client config, returning empty string
	return ""
}

func (c *KeycloakClientConfig) mapRequestURIs() []string {
	if val, ok := c.Attributes["request.uris"]; ok && val != "" {
		return strings.Split(val, ",")
	}
	return nil
}

func (c *KeycloakClientConfig) mapEnabled() bool {
	return c.Enabled
}

func (c *KeycloakClientConfig) mapPKCERequired() bool {
	if val, ok := c.Attributes["pkce.code.challenge.method"]; ok {
		return val != "none" && val != ""
	}
	return false
}

func (c *KeycloakClientConfig) mapDPoPRequired() bool {
	if val, ok := c.Attributes["dpop.bound.access.tokens"]; ok {
		return val == "true"
	}
	return false
}

func (c *KeycloakClientConfig) mapPARRequired() bool {
	if val, ok := c.Attributes["require.pushed.authorization.requests"]; ok {
		return val == "true"
	}
	return false
}

func (c *KeycloakClientConfig) mapAccessTokenFormat() string {
	// Keycloak issues JWTs by default
	return authserver.AccessTokenFormatJwt
}

func (c *KeycloakClientConfig) mapAccessTokenLifetimeSeconds() int {
	// Check attributes for access token lifespan
	if val, ok := strconv.Atoi(c.Attributes["access.token.lifespan"]); ok == nil {
		return val
	}
	return authserver.AccessTokenLifetimeDefault
}

func (c *KeycloakClientConfig) mapRefreshTokenLifetimeSeconds() int {
	// Not available in Keycloak client config, returning default
	return authserver.RefreshTokenLifetimeDefault
}

func (c *KeycloakClientConfig) mapRefreshTokenIdleTimeoutSeconds() int {
	// Not available in Keycloak client config, returning default
	return authserver.RefreshTokenIdleTimeoutDefault
}

func (c *KeycloakClientConfig) mapRotateRefreshTokens() bool {
	// Not available in Keycloak client config, returning default
	if val, ok := c.Attributes["use.refresh.tokens"]; ok && val == "true" {
		return KeycloakRefreshRollingDefault
	}
	return false
}

func (c *KeycloakClientConfig) mapMinimumACRValue() string {
	if val, ok := c.Attributes["minimum.acr.value"]; ok && val != "" {
		return val
	}
	return ""
}

func (c *KeycloakClientConfig) mapRequireTermsAndConditionsApproval() bool {
	// Check attributes for terms and conditions requirement.
	if val, ok := c.Attributes["require-tnc"]; ok {
		return val == "true"
	}
	return false
}

func (c *KeycloakClientConfig) mapPlainSecret() string {
	return c.Secret
}
