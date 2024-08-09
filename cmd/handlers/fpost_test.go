package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFpost(t *testing.T) {
	type Want struct {
		contentType string
		statusCode  int
	}

	type Request struct {
		path string
		body string
	}

	testCases := []struct {
		name    string
		want    Want
		request Request
	}{
		{
			name: "#1 Valid URL",
			want: Want{
				contentType: "text/plain",
				statusCode:  201},
			request: Request{
				path: "/",
				body: "https://practicum.yandex.ru/"},
		},
		{
			name: "#2 Valid URL",
			want: Want{
				contentType: "text/plain",
				statusCode:  201},
			request: Request{
				path: "/",
				body: "https://sdfsdfsdxcxcv.yandex.ru/"},
		},
		{
			name: "#3 Valid URL",
			want: Want{
				contentType: "text/plain",
				statusCode:  201},
			request: Request{
				path: "/",
				body: "https://aaa.yandex.ru/"},
		},
		{
			name: "#4 Invalid URL in Body",
			want: Want{
				contentType: "",
				statusCode:  http.StatusBadRequest,
			},
			request: Request{
				path: "/",
				body: "invalid-url",
			},
		},
		{
			name: "#5 Empty Body",
			want: Want{
				contentType: "",
				statusCode:  http.StatusBadRequest,
			},
			request: Request{
				path: "/",
				body: "",
			},
		},
		{
			name: "#6 Invalid URL in body",
			want: Want{
				contentType: "",
				statusCode:  http.StatusBadRequest},
			request: Request{
				path: "/",
				body: "https:///"},
		},
		{
			name: "#7 Valid URL",
			want: Want{
				contentType: "text/plain",
				statusCode:  201},
			request: Request{
				path: "/",
				body: "https://.ru/"},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := NewMockRepository()
			handler := NewHandler(mockRepo)

			request := httptest.NewRequest(http.MethodPost, test.request.path, strings.NewReader(test.request.body))
			w := httptest.NewRecorder()
			handler.Fpost(w, request)

			result := w.Result()

			assert.Equal(t, test.want.statusCode, result.StatusCode)

			if test.want.contentType != "" {
				assert.Equal(t, test.want.contentType, result.Header.Get("Content-Type"))
			}

			if test.want.statusCode == http.StatusCreated {
				shortURL := strings.TrimPrefix(w.Body.String(), "http://localhost:8080/")
				_, exists := mockRepo.Find(shortURL)
				if !exists {
					t.Errorf("short URL was not saved in the repository")
				}
			}
		})
	}

}
