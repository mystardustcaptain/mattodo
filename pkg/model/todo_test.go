package model_test

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mystardustcaptain/mattodo/pkg/model"
	"github.com/stretchr/testify/assert"
)

// TestGetAllTodoItems_ExecuteCorrectQuery tests that GetAllTodoItems executes the correct query,
// returns the correct number of TodoItems,
// and returns the correct TodoItems data,
// and handle errors correctly.
func TestGetAllTodoItems_ExecuteCorrectQuery(t *testing.T) {
	/// Arrange
	///
	// Create a new instance of sqlmock
	db, mock, errdb := sqlmock.New()
	if errdb != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", errdb)
	}
	defer db.Close()

	columns := []string{"id", "user_id", "title", "completed", "created_at", "updated_at"}
	mock.ExpectQuery("SELECT id, user_id, title, completed, created_at, updated_at FROM todos WHERE user_id = ?").
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows(columns).
			AddRow(2, 2, "Todo 2", false, time.Now(), time.Now()).
			AddRow(3, 2, "Todo 3", false, time.Now(), time.Now()))

	/// Act
	///
	// call GetAllTodoItems and pass the mocked db instance
	todos, err := model.GetAllTodoItems(db, 2)
	if err != nil {
		t.Errorf("error was not expected while getting todo items: %s", err)
	}

	/// Assert
	///
	assert.Equal(t, 2, len(todos))
	assert.NoError(t, err, "Expected no error but got one")
	assert.Equal(t, 2, todos[0].ID)
	assert.Equal(t, 2, todos[0].UserID)
	assert.Equal(t, "Todo 2", todos[0].Title)
	assert.Equal(t, false, todos[0].Completed)
	assert.Equal(t, 3, todos[1].ID)
	assert.Equal(t, 2, todos[1].UserID)
	assert.Equal(t, "Todo 3", todos[1].Title)
	assert.Equal(t, false, todos[1].Completed)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestGetAllTodoItems_ReturnNothingWhenNoTodoItems tests that GetAllTodoItems returns nothing when there are error with the query.
// It is to ensure that the error is handled correctly,
// and that the function does not return any TodoItems,
// and that the error is the expected error.
func TestGetAllTodoItems_ReturnNothingWhenQueryError(t *testing.T) {
	/// Arrange
	///
	// Create a new instance of sqlmock
	db, mock, errdb := sqlmock.New()
	if errdb != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", errdb)
	}
	defer db.Close()

	// Define a custom error
	customErr := errors.New("mock database connection error")

	mock.ExpectQuery("SELECT id, user_id, title, completed, created_at, updated_at FROM todos WHERE user_id = ?").
		WithArgs(2).
		WillReturnError(customErr)

	/// Act
	///
	// call GetAllTodoItems and pass the mocked db instance
	todos, err := model.GetAllTodoItems(db, 2)

	/// Assert
	///
	assert.Error(t, err, "Expected an error but got none")
	assert.Equal(t, customErr, err, "Expected a different error")
	assert.Nil(t, todos, "Expected todos to be nil on error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestGetAllTodoItems_ReturnNothingWhenScanError tests that GetAllTodoItems returns nothing when there are error during scanning rows.
// It is to ensure that the error is handled correctly,
// and that the function does not return any TodoItems,
// and that the error is the expected error.
func TestGetAllTodoItems_ReturnNothingWhenScanError(t *testing.T) {
	/// Arrange
	// Create a new instance of sqlmock
	db, mock, errdb := sqlmock.New()
	if errdb != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", errdb)
	}
	defer db.Close()

	// Define a custom error
	customErr := errors.New("sql: Scan error on column index 5, name \"updated_at\": unsupported Scan, storing driver.Value type string into type *time.Time")

	columns := []string{"id", "user_id", "title", "completed", "created_at", "updated_at"}
	mock.ExpectQuery("SELECT id, user_id, title, completed, created_at, updated_at FROM todos WHERE user_id = ?").
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows(columns).
			AddRow(2, 2, "Todo 2", false, time.Now(), time.Now()).
			AddRow(3, 2, "Todo 3", false, time.Now(), "hi")) // This will cause an error due to the wrong type

	/// Act
	// call GetAllTodoItems and pass the mocked db instance
	todos, err := model.GetAllTodoItems(db, 2)

	/// Assert
	///
	assert.Error(t, err, "Expected an error but got none")
	assert.Equal(t, customErr.Error(), err.Error(), "Expected a different error")
	assert.Nil(t, todos, "Expected todos to be nil on error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
