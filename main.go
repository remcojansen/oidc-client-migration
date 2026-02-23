package main

import (
	"flag"
	"fmt"
	"log"
	"ocm/authserver"
	"os"
)

const (
	ExitError = 1
	ExitOK    = 0
)

func main() {
	var sourceSystem, outputDir, format string
	flag.StringVar(&sourceSystem, "source", "keycloak", "Source system: pingfederate or keycloak")
	flag.StringVar(&outputDir, "dir", ".", "Directory to save client configuration files")
	flag.StringVar(&format, "format", "yaml", "Output format: json or yaml")
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
		c = authserver.CreatePingFederateClient().
			WithBaseURL(os.Getenv("AUTH_SERVER_BASE_URL")).
			WithUsernamePassword(os.Getenv("AUTH_SERVER_USERNAME"), os.Getenv("AUTH_SERVER_PASSWORD"))
	case "keycloak":
		c = authserver.CreateKeycloakClient().
			WithBaseURL(os.Getenv("AUTH_SERVER_BASE_URL")).
			WithAccessToken(os.Getenv("AUTH_SERVER_ACCESS_TOKEN"))
	default:
		log.Fatalf("Unsupported source system: %s\n", sourceSystem)
	}

	oauthClients, err := c.FetchClientConfigurations()
	if err != nil {
		log.Fatalf("Error fetching oauthClients: %v\n", err)
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
