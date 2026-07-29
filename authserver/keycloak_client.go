package authserver

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type KeycloakClient struct {
	baseURL string
	realm   string
	client  *http.Client
	header  *http.Header
}

func CreateKeycloakClient() *KeycloakClient {
	c := &KeycloakClient{
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
	c.header.Set("Accept", "application/json")
	c.header.Set("Content-Type", "application/json")

	if c.realm == "" {
		c.realm = "master"
	}
	return c
}

func (c *KeycloakClient) WithUsernamePassword(username, password string) AuthServerClient {
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
	c.header.Set("Authorization", basicAuth)
	return c
}

func (c *KeycloakClient) WithAccessToken(token string) AuthServerClient {
	c.header.Set("Authorization", "Bearer "+token)
	return c
}

func (c *KeycloakClient) WithBaseURL(baseURL string) AuthServerClient {
	baseURL = strings.TrimSuffix(baseURL, "/")

	if before, after, ok := strings.Cut(baseURL, "/realms/"); ok {
		c.baseURL = before
		c.realm = after
	} else {
		c.baseURL = baseURL
	}

	return c
}

func (c *KeycloakClient) FetchClientConfigurations() ([]OAuthClientConfig, error) {
	// Keycloak Admin API endpoint for fetching clients
	url := fmt.Sprintf("%s/admin/realms/%s/clients", c.baseURL, c.realm)

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
		return nil, fmt.Errorf("keycloak API returned status: %s", resp.Status)
	}

	var clientList []KeycloakClientConfig
	if err := json.NewDecoder(resp.Body).Decode(&clientList); err != nil {
		return nil, err
	}

	oauthClients := make([]OAuthClientConfig, len(clientList))
	for i := range clientList {
		oauthClients[i] = &clientList[i]
	}

	return oauthClients, nil
}

func (c *KeycloakClient) FetchClientConfigurationByClientId(clientId string) (OAuthClientConfig, error) {
	// Keycloak Admin API endpoint for fetching a specific client by ID
	url := fmt.Sprintf("%s/admin/realms/%s/clients/%s", c.baseURL, c.realm, clientId)

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
		return nil, fmt.Errorf("keycloak API returned status: %s", resp.Status)
	}

	var client KeycloakClientConfig
	if err := json.NewDecoder(resp.Body).Decode(&client); err != nil {
		return nil, err
	}

	return &client, nil
}
