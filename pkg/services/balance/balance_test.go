package bankly_test

import (
	"context"
	"github.com/contbank/bankly-sdk/pkg/services/authentication"
	balance "github.com/contbank/bankly-sdk/pkg/services/balance"
	utils "github.com/contbank/bankly-sdk/pkg/utils"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BalanceTestSuite struct {
	suite.Suite
	ctx     context.Context
	assert  *assert.Assertions
	session *bankly.Session
	balance *balance.Balance
}

func TestBalanceTestSuite(t *testing.T) {
	suite.Run(t, new(BalanceTestSuite))
}

func (s *BalanceTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.ctx = context.Background()

	session, err := bankly.NewSession(bankly.Config{
		ClientID:     bankly.String(*utils.GetEnvBanklyClientID()),
		ClientSecret: bankly.String(*utils.GetEnvBanklyClientSecret()),
		Scopes:       bankly.String("account.read"),
	})

	s.assert.NoError(err)

	httpClient := &http.Client{
		Timeout:   30 * time.Second,
		Transport: utils.LoggingRoundTripper{Proxied: http.DefaultTransport},
	}

	s.session = session
	s.balance = balance.NewBalance(httpClient, *s.session)
}

func (s *BalanceTestSuite) TestBalance() {
	balance, err := s.balance.Balance(s.ctx, "184152")

	s.assert.NoError(err)
	s.assert.NotNil(balance)
}