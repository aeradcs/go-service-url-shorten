package http

import (
	"encoding/json"
	"net/http"
	"service-url-shortener/internal/usecases"
)

type UrlHandler struct {
	useCase usecases.UrlShortener
}

func NewUrlHandler(useCase usecases.UrlShortener) *UrlHandler {
	return &UrlHandler{useCase: useCase}
}

func (h *UrlHandler) PostReduceUrl(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		Original string `json:"original"`
	}
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	shortUrl, err := h.useCase.GetShortUrl(input.Original)
	if err != nil {
		shortUrl, err = h.useCase.GenerateShortUrl(input.Original)
		if err != nil {
			http.Error(w, "Error generating short URL", http.StatusInternalServerError)
			return
		}
	}

	response := struct {
		Original  string `json:"original"`
		Shortened string `json:"shortened"`
	}{
		Original:  input.Original,
		Shortened: shortUrl,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
