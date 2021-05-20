package bankly_test

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	bankly "github.com/contbank/bankly-sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BoletosTestSuite struct {
	suite.Suite
	assert    *assert.Assertions
	session   *bankly.Session
	boletos   *bankly.Boletos
	customers *bankly.Customers
}

func TestBoletosTestSuite(t *testing.T) {
	suite.Run(t, new(BoletosTestSuite))
}

func (s *BoletosTestSuite) SetupTest() {
	s.assert = assert.New(s.T())

	session, err := bankly.NewSession(bankly.Config{
		ClientID:     bankly.String(os.Getenv("BANKLY_CLIENT_ID")),
		ClientSecret: bankly.String(os.Getenv("BANKLY_CLIENT_SECRET")),
	})

	s.assert.NoError(err)

	s.session = session
	s.boletos = bankly.NewBoletos(*s.session)
	s.customers = bankly.NewCustomers(*s.session)
}

func (s *BoletosTestSuite) TestCreateBoleto() {
	document, account := s.createCustomerAndAccount()
	boleto := s.createBoleto(document, account)

	res, err := s.boletos.CreateBoleto(boleto)

	s.assert.NoError(err)
	s.assert.NotEmpty(res.Account)
	s.assert.NotEmpty(res.Account.Branch)
	s.assert.NotEmpty(res.Account.Number)
	s.assert.NotEmpty(res.AuthenticationCode)
}

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

	findReq := &bankly.FindBoletoRequest{
		AuthenticationCode: res.AuthenticationCode,
		Account: &bankly.Account{
			Branch: res.Account.Branch,
			Number: res.Account.Number,
		},
	}

	r, err := s.boletos.FindBoleto(findReq)
	s.assert.NoError(err)
	s.assert.Equal("Cancelled", r.Status)
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

	findReq := &bankly.FindBoletoRequest{
		AuthenticationCode: res.AuthenticationCode,
		Account: &bankly.Account{
			Branch: res.Account.Branch,
			Number: res.Account.Number,
		},
	}

	time.Sleep(3 * time.Second)

	r, err := s.boletos.FindBoleto(findReq)
	s.assert.NoError(err)
	s.assert.NotNil(r)
	s.assert.Equal("Settled", r.Status)
	s.assert.NotEmpty(r.Payments)
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
