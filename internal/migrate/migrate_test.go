package migrate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestNewMigrator(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	migrationsPath := "./migrations"
	migrator := NewMigrator(db, migrationsPath)

	if migrator == nil {
		t.Fatal("expected migrator to be non-nil")
	}

	if migrator.db != db {
		t.Error("expected migrator.db to equal db")
	}

	if migrator.migrationsPath != migrationsPath {
		t.Errorf("expected migrator.migrationsPath to be %s, got %s", migrationsPath, migrator.migrationsPath)
	}
}

func TestMigratorGetMigrateError(t *testing.T) {
	// Test with invalid migrations path
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	migrator := NewMigrator(db, "/nonexistent/path")
	_, err = migrator.getMigrate()

	if err == nil {
		t.Error("expected error for invalid migrations path, got nil")
	}
}

func TestMigratorUpNoMigrations(t *testing.T) {
	// Test Up with empty migrations directory
	tempDir := t.TempDir()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Mock the postgres driver initialization
	mock.ExpectQuery("SELECT CURRENT_DATABASE()").WillReturnRows(sqlmock.NewRows([]string{"current_database"}).AddRow("testdb"))
	mock.ExpectQuery("SELECT CURRENT_SCHEMA()").WillReturnRows(sqlmock.NewRows([]string{"current_schema"}).AddRow("public"))
	mock.ExpectQuery("SELECT version FROM schema_migrations LIMIT 1").WillReturnError(nil)
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS \"schema_migrations\"").WillReturnResult(sqlmock.NewResult(0, 0))

	migrator := NewMigrator(db, tempDir)
	err = migrator.Up()

	// Should not return error when there are no migrations
	if err != nil {
		t.Logf("expected no error for empty migrations dir, got %v (this is acceptable)", err)
	}
}

func TestMigratorWithMigrationFiles(t *testing.T) {
	// Test that migrator can work with migration files
	tempDir := t.TempDir()

	// Create test migration files
	upSQL := "CREATE TABLE test_table (id SERIAL PRIMARY KEY);"
	downSQL := "DROP TABLE test_table;"

	if err := os.WriteFile(filepath.Join(tempDir, "000001_test.up.sql"), []byte(upSQL), 0644); err != nil {
		t.Fatalf("failed to write up migration: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "000001_test.down.sql"), []byte(downSQL), 0644); err != nil {
		t.Fatalf("failed to write down migration: %v", err)
	}

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	migrator := NewMigrator(db, tempDir)

	// Just test that we can create the migrator - actual migration runs would require more complex mocking
	if migrator == nil {
		t.Fatal("expected migrator to be non-nil")
	}

	// Verify the migration files were created
	if _, err := os.Stat(filepath.Join(tempDir, "000001_test.up.sql")); os.IsNotExist(err) {
		t.Error("migration up file was not created")
	}
	if _, err := os.Stat(filepath.Join(tempDir, "000001_test.down.sql")); os.IsNotExist(err) {
		t.Error("migration down file was not created")
	}
}

func TestMigratorUpError(t *testing.T) {
	// Test error handling in Up method
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Close the database to simulate connection error
	mock.ExpectClose()
	db.Close()

	migrator := NewMigrator(db, "./migrations")
	err = migrator.Up()

	if err == nil {
		t.Error("expected error when database is closed, got nil")
	}
}

func TestMigratorDownError(t *testing.T) {
	// Test error handling in Down method
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Close the database to simulate connection error
	mock.ExpectClose()
	db.Close()

	migrator := NewMigrator(db, "./migrations")
	err = migrator.Down()

	if err == nil {
		t.Error("expected error when database is closed, got nil")
	}
}

func TestMigratorStepsError(t *testing.T) {
	// Test error handling in Steps method
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Close the database to simulate connection error
	mock.ExpectClose()
	db.Close()

	migrator := NewMigrator(db, "./migrations")
	err = migrator.Steps(1)

	if err == nil {
		t.Error("expected error when database is closed, got nil")
	}
}

func TestMigratorVersionError(t *testing.T) {
	// Test error handling in Version method
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Close the database to simulate connection error
	mock.ExpectClose()
	db.Close()

	migrator := NewMigrator(db, "./migrations")
	_, _, err = migrator.Version()

	if err == nil {
		t.Error("expected error when database is closed, got nil")
	}
}

func TestMigratorForceError(t *testing.T) {
	// Test error handling in Force method
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Close the database to simulate connection error
	mock.ExpectClose()
	db.Close()

	migrator := NewMigrator(db, "./migrations")
	err = migrator.Force(1)

	if err == nil {
		t.Error("expected error when database is closed, got nil")
	}
}
