package docs

import (
	"os"
	"strings"
	"testing"
)

func TestArchitectureDocumentationExists(t *testing.T) {
	_, err := os.Stat("../../docs/ARCHITECTURE.md")
	if err != nil {
		t.Fatalf("ARCHITECTURE.md does not exist: %v", err)
	}
}

func TestArchitectureDocumentationContent(t *testing.T) {
	content, err := os.ReadFile("../../docs/ARCHITECTURE.md")
	if err != nil {
		t.Fatalf("Failed to read ARCHITECTURE.md: %v", err)
	}

	doc := string(content)

	t.Run("HasTableOfContents", func(t *testing.T) {
		if !strings.Contains(doc, "## Table of Contents") {
			t.Error("ARCHITECTURE.md missing Table of Contents")
		}
	})

	t.Run("HasSystemOverview", func(t *testing.T) {
		if !strings.Contains(doc, "## System Overview") {
			t.Error("ARCHITECTURE.md missing System Overview section")
		}
		if !strings.Contains(doc, "12-factor app") {
			t.Error("ARCHITECTURE.md should mention 12-factor app methodology")
		}
	})

	t.Run("HasHighLevelArchitecture", func(t *testing.T) {
		if !strings.Contains(doc, "## High-Level Architecture") {
			t.Error("ARCHITECTURE.md missing High-Level Architecture section")
		}
		// Check for ASCII diagram elements
		if !strings.Contains(doc, "┌") || !strings.Contains(doc, "└") || !strings.Contains(doc, "│") {
			t.Error("ARCHITECTURE.md should contain ASCII art diagrams")
		}
	})

	t.Run("HasComponentArchitecture", func(t *testing.T) {
		if !strings.Contains(doc, "## Component Architecture") {
			t.Error("ARCHITECTURE.md missing Component Architecture section")
		}
		if !strings.Contains(doc, "cmd/server/") {
			t.Error("ARCHITECTURE.md should document command layer")
		}
		if !strings.Contains(doc, "internal/") {
			t.Error("ARCHITECTURE.md should document internal layer")
		}
	})

	t.Run("HasDataFlow", func(t *testing.T) {
		if !strings.Contains(doc, "## Data Flow") {
			t.Error("ARCHITECTURE.md missing Data Flow section")
		}
		if !strings.Contains(doc, "Request Processing Flow") {
			t.Error("ARCHITECTURE.md should document request processing flow")
		}
		if !strings.Contains(doc, "Configuration Loading Flow") {
			t.Error("ARCHITECTURE.md should document configuration loading flow")
		}
	})

	t.Run("HasDeploymentArchitecture", func(t *testing.T) {
		if !strings.Contains(doc, "## Deployment Architecture") {
			t.Error("ARCHITECTURE.md missing Deployment Architecture section")
		}
		if !strings.Contains(doc, "Docker Deployment") {
			t.Error("ARCHITECTURE.md should document Docker deployment")
		}
		if !strings.Contains(doc, "Kubernetes Deployment") {
			t.Error("ARCHITECTURE.md should document Kubernetes deployment")
		}
	})

	t.Run("HasConfigurationFlow", func(t *testing.T) {
		if !strings.Contains(doc, "## Configuration Flow") {
			t.Error("ARCHITECTURE.md missing Configuration Flow section")
		}
		if !strings.Contains(doc, "CLI Flags") && !strings.Contains(doc, "Highest Priority") {
			t.Error("ARCHITECTURE.md should document configuration precedence")
		}
	})

	t.Run("HasHealthCheckFlow", func(t *testing.T) {
		if !strings.Contains(doc, "## Health Check Flow") {
			t.Error("ARCHITECTURE.md missing Health Check Flow section")
		}
		if !strings.Contains(doc, "/healthz/liveness") {
			t.Error("ARCHITECTURE.md should document liveness probe")
		}
		if !strings.Contains(doc, "/healthz/readiness") {
			t.Error("ARCHITECTURE.md should document readiness probe")
		}
	})

	t.Run("HasDatabaseMigrationFlow", func(t *testing.T) {
		if !strings.Contains(doc, "## Database Migration Flow") {
			t.Error("ARCHITECTURE.md missing Database Migration Flow section")
		}
		if !strings.Contains(doc, "golang-migrate") {
			t.Error("ARCHITECTURE.md should mention golang-migrate")
		}
		if !strings.Contains(doc, "schema_migrations") {
			t.Error("ARCHITECTURE.md should document schema_migrations table")
		}
	})
}

