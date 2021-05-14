package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/testingallthethings/033-go-rest/book"
	"github.com/testingallthethings/033-go-rest/rest"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"

	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
)

type health struct {
	Status string `json:"status"`
	Messages []string `json:"messages"`
}

type jsonError struct {
	Code string `json:"code"`
	Msg string `json:"msg"`
}

type Book struct {
	ISBN string `json:"isbn"`
	Title string `json:"title"`
	Image string `json:"image"`
	Genre string `json:"genre"`
	YearPublished int `json:"year_published"`
}

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error making DB connected: %s", err.Error())
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Error making DB driver: %s", err.Error())
	}

	migrator, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("Error making migration engine: %s", err.Error())
	}
	migrator.Steps(2)

	r := mux.NewRouter()

	dbRetriever := book.NewDBRetriever(db)
	retriever := book.NewRetriever(dbRetriever)

	r.Handle("/book/{isbn}", rest.NewGetBookHandler(retriever))
	r.HandleFunc(
		"/healthcheck",
		func(w http.ResponseWriter, r *http.Request) {
			h := health{
				Status:   "OK",
				Messages: []string{},
			}

			b, _ := json.Marshal(h)

			w.WriteHeader(http.StatusOK)
			w.Write(b)
	})

	s := http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}
	s.ListenAndServe()
}
