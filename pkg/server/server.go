package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)
// NewServer create & config HTTP-router
func NewServer(port int, webDir string) (*http.ServeMux, error) {
	// Create new router
	router := http.NewServeMux()

	// Check exist dir web & file index.html
	webPath := filepath.Join(webDir, "index.html")
	if _, err := os.Stat(webPath); os.IsNotExist(err) {
		log.Printf("Внимание: файл %s не найден\n", webPath)
		log.Println("Создайте директорию 'web' поместите в неё index.html")
	}

	// Настраиваем обработчик для статических файлов
	router.Handle("/", http.FileServer(http.Dir(webDir)))

	return router, nil
}

func Start (router *http.ServeMux, port int) error {
	addr := fmt.Sprintf(":%d", port)
	log.Printf("Сервер запущен на порту %d", port)
	return http.ListenAndServe(addr, router)
}