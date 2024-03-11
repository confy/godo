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
