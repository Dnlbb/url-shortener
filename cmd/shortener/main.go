package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sync"
)

var (
	urlStorage = make(map[string]string)
	mu         sync.RWMutex
)

func main() {
	http.HandleFunc("/", master)
	err := http.ListenAndServe(`:8080`, nil)
	if err != nil {
		panic(err)
	}
}

func master(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost:
		fpost(w, r)
	case r.Method == http.MethodGet:
		fget(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func fpost(w http.ResponseWriter, r *http.Request) {

	// contentType := r.Header.Get("Content-Type")
	// if contentType != "text/plain" {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	originalUrl := string(body)
	shortURL := generateShortUrl(originalUrl)

	mu.Lock()
	urlStorage[shortURL] = originalUrl
	mu.Unlock()

	response := fmt.Sprintf("http://localhost:8080/%s", shortURL)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
}

func generateShortUrl(url string) string {
	hash := sha1.New()
	hash.Write([]byte(url))
	return hex.EncodeToString(hash.Sum(nil))[:8]
}

func fget(w http.ResponseWriter, r *http.Request) {

	// contentType := r.Header.Get("Content-Type")
	// if contentType != "text/plain" {
	// 	w.WriteHeader(http.StatusBadRequest)
	// }

	re := regexp.MustCompile(`^/([a-zA-Z0-9]+)$`)
	matches := re.FindStringSubmatch(r.URL.Path)
	if len(matches) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	shortURL := matches[1]

	mu.RLock()
	originalURL, exists := urlStorage[shortURL]
	mu.RUnlock()

	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
