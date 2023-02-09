package bankly_test

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/contbank/grok"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"testing"
	"time"

	bankly "github.com/contbank/bankly-sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BoletosTestSuite struct {
	suite.Suite
	ctx       context.Context
	assert    *assert.Assertions
	session   *bankly.Session
	boletos   *bankly.Boletos
	customers *bankly.Customers
}

func TestBoletosTestSuite(t *testing.T) {
	suite.Run(t, new(BoletosTestSuite))
}

func (s *BoletosTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.assert = assert.New(s.T())

	session, err := bankly.NewSession(bankly.Config{
		ClientID:     bankly.String(*bankly.GetEnvBanklyClientID()),
		ClientSecret: bankly.String(*bankly.GetEnvBanklyClientSecret()),
		Scopes:       bankly.String("account.read boleto.create boleto.read boleto.delete"),
	})

	s.assert.NoError(err)

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	s.session = session
	s.boletos = bankly.NewBoletos(httpClient, *s.session)
	s.customers = bankly.NewCustomers(httpClient, *s.session)
}

// TestCreateBoleto_TypeLevy ...
func (s *BoletosTestSuite) TestCreateBoleto_TypeLevy() {

	// TODO Mockar teste
	s.T().Skip("Mockar teste. Erro de SCOUTER_QUANTITY devido ao limite de geração de boletos")

	s.ctx = context.WithValue(s.ctx, "Request-Id", primitive.NewObjectID().String())

	request := createBoletoRequest(bankly.Levy)
	response, err := s.boletos.CreateBoleto(s.ctx, request)

	s.assert.NoError(err)
	s.assert.NotEmpty(response)
	s.assert.NotEmpty(response.AuthenticationCode)
	s.assert.Equal(request.Account.Branch, response.Account.Branch)
	s.assert.Equal(request.Account.Number, response.Account.Number)
}

// TestCreateBoleto_TypeDeposit ...
func (s *BoletosTestSuite) TestCreateBoleto_TypeDeposit() {

	// TODO Mockar teste
	s.T().Skip("Mockar teste. Erro de SCOUTER_QUANTITY devido ao limite de geração de boletos")

	s.ctx = context.WithValue(s.ctx, "Request-Id", primitive.NewObjectID().String())

	request := createBoletoRequest(bankly.Deposit)
	response, err := s.boletos.CreateBoleto(s.ctx, request)

	s.assert.NoError(err)
	s.assert.NotEmpty(response)
	s.assert.NotEmpty(response.AuthenticationCode)
	s.assert.Equal(request.Account.Branch, response.Account.Branch)
	s.assert.Equal(request.Account.Number, response.Account.Number)
}

// TestCreateBoleto_InvalidNameLength ...
func (s *BoletosTestSuite) TestCreateBoleto_InvalidNameLength() {
	s.ctx = context.WithValue(s.ctx, "Request-Id", primitive.NewObjectID().String())

	request := createBoletoRequest(bankly.Levy)
	request.Payer.Name = grok.GeneratorIDBase(51)
	response, err := s.boletos.CreateBoleto(s.ctx, request)

	s.assert.Error(err)
	s.assert.Empty(response)
}

// TestCreateBoleto_InvalidTradeNameLength ...
func (s *BoletosTestSuite) TestCreateBoleto_InvalidTradeNameLength() {
	s.ctx = context.WithValue(s.ctx, "Request-Id", primitive.NewObjectID().String())

	request := createBoletoRequest(bankly.Levy)
	request.Payer.Name = grok.GeneratorIDBase(81)
	response, err := s.boletos.CreateBoleto(s.ctx, request)

	s.assert.Error(err)
	s.assert.Empty(response)
}

