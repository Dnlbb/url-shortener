package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFget(t *testing.T) {
	type Want struct {
		statusCode int
		location   string
	}

	type Request struct {
		path string
	}

	testCases := []struct {
		name        string
		want        Want
		request     Request
		originalURL string
	}{
		{
			name: "#1 Valid short URL",
			want: Want{
				statusCode: http.StatusTemporaryRedirect,
				location:   "https://practicum.yandex.ru/",
			},
			request: Request{
				path: "/" + GenerateShortURL("https://practicum.yandex.ru/"),
			},
			originalURL: "https://practicum.yandex.ru/",
		},
		{
			name: "#2 Invalid short URL",
			want: Want{
				statusCode: http.StatusBadRequest,
				location:   "",
			},
			request: Request{
				path: "/invalidURL",
			},
			originalURL: "",
		},
		{
			name: "#3 Valid test",
			want: Want{
				statusCode: http.StatusTemporaryRedirect,
				location:   "https://aaaa.ru/",
			},
			request: Request{
				path: "/" + GenerateShortURL("https://aaaa.ru/"),
			},
			originalURL: "https://aaaa.ru/",
		},
		{
			name: "#4 Valid test",
			want: Want{
				statusCode: http.StatusTemporaryRedirect,
				location:   "https://asfasfsafsaf.ru/",
			},
			request: Request{
				path: "/" + GenerateShortURL("https://asfasfsafsaf.ru/"),
			},
			originalURL: "https://asfasfsafsaf.ru/",
		},
		{
			name: "#5 Valid test",
			want: Want{
				statusCode: http.StatusTemporaryRedirect,
				location:   "https://tegeregergafsaf.ru/",
			},
			request: Request{
				path: "/" + GenerateShortURL("https://tegeregergafsaf.ru/"),
			},
			originalURL: "https://tegeregergafsaf.ru/",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := NewMockRepository()
			handler := NewHandler(mockRepo)
			shortURL := GenerateShortURL(test.originalURL)

			if test.originalURL != "" {
				mockRepo.Save(shortURL, test.originalURL)
			}

			request := httptest.NewRequest(http.MethodGet, test.request.path, nil)
			w := httptest.NewRecorder()
			handler.Fget(w, request)

			result := w.Result()

			assert.Equal(t, test.want.statusCode, result.StatusCode)

			location := result.Header.Get("Location")
			assert.Equal(t, test.want.location, location)
			err := result.Body.Close()
			require.NoError(t, err)
		})
	}
}
