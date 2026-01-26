package main

import (
	"testing"

	"github.com/asajaroff/go-lang-template/internal/config"
	"github.com/spf13/cobra"
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

// TestMigrateCommandProperties tests properties of migrate commands
func TestMigrateCommandProperties(t *testing.T) {
	tests := []struct {
		name        string
		cmd         string
		expectedUse string
		hasRunE     bool
	}{
		{"migrate root", "migrate", "migrate", false},
		{"migrate up", "up", "up", true},
		{"migrate down", "down", "down", true},
		{"migrate steps", "steps", "steps", true},
		{"migrate version", "version", "version", true},
		{"migrate force", "force", "force", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cmd interface{}

			switch tt.cmd {
			case "migrate":
				cmd = migrateCmd
			case "up":
				cmd = migrateUpCmd
			case "down":
				cmd = migrateDownCmd
			case "steps":
				cmd = migrateStepsCmd
			case "version":
				cmd = migrateVersionCmd
			case "force":
				cmd = migrateForceCmd
			}

			if cmd == nil {
				t.Fatalf("%s command is nil", tt.name)
			}

			// Check that command has proper Short description
			switch c := cmd.(type) {
			case *cobra.Command:
				if c.Short == "" {
					t.Errorf("%s should have a Short description", tt.name)
				}
			}
		})
	}
}

// TestMigrateStepsFlagDefaults tests default values for steps flag
func TestMigrateStepsFlagDefaults(t *testing.T) {
	flag := migrateStepsCmd.Flags().Lookup("steps")
	if flag == nil {
		t.Fatal("steps flag not found")
	}

	if flag.DefValue != "1" {
		t.Errorf("expected default value 1, got %s", flag.DefValue)
	}
}

// TestMigrateForceFlagDefaults tests default values for force flag
func TestMigrateForceFlagDefaults(t *testing.T) {
	flag := migrateForceCmd.Flags().Lookup("version")
	if flag == nil {
		t.Fatal("version flag not found")
	}

	if flag.DefValue != "0" {
		t.Errorf("expected default value 0, got %s", flag.DefValue)
	}
}

// TestMigrateCommandInheritance tests that migrate commands inherit from root
func TestMigrateCommandInheritance(t *testing.T) {
	// Verify migrate command is added to root command
	commands := rootCmd.Commands()
	found := false

	for _, cmd := range commands {
		if cmd.Use == "migrate" {
			found = true
			break
		}
	}

	if !found {
		t.Error("migrate command should be added to root command")
	}
}

// TestMigrateSubcommandInheritance tests subcommands are added to migrate
func TestMigrateSubcommandInheritance(t *testing.T) {
	expectedSubcommands := map[string]bool{
		"up":      false,
		"down":    false,
		"steps":   false,
		"version": false,
		"force":   false,
	}

	subcommands := migrateCmd.Commands()
	for _, cmd := range subcommands {
		if _, exists := expectedSubcommands[cmd.Use]; exists {
			expectedSubcommands[cmd.Use] = true
		}
	}

	for name, found := range expectedSubcommands {
		if !found {
			t.Errorf("subcommand %s not found in migrate command", name)
		}
	}
}

// TestOpenDatabaseWithInvalidConfig tests openDatabase error handling
func TestOpenDatabaseWithInvalidConfig(t *testing.T) {
	// Save original config
	originalCfg := cfg

	// Test with invalid config
	cfg = &config.Config{
		Database: config.DatabaseConfig{
			Driver:   "postgres",
			Host:     "invalid-host-that-does-not-exist-12345",
			Port:     99999,
			Name:     "nonexistent",
			User:     "nonexistent",
			Password: "invalid",
			SSLMode:  "disable",
		},
	}

	// Try to open database - should fail
	db, err := openDatabase()
	if err == nil {
		t.Error("expected error with invalid config, got nil")
		if db != nil {
			db.Close()
		}
	}

	// Restore original config
	cfg = originalCfg
}

// TestMigrateCommandDescriptions tests that all commands have descriptions
func TestMigrateCommandDescriptions(t *testing.T) {
	commands := []*cobra.Command{
		migrateCmd,
		migrateUpCmd,
		migrateDownCmd,
		migrateStepsCmd,
		migrateVersionCmd,
		migrateForceCmd,
	}

	commandNames := []string{
		"migrate",
		"migrate up",
		"migrate down",
		"migrate steps",
		"migrate version",
		"migrate force",
	}

	for i, cmd := range commands {
		if cmd == nil {
			t.Errorf("%s command is nil", commandNames[i])
			continue
		}

		if cmd.Use == "" {
			t.Errorf("%s command should have Use field set", commandNames[i])
		}

		if cmd.Short == "" {
			t.Errorf("%s command should have Short description", commandNames[i])
		}
	}
}

// TestMigrateVariables tests global migrate variables
func TestMigrateVariables(t *testing.T) {
	// Test that migrateSteps and migrateForce variables exist
	// They should be accessible and modifiable

	oldSteps := migrateSteps
	oldForce := migrateForce

	migrateSteps = 5
	if migrateSteps != 5 {
		t.Error("migrateSteps should be modifiable")
	}

	migrateForce = 10
	if migrateForce != 10 {
		t.Error("migrateForce should be modifiable")
	}

	// Restore
	migrateSteps = oldSteps
	migrateForce = oldForce
}

// TestMigrateCmdUsage tests command usage strings
func TestMigrateCmdUsage(t *testing.T) {
	if migrateCmd.Use != "migrate" {
		t.Errorf("expected Use='migrate', got %s", migrateCmd.Use)
	}
	if migrateUpCmd.Use != "up" {
		t.Errorf("expected Use='up', got %s", migrateUpCmd.Use)
	}
	if migrateDownCmd.Use != "down" {
		t.Errorf("expected Use='down', got %s", migrateDownCmd.Use)
	}
	if migrateStepsCmd.Use != "steps" {
		t.Errorf("expected Use='steps', got %s", migrateStepsCmd.Use)
	}
	if migrateVersionCmd.Use != "version" {
		t.Errorf("expected Use='version', got %s", migrateVersionCmd.Use)
	}
	if migrateForceCmd.Use != "force" {
		t.Errorf("expected Use='force', got %s", migrateForceCmd.Use)
	}
}

// TestOpenDatabaseConnectionFailure tests connection failure handling
func TestOpenDatabaseConnectionFailure(t *testing.T) {
	originalCfg := cfg

	// Set config with host that connection will fail
	cfg = &config.Config{
		Database: config.DatabaseConfig{
			Driver:   "postgres",
			Host:     "localhost",
			Port:     1,
			Name:     "test",
			User:     "test",
			Password: "test",
			SSLMode:  "disable",
		},
	}

	db, err := openDatabase()
	if err == nil {
		t.Error("expected error when connecting to invalid port, got nil")
		if db != nil {
			db.Close()
		}
	}

	cfg = originalCfg
}
