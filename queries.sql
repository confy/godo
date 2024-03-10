-- name: GetUsers :many
SELECT * FROM users;

-- name: GetUserById :one 
SELECT * FROM users WHERE id = ?;

-- name: CreateUser :one
INSERT INTO users (username, email) VALUES (?, ?) RETURNING *;

-- name: UpdateUser :exec
UPDATE users SET username = ?, email = ? WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;

-- name: GetTodosByUserId :many
SELECT * FROM todos WHERE user_id = ?;

-- name: GetTodoById :one
SELECT * FROM todos WHERE id = ?;

-- name: CreateTodo :one
INSERT INTO todos (user_id, title, description, done) VALUES (?, ?, ?, ?) RETURNING *;

-- name: UpdateTodo :exec
UPDATE todos SET title = ?, description = ?, done = ? WHERE id = ?;

-- name: DeleteTodo :exec
DELETE FROM todos WHERE id = ?;
