package bankly_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	bankly "github.com/contbank/bankly-sdk"
	"github.com/contbank/grok"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BusinessTestSuite struct {
	suite.Suite
	assert   *assert.Assertions
	session  *bankly.Session
	business *bankly.Business
}

func TestBusinessTestSuite(t *testing.T) {
	suite.Run(t, new(BusinessTestSuite))
}

func (s *BusinessTestSuite) SetupTest() {
	s.assert = assert.New(s.T())

	session, err := bankly.NewSession(bankly.Config{
		ClientID:     bankly.String(*bankly.GetEnvBanklyClientID()),
		ClientSecret: bankly.String(*bankly.GetEnvBanklyClientSecret()),
		Scopes:       bankly.String("business.write business.read business.cancel account.read account.create account.close"),
	})

	s.assert.NoError(err)

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	s.session = session
	s.business = bankly.NewBusiness(httpClient, *s.session)
}

func (s *BusinessTestSuite) TestCreateBusiness_TypeEI_SizeME() {
	// TODO Mockar teste
	s.T().Skip("Bankly está retornando erro 403 no endpoint de business. Mockar este teste.")

	businessRequest := createBusinessRequest(grok.GeneratorCNPJ(), bankly.BusinessTypeEI, bankly.BusinessSizeME)

	err := s.business.CreateBusinessRegistration(context.Background(), businessRequest)

	s.assert.NoError(err)
	s.assert.Nil(err)
}

func (s *BusinessTestSuite) TestCreateBusiness_TypeMEI_SizeMEI() {
	// TODO Mockar teste
	s.T().Skip("Bankly está retornando erro 403 no endpoint de business. Mockar este teste.")

	businessRequest := createBusinessRequest(grok.GeneratorCNPJ(), bankly.BusinessTypeMEI, bankly.BusinessSizeMEI)

	err := s.business.CreateBusinessRegistration(context.Background(), businessRequest)

	s.assert.NoError(err)
	s.assert.Nil(err)
}

func (s *BusinessTestSuite) TestCreateBusiness_TypeEIRELI_SizeEPP() {
	// TODO Mockar teste
	s.T().Skip("Bankly está retornando erro 403 no endpoint de business. Mockar este teste.")

	businessRequest := createBusinessRequest(grok.GeneratorCNPJ(), bankly.BusinessTypeEIRELI, bankly.BusinessSizeEPP)

	err := s.business.CreateBusinessRegistration(context.Background(), businessRequest)

	s.assert.NoError(err)
	s.assert.Nil(err)
}

func (s *BusinessTestSuite) TestCreateBusinessErrorInvalidTypeMEIAndSizeEPP() {
	// TODO Mockar teste
	s.T().Skip("Bankly está retornando erro 403 no endpoint de business. Mockar este teste.")

	businessRequest := createBusinessRequest(grok.GeneratorCNPJ(), bankly.BusinessTypeMEI, bankly.BusinessSizeEPP)

	err := s.business.CreateBusinessRegistration(context.Background(), businessRequest)

	s.assert.Error(err)

	banklyErr, ok := bankly.ParseErr(err)

	s.assert.True(ok)
	s.assert.Equal(bankly.ErrInvalidBusinessSize, banklyErr.GrokError)
}

func (s *BusinessTestSuite) TestCreateBusinessErrorInvalidTypeMEIAndSizeME() {
	// TODO Mockar teste
	s.T().Skip("Bankly está retornando erro 403 no endpoint de business. Mockar este teste.")

	businessRequest := createBusinessRequest(grok.GeneratorCNPJ(), bankly.BusinessTypeMEI, bankly.BusinessSizeME)

	err := s.business.CreateBusinessRegistration(context.Background(), businessRequest)

	s.assert.Error(err)

	banklyErr, ok := bankly.ParseErr(err)

	s.assert.True(ok)
	s.assert.Equal(bankly.ErrInvalidBusinessSize, banklyErr.GrokError)
}

func (s *BusinessTestSuite) TestUpdateBusinessName() {

	// TODO O PATCH da Business não está funcionando na Bankly. Aguardando retorno deles para ativar este teste.
	s.T().Skip("O PATCH da Business não está funcionando na Bankly. Aguardando retorno deles para ativar este teste.")

	// create business
	businessRequest := createBusinessRequest(grok.GeneratorCNPJ(), bankly.BusinessTypeMEI, bankly.BusinessSizeMEI)
	err := s.business.CreateBusinessRegistration(context.Background(), businessRequest)
	s.assert.NoError(err)
	s.assert.Nil(err)

	// update business
	newBusinessName := "NOVO EMAIL DA EMPRESA VIA UPDATE REQUEST"
	businessUpdateRequest := &bankly.BusinessUpdateRequest{
		BusinessName: newBusinessName,
	}
	err = s.business.UpdateBusiness(context.Background(), businessRequest.DocumentNumber, *businessUpdateRequest)
	s.assert.NoError(err)
	s.assert.Nil(err)

	// find updated business
	updatedAccount, err := s.business.FindBusiness(context.Background(), businessRequest.DocumentNumber)
	s.assert.NoError(err)
	s.assert.NotNil(updatedAccount)
	s.assert.Equal(businessRequest.DocumentNumber, updatedAccount.Document)
	s.assert.Equal(newBusinessName, updatedAccount.BusinessName)
}

