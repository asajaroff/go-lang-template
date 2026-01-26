package docs

import (
	"os"
	"strings"
	"testing"
)

func TestDatabaseDocumentationExists(t *testing.T) {
	content, err := os.ReadFile("../../docs/DATABASE.md")
	if err != nil {
		t.Fatalf("DATABASE.md does not exist: %v", err)
	}

	if len(content) == 0 {
		t.Fatal("DATABASE.md is empty")
	}

	// Check minimum content length (comprehensive documentation should be substantial)
	if len(content) < 10000 {
		t.Errorf("DATABASE.md seems too short (%d bytes), expected comprehensive documentation", len(content))
	}
}

func TestDatabaseDocumentationStructure(t *testing.T) {
	content, err := os.ReadFile("../../docs/DATABASE.md")
	if err != nil {
		t.Fatalf("Failed to read DATABASE.md: %v", err)
	}

	doc := string(content)

	// Check for required sections
	requiredSections := []string{
		"# Database Schema Documentation",
		"## Table of Contents",
		"## Overview",
		"## Database Schema",
		"## Tables",
		"## Indexes",
		"## Relationships",
		"## Migration Strategy",
		"## Data Types",
		"## Best Practices",
		"## Examples",
	}

	for _, section := range requiredSections {
		if !strings.Contains(doc, section) {
			t.Errorf("DATABASE.md missing required section: %s", section)
		}
	}
}

func TestDatabaseDocumentationUsersTable(t *testing.T) {
	content, err := os.ReadFile("../../docs/DATABASE.md")
	if err != nil {
		t.Fatalf("Failed to read DATABASE.md: %v", err)
	}

	doc := string(content)

	// Check for users table documentation
	requiredElements := []string{
		"### users",
		"id",
		"username",
		"email",
		"created_at",
		"updated_at",
		"SERIAL",
		"VARCHAR(255)",
		"TIMESTAMP",
		"PRIMARY KEY",
		"UNIQUE",
	}

	for _, element := range requiredElements {
		if !strings.Contains(doc, element) {
			t.Errorf("DATABASE.md missing users table element: %s", element)
		}
	}
}

func TestDatabaseDocumentationIndexes(t *testing.T) {
	content, err := os.ReadFile("../../docs/DATABASE.md")
	if err != nil {
		t.Fatalf("Failed to read DATABASE.md: %v", err)
	}

	doc := string(content)

	// Check for index documentation
	requiredIndexes := []string{
		"idx_users_username",
		"idx_users_email",
		"B-tree",
	}

	for _, index := range requiredIndexes {
		if !strings.Contains(doc, index) {
			t.Errorf("DATABASE.md missing index documentation: %s", index)
		}
	}
}

func TestDatabaseDocumentationMigrations(t *testing.T) {
	content, err := os.ReadFile("../../docs/DATABASE.md")
	if err != nil {
		t.Fatalf("Failed to read DATABASE.md: %v", err)
	}

	doc := string(content)

	// Check for migration documentation
	requiredMigrations := []string{
		"golang-migrate",
		"migrations/000001_create_users_table.up.sql",
		"migrations/000001_create_users_table.down.sql",
		"migrate up",
		"migrate down",
		"migrate steps",
		"migrate version",
		"migrate force",
	}

	for _, migration := range requiredMigrations {
		if !strings.Contains(doc, migration) {
			t.Errorf("DATABASE.md missing migration documentation: %s", migration)
		}
	}
}

func TestDatabaseDocumentationMigrationCommands(t *testing.T) {
	content, err := os.ReadFile("../../docs/DATABASE.md")
	if err != nil {
		t.Fatalf("Failed to read DATABASE.md: %v", err)
	}

	doc := string(content)

	// Check for specific migration commands
	commands := []string{
		"./bin/go-lang-template migrate up",
		"./bin/go-lang-template migrate down",
		"./bin/go-lang-template migrate steps",
		"./bin/go-lang-template migrate version",
		"./bin/go-lang-template migrate force",
	}

	for _, cmd := range commands {
		if !strings.Contains(doc, cmd) {
			t.Errorf("DATABASE.md missing migration command: %s", cmd)
		}
	}
}

func TestDatabaseDocumentationDataTypes(t *testing.T) {
	content, err := os.ReadFile("../../docs/DATABASE.md")
	if err != nil {
		t.Fatalf("Failed to read DATABASE.md: %v", err)
	}

	doc := string(content)

	// Check for data type documentation
	dataTypes := []string{
		"SERIAL",
		"VARCHAR",
		"TEXT",
		"INTEGER",
		"BIGINT",
		"BOOLEAN",
		"TIMESTAMP",
		"TIMESTAMPTZ",
		"NUMERIC",
		"JSON",
		"JSONB",
		"UUID",
	}

	for _, dataType := range dataTypes {
		if !strings.Contains(doc, dataType) {
			t.Errorf("DATABASE.md missing data type documentation: %s", dataType)
		}
	}
}

