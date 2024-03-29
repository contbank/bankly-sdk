package bankly_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/contbank/grok"

	bankly "github.com/contbank/bankly-sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CardTestSuite struct {
	suite.Suite
	assert *assert.Assertions
	ctx    context.Context
	card   *bankly.Card
}

func (c *CardTestSuite) mockAlterCardCanceled(proxy string) {
	alterCardBody := bankly.CardUpdateStatusDTO{
		Status:           "CanceledByCustomer",
		Password:         "1234",
		UpdateCardBinded: true,
	}

	err := c.card.UpdateStatusCardByProxy(context.Background(), &proxy, &alterCardBody)

	c.assert.NoError(err)
}

func TestCardTestSuite(t *testing.T) {
	suite.Run(t, new(CardTestSuite))
}

func (s *CardTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.ctx = context.Background()

	session, err := bankly.NewSession(bankly.Config{
		ClientID:     bankly.String(*bankly.GetEnvBanklyClientID()),
		ClientSecret: bankly.String(*bankly.GetEnvBanklyClientSecret()),
		Scopes:       bankly.String("card.create card.update card.read card.pci.password.update"),
	})
	s.assert.NoError(err)

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	newHttpClient := bankly.NewBanklyHttpClient(*session, httpClient, bankly.NewAuthentication(httpClient, *session))

	s.card = bankly.NewCard(newHttpClient)
}

func (c *CardTestSuite) TestGetCardsByIdentifier_OK() {
	// TODO Mockar teste
	c.T().Skip("O bankly está retornando 500 com o body null. Mockar teste.")

	card, err := c.card.GetCardsByIdentifier(c.ctx, "93707422046")
	c.assert.NoError(err)
	c.assert.NotNil(card)
}

func (c *CardTestSuite) TestGetCardsByIdentifier_NOT_FOUND() {
	// TODO Mockar teste
	c.T().Skip("O bankly está retornando 500 com o body null. Mockar teste.")

	card, err := c.card.GetCardsByIdentifier(c.ctx, "00000000000000")

	c.assert.Error(err)
	c.assert.Nil(card)
	c.assert.Equal("Code: 404 - Messages: not found", err.Error())
}

func (c *CardTestSuite) TestGetTransactionByProxy_OK() {
	// TODO Mockar teste
	c.T().Skip("O bankly está retornando 500 com o body null. Mockar teste.")

	proxy := "2229041000054459062"
	card, err := c.card.GetTransactionsByProxy(c.ctx, &proxy, "1", "2021-01-01", "2021-01-08", "10")
	c.assert.NoError(err)
	c.assert.NotNil(card)
}

func (c *CardTestSuite) TestGetTransactionByProxy_INTERVAL_DATE_NOT_OK() {
	// TODO Mockar teste
	c.T().Skip("O bankly está retornando 500 com o body null. Mockar teste.")

	proxy := "2229041000054459062"
	card, err := c.card.GetTransactionsByProxy(c.ctx, &proxy, "1", "2021-01-01", "2021-01-09", "10")
	c.assert.Error(err)
	c.assert.Nil(card)
}

func (c *CardTestSuite) TestGetTransactionByProxy_ENDDATE_NOTFOUND_NOT_OK() {
	// TODO Mockar teste
	c.T().Skip("O bankly está retornando 500 com o body null. Mockar teste.")

	proxy := "2229041000054459062"
	card, err := c.card.GetTransactionsByProxy(c.ctx, &proxy, "1", "2021-01-01", "", "10")
	c.assert.Error(err)
	c.assert.Nil(card)
}

func (c *CardTestSuite) TestGetCardByProxy_OK() {
	// TODO Mockar teste
	c.T().Skip("O bankly está retornando 500 com o body null. Mockar teste.")

	card, err := c.card.GetCardByProxy(c.ctx, "2229041000054459062")
	c.assert.NoError(err)
	c.assert.NotNil(card)
	c.assert.NotEmpty(card.Proxy)
	c.assert.NotEmpty(card.Status)
	c.assert.NotEmpty(card.Address)
	c.assert.NotNil(card.IsActivated)
	c.assert.NotNil(card.IsLocked)
	c.assert.NotNil(card.IsCanceled)
	c.assert.NotNil(card.IsPos)
	c.assert.NotNil(card.IsBuilding)
}

