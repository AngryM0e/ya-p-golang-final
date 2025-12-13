package api

import (
	"net/http"
	"strconv"

	"github.com/AngryM0e/ya-p-golang-final/pkg/db"
)
func GetTaskHandler(w http.ResponseWriter, r *http.Request, database *db.DB) {
	// Get ID from URL
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSONError(w, "ID is required", http.StatusBadRequest)
		return
	}
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, "Invalid task ID", http.StatusBadRequest)
		return
	}
	
	// Get task from database
	task, err := database.GetTaskByID(id)
	if err != nil {
		// Check if task not found
		if err.Error() == "task not found" {
			writeJSONError(w, "Task not found", http.StatusNotFound)
		} else {
			writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	
	// Create succsesfull response
	response := map[string]string{
		"id":      strconv.Itoa(task.ID),
		"date":    task.Date,
		"title":   task.Title,
		"comment": task.Comment,
		"repeat":  task.Repeat,
	}
	
	writeJSONSuccess(w, response)
}
