package main

import (
	"database/sql"
	"fmt"

	"github.com/asajaroff/go-lang-template/internal/migrate"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var (
	migrateSteps int
	migrateForce int
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `Run database migrations using golang-migrate.`,
}

// migrateUpCmd runs all available migrations
var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Run all available migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := openDatabase()
		if err != nil {
			return err
		}
		defer db.Close()

		migrator := migrate.NewMigrator(db, cfg.Database.MigrationsPath)
		if err := migrator.Up(); err != nil {
			return fmt.Errorf("failed to run migrations: %w", err)
		}

		version, dirty, err := migrator.Version()
		if err != nil {
			return fmt.Errorf("failed to get migration version: %w", err)
		}

		fmt.Printf("Migrations completed successfully. Current version: %d (dirty: %t)\n", version, dirty)
		return nil
	},
}

// migrateDownCmd rolls back the last migration
var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Roll back the last migration",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := openDatabase()
		if err != nil {
			return err
		}
		defer db.Close()

		migrator := migrate.NewMigrator(db, cfg.Database.MigrationsPath)
		if err := migrator.Down(); err != nil {
			return fmt.Errorf("failed to rollback migration: %w", err)
		}

		version, dirty, err := migrator.Version()
		if err != nil {
			return fmt.Errorf("failed to get migration version: %w", err)
		}

		fmt.Printf("Rollback completed successfully. Current version: %d (dirty: %t)\n", version, dirty)
		return nil
	},
}

// migrateStepsCmd runs n migrations
var migrateStepsCmd = &cobra.Command{
	Use:   "steps",
	Short: "Run n migrations (positive for forward, negative for backward)",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := openDatabase()
		if err != nil {
			return err
		}
		defer db.Close()

		migrator := migrate.NewMigrator(db, cfg.Database.MigrationsPath)
		if err := migrator.Steps(migrateSteps); err != nil {
			return fmt.Errorf("failed to run migration steps: %w", err)
		}

		version, dirty, err := migrator.Version()
		if err != nil {
			return fmt.Errorf("failed to get migration version: %w", err)
		}

		fmt.Printf("Migration steps completed successfully. Current version: %d (dirty: %t)\n", version, dirty)
		return nil
	},
}

// migrateVersionCmd shows the current migration version
var migrateVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the current migration version",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := openDatabase()
		if err != nil {
			return err
		}
		defer db.Close()

		migrator := migrate.NewMigrator(db, cfg.Database.MigrationsPath)
		version, dirty, err := migrator.Version()
		if err != nil {
			return fmt.Errorf("failed to get migration version: %w", err)
		}

		fmt.Printf("Current migration version: %d (dirty: %t)\n", version, dirty)
		return nil
	},
}

// migrateForceCmd forces the migration version
var migrateForceCmd = &cobra.Command{
	Use:   "force",
	Short: "Force set the migration version without running migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := openDatabase()
		if err != nil {
			return err
		}
		defer db.Close()

		migrator := migrate.NewMigrator(db, cfg.Database.MigrationsPath)
		if err := migrator.Force(migrateForce); err != nil {
			return fmt.Errorf("failed to force migration version: %w", err)
		}

		fmt.Printf("Forced migration version to: %d\n", migrateForce)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateStepsCmd)
	migrateCmd.AddCommand(migrateVersionCmd)
	migrateCmd.AddCommand(migrateForceCmd)

	migrateStepsCmd.Flags().IntVarP(&migrateSteps, "steps", "n", 1, "number of migration steps to run")
	migrateForceCmd.Flags().IntVarP(&migrateForce, "version", "v", 0, "version to force")
}

// openDatabase opens a database connection for running migrations
// This function is separate from the main database package to avoid circular dependencies
func openDatabase() (*sql.DB, error) {
	// Build PostgreSQL connection string from config
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Verify the connection is actually working before proceeding
	// This catches configuration errors early
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
