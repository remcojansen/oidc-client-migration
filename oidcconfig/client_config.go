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
}

type CanonicalClientConfigClient struct {
	// Standard RFC 7591 fields
	ClientID                string   `json:"client_id" yaml:"client_id"`
	Name                    string   `json:"client_name" yaml:"client_name"`
	Description             string   `json:"description,omitempty" yaml:"description,omitempty"`
	RedirectURIs            []string `json:"redirect_uris,omitempty" yaml:"redirect_uris,omitempty"`
	ResponseTypes           []string `json:"response_types,omitempty" yaml:"response_types,omitempty"`
	GrantTypes              []string `json:"grant_types,omitempty" yaml:"grant_types,omitempty"`
	TokenEndpointAuthMethod string   `json:"token_endpoint_auth_method,omitempty" yaml:"token_endpoint_auth_method,omitempty"`
	Scope                   string   `json:"scope,omitempty" yaml:"scope,omitempty"`

	// OIDC specific fields
	BackchannelLogoutURIs  []string `json:"backchannel_logout_uri,omitempty" yaml:"backchannel_logout_uri,omitempty"`
	FrontchannelLogoutURIs []string `json:"frontchannel_logout_uri,omitempty" yaml:"frontchannel_logout_uri,omitempty"`
	PostLogoutRedirectURIs []string `json:"post_logout_redirect_uris,omitempty" yaml:"post_logout_redirect_uris,omitempty"`

	// Additional fields for common settings
	Enabled             bool   `json:"enabled" yaml:"enabled"`
	ClientType          string `json:"client_type,omitempty" yaml:"client_type,omitempty"` // "public" or "confidential"
	PKCERequired        bool   `json:"pkce_required,omitempty" yaml:"pkce_required,omitempty"`
	DPoPRequired        bool   `json:"dpop_required,omitempty" yaml:"dpop_required,omitempty"`
	PARRequired         bool   `json:"par_required,omitempty" yaml:"par_required,omitempty"`
	PromptScopeApproval bool   `json:"prompt_scope_approval,omitempty" yaml:"prompt_scope_approval,omitempty"`

	// Token settings
	AccessTokenFormat   string `json:"access_token_format,omitempty" yaml:"access_token_format,omitempty"`     // "jwt" or "opaque"
	AccessTokenLifetime int    `json:"access_token_lifetime,omitempty" yaml:"access_token_lifetime,omitempty"` // lifetime in seconds
	RotateRefreshTokens bool   `json:"rotate_refresh_tokens,omitempty" yaml:"rotate_refresh_tokens,omitempty"`
}

type CanonicalClientConfigExtensions struct {
	// Add custom extension fields here if needed
	SSO                      bool `json:"sso,omitempty" yaml:"sso,omitempty"`
	Enforce2SV               bool `json:"enforce_2sv,omitempty" yaml:"enforce_2sv,omitempty"`
	PromptTermsAndConditions bool `json:"prompt_terms_and_conditions,omitempty" yaml:"prompt_terms_and_conditions,omitempty"`
}

type CanonicalClientConfigSecrets struct {
	// Client secrets or certificates can be stored here if needed
	PlainSecret     string `json:"plain_secret,omitempty" yaml:"plain_secret,omitempty"`
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
		return fmt.Errorf("error marshalling cl %s: %v", c.Client.ClientID, err)
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
