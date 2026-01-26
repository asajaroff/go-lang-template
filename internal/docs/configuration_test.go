package docs

import (
	"os"
	"strings"
	"testing"
)

func TestConfigurationDocumentationExists(t *testing.T) {
	_, err := os.Stat("../../docs/CONFIGURATION.md")
	if err != nil {
		t.Errorf("CONFIGURATION.md file should exist: %v", err)
	}
}

func TestConfigurationDocumentationContent(t *testing.T) {
	content, err := os.ReadFile("../../docs/CONFIGURATION.md")
	if err != nil {
		t.Fatalf("Failed to read CONFIGURATION.md: %v", err)
	}

	contentStr := string(content)

	// Required sections
	requiredSections := []string{
		"# Configuration Guide",
		"## Table of Contents",
		"## Configuration Sources",
		"## Priority Order",
		"## Server Configuration",
		"## Database Configuration",
		"## Logging Configuration",
		"## Configuration File Format",
		"## Environment Variables Reference",
		"## CLI Flags Reference",
		"## Configuration Examples",
	}

	for _, section := range requiredSections {
		if !strings.Contains(contentStr, section) {
			t.Errorf("CONFIGURATION.md should contain section: %s", section)
		}
	}
}

func TestConfigurationDocumentationServerOptions(t *testing.T) {
	content, err := os.ReadFile("../../docs/CONFIGURATION.md")
	if err != nil {
		t.Fatalf("Failed to read CONFIGURATION.md: %v", err)
	}

	contentStr := string(content)

	// Server configuration options
	serverOptions := []string{
		"server.host",
		"server.port",
		"--server-host",
		"--server-port",
		"APP_SERVER_HOST",
		"APP_SERVER_PORT",
	}

	for _, option := range serverOptions {
		if !strings.Contains(contentStr, option) {
			t.Errorf("CONFIGURATION.md should document server option: %s", option)
		}
	}

	// Server default values
	if !strings.Contains(contentStr, "0.0.0.0") {
		t.Error("CONFIGURATION.md should document default server host")
	}
	if !strings.Contains(contentStr, "8080") {
		t.Error("CONFIGURATION.md should document default server port")
	}
}

func TestConfigurationDocumentationDatabaseOptions(t *testing.T) {
	content, err := os.ReadFile("../../docs/CONFIGURATION.md")
	if err != nil {
		t.Fatalf("Failed to read CONFIGURATION.md: %v", err)
	}

	contentStr := string(content)

	// Database configuration options
	databaseOptions := []string{
		"database.enabled",
		"database.driver",
		"database.host",
		"database.port",
		"database.name",
		"database.user",
		"database.password",
		"database.sslmode",
		"database.migrationspath",
		"--db-enabled",
		"--db-driver",
		"--db-host",
		"--db-port",
		"--db-name",
		"--db-user",
		"--db-password",
		"--db-sslmode",
		"--db-migrations-path",
		"APP_DATABASE_ENABLED",
		"APP_DATABASE_DRIVER",
		"APP_DATABASE_HOST",
		"APP_DATABASE_PORT",
		"APP_DATABASE_NAME",
		"APP_DATABASE_USER",
		"APP_DATABASE_PASSWORD",
		"APP_DATABASE_SSLMODE",
		"APP_DATABASE_MIGRATIONSPATH",
	}

	for _, option := range databaseOptions {
		if !strings.Contains(contentStr, option) {
			t.Errorf("CONFIGURATION.md should document database option: %s", option)
		}
	}

	// Database default values
	if !strings.Contains(contentStr, "postgres") {
		t.Error("CONFIGURATION.md should document default database driver")
	}
	if !strings.Contains(contentStr, "localhost") {
		t.Error("CONFIGURATION.md should document default database host")
	}
	if !strings.Contains(contentStr, "5432") {
		t.Error("CONFIGURATION.md should document default database port")
	}
}

func TestConfigurationDocumentationSSLModes(t *testing.T) {
	content, err := os.ReadFile("../../docs/CONFIGURATION.md")
	if err != nil {
		t.Fatalf("Failed to read CONFIGURATION.md: %v", err)
	}

	contentStr := string(content)

	// SSL modes
	sslModes := []string{
		"disable",
		"require",
		"verify-ca",
		"verify-full",
	}

	for _, mode := range sslModes {
		if !strings.Contains(contentStr, mode) {
			t.Errorf("CONFIGURATION.md should document SSL mode: %s", mode)
		}
	}

	// SSL mode descriptions
	if !strings.Contains(contentStr, "SSL Modes") {
		t.Error("CONFIGURATION.md should have SSL Modes section")
	}
}

