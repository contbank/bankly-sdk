package bankly_test

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/contbank/bankly-sdk"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BalanceTestSuite struct {
	suite.Suite
	assert  *assert.Assertions
	session *bankly.Session
	balance *bankly.Balance
}

func TestBalanceTestSuite(t *testing.T) {
	suite.Run(t, new(BalanceTestSuite))
}

func (s *BalanceTestSuite) SetupTest() {
	s.assert = assert.New(s.T())

	session, err := bankly.NewSession(bankly.Config{
		ClientID:     bankly.String(os.Getenv("BANKLY_CLIENT_ID")),
		ClientSecret: bankly.String(os.Getenv("BANKLY_CLIENT_SECRET")),
	})

	s.assert.NoError(err)

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	s.session = session
	s.balance = bankly.NewBalance(httpClient, *s.session)
}

func (s *BalanceTestSuite) TestBalance() {
	balance, err := s.balance.Balance("184152")

	s.assert.NoError(err)
	s.assert.NotNil(balance)
}