func (c *CardTestSuite) TestGetCardByProxy_NOT_FOUND() {
	// TODO Mockar teste
	c.T().Skip("O bankly está retornando 500 com o body null. Mockar teste.")

	card, err := c.card.GetCardByProxy(c.ctx, "00000000000000")
	c.assert.Error(err)
	c.assert.Nil(card)
	c.assert.Equal("Code: 404 - Messages: not found", err.Error())
}

func (c *CardTestSuite) TestGetCardByActivateCode_OK() {
	// TODO Mockar teste
	c.T().Skip("O bankly está retornando 500 com o body null. Mockar teste.")

	card, err := c.card.GetCardByActivateCode(c.ctx, "F6D8B3C1269D")
	c.assert.NoError(err)
	c.assert.NotNil(card)
	c.assert.GreaterOrEqual(1, len(card))
	c.assert.NotEmpty(card[0].Proxy)
	c.assert.NotEmpty(card[0].Status)
	c.assert.NotEmpty(card[0].Address)
	c.assert.NotNil(card[0].IsActivated)
	c.assert.NotNil(card[0].IsLocked)
	c.assert.NotNil(card[0].IsCanceled)
	c.assert.NotNil(card[0].IsPos)
	c.assert.NotNil(card[0].IsBuilding)
}

func (c *CardTestSuite) TestGetCardByActivateCode_NOT_FOUND() {
	// TODO Mockar teste
	c.T().Skip("O bankly está retornando 500 com o body null. Mockar teste.")

	card, err := c.card.GetCardByActivateCode(c.ctx, "000000000000")
	c.assert.Error(err)
	c.assert.Nil(card)
	c.assert.Equal("Code: 404 - Messages: not found", err.Error())
}

func (c *CardTestSuite) TestGetNextStatusByProxy_OK() {
	// TODO Mockar teste
	c.T().Skip("O bankly está retornando 500 com o body null. Mockar teste.")

	card, err := c.card.GetNextStatusByProxy(c.ctx, "2229131000063855144")
	c.assert.NoError(err)
	c.assert.NotNil(card)
}

func (c *CardTestSuite) TestGetNextStatusByProxy_NOT_FOUND() {
	// TODO Mockar teste
	c.T().Skip("O bankly está retornando 500 com o body null. Mockar teste.")

	card, err := c.card.GetNextStatusByProxy(c.ctx, "00000000000000")
	c.assert.Error(err)
	c.assert.Nil(card)
	c.assert.Equal("Code: 404 - Messages: not found", err.Error())
}

func (c *CardTestSuite) TestGetCardByAccount_OK() {
	// TODO Mockar teste
	c.T().Skip("Bankly está retornando erro 500. Mockar teste.")

	card, err := c.card.GetCardByAccount(c.ctx, "236268", "0001", "93707422046")
	c.assert.NoError(err)
	c.assert.NotNil(card)
}

func (c *CardTestSuite) TestGetCardByAccount_NOT_FOUND() {
	// TODO Mockar teste
	c.T().Skip("Bankly está retornando erro 500. Mockar teste.")

	card, err := c.card.GetCardByAccount(c.ctx, "0", "0", "0")
	c.assert.Error(err)
	c.assert.Nil(card)
	c.assert.Equal("Code: 404 - Messages: not found", err.Error())
}

