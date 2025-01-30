package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"

	"github.com/example-module/url-shortener/internal/handler"
	"github.com/example-module/url-shortener/internal/repository"
	"github.com/example-module/url-shortener/internal/usecase"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	repo, err := repository.NewPostgresRepository()
	if err != nil {
		panic(err)
	}
	defer repo.Db.Close()

	useCase := &usecase.UrlShortenerUseCase{Repo: repo}
	handler := &handler.HttpHandler{UseCase: useCase}

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, handler.PostReduceUrl)

	fmt.Println("Server is running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
