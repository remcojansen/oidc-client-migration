package authserver

import "ocm/oidcconfig"

type AuthServerClient interface {
	FetchClientConfigurations() ([]OAuthClientConfig, error)
	WithBaseURL(baseURL string) AuthServerClient
	WithUsernamePassword(username, password string) AuthServerClient
	WithAccessToken(token string) AuthServerClient
}

type OAuthClientConfig interface {
	GetCanonicalClientConfig() *oidcconfig.CanonicalClientConfig
	GetClientID() string
}
