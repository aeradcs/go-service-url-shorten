package handler

import (
	"errors"
	"github.com/gorilla/mux"
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
		return "http://localhost:8080/abcd123", nil
	}
	return "", errors.New("url not found")
}

func (m *mockUrlRepository) SaveUrl(original, shortKey string) (string, error) {
	if original == "http://example.com" {
		return "http://localhost:8080/abcd123", nil
	}
	return "", errors.New("error saving URL")
}

func (m *mockUrlRepository) GetOriginalUrl(short string) (string, error) {
	if short == "abcd123" {
		return "http://example.com", nil
	}
	return "", errors.New("url not found")
}

func TestPostReduceUrlSuccess(t *testing.T) {
	mockRepo := &mockUrlRepository{}
	useCase := &usecase.UrlShortenerUseCase{Repo: mockRepo}
	handler := &HttpHandler{UseCase: useCase}

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("http://example.com"))
	w := httptest.NewRecorder()
	handler.PostReduceUrl(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %v", resp.StatusCode)
	}
	if string(body) != "http://localhost:8080/abcd123" {
		t.Errorf("expected response body <http://localhost:8080/abcd123>, got %s", string(body))
	}
}

func TestPostReduceUrlMethodNotAllowed(t *testing.T) {
	mockRepo := &mockUrlRepository{}
	useCase := &usecase.UrlShortenerUseCase{Repo: mockRepo}
	handler := &HttpHandler{UseCase: useCase}

	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader("http://example.com"))
	w := httptest.NewRecorder()
	handler.PostReduceUrl(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %v", resp.StatusCode)
	}
	if string(body) != "Only POST requests are allowed!\n" {
		t.Errorf("expected response body <Only POST requests are allowed!\n>, got %s", string(body))
	}
}

func TestGerOriginalUrlSuccess(t *testing.T) {
	mockRepo := &mockUrlRepository{}
	useCase := &usecase.UrlShortenerUseCase{Repo: mockRepo}
	handler := &HttpHandler{UseCase: useCase}
	router := mux.NewRouter()
	router.HandleFunc("/{url_short}", handler.GerOriginalUrl).Methods(http.MethodGet)

	req := httptest.NewRequest(http.MethodGet, "/abcd123", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusTemporaryRedirect {
		t.Errorf("expected status 307, got %v", resp.StatusCode)
	}
	if string(body) != "http://example.com" {
		t.Errorf("expected response body <http://example.com>, got %s", string(body))
	}
}

func TestGerOriginalUrlNotFound(t *testing.T) {
	mockRepo := &mockUrlRepository{}
	useCase := &usecase.UrlShortenerUseCase{Repo: mockRepo}
	handler := &HttpHandler{UseCase: useCase}
	router := mux.NewRouter()
	router.HandleFunc("/{url_short}", handler.GerOriginalUrl).Methods(http.MethodGet)

	req := httptest.NewRequest(http.MethodGet, "/abcd123aaa", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %v", resp.StatusCode)
	}
	if string(body) != "url not found\n" {
		t.Errorf("expected response body <url not found\n>, got %s", string(body))
	}
}

func TestGerOriginalUrlMethodNotAllowed(t *testing.T) {
	mockRepo := &mockUrlRepository{}
	useCase := &usecase.UrlShortenerUseCase{Repo: mockRepo}
	handler := &HttpHandler{UseCase: useCase}
	router := mux.NewRouter()
	router.HandleFunc("/{url_short}", handler.GerOriginalUrl).Methods(http.MethodPost)

	req := httptest.NewRequest(http.MethodPost, "/abcd123aaa", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %v", resp.StatusCode)
	}
	if string(body) != "Only GET requests are allowed!\n" {
		t.Errorf("expected response body <Only GET requests are allowed!\n>, got %s", string(body))
	}
}
