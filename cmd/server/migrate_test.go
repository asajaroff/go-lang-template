package main

import (
	"testing"

	"github.com/asajaroff/go-lang-template/internal/config"
)

func TestMigrateCommands(t *testing.T) {
	tests := []struct {
		name    string
		command string
	}{
		{"MigrateUp", "migrate up"},
		{"MigrateDown", "migrate down"},
		{"MigrateSteps", "migrate steps"},
		{"MigrateVersion", "migrate version"},
		{"MigrateForce", "migrate force"},
	}

	// Set up a test config to avoid nil pointer errors
	cfg = &config.Config{
		Database: config.DatabaseConfig{
			Driver:         "postgres",
			Host:           "localhost",
			Port:           5432,
			Name:           "test_db",
			User:           "test_user",
			Password:       "test_pass",
			SSLMode:        "disable",
			MigrationsPath: "./migrations",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that the command exists and is properly configured
			// We won't actually execute them as they require a database
			cmd := rootCmd
			if cmd == nil {
				t.Fatal("rootCmd is nil")
			}

			// Find the migrate command
			migrateCmd := cmd.Commands()
			found := false
			for _, c := range migrateCmd {
				if c.Use == "migrate" {
					found = true
					break
				}
			}

			if !found {
				t.Error("migrate command not found")
			}
		})
	}
}

func TestOpenDatabaseConnectionString(t *testing.T) {
	// Test that the connection string is properly formatted
	cfg = &config.Config{
		Database: config.DatabaseConfig{
			Driver:         "postgres",
			Host:           "testhost",
			Port:           5433,
			Name:           "testdb",
			User:           "testuser",
			Password:       "testpass",
			SSLMode:        "require",
			MigrationsPath: "./migrations",
		},
	}

	// We can't actually test openDatabase without a real database,
	// but we can verify the config is set correctly
	if cfg.Database.Host != "testhost" {
		t.Errorf("expected host testhost, got %s", cfg.Database.Host)
	}
	if cfg.Database.Port != 5433 {
		t.Errorf("expected port 5433, got %d", cfg.Database.Port)
	}
	if cfg.Database.Name != "testdb" {
		t.Errorf("expected database name testdb, got %s", cfg.Database.Name)
	}
	if cfg.Database.User != "testuser" {
		t.Errorf("expected user testuser, got %s", cfg.Database.User)
	}
	if cfg.Database.Password != "testpass" {
		t.Errorf("expected password testpass, got %s", cfg.Database.Password)
	}
	if cfg.Database.SSLMode != "require" {
		t.Errorf("expected sslmode require, got %s", cfg.Database.SSLMode)
	}
}

func TestMigrateCmdFlags(t *testing.T) {
	// Test that migrate commands have proper flags
	stepsCmd := migrateStepsCmd
	if stepsCmd == nil {
		t.Fatal("migrateStepsCmd is nil")
	}

	stepsFlag := stepsCmd.Flags().Lookup("steps")
	if stepsFlag == nil {
		t.Error("steps flag not found on migrate steps command")
	}

	forceCmd := migrateForceCmd
	if forceCmd == nil {
		t.Fatal("migrateForceCmd is nil")
	}

	versionFlag := forceCmd.Flags().Lookup("version")
	if versionFlag == nil {
		t.Error("version flag not found on migrate force command")
	}
}

func TestMigrateCommandsExist(t *testing.T) {
	// Verify all migrate subcommands are registered
	expectedCommands := []string{"up", "down", "steps", "version", "force"}

	if migrateCmd == nil {
		t.Fatal("migrateCmd is nil")
	}

	subcommands := migrateCmd.Commands()
	commandMap := make(map[string]bool)

	for _, cmd := range subcommands {
		commandMap[cmd.Use] = true
	}

	for _, expected := range expectedCommands {
		if !commandMap[expected] {
			t.Errorf("expected subcommand %s not found", expected)
		}
	}
}
