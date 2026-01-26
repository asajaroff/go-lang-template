package migrate

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Migrator handles database migrations
type Migrator struct {
	db             *sql.DB
	migrationsPath string
}

// NewMigrator creates a new Migrator instance
func NewMigrator(db *sql.DB, migrationsPath string) *Migrator {
	return &Migrator{
		db:             db,
		migrationsPath: migrationsPath,
	}
}

// Up runs all available migrations
func (m *Migrator) Up() error {
	mig, err := m.getMigrate()
	if err != nil {
		return err
	}
	defer mig.Close()

	// migrate.ErrNoChange is not an error - it just means the database is already up to date
	// We only return an error if something actually went wrong
	if err := mig.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// Down rolls back the last migration
func (m *Migrator) Down() error {
	mig, err := m.getMigrate()
	if err != nil {
		return err
	}
	defer mig.Close()

	// migrate.ErrNoChange is not an error - it just means there are no migrations to roll back
	// We only return an error if something actually went wrong
	if err := mig.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	return nil
}

// Steps runs n migrations. If n is positive, it runs n migrations forward.
// If n is negative, it rolls back n migrations.
func (m *Migrator) Steps(n int) error {
	mig, err := m.getMigrate()
	if err != nil {
		return err
	}
	defer mig.Close()

	// migrate.ErrNoChange is not an error - it just means we can't move the requested number of steps
	// This can happen if we're already at the latest version or there aren't enough migrations
	if err := mig.Steps(n); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run %d migration steps: %w", n, err)
	}

	return nil
}

// Version returns the current migration version
// The dirty flag indicates whether a migration is in progress or failed mid-execution
func (m *Migrator) Version() (uint, bool, error) {
	mig, err := m.getMigrate()
	if err != nil {
		return 0, false, err
	}
	defer mig.Close()

	// migrate.ErrNilVersion is returned when no migrations have been run yet
	// This is a normal state for a fresh database, not an error condition
	version, dirty, err := mig.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return 0, false, fmt.Errorf("failed to get migration version: %w", err)
	}

	return version, dirty, nil
}

// Force sets the migration version without running migrations
// This is useful for recovering from failed migrations or manual database changes
func (m *Migrator) Force(version int) error {
	mig, err := m.getMigrate()
	if err != nil {
		return err
	}
	defer mig.Close()

	// Force the version in the schema_migrations table without running any migrations
	// Use with caution - this can leave your database in an inconsistent state
	if err := mig.Force(version); err != nil {
		return fmt.Errorf("failed to force migration version: %w", err)
	}

	return nil
}

// getMigrate creates a new migrate instance
func (m *Migrator) getMigrate() (*migrate.Migrate, error) {
	// Create a postgres driver instance from our sql.DB connection
	// This allows the migrate library to interact with our database
	driver, err := postgres.WithInstance(m.db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Create the migrate instance with file-based migrations
	// The file:// protocol indicates migrations are stored as local files
	mig, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", m.migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	return mig, nil
}
