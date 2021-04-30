package e2e_test

import (
	"database/sql"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

type GetSingleBookSuite struct {
	suite.Suite
}

func TestGetSingleBookSuite(t *testing.T) {
	suite.Run(t, new(GetSingleBookSuite))
}

func (s *GetSingleBookSuite) TestGetBookThatDoesNotExist() {
	c := http.Client{}

	r, _ := c.Get("http://localhost:8080/book/123456789")
	body, _ := ioutil.ReadAll(r.Body)

	s.Equal(http.StatusNotFound, r.StatusCode)
 	s.JSONEq(`{"code": "001", "msg": "No book with ISBN 123456789"}`, string(body))
}

func (s *GetSingleBookSuite) TestGetBookThatDoesExist() {
	db, _ := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	db.Exec("INSERT INTO book (isbn, name, image, genre, year_published) VALUES ('987654321', 'Testing All The Things', 'testing.jpg', 'Computing', 2021)")

	c := http.Client{}

	r, _ := c.Get("http://localhost:8080/book/987654321")
	body, _ := ioutil.ReadAll(r.Body)

	s.Equal(http.StatusOK, r.StatusCode)

	expBody := `{
	"isbn": "987654321",
	"title": "Testing All The Things",
	"image": "testing.jpg",
	"genre": "Computing",
	"year_published": 2021
}`

	s.JSONEq(expBody, string(body))
}
