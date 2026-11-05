package main

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"

	"github.com/asajaroff/go-lang-template/internal/config"
)

func newMigrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run database schema migrations (Factor XII: one-off admin process)",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "up",
		Short: "Apply all pending migrations",
		RunE: func(_ *cobra.Command, _ []string) error {
			return runMigrate(func(m *migrate.Migrate) error { return m.Up() })
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "down",
		Short: "Roll back the most recent migration",
		RunE: func(_ *cobra.Command, _ []string) error {
			return runMigrate(func(m *migrate.Migrate) error { return m.Steps(-1) })
		},
	})
	return cmd
}

func runMigrate(step func(*migrate.Migrate) error) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	m, err := migrate.New("file://"+cfg.MigrationsPath, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("initializing migrator: %w", err)
	}
	defer func() { _, _ = m.Close() }()

	if err := step(m); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("running migration: %w", err)
	}
	return nil
}
