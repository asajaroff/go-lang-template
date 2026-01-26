package docs

import (
	"os"
	"strings"
	"testing"
)

func TestDeploymentDocumentationExists(t *testing.T) {
	deploymentDoc := "../../docs/DEPLOYMENT.md"
	if _, err := os.Stat(deploymentDoc); os.IsNotExist(err) {
		t.Fatalf("DEPLOYMENT.md does not exist at %s", deploymentDoc)
	}
}

func TestDeploymentDocumentationContent(t *testing.T) {
	deploymentDoc := "../../docs/DEPLOYMENT.md"
	content, err := os.ReadFile(deploymentDoc)
	if err != nil {
		t.Fatalf("Failed to read DEPLOYMENT.md: %v", err)
	}

	contentStr := string(content)
	if len(contentStr) == 0 {
		t.Fatal("DEPLOYMENT.md is empty")
	}

	// Check for main sections
	requiredSections := []string{
		"# Production Deployment Guide",
		"## Prerequisites",
		"## Build and Release",
		"## Docker Deployment",
		"## Docker Compose Deployment",
		"## Kubernetes Deployment",
		"## Production Configuration",
		"## Security Best Practices",
		"## Monitoring and Logging",
		"## Troubleshooting",
	}

	for _, section := range requiredSections {
		if !strings.Contains(contentStr, section) {
			t.Errorf("DEPLOYMENT.md is missing required section: %s", section)
		}
	}
}

func TestDeploymentDockerInstructions(t *testing.T) {
	deploymentDoc := "../../docs/DEPLOYMENT.md"
	content, err := os.ReadFile(deploymentDoc)
	if err != nil {
		t.Fatalf("Failed to read DEPLOYMENT.md: %v", err)
	}

	contentStr := string(content)

	// Check for Docker-related content
	dockerElements := []string{
		"FROM golang:",
		"FROM alpine:",
		"docker build",
		"docker run",
		"docker push",
		"Dockerfile",
		"HEALTHCHECK",
		"ENTRYPOINT",
	}

	for _, element := range dockerElements {
		if !strings.Contains(contentStr, element) {
			t.Errorf("DEPLOYMENT.md is missing Docker instruction: %s", element)
		}
	}
}

func TestDeploymentDockerComposeInstructions(t *testing.T) {
	deploymentDoc := "../../docs/DEPLOYMENT.md"
	content, err := os.ReadFile(deploymentDoc)
	if err != nil {
		t.Fatalf("Failed to read DEPLOYMENT.md: %v", err)
	}

	contentStr := string(content)

	// Check for Docker Compose-related content
	composeElements := []string{
		"version:",
		"services:",
		"docker-compose",
		"depends_on:",
		"healthcheck:",
		"networks:",
		"volumes:",
	}

	for _, element := range composeElements {
		if !strings.Contains(contentStr, element) {
			t.Errorf("DEPLOYMENT.md is missing Docker Compose instruction: %s", element)
		}
	}
}

func TestDeploymentKubernetesInstructions(t *testing.T) {
	deploymentDoc := "../../docs/DEPLOYMENT.md"
	content, err := os.ReadFile(deploymentDoc)
	if err != nil {
		t.Fatalf("Failed to read DEPLOYMENT.md: %v", err)
	}

	contentStr := string(content)

	// Check for Kubernetes-related content
	k8sElements := []string{
		"apiVersion:",
		"kind: Deployment",
		"kind: Service",
		"kind: Secret",
		"kind: ConfigMap",
		"kind: Ingress",
		"kind: Job",
		"kubectl apply",
		"livenessProbe:",
		"readinessProbe:",
		"HorizontalPodAutoscaler",
	}

	for _, element := range k8sElements {
		if !strings.Contains(contentStr, element) {
			t.Errorf("DEPLOYMENT.md is missing Kubernetes instruction: %s", element)
		}
	}
}