func TestArchitectureDocumentationComponents(t *testing.T) {
	content, err := os.ReadFile("../../docs/ARCHITECTURE.md")
	if err != nil {
		t.Fatalf("Failed to read ARCHITECTURE.md: %v", err)
	}

	doc := string(content)

	t.Run("DocumentsCommandLayer", func(t *testing.T) {
		components := []string{
			"main.go",
			"root.go",
			"migrate.go",
			"Cobra",
		}
		for _, comp := range components {
			if !strings.Contains(doc, comp) {
				t.Errorf("ARCHITECTURE.md should document %s in command layer", comp)
			}
		}
	})

	t.Run("DocumentsInternalPackages", func(t *testing.T) {
		packages := []string{
			"config/",
			"health/",
			"logger/",
			"database/",
			"migrate/",
		}
		for _, pkg := range packages {
			if !strings.Contains(doc, pkg) {
				t.Errorf("ARCHITECTURE.md should document internal package: %s", pkg)
			}
		}
	})

	t.Run("DocumentsConfigPackage", func(t *testing.T) {
		features := []string{
			"Viper",
			"Configuration struct",
			"environment variables",
		}
		for _, feature := range features {
			if !strings.Contains(doc, feature) {
				t.Errorf("ARCHITECTURE.md should document config feature: %s", feature)
			}
		}
	})

	t.Run("DocumentsHealthPackage", func(t *testing.T) {
		features := []string{
			"Checker interface",
			"Liveness probe",
			"Readiness probe",
		}
		for _, feature := range features {
			if !strings.Contains(doc, feature) {
				t.Errorf("ARCHITECTURE.md should document health feature: %s", feature)
			}
		}
	})

	t.Run("DocumentsLoggerPackage", func(t *testing.T) {
		features := []string{
			"log/slog",
			"JSON",
			"text format",
			"stdout",
			"stderr",
		}
		for _, feature := range features {
			if !strings.Contains(doc, feature) {
				t.Errorf("ARCHITECTURE.md should document logger feature: %s", feature)
			}
		}
	})

	t.Run("DocumentsDatabasePackage", func(t *testing.T) {
		features := []string{
			"Database connection management",
			"DSN",
			"sql.DB",
			"Connection pooling",
		}
		for _, feature := range features {
			if !strings.Contains(doc, feature) {
				t.Errorf("ARCHITECTURE.md should document database feature: %s", feature)
			}
		}
	})

	t.Run("DocumentsMigratePackage", func(t *testing.T) {
		features := []string{
			"golang-migrate",
			"up",
			"down",
			"steps",
		}
		for _, feature := range features {
			if !strings.Contains(doc, feature) {
				t.Errorf("ARCHITECTURE.md should document migrate feature: %s", feature)
			}
		}
	})
}

