package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
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
	if cfg.Database.Enabled != false {
		t.Errorf("expected database.enabled = false, got %v", cfg.Database.Enabled)
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
	if cfg.Database.MigrationsPath != "./migrations" {
		t.Errorf("expected database.migrationspath = './migrations', got '%s'", cfg.Database.MigrationsPath)
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

func TestLoadWithCmdFlags(t *testing.T) {
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

	// Create a test command with flags
	cmd := &cobra.Command{
		Use: "test",
	}

	// Add flags (matching what's in root.go)
	cmd.Flags().String("server-host", "", "server host")
	cmd.Flags().Int("server-port", 0, "server port")
	cmd.Flags().Bool("db-enabled", false, "enable database connection")
	cmd.Flags().String("db-driver", "", "database driver")
	cmd.Flags().String("db-host", "", "database host")
	cmd.Flags().Int("db-port", 0, "database port")
	cmd.Flags().String("db-name", "", "database name")
	cmd.Flags().String("db-user", "", "database user")
	cmd.Flags().String("db-password", "", "database password")
	cmd.Flags().String("db-sslmode", "", "database SSL mode")
	cmd.Flags().String("db-migrations-path", "", "database migrations path")
	cmd.Flags().String("log-level", "", "log level")
	cmd.Flags().String("log-format", "", "log format")

	// Set flag values
	cmd.Flags().Set("server-host", "flag-host.example.com")
	cmd.Flags().Set("server-port", "7000")
	cmd.Flags().Set("db-host", "flag-db.example.com")
	cmd.Flags().Set("db-port", "5435")
	cmd.Flags().Set("db-name", "flagdb")
	cmd.Flags().Set("db-user", "flaguser")
	cmd.Flags().Set("db-password", "flagpass")
	cmd.Flags().Set("db-sslmode", "verify-full")
	cmd.Flags().Set("log-level", "error")
	cmd.Flags().Set("log-format", "text")

	cfg, err := LoadWithCmd("", cmd)
	if err != nil {
		t.Fatalf("LoadWithCmd() error = %v", err)
	}

	// Check values from CLI flags
	if cfg.Server.Host != "flag-host.example.com" {
		t.Errorf("expected server.host = 'flag-host.example.com', got '%s'", cfg.Server.Host)
	}
	if cfg.Server.Port != 7000 {
		t.Errorf("expected server.port = 7000, got %d", cfg.Server.Port)
	}
	if cfg.Database.Host != "flag-db.example.com" {
		t.Errorf("expected database.host = 'flag-db.example.com', got '%s'", cfg.Database.Host)
	}
	if cfg.Database.Port != 5435 {
		t.Errorf("expected database.port = 5435, got %d", cfg.Database.Port)
	}
	if cfg.Database.Name != "flagdb" {
		t.Errorf("expected database.name = 'flagdb', got '%s'", cfg.Database.Name)
	}
	if cfg.Database.User != "flaguser" {
		t.Errorf("expected database.user = 'flaguser', got '%s'", cfg.Database.User)
	}
	if cfg.Database.Password != "flagpass" {
		t.Errorf("expected database.password = 'flagpass', got '%s'", cfg.Database.Password)
	}
	if cfg.Database.SSLMode != "verify-full" {
		t.Errorf("expected database.sslmode = 'verify-full', got '%s'", cfg.Database.SSLMode)
	}
	if cfg.Log.Level != "error" {
		t.Errorf("expected log.level = 'error', got '%s'", cfg.Log.Level)
	}
	if cfg.Log.Format != "text" {
		t.Errorf("expected log.format = 'text', got '%s'", cfg.Log.Format)
	}
}

func TestConfigPriorityOrder(t *testing.T) {
	// Test that CLI flags > Environment variables > Config file > Defaults
	tmpDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	// Create a config file
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
server:
  host: "config-host"
  port: 8000
log:
  level: "info"
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Set environment variables (should override config file)
	os.Setenv("APP_SERVER_HOST", "env-host")
	os.Setenv("APP_LOG_LEVEL", "debug")
	defer func() {
		os.Unsetenv("APP_SERVER_HOST")
		os.Unsetenv("APP_LOG_LEVEL")
	}()

	// Create command with CLI flags (should override everything)
	cmd := &cobra.Command{
		Use: "test",
	}
	cmd.Flags().String("server-host", "", "server host")
	cmd.Flags().Int("server-port", 0, "server port")
	cmd.Flags().Bool("db-enabled", false, "enable database connection")
	cmd.Flags().String("db-driver", "", "database driver")
	cmd.Flags().String("db-host", "", "database host")
	cmd.Flags().Int("db-port", 0, "database port")
	cmd.Flags().String("db-name", "", "database name")
	cmd.Flags().String("db-user", "", "database user")
	cmd.Flags().String("db-password", "", "database password")
	cmd.Flags().String("db-sslmode", "", "database SSL mode")
	cmd.Flags().String("db-migrations-path", "", "database migrations path")
	cmd.Flags().String("log-level", "", "log level")
	cmd.Flags().String("log-format", "", "log format")

	// Set only server-port via CLI flag
	cmd.Flags().Set("server-port", "9000")

	cfg, err := LoadWithCmd(configPath, cmd)
	if err != nil {
		t.Fatalf("LoadWithCmd() error = %v", err)
	}

	// server.host: env var should override config file
	if cfg.Server.Host != "env-host" {
		t.Errorf("expected server.host = 'env-host' (from env), got '%s'", cfg.Server.Host)
	}

	// server.port: CLI flag should override config file
	if cfg.Server.Port != 9000 {
		t.Errorf("expected server.port = 9000 (from CLI), got %d", cfg.Server.Port)
	}

	// log.level: env var should override config file
	if cfg.Log.Level != "debug" {
		t.Errorf("expected log.level = 'debug' (from env), got '%s'", cfg.Log.Level)
	}

	// log.format: should use default (no override)
	if cfg.Log.Format != "json" {
		t.Errorf("expected log.format = 'json' (from default), got '%s'", cfg.Log.Format)
	}
}

func TestLoadWithCmdNilCommand(t *testing.T) {
	// Test that LoadWithCmd works when cmd is nil (should behave like Load)
	tmpDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadWithCmd("", nil)
	if err != nil {
		t.Fatalf("LoadWithCmd() with nil cmd error = %v", err)
	}

	// Should return defaults
	if cfg.Server.Port != 8080 {
		t.Errorf("expected server.port = 8080, got %d", cfg.Server.Port)
	}
}
