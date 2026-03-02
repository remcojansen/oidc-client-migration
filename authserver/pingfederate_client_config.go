package authserver

import (
	"encoding/json"
	"ocm/oidcconfig"
	"strings"
)

func CreatePingFederateClientConfig() *PingFederateClientConfig {
	return &PingFederateClientConfig{}
}

type PingFederateClientConfig struct {
	ClientID                       string             `json:"clientId"`
	ClientAuth                     AuthConfig         `json:"clientAuth"`
	Enabled                        bool               `json:"enabled"`
	Name                           string             `json:"name"`
	Description                    string             `json:"description"`
	RedirectUris                   []string           `json:"redirectUris"`
	GrantTypes                     []string           `json:"grantTypes"`
	BypassApprovalPage             bool               `json:"bypassApprovalPage"`
	RequireProofKeyForCodeExchange bool               `json:"requireProofKeyForCodeExchange"`
	RequireDPoP                    bool               `json:"requireDpop"`
	RequirePushedAuthorization     bool               `json:"requirePushedAuthorizationRequests"`
	RestrictScopes                 bool               `json:"restrictScopes"`
	RestrictedScopes               []string           `json:"restrictedScopes"`
	ExclusiveScopes                []string           `json:"exclusiveScopes"`
	DefaultAccessTokenManagerRef   AccessTokenManager `json:"defaultAccessTokenManagerRef"`
	OidcPolicy                     OIDCPolicy         `json:"oidcPolicy"`
	ExtendedParameters             ExtendedParameters `json:"extendedParameters"`
	RefreshRolling                 string             `json:"refreshRolling"`
	// Other fields are omitted as they are unused in our mapping
}

type PingFederateClientList struct {
	Items []PingFederateClientConfig `json:"items"`
}

type AuthConfig struct {
	Type            string `json:"type"`
	Secret          string `json:"secret"`
	EncryptedSecret string `json:"encryptedSecret"`
	// Other fields are omitted as they are unused in our mapping
}

type AccessTokenManager struct {
	Id string `json:"id"`
	// Other fields are omitted as they are unused in our mapping
}

type OIDCPolicy struct {
	PolicyGroup            PolicyGroup `json:"policyGroup"`
	LogoutMode             string      `json:"logoutMode"`
	LogoutURIs             []string    `json:"logoutUris"`
	BackChannelLogoutURI   string      `json:"backChannelLogoutUri"`
	PostLogoutRedirectURIs []string    `json:"postLogoutRedirectURIs"`
	// Other fields are omitted as they are unused in our mapping
}

type PolicyGroup struct {
	Id string `json:"id"`
	// Other fields are omitted as they are unused in our mapping
}
type ExtendedParameters struct {
	ExcludeTnC  ExtendedParameterValue `json:"exclude_tnc"`
	Enforce2SV  ExtendedParameterValue `json:"enforce_2sv"`
	AdapterType ExtendedParameterValue `json:"adapter_type"`
	// Other fields are omitted as they are unused in our mapping
}

type ExtendedParameterValue struct {
	Value []string `json:"values"`
}

