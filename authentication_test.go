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

type AuthenticationTestSuite struct {
	suite.Suite
	assert         *assert.Assertions
	ctx            context.Context
	session        *bankly.Session
	authentication *bankly.Authentication
}

func TestAuthenticationTestSuite(t *testing.T) {
	suite.Run(t, new(AuthenticationTestSuite))
}

func (s *AuthenticationTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.ctx = context.Background()

	session, err := bankly.NewSession(bankly.Config{
		ClientID : bankly.String(*bankly.GetEnvBanklyClientID()),
		ClientSecret : bankly.String(*bankly.GetEnvBanklyClientSecret()),
	})

	s.assert.NoError(err)

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	s.session = session
	s.authentication = bankly.NewAuthentication(httpClient, *s.session)
}

func (s *AuthenticationTestSuite) TestToken() {
	token, err := s.authentication.Token(s.ctx)

	s.assert.NoError(err)
	s.assert.Contains(token, "Bearer")
}
