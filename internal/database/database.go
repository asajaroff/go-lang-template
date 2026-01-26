package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// Config holds database configuration
type Config struct {
	Driver   string
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	SSLMode  string
}

// New creates a new database connection
func New(cfg Config) (*sql.DB, error) {
	var dsn string

	switch cfg.Driver {
	case "postgres", "postgresql":
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	db, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// Close closes the database connection
func Close(db *sql.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}
