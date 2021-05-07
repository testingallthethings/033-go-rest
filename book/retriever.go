package book

import (
	"database/sql"
	"errors"
	"github.com/testingallthethings/033-go-rest/rest"
)

type Retriever struct {
	db *sql.DB
}

var (
	ErrFailedToRetrieveBook = errors.New("error occurred retrieving book")
)

func (r Retriever) GetBook(isbn string) (rest.Book, error) {
	b := rest.Book{}

	row := r.db.QueryRow("SELECT isbn, name, image, genre, year_published FROM book WHERE isbn = $1", isbn)
	err := row.Scan(
		&b.ISBN,
		&b.Title,
		&b.Image,
		&b.Genre,
		&b.YearPublished,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return b, rest.ErrBookNotFound
		}

		return b, ErrFailedToRetrieveBook
	}

	return b, nil
}

func NewRetriever(db *sql.DB) Retriever {
	return Retriever{db}
}