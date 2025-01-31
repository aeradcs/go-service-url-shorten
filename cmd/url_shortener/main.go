package main

import (
	"flag"
	"fmt"
	"github.com/example-module/url-shortener/config"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"

	"github.com/example-module/url-shortener/internal/handler"
	"github.com/example-module/url-shortener/internal/repository"
	"github.com/example-module/url-shortener/internal/usecase"
	_ "github.com/lib/pq"
)

func main() {
	flag.Parse()
	fmt.Printf("Parsed args : a = %s, b = %s\n", *config.Port, *config.BaseUrl)

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

	router := mux.NewRouter()
	router.HandleFunc(`/{url_short}`, handler.GerOriginalUrl)
	router.HandleFunc(`/`, handler.PostReduceUrl)

	fmt.Printf("Server is running on :%s\n", *config.Port)
	if err := http.ListenAndServe(":"+*config.Port, router); err != nil {
		panic(err)
	}
}
