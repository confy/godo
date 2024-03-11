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
