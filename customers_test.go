package bankly_test

import (
	"github.com/contbank/grok"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/contbank/bankly-sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CustomersTestSuite struct {
	suite.Suite
	assert    *assert.Assertions
	session   *bankly.Session
	customers *bankly.Customers
}

func TestCustomersTestSuite(t *testing.T) {
	suite.Run(t, new(CustomersTestSuite))
}

func (s *CustomersTestSuite) SetupTest() {
	s.assert = assert.New(s.T())

	session, err := bankly.NewSession(bankly.Config{
		ClientID:     bankly.String(os.Getenv("BANKLY_CLIENT_ID")),
		ClientSecret: bankly.String(os.Getenv("BANKLY_CLIENT_SECRET")),
	})

	s.assert.NoError(err)

	s.session = session
	s.customers = bankly.NewCustomers(*s.session)
}

func (s *CustomersTestSuite) TestCreateRegistration() {
	randSurname := randStringBytes(10)
	email := "email_de_teste_" + randSurname + "@contbank.com"

	err := s.createRegistrationWithParams(randSurname, grok.GeneratorCPF(), grok.GeneratorCellphone(), email)

	s.assert.NoError(err)
}

func (s *CustomersTestSuite) TestFindRegistration() {
	response, err := s.customers.FindRegistration("36183588814")

	s.assert.NoError(err)
	s.assert.NotNil(response)
}

func (s *CustomersTestSuite) TestFindRegistrationErrorNotFound() {
	response, err := s.customers.FindRegistration(grok.GeneratorCPF())

	s.assert.Error(err)
	s.assert.EqualError(err, "not found")
	s.assert.Nil(response)
}

func (s *CustomersTestSuite) TestCreateAccountErrorMoreThanOneAccountPerHolder() {
	account, err := s.customers.CreateAccount("54012948083", bankly.PaymentAccount)

	s.assert.Error(err)
	s.assert.Equal(err, bankly.ErrHolderAlreadyHaveAAccount)
	s.assert.Nil(account)
}

func (s *CustomersTestSuite) TestCreateAccountErrorDoesntHaveAnApprovedRegistrationYet() {
	account, err := s.customers.CreateAccount(grok.GeneratorCPF(), bankly.PaymentAccount)

	s.assert.Error(err)
	s.assert.Equal(err, bankly.ErrAccountHolderNotExists)
	s.assert.Nil(account)
}

func (s *CustomersTestSuite) TestFindAccounts() {
	account, err := s.customers.FindAccounts("36183588814")

	s.assert.NoError(err)
	s.assert.NotNil(account)
}

func (s *CustomersTestSuite) TestCreateAndFindAccount() {
	document := grok.GeneratorCPF()
	cellphone := grok.GeneratorCellphone()
	randSurname := randStringBytes(10)
	email := "email_de_teste_" + randSurname + "@contbank.com"

	err := s.createRegistrationWithParams(randSurname, document, cellphone, email)
	s.assert.NoError(err)

	time.Sleep(time.Second)

	account, err := s.customers.CreateAccount(document, bankly.PaymentAccount)
	s.assert.NoError(err)
	s.assert.NotNil(account)

	accountResponse, err := s.customers.FindAccounts(document)
	s.assert.NoError(err)
	s.assert.NotNil(accountResponse)
}

func (s *CustomersTestSuite) createRegistrationWithParams(surname string, document string, cellphone string, email string) error {
	return s.customers.CreateRegistration(bankly.CustomersRequest{
		Document: document,
		Phone: &bankly.Phone{
			CountryCode: "55",
			Number:      cellphone,
		},
		Address: &bankly.Address{
			ZipCode:        "03503030",
			City:           "São Paulo",
			AddressLine:    "Rua Fulano de Tal",
			BuildingNumber: "1000",
			Neighborhood:   "Chácara Califórnia",
			State:          "SP",
			Country:        "BR",
		},
		RegisterName: "Nome da Pessoa " + surname,
		BirthDate:    time.Date(1993, time.March, 25, 0, 0, 0, 0, time.UTC),
		MotherName:   "Nome da Mãe da Pessoa " + surname,
		Email:        email,
	})
}

func randStringBytes(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}