func TestDatabaseDocumentationBestPractices(t *testing.T) {
	content, err := os.ReadFile("../../docs/DATABASE.md")
	if err != nil {
		t.Fatalf("Failed to read DATABASE.md: %v", err)
	}

	doc := string(content)

	// Check for best practices documentation
	bestPractices := []string{
		"Normalization",
		"Primary Keys",
		"Foreign Keys",
		"Indexes",
		"Timestamps",
		"Constraints",
		"Naming Conventions",
		"Security Considerations",
	}

	for _, practice := range bestPractices {
		if !strings.Contains(doc, practice) {
			t.Errorf("DATABASE.md missing best practice: %s", practice)
		}
	}
}

func TestDatabaseDocumentationRelationships(t *testing.T) {
	content, err := os.ReadFile("../../docs/DATABASE.md")
	if err != nil {
		t.Fatalf("Failed to read DATABASE.md: %v", err)
	}

	doc := string(content)

	// Check for relationship pattern examples
	relationships := []string{
		"One-to-Many",
		"Many-to-Many",
		"One-to-One",
		"FOREIGN KEY",
		"ON DELETE CASCADE",
	}

	for _, relationship := range relationships {
		if !strings.Contains(doc, relationship) {
			t.Errorf("DATABASE.md missing relationship pattern: %s", relationship)
		}
	}
}

func TestDatabaseDocumentationExamples(t *testing.T) {
	content, err := os.ReadFile("../../docs/DATABASE.md")
	if err != nil {
		t.Fatalf("Failed to read DATABASE.md: %v", err)
	}

	doc := string(content)

	// Check for code examples
	examples := []string{
		"CREATE TABLE",
		"ALTER TABLE",
		"DROP TABLE",
		"CREATE INDEX",
		"DROP INDEX",
		"SELECT",
		"INSERT",
		"UPDATE",
		"DELETE",
	}

	for _, example := range examples {
		if !strings.Contains(doc, example) {
			t.Errorf("DATABASE.md missing SQL example: %s", example)
		}
	}
}

func TestDatabaseDocumentationGoExamples(t *testing.T) {
	content, err := os.ReadFile("../../docs/DATABASE.md")
	if err != nil {
		t.Fatalf("Failed to read DATABASE.md: %v", err)
	}

	doc := string(content)

	// Check for Go code examples
	goExamples := []string{
		"func GetUserByUsername",
		"func CreateUser",
		"db.QueryRow",
		"Scan",
	}

	for _, example := range goExamples {
		if !strings.Contains(doc, example) {
			t.Errorf("DATABASE.md missing Go code example: %s", example)
		}
	}
}

func TestDatabaseDocumentationSecurity(t *testing.T) {
	content, err := os.ReadFile("../../docs/DATABASE.md")
	if err != nil {
		t.Fatalf("Failed to read DATABASE.md: %v", err)
	}

	doc := string(content)

	// Check for security best practices
	security := []string{
		"SQL injection",
		"prepared statements",
		"SSL/TLS",
		"sslmode",
		"bcrypt",
		"argon2",
	}

	for _, item := range security {
		if !strings.Contains(doc, item) {
			t.Errorf("DATABASE.md missing security consideration: %s", item)
		}
	}
}

func TestDatabaseDocumentationMaintenance(t *testing.T) {
	content, err := os.ReadFile("../../docs/DATABASE.md")
	if err != nil {
		t.Fatalf("Failed to read DATABASE.md: %v", err)
	}

	doc := string(content)

	// Check for maintenance documentation
	maintenance := []string{
		"Maintenance",
		"VACUUM",
		"ANALYZE",
		"REINDEX",
		"pg_size_pretty",
		"pg_stat_user_indexes",
	}

	for _, item := range maintenance {
		if !strings.Contains(doc, item) {
			t.Errorf("DATABASE.md missing maintenance documentation: %s", item)
		}
	}
}

func TestDatabaseDocumentationERD(t *testing.T) {
	content, err := os.ReadFile("../../docs/DATABASE.md")
	if err != nil {
		t.Fatalf("Failed to read DATABASE.md: %v", err)
	}

	doc := string(content)

	// Check for Entity Relationship Diagram
	if !strings.Contains(doc, "Entity Relationship Diagram") {
		t.Error("DATABASE.md missing Entity Relationship Diagram section")
	}

	// Check for ASCII diagram elements
	diagramElements := []string{
		"┌",
		"└",
		"│",
		"├",
	}

	hasAsciiDiagram := false
	for _, element := range diagramElements {
		if strings.Contains(doc, element) {
			hasAsciiDiagram = true
			break
		}
	}

	if !hasAsciiDiagram {
		t.Error("DATABASE.md missing ASCII diagram representation")
	}
}

