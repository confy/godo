package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/alexedwards/scs/v2"
	"github.com/confy/godo/internal/db"
)


func extractUserAndTodoID(session *scs.SessionManager, r *http.Request) (userID, todoID int64, err error){
    userID = session.GetInt64(r.Context(), "userID")
    if userID == 0 {
        return 0, 0, errors.New("UserID not found in session")
    }
	todoID, err = strconv.ParseInt(r.PathValue("id"), 10, 64)
    if err != nil {
        return 0, 0, errors.New("Invalid todo ID")
    }
    return userID, todoID, nil
}


func HandleDeleteTodo(database *db.Queries, session *scs.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        userID, todoID, err := extractUserAndTodoID(session, r)
        if err != nil {
            errorHandler(w, r, err.Error(), http.StatusBadRequest)
            return
        }

        todo, err := database.GetTodoById(r.Context(), todoID)
        if err != nil {
            errorHandler(w, r, "Failed to get todo", http.StatusInternalServerError)
            return
        }

        if todo.UserID != userID {
            errorHandler(w, r, "Unauthorized", http.StatusUnauthorized)
            return
        }


        err = database.DeleteTodo(r.Context(), todoID)
        if err != nil {
            errorHandler(w, r, "Failed to delete todo", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
        w.Write([]byte(""))
	}
}
