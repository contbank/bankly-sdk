package bankly_test

import (
	"context"
	models "github.com/contbank/bankly-sdk/pkg/models"
	"github.com/contbank/bankly-sdk/pkg/services/authentication"
	cards "github.com/contbank/bankly-sdk/pkg/services/cards"
	utils "github.com/contbank/bankly-sdk/pkg/utils"
	"net/http"
	"testing"
	"time"

	"github.com/contbank/grok"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CardTestSuite struct {
	suite.Suite
	assert *assert.Assertions
	ctx    context.Context
	card   *cards.Card
}

func (c *CardTestSuite) mockAlterCardCanceled(proxy string) {
	alterCardBody := models.CardUpdateStatusDTO{
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
		ClientID:     bankly.String(*utils.GetEnvBanklyClientID()),
		ClientSecret: bankly.String(*utils.GetEnvBanklyClientSecret()),
		Scopes:       bankly.String("card.create card.update card.read card.pci.password.update"),
	})
	s.assert.NoError(err)

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	newHttpClient := utils.BanklyHttpClient{
		Session:        *session,
		HttpClient:     httpClient,
		Authentication: bankly.NewAuthentication(httpClient, *session),
	}

	s.card = cards.NewCard(newHttpClient)
}

func (c *CardTestSuite) TestGetCardsByIdentifier_OK() {
	card, err := c.card.GetCardsByIdentifier(c.ctx, "93707422046")
	c.assert.NoError(err)
	c.assert.NotNil(card)
}

func (c *CardTestSuite) TestGetCardsByIdentifier_NOT_FOUND() {
	card, err := c.card.GetCardsByIdentifier(c.ctx, "00000000000000")

	c.assert.Error(err)
	c.assert.Nil(card)
	c.assert.Equal("Code: 404 - Messages: not found", err.Error())
}

func (c *CardTestSuite) TestGetTransactionByProxy_OK() {
	proxy := "2229041000054459062"
	card, err := c.card.GetTransactionsByProxy(c.ctx, &proxy, "1", "2021-01-01", "2021-01-08", "10")
	c.assert.NoError(err)
	c.assert.NotNil(card)
}

func (c *CardTestSuite) TestGetTransactionByProxy_INTERVAL_DATE_NOT_OK() {
	proxy := "2229041000054459062"
	card, err := c.card.GetTransactionsByProxy(c.ctx, &proxy, "1", "2021-01-01", "2021-01-09", "10")
	c.assert.Error(err)
	c.assert.Nil(card)
}

func (c *CardTestSuite) TestGetTransactionByProxy_ENDDATE_NOTFOUND_NOT_OK() {
	proxy := "2229041000054459062"
	card, err := c.card.GetTransactionsByProxy(c.ctx, &proxy, "1", "2021-01-01", "", "10")
	c.assert.Error(err)
	c.assert.Nil(card)
}

func (c *CardTestSuite) TestGetCardByProxy_OK() {
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
	card, err := c.card.GetCardByProxy(c.ctx, "00000000000000")
	c.assert.Error(err)
	c.assert.Nil(card)
	c.assert.Equal("Code: 404 - Messages: not found", err.Error())
}

func (c *CardTestSuite) TestGetCardByActivateCode_OK() {
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
	card, err := c.card.GetCardByActivateCode(c.ctx, "000000000000")
	c.assert.Error(err)
	c.assert.Nil(card)
	c.assert.Equal("Code: 404 - Messages: not found", err.Error())
}

func (c *CardTestSuite) TestGetNextStatusByProxy_OK() {
	card, err := c.card.GetNextStatusByProxy(c.ctx, "2229131000063855144")
	c.assert.NoError(err)
	c.assert.NotNil(card)
}

func (c *CardTestSuite) TestGetNextStatusByProxy_NOT_FOUND() {
	card, err := c.card.GetNextStatusByProxy(c.ctx, "00000000000000")
	c.assert.Error(err)
	c.assert.Nil(card)
	c.assert.Equal("Code: 404 - Messages: not found", err.Error())
}

func (c *CardTestSuite) TestGetCardByAccount_OK() {
	card, err := c.card.GetCardByAccount(c.ctx, "236268", "0001", "93707422046")
	c.assert.NoError(err)
	c.assert.NotNil(card)
}

func (c *CardTestSuite) TestGetCardByAccount_NOT_FOUND() {
	card, err := c.card.GetCardByAccount(c.ctx, "0", "0", "0")
	c.assert.Error(err)
	c.assert.Nil(card)
	c.assert.Equal("Code: 404 - Messages: not found", err.Error())
}

