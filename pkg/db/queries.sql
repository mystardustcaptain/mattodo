--- tables used in the database

CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    oauth_provider TEXT,
    oauth_id TEXT,
    name TEXT,
    email TEXT
);

CREATE TABLE todos (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    title TEXT,
    completed BOOLEAN,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
