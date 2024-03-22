package models

import (
	"strconv"

	"github.com/confy/godo/internal/db"
)

type DisplayUser struct {
	LoggedIn  bool
	Login     string
	Email     string
	AvatarURL string
}

func DisplayUserFromUser(user db.User) DisplayUser {
	return DisplayUser{
		LoggedIn:  true,
		Login:     user.Login,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
	}
}

type DisplayTodo struct {
	ID          int
	DOMID       string
	Route       string
	Target      string
	UserID      int
	CreatedAt   string
	UpdatedAt   string
	CompletedAt string
	Title       string
	Description string
	Done        bool
}

func DisplayTodoFromTodo(todo db.Todo) DisplayTodo {
	strID := strconv.Itoa(int(todo.ID))

	return DisplayTodo{
		ID:          int(todo.ID),
		DOMID:       "todo-" + strID,
		Route:       "/todo/" + strID,
		Target:      "#" + strID,
		UserID:      int(todo.UserID),
		CreatedAt:   todo.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:   todo.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		CompletedAt: todo.CompletedAt.Time.Format("2006-01-02 15:04:05"),
		Title:       todo.Title,
		Description: todo.Description.String,
		Done:        todo.Done.Bool,
	}
}
