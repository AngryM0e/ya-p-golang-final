package main

import (
	"log"

	"github.com/AngryM0e/ya-p-golang-final/pkg/server"
	"github.com/AngryM0e/ya-p-golang-final/pkg/filechecker"
)

const (
	defaultPort = 7540
	webDir      = "web"
)

func main() {
	// Check exist scheduler.db
	dbFilePath := "scheduler.db"
	if !filechecker.CheckFileExists(dbFilePath) {
		log.Fatalf("Файл %s не найден. Сервер не может быть запущен.", dbFilePath)
	}

	// Create & config server
	router, err := server.NewServer(defaultPort, webDir)
	if err != nil {
		log.Fatal("Ошибка создания сервера:", err)
	}

	// Server start
	if err := server.Start(router, defaultPort); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