func (s *BusinessTestSuite) TestUpdateBusinessEmailAndBusinessTypeAndBusinessType() {

	// TODO O PATCH da Business não está funcionando na Bankly. Aguardando retorno deles para ativar este teste.
	s.T().Skip("O PATCH da Business não está funcionando na Bankly. Aguardando retorno deles para ativar este teste.")

	// create business
	businessRequest := createBusinessRequest(grok.GeneratorCNPJ(), bankly.BusinessTypeMEI, bankly.BusinessSizeMEI)
	err := s.business.CreateBusinessRegistration(context.Background(), businessRequest)
	s.assert.NoError(err)
	s.assert.Nil(err)

	// update business
	newEmail := "novo_email_" + businessRequest.DocumentNumber + "@contbank.com"
	businessUpdateRequest := &bankly.BusinessUpdateRequest{
		BusinessEmail: newEmail,
		BusinessType:  bankly.BusinessTypeEI,
		BusinessSize:  bankly.BusinessSizeME,
	}
	err = s.business.UpdateBusiness(context.Background(), businessRequest.DocumentNumber, *businessUpdateRequest)
	s.assert.NoError(err)
	s.assert.Nil(err)

	// find business updated
	updatedAccount, err := s.business.FindBusiness(context.Background(), "59619372000143")
	s.assert.NoError(err)
	s.assert.NotNil(updatedAccount)
	s.assert.Equal(businessRequest.DocumentNumber, updatedAccount.Document)
	s.assert.Equal(bankly.BusinessTypeEI, updatedAccount.BusinessType)
	s.assert.Equal(bankly.BusinessSizeME, updatedAccount.BusinessSize)
	s.assert.Equal(newEmail, updatedAccount.BusinessEmail)

}

func (s *BusinessTestSuite) TestCreateBusinessAccount() {
	// TODO Mockar teste
	s.T().Skip("Bankly está retornando erro 403 no endpoint de business. Mockar este teste.")

	businessRequest := createBusinessRequest(grok.GeneratorCNPJ(), bankly.BusinessTypeEI, bankly.BusinessSizeME)

	err := s.business.CreateBusinessRegistration(context.Background(), businessRequest)
	s.assert.NoError(err)
	s.assert.Nil(err)

	businessAccountRequest := bankly.BusinessAccountRequest{
		Document:    businessRequest.DocumentNumber,
		AccountType: bankly.PaymentAccount,
	}

	time.Sleep(time.Second)

	account, err := s.business.CreateBusinessAccount(context.Background(), businessAccountRequest)
	s.assert.NoError(err)
	s.assert.NotNil(account)
}

func (s *BusinessTestSuite) TestCreateBusinessAccountErrorDoesntHaveAnApprovedRegistration() {
	// TODO Mockar teste
	s.T().Skip("Bankly está retornando erro 403 no endpoint de business. Mockar este teste.")

	businessAccountRequest := bankly.BusinessAccountRequest{
		Document:    bankly.GeneratorCNPJ(),
		AccountType: bankly.PaymentAccount,
	}
	account, err := s.business.CreateBusinessAccount(context.Background(), businessAccountRequest)

	s.assert.Error(err)

	banklyErr, ok := bankly.ParseErr(err)

	s.assert.True(ok)
	s.assert.Equal(bankly.ErrAccountHolderNotExists, banklyErr.GrokError)
	s.assert.Nil(account)
}

func (s *BusinessTestSuite) TestFindBusiness_APPROVED() {
	// TODO Mockar teste
	s.T().Skip("Bankly está retornando erro 403 no endpoint de business. Mockar este teste.")

	identifier := "59619372000143"
	account, err := s.business.FindBusiness(context.Background(), identifier)

	s.assert.NoError(err)
	s.assert.NotNil(account)
	s.assert.Equal(identifier, account.Document)
}

func (s *BusinessTestSuite) TestFindBusiness_PENDING() {
	// TODO Mockar teste
	s.T().Skip("Bankly está retornando erro 403 no endpoint de business. Mockar este teste.")

	identifier := "88427552000121"
	account, err := s.business.FindBusiness(context.Background(), identifier)

	s.assert.NoError(err)
	s.assert.NotNil(account)
	s.assert.Equal(identifier, account.Document)
}

func (s *BusinessTestSuite) TestFindBusiness_NOT_FOUND() {
	// TODO Mockar teste
	s.T().Skip("Bankly está retornando erro 403 no endpoint de business. Mockar este teste.")

	account, err := s.business.FindBusiness(context.Background(), "00000000000000")

	s.assert.Error(err)
	s.assert.Nil(account)
	s.assert.Equal("Code: 404 - Messages: not found", err.Error())
}

