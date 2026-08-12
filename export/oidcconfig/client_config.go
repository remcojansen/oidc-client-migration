package oidcconfig

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type CanonicalClientConfig struct {
	Metadata   CanonicalClientConfigMetadata   `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Client     CanonicalClientConfigClient     `json:"client" yaml:"client"`
	Extensions CanonicalClientConfigExtensions `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	Secrets    CanonicalClientConfigSecrets    `json:"secrets,omitempty" yaml:"secrets,omitempty"`
}

type CanonicalClientConfigMetadata struct {
	// Metadata fields for tracking client configuration source, owner, etc.
	OwnerEmail        string `json:"owner_email,omitempty" yaml:"owner_email,omitempty"`
	OwnerSlackChannel string `json:"owner_slack_channel,omitempty" yaml:"owner_slack_channel,omitempty"`
	OwnerTeamADGroup  string `json:"owner_team_ad_group,omitempty" yaml:"owner_team_ad_group,omitempty"`
}

type CanonicalClientConfigClient struct {
	// OAuth 2.0 and OIDC client metadata fields
	ClientID                    string   `json:"client_id" yaml:"client_id"`
	Name                        string   `json:"client_name" yaml:"client_name"`
	Description                 string   `json:"description,omitempty" yaml:"description,omitempty"`
	Contacts                    []string `json:"contacts,omitempty" yaml:"contacts,omitempty"`
	ClientURI                   string   `json:"client_uri,omitempty" yaml:"client_uri,omitempty"`
	LogoURI                     string   `json:"logo_uri,omitempty" yaml:"logo_uri,omitempty"`
	TosURI                      string   `json:"tos_uri,omitempty" yaml:"tos_uri,omitempty"`
	PolicyURI                   string   `json:"policy_uri,omitempty" yaml:"policy_uri,omitempty"`
	JWKSURI                     string   `json:"jwks_uri,omitempty" yaml:"jwks_uri,omitempty"`
	RedirectURIs                []string `json:"redirect_uris,omitempty" yaml:"redirect_uris,omitempty"`
	ResponseTypes               []string `json:"response_types,omitempty" yaml:"response_types,omitempty"`
	GrantTypes                  []string `json:"grant_types,omitempty" yaml:"grant_types,omitempty"`
	TokenEndpointAuthMethod     string   `json:"token_endpoint_auth_method,omitempty" yaml:"token_endpoint_auth_method,omitempty"`
	TokenEndpointAuthSigningAlg string   `json:"token_endpoint_auth_signing_alg,omitempty" yaml:"token_endpoint_auth_signing_alg,omitempty"`
	IdTokenSignedResponseAlg    string   `json:"id_token_signed_response_alg,omitempty" yaml:"id_token_signed_response_alg,omitempty"`
	RequestObjectSigningAlg     string   `json:"request_object_signing_alg,omitempty" yaml:"request_object_signing_alg,omitempty"`
	Scopes                      []string `json:"scopes,omitempty" yaml:"scopes,omitempty"`
	ConsentRequired             bool     `json:"consent_required" yaml:"consent_required"`
	ApplicationType             string   `json:"application_type,omitempty" yaml:"application_type,omitempty"` // "web" (confidential client) or "native" (public client)
	SubjectType                 string   `json:"subject_type,omitempty" yaml:"subject_type,omitempty"`
	SectorIdentifierURI         string   `json:"sector_identifier_uri,omitempty" yaml:"sector_identifier_uri,omitempty"`
	BackchannelLogoutURI        string   `json:"backchannel_logout_uri,omitempty" yaml:"backchannel_logout_uri,omitempty"`
	FrontchannelLogoutURI       string   `json:"frontchannel_logout_uri,omitempty" yaml:"frontchannel_logout_uri,omitempty"`
	PostLogoutRedirectURIs      []string `json:"post_logout_redirect_uris,omitempty" yaml:"post_logout_redirect_uris,omitempty"`
	DefaultACRValues            []string `json:"default_acr_values,omitempty" yaml:"default_acr_values,omitempty"` // default ACR values to be used if the authorization request does not specify any ACR values
	InitiateLoginURI            string   `json:"initiate_login_uri,omitempty" yaml:"initiate_login_uri,omitempty"`
	RequestURIs                 []string `json:"request_uris,omitempty" yaml:"request_uris,omitempty"`
}

