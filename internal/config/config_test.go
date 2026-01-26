package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Check default values
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("expected server.host = '0.0.0.0', got '%s'", cfg.Server.Host)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("expected server.port = 8080, got %d", cfg.Server.Port)
	}
	if cfg.Database.Driver != "postgres" {
		t.Errorf("expected database.driver = 'postgres', got '%s'", cfg.Database.Driver)
	}
	if cfg.Database.Host != "localhost" {
		t.Errorf("expected database.host = 'localhost', got '%s'", cfg.Database.Host)
	}
	if cfg.Database.Port != 5432 {
		t.Errorf("expected database.port = 5432, got %d", cfg.Database.Port)
	}
	if cfg.Log.Level != "info" {
		t.Errorf("expected log.level = 'info', got '%s'", cfg.Log.Level)
	}
}

func TestLoadFromConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
server:
  host: "127.0.0.1"
  port: 9000
database:
  driver: "postgres"
  host: "db.example.com"
  port: 5433
  name: "testdb"
  user: "testuser"
  password: "testpass"
  sslmode: "require"
log:
  level: "debug"
  format: "text"
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Check values from config file
	if cfg.Server.Host != "127.0.0.1" {
		t.Errorf("expected server.host = '127.0.0.1', got '%s'", cfg.Server.Host)
	}
	if cfg.Server.Port != 9000 {
		t.Errorf("expected server.port = 9000, got %d", cfg.Server.Port)
	}
	if cfg.Database.Host != "db.example.com" {
		t.Errorf("expected database.host = 'db.example.com', got '%s'", cfg.Database.Host)
	}
	if cfg.Database.Port != 5433 {
		t.Errorf("expected database.port = 5433, got %d", cfg.Database.Port)
	}
	if cfg.Database.Name != "testdb" {
		t.Errorf("expected database.name = 'testdb', got '%s'", cfg.Database.Name)
	}
	if cfg.Database.User != "testuser" {
		t.Errorf("expected database.user = 'testuser', got '%s'", cfg.Database.User)
	}
	if cfg.Database.Password != "testpass" {
		t.Errorf("expected database.password = 'testpass', got '%s'", cfg.Database.Password)
	}
	if cfg.Database.SSLMode != "require" {
		t.Errorf("expected database.sslmode = 'require', got '%s'", cfg.Database.SSLMode)
	}
	if cfg.Log.Level != "debug" {
		t.Errorf("expected log.level = 'debug', got '%s'", cfg.Log.Level)
	}
	if cfg.Log.Format != "text" {
		t.Errorf("expected log.format = 'text', got '%s'", cfg.Log.Format)
	}
}

func TestLoadFromEnvironment(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	// Set environment variables
	os.Setenv("APP_SERVER_HOST", "192.168.1.1")
	os.Setenv("APP_SERVER_PORT", "3000")
	os.Setenv("APP_DATABASE_HOST", "env-db.example.com")
	os.Setenv("APP_DATABASE_PORT", "5434")
	os.Setenv("APP_DATABASE_NAME", "envdb")
	os.Setenv("APP_DATABASE_USER", "envuser")
	os.Setenv("APP_DATABASE_PASSWORD", "envpass")
	os.Setenv("APP_LOG_LEVEL", "warn")

	defer func() {
		os.Unsetenv("APP_SERVER_HOST")
		os.Unsetenv("APP_SERVER_PORT")
		os.Unsetenv("APP_DATABASE_HOST")
		os.Unsetenv("APP_DATABASE_PORT")
		os.Unsetenv("APP_DATABASE_NAME")
		os.Unsetenv("APP_DATABASE_USER")
		os.Unsetenv("APP_DATABASE_PASSWORD")
		os.Unsetenv("APP_LOG_LEVEL")
	}()

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Check values from environment
	if cfg.Server.Host != "192.168.1.1" {
		t.Errorf("expected server.host = '192.168.1.1', got '%s'", cfg.Server.Host)
	}
	if cfg.Server.Port != 3000 {
		t.Errorf("expected server.port = 3000, got %d", cfg.Server.Port)
	}
	if cfg.Database.Host != "env-db.example.com" {
		t.Errorf("expected database.host = 'env-db.example.com', got '%s'", cfg.Database.Host)
	}
	if cfg.Database.Port != 5434 {
		t.Errorf("expected database.port = 5434, got %d", cfg.Database.Port)
	}
	if cfg.Database.Name != "envdb" {
		t.Errorf("expected database.name = 'envdb', got '%s'", cfg.Database.Name)
	}
	if cfg.Database.User != "envuser" {
		t.Errorf("expected database.user = 'envuser', got '%s'", cfg.Database.User)
	}
	if cfg.Database.Password != "envpass" {
		t.Errorf("expected database.password = 'envpass', got '%s'", cfg.Database.Password)
	}
	if cfg.Log.Level != "warn" {
		t.Errorf("expected log.level = 'warn', got '%s'", cfg.Log.Level)
	}
}

func TestLoadInvalidConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.yaml")

	// Create invalid YAML
	invalidContent := `
server:
  host: [invalid yaml structure
`

	if err := os.WriteFile(configPath, []byte(invalidContent), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Error("expected error for invalid config file, got nil")
	}
}