// TestFindBoleto ...
func (s *BoletosTestSuite) TestFindBoleto() {

	// TODO Mockar teste
	s.T().Skip("Mockar teste. Erro de SCOUTER_QUANTITY devido ao limite de geração de boletos")
	
	s.ctx = context.WithValue(s.ctx, "Request-Id", primitive.NewObjectID().String())

	// create boleto
	request := createBoletoRequest(bankly.Deposit)
	response, err := s.boletos.CreateBoleto(s.ctx, request)

	s.assert.NoError(err)
	s.assert.NotEmpty(response)

	findRequest := &bankly.FindBoletoRequest{
		AuthenticationCode: response.AuthenticationCode,
		Account: &bankly.Account{
			Branch: response.Account.Branch,
			Number: response.Account.Number,
		},
	}

	// find boleto
	r, err := s.boletos.FindBoleto(s.ctx, findRequest)

	s.assert.NoError(err)
	s.assert.NotNil(r)
	s.assert.Equal(request.Account.Number, r.Account.Number)
	s.assert.Equal(findRequest.AuthenticationCode, r.AuthenticationCode)
}

// TestCancelBoleto
func (s *BoletosTestSuite) TestCancelBoleto() {
	s.T().Skip("Necessário mockar")

	s.ctx = context.WithValue(s.ctx, "Request-Id", primitive.NewObjectID().String())

	// create boleto
	request := createBoletoRequest(bankly.Deposit)
	response, err := s.boletos.CreateBoleto(s.ctx, request)

	s.assert.NoError(err)
	s.assert.NotEmpty(response)

	// time.Sleep(time.Second * 3)

	// cancel boleto
	err = s.boletos.CancelBoleto(s.ctx,
		&bankly.CancelBoletoRequest{
			AuthenticationCode: response.AuthenticationCode,
			Account: &bankly.Account{
				Number: response.Account.Number,
				Branch: response.Account.Branch,
			},
		})

	s.assert.NoError(err)
}

// createBoletoRequest create BoletoRequest object to tests
func createBoletoRequest(boletoType bankly.BoletoType) *bankly.BoletoRequest {
	dueDate := time.Now().Add(24 * time.Hour * 10)           // today + 10 days
	limitePaymentDate := time.Now().Add(24 * time.Hour * 60) // today + 60 days
	request := &bankly.BoletoRequest{
		Document: "36183588814",
		Amount:   21.61,
		DueDate:  dueDate,
		Type:     boletoType,
		Account: &bankly.Account{
			Branch: "0001",
			Number: "184152",
		},
		Alias:        aws.String(grok.GeneratorIDBase(15)),
		ClosePayment: limitePaymentDate,
	}
	if boletoType == bankly.Levy || boletoType == bankly.Invoice {
		request.Payer = &bankly.BoletoPayer{
			Name:      grok.GeneratorIDBase(25) + grok.GeneratorIDBase(25),
			TradeName: grok.GeneratorIDBase(80),
			Document:  grok.GeneratorCPF(),
			Address: &bankly.BoletoAddress{
				AddressLine: "Rua da Consolação, 1234",
				ZipCode:     grok.OnlyDigits("01301100"),
				State:       "SP",
				City:        "São Paulo",
			},
		}
		request.Interest = createInterest(dueDate)
		request.Fine = createFine(dueDate)
		request.Discounts = createDiscount(dueDate)
	}
	return request
}

// createInterest create interest object to test
func createInterest(dueDate time.Time) *bankly.BoletoInterest {
	dueDate1D := dueDate.Add(24 * time.Hour)
	return &bankly.BoletoInterest{
		StartDate: *bankly.OnlyDate(&dueDate1D),
		Value:     2.00,
		Type:      bankly.PercentPerMonth,
	}
}

// createFine create fine object to test
func createFine(dueDate time.Time) *bankly.BoletoFine {
	dueDate1D := dueDate.Add(24 * time.Hour)
	return &bankly.BoletoFine{
		StartDate: *bankly.OnlyDate(&dueDate1D),
		Value:     1.75,
		Type:      bankly.Percent,
	}
}

// createDiscount create discount object to test
func createDiscount(dueDate time.Time) *bankly.BoletoDiscounts {
	previousDate := dueDate.Add(24 * time.Hour * -1)
	return &bankly.BoletoDiscounts{
		LimitDate: *bankly.OnlyDate(&previousDate),
		Value:     1.75,
		Type:      bankly.FixedPercentUntilLimitDate,
	}
}
