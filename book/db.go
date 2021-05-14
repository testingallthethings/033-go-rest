package book

import (
	"database/sql"
	"errors"
	"github.com/testingallthethings/033-go-rest/rest"
)

type DBRetriever struct {
	db *sql.DB
}

var (
	ErrFailedToRetrieveBook = errors.New("error occurred retrieving book")
)

func (r DBRetriever) FindBookBy(isbn string) (rest.Book, error) {
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

func NewDBRetriever(db *sql.DB) DBRetriever {
	return DBRetriever{db}
}