func (cc *PingFederateClientConfig) ToJSON() ([]byte, error) {
	jsonData, err := json.MarshalIndent(cc, "", "  ")
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func (c *PingFederateClientConfig) GetCanonicalClientConfig() *oidcconfig.CanonicalClientConfig {
	return &oidcconfig.CanonicalClientConfig{
		Metadata: oidcconfig.CanonicalClientConfigMetadata{
			// Not provided by PingFederate API
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
			PlainSecret:     c.mapPlainSecret(),
			EncryptedSecret: c.mapEncryptedSecret(),
		},
	}
}

func (c *PingFederateClientConfig) GetClientID() string {
	return c.ClientID
}

func (c *PingFederateClientConfig) mapClientID() string {
	return c.ClientID
}

func (c *PingFederateClientConfig) mapName() string {
	return c.Name
}

func (c *PingFederateClientConfig) mapDescription() string {
	return c.Description
}

func (c *PingFederateClientConfig) mapRedirectURIs() []string {
	return c.RedirectUris
}

func (cc *PingFederateClientConfig) mapResponseTypes() []string {
	// Not provided by PingFederate API
	return nil
}

func (c *PingFederateClientConfig) mapGrantTypes() []string {
	var g []string
	for _, gt := range c.GrantTypes {
		switch gt {
		case "AUTHORIZATION_CODE":
			g = append(g, "authorization_code")
		case "IMPLICIT":
			g = append(g, "implicit")
		case "REFRESH_TOKEN":
			g = append(g, "refresh_token")
		case "CLIENT_CREDENTIALS":
			g = append(g, "client_credentials")
		case "RESOURCE_OWNER_PASSWORD_CREDENTIALS":
			g = append(g, "password")
		case "DEVICE_CODE":
			g = append(g, "device_code")
		case "TOKEN_EXCHANGE":
			g = append(g, "token_exchange")
		case "ACCESS_TOKEN_VALIDATION":
			g = append(g, "introspect")
		}
	}
	return g
}

func (cc *PingFederateClientConfig) mapTokenEndpointAuthMethod() string {
	switch cc.ClientAuth.Type {
	case "NONE":
		return AuthMethodNone
	case "CLIENT_SECRET_JWT":
		return AuthMethodClientSecretJwt
	case "PRIVATE_KEY_JWT":
		return AuthMethodPrivateKeyJwt
	case "CLIENT_SECRET_BASIC":
		return AuthMethodClientSecretBasic
	default:
		return AuthMethodClientSecretBasic
	}
}

func (c *PingFederateClientConfig) mapScopes() string {
	var scopes []string
	scopes = append(scopes, c.RestrictedScopes...)
	scopes = append(scopes, c.ExclusiveScopes...)
	return strings.Join(scopes, ",")
}

func (cc *PingFederateClientConfig) mapBackchannelLogoutURIs() []string {
	if cc.OidcPolicy.LogoutMode == "BACK_CHANNEL" {
		return []string{cc.OidcPolicy.BackChannelLogoutURI}
	}
	return nil
}

func (cc *PingFederateClientConfig) mapFrontchannelLogoutURIs() []string {
	if cc.OidcPolicy.LogoutMode == "FRONT_CHANNEL" {
		return cc.OidcPolicy.LogoutURIs
	}
	return nil
}

func (cc *PingFederateClientConfig) mapPostLogoutRedirectURIs() []string {
	return cc.OidcPolicy.PostLogoutRedirectURIs
}

func (c *PingFederateClientConfig) mapClientType() string {
	if c.ClientAuth.Type == "NONE" {
		return ClientTypePublic
	}
	return ClientTypeConfidential
}

func (c *PingFederateClientConfig) mapPKCERequired() bool {
	return c.RequireProofKeyForCodeExchange
}

func (c *PingFederateClientConfig) mapDPoPRequired() bool {
	return c.RequireDPoP
}

func (c *PingFederateClientConfig) mapPARRequired() bool {
	return c.RequirePushedAuthorization
}

func (c *PingFederateClientConfig) mapPromptScopeApproval() bool {
	return !c.BypassApprovalPage
}

func (c *PingFederateClientConfig) mapAccessTokenFormat() string {
	if strings.HasPrefix(c.DefaultAccessTokenManagerRef.Id, "jwt") {
		return TokenFormatJwt
	}
	return TokenFormatOpaque
}

func (c *PingFederateClientConfig) mapAccessTokenLifetime() int {
	if strings.HasSuffix(c.DefaultAccessTokenManagerRef.Id, "long") {
		return TokenLifetimeLong
	}
	return TokenLifetimeShort
}

func (cc *PingFederateClientConfig) mapRotateRefreshTokens() bool {
	// When refreshRolling is "SERVER_DEFAULT", PingFederate will rotate refresh tokens based on the server's default behavior, which is typically to rotate them. When it's "ROLL", it will always rotate refresh tokens. In both cases, we can consider that refresh token rotation is enabled.
	if cc.RefreshRolling == "SERVER_DEFAULT" || cc.RefreshRolling == "ROLL" {
		return true
	}
	return false
}

func (c *PingFederateClientConfig) mapSSO() bool {
	return contains(c.ExtendedParameters.AdapterType.Value, "sso")
}

func (c *PingFederateClientConfig) mapEnforce2SV() bool {
	return contains(c.ExtendedParameters.Enforce2SV.Value, "true")
}

func (c *PingFederateClientConfig) mapPromptTermsAndConditions() bool {
	return !contains(c.ExtendedParameters.ExcludeTnC.Value, "true")
}

func (c *PingFederateClientConfig) mapPlainSecret() string {
	return c.ClientAuth.Secret
}

func (c *PingFederateClientConfig) mapEncryptedSecret() string {
	return c.ClientAuth.EncryptedSecret
}
