package handler

import (
	"errors"
	"github.com/example-module/url-shortener/config"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/example-module/url-shortener/internal/usecase"
)

type mockUrlRepository struct{}

func (m *mockUrlRepository) GetShortUrl(original string) (string, error) {
	if original == "http://example.com" {
		return *config.BaseUrl + "abcd123", nil
	}
	return "", errors.New("url not found")
}

func (m *mockUrlRepository) SaveUrl(original, shortKey string) (string, error) {
	if original == "http://example.com" {
		return *config.BaseUrl + "abcd123", nil
	}
	return "", errors.New("error saving URL")
}

func (m *mockUrlRepository) GetOriginalUrl(short string) (string, error) {
	if short == "abcd123" {
		return "http://example.com", nil
	}
	return "", errors.New("url not found")
}

func TestShortenUrlAPI(t *testing.T) {
	tests := []struct {
		name         string
		handlerFunc  func(*HttpHandler) http.HandlerFunc
		method       string
		urlParam     string
		url          string
		body         string
		responseCode int
		responseBody string
	}{
		{
			name:         "Post Reduce Url Success",
			handlerFunc:  func(h *HttpHandler) http.HandlerFunc { return h.PostReduceUrl },
			method:       http.MethodPost,
			urlParam:     "/",
			url:          "/",
			body:         "http://example.com",
			responseCode: http.StatusOK,
			responseBody: *config.BaseUrl + "abcd123",
		},
		{
			name:         "Post Reduce Url Method Not Allowed",
			handlerFunc:  func(h *HttpHandler) http.HandlerFunc { return h.PostReduceUrl },
			method:       http.MethodGet,
			urlParam:     "/",
			url:          "/",
			body:         "http://example.com",
			responseCode: http.StatusMethodNotAllowed,
			responseBody: "Only POST requests are allowed!\n",
		},
		{
			name:         "Get Original Url Success",
			handlerFunc:  func(h *HttpHandler) http.HandlerFunc { return h.GerOriginalUrl },
			method:       http.MethodGet,
			urlParam:     "/{url_short}",
			url:          "/abcd123",
			body:         "",
			responseCode: http.StatusTemporaryRedirect,
			responseBody: "http://example.com",
		},
		{
			name:         "Get Original Url Not Found",
			handlerFunc:  func(h *HttpHandler) http.HandlerFunc { return h.GerOriginalUrl },
			method:       http.MethodGet,
			urlParam:     "/{url_short}",
			url:          "/abcd123aaa",
			body:         "",
			responseCode: http.StatusNotFound,
			responseBody: "url not found\n",
		},
		{
			name:         "Get Original Url Method Not Allowed",
			handlerFunc:  func(h *HttpHandler) http.HandlerFunc { return h.GerOriginalUrl },
			method:       http.MethodPost,
			urlParam:     "/{url_short}",
			url:          "/abcd123",
			body:         "",
			responseCode: http.StatusMethodNotAllowed,
			responseBody: "Only GET requests are allowed!\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := &mockUrlRepository{}
			useCase := &usecase.UrlShortenerUseCase{Repo: mockRepo}
			handler := &HttpHandler{UseCase: useCase}
			router := mux.NewRouter()
			router.HandleFunc(test.urlParam, test.handlerFunc(handler)).Methods(test.method)

			req := httptest.NewRequest(test.method, test.url, strings.NewReader(test.body))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			resp := w.Result()
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("failed to read response body: %v", err)
			}
			assert.Equal(t, resp.StatusCode, test.responseCode)
			assert.Equal(t, string(body), test.responseBody)
		})
	}
}
