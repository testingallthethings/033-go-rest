// +build pact

package e2e_test

import (
	"database/sql"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type PactSuite struct {
	suite.Suite
}

func TestPactSuite(t *testing.T) {
	suite.Run(t, new(PactSuite))
}

func (s *PactSuite) SetupTest() {
	db, _ = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	db.Exec("TRUNCATE book")
}

func (s *PactSuite) TestPact() {
	pact := dsl.Pact{
		Provider: "BookApi",
	}

	_, err := pact.VerifyProvider(
		s.T(),
		types.VerifyRequest{
			ProviderBaseURL: "http://localhost:8080",
			PactURLs: []string{"https://braddle.pact.dius.com.au/pacts/provider/BookApi/consumer/MarksBookClient/latest"},
			BrokerToken: os.Getenv("PACT_TOKEN"),
			PublishVerificationResults: true,
			ProviderVersion: "1.0.0",
			StateHandlers: types.StateHandlers{
				"Book with ISBN 987654321 exists": func() error {
					_, err := db.Exec("INSERT INTO book (isbn, name, image, genre, year_published) VALUES ('987654321', 'Testing All The Things', 'testing.jpg', 'Computing', 2021)")
					return err
				},
			},
		},
	)

	s.NoError(err)

}
