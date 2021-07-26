package bankly_test

import (
	"math/rand"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	bankly "github.com/contbank/bankly-sdk"
	"github.com/contbank/grok"
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

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	s.session = session
	s.customers = bankly.NewCustomers(httpClient, *s.session)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

func (s *CustomersTestSuite) TestCreateRegistration() {
	randSurname := bankly.RandStringBytes(10)
	email := "email_de_teste_" + randSurname + "@contbank.com"

	err := s.createRegistrationWithParams(randSurname, bankly.GeneratorCPF(), bankly.GeneratorCellphone(), email)

	s.assert.NoError(err)
}

func (s *CustomersTestSuite) TestFindRegistration() {

	// TODO corrigir este teste. Pode ser que não tenha esta conta.
	s.T().Skip("Criar a conta e depois dar um filter.")

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
	response, err := s.customers.FindRegistration(bankly.GeneratorCPF())

	s.assert.Error(err)
	s.assert.Contains(err.Error(), "not found")
	s.assert.Nil(response)
}

func (s *CustomersTestSuite) TestCreateAccountErrorMoreThanOneAccountPerHolder() {
	account, err := s.customers.CreateAccount("54012948083", bankly.PaymentAccount)

	s.assert.Error(err)

	banklyErr, ok := bankly.ParseErr(err)

	s.assert.True(ok)
	s.assert.Equal(bankly.ErrHolderAlreadyHaveAAccount, banklyErr.GrokError)
	s.assert.Nil(account)
}

func (s *CustomersTestSuite) TestCreateAccountErrorDoesntHaveAnApprovedRegistrationYet() {
	account, err := s.customers.CreateAccount(bankly.GeneratorCPF(), bankly.PaymentAccount)

	s.assert.Error(err)

	banklyErr, ok := bankly.ParseErr(err)

	s.assert.True(ok)
	s.assert.Equal(bankly.ErrAccountHolderNotExists, banklyErr.GrokError)
	s.assert.Nil(account)
}

func (s *CustomersTestSuite) TestFindAccounts() {
	account, err := s.customers.FindAccounts("36183588814")

	s.assert.NoError(err)
	s.assert.NotNil(account)
}

func (s *CustomersTestSuite) TestCreateAndFindRegistration() {

	// TODO corrigir birthdate
	s.T().Skip("Erro de parse no birthdate : parsing time \"\"1993-03-25T00:00:00\"\" as \"\"2006-01-02T15:04:05Z07:00\"\": cannot parse \"\"\" as \"Z07:00\"")

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
	document := bankly.GeneratorCPF()
	cellphone := bankly.GeneratorCellphone()
	randSurname := bankly.RandStringBytes(10)
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

	// TODO O PATCH da Customers não está funcionando na Bankly. Aguardando retorno deles para ativar este teste.
	s.T().Skip("O PATCH da Customers não está funcionando na Bankly. Aguardando retorno deles para ativar este teste.")

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
	customerUpdateRequest := &bankly.CustomerUpdateRequest{
		RegisterName: newRegisterName,
		SocialName:   registrationResponse.SocialName,
		BirthDate:    registrationResponse.BirthDate,
		MotherName:   registrationResponse.MotherName,
		Phone:        &registrationResponse.Phone,
		Email:        registrationResponse.Email,
		Address:      &registrationResponse.Address,
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
