package main

import (
	"fmt"
	"net/http"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/AngryM0e/ya-p-golang-final/pkg/db"
	"github.com/AngryM0e/ya-p-golang-final/pkg/server"
)

const (
	defaultPort   = 7540
	webDir        = "web"
	defaultDBfile = "scheduler.db"
)

func main() {
	port := getPort()
	dbPath := getAbsolutePath(defaultDBfile)
	
	log.Printf("Starting server on port %d", port)
	log.Printf("Using database: %s", dbPath)
	
	cfg := server.Config{
		Port:   port,
		WebDir: webDir,
		DBPath: dbPath,
	}

	// Create & config server
	router, database, err := server.NewServer(cfg)
	if err != nil {
		log.Fatal("Create server error:", err)
	}

	// Launch server
	if err := Start(router, database, port); err != nil {
		log.Fatal("Server start error:", err)
	}
}

// getPort - getting port from env or using default
func getPort() int {
	portStr := os.Getenv("TODO_PORT")
	if portStr == "" {
		return defaultPort
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Printf("Invalid TODO_PORT, using default: %d", defaultPort)
		return defaultPort
	}

	if port < 1 || port > 9999 {
		log.Printf("Port %d out of range, using default: %d", port, defaultPort)
		return defaultPort
	}

	return port
}

// getAbsolutePath - getting absolute path from env or using default
func getAbsolutePath(filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}
	
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return filename
	}
	return absPath
}

func Start (router *http.ServeMux, database *db.DB, port int) error {
	// Close database connection on exit
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		} else {
			log.Println("Database connection closed")
		}
	}()

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Server launch on http://localhost:%d", port)
	
	return http.ListenAndServe(addr, router)
}