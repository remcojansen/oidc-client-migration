package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"ocm/export/authserver"
	"ocm/export/authserver/keycloak"
	"ocm/export/authserver/pingfederate"
	"os"
)

const (
	ExitError = 1
	ExitOK    = 0
)

// parseAccessTokenManagerMapping parses AUTH_SERVER_ACCESS_TOKEN_MANAGER_MAP, a JSON object
// mapping PingFederate access token manager IDs to their format/lifetime, e.g.
// {"my-jwt-manager-short": {"format": "jwt", "lifetime_seconds": 300}}. An empty string is treated
// as an empty mapping (no access token managers known).
func parseAccessTokenManagerMapping(raw string) (map[string]pingfederate.AccessTokenManagerInfo, error) {
	if raw == "" {
		return map[string]pingfederate.AccessTokenManagerInfo{}, nil
	}
	var mapping map[string]pingfederate.AccessTokenManagerInfo
	if err := json.Unmarshal([]byte(raw), &mapping); err != nil {
		return nil, err
	}
	return mapping, nil
}

func main() {
	var sourceSystem, outputDir, format, clientId string
	var verbose bool
	flag.StringVar(&sourceSystem, "source", "keycloak", "Source system: pingfederate or keycloak")
	flag.StringVar(&outputDir, "dir", ".", "Directory to save client configuration files")
	flag.StringVar(&format, "format", "yaml", "Output format: json or yaml")
	flag.StringVar(&clientId, "client-id", "", "Client ID to fetch (optional)")
	flag.BoolVar(&verbose, "verbose", false, "Print the URL and response status code for each HTTP request made")
	flag.Parse()
	if format != "json" && format != "yaml" {
		fmt.Println("Invalid format. Use 'json' or 'yaml'.")
		os.Exit(ExitError)
	}

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create output directory: %v\n", err)
	}

	var c authserver.AuthServerClient
	switch sourceSystem {
	case "pingfederate":
		atmMapping, err := parseAccessTokenManagerMapping(os.Getenv("AUTH_SERVER_ACCESS_TOKEN_MANAGER_MAP"))
		if err != nil {
			log.Fatalf("Failed to parse AUTH_SERVER_ACCESS_TOKEN_MANAGER_MAP: %v\n", err)
		}
		c = pingfederate.CreatePingFederateClient().
			WithAccessTokenManagerMapping(atmMapping).
			WithBaseURL(os.Getenv("AUTH_SERVER_BASE_URL")).
			WithUsernamePassword(os.Getenv("AUTH_SERVER_USERNAME"), os.Getenv("AUTH_SERVER_PASSWORD")).
			WithVerbose(verbose)
	case "keycloak":
		c = keycloak.CreateKeycloakClient().
			WithBaseURL(os.Getenv("AUTH_SERVER_BASE_URL")).
			WithAccessToken(os.Getenv("AUTH_SERVER_ACCESS_TOKEN")).
			WithVerbose(verbose)
	default:
		log.Fatalf("Unsupported source system: %s\n", sourceSystem)
	}

	oauthClients := []authserver.OAuthClientConfig{}
	if clientId != "" {
		clientConfig, err := c.FetchClientConfigurationByClientId(clientId)
		if err != nil {
			log.Fatalf("Error fetching client configuration for %s: %v\n", clientId, err)
		}
		if clientConfig == nil {
			log.Fatalf("Client configuration for %s not found\n", clientId)
		}
		oauthClients = append(oauthClients, clientConfig)
	} else {
		var err error
		oauthClients, err = c.FetchClientConfigurations()
		if err != nil {
			log.Fatalf("Error fetching oauthClients: %v\n", err)
		}
	}

	exitCode := ExitOK
	for _, cl := range oauthClients {
		clientConfig := cl.GetCanonicalClientConfig()
		if clientConfig == nil {
			log.Printf("Error mapping client %s: unsupported config", cl.GetClientID())
			exitCode = ExitError
			continue
		}

		if err := clientConfig.WriteConfigFile(outputDir, format); err != nil {
			log.Printf("Error writing config for client %s: %v", cl.GetClientID(), err)
			exitCode = ExitError
			continue
		}
	}
	os.Exit(exitCode)
}