/*
func (c *CardTestSuite) TestCreateCardVirtual_OK() {
	// TODO Mockar teste
	c.T().Skip("Bankly está retornando erro 503 no cancelamento do cartão. Mockar este teste.")

	c.CancelCard("21632071000187")

	createCardModel := createCardModel("21632071000187", "202142", models.VirtualCardType)

	card, err := c.card.CreateCard(c.ctx, &createCardModel)

	c.assert.NoError(err)
	c.assert.NotNil(card)
}
*/
func (c *CardTestSuite) TestCreateCardVirtual_INVALID_PARAMETER_EMPTY() {
	createCardModel := createCardModel("1234567", "202142", models.VirtualCardType)

	card, err := c.card.CreateCard(c.ctx, &createCardModel)

	c.assert.Error(err)
	c.assert.Nil(card)
}

/*
func (c *CardTestSuite) TestCreateCardPhysical_OK() {
	// TODO Mockar teste
	c.T().Skip("Bankly está retornando erro 503 no cancelamento do cartão. Mockar este teste.")

	c.CancelCard("93707422046")

	createCardModel := createCardModel("93707422046", "236268", models.PhysicalCardType)

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

<<<<<<< HEAD:pkg/services/cards/card_test.go
func (c *CardTestSuite) TestCreateCardPhysical_INVALID_PARAMETER_EMPTY() {
	createCardModel := createCardModel("123456", "202142", models.PhysicalCardType)

	card, err := c.card.CreateCard(c.ctx, &createCardModel)

	c.assert.Error(err)
	c.assert.Nil(card)
}
=======
>>>>>>> master:card_test.go

func (c *CardTestSuite) TestGetPCIByProxy_OK() {
	// TODO Mockar teste
	c.T().Skip("Bankly está retornando erro 403 no endpoint do PCI. Mockar este teste.")

	proxy := "2229041000006610201"
	cardPCIDTO := models.CardPCIDTO{
		Password: "1234",
	}

	card, err := c.card.GetPCIByProxy(c.ctx, &proxy, &cardPCIDTO)

	c.assert.NoError(err)
	c.assert.NotNil(card)
}

func (c *CardTestSuite) TestGetPCIByProxy_PASSWORD_INVALID_NOT_OK() {
	// TODO Mockar teste
	c.T().Skip("Bankly está retornando erro 403 no endpoint do PCI. Mockar este teste.")

	proxy := "2229041000006610201"
	cardPCIDTO := models.CardPCIDTO{
		Password: "1233",
	}

	card, err := c.card.GetPCIByProxy(c.ctx, &proxy, &cardPCIDTO)

	c.assert.Error(err)
	c.assert.Contains(err.Error(), "invalid password")
	c.assert.Nil(card)
}
*/

func (c *CardTestSuite) TestCreateCardPhysical_INVALID_PARAMETER_EMPTY() {
	createCardModel := createCardModel("123456", "202142", models.PhysicalCardType)

	card, err := c.card.CreateCard(c.ctx, &createCardModel)

	c.assert.Error(err)
	c.assert.Nil(card)
}

func (c *CardTestSuite) TestGetPCIByProxy_PROXY_INVALID_NOT_OK() {
	proxy := "2229041000006610200"
	cardPCIDTO := models.CardPCIDTO{
		Password: "1234",
	}

	card, err := c.card.GetPCIByProxy(c.ctx, &proxy, &cardPCIDTO)

	c.assert.Error(err)
	c.assert.Nil(card)
}

/*
func (c *CardTestSuite) TestAlteredStatusCard_OK() {
	// TODO Mockar teste
	c.T().Skip("Bankly está retornando erro 503 no cancelamento do cartão. Mockar este teste.")

	c.CancelCard("93707422046")
}
*/
func (c *CardTestSuite) CancelCard(identifier string) {
	card, err := c.card.GetCardsByIdentifier(context.Background(), identifier)

	c.assert.NoError(err)

	for _, elm := range card {
		if elm.Status != "CanceledByCustomer" {
			c.mockAlterCardCanceled(elm.Proxy)
		}
	}
}

func (c *CardTestSuite) TestGetTrackingByProxy_OK() {
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

// createCardModel ...
func createCardModel(identifier, bankAccountNumber string, cardType models.CardType) models.CardCreateDTO {
	return models.CardCreateDTO{
		CardData: models.CardCreateRequest{
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
func createCardAddressModel() *models.CardAddress {
	complement := "Apto 1231"
	return &models.CardAddress{
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