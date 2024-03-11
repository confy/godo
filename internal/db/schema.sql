CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS todos(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id     INTEGER NOT NULL,
    created_at  DATETIME default current_timestamp, 
    updated_at  DATETIME default current_timestamp,
    deleted_at  DATETIME default null,
    title       TEXT,
    description TEXT,
    done        INTEGER(1),
    FOREIGN KEY (user_id) REFERENCES users(id)
);