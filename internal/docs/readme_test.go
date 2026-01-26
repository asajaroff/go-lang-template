package docs

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadmeExists(t *testing.T) {
	// Get the project root (two directories up from internal/docs)
	projectRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("Failed to get project root: %v", err)
	}

	readmePath := filepath.Join(projectRoot, "README.md")
	if _, err := os.Stat(readmePath); os.IsNotExist(err) {
		t.Errorf("README.md does not exist at %s", readmePath)
	}
}

func TestReadmeContent(t *testing.T) {
	// Get the project root
	projectRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("Failed to get project root: %v", err)
	}

	readmePath := filepath.Join(projectRoot, "README.md")
	content, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("Failed to read README.md: %v", err)
	}

	readmeText := string(content)

	// Test for essential sections
	requiredSections := []string{
		"# Go Lang Template",
		"## Features",
		"## Prerequisites",
		"## Quick Start",
		"## Configuration",
		"## API Endpoints",
		"## Development",
		"## Contributing",
	}

	for _, section := range requiredSections {
		if !strings.Contains(readmeText, section) {
			t.Errorf("README.md is missing required section: %s", section)
		}
	}
}

func TestReadmeHasSetupInstructions(t *testing.T) {
	projectRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("Failed to get project root: %v", err)
	}

	readmePath := filepath.Join(projectRoot, "README.md")
	content, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("Failed to read README.md: %v", err)
	}

	readmeText := string(content)

	// Test for setup instructions
	setupInstructions := []string{
		"git clone",
		"go mod download",
		"make install-tools",
		"make build",
		"./bin/server",
	}

	for _, instruction := range setupInstructions {
		if !strings.Contains(readmeText, instruction) {
			t.Errorf("README.md is missing setup instruction: %s", instruction)
		}
	}
}

func TestReadmeHasConfigurationExamples(t *testing.T) {
	projectRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("Failed to get project root: %v", err)
	}

	readmePath := filepath.Join(projectRoot, "README.md")
	content, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("Failed to read README.md: %v", err)
	}

	readmeText := string(content)

	// Test for configuration examples
	configExamples := []string{
		"--server-port",
		"--db-enabled",
		"APP_SERVER_PORT",
		"APP_DB_ENABLED",
		"config.yaml",
	}

	for _, example := range configExamples {
		if !strings.Contains(readmeText, example) {
			t.Errorf("README.md is missing configuration example: %s", example)
		}
	}
}

func TestReadmeHasHealthCheckEndpoints(t *testing.T) {
	projectRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("Failed to get project root: %v", err)
	}

	readmePath := filepath.Join(projectRoot, "README.md")
	content, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("Failed to read README.md: %v", err)
	}

	readmeText := string(content)

	// Test for health check endpoints
	healthCheckInfo := []string{
		"/healthz/liveness",
		"/healthz/readiness",
		"200 OK",
	}

	for _, info := range healthCheckInfo {
		if !strings.Contains(readmeText, info) {
			t.Errorf("README.md is missing health check information: %s", info)
		}
	}
}

func TestReadmeHasDatabaseMigrationInfo(t *testing.T) {
	projectRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("Failed to get project root: %v", err)
	}

	readmePath := filepath.Join(projectRoot, "README.md")
	content, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("Failed to read README.md: %v", err)
	}

	readmeText := string(content)

	// Test for migration commands
	migrationCommands := []string{
		"migrate up",
		"migrate down",
		"migrate steps",
		"migrate version",
	}

	for _, cmd := range migrationCommands {
		if !strings.Contains(readmeText, cmd) {
			t.Errorf("README.md is missing migration command: %s", cmd)
		}
	}
}

func TestReadmeHasDevelopmentTools(t *testing.T) {
	projectRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("Failed to get project root: %v", err)
	}

	readmePath := filepath.Join(projectRoot, "README.md")
	content, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("Failed to read README.md: %v", err)
	}

	readmeText := string(content)

	// Test for development tools
	devTools := []string{
		"make test",
		"make lint",
		"make build",
		"make swagger",
	}

	for _, tool := range devTools {
		if !strings.Contains(readmeText, tool) {
			t.Errorf("README.md is missing development tool: %s", tool)
		}
	}
}

func TestReadmeHasSwaggerDocumentation(t *testing.T) {
	projectRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("Failed to get project root: %v", err)
	}

	readmePath := filepath.Join(projectRoot, "README.md")
	content, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("Failed to read README.md: %v", err)
	}

	readmeText := string(content)

	// Test for Swagger documentation
	if !strings.Contains(readmeText, "/swagger/index.html") {
		t.Error("README.md is missing Swagger documentation endpoint")
	}
}

func TestReadmeHasProjectStructure(t *testing.T) {
	projectRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("Failed to get project root: %v", err)
	}

	readmePath := filepath.Join(projectRoot, "README.md")
	content, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("Failed to read README.md: %v", err)
	}

	readmeText := string(content)

	// Test for project structure
	structureElements := []string{
		"cmd/",
		"internal/",
		"migrations/",
		"docs/",
	}

	for _, element := range structureElements {
		if !strings.Contains(readmeText, element) {
			t.Errorf("README.md is missing project structure element: %s", element)
		}
	}
}

func TestReadmeHasDeploymentInstructions(t *testing.T) {
	projectRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("Failed to get project root: %v", err)
	}

	readmePath := filepath.Join(projectRoot, "README.md")
	content, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("Failed to read README.md: %v", err)
	}

	readmeText := string(content)

	// Test for deployment instructions
	deploymentInfo := []string{
		"Docker",
		"Kubernetes",
		"docker build",
		"livenessProbe",
		"readinessProbe",
	}

	for _, info := range deploymentInfo {
		if !strings.Contains(readmeText, info) {
			t.Errorf("README.md is missing deployment information: %s", info)
		}
	}
}

func TestReadmeNotEmpty(t *testing.T) {
	projectRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("Failed to get project root: %v", err)
	}

	readmePath := filepath.Join(projectRoot, "README.md")
	content, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("Failed to read README.md: %v", err)
	}

	if len(content) < 100 {
		t.Error("README.md appears to be too short or empty")
	}

	// Check for minimum number of lines
	lines := strings.Split(string(content), "\n")
	if len(lines) < 50 {
		t.Errorf("README.md should have at least 50 lines, found %d", len(lines))
	}
}

func TestReadmeHasContributingLink(t *testing.T) {
	projectRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("Failed to get project root: %v", err)
	}

	readmePath := filepath.Join(projectRoot, "README.md")
	content, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("Failed to read README.md: %v", err)
	}

	readmeText := string(content)

	if !strings.Contains(readmeText, "CONTRIBUTING.md") {
		t.Error("README.md should link to CONTRIBUTING.md")
	}
}
