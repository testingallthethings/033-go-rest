package e2e_test

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type EndToEndSuite struct {
	suite.Suite
}

func TestEndToEndSuite(t *testing.T) {
	suite.Run(t, new(EndToEndSuite))
}

func (s *EndToEndSuite) TestPlaceholder() {
	s.True(true)
}
