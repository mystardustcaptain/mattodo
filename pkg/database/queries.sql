--- tables used in the database

CREATE TABLE todos (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    completed BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
)
CREATE INDEX idx_todos_user_id ON todos(user_id);

CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    oauth_provider TEXT NOT NULL,
    oauth_id TEXT NOT NULL,
    name TEXT NOT NULL,
    email TEXT NOT NULL
)