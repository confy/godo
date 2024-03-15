-- name: GetUsers :many
SELECT * FROM users;

-- name: GetUserById :one 
SELECT * FROM users WHERE id = ?;

-- name: GetUserByLogin :one
SELECT * FROM users WHERE login = ?;

-- name: CreateUser :one
INSERT INTO users (login, email) VALUES (?, ?) RETURNING *;

-- name: UpdateUser :exec
UPDATE users SET login = ?, email = ? WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;
