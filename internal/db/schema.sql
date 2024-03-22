CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    login VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS todos(
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id         INTEGER NOT NULL,
    created_at      DATETIME default current_timestamp, 
    updated_at      DATETIME default current_timestamp,
    completed_at    DATETIME default null,
    title           TEXT NOT NULL,
    description     TEXT,
    done            BOOLEAN default false,
    FOREIGN KEY (user_id) REFERENCES users(id)
);


-- This table should be seperated into a different file
CREATE TABLE IF NOT EXISTS sessions (
	token TEXT PRIMARY KEY,
	data BLOB NOT NULL,
	expiry REAL NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions(expiry);