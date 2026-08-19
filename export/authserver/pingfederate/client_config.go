package pingfederate

import (
	"encoding/json"
	"fmt"
	"ocm/export/authserver"
	"ocm/export/oidcconfig"
)

// PingFederateRefreshRollingDefault is PingFederate's server-default refresh-token rotation setting.
const PingFederateRefreshRollingDefault = "ROLL"

func contains[T comparable](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func CreatePingFederateClientConfig() *PingFederateClientConfig {
	return &PingFederateClientConfig{}
}

type PingFederateClientConfig struct {
	ClientID                           string             `json:"clientId"`
	Enabled                            bool               `json:"enabled"`
	Name                               string             `json:"name"`
	Description                        string             `json:"description"`
	LogoURL                            string             `json:"logoUrl"`
	RedirectUris                       []string           `json:"redirectUris"`
	GrantTypes                         []string           `json:"grantTypes"`
	BypassApprovalPage                 bool               `json:"bypassApprovalPage"`
	RequireProofKeyForCodeExchange     bool               `json:"requireProofKeyForCodeExchange"`
	RequireDPoP                        bool               `json:"requireDpop"`
	RequirePushedAuthorizationRequests bool               `json:"requirePushedAuthorizationRequests"`
	RestrictScopes                     bool               `json:"restrictScopes"`
	RestrictedScopes                   []string           `json:"restrictedScopes"`
	ExclusiveScopes                    []string           `json:"exclusiveScopes"`
	ClientAuth                         AuthConfig         `json:"clientAuth"`
	JWKSSettings                       JWKSSettings       `json:"jwksSettings"`
	DefaultAccessTokenManagerRef       AccessTokenManager `json:"defaultAccessTokenManagerRef"`
	OidcPolicy                         OIDCPolicy         `json:"oidcPolicy"`
	ExtendedParameters                 map[string]ExtendedParameterValue `json:"extendedParameters"`
	RefreshRolling                     string             `json:"refreshRolling"`
	PersistentGrantIdleTimeoutType     string             `json:"persistentGrantIdleTimeoutType"`
	PersistentGrantIdleTimeout         int                `json:"persistentGrantIdleTimeout"`
	PersistentGrantIdleTimeoutTimeUnit string             `json:"persistentGrantIdleTimeoutTimeUnit"`
	PersistentGrantExpirationType      string             `json:"persistentGrantExpirationType"`
	PersistentGrantExpirationTime      int                `json:"persistentGrantExpirationTime"`
	PersistentGrantExpirationTimeUnit  string             `json:"persistentGrantExpirationTimeUnit"`
	// Other fields are omitted as they are unused in our mapping

	// accessTokenManagerMapping maps DefaultAccessTokenManagerRef.Id to its format/lifetime; it is
	// injected by the client (see PingFederateClient.WithAccessTokenManagerMapping) rather than
	// read from the PingFederate API response.
	accessTokenManagerMapping map[string]AccessTokenManagerInfo

	minimumAcrValueParamName string
}

// AccessTokenManagerInfo describes the access token format and lifetime produced by a specific
// PingFederate access token manager, as configured by the caller via
// PingFederateClient.WithAccessTokenManagerMapping.
type AccessTokenManagerInfo struct {
	Format          string `json:"format"`
	LifetimeSeconds int    `json:"lifetime_seconds"`
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

type JWKSSettings struct {
	JWKSURL string `json:"jwksUrl"`
	// Other fields are omitted as they are unused in our mapping
}

type AccessTokenManager struct {
	Id string `json:"id"`
	// Other fields are omitted as they are unused in our mapping
}

type OIDCPolicy struct {
	PolicyGroup                PolicyGroup `json:"policyGroup"`
	IdTokenSigningAlgorithm    string      `json:"idTokenSigningAlgorithm"`
	LogoutMode                 string      `json:"logoutMode"`
	LogoutURIs                 []string    `json:"logoutUris"`
	BackChannelLogoutURI       string      `json:"backChannelLogoutUri"`
	PostLogoutRedirectURIs     []string    `json:"postLogoutRedirectURIs"`
	PairwiseIdentifierUserType bool        `json:"pairwiseIdentifierUserType"`
	// Other fields are omitted as they are unused in our mapping
}

type PolicyGroup struct {
	Id string `json:"id"`
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
			Enabled:                          c.mapEnabled(),
			PKCERequired:                     c.mapPKCERequired(),
			DPoPRequired:                     c.mapDPoPRequired(),
			PARRequired:                      c.mapPARRequired(),
			AccessTokenFormat:                c.mapAccessTokenFormat(),
			AccessTokenLifetimeSeconds:       c.mapAccessTokenLifetimeSeconds(),
			OfflineSessionMaxLifetimeSeconds: c.mapOfflineSessionMaxLifetimeSeconds(),
			OfflineSessionIdleTimeoutSeconds: c.mapOfflineSessionIdleTimeoutSeconds(),
			RotateRefreshTokens:              c.mapRotateRefreshTokens(),
			MinimumACRValue:                  c.mapMinimumACRValue(),
			IntrospectionEnabled:             c.mapIntrospectionEnabled(),
		},
		Secrets: oidcconfig.CanonicalClientConfigSecrets{
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

func (c *PingFederateClientConfig) mapContacts() []string {
	// Not available in PingFederate client configuration, so return nil.
	return nil
}

func (c *PingFederateClientConfig) mapClientURI() string {
	// Not available in PingFederate client configuration, so return empty string.
	return ""
}

func (c *PingFederateClientConfig) mapLogoURI() string {
	return c.LogoURL
}

func (c *PingFederateClientConfig) mapTosURI() string {
	// Not available in PingFederate client configuration, so return empty string.
	return ""
}

func (c *PingFederateClientConfig) mapPolicyURI() string {
	// Not available in PingFederate client configuration, so return empty string.
	return ""
}

func (c *PingFederateClientConfig) mapJWKSURI() string {
	return c.JWKSSettings.JWKSURL
}

func (c *PingFederateClientConfig) mapRedirectURIs() []string {
	return c.RedirectUris
}

func (cc *PingFederateClientConfig) mapResponseTypes() []string {
	var responseTypes []string
	for _, grantType := range cc.GrantTypes {
		switch grantType {
		case "AUTHORIZATION_CODE":
			responseTypes = append(responseTypes, "code")
		case "IMPLICIT":
			responseTypes = append(responseTypes, "token", "id_token")
		}
	}
	return responseTypes
}

func (c *PingFederateClientConfig) mapGrantTypes() []string {
	var g []string
	for _, gt := range c.GrantTypes {
		switch gt {
		case "AUTHORIZATION_CODE":
			g = append(g, authserver.GrantTypeAuthorizationCode)
		case "IMPLICIT":
			g = append(g, authserver.GrantTypeImplicit)
		case "REFRESH_TOKEN":
			g = append(g, authserver.GrantTypeRefreshToken)
		case "CLIENT_CREDENTIALS":
			g = append(g, authserver.GrantTypeClientCredentials)
		case "RESOURCE_OWNER_PASSWORD_CREDENTIALS":
			g = append(g, authserver.GrantTypeResourceOwnerPassword)
		case "DEVICE_CODE":
			g = append(g, authserver.GrantTypeDeviceCode)
		case "TOKEN_EXCHANGE":
			g = append(g, authserver.GrantTypeTokenExchange)
		}
	}
	return g
}

// mapIntrospectionEnabled maps PingFederate's ACCESS_TOKEN_VALIDATION grant type. It is not a
// real OAuth grant type, so it is surfaced as extensions.introspection_enabled instead of being
// included in GrantTypes.
func (c *PingFederateClientConfig) mapIntrospectionEnabled() bool {
	return contains(c.GrantTypes, "ACCESS_TOKEN_VALIDATION")
}

func (cc *PingFederateClientConfig) mapTokenEndpointAuthMethod() string {
	switch cc.ClientAuth.Type {
	case "NONE":
		return authserver.AuthMethodNone
	case "SECRET":
		return authserver.AuthMethodClientSecretBasic
	case "CLIENT_SECRET_JWT":
		return authserver.AuthMethodClientSecretJwt
	case "PRIVATE_KEY_JWT":
		return authserver.AuthMethodPrivateKeyJwt
	case "CLIENT_SECRET_BASIC":
		return authserver.AuthMethodClientSecretBasic
	default:
		return authserver.AuthMethodClientSecretBasic
	}
}

func (c *PingFederateClientConfig) mapTokenEndpointAuthSigningAlg() string {
	// Not available in PingFederate client configuration, so return empty string.
	return ""
}

func (c *PingFederateClientConfig) mapIdTokenSignedResponseAlg() string {
	return c.OidcPolicy.IdTokenSigningAlgorithm
}

func (c *PingFederateClientConfig) mapRequestObjectSigningAlg() string {
	// Not available in PingFederate client configuration, so return empty string.
	return ""
}

func (c *PingFederateClientConfig) mapScopes() []string {
	scopes := append([]string{}, c.RestrictedScopes...)
	scopes = append(scopes, c.ExclusiveScopes...)
	return scopes
}

func (c *PingFederateClientConfig) mapConsentRequired() bool {
	return !c.BypassApprovalPage
}

func (c *PingFederateClientConfig) mapApplicationType() string {
	if c.ClientAuth.Type == "NONE" {
		return authserver.ApplicationTypeNative
	}
	return authserver.ApplicationTypeWeb
}

func (c *PingFederateClientConfig) mapSubjectType() string {
	if c.OidcPolicy.PairwiseIdentifierUserType == true {
		return authserver.SubjectTypePairwise
	}
	return authserver.SubjectTypePublic
}

func (c *PingFederateClientConfig) mapSectorIdentifierURI() string {
	// Not available in PingFederate client configuration, so return empty string.
	return ""
}

func (c *PingFederateClientConfig) mapBackchannelLogoutURI() string {
	if c.OidcPolicy.LogoutMode == "BACK_CHANNEL" && c.OidcPolicy.BackChannelLogoutURI != "" {
		return c.OidcPolicy.BackChannelLogoutURI
	}
	return ""
}

func (c *PingFederateClientConfig) mapFrontchannelLogoutURI() string {
	if c.OidcPolicy.LogoutMode == "FRONT_CHANNEL" && len(c.OidcPolicy.LogoutURIs) > 0 {
		return c.OidcPolicy.LogoutURIs[0]
	}
	return ""
}

func (c *PingFederateClientConfig) mapPostLogoutRedirectURIs() []string {
	return c.OidcPolicy.PostLogoutRedirectURIs
}

func (c *PingFederateClientConfig) mapDefaultACRValues() []string {
	// Not available in PingFederate client configuration, so return nil.
	return nil
}

func (c *PingFederateClientConfig) mapInitiateLoginURI() string {
	// Not available in PingFederate client configuration, so return empty string.
	return ""
}

func (c *PingFederateClientConfig) mapRequestURIs() []string {
	// Not available in PingFederate client configuration, so return nil.
	return nil
}

func (c *PingFederateClientConfig) mapEnabled() bool {
	return c.Enabled
}

func (c *PingFederateClientConfig) mapPKCERequired() bool {
	return c.RequireProofKeyForCodeExchange
}

func (c *PingFederateClientConfig) mapDPoPRequired() bool {
	return c.RequireDPoP
}

func (c *PingFederateClientConfig) mapPARRequired() bool {
	return c.RequirePushedAuthorizationRequests
}

func (c *PingFederateClientConfig) mapRotateRefreshTokens() bool {
	val := c.RefreshRolling
	if val == "SERVER_DEFAULT" {
		val = PingFederateRefreshRollingDefault
	}
	return val == "ROLL"
}

func (c *PingFederateClientConfig) mapAccessTokenFormat() string {
	if info, ok := c.accessTokenManagerMapping[c.DefaultAccessTokenManagerRef.Id]; ok {
		return info.Format
	}
	return ""
}

func (c *PingFederateClientConfig) mapAccessTokenLifetimeSeconds() int {
	if info, ok := c.accessTokenManagerMapping[c.DefaultAccessTokenManagerRef.Id]; ok {
		return info.LifetimeSeconds
	}
	return authserver.AccessTokenLifetimeDefault
}

func (c *PingFederateClientConfig) mapOfflineSessionMaxLifetimeSeconds() int {
	if c.PersistentGrantExpirationType == "OVERRIDE_SERVER_DEFAULT" {
		timeSec, err := calculateSeconds(c.PersistentGrantExpirationTime, c.PersistentGrantExpirationTimeUnit)
		if err != nil {
			return authserver.OfflineSessionMaxLifetimeDefault
		}
		return timeSec
	}
	return authserver.OfflineSessionMaxLifetimeDefault
}

func (c *PingFederateClientConfig) mapOfflineSessionIdleTimeoutSeconds() int {
	if c.PersistentGrantIdleTimeoutType == "OVERRIDE_SERVER_DEFAULT" {
		timeSec, err := calculateSeconds(c.PersistentGrantIdleTimeout, c.PersistentGrantIdleTimeoutTimeUnit)
		if err != nil {
			return authserver.OfflineSessionIdleTimeoutDefault
		}
		return timeSec
	}
	return authserver.OfflineSessionIdleTimeoutDefault
}

func (c *PingFederateClientConfig) mapMinimumACRValue() string {
	param, ok := c.ExtendedParameters[c.minimumAcrValueParamName]
	if !ok || len(param.Value) == 0 {
		return ""
	}
	return param.Value[0]
}

func (c *PingFederateClientConfig) mapEncryptedSecret() string {
	return c.ClientAuth.EncryptedSecret
}

func calculateSeconds(value int, unit string) (int, error) {
	switch unit {
	case "SECONDS":
		return value, nil
	case "MINUTES":
		return value * 60, nil
	case "HOURS":
		return value * 3600, nil
	case "DAYS":
		return value * 86400, nil
	default:
		return 0, fmt.Errorf("unknown unit: %s", unit)
	}
}
