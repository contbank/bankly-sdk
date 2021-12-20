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

const (
	Proxy          = "2229041000054459062"
	DocumentNumber = "21632071000187"
	Status         = "CanceledByCustomer"
)

type CardTestSuite struct {
	suite.Suite
	assert *assert.Assertions
	card   *bankly.Card
}

func mockCreateCard(documentNumber string, cardType bankly.CardType) bankly.CardCreateDTO {
	return bankly.CardCreateDTO{
		CardData: bankly.CardCreateRequest{
			DocumentNumber: documentNumber,
			CardName:       "NOME DA PESSOA",
			Alias:          "NOME PESSOA",
			BankAgency:     "0001",
			BankAccount:    "202142",
			Password:       "1234",
		},
		CardType: cardType,
	}
}

func (c *CardTestSuite) mockAlterCardCanceled(proxy string) {
	alterCardBody := bankly.CardUpdateStatusDTO{
		Status:           "CanceledByCustomer",
		Password:         "1234",
		UpdateCardBinded: true,
	}

	altered, err := c.card.UpdateStatusCard(context.Background(), proxy, alterCardBody)

	c.assert.NoError(err)
	c.assert.NotNil(altered)
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

	newHttpClient := bankly.NewHttpClient{
		Session:        *session,
		HttpClient:     httpClient,
		Authentication: bankly.NewAuthentication(httpClient, *session),
	}

	s.card = bankly.NewCard(newHttpClient)
}

func (c *CardTestSuite) TestGetCardsByIdentifier_OK() {
	card, err := c.card.GetCardsByIdentifier(context.Background(), DocumentNumber)
	c.assert.NoError(err)
	c.assert.NotNil(card)
}

func (c *CardTestSuite) TestGetCardsByIdentifier_NOT_FOUND() {
	card, err := c.card.GetCardsByIdentifier(context.Background(), "00000000000000")

	c.assert.Error(err)
	c.assert.Nil(card)
	c.assert.Equal("Code: 404 - Messages: not found", err.Error())
}

func (c *CardTestSuite) TestGetCardsByProxy_OK() {
	card, err := c.card.GetCardByProxy(context.Background(), Proxy)
	c.assert.NoError(err)
	c.assert.NotNil(card)
}

func (c *CardTestSuite) TestGetCardsByProxy_NOT_FOUND() {
	card, err := c.card.GetCardByProxy(context.Background(), "00000000000000")
	c.assert.Error(err)
	c.assert.Nil(card)
	c.assert.Equal("Code: 404 - Messages: not found", err.Error())
}

func (c *CardTestSuite) TestCreateCardVirtual_OK() {
	body := mockCreateCard(DocumentNumber, bankly.VirtualCardType)

	card, err := c.card.CreateCard(context.Background(), body)

	c.assert.NoError(err)
	c.assert.NotNil(card)
}

func (c *CardTestSuite) TestCreateCardVirtual_INVALID_PARAMETER_EMPTY() {
	body := mockCreateCard("1234567", bankly.VirtualCardType)

	card, err := c.card.CreateCard(context.Background(), body)

	c.assert.Error(err)
	c.assert.Nil(card)
}

func (c *CardTestSuite) TestCreateCardPhysical_OK() {
	body := mockCreateCard(DocumentNumber, bankly.PhysicalCardType)

	card, err := c.card.CreateCard(context.Background(), body)

	c.assert.NoError(err)
	c.assert.NotNil(card)
}

func (c *CardTestSuite) TestCreateCardPhysical_INVALID_PARAMETER_EMPTY() {
	body := mockCreateCard("123456", bankly.PhysicalCardType)

	card, err := c.card.CreateCard(context.Background(), body)

	c.assert.Error(err)
	c.assert.Nil(card)
}

func (c *CardTestSuite) TestAlteredStatusCard_OK() {
	card, err := c.card.GetCardsByIdentifier(context.Background(), DocumentNumber)

	c.assert.NoError(err)
	c.assert.NotNil(card)

	c.mockAlterCardCanceled(card[len(card)-1].Proxy)
	c.mockAlterCardCanceled(card[len(card)-2].Proxy)
}