func TestDeploymentHealthCheckConfiguration(t *testing.T) {
	deploymentDoc := "../../docs/DEPLOYMENT.md"
	content, err := os.ReadFile(deploymentDoc)
	if err != nil {
		t.Fatalf("Failed to read DEPLOYMENT.md: %v", err)
	}

	contentStr := string(content)

	// Check for health check endpoints
	healthCheckElements := []string{
		"/healthz/liveness",
		"/healthz/readiness",
		"httpGet:",
		"initialDelaySeconds:",
		"periodSeconds:",
	}

	for _, element := range healthCheckElements {
		if !strings.Contains(contentStr, element) {
			t.Errorf("DEPLOYMENT.md is missing health check configuration: %s", element)
		}
	}
}

func TestDeploymentSecurityBestPractices(t *testing.T) {
	deploymentDoc := "../../docs/DEPLOYMENT.md"
	content, err := os.ReadFile(deploymentDoc)
	if err != nil {
		t.Fatalf("Failed to read DEPLOYMENT.md: %v", err)
	}

	contentStr := string(content)

	// Check for security-related content
	securityElements := []string{
		"Security Best Practices",
		"non-root user",
		"SSL/TLS",
		"sslmode: require",
		"secrets",
		"strong passwords",
		"runAsNonRoot:",
		"allowPrivilegeEscalation: false",
		"capabilities:",
	}

	for _, element := range securityElements {
		if !strings.Contains(contentStr, element) {
			t.Errorf("DEPLOYMENT.md is missing security best practice: %s", element)
		}
	}
}

func TestDeploymentSecretsManagement(t *testing.T) {
	deploymentDoc := "../../docs/DEPLOYMENT.md"
	content, err := os.ReadFile(deploymentDoc)
	if err != nil {
		t.Fatalf("Failed to read DEPLOYMENT.md: %v", err)
	}

	contentStr := string(content)

	// Check for secrets management content
	secretsElements := []string{
		"kind: Secret",
		"APP_DB_PASSWORD",
		"secretKeyRef:",
		"Never commit secrets",
		"Rotate secrets",
	}

	for _, element := range secretsElements {
		if !strings.Contains(contentStr, element) {
			t.Errorf("DEPLOYMENT.md is missing secrets management instruction: %s", element)
		}
	}
}

func TestDeploymentDatabaseMigrations(t *testing.T) {
	deploymentDoc := "../../docs/DEPLOYMENT.md"
	content, err := os.ReadFile(deploymentDoc)
	if err != nil {
		t.Fatalf("Failed to read DEPLOYMENT.md: %v", err)
	}

	contentStr := string(content)

	// Check for database migration content
	migrationElements := []string{
		"migrate up",
		"migrate down",
		"migration-job.yaml",
		"kind: Job",
		"Run Database Migrations",
	}

	for _, element := range migrationElements {
		if !strings.Contains(contentStr, element) {
			t.Errorf("DEPLOYMENT.md is missing database migration instruction: %s", element)
		}
	}
}

func TestDeploymentMonitoringAndLogging(t *testing.T) {
	deploymentDoc := "../../docs/DEPLOYMENT.md"
	content, err := os.ReadFile(deploymentDoc)
	if err != nil {
		t.Fatalf("Failed to read DEPLOYMENT.md: %v", err)
	}

	contentStr := string(content)

	// Check for monitoring and logging content
	monitoringElements := []string{
		"Monitoring and Logging",
		"Prometheus",
		"Grafana",
		"JSON logs",
		"Health Check Monitoring",
		"Log Aggregation",
		"Metrics Collection",
	}

	for _, element := range monitoringElements {
		if !strings.Contains(contentStr, element) {
			t.Errorf("DEPLOYMENT.md is missing monitoring/logging instruction: %s", element)
		}
	}
}

