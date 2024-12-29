package main

import (
	"net/http"
	"io"
	"database/sql"
	_ "github.com/lib/pq"
	"math/rand"
    "time"
	"strings"
)

func generateShortKey() string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    const keyLength = 10

    rand.Seed(time.Now().UnixNano())
    shortKey := make([]byte, keyLength)
    for i := range shortKey {
        shortKey[i] = charset[rand.Intn(len(charset))]
    }
    return string(shortKey)
}

func postReduceUrl(w http.ResponseWriter, req *http.Request) {
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
	urlShort := getShortUrlFromDb(urlOriginal)
	if urlShort == "" {
		shortKey := generateShortKey()
		urlShort =	insertUrlInDb(urlOriginal, shortKey)
	}
	resp := generateResponse(urlOriginal, urlShort)
	w.Write(resp)
}

func getShortUrlFromDb(original string) string {
	connStr := "user=nemo dbname=urlreducedb password=1101 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	query := "SELECT url_original, url_shortened FROM urls WHERE url_original = '" + original + "'"
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var url_original, url_shortened string
		if err := rows.Scan(&url_original, &url_shortened); err != nil {
			panic(err)
		}
		return url_shortened
	}
	return ""
}

func insertUrlInDb(original string, shortKey string) string {
	var builder strings.Builder
	builder.WriteString("INSERT INTO urls VALUES ('")
	builder.WriteString(string(original))
	builder.WriteString("', 'http://localhost:8080/")
	builder.WriteString(shortKey)
	builder.WriteString("')")
	insertQueryStr := builder.String()

	connStr := "user=nemo dbname=urlreducedb password=1101 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query(insertQueryStr)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	return "http://localhost:8080/" + shortKey
}

func generateResponse(original string, short string) []byte {
	var builder strings.Builder
	builder.WriteString(original)
	builder.WriteString("\n")
	builder.WriteString(short)
	return []byte(builder.String())
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, postReduceUrl)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
