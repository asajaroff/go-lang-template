package config_test

import (
	"testing"

	"github.com/asajaroff/go-lang-template/internal/config"
)

func TestLoadRequiresDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "")

	if _, err := config.Load(); err == nil {
		t.Fatal("expected an error when DATABASE_URL is unset, got nil")
	}
}

func TestLoadAppliesDefaults(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Port != "8000" {
		t.Errorf("got port %q, want 8000", cfg.Port)
	}
	if cfg.AppName != "go-lang-template" {
		t.Errorf("got app name %q, want go-lang-template", cfg.AppName)
	}
}

func TestLoadReadsOverrides(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")
	t.Setenv("PORT", "9090")
	t.Setenv("LOG_LEVEL", "debug")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Port != "9090" {
		t.Errorf("got port %q, want 9090", cfg.Port)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("got log level %q, want debug", cfg.LogLevel)
	}
}
