package health

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestDatabaseChecker_Check_Success(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	mock.ExpectPing()

	checker := NewDatabaseChecker(db)

	if err := checker.Check(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestDatabaseChecker_Check_Error(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	mock.ExpectPing().WillReturnError(sql.ErrConnDone)

	checker := NewDatabaseChecker(db)

	if err := checker.Check(); err == nil {
		t.Error("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestDatabaseChecker_Check_NilDB(t *testing.T) {
	checker := NewDatabaseChecker(nil)

	if err := checker.Check(); err == nil {
		t.Error("expected error for nil db, got nil")
	}
}

func TestDatabaseChecker_Name(t *testing.T) {
	checker := NewDatabaseChecker(nil)

	if name := checker.Name(); name != "database" {
		t.Errorf("expected name 'database', got '%s'", name)
	}
}
