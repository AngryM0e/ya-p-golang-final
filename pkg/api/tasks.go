package api

import (
	"net/http"
	"strconv"

	"github.com/AngryM0e/ya-p-golang-final/pkg/db"
)

// TaskResponse - struct for API response
type TaskResponse struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title,omitempty"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

// TasksResponse - struct for API response
type TasksResponse struct {
	Tasks []TaskResponse `json:"tasks"`
}

// ErrorResponse struct for error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// convertTask - convert task from DB to API format
func convertTask(task db.Task) TaskResponse {
	return TaskResponse{
		ID:      strconv.Itoa(task.ID),
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}
}

// tasksHandler - handler for GET /api/tasks (без поиска)
func tasksHandler(w http.ResponseWriter, r *http.Request, database *db.DB) {
	// Check request method
	if r.Method != http.MethodGet {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Getting tasks from database with limit 50
	tasks, err := database.GetAllTasks(50)
	if err != nil {
		writeJSONError(w, "Error while getting tasks", http.StatusInternalServerError)
		return
	}

	// Convert tasks to API format
	taskResponse := make([]TaskResponse, 0, len(tasks))
	for _, task := range tasks {
		taskResponse = append(taskResponse, convertTask(task))
	}

	// Create response
	response := TasksResponse{
		Tasks: taskResponse,
	}

	if response.Tasks == nil {
		response.Tasks = []TaskResponse{}
	}
	
	// Send response
	writeJSONSuccess(w, response)
}