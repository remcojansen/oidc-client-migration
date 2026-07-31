package keycloak

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchClientConfigurationByClientId_LooksUpByExternalClientId(t *testing.T) {
	const externalClientId = "my-app"
	const internalId = "11111111-1111-1111-1111-111111111111"

	var lookupCalled bool
	var requestCount int

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if r.URL.Path == "/admin/realms/master/clients" && r.URL.Query().Get("clientId") == externalClientId {
			lookupCalled = true
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode([]KeycloakClientConfig{
				{ID: internalId, ClientID: externalClientId},
			})
			return
		}
		t.Errorf("unexpected request: %s %s", r.Method, r.URL.String())
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	c := CreateKeycloakClient().WithBaseURL(server.URL)

	config, err := c.FetchClientConfigurationByClientId(externalClientId)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if config == nil {
		t.Fatal("expected a client configuration, got nil")
	}
	if config.GetClientID() != externalClientId {
		t.Errorf("expected client ID %q, got %q", externalClientId, config.GetClientID())
	}
	if !lookupCalled {
		t.Error("expected the clientId lookup endpoint to be called")
	}
	if requestCount != 1 {
		t.Errorf("expected exactly 1 HTTP request, got %d", requestCount)
	}
}

func TestFetchClientConfigurationByClientId_ReturnsNilWhenNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]KeycloakClientConfig{})
	}))
	defer server.Close()

	c := CreateKeycloakClient().WithBaseURL(server.URL)

	config, err := c.FetchClientConfigurationByClientId("does-not-exist")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if config != nil {
		t.Errorf("expected nil configuration, got %+v", config)
	}
}
