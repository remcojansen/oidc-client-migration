package oidcconfig

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testConfig() *CanonicalClientConfig {
	return &CanonicalClientConfig{
		Client: CanonicalClientConfigClient{
			ClientID: "my-client",
			Name:     "My Client",
		},
		Extensions: CanonicalClientConfigExtensions{
			Enabled: true,
		},
		Secrets: CanonicalClientConfigSecrets{
			PlainSecret: "s3cr3t",
		},
	}
}

func TestWriteConfigFile_JSON(t *testing.T) {
	dir := t.TempDir()
	cfg := testConfig()

	if err := cfg.WriteConfigFile(dir, "json"); err != nil {
		t.Fatalf("WriteConfigFile returned error: %v", err)
	}

	path := filepath.Join(dir, "my-client.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("expected file %s to exist: %v", path, err)
	}
	if !strings.Contains(string(data), `"client_id": "my-client"`) {
		t.Errorf("expected JSON output to contain client_id, got: %s", data)
	}
}

func TestWriteConfigFile_YAML(t *testing.T) {
	dir := t.TempDir()
	cfg := testConfig()

	if err := cfg.WriteConfigFile(dir, "yaml"); err != nil {
		t.Fatalf("WriteConfigFile returned error: %v", err)
	}

	path := filepath.Join(dir, "my-client.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("expected file %s to exist: %v", path, err)
	}
	if !strings.Contains(string(data), "client_id: my-client") {
		t.Errorf("expected YAML output to contain client_id, got: %s", data)
	}
}

func TestWriteConfigFile_InvalidOutputDir(t *testing.T) {
	cfg := testConfig()
	err := cfg.WriteConfigFile(filepath.Join(t.TempDir(), "does-not-exist"), "json")
	if err == nil {
		t.Fatal("expected an error when writing to a non-existent directory, got nil")
	}
}

func TestWriteConfigFile_ExplicitFalseBooleanIsNotOmitted(t *testing.T) {
	// Regression test: consent_required, pkce_required, dpop_required, par_required,
	// rotate_refresh_tokens, and terms_and_conditions_required must always be serialized,
	// even when false, rather than being dropped by `omitempty` (which previously made it
	// impossible to tell "explicitly not required" apart from "not captured at all").
	dir := t.TempDir()
	cfg := testConfig()
	cfg.Client.ConsentRequired = false
	cfg.Extensions.PKCERequired = false

	if err := cfg.WriteConfigFile(dir, "yaml"); err != nil {
		t.Fatalf("WriteConfigFile returned error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "my-client.yaml"))
	if err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
	if !strings.Contains(string(data), "consent_required: false") {
		t.Errorf("expected YAML output to contain explicit consent_required: false, got: %s", data)
	}
	if !strings.Contains(string(data), "pkce_required: false") {
		t.Errorf("expected YAML output to contain explicit pkce_required: false, got: %s", data)
	}
}
