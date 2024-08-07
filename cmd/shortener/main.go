package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

var (
	urlStorage = make(map[string]string)
	mu         sync.RWMutex
)

func main() {
	mux := http.NewServeMux()
	http.HandleFunc("/", fpost)
	http.HandleFunc("/{id}", fget)
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}

func fpost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	originalUrl := string(body)
	shortURL := generateShortUrl(originalUrl)

	mu.Lock()
	urlStorage[shortURL] = originalUrl
	mu.Unlock()

	response := fmt.Sprintf("http://localhost:8080/%s", shortURL)
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(response)))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
}

func generateShortUrl(url string) string {
	hash := sha1.New()
	hash.Write([]byte(url))
	return hex.EncodeToString(hash.Sum(nil))
}

func fget(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
	}

	vars := mux.Vars(r)
	shortURL := vars["id"]

	mu.RLock()
	originalURL, exists := urlStorage[shortURL]
	mu.RUnlock()

	if !exists {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
