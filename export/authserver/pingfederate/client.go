package pingfederate

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"ocm/export/authserver"
	"strings"
	"time"
)

type PingFederateClient struct {
	apiURL                    string
	client                    *http.Client
	header                    *http.Header
	accessTokenManagerMapping map[string]AccessTokenManagerInfo
}

func CreatePingFederateClient() *PingFederateClient {
	c := &PingFederateClient{
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
			Timeout: 10 * time.Second,
		},
		header: &http.Header{},
	}
	c.header.Set("X-XSRF-Header", "PingFederate")
	c.header.Set("Accept", "application/json")
	return c
}

// WithAccessTokenManagerMapping configures the access token manager ID -> format/lifetime mapping,
// since access token managers are user-created, deployment-specific resources. This is
// PingFederate-specific and not part of the authserver.AuthServerClient interface, so it must be
// called first in the builder chain, before any other With... call.
func (c *PingFederateClient) WithAccessTokenManagerMapping(mapping map[string]AccessTokenManagerInfo) *PingFederateClient {
	c.accessTokenManagerMapping = mapping
	return c
}

func (c *PingFederateClient) WithUsernamePassword(username, password string) authserver.AuthServerClient {
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
	c.header.Set("Authorization", basicAuth)
	return c
}

func (c *PingFederateClient) WithAccessToken(token string) authserver.AuthServerClient {
	c.header.Set("Authorization", "Bearer "+token)
	return c
}

func (c *PingFederateClient) WithBaseURL(baseURL string) authserver.AuthServerClient {
	c.apiURL = strings.TrimSuffix(baseURL, "/") + "/pf-admin-api/v1"
	return c
}

func (c *PingFederateClient) WithVerbose(verbose bool) authserver.AuthServerClient {
	if verbose {
		c.client.Transport = &authserver.VerboseRoundTripper{Transport: c.client.Transport}
	}
	return c
}

func (c *PingFederateClient) FetchClientConfigurations() ([]authserver.OAuthClientConfig, error) {
	req, err := http.NewRequest("GET", c.apiURL+"/oauth/clients", nil)
	if err != nil {
		return nil, err
	}
	req.Header = *c.header

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("error closing response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pingFederate API returned status: %s", resp.Status)
	}

	var clientList PingFederateClientList
	if err := json.NewDecoder(resp.Body).Decode(&clientList); err != nil {
		return nil, err
	}

	oauthClients := make([]authserver.OAuthClientConfig, len(clientList.Items))
	for i := range clientList.Items {
		clientList.Items[i].accessTokenManagerMapping = c.accessTokenManagerMapping
		oauthClients[i] = &clientList.Items[i]
	}

	return oauthClients, nil
}

func (c *PingFederateClient) FetchClientConfigurationByClientId(clientId string) (authserver.OAuthClientConfig, error) {
	url := fmt.Sprintf("%s/oauth/clients/%s", c.apiURL, clientId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = *c.header

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("error closing response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pingFederate API returned status: %s", resp.Status)
	}

	var client PingFederateClientConfig
	if err := json.NewDecoder(resp.Body).Decode(&client); err != nil {
		return nil, err
	}
	client.accessTokenManagerMapping = c.accessTokenManagerMapping

	return &client, nil
}
