package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type Book struct {
	ISBN          string `json:"isbn"`
	Title         string `json:"title"`
	Image         string `json:"image"`
	Genre         string `json:"genre"`
	YearPublished int    `json:"year_published"`
}

type jsonError struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

var (
	ErrBookNotFound = errors.New("Book not found")
)

type BookRetriever interface {
	GetBook(isbn string) (Book, error)
}

type GetBookHandler struct {
	br BookRetriever
}

func (g GetBookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	isbn := v["isbn"]

	book, err := g.br.GetBook(isbn)

	if err != nil {
		var e jsonError
		if err == ErrBookNotFound {
			e.Code = "001"
			e.Msg = fmt.Sprintf("No book with ISBN %s", isbn)
			w.WriteHeader(http.StatusNotFound)
		} else {
			e.Code = "002"
			e.Msg = "error attempting to get book"
			w.WriteHeader(http.StatusInternalServerError)
		}

		body, _ := json.Marshal(e)
		w.Write(body)
		return
	}

	body, _ := json.Marshal(book)

	w.WriteHeader(http.StatusOK)
	w.Write(body)

}

func NewGetBookHandler(br BookRetriever) GetBookHandler {
	return GetBookHandler{br}
}