/*
	func (c *CardTestSuite) TestCreateCardVirtual_OK() {
		c.CancelCard("21632071000187")

		createCardModel := createCardModel("21632071000187", "202142", bankly.VirtualCardType)

		card, err := c.card.CreateCard(c.ctx, &createCardModel)

		c.assert.NoError(err)
		c.assert.NotNil(card)
	}
*/
func (c *CardTestSuite) TestCreateCardVirtual_INVALID_PARAMETER_EMPTY() {
	createCardModel := createCardModel("1234567", "202142", bankly.VirtualCardType)

	card, err := c.card.CreateCard(c.ctx, &createCardModel)

	c.assert.Error(err)
	c.assert.Nil(card)
}

/*
func (c *CardTestSuite) TestCreateCardPhysical_OK() {
	c.CancelCard("93707422046")

	createCardModel := createCardModel("93707422046", "236268", bankly.PhysicalCardType)

	card, err := c.card.CreateCard(c.ctx, &createCardModel)

	c.assert.NoError(err)
	c.assert.NotNil(card)
}

// TestActivateCardByProxy_OK ...
func (c *CardTestSuite) TestActivateCardByProxy_OK() {
	bankAccount := "202134"
	identifier := "04272617000117"

	// cancel cards
	c.CancelCard(identifier)

	createCardModel := createCardModel(identifier, bankAccount, bankly.PhysicalCardType)

	// TODO mockar o response do bankly de cartão criado com sucesso
	card, err := c.card.CreateCard(c.ctx, &createCardModel)

	c.assert.NoError(err)
	c.assert.NotNil(card)

	time.Sleep(time.Second * 10)

	cardActivateModel := bankly.CardActivateDTO{
		ActivateCode : card.ActivateCode,
		Password : createCardModel.CardData.Password,
	}

	// activate card by proxy
	activateErr := c.card.ActivateCardByProxy(c.ctx, &card.Proxy, &cardActivateModel)

	c.assert.NoError(activateErr)
}
*/

/*
func (c *CardTestSuite) TestUpdatePasswordByProxy_OK() {
	card, err := c.card.GetCardsByIdentifier(c.ctx, "93707422046")

	c.assert.NoError(err)
	c.assert.NotNil(card)

	cardAlterPasswordDTO := bankly.CardAlterPasswordDTO{
		Password: "1235",
	}
	// TODO - Rever sleep devido a criação de cartão nao ser sincrono.
	time.Sleep(time.Second * 10)
	cardActivated, err := c.card.UpdatePasswordByProxy(c.ctx, card[0].Proxy, cardAlterPasswordDTO)

	c.assert.NoError(err)
	c.assert.NotNil(cardActivated)
}


func (c *CardTestSuite) TestGetPCIByProxy_OK() {
	proxy := "2229041000006610201"
	cardPCIDTO := bankly.CardPCIDTO{
		Password: "1234",
	}

	card, err := c.card.GetPCIByProxy(c.ctx, &proxy, &cardPCIDTO)

	c.assert.NoError(err)
	c.assert.NotNil(card)
}

func (c *CardTestSuite) TestGetPCIByProxy_PASSWORD_INVALID_NOT_OK() {
	proxy := "2229041000006610201"
	cardPCIDTO := bankly.CardPCIDTO{
		Password: "1233",
	}

	card, err := c.card.GetPCIByProxy(c.ctx, &proxy, &cardPCIDTO)

	c.assert.Error(err)
	c.assert.Contains(err.Error(), "invalid password")
	c.assert.Nil(card)
}
*/

func (c *CardTestSuite) TestCreateCardPhysical_INVALID_PARAMETER_EMPTY() {
	createCardModel := createCardModel("123456", "202142", bankly.PhysicalCardType)

	card, err := c.card.CreateCard(c.ctx, &createCardModel)

	c.assert.Error(err)
	c.assert.Nil(card)
}

func (c *CardTestSuite) TestGetPCIByProxy_PROXY_INVALID_NOT_OK() {
	proxy := "2229041000006610200"
	cardPCIDTO := bankly.CardPCIDTO{
		Password: "1234",
	}

	card, err := c.card.GetPCIByProxy(c.ctx, &proxy, &cardPCIDTO)

	c.assert.Error(err)
	c.assert.Nil(card)
}