func (s *BusinessTestSuite) TestFindBusinessAccounts() {
	// TODO Mockar teste
	s.T().Skip("Bankly está retornando erro 403 no endpoint de business. Mockar este teste.")

	account, err := s.business.FindBusinessAccounts(context.Background(), "59619372000143")

	s.assert.NoError(err)
	s.assert.NotNil(account)
}

func (s *BusinessTestSuite) TestBusinessName_TypeMEI() {
	// TODO Mockar teste
	s.T().Skip("Bankly está retornando erro 403 no endpoint de business. Mockar este teste.")

	businessRequest := createBusinessRequest(grok.GeneratorCNPJ(), bankly.BusinessTypeMEI, bankly.BusinessSizeMEI)
	businessName := businessRequest.LegalRepresentatives[0].RegisterName + " " +
		businessRequest.LegalRepresentatives[0].DocumentNumber
	err := s.business.CreateBusinessRegistration(context.Background(), businessRequest)

	s.assert.NoError(err)
	s.assert.Nil(err)

	time.Sleep(time.Millisecond)

	account, err := s.business.FindBusiness(context.Background(), businessRequest.DocumentNumber)

	s.assert.NoError(err)
	s.assert.NotNil(account)
	s.assert.Equal(businessName, account.BusinessName)
}

func (s *BusinessTestSuite) TestBusinessName_TypeEI() {
	// TODO Mockar teste
	s.T().Skip("Bankly está retornando erro 403 no endpoint de business. Mockar este teste.")

	businessRequest := createBusinessRequest(grok.GeneratorCNPJ(), bankly.BusinessTypeEI, bankly.BusinessSizeEPP)
	err := s.business.CreateBusinessRegistration(context.Background(), businessRequest)

	s.assert.NoError(err)
	s.assert.Nil(err)

	time.Sleep(time.Millisecond)

	account, err := s.business.FindBusiness(context.Background(), businessRequest.DocumentNumber)

	s.assert.NoError(err)
	s.assert.NotNil(account)
	s.assert.Equal(businessRequest.BusinessName, account.BusinessName)
}

func (s *BusinessTestSuite) TestBusinessName_TypeEIRELI() {
	// TODO Mockar teste
	s.T().Skip("Bankly está retornando erro 403 no endpoint de business. Mockar este teste.")

	businessRequest := createBusinessRequest(grok.GeneratorCNPJ(), bankly.BusinessTypeEIRELI, bankly.BusinessSizeEPP)
	err := s.business.CreateBusinessRegistration(context.Background(), businessRequest)

	s.assert.NoError(err)
	s.assert.Nil(err)

	time.Sleep(time.Millisecond)

	account, err := s.business.FindBusiness(context.Background(), businessRequest.DocumentNumber)

	s.assert.NoError(err)
	s.assert.NotNil(account)
	s.assert.Equal(businessRequest.BusinessName, account.BusinessName)
}

// createBusinessRequest ...
func createBusinessRequest(document string, businessType bankly.BusinessType,
	businessSize bankly.BusinessSize) bankly.BusinessRequest {

	randomName := bankly.RandStringBytes(10)
	email := "email_de_teste_" + randomName + "@contbank.com"

	return bankly.BusinessRequest{
		DocumentNumber:  bankly.OnlyDigits(document),
		BusinessName:    "Nome da Empresa " + randomName,
		TradingName:     "Nome Fantasia " + randomName,
		BusinessEmail:   email,
		BusinessType:    businessType,
		BusinessSize:    businessSize,
		BusinessAddress: createAddress(),
		LegalRepresentatives: []*bankly.LegalRepresentative{
			createLegalRepresentative(),
		},
		DeclaredAnnualBilling: bankly.NOT_DECLARED,
	}
}

func createLegalRepresentative() *bankly.LegalRepresentative {
	randSurname := randStringBytes(10)
	return &bankly.LegalRepresentative{
		DocumentNumber: grok.GeneratorCPF(),
		RegisterName:   "Nome do Representante Legal " + randSurname,
		Phone: &bankly.Phone{
			CountryCode: "55",
			Number:      bankly.GeneratorCellphone(),
		},
		Address:    createAddress(),
		BirthDate:  time.Date(1990, time.June, 29, 0, 0, 0, 0, time.UTC),
		MotherName: "Nome da Mãe do Representante Legal " + randSurname,
		Email:      "email_legal_representative_" + randSurname + "@contbank.com",
	}
}

func createAddress() *bankly.Address {
	return &bankly.Address{
		ZipCode:        "05410900",
		City:           "São Paulo",
		AddressLine:    "Rua Qualquer XYZ",
		BuildingNumber: "2000",
		Neighborhood:   "Pinheiros",
		State:          "SP",
		Country:        "BR",
	}
}
