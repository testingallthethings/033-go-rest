// +build integration

package book_test

import (
	"database/sql"
	"github.com/stretchr/testify/suite"
	"github.com/testingallthethings/033-go-rest/book"
	"github.com/testingallthethings/033-go-rest/rest"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

type DBRetrieverSuite struct {
	suite.Suite
}

func TestDBRetrieverSuite(t *testing.T) {
	suite.Run(t, new(DBRetrieverSuite))
}

var (
	db *sql.DB
	r book.DBRetriever
)

func (s *DBRetrieverSuite) SetupTest() {
	db, _ = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	r = book.NewDBRetriever(db)
}

func (s *DBRetrieverSuite) TearDownTest() {
	db.Close()
}

func (s *DBRetrieverSuite) TestRetrievingBookThatDoesNotExist() {
	_, err := r.FindBookBy("123456789")

	s.Equal(rest.ErrBookNotFound, err)
}

func (s *DBRetrieverSuite) TestRetrievingBookThatExists() {
	db.Exec("INSERT INTO book (isbn, name, image, genre, year_published) VALUES ('987654321', 'Testing All The Things', 'testing.jpg', 'Computing', 2021)")

	b, err := r.FindBookBy("987654321")

	s.NoError(err)

	book := rest.Book{
		ISBN:          "987654321",
		Title:         "Testing All The Things",
		Image:         "testing.jpg",
		Genre:         "Computing",
		YearPublished: 2021,
	}

	s.Equal(book, b)
}

func (s *DBRetrieverSuite) TestWhenUnexpectedErrorRetrievingBook() {
	db.Close()

	_, err := r.FindBookBy("123456789")

	s.Equal(book.ErrFailedToRetrieveBook, err)
}
