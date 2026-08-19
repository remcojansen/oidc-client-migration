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
		Annotations: oidcconfig.CanonicalClientConfigAnnotations{
			// Not provided by Keycloak API
		},
		Client: oidcconfig.CanonicalClientConfigClient{
			ClientID:                    c.mapClientID(),
			Name:                        c.mapName(),
			Description:                 c.mapDescription(),
			Enabled:                     c.mapEnabled(),
			ApplicationType:             c.mapApplicationType(),
			TokenEndpointAuthMethod:     c.mapTokenEndpointAuthMethod(),
			GrantTypes:                  c.mapGrantTypes(),
			ResponseTypes:               c.mapResponseTypes(),
			RedirectURIs:                c.mapRedirectURIs(),
			PostLogoutRedirectURIs:      c.mapPostLogoutRedirectURIs(),
			Scopes:                      c.mapScopes(),
			ConsentRequired:             c.mapConsentRequired(),
			PKCERequired:                c.mapPKCERequired(),
			DPoPRequired:                c.mapDPoPRequired(),
			PARRequired:                 c.mapPARRequired(),
			AccessTokenFormat:           c.mapAccessTokenFormat(),
			AccessTokenLifetimeSeconds:  c.mapAccessTokenLifetimeSeconds(),
			RotateRefreshTokens:         c.mapRotateRefreshTokens(),
			OfflineSessionMaxLifetimeSeconds: c.mapOfflineSessionMaxLifetimeSeconds(),
			OfflineSessionIdleTimeoutSeconds: c.mapOfflineSessionIdleTimeoutSeconds(),
			SessionMaxLifetimeSeconds:        c.mapSessionMaxLifetimeSeconds(),
			SessionIdleTimeoutSeconds:        c.mapSessionIdleTimeoutSeconds(),
			IntrospectionEnabled:             c.mapIntrospectionEnabled(),
			MinimumACRValue:                  c.mapMinimumACRValue(),
			DefaultACRValues:                 c.mapDefaultACRValues(),
			SubjectType:                      c.mapSubjectType(),
			SectorIdentifierURI:              c.mapSectorIdentifierURI(),
			IdTokenSignedResponseAlg:         c.mapIdTokenSignedResponseAlg(),
			TokenEndpointAuthSigningAlg:      c.mapTokenEndpointAuthSigningAlg(),
			RequestObjectSigningAlg:          c.mapRequestObjectSigningAlg(),
			JWKSURI:                          c.mapJWKSURI(),
			Contacts:                         c.mapContacts(),
			ClientURI:                        c.mapClientURI(),
			LogoURI:                          c.mapLogoURI(),
			TosURI:                           c.mapTosURI(),
			PolicyURI:                        c.mapPolicyURI(),
			BackchannelLogoutURI:             c.mapBackchannelLogoutURI(),
			FrontchannelLogoutURI:            c.mapFrontchannelLogoutURI(),
			InitiateLoginURI:                 c.mapInitiateLoginURI(),
			RequestURIs:                      c.mapRequestURIs(),
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
	// A public client never authenticates with the token endpoint, regardless of what
	// clientAuthenticatorType Keycloak reports — Keycloak sets a default of "client-secret"
	// on virtually every client record, including public ones, where it's simply ignored.
	if c.PublicClient {
		return authserver.AuthMethodNone
	}

	// Map client authenticator type to standard method
	switch c.ClientAuthenticatorType {
	case "client-secret":
		return authserver.AuthMethodClientSecretBasic
	case "client-secret-jwt":
		return authserver.AuthMethodClientSecretJwt
	case "client-jwt":
		return authserver.AuthMethodPrivateKeyJwt
	default:
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

func (c *KeycloakClientConfig) mapOfflineSessionMaxLifetimeSeconds() int {
	// Keycloak's "offline session max lifespan" client attribute governs the lifetime of
	// persistent/offline refresh tokens, which are only issued to clients that request the
	// offline_access scope (see import/modules/oauth-client-keycloak, which adds that scope
	// automatically whenever this value is set).
	if val, ok := strconv.Atoi(c.Attributes["client.offline.session.max.lifespan"]); ok == nil {
		return val
	}
	return authserver.OfflineSessionMaxLifetimeDefault
}

func (c *KeycloakClientConfig) mapOfflineSessionIdleTimeoutSeconds() int {
	if val, ok := strconv.Atoi(c.Attributes["client.offline.session.idle.timeout"]); ok == nil {
		return val
	}
	return authserver.OfflineSessionIdleTimeoutDefault
}

func (c *KeycloakClientConfig) mapSessionMaxLifetimeSeconds() int {
	// Keycloak's "session max lifespan" client attribute governs the lifetime of refresh
	// tokens tied to the browser SSO session (i.e. clients that do not request offline_access).
	// There is no PingFederate equivalent for this concept.
	if val, ok := strconv.Atoi(c.Attributes["client.session.max.lifespan"]); ok == nil {
		return val
	}
	return authserver.SessionMaxLifetimeDefault
}

func (c *KeycloakClientConfig) mapSessionIdleTimeoutSeconds() int {
	if val, ok := strconv.Atoi(c.Attributes["client.session.idle.timeout"]); ok == nil {
		return val
	}
	return authserver.SessionIdleTimeoutDefault
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

func (c *KeycloakClientConfig) mapIntrospectionEnabled() bool {
	// Keycloak has no per-client toggle for introspection: any client that can authenticate
	// itself may call the introspection endpoint, so this is always true.
	return true
}

func (c *KeycloakClientConfig) mapPlainSecret() string {
	return c.Secret
}
