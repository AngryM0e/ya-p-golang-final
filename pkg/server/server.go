package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/AngryM0e/ya-p-golang-final/pkg/api"
	"github.com/AngryM0e/ya-p-golang-final/pkg/db"
)

type Config struct {
	Port int
	WebDir string
	DBPath string
}

// NewServer create & config HTTP-router
func NewServer(cfg Config) (*http.ServeMux, *db.DB, error) {
	// Create new router
	router := http.NewServeMux()

	// Check exist dir web & file index.html
	webPath := filepath.Join(cfg.WebDir, "index.html")
	if _, err := os.Stat(webPath); os.IsNotExist(err) {
		log.Printf("Directory doesn't exist", cfg.WebDir)
		log.Println("Create directory 's' and put static files there", cfg.WebDir)
	} else {
		// Check index.html
		webPath := filepath.Join(cfg.WebDir, "index.html")
		if _, err := os.Stat(webPath); os.IsNotExist(err) {
			log.Printf("File dosn't exist: %s\n", webPath)
			log.Printf("Create file index.html in directory %s", cfg.WebDir)
		}
	}

	// Configure static files handler
	router.Handle("/", http.FileServer(http.Dir(cfg.WebDir)))

	database, err := db.Init(cfg.DBPath)
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка инициализации БД: %w", err)
	}

	// Initialize API
	api.Init(router, database)

	return router, database, nil
}
