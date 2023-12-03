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

// CreateTodoItem function to easily create a new TodoItem for a User.
func NewTodoItem(userID int, title string) *TodoItem {
	return &TodoItem{
		UserID:    userID,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func GetAllTodoItems(db *sql.DB, userID int) ([]*TodoItem, error) {
	var todoItems []*TodoItem

	query := "SELECT id, user_id, title, completed, created_at, updated_at FROM todos WHERE user_id = ?"

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var todoItem TodoItem
		err := rows.Scan(&todoItem.ID, &todoItem.UserID, &todoItem.Title, &todoItem.Completed, &todoItem.CreatedAt, &todoItem.UpdatedAt)
		if err != nil {
			return nil, err
		}
		todoItems = append(todoItems, &todoItem)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return todoItems, nil
}

// MarkComplete function to mark a TodoItem as completed.
func (t *TodoItem) MarkComplete() {
	t.Completed = true
	t.UpdatedAt = time.Now()
}

// UpdateTitle function to update the title of a TodoItem.
func (t *TodoItem) UpdateTitle(newTitle string) {
	t.Title = newTitle
	t.UpdatedAt = time.Now()
}
