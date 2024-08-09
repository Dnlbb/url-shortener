package handlers

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"

	"github.com/Dnlbb/url-shortener/cmd/storage"
)

type Handler struct {
	repo storage.Repository
}

func NewHandler(repo storage.Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Master(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost:
		h.Fpost(w, r)
	case r.Method == http.MethodGet:
		h.Fget(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h *Handler) Fpost(w http.ResponseWriter, r *http.Request) {

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

	originalURL := string(body)

	parsedURL, err := url.ParseRequestURI(originalURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURL := GenerateShortURL(originalURL)

	h.repo.Save(shortURL, originalURL)

	response := fmt.Sprintf("http://localhost:8080/%s", shortURL)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
}

func GenerateShortURL(url string) string {
	hash := sha1.New()
	hash.Write([]byte(url))
	return hex.EncodeToString(hash.Sum(nil))[:8]
}

func (h *Handler) Fget(w http.ResponseWriter, r *http.Request) {

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

	originalURL, exists := h.repo.Find(shortURL)

	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
