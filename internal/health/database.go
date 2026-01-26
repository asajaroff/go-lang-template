package health

import (
	"database/sql"
	"fmt"
)

// DatabaseChecker checks database connectivity
type DatabaseChecker struct {
	db *sql.DB
}

// NewDatabaseChecker creates a new database health checker
func NewDatabaseChecker(db *sql.DB) *DatabaseChecker {
	return &DatabaseChecker{
		db: db,
	}
}

// Check verifies database connectivity
func (d *DatabaseChecker) Check() error {
	if d.db == nil {
		return fmt.Errorf("database connection is nil")
	}
	return d.db.Ping()
}

// Name returns the checker name
func (d *DatabaseChecker) Name() string {
	return "database"
}