func TestConfigurationDocumentationLoggingOptions(t *testing.T) {
	content, err := os.ReadFile("../../docs/CONFIGURATION.md")
	if err != nil {
		t.Fatalf("Failed to read CONFIGURATION.md: %v", err)
	}

	contentStr := string(content)

	// Logging configuration options
	loggingOptions := []string{
		"log.level",
		"log.format",
		"--log-level",
		"--log-format",
		"APP_LOG_LEVEL",
		"APP_LOG_FORMAT",
	}

	for _, option := range loggingOptions {
		if !strings.Contains(contentStr, option) {
			t.Errorf("CONFIGURATION.md should document logging option: %s", option)
		}
	}

	// Log levels
	logLevels := []string{
		"debug",
		"info",
		"warn",
		"error",
	}

	for _, level := range logLevels {
		if !strings.Contains(contentStr, level) {
			t.Errorf("CONFIGURATION.md should document log level: %s", level)
		}
	}

	// Log formats
	logFormats := []string{
		"json",
		"text",
	}

	for _, format := range logFormats {
		if !strings.Contains(contentStr, format) {
			t.Errorf("CONFIGURATION.md should document log format: %s", format)
		}
	}

	// Log output routing
	if !strings.Contains(contentStr, "stdout") {
		t.Error("CONFIGURATION.md should document stdout log routing")
	}
	if !strings.Contains(contentStr, "stderr") {
		t.Error("CONFIGURATION.md should document stderr log routing")
	}
}

func TestConfigurationDocumentationPriorityOrder(t *testing.T) {
	content, err := os.ReadFile("../../docs/CONFIGURATION.md")
	if err != nil {
		t.Fatalf("Failed to read CONFIGURATION.md: %v", err)
	}

	contentStr := string(content)

	// Priority order documentation
	priorityElements := []string{
		"CLI Flags",
		"Environment Variables",
		"Configuration File",
		"Default Values",
	}

	for _, element := range priorityElements {
		if !strings.Contains(contentStr, element) {
			t.Errorf("CONFIGURATION.md should document priority element: %s", element)
		}
	}

	// Should mention priority/precedence
	if !strings.Contains(contentStr, "priority") && !strings.Contains(contentStr, "Priority") {
		t.Error("CONFIGURATION.md should explain configuration priority")
	}
}

func TestConfigurationDocumentationExamples(t *testing.T) {
	content, err := os.ReadFile("../../docs/CONFIGURATION.md")
	if err != nil {
		t.Fatalf("Failed to read CONFIGURATION.md: %v", err)
	}

	contentStr := string(content)

	// Should contain practical examples
	examples := []string{
		"Development Environment",
		"Production Environment",
		"Docker Environment",
		"Kubernetes Environment",
	}

	for _, example := range examples {
		if !strings.Contains(contentStr, example) {
			t.Errorf("CONFIGURATION.md should contain example: %s", example)
		}
	}

	// Should contain code blocks
	if !strings.Contains(contentStr, "```bash") {
		t.Error("CONFIGURATION.md should contain bash code examples")
	}
	if !strings.Contains(contentStr, "```yaml") {
		t.Error("CONFIGURATION.md should contain YAML code examples")
	}
}

func TestConfigurationDocumentationYAMLFormat(t *testing.T) {
	content, err := os.ReadFile("../../docs/CONFIGURATION.md")
	if err != nil {
		t.Fatalf("Failed to read CONFIGURATION.md: %v", err)
	}

	contentStr := string(content)

	// Should document YAML file format
	if !strings.Contains(contentStr, "config.yaml") {
		t.Error("CONFIGURATION.md should reference config.yaml file")
	}

	// Should show complete YAML example
	yamlElements := []string{
		"server:",
		"database:",
		"log:",
	}

	for _, element := range yamlElements {
		if !strings.Contains(contentStr, element) {
			t.Errorf("CONFIGURATION.md YAML example should contain: %s", element)
		}
	}
}

func TestConfigurationDocumentationConfigFileLocations(t *testing.T) {
	content, err := os.ReadFile("../../docs/CONFIGURATION.md")
	if err != nil {
		t.Fatalf("Failed to read CONFIGURATION.md: %v", err)
	}

	contentStr := string(content)

	// Should document config file search locations
	locations := []string{
		"./config.yaml",
		"./config/config.yaml",
		"/etc/go-lang-template/config.yaml",
	}

	for _, location := range locations {
		if !strings.Contains(contentStr, location) {
			t.Errorf("CONFIGURATION.md should document config file location: %s", location)
		}
	}
}

func TestConfigurationDocumentationSecuritySection(t *testing.T) {
	content, err := os.ReadFile("../../docs/CONFIGURATION.md")
	if err != nil {
		t.Fatalf("Failed to read CONFIGURATION.md: %v", err)
	}

	contentStr := string(content)

	// Should have security considerations
	if !strings.Contains(contentStr, "Security") {
		t.Error("CONFIGURATION.md should have security section")
	}

	// Security topics
	securityTopics := []string{
		"password",
		"secret",
		"SSL",
	}

	securityMentionCount := 0
	for _, topic := range securityTopics {
		if strings.Contains(strings.ToLower(contentStr), strings.ToLower(topic)) {
			securityMentionCount++
		}
	}

	if securityMentionCount < 2 {
		t.Error("CONFIGURATION.md should discuss security considerations")
	}
}