type CanonicalClientConfigExtensions struct {
	// Non-standard extensions for additional client configuration options
	Enabled                          bool   `json:"enabled" yaml:"enabled"` // default: true
	PKCERequired                     bool   `json:"pkce_required" yaml:"pkce_required"`
	DPoPRequired                     bool   `json:"dpop_required" yaml:"dpop_required"`
	PARRequired                      bool   `json:"par_required" yaml:"par_required"`
	AccessTokenFormat                string `json:"access_token_format,omitempty" yaml:"access_token_format,omitempty"` // "jwt" or "opaque"
	AccessTokenLifetimeSeconds       int    `json:"access_token_lifetime_seconds,omitempty" yaml:"access_token_lifetime_seconds,omitempty"`
	OfflineSessionMaxLifetimeSeconds int    `json:"offline_session_max_lifetime_seconds,omitempty" yaml:"offline_session_max_lifetime_seconds,omitempty"` // persistent/offline refresh token lifetime, independent of any browser SSO session
	OfflineSessionIdleTimeoutSeconds int    `json:"offline_session_idle_timeout_seconds,omitempty" yaml:"offline_session_idle_timeout_seconds,omitempty"`
	SessionMaxLifetimeSeconds        int    `json:"session_max_lifetime_seconds,omitempty" yaml:"session_max_lifetime_seconds,omitempty"` // refresh token lifetime tied to the browser SSO session; no PingFederate equivalent
	SessionIdleTimeoutSeconds        int    `json:"session_idle_timeout_seconds,omitempty" yaml:"session_idle_timeout_seconds,omitempty"`
	RotateRefreshTokens              bool   `json:"rotate_refresh_tokens" yaml:"rotate_refresh_tokens"`
	MinimumACRValue                  string `json:"minimum_acr_value,omitempty" yaml:"minimum_acr_value,omitempty"`     // minimum ACR values to be used regardless of the requested ACR values in the authorization request
	TermsAndConditionsRequired       bool   `json:"terms_and_conditions_required" yaml:"terms_and_conditions_required"` // prompt users to accept authorization server's terms and conditions
}

type CanonicalClientConfigSecrets struct {
	// Client secrets or certificates can be provided here if needed
	PlainSecret     string `json:"plain_secret,omitempty" yaml:"plain_secret,omitempty"` // insecure, do not use for production clients
	EncryptedSecret string `json:"encrypted_secret,omitempty" yaml:"encrypted_secret,omitempty"`
}

func (c *CanonicalClientConfig) ToJSON() ([]byte, error) {
	return json.MarshalIndent(c, "", "  ")
}

func (c *CanonicalClientConfig) ToYAML() ([]byte, error) {
	return yaml.Marshal(c)
}

func (c *CanonicalClientConfig) ToJSONString() (string, error) {
	data, err := c.ToJSON()
	if err != nil {
		return "", err
	}
	return string(data) + "\n", nil
}

func (c *CanonicalClientConfig) WriteConfigFile(outputDir, format string) error {
	var err error
	var configData []byte
	var fileSuffix string

	if format == "json" {
		configData, err = c.ToJSON()
	} else {
		configData, err = c.ToYAML()
	}
	if err != nil {
		return fmt.Errorf("error marshalling client %s: %v", c.Client.ClientID, err)
	}

	if format == "json" {
		fileSuffix = ".json"
	} else {
		fileSuffix = ".yaml"
	}

	filePath := outputDir + "/" + c.Client.ClientID + fileSuffix
	if err = os.WriteFile(filePath, append(configData, '\n'), 0644); err != nil {
		return fmt.Errorf("error writing file %s: %v", filePath, err)
	}

	return nil
}
