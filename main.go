package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func main() {
	r := mux.NewRouter()

	s := http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}


	s.ListenAndServe()

}
