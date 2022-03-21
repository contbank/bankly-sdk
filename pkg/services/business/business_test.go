package bankly_test

import (
	"context"
	errors "github.com/contbank/bankly-sdk/pkg/errors"
	models "github.com/contbank/bankly-sdk/pkg/models"
	"github.com/contbank/bankly-sdk/pkg/services/authentication"
	business "github.com/contbank/bankly-sdk/pkg/services/business"
	utils "github.com/contbank/bankly-sdk/pkg/utils"
	"net/http"
	"testing"
	"time"

	"github.com/contbank/grok"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BusinessTestSuite struct {
	suite.Suite
	assert   *assert.Assertions
	session  *bankly.Session
	business *business.Business
}

func TestBusinessTestSuite(t *testing.T) {
	suite.Run(t, new(BusinessTestSuite))
}

func (s *BusinessTestSuite) SetupTest() {
	s.assert = assert.New(s.T())

	session, err := bankly.NewSession(bankly.Config{
		ClientID:     bankly.String(*utils.GetEnvBanklyClientID()),
		ClientSecret: bankly.String(*utils.GetEnvBanklyClientSecret()),
		Scopes:       bankly.String("business.write business.read business.cancel account.read account.create account.close"),
	})

	s.assert.NoError(err)

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	s.session = session
	s.business = business.NewBusiness(httpClient, *s.session)
}

func (s *BusinessTestSuite) TestCreateBusiness_TypeEI_SizeME() {
	// TODO Mockar teste
	s.T().Skip("Bankly está retornando erro 403 no endpoint de business. Mockar este teste.")

	businessRequest := createBusinessRequest(grok.GeneratorCNPJ(), models.BusinessTypeEI, models.BusinessSizeME)

	err := s.business.CreateBusiness(context.Background(), businessRequest)

	s.assert.NoError(err)
	s.assert.Nil(err)
}

func (s *BusinessTestSuite) TestCreateBusiness_TypeMEI_SizeMEI() {
	// TODO Mockar teste
	s.T().Skip("Bankly está retornando erro 403 no endpoint de business. Mockar este teste.")

	businessRequest := createBusinessRequest(grok.GeneratorCNPJ(), models.BusinessTypeMEI, models.BusinessSizeMEI)

	err := s.business.CreateBusiness(context.Background(), businessRequest)

	s.assert.NoError(err)
	s.assert.Nil(err)
}

func (s *BusinessTestSuite) TestCreateBusiness_TypeEIRELI_SizeEPP() {
	// TODO Mockar teste
	s.T().Skip("Bankly está retornando erro 403 no endpoint de business. Mockar este teste.")

	businessRequest := createBusinessRequest(grok.GeneratorCNPJ(), models.BusinessTypeEIRELI, models.BusinessSizeEPP)

	err := s.business.CreateBusiness(context.Background(), businessRequest)

	s.assert.NoError(err)
	s.assert.Nil(err)
}

func (s *BusinessTestSuite) TestCreateBusinessErrorInvalidTypeMEIAndSizeEPP() {
	// TODO Mockar teste
	s.T().Skip("Bankly está retornando erro 403 no endpoint de business. Mockar este teste.")

	businessRequest := createBusinessRequest(grok.GeneratorCNPJ(), models.BusinessTypeMEI, models.BusinessSizeEPP)

	err := s.business.CreateBusiness(context.Background(), businessRequest)

	s.assert.Error(err)

	banklyErr, ok := errors.ParseErr(err)

	s.assert.True(ok)
	s.assert.Equal(errors.ErrInvalidBusinessSize, banklyErr.GrokError)
}

func (s *BusinessTestSuite) TestCreateBusinessErrorInvalidTypeMEIAndSizeME() {
	// TODO Mockar teste
	s.T().Skip("Bankly está retornando erro 403 no endpoint de business. Mockar este teste.")

	businessRequest := createBusinessRequest(grok.GeneratorCNPJ(), models.BusinessTypeMEI, models.BusinessSizeME)

	err := s.business.CreateBusiness(context.Background(), businessRequest)

	s.assert.Error(err)

	banklyErr, ok := errors.ParseErr(err)

	s.assert.True(ok)
	s.assert.Equal(errors.ErrInvalidBusinessSize, banklyErr.GrokError)
}

func (s *BusinessTestSuite) TestUpdateBusinessName() {

	// TODO O PATCH da Business não está funcionando na Bankly. Aguardando retorno deles para ativar este teste.
	s.T().Skip("O PATCH da Business não está funcionando na Bankly. Aguardando retorno deles para ativar este teste.")

	// create business
	businessRequest := createBusinessRequest(grok.GeneratorCNPJ(), models.BusinessTypeMEI, models.BusinessSizeMEI)
	err := s.business.CreateBusiness(context.Background(), businessRequest)
	s.assert.NoError(err)
	s.assert.Nil(err)

	// update business
	newBusinessName := "NOVO EMAIL DA EMPRESA VIA UPDATE REQUEST"
	businessUpdateRequest := &models.BusinessUpdateRequest{
		BusinessName: newBusinessName,
	}
	err = s.business.UpdateBusiness(context.Background(), businessRequest.Document, *businessUpdateRequest)
	s.assert.NoError(err)
	s.assert.Nil(err)

	// find updated business
	updatedAccount, err := s.business.FindBusiness(context.Background(), businessRequest.Document)
	s.assert.NoError(err)
	s.assert.NotNil(updatedAccount)
	s.assert.Equal(businessRequest.Document, updatedAccount.Document)
	s.assert.Equal(newBusinessName, updatedAccount.BusinessName)
}

func (s *BusinessTestSuite) TestUpdateBusinessEmailAndBusinessTypeAndBusinessType() {

	// TODO O PATCH da Business não está funcionando na Bankly. Aguardando retorno deles para ativar este teste.
	s.T().Skip("O PATCH da Business não está funcionando na Bankly. Aguardando retorno deles para ativar este teste.")

	// create business
	businessRequest := createBusinessRequest(grok.GeneratorCNPJ(), models.BusinessTypeMEI, models.BusinessSizeMEI)
	err := s.business.CreateBusiness(context.Background(), businessRequest)
	s.assert.NoError(err)
	s.assert.Nil(err)

	// update business
	newEmail := "novo_email_" + businessRequest.Document + "@contbank.com"
	businessUpdateRequest := &models.BusinessUpdateRequest{
		BusinessEmail: newEmail,
		BusinessType:  models.BusinessTypeEI,
		BusinessSize:  models.BusinessSizeME,
	}
	err = s.business.UpdateBusiness(context.Background(), businessRequest.Document, *businessUpdateRequest)
	s.assert.NoError(err)
	s.assert.Nil(err)

	// find business updated
	updatedAccount, err := s.business.FindBusiness(context.Background(), "59619372000143")
	s.assert.NoError(err)
	s.assert.NotNil(updatedAccount)
	s.assert.Equal(businessRequest.Document, updatedAccount.Document)
	s.assert.Equal(models.BusinessTypeEI, updatedAccount.BusinessType)
	s.assert.Equal(models.BusinessSizeME, updatedAccount.BusinessSize)
	s.assert.Equal(newEmail, updatedAccount.BusinessEmail)

}

func (s *BusinessTestSuite) TestCreateBusinessAccount() {
	// TODO Mockar teste
	s.T().Skip("Bankly está retornando erro 403 no endpoint de business. Mockar este teste.")

	businessRequest := createBusinessRequest(grok.GeneratorCNPJ(), models.BusinessTypeEI, models.BusinessSizeME)

	err := s.business.CreateBusiness(context.Background(), businessRequest)
	s.assert.NoError(err)
	s.assert.Nil(err)

	businessAccountRequest := models.BusinessAccountRequest{
		Document:    businessRequest.Document,
		AccountType: models.PaymentAccount,
	}

	time.Sleep(time.Second)

	account, err := s.business.CreateBusinessAccount(context.Background(), businessAccountRequest)
	s.assert.NoError(err)
	s.assert.NotNil(account)
}

func (s *BusinessTestSuite) TestCreateBusinessAccountErrorDoesntHaveAnApprovedRegistration() {
	// TODO Mockar teste
	s.T().Skip("Bankly está retornando erro 403 no endpoint de business. Mockar este teste.")

	businessAccountRequest := models.BusinessAccountRequest{
		Document:    utils.GeneratorCNPJ(),
		AccountType: models.PaymentAccount,
	}
	account, err := s.business.CreateBusinessAccount(context.Background(), businessAccountRequest)

	s.assert.Error(err)

	banklyErr, ok := errors.ParseErr(err)

	s.assert.True(ok)
	s.assert.Equal(errors.ErrAccountHolderNotExists, banklyErr.GrokError)
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

	businessRequest := createBusinessRequest(grok.GeneratorCNPJ(), models.BusinessTypeMEI, models.BusinessSizeMEI)
	businessName := businessRequest.LegalRepresentative.RegisterName + " " + businessRequest.LegalRepresentative.Document
	err := s.business.CreateBusiness(context.Background(), businessRequest)

	s.assert.NoError(err)
	s.assert.Nil(err)

	time.Sleep(time.Millisecond)

	account, err := s.business.FindBusiness(context.Background(), businessRequest.Document)

	s.assert.NoError(err)
	s.assert.NotNil(account)
	s.assert.Equal(businessName, account.BusinessName)
}

func (s *BusinessTestSuite) TestBusinessName_TypeEI() {
	// TODO Mockar teste
	s.T().Skip("Bankly está retornando erro 403 no endpoint de business. Mockar este teste.")

	businessRequest := createBusinessRequest(grok.GeneratorCNPJ(), models.BusinessTypeEI, models.BusinessSizeEPP)
	err := s.business.CreateBusiness(context.Background(), businessRequest)

	s.assert.NoError(err)
	s.assert.Nil(err)

	time.Sleep(time.Millisecond)

	account, err := s.business.FindBusiness(context.Background(), businessRequest.Document)

	s.assert.NoError(err)
	s.assert.NotNil(account)
	s.assert.Equal(businessRequest.BusinessName, account.BusinessName)
}

func (s *BusinessTestSuite) TestBusinessName_TypeEIRELI() {
	// TODO Mockar teste
	s.T().Skip("Bankly está retornando erro 403 no endpoint de business. Mockar este teste.")

	businessRequest := createBusinessRequest(grok.GeneratorCNPJ(), models.BusinessTypeEIRELI, models.BusinessSizeEPP)
	err := s.business.CreateBusiness(context.Background(), businessRequest)

	s.assert.NoError(err)
	s.assert.Nil(err)

	time.Sleep(time.Millisecond)

	account, err := s.business.FindBusiness(context.Background(), businessRequest.Document)

	s.assert.NoError(err)
	s.assert.NotNil(account)
	s.assert.Equal(businessRequest.BusinessName, account.BusinessName)
}

func createBusinessRequest(document string, businessType models.BusinessType, businessSize models.BusinessSize) models.BusinessRequest {

	randomName := utils.RandStringBytes(10)
	email := "email_de_teste_" + randomName + "@contbank.com"

	return models.BusinessRequest{
		Document:            utils.OnlyDigits(document),
		BusinessName:        "Nome da Empresa " + randomName,
		TradingName:         "Nome Fantasia " + randomName,
		BusinessEmail:       email,
		BusinessType:        businessType,
		BusinessSize:        businessSize,
		BusinessAddress:     createAddress(),
		LegalRepresentative: createLegalRepresentative(),
	}
}

func createLegalRepresentative() *models.LegalRepresentative {
	randSurname := utils.RandStringBytes(10)
	return &models.LegalRepresentative{
		Document:     grok.GeneratorCPF(),
		RegisterName: "Nome do Representante Legal " + randSurname,
		Phone: &models.Phone{
			CountryCode: "55",
			Number:      utils.GeneratorCellphone(),
		},
		Address:    createAddress(),
		BirthDate:  time.Date(1990, time.June, 29, 0, 0, 0, 0, time.UTC),
		MotherName: "Nome da Mãe do Representante Legal " + randSurname,
		Email:      "email_legal_representative_" + randSurname + "@contbank.com",
	}
}

func createAddress() *models.Address {
	return &models.Address{
		ZipCode:        "05410900",
		City:           "São Paulo",
		AddressLine:    "Rua Qualquer XYZ",
		BuildingNumber: "2000",
		Neighborhood:   "Pinheiros",
		State:          "SP",
		Country:        "BR",
	}
}
