package api

import (
	"net/http"
	"strconv"

	"github.com/AngryM0e/ya-p-golang-final/pkg/db"
)

// DeleteTaskHandler handles the task deletion endpoint
func DeleteTaskHandler (w http.ResponseWriter, r *http.Request, database *db.DB) {
	// Check if the method is DELETE
	if r.Method != http.MethodDelete {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get ID from URL
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSONError(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	// Convert ID to int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	// Check if task exists
	if _, err := database.GetTaskByID(id); err != nil {
		if err.Error() == "task not found" {
			writeJSONError(w, "Task not found", http.StatusNotFound)
			return
		} else {
			writeJSONError(w, "Error with checking task existence in DB with ID", http.StatusInternalServerError)
			return
		}
	}
	
	// Delete task from DB
	if err := database.DeleteTask(id); err != nil {
		writeJSONError(w, "Error with deleting task from DB", http.StatusInternalServerError)
		return
	}

	// Response with status OK
	writeJSONEmpty(w)
}
