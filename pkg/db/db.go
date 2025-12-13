package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

// Task - struct for database operations
type Task struct {
	ID      int    `json:"id"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title,omitempty"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

// DB - struct for DB connection
type DB struct {
	db *sql.DB
}

// Close - close connection with DB
func (d *DB) Close() error {
	return d.db.Close()
}

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT "19700101",
	title VARCHAR(256) NOT NULL DEFAULT "",
	comment TEXT NOT NULL DEFAULT "",
	repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX IF NOT EXISTS idx_date_scheduler ON scheduler(date);
`

// Init - initialize DB connection
func Init(dbFile string) (*DB, error) {
	log.Printf("Initializing database: %s", dbFile)
	
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}

	// Всегда создаем таблицы если их нет
	_, err = db.Exec(schema)
	if err != nil {
		db.Close()
		return nil, err
	}
	
	log.Printf("DB initialized successfully")
	return &DB{db: db}, nil
}

// AddTask add task to database
func (d *DB) AddTask(task Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	result, err := d.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetTaskByID get task by id from database
func (d *DB) GetTaskByID(id int) (Task, error) {
	var task Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err := d.db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return Task{}, fmt.Errorf("task not found")
		}
		return Task{}, err
	}
	return task, nil
}

// UpdateTask update task in database
func (d *DB) UpdateTask(task Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	result, err := d.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

// GetAllTasks gets all tasks
func (d *DB) GetAllTasks(limit int) ([]Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler
		ORDER BY date ASC, id ASC LIMIT ?`
	
	rows, err := d.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	if tasks == nil {
		tasks = []Task{}
	}
	return tasks, nil
}

// DeleteTask delete task by ID
func (d *DB) DeleteTask(id int) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	result, err := d.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no task found with ID %d", id)
	}
	return nil
}

// UpdateTaskDate update only task date
func (d *DB) UpdateTaskDate(id int, date string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	result, err := d.db.Exec(query, date, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no task found with ID %d", id)
	}
	return nil
}