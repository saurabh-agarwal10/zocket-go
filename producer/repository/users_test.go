package repository

import (
	"errors"
	"producer/config"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"go.uber.org/zap"
)

func TestGetUserByID(t *testing.T) {
	// Create a new mock DB and a logger for testing
	db, mock, _ := sqlmock.New()
	defer db.Close()

	zapLogger := zap.NewNop()

	// Replace the global MySQL connection and logger with the mock instances
	config.MysqlConnection = db
	zap.ReplaceGlobals(zapLogger)

	// Test case for user exists
	query := "SELECT EXISTS(.+)"

	mock.ExpectQuery(query).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"EXISTS"}).AddRow(1))

	userExists, err := GetUserByID(1)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	if !userExists {
		t.Fatal("Expected user to exist, but it doesn't.")
	}

	// Test case for user does not exist
	mock.ExpectQuery(query).WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"EXISTS"}).AddRow(0))

	userExists, err = GetUserByID(2)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	if userExists {
		t.Fatal("Expected user not to exist, but it does.")
	}

	// Test case for database error
	mock.ExpectQuery(query).WithArgs(3).WillReturnError(errors.New("Database error"))

	_, err = GetUserByID(3)
	if err == nil {
		t.Fatal("Expected an error, but got nil")
	}
}

func TestMain(m *testing.M) {
	// Hook up zap logging for testing
	zapLogger := zap.NewNop()
	zap.ReplaceGlobals(zapLogger)

	m.Run()
}