func TestArchitectureDocumentationDeployment(t *testing.T) {
	content, err := os.ReadFile("../../docs/ARCHITECTURE.md")
	if err != nil {
		t.Fatalf("Failed to read ARCHITECTURE.md: %v", err)
	}

	doc := string(content)

	t.Run("DocumentsDockerDeployment", func(t *testing.T) {
		elements := []string{
			"Docker Host",
			"Container",
			"port 8080",
		}
		for _, elem := range elements {
			if !strings.Contains(doc, elem) {
				t.Errorf("ARCHITECTURE.md should document Docker element: %s", elem)
			}
		}
	})

	t.Run("DocumentsKubernetesDeployment", func(t *testing.T) {
		elements := []string{
			"Kubernetes Cluster",
			"Ingress",
			"Service",
			"Deployment",
			"Pod",
			"ConfigMap",
			"Secret",
			"StatefulSet",
		}
		for _, elem := range elements {
			if !strings.Contains(doc, elem) {
				t.Errorf("ARCHITECTURE.md should document Kubernetes element: %s", elem)
			}
		}
	})

	t.Run("DocumentsHealthProbes", func(t *testing.T) {
		if !strings.Contains(doc, "Liveness") {
			t.Error("ARCHITECTURE.md should document liveness probes")
		}
		if !strings.Contains(doc, "Readiness") {
			t.Error("ARCHITECTURE.md should document readiness probes")
		}
	})

	t.Run("DocumentsPersistentStorage", func(t *testing.T) {
		if !strings.Contains(doc, "Persistent Volume") {
			t.Error("ARCHITECTURE.md should document persistent storage for PostgreSQL")
		}
	})
}

func TestArchitectureDocumentationFlows(t *testing.T) {
	content, err := os.ReadFile("../../docs/ARCHITECTURE.md")
	if err != nil {
		t.Fatalf("Failed to read ARCHITECTURE.md: %v", err)
	}

	doc := string(content)

	t.Run("DocumentsConfigurationPrecedence", func(t *testing.T) {
		if !strings.Contains(doc, "CLI Flags") {
			t.Error("ARCHITECTURE.md should document CLI flags priority")
		}
		if !strings.Contains(doc, "Environment Variables") {
			t.Error("ARCHITECTURE.md should document environment variables priority")
		}
		if !strings.Contains(doc, "Configuration File") {
			t.Error("ARCHITECTURE.md should document config file priority")
		}
		if !strings.Contains(doc, "Default Values") {
			t.Error("ARCHITECTURE.md should document default values priority")
		}
	})

	t.Run("DocumentsReadinessProbe", func(t *testing.T) {
		elements := []string{
			"GET /healthz/readiness",
			"Checker Interface",
			"200 OK",
			"503 Service Unavailable",
		}
		for _, elem := range elements {
			if !strings.Contains(doc, elem) {
				t.Errorf("ARCHITECTURE.md should document readiness element: %s", elem)
			}
		}
	})

	t.Run("DocumentsMigrationExecution", func(t *testing.T) {
		steps := []string{
			"Load Configuration",
			"Connect to Database",
			"Initialize Migrator",
			"Check Current Version",
			"Execute Migrations",
		}
		for _, step := range steps {
			if !strings.Contains(doc, step) {
				t.Errorf("ARCHITECTURE.md should document migration step: %s", step)
			}
		}
	})

	t.Run("DocumentsGracefulShutdown", func(t *testing.T) {
		if !strings.Contains(doc, "Graceful Shutdown Flow") {
			t.Error("ARCHITECTURE.md should document graceful shutdown")
		}
		steps := []string{
			"SIGINT",
			"SIGTERM",
			"Stop Accepting New Connections",
			"Wait for In-Flight Requests",
			"Close Database Connections",
		}
		for _, step := range steps {
			if !strings.Contains(doc, step) {
				t.Errorf("ARCHITECTURE.md should document shutdown step: %s", step)
			}
		}
	})
}

