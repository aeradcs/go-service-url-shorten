package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http" // Standard library
	"service-url-shortener/internal/adapters/db"
	adapterHTTP "service-url-shortener/internal/adapters/http" // Aliased to avoid conflict
	"service-url-shortener/internal/usecases"

	_ "github.com/lib/pq"
)

func main() {
	// Connect to the database
	connStr := "user=nemo dbname=urlreducedb password=1101 sslmode=disable"
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer dbConn.Close()

	// Repositories and use case initialization
	repo := db.NewUrlRepository(dbConn)
	useCase := usecases.NewUrlShortener(repo)

	// HTTP Handler
	handler := adapterHTTP.NewUrlHandler(useCase)

	// HTTP Routes
	http.HandleFunc("/", handler.PostReduceUrl)

	// Start the server
	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