/*
func (c *CardTestSuite) TestAlteredStatusCard_OK() {
	c.CancelCard("93707422046")
}
*/

func (c *CardTestSuite) TestUpdatePasswordByProxy_OK() {

	// TODO Mockar teste
	c.T().Skip("Bankly está retornando erro 500. Mockar teste.")

	model := bankly.CardUpdatePasswordDTO{
		Password: "1234",
	}

	err := c.card.UpdatePasswordByProxy(context.Background(), "2229041000032315297", model)

	c.assert.NoError(err)
}

func (c *CardTestSuite) TestUpdatePasswordByProxy_InvalidProxy() {

	// TODO Mockar teste
	c.T().Skip("Bankly está retornando erro 500. Mockar teste.")

	model := bankly.CardUpdatePasswordDTO{
		Password: "1234",
	}

	err := c.card.UpdatePasswordByProxy(context.Background(), "2200000000000000000", model)

	c.assert.Error(err)
	c.assert.Contains(err.Error(), "not found")
}

func (c *CardTestSuite) TestUpdatePasswordByProxy_PasswordIsEmpty() {
	model := bankly.CardUpdatePasswordDTO{}

	err := c.card.UpdatePasswordByProxy(context.Background(), "2229041000032315297", model)

	c.assert.Error(err)
	c.assert.Contains(err.Error(), "validation failed for password")

	model.Password = ""

	err = c.card.UpdatePasswordByProxy(context.Background(), "2229041000032315297", model)

	c.assert.Error(err)
	c.assert.Contains(err.Error(), "validation failed for password")

	model.Password = "123"

	err = c.card.UpdatePasswordByProxy(context.Background(), "2229041000032315297", model)

	c.assert.Error(err)
	c.assert.Contains(err.Error(), "validation failed for password")
}

func (c *CardTestSuite) TestGetTrackingByProxy_OK() {
	// TODO Mockar teste
	c.T().Skip("O bankly está retornando 500 com o body null. Mockar teste.")

	proxy := "2229041000088256610"
	tracking, err := c.card.GetTrackingByProxy(c.ctx, &proxy)

	c.assert.NoError(err)
	c.assert.NotEmpty(tracking)
}

func (c *CardTestSuite) TestGetTrackingByProxy_NOT_OK() {
	proxy := "22290410000882566101"
	tracking, err := c.card.GetTrackingByProxy(c.ctx, &proxy)

	c.assert.Error(err)
	c.assert.Nil(tracking)
}

func (c *CardTestSuite) CancelCard(identifier string) {
	card, err := c.card.GetCardsByIdentifier(context.Background(), identifier)

	c.assert.NoError(err)

	for _, elm := range card {
		if elm.Status != "CanceledByCustomer" {
			c.mockAlterCardCanceled(elm.Proxy)
		}
	}
}

// createCardModel ...
func createCardModel(identifier, bankAccountNumber string, cardType bankly.CardType) bankly.CardCreateDTO {
	return bankly.CardCreateDTO{
		CardData: bankly.CardCreateRequest{
			DocumentNumber: grok.OnlyDigits(identifier),
			CardName:       grok.ToTitle("NOME DA PESSOA"),
			Alias:          grok.ToTitle("CONTBANK"),
			BankAgency:     "0001",
			BankAccount:    grok.OnlyLettersOrDigits(bankAccountNumber),
			Password:       "1234",
			Address:        *createCardAddressModel(),
		},
		CardType: cardType,
	}
}

// createCardAddressModel ...
func createCardAddressModel() *bankly.CardAddress {
	complement := "Apto 1231"
	return &bankly.CardAddress{
		ZipCode:      "01307012",
		Address:      "Rua Dona Antônia de Queirós",
		Number:       "888",
		Complement:   &complement,
		Neighborhood: "Consolação",
		City:         "São Paulo",
		State:        "SP",
		Country:      "BR",
	}
}
