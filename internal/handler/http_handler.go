package handler

import (
	"fmt"
	"github.com/example-module/url-shortener/internal/usecase"
	"github.com/gorilla/mux"
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

	response := fmt.Sprint(urlShort)
	w.Write([]byte(response))
}

func (h *HttpHandler) GerOriginalUrl(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	urlShort := mux.Vars(request)["url_short"]
	urlOriginal, err := h.UseCase.GerOriginalUrl(urlShort)
	if err != nil {
		if err.Error() == "url not found" {
			http.Error(writer, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusTemporaryRedirect)
	writer.Write([]byte(urlOriginal))
}
