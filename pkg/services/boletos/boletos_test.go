package bankly_test

import (
	"github.com/contbank/bankly-sdk/pkg/services/authentication"
	boletos "github.com/contbank/bankly-sdk/pkg/services/boletos"
	customers "github.com/contbank/bankly-sdk/pkg/services/customers"
	utils "github.com/contbank/bankly-sdk/pkg/utils"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BoletosTestSuite struct {
	suite.Suite
	assert    *assert.Assertions
	session   *bankly.Session
	boletos   *boletos.Boletos
	customers *customers.Customers
}

func TestBoletosTestSuite(t *testing.T) {
	suite.Run(t, new(BoletosTestSuite))
}

func (s *BoletosTestSuite) SetupTest() {
	s.assert = assert.New(s.T())

	session, err := bankly.NewSession(bankly.Config{
		ClientID:     bankly.String(*utils.GetEnvBanklyClientID()),
		ClientSecret: bankly.String(*utils.GetEnvBanklyClientSecret()),
		Scopes:       bankly.String("account.read boleto.create boleto.read boleto.delete"),
	})

	s.assert.NoError(err)

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	s.session = session
	s.boletos = boletos.NewBoletos(httpClient, *s.session)
	s.customers = customers.NewCustomers(httpClient, *s.session)
}

/*
func (s *BoletosTestSuite) TestCreateBoleto() {

	res, err := s.boletos.CreateBoleto(context.Background(), &bankly.BoletoRequest{
		Document: "36183588814",
		Amount:   100,
		DueDate:  time.Now().Add(24 * time.Hour),
		Type:     bankly.Levy,
		Account: &bankly.Account{
			Branch: "0001",
			Number: "184152",
		},
		Payer: &bankly.Payer{
			Name:      grok.GeneratorIDBase(50) + grok.GeneratorIDBase(50),
			TradeName: grok.GeneratorIDBase(500),
			Document:  "36183588814",
			Address: &bankly.BoletoAddress{
				AddressLine: grok.GeneratorIDBase(24) + " " + grok.GeneratorIDBase(34),
				ZipCode:     grok.GeneratorIDBase(500),
				State:       grok.GeneratorIDBase(500),
				City:        grok.GeneratorIDBase(500),
			},
		},
	})

	s.assert.NoError(err)
	s.assert.NotEmpty(res.Account)
	s.assert.NotEmpty(res.Account.Branch)
	s.assert.NotEmpty(res.Account.Number)
	s.assert.NotEmpty(res.AuthenticationCode)
}
*/

/*

func (s *BoletosTestSuite) TestFindBoleto() {
	document, account := s.createCustomerAndAccount()
	boleto := s.createBoleto(document, account)

	res, err := s.boletos.CreateBoleto(boleto)
	s.assert.NoError(err)

	req := &bankly.FindBoletoRequest{
		AuthenticationCode: res.AuthenticationCode,
		Account: &bankly.Account{
			Branch: res.Account.Branch,
			Number: res.Account.Number,
		},
	}

	r, err := s.boletos.FindBoleto(req)
	s.assert.NoError(err)
	s.assert.NotNil(r)
}

// func (s *BoletosTestSuite) TestFilterBoleto() {
// 	document, account := s.createCustomerAndAccount()
// 	boleto := s.createBoleto(document, account)

// 	_, err := s.boletos.CreateBoleto(boleto)
// 	s.assert.NoError(err)

// 	r, err := s.boletos.FilterBoleto(time.Now().Add(-24 * time.Hour))
// 	s.assert.NoError(err)
// 	s.assert.NotNil(r)
// }

// func (s *BoletosTestSuite) TestFindBoletoByBarCode() {
// 	document, account := s.createCustomerAndAccount()
// 	boleto := s.createBoleto(document, account)

// 	res, err := s.boletos.CreateBoleto(boleto)
// 	s.assert.NoError(err)

// 	req := &bankly.FindBoletoRequest{
// 		AuthenticationCode: res.AuthenticationCode,
// 		Account: &bankly.Account{
// 			Branch: res.Account.Branch,
// 			Number: res.Account.Number,
// 		},
// 	}

// 	r, err := s.boletos.FindBoleto(req)
// 	s.assert.NoError(err)

// 	r2, err := s.boletos.FindBoletoByBarCode(r.Digitable)
// 	s.assert.NoError(err)
// 	s.assert.NotNil(r2)
// }

func (s *BoletosTestSuite) TestDownloadBoleto() {
	document, account := s.createCustomerAndAccount()
	boleto := s.createBoleto(document, account)

	res, err := s.boletos.CreateBoleto(boleto)
	s.assert.NoError(err)

	f, err := ioutil.TempFile(os.TempDir(), "temp-")
	s.assert.NoError(err)

	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()

	err = s.boletos.DownloadBoleto(res.AuthenticationCode, f)
	s.assert.NoError(err)

	stats, err := f.Stat()
	s.assert.NoError(err)
	s.assert.NotZero(stats.Size())
}


func (s *BoletosTestSuite) TestCancelBoleto() {
	document, account := s.createCustomerAndAccount()
	boleto := s.createBoleto(document, account)

	res, err := s.boletos.CreateBoleto(boleto)
	s.assert.NoError(err)

	req := &bankly.CancelBoletoRequest{
		AuthenticationCode: res.AuthenticationCode,
		Account: &bankly.Account{
			Number: res.Account.Number,
			Branch: res.Account.Branch,
		},
	}

	err = s.boletos.CancelBoleto(req)
	s.assert.NoError(err)

}

func (s *BoletosTestSuite) TestPayBoleto() {
	document, account := s.createCustomerAndAccount()
	boleto := s.createBoleto(document, account)

	res, err := s.boletos.CreateBoleto(boleto)
	s.assert.NoError(err)

	req := &bankly.PayBoletoRequest{
		AuthenticationCode: res.AuthenticationCode,
		Account: &bankly.Account{
			Number: res.Account.Number,
			Branch: res.Account.Branch,
		},
	}

	err = s.boletos.PayBoleto(req)
	s.assert.NoError(err)
}

func (s *BoletosTestSuite) createBoleto(document string, account *bankly.Account) *bankly.BoletoRequest {
	return &bankly.BoletoRequest{
		Alias:       bankly.String(bankly.RandStringBytes(12)),
		Document:    document,
		Amount:      1547.55,
		DueDate:     time.Now().Add(48 * time.Hour),
		EmissionFee: false,
		Type:        bankly.Levy,
		Account:     account,
		Payer:       s.createPayer(),
	}
}

func (s *BoletosTestSuite) createPayer() *bankly.Payer {
	return &bankly.Payer{
		Name:      bankly.RandStringBytes(9),
		TradeName: bankly.RandStringBytes(15),
		Document:  bankly.GeneratorCPF(),
		Address:   s.createAddress(),
	}
}

func (s *BoletosTestSuite) createAddress() *bankly.Address {
	return &bankly.Address{
		ZipCode:        "03503030",
		City:           "São Paulo",
		AddressLine:    "Rua Fulano de Tal",
		BuildingNumber: "1000",
		Neighborhood:   "Chácara Califórnia",
		State:          "SP",
		Country:        "BR",
	}
}

func (s *BoletosTestSuite) createCustomerAndAccount() (string, *bankly.Account) {
	document := bankly.GeneratorCPF()
	cellphone := bankly.GeneratorCellphone()
	surname := bankly.RandStringBytes(10)
	email := "email_de_teste_" + surname + "@contbank.com"

	req := bankly.CustomersRequest{
		Document: document,
		Phone: &bankly.Phone{
			CountryCode: "55",
			Number:      cellphone,
		},
		Address:      s.createAddress(),
		RegisterName: "Nome da Pessoa " + surname,
		BirthDate:    time.Date(1993, time.March, 25, 0, 0, 0, 0, time.UTC),
		MotherName:   "Nome da Mãe da Pessoa " + surname,
		Email:        email,
	}

	err := s.customers.CreateRegistration(req)
	s.assert.NoError(err)

	time.Sleep(2 * time.Second)

	account, err := s.customers.CreateAccount(document, bankly.PaymentAccount)
	s.assert.NoError(err)
	s.assert.NotNil(account)

	return document, &bankly.Account{
		Branch: account.Branch,
		Number: account.Number,
	}
}
*/
