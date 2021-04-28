package bankly_test

import (
	"os"
	"testing"

	"github.com/contbank/bankly-sdk"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthenticationTestSuite struct {
	suite.Suite
	assert         *assert.Assertions
	session        *bankly.Session
	authentication *bankly.Authentication
}

func TestAuthenticationTestSuite(t *testing.T) {
	suite.Run(t, new(AuthenticationTestSuite))
}

func (s *AuthenticationTestSuite) SetupTest() {
	s.assert = assert.New(s.T())

	session, err := bankly.NewSession(bankly.Config{
		ClientID:     bankly.String(os.Getenv("BANKLY_CLIENT_ID")),
		ClientSecret: bankly.String(os.Getenv("BANKLY_CLIENT_SECRET")),
	})

	s.assert.NoError(err)

	s.session = session
	s.authentication = bankly.NewAuthentication(*s.session)
}

func (s *AuthenticationTestSuite) TestToken() {
	token, err := s.authentication.Token()

	s.assert.NoError(err)
	s.assert.Contains(token, "Bearer")
}