func TestArchitectureDocumentationSecurity(t *testing.T) {
	content, err := os.ReadFile("../../docs/ARCHITECTURE.md")
	if err != nil {
		t.Fatalf("Failed to read ARCHITECTURE.md: %v", err)
	}

	doc := string(content)

	t.Run("HasSecurityArchitecture", func(t *testing.T) {
		if !strings.Contains(doc, "## Security Architecture") {
			t.Error("ARCHITECTURE.md missing Security Architecture section")
		}
	})

	t.Run("DocumentsSecurityLayers", func(t *testing.T) {
		layers := []string{
			"Network Security",
			"Application Security",
			"Data Security",
			"Container Security",
		}
		for _, layer := range layers {
			if !strings.Contains(doc, layer) {
				t.Errorf("ARCHITECTURE.md should document security layer: %s", layer)
			}
		}
	})

	t.Run("DocumentsNetworkSecurity", func(t *testing.T) {
		if !strings.Contains(doc, "TLS") || !strings.Contains(doc, "Ingress") {
			t.Error("ARCHITECTURE.md should document network security measures")
		}
	})

	t.Run("DocumentsContainerSecurity", func(t *testing.T) {
		measures := []string{
			"Non-root user",
			"Read-only",
			"Minimal base image",
		}
		for _, measure := range measures {
			if !strings.Contains(doc, measure) {
				t.Errorf("ARCHITECTURE.md should document container security: %s", measure)
			}
		}
	})
}

func TestArchitectureDocumentationObservability(t *testing.T) {
	content, err := os.ReadFile("../../docs/ARCHITECTURE.md")
	if err != nil {
		t.Fatalf("Failed to read ARCHITECTURE.md: %v", err)
	}

	doc := string(content)

	t.Run("HasObservabilityArchitecture", func(t *testing.T) {
		if !strings.Contains(doc, "## Observability Architecture") {
			t.Error("ARCHITECTURE.md missing Observability Architecture section")
		}
	})

	t.Run("DocumentsLogging", func(t *testing.T) {
		if !strings.Contains(doc, "Logs") {
			t.Error("ARCHITECTURE.md should document logging")
		}
		if !strings.Contains(doc, "stdout") && !strings.Contains(doc, "stderr") {
			t.Error("ARCHITECTURE.md should document log output streams")
		}
	})

	t.Run("MentionsFutureEnhancements", func(t *testing.T) {
		// Document mentions future observability enhancements
		if !strings.Contains(doc, "Metrics") || !strings.Contains(doc, "Traces") {
			t.Error("ARCHITECTURE.md should mention metrics and traces for future")
		}
	})
}

func TestArchitectureDocumentationTechnologyStack(t *testing.T) {
	content, err := os.ReadFile("../../docs/ARCHITECTURE.md")
	if err != nil {
		t.Fatalf("Failed to read ARCHITECTURE.md: %v", err)
	}

	doc := string(content)

	t.Run("HasTechnologyStack", func(t *testing.T) {
		if !strings.Contains(doc, "## Technology Stack Summary") {
			t.Error("ARCHITECTURE.md missing Technology Stack Summary section")
		}
	})

	t.Run("DocumentsLanguageAndRuntime", func(t *testing.T) {
		if !strings.Contains(doc, "Go") {
			t.Error("ARCHITECTURE.md should document Go language")
		}
	})

	t.Run("DocumentsCoreLibraries", func(t *testing.T) {
		libraries := []string{
			"net/http",
			"log/slog",
			"database/sql",
		}
		for _, lib := range libraries {
			if !strings.Contains(doc, lib) {
				t.Errorf("ARCHITECTURE.md should document core library: %s", lib)
			}
		}
	})

	t.Run("DocumentsThirdPartyDependencies", func(t *testing.T) {
		dependencies := []string{
			"viper",
			"cobra",
			"pq",
			"migrate",
			"swag",
		}
		for _, dep := range dependencies {
			if !strings.Contains(doc, dep) {
				t.Errorf("ARCHITECTURE.md should document dependency: %s", dep)
			}
		}
	})

	t.Run("DocumentsDatabase", func(t *testing.T) {
		if !strings.Contains(doc, "PostgreSQL") {
			t.Error("ARCHITECTURE.md should document PostgreSQL")
		}
	})

	t.Run("DocumentsContainerRuntime", func(t *testing.T) {
		if !strings.Contains(doc, "Docker") {
			t.Error("ARCHITECTURE.md should document Docker")
		}
	})

	t.Run("DocumentsOrchestration", func(t *testing.T) {
		if !strings.Contains(doc, "Kubernetes") {
			t.Error("ARCHITECTURE.md should document Kubernetes")
		}
	})
}

