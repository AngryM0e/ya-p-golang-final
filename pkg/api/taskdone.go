package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/AngryM0e/ya-p-golang-final/pkg/db"
)

// TaskDoneHandler handles the task done endpoint
func TaskDoneHandler(w http.ResponseWriter, r *http.Request, database *db.DB) {
	if r.Method != http.MethodPost {
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

	// Get task from DB
	task, err := database.GetTaskByID(id)
	if err != nil {
		if err.Error() == "Task not found" {
			writeJSONError(w, "Task not found", http.StatusNotFound)
			return
		} else {
			writeJSONError(w, "Getting task error", http.StatusInternalServerError)
			return
		}
	}

	// If task is not repeatable, delete it
	if task.Repeat == "" {
		if err := database.DeleteTask(id); err != nil {
			writeJSONError(w, "Deleting task error", http.StatusInternalServerError)
			return
		}
		// Return empty response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{})
		return
	}

	// If task is repeatable, update it
	now := time.Now()
	nextDate, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeJSONError(w, "Error calculating next date", http.StatusInternalServerError)
		return
	}

	// Update task date
	if err := database.UpdateTaskDate(id, nextDate); err != nil {
		writeJSONError(w, "Updating task error", http.StatusInternalServerError)
		return
	}
	writeJSONEmpty(w)
}