func TestConfigurationDocumentationTroubleshooting(t *testing.T) {
	content, err := os.ReadFile("../../docs/CONFIGURATION.md")
	if err != nil {
		t.Fatalf("Failed to read CONFIGURATION.md: %v", err)
	}

	contentStr := string(content)

	// Should have troubleshooting section
	if !strings.Contains(contentStr, "Troubleshooting") {
		t.Error("CONFIGURATION.md should have troubleshooting section")
	}

	// Common issues
	issues := []string{
		"Configuration Not Loading",
		"Database Connection",
	}

	for _, issue := range issues {
		if !strings.Contains(contentStr, issue) {
			t.Errorf("CONFIGURATION.md should address issue: %s", issue)
		}
	}
}

func TestConfigurationDocumentationEnvPrefix(t *testing.T) {
	content, err := os.ReadFile("../../docs/CONFIGURATION.md")
	if err != nil {
		t.Fatalf("Failed to read CONFIGURATION.md: %v", err)
	}

	contentStr := string(content)

	// Should document APP_ prefix for environment variables
	if !strings.Contains(contentStr, "APP_") {
		t.Error("CONFIGURATION.md should document APP_ prefix for environment variables")
	}

	// Should explain the prefix usage
	if !strings.Contains(strings.ToLower(contentStr), "prefix") {
		t.Error("CONFIGURATION.md should explain environment variable prefix")
	}
}

func TestConfigurationDocumentationCompleteReference(t *testing.T) {
	content, err := os.ReadFile("../../docs/CONFIGURATION.md")
	if err != nil {
		t.Fatalf("Failed to read CONFIGURATION.md: %v", err)
	}

	contentStr := string(content)

	// Should have complete reference sections
	references := []string{
		"Environment Variables Reference",
		"CLI Flags Reference",
	}

	for _, ref := range references {
		if !strings.Contains(contentStr, ref) {
			t.Errorf("CONFIGURATION.md should have section: %s", ref)
		}
	}

	// Should list all available options
	if strings.Count(contentStr, "APP_") < 10 {
		t.Error("CONFIGURATION.md should comprehensively list environment variables")
	}

	if strings.Count(contentStr, "--") < 15 {
		t.Error("CONFIGURATION.md should comprehensively list CLI flags")
	}
}

func TestConfigurationDocumentationTables(t *testing.T) {
	content, err := os.ReadFile("../../docs/CONFIGURATION.md")
	if err != nil {
		t.Fatalf("Failed to read CONFIGURATION.md: %v", err)
	}

	contentStr := string(content)

	// Should contain markdown tables for better readability
	if !strings.Contains(contentStr, "| Option |") {
		t.Error("CONFIGURATION.md should use tables to document options")
	}

	// Should have table headers
	tableHeaders := []string{
		"| Option |",
		"| Type |",
		"| Default |",
		"| Description |",
	}

	for _, header := range tableHeaders {
		if !strings.Contains(contentStr, header) {
			t.Errorf("CONFIGURATION.md tables should have header: %s", header)
		}
	}
}

func TestConfigurationDocumentationDocker(t *testing.T) {
	content, err := os.ReadFile("../../docs/CONFIGURATION.md")
	if err != nil {
		t.Fatalf("Failed to read CONFIGURATION.md: %v", err)
	}

	contentStr := string(content)

	// Should have Docker-specific examples
	if !strings.Contains(contentStr, "docker run") {
		t.Error("CONFIGURATION.md should contain Docker run examples")
	}

	// Should show how to pass environment variables to Docker
	if !strings.Contains(contentStr, "-e APP_") {
		t.Error("CONFIGURATION.md should show how to pass environment variables to Docker")
	}
}

func TestConfigurationDocumentationKubernetes(t *testing.T) {
	content, err := os.ReadFile("../../docs/CONFIGURATION.md")
	if err != nil {
		t.Fatalf("Failed to read CONFIGURATION.md: %v", err)
	}

	contentStr := string(content)

	// Should have Kubernetes-specific examples
	k8sElements := []string{
		"ConfigMap",
		"Secret",
		"Deployment",
	}

	for _, element := range k8sElements {
		if !strings.Contains(contentStr, element) {
			t.Errorf("CONFIGURATION.md should contain Kubernetes example with: %s", element)
		}
	}

	// Should show proper Kubernetes YAML
	if !strings.Contains(contentStr, "apiVersion:") {
		t.Error("CONFIGURATION.md should contain valid Kubernetes YAML examples")
	}
}

func TestConfigurationDocumentationReadability(t *testing.T) {
	content, err := os.ReadFile("../../docs/CONFIGURATION.md")
	if err != nil {
		t.Fatalf("Failed to read CONFIGURATION.md: %v", err)
	}

	contentStr := string(content)

	// Should be substantial documentation
	if len(contentStr) < 5000 {
		t.Error("CONFIGURATION.md should be comprehensive (at least 5000 characters)")
	}

	// Should have good structure with headers
	headerCount := strings.Count(contentStr, "##")
	if headerCount < 10 {
		t.Error("CONFIGURATION.md should have good section structure with multiple headers")
	}

	// Should have examples
	exampleCount := strings.Count(contentStr, "```")
	if exampleCount < 20 {
		t.Error("CONFIGURATION.md should have multiple code examples")
	}
}