func TestArchitectureDocumentationDesignPatterns(t *testing.T) {
	content, err := os.ReadFile("../../docs/ARCHITECTURE.md")
	if err != nil {
		t.Fatalf("Failed to read ARCHITECTURE.md: %v", err)
	}

	doc := string(content)

	t.Run("HasDesignPatternsSection", func(t *testing.T) {
		if !strings.Contains(doc, "## Design Patterns and Best Practices") {
			t.Error("ARCHITECTURE.md missing Design Patterns section")
		}
	})

	t.Run("DocumentsDesignPatterns", func(t *testing.T) {
		patterns := []string{
			"Dependency Injection",
			"Interface Segregation",
			"Factory Pattern",
			"Strategy Pattern",
		}
		for _, pattern := range patterns {
			if !strings.Contains(doc, pattern) {
				t.Errorf("ARCHITECTURE.md should document design pattern: %s", pattern)
			}
		}
	})

	t.Run("Documents12FactorCompliance", func(t *testing.T) {
		if !strings.Contains(doc, "12-Factor App") {
			t.Error("ARCHITECTURE.md should document 12-Factor App compliance")
		}
	})

	t.Run("DocumentsCloudNativePrinciples", func(t *testing.T) {
		principles := []string{
			"Health checks",
			"Graceful shutdown",
			"Configuration externalization",
		}
		for _, principle := range principles {
			if !strings.Contains(doc, principle) {
				t.Errorf("ARCHITECTURE.md should document cloud-native principle: %s", principle)
			}
		}
	})
}

func TestArchitectureDocumentationPerformance(t *testing.T) {
	content, err := os.ReadFile("../../docs/ARCHITECTURE.md")
	if err != nil {
		t.Fatalf("Failed to read ARCHITECTURE.md: %v", err)
	}

	doc := string(content)

	t.Run("HasPerformanceSection", func(t *testing.T) {
		if !strings.Contains(doc, "## Performance Considerations") {
			t.Error("ARCHITECTURE.md missing Performance Considerations section")
		}
	})

	t.Run("DocumentsConnectionPooling", func(t *testing.T) {
		if !strings.Contains(doc, "Connection Pool") || !strings.Contains(doc, "Connection Pooling") {
			t.Error("ARCHITECTURE.md should document connection pooling")
		}
		if !strings.Contains(doc, "MaxOpenConns") || !strings.Contains(doc, "Max Open Connections") {
			t.Error("ARCHITECTURE.md should document connection pool configuration")
		}
	})

	t.Run("DocumentsHTTPServerTuning", func(t *testing.T) {
		if !strings.Contains(doc, "HTTP Server") {
			t.Error("ARCHITECTURE.md should document HTTP server configuration")
		}
	})
}

func TestArchitectureDocumentationFutureEnhancements(t *testing.T) {
	content, err := os.ReadFile("../../docs/ARCHITECTURE.md")
	if err != nil {
		t.Fatalf("Failed to read ARCHITECTURE.md: %v", err)
	}

	doc := string(content)

	t.Run("HasFutureEnhancements", func(t *testing.T) {
		if !strings.Contains(doc, "## Future Enhancements") {
			t.Error("ARCHITECTURE.md missing Future Enhancements section")
		}
	})

	t.Run("MentionsMetricsCollection", func(t *testing.T) {
		if !strings.Contains(doc, "Metrics") && !strings.Contains(doc, "Prometheus") {
			t.Error("ARCHITECTURE.md should mention future metrics collection")
		}
	})

	t.Run("MentionsDistributedTracing", func(t *testing.T) {
		if !strings.Contains(doc, "Tracing") || !strings.Contains(doc, "OpenTelemetry") {
			t.Error("ARCHITECTURE.md should mention future distributed tracing")
		}
	})
}

