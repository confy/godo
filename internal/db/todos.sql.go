// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: todos.sql

package db

import (
	"context"
	"database/sql"
)

const createTodo = `-- name: CreateTodo :one
INSERT INTO todos (user_id, title, description, done) VALUES (?, ?, ?, ?) RETURNING id, user_id, created_at, updated_at, deleted_at, title, description, done
`

type CreateTodoParams struct {
	UserID      int64          `json:"user_id"`
	Title       string         `json:"title"`
	Description sql.NullString `json:"description"`
	Done        sql.NullBool   `json:"done"`
}

func (q *Queries) CreateTodo(ctx context.Context, arg CreateTodoParams) (Todo, error) {
	row := q.db.QueryRowContext(ctx, createTodo,
		arg.UserID,
		arg.Title,
		arg.Description,
		arg.Done,
	)
	var i Todo
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Title,
		&i.Description,
		&i.Done,
	)
	return i, err
}

const deleteTodo = `-- name: DeleteTodo :exec
DELETE FROM todos WHERE id = ?
`

func (q *Queries) DeleteTodo(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteTodo, id)
	return err
}

const getTodoById = `-- name: GetTodoById :one
SELECT id, user_id, created_at, updated_at, deleted_at, title, description, done FROM todos WHERE id = ?
`

func (q *Queries) GetTodoById(ctx context.Context, id int64) (Todo, error) {
	row := q.db.QueryRowContext(ctx, getTodoById, id)
	var i Todo
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Title,
		&i.Description,
		&i.Done,
	)
	return i, err
}

const getTodosByUserId = `-- name: GetTodosByUserId :many
SELECT id, user_id, created_at, updated_at, deleted_at, title, description, done FROM todos WHERE user_id = ?
`

func (q *Queries) GetTodosByUserId(ctx context.Context, userID int64) ([]Todo, error) {
	rows, err := q.db.QueryContext(ctx, getTodosByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Todo
	for rows.Next() {
		var i Todo
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.Title,
			&i.Description,
			&i.Done,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateTodo = `-- name: UpdateTodo :exec
UPDATE todos SET title = ?, description = ?, done = ? WHERE id = ?
`

type UpdateTodoParams struct {
	Title       string         `json:"title"`
	Description sql.NullString `json:"description"`
	Done        sql.NullBool   `json:"done"`
	ID          int64          `json:"id"`
}

func (q *Queries) UpdateTodo(ctx context.Context, arg UpdateTodoParams) error {
	_, err := q.db.ExecContext(ctx, updateTodo,
		arg.Title,
		arg.Description,
		arg.Done,
		arg.ID,
	)
	return err
}
