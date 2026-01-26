package database

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		expectError bool
		errorMsg    string
	}{
		{
			name: "unsupported driver",
			config: Config{
				Driver: "mysql",
				Host:   "localhost",
				Port:   3306,
			},
			expectError: true,
			errorMsg:    "unsupported database driver",
		},
		{
			name: "postgres driver name",
			config: Config{
				Driver:   "postgres",
				Host:     "localhost",
				Port:     5432,
				Name:     "testdb",
				User:     "testuser",
				Password: "testpass",
				SSLMode:  "disable",
			},
			expectError: true, // Will fail because we can't connect to real DB in unit tests
		},
		{
			name: "postgresql driver name",
			config: Config{
				Driver:   "postgresql",
				Host:     "localhost",
				Port:     5432,
				Name:     "testdb",
				User:     "testuser",
				Password: "testpass",
				SSLMode:  "disable",
			},
			expectError: true, // Will fail because we can't connect to real DB in unit tests
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := New(tt.config)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
					if db != nil {
						db.Close()
					}
					return
				}
				if tt.errorMsg != "" && err.Error()[:len(tt.errorMsg)] != tt.errorMsg {
					t.Errorf("Expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if db != nil {
					db.Close()
				}
			}
		})
	}
}

func TestClose(t *testing.T) {
	t.Run("close nil database", func(t *testing.T) {
		err := Close(nil)
		if err != nil {
			t.Errorf("Expected no error when closing nil database, got %v", err)
		}
	})

	t.Run("close mock database", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock database: %v", err)
		}

		// Expect close to be called
		mock.ExpectClose()

		err = Close(db)
		if err != nil {
			t.Errorf("Unexpected error when closing database: %v", err)
		}

		// Verify expectations
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})
}

func TestConfig(t *testing.T) {
	t.Run("config struct initialization", func(t *testing.T) {
		cfg := Config{
			Driver:   "postgres",
			Host:     "localhost",
			Port:     5432,
			Name:     "mydb",
			User:     "myuser",
			Password: "mypass",
			SSLMode:  "require",
		}

		if cfg.Driver != "postgres" {
			t.Errorf("Expected driver 'postgres', got %q", cfg.Driver)
		}
		if cfg.Host != "localhost" {
			t.Errorf("Expected host 'localhost', got %q", cfg.Host)
		}
		if cfg.Port != 5432 {
			t.Errorf("Expected port 5432, got %d", cfg.Port)
		}
		if cfg.Name != "mydb" {
			t.Errorf("Expected name 'mydb', got %q", cfg.Name)
		}
		if cfg.User != "myuser" {
			t.Errorf("Expected user 'myuser', got %q", cfg.User)
		}
		if cfg.Password != "mypass" {
			t.Errorf("Expected password 'mypass', got %q", cfg.Password)
		}
		if cfg.SSLMode != "require" {
			t.Errorf("Expected sslmode 'require', got %q", cfg.SSLMode)
		}
	})
}