func TestDatabaseDocumentationPostgreSQLSupport(t *testing.T) {
	content, err := os.ReadFile("../../docs/DATABASE.md")
	if err != nil {
		t.Fatalf("Failed to read DATABASE.md: %v", err)
	}

	doc := string(content)

	// Check for PostgreSQL specific documentation
	postgres := []string{
		"PostgreSQL",
		"Database.Enabled",
		"optional",
	}

	for _, item := range postgres {
		if !strings.Contains(doc, item) {
			t.Errorf("DATABASE.md missing PostgreSQL documentation: %s", item)
		}
	}
}

func TestDatabaseDocumentationMigrationBestPractices(t *testing.T) {
	content, err := os.ReadFile("../../docs/DATABASE.md")
	if err != nil {
		t.Fatalf("Failed to read DATABASE.md: %v", err)
	}

	doc := string(content)

	// Check for migration best practices
	bestPractices := []string{
		"reversible migrations",
		"Test migrations",
		"atomic",
		"IF EXISTS",
		"IF NOT EXISTS",
		"Backup",
		"transactions",
	}

	for _, practice := range bestPractices {
		if !strings.Contains(doc, practice) {
			t.Errorf("DATABASE.md missing migration best practice: %s", practice)
		}
	}
}

func TestDatabaseDocumentationMigrationWorkflow(t *testing.T) {
	content, err := os.ReadFile("../../docs/DATABASE.md")
	if err != nil {
		t.Fatalf("Failed to read DATABASE.md: %v", err)
	}

	doc := string(content)

	// Check for migration workflow steps
	if !strings.Contains(doc, "Migration Workflow") && !strings.Contains(doc, "Migration workflow") {
		t.Error("DATABASE.md missing migration workflow section")
	}

	// Check for workflow commands
	workflow := []string{
		"touch migrations/",
		".up.sql",
		".down.sql",
		"pg_dump",
	}

	for _, step := range workflow {
		if !strings.Contains(doc, step) {
			t.Errorf("DATABASE.md missing workflow step: %s", step)
		}
	}
}

func TestDatabaseDocumentationNamingConventions(t *testing.T) {
	content, err := os.ReadFile("../../docs/DATABASE.md")
	if err != nil {
		t.Fatalf("Failed to read DATABASE.md: %v", err)
	}

	doc := string(content)

	// Check for naming convention documentation
	conventions := []string{
		"idx_",
		"fk_",
		"uq_",
		"Plural lowercase",
		"underscores",
	}

	for _, convention := range conventions {
		if !strings.Contains(doc, convention) {
			t.Errorf("DATABASE.md missing naming convention: %s", convention)
		}
	}
}

func TestDatabaseDocumentationResources(t *testing.T) {
	content, err := os.ReadFile("../../docs/DATABASE.md")
	if err != nil {
		t.Fatalf("Failed to read DATABASE.md: %v", err)
	}

	doc := string(content)

	// Check for external resources
	resources := []string{
		"Resources",
		"PostgreSQL Documentation",
		"golang-migrate",
	}

	for _, resource := range resources {
		if !strings.Contains(doc, resource) {
			t.Errorf("DATABASE.md missing resource reference: %s", resource)
		}
	}
}

func TestDatabaseDocumentationTables(t *testing.T) {
	content, err := os.ReadFile("../../docs/DATABASE.md")
	if err != nil {
		t.Fatalf("Failed to read DATABASE.md: %v", err)
	}

	doc := string(content)

	// Check for markdown tables
	if !strings.Contains(doc, "| Column") || !strings.Contains(doc, "|---") {
		t.Error("DATABASE.md should contain formatted markdown tables")
	}

	// Check for multiple table sections
	tableCount := strings.Count(doc, "|---")
	if tableCount < 3 {
		t.Errorf("DATABASE.md should contain multiple tables, found only %d", tableCount)
	}
}

func TestDatabaseDocumentationCodeBlocks(t *testing.T) {
	content, err := os.ReadFile("../../docs/DATABASE.md")
	if err != nil {
		t.Fatalf("Failed to read DATABASE.md: %v", err)
	}

	doc := string(content)

	// Check for code blocks with language specifications
	codeBlocks := []string{
		"```sql",
		"```go",
		"```bash",
	}

	for _, block := range codeBlocks {
		if !strings.Contains(doc, block) {
			t.Errorf("DATABASE.md missing code block type: %s", block)
		}
	}

	// Ensure there are multiple code examples
	sqlBlockCount := strings.Count(doc, "```sql")
	if sqlBlockCount < 5 {
		t.Errorf("DATABASE.md should contain multiple SQL examples, found only %d", sqlBlockCount)
	}
}

func TestDatabaseDocumentationConstraints(t *testing.T) {
	content, err := os.ReadFile("../../docs/DATABASE.md")
	if err != nil {
		t.Fatalf("Failed to read DATABASE.md: %v", err)
	}

	doc := string(content)

	// Check for constraint documentation
	constraints := []string{
		"NOT NULL",
		"UNIQUE",
		"CHECK",
		"DEFAULT",
	}

	for _, constraint := range constraints {
		if !strings.Contains(doc, constraint) {
			t.Errorf("DATABASE.md missing constraint type: %s", constraint)
		}
	}
}
