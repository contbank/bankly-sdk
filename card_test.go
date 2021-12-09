package bankly_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	bankly "github.com/contbank/bankly-sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CardTestSuite struct {
	suite.Suite
	assert  *assert.Assertions
	session *bankly.Session
	card    *bankly.Card
}

func TestCardTestSuite(t *testing.T) {
	suite.Run(t, new(CardTestSuite))
}

func (s *CardTestSuite) SetupTest() {
	s.assert = assert.New(s.T())

	session, err := bankly.NewSession(bankly.Config{
		ClientID:     bankly.String(*bankly.GetEnvBanklyClientID()),
		ClientSecret: bankly.String(*bankly.GetEnvBanklyClientSecret()),
	})

	s.assert.NoError(err)

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	s.session = session
	s.card = bankly.NewCard(httpClient, *s.session)
}

func (c *CardTestSuite) TestGetCardsByIdentifier_OK() {
	card, err := c.card.GetCardsByIdentifier(context.Background(), "21632071000187")
	c.assert.NoError(err)
	c.assert.NotNil(card)
}

func (c *CardTestSuite) TestGetCardsByIdentifier_NOT_FOUND() {
	card, err := c.card.GetCardsByIdentifier(context.Background(), "00000000000000")

	c.assert.Error(err)
	c.assert.Nil(card)
	c.assert.Equal("Code: 404 - Messages: not found", err.Error())
}
