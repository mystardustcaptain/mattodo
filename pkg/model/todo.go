package model

import "time"

// TodoItem with ID, title, completed status, and timestamps.
type TodoItem struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"` // Foreign key to User
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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
