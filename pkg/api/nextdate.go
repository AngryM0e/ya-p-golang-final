package api

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// NextDate calculates the next execution date for a task
func NextDate(now time.Time, dateStr string, repeat string) (string, error) {
	// Parse the start date
	date, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		return "", errors.New("invalid start date format")
	}

	// Check if repeat rule is empty
	if repeat == "" {
		return "", errors.New("empty repeat rule")
	}

	// Split the repeat rule into parts
	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("invalid repeat rule format")
	}

	// Process the repeat rule based on its type
	ruleType := parts[0]

	switch ruleType {
	case "d":
		// d <number> - repeat every <number> days
		return handleDailyRule(now, date, parts)
	case "y":
		// y - repeat yearly
		return handleYearlyRule(now, date, parts)
	case "w":
		// w <weekdays> - repeat weekly on specified weekdays
		return handleWeeklyRule(now, date, parts)
	case "m":
		// m <monthdays> [<months>] - repeat monthly on specified days
		return handleMonthlyRule(now, date, parts)
	default:
		return "", errors.New("unsupported repeat rule")
	}
}

// handleDailyRule processes daily repeat rule: d <number>
func handleDailyRule(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", errors.New("invalid daily rule format")
	}

	interval, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", errors.New("invalid day interval")
	}

	if interval <= 0 || interval > 400 {
		return "", errors.New("invalid day interval")
	}

	// Start from the given date
	currentDate := date

	// Find the next date after 'now' by adding intervals
	for {
		currentDate = currentDate.AddDate(0, 0, interval)
		if afterNow(currentDate, now) {
			break
		}
	}

	return currentDate.Format(dateFormat), nil
}

// handleYearlyRule processes yearly repeat rule: y
func handleYearlyRule(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 1 {
		return "", errors.New("invalid yearly rule format")
	}

	// Start from the given date
	currentDate := date

	// Find the next date after 'now' by adding years
	for {
		currentDate = currentDate.AddDate(1, 0, 0)
		if afterNow(currentDate, now) {
			break
		}
	}

	return currentDate.Format(dateFormat), nil
}

// handleWeeklyRule processes weekly repeat rule: w <weekdays>
func handleWeeklyRule(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", errors.New("invalid weekly rule format")
	}

	// Parse weekdays
	dayStrs := strings.Split(parts[1], ",")
	if len(dayStrs) == 0 {
		return "", errors.New("weekdays are required")
	}

	// Create array for weekdays (1-7, where 1=Monday, 7=Sunday)
	var weekdays [8]bool // indices 1-7, 0 is unused

	for _, dayStr := range dayStrs {
		day, err := strconv.Atoi(dayStr)
		if err != nil {
			return "", errors.New("invalid weekday")
		}
		if day < 1 || day > 7 {
			return "", errors.New("weekday must be between 1 and 7")
		}
		weekdays[day] = true
	}

	// Start from the given date
	currentDate := date

	// Check day by day using AddDate(0, 0, 1)
	for i := 0; i < 1000; i++ { // protection against infinite loop
		// Get weekday (1-7)
		weekday := int(currentDate.Weekday())
		if weekday == 0 {
			weekday = 7 // Sunday
		}

		// Check if weekday matches and date is after 'now'
		if afterNow(currentDate, now) && weekdays[weekday] {
			return currentDate.Format(dateFormat), nil
		}

		// Move to next day
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return "", errors.New("cannot find next date for weekdays")
}

// handleMonthlyRule processes monthly repeat rule: m <days> [<months>]
func handleMonthlyRule(now, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 || len(parts) > 3 {
		return "", errors.New("invalid monthly rule format")
	}

	// Parse days of the month
	dayStrs := strings.Split(parts[1], ",")
	if len(dayStrs) == 0 {
		return "", errors.New("days of month are required")
	}

	// Create array for days of month (1-31)
	var days [32]bool // indices 1-31, 0 is unused
	var hasNegativeDays bool
	var negativeDays []int

	for _, dayStr := range dayStrs {
		day, err := strconv.Atoi(dayStr)
		if err != nil {
			return "", errors.New("invalid day of month")
		}

		// Check valid day range
		if day < -2 || day == 0 || day > 31 {
			return "", errors.New("invalid day of month")
		}

		if day > 0 {
			// Positive days: mark in array
			days[day] = true
		} else {
			// Negative days: store for special handling
			hasNegativeDays = true
			negativeDays = append(negativeDays, day)
		}
	}

	// Create array for months (1-12)
	var months [13]bool // indices 1-12, 0 is unused

	if len(parts) == 3 {
		// Parse specified months
		monthStrs := strings.Split(parts[2], ",")
		if len(monthStrs) == 0 {
			return "", errors.New("months are required")
		}
		for _, monthStr := range monthStrs {
			month, err := strconv.Atoi(monthStr)
			if err != nil || month < 1 || month > 12 {
				return "", errors.New("invalid month")
			}
			months[month] = true
		}
	} else {
		// If months not specified, use all months
		for i := 1; i <= 12; i++ {
			months[i] = true
		}
	}

	// Start from the given date
	currentDate := date

	// Check day by day using AddDate(0, 0, 1)
	for i := 0; i < 2000; i++ { // protection against infinite loop (~5.5 years)
		currentMonth := int(currentDate.Month())
		currentDay := currentDate.Day()

		// Check month
		if !months[currentMonth] {
			currentDate = currentDate.AddDate(0, 0, 1)
			continue
		}

		// Check day
		dayValid := false

		// Check positive days from array
		if currentDay >= 1 && currentDay <= 31 && days[currentDay] {
			dayValid = true
		}

		// Check negative days if any
		if !dayValid && hasNegativeDays {
			// Get last day of current month
			lastDayOfMonth := time.Date(
				currentDate.Year(),
				currentDate.Month()+1,
				0, 0, 0, 0, 0,
				currentDate.Location(),
			).Day()

			// Check each negative day
			for _, negDay := range negativeDays {
				// Calculate actual day for negative value
				actualDay := lastDayOfMonth + negDay + 1
				if actualDay >= 1 && actualDay == currentDay {
					dayValid = true
					break
				}
			}
		}

		// If day is valid and date is after 'now'
		if dayValid && afterNow(currentDate, now) {
			return currentDate.Format(dateFormat), nil
		}

		// Move to next day
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return "", errors.New("cannot find next date")
}

// nextDayHandler handles GET requests to /api/nextdate
func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get parameters from query string
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	// If 'now' is not specified, use current date
	if nowStr == "" {
		nowStr = time.Now().Format(dateFormat)
	}

	// Validate required parameters
	if dateStr == "" {
		writeJSONError(w, "Date parameter is required", http.StatusBadRequest)
		return
	}

	if repeat == "" {
		writeJSONError(w, "Repeat parameter is required", http.StatusBadRequest)
		return
	}

	// Parse dates
	nowTime, err := time.Parse(dateFormat, nowStr)
	if err != nil {
		writeJSONError(w, "Invalid date format for 'now' parameter", http.StatusBadRequest)
		return
	}

	// Calculate next date
	nextDate, err := NextDate(nowTime, dateStr, repeat)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return result as plain text
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := w.Write([]byte(nextDate)); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
	w.Write([]byte(nextDate))
}