func TestDeploymentTroubleshootingSection(t *testing.T) {
	deploymentDoc := "../../docs/DEPLOYMENT.md"
	content, err := os.ReadFile(deploymentDoc)
	if err != nil {
		t.Fatalf("Failed to read DEPLOYMENT.md: %v", err)
	}

	contentStr := string(content)

	// Check for troubleshooting content
	troubleshootingElements := []string{
		"## Troubleshooting",
		"Application Not Starting",
		"Database Connection Errors",
		"Readiness Probe Failing",
		"Migration Failures",
		"Check logs:",
		"kubectl logs",
		"docker logs",
	}

	for _, element := range troubleshootingElements {
		if !strings.Contains(contentStr, element) {
			t.Errorf("DEPLOYMENT.md is missing troubleshooting information: %s", element)
		}
	}
}

func TestDeploymentEnvironmentConfiguration(t *testing.T) {
	deploymentDoc := "../../docs/DEPLOYMENT.md"
	content, err := os.ReadFile(deploymentDoc)
	if err != nil {
		t.Fatalf("Failed to read DEPLOYMENT.md: %v", err)
	}

	contentStr := string(content)

	// Check for environment-specific configuration
	envElements := []string{
		"Production Configuration",
		"Environment-Specific Settings",
		"Production",
		"Staging",
		"APP_SERVER_HOST",
		"APP_SERVER_PORT",
		"APP_DB_ENABLED",
		"APP_LOG_LEVEL",
		"APP_LOG_FORMAT",
	}

	for _, element := range envElements {
		if !strings.Contains(contentStr, element) {
			t.Errorf("DEPLOYMENT.md is missing environment configuration: %s", element)
		}
	}
}

func TestDeploymentResourceManagement(t *testing.T) {
	deploymentDoc := "../../docs/DEPLOYMENT.md"
	content, err := os.ReadFile(deploymentDoc)
	if err != nil {
		t.Fatalf("Failed to read DEPLOYMENT.md: %v", err)
	}

	contentStr := string(content)

	// Check for resource management content
	resourceElements := []string{
		"resources:",
		"requests:",
		"limits:",
		"memory:",
		"cpu:",
		"HorizontalPodAutoscaler",
		"minReplicas:",
		"maxReplicas:",
	}

	for _, element := range resourceElements {
		if !strings.Contains(contentStr, element) {
			t.Errorf("DEPLOYMENT.md is missing resource management configuration: %s", element)
		}
	}
}

func TestDeploymentRollbackProcedures(t *testing.T) {
	deploymentDoc := "../../docs/DEPLOYMENT.md"
	content, err := os.ReadFile(deploymentDoc)
	if err != nil {
		t.Fatalf("Failed to read DEPLOYMENT.md: %v", err)
	}

	contentStr := string(content)

	// Check for rollback procedures
	rollbackElements := []string{
		"Rollback Procedures",
		"Docker Rollback",
		"Kubernetes Rollback",
		"Database Migration Rollback",
		"kubectl rollout undo",
		"rollout history",
	}

	for _, element := range rollbackElements {
		if !strings.Contains(contentStr, element) {
			t.Errorf("DEPLOYMENT.md is missing rollback procedure: %s", element)
		}
	}
}

func TestDeploymentMaintenanceProcedures(t *testing.T) {
	deploymentDoc := "../../docs/DEPLOYMENT.md"
	content, err := os.ReadFile(deploymentDoc)
	if err != nil {
		t.Fatalf("Failed to read DEPLOYMENT.md: %v", err)
	}

	contentStr := string(content)

	// Check for maintenance procedures
	maintenanceElements := []string{
		"## Maintenance",
		"Regular Tasks",
		"Backup Procedures",
		"Upgrade Procedures",
		"Weekly:",
		"Monthly:",
		"pg_dump",
	}

	for _, element := range maintenanceElements {
		if !strings.Contains(contentStr, element) {
			t.Errorf("DEPLOYMENT.md is missing maintenance procedure: %s", element)
		}
	}
}

