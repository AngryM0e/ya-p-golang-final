package api

import (
	"net/http"

	"github.com/AngryM0e/ya-p-golang-final/pkg/db"
)

const (
	dateFormat = "20060102"
	limit = 50
)

// Init initializes the API routes and handlers
func Init(router *http.ServeMux, database *db.DB) {
	router.HandleFunc("/api/signin", SignInHandler)
	router.HandleFunc("/api/nextdate", nextDayHandler)
	router.HandleFunc("/api/task", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		TaskHandler(w, r, database)
	}))
	router.HandleFunc("/api/tasks", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		tasksHandler(w, r, database)
	}))
	router.HandleFunc("/api/task/done", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		TaskDoneHandler(w, r, database)
	}))
}

// TaskHandler handle requests to /api/task 
func TaskHandler(w http.ResponseWriter, r *http.Request, database *db.DB) {
	switch r.Method {
	case http.MethodPost:
		AddTask(w, r, database)
	case http.MethodGet:
		GetTaskHandler(w, r, database)
	case http.MethodPut:
		UpdateTaskHandler(w, r, database)
	case http.MethodDelete:
		DeleteTaskHandler(w, r, database)
	default:
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}