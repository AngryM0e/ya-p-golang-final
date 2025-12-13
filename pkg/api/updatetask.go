package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/AngryM0e/ya-p-golang-final/pkg/db"
)

type UpdateTaskRequest struct {
	ID string `json:"id"`
	Date string `json:"date"`
	Title string `json:"title"`
	Comment string `json:"comment"`
	Repeat string `json:"repeat"`
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request, database *db.DB) {
	// Check method
	if r.Method != http.MethodPut {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON request
	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validation required fields
	if req.ID == "" {
		writeJSONError(w, "ID is required", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		writeJSONError(w, "Title is required", http.StatusBadRequest)
		return
	}

	// Convert ID into int
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		writeJSONError(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Check task existing
	if _, err := database.GetTaskByID(id); err != nil {
		writeJSONError(w, "Task not found", http.StatusNotFound)
		return
	}
	
	//Format date
	now := time.Now()
	today := now.Format("20060102")

	dateToUse := req.Date
	if dateToUse == "" {
		dateToUse = today
	}

	// Check date format
	parsedDate, err := time.Parse("20060102", dateToUse)
	if err != nil {
		writeJSONError(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	// Check repeat rules
	if req.Repeat != "" {
		nextDate, err := NextDate(now, dateToUse, req.Repeat)
		if err != nil {
			writeJSONError(w, "Invalid repeat rule", http.StatusBadRequest)
			return
		}
		dateToUse = nextDate
	} else {
		if parsedDate.Before(now) {
			dateToUse = today
		}
	}

	task := db.Task {
		ID: id,
		Date: dateToUse,
		Title: req.Title,
		Comment: req.Comment,
		Repeat: req.Repeat,
	}

	// Update task in BD
	if err := database.UpdateTask(task); err != nil {
		writeJSONError(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	// Send successful empty response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{})
}
