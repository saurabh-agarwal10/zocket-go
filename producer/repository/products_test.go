package repository

import (
	"errors"
	"producer/config"
	"producer/dto"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"go.uber.org/zap"
)

func TestInsertProduct(t *testing.T) {
	// Create a new mock DB and a logger for testing
	db, mock, _ := sqlmock.New()
	defer db.Close()

	zapLogger := zap.NewNop()

	// Replace the global MySQL connection and logger with the mock instances
	config.MysqlConnection = db
	zap.ReplaceGlobals(zapLogger)

	// Test data
	testProduct := &dto.RequestData{
		ProductName:        "Test Product",
		ProductDescription: "Description",
		ProductImages:      "image1.jpg,image2.jpg",
		ProductPrice:       10.0,
	}

	// Expected SQL query and result
	query := "INSERT INTO products(.+) VALUES(.+)"
	mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))

	// Run the function
	productID, err := InsertProduct(testProduct)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	// Check the result
	if productID != 1 {
		t.Fatalf("Expected productID 1, but got: %d", productID)
	}

	// Test case for database error
	mock.ExpectExec(query).WillReturnError(errors.New("Database error"))

	_, err = InsertProduct(testProduct)
	if err == nil {
		t.Fatalf("Expected an error, but got nil")
	}
}
