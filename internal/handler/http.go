package handler

import (
	"fmt"
	"github.com/example-module/url-shortener/internal/usecase"
	"io"
	"net/http"
)

type HttpHandler struct {
	UseCase *usecase.UrlShortenerUseCase
}

func (h *HttpHandler) PostReduceUrl(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusInternalServerError)
		return
	}

	urlOriginal := string(body)
	urlShort, err := h.UseCase.ShortenUrl(urlOriginal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf("%s\n%s", urlOriginal, urlShort)
	w.Write([]byte(response))
}