func TestDeploymentTableOfContents(t *testing.T) {
	deploymentDoc := "../../docs/DEPLOYMENT.md"
	content, err := os.ReadFile(deploymentDoc)
	if err != nil {
		t.Fatalf("Failed to read DEPLOYMENT.md: %v", err)
	}

	contentStr := string(content)

	// Check for table of contents
	if !strings.Contains(contentStr, "## Table of Contents") {
		t.Error("DEPLOYMENT.md is missing table of contents")
	}

	tocElements := []string{
		"[Prerequisites]",
		"[Docker Deployment]",
		"[Kubernetes Deployment]",
		"[Security Best Practices]",
		"[Troubleshooting]",
	}

	for _, element := range tocElements {
		if !strings.Contains(contentStr, element) {
			t.Errorf("DEPLOYMENT.md table of contents is missing link: %s", element)
		}
	}
}

func TestDeploymentCodeExamples(t *testing.T) {
	deploymentDoc := "../../docs/DEPLOYMENT.md"
	content, err := os.ReadFile(deploymentDoc)
	if err != nil {
		t.Fatalf("Failed to read DEPLOYMENT.md: %v", err)
	}

	contentStr := string(content)

	// Check for code blocks
	if !strings.Contains(contentStr, "```bash") {
		t.Error("DEPLOYMENT.md should contain bash code examples")
	}

	if !strings.Contains(contentStr, "```yaml") {
		t.Error("DEPLOYMENT.md should contain YAML code examples")
	}

	if !strings.Contains(contentStr, "```dockerfile") {
		t.Error("DEPLOYMENT.md should contain Dockerfile examples")
	}
}

func TestDeploymentNetworkConfiguration(t *testing.T) {
	deploymentDoc := "../../docs/DEPLOYMENT.md"
	content, err := os.ReadFile(deploymentDoc)
	if err != nil {
		t.Fatalf("Failed to read DEPLOYMENT.md: %v", err)
	}

	contentStr := string(content)

	// Check for network configuration
	networkElements := []string{
		"kind: Service",
		"kind: Ingress",
		"type: LoadBalancer",
		"type: ClusterIP",
		"ingressClassName:",
		"Network Security",
	}

	for _, element := range networkElements {
		if !strings.Contains(contentStr, element) {
			t.Errorf("DEPLOYMENT.md is missing network configuration: %s", element)
		}
	}
}

func TestDeploymentPersistentStorage(t *testing.T) {
	deploymentDoc := "../../docs/DEPLOYMENT.md"
	content, err := os.ReadFile(deploymentDoc)
	if err != nil {
		t.Fatalf("Failed to read DEPLOYMENT.md: %v", err)
	}

	contentStr := string(content)

	// Check for persistent storage configuration
	storageElements := []string{
		"PersistentVolumeClaim",
		"volumeMounts:",
		"volumes:",
		"persistentVolumeClaim:",
		"storage:",
	}

	for _, element := range storageElements {
		if !strings.Contains(contentStr, element) {
			t.Errorf("DEPLOYMENT.md is missing persistent storage configuration: %s", element)
		}
	}
}

func TestDeploymentDocumentationLength(t *testing.T) {
	deploymentDoc := "../../docs/DEPLOYMENT.md"
	content, err := os.ReadFile(deploymentDoc)
	if err != nil {
		t.Fatalf("Failed to read DEPLOYMENT.md: %v", err)
	}

	// Check that documentation is comprehensive (at least 20KB)
	if len(content) < 20000 {
		t.Errorf("DEPLOYMENT.md seems too short (%d bytes). Expected comprehensive deployment guide (>20KB)", len(content))
	}

	// Check for reasonable number of sections (at least 10 major sections)
	contentStr := string(content)
	sectionCount := strings.Count(contentStr, "\n## ")
	if sectionCount < 10 {
		t.Errorf("DEPLOYMENT.md has only %d major sections. Expected at least 10 comprehensive sections", sectionCount)
	}
}
