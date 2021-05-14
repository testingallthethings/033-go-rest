package book_test

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/testingallthethings/033-go-rest/book"
	"github.com/testingallthethings/033-go-rest/rest"
	"testing"
)

type RetrieverSuite struct {
	suite.Suite
}

func TestRetrieverSuite(t *testing.T) {
	suite.Run(t, new(RetrieverSuite))
}

func (s *RetrieverSuite) TestRetrieverError() {
	expectedErr := errors.New("broken")
 	mbr := new(MockBookRetriever)
 	mbr.On("FindBookBy", "123456789").Return(rest.Book{}, expectedErr)

 	br := book.NewRetriever(mbr)
 	_, err := br.GetBook("123456789")

 	s.Error(err)
 	s.Equal(expectedErr, err)
}

func (s *RetrieverSuite) TestRetrieverSuccessful() {
	mbr := new(MockBookRetriever)
	expBook := rest.Book{ISBN: "123456789"}
	mbr.On("FindBookBy", "123456789").Return(expBook, nil)

	br := book.NewRetriever(mbr)
	book, err := br.GetBook("123456789")

	s.NoError(err)
	s.Equal(expBook, book)
}

func (s *RetrieverSuite) TestInvalidISBN() {
	mbr := new(MockBookRetriever)

	br := book.NewRetriever(mbr)
	_, err := br.GetBook("1234C6789")

	s.Error(err)
	s.Equal(rest.ErrInvalidISBN, err)
}

type MockBookRetriever struct {
	mock.Mock
}

func (m *MockBookRetriever) FindBookBy(isbn string) (rest.Book, error) {
	args := m.Called(isbn)

	return args.Get(0).(rest.Book), args.Error(1)
}
