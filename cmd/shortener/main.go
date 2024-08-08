package main

import (
	"net/http"

	"github.com/Dnlbb/url-shortener/cmd/handlers"
	"github.com/Dnlbb/url-shortener/cmd/storage"
)

func main() {
	repo := storage.NewInMemoryStorage() // Создаем новое хранилище
	handler := handlers.NewHandler(repo)
	http.HandleFunc("/", handler.Master)
	err := http.ListenAndServe(`:8080`, nil)
	if err != nil {
		panic(err)
	}
}
