package bankly_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/contbank/bankly-sdk"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BalanceTestSuite struct {
	suite.Suite
	ctx     context.Context
	assert  *assert.Assertions
	session *bankly.Session
	balance *bankly.Balance
}

func TestBalanceTestSuite(t *testing.T) {
	suite.Run(t, new(BalanceTestSuite))
}

func (s *BalanceTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.ctx = context.Background()

	session, err := bankly.NewSession(bankly.Config{
		ClientID : bankly.String(*bankly.GetEnvBanklyClientID()),
		ClientSecret : bankly.String(*bankly.GetEnvBanklyClientSecret()),
	})

	s.assert.NoError(err)

	httpClient := &http.Client{
		Timeout:   30 * time.Second,
		Transport: bankly.LoggingRoundTripper{Proxied: http.DefaultTransport},
	}

	s.session = session
	s.balance = bankly.NewBalance(httpClient, *s.session)
}

func (s *BalanceTestSuite) TestBalance() {
	balance, err := s.balance.Balance(s.ctx, "184152")

	s.assert.NoError(err)
	s.assert.NotNil(balance)
}
