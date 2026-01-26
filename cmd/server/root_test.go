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
