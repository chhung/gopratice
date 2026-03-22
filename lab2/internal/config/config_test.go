package config

import (
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("APP_PORT", "")
	t.Setenv("SHUTDOWN_TOKEN", "")
	t.Setenv("HTTP_READ_HEADER_TIMEOUT", "")
	t.Setenv("HTTP_READ_TIMEOUT", "")
	t.Setenv("HTTP_WRITE_TIMEOUT", "")
	t.Setenv("HTTP_IDLE_TIMEOUT", "")
	t.Setenv("HTTP_SHUTDOWN_TIMEOUT", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Port != defaultPort {
		t.Fatalf("Load().Port = %q, want %q", cfg.Port, defaultPort)
	}

	if cfg.ReadTimeout != defaultReadTimeout {
		t.Fatalf("Load().ReadTimeout = %s, want %s", cfg.ReadTimeout, defaultReadTimeout)
	}
}

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("APP_PORT", "9090")
	t.Setenv("SHUTDOWN_TOKEN", "secret-token")
	t.Setenv("HTTP_READ_HEADER_TIMEOUT", "1s")
	t.Setenv("HTTP_READ_TIMEOUT", "3s")
	t.Setenv("HTTP_WRITE_TIMEOUT", "7s")
	t.Setenv("HTTP_IDLE_TIMEOUT", "45s")
	t.Setenv("HTTP_SHUTDOWN_TIMEOUT", "9s")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Port != "9090" {
		t.Fatalf("Load().Port = %q, want %q", cfg.Port, "9090")
	}

	if cfg.ShutdownToken != "secret-token" {
		t.Fatalf("Load().ShutdownToken = %q, want %q", cfg.ShutdownToken, "secret-token")
	}

	if cfg.ReadHeaderTimeout != time.Second {
		t.Fatalf("Load().ReadHeaderTimeout = %s, want %s", cfg.ReadHeaderTimeout, time.Second)
	}

	if cfg.ReadTimeout != 3*time.Second {
		t.Fatalf("Load().ReadTimeout = %s, want %s", cfg.ReadTimeout, 3*time.Second)
	}

	if cfg.WriteTimeout != 7*time.Second {
		t.Fatalf("Load().WriteTimeout = %s, want %s", cfg.WriteTimeout, 7*time.Second)
	}

	if cfg.IdleTimeout != 45*time.Second {
		t.Fatalf("Load().IdleTimeout = %s, want %s", cfg.IdleTimeout, 45*time.Second)
	}

	if cfg.ShutdownTimeout != 9*time.Second {
		t.Fatalf("Load().ShutdownTimeout = %s, want %s", cfg.ShutdownTimeout, 9*time.Second)
	}

	if cfg.HTTPAddress() != ":9090" {
		t.Fatalf("HTTPAddress() = %q, want %q", cfg.HTTPAddress(), ":9090")
	}
}

func TestLoadRejectsInvalidDuration(t *testing.T) {
	t.Setenv("HTTP_READ_TIMEOUT", "not-a-duration")

	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want non-nil")
	}
}
