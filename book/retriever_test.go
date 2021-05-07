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

type RetrieverSuite struct {
	suite.Suite
}

func TestRetrieverSuite(t *testing.T) {
	suite.Run(t, new(RetrieverSuite))
}

var (
	db *sql.DB
	r book.Retriever
)

func (s *RetrieverSuite) SetupTest() {
	db, _ = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	r = book.NewRetriever(db)
}

func (s *RetrieverSuite) TearDownTest() {
	db.Close()
}

func (s *RetrieverSuite) TestRetrievingBookThatDoesNotExist() {
	_, err := r.GetBook("123456789")

	s.Equal(rest.ErrBookNotFound, err)
}

func (s *RetrieverSuite) TestRetrievingBookThatExists() {
	db.Exec("INSERT INTO book (isbn, name, image, genre, year_published) VALUES ('987654321', 'Testing All The Things', 'testing.jpg', 'Computing', 2021)")

	b, err := r.GetBook("987654321")

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

func (s *RetrieverSuite) TestWhenUnexpectedErrorRetrievingBook() {
	db.Close()

	_, err := r.GetBook("123456789")

	s.Equal(book.ErrFailedToRetrieveBook, err)
}
