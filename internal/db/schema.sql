CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    login VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS todos(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id     INTEGER NOT NULL,
    created_at  DATETIME default current_timestamp, 
    updated_at  DATETIME default current_timestamp,
    deleted_at  DATETIME default null,
    title       TEXT NOT NULL,
    description TEXT,
    done        INTEGER(1),
    FOREIGN KEY (user_id) REFERENCES users(id)
);