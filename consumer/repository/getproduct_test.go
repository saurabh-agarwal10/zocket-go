package repository

import (
	"consumer/config"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"go.uber.org/zap"
)

func TestGetProductFromID(t *testing.T) {
	// Initialize a new logger for testing
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	// Set up the database connection in the config
	config.MysqlConnection = db

	// Create a test product ID
	productID := 1

	// Define the expected query and results for the mock
	query := "SELECT product_images FROM products WHERE product_id=?"

	// Test case 1: No error
	rows := sqlmock.NewRows([]string{"product_images"}).AddRow("image1,image2,image3")
	mock.ExpectQuery(query).WithArgs(productID).WillReturnRows(rows)

	images, err := GetProductFromID(productID)
	if err != nil {
		t.Errorf("Expected no error, got error: %v", err)
	}
	if images != "image1,image2,image3" {
		t.Errorf("Expected 'image1,image2,image3', got %s", images)
	}

	// Test case 2: SQL no rows error
	mock.ExpectQuery(query).WithArgs(productID).WillReturnError(sql.ErrNoRows)
	images, err = GetProductFromID(productID)
	if err == nil {
		t.Error("Expected an error, got nil")
	}
	if err.Error() != "no matching id found" {
		t.Errorf("Expected error message 'no matching id found', got '%s'", err.Error())
	}

	// Test case 3: Database error
	mock.ExpectQuery(query).WithArgs(productID).WillReturnError(errors.New("Database error"))
	images, err = GetProductFromID(productID)
	if err == nil {
		t.Error("Expected an error, got nil")
	}
	if err.Error() != "Database error" {
		t.Errorf("Expected error message 'Database error', got '%s'", err.Error())
	}
}

func TestUpdateProductImage(t *testing.T) {
	// Initialize a new logger for testing
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	// Set up the database connection in the config
	config.MysqlConnection = db

	// Create a test product ID and path
	productID := 1
	path := "path/to/images"

	// Define the expected query and results for the mock
	query := "UPDATE products SET compressed_product_images=? WHERE product_id=?"

	// Test case 1: No error
	mock.ExpectExec(query).WithArgs(path, productID).WillReturnResult(sqlmock.NewResult(1, 1))

	err = UpdateProductImage(productID, path)
	if err != nil {
		t.Errorf("Expected no error, got error: %v", err)
	}

	// Test case 2: Database error
	mock.ExpectExec(query).WithArgs(path, productID).WillReturnError(errors.New("Database error"))

	err = UpdateProductImage(productID, path)
	if err == nil {
		t.Error("Expected an error, got nil")
	}
	if err.Error() != "Database error" {
		t.Errorf("Expected error message 'Database error', got '%s'", err.Error())
	}
}
