package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadValidConfig(t *testing.T) {
	t.Parallel()

	path := writeTempConfig(t, `
mongodb:
  uri: "mongodb://localhost:27017"
  database: "testdb"
  timeout_seconds: 5
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.MongoDB.URI != "mongodb://localhost:27017" {
		t.Fatalf("MongoDB.URI = %q, want %q", cfg.MongoDB.URI, "mongodb://localhost:27017")
	}

	if cfg.MongoDB.Database != "testdb" {
		t.Fatalf("MongoDB.Database = %q, want %q", cfg.MongoDB.Database, "testdb")
	}

	if cfg.MongoDB.TimeoutSeconds != 5 {
		t.Fatalf("MongoDB.TimeoutSeconds = %d, want %d", cfg.MongoDB.TimeoutSeconds, 5)
	}

	if got := cfg.MongoTimeout(); got != 5*time.Second {
		t.Fatalf("MongoTimeout() = %v, want %v", got, 5*time.Second)
	}
}

func TestLoadUsesDefaultTimeout(t *testing.T) {
	t.Parallel()

	path := writeTempConfig(t, `
mongodb:
  uri: "mongodb://localhost:27017"
  database: "testdb"
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.MongoDB.TimeoutSeconds != 10 {
		t.Fatalf("MongoDB.TimeoutSeconds = %d, want default %d", cfg.MongoDB.TimeoutSeconds, 10)
	}
}

func TestLoadRejectsMissingURI(t *testing.T) {
	t.Parallel()

	path := writeTempConfig(t, `
mongodb:
  database: "testdb"
`)

	_, err := Load(path)
	if err == nil {
		t.Fatal("Load() error = nil, want validation error")
	}

	if !strings.Contains(err.Error(), "mongodb.uri is required") {
		t.Fatalf("Load() error = %v, want URI validation message", err)
	}
}

func TestLoadRejectsMissingDatabase(t *testing.T) {
	t.Parallel()

	path := writeTempConfig(t, `
mongodb:
  uri: "mongodb://localhost:27017"
`)

	_, err := Load(path)
	if err == nil {
		t.Fatal("Load() error = nil, want validation error")
	}

	if !strings.Contains(err.Error(), "mongodb.database is required") {
		t.Fatalf("Load() error = %v, want database validation message", err)
	}
}

func TestLoadRejectsInvalidTimeout(t *testing.T) {
	t.Parallel()

	path := writeTempConfig(t, `
mongodb:
  uri: "mongodb://localhost:27017"
  database: "testdb"
  timeout_seconds: 0
`)

	_, err := Load(path)
	if err == nil {
		t.Fatal("Load() error = nil, want validation error")
	}

	if !strings.Contains(err.Error(), "mongodb.timeout_seconds must be at least 1") {
		t.Fatalf("Load() error = %v, want timeout validation message", err)
	}
}

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	return configPath
}
