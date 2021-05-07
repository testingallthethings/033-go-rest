package rest_test

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/testingallthethings/033-go-rest/rest"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type GetBookSuite struct {
	suite.Suite
}

func TestGetBookSuite(t *testing.T) {
	suite.Run(t, new(GetBookSuite))
}

var (
	req *http.Request
	resp *httptest.ResponseRecorder
	br *MockBookRetriever
	h rest.GetBookHandler
)

func (s *GetBookSuite) SetupTest() {
	req, _ = http.NewRequest(http.MethodGet, "/book/123456789", nil)
	req = mux.SetURLVars(req, map[string]string{"isbn": "123456789"})

	resp = httptest.NewRecorder()

	br = new(MockBookRetriever)

	h = rest.NewGetBookHandler(br)
}

func (s *GetBookSuite) TestGettingBookThatDoesNotExist() {
 	br.On("GetBook", "123456789").Return(rest.Book{}, rest.ErrBookNotFound)

 	h.ServeHTTP(resp, req)

 	body, _ := ioutil.ReadAll(resp.Body)

 	s.Equal(http.StatusNotFound, resp.Code)
	s.JSONEq(`{"code": "001", "msg": "No book with ISBN 123456789"}`, string(body))
}

func (s *GetBookSuite) TestGettingABookThatDoesExist() {
	book := rest.Book{
		ISBN:          "123456789",
		Title:         "Testing All The Things",
		Image:         "testing.jpg",
		Genre:         "Computing",
		YearPublished: 2021,
	}
	br.On("GetBook", "123456789").Return(book, nil)

	h.ServeHTTP(resp, req)

	s.Equal(http.StatusOK, resp.Code)
	body, _ := ioutil.ReadAll(resp.Body)
	expBody := `{
	"isbn": "123456789",
	"title": "Testing All The Things",
	"image": "testing.jpg",
	"genre": "Computing",
	"year_published": 2021
}`
	s.JSONEq(expBody, string(body))
}

func (s *GetBookSuite) TestGetBookReturnsUnexpectedError() {
	br.On("GetBook", "123456789").Return(rest.Book{}, errors.New("broken"))

	h.ServeHTTP(resp, req)

	body, _ := ioutil.ReadAll(resp.Body)
	s.Equal(http.StatusInternalServerError, resp.Code)
	s.JSONEq(`{"code": "002", "msg": "error attempting to get book"}`, string(body))
}

type MockBookRetriever struct {
	mock.Mock
}

func (m *MockBookRetriever) GetBook(isbn string) (rest.Book, error) {
	args := m.Called(isbn)

	return args.Get(0).(rest.Book), args.Error(1)
}

