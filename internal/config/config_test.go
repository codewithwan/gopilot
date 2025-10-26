package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	// Test default values
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Server.Port != "8080" {
		t.Errorf("Expected default port 8080, got %s", cfg.Server.Port)
	}

	if cfg.Database.Host != "localhost" {
		t.Errorf("Expected default database host localhost, got %s", cfg.Database.Host)
	}

	if cfg.JWT.Expiration != 24*time.Hour {
		t.Errorf("Expected default JWT expiration 24h, got %v", cfg.JWT.Expiration)
	}
}

func TestLoadWithEnv(t *testing.T) {
	// Test environment variable override
	// Viper uses nested keys with dots
	os.Setenv("SERVER_PORT", "9090")
	defer os.Unsetenv("SERVER_PORT")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Server.Port != "9090" {
		t.Errorf("Expected port 9090 from env, got %s", cfg.Server.Port)
	}
}

func TestDatabaseDSN(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	dsn := cfg.Database.DSN()
	expectedDSN := "host=localhost port=5432 user=postgres password=postgres dbname=gopilot sslmode=disable"
	if dsn != expectedDSN {
		t.Errorf("Expected DSN %s, got %s", expectedDSN, dsn)
	}
}
