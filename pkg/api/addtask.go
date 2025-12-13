package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/AngryM0e/ya-p-golang-final/pkg/db"
)

func checkDate(task *db.Task) error {
	now := time.Now()
	today := now.Format("20060102")
	
	if task.Date == "" {
		task.Date = today
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return err
	}

	if task.Repeat != "" {
		parts := strings.Fields(task.Repeat)
		if len(parts) == 0 {
			return errors.New("invalid format repeat rule")
		}

		// Check dayli rule
		if parts[0] == "d" {
			if len(parts) != 2 {
				return errors.New("invalid format dayli rule")
			}
			
			// Parse days interval
			interval, err := strconv.Atoi(parts[1])
			if err != nil {
				return errors.New("invalid days interval")
			}
			
			if interval <= 0 {
				return errors.New("interval must be more 0")
			}
			
			if interval > 400 {
				return errors.New("interval must be less 400")
			}
		}
		
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
		
		if afterNow(now, t) {
			task.Date = next
		}
	} else {
		// If task not repeatable, check if it's in the past
		if afterNow(now, t) {
			task.Date = today
		}
	}

	return nil
}

func AddTask(w http.ResponseWriter, r *http.Request, database *db.DB) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONError(w,"JSON decode error", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJSONError(w, "Title is required", http.StatusBadRequest)
		return
	}
	if err := checkDate(&task); err != nil {
		writeJSONError(w, "Invalid date", http.StatusBadRequest)
		return
	}
	
	id, err := database.AddTask(task)
	if err != nil {
		writeJSONError(w, "Adding task error", http.StatusInternalServerError)
		return
	}
	writeJSONSuccess(w, map[string]int64{"id": id})
}

// afterNow checks if a date is strictly after another date (ignoring time)
func afterNow(date, now time.Time) bool {
	// Compare only dates (year, month, day)
	dateYear, dateMonth, dateDay := date.Date()
	nowYear, nowMonth, nowDay := now.Date()

	if dateYear > nowYear {
		return true
	}
	if dateYear < nowYear {
		return false
	}

	if dateMonth > nowMonth {
		return true
	}
	if dateMonth < nowMonth {
		return false
	}

	return dateDay > nowDay
}