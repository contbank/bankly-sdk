package bankly_test

import (
	"context"
	"github.com/contbank/bankly-sdk/pkg/services/authentication"
	utils "github.com/contbank/bankly-sdk/pkg/utils"
	"net/http"
	"testing"
	"time"

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
		ClientID:     utils.String(*utils.GetEnvBanklyClientID()),
		ClientSecret: utils.String(*utils.GetEnvBanklyClientSecret()),
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