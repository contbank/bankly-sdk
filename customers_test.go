package bankly_test

import (
	"github.com/contbank/grok"
	"math/rand"
	"os"
	"strings"
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

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

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
	s.assert.NotNil(response.DocumentNumber)
	s.assert.NotNil(response.RegisterName)
	s.assert.NotNil(response.Phone)
	s.assert.NotNil(response.Address)
	s.assert.NotNil(response.Email)
	s.assert.NotNil(response.MotherName)
	s.assert.NotNil(response.BirthDate)
}

func (s *CustomersTestSuite) TestFindRegistrationErrorNotFound() {
	response, err := s.customers.FindRegistration(grok.GeneratorCPF())

	s.assert.Error(err)
	s.assert.Contains(err.Error(), "not found")
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

func (s *CustomersTestSuite) TestCreateAndFindRegistration() {
	document := grok.GeneratorCPF()
	cellphone := grok.GeneratorCellphone()
	randSurname := randStringBytes(10)
	email := "email_de_teste_" + randSurname + "@contbank.com"

	err := s.createRegistrationWithParams(randSurname, document, cellphone, email)
	s.assert.NoError(err)

	time.Sleep(time.Millisecond)

	registrationResponse, err := s.customers.FindRegistration(document)
	s.assert.NoError(err)
	s.assert.NotNil(registrationResponse)
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

	time.Sleep(time.Millisecond)

	accountResponse, err := s.customers.FindAccounts(document)
	s.assert.NoError(err)
	s.assert.NotNil(accountResponse)
}

func (s *CustomersTestSuite) TestUpdateRegistration() {

	// TODO Verificar com o Bankly o motivo do FindRegistration estar retornando apenas alguns dados.
	s.T().Skip("Aguardar retorno do Bankly para o FindRegistration, que está retornando dados incompletos.")

	document := grok.GeneratorCPF()
	cellphone := grok.GeneratorCellphone()
	randSurname := randStringBytes(10)
	email := "email_de_teste_" + randSurname + "@contbank.com"

	err := s.createRegistrationWithParams(randSurname, document, cellphone, email)
	s.assert.NoError(err)

	time.Sleep(time.Millisecond)

	registrationResponse, err := s.customers.FindRegistration(document)
	s.assert.NoError(err)
	s.assert.NotNil(registrationResponse)

	// update customer
	newRegisterName := "NOVO NOME DA PESSOA VIA UPDATE REQUEST"
	customerUpdateRequest := &bankly.CustomerUpdateRequest {
		RegisterName: newRegisterName,
		SocialName: registrationResponse.SocialName,
		BirthDate: registrationResponse.BirthDate,
		MotherName: registrationResponse.MotherName,
		Phone: &registrationResponse.Phone,
		Email: registrationResponse.Email,
		Address: &registrationResponse.Address,
	}
	s.customers.UpdateRegistration(document, *customerUpdateRequest)
	s.assert.NoError(err)
	s.assert.Nil(err)

	time.Sleep(time.Millisecond)

	// find updated customer
	updatedAccount, err := s.customers.FindRegistration(registrationResponse.DocumentNumber)
	s.assert.NoError(err)
	s.assert.NotNil(updatedAccount)
	s.assert.Equal(registrationResponse.DocumentNumber, updatedAccount.DocumentNumber)
	s.assert.Equal(newRegisterName, updatedAccount.RegisterName)
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
	var src = rand.NewSource(time.Now().UnixNano())
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return sb.String()
}