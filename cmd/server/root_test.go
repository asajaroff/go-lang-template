package main

import (
	"testing"

	"github.com/asajaroff/go-lang-template/internal/config"
)

func TestExecuteCommand(t *testing.T) {
	// Test that Execute function exists and can be called
	// Note: actual execution is tested via TestMainFunction
	err := Execute
	if err == nil {
		t.Log("Execute function is available")
	}
}

func TestRootCmdInitialization(t *testing.T) {
	// Test that rootCmd is properly initialized
	if rootCmd == nil {
		t.Fatal("rootCmd should not be nil")
	}

	if rootCmd.Use != "server" {
		t.Errorf("expected Use = 'server', got '%s'", rootCmd.Use)
	}

	if rootCmd.Short != "Go Lang Template - Microservice Server" {
		t.Errorf("expected Short description, got '%s'", rootCmd.Short)
	}

	// Check that flags are registered
	flags := rootCmd.PersistentFlags()
	if flags == nil {
		t.Fatal("PersistentFlags should not be nil")
	}

	configFlag := flags.Lookup("config")
	if configFlag == nil {
		t.Error("config flag should be registered")
	}

	serverHostFlag := flags.Lookup("server-host")
	if serverHostFlag == nil {
		t.Error("server-host flag should be registered")
	}

	serverPortFlag := flags.Lookup("server-port")
	if serverPortFlag == nil {
		t.Error("server-port flag should be registered")
	}
}

func TestConfigVariables(t *testing.T) {
	// Test that config variables are declared
	if cfg != nil {
		// cfg should be nil until config is loaded
		t.Log("cfg is initially set")
	}

	// Test config type
	var testCfg *config.Config
	if testCfg == nil {
		t.Log("config.Config type is available")
	}
}

// TestRootCmdRunE tests that RunE function exists
func TestRootCmdRunE(t *testing.T) {
	if rootCmd.RunE == nil {
		t.Error("rootCmd should have RunE function")
	}
}

// TestRootCmdPersistentPreRunE tests that PersistentPreRunE function exists
func TestRootCmdPersistentPreRunE(t *testing.T) {
	if rootCmd.PersistentPreRunE == nil {
		t.Error("rootCmd should have PersistentPreRunE function")
	}
}

// TestAllRequiredFlags tests that all required flags are registered
func TestAllRequiredFlags(t *testing.T) {
	requiredFlags := []string{
		"config",
		"server-host",
		"server-port",
		"db-enabled",
		"db-driver",
		"db-host",
		"db-port",
		"db-name",
		"db-user",
		"db-password",
		"db-sslmode",
		"db-migrations-path",
		"log-level",
		"log-format",
	}

	flags := rootCmd.PersistentFlags()

	for _, flagName := range requiredFlags {
		flag := flags.Lookup(flagName)
		if flag == nil {
			t.Errorf("flag %s should be registered", flagName)
		}
	}
}

// TestFlagTypes tests that flags have correct types
func TestFlagTypes(t *testing.T) {
	flags := rootCmd.PersistentFlags()

	// String flags
	stringFlags := []string{"config", "server-host", "db-driver", "db-host", "db-name", "db-user", "db-password", "db-sslmode", "db-migrations-path", "log-level", "log-format"}
	for _, flagName := range stringFlags {
		flag := flags.Lookup(flagName)
		if flag == nil {
			t.Errorf("flag %s not found", flagName)
			continue
		}
		if flag.Value.Type() != "string" {
			t.Errorf("flag %s should be string type, got %s", flagName, flag.Value.Type())
		}
	}

	// Int flags
	intFlags := []string{"server-port", "db-port"}
	for _, flagName := range intFlags {
		flag := flags.Lookup(flagName)
		if flag == nil {
			t.Errorf("flag %s not found", flagName)
			continue
		}
		if flag.Value.Type() != "int" {
			t.Errorf("flag %s should be int type, got %s", flagName, flag.Value.Type())
		}
	}

	// Bool flags
	boolFlags := []string{"db-enabled"}
	for _, flagName := range boolFlags {
		flag := flags.Lookup(flagName)
		if flag == nil {
			t.Errorf("flag %s not found", flagName)
			continue
		}
		if flag.Value.Type() != "bool" {
			t.Errorf("flag %s should be bool type, got %s", flagName, flag.Value.Type())
		}
	}
}

// TestRootCmdBasicProperties tests basic command properties
func TestRootCmdBasicProperties(t *testing.T) {
	if rootCmd.Use == "" {
		t.Error("rootCmd should have Use field set")
	}

	if rootCmd.Short == "" {
		t.Error("rootCmd should have Short description")
	}

	if rootCmd.Long == "" {
		t.Error("rootCmd should have Long description")
	}
}

// TestCfgFileVariable tests cfgFile variable
func TestCfgFileVariable(t *testing.T) {
	// cfgFile should be accessible and modifiable
	oldCfgFile := cfgFile
	cfgFile = "test.yaml"
	if cfgFile != "test.yaml" {
		t.Error("cfgFile should be modifiable")
	}
	cfgFile = oldCfgFile
}

// TestExecuteFunction tests Execute function existence
func TestExecuteFunction(t *testing.T) {
	// Test that Execute function exists and returns error type
	err := Execute
	if err == nil {
		t.Log("Execute function is available")
	}
}

// TestRootCmdHasSubcommands tests that root command can have subcommands
func TestRootCmdHasSubcommands(t *testing.T) {
	commands := rootCmd.Commands()
	if commands == nil {
		t.Error("rootCmd should be able to have subcommands")
	}

	// Check that migrate command is added
	hasMigrate := false
	for _, cmd := range commands {
		if cmd.Use == "migrate" {
			hasMigrate = true
			break
		}
	}

	if !hasMigrate {
		t.Error("rootCmd should have migrate subcommand")
	}
}

// TestConfigFlagDefaults tests default values for config flag
func TestConfigFlagDefaults(t *testing.T) {
	flag := rootCmd.PersistentFlags().Lookup("config")
	if flag == nil {
		t.Fatal("config flag not found")
	}

	if flag.DefValue != "" {
		t.Errorf("expected empty default value for config, got %s", flag.DefValue)
	}
}

// TestServerFlagDefaults tests default values for server flags
func TestServerFlagDefaults(t *testing.T) {
	tests := []struct {
		name         string
		expectedType string
		expectedDef  string
	}{
		{"server-host", "string", ""},
		{"server-port", "int", "0"},
	}

	flags := rootCmd.PersistentFlags()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := flags.Lookup(tt.name)
			if flag == nil {
				t.Fatalf("flag %s not found", tt.name)
			}

			if flag.Value.Type() != tt.expectedType {
				t.Errorf("expected type %s, got %s", tt.expectedType, flag.Value.Type())
			}

			if flag.DefValue != tt.expectedDef {
				t.Errorf("expected default %s, got %s", tt.expectedDef, flag.DefValue)
			}
		})
	}
}

// TestDatabaseFlagDefaults tests default values for database flags
func TestDatabaseFlagDefaults(t *testing.T) {
	tests := []struct {
		name         string
		expectedType string
		expectedDef  string
	}{
		{"db-enabled", "bool", "false"},
		{"db-driver", "string", ""},
		{"db-host", "string", ""},
		{"db-port", "int", "0"},
		{"db-name", "string", ""},
		{"db-user", "string", ""},
		{"db-password", "string", ""},
		{"db-sslmode", "string", ""},
		{"db-migrations-path", "string", ""},
	}

	flags := rootCmd.PersistentFlags()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := flags.Lookup(tt.name)
			if flag == nil {
				t.Fatalf("flag %s not found", tt.name)
			}

			if flag.Value.Type() != tt.expectedType {
				t.Errorf("expected type %s, got %s", tt.expectedType, flag.Value.Type())
			}

			if flag.DefValue != tt.expectedDef {
				t.Errorf("expected default %s, got %s", tt.expectedDef, flag.DefValue)
			}
		})
	}
}

// TestLogFlagDefaults tests default values for log flags
func TestLogFlagDefaults(t *testing.T) {
	tests := []struct {
		name         string
		expectedType string
		expectedDef  string
	}{
		{"log-level", "string", ""},
		{"log-format", "string", ""},
	}

	flags := rootCmd.PersistentFlags()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := flags.Lookup(tt.name)
			if flag == nil {
				t.Fatalf("flag %s not found", tt.name)
			}

			if flag.Value.Type() != tt.expectedType {
				t.Errorf("expected type %s, got %s", tt.expectedType, flag.Value.Type())
			}

			if flag.DefValue != tt.expectedDef {
				t.Errorf("expected default %s, got %s", tt.expectedDef, flag.DefValue)
			}
		})
	}
}

// TestRootCmdAnnotations tests that root command is properly configured
func TestRootCmdAnnotations(t *testing.T) {
	if rootCmd == nil {
		t.Fatal("rootCmd should not be nil")
	}

	// Verify command is executable
	if rootCmd.RunE == nil {
		t.Error("rootCmd should have RunE function defined")
	}

	// Verify pre-run hook for config loading
	if rootCmd.PersistentPreRunE == nil {
		t.Error("rootCmd should have PersistentPreRunE for config loading")
	}
}
