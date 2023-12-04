package model

import (
	"database/sql"
	"time"
)

// TodoItem with ID, title, completed status, and timestamps.
type TodoItem struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"` // Foreign key to User
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetAllTodoItems function to get all TodoItems for a User of a given userID.
func GetAllTodoItems(db *sql.DB, userID int) ([]*TodoItem, error) {
	var todoItems []*TodoItem

	query := "SELECT id, user_id, title, completed, created_at, updated_at FROM todos WHERE user_id = ?"

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows
	for rows.Next() {
		var todoItem TodoItem
		err := rows.Scan(&todoItem.ID, &todoItem.UserID, &todoItem.Title, &todoItem.Completed, &todoItem.CreatedAt, &todoItem.UpdatedAt)
		if err != nil {
			return nil, err
		}
		todoItems = append(todoItems, &todoItem)
	}

	// Check for errors after we are done iterating over the rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return todoItems, nil
}

// CreateTodoItem function to create a new TodoItem in the database.
// Takes in a userID to ensure that the TodoItem belongs to the User.
func (t *TodoItem) CreateTodoItem(db *sql.DB, userID int) error {
	query := "INSERT INTO todos (user_id, title, completed, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"

	// You can only create a todo item for yourself
	// ? Should we return an error if the user tries to create a todo item for someone else?
	// ? Or should we just ignore the userID in the request body?
	// ? Or should we just return an error if the userID in the request body is not the same as the userID in the request context?
	// Simple approach for now
	t.UserID = userID
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()

	result, err := db.Exec(query, t.UserID, t.Title, t.Completed, t.CreatedAt, t.UpdatedAt)
	if err != nil {
		return err
	}

	// Get the ID of the newly created TodoItem
	todoItemID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Set the ID of the TodoItem to the receiver
	t.ID = int(todoItemID)

	return nil
}

// MarkComplete function to mark a TodoItem as completed.
// Takes in a userID to ensure that the TodoItem belongs to the User.
func (t *TodoItem) MarkComplete(db *sql.DB, userID int) error {
	query := "UPDATE todos SET completed = ?, updated_at = ? WHERE id = ? AND user_id = ?"

	// Update the TodoItem
	// Mark it as completed and update the timestamp
	t.Completed = true
	t.UpdatedAt = time.Now()

	result, err := db.Exec(query, t.Completed, t.UpdatedAt, t.ID, userID)
	if err != nil {
		// Reset the TodoItem to its original state
		t.Completed = false
		t.UpdatedAt = time.Time{}
		return err
	}

	// Check if the TodoItem was actually updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// If no rows were affected, then the TodoItem was not found
	// or it does not belong to the user
	// or it was already deleted
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetTodoItem function to get a TodoItem by its ID for a User of a given userID.
func (t *TodoItem) GetTodoItem(db *sql.DB, userID int) error {
	query := "SELECT id, user_id, title, completed, created_at, updated_at FROM todos WHERE id = ? AND user_id = ?"

	err := db.QueryRow(query, t.ID, userID).Scan(&t.ID, &t.UserID, &t.Title, &t.Completed, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		// Error returned might not be clear enough to the user
		// whether the TodoItem was not found (does not exist or does not belong to the user)
		// or if there was some other error
		return err
	}

	return nil
}

// DeleteTodoItem function to delete a TodoItem by its ID for a User of a given userID.
func (t *TodoItem) DeleteTodoItem(db *sql.DB, userID int) error {
	query := "DELETE FROM todos WHERE id = ? AND user_id = ?"

	result, err := db.Exec(query, t.ID, userID)
	if err != nil {
		return err
	}

	// Check if the TodoItem was actually deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// If no rows were affected, then the TodoItem was not found
	// or it does not belong to the user
	// or it was already deleted
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
