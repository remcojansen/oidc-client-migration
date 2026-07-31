package authserver

import "ocm/export/oidcconfig"

type AuthServerClient interface {
	FetchClientConfigurations() ([]OAuthClientConfig, error)
	FetchClientConfigurationByClientId(clientId string) (OAuthClientConfig, error)
	WithBaseURL(baseURL string) AuthServerClient
	WithUsernamePassword(username, password string) AuthServerClient
	WithAccessToken(token string) AuthServerClient
	// WithVerbose enables logging of each HTTP request's URL and response status code.
	WithVerbose(verbose bool) AuthServerClient
}

type OAuthClientConfig interface {
	GetCanonicalClientConfig() *oidcconfig.CanonicalClientConfig
	GetClientID() string
}