func TestArchitectureDocumentationReferences(t *testing.T) {
	content, err := os.ReadFile("../../docs/ARCHITECTURE.md")
	if err != nil {
		t.Fatalf("Failed to read ARCHITECTURE.md: %v", err)
	}

	doc := string(content)

	t.Run("HasReferencesSection", func(t *testing.T) {
		if !strings.Contains(doc, "## References") {
			t.Error("ARCHITECTURE.md missing References section")
		}
	})

	t.Run("HasExternalLinks", func(t *testing.T) {
		if !strings.Contains(doc, "http") {
			t.Error("ARCHITECTURE.md should contain external reference links")
		}
	})
}

func TestArchitectureDocumentationQuality(t *testing.T) {
	content, err := os.ReadFile("../../docs/ARCHITECTURE.md")
	if err != nil {
		t.Fatalf("Failed to read ARCHITECTURE.md: %v", err)
	}

	doc := string(content)

	t.Run("HasSubstantialContent", func(t *testing.T) {
		// Architecture documentation should be comprehensive (at least 30KB)
		if len(content) < 30000 {
			t.Errorf("ARCHITECTURE.md seems too short (%d bytes), expected at least 30000 bytes", len(content))
		}
	})

	t.Run("HasMultipleDiagrams", func(t *testing.T) {
		// Count number of ASCII diagrams (boxes with corners)
		diagramCount := strings.Count(doc, "┌─────")
		if diagramCount < 10 {
			t.Errorf("ARCHITECTURE.md should contain at least 10 diagrams, found %d", diagramCount)
		}
	})

	t.Run("HasCodeBlocks", func(t *testing.T) {
		codeBlockCount := strings.Count(doc, "```")
		if codeBlockCount < 4 {
			t.Errorf("ARCHITECTURE.md should contain code examples, found %d code blocks", codeBlockCount/2)
		}
	})

	t.Run("HasHeaders", func(t *testing.T) {
		headerCount := strings.Count(doc, "## ")
		if headerCount < 15 {
			t.Errorf("ARCHITECTURE.md should have comprehensive sections, found only %d level-2 headers", headerCount)
		}
	})
}

func TestArchitectureDocumentationDiagrams(t *testing.T) {
	content, err := os.ReadFile("../../docs/ARCHITECTURE.md")
	if err != nil {
		t.Fatalf("Failed to read ARCHITECTURE.md: %v", err)
	}

	doc := string(content)

	t.Run("HasHighLevelDiagram", func(t *testing.T) {
		// Check for external clients and HTTP server in diagrams
		if !strings.Contains(doc, "External Clients") {
			t.Error("ARCHITECTURE.md should have high-level system diagram with external clients")
		}
		if !strings.Contains(doc, "HTTP Server") {
			t.Error("ARCHITECTURE.md should show HTTP server in diagrams")
		}
	})

	t.Run("HasComponentDiagram", func(t *testing.T) {
		// Check for component structure
		if !strings.Contains(doc, "cmd/server/") {
			t.Error("ARCHITECTURE.md should have component diagram showing cmd structure")
		}
	})

	t.Run("HasDeploymentDiagram", func(t *testing.T) {
		// Check for deployment topology
		if !strings.Contains(doc, "Kubernetes Cluster") || !strings.Contains(doc, "Pod") {
			t.Error("ARCHITECTURE.md should have deployment diagram showing K8s topology")
		}
	})

	t.Run("HasFlowDiagrams", func(t *testing.T) {
		// Check for flow diagrams with arrows
		if !strings.Contains(doc, "│") || !strings.Contains(doc, "▼") {
			t.Error("ARCHITECTURE.md should have flow diagrams with directional arrows")
		}
	})
}
