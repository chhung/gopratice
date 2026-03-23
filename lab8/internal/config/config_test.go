package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadUsesDefaultPortAndPrefixes(t *testing.T) {
	t.Parallel()

	configPath := writeTempConfig(t, `{
		"broker": {
			"host": "127.0.0.1"
		},
		"destinations": {
			"queue": "orders.in",
			"topic": "events.broadcast"
		}
	}`)

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if got := cfg.Broker.Port; got != defaultPort {
		t.Fatalf("Broker.Port = %d, want %d", got, defaultPort)
	}

	if got := cfg.QueueDestination(); got != "/queue/orders.in" {
		t.Fatalf("QueueDestination() = %q, want %q", got, "/queue/orders.in")
	}

	if got := cfg.TopicDestination(); got != "/topic/events.broadcast" {
		t.Fatalf("TopicDestination() = %q, want %q", got, "/topic/events.broadcast")
	}

	if got := cfg.Address(); got != "127.0.0.1:61613" {
		t.Fatalf("Address() = %q, want %q", got, "127.0.0.1:61613")
	}
}

func TestLoadRejectsMissingDestinations(t *testing.T) {
	t.Parallel()

	configPath := writeTempConfig(t, `{
		"broker": {
			"host": "127.0.0.1",
			"port": 61613
		},
		"destinations": {}
	}`)

	_, err := Load(configPath)
	if err == nil {
		t.Fatal("Load() error = nil, want validation error")
	}

	if !strings.Contains(err.Error(), "at least one of destinations.queue or destinations.topic is required") {
		t.Fatalf("Load() error = %v, want destination validation message", err)
	}
}

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")
	if err := os.WriteFile(configPath, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	return configPath